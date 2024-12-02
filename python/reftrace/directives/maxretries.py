from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class MaxRetries(Directive):
    """The 'maxRetries' directive specifies the maximum number of retry attempts."""
    _value: module_pb2.MaxRetriesDirective

    @property
    def num(self) -> int:
        """The maximum number of retry attempts allowed."""
        return self._value.num

__all__ = ['MaxRetries']