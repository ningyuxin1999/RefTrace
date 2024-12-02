from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class ErrorStrategy(Directive):
    """The 'errorStrategy' directive specifies how to handle process errors."""
    _value: module_pb2.ErrorStrategyDirective

    @property
    def strategy(self) -> str:
        """The error handling strategy (e.g., 'terminate', 'ignore', 'retry')."""
        return self._value.strategy

__all__ = ['ErrorStrategy']