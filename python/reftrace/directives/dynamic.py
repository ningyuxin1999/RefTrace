from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class Dynamic(Directive):
    """The 'dynamic' directive enables or disables dynamic input handling."""
    _value: module_pb2.DynamicDirective

    @property
    def enabled(self) -> bool:
        """Whether dynamic input handling is enabled."""
        return self._value.enabled

__all__ = ['Dynamic']