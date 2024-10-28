package nf

import (
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
