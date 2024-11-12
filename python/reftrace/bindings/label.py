from .lib import _lib

class Label:
    def __init__(self, handle: int):
        self._handle = handle

    def __del__(self):
        if hasattr(self, '_handle'):
            _lib.Directive_Free(self._handle)
        
    @property
    def value(self) -> str:
        result = _lib.Label_GetValue(self._handle)
        if result:
            return result.decode('utf-8')
        return ""
    
    @property
    def line(self) -> int:
        return _lib.Label_GetLine(self._handle)
