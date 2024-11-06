package corelint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"reft-go/nf"
)

func TestruleContainerWithSpace(t *testing.T) {
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
	result := ruleContainerWithSpace(module)
	if result.Error == nil {
		t.Fatal("Expected error due to space in container name, but got none")
	}

	expectedError := "container name 'ubuntu latest' contains spaces, which is not allowed"
	if result.Error.Error.Error() != expectedError {
		t.Errorf("Expected error message %q but got %q", expectedError, result.Error.Error.Error())
	}
	if result.Error.ModulePath != processFile {
		t.Errorf("Expected module path %q but got %q", processFile, result.Error.ModulePath)
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
	result := ruleContainerWithSpace(module)
	if result.Error != nil {
		t.Errorf("Expected no error for valid container name, but got: %v", result.Error.Error)
	}
	if result.Warning != nil {
		t.Errorf("Expected no warning for valid container name, but got: %v", result.Warning.Warning)
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
	result := ruleContainerWithSpace(module)
	if result.Error != nil {
		t.Errorf("Expected no error for valid ternary container names, but got: %v", result.Error.Error)
	}
	if result.Warning != nil {
		t.Errorf("Expected no warning for valid ternary container names, but got: %v", result.Warning.Warning)
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
	result := ruleContainerWithSpace(module)
	if result.Error != nil {
		t.Errorf("Expected no error for container without spaces, but got: %v", result.Error.Error)
	}
	if result.Warning != nil {
		t.Errorf("Expected no warning for container without spaces, but got: %v", result.Warning.Warning)
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
	result := ruleMultipleContainers(module)
	if result.Warning == nil {
		t.Fatal("Expected warning for mixed container syntax, but got none")
	}

	expectedWarning := "Docker and Singularity containers specified on the same line"
	if result.Warning.Warning != expectedWarning {
		t.Errorf("Expected warning message %q but got %q", expectedWarning, result.Warning.Warning)
	}
	if result.Warning.ModulePath != processFile {
		t.Errorf("Expected module path %q but got %q", processFile, result.Warning.ModulePath)
	}
}
