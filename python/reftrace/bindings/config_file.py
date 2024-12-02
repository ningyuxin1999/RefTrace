from ..proto import config_file_pb2
from .lib import _lib
import ctypes
import base64
from typing import List
from dataclasses import dataclass
from functools import cached_property
from ..configfile import ProcessScope

@dataclass
class ConfigFile:
    """Represents a Nextflow config file that contains process configurations."""
    _proto: config_file_pb2.ConfigFile  # Internal protobuf representation

    @classmethod
    def from_file(cls, filepath: str) -> 'ConfigFile':
        encoded_path = filepath.encode('utf-8')
        result_ptr = _lib.ConfigFile_New(encoded_path)
        if not result_ptr:
            raise RuntimeError("Failed to create config file")
            
        try:
            # Get base64 string from pointer and decode it
            encoded_str = ctypes.cast(result_ptr, ctypes.c_char_p).value.decode('utf-8')
            bytes_data = base64.b64decode(encoded_str)
            
            result = config_file_pb2.ConfigFileResult()
            result.ParseFromString(bytes_data)
            
            if result.HasField('error'):
                raise RuntimeError(result.error)
                
            return cls(_proto=result.config_file)
        finally:
            _lib.ConfigFile_Free(result_ptr)

    @cached_property
    def path(self) -> str:
        """The file path of the config file."""
        return self._proto.path

    @cached_property
    def process_scopes(self) -> list[ProcessScope]:
        """All process scopes defined in this config file."""
        return [ProcessScope(scope) for scope in self._proto.process_scopes]
