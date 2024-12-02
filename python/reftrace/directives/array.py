from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

__all__ = ['Array']

@dataclass(frozen=True)
class Array(Directive):
    """The 'array' directive specifies array job size."""
    _value: module_pb2.ArrayDirective

    @property
    def size(self) -> int:
        """The size of the array job."""
        return self._value.size