from .lib import _lib
from .container import Container

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
    def containers(self) -> list[Container]:
        count = _lib.Process_GetDirectiveCount(self._handle)
        result = []
        for i in range(count):
            handle = _lib.Process_GetDirective(self._handle, i)
            if handle and _lib.Directive_IsContainer(handle):
                result.append(Container(handle))
        return result