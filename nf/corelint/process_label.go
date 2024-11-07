package corelint

import (
	"fmt"
	"reft-go/nf"
	"reft-go/nf/directives"
	"slices"
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

func ruleConflictingLabels(module *nf.Module) LintResults {
	results := LintResults{
		ModulePath: module.Path,
		Errors:     []ModuleError{},
		Warnings:   []ModuleWarning{},
	}

	for _, process := range module.Processes {
		var goodLabels []*directives.LabelDirective
		labels := getLabels(process)
		for _, label := range labels {
			if slices.Contains(correctProcessLabels, label.Label) {
				goodLabels = append(goodLabels, label)
			}
		}
		if len(goodLabels) > 1 {
			labelNames := make([]string, len(goodLabels))
			for i, label := range goodLabels {
				labelNames[i] = label.Label
			}
			results.Warnings = append(results.Warnings, ModuleWarning{
				Warning: fmt.Sprintf("process '%s' has conflicting labels: %v", process.Name, labelNames),
				Line:    goodLabels[len(goodLabels)-1].Line(),
			})
		}
	}

	return results
}

func ruleNoStandardLabel(module *nf.Module) LintResults {
	results := LintResults{
		ModulePath: module.Path,
		Errors:     []ModuleError{},
		Warnings:   []ModuleWarning{},
	}

	for _, process := range module.Processes {
		var goodLabels []*directives.LabelDirective
		labels := getLabels(process)
		for _, label := range labels {
			if slices.Contains(correctProcessLabels, label.Label) {
				goodLabels = append(goodLabels, label)
			}
		}
		if len(goodLabels) == 0 {
			results.Warnings = append(results.Warnings, ModuleWarning{
				Warning: fmt.Sprintf("process '%s' has no standard label", process.Name),
				Line:    process.Line(),
			})
		}
	}

	return results
}

func ruleNonStandardLabel(module *nf.Module) LintResults {
	results := LintResults{
		ModulePath: module.Path,
		Errors:     []ModuleError{},
		Warnings:   []ModuleWarning{},
	}

	for _, process := range module.Processes {
		var badLabels []*directives.LabelDirective
		labels := getLabels(process)
		for _, label := range labels {
			if !slices.Contains(correctProcessLabels, label.Label) {
				badLabels = append(badLabels, label)
			}
		}
		if len(badLabels) > 0 {
			labelNames := make([]string, len(badLabels))
			for i, label := range badLabels {
				labelNames[i] = label.Label
			}
			results.Warnings = append(results.Warnings, ModuleWarning{
				Warning: fmt.Sprintf("process '%s' has non-standard labels: %v", process.Name, labelNames),
				Line:    badLabels[0].Line(),
			})
		}
	}

	return results
}

func ruleDuplicateLabels(module *nf.Module) LintResults {
	results := LintResults{
		ModulePath: module.Path,
		Errors:     []ModuleError{},
		Warnings:   []ModuleWarning{},
	}

	for _, process := range module.Processes {
		// Create a map to track label occurrences
		labelCount := make(map[string]int)
		labels := getLabels(process)

		// Count occurrences of each label
		for _, label := range labels {
			labelCount[label.Label]++
		}

		// Check for duplicates
		for labelName, count := range labelCount {
			if count > 1 {
				results.Warnings = append(results.Warnings, ModuleWarning{
					Warning: fmt.Sprintf("process '%s' has duplicate label '%s' (%d times)",
						process.Name, labelName, count),
					Line: process.Line(),
				})
			}
		}
	}

	return results
}

func ruleNoLabels(module *nf.Module) LintResults {
	results := LintResults{
		ModulePath: module.Path,
		Errors:     []ModuleError{},
		Warnings:   []ModuleWarning{},
	}

	for _, process := range module.Processes {
		labels := getLabels(process)
		if len(labels) == 0 {
			results.Warnings = append(results.Warnings, ModuleWarning{
				Warning: fmt.Sprintf("process '%s' has no labels", process.Name),
				Line:    process.Line(),
			})
		}
	}

	return results
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
