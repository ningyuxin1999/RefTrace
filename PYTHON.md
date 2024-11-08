# Using From Python

## Example

```python
from reftrace import Module

module = Module("path/to/nextflow.nf")
print(f"File: {module.path}")
print(f"DSL Version: {module.dsl_version}")

for process in module.processes:
    print(f"\nProcess: {process.name}")
    for directive in process.directives:
        if directive.format == "simple":
            print(f"  Container: {directive.simple_name}")
        else:  # ternary
            print(f"  Container: {directive.condition} ? {directive.true_name} : {directive.false_name}")
```

You can access the python module by running `python` from the root of the RefTrace repo.  
You need to first build the shared library.

## Building the shared library

```
go build -buildmode=c-shared -mod=vendor -o reftrace/libreftrace.so reft-go/reftrace
```