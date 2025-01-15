package nf

import (
	"path/filepath"
	"testing"
)

func TestIncludes(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-testdata", "includes.nf")
	module, err, _ := BuildModule(filePath)
	if err != nil {
		t.Fatalf("Failed to build module: %v", err)
	}
	includes := module.Includes
	if len(includes) != 2 {
		t.Fatalf("Expected 2 includes, got %d", len(includes))
	}
}

func TestIncludes2(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "foo.nf")
	module, err, _ := BuildModule(filePath)
	if err != nil {
		t.Fatalf("Failed to build module: %v", err)
	}
	includes := module.Includes
	if len(includes) != 2 {
		t.Fatalf("Expected 2 includes, got %d", len(includes))
	}
}
