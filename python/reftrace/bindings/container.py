from .lib import _lib
from typing import Optional

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
    def line(self) -> int:
        return _lib.Container_GetLine(self._handle)
    
    @property
    def names(self) -> list[str]:
        """Get all container names based on the format"""
        if self.format == "simple":
            return [self.simple_name] if self.simple_name else []
        elif self.format == "ternary":
            names = []
            if self.true_name:
                names.append(self.true_name)
            if self.false_name:
                names.append(self.false_name)
            return names
        raise ValueError(f"invalid container format: {self.format}")
        
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