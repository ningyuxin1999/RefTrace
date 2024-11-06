package nf

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.starlark.net/starlark"
)

func TestMinMaxCpus(t *testing.T) {
	baseDir := filepath.Join(getTestDataDir(), "linting-tests", "min_max_cpus")
	config := LintConfig{
		RulesFile: filepath.Join(baseDir, "rules.py"),
		Directory: filepath.Join(baseDir, "process.nf"),
		RuleToRun: "",
	}

	// Create a writer to capture all output
	var output strings.Builder

	// Capture the Starlark print output
	var starlarkOutput strings.Builder
	originalPrint := starlark.Universe["print"].(*starlark.Builtin)
	starlark.Universe["print"] = starlark.NewBuiltin("print", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		msg := args[0].String()
		starlarkOutput.WriteString(msg + "\n")
		return originalPrint.CallInternal(thread, args, kwargs)
	})
	defer func() {
		starlark.Universe["print"] = originalPrint
	}()

	// Run the linting process
	err := RunLintWithConfig(config, &output)
	if err == nil {
		t.Fatal("Expected linting to fail due to invalid CPU value, but it succeeded")
	}

	// Check both outputs
	combinedOutput := output.String() + starlarkOutput.String()
	expectedError := "Error: Process FOO has an invalid CPU value. It should be >= 2 and <= 96, but it is 100"
	if !strings.Contains(combinedOutput, expectedError) {
		t.Errorf("Expected output to contain error message about invalid CPU value, but got:\n%s", combinedOutput)
	}
}

func TestUnknownDirective(t *testing.T) {
	baseDir := filepath.Join(getTestDataDir(), "linting-tests", "unknown_directive")
	config := LintConfig{
		RulesFile: filepath.Join(baseDir, "rules.py"),
		Directory: filepath.Join(baseDir, "process.nf"),
		RuleToRun: "",
	}

	// Create a writer to capture all output
	var output strings.Builder

	// Capture the Starlark print output
	var starlarkOutput strings.Builder
	originalPrint := starlark.Universe["print"].(*starlark.Builtin)
	starlark.Universe["print"] = starlark.NewBuiltin("print", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		msg := args[0].String()
		starlarkOutput.WriteString(msg + "\n")
		return originalPrint.CallInternal(thread, args, kwargs)
	})
	defer func() {
		starlark.Universe["print"] = originalPrint
	}()

	// Run the linting process
	err := RunLintWithConfig(config, &output)
	if err == nil {
		t.Fatal("Expected linting to fail due to unknown directive, but it succeeded")
	}

	// Check both outputs
	combinedOutput := output.String() + starlarkOutput.String()
	expectedError := "Error: Unknown directive found: bar"
	if !strings.Contains(combinedOutput, expectedError) {
		t.Errorf("Expected output to contain error message about unknown directive, but got:\n%s", combinedOutput)
	}
}

func TestNoLabelOrMemory(t *testing.T) {
	baseDir := filepath.Join(getTestDataDir(), "linting-tests", "no_label_or_memory")
	config := LintConfig{
		RulesFile: filepath.Join(baseDir, "rules.py"),
		Directory: filepath.Join(baseDir, "process.nf"),
		RuleToRun: "",
	}

	// Create a writer to capture all output
	var output strings.Builder

	// Capture the Starlark print output
	var starlarkOutput strings.Builder
	originalPrint := starlark.Universe["print"].(*starlark.Builtin)
	starlark.Universe["print"] = starlark.NewBuiltin("print", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		msg := args[0].String()
		starlarkOutput.WriteString(msg + "\n")
		return originalPrint.CallInternal(thread, args, kwargs)
	})
	defer func() {
		starlark.Universe["print"] = originalPrint
	}()

	// Run the linting process
	err := RunLintWithConfig(config, &output)
	if err == nil {
		t.Fatal("Expected linting to fail due to missing label and memory directives, but it succeeded")
	}

	// Check both outputs
	combinedOutput := output.String() + starlarkOutput.String()
	expectedError := "Error: process FOO has no label or memory directive"
	if !strings.Contains(combinedOutput, expectedError) {
		t.Errorf("Expected output to contain error message about missing label and memory directives, but got:\n%s", combinedOutput)
	}
}

func TestExecEnvironment(t *testing.T) {
	baseDir := filepath.Join(getTestDataDir(), "linting-tests", "exec_environment")
	config := LintConfig{
		RulesFile: filepath.Join(baseDir, "rules.py"),
		Directory: filepath.Join(baseDir, "process.nf"),
		RuleToRun: "",
	}

	// Create a writer to capture all output
	var output strings.Builder

	// Capture the Starlark print output
	var starlarkOutput strings.Builder
	originalPrint := starlark.Universe["print"].(*starlark.Builtin)
	starlark.Universe["print"] = starlark.NewBuiltin("print", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		msg := args[0].String()
		starlarkOutput.WriteString(msg + "\n")
		return originalPrint.CallInternal(thread, args, kwargs)
	})
	defer func() {
		starlark.Universe["print"] = originalPrint
	}()

	// Run the linting process
	err := RunLintWithConfig(config, &output)
	if err == nil {
		t.Fatal("Expected linting to fail due to missing execution environment, but it succeeded")
	}

	// Check both outputs
	combinedOutput := output.String() + starlarkOutput.String()
	lines := strings.Split(strings.TrimSpace(combinedOutput), "\n")

	errorLines := 0
	expectedError := "Error: Process BAD must specify at least one of: container, label, conda, module, spack"
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Error:") {
			errorLines++
			if line != expectedError {
				t.Fatalf("Expected error message:\n%s\n\nbut got:\n%s", expectedError, line)
			}
		}
	}

	if errorLines != 1 {
		t.Fatalf("Expected exactly one error line, but got %d error lines", errorLines)
	}
}

func TestNoInput(t *testing.T) {
	baseDir := filepath.Join(getTestDataDir(), "linting-tests", "no_input_rule")
	config := LintConfig{
		RulesFile: filepath.Join(baseDir, "rules.py"),
		Directory: filepath.Join(baseDir, "process.nf"),
		RuleToRun: "",
	}

	// Create a writer to capture all output
	var output strings.Builder

	// Capture the Starlark print output
	var starlarkOutput strings.Builder
	originalPrint := starlark.Universe["print"].(*starlark.Builtin)
	starlark.Universe["print"] = starlark.NewBuiltin("print", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		msg := args[0].String()
		starlarkOutput.WriteString(msg + "\n")
		return originalPrint.CallInternal(thread, args, kwargs)
	})
	defer func() {
		starlark.Universe["print"] = originalPrint
	}()

	// Run the linting process
	err := RunLintWithConfig(config, &output)
	if err == nil {
		t.Fatal("Expected linting to fail due to missing inputs, but it succeeded")
	}

	// Check both outputs
	combinedOutput := output.String() + starlarkOutput.String()
	lines := strings.Split(strings.TrimSpace(combinedOutput), "\n")

	errorLines := 0
	expectedError := "Error: Process 'FOO' has no inputs defined"
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Error:") {
			errorLines++
			if line != expectedError {
				t.Fatalf("Expected error message:\n%s\n\nbut got:\n%s", expectedError, line)
			}
		}
	}

	if errorLines != 1 {
		t.Fatalf("Expected exactly one error line, but got %d error lines", errorLines)
	}
}

func TestNoOutput(t *testing.T) {
	baseDir := filepath.Join(getTestDataDir(), "linting-tests", "no_output_rule")
	config := LintConfig{
		RulesFile: filepath.Join(baseDir, "rules.py"),
		Directory: filepath.Join(baseDir, "process.nf"),
		RuleToRun: "",
	}

	// Create a writer to capture all output
	var output strings.Builder

	// Capture the Starlark print output
	var starlarkOutput strings.Builder
	originalPrint := starlark.Universe["print"].(*starlark.Builtin)
	starlark.Universe["print"] = starlark.NewBuiltin("print", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		msg := args[0].String()
		starlarkOutput.WriteString(msg + "\n")
		return originalPrint.CallInternal(thread, args, kwargs)
	})
	defer func() {
		starlark.Universe["print"] = originalPrint
	}()

	// Run the linting process
	err := RunLintWithConfig(config, &output)
	if err == nil {
		t.Fatal("Expected linting to fail due to missing outputs, but it succeeded")
	}

	// Check both outputs
	combinedOutput := output.String() + starlarkOutput.String()
	lines := strings.Split(strings.TrimSpace(combinedOutput), "\n")

	errorLines := 0
	expectedError := "Error: Process 'FOO' has no outputs defined"
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Error:") {
			errorLines++
			if line != expectedError {
				t.Fatalf("Expected error message:\n%s\n\nbut got:\n%s", expectedError, line)
			}
		}
	}

	if errorLines != 1 {
		t.Fatalf("Expected exactly one error line, but got %d error lines", errorLines)
	}
}

func TestCheckGpu(t *testing.T) {
	baseDir := filepath.Join(getTestDataDir(), "linting-tests", "gpu_rule")
	config := LintConfig{
		RulesFile: filepath.Join(baseDir, "rules.py"),
		Directory: filepath.Join(baseDir, "process.nf"),
		RuleToRun: "",
	}

	// Create a writer to capture all output
	var output strings.Builder

	// Capture the Starlark print output
	var starlarkOutput strings.Builder
	originalPrint := starlark.Universe["print"].(*starlark.Builtin)
	starlark.Universe["print"] = starlark.NewBuiltin("print", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		msg := args[0].String()
		starlarkOutput.WriteString(msg + "\n")
		return originalPrint.CallInternal(thread, args, kwargs)
	})
	defer func() {
		starlark.Universe["print"] = originalPrint
	}()

	// Run the linting process
	err := RunLintWithConfig(config, &output)
	if err == nil {
		t.Fatal("Expected linting to fail due to GPU validation errors, but it succeeded")
	}

	// Check both outputs
	combinedOutput := output.String() + starlarkOutput.String()
	lines := strings.Split(strings.TrimSpace(combinedOutput), "\n")

	errorLines := 0
	expectedErrors := []string{
		"Error: Process TOO_SMALL requests invalid number of GPUs: 1. Value must be between 1 and 8",
		"Error: Process TOO_BIG requests invalid number of GPUs: 100. Value must be between 1 and 8",
		"Error: Process FOO specifies invalid GPU type 'foo'. Allowed types are: nvidia-tesla-v100, nvidia-tesla-p100, nvidia-tesla-k80, nvidia-tesla-a100",
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Error:") {
			errorLines++
			found := false
			for _, expectedError := range expectedErrors {
				if line == expectedError {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Unexpected error message: %s", line)
			}
		}
	}

	if errorLines != 3 {
		t.Errorf("Expected exactly three error lines, but got %d error lines", errorLines)
	}
}

func TestBadMemory(t *testing.T) {
	baseDir := filepath.Join(getTestDataDir(), "linting-tests", "bad_memory")
	config := LintConfig{
		RulesFile: filepath.Join(baseDir, "rules.py"),
		Directory: filepath.Join(baseDir, "process.nf"),
		RuleToRun: "",
	}

	// Create a writer to capture all output
	var output strings.Builder

	// Run the linting process
	err := RunLintWithConfig(config, &output)
	if err == nil {
		t.Fatal("Expected linting to fail due to invalid memory format, but it succeeded")
	}

	expectedErrors := []string{
		"error processing directory: encountered 1 errors",
		"errors found in processes",
		"unknown memory unit: GBB",
	}

	errMsg := err.Error()
	for _, expected := range expectedErrors {
		if !strings.Contains(errMsg, expected) {
			t.Errorf("Expected error to contain: %q\nGot: %v", expected, errMsg)
		}
	}
}

func TestContainerWithSpaceBad(t *testing.T) {
	// Create temporary files for rules and process
	rulesContent := `
def rule_check_container(module):
    for process in module.processes:
        container_directives = process.directives.container
        if not container_directives:
            fatal("Process '%s' must have a container directive" % process.name)
        for container in container_directives:
            if " " in container.simple_name:
                fatal("Container name '%s' contains spaces, which is not allowed" % container.simple_name)
`
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
	// Create a writer to capture all output
	var output strings.Builder

	// Create temporary directory and files
	tmpDir := t.TempDir()
	rulesFile := filepath.Join(tmpDir, "rules.py")
	processFile := filepath.Join(tmpDir, "process.nf")

	if err := os.WriteFile(rulesFile, []byte(rulesContent), 0644); err != nil {
		t.Fatal("Failed to write rules file:", err)
	}
	if err := os.WriteFile(processFile, []byte(processContent), 0644); err != nil {
		t.Fatal("Failed to write process file:", err)
	}

	config := LintConfig{
		RulesFile: rulesFile,
		Directory: processFile,
		RuleToRun: "",
	}

	// Run the linting process
	err := RunLintWithConfig(config, &output)
	if err == nil {
		t.Fatal("Expected linting to fail due to missing container directive, but it succeeded")
	}

	// Check both outputs
	//combinedOutput := output.String() + starlarkOutput.String()
	expectedError := "Error: Container name 'ubuntu latest' contains spaces, which is not allowed"
	if !strings.Contains(output.String(), expectedError) {
		t.Errorf("Expected output to contain error message about container name with spaces, but got:\n%s", output.String())
	}
}

func TestContainerWithoutSpaceGood(t *testing.T) {
	// Create temporary files for rules and process
	rulesContent := `
def rule_check_container(module):
    for process in module.processes:
        container_directives = process.directives.container
        if not container_directives:
            fatal("Process '%s' must have a container directive" % process.name)
        for container in container_directives:
            if " " in container.simple_name:
                fatal("Container name '%s' contains spaces, which is not allowed" % container.simple_name)
`
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
	// Create a writer to capture all output
	var output strings.Builder

	// Create temporary directory and files
	tmpDir := t.TempDir()
	rulesFile := filepath.Join(tmpDir, "rules.py")
	processFile := filepath.Join(tmpDir, "process.nf")

	if err := os.WriteFile(rulesFile, []byte(rulesContent), 0644); err != nil {
		t.Fatal("Failed to write rules file:", err)
	}
	if err := os.WriteFile(processFile, []byte(processContent), 0644); err != nil {
		t.Fatal("Failed to write process file:", err)
	}

	config := LintConfig{
		RulesFile: rulesFile,
		Directory: processFile,
		RuleToRun: "",
	}

	// Run the linting process
	err := RunLintWithConfig(config, &output)
	if err != nil {
		t.Fatalf("Expected linting to succeed with valid container name, but got error: %v\nOutput: %s", err, output.String())
	}

	// Check the output
	outputStr := output.String()
	expectedOutput := "\nRule: check_container\n\n"
	if outputStr != expectedOutput {
		t.Errorf("Expected output:\n%q\nbut got:\n%q", expectedOutput, outputStr)
	}
}

func TestContainerWithTooManyQuotes(t *testing.T) {
	rulesContent := `
def rule_check_container(module):
    for process in module.processes:
        container_directives = process.directives.container
        if not container_directives:
            fatal("Process '%s' must have a container directive" % process.name)
        for container in container_directives:
            fatal(container.simple_name)
            if container.simple_name.count('"') > 2:
                fatal("Too many double quotes found when specifying container: %s" % container.simple_name)
`
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
	// Create a writer to capture all output
	var output strings.Builder

	// Create temporary directory and files
	tmpDir := t.TempDir()
	rulesFile := filepath.Join(tmpDir, "rules.py")
	processFile := filepath.Join(tmpDir, "process.nf")

	if err := os.WriteFile(rulesFile, []byte(rulesContent), 0644); err != nil {
		t.Fatal("Failed to write rules file:", err)
	}
	if err := os.WriteFile(processFile, []byte(processContent), 0644); err != nil {
		t.Fatal("Failed to write process file:", err)
	}

	config := LintConfig{
		RulesFile: rulesFile,
		Directory: processFile,
		RuleToRun: "",
	}

	// Run the linting process
	err := RunLintWithConfig(config, &output)
	if err == nil {
		t.Fatal("Expected linting to fail due to too many double quotes, but it succeeded")
	}

	expectedErrors := []string{
		"error processing directory",
		"too many quotes found when specifying container",
	}

	errMsg := err.Error()
	for _, expected := range expectedErrors {
		if !strings.Contains(errMsg, expected) {
			t.Errorf("Expected error to contain: %q\nGot: %v", expected, errMsg)
		}
	}
}

func TestContainerWithMultipleDefinitions(t *testing.T) {
	rulesContent := `
def rule_check_container(module):
    for process in module.processes:
        container_directives = process.directives.container
        if not container_directives:
            fatal("Process '%s' must have a container directive" % process.name)
        for container in container_directives:
            if "https://containers" in container.simple_name and "biocontainers/" in container.simple_name:
                fatal("Docker and Singularity containers specified in the same line: %s" % container.simple_name)
`
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
	// Create a writer to capture all output
	var output strings.Builder

	// Create temporary directory and files
	tmpDir := t.TempDir()
	rulesFile := filepath.Join(tmpDir, "rules.py")
	processFile := filepath.Join(tmpDir, "process.nf")

	if err := os.WriteFile(rulesFile, []byte(rulesContent), 0644); err != nil {
		t.Fatal("Failed to write rules file:", err)
	}
	if err := os.WriteFile(processFile, []byte(processContent), 0644); err != nil {
		t.Fatal("Failed to write process file:", err)
	}

	config := LintConfig{
		RulesFile: rulesFile,
		Directory: processFile,
		RuleToRun: "",
	}

	// Run the linting process
	err := RunLintWithConfig(config, &output)
	if err == nil {
		t.Fatal("Expected linting to fail due to mixing Docker and Singularity containers, but it succeeded")
	}

	// Check the output
	outputStr := output.String()
	expectedError := "Error: Docker and Singularity containers specified in the same line: https://containers/biocontainers/ubuntu"
	if !strings.Contains(outputStr, expectedError) {
		t.Errorf("Expected output to contain error message about mixed container types, but got:\n%s", outputStr)
	}
}

func TestContainerWithConditionalSpecification(t *testing.T) {
	rulesContent := `
def rule_check_container(module):
    for process in module.processes:
        container_directives = process.directives.container
        if not container_directives:
            fatal("Process '%s' must have a container directive" % process.name)
        for container in container_directives:
            if container.format == "simple":
                if " " in container.simple_name:
                    fatal("Container name '%s' contains spaces, which is not allowed" % container.simple_name)
            elif container.format == "ternary":
                if " " in container.true_name:
                    fatal("Container true_name '%s' contains spaces, which is not allowed" % container.true_name)
                if " " in container.false_name:
                    fatal("Container false_name '%s' contains spaces, which is not allowed" % container.false_name)
`
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
	// Create a writer to capture all output
	var output strings.Builder

	// Create temporary directory and files
	tmpDir := t.TempDir()
	rulesFile := filepath.Join(tmpDir, "rules.py")
	processFile := filepath.Join(tmpDir, "process.nf")

	if err := os.WriteFile(rulesFile, []byte(rulesContent), 0644); err != nil {
		t.Fatal("Failed to write rules file:", err)
	}
	if err := os.WriteFile(processFile, []byte(processContent), 0644); err != nil {
		t.Fatal("Failed to write process file:", err)
	}

	config := LintConfig{
		RulesFile: rulesFile,
		Directory: processFile,
		RuleToRun: "",
	}

	// Run the linting process
	err := RunLintWithConfig(config, &output)
	if err != nil {
		t.Fatalf("Expected linting to succeed with conditional container specification, but got error: %v\nOutput: %s", err, output.String())
	}

	// Check the output
	outputStr := output.String()
	expectedOutput := "\nRule: check_container\n\n"
	if outputStr != expectedOutput {
		t.Errorf("Expected output:\n%q\nbut got:\n%q", expectedOutput, outputStr)
	}
}
