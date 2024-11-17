## Development

Source is at [python/reftrace](python/reftrace). These are Python bindings to the Go shared library via the C FFI. Functionality is exposed on the Go side in [pkg/capi](pkg/capi).

There are tests.

The native Go parser is tested. The C FFI-specific part is tested at [capi_test.go](pkg/capi/capi_test.go). The Python bindings are tested at [test_module.py](python/tests/test_module.py).

## Building

You need Go. You also need Java to run ANTLR.  
You do not need Java to run the built binary.
Download ANTLR from [here](https://www.antlr.org/download/antlr-4.13.1-complete.jar) and put it in the `parser` directory.

```
go generate ./...
go build -o reft
```

Dependencies are vendored. The Go ANTLR target is patched to fix a bug. The ANTLR-generated parser is also patched by [generate_parser.go](cmd/generate_parser.go).

Getting licenses of dependencies:

```
go-licenses save . --save_path="licenses"
```

### Adding a new dependency

```
go get <package>
go mod vendor
git restore vendor/github.com/antlr4-go/antlr/v4/lexer.go  # dependency we patched
```

## Testing

The test data is in a separate repository: [reftrace/reft-testdata](https://github.com/reftrace/reft-testdata).
The Go tests assume you've cloned that to `~/reft-testdata`.

```
go test ./...
```

There are two sets of tests: `reft-go/nf` tests the exposing of the Nextflow DSL to linting rules. `reft-go/parser` tests the underlying Groovy parser.