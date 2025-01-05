package parser

import "github.com/antlr4-go/antlr/v4"

func (v *ASTBuilder) VisitPathElement(ctx *PathElementContext) interface{} {
	baseExpr := ctx.GetNodeMetaData(PATH_EXPRESSION_BASE_EXPR).(Expression)
	if baseExpr == nil {
		panic("baseExpr is required!")
	}

	if ctx.NamePart() != nil {
		namePartExpr := v.VisitNamePart(ctx.NamePart().(*NamePartContext)).(Expression)
		var nonWildcardTypeArgumentsContext *NonWildcardTypeArgumentsContext
		if ctx.NonWildcardTypeArguments() != nil {
			nonWildcardTypeArgumentsContext = ctx.NonWildcardTypeArguments().(*NonWildcardTypeArgumentsContext)
		}
		genericsTypes := (v.VisitNonWildcardTypeArguments(nonWildcardTypeArgumentsContext)).([]*GenericsType)

		if ctx.DOT() != nil {
			isSafeChain := isTrue(baseExpr, PATH_EXPRESSION_BASE_EXPR_SAFE_CHAIN)
			return v.createDotExpression(ctx, baseExpr, namePartExpr, genericsTypes, isSafeChain)
		} else if ctx.SAFE_DOT() != nil {
			return v.createDotExpression(ctx, baseExpr, namePartExpr, genericsTypes, true)
		} else if ctx.SAFE_CHAIN_DOT() != nil {
			expression := v.createDotExpression(ctx, baseExpr, namePartExpr, genericsTypes, true)
			expression.SetNodeMetaData(PATH_EXPRESSION_BASE_EXPR_SAFE_CHAIN, true)
			return expression
		} else if ctx.METHOD_POINTER() != nil {
			return configureAST(NewMethodPointerExpression(baseExpr, namePartExpr), ctx)
		} else if ctx.METHOD_REFERENCE() != nil {
			return configureAST(NewMethodReferenceExpression(baseExpr, namePartExpr), ctx)
		} else if ctx.SPREAD_DOT() != nil {
			if ctx.AT() != nil {
				attributeExpression := NewAttributeExpressionWithSafe(baseExpr, namePartExpr, true)
				attributeExpression.SetSpreadSafe(true)
				return configureAST(attributeExpression, ctx)
			} else {
				propertyExpression := NewPropertyExpressionWithSafe(baseExpr, namePartExpr, true)
				propertyExpression.SetNodeMetaData(PATH_EXPRESSION_BASE_EXPR_GENERICS_TYPES, genericsTypes)
				propertyExpression.SetSpreadSafe(true)
				return configureAST(propertyExpression, ctx)
			}
		}
	} else if ctx.Creator() != nil {
		creatorContext := ctx.Creator().(*CreatorContext)
		creatorContext.PutNodeMetaData(ENCLOSING_INSTANCE_EXPRESSION, baseExpr)
		return configureAST(v.VisitCreator(creatorContext).(Expression), ctx)
	} else if ctx.IndexPropertyArgs() != nil {
		indexPropertyArgs := ctx.IndexPropertyArgs().(*IndexPropertyArgsContext)
		if indexPropertyArgs.ExpressionList() == nil {
			// Handle empty index case
			return configureAST(
				NewBinaryExpressionWithSafe(
					baseExpr,
					v.createGroovyToken(indexPropertyArgs.LBRACK().GetSymbol()),
					EMPTY_EXPRESSION,
					false,
				),
				ctx,
			)
		}
		tuple := v.VisitIndexPropertyArgs(indexPropertyArgs).(Tuple2[antlr.Token, Expression])
		isSafeChain := isTrue(baseExpr, PATH_EXPRESSION_BASE_EXPR_SAFE_CHAIN)
		return configureAST(
			NewBinaryExpressionWithSafe(
				baseExpr,
				v.createGroovyToken(tuple.V1),
				tuple.V2.(Expression),
				isSafeChain || ctx.IndexPropertyArgs().(*IndexPropertyArgsContext).SAFE_INDEX() != nil,
			),
			ctx,
		)
	} else if ctx.NamedPropertyArgs() != nil {
		mapEntryExpressionList := v.VisitNamedPropertyArgs(ctx.NamedPropertyArgs().(*NamedPropertyArgsContext)).([]MapEntryExpression)

		expressions := make([]Expression, len(mapEntryExpressionList))
		for i, v := range mapEntryExpressionList {
			expressions[i] = &v
		}
		listExpression := configureAST(
			NewListExpressionWithExpressions(expressions),
			ctx.NamedPropertyArgs(),
		)

		namedPropertyArgsContext := ctx.NamedPropertyArgs().(*NamedPropertyArgsContext)
		var token antlr.Token
		if namedPropertyArgsContext.LBRACK() == nil {
			token = namedPropertyArgsContext.SAFE_INDEX().GetSymbol()
		} else {
			token = namedPropertyArgsContext.LBRACK().GetSymbol()
		}
		return configureAST(
			NewBinaryExpression(baseExpr, v.createGroovyToken(token), listExpression),
			ctx,
		)
	} else if ctx.Arguments() != nil {
		argumentsExpr := v.VisitArguments(ctx.Arguments().(*ArgumentsContext)).(Expression)
		configureAST(argumentsExpr, ctx)

		if v.isInsideParentheses(baseExpr) {
			return configureAST(v.createCallMethodCallExpression(baseExpr, argumentsExpr), ctx)
		}

		if attributeExpression, ok := baseExpr.(*AttributeExpression); ok {
			attributeExpression.SetSpreadSafe(false)
			return configureAST(v.createCallMethodCallExpressionWithImplicitThis(attributeExpression, argumentsExpr, true), ctx)
		}

		if propertyExpression, ok := baseExpr.(*PropertyExpression); ok {
			methodCallExpression := v.createMethodCallExpression(propertyExpression, argumentsExpr)
			return configureAST(methodCallExpression, ctx)
		}

		if variableExpression, ok := baseExpr.(*VariableExpression); ok {
			baseExprText := variableExpression.GetText()
			if baseExprText == VOID_STR {
				return configureAST(v.createCallMethodCallExpression(v.createConstantExpression(baseExpr), argumentsExpr), ctx)
			} else if isPrimitiveType(baseExprText) {
				panic(createParsingFailedException("Primitive type literal: "+baseExprText+" cannot be used as a method name", parserRuleContextAdapter{ctx}))
			}
		}

		_, isVar := baseExpr.(*VariableExpression)
		_, isGString := baseExpr.(*GStringExpression)
		if isVar || isGString || (IsInstanceOf(baseExpr, (*ConstantExpression)(nil)) && isTrue(baseExpr, IS_STRING)) {
			baseExprText := baseExpr.GetText()
			if baseExprText == THIS_STR || baseExprText == SUPER_STR {
				if v.visitingClosureCount > 0 {
					return configureAST(
						NewMethodCallExpression(
							baseExpr,
							baseExprText,
							argumentsExpr,
						),
						ctx,
					)
				}

				var classNode IClassNode
				if baseExprText == SUPER_STR {
					classNode = SUPER
				} else {
					classNode = THIS
				}
				return configureAST(
					NewConstructorCallExpression(classNode, argumentsExpr),
					ctx,
				)
			}

			methodCallExpression := v.createMethodCallExpression(baseExpr, argumentsExpr)
			return configureAST(methodCallExpression, ctx)
		}

		return configureAST(v.createCallMethodCallExpression(baseExpr, argumentsExpr), ctx)
	} else if ctx.ClosureOrLambdaExpression() != nil {
		closureExpression := v.VisitClosureOrLambdaExpression(ctx.ClosureOrLambdaExpression().(*ClosureOrLambdaExpressionContext)).(*ClosureExpression)

		if methodCallExpression, ok := baseExpr.(*MethodCallExpression); ok {
			argumentsExpression := methodCallExpression.GetArguments()

			if argumentListExpression, ok := argumentsExpression.(*ArgumentListExpression); ok {
				argumentListExpression.AddExpression(closureExpression)
				return configureAST(methodCallExpression, ctx)
			}

			if tupleExpression, ok := argumentsExpression.(*TupleExpression); ok {
				namedArgumentListExpression := tupleExpression.GetExpression(0).(*NamedArgumentListExpression)

				if len(tupleExpression.GetExpressions()) > 0 {
					methodCallExpression.SetArguments(
						configureASTFromSource(
							NewArgumentListExpressionFromSlice(
								configureASTFromSource(
									NewMapExpressionWithEntries(namedArgumentListExpression.GetMapEntryExpressions()),
									namedArgumentListExpression,
								),
								closureExpression,
							),
							tupleExpression,
						),
					)
				} else {
					methodCallExpression.SetArguments(
						configureASTFromSource(
							NewArgumentListExpressionFromSlice(closureExpression),
							tupleExpression,
						),
					)
				}

				return configureAST(methodCallExpression, ctx)
			}
		}

		if propertyExpression, ok := baseExpr.(*PropertyExpression); ok {
			methodCallExpression := v.createMethodCallExpression(
				propertyExpression,
				configureASTFromSource(
					NewArgumentListExpressionFromSlice(closureExpression),
					closureExpression,
				),
			)

			return configureAST(methodCallExpression, ctx)
		}

		if ve, ok := baseExpr.(*VariableExpression); ok {
			// Handle VariableExpression
			methodCallExpression := v.createMethodCallExpression(
				ve,
				configureASTFromSource(NewArgumentListExpressionFromSlice(closureExpression), closureExpression),
			)
			return configureAST(methodCallExpression, ctx)
		} else if gse, ok := baseExpr.(*GStringExpression); ok {
			// Handle GStringExpression
			methodCallExpression := v.createMethodCallExpression(
				gse,
				configureASTFromSource(NewArgumentListExpressionFromSlice(closureExpression), closureExpression),
			)
			return configureAST(methodCallExpression, ctx)
		} else if ce, ok := baseExpr.(*ConstantExpression); ok && isTrue(ce, IS_STRING) {
			// Handle ConstantExpression that is a string
			methodCallExpression := v.createMethodCallExpression(
				ce,
				configureASTFromSource(NewArgumentListExpressionFromSlice(closureExpression), closureExpression),
			)
			return configureAST(methodCallExpression, ctx)
		}

		methodCallExpression := v.createMethodCallExpression(
			baseExpr,
			configureASTFromSource(
				NewArgumentListExpressionFromSlice(closureExpression),
				closureExpression,
			),
		)

		return configureAST(methodCallExpression, ctx)
	}

	panic(createParsingFailedException("Unsupported path element: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}
