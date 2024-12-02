from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class StageOutMode(Directive):
    """The 'stageOutMode' directive specifies how output files should be staged."""
    _value: module_pb2.StageOutModeDirective

    @property
    def mode(self) -> str:
        """The staging mode for output files (e.g., 'copy', 'move', 'rsync')."""
        return self._value.mode

__all__ = ['StageOutMode']