from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class PenvDirective(Directive):
    """The 'penv' directive specifies the parallel environment to use."""
    _value: module_pb2.PenvDirective

    @property
    def environment(self) -> str:
        """The parallel environment name."""
        return self._value.environment