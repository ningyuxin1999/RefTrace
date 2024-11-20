from .lib import _lib
import ctypes
from ctypes import c_char_p
from .process import Process

class Module:
    """Represents a module that contains process definitions.
    
    This class provides an interface to load and interact with module files that contain
    process definitions for reference tracing.
    """

    def __init__(self, filepath: str):
        """Initialize a new Module instance.
        
        Args:
            filepath: Path to the module file to load.
            
        Raises:
            RuntimeError: If the module file cannot be loaded or is invalid.
        """
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
        """Get the file path of the loaded module.
        
        Returns:
            str: The file path of the module. Returns empty string if path cannot be retrieved.
        """
        result = _lib.Module_GetPath(self._handle)
        if result:
            return result.decode('utf-8')
        return ""

    @property
    def dsl_version(self) -> int:
        """Get the DSL (Domain Specific Language) version of the module.
        
        Returns:
            int: The version number of the DSL used in this module.
        """
        return _lib.Module_GetDSLVersion(self._handle)

    @property
    def processes(self) -> list[Process]:
        """Get all processes defined in this module.
        
        Returns:
            list[Process]: A list of Process objects defined in this module.
        """
        count = _lib.Module_GetProcessCount(self._handle)
        result = []
        for i in range(count):
            handle = _lib.Module_GetProcess(self._handle, i)
            if handle:
                result.append(Process(handle))
        return result