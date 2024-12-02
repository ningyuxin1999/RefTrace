from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

__all__ = ['ClusterOptions']
@dataclass(frozen=True)
class ClusterOptions(Directive):
    """The 'clusterOptions' directive specifies additional cluster submission options."""
    _value: module_pb2.ClusterOptionsDirective

    @property
    def options(self) -> str:
        """The cluster submission options string."""
        return self._value.options