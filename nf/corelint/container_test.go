package corelint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"reft-go/nf"
	"reft-go/nf/directives"
)

func TestRuleContainerWithSpace(t *testing.T) {
	processContent := `
process FOO {
    container 'ubuntu latest'  // Invalid container name with space
    cpus 4
    memory '8 GB'

    input:
    path x

    output:
    path 'y'

    script:
    """
    echo "test"
    """
}
`
	// Create temporary file
	tmpDir := t.TempDir()
	processFile := filepath.Join(tmpDir, "process.nf")
	if err := os.WriteFile(processFile, []byte(processContent), 0644); err != nil {
		t.Fatal("Failed to write process file:", err)
	}

	// Parse the file
	module, err := nf.BuildModule(processFile)
	if err != nil {
		t.Fatal("Failed to parse process file:", err)
	}

	// Run the container space check
	results := ruleContainerWithSpace(module)
	if len(results.Errors) == 0 {
		t.Fatal("Expected error due to space in container name, but got none")
	}

	expectedError := "container name 'ubuntu latest' contains spaces, which is not allowed"
	if results.Errors[0].Error.Error() != expectedError {
		t.Errorf("Expected error message %q but got %q", expectedError, results.Errors[0].Error.Error())
	}
	if results.ModulePath != processFile {
		t.Errorf("Expected module path %q but got %q", processFile, results.ModulePath)
	}
}

func TestContainerWithoutSpace(t *testing.T) {
	processContent := `
process FOO {
    container 'ubuntu-latest'  // Valid container name without space
    cpus 4
    memory '8 GB'

    input:
    path x

    output:
    path 'y'

    script:
    """
    echo "test"
    """
}
`
	// Create temporary file
	tmpDir := t.TempDir()
	processFile := filepath.Join(tmpDir, "process.nf")
	if err := os.WriteFile(processFile, []byte(processContent), 0644); err != nil {
		t.Fatal("Failed to write process file:", err)
	}

	// Parse the file
	module, err := nf.BuildModule(processFile)
	if err != nil {
		t.Fatal("Failed to parse process file:", err)
	}

	// Run the container space check
	results := ruleContainerWithSpace(module)
	if len(results.Errors) > 0 {
		t.Errorf("Expected no errors for valid container name, but got: %v", results.Errors[0].Error)
	}
	if len(results.Warnings) > 0 {
		t.Errorf("Expected no warnings for valid container name, but got: %v", results.Warnings[0].Warning)
	}
}

func TestContainerWithTooManyQuotes(t *testing.T) {
	processContent := `
process FOO {
    container "ubuntu"latest""  // Invalid: too many double quotes

    input:
    path x

    output:
    path 'y'

    script:
    """
    echo "test"
    """
}
`
	// Create temporary file
	tmpDir := t.TempDir()
	processFile := filepath.Join(tmpDir, "process.nf")
	if err := os.WriteFile(processFile, []byte(processContent), 0644); err != nil {
		t.Fatal("Failed to write process file:", err)
	}

	// Parse the file - this should fail
	_, err := nf.BuildModule(processFile)
	if err == nil {
		t.Fatal("Expected parsing to fail due to invalid container specification, but got no error")
	}

	expectedError := "too many quotes found when specifying container"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error message to contain %q but got %q", expectedError, err.Error())
	}
}

func TestContainerWithTernaryOperator(t *testing.T) {
	processContent := `
process FOO {
    container "${ workflow.containerEngine == 'singularity' && !task.ext.singularity_pull_docker_container ?
            'https://depot.galaxyproject.org/singularity/gatk4:4.4.0.0--foo':
            'biocontainers/gatk4:4.4.0.0--foo' }"
    cpus 4
    memory '8 GB'

    input:
    path x

    output:
    path 'y'

    script:
    """
    echo "test"
    """
}
`
	// Create temporary file
	tmpDir := t.TempDir()
	processFile := filepath.Join(tmpDir, "process.nf")
	if err := os.WriteFile(processFile, []byte(processContent), 0644); err != nil {
		t.Fatal("Failed to write process file:", err)
	}

	// Parse the file
	module, err := nf.BuildModule(processFile)
	if err != nil {
		t.Fatal("Failed to parse process file:", err)
	}

	// Run the container space check
	results := ruleContainerWithSpace(module)
	if len(results.Errors) > 0 {
		t.Errorf("Expected no error for valid ternary container names, but got: %v", results.Errors[0].Error)
	}
	if len(results.Warnings) > 0 {
		t.Errorf("Expected no warning for valid ternary container names, but got: %v", results.Warnings[0].Warning)
	}
}

func TestContainerWithMixedSyntax(t *testing.T) {
	processContent := `
process FOO {
    container 'https://containers/biocontainers/ubuntu'  // Invalid: mixing Docker and Singularity
    cpus 4
    memory '8 GB'

    input:
    path x

    output:
    path 'y'

    script:
    """
    echo "test"
    """
}
`
	// Create temporary file
	tmpDir := t.TempDir()
	processFile := filepath.Join(tmpDir, "process.nf")
	if err := os.WriteFile(processFile, []byte(processContent), 0644); err != nil {
		t.Fatal("Failed to write process file:", err)
	}

	// Parse the file
	module, err := nf.BuildModule(processFile)
	if err != nil {
		t.Fatal("Failed to parse process file:", err)
	}

	// Run the container space check
	results := ruleContainerWithSpace(module)
	if len(results.Errors) > 0 {
		t.Errorf("Expected no error for container without spaces, but got: %v", results.Errors[0].Error)
	}
	if len(results.Warnings) > 0 {
		t.Errorf("Expected no warning for container without spaces, but got: %v", results.Warnings[0].Warning)
	}
}

func TestContainerWithMixedSyntaxWarning(t *testing.T) {
	processContent := `
process FOO {
    container "biocontainers/ubuntu https://containers/something"
    script:
    """
    echo "test"
    """
}
`
	// Create temporary file
	tmpDir := t.TempDir()
	processFile := filepath.Join(tmpDir, "process.nf")
	if err := os.WriteFile(processFile, []byte(processContent), 0644); err != nil {
		t.Fatal("Failed to write process file:", err)
	}

	// Parse the file
	module, err := nf.BuildModule(processFile)
	if err != nil {
		t.Fatal("Failed to parse process file:", err)
	}

	// Run the multiple containers check
	results := ruleMultipleContainers(module)
	if len(results.Warnings) == 0 {
		t.Fatal("Expected warning for mixed container syntax, but got none")
	}

	expectedWarning := "Docker and Singularity containers specified on the same line"
	if results.Warnings[0].Warning != expectedWarning {
		t.Errorf("Expected warning message %q but got %q", expectedWarning, results.Warnings[0].Warning)
	}
	if results.ModulePath != processFile {
		t.Errorf("Expected module path %q but got %q", processFile, results.ModulePath)
	}
}

func TestSingularityContainerTags(t *testing.T) {
	testCases := []struct {
		name          string
		containerURL  string
		expectedTag   string
		shouldSucceed bool
	}{
		{
			name:          "biocontainers URL with version",
			containerURL:  "https://containers.biocontainers.pro/s3/SingImgsRepo/biocontainers/v1.2.0_cv1/biocontainers_v1.2.0_cv1.img",
			expectedTag:   "v1.2.0_cv1",
			shouldSucceed: true,
		},
		{
			name:          "galaxy depot URL with version",
			containerURL:  "https://depot.galaxyproject.org/singularity/fastqc:0.11.9--0",
			expectedTag:   "0.11.9--0",
			shouldSucceed: true,
		},
		{
			name:          "invalid URL without tag",
			containerURL:  "https://depot.galaxyproject.org/singularity/fastqc",
			expectedTag:   "",
			shouldSucceed: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tag, err := getSingularityTag(tc.containerURL)

			if tc.shouldSucceed {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
				if tag != tc.expectedTag {
					t.Errorf("Expected tag %q but got %q", tc.expectedTag, tag)
				}
			} else {
				if err == nil {
					t.Error("Expected error but got success")
				}
			}
		})
	}
}

func TestValidDockerContainerWithTag(t *testing.T) {
	processContent := `
process FOO {
    container "biocontainers/fastqc:0.11.9"
}`
	runRuleMustBeTaggedTest(t,
		"valid docker container with tag",
		processContent,
		0,
		nil,
	)
}

func TestValidSingularityContainerWithTag(t *testing.T) {
	processContent := `
process FOO {
    container "https://depot.galaxyproject.org/singularity/fastqc:0.11.9--0"
}`
	runRuleMustBeTaggedTest(t,
		"valid singularity container with tag",
		processContent,
		0,
		nil,
	)
}

func TestDockerContainerMissingTag(t *testing.T) {
	processContent := `
process FOO {
    container "biocontainers/fastqc"
}`
	runRuleMustBeTaggedTest(t,
		"docker container missing tag",
		processContent,
		1,
		[]string{"docker container must specify a tag"},
	)
}

func TestSingularityContainerMissingTag(t *testing.T) {
	processContent := `
process FOO {
    container "https://depot.galaxyproject.org/singularity/fastqc"
}`
	runRuleMustBeTaggedTest(t,
		"singularity container missing tag",
		processContent,
		1,
		[]string{"singularity container must specify a tag"},
	)
}

func TestQuayIoContainerWithoutTag(t *testing.T) {
	processContent := `
process FOO {
    container "quay.io/biocontainers/fastqc"
}`
	runRuleMustBeTaggedTest(t,
		"quay.io container without tag",
		processContent,
		2,
		[]string{
			"docker container must specify a tag",
			"please use 'organisation/container:tag' format instead of full registry URL",
		},
	)
}

func TestMultipleContainersWithTernary(t *testing.T) {
	processContent := `
process FOO {
    container "${workflow.containerEngine == 'singularity' ? 
        'https://depot.galaxyproject.org/singularity/fastqc' : 
        'biocontainers/fastqc'}"
}`
	runRuleMustBeTaggedTest(t,
		"multiple containers with ternary",
		processContent,
		2,
		[]string{"container must specify a tag"},
	)
}

func runRuleMustBeTaggedTest(t *testing.T, testName, processContent string, wantErrors int, errorContains []string) {
	t.Run(testName, func(t *testing.T) {
		// Create temporary file
		tmpDir := t.TempDir()
		processFile := filepath.Join(tmpDir, "process.nf")
		if err := os.WriteFile(processFile, []byte(processContent), 0644); err != nil {
			t.Fatal("Failed to write process file:", err)
		}

		// Parse the file
		module, err := nf.BuildModule(processFile)
		if err != nil {
			t.Fatal("Failed to parse process file:", err)
		}

		// Debug: Print process and container information
		t.Logf("Test case: %s", testName)
		for _, process := range module.Processes {
			t.Logf("  Process name: %s", process.Name)
			for _, directive := range process.Directives {
				if container, ok := directive.(*directives.Container); ok {
					names := container.GetNames()
					t.Logf("    Container names: %v", names)
					for _, name := range names {
						containerType, _ := dockerOrSingularity(name)
						t.Logf("    Container %q detected as %s", name, containerType)
					}
				}
			}
		}

		// Run the tag check
		results := ruleMustBeTagged(module)

		// Debug: Print results
		t.Logf("  Got %d errors (expected %d):", len(results.Errors), wantErrors)
		for i, err := range results.Errors {
			t.Logf("    Error %d: %v (line %d)", i+1, err.Error, err.Line)
		}

		// Check number of errors
		if got := len(results.Errors); got != wantErrors {
			t.Errorf("ruleMustBeTagged() got %v errors, want %v", got, wantErrors)
		}

		if errorContains != nil {
			for _, expectedErr := range errorContains {
				found := false
				for _, gotErr := range results.Errors {
					if strings.Contains(gotErr.Error.Error(), expectedErr) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Missing expected error %q\nGot errors:", expectedErr)
					for _, err := range results.Errors {
						t.Errorf("  - %v", err.Error)
					}
				}
			}
		}
	})
}

func TestUnknownContainerType(t *testing.T) {
	processContent := `
process FOO {
    container "just-a-name"  // Invalid: no URL prefix or org/image format
    script:
    """
    echo "test"
    """
}
`
	// Create temporary file
	tmpDir := t.TempDir()
	processFile := filepath.Join(tmpDir, "process.nf")
	if err := os.WriteFile(processFile, []byte(processContent), 0644); err != nil {
		t.Fatal("Failed to write process file:", err)
	}

	// Parse the file
	module, err := nf.BuildModule(processFile)
	if err != nil {
		t.Fatal("Failed to parse process file:", err)
	}

	// Get the container directive from the process
	process := module.Processes[0]
	container := process.Directives[0].(*directives.Container)
	containerName := container.GetNames()[0]

	// Test the dockerOrSingularity function directly
	_, err = dockerOrSingularity(containerName)
	if err == nil {
		t.Error("Expected error for unknown container type, but got none")
	}
	if err.Error() != "unknown container type" {
		t.Errorf("Expected error message 'unknown container type' but got %q", err.Error())
	}
}

func TestInvalidSingularityContainerURL(t *testing.T) {
	processContent := `
process FOO {
    container "https://[invalid-url"  // Invalid: malformed URL
    script:
    """
    echo "test"
    """
}
`
	// Create temporary file
	tmpDir := t.TempDir()
	processFile := filepath.Join(tmpDir, "process.nf")
	if err := os.WriteFile(processFile, []byte(processContent), 0644); err != nil {
		t.Fatal("Failed to write process file:", err)
	}

	// Parse the file
	module, err := nf.BuildModule(processFile)
	if err != nil {
		t.Fatal("Failed to parse process file:", err)
	}

	// Get the container directive from the process
	process := module.Processes[0]
	container := process.Directives[0].(*directives.Container)
	containerName := container.GetNames()[0]

	// Test the dockerOrSingularity function directly
	_, err = dockerOrSingularity(containerName)
	if err == nil {
		t.Error("Expected error for invalid Singularity URL, but got none")
	}
	if err.Error() != "invalid singularity container URL" {
		t.Errorf("Expected error message 'invalid singularity container URL' but got %q", err.Error())
	}
}

func TestInvalidDockerTagFormat(t *testing.T) {
	testCases := []struct {
		name          string
		containerName string
		wantErr       string
	}{
		{
			name:          "invalid characters",
			containerName: "biocontainers/fastqc:1.0$",
			wantErr:       "invalid docker tag format",
		},
		{
			name:          "empty tag",
			containerName: "biocontainers/fastqc:",
			wantErr:       "invalid docker tag format",
		},
		{
			name:          "space in tag",
			containerName: "biocontainers/fastqc:1.0 2",
			wantErr:       "invalid docker tag format",
		},
		{
			name:          "valid tag",
			containerName: "biocontainers/fastqc:v1.0.2-alpha_beta.3",
			wantErr:       "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tag, err := getDockerTag(tc.containerName)
			if tc.wantErr == "" {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
				if !isValidTag(tag) {
					t.Errorf("Got invalid tag: %s", tag)
				}
			} else {
				if err == nil {
					t.Error("Expected error but got success")
				}
				if err.Error() != tc.wantErr {
					t.Errorf("Expected error %q but got %q", tc.wantErr, err.Error())
				}
			}
		})
	}
}

func TestGetSingularityTag(t *testing.T) {
	testCases := []struct {
		name          string
		containerURL  string
		expectedTag   string
		expectedError string
	}{
		{
			name:          "invalid URL format",
			containerURL:  "https://example.com/%ZZ",
			expectedError: "invalid container URL: parse",
		},
		{
			name:          "URL with no path",
			containerURL:  "https://example.com",
			expectedError: "invalid container URL: no path segments",
		},
		{
			name:          "URL ending in slash",
			containerURL:  "https://example.com/",
			expectedError: "invalid container URL: no path segments",
		},
		{
			name:          "URL with dot",
			containerURL:  "https://example.com/.",
			expectedError: "invalid container URL: no path segments",
		},
		{
			name:          "valid URL with colon tag",
			containerURL:  "https://depot.galaxyproject.org/singularity/fastqc:0.11.9--0",
			expectedTag:   "0.11.9--0",
			expectedError: "",
		},
		{
			name:          "valid URL with _v tag",
			containerURL:  "https://containers.biocontainers.pro/s3/SingImgsRepo/biocontainers_v1.2.0_cv1.img",
			expectedTag:   "v1.2.0_cv1",
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tag, err := getSingularityTag(tc.containerURL)

			if tc.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error containing %q but got nil", tc.expectedError)
					return
				}
				if !strings.Contains(err.Error(), tc.expectedError) {
					t.Errorf("Expected error containing %q but got %q", tc.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
					return
				}
				if tag != tc.expectedTag {
					t.Errorf("Expected tag %q but got %q", tc.expectedTag, tag)
				}
			}
		})
	}
}

func TestRuleMustBeTaggedUnknownContainerType(t *testing.T) {
	processContent := `
process FOO {
    container "just-a-name"  // Invalid: no URL prefix or org/image format
    script:
    """
    echo "test"
    """
}
`
	// Create temporary file
	tmpDir := t.TempDir()
	processFile := filepath.Join(tmpDir, "process.nf")
	if err := os.WriteFile(processFile, []byte(processContent), 0644); err != nil {
		t.Fatal("Failed to write process file:", err)
	}

	// Parse the file
	module, err := nf.BuildModule(processFile)
	if err != nil {
		t.Fatal("Failed to parse process file:", err)
	}

	// Run the tag check
	results := ruleMustBeTagged(module)

	// Should have exactly one error
	if len(results.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(results.Errors))
	}

	// Verify error message
	expectedError := "unknown container type"
	if !strings.Contains(results.Errors[0].Error.Error(), expectedError) {
		t.Errorf("Expected error containing %q, got %q", expectedError, results.Errors[0].Error.Error())
	}

	// Verify module path
	if results.ModulePath != processFile {
		t.Errorf("Expected module path %q but got %q", processFile, results.ModulePath)
	}
}
