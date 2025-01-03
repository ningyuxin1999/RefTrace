package parser

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"
)

func testGroovyFile(t *testing.T, filePath string) {
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST for %s: %v", filePath, err)
	}

	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	if ast == nil {
		t.Fatalf("Failed to build AST for %s", filePath)
	}
}

func TestGroovyCoreAnnotation01(t *testing.T) {
	filePath := filepath.Join("testdata", "groovy_core", "Annotation_01.groovy")
	testGroovyFile(t, filePath)
}

func TestArray01x(t *testing.T) {
	filePath := filepath.Join("testdata", "groovy_core", "Array_01x.groovy")
	testGroovyFile(t, filePath)
}

func TestBreakingChange01x(t *testing.T) {
	filePath := filepath.Join("testdata", "groovy_core", "BreakingChange_01x.groovy")
	testGroovyFile(t, filePath)
}

func TestClassDeclaration01(t *testing.T) {
	filePath := filepath.Join("testdata", "groovy_core", "ClassDeclaration_01.groovy")
	testGroovyFile(t, filePath)
}

func TestClassDeclarationMinimal(t *testing.T) {
	filePath := filepath.Join("testdata", "groovy_core", "ClassDeclaration_minimal.groovy")
	testGroovyFile(t, filePath)
}

func TestClassDeclaration06(t *testing.T) {
	filePath := filepath.Join("testdata", "groovy_core", "ClassDeclaration_06.groovy")
	testGroovyFile(t, filePath)
}

func TestClassDeclaration07(t *testing.T) {
	filePath := filepath.Join("testdata", "groovy_core", "ClassDeclaration_07.groovy")
	testGroovyFile(t, filePath)
}

func TestEnumDeclaration02(t *testing.T) {
	filePath := filepath.Join("testdata", "groovy_core", "EnumDeclaration_02.groovy")
	testGroovyFile(t, filePath)
}

func TestEnumDeclaration03(t *testing.T) {
	filePath := filepath.Join("testdata", "groovy_core", "EnumDeclaration_03.groovy")
	testGroovyFile(t, filePath)
}

func TestCommand02(t *testing.T) {
	filePath := filepath.Join("testdata", "groovy_core", "Command_02.groovy")
	testGroovyFile(t, filePath)
}

func TestCommand03x(t *testing.T) {
	filePath := filepath.Join("testdata", "groovy_core", "Command_03x.groovy")
	testGroovyFile(t, filePath)
}

func TestGstring02(t *testing.T) {
	filePath := filepath.Join("testdata", "groovy_core", "Gstring_02x.groovy")
	testGroovyFile(t, filePath)
}

func TestInterfaceDeclaration02(t *testing.T) {
	filePath := filepath.Join("testdata", "groovy_core", "InterfaceDeclaration_02.groovy")
	testGroovyFile(t, filePath)
}

func TestWhile01(t *testing.T) {
	filePath := filepath.Join("testdata", "groovy_core", "While_01.groovy")
	testGroovyFile(t, filePath)
}

func TestTryWithResources02x(t *testing.T) {
	filePath := filepath.Join("testdata", "groovy_core", "TryWithResources_02x.groovy")
	testGroovyFile(t, filePath)
}

func TestParseAllTestdataFiles(t *testing.T) {
	// Get testdata directory
	testDir := "testdata"

	totalTests := 0
	passedTests := 0
	failedTests := 0
	panickedTests := 0

	// Walk through all files in testdata directory
	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-.nf/.groovy files
		if info.IsDir() || (!strings.HasSuffix(path, ".nf") && !strings.HasSuffix(path, ".groovy")) {
			return nil
		}

		totalTests++

		// Run as a subtest to get better reporting
		t.Run(path, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					panickedTests++
					t.Errorf("Test panicked for %s: %v\nStack trace:\n%s", path, r, debug.Stack())
				}
			}()

			result, err := BuildCST(path)
			if err != nil {
				failedTests++
				t.Errorf("Failed to build CST for %s: %v", path, err)
				return
			}

			// Build AST to verify full parsing
			builder := NewASTBuilder(path)
			ast := builder.Visit(result.Tree).(*ModuleNode)
			if ast == nil {
				failedTests++
				t.Errorf("Failed to build AST for %s", path)
				return
			}

			passedTests++
			// Log successful parse
			t.Logf("Successfully parsed %s in %s mode", path, result.Mode)
		})

		return nil
	})

	if err != nil {
		t.Fatalf("Error walking testdata directory: %v", err)
	}

	t.Logf("Test Summary: Total: %d, Passed: %d, Failed: %d, Panicked: %d", totalTests, passedTests, failedTests, panickedTests)
}
