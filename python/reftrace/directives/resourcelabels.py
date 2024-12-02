from dataclasses import dataclass
from typing import Dict
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class ResourceLabels(Directive):
    """The 'resourceLabels' directive specifies labels for resources."""
    _value: module_pb2.ResourceLabelsDirective

    @property
    def labels(self) -> Dict[str, str]:
        """The resource labels as a dictionary."""
        return dict(self._value.labels)

__all__ = ['ResourceLabels']