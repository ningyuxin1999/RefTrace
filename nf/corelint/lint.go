package corelint

import (
	"fmt"
	"net/url"
	"reft-go/nf"
	"reft-go/nf/directives"
	"strings"
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
func NFCoreLint(directory string) LintResults {
	results := LintResults{
		Errors:   make([]ModuleError, 0),
		Warnings: make([]ModuleWarning, 0),
	}

	modules, err := nf.ProcessDirectory(directory)
	if err != nil {
		results.Errors = append(results.Errors, ModuleError{
			ModulePath: directory,
			Error:      err,
		})
		return results
	}

	for _, module := range modules {
		for _, rule := range moduleRules {
			result := rule(module)
			if result.Error != nil {
				results.Errors = append(results.Errors, *result.Error)
			}
			if result.Warning != nil {
				results.Warnings = append(results.Warnings, *result.Warning)
			}
		}
	}
	return results
}

// Boilerplate

type ModuleError struct {
	ModulePath string
	Error      error
	Line       int
}

type ModuleWarning struct {
	ModulePath string
	Warning    string
	Line       int
}

type LintResults struct {
	Errors   []ModuleError
	Warnings []ModuleWarning
}

type LintResult struct {
	Error   *ModuleError
	Warning *ModuleWarning
}

type ModuleRule func(*nf.Module) LintResult

var moduleRules []ModuleRule

func init() {
	moduleRules = []ModuleRule{ruleContainerWithSpace, ruleMultipleContainers}
}

// Rules

func ruleContainerWithSpace(module *nf.Module) LintResult {
	for _, process := range module.Processes {
		for _, directive := range process.Directives {
			if container, ok := directive.(*directives.Container); ok {
				names := container.GetNames()
				for _, name := range names {
					if strings.Contains(name, " ") {
						return LintResult{
							Error: &ModuleError{
								ModulePath: module.Path,
								Error:      fmt.Errorf("container name '%s' contains spaces, which is not allowed", container.SimpleName),
								Line:       container.Line(),
							},
						}
					}
				}
			}
		}
	}
	return LintResult{}
}

func ruleMultipleContainers(module *nf.Module) LintResult {
	for _, process := range module.Processes {
		for _, directive := range process.Directives {
			if container, ok := directive.(*directives.Container); ok {
				names := container.GetNames()
				for _, name := range names {
					if strings.Contains(name, "biocontainers/") && (strings.Contains(name, "https://containers") || strings.Contains(name, "https://depot")) {
						return LintResult{
							Warning: &ModuleWarning{
								ModulePath: module.Path,
								Warning:    "Docker and Singularity containers specified on the same line",
								Line:       container.Line(),
							},
						}
					}
				}
			}
		}
	}
	return LintResult{}
}

func dockerOrSingularity(containerName string) string {
	// Check for Singularity container URLs
	if strings.HasPrefix(containerName, "https://") || strings.HasPrefix(containerName, "https://depot") {
		// Try parsing as URL to validate
		_, err := url.Parse(containerName)
		if err == nil {
			return "singularity"
		}
		return ""
	}

	// Check for Docker container format (org/image:tag)
	if strings.Count(containerName, "/") >= 1 &&
		strings.Count(containerName, ":") == 1 &&
		strings.Count(containerName, " ") == 0 &&
		!strings.Contains(containerName, "https://") {
		return "docker"
	}

	return ""
}
