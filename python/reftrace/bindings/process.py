from dataclasses import dataclass
from typing import Optional, List, ClassVar
from functools import cached_property
from ..proto import module_pb2
from ..directives import Container
from ..directives import Label

@dataclass
class DirectiveValue:
    """Wrapper for directive values that includes the line number while passing through
    all other attribute access to the underlying protobuf object."""
    _value: any
    line: int

    def __getattr__(self, name: str) -> any:
        """Pass through any attribute access to the underlying protobuf object,
        except for 'line' which is handled by the dataclass."""
        return getattr(self._value, name)

@dataclass
class Process:
    """Wrapper for protobuf Process message that provides easier access to directives."""
    _proto: module_pb2.Process

    # All available directive types
    DIRECTIVE_TYPES: ClassVar[List[str]] = [
        'accelerator',
        'after_script',
        'arch',
        'array',
        'before_script',
        'cache',
        'cluster_options',
        'conda',
        'container',
        'container_options',
        'cpus',
        'debug',
        'disk',
        'echo',
        'error_strategy',
        'executor',
        'ext',
        'fair',
        'label',
        'machine_type',
        'max_submit_await',
        'max_errors',
        'max_forks',
        'max_retries',
        'memory',
        'module',
        'penv',
        'pod',
        'publish_dir',
        'queue',
        'resource_labels',
        'resource_limits',
        'scratch',
        'shell',
        'spack',
        'stage_in_mode',
        'stage_out_mode',
        'store_dir',
        'tag',
        'time',
        'dynamic',
        'unknown'
    ]

    @property
    def name(self) -> str:
        """The name of the process."""
        return self._proto.name

    @property
    def line(self) -> int:
        """The line number where this process is defined."""
        return self._proto.line

    def get_directives(self, directive_type: str) -> List[DirectiveValue]:
        """Get all directives of the specified type.
        
        Args:
            directive_type: The type of directive to find (e.g., 'tag', 'cpus', etc.)
            
        Returns:
            List of DirectiveValue objects that wrap the directive values and include line numbers
            
        Raises:
            ValueError: If the directive type is not recognized
        """
        if directive_type not in self.DIRECTIVE_TYPES:
            raise ValueError(f"Unknown directive type: {directive_type}. "
                           f"Available types: {', '.join(self.DIRECTIVE_TYPES)}")
            
        results = []
        for directive in self._proto.directives:
            if directive.WhichOneof('directive') == directive_type:
                value = getattr(directive, directive_type)
                results.append(DirectiveValue(_value=value, line=directive.line))
        return results

    @property
    def directives(self) -> List[module_pb2.Directive]:
        """All directives defined in this process."""
        return list(self._proto.directives)

    # Label directives
    @cached_property
    def labels(self) -> List[Label]:
        """Labels attached to this process."""
        return [
            Label(_value=getattr(d, 'label'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'label'
        ]

    # Container directives
    @cached_property
    def containers(self) -> List[Container]:  # renamed from container for consistency
        """Container specifications for this process."""
        return [
            Container(_value=getattr(d, 'container'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'container'
        ]

    # Convenience methods for single directives
    @property
    def first_label(self) -> Optional[Label]:
        """First label directive or None."""
        return self.labels[0] if self.labels else None

    @property
    def first_container(self) -> Optional[Container]:
        """First container directive or None."""
        return self.containers[0] if self.containers else None