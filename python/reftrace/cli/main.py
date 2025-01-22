import click
import os
import importlib.util
import sys
from pathlib import Path
from typing import List, Callable
from importlib.metadata import version
import pkgutil
import json
from reftrace import Module, ConfigFile, parse_modules
from reftrace.linting import LintError, LintWarning, LintResults, rule, configrule
import networkx as nx
import matplotlib.pyplot as plt

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
            if module_result.error:
                if module_result.error.likely_rt_bug:
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
                module = module_result.module

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
                
        module_results = parse_modules(directory, progress_callback)

        user_errors = False
        
        for module_result in module_results:
            if module_result.error:
                if module_result.error.likely_rt_bug:
                    # Internal error - should be reported as a bug
                    click.secho(f"\nInternal error parsing {module_result.filepath}:", fg="red", err=True)
                    click.secho(f"  {module_result.error.error}", fg="red", err=True)
                    click.secho("This is likely a bug in reftrace. Please file an issue at https://github.com/RefTrace/RefTrace/issues/new", fg="yellow", err=True)
                    sys.exit(1)
                else:
                    # User error - malformed Nextflow file
                    click.secho(f"\nFailed to parse {module_result.filepath}:", fg="red")
                    click.secho(f"  {module_result.error.error}", fg="red")
                    user_errors = True
                    continue

            lint_results = LintResults(
                module_path=module_result.filepath,
                errors=[],
                warnings=[]
            )

            for rule in module_rules:
                if debug:
                    click.echo(f"Running {rule.__name__} on {module_result.filepath}")

                rule_result = rule(module_result.module)
                lint_results.errors.extend(rule_result.errors)
                lint_results.warnings.extend(rule_result.warnings)

            results.append(lint_results)

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

@cli.command(name="json")
@click.option('--directory', '-d', 
              type=click.Path(exists=True),
              default='.',
              help="Directory containing .nf files (default: current directory)")
@click.option('--pretty/--compact', default=True,
              help="Pretty print the JSON output (default: pretty)")
@click.option('--only-paths', is_flag=True, default=False, help="Only show the module path for each module")
def show_json(directory: str, pretty: bool, only_paths: bool):
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

            modules_info.append(module_result.module.to_dict(only_paths=only_paths))
    
    # Sort modules by path
    modules_info.sort(key=lambda x: x['path'])
    resolved_includes.sort(key=lambda x: x.module_path)
    unresolved_includes.sort(key=lambda x: x.module_path)

    ret = {
        "modules": modules_info,
        "resolved_includes": [i.to_dict() for i in resolved_includes],
        "unresolved_includes": [i.to_dict() for i in unresolved_includes]
    }

    # Print JSON output
    indent = 2 if pretty else None
    click.echo(json.dumps(ret, indent=indent))


@cli.command()
@click.option('--directory', '-d', 
              type=click.Path(exists=True),
              default='.',
              help="Directory containing .nf files (default: current directory)")
@click.option('--inline', is_flag=True, default=False,
              help="Display graph inline in terminal (requires terminal with Kitty image protocol support)")
def graph(directory: str, inline: bool):
    """Generate a dependency graph for the pipeline."""
    modules: List[Module] = []
    
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

        has_errors = False

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
                    has_errors = True
                    continue
            else:
                modules.append(module_result.module)
        if has_errors:
            click.echo("Please fix parsing errors before generating a graph.")
            sys.exit(1)

    module_names = [m.path for m in modules]
    resolved_includes = module_list_result.resolved_includes

    def split_into_lines(text: str, max_length: int = 20) -> str:
        """Split text into multiple lines, each no longer than max_length characters.
        Attempts to split on word boundaries (/ or _) when possible."""
        if len(text) <= max_length:
            return text
            
        # First try splitting on slashes
        if '/' in text:
            parts = text.split('/')
            return '\n'.join(parts)
            
        words = text.replace('_', ' ').split()
        lines = []
        current_line = []
        current_length = 0
        
        for word in words:
            # +1 for the space we'll add between words
            word_length = len(word) + (1 if current_line else 0)
            
            if current_length + word_length <= max_length:
                current_line.append(word)
                current_length += word_length
            else:
                if current_line:
                    lines.append('_'.join(current_line))
                current_line = [word]
                current_length = len(word)
                
        if current_line:
            lines.append('_'.join(current_line))
            
        return '\n'.join(lines)

    def simplify_path(path):
        """Extract meaningful name from path"""
        parts = path.split('/')
        if len(parts) < 2:
            # If there's no parent folder, just return the filename without extension
            return parts[0].replace('.nf', '')
        # Return the parent folder and filename without extension
        simplified = f"{parts[-2]}/{parts[-1].replace('.nf', '')}"
        simplified = simplified.replace('/main', '')
        return split_into_lines(simplified)

    G = nx.DiGraph()
    labels = {path: simplify_path(path) for path in module_names}
    G.add_nodes_from(module_names)

    # Calculate node size based on number of nodes
    num_nodes = len(G.nodes())
    if num_nodes > 50:
        node_size = 2000
    elif num_nodes > 30:
        node_size = 2500
    else:
        node_size = 3000

    # Adjust font size based on number of nodes
    if num_nodes > 50:
        font_size = 6
    elif num_nodes > 30:
        font_size = 7
    else:
        font_size = 8

    for include in resolved_includes:
        for module in include.includes:
            G.add_edge(include.module_path, module)

    def hierarchical_layout(G, root=None, max_nodes_per_row=10):
        """Create a hierarchical layout using networkx"""
        pos = {}
        
        # Find all connected components
        components = list(nx.weakly_connected_components(G))
        
        # Process each component separately
        y_offset = 0
        for component in components:
            # Create subgraph for this component
            subG = G.subgraph(component)
            
            # If no root specified for this component, find node with minimum in-degree
            if root is None or root not in component:
                component_root = min(component, key=lambda n: G.in_degree(n))
            else:
                component_root = root
            
            # Get all layers using BFS for this component
            layers = []
            nodes_seen = set()
            current_layer = {component_root}
            
            while current_layer:
                layers.append(list(current_layer))
                nodes_seen.update(current_layer)
                # Get all neighbors of the current layer that haven't been seen
                next_layer = set()
                for node in current_layer:
                    next_layer.update(n for n in subG.neighbors(node) if n not in nodes_seen)
                current_layer = next_layer
            
            # Add any remaining nodes that weren't reached by BFS
            remaining = component - nodes_seen
            if remaining:
                layers.append(list(remaining))
            
            # Split large layers into sub-layers
            split_layers = []
            for layer in layers:
                if len(layer) > max_nodes_per_row:
                    # Split into multiple rows
                    num_rows = (len(layer) + max_nodes_per_row - 1) // max_nodes_per_row
                    for i in range(num_rows):
                        start_idx = i * max_nodes_per_row
                        end_idx = start_idx + max_nodes_per_row
                        split_layers.append(layer[start_idx:end_idx])
                else:
                    split_layers.append(layer)
            
            # Find the widest layer to scale horizontal spacing
            max_layer_width = max(len(layer) for layer in split_layers)
            
            # Position nodes by layer
            for y, layer in enumerate(split_layers):
                # Calculate x position for each node in the layer
                for x, node in enumerate(layer):
                    # Center each layer and scale x positions based on max width
                    x_pos = (x - len(layer)/2) * (3 * max_layer_width / len(layer))
                    y_pos = -(y + y_offset) * 2  # Multiply by 2 for less vertical space
                    pos[node] = (x_pos, y_pos)
            
            # Update y_offset for next component
            y_offset += len(split_layers)
        
        return pos

    # First get the layout positions
    pos = hierarchical_layout(G, max_nodes_per_row=10)

    # Import color utilities
    from matplotlib.colors import rgb2hex

    def generate_distinct_colors(n):
        """Generate n visually distinct colors using HSV color space"""
        colors = []
        for i in range(n):
            # Use golden ratio to maximize difference between hues
            hue = i * 0.618033988749895
            hue = hue - int(hue)
            # Vary saturation and value for more distinction
            saturation = 0.6 + (i % 3) * 0.2  # Varies between 0.6, 0.8, and 1.0
            value = 0.8 + (i % 2) * 0.2       # Varies between 0.8 and 1.0
            
            # Convert HSV to RGB
            import colorsys
            rgb = colorsys.hsv_to_rgb(hue, saturation, value)
            # Convert to hex
            color = rgb2hex(rgb)
            colors.append(color)
        return colors

    # Generate distinct colors for all nodes instead of just top-level modules
    num_colors = len(G.nodes())
    # color_palette = sns.color_palette("husl", num_colors)  # Generate distinct colors for all nodes
    # node_color_map = {node: rgb2hex(color) for node, color in zip(G.nodes(), color_palette)}
    colors = generate_distinct_colors(num_colors)
    node_color_map = {node: color for node, color in zip(G.nodes(), colors)}
    
    # Create node colors list in the same order as nodes
    node_colors = [node_color_map[node] for node in G.nodes()]

    # Create edge colors based on source node
    edge_colors = [node_color_map[edge[0]] for edge in G.edges()]

    current_dir = os.path.basename(os.path.abspath(directory))

    # Get git commit hash (short version)
    try:
        import subprocess
        git_hash = subprocess.check_output(['git', 'rev-parse', '--short', 'HEAD'], 
                                        cwd=directory,
                                        stderr=subprocess.DEVNULL).decode().strip()
    except:
        git_hash = "unknown"

    # Draw the graph
    plt.style.use('dark_background')
    fig = plt.figure(figsize=(15, 15))
    ax = fig.add_subplot(111)
    ax.set_facecolor('#1a1a1a')  # Dark gray background
    fig.set_facecolor('#1a1a1a')
    
    # Draw edges first (so they're behind nodes)
    nx.draw_networkx_edges(G, pos, 
                          edge_color=edge_colors,
                          arrows=True,
                          arrowsize=10,
                          min_target_margin=25,
                          connectionstyle="arc3,rad=0.2",
                          alpha=0.7)  # Added some transparency
    
    # Draw nodes
    nx.draw_networkx_nodes(G, pos, 
                          node_color=node_colors,
                          node_size=node_size)
    
    # Draw labels with white text and black background box
    nx.draw_networkx_labels(G, pos,
                           labels=labels,
                           font_size=font_size,
                           font_color='white',
                           bbox=dict(facecolor='black', 
                                   edgecolor='none',
                                   alpha=0.7,
                                   pad=2))
    
    # plt.title(f"{current_dir} module dependencies", pad=100, size=16, color="white")

    # Remove legend since each node is unique
    # plt.legend(handles=legend_elements, loc='upper right')

    # Add title with directory name and commit hash
    plt.figtext(0.5, 0.95,
                f"{current_dir}", 
                ha='center',
                color='white',
                size=20)
    
    # Add subtitle below the title
    plt.figtext(0.5, 0.91,
                f"commit {git_hash}\ngenerated with RefTrace",
                ha='center',
                color='white',
                alpha=0.7,
                fontsize=12)
    
    # Remove axes
    plt.axis('off')
    
    plt.tight_layout()
    plt.savefig("graph.png", 
                dpi=300, 
                bbox_inches='tight',
                facecolor='#1a1a1a',  # Ensure dark background is saved
                edgecolor='none')
    click.echo("Graph saved to graph.png")
    plt.close()


if __name__ == "__main__":
    cli()
