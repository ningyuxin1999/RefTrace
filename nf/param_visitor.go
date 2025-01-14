package nf

import (
	pb "reft-go/nf/proto"
	"reft-go/parser"
	"sort"
)

var _ parser.GroovyCodeVisitor = (*ParamVisitor)(nil)

type ParamInfo struct {
	Name       string
	LineNumber int
}

func (p *ParamInfo) ToProto() *pb.Param {
	return &pb.Param{
		Line: int32(p.LineNumber),
		Name: p.Name,
	}
}

type ParamVisitor struct {
	*BaseVisitor
	params map[string]int // Use a map to represent a set of strings with line numbers
}

// NewParamVisitor creates a new ParamVisitor
func NewParamVisitor() *ParamVisitor {
	v := &ParamVisitor{
		BaseVisitor: NewBaseVisitor(),
		params:      make(map[string]int),
	}
	v.VisitPropertyExpressionHook = func(expression *parser.PropertyExpression) {
		varExpr, isVarExpr := expression.GetObjectExpression().(*parser.VariableExpression)
		constExpr, isConstExpr := expression.GetProperty().(*parser.ConstantExpression)
		if isVarExpr && isConstExpr && varExpr.GetName() == "params" {
			paramName := constExpr.GetText()
			lineNumber := constExpr.GetLineNumber()
			// Only store the line number if it's the first occurrence of this param
			if _, exists := v.params[paramName]; !exists {
				v.params[paramName] = lineNumber
			}
		}
		v.VisitExpression(expression.GetObjectExpression())
		v.VisitExpression(expression.GetProperty())
	}
	return v
}

func (v *ParamVisitor) GetSortedParams() []ParamInfo {
	sortedParams := make([]ParamInfo, 0, len(v.params))
	for param, lineNumber := range v.params {
		sortedParams = append(sortedParams, ParamInfo{Name: param, LineNumber: lineNumber})
	}
	sort.Slice(sortedParams, func(i, j int) bool {
		return sortedParams[i].LineNumber < sortedParams[j].LineNumber
	})
	return sortedParams
}
