from dataclasses import dataclass
from typing import List, Callable
from functools import wraps

@dataclass
class ModuleError:
    line: int
    error: str

@dataclass
class ModuleWarning:
    line: int
    warning: str

@dataclass
class LintResults:
    module_path: str
    errors: List[ModuleError]
    warnings: List[ModuleWarning]

def rule(func: Callable):
    @wraps(func)  # Preserve the original function's metadata
    def wrapper(module):
        results = LintResults(
            module_path=module.path,
            errors=[],
            warnings=[]
        )
        func(module, results)
        return results
    return wrapper