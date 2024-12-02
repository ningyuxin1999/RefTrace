from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class Cpus(Directive):
    """The 'cpus' directive specifies CPU requirements."""
    _value: module_pb2.CpusDirective

    @property
    def value(self) -> int:
        """Number of CPUs requested."""
        return self._value.value

__all__ = ['Cpus']