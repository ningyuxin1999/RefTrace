from reftrace import Module
from reftrace.linting import ModuleWarning, LintResults, rule


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
