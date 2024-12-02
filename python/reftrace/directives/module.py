from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class Module(Directive):
    """The 'module' directive specifies environment modules to load."""
    _value: module_pb2.ModuleDirective

    @property
    def name(self) -> str:
        """The name of the environment module to load."""
        return self._value.name

__all__ = ['Module']