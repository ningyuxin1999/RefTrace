package nf

import (
	"os"
	"path/filepath"
	"reft-go/parser"
	"strings"
	"testing"
)

func getTestDataDir() string {
	if dir := os.Getenv("NF_CORE_TEST_DATA"); dir != "" {
		return dir
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		// If we can't get the home directory, fall back to "testdata"
		return "testdata"
	}

	return filepath.Join(homeDir, "reft-testdata")
}

func TestSarekEntireMain(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "sarek_entire_main.nf")
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}
	paramVisitor := NewParamVisitor()
	paramVisitor.VisitBlockStatement(ast.StatementBlock)
	params := paramVisitor.GetSortedParams()
	if len(params) != 65 {
		t.Fatalf("Expected 65 params, got %d", len(params))
	}
	includeVisitor := NewIncludeVisitor()
	includeVisitor.VisitBlockStatement(ast.StatementBlock)
	includes := includeVisitor.GetSortedIncludes()
	if len(includes) != 8 {
		t.Fatalf("Expected 8 includes, got %d", len(includes))
	}
	workflowVisitor := NewWorkflowVisitor()
	workflowVisitor.VisitBlockStatement(ast.StatementBlock)
	workflows := workflowVisitor.workflows
	if len(workflows) != 2 {
		t.Fatalf("Expected 2 workflows, got %d", len(workflows))
	}
}

func TestSimpleWorkflow(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-testdata", "simple_workflow.nf")

	module, err, _ := BuildModule(filePath)
	if err != nil {
		t.Fatalf("Failed to build module: %v", err)
	}

	if len(module.Workflows) != 2 {
		t.Fatalf("Expected 2 workflows, got %d", len(module.Workflows))
	}
	workflow := module.Workflows[0]
	if len(workflow.Takes) != 2 {
		t.Fatalf("Expected 2 takes, got %d", len(workflow.Takes))
	}
	if len(workflow.Emits) != 2 {
		t.Fatalf("Expected 2 emits, got %d", len(workflow.Emits))
	}
}

func TestSimpleProcess(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-testdata", "simple_process.nf")
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}

	processVisitor := NewProcessVisitor()
	processVisitor.VisitBlockStatement(ast.StatementBlock)
	processes := processVisitor.processes
	if len(processes) != 1 {
		t.Fatalf("Expected 1 process, got %d", len(processes))
	}
	directives := processes[0].Directives
	if len(directives) != 51 {
		t.Fatalf("Expected 51 directives, got %d", len(directives))
	}
}

func TestClosureDirective(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-testdata", "closure_directive.nf")
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}

	processVisitor := NewProcessVisitor()
	processVisitor.VisitBlockStatement(ast.StatementBlock)
	processes := processVisitor.processes
	if len(processes) != 1 {
		t.Fatalf("Expected 1 process, got %d", len(processes))
	}
	directives := processes[0].Directives
	if len(directives) != 2 {
		t.Fatalf("Expected 2 directives, got %d", len(directives))
	}
}

func TestChannelFromPath(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-testdata", "channel_frompath.nf")
	_, err := parser.BuildAST(filePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}
	// cls := ast.GetClasses()[0]
	// TODO: fix this test
}

func TestCountProcesses(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-core")
	_, err := ProcessDirectory(filePath)
	if err == nil {
		t.Fatal("Expected error when processing directory, but got none")
	}

	testDataDir := getTestDataDir()
	expectedError := "encountered 2 errors: [" +
		filepath.Join(testDataDir, "nf-core/clipseq/main.nf") + ": errors found in processes in " +
		filepath.Join(testDataDir, "nf-core/clipseq/main.nf") + ": process 'generate_star_index': invalid publish dir directive: no valid path specified; process 'generate_star_index_no_gtf': invalid publish dir directive: no valid path specified " +
		filepath.Join(testDataDir, "nf-core/eager/main.nf") + ": only DSL2 scripts are supported. Found explicit DSL1 declaration in " +
		filepath.Join(testDataDir, "nf-core/eager/main.nf") + "]"

	if err.Error() != expectedError {
		t.Errorf("Expected error:\n%s\n\nGot:\n%s", expectedError, err.Error())
	}
}

func TestCountProcessesAirrflow(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-core", "airrflow")
	modules, err := ProcessDirectory(filePath)
	if err != nil {
		t.Fatalf("Failed to process directory: %v", err)
	}

	totalProcesses := 0
	for _, module := range modules {
		totalProcesses += len(module.Processes)
	}

	t.Logf("Found %d modules with a total of %d processes", len(modules), totalProcesses)
}

func TestCountProcessesGWAS(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-core", "gwas")
	modules, err := ProcessDirectory(filePath)
	if err != nil {
		t.Fatalf("Failed to process directory: %v", err)
	}

	totalProcesses := 0
	for _, module := range modules {
		totalProcesses += len(module.Processes)
	}

	t.Logf("Found %d modules with a total of %d processes", len(modules), totalProcesses)
}

func TestCountProcessesPathogenSurveillance(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-core", "pathogensurveillance")
	modules, err := ProcessDirectory(filePath)
	if err != nil {
		t.Fatalf("Failed to process directory: %v", err)
	}

	totalProcesses := 0
	for _, module := range modules {
		totalProcesses += len(module.Processes)
	}

	t.Logf("Found %d modules with a total of %d processes", len(modules), totalProcesses)
}

func TestFetchNGS(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-core", "fetchngs")
	modules, err := ProcessDirectory(filePath)
	if err != nil {
		t.Fatalf("Failed to process directory: %v", err)
	}

	totalProcesses := 0
	for _, module := range modules {
		totalProcesses += len(module.Processes)
	}

	t.Logf("Found %d modules with a total of %d processes", len(modules), totalProcesses)
}

func TestPGDB(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-core", "pgdb")
	modules, err := ProcessDirectory(filePath)
	if err != nil {
		t.Fatalf("Failed to process directory: %v", err)
	}

	totalProcesses := 0
	for _, module := range modules {
		totalProcesses += len(module.Processes)
	}

	t.Logf("Found %d modules with a total of %d processes", len(modules), totalProcesses)
}

func TestIfStatementProcess(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-core", "airrflow/modules/local/enchantr/report_file_size.nf")
	module, err, _ := BuildModule(filePath)
	if err != nil {
		t.Fatalf("Failed to process file: %v", err)
	}
	if len(module.Processes) != 1 {
		t.Fatalf("Expected 1 process, got %d", len(module.Processes))
	}
	process := module.Processes[0]
	if len(process.Directives) != 4 {
		t.Fatalf("Expected 4 directives, got %d", len(process.Directives))
	}
}

func TestRejectDSL1(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-testdata", "dsl1.nf")
	_, err, _ := BuildModule(filePath)
	if err == nil {
		t.Fatal("Expected error when processing explicit DSL1 script, but got none")
	}

	expectedError := "only DSL2 scripts are supported"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error containing '%s', but got: %v", expectedError, err)
	}
}
