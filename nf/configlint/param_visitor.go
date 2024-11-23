package configlint

import (
	"reft-go/nf"
	"reft-go/parser"
	"sort"
)

var _ parser.GroovyCodeVisitor = (*ParamVisitor)(nil)

type ParamInfo struct {
	Name          string
	LineNumber    int
	InDirective   bool
	DirectiveName string
	InClosure     bool
}

type ParamVisitor struct {
	*nf.BaseVisitor
	params           map[string]ParamInfo
	methodCallStack  []string
	currentDirective string
	inMapExpression  bool
	closureDepth     int
}

// NewParamVisitor creates a new ParamVisitor
func NewParamVisitor() *ParamVisitor {
	return &ParamVisitor{
		BaseVisitor:     nf.NewBaseVisitor(),
		params:          make(map[string]ParamInfo),
		methodCallStack: make([]string, 0),
		closureDepth:    0,
	}
}

func (v *ParamVisitor) GetSortedParams() []ParamInfo {
	sortedParams := make([]ParamInfo, 0, len(v.params))
	for _, info := range v.params {
		sortedParams = append(sortedParams, info)
	}
	sort.Slice(sortedParams, func(i, j int) bool {
		return sortedParams[i].LineNumber < sortedParams[j].LineNumber
	})
	return sortedParams
}

// Override only the methods that need custom behavior
func (v *ParamVisitor) VisitMethodCallExpression(call *parser.MethodCallExpression) {
	methodName := call.GetMethod().GetText()
	v.methodCallStack = append(v.methodCallStack, methodName)

	_, isDirective := nf.DirectiveSet[methodName]
	if isDirective {
		v.currentDirective = methodName
	}

	v.BaseVisitor.VisitMethodCallExpression(call)

	if isDirective {
		v.currentDirective = ""
	}
	v.methodCallStack = v.methodCallStack[:len(v.methodCallStack)-1]
}

func (v *ParamVisitor) VisitMapExpression(expression *parser.MapExpression) {
	prevInMap := v.inMapExpression
	v.inMapExpression = true
	v.BaseVisitor.VisitMapExpression(expression)
	v.inMapExpression = prevInMap
}

func (v *ParamVisitor) VisitClosureExpression(expression *parser.ClosureExpression) {
	v.closureDepth++
	v.BaseVisitor.VisitClosureExpression(expression)
	v.closureDepth--
}

func (v *ParamVisitor) VisitPropertyExpression(expression *parser.PropertyExpression) {
	varExpr, isVarExpr := expression.GetObjectExpression().(*parser.VariableExpression)
	constExpr, isConstExpr := expression.GetProperty().(*parser.ConstantExpression)

	if isVarExpr && isConstExpr && varExpr.GetName() == "params" {
		paramName := constExpr.GetText()
		lineNumber := constExpr.GetLineNumber()

		inDirective := v.currentDirective != "" && v.inMapExpression
		v.params[paramName] = ParamInfo{
			Name:          paramName,
			LineNumber:    lineNumber,
			InDirective:   inDirective,
			DirectiveName: v.currentDirective,
			InClosure:     v.closureDepth > 0,
		}
	}

	v.BaseVisitor.VisitPropertyExpression(expression)
}
