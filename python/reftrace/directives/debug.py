from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class Debug(Directive):
    """The 'debug' directive enables or disables debug mode for a process."""
    _value: module_pb2.DebugDirective

    @property
    def enabled(self) -> bool:
        """Whether debug mode is enabled."""
        return self._value.enabled

__all__ = ['Debug']