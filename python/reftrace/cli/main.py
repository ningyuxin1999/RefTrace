import click
import os
import importlib.util
import sys
from pathlib import Path
from typing import List, Callable
from importlib.metadata import version
import pkgutil

from reftrace import Module, ConfigFile
from reftrace.linting import LintError, LintWarning, LintResults, rule, configrule

def load_rules(rules_file: str = "rules.py") -> tuple[List[Callable], List[Callable]]:
    """Load rules from rules.py using the decorators"""
    if not os.path.exists(rules_file):
        click.secho(f"{rules_file} not found", fg="red")
        sys.exit(1)

    spec = importlib.util.spec_from_file_location("rules", rules_file)
    rules_module = importlib.util.module_from_spec(spec)

    # Inject necessary classes and decorators into the module's namespace
    rules_module.Module = Module
    rules_module.ConfigFile = ConfigFile
    rules_module.LintError = LintError
    rules_module.LintWarning = LintWarning
    rules_module.LintResults = LintResults
    rules_module.rule = rule
    rules_module.configrule = configrule

    spec.loader.exec_module(rules_module)

    # Find all functions decorated with @rule or @configrule
    module_rules = []
    config_rules = []
    for name in dir(rules_module):
        obj = getattr(rules_module, name)
        if callable(obj) and hasattr(obj, '__wrapped__'):
            if hasattr(obj, '_is_config_rule'):
                config_rules.append(obj)
            else:
                module_rules.append(obj)

    if not (module_rules or config_rules):
        click.secho(f"No rules registered in {rules_file}", fg="yellow")

    return module_rules, config_rules

def find_nf_files(directory: str) -> List[str]:
    """Recursively find all .nf files in directory"""
    return [str(p) for p in Path(directory).rglob("*.nf")]

def find_config_files(directory: str) -> List[str]:
    """Recursively find all .config files in directory"""
    return [str(p) for p in Path(directory).rglob("*.config")]

def run_lint(directory: str, rules_file: str, debug: bool = False) -> List[LintResults]:
    """Main linting function with optional debug"""
    results = []
    module_rules, config_rules = load_rules(rules_file)
    
    # Lint Nextflow files
    nf_files = find_nf_files(directory)
    with click.progressbar(nf_files, label='Linting Nextflow files', show_pos=True) as files:
        for nf_file in files:
            module = Module.from_file(nf_file)
            module_results = LintResults(
                module_path=nf_file,
                errors=[],
                warnings=[]
            )

            for rule in module_rules:
                if debug:
                    click.echo(f"Running {rule.__name__} on {nf_file}")

                rule_result = rule(module)
                module_results.errors.extend(rule_result.errors)
                module_results.warnings.extend(rule_result.warnings)

            results.append(module_results)

    # Lint config files
    config_files = find_config_files(directory)
    with click.progressbar(config_files, label='Linting config files', show_pos=True) as files:
        for config_file in files:
            config = ConfigFile.from_file(config_file)
            config_results = LintResults(
                module_path=config_file,
                errors=[],
                warnings=[]
            )

            for rule in config_rules:
                if debug:
                    click.echo(f"Running {rule.__name__} on {config_file}")

                rule_result = rule(config)
                config_results.errors.extend(rule_result.errors)
                config_results.warnings.extend(rule_result.warnings)

            results.append(config_results)

    return results

@click.group()
@click.version_option(version=version("reftrace"))
def cli():
    """reftrace - A linting tool for Nextflow files"""
    pass

@cli.command()
@click.option('--rules', '-r', 'rules_file', 
              type=click.Path(),
              default='rules.py',
              help="Path to rules file (default: rules.py in current directory)")
@click.option('--directory', '-d', 
              type=click.Path(exists=True),
              default='.',
              help="Directory containing .nf files (default: current directory)")
@click.option('--debug', is_flag=True, 
              help="Enable debug output")
@click.option('--quiet', '-q', is_flag=True,
              help="Only show errors, not warnings")
def lint(rules_file: str, directory: str, debug: bool, quiet: bool):
    """Lint Nextflow (.nf) files using custom rules."""
    if not os.path.exists(rules_file):
        click.secho(f"No {rules_file} found in current directory", fg="red")
        click.echo("\nTo get started:")
        click.echo("1. Run 'reftrace generate' to create a template rules file")
        click.echo("2. Edit rules.py to customize the linting rules")
        click.echo("3. Run 'reftrace lint' to check your Nextflow files")
        sys.exit(1)

    # Add initial feedback
    click.secho(f"Loading rules from {rules_file}...", fg="cyan")
    results = run_lint(directory, rules_file, debug)

    has_errors = False
    error_count = 0
    warning_count = 0

    for result in results:
        if result.warnings or result.errors:
            click.echo(f"\nModule: {click.style(result.module_path, fg='cyan')}")

        if not quiet:
            for warning in result.warnings:
                warning_count += 1
                click.secho(f"  Warning on line {warning.line}: {warning.warning}", fg="yellow")

        for error in result.errors:
            error_count += 1
            has_errors = True
            click.secho(f"  Error on line {error.line}: {error.error}", fg="red")

    # Add summary at the end
    click.echo("\nSummary:")
    if error_count:
        click.secho(f"Found {error_count} errors", fg="red")
    if warning_count and not quiet:
        click.secho(f"Found {warning_count} warnings", fg="yellow")
    if not (error_count or warning_count):
        click.secho("No issues found!", fg="green")

    if has_errors:
        sys.exit(1)

@cli.command()
@click.option('--force', '-f', is_flag=True,
              help="Overwrite existing rules.py file")
def generate(force: bool):
    """Generate a template rules.py file with example rules."""
    if os.path.exists('rules.py') and not force:
        click.secho("rules.py already exists. Use --force to overwrite.", fg="red")
        sys.exit(1)
    
    # Read the template from the fixtures
    template = pkgutil.get_data('reftrace', 'fixtures/rules.py').decode('utf-8')
    
    with open('rules.py', 'w') as f:
        f.write(template)
    
    click.secho("Created rules.py with example rules!", fg="green")
    click.echo("\nTo get started:")
    click.echo("1. Edit rules.py to customize the linting rules")
    click.echo("2. Run 'reftrace lint' to check your Nextflow files")

if __name__ == "__main__":
    cli()
