"""Directive types for Nextflow process analysis."""

from .accelerator import AcceleratorDirective as Accelerator
from .afterscript import AfterScriptDirective as AfterScript
from .arch import ArchDirective as Arch
from .array import ArrayDirective as Array
from .beforescript import BeforeScriptDirective as BeforeScript
from .cache import CacheDirective as Cache
from .clusteroptions import ClusterOptionsDirective as ClusterOptions
from .conda import CondaDirective as Conda
from .container import ContainerDirective as Container, ContainerFormat
from .containeroptions import ContainerOptionsDirective as ContainerOptions
from .cpus import CpusDirective as Cpus
from .debug import DebugDirective as Debug
from .disk import DiskDirective as Disk
from .dynamic import DynamicDirective as Dynamic
from .echo import EchoDirective as Echo
from .errorstrategy import ErrorStrategyDirective as ErrorStrategy
from .executor import ExecutorDirective as Executor
from .ext import ExtDirective as Ext
from .fair import FairDirective as Fair
from .label import LabelDirective as Label
from .machinetype import MachineTypeDirective as MachineType
from .maxerrors import MaxErrorsDirective as MaxErrors
from .maxforks import MaxForksDirective as MaxForks
from .maxretries import MaxRetriesDirective as MaxRetries
from .maxsubmitawait import MaxSubmitAwaitDirective as MaxSubmitAwait
from .memory import MemoryDirective as Memory
from .module import ModuleDirective as Module
from .penv import PenvDirective as Penv
from .pod import PodDirective as Pod
from .publishdir import PublishDirDirective as PublishDir
from .queue import QueueDirective as Queue
from .resourcelabels import ResourceLabelsDirective as ResourceLabels
from .resourcelimits import ResourceLimitsDirective as ResourceLimits
from .scratch import ScratchDirective as Scratch
from .shell import ShellDirective as Shell
from .spack import SpackDirective as Spack
from .stageinmode import StageInModeDirective as StageInMode
from .stageoutmode import StageOutModeDirective as StageOutMode
from .storedir import StoreDirDirective as StoreDir
from .tag import TagDirective as Tag
from .time import TimeDirective as Time
from .unknown import UnknownDirective as Unknown


__all__ = [
    # Directive types
    'Accelerator',
    'AfterScript',
    'Arch',
    'Array',
    'BeforeScript',
    'Cache',
    'ClusterOptions',
    'Conda',
    'Container',
    'ContainerFormat',
    'ContainerOptions',
    'Cpus',
    'Debug',
    'Disk',
    'Dynamic',
    'Echo',
    'ErrorStrategy',
    'Executor',
    'Ext',
    'Fair',
    'Label',
    'MachineType',
    'MaxErrors',
    'MaxForks',
    'MaxRetries',
    'MaxSubmitAwait',
    'Memory',
    'Module',
    'Penv',
    'Pod',
    'PublishDir',
    'Queue',
    'ResourceLabels',
    'ResourceLimits',
    'Scratch',
    'Shell',
    'Spack',
    'StageInMode',
    'StageOutMode',
    'StoreDir',
    'Tag',
    'Time',
    'Unknown',
]
