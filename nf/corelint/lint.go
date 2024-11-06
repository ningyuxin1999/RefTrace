package corelint

import (
	"fmt"
	"reft-go/nf"
	"reft-go/nf/directives"
	"strings"
)

/*
NFCoreLint runs the nf-core linting rules on the given directory.

ProcessDirectory() parses the Nextflow DSL. The errors returned by
this function are not linting errors, but errors parsing the DSL.
For instance, a container named "ubuntu " would be rejected as a matter
of taste. But a container named "ubuntu"latest"" is always wrong at the DSL level.

Some of these linting rules could be implemented in user-space (rules.py).
But they are implemented here to avoid forcing users to define a rules.py
for common cases.
*/
func NFCoreLint(directory string) LintResults {
	results := LintResults{
		Errors:   make([]error, 0),
		Warnings: make([]string, 0),
	}

	modules, err := nf.ProcessDirectory(directory)
	if err != nil {
		results.Errors = append(results.Errors, err)
		return results
	}

	for _, module := range modules {
		for _, rule := range moduleRules {
			result := rule(module)
			if result.Error != nil {
				results.Errors = append(results.Errors, result.Error)
			}
			if result.Warning != "" {
				results.Warnings = append(results.Warnings, result.Warning)
			}
		}
	}
	return results
}

type LintResults struct {
	Errors   []error
	Warnings []string
}

type LintResult struct {
	Error   error
	Warning string
}

type ModuleRule func(*nf.Module) LintResult

var moduleRules []ModuleRule

func init() {
	moduleRules = []ModuleRule{containerWithSpace}
}

func containerWithSpace(module *nf.Module) LintResult {
	for _, process := range module.Processes {
		for _, directive := range process.Directives {
			if container, ok := directive.(*directives.Container); ok {
				if container.Format == directives.Simple {
					if strings.Contains(container.SimpleName, " ") {
						return LintResult{
							Error: fmt.Errorf("container name '%s' contains spaces, which is not allowed", container.SimpleName),
						}
					}
				}
				if container.Format == directives.Ternary {
					if strings.Contains(container.TrueName, " ") {
						return LintResult{
							Error: fmt.Errorf("container true_name '%s' contains spaces, which is not allowed", container.TrueName),
						}
					}
					if strings.Contains(container.FalseName, " ") {
						return LintResult{
							Error: fmt.Errorf("container false_name '%s' contains spaces, which is not allowed", container.FalseName),
						}
					}
				}
			}
		}
	}
	return LintResult{}
}
