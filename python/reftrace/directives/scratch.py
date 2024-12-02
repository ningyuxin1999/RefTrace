from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

__all__ = ['Scratch']

@dataclass(frozen=True)
class Scratch(Directive):
    """The 'scratch' directive specifies the scratch directory path."""
    _value: module_pb2.ScratchDirective

    @property
    def path(self) -> str:
        """The scratch directory path."""
        return self._value.path