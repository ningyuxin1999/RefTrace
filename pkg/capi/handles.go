package main

import "C"
import (
	"reft-go/nf"
	"reft-go/nf/directives"
)

func main() {} // Required for C shared library

// Handle types for opaque pointers
type ModuleHandle C.ulonglong
type ProcessHandle C.ulonglong
type DirectiveHandle C.ulonglong

// Stores for managing objects
var (
	moduleStore                         = make(map[ModuleHandle]*Module)
	processStore                        = make(map[ProcessHandle]*nf.Process)
	directiveStore                      = make(map[DirectiveHandle]directives.Directive)
	nextModuleHandle    ModuleHandle    = 1
	nextProcessHandle   ProcessHandle   = 1
	nextDirectiveHandle DirectiveHandle = 1
)
