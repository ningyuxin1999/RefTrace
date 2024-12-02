from dataclasses import dataclass
from typing import List
from ..proto import module_pb2
from .base import Directive

__all__ = ['Conda']
@dataclass(frozen=True)
class Conda(Directive):
    """The 'conda' directive specifies Conda environment requirements."""
    _value: module_pb2.CondaDirective

    @property
    def possible_values(self) -> List[str]:
        """List of possible Conda environment specifications."""
        return list(self._value.possible_values)