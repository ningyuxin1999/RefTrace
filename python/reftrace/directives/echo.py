from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class Echo(Directive):
    """The 'echo' directive enables or disables command echoing."""
    _value: module_pb2.EchoDirective

    @property
    def enabled(self) -> bool:
        """Whether command echoing is enabled."""
        return self._value.enabled

__all__ = ['Echo']