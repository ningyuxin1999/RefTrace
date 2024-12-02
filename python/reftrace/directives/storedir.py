from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class StoreDir(Directive):
    """The 'storeDir' directive specifies where to store cached task outputs."""
    _value: module_pb2.StoreDirDirective

    @property
    def path(self) -> str:
        """The directory path where cached outputs should be stored."""
        return self._value.path

__all__ = ['StoreDir']