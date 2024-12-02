from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class MaxSubmitAwait(Directive):
    """The 'maxSubmitAwait' directive specifies the maximum time to wait for job submission."""
    _value: module_pb2.MaxSubmitAwaitDirective

    @property
    def max_submit_await(self) -> str:
        """The maximum submission wait time (e.g., '1h', '30m')."""
        return self._value.max_submit_await

__all__ = ['MaxSubmitAwait']