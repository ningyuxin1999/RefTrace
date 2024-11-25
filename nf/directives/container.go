package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"

	pb "reft-go/nf/proto"
)

func (c *Container) ToProto() *pb.Directive {
	format := pb.ContainerDirective_SIMPLE
	if c.Format == Ternary {
		format = pb.ContainerDirective_TERNARY
	}

	return &pb.Directive{
		Line: int32(c.Line()),
		Directive: &pb.Directive_Container{
			Container: &pb.ContainerDirective{
				Format:     format,
				SimpleName: c.SimpleName,
				Condition:  c.Condition,
				TrueName:   c.TrueName,
				FalseName:  c.FalseName,
			},
		},
	}
}

var _ Directive = (*Container)(nil)

func (c *Container) String() string {
	return fmt.Sprintf("Container(%q)", c.GetName())
}

func (c *Container) Type() string {
	return "container"
}

func (c *Container) Freeze() {
	// No mutable fields, so no action needed
}

func (c *Container) Truth() starlark.Bool {
	return starlark.Bool(c.GetName() != "")
}

func (c *Container) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(c.GetName()))
	return h.Sum32(), nil
}

var _ starlark.Value = (*Container)(nil)
var _ starlark.HasAttrs = (*Container)(nil)

// Starlark attribute access
func (c *Container) Attr(name string) (starlark.Value, error) {
	switch name {
	case "format":
		return starlark.String(c.Format), nil
	case "simple_name":
		if c.Format == Simple {
			return starlark.String(c.SimpleName), nil
		}
		return starlark.None, nil
	case "condition":
		if c.Format == Ternary {
			return starlark.String(c.Condition), nil
		}
		return starlark.None, nil
	case "true_name":
		if c.Format == Ternary {
			return starlark.String(c.TrueName), nil
		}
		return starlark.None, nil
	case "false_name":
		if c.Format == Ternary {
			return starlark.String(c.FalseName), nil
		}
		return starlark.None, nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("container has no attribute %q", name))
	}
}

func (c *Container) AttrNames() []string {
	switch c.Format {
	case Simple:
		return []string{"format", "simple_name"}
	case Ternary:
		return []string{"format", "condition", "true_name", "false_name"}
	default:
		return []string{"format"}
	}
}

// either "simple" or "ternary"
type Format string

const (
	Simple  Format = "simple"
	Ternary Format = "ternary"
)

type Container struct {
	Format     Format
	SimpleName string
	Condition  string
	TrueName   string
	FalseName  string
	line       int
}

func (c *Container) GetName() string {
	switch c.Format {
	case Simple:
		return c.SimpleName
	case Ternary:
		return fmt.Sprintf("%s ? %s : %s", c.Condition, c.TrueName, c.FalseName)
	}
	panic("invalid container format")
}

func (c *Container) GetNames() []string {
	switch c.Format {
	case Simple:
		return []string{c.SimpleName}
	case Ternary:
		return []string{c.TrueName, c.FalseName}
	}
	panic("invalid container format")
}

func (c *Container) Line() int {
	return c.line
}

func NewSimpleContainer(name string, line int) *Container {
	return &Container{
		Format:     Simple,
		SimpleName: name,
		line:       line,
	}
}

func NewTernaryContainer(condition string, trueName string, falseName string, line int) *Container {
	return &Container{
		Format:    Ternary,
		Condition: condition,
		TrueName:  trueName,
		FalseName: falseName,
		line:      line,
	}
}

func MakeContainer(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				if constantExpr.GetText() == "null" {
					return NewSimpleContainer("", mce.GetLineNumber()), nil
				}
				if value, ok := constantExpr.GetValue().(string); ok {
					return NewSimpleContainer(value, mce.GetLineNumber()), nil
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				if len(gstringExpr.GetValues()) == 1 {
					if ternaryExpr, ok := gstringExpr.GetValues()[0].(*parser.TernaryExpression); ok {
						return NewTernaryContainer(ternaryExpr.GetBooleanExpression().GetText(), ternaryExpr.GetTrueExpression().GetText(), ternaryExpr.GetFalseExpression().GetText(), mce.GetLineNumber()), nil
					}
				}
				return NewSimpleContainer(gstringExpr.GetText(), mce.GetLineNumber()), nil
			}
		}
	}
	return nil, errors.New("invalid container directive")
}
