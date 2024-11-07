package corelint

import (
	"reft-go/nf"
)

/*
NFCoreLint runs the nf-core linting rules on the given directory.

ProcessDirectory() parses the Nextflow DSL. The errors returned by
this function are not linting errors, but errors parsing the DSL.
For instance, a container named "ubuntu " would be rejected as a matter
of policy. But a container named "ubuntu"latest"" is always wrong at the DSL level.

Some of these linting rules could be implemented in user-space (rules.py).
But they are implemented here to avoid forcing users to define a rules.py
for common cases.
*/
func NFCoreLint(directory string) ([]LintResults, error) {
	results := []LintResults{}

	modules, err := nf.ProcessDirectory(directory)
	if err != nil {
		return nil, err
	}

	// Run all rules and merge their results
	for _, module := range modules {
		// Create a moduleResults to collect all rule results for this module
		moduleResults := LintResults{
			ModulePath: module.Path,
			Errors:     make([]ModuleError, 0),
			Warnings:   make([]ModuleWarning, 0),
		}

		// Run all rules and merge their results
		for _, rule := range moduleRules {
			ruleResult := rule(module)
			moduleResults.Errors = append(moduleResults.Errors, ruleResult.Errors...)
			moduleResults.Warnings = append(moduleResults.Warnings, ruleResult.Warnings...)
		}

		results = append(results, moduleResults)
	}

	return results, nil
}

// Boilerplate

type ModuleError struct {
	Error error
	Line  int
}

type ModuleWarning struct {
	Warning string
	Line    int
}

type LintResults struct {
	ModulePath string
	Errors     []ModuleError
	Warnings   []ModuleWarning
}

type ModuleRule func(*nf.Module) LintResults

var moduleRules []ModuleRule

func init() {
	moduleRules = []ModuleRule{
		ruleContainerWithSpace,
		ruleMultipleContainers,
		ruleMustBeTagged,
		ruleAlphanumerics,
	}
}
