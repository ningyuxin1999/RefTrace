from .lib import _lib
from .container import Container
from .label import Label
class Process:
    """Represents a process definition within a module.
    
    This class provides access to process properties and its contained directives
    such as containers and labels.
    """

    def __init__(self, handle: int):
        """Initialize a new Process instance.
        
        Args:
            handle: Internal handle identifier for the process.
        """
        self._handle = handle

    def __del__(self):
        """Cleanup method to free the process handle when the object is destroyed."""
        if hasattr(self, '_handle'):
            _lib.Process_Free(self._handle)

    @property
    def name(self) -> str:
        """Get the name of the process.
        
        Returns:
            str: The name of the process. Returns empty string if name cannot be retrieved.
        """
        result = _lib.Process_GetName(self._handle)
        if result:
            return result.decode('utf-8')
        return ""
    
    @property
    def line(self) -> int:
        """Get the line number where this process is defined.
        
        Returns:
            int: The line number in the source file where this process is defined.
        """
        return _lib.Process_GetLine(self._handle)

    @property
    def containers(self) -> list[Container]:
        """Get all container directives within this process.
        
        Returns:
            list[Container]: A list of Container objects defined within this process.
        """
        count = _lib.Process_GetDirectiveCount(self._handle)
        result = []
        for i in range(count):
            handle = _lib.Process_GetDirective(self._handle, i)
            if handle and _lib.Directive_IsContainer(handle):
                result.append(Container(handle))
        return result
    
    @property
    def labels(self) -> list[Label]:
        """Get all label directives within this process.
        
        Returns:
            list[Label]: A list of Label objects defined within this process.
        """
        count = _lib.Process_GetDirectiveCount(self._handle)
        result = []
        for i in range(count):
            directive_handle = _lib.Process_GetDirective(self._handle, i)
            if directive_handle and _lib.Directive_IsLabel(directive_handle):
                result.append(Label(directive_handle))
        return result
