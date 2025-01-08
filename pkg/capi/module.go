package main

// #include <stdlib.h>
// typedef void (*callback_func)(int32_t, int32_t);
// static inline void CallbackFunc(void* f, int32_t current, int32_t total) {
//     ((callback_func)f)(current, total);
// }
import "C"
import (
	"encoding/base64"
	"fmt"
	"io/fs"
	"path/filepath"
	"reft-go/nf"
	pb "reft-go/nf/proto"
	"reft-go/parser"
	"strings"
	"sync"
	"unsafe"

	"google.golang.org/protobuf/proto"
)

func main() {} // Required for C shared library

type Module struct {
	Path       string
	Processes  []nf.Process
	Includes   []nf.IncludeStatement
	DSLVersion int
}

func (m *Module) ToProto() *pb.Module {
	protoModule := &pb.Module{
		Path:       m.Path,
		DslVersion: int32(m.DSLVersion),
	}

	for _, p := range m.Processes {
		protoModule.Processes = append(protoModule.Processes, p.ToProto())
	}

	return protoModule
}

//export Module_New
func Module_New(filePath *C.char) *C.char {
	goPath := C.GoString(filePath)
	module, err, likelyBug := BuildModuleInternal(goPath)

	result := &pb.ModuleResult{}
	parseError := &pb.ParseError{}
	if err != nil {
		parseError.LikelyRtBug = likelyBug
		parseError.Error = err.Error()
		result.Result = &pb.ModuleResult_Error{Error: parseError}
	} else {
		result.Result = &pb.ModuleResult_Module{Module: module.ToProto()}
	}

	bytes, err := proto.Marshal(result)
	if err != nil {
		panic("serialization error: " + err.Error())
	}

	return C.CString(base64.StdEncoding.EncodeToString(bytes))
}

type ProgressCallback func(int32, int32)

type ModuleResult struct {
	Module *Module
	Error  error
	Path   string
}

func ProcessDirectory(dir string, callback ProgressCallback) ([]ModuleResult, error) {
	var results []ModuleResult
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Count total .nf files first
	var totalFiles int32 = 0
	var processedFiles int32 = 0

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".nf" {
			totalFiles++
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %v", err)
	}

	// Now process the files
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".nf" {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				module, err, _ := BuildModuleInternal(path)

				mu.Lock()
				results = append(results, ModuleResult{
					Module: module,
					Error:  err,
					Path:   path,
				})
				processedFiles++
				if callback != nil {
					callback(processedFiles, totalFiles)
				}
				mu.Unlock()
			}(path)
		}
		return nil
	})

	wg.Wait()

	if err != nil {
		return nil, fmt.Errorf("error during processing: %v", err)
	}

	return results, nil
}

var totalFiles int32
var processedFiles int32

// processes nextflow modules in a directory

//export Parse_Modules
func Parse_Modules(dir *C.char, callback unsafe.Pointer) *C.char {
	goDir := C.GoString(dir)

	var progressCallback ProgressCallback
	if callback != nil {
		progressCallback = func(current, total int32) {
			C.CallbackFunc(callback, C.int32_t(current), C.int32_t(total))
		}
	}

	results, err := ProcessDirectory(goDir, progressCallback)

	listResult := &pb.ModuleListResult{}
	if err != nil {
		// Only if we couldn't even process the directory
		return C.CString(fmt.Sprintf("error processing directory: %v", err))
	}

	for _, res := range results {
		moduleResult := &pb.ModuleResult{
			FilePath: res.Path,
		}

		if res.Error != nil {
			moduleResult.Result = &pb.ModuleResult_Error{
				Error: &pb.ParseError{
					Error:       res.Error.Error(),
					LikelyRtBug: false,
				},
			}
		} else {
			moduleResult.Result = &pb.ModuleResult_Module{
				Module: res.Module.ToProto(),
			}
		}

		listResult.Results = append(listResult.Results, moduleResult)
	}

	bytes, err := proto.Marshal(listResult)
	if err != nil {
		panic("serialization error: " + err.Error())
	}

	return C.CString(base64.StdEncoding.EncodeToString(bytes))
}

//export Module_Free
func Module_Free(ptr *C.char) {
	C.free(unsafe.Pointer(ptr))
}

func BuildModuleInternal(filePath string) (*Module, error, bool) {
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		if _, ok := err.(*parser.SyntaxException); ok {
			return nil, err, false
		}
		return nil, err, true
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
		return nil, fmt.Errorf("only DSL2 scripts are supported. Found explicit DSL1 declaration in %s", filePath), false
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
		return nil, fmt.Errorf("errors found in processes in %s: %s", filePath, strings.Join(processErrors, "; ")), false
	}

	return &Module{
		Path:       filePath,
		Processes:  processes,
		Includes:   includes,
		DSLVersion: dslVersion,
	}, nil, false
}
