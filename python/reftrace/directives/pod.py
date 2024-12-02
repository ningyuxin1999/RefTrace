from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class Pod(Directive):
    """The 'pod' directive specifies Kubernetes pod configuration."""
    _value: module_pb2.PodDirective

    @property
    def env(self) -> str:
        """The pod environment variable name."""
        return self._value.env

    @property
    def value(self) -> str:
        """The pod environment variable value."""
        return self._value.value

__all__ = ['Pod']