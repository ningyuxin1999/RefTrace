from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class BeforeScriptDirective(Directive):
    """The 'beforeScript' directive specifies a script to run before the main process."""
    _value: module_pb2.BeforeScriptDirective

    @property
    def script(self) -> str:
        """The script to execute before the main process."""
        return self._value.script