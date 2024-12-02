"""Directive types for Nextflow process analysis."""

# Import submodules
from . import accelerator
from . import afterscript
from . import arch
from . import array
from . import beforescript
from . import cache
from . import clusteroptions
from . import conda
from . import container
from . import containeroptions
from . import cpus
from . import debug
from . import disk
from . import dynamic
from . import echo
from . import errorstrategy
from . import executor
from . import ext
from . import fair
from . import label
from . import machinetype
from . import maxerrors
from . import maxforks
from . import maxretries
from . import maxsubmitawait
from . import memory
from . import module
from . import penv
from . import pod
from . import publishdir
from . import queue
from . import resourcelabels
from . import resourcelimits
from . import scratch
from . import shell
from . import spack
from . import stageinmode
from . import stageoutmode
from . import storedir
from . import tag
from . import time
from . import unknown

# Import all public names from submodules
from .accelerator import *
from .afterscript import *
from .arch import *
from .array import *
from .beforescript import *
from .cache import *
from .clusteroptions import *
from .conda import *
from .container import *
from .containeroptions import *
from .cpus import *
from .debug import *
from .disk import *
from .dynamic import *
from .echo import *
from .errorstrategy import *
from .executor import *
from .ext import *
from .fair import *
from .label import *
from .machinetype import *
from .maxerrors import *
from .maxforks import *
from .maxretries import *
from .maxsubmitawait import *
from .memory import *
from .module import *
from .penv import *
from .pod import *
from .publishdir import *
from .queue import *
from .resourcelabels import *
from .resourcelimits import *
from .scratch import *
from .shell import *
from .spack import *
from .stageinmode import *
from .stageoutmode import *
from .storedir import *
from .tag import *
from .time import *
from .unknown import *

# Collect __all__ from submodules
__all__ = (
    accelerator.__all__ +
    afterscript.__all__ +
    arch.__all__ +
    array.__all__ +
    beforescript.__all__ +
    cache.__all__ +
    clusteroptions.__all__ +
    conda.__all__ +
    container.__all__ +
    containeroptions.__all__ +
    cpus.__all__ +
    debug.__all__ +
    disk.__all__ +
    dynamic.__all__ +
    echo.__all__ +
    errorstrategy.__all__ +
    executor.__all__ +
    ext.__all__ +
    fair.__all__ +
    label.__all__ +
    machinetype.__all__ +
    maxerrors.__all__ +
    maxforks.__all__ +
    maxretries.__all__ +
    maxsubmitawait.__all__ +
    memory.__all__ +
    module.__all__ +
    penv.__all__ +
    pod.__all__ +
    publishdir.__all__ +
    queue.__all__ +
    resourcelabels.__all__ +
    resourcelimits.__all__ +
    scratch.__all__ +
    shell.__all__ +
    spack.__all__ +
    stageinmode.__all__ +
    stageoutmode.__all__ +
    storedir.__all__ +
    tag.__all__ +
    time.__all__ +
    unknown.__all__
)
