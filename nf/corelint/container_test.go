package corelint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"reft-go/nf"
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
