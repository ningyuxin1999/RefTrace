package main

import "C"

//export Process_Free
func Process_Free(handle ProcessHandle) {
	if process, ok := processStore[handle]; ok {
		// Free all directives associated with this process
		for dirHandle, dir := range directiveStore {
			for _, procDir := range process.Directives {
				if dir == procDir {
					delete(directiveStore, dirHandle)
				}
			}
		}
		delete(processStore, handle)
	}
}

//export Process_GetName
func Process_GetName(handle ProcessHandle) *C.char {
	if process, ok := processStore[handle]; ok {
		return C.CString(process.Name)
	}
	return nil
}

//export Process_GetDirectiveCount
func Process_GetDirectiveCount(handle ProcessHandle) C.int {
	if process, ok := processStore[handle]; ok {
		return C.int(len(process.Directives))
	}
	return 0
}

//export Process_GetDirective
func Process_GetDirective(handle ProcessHandle, index C.int) DirectiveHandle {
	if process, ok := processStore[handle]; ok {
		if idx := int(index); idx >= 0 && idx < len(process.Directives) {
			// Check if this directive already has a handle
			directive := process.Directives[idx]
			for existingHandle, existingDirective := range directiveStore {
				if existingDirective == directive {
					return existingHandle
				}
			}
			// If not found, create new handle
			handle := nextDirectiveHandle
			nextDirectiveHandle++
			directiveStore[handle] = directive
			return handle
		}
	}
	return 0
}
