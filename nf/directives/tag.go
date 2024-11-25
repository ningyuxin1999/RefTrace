package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"

	pb "reft-go/nf/proto"
)

func (t *TagDirective) ToProto() *pb.Directive {
	return &pb.Directive{
		Line: int32(t.Line()),
		Directive: &pb.Directive_Tag{
			Tag: &pb.TagDirective{
				Tag: t.Tag,
			},
		},
	}
}

var _ Directive = (*TagDirective)(nil)
var _ starlark.Value = (*TagDirective)(nil)
var _ starlark.HasAttrs = (*TagDirective)(nil)

func (t *TagDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "tag":
		return starlark.String(t.Tag), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("tag directive has no attribute %q", name))
	}
}

func (t *TagDirective) AttrNames() []string {
	return []string{"tag"}
}

type TagDirective struct {
	Tag  string
	line int
}

func (t *TagDirective) Line() int {
	return t.line
}

func (t *TagDirective) String() string {
	return fmt.Sprintf("TagDirective(Tag: %q)", t.Tag)
}

func (t *TagDirective) Type() string {
	return "tag_directive"
}

func (t *TagDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (t *TagDirective) Truth() starlark.Bool {
	return starlark.Bool(t.Tag != "")
}

func (t *TagDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(t.Tag))
	return h.Sum32(), nil
}

func MakeTagDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					return &TagDirective{Tag: strValue, line: mce.GetLineNumber()}, nil
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				return &TagDirective{Tag: gstringExpr.GetText(), line: mce.GetLineNumber()}, nil
			}
		}
	}
	return nil, errors.New("invalid TagDirective directive")
}
