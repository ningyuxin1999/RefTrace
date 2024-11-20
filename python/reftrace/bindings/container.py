from .lib import _lib
from typing import Optional

class Container:
    """Represents a container directive.
    
    A container can be in different formats:
    - "simple": A basic container with a single name
    - "ternary": A conditional container with true/false branches and associated names
    """

    def __init__(self, handle: int):
        """Initialize a new Container instance.
        
        Args:
            handle: Internal handle identifier for the container.
        """
        self._handle = handle

    def __del__(self):
        """Cleanup method to free the container handle when the object is destroyed."""
        if hasattr(self, '_handle'):
            _lib.Directive_Free(self._handle)
        
    @property
    def format(self) -> str:
        """Get the format type of the container.
        
        Returns:
            str: The container format ("simple" or "ternary"). Returns empty string if format cannot be retrieved.
        """
        result = _lib.Container_GetFormat(self._handle)
        if result:
            return result.decode('utf-8')
        return ""
    
    @property
    def line(self) -> int:
        """Get the line number where this container is defined.
        
        Returns:
            int: The line number in the source file where this container is defined.
        """
        return _lib.Container_GetLine(self._handle)
    
    @property
    def names(self) -> list[str]:
        """Get all container names based on the format.
        
        Returns:
            list[str]: A list of names associated with this container. For simple containers,
                      returns a single-item list. For ternary containers, returns both true
                      and false branch names if they exist.
        
        Raises:
            ValueError: If the container format is invalid.
        """
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
        """Get the name of a simple format container.
        
        Returns:
            Optional[str]: The container name if it exists and is in simple format, None otherwise.
        """
        result = _lib.Container_GetSimpleName(self._handle)
        if result:
            return result.decode('utf-8')
        return None
        
    @property
    def condition(self) -> Optional[str]:
        """Get the condition expression for a ternary container.
        
        Returns:
            Optional[str]: The condition expression if it exists, None otherwise.
        """
        result = _lib.Container_GetCondition(self._handle)
        if result:
            return result.decode('utf-8')
        return None
        
    @property
    def true_name(self) -> Optional[str]:
        """Get the name associated with the true branch of a ternary container.
        
        Returns:
            Optional[str]: The name for the true condition if it exists, None otherwise.
        """
        result = _lib.Container_GetTrueName(self._handle)
        if result:
            return result.decode('utf-8')
        return None
        
    @property
    def false_name(self) -> Optional[str]:
        """Get the name associated with the false branch of a ternary container.
        
        Returns:
            Optional[str]: The name for the false condition if it exists, None otherwise.
        """
        result = _lib.Container_GetFalseName(self._handle)
        if result:
            return result.decode('utf-8')
        return None