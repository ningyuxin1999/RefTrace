from ..proto import common_pb2, module_pb2
from .lib import _lib
import ctypes
import base64
from typing import List, Optional
from dataclasses import dataclass
from functools import cached_property
from .process import Process

@dataclass
class Module:
    """Represents a Nextflow module that contains process definitions."""
    _proto: module_pb2.Module  # Internal protobuf representation

    @classmethod
    def from_file(cls, filepath: str) -> 'ModuleResult':
        encoded_path = filepath.encode('utf-8')
        result_ptr = _lib.Module_New(encoded_path)
        if not result_ptr:
            return ModuleResult(filepath=filepath, module=None, error=common_pb2.ParseError(likely_rt_bug=True, error="Failed to create module"))
            
        try:
            # Get base64 string from pointer and decode it
            encoded_str = ctypes.cast(result_ptr, ctypes.c_char_p).value.decode('utf-8')
            bytes_data = base64.b64decode(encoded_str)
            
            result: module_pb2.ModuleResult = module_pb2.ModuleResult()
            result.ParseFromString(bytes_data)
            
            if result.HasField('error'):
                return ModuleResult(filepath=filepath, module=None, error=result.error)
                
            return ModuleResult(filepath=filepath, module=cls(_proto=result.module), error=None)
        finally:
            _lib.Module_Free(result_ptr)

    @cached_property
    def path(self) -> str:
        """The file path of the module."""
        return self._proto.path

    @cached_property
    def dsl_version(self) -> int:
        """The DSL (Domain Specific Language) version of the module."""
        return self._proto.dsl_version

    @cached_property
    def processes(self) -> List[Process]:
        """All processes defined in this module."""
        return [Process(_proto=p) for p in self._proto.processes]
    
@dataclass
class ParseError:
    """An error can either come from user input (malformed Nextflow files) or may be a bug in RefTrace."""
    error: str
    likely_rt_bug: bool

@dataclass
class ModuleResult:
    """Result type for Module creation that can contain either a Module or an error."""
    filepath: str
    module: Optional[Module]
    error: Optional[ParseError]
       
def parse_modules(directory, progress_callback=None) -> List[ModuleResult]:
    """
    Parse all Nextflow modules in a directory.
    
    Args:
        directory (str): Path to directory containing .nf files
        progress_callback (callable): Optional callback function(current, total)
    
    Returns:
        List[ModuleResult]: List of results, each containing either a Module or ParseError
    """
    if progress_callback is None:
        callback_ptr = None
    else:
        CALLBACK_TYPE = ctypes.CFUNCTYPE(None, ctypes.c_int32, ctypes.c_int32)
        callback_ptr = CALLBACK_TYPE(progress_callback)
    
    result_ptr = _lib.Parse_Modules(
        directory.encode('utf-8'),
        callback_ptr
    )
    
    if not result_ptr:
        raise RuntimeError("Failed to process modules")
        
    try:
        # Get base64 string from pointer and decode it
        encoded_str = ctypes.cast(result_ptr, ctypes.c_char_p).value.decode('utf-8')
        bytes_data = base64.b64decode(encoded_str)
        
        proto_result = module_pb2.ModuleListResult()
        proto_result.ParseFromString(bytes_data)
        
        results = []
        for result in proto_result.results:
            if result.HasField('module'):
                # Success case - create Module instance
                results.append(ModuleResult(
                    filepath=result.file_path,
                    module=Module(_proto=result.module),
                    error=None
                ))
            else:
                # Error case - create ParseError
                results.append(ModuleResult(
                    filepath=result.file_path,
                    module=None,
                    error=ParseError(
                        error=result.error.error,
                        likely_rt_bug=result.error.likely_rt_bug
                    )
                ))
        
        return results
    finally:
        _lib.Module_Free(result_ptr)
