package nf

import (
	"reft-go/parser"
)

var _ parser.GroovyCodeVisitor = (*BaseVisitor)(nil)

type BaseVisitor struct {
}

// NewBaseVisitor creates a new BaseVisitor
func NewBaseVisitor() *BaseVisitor {
	return &BaseVisitor{}
}

// Statements
func (v *BaseVisitor) VisitBlockStatement(block *parser.BlockStatement) {
	for _, statement := range block.GetStatements() {
		v.VisitStatement(statement)
	}
}

func (v *BaseVisitor) VisitForLoop(statement *parser.ForStatement) {
	v.VisitExpression(statement.GetCollectionExpression())
	v.VisitStatement(statement.GetLoopBlock())
}

func (v *BaseVisitor) VisitWhileLoop(statement *parser.WhileStatement) {
	v.VisitExpression(statement.GetBooleanExpression())
	v.VisitStatement(statement.GetLoopBlock())
}

func (v *BaseVisitor) VisitDoWhileLoop(statement *parser.DoWhileStatement) {
	v.VisitStatement(statement.GetLoopBlock())
	v.VisitExpression(statement.GetBooleanExpression())
}

func (v *BaseVisitor) VisitIfElse(statement *parser.IfStatement) {
	v.VisitExpression(statement.GetBooleanExpression())
	v.VisitStatement(statement.GetIfBlock())
	v.VisitStatement(statement.GetElseBlock())
}

func (v *BaseVisitor) VisitExpressionStatement(statement *parser.ExpressionStatement) {
	v.VisitExpression(statement.GetExpression())
}

func (v *BaseVisitor) VisitReturnStatement(statement *parser.ReturnStatement) {
	v.VisitExpression(statement.GetExpression())
}

func (v *BaseVisitor) VisitAssertStatement(statement *parser.AssertStatement) {
	v.VisitExpression(statement.GetBooleanExpression())
	v.VisitExpression(statement.GetMessageExpression())
}

func (v *BaseVisitor) VisitTryCatchFinally(statement *parser.TryCatchStatement) {
	for _, resource := range statement.GetResourceStatements() {
		v.VisitStatement(resource)
	}
	v.VisitStatement(statement.GetTryStatement())
	for _, catchStatement := range statement.GetCatchStatements() {
		v.VisitStatement(catchStatement)
	}
	v.VisitStatement(statement.GetFinallyStatement())
}

func (v *BaseVisitor) VisitSwitch(statement *parser.SwitchStatement) {
	v.VisitExpression(statement.GetExpression())
	for _, caseStatement := range statement.GetCaseStatements() {
		v.VisitStatement(caseStatement)
	}
	v.VisitStatement(statement.GetDefaultStatement())
}

func (v *BaseVisitor) VisitCaseStatement(statement *parser.CaseStatement) {
	v.VisitExpression(statement.GetExpression())
	v.VisitStatement(statement.GetCode())
}

func (v *BaseVisitor) VisitBreakStatement(statement *parser.BreakStatement) {}

func (v *BaseVisitor) VisitContinueStatement(statement *parser.ContinueStatement) {}

func (v *BaseVisitor) VisitThrowStatement(statement *parser.ThrowStatement) {
	v.VisitExpression(statement.GetExpression())
}

func (v *BaseVisitor) VisitSynchronizedStatement(statement *parser.SynchronizedStatement) {
	v.VisitExpression(statement.GetExpression())
	v.VisitStatement(statement.GetCode())
}

func (v *BaseVisitor) VisitCatchStatement(statement *parser.CatchStatement) {
	v.VisitStatement(statement.GetCode())
}

func (v *BaseVisitor) VisitEmptyStatement(statement *parser.EmptyStatement) {}

func (v *BaseVisitor) VisitStatement(statement parser.Statement) {
	statement.Visit(v)
}

// Expressions
func (v *BaseVisitor) VisitMethodCallExpression(call *parser.MethodCallExpression) {
	v.VisitExpression(call.GetObjectExpression())
	v.VisitExpression(call.GetMethod())
	v.VisitExpression(call.GetArguments())
}

func (v *BaseVisitor) VisitStaticMethodCallExpression(call *parser.StaticMethodCallExpression) {
	v.VisitExpression(call.GetArguments())
}

func (v *BaseVisitor) VisitConstructorCallExpression(call *parser.ConstructorCallExpression) {
	v.VisitExpression(call.GetArguments())
}

func (v *BaseVisitor) VisitTernaryExpression(expression *parser.TernaryExpression) {
	booleanExpr := expression.GetBooleanExpression()
	v.VisitExpression(booleanExpr)
	v.VisitExpression(expression.GetTrueExpression())
	v.VisitExpression(expression.GetFalseExpression())
}

func (v *BaseVisitor) VisitShortTernaryExpression(expression *parser.ElvisOperatorExpression) {
	v.VisitTernaryExpression(expression.TernaryExpression)
}

func (v *BaseVisitor) VisitBinaryExpression(expression *parser.BinaryExpression) {
	v.VisitExpression(expression.GetLeftExpression())
	v.VisitExpression(expression.GetRightExpression())
}

func (v *BaseVisitor) VisitPrefixExpression(expression *parser.PrefixExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *BaseVisitor) VisitPostfixExpression(expression *parser.PostfixExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *BaseVisitor) VisitBooleanExpression(expression *parser.BooleanExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *BaseVisitor) VisitClosureExpression(expression *parser.ClosureExpression) {
	if expression.IsParameterSpecified() {
		for _, parameter := range expression.GetParameters() {
			if parameter.HasInitialExpression() {
				v.VisitExpression(parameter.GetInitialExpression())
			}
		}
	}
	v.VisitStatement(expression.GetCode())
}

func (v *BaseVisitor) VisitLambdaExpression(expression *parser.LambdaExpression) {
	v.VisitClosureExpression(expression.ClosureExpression)
}

func (v *BaseVisitor) VisitTupleExpression(expression parser.ITupleExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
}

func (v *BaseVisitor) VisitMapExpression(expression *parser.MapExpression) {
	entries := expression.GetMapEntryExpressions()
	exprs := make([]parser.Expression, len(entries))
	for i, entry := range entries {
		exprs[i] = entry
	}
	v.VisitListOfExpressions(exprs)
}

func (v *BaseVisitor) VisitMapEntryExpression(expression *parser.MapEntryExpression) {
	v.VisitExpression(expression.GetKeyExpression())
	v.VisitExpression(expression.GetValueExpression())
}

func (v *BaseVisitor) VisitListExpression(expression *parser.ListExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
}

func (v *BaseVisitor) VisitRangeExpression(expression *parser.RangeExpression) {
	v.VisitExpression(expression.GetFrom())
	v.VisitExpression(expression.GetTo())
}

func (v *BaseVisitor) VisitPropertyExpression(expression *parser.PropertyExpression) {
	v.VisitExpression(expression.GetObjectExpression())
	v.VisitExpression(expression.GetProperty())
}

func (v *BaseVisitor) VisitAttributeExpression(expression *parser.AttributeExpression) {
	v.VisitExpression(expression.GetObjectExpression())
	v.VisitExpression(expression.GetProperty())
}

func (v *BaseVisitor) VisitFieldExpression(expression *parser.FieldExpression) {}

func (v *BaseVisitor) VisitMethodPointerExpression(expression *parser.MethodPointerExpression) {
	v.VisitExpression(expression.GetExpression())
	v.VisitExpression(expression.GetMethodName())
}

func (v *BaseVisitor) VisitMethodReferenceExpression(expression *parser.MethodReferenceExpression) {
	v.VisitMethodPointerExpression(expression.MethodPointerExpression)
}

func (v *BaseVisitor) VisitConstantExpression(expression *parser.ConstantExpression) {}

func (v *BaseVisitor) VisitClassExpression(expression *parser.ClassExpression) {}

func (v *BaseVisitor) VisitVariableExpression(expression *parser.VariableExpression) {}

func (v *BaseVisitor) VisitDeclarationExpression(expression *parser.DeclarationExpression) {
	v.VisitBinaryExpression(expression.BinaryExpression)
}

func (v *BaseVisitor) VisitGStringExpression(expression *parser.GStringExpression) {
	v.VisitListOfExpressions(convertToExpressionSlice(expression.GetStrings()))
	v.VisitListOfExpressions(expression.GetValues())
}

func (v *BaseVisitor) VisitArrayExpression(expression *parser.ArrayExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
	v.VisitListOfExpressions(expression.GetSizeExpression())
}

func (v *BaseVisitor) VisitSpreadExpression(expression *parser.SpreadExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *BaseVisitor) VisitSpreadMapExpression(expression *parser.SpreadMapExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *BaseVisitor) VisitNotExpression(expression *parser.NotExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *BaseVisitor) VisitUnaryMinusExpression(expression *parser.UnaryMinusExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *BaseVisitor) VisitUnaryPlusExpression(expression *parser.UnaryPlusExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *BaseVisitor) VisitBitwiseNegationExpression(expression *parser.BitwiseNegationExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *BaseVisitor) VisitCastExpression(expression *parser.CastExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *BaseVisitor) VisitArgumentlistExpression(expression *parser.ArgumentListExpression) {
	v.VisitTupleExpression(expression)
}

func (v *BaseVisitor) VisitClosureListExpression(expression *parser.ClosureListExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
}

func (v *BaseVisitor) VisitEmptyExpression(expression *parser.EmptyExpression) {}

func (v *BaseVisitor) VisitListOfExpressions(expressions []parser.Expression) {
	for _, expr := range expressions {
		v.VisitExpression(expr)
	}
}

func (v *BaseVisitor) VisitExpression(expression parser.Expression) {
	expression.Visit(v)
}
