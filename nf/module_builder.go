package nf

import (
	"fmt"
	"reft-go/parser"
	"strings"

	"go.starlark.net/starlark"
)

type Module struct {
	Path       string
	Processes  []Process
	Includes   []IncludeStatement
	DSLVersion int
}

func BuildModule(filePath string) (*Module, error) {
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

	includeVisitor := NewIncludeVisitor()
	includeVisitor.VisitBlockStatement(ast.StatementBlock)
	includes := includeVisitor.includes

	processVisitor := NewProcessVisitor()
	processVisitor.VisitBlockStatement(ast.StatementBlock)
	processes := processVisitor.processes

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

func ConvertToStarlarkModule(m *Module) *StarlarkModule {
	starlarkProcesses := make([]*StarlarkProcess, len(m.Processes))
	for i, process := range m.Processes {
		starlarkProcesses[i] = ConvertToStarlarkProcess(process)
	}

	return &StarlarkModule{
		Path:      m.Path,
		Processes: starlarkProcesses,
		Includes:  m.Includes,
	}
}

var _ starlark.Value = (*StarlarkModule)(nil)
var _ starlark.HasAttrs = (*StarlarkModule)(nil)

type StarlarkModule struct {
	Path      string
	Processes []*StarlarkProcess
	Includes  []IncludeStatement
}

func (m *StarlarkModule) String() string {
	return fmt.Sprintf("Module(%s)", m.Path)
}

func (m *StarlarkModule) Type() string {
	return "module"
}

func (m *StarlarkModule) Freeze() {}

func (m *StarlarkModule) Truth() starlark.Bool {
	return starlark.Bool(true)
}

func (m *StarlarkModule) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: module")
}

func (m *StarlarkModule) Attr(name string) (starlark.Value, error) {
	switch name {
	case "path":
		return starlark.String(m.Path), nil
	case "processes":
		processes := make([]starlark.Value, len(m.Processes))
		for i, p := range m.Processes {
			processes[i] = p
		}
		return starlark.NewList(processes), nil
	case "includes":
		includes := make([]starlark.Value, len(m.Includes))
		for i, inc := range m.Includes {
			includes[i] = inc
		}
		return starlark.NewList(includes), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("module has no attribute %q", name))
	}
}

func (m *StarlarkModule) AttrNames() []string {
	return []string{"path", "processes", "includes"}
}
