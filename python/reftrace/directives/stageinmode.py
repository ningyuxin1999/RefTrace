from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class StageInModeDirective(Directive):
    """The 'stageInMode' directive specifies how input files should be staged."""
    _value: module_pb2.StageInModeDirective

    @property
    def mode(self) -> str:
        """The staging mode for input files (e.g., 'copy', 'link', 'symlink')."""
        return self._value.mode