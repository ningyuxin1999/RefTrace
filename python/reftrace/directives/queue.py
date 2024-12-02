from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class Queue(Directive):
    """The 'queue' directive specifies which job queue to use."""
    _value: module_pb2.QueueDirective

    @property
    def name(self) -> str:
        """The queue name."""
        return self._value.name

__all__ = ['Queue']