package main

import "C"
import "reft-go/nf/directives"

//export Label_GetValue
func Label_GetValue(handle DirectiveHandle) *C.char {
	if directive, ok := directiveStore[handle]; ok {
		if label, ok := directive.(*directives.LabelDirective); ok {
			return C.CString(label.Label)
		}
	}
	return nil
}

//export Directive_IsLabel
func Directive_IsLabel(handle DirectiveHandle) C.int {
	if directive, ok := directiveStore[handle]; ok {
		_, isLabel := directive.(*directives.LabelDirective)
		if isLabel {
			return 1
		}
	}
	return 0
}
