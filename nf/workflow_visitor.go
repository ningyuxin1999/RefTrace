package nf

import (
	"fmt"
	pb "reft-go/nf/proto"
	"reft-go/parser"
)

var _ parser.GroovyCodeVisitor = (*WorkflowVisitor)(nil)
var _ parser.GroovyCodeVisitor = (*WorkflowBodyVisitor)(nil)

// WorkflowMode represents the different modes of a workflow
type WorkflowMode int

const (
	// InitialMode represents the initial mode before any labels
	InitialMode WorkflowMode = iota
	// TakeMode represents the 'take' mode of a workflow
	TakeMode
	// MainMode represents the 'main' mode of a workflow
	MainMode
	// EmitMode represents the 'emit' mode of a workflow
	EmitMode
)

type WorkflowBodyVisitor struct {
	*BaseVisitor
	mode    WorkflowMode
	Takes   []string
	Emits   []string
	hasMain bool
	hasTake bool
	errors  []string
}

// NewWorkflowBodyVisitor creates a new WorkflowBodyVisitor
func NewWorkflowBodyVisitor() *WorkflowBodyVisitor {
	v := &WorkflowBodyVisitor{
		BaseVisitor: NewBaseVisitor(),
		mode:        InitialMode,
	}
	v.VisitBlockStatementHook = func(block *parser.BlockStatement) {
		for _, statement := range block.GetStatements() {
			label := statement.GetStatementLabel()
			switch label {
			case "take":
				if v.mode != InitialMode {
					v.errors = append(v.errors, "take: must be the first section in the workflow")
				}
				v.mode = TakeMode
				v.hasTake = true
			case "main":
				if v.mode == EmitMode {
					v.errors = append(v.errors, "main: cannot come after emit:")
				}
				v.mode = MainMode
				v.hasMain = true
			case "emit":
				v.mode = EmitMode
			case "":
				if v.mode == InitialMode {
					v.mode = MainMode
				}
			default:
				v.errors = append(v.errors, fmt.Sprintf("Unknown label: %s", label))
			}
			v.VisitStatement(statement)
		}

		// Validate workflow structure
		if v.hasTake && !v.hasMain {
			v.errors = append(v.errors, "When take: is used, main: must also be present")
		}
	}
	v.VisitExpressionStatementHook = func(statement *parser.ExpressionStatement) {
		expr := statement.GetExpression()
		if variable, ok := expr.(*parser.VariableExpression); ok {
			if v.mode == TakeMode {
				v.Takes = append(v.Takes, variable.GetText())
				return
			} else if v.mode == EmitMode {
				v.Emits = append(v.Emits, variable.GetText())
				return
			}
		}
		// TODO: check property .out
		// TODO: check binary expr
		if propExpr, ok := expr.(*parser.PropertyExpression); ok {
			prop := propExpr.GetText()
			if v.mode == EmitMode {
				v.Emits = append(v.Emits, prop)
				return
			}
		}
		if binaryExpr, ok := expr.(*parser.BinaryExpression); ok {
			if binaryExpr.GetOperation().GetText() == "=" {
				if leftVar, ok := binaryExpr.GetLeftExpression().(*parser.VariableExpression); ok {
					if v.mode == EmitMode {
						v.Emits = append(v.Emits, leftVar.GetText())
					}
				}
				return
			}
		}
		v.VisitExpression(statement.GetExpression())
	}
	return v
}

type Workflow struct {
	Name    string
	Takes   []string
	Emits   []string
	Closure *parser.ClosureExpression
}

func (w *Workflow) ToProto() *pb.Workflow {
	return &pb.Workflow{
		Name:  w.Name,
		Takes: w.Takes,
		Emits: w.Emits,
	}
}

type WorkflowVisitor struct {
	*BaseVisitor
	workflows []Workflow
}

func NewWorkflowVisitor() *WorkflowVisitor {
	v := &WorkflowVisitor{
		BaseVisitor: NewBaseVisitor(),
		workflows:   []Workflow{},
	}
	v.VisitMethodCallExpressionHook = func(call *parser.MethodCallExpression) {
		v.VisitExpression(call.GetObjectExpression())
		v.VisitExpression(call.GetMethod())
		v.VisitExpression(call.GetArguments())
		mce, ok := call.GetMethod().(*parser.ConstantExpression)
		if !ok {
			return
		}
		if mce.GetText() != "workflow" {
			return
		}
		args, ok := call.GetArguments().(*parser.ArgumentListExpression)
		if !ok {
			return
		}
		if len(args.GetExpressions()) != 1 {
			return
		}
		arg := args.GetExpression(0)
		switch argExpr := arg.(type) {
		case *parser.ClosureExpression:
			v.workflows = append(v.workflows, makeWorkflow("", argExpr))
		case *parser.MethodCallExpression:
			if methodName, ok := argExpr.GetMethod().(*parser.ConstantExpression); ok {
				if args, ok := argExpr.GetArguments().(*parser.ArgumentListExpression); ok && len(args.GetExpressions()) > 0 {
					if closure, ok := args.GetExpression(0).(*parser.ClosureExpression); ok {
						v.workflows = append(v.workflows, makeWorkflow(methodName.GetText(), closure))
					}
				}
			}
		}
	}
	return v

}

func makeWorkflow(name string, closure *parser.ClosureExpression) Workflow {
	visitor := NewWorkflowBodyVisitor()
	visitor.VisitClosureExpression(closure)
	return Workflow{
		Name:    name,
		Takes:   visitor.Takes,
		Emits:   visitor.Emits,
		Closure: closure,
	}
}
