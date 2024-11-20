from .lib import _lib

class Label:
    """Represents a label directive.
    
    Labels are used to mark specific positions in the process definition that can be
    referenced by other directives.
    """

    def __init__(self, handle: int):
        """Initialize a new Label instance.
        
        Args:
            handle: Internal handle identifier for the label.
        """
        self._handle = handle

    def __del__(self):
        """Cleanup method to free the label handle when the object is destroyed."""
        if hasattr(self, '_handle'):
            _lib.Directive_Free(self._handle)
        
    @property
    def value(self) -> str:
        """Get the name/value of the label.
        
        Returns:
            str: The label's name. Returns empty string if the value cannot be retrieved.
        """
        result = _lib.Label_GetValue(self._handle)
        if result:
            return result.decode('utf-8')
        return ""
    
    @property
    def line(self) -> int:
        """Get the line number where this label is defined.
        
        Returns:
            int: The line number in the source file where this label is defined.
        """
        return _lib.Label_GetLine(self._handle)
