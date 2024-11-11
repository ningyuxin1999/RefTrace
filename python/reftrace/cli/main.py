# linter.py
import click
import os
import importlib.util
import sys
from dataclasses import dataclass
from typing import List, Callable
from pathlib import Path

from reftrace import Module

registered_rules = []

def register_rule(func: Callable):
    registered_rules.append(func)
    return func  # Return the function unmodified

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

def load_rules(rules_file: str = "rules.py") -> List[Callable]:
    """Load rules from rules.py using the decorator"""
    if not os.path.exists(rules_file):
        click.secho(f"No {rules_file} found in current directory", fg="red")
        sys.exit(1)

    # Prepare the module spec and inject necessary attributes
    spec = importlib.util.spec_from_file_location("rules", rules_file)
    rules_module = importlib.util.module_from_spec(spec)

    # Inject necessary classes and the decorator into the module's namespace
    rules_module.Module = Module
    rules_module.ModuleError = ModuleError
    rules_module.ModuleWarning = ModuleWarning
    rules_module.LintResults = LintResults
    rules_module.register_rule = register_rule

    spec.loader.exec_module(rules_module)

    if not registered_rules:
        click.secho(f"No rules registered in {rules_file}", fg="yellow")

    return registered_rules

def find_nf_files(directory: str) -> List[str]:
    """Recursively find all .nf files in directory"""
    return [str(p) for p in Path(directory).rglob("*.nf")]

def run_lint(directory: str, rules_file: str, debug: bool = False) -> List[LintResults]:
    """Main linting function with optional debug"""
    results = []
    rules = load_rules(rules_file)
    nf_files = find_nf_files(directory)

    for nf_file in nf_files:
        module = Module(nf_file)
        module_results = LintResults(
            module_path=nf_file,
            errors=[],
            warnings=[]
        )

        for rule in rules:
            if debug:
                click.echo(f"Running {rule.__name__} on {nf_file}")

            rule_result = rule(module)
            module_results.errors.extend(rule_result.errors)
            module_results.warnings.extend(rule_result.warnings)

        results.append(module_results)

    return results

@click.command()
@click.argument('rules_file', type=click.Path(exists=True))
@click.argument('directory', type=click.Path(exists=True))
@click.option('--debug', is_flag=True, help="Enable debug output")
def main(rules_file: str, directory: str, debug: bool):
    """Lint .nf files using specified rules file with optional debug mode"""
    results = run_lint(directory, rules_file, debug)

    has_errors = False

    for result in results:
        if result.warnings or result.errors:
            click.echo(f"\nModule: {click.style(result.module_path, fg='cyan')}")

        for warning in result.warnings:
            click.secho(f"  Warning on line {warning.line}: {warning.warning}", fg="yellow")

        for error in result.errors:
            has_errors = True
            click.secho(f"  Error on line {error.line}: {error.error}", fg="red")

    if has_errors:
        sys.exit(1)

if __name__ == "__main__":
    main()
