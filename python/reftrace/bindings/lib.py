import os
import sys
import ctypes
from ctypes import c_char_p, c_int, c_ulonglong, c_void_p

def load_library():
    _lib_dir = os.path.dirname(os.path.abspath(__file__))

    # Platform-specific library names
    if sys.platform == "darwin":
        lib_name = "libreftrace.dylib"
    elif sys.platform == "win32":
        lib_name = "libreftrace.dll"
    else:  # Linux and others
        lib_name = "libreftrace.so"

    _lib_path = os.path.join(_lib_dir, lib_name)
    
    if not os.path.exists(_lib_path):
        raise RuntimeError(f"Library not found at {_lib_path}")
        
    try:
        return ctypes.CDLL(_lib_path)
    except OSError as e:
        raise RuntimeError(f"Failed to load library: {e}") from e

# Load the shared library
_lib = load_library()

# Module function signatures
class ModuleNewResult(ctypes.Structure):
    _fields_ = [
        ("handle", c_ulonglong),
        ("error", c_void_p)
    ]

_lib.Module_New.argtypes = [c_char_p]
_lib.Module_New.restype = ModuleNewResult
_lib.Module_Free.argtypes = [c_ulonglong]
_lib.Module_GetPath.argtypes = [c_ulonglong]
_lib.Module_GetPath.restype = c_char_p
_lib.Module_GetDSLVersion.argtypes = [c_ulonglong]
_lib.Module_GetDSLVersion.restype = c_int
_lib.Module_GetProcessCount.argtypes = [c_ulonglong]
_lib.Module_GetProcessCount.restype = c_int
_lib.Module_GetProcess.argtypes = [c_ulonglong, c_int]
_lib.Module_GetProcess.restype = c_ulonglong
_lib.Module_Free_Error.argtypes = [c_void_p]

# Process function signatures
_lib.Process_Free.argtypes = [c_ulonglong]
_lib.Process_GetName.argtypes = [c_ulonglong]
_lib.Process_GetName.restype = c_char_p
_lib.Process_GetDirectiveCount.argtypes = [c_ulonglong]
_lib.Process_GetDirectiveCount.restype = c_int
_lib.Process_GetDirective.argtypes = [c_ulonglong, c_int]
_lib.Process_GetDirective.restype = c_ulonglong

# Container function signatures
_lib.Container_GetFormat.argtypes = [c_ulonglong]
_lib.Container_GetFormat.restype = c_char_p
_lib.Container_GetSimpleName.argtypes = [c_ulonglong]
_lib.Container_GetSimpleName.restype = c_char_p
_lib.Container_GetCondition.argtypes = [c_ulonglong]
_lib.Container_GetCondition.restype = c_char_p
_lib.Container_GetTrueName.argtypes = [c_ulonglong]
_lib.Container_GetTrueName.restype = c_char_p
_lib.Container_GetFalseName.argtypes = [c_ulonglong]
_lib.Container_GetFalseName.restype = c_char_p
_lib.Container_GetLine.argtypes = [c_ulonglong]
_lib.Container_GetLine.restype = c_int

# Label function signatures
_lib.Directive_IsLabel.argtypes = [c_ulonglong]
_lib.Directive_IsLabel.restype = c_int
_lib.Label_GetValue.argtypes = [c_ulonglong]
_lib.Label_GetValue.restype = c_char_p
_lib.Label_GetLine.argtypes = [c_ulonglong]
_lib.Label_GetLine.restype = c_int

# Add free function signature
_lib.free.argtypes = [c_void_p]

# Add to directive function signatures
_lib.Directive_IsContainer.argtypes = [c_ulonglong]
_lib.Directive_IsContainer.restype = c_int
_lib.Directive_Free.argtypes = [c_ulonglong]