package configlint

import (
	"os"
	"path/filepath"
	"reft-go/parser"
	"testing"
)

// Tests the parsing of Nextflow config files
func TestParseConfig(t *testing.T) {
	// Read the fixture file
	testCase, err := os.ReadFile("testdata/modules.config")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	// Write to temporary test file
	testFilePath := filepath.Join(t.TempDir(), "test.config")
	if err := os.WriteFile(testFilePath, testCase, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Parse the file using BuildAST
	ast, err := parser.BuildAST(testFilePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}

	processScopes := parseConfig(ast.StatementBlock)

	// Basic validation
	if len(processScopes) == 0 {
		t.Error("Expected to find process scope declarations, but found none")
	}

	// Expected process declarations (line numbers where 'process' appears)
	expectedLines := []int{19, 47, 156, 173, 202, 210, 229, 260, 287, 317, 338, 363, 403, 443, 497, 543, 564, 636, 758, 783, 822, 847, 879, 899, 911, 945, 970, 982, 994, 1028, 1047, 1071, 1083, 1108, 1121, 1141, 1177}

	// Compare results
	if len(processScopes) != len(expectedLines) {
		t.Errorf("Expected %d process scope declarations, got %d", len(expectedLines), len(processScopes))
	}

	// Check line numbers
	foundLines := make([]int, len(processScopes))
	for i, processScope := range processScopes {
		foundLines[i] = processScope.LineNumber
	}

	// Compare line numbers (they should match exactly)
	for i, expected := range expectedLines {
		if i >= len(foundLines) {
			t.Errorf("Missing expected process scope declaration at line %d", expected)
			continue
		}
		if foundLines[i] != expected {
			t.Errorf("Expected process scope declaration at line %d, got line %d", expected, foundLines[i])
		}
	}

}

func TestParseConfigConditional(t *testing.T) {
	// Create test file content
	testCase := `
if (!params.skip_bbsplit && params.bbsplit_fasta_list) {
    process {
        withName: '.*:PREPARE_GENOME:BBMAP_BBSPLIT' {
            ext.args   = 'build=1'
            publishDir = [
                path: { "${params.outdir}/genome/index" },
                mode: params.publish_dir_mode,
                saveAs: { filename -> filename.equals('versions.yml') ? null : filename },
                enabled: params.save_reference
            ]
        }
    }
}
`
	// Write test file and get its path
	testFilePath := filepath.Join(t.TempDir(), "test.config")
	if err := os.WriteFile(testFilePath, []byte(testCase), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Parse the file using BuildAST
	ast, err := parser.BuildAST(testFilePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}

	processScopes := parseConfig(ast.StatementBlock)
	if len(processScopes) != 1 {
		t.Errorf("Expected 1 process scope declaration, got %d", len(processScopes))
	}
}

func TestParseConfigDuplicateArg(t *testing.T) {
	testCase := `
        process {
            withName: '.*:FASTQ_FASTQC_UMITOOLS_TRIMGALORE:TRIMGALORE' {
                ext.args   = {
                    [
                        "--fastqc_args '-t ${task.cpus}'",
                        params.extra_trimgalore_args ? params.extra_trimgalore_args.split("\\s(?=--)") : ''
                    ].flatten().unique(false).join(' ').trim()
                }
        }
    }`

	// Write test file and get its path
	testFilePath := filepath.Join(t.TempDir(), "test.config")
	if err := os.WriteFile(testFilePath, []byte(testCase), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Parse the file using BuildAST
	ast, err := parser.BuildAST(testFilePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}

	processScopes := parseConfig(ast.StatementBlock)
	if len(processScopes) != 1 {
		t.Fatalf("Expected 1 process scope declaration, got %d", len(processScopes))
	}
	ps := processScopes[0]
	if len(ps.NamedScopes) != 1 {
		t.Fatalf("Expected 1 named scope, got %d", len(ps.NamedScopes))
	}
	namedScope := ps.NamedScopes[0]
	if namedScope.Name != ".*:FASTQ_FASTQC_UMITOOLS_TRIMGALORE:TRIMGALORE" {
		t.Fatalf("Expected named scope name to be '.*:FASTQ_FASTQC_UMITOOLS_TRIMGALORE:TRIMGALORE', got %s", namedScope.Name)
	}
	if len(namedScope.Directives) != 1 {
		t.Fatalf("Expected 1 directive, got %d", len(namedScope.Directives))
	}
	directive := namedScope.Directives[0]
	if directive.Name != "ext.args" {
		t.Fatalf("Expected directive name to be 'ext.args', got %s", directive.Name)
	}
	value := directive.Value
	if len(value.Params) != 1 {
		t.Fatalf("Expected 1 param, got %d", len(value.Params))
	}
	if value.Params[0] != "extra_trimgalore_args" {
		t.Fatalf("Expected param to be 'extra_trimgalore_args', got %s", value.Params[0])
	}
}
