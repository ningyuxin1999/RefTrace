from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class MaxForks(Directive):
    """The 'maxForks' directive specifies the maximum number of parallel process instances."""
    _value: module_pb2.MaxForksDirective

    @property
    def num(self) -> int:
        """The maximum number of parallel forks allowed."""
        return self._value.num

__all__ = ['MaxForks']