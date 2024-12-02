from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

__all__ = ['Arch']

@dataclass(frozen=True)
class Arch(Directive):
    """The 'arch' directive specifies architecture requirements."""
    _value: module_pb2.ArchDirective

    @property
    def name(self) -> str:
        """The architecture name."""
        return self._value.name

    @property
    def target(self) -> str:
        """The architecture target."""
        return self._value.target