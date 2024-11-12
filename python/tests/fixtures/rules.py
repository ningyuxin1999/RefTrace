import os
from urllib.parse import urlparse
from reftrace import Module
from reftrace.linting import ModuleError, ModuleWarning, LintResults, rule


# Process label rules

CORRECT_PROCESS_LABELS = [
    "process_single",
    "process_low",
    "process_medium",
    "process_high",
    "process_long",
    "process_high_memory",
]

@rule
def conflicting_labels(module: Module, results: LintResults):
    for process in module.processes:
        good_labels = [label for label in process.labels 
                      if label.value in CORRECT_PROCESS_LABELS]
        if len(good_labels) > 1:
            label_names = [label.value for label in good_labels]
            results.warnings.append(
                ModuleWarning(
                    line=process.line,
                    warning=f"process '{process.name}' has conflicting labels: {label_names}"
                )
            )

@rule
def no_standard_label(module: Module, results: LintResults):
    for process in module.processes:
        good_labels = [label for label in process.labels 
                      if label.value in CORRECT_PROCESS_LABELS]
        if len(good_labels) == 0:
            results.warnings.append(
                ModuleWarning(
                    line=process.line,
                    warning=f"process '{process.name}' has no standard label"
                )
            )

@rule
def non_standard_label(module: Module, results: LintResults):
    for process in module.processes:
        bad_labels = [label for label in process.labels 
                     if label.value not in CORRECT_PROCESS_LABELS]
        if bad_labels:
            label_names = [label.value for label in bad_labels]
            results.warnings.append(
                ModuleWarning(
                    line=process.line,
                    warning=f"process '{process.name}' has non-standard labels: {label_names}"
                )
            )

@rule
def duplicate_labels(module: Module, results: LintResults):
    for process in module.processes:
        label_count = {}
        for label in process.labels:
            label_count[label.value] = label_count.get(label.value, 0) + 1
            
        for label_name, count in label_count.items():
            if count > 1:
                results.warnings.append(
                    ModuleWarning(
                        line=process.line,
                        warning=f"process '{process.name}' has duplicate label '{label_name}' ({count} times)"
                    )
                )

@rule
def no_labels(module: Module, results: LintResults):
    for process in module.processes:
        if not process.labels:
            results.warnings.append(
                ModuleWarning(
                    line=process.line,
                    warning=f"process '{process.name}' has no labels"
                )
            )

@rule
def alphanumerics(module: Module, results: LintResults):
    
    def check_fn(label: str) -> str:
        for char in label:
            if not (char.isalnum() or char == '_'):
                return f"process label '{label}' contains non-alphanumeric characters (only letters, numbers and underscores recommended)"
        return ""
    
    for process in module.processes:
        for label in process.labels:
            if msg := check_fn(label.value):
                results.warnings.append(
                    ModuleWarning(
                        line=label.line,
                        warning=msg
                    )
                )

# Container rules

@rule
def container_with_space(module: Module, results: LintResults):
    for process in module.processes:
        for container in process.containers:
            for name in container.names:
                if " " in name:
                    results.errors.append(
                        ModuleError(
                            line=container.line,
                            error=f"container name '{name}' contains spaces, which is not allowed"
                        )
                    )

@rule
def multiple_containers(module: Module, results: LintResults):
    for process in module.processes:
        for container in process.containers:
            for name in container.names:
                if "biocontainers/" in name and ("https://containers" in name or "https://depot" in name):
                    results.warnings.append(
                        ModuleWarning(
                            line=container.line,
                            warning="Docker and Singularity containers specified on the same line"
                        )
                    )

def is_valid_tag(tag: str) -> bool:
    if not tag:
        return False
    return all(c.isalnum() or c in '-_.' for c in tag)

def get_singularity_tag(container_name: str) -> tuple[str, str | None]:
    try:
        parsed_url = urlparse(container_name)
        last_segment = os.path.basename(parsed_url.path)
        
        if last_segment in (".", "/"):
            return "", "invalid container URL: no path segments"
            
        last_segment = last_segment.removesuffix(".img").removesuffix(".sif")
        
        # Check for colon-separated tag
        if ":" in last_segment:
            tag = last_segment.split(":")[-1]
            if is_valid_tag(tag):
                return tag, None
                
        # Check for _v<digit> pattern
        if "_v" in last_segment:
            idx = last_segment.rindex("_v")
            if len(last_segment) > idx + 2 and last_segment[idx + 2].isdigit():
                tag = last_segment[idx + 1:]
                if is_valid_tag(tag):
                    return tag, None
                    
        return "", f"singularity container '{container_name}' must specify a tag"
    except Exception as e:
        return "", f"invalid container URL '{container_name}': {str(e)}"

def get_docker_tag(container_name: str) -> tuple[str, str | None]:
    if ":" in container_name:
        tag = container_name.split(":")[-1]
        if not is_valid_tag(tag):
            return "", f"invalid docker tag format for container '{container_name}'"
        return tag, None
    return "", f"docker container '{container_name}' must specify a tag"

def docker_or_singularity(container_name: str) -> tuple[str, str | None]:
    if container_name.startswith(("https://", "https://depot")):
        try:
            urlparse(container_name)
            return "singularity", None
        except Exception:
            return "", f"invalid singularity container URL '{container_name}'"
            
    if "/" in container_name or ":" in container_name:
        return "docker", None
        
    return "", f"unknown container type '{container_name}'"

@rule
def must_be_tagged(module: Module, results: LintResults):
    for process in module.processes:
        for container in process.containers:
            for name in container.names:
                container_type, error = docker_or_singularity(name)
                if error:
                    results.errors.append(
                        ModuleError(
                            line=container.line,
                            error=error
                        )
                    )
                    continue
                    
                if container_type == "singularity":
                    _, error = get_singularity_tag(name)
                    if error:
                        results.errors.append(
                            ModuleError(
                                line=container.line,
                                error=error
                            )
                        )
                        
                elif container_type == "docker":
                    _, error = get_docker_tag(name)
                    if error:
                        results.errors.append(
                            ModuleError(
                                line=container.line,
                                error=error
                            )
                        )
                    
                    if name.startswith("quay.io"):
                        results.errors.append(
                            ModuleError(
                                line=container.line,
                                error=f"container '{name}': please use 'organization/container:tag' format instead of full registry URL"
                            )
                        )
