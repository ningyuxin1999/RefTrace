# Using From Python

This API is a low level API that allows you to access the Nextflow DSL from Python.
It does parsing of the DSL. It's intended to be consumed by other tools.

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

## API Docs

See [docs/reftrace.md](docs/reftrace.md).

## Building

It's not yet published to PyPI.

You can build it and run the tests by running `make test-python`.

## Development

Source is at [python/reftrace](python/reftrace). These are Python bindings to the Go shared library via the C FFI. Functionality is exposed on the Go side in [pkg/capi](pkg/capi).

There are tests.

The native Go parser is [tested](README.md#testing). The C FFI-specific part is tested at [capi_test.go](pkg/capi/capi_test.go). The Python bindings are tested at [test_module.py](python/tests/test_module.py).