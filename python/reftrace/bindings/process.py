from dataclasses import dataclass
from typing import Optional, List, ClassVar
from functools import cached_property
from ..proto import module_pb2
from ..directives import Container
from ..directives import Label
from ..directives import (
    Accelerator, AfterScript, Arch, Array, BeforeScript, Cache, ClusterOptions,
    Conda, ContainerOptions, Cpus, Debug, Disk, Dynamic, Echo, ErrorStrategy,
    Executor, Ext, Fair, Label, MachineType, MaxErrors, MaxForks, MaxRetries,
    MaxSubmitAwait, Memory, Module, Penv, Pod, PublishDir, Queue, ResourceLabels,
    ResourceLimits, Scratch, Shell, Spack, StageInMode, StageOutMode, StoreDir,
    Tag, Time, Unknown
)

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

    # Directive accessors
    @cached_property
    def accelerators(self) -> List[Accelerator]:
        """Accelerator directives for this process."""
        return [
            Accelerator(_value=getattr(d, 'accelerator'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'accelerator'
        ]

    @cached_property
    def after_scripts(self) -> List[AfterScript]:
        """After script directives for this process."""
        return [
            AfterScript(_value=getattr(d, 'after_script'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'after_script'
        ]

    @cached_property
    def arches(self) -> List[Arch]:
        """Architecture directives for this process."""
        return [
            Arch(_value=getattr(d, 'arch'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'arch'
        ]

    @cached_property
    def arrays(self) -> List[Array]:
        """Array directives for this process."""
        return [
            Array(_value=getattr(d, 'array'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'array'
        ]

    @cached_property
    def before_scripts(self) -> List[BeforeScript]:
        """Before script directives for this process."""
        return [
            BeforeScript(_value=getattr(d, 'before_script'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'before_script'
        ]

    @cached_property
    def caches(self) -> List[Cache]:
        """Cache directives for this process."""
        return [
            Cache(_value=getattr(d, 'cache'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'cache'
        ]

    @cached_property
    def cluster_options(self) -> List[ClusterOptions]:
        """Cluster options directives for this process."""
        return [
            ClusterOptions(_value=getattr(d, 'cluster_options'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'cluster_options'
        ]

    @cached_property
    def condas(self) -> List[Conda]:
        """Conda directives for this process."""
        return [
            Conda(_value=getattr(d, 'conda'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'conda'
        ]
    
    @cached_property
    def containers(self) -> List[Container]:
        """Container specifications for this process."""
        return [
            Container(_value=getattr(d, 'container'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'container'
        ]

    @cached_property
    def container_options(self) -> List[ContainerOptions]:
        """Container options directives for this process."""
        return [
            ContainerOptions(_value=getattr(d, 'container_options'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'container_options'
        ]

    @cached_property
    def cpus(self) -> List[Cpus]:
        """CPU directives for this process."""
        return [
            Cpus(_value=getattr(d, 'cpus'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'cpus'
        ]

    @cached_property
    def debugs(self) -> List[Debug]:
        """Debug directives for this process."""
        return [
            Debug(_value=getattr(d, 'debug'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'debug'
        ]

    @cached_property
    def disks(self) -> List[Disk]:
        """Disk directives for this process."""
        return [
            Disk(_value=getattr(d, 'disk'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'disk'
        ]

    @cached_property
    def dynamics(self) -> List[Dynamic]:
        """Dynamic directives for this process."""
        return [
            Dynamic(_value=getattr(d, 'dynamic'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'dynamic'
        ]

    @cached_property
    def echos(self) -> List[Echo]:
        """Echo directives for this process."""
        return [
            Echo(_value=getattr(d, 'echo'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'echo'
        ]

    @cached_property
    def error_strategies(self) -> List[ErrorStrategy]:
        """Error strategy directives for this process."""
        return [
            ErrorStrategy(_value=getattr(d, 'error_strategy'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'error_strategy'
        ]

    @cached_property
    def executors(self) -> List[Executor]:
        """Executor directives for this process."""
        return [
            Executor(_value=getattr(d, 'executor'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'executor'
        ]

    @cached_property
    def exts(self) -> List[Ext]:
        """Extension directives for this process."""
        return [
            Ext(_value=getattr(d, 'ext'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'ext'
        ]

    @cached_property
    def fairs(self) -> List[Fair]:
        """Fair scheduling directives for this process."""
        return [
            Fair(_value=getattr(d, 'fair'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'fair'
        ]
    
    @cached_property
    def labels(self) -> List[Label]:
        """Labels attached to this process."""
        return [
            Label(_value=getattr(d, 'label'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'label'
        ]

    @cached_property
    def machine_types(self) -> List[MachineType]:
        """Machine type directives for this process."""
        return [
            MachineType(_value=getattr(d, 'machine_type'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'machine_type'
        ]

    @cached_property
    def max_errors(self) -> List[MaxErrors]:
        """Maximum errors directives for this process."""
        return [
            MaxErrors(_value=getattr(d, 'max_errors'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'max_errors'
        ]

    @cached_property
    def max_forks(self) -> List[MaxForks]:
        """Maximum forks directives for this process."""
        return [
            MaxForks(_value=getattr(d, 'max_forks'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'max_forks'
        ]

    @cached_property
    def max_retries(self) -> List[MaxRetries]:
        """Maximum retries directives for this process."""
        return [
            MaxRetries(_value=getattr(d, 'max_retries'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'max_retries'
        ]

    @cached_property
    def max_submit_awaits(self) -> List[MaxSubmitAwait]:
        """Maximum submit await directives for this process."""
        return [
            MaxSubmitAwait(_value=getattr(d, 'max_submit_await'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'max_submit_await'
        ]

    @cached_property
    def memories(self) -> List[Memory]:
        """Memory directives for this process."""
        return [
            Memory(_value=getattr(d, 'memory'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'memory'
        ]

    @cached_property
    def modules(self) -> List[Module]:
        """Module directives for this process."""
        return [
            Module(_value=getattr(d, 'module'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'module'
        ]

    @cached_property
    def penvs(self) -> List[Penv]:
        """Parallel environment directives for this process."""
        return [
            Penv(_value=getattr(d, 'penv'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'penv'
        ]

    @cached_property
    def pods(self) -> List[Pod]:
        """Pod directives for this process."""
        return [
            Pod(_value=getattr(d, 'pod'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'pod'
        ]

    @cached_property
    def publish_dirs(self) -> List[PublishDir]:
        """Publish directory directives for this process."""
        return [
            PublishDir(_value=getattr(d, 'publish_dir'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'publish_dir'
        ]

    @cached_property
    def queues(self) -> List[Queue]:
        """Queue directives for this process."""
        return [
            Queue(_value=getattr(d, 'queue'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'queue'
        ]

    @cached_property
    def resource_labels(self) -> List[ResourceLabels]:
        """Resource labels directives for this process."""
        return [
            ResourceLabels(_value=getattr(d, 'resource_labels'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'resource_labels'
        ]

    @cached_property
    def resource_limits(self) -> List[ResourceLimits]:
        """Resource limits directives for this process."""
        return [
            ResourceLimits(_value=getattr(d, 'resource_limits'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'resource_limits'
        ]

    @cached_property
    def scratches(self) -> List[Scratch]:
        """Scratch directives for this process."""
        return [
            Scratch(_value=getattr(d, 'scratch'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'scratch'
        ]

    @cached_property
    def shells(self) -> List[Shell]:
        """Shell directives for this process."""
        return [
            Shell(_value=getattr(d, 'shell'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'shell'
        ]

    @cached_property
    def spacks(self) -> List[Spack]:
        """Spack directives for this process."""
        return [
            Spack(_value=getattr(d, 'spack'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'spack'
        ]

    @cached_property
    def stage_in_modes(self) -> List[StageInMode]:
        """Stage in mode directives for this process."""
        return [
            StageInMode(_value=getattr(d, 'stage_in_mode'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'stage_in_mode'
        ]

    @cached_property
    def stage_out_modes(self) -> List[StageOutMode]:
        """Stage out mode directives for this process."""
        return [
            StageOutMode(_value=getattr(d, 'stage_out_mode'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'stage_out_mode'
        ]

    @cached_property
    def store_dirs(self) -> List[StoreDir]:
        """Store directory directives for this process."""
        return [
            StoreDir(_value=getattr(d, 'store_dir'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'store_dir'
        ]

    @cached_property
    def tags(self) -> List[Tag]:
        """Tag directives for this process."""
        return [
            Tag(_value=getattr(d, 'tag'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'tag'
        ]

    @cached_property
    def times(self) -> List[Time]:
        """Time directives for this process."""
        return [
            Time(_value=getattr(d, 'time'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'time'
        ]

    @cached_property
    def unknowns(self) -> List[Unknown]:
        """Unknown directives for this process."""
        return [
            Unknown(_value=getattr(d, 'unknown'), line=d.line)
            for d in self._proto.directives
            if d.WhichOneof('directive') == 'unknown'
        ]

    # Convenience methods for single directives
    @property
    def first_accelerator(self) -> Optional[Accelerator]:
        """First accelerator directive or None."""
        return self.accelerators[0] if self.accelerators else None

    @property
    def first_after_script(self) -> Optional[AfterScript]:
        """First after script directive or None."""
        return self.after_scripts[0] if self.after_scripts else None

    @property
    def first_arch(self) -> Optional[Arch]:
        """First architecture directive or None."""
        return self.arches[0] if self.arches else None

    @property
    def first_array(self) -> Optional[Array]:
        """First array directive or None."""
        return self.arrays[0] if self.arrays else None

    @property
    def first_before_script(self) -> Optional[BeforeScript]:
        """First before script directive or None."""
        return self.before_scripts[0] if self.before_scripts else None

    @property
    def first_cache(self) -> Optional[Cache]:
        """First cache directive or None."""
        return self.caches[0] if self.caches else None

    @property
    def first_cluster_options(self) -> Optional[ClusterOptions]:
        """First cluster options directive or None."""
        return self.cluster_options[0] if self.cluster_options else None

    @property
    def first_conda(self) -> Optional[Conda]:
        """First conda directive or None."""
        return self.condas[0] if self.condas else None
    
    @property
    def first_container(self) -> Optional[Container]:
        """First container directive or None."""
        return self.containers[0] if self.containers else None

    @property
    def first_container_options(self) -> Optional[ContainerOptions]:
        """First container options directive or None."""
        return self.container_options[0] if self.container_options else None

    @property
    def first_cpus(self) -> Optional[Cpus]:
        """First CPU directive or None."""
        return self.cpus[0] if self.cpus else None

    @property
    def first_debug(self) -> Optional[Debug]:
        """First debug directive or None."""
        return self.debugs[0] if self.debugs else None

    @property
    def first_disk(self) -> Optional[Disk]:
        """First disk directive or None."""
        return self.disks[0] if self.disks else None

    @property
    def first_dynamic(self) -> Optional[Dynamic]:
        """First dynamic directive or None."""
        return self.dynamics[0] if self.dynamics else None

    @property
    def first_echo(self) -> Optional[Echo]:
        """First echo directive or None."""
        return self.echos[0] if self.echos else None

    @property
    def first_error_strategy(self) -> Optional[ErrorStrategy]:
        """First error strategy directive or None."""
        return self.error_strategies[0] if self.error_strategies else None

    @property
    def first_executor(self) -> Optional[Executor]:
        """First executor directive or None."""
        return self.executors[0] if self.executors else None

    @property
    def first_ext(self) -> Optional[Ext]:
        """First extension directive or None."""
        return self.exts[0] if self.exts else None

    @property
    def first_fair(self) -> Optional[Fair]:
        """First fair scheduling directive or None."""
        return self.fairs[0] if self.fairs else None
    
    @property
    def first_label(self) -> Optional[Label]:
        """First label directive or None."""
        return self.labels[0] if self.labels else None

    @property
    def first_machine_type(self) -> Optional[MachineType]:
        """First machine type directive or None."""
        return self.machine_types[0] if self.machine_types else None

    @property
    def first_max_errors(self) -> Optional[MaxErrors]:
        """First maximum errors directive or None."""
        return self.max_errors[0] if self.max_errors else None

    @property
    def first_max_forks(self) -> Optional[MaxForks]:
        """First maximum forks directive or None."""
        return self.max_forks[0] if self.max_forks else None

    @property
    def first_max_retries(self) -> Optional[MaxRetries]:
        """First maximum retries directive or None."""
        return self.max_retries[0] if self.max_retries else None

    @property
    def first_max_submit_await(self) -> Optional[MaxSubmitAwait]:
        """First maximum submit await directive or None."""
        return self.max_submit_awaits[0] if self.max_submit_awaits else None

    @property
    def first_memory(self) -> Optional[Memory]:
        """First memory directive or None."""
        return self.memories[0] if self.memories else None

    @property
    def first_module(self) -> Optional[Module]:
        """First module directive or None."""
        return self.modules[0] if self.modules else None

    @property
    def first_penv(self) -> Optional[Penv]:
        """First parallel environment directive or None."""
        return self.penvs[0] if self.penvs else None

    @property
    def first_pod(self) -> Optional[Pod]:
        """First pod directive or None."""
        return self.pods[0] if self.pods else None

    @property
    def first_publish_dir(self) -> Optional[PublishDir]:
        """First publish directory directive or None."""
        return self.publish_dirs[0] if self.publish_dirs else None

    @property
    def first_queue(self) -> Optional[Queue]:
        """First queue directive or None."""
        return self.queues[0] if self.queues else None

    @property
    def first_resource_labels(self) -> Optional[ResourceLabels]:
        """First resource labels directive or None."""
        return self.resource_labels[0] if self.resource_labels else None

    @property
    def first_resource_limits(self) -> Optional[ResourceLimits]:
        """First resource limits directive or None."""
        return self.resource_limits[0] if self.resource_limits else None

    @property
    def first_scratch(self) -> Optional[Scratch]:
        """First scratch directive or None."""
        return self.scratches[0] if self.scratches else None

    @property
    def first_shell(self) -> Optional[Shell]:
        """First shell directive or None."""
        return self.shells[0] if self.shells else None

    @property
    def first_spack(self) -> Optional[Spack]:
        """First spack directive or None."""
        return self.spacks[0] if self.spacks else None

    @property
    def first_stage_in_mode(self) -> Optional[StageInMode]:
        """First stage in mode directive or None."""
        return self.stage_in_modes[0] if self.stage_in_modes else None

    @property
    def first_stage_out_mode(self) -> Optional[StageOutMode]:
        """First stage out mode directive or None."""
        return self.stage_out_modes[0] if self.stage_out_modes else None

    @property
    def first_store_dir(self) -> Optional[StoreDir]:
        """First store directory directive or None."""
        return self.store_dirs[0] if self.store_dirs else None

    @property
    def first_tag(self) -> Optional[Tag]:
        """First tag directive or None."""
        return self.tags[0] if self.tags else None

    @property
    def first_time(self) -> Optional[Time]:
        """First time directive or None."""
        return self.times[0] if self.times else None

    @property
    def first_unknown(self) -> Optional[Unknown]:
        """First unknown directive or None."""
        return self.unknowns[0] if self.unknowns else None