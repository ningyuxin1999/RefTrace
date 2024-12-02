from dataclasses import dataclass
from typing import List
from ..proto import module_pb2
from .base import Directive

__all__ = ['Spack']

@dataclass(frozen=True)
class Spack(Directive):
    """The 'spack' directive specifies Spack package requirements."""
    _value: module_pb2.SpackDirective

    @property
    def possible_values(self) -> List[str]:
        """List of possible Spack package specifications."""
        return list(self._value.possible_values)