from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class Label(Directive):
    """The 'label' directive adds a label to the process for organization."""
    _value: module_pb2.LabelDirective

    @property
    def value(self) -> str:
        """The label text."""
        return self._value.label

__all__ = ['Label']