from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class MaxErrors(Directive):
    """The 'maxErrors' directive specifies the maximum number of errors allowed."""
    _value: module_pb2.MaxErrorsDirective

    @property
    def num(self) -> int:
        """The maximum number of errors allowed."""
        return self._value.num

__all__ = ['MaxErrors']