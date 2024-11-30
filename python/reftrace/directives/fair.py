from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class FairDirective(Directive):
    """The 'fair' directive enables or disables fair scheduling."""
    _value: module_pb2.FairDirective

    @property
    def enabled(self) -> bool:
        """Whether fair scheduling is enabled."""
        return self._value.enabled