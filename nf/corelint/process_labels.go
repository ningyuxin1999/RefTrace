package corelint

import (
	"fmt"
	"reft-go/nf"
	"reft-go/nf/directives"
	"unicode"
)

var correctProcessLabels = []string{
	"process_single",
	"process_low",
	"process_medium",
	"process_high",
	"process_long",
	"process_high_memory",
}

func getLabels(process nf.Process) []*directives.LabelDirective {
	var labels []*directives.LabelDirective
	for _, directive := range process.Directives {
		if label, ok := directive.(*directives.LabelDirective); ok {
			labels = append(labels, label)
		}
	}
	return labels
}

func ruleAlphanumerics(module *nf.Module) LintResults {
	results := LintResults{
		ModulePath: module.Path,
		Errors:     []ModuleError{},
		Warnings:   []ModuleWarning{},
	}

	checkFn := func(label string) string {
		for _, char := range label {
			if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' {
				return fmt.Sprintf("process label '%s' contains non-alphanumeric characters (only letters, numbers and underscores recommended)", label)
			}
		}
		return ""
	}

	for _, process := range module.Processes {
		labels := getLabels(process)
		for _, label := range labels {
			if msg := checkFn(label.Label); msg != "" {
				results.Warnings = append(results.Warnings, ModuleWarning{
					Warning: msg,
					Line:    label.Line(),
				})
			}
		}
	}

	return results
}
