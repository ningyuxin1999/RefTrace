import os
import ctypes
from ctypes import c_char_p, c_int, c_ulonglong, c_void_p
from typing import Optional

# Load the shared library
_lib_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
_lib = ctypes.CDLL(os.path.join(_lib_dir, 'libreftrace.so'))

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

# Add free function signature
_lib.free.argtypes = [c_void_p]

# Add to directive function signatures
_lib.Directive_IsContainer.argtypes = [c_ulonglong]
_lib.Directive_IsContainer.restype = c_int

class Container:
    def __init__(self, handle: int):
        self._handle = handle

    def __del__(self):
        if hasattr(self, '_handle'):
            _lib.Directive_Free(self._handle)
        
    @property
    def format(self) -> str:
        result = _lib.Container_GetFormat(self._handle)
        if result:
            return result.decode('utf-8')
        return ""
        
    @property
    def simple_name(self) -> Optional[str]:
        result = _lib.Container_GetSimpleName(self._handle)
        if result:
            return result.decode('utf-8')
        return None
        
    @property
    def condition(self) -> Optional[str]:
        result = _lib.Container_GetCondition(self._handle)
        if result:
            return result.decode('utf-8')
        return None
        
    @property
    def true_name(self) -> Optional[str]:
        result = _lib.Container_GetTrueName(self._handle)
        if result:
            return result.decode('utf-8')
        return None
        
    @property
    def false_name(self) -> Optional[str]:
        result = _lib.Container_GetFalseName(self._handle)
        if result:
            return result.decode('utf-8')
        return None

class Process:
    def __init__(self, handle: int):
        self._handle = handle

    def __del__(self):
        if hasattr(self, '_handle'):
            _lib.Process_Free(self._handle)

    @property
    def name(self) -> str:
        result = _lib.Process_GetName(self._handle)
        if result:
            return result.decode('utf-8')
        return ""

    @property
    def directives(self) -> list[Container]:
        count = _lib.Process_GetDirectiveCount(self._handle)
        result = []
        for i in range(count):
            handle = _lib.Process_GetDirective(self._handle, i)
            if handle and _lib.Directive_IsContainer(handle):
                result.append(Container(handle))
        return result

class Module:
    def __init__(self, filepath: str):
        encoded_path = filepath.encode('utf-8')
        result = _lib.Module_New(encoded_path)
        if result.error:
            # if we want exact addresses, we have to use c_void_p
            # python's ctypes will auto convert c_char_p to a python bytes object
            error_msg = ctypes.cast(result.error, c_char_p).value.decode('utf-8')
            _lib.Module_Free_Error(result.error)  # Pass the raw pointer
            raise RuntimeError(error_msg)
        self._handle = result.handle

    def __del__(self):
        if hasattr(self, '_handle'):
            _lib.Module_Free(self._handle)

    @property
    def path(self) -> str:
        result = _lib.Module_GetPath(self._handle)
        if result:
            return result.decode('utf-8')
        return ""

    @property
    def dsl_version(self) -> int:
        return _lib.Module_GetDSLVersion(self._handle)

    @property
    def processes(self) -> list[Process]:
        count = _lib.Module_GetProcessCount(self._handle)
        result = []
        for i in range(count):
            handle = _lib.Module_GetProcess(self._handle, i)
            if handle:
                result.append(Process(handle))
        return result