from ..proto import module_pb2
from .lib import _lib
import ctypes
import base64
from typing import List
from dataclasses import dataclass
from functools import cached_property
from .process import Process

@dataclass
class Module:
    """Represents a Nextflow module that contains process definitions."""
    _proto: module_pb2.Module  # Internal protobuf representation

    @classmethod
    def from_file(cls, filepath: str) -> 'Module':
        encoded_path = filepath.encode('utf-8')
        result_ptr = _lib.Module_New(encoded_path)
        if not result_ptr:
            raise RuntimeError("Failed to create module")
            
        try:
            # Get base64 string from pointer and decode it
            encoded_str = ctypes.cast(result_ptr, ctypes.c_char_p).value.decode('utf-8')
            bytes_data = base64.b64decode(encoded_str)
            
            result = module_pb2.ModuleResult()
            result.ParseFromString(bytes_data)
            
            if result.HasField('error'):
                raise RuntimeError(result.error)
                
            return cls(_proto=result.module)
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