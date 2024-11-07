package corelint

import (
	"fmt"
	"net/url"
	"path"
	"reft-go/nf"
	"reft-go/nf/directives"
	"strings"
	"unicode"
)

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

func ruleMustBeTagged(module *nf.Module) LintResult {
	for _, process := range module.Processes {
		for _, directive := range process.Directives {
			if container, ok := directive.(*directives.Container); ok {
				names := container.GetNames()
				for _, name := range names {
					containerType := dockerOrSingularity(name)
					if containerType == "singularity" {
						_, err := getSingularityTag(name)
						if err != nil {
							return LintResult{
								Error: &ModuleError{
									ModulePath: module.Path,
									Error:      err,
									Line:       container.Line(),
								},
							}
						}
					}
					if containerType == "docker" {
						_, err := getDockerTag(name)
						if err != nil {
							return LintResult{
								Error: &ModuleError{
									ModulePath: module.Path,
									Error:      err,
									Line:       container.Line(),
								},
							}
						}
						err = quayPrefix(name)
						if err != nil {
							return LintResult{
								Error: &ModuleError{
									ModulePath: module.Path,
									Error:      err,
									Line:       container.Line(),
								},
							}
						}
					}
				}
			}
		}
	}
	return LintResult{}
}

func getSingularityTag(containerName string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(containerName)
	if err != nil {
		return "", fmt.Errorf("invalid container URL: %v", err)
	}

	// Get the last segment of the path
	lastSegment := path.Base(parsedURL.Path)
	if lastSegment == "." || lastSegment == "/" {
		return "", fmt.Errorf("invalid container URL: no path segments")
	}

	lastSegment = strings.TrimSuffix(lastSegment, ".img")
	lastSegment = strings.TrimSuffix(lastSegment, ".sif")

	// Check for colon-separated tag
	if idx := strings.LastIndex(lastSegment, ":"); idx != -1 {
		tag := lastSegment[idx+1:]
		if isValidTag(tag) {
			return tag, nil
		}
	}

	// Check for _v<digit> pattern
	if idx := strings.LastIndex(lastSegment, "_v"); idx != -1 && len(lastSegment) > idx+2 {
		if unicode.IsDigit(rune(lastSegment[idx+2])) {
			tag := lastSegment[idx+1:]
			if isValidTag(tag) {
				return tag, nil
			}
		}
	}

	return "", fmt.Errorf("unsupported singularity container URL format")
}

// Helper function to validate the tag against allowed characters
func isValidTag(tag string) bool {
	if tag == "" {
		return false
	}
	for _, c := range tag {
		if !((c >= 'A' && c <= 'Z') ||
			(c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') ||
			c == '-' || c == '_' || c == '.') {
			return false
		}
	}
	return true
}

func getDockerTag(containerName string) (string, error) {
	// Look for the tag after the last colon
	if idx := strings.LastIndex(containerName, ":"); idx != -1 {
		tag := containerName[idx+1:]
		if !isValidTag(tag) {
			return "", fmt.Errorf("invalid docker tag format")
		}
		return tag, nil
	}
	return "", fmt.Errorf("no docker tag found")
}

func quayPrefix(containerName string) error {
	if strings.HasPrefix(containerName, "quay.io") {
		return fmt.Errorf("please use 'organisation/container:tag' format instead of full registry URL")
	}
	return nil
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
