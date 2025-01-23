import click
import json
import sys
from reftrace import parse_modules
from reftrace.graph import make_graph
import networkx as nx
from networkx.algorithms.dag import transitive_closure

@click.group()
def info():
    """Display detailed information about the pipeline."""
    pass

@info.command()
@click.option('--directory', '-d', 
              type=click.Path(exists=True),
              default='.',
              help="Directory containing .nf files (default: current directory)")
@click.option('--pretty/--compact', default=True,
              help="Pretty print the JSON output (default: pretty)")
@click.option('--isolated', is_flag=True, default=False,
              help="Exit with an error if there are any isolated nodes (possibly unused modules)")
def rdeps(directory: str, pretty: bool, isolated: bool):
    """Display reverse dependencies for each module in JSON format."""
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

        G = make_graph(directory, progress_callback)
        if not isinstance(G, nx.DiGraph):
            for error in G:
                if error.likely_rt_bug:
                    click.secho(f"\nInternal error parsing {error.path}:", fg="red", err=True)
                    click.secho(f"  {error.error}", fg="red", err=True)
                    click.secho("This is likely a bug in reftrace. Please file an issue at https://github.com/RefTrace/RefTrace/issues/new", fg="yellow", err=True)
                    sys.exit(1)
                else:
                    click.secho(f"\nFailed to parse {error.path}:", fg="red")
                    click.secho(f"  {error.error}", fg="red")
                    continue
            click.echo("Please fix parsing errors before generating a graph.")
            sys.exit(1)
        
        if len(G.nodes()) == 0:
            click.echo("No Nextflow files found to generate graph.")
            sys.exit(1)
    
    if isolated:
        isolated_nodes = [node for node in G.nodes() 
                         if G.in_degree(node) == 0 and G.out_degree(node) == 0]
        if isolated_nodes:
            click.secho("\nWarning: Found isolated nodes:", fg="yellow", err=True)
            for node in sorted(isolated_nodes):
                click.secho(f"  {node}", fg="yellow", err=True)
            sys.exit(1)
        sys.exit(0)
    
    # Compute transitive closure
    closure = transitive_closure(G)

    # Build the reverse dependencies list
    rdeps_list = []
    for node in sorted(G.nodes()):  # Sort nodes for consistent output
        # Get direct predecessors from original graph
        direct_predecessors = list(G.predecessors(node))
        direct_predecessors.sort()
        
        # Get all predecessors from closure (transitive)
        all_predecessors = list(closure.predecessors(node))
        # Remove direct predecessors to get only transitive ones
        transitive_predecessors = list(set(all_predecessors) - set(direct_predecessors))
        transitive_predecessors.sort()

        rdeps_list.append({
            "path": node,
            "direct_rdeps": direct_predecessors,
            "transitive_rdeps": transitive_predecessors
        })

    # Print JSON output
    indent = 2 if pretty else None
    click.echo(json.dumps(rdeps_list, indent=indent))

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
        modules = module_list_result.results
        resolved_includes = module_list_result.resolved_includes
        unresolved_includes = module_list_result.unresolved_includes
        errors = module_list_result.errors

        for error in errors:
            if error.likely_rt_bug:
                click.secho(f"\nInternal error parsing {error.path}:", fg="red", err=True)
                click.secho(f"  {error.error}", fg="red", err=True)
                click.secho("This is likely a bug in reftrace. Please file an issue at https://github.com/RefTrace/RefTrace/issues/new", fg="yellow", err=True)
                sys.exit(1)
            else:
                click.secho(f"\nFailed to parse {error.path}:", fg="red", err=True)
                click.secho(f"  {error.error}", fg="red", err=True)
                continue

        modules_info = []
        for module in modules:
            modules_info.append(module.to_dict(only_paths=paths))
    
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