from dataclasses import dataclass
from enum import Enum
from ..proto import module_pb2
from .base import Directive

__all__ = ['Container', 'ContainerFormat']

class ContainerFormat(Enum):
    """Container directive format types."""
    SIMPLE = 0
    TERNARY = 1

@dataclass(frozen=True)
class Container(Directive):
    """The 'container' directive specifies the container image to use."""
    _value: module_pb2.ContainerDirective

    @property
    def format(self) -> ContainerFormat:
        """The format of the container specification (simple or ternary)."""
        return ContainerFormat(self._value.format)

    @property
    def name(self) -> str:
        """The container name/image."""
        if self.format == ContainerFormat.SIMPLE:
            return self._value.simple_name
        return self.true_name if self.condition else self.false_name

    @property
    def condition(self) -> str:
        """The condition for ternary format."""
        return self._value.condition

    @property
    def true_name(self) -> str:
        """Container to use when condition is true."""
        return self._value.true_name

    @property
    def false_name(self) -> str:
        """Container to use when condition is false."""
        return self._value.false_name