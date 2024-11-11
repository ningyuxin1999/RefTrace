package main

import (
	"C"
	"fmt"
)

// String conversion helpers
func cString(s string) *C.char {
	return C.CString(s)
}

func goString(s *C.char) string {
	return C.GoString(s)
}

// Integer conversion helper
func cInt(i int) C.int {
	return C.int(i)
}

// Helper functions that don't use cgo directly
func createModuleFromFile(path string) (ModuleHandle, error) {
	result := Module_New(cString(path))
	if result.error != nil {
		err := goString(result.error)
		Module_Free_Error(result.error)
		return 0, fmt.Errorf(err)
	}
	return ModuleHandle(result.handle), nil
}

func getProcessName(handle ProcessHandle) string {
	if name := Process_GetName(handle); name != nil {
		return goString(name)
	}
	return ""
}

func getContainerFormat(handle DirectiveHandle) string {
	if format := Container_GetFormat(handle); format != nil {
		return goString(format)
	}
	return ""
}

func getContainerCondition(handle DirectiveHandle) string {
	if cond := Container_GetCondition(handle); cond != nil {
		return goString(cond)
	}
	return ""
}

func getContainerTrueName(handle DirectiveHandle) string {
	if name := Container_GetTrueName(handle); name != nil {
		return goString(name)
	}
	return ""
}

func getContainerFalseName(handle DirectiveHandle) string {
	if name := Container_GetFalseName(handle); name != nil {
		return goString(name)
	}
	return ""
}
