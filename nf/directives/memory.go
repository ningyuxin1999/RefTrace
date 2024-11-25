package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"

	pb "reft-go/nf/proto"
)

func (m *MemoryDirective) ToProto() *pb.Directive {
	return &pb.Directive{
		Line: int32(m.Line()),
		Directive: &pb.Directive_Memory{
			Memory: &pb.MemoryDirective{
				MemoryGb: m.MemoryGB,
			},
		},
	}
}

var _ Directive = (*MemoryDirective)(nil)
var _ starlark.Value = (*MemoryDirective)(nil)
var _ starlark.HasAttrs = (*MemoryDirective)(nil)

func (m *MemoryDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "memory_gb":
		return starlark.Float(m.MemoryGB), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("memory directive has no attribute %q", name))
	}
}

func (m *MemoryDirective) AttrNames() []string {
	return []string{"memory_gb"}
}

type MemoryDirective struct {
	MemoryGB float64
	line     int
}

func (m *MemoryDirective) Line() int {
	return m.line
}

func (m *MemoryDirective) String() string {
	return fmt.Sprintf("MemoryDirective(MemoryGB: %f)", m.MemoryGB)
}

func (m *MemoryDirective) Type() string {
	return "memory_directive"
}

func (m *MemoryDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (m *MemoryDirective) Truth() starlark.Bool {
	return starlark.Bool(m.MemoryGB != 0)
}

func (m *MemoryDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%f", m.MemoryGB)))
	return h.Sum32(), nil
}

func (m *MemoryDirective) GetGB() (float64, error) {
	return m.MemoryGB, nil
}

func MakeMemoryDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid memory directive")
		}
		expr := exprs[0]

		// Handle string literal case
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				var value float64
				var unit string

				_, err := fmt.Sscanf(strValue, "%f%s", &value, &unit)
				if err != nil {
					return nil, fmt.Errorf("invalid memory format: %s", strValue)
				}
				return convertToMemoryDirective(value, unit, mce.GetLineNumber())
			}
		}

		// Handle numeric literal with unit property (e.g., 100.MB)
		if propExpr, ok := expr.(*parser.PropertyExpression); ok {
			if objExpr, ok := propExpr.GetObjectExpression().(*parser.ConstantExpression); ok {
				if numValue, ok := objExpr.GetValue().(int); ok {
					unit := propExpr.GetProperty()
					return convertToMemoryDirective(float64(numValue), unit.GetText(), mce.GetLineNumber())
				}
			}
		}
	}
	return nil, errors.New("invalid memory directive")
}

func convertToMemoryDirective(value float64, unit string, line int) (*MemoryDirective, error) {
	var memoryGB float64
	switch unit {
	case "B":
		memoryGB = value / (1024 * 1024 * 1024)
	case "KB":
		memoryGB = value / (1024 * 1024)
	case "MB":
		memoryGB = value / 1024
	case "GB":
		memoryGB = value
	case "TB":
		memoryGB = value * 1024
	case "PB":
		memoryGB = value * 1024 * 1024
	case "EB":
		memoryGB = value * 1024 * 1024 * 1024
	case "ZB":
		memoryGB = value * 1024 * 1024 * 1024 * 1024
	default:
		return nil, fmt.Errorf("unknown memory unit: %s", unit)
	}
	return &MemoryDirective{MemoryGB: memoryGB, line: line}, nil
}
