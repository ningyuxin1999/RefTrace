package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*PublishDirDirective)(nil)
var _ starlark.Value = (*PublishDirDirective)(nil)
var _ starlark.HasAttrs = (*PublishDirDirective)(nil)

func (p *PublishDirDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "path":
		if p.Path != "" {
			return starlark.String(p.Path), nil
		}
		return starlark.None, nil
	case "params":
		if p.Params != "" {
			return starlark.String(p.Params), nil
		}
		return starlark.None, nil
	case "contentType":
		if p.ContentType != nil {
			return starlark.Bool(*p.ContentType), nil
		}
		return starlark.None, nil
	case "enabled":
		if p.Enabled != nil {
			return starlark.Bool(*p.Enabled), nil
		}
		return starlark.None, nil
	case "failOnError":
		if p.FailOnError != nil {
			return starlark.Bool(*p.FailOnError), nil
		}
		return starlark.None, nil
	case "mode":
		return starlark.String(p.Mode), nil
	case "overwrite":
		if p.Overwrite != nil {
			return starlark.Bool(*p.Overwrite), nil
		}
		return starlark.None, nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("publish_dir directive has no attribute %q", name))
	}
}

func (p *PublishDirDirective) AttrNames() []string {
	return []string{"path", "params", "contentType", "enabled", "failOnError", "mode", "overwrite"}
}

type PublishDirDirective struct {
	Path        string
	Params      string
	ContentType *bool
	Enabled     *bool
	FailOnError *bool
	Mode        string
	Overwrite   *bool
}

func (p *PublishDirDirective) String() string {
	pathStr := p.Path
	if p.Params != "" {
		pathStr = fmt.Sprintf("params.%s", p.Params)
	}

	return fmt.Sprintf("PublishDirDirective(Path: %q, Params: %q, ContentType: %v, Enabled: %v, FailOnError: %v, Mode: %q, Overwrite: %v)",
		pathStr,
		p.Params,
		boolPtrToString(p.ContentType),
		boolPtrToString(p.Enabled),
		boolPtrToString(p.FailOnError),
		p.Mode,
		boolPtrToString(p.Overwrite))
}

func (p *PublishDirDirective) Type() string {
	return "publish_dir_directive"
}

func (p *PublishDirDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (p *PublishDirDirective) Truth() starlark.Bool {
	return starlark.Bool(p.Path != "")
}

func (p *PublishDirDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(p.Path))
	h.Write([]byte(boolPtrToString(p.ContentType)))
	h.Write([]byte(boolPtrToString(p.Enabled)))
	h.Write([]byte(boolPtrToString(p.FailOnError)))
	h.Write([]byte(p.Mode))
	h.Write([]byte(boolPtrToString(p.Overwrite)))
	return h.Sum32(), nil
}

// Helper function to convert *bool to string
func boolPtrToString(b *bool) string {
	if b == nil {
		return "nil"
	}
	return fmt.Sprintf("%t", *b)
}

func MakePublishDirDirective(mce *parser.MethodCallExpression) (Directive, error) {
	var dir string = ""
	var paramName string = ""
	var contentType *bool = nil
	var enabled *bool = nil
	var failOnError *bool = nil
	var mode string = ""
	var overwrite *bool = nil

	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		for _, expr := range exprs {
			// Handle direct arguments
			switch e := expr.(type) {
			case *parser.ConstantExpression:
				if value, ok := e.GetValue().(string); ok {
					dir = value
				}
			case *parser.GStringExpression:
				dir = e.GetText()
			case *parser.PropertyExpression:
				// Check if this is a params reference
				if ve, ok := e.GetObjectExpression().(*parser.VariableExpression); ok &&
					ve.GetText() == "params" {
					if prop, ok := e.GetProperty().(*parser.ConstantExpression); ok {
						paramName = prop.GetText()
						dir = "" // Clear dir as we're using params
					}
				}
			case *parser.MapExpression:
				entries := e.GetMapEntryExpressions()
				for _, entry := range entries {
					key := entry.GetKeyExpression().GetText()
					valueExpr := entry.GetValueExpression()

					switch key {
					case "path":
						switch ve := valueExpr.(type) {
						case *parser.ConstantExpression:
							dir = ve.GetText()
						case *parser.GStringExpression:
							dir = ve.GetText()
						case *parser.PropertyExpression:
							// Handle cases like params.output_dir in the path option
							if obj, ok := ve.GetObjectExpression().(*parser.VariableExpression); ok &&
								obj.GetText() == "params" {
								if prop, ok := ve.GetProperty().(*parser.ConstantExpression); ok {
									paramName = prop.GetText()
									dir = "" // Clear dir as we're using params
								}
							}
						}
					case "contentType":
						if ce, ok := valueExpr.(*parser.ConstantExpression); ok {
							if v, ok := ce.GetValue().(bool); ok {
								contentType = &v
							}
						}
					case "enabled":
						if ce, ok := valueExpr.(*parser.ConstantExpression); ok {
							if v, ok := ce.GetValue().(bool); ok {
								enabled = &v
							}
						}
					case "failOnError":
						if ce, ok := valueExpr.(*parser.ConstantExpression); ok {
							if v, ok := ce.GetValue().(bool); ok {
								failOnError = &v
							}
						}
					case "mode":
						if ce, ok := valueExpr.(*parser.ConstantExpression); ok {
							mode = ce.GetText()
						}
					case "overwrite":
						if ce, ok := valueExpr.(*parser.ConstantExpression); ok {
							if v, ok := ce.GetValue().(bool); ok {
								overwrite = &v
							}
						}
					}
				}
			}
		}
	}

	// Validate that we have either a path or a params reference
	if dir == "" && paramName == "" {
		return nil, errors.New("invalid publish dir directive: no valid path specified")
	}

	return &PublishDirDirective{
		Path:        dir,
		Params:      paramName,
		ContentType: contentType,
		Enabled:     enabled,
		FailOnError: failOnError,
		Mode:        mode,
		Overwrite:   overwrite,
	}, nil
}
