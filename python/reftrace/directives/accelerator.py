from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class AcceleratorDirective(Directive):
    """The 'accelerator' directive specifies GPU requirements."""
    _value: module_pb2.AcceleratorDirective

    @property
    def num_gpus(self) -> int:
        """Number of GPUs requested."""
        return self._value.num_gpus

    @property
    def gpu_type(self) -> str:
        """Type of GPU requested."""
        return self._value.gpu_type