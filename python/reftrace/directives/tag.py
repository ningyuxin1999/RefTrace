from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class Tag(Directive):
    """The 'tag' directive specifies a tag for the process execution."""
    _value: module_pb2.TagDirective

    @property
    def tag(self) -> str:
        """The tag value."""
        return self._value.tag

__all__ = ['Tag']