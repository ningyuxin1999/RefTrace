from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class Disk(Directive):
    """The 'disk' directive specifies disk space requirements."""
    _value: module_pb2.DiskDirective

    @property
    def space(self) -> str:
        """The disk space requirement (e.g., '2 GB', '1 TB')."""
        return self._value.space

__all__ = ['Disk']