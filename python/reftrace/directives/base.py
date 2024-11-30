from dataclasses import dataclass
from typing import Any

@dataclass(frozen=True)
class Directive:
    """Base class for all directive wrappers."""
    _value: Any  # The protobuf directive
    line: int

    def __getattr__(self, name: str) -> Any:
        """Pass through any unknown attributes to the protobuf."""
        return getattr(self._value, name)
