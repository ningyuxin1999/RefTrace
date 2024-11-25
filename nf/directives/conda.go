package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"

	pb "reft-go/nf/proto"
)

func (c *Conda) ToProto() *pb.Directive {
	return &pb.Directive{
		Line: int32(c.Line()),
		Directive: &pb.Directive_Conda{
			Conda: &pb.CondaDirective{
				PossibleValues: c.PossibleValues,
			},
		},
	}
}

var _ Directive = (*Conda)(nil)

func (c *Conda) String() string {
	return fmt.Sprintf("Conda(possible_values=%q)", c.PossibleValues)
}

func (c *Conda) Type() string {
	return "conda"
}

func (c *Conda) Freeze() {
	// No mutable fields, so no action needed
}

func (c *Conda) Truth() starlark.Bool {
	return starlark.Bool(len(c.PossibleValues) > 0)
}

func (c *Conda) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%v", c.PossibleValues)))
	return h.Sum32(), nil
}

var _ starlark.Value = (*Conda)(nil)
var _ starlark.HasAttrs = (*Conda)(nil)

func (c *Conda) Attr(name string) (starlark.Value, error) {
	switch name {
	case "possible_values":
		values := make([]starlark.Value, len(c.PossibleValues))
		for i, v := range c.PossibleValues {
			values[i] = starlark.String(v)
		}
		return starlark.NewList(values), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("conda directive has no attribute %q", name))
	}
}

func (c *Conda) AttrNames() []string {
	return []string{"possible_values"}
}

type Conda struct {
	PossibleValues []string // List of possible values (e.g., ["bioconda::bcftools=1.14", ""])
	line           int
}

func (c *Conda) Line() int {
	return c.line
}

func MakeConda(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			switch expr := exprs[0].(type) {
			case *parser.ConstantExpression:
				if expr.GetText() == "null" {
					return &Conda{PossibleValues: []string{""}, line: mce.GetLineNumber()}, nil
				}
				if value, ok := expr.GetValue().(string); ok {
					return &Conda{PossibleValues: []string{value}, line: mce.GetLineNumber()}, nil
				}
			case *parser.GStringExpression:
				return &Conda{PossibleValues: []string{expr.GetText()}, line: mce.GetLineNumber()}, nil
			case *parser.TernaryExpression:
				var values []string

				// Get true value
				if tv, ok := expr.GetTrueExpression().(*parser.ConstantExpression); ok {
					if tv.GetText() == "null" {
						values = append(values, "")
					} else {
						values = append(values, tv.GetText())
					}
				}
				// Get false value
				if fv, ok := expr.GetFalseExpression().(*parser.ConstantExpression); ok {
					if fv.GetText() == "null" {
						values = append(values, "")
					} else {
						values = append(values, fv.GetText())
					}
				}
				return &Conda{PossibleValues: values, line: mce.GetLineNumber()}, nil
			}
		}
	}
	return nil, errors.New("invalid conda directive")
}
