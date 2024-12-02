from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class Ext(Directive):
    """The 'ext' directive specifies extension configuration."""
    _value: module_pb2.ExtDirective

    @property
    def version(self) -> str:
        """The extension version."""
        return self._value.version

    @property
    def args(self) -> str:
        """The extension arguments."""
        return self._value.args

__all__ = ['Ext']