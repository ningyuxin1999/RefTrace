# RefTrace

<p align="center">
  <a href="https://github.com/reftrace/reftrace/releases"><img src="https://img.shields.io/github/release/reftrace/reftrace" alt="GitHub release"></a>
  <a href="https://discord.gg/kK7hVKXHQ3"><img src="https://img.shields.io/discord/1299076437571010580?color=blue" alt="Discord chat"></a>
</p>

Lint Nextflow pipelines in Python.

```
pip install reftrace
```

## Features

- Write custom linting rules in Python.
- Avoid parsing the Nextflow DSL.
- Access the Nextflow DSL from Python.

## Example Rule

```python
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
```

## Getting Started

Run in your pipeline directory:

```
reftrace generate  # Create rules.py starter rules
reftrace lint
```

You'll get output like this:

```
Module: /home/andrew/nf-core/rnaseq/modules/nf-core/preseq/lcextrap/main.nf
  Warning on line 1: process 'PRESEQ_LCEXTRAP' has non-standard labels: ['error_ignore']

Module: /home/andrew/nf-core/rnaseq/modules/nf-core/hisat2/build/main.nf
  Warning on line 1: process 'HISAT2_BUILD' has conflicting labels: ['process_high', 'process_high_memory']
```

RefTrace automatically runs the linting rules in `rules.py` on all modules in the pipeline.

## Access Nextflow from Python

```python
from reftrace import Module

module = Module("path/to/nextflow.nf")
```

## Limitations

- Not all parts of the Nextflow DSL are yet exposed. Specifically, only processes are handled. Only directives, process inputs, and process outputs are exposed to linting rules.

- The parser is not perfect. It doesn't seek to handle all of Groovy, but enough to work in practice. Even so, test coverage could be better. If you encounter a
parsing bug, please open an issue.

## Test Data

The files in the `testdata/groovy_core` directory are derived from the Apache Groovy project's test suite and are licensed under the Apache License 2.0. The original source files can be found in the [Apache Groovy repository](https://github.com/apache/groovy/tree/master/src/test-resources/core).

## Acknowledgements

We would like to express our gratitude to the following:

- The [Apache Groovy](https://groovy-lang.org/) project.
- The [ANTLR](https://www.antlr.org/) project, for providing the parser generator used in this tool.
- The Go programming language and its standard library.
- The [Nextflow](https://www.nextflow.io/) project and community, for being so welcoming and helpful.

## License

This project is licensed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file for details.

