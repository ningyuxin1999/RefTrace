import click
import sys
import os
from reftrace import parse_modules
import networkx as nx
import matplotlib.pyplot as plt
from reftrace import Module, ParseError
from reftrace.graph import make_graph
from typing import List

@click.command()
@click.option('--directory', '-d', 
              type=click.Path(exists=True),
              default='.',
              help="Directory containing .nf files (default: current directory)")
def graph(directory: str):
    """Generate a dependency graph for the pipeline."""
    
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

    labels = {node: simplify_path(node) for node in G.nodes()}

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