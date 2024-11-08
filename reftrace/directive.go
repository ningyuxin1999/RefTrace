package main

import "C"
import "reft-go/nf/directives"

//export Directive_Free
func Directive_Free(handle DirectiveHandle) {
	delete(directiveStore, handle)
}

//export Directive_IsContainer
func Directive_IsContainer(handle DirectiveHandle) C.int {
	if directive, ok := directiveStore[handle]; ok {
		_, isContainer := directive.(*directives.Container)
		if isContainer {
			return 1
		}
	}
	return 0
}
