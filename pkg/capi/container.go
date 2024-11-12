package main

import "C"
import "reft-go/nf/directives"

//export Container_GetFormat
func Container_GetFormat(handle DirectiveHandle) *C.char {
	if directive, ok := directiveStore[handle]; ok {
		if container, ok := directive.(*directives.Container); ok {
			return C.CString(string(container.Format))
		}
	}
	return nil
}

//export Container_GetLine
func Container_GetLine(handle DirectiveHandle) C.int {
	if directive, ok := directiveStore[handle]; ok {
		if container, ok := directive.(*directives.Container); ok {
			return C.int(container.Line())
		}
	}
	return 0
}

//export Container_GetSimpleName
func Container_GetSimpleName(handle DirectiveHandle) *C.char {
	if directive, ok := directiveStore[handle]; ok {
		if container, ok := directive.(*directives.Container); ok {
			if container.Format == directives.Simple {
				return C.CString(container.SimpleName)
			}
		}
	}
	return nil
}

//export Container_GetCondition
func Container_GetCondition(handle DirectiveHandle) *C.char {
	if directive, ok := directiveStore[handle]; ok {
		if container, ok := directive.(*directives.Container); ok {
			if container.Format == directives.Ternary {
				return C.CString(container.Condition)
			}
		}
	}
	return nil
}

//export Container_GetTrueName
func Container_GetTrueName(handle DirectiveHandle) *C.char {
	if directive, ok := directiveStore[handle]; ok {
		if container, ok := directive.(*directives.Container); ok {
			if container.Format == directives.Ternary {
				return C.CString(container.TrueName)
			}
		}
	}
	return nil
}

//export Container_GetFalseName
func Container_GetFalseName(handle DirectiveHandle) *C.char {
	if directive, ok := directiveStore[handle]; ok {
		if container, ok := directive.(*directives.Container); ok {
			if container.Format == directives.Ternary {
				return C.CString(container.FalseName)
			}
		}
	}
	return nil
}
