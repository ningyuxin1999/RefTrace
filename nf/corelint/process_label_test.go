package corelint

import (
	"os"
	"path/filepath"
	"reft-go/nf"
	"testing"
)

func TestRuleConflictingLabels(t *testing.T) {
	tests := []struct {
		name           string
		processContent string
		wantWarn       bool
	}{
		{
			name: "no conflict with single label",
			processContent: `
process FOO {
    label 'process_single'
    script:
    """
    echo "test"
    """
}`,
			wantWarn: false,
		},
		{
			name: "conflict with multiple standard labels",
			processContent: `
process FOO {
    label 'process_single'
    label 'process_high'
    script:
    """
    echo "test"
    """
}`,
			wantWarn: true,
		},
		{
			name: "no conflict with non-standard labels",
			processContent: `
process FOO {
    label 'custom_label'
    label 'another_label'
    script:
    """
    echo "test"
    """
}`,
			wantWarn: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			processFile := filepath.Join(tmpDir, "process.nf")
			if err := os.WriteFile(processFile, []byte(tt.processContent), 0644); err != nil {
				t.Fatal("Failed to write process file:", err)
			}

			// Parse the file
			module, err, _ := nf.BuildModule(processFile)
			if err != nil {
				t.Fatal("Failed to parse process file:", err)
			}

			// Run the rule
			results := ruleConflictingLabels(module)
			if (len(results.Warnings) > 0) != tt.wantWarn {
				t.Errorf("ruleConflictingLabels() got warnings = %v, want warnings = %v", len(results.Warnings) > 0, tt.wantWarn)
			}
		})
	}
}

func TestRuleNoStandardLabel(t *testing.T) {
	tests := []struct {
		name           string
		processContent string
		wantWarn       bool
	}{
		{
			name: "has standard label",
			processContent: `
process FOO {
    label 'process_single'
    script:
    """
    echo "test"
    """
}`,
			wantWarn: false,
		},
		{
			name: "no standard label",
			processContent: `
process FOO {
    label 'custom_label'
    script:
    """
    echo "test"
    """
}`,
			wantWarn: true,
		},
		{
			name: "no labels at all",
			processContent: `
process FOO {
    script:
    """
    echo "test"
    """
}`,
			wantWarn: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			processFile := filepath.Join(tmpDir, "process.nf")
			if err := os.WriteFile(processFile, []byte(tt.processContent), 0644); err != nil {
				t.Fatal("Failed to write process file:", err)
			}

			module, err, _ := nf.BuildModule(processFile)
			if err != nil {
				t.Fatal("Failed to parse process file:", err)
			}

			results := ruleNoStandardLabel(module)
			if (len(results.Warnings) > 0) != tt.wantWarn {
				t.Errorf("ruleNoStandardLabel() got warnings = %v, want warnings = %v", len(results.Warnings) > 0, tt.wantWarn)
			}
		})
	}
}

func TestRuleNonStandardLabel(t *testing.T) {
	tests := []struct {
		name           string
		processContent string
		wantWarn       bool
	}{
		{
			name: "only standard labels",
			processContent: `
process FOO {
    label 'process_single'
    script:
    """
    echo "test"
    """
}`,
			wantWarn: false,
		},
		{
			name: "has non-standard label",
			processContent: `
process FOO {
    label 'process_single'
    label 'custom_label'
    script:
    """
    echo "test"
    """
}`,
			wantWarn: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			processFile := filepath.Join(tmpDir, "process.nf")
			if err := os.WriteFile(processFile, []byte(tt.processContent), 0644); err != nil {
				t.Fatal("Failed to write process file:", err)
			}

			module, err, _ := nf.BuildModule(processFile)
			if err != nil {
				t.Fatal("Failed to parse process file:", err)
			}

			results := ruleNonStandardLabel(module)
			if (len(results.Warnings) > 0) != tt.wantWarn {
				t.Errorf("ruleNonStandardLabel() got warnings = %v, want warnings = %v", len(results.Warnings) > 0, tt.wantWarn)
			}
		})
	}
}

func TestRuleDuplicateLabels(t *testing.T) {
	tests := []struct {
		name           string
		processContent string
		wantWarn       bool
	}{
		{
			name: "no duplicate labels",
			processContent: `
process FOO {
    label 'process_single'
    label 'process_high'
    script:
    """
    echo "test"
    """
}`,
			wantWarn: false,
		},
		{
			name: "has duplicate label",
			processContent: `
process FOO {
    label 'process_single'
    label 'process_single'
    script:
    """
    echo "test"
    """
}`,
			wantWarn: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			processFile := filepath.Join(tmpDir, "process.nf")
			if err := os.WriteFile(processFile, []byte(tt.processContent), 0644); err != nil {
				t.Fatal("Failed to write process file:", err)
			}

			module, err, _ := nf.BuildModule(processFile)
			if err != nil {
				t.Fatal("Failed to parse process file:", err)
			}

			results := ruleDuplicateLabels(module)
			if (len(results.Warnings) > 0) != tt.wantWarn {
				t.Errorf("ruleDuplicateLabels() got warnings = %v, want warnings = %v", len(results.Warnings) > 0, tt.wantWarn)
			}
		})
	}
}

func TestRuleNoLabels(t *testing.T) {
	tests := []struct {
		name           string
		processContent string
		wantWarn       bool
	}{
		{
			name: "has labels",
			processContent: `
process FOO {
    label 'process_single'
    script:
    """
    echo "test"
    """
}`,
			wantWarn: false,
		},
		{
			name: "no labels",
			processContent: `
process FOO {
    script:
    """
    echo "test"
    """
}`,
			wantWarn: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			processFile := filepath.Join(tmpDir, "process.nf")
			if err := os.WriteFile(processFile, []byte(tt.processContent), 0644); err != nil {
				t.Fatal("Failed to write process file:", err)
			}

			module, err, _ := nf.BuildModule(processFile)
			if err != nil {
				t.Fatal("Failed to parse process file:", err)
			}

			results := ruleNoLabels(module)
			if (len(results.Warnings) > 0) != tt.wantWarn {
				t.Errorf("ruleNoLabels() got warnings = %v, want warnings = %v", len(results.Warnings) > 0, tt.wantWarn)
			}
		})
	}
}

func TestRuleAlphanumerics(t *testing.T) {
	tests := []struct {
		name           string
		processContent string
		wantWarn       bool
	}{
		{
			name: "valid alphanumeric labels",
			processContent: `
process FOO {
    label 'process_single'
    label 'label123'
    script:
    """
    echo "test"
    """
}`,
			wantWarn: false,
		},
		{
			name: "invalid characters in label",
			processContent: `
process FOO {
    label 'process-single'
    label 'label@123'
    script:
    """
    echo "test"
    """
}`,
			wantWarn: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			processFile := filepath.Join(tmpDir, "process.nf")
			if err := os.WriteFile(processFile, []byte(tt.processContent), 0644); err != nil {
				t.Fatal("Failed to write process file:", err)
			}

			module, err, _ := nf.BuildModule(processFile)
			if err != nil {
				t.Fatal("Failed to parse process file:", err)
			}

			results := ruleAlphanumerics(module)
			if (len(results.Warnings) > 0) != tt.wantWarn {
				t.Errorf("ruleAlphanumerics() got warnings = %v, want warnings = %v", len(results.Warnings) > 0, tt.wantWarn)
			}
		})
	}
}
