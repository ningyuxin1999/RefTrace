from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class Executor(Directive):
    """The 'executor' directive specifies which executor to use for the process."""
    _value: module_pb2.ExecutorDirective

    @property
    def executor(self) -> str:
        """The executor name (e.g., 'local', 'sge', 'slurm')."""
        return self._value.executor

__all__ = ['Executor']