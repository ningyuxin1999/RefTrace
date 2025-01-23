import networkx as nx
from typing import Union, List, Optional, Callable
from reftrace import ParseError, ModuleListResult, parse_modules

def make_graph(directory: str, progress_callback: Optional[Callable[[int, int], None]] = None) -> Union[nx.DiGraph, List[ParseError]]:
    module_list_result: ModuleListResult = parse_modules(directory, progress_callback)
    modules = module_list_result.results
    errors = module_list_result.errors

    if errors:
        return errors
    
    module_names = [m.path for m in modules]
    resolved_includes = module_list_result.resolved_includes

    G = nx.DiGraph()
    G.add_nodes_from(module_names)

    for include in resolved_includes:
        for module in include.includes:
            G.add_edge(include.module_path, module)

    return G
