from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class Memory(Directive):
    """The 'memory' directive specifies memory requirements."""
    _value: module_pb2.MemoryDirective

    @property
    def memory_gb(self) -> float:
        """The memory requirement in gigabytes."""
        return self._value.memory_gb

__all__ = ['Memory']