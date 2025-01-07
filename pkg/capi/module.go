package main

// #include <stdlib.h>
import "C"
import (
	"encoding/base64"
	"fmt"
	"reft-go/nf"
	pb "reft-go/nf/proto"
	"reft-go/parser"
	"strings"
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

//export Module_Free
func Module_Free(ptr *C.char) {
	C.free(unsafe.Pointer(ptr))
}

func BuildModuleInternal(filePath string) (*Module, error, bool) {
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		if _, ok := err.(*parser.SyntaxException); ok {
			return nil, err, true
		}
		return nil, err, false
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
