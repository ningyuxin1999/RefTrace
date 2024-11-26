from dataclasses import dataclass
from typing import List, Callable
from functools import wraps

@dataclass
class LintError:
    line: int
    error: str

@dataclass
class LintWarning:
    line: int
    warning: str

@dataclass
class LintResults:
    module_path: str
    errors: List[LintError]
    warnings: List[LintWarning]

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

def configrule(func):
    """Decorator for config file linting rules"""
    @wraps(func)
    def wrapper(config: 'ConfigFile') -> LintResults:
        results = LintResults(
            module_path=config.path,
            errors=[],
            warnings=[]
        )
        func(config, results)
        return results
    wrapper._is_config_rule = True
    return wrapper