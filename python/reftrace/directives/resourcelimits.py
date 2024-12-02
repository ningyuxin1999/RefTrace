from dataclasses import dataclass
from typing import Dict
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class ResourceLimits(Directive):
    """The 'resourceLimits' directive specifies resource limits."""
    _value: module_pb2.ResourceLimitsDirective

    @property
    def limits(self) -> Dict[str, str]:
        """The resource limits as a dictionary."""
        return dict(self._value.limits)

__all__ = ['ResourceLimits']