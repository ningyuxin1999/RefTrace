package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"

	pb "reft-go/nf/proto"
)

func (a *AfterScript) ToProto() *pb.Directive {
	return &pb.Directive{
		Line: int32(a.Line()),
		Directive: &pb.Directive_AfterScript{
			AfterScript: &pb.AfterScriptDirective{
				Script: a.Script,
			},
		},
	}
}

var _ Directive = (*AfterScript)(nil)
var _ starlark.Value = (*AfterScript)(nil)
var _ starlark.HasAttrs = (*AfterScript)(nil)

func (a *AfterScript) Attr(name string) (starlark.Value, error) {
	switch name {
	case "script":
		return starlark.String(a.Script), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("afterscript directive has no attribute %q", name))
	}
}

func (a *AfterScript) AttrNames() []string {
	return []string{"script"}
}

func (a *AfterScript) String() string {
	return fmt.Sprintf("AfterScript(%q)", a.Script)
}

func (a *AfterScript) Type() string {
	return "afterscript"
}

func (a *AfterScript) Freeze() {
	// No mutable fields, so no action needed
}

func (a *AfterScript) Truth() starlark.Bool {
	return starlark.Bool(a.Script != "")
}

func (a *AfterScript) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(a.Script))
	return h.Sum32(), nil
}

type AfterScript struct {
	Script string
	line   int
}

func (a *AfterScript) Line() int {
	return a.line
}

func MakeAfterScript(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					return &AfterScript{Script: strValue, line: mce.GetLineNumber()}, nil
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				return &AfterScript{Script: gstringExpr.GetText(), line: mce.GetLineNumber()}, nil
			}
		}
	}
	return nil, errors.New("invalid afterScript directive")
}
