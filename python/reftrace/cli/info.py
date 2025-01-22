import click
import json
import sys
from reftrace import parse_modules

@click.group()
def info():
    """Display detailed information about the pipeline."""
    pass

@info.command(name="modules")
@click.option('--directory', '-d', 
              type=click.Path(exists=True),
              default='.',
              help="Directory containing .nf files (default: current directory)")
@click.option('--pretty/--compact', default=True,
              help="Pretty print the JSON output (default: pretty)")
@click.option('--paths', is_flag=True, default=False, help="Only show the module path for each module")
def show_modules_info(directory: str, pretty: bool, paths: bool):
    """Display detailed information about Nextflow modules in JSON format."""
    modules_info = []
    
    with click.progressbar(length=0, label='Parsing Nextflow files', 
                         show_pos=True, 
                         show_percent=True,
                         show_eta=False,
                         width=40,
                         file=sys.stderr) as bar:
        def progress_callback(current: int, total: int):
            if bar.length == 0:
                bar.length = total
            bar.update(current - bar.pos)
                
        module_list_result = parse_modules(directory, progress_callback)
        module_results = module_list_result.results
        resolved_includes = module_list_result.resolved_includes
        unresolved_includes = module_list_result.unresolved_includes

        for module_result in module_results:
            if module_result.error:
                if module_result.error.likely_rt_bug:
                    click.secho(f"\nInternal error parsing {module_result.filepath}:", fg="red", err=True)
                    click.secho(f"  {module_result.error.error}", fg="red", err=True)
                    click.secho("This is likely a bug in reftrace. Please file an issue at https://github.com/RefTrace/RefTrace/issues/new", fg="yellow", err=True)
                    sys.exit(1)
                else:
                    click.secho(f"\nFailed to parse {module_result.filepath}:", fg="red")
                    click.secho(f"  {module_result.error.error}", fg="red")
                    continue

            modules_info.append(module_result.module.to_dict(only_paths=paths))
    
    # Sort modules by path
    modules_info.sort(key=lambda x: x['path'])
    resolved_includes.sort(key=lambda x: x.module_path)
    unresolved_includes.sort(key=lambda x: x.module_path)

    if paths:
        ret = [m['path'] for m in modules_info]
    else:
        ret = {
            "modules": modules_info,
            "resolved_includes": [i.to_dict() for i in resolved_includes],
            "unresolved_includes": [i.to_dict() for i in unresolved_includes]
        }

    # Print JSON output
    indent = 2 if pretty else None
    click.echo(json.dumps(ret, indent=indent))