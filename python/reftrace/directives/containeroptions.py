from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class ContainerOptionsDirective(Directive):
    """The 'containerOptions' directive specifies additional container runtime options."""
    _value: module_pb2.ContainerOptionsDirective

    @property
    def options(self) -> str:
        """The container runtime options string."""
        return self._value.options