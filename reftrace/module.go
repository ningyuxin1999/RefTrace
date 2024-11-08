package main

/*
#include <stdlib.h>

struct ModuleNewResult {
    unsigned long long handle;
    char* error;
};
*/
import "C"
import (
	"fmt"
	"reft-go/nf"
	"reft-go/parser"
	"strings"
	"unsafe"
)

type Module struct {
	Path       string
	Processes  []nf.Process
	Includes   []nf.IncludeStatement
	DSLVersion int
}

//export Module_New
func Module_New(filePath *C.char) C.struct_ModuleNewResult {
	goPath := C.GoString(filePath)
	module, err := BuildModuleInternal(goPath)
	if err != nil {
		errStr := C.CString(err.Error())
		return C.struct_ModuleNewResult{
			handle: 0,
			error:  errStr,
		}
	}

	handle := nextModuleHandle
	nextModuleHandle++
	moduleStore[handle] = module

	return C.struct_ModuleNewResult{
		handle: C.ulonglong(handle),
		error:  nil,
	}
}

//export Module_Free_Error
func Module_Free_Error(cstr *C.char) {
	C.free(unsafe.Pointer(cstr))
}

//export Module_Free
func Module_Free(handle ModuleHandle) {
	if module, ok := moduleStore[handle]; ok {
		// Free all processes associated with this module
		for procHandle, proc := range processStore {
			for i := range module.Processes {
				if proc == &module.Processes[i] {
					Process_Free(procHandle)
					break
				}
			}
		}
		delete(moduleStore, handle)
	}
}

//export Module_GetPath
func Module_GetPath(handle ModuleHandle) *C.char {
	if module, ok := moduleStore[handle]; ok {
		return C.CString(module.Path)
	}
	return nil
}

//export Module_GetDSLVersion
func Module_GetDSLVersion(handle ModuleHandle) C.int {
	if module, ok := moduleStore[handle]; ok {
		return C.int(module.DSLVersion)
	}
	return 0
}

//export Module_GetProcessCount
func Module_GetProcessCount(handle ModuleHandle) C.int {
	if module, ok := moduleStore[handle]; ok {
		return C.int(len(module.Processes))
	}
	return 0
}

//export Module_GetProcess
func Module_GetProcess(moduleHandle ModuleHandle, index C.int) ProcessHandle {
	if module, ok := moduleStore[moduleHandle]; ok {
		if idx := int(index); idx >= 0 && idx < len(module.Processes) {
			// Check if this process already has a handle
			process := &module.Processes[idx]
			for existingHandle, existingProcess := range processStore {
				if existingProcess == process {
					return existingHandle
				}
			}
			// If not found, create new handle
			handle := nextProcessHandle
			nextProcessHandle++
			processStore[handle] = process
			return handle
		}
	}
	return 0
}

func BuildModuleInternal(filePath string) (*Module, error) {
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		return nil, err
	}

	dslVersion := 2
	for _, stmt := range ast.StatementBlock.GetStatements() {
		if expr, ok := stmt.(*parser.ExpressionStatement); ok {
			if binExpr, ok := expr.GetExpression().(*parser.BinaryExpression); ok {
				if binExpr.GetLeftExpression().GetText() == "nextflow.enable.dsl" {
					if constExpr, ok := binExpr.GetRightExpression().(*parser.ConstantExpression); ok {
						if constExpr.GetValue() == 1 {
							dslVersion = 1
							break
						}
					}
				}
			}
		}
	}

	if dslVersion == 1 {
		return nil, fmt.Errorf("only DSL2 scripts are supported. Found explicit DSL1 declaration in %s", filePath)
	}

	includeVisitor := nf.NewIncludeVisitor()
	includeVisitor.VisitBlockStatement(ast.StatementBlock)
	includes := includeVisitor.Includes()

	processVisitor := nf.NewProcessVisitor()
	processVisitor.VisitBlockStatement(ast.StatementBlock)
	processes := processVisitor.Processes()

	// Collect process errors into error message
	var processErrors []string
	hasErrors := false
	for _, process := range processes {
		if len(process.Errors) > 0 {
			hasErrors = true
			for _, err := range process.Errors {
				processErrors = append(processErrors, fmt.Sprintf("process '%s': %v", process.Name, err))
			}
		}
	}

	if hasErrors {
		return nil, fmt.Errorf("errors found in processes in %s: %s", filePath, strings.Join(processErrors, "; "))
	}

	return &Module{
		Path:       filePath,
		Processes:  processes,
		Includes:   includes,
		DSLVersion: dslVersion,
	}, nil
}
