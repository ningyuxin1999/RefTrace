from ..proto import common_pb2, module_pb2
from .lib import _lib
import ctypes
import base64
from typing import List, Optional, Union
from dataclasses import dataclass
from functools import cached_property
from .process import Process
import json

@dataclass
class ResolvedInclude:
    _proto: module_pb2.ResolvedInclude

    @property
    def module_path(self) -> str:
        return self._proto.module_path

    @property
    def includes(self) -> List[str]:
        return list(self._proto.includes)

    def to_dict(self) -> dict:
        return {
            "module_path": self.module_path,
            "includes": self.includes
        }

@dataclass
class UnresolvedInclude:
    _proto: module_pb2.UnresolvedInclude

    @property
    def module_path(self) -> str:
        return self._proto.module_path

    @property
    def includes(self) -> List[str]:
        return list(self._proto.includes)

    def to_dict(self) -> dict:
        return {
            "module_path": self.module_path,
            "includes": self.includes
        }

@dataclass
class Workflow:
    _proto: module_pb2.Workflow

    @property
    def name(self) -> str:
        """The name of the workflow."""
        return self._proto.name

    @property
    def takes(self) -> List[str]:
        """The takes of the workflow."""
        return list(self._proto.takes)

    @property
    def emits(self) -> List[str]:
        """The emits of the workflow."""
        return list(self._proto.emits)

@dataclass
class Param:
    _proto: module_pb2.Param

    @property
    def line(self) -> int:
        """The line number where this param is defined."""
        return self._proto.line

    @property
    def name(self) -> str:
        """The name of the param."""
        return self._proto.name

@dataclass
class IncludedItem:
    _proto: module_pb2.IncludedItem

    @property
    def name(self) -> str:
        """The name of the included item."""
        return self._proto.name

    @property
    def alias(self) -> Optional[str]:
        """The alias of the included item, if any."""
        return self._proto.alias or None

@dataclass
class IncludeStatement:
    _proto: module_pb2.IncludeStatement

    @property
    def line(self) -> int:
        """The line number where this include statement is defined."""
        return self._proto.line

    @property
    def items(self) -> List[IncludedItem]:
        """The items included in this include statement."""
        return [IncludedItem(_proto=i) for i in self._proto.items]

    @property
    def from_module(self) -> str:
        """The module path that this include statement is from."""
        return self._proto.from_module

@dataclass
class ParseError:
    """An error can either come from user input (malformed Nextflow files) or may be a bug in RefTrace."""
    error: str
    likely_rt_bug: bool
    path: str

@dataclass
class Module:
    """Represents a Nextflow module that contains process definitions."""
    _proto: module_pb2.Module  # Internal protobuf representation

    @classmethod
    def from_file(cls, filepath: str) -> Union['Module', ParseError]:
        encoded_path = filepath.encode('utf-8')
        result_ptr = _lib.Module_New(encoded_path)
        if not result_ptr:
            # would be an issue with the C API
            return ParseError(path=filepath, likely_rt_bug=True, error="Failed to create module")
            
        try:
            # Get base64 string from pointer and decode it
            encoded_str = ctypes.cast(result_ptr, ctypes.c_char_p).value.decode('utf-8')
            bytes_data = base64.b64decode(encoded_str)
            
            result: module_pb2.ModuleResult = module_pb2.ModuleResult()
            result.ParseFromString(bytes_data)
            
            if result.HasField('error'):
                return ParseError(path=filepath, likely_rt_bug=result.error.likely_rt_bug, error=result.error.error)
                
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
    
    @property
    def includes(self) -> List[IncludeStatement]:
        """All includes defined in this module."""
        return [IncludeStatement(_proto=i) for i in self._proto.includes]
    
    @property
    def params(self) -> List[Param]:
        """All params defined in this module."""
        return [Param(_proto=p) for p in self._proto.params]
    
    @property
    def workflows(self) -> List[Workflow]:
        """All workflows defined in this module."""
        return [Workflow(_proto=w) for w in self._proto.workflows]
    
    def to_dict(self, only_paths: bool = False) -> dict:
        """Convert the module to a dictionary representation."""
        if only_paths:
            return {
                "path": self.path,
            }
        return {
            "path": self.path,
            "dsl_version": self.dsl_version,
            "includes": [
                {
                    "line": inc.line,
                    "from_module": inc.from_module,
                    "items": [
                        {"name": item.name, "alias": item.alias}
                        for item in inc.items
                    ]
                }
                for inc in self.includes
            ],
            "params": [
                {"name": param.name, "line": param.line}
                for param in self.params
            ],
            "workflows": [
                {
                    "name": wf.name,
                    "takes": wf.takes,
                    "emits": wf.emits
                }
                for wf in self.workflows
            ],
            "processes": [
                {
                    "name": proc.name
                }
                for proc in self.processes
            ]
        }
    
    def to_json(self, indent: int = 2) -> str:
        """Convert the module to JSON format."""
        return json.dumps(self.to_dict(), indent=indent)

@dataclass
class ModuleListResult:
    results: List[Module]
    errors: List[ParseError]
    resolved_includes: List[ResolvedInclude]
    unresolved_includes: List[UnresolvedInclude]
       
def parse_modules(directory, progress_callback=None) -> ModuleListResult:
    """
    Parse all Nextflow modules in a directory.
    
    Args:
        directory (str): Path to directory containing .nf files
        progress_callback (callable): Optional callback function(current, total)
    
    Returns:
        ModuleListResult: List of results, each containing either a Module or ParseError
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
        errors = []
        for result in proto_result.results:
            if result.HasField('module'):
                # Success case - create Module instance
                results.append(Module(_proto=result.module))
            else:
                # Error case - create ParseError
                parse_error = ParseError(
                    path=result.file_path,
                    likely_rt_bug=result.error.likely_rt_bug,
                    error=result.error.error
                )
                errors.append(parse_error)

        return ModuleListResult(
            results=results,
            errors=errors,
            resolved_includes=[ResolvedInclude(_proto=i) for i in proto_result.resolved_includes],
            unresolved_includes=[UnresolvedInclude(_proto=i) for i in proto_result.unresolved_includes]
        )
    finally:
        _lib.Module_Free(result_ptr)
