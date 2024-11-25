package configlint

import (
	"reft-go/nf"
	"reft-go/parser"
	"slices"
)

func ParseConfig(block *parser.BlockStatement) []ProcessScope {
	// Create and use ProcessVisitor
	processVisitor := NewProcessScopeVisitor()
	processVisitor.VisitBlockStatement(block)

	// Get the process scopes
	processScopes := processVisitor.processScopes

	var scopes []ProcessScope
	for _, processScope := range processScopes {
		scopes = append(scopes, makeProcessScope(processScope.First, processScope.Second))
	}
	return scopes
}

func makeProcessScope(lineNumber int, closure *parser.ClosureExpression) ProcessScope {
	directives := getDirectives(closure)
	namedScopes := getNamedScopes(closure)
	return ProcessScope{
		LineNumber:  lineNumber,
		Directives:  directives,
		NamedScopes: namedScopes,
	}
}

func getDirectives(closure *parser.ClosureExpression) []Directive {
	directiveVisitor := NewDirectiveVisitor()
	directiveVisitor.VisitClosureExpression(closure)
	var directives []Directive
	for name, pair := range directiveVisitor.directives {
		directive := Directive{
			LineNumber: pair.First,
			Name:       name,
		}
		directiveBodyVisitor := NewDirectiveBodyVisitor()
		directiveBodyVisitor.VisitExpression(pair.Second)
		directive.Options = directiveBodyVisitor.namedOptions
		directive.Value = directiveBodyVisitor.value
		directives = append(directives, directive)
	}
	return directives
}

func getNamedScopes(closure *parser.ClosureExpression) []NamedScope {
	namedScopeVisitor := NewNamedScopeVisitor()
	namedScopeVisitor.VisitClosureExpression(closure)
	var namedScopes []NamedScope
	for _, triple := range namedScopeVisitor.namedScopes {
		namedScopes = append(namedScopes, NamedScope{
			LineNumber: triple.First,
			Name:       triple.Second,
			Directives: getDirectives(triple.Third),
		})
	}
	return namedScopes
}

type DirectiveValue struct {
	Params     []string
	InClosure  bool
	Expression parser.Expression
}

type NamedOption struct {
	LineNumber int
	Name       string
	Value      DirectiveValue
}

type Directive struct {
	LineNumber int
	Name       string
	Options    []NamedOption
	Value      DirectiveValue
}

type NamedScope struct {
	LineNumber int
	Name       string
	Directives []Directive
}

type ProcessScope struct {
	LineNumber  int
	Directives  []Directive
	NamedScopes []NamedScope
}

type Pair[F any, S any] struct {
	First  F
	Second S
}

type Triple[F any, S any, T any] struct {
	First  F
	Second S
	Third  T
}

type ProcessScopeVisitor struct {
	*nf.BaseVisitor
	processScopes []Pair[int, *parser.ClosureExpression]
}

func NewProcessScopeVisitor() *ProcessScopeVisitor {
	v := &ProcessScopeVisitor{BaseVisitor: nf.NewBaseVisitor()}
	v.VisitMethodCallExpressionHook = func(call *parser.MethodCallExpression) {
		if call.GetMethod().GetText() == "process" {
			_ = call.GetArguments()
			if args, ok := call.GetArguments().(*parser.ArgumentListExpression); ok {
				exprs := args.GetExpressions()
				if len(exprs) == 1 {
					expr := exprs[0]
					if closure, ok := expr.(*parser.ClosureExpression); ok {
						v.processScopes = append(v.processScopes, Pair[int, *parser.ClosureExpression]{
							First:  expr.GetLineNumber(),
							Second: closure,
						})
					}
				}
			}
		}
		v.VisitExpression(call.GetObjectExpression())
		v.VisitExpression(call.GetMethod())
		v.VisitExpression(call.GetArguments())
	}
	return v
}

type DirectiveVisitor struct {
	*nf.BaseVisitor
	// directive assignments should be unique
	directives map[string]Pair[int, parser.Expression]
}

func NewDirectiveVisitor() *DirectiveVisitor {
	v := &DirectiveVisitor{BaseVisitor: nf.NewBaseVisitor(), directives: make(map[string]Pair[int, parser.Expression])}
	v.VisitBinaryExpressionHook = func(expr *parser.BinaryExpression) {
		// find assignments
		if expr.GetOperation().GetText() == "=" {
			name := expr.GetLeftExpression().GetText()
			if _, exists := nf.DirectiveSet[name]; exists {
				// TODO: raise an error if the directive is already defined
				v.directives[name] = Pair[int, parser.Expression]{expr.GetLineNumber(), expr.GetRightExpression()}
			} else if len(name) > 4 && name[:4] == "ext." {
				// TODO: raise an error if the directive is already defined
				v.directives[name] = Pair[int, parser.Expression]{expr.GetLineNumber(), expr.GetRightExpression()}
			}
		}
	}
	// only visit top level assignments
	v.VisitClosureExpressionHook = func(expr *parser.ClosureExpression) {
		if bs, ok := expr.GetCode().(*parser.BlockStatement); ok {
			for _, stmt := range bs.GetStatements() {
				if exprStmt, ok := stmt.(*parser.ExpressionStatement); ok {
					if binaryExpr, ok := exprStmt.GetExpression().(*parser.BinaryExpression); ok {
						v.VisitBinaryExpression(binaryExpr)
					}
				}
			}
		}
	}
	return v
}

// guaranteed to not have both namedOptions and value set
type DirectiveBodyVisitor struct {
	*nf.BaseVisitor
	namedOptions []NamedOption
	value        DirectiveValue
}

func NewDirectiveBodyVisitor() *DirectiveBodyVisitor {
	v := &DirectiveBodyVisitor{BaseVisitor: nf.NewBaseVisitor()}
	v.VisitExpressionHook = func(expr parser.Expression) {
		if mapExpr, ok := expr.(*parser.MapExpression); ok {
			exprs := mapExpr.GetMapEntryExpressions()
			for _, entry := range exprs {
				name := entry.GetKeyExpression().GetText()
				value := entry.GetValueExpression()
				directiveValue := getDirectiveValue(value)
				v.namedOptions = append(v.namedOptions, NamedOption{
					LineNumber: entry.GetLineNumber(),
					Name:       name,
					Value:      directiveValue,
				})
			}
		} else {
			v.value = getDirectiveValue(expr)
		}
	}
	return v
}

func getDirectiveValue(expr parser.Expression) DirectiveValue {
	visitor := NewDirectiveValueVisitor()
	visitor.VisitExpression(expr)
	inClosure := true
	var params []string
	for _, pair := range visitor.params {
		if !pair.First {
			// if any of the params are not in a closure,
			// then the entire directive is considered not in a closure
			inClosure = false
		}
		params = append(params, pair.Second)
	}
	// if there are no params, we'll reset to the false zero-value
	if len(params) == 0 {
		inClosure = false
	}
	// De-duplicate params
	slices.Sort(params)
	params = slices.Compact(params)
	directiveValue := DirectiveValue{
		Params:     params,
		InClosure:  inClosure,
		Expression: expr,
	}
	return directiveValue
}

type DirectiveValueVisitor struct {
	*nf.BaseVisitor
	params    []Pair[bool, string]
	inClosure bool
}

func NewDirectiveValueVisitor() *DirectiveValueVisitor {
	v := &DirectiveValueVisitor{BaseVisitor: nf.NewBaseVisitor()}
	v.VisitClosureExpressionHook = func(expr *parser.ClosureExpression) {
		v.inClosure = true
		v.VisitStatement(expr.GetCode())
		v.inClosure = false
	}
	v.VisitPropertyExpressionHook = func(expr *parser.PropertyExpression) {
		if expr.GetObjectExpression().GetText() == "params" {
			v.params = append(v.params, Pair[bool, string]{v.inClosure, expr.GetProperty().GetText()})
		}
	}
	return v
}

type NamedScopeVisitor struct {
	*nf.BaseVisitor
	namedScopes []Triple[int, string, *parser.ClosureExpression]
}

func NewNamedScopeVisitor() *NamedScopeVisitor {
	v := &NamedScopeVisitor{BaseVisitor: nf.NewBaseVisitor()}
	v.VisitExpressionStatementHook = func(expr *parser.ExpressionStatement) {
		label := expr.GetStatementLabel()
		if label == "withName" {
			if mce, ok := expr.GetExpression().(*parser.MethodCallExpression); ok {
				name := mce.GetMethod().GetText()
				if argList, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
					args := argList.GetExpressions()
					if len(args) == 1 {
						arg := args[0]
						if closure, ok := arg.(*parser.ClosureExpression); ok {
							v.namedScopes = append(v.namedScopes, Triple[int, string, *parser.ClosureExpression]{
								First:  expr.GetLineNumber(),
								Second: name,
								Third:  closure,
							})
						}
					}
				}
			}
		}
	}
	return v
}
