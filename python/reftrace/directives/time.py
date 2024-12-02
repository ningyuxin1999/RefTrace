from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class Time(Directive):
    """The 'time' directive specifies the maximum execution time allowed."""
    _value: module_pb2.TimeDirective

    @property
    def time(self) -> str:
        """The time limit (e.g., '1h', '30m', '2d')."""
        return self._value.time

__all__ = ['Time']