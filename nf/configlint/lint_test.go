package configlint

import (
	"os"
	"path/filepath"
	"reft-go/parser"
	"testing"
)

func TestParamsVariableDetection(t *testing.T) {
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

	// Create and use ParamVisitor to find params
	paramVisitor := NewParamVisitor()
	paramVisitor.VisitBlockStatement(ast.StatementBlock)
	params := paramVisitor.GetSortedParams()

	// Expected params with their line numbers and directive context
	expected := []ParamInfo{
		{Name: "skip_bbsplit", LineNumber: 2, InDirective: false},
		{Name: "bbsplit_fasta_list", LineNumber: 2, InDirective: false},
		{Name: "outdir", LineNumber: 7, InDirective: true, DirectiveName: "publishDir"},
		{Name: "publish_dir_mode", LineNumber: 8, InDirective: true, DirectiveName: "publishDir"},
		{Name: "save_reference", LineNumber: 10, InDirective: true, DirectiveName: "publishDir"},
	}

	// Compare results with additional directive context checks
	for i, param := range params {
		if i >= len(expected) {
			break
		}
		if param.Name != expected[i].Name {
			t.Errorf("Expected param name %s, got %s", expected[i].Name, param.Name)
		}
		if param.LineNumber != expected[i].LineNumber {
			t.Errorf("For param %s: expected line number %d, got %d",
				param.Name, expected[i].LineNumber, param.LineNumber)
		}
		if param.InDirective != expected[i].InDirective {
			t.Errorf("For param %s: expected InDirective %v, got %v",
				param.Name, expected[i].InDirective, param.InDirective)
		}
		if param.DirectiveName != expected[i].DirectiveName {
			t.Errorf("For param %s: expected DirectiveName %s, got %s",
				param.Name, expected[i].DirectiveName, param.DirectiveName)
		}
	}
}
