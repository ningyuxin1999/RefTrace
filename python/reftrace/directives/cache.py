from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

__all__ = ['Cache']
@dataclass(frozen=True)
class Cache(Directive):
    """The 'cache' directive controls process-level caching behavior."""
    _value: module_pb2.CacheDirective

    @property
    def enabled(self) -> bool:
        """Whether caching is enabled."""
        return self._value.enabled

    @property
    def deep(self) -> bool:
        """Whether deep caching is enabled."""
        return self._value.deep

    @property
    def lenient(self) -> bool:
        """Whether lenient caching is enabled."""
        return self._value.lenient