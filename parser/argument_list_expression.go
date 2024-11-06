package parser

import "strings"

var _ Expression = (*ArgumentListExpression)(nil)
var _ ITupleExpression = (*ArgumentListExpression)(nil)

// ArgumentListExpression represents one or more arguments being passed into a method
type ArgumentListExpression struct {
	*BaseExpression
	expressions []Expression
}

var (
	EmptyArray     = []interface{}{}
	EmptyArguments = NewArgumentListExpression()
)

func NewArgumentListExpression() *ArgumentListExpression {
	return &ArgumentListExpression{
		BaseExpression: NewBaseExpression(),
		expressions:    []Expression{},
	}
}

func NewArgumentListExpressionFromList(expressions []Expression) *ArgumentListExpression {
	return &ArgumentListExpression{
		BaseExpression: NewBaseExpression(),
		expressions:    expressions,
	}
}

func NewArgumentListExpressionFromSlice(expressions ...Expression) *ArgumentListExpression {
	return &ArgumentListExpression{
		BaseExpression: NewBaseExpression(),
		expressions:    expressions,
	}
}

func NewArgumentListExpressionFromParameters(parameters []*Parameter) *ArgumentListExpression {
	ale := NewArgumentListExpression()
	for _, param := range parameters {
		ale.AddExpression(NewVariableExpressionWithVariable(param))
	}
	return ale
}

func (a *ArgumentListExpression) AddExpression(expr Expression) {
	a.expressions = append(a.expressions, expr)
}

/*
func (a *ArgumentListExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewArgumentListExpressionFromList(TransformExpressions(a.expressions, transformer))
	ret.SetSourcePosition(a)
	ret.CopyNodeMetaData(a)
	return ret
}
*/

func (a *ArgumentListExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitArgumentlistExpression(a)
}

func (a *ArgumentListExpression) GetExpressions() []Expression {
	return a.expressions
}

func (a *ArgumentListExpression) GetExpression(i int) Expression {
	return a.expressions[i]
}

func (a *ArgumentListExpression) PrependExpression(expression Expression) ITupleExpression {
	a.expressions = append([]Expression{expression}, a.expressions...)
	return a
}

func (a *ArgumentListExpression) GetText() string {
	var buffer strings.Builder
	buffer.WriteString("(")
	for i, expression := range a.GetExpressions() {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(expression.GetText())
	}
	buffer.WriteString(")")
	return buffer.String()
}
