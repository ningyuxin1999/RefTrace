package nf

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/fatih/color"
	re "github.com/magnetde/starlark-re"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

type LintConfig struct {
	RulesFile string
	Directory string
	RuleToRun string
}

type RuleModuleOutput struct {
	Errors  []string
	Outputs []string
}

// rule -> module -> output
type GroupedOutput map[string]map[string]RuleModuleOutput

func RunLintWithConfig(config LintConfig, output io.Writer) error {
	if output == nil {
		output = os.Stdout
	}

	// Check if rules file exists
	if _, err := os.Stat(config.RulesFile); os.IsNotExist(err) {
		return fmt.Errorf("rules file not found: %s", config.RulesFile)
	}

	// Convert relative paths to absolute paths
	rulesFile, _ := filepath.Abs(config.RulesFile)
	dir, _ := filepath.Abs(config.Directory)

	// Read the rules.py file
	rulesContent, err := os.ReadFile(rulesFile)
	if err != nil {
		return fmt.Errorf("error reading rules.py file: %v", err)
	}

	// Remove the 'fail' function from the Universe
	delete(starlark.Universe, "fail")

	var outputMutex sync.Mutex
	groupedOutput := make(GroupedOutput)

	// Create a new Starlark thread with a custom print function
	thread := &starlark.Thread{
		Name: "lint_thread",
		Print: func(thread *starlark.Thread, msg string) {
			outputMutex.Lock()
			defer outputMutex.Unlock()

			ruleName := thread.Local("current_rule").(string)
			moduleName := thread.Local("current_module").(string)

			entry := groupedOutput[ruleName][moduleName]
			entry.Outputs = append(entry.Outputs, msg)
			groupedOutput[ruleName][moduleName] = entry
		},
	}

	fo := &syntax.FileOptions{}

	// Parse the Starlark code without executing it
	f, err := fo.Parse(rulesFile, rulesContent, 0)
	if err != nil {
		log.Fatalf("Error parsing rules program: %v", err)
	}

	// Compile the parsed code
	prog, err := starlark.FileProgram(f, func(name string) bool {
		if name == "fatal" || name == "error" || name == "re" {
			return true
		}
		return false
	})
	if err != nil {
		log.Fatalf("Error compiling rules program: %v", err)
	}

	// Create predefined variables for the Starlark environment
	predefined := starlark.StringDict{
		"fatal": starlark.NewBuiltin("fatal", fatalFunc),
		"error": starlark.NewBuiltin("error", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			sep := " "
			if err := starlark.UnpackArgs("error", nil, kwargs, "sep?", &sep); err != nil {
				return nil, err
			}
			buf := new(strings.Builder)
			for i, v := range args {
				if i > 0 {
					buf.WriteString(sep)
				}
				if s, ok := starlark.AsString(v); ok {
					buf.WriteString(s)
				} else {
					buf.WriteString(v.String())
				}
			}

			outputMutex.Lock()
			defer outputMutex.Unlock()

			ruleName := thread.Local("current_rule").(string)
			moduleName := thread.Local("current_module").(string)

			entry := groupedOutput[ruleName][moduleName]
			entry.Errors = append(entry.Errors, buf.String())
			groupedOutput[ruleName][moduleName] = entry

			return starlark.None, nil
		}),
		"re": re.NewModule(), // Add the regex module
	}

	// Execute the compiled program
	globals, err := prog.Init(thread, predefined)
	if err != nil {
		log.Fatalf("Error initializing rules program: %v", err)
	}

	// Collect rules (functions starting with "rule_")
	rules := make(map[string]starlark.Callable)
	for name, value := range globals {
		if strings.HasPrefix(name, "rule_") {
			if callable, ok := value.(starlark.Callable); ok {
				strippedRuleName := strings.TrimPrefix(name, "rule_")
				rules[strippedRuleName] = callable
			}
		}
	}

	// Create the params list
	/*
		params, includes := parse(rulesFile)
		starlarkParams := make([]starlark.Value, len(params))
		for i, param := range params {
			starlarkParams[i] = NewStarlarkParamInfo(param)
		}
		_ = starlark.NewList(starlarkParams)
	*/

	/*
		starlarkIncludes := make([]starlark.Value, len(includes))
		for i, include := range includes {
			starlarkIncludes[i] = NewStarlarkIncludeInfo(include)
		}
		_ = starlark.NewList(starlarkIncludes)
	*/

	// Parse the directory and get the modules
	modules, err := ProcessDirectory(dir)
	if err != nil {
		return fmt.Errorf("error processing directory: %v", err)
	}

	// Execute each rule
	for ruleName, ruleFunc := range rules {
		if config.RuleToRun != "" && ruleName != config.RuleToRun {
			continue
		}
		groupedOutput[ruleName] = make(map[string]RuleModuleOutput)
		for _, module := range modules {
			groupedOutput[ruleName][module.Path] = RuleModuleOutput{}
			starlarkModule := ConvertToStarlarkModule(module)

			// Set the current rule and module context
			thread.SetLocal("current_rule", ruleName)
			thread.SetLocal("current_module", module.Path)

			_, err := starlark.Call(thread, ruleFunc, starlark.Tuple{starlarkModule}, nil)
			if err != nil {
				if evalErr, ok := err.(*starlark.EvalError); ok {
					entry := groupedOutput[ruleName][module.Path]
					entry.Errors = append(entry.Errors, evalErr.Msg)
					groupedOutput[ruleName][module.Path] = entry
					//fmt.Printf("Rule %s execution failed: %s\n", ruleName, evalErr.Msg)
				} else {
					log.Fatalf("Error calling rule %s: %v\n", ruleName, err)
				}
			}
		}
	}

	hasErrors := printGroupedOutput(groupedOutput, output)
	if hasErrors {
		return fmt.Errorf("Linting failed")
	}

	return nil
}

func printGroupedOutput(groupedOutput GroupedOutput, output io.Writer) bool {
	hasErrors := false
	ruleNames := make([]string, 0, len(groupedOutput))
	for ruleName := range groupedOutput {
		ruleNames = append(ruleNames, ruleName)
	}
	sort.Strings(ruleNames)

	rulePrinter := color.New(color.FgCyan, color.Bold)
	modulePrinter := color.New(color.FgYellow)
	errorPrinter := color.New(color.FgRed)
	outputPrinter := color.New(color.FgGreen)

	fmt.Fprintln(output) // Start with a blank line

	for _, ruleName := range ruleNames {
		rulePrinter.Fprintf(output, "Rule: %s\n", ruleName)

		moduleNames := make([]string, 0, len(groupedOutput[ruleName]))
		for moduleName := range groupedOutput[ruleName] {
			moduleNames = append(moduleNames, moduleName)
		}
		sort.Strings(moduleNames)

		for _, moduleName := range moduleNames {
			entry := groupedOutput[ruleName][moduleName]
			if len(entry.Errors) > 0 || len(entry.Outputs) > 0 {
				modulePrinter.Fprintf(output, "  Module: %s\n", moduleName)
				for _, e := range entry.Errors {
					hasErrors = true
					errorPrinter.Fprintf(output, "    Error: %s\n", e)
				}
				for _, o := range entry.Outputs {
					outputPrinter.Fprintf(output, "    Output: %s\n", o)
				}
			}
		}
		fmt.Fprintln(output) // Add a blank line between rules
	}
	return hasErrors
}

// failnowFunc is the implementation of the failnow function for Starlark
func fatalFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	sep := " "
	if err := starlark.UnpackArgs("failnow", nil, kwargs, "sep?", &sep); err != nil {
		return nil, err
	}
	buf := new(strings.Builder)
	for i, v := range args {
		if i > 0 {
			buf.WriteString(sep)
		}
		if s, ok := starlark.AsString(v); ok {
			buf.WriteString(s)
		} else {
			buf.WriteString(v.String())
		}
	}

	return nil, errors.New(buf.String())
}
