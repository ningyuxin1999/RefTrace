from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class MachineType(Directive):
    """The 'machineType' directive specifies the cloud instance type."""
    _value: module_pb2.MachineTypeDirective

    @property
    def machine_type(self) -> str:
        """The machine/instance type (e.g., 'n1-standard-2', 't2.micro')."""
        return self._value.machine_type

__all__ = ['MachineType']