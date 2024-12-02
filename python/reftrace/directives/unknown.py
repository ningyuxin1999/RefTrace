from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class Unknown(Directive):
    """A directive that isn't recognized or supported."""
    _value: module_pb2.UnknownDirective

    @property
    def name(self) -> str:
        """The name of the unknown directive."""
        return self._value.name

    @property
    def value(self) -> str:
        """The raw value of the unknown directive."""
        return self._value.value

__all__ = ['Unknown']