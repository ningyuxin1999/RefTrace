from dataclasses import dataclass
from typing import Optional
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class PublishDir(Directive):
    """The 'publishDir' directive specifies where to publish output files."""
    _value: module_pb2.PublishDirDirective

    @property
    def path(self) -> str:
        """The directory path where files should be published."""
        return self._value.path

    @property
    def params(self) -> str:
        """Additional publishing parameters."""
        return self._value.params

    @property
    def content_type(self) -> Optional[bool]:
        """Whether to determine content type."""
        return self._value.content_type if self._value.HasField('content_type') else None

    @property
    def enabled(self) -> Optional[bool]:
        """Whether publishing is enabled."""
        return self._value.enabled if self._value.HasField('enabled') else None

    @property
    def fail_on_error(self) -> Optional[bool]:
        """Whether to fail on publishing errors."""
        return self._value.fail_on_error if self._value.HasField('fail_on_error') else None

    @property
    def mode(self) -> str:
        """The publishing mode."""
        return self._value.mode

    @property
    def overwrite(self) -> Optional[bool]:
        """Whether to overwrite existing files."""
        return self._value.overwrite if self._value.HasField('overwrite') else None

__all__ = ['PublishDir']