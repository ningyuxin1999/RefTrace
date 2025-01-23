import click
import os
import importlib.util
import sys
from pathlib import Path
from typing import List, Callable
from importlib.metadata import version
import pkgutil
from reftrace import Module, ConfigFile, parse_modules, ParseError
from reftrace.linting import LintError, LintWarning, LintResults, rule, configrule
from .graph import graph
from .info import info

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
            module_result = Module.from_file(nf_file)
            if isinstance(module_result, ParseError):
                if module_result.likely_rt_bug:
                    # Internal error - should be reported as a bug
                    click.secho(f"Internal error parsing {nf_file}:", fg="red", err=True)
                    click.secho(f"  {module_result.error}", fg="red", err=True)
                    click.secho("This is likely a bug in reftrace. Please file an issue at https://github.com/RefTrace/RefTrace/issues/new", fg="yellow", err=True)
                    sys.exit(1)
                else:
                    # User error - malformed Nextflow file
                    click.secho(f"Failed to parse {nf_file}:", fg="red")
                    click.secho(f"  {module_result.error}", fg="red")
                    continue
            else:
                module = module_result

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
            config_result = ConfigFile.from_file(config_file)
            if config_result.error:
                if config_result.error.likely_rt_bug:
                    # Internal error - should be reported as a bug
                    click.secho(f"Internal error parsing {nf_file}:", fg="red", err=True)
                    click.secho(f"  {module_result.error}", fg="red", err=True)
                    click.secho("This is likely a bug in reftrace. Please file an issue at https://github.com/RefTrace/RefTrace/issues/new", fg="yellow", err=True)
                    sys.exit(1)
                else:
                    # User error - malformed Nextflow file
                    click.secho(f"Failed to parse {nf_file}:", fg="red")
                    click.secho(f"  {module_result.error}", fg="red")
                    continue
            else:
                config = config_result.config_file
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

def run_quicklint(directory: str, rules_file: str, debug: bool = False) -> List[LintResults]:
    """Main linting function with optional debug"""
    results = []
    module_rules, config_rules = load_rules(rules_file)
    
    # Lint Nextflow files using new parse_modules API with progress bar
    with click.progressbar(length=0, label='Linting Nextflow files', 
                         show_pos=True, 
                         show_percent=True,
                         show_eta=False,
                         width=40) as bar:
        def progress_callback(current: int, total: int):
            # First time we get the total, set up the progress bar
            if bar.length == 0:
                bar.length = total
            bar.update(current - bar.pos)
                
        result = parse_modules(directory, progress_callback)
        errors = result.errors
        modules = result.results

        user_errors = False

        for module in modules:
            lint_results = LintResults(
                module_path=module.path,
                errors=[],
                warnings=[]
            )

            for rule in module_rules:
                if debug:
                    click.echo(f"Running {rule.__name__} on {module.path}")

                rule_result = rule(module)
                lint_results.errors.extend(rule_result.errors)
                lint_results.warnings.extend(rule_result.warnings)

            results.append(lint_results)

        for error in errors:
            if error.likely_rt_bug:
                # Internal error - should be reported as a bug
                click.secho(f"\nInternal error parsing {error.error}:", fg="red", err=True)
                click.secho(f"  {error.error}", fg="red", err=True)
                click.secho("This is likely a bug in reftrace. Please file an issue at https://github.com/RefTrace/RefTrace/issues/new", fg="yellow", err=True)
                sys.exit(1)
            else:
                # User error - malformed Nextflow file
                click.secho(f"\nFailed to parse {error.error}:", fg="red")
                click.secho(f"  {error.error}", fg="red")
                user_errors = True
                continue

    if user_errors:
        sys.exit(1)

    # Lint config files
    config_files = find_config_files(directory)
    with click.progressbar(config_files, label='Linting config files', show_pos=True) as files:
        for config_file in files:
            config_result = ConfigFile.from_file(config_file)
            if config_result.error:
                if config_result.error.likely_rt_bug:
                    # Internal error - should be reported as a bug
                    click.secho(f"Internal error parsing {config_file}:", fg="red", err=True)
                    click.secho(f"  {config_result.error}", fg="red", err=True)
                    click.secho("This is likely a bug in reftrace. Please file an issue at https://github.com/RefTrace/RefTrace/issues/new", fg="yellow", err=True)
                    sys.exit(1)
                else:
                    # User error - malformed Nextflow file
                    click.secho(f"Failed to parse {config_file}:", fg="red")
                    click.secho(f"  {config_result.error}", fg="red")
                    user_errors = True
                    continue

            config_results = LintResults(
                module_path=config_file,
                errors=[],
                warnings=[]
            )

            for rule in config_rules:
                if debug:
                    click.echo(f"Running {rule.__name__} on {config_file}")

                rule_result = rule(config_result.config_file)
                config_results.errors.extend(rule_result.errors)
                config_results.warnings.extend(rule_result.warnings)

            results.append(config_results)

    if user_errors:
        sys.exit(1)

    return results

@click.group()
@click.version_option(version=version("reftrace"))
def cli():
    """reftrace - A linting tool for Nextflow files"""
    pass

cli.add_command(graph)
cli.add_command(info)

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
def slowlint(rules_file: str, directory: str, debug: bool, quiet: bool):
    """(Deprecated) Lint Nextflow (.nf) files using custom rules."""
    if not os.path.exists(rules_file):
        click.secho(f"No {rules_file} found. Generating default rules file...", fg="yellow")
        # Read the template from the fixtures
        template = pkgutil.get_data('reftrace', 'fixtures/rules.py').decode('utf-8')
        
        with open(rules_file, 'w') as f:
            f.write(template)
        
        click.secho(f"Created {rules_file} with default rules!", fg="green")

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
        click.secho(f"No {rules_file} found. Generating default rules file...", fg="yellow")
        # Read the template from the fixtures
        template = pkgutil.get_data('reftrace', 'fixtures/rules.py').decode('utf-8')
        
        with open(rules_file, 'w') as f:
            f.write(template)
        
        click.secho(f"Created {rules_file} with default rules!", fg="green")

    # Add initial feedback
    click.secho(f"Loading rules from {rules_file}...", fg="cyan")
    results = run_quicklint(directory, rules_file, debug)

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

if __name__ == "__main__":
    cli()
