from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class ClusterOptionsDirective(Directive):
    """The 'clusterOptions' directive specifies additional cluster submission options."""
    _value: module_pb2.ClusterOptionsDirective

    @property
    def options(self) -> str:
        """The cluster submission options string."""
        return self._value.options