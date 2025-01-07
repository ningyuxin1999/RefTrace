.PHONY: build test clean lib venv docs proto setup-deps

# Main targets
build: proto lib build-python

# Build the shared library
lib:
	@if [ "$$(uname)" = "Darwin" ]; then \
		go build -buildmode=c-shared \
			-o python/reftrace/bindings/libreftrace.dylib \
			./pkg/capi; \
	else \
		go build -buildmode=c-shared \
			-o python/reftrace/bindings/libreftrace.so \
			./pkg/capi; \
	fi

setup-deps:
	python3 scripts/setup_deps.py

# Build the CLI tool
build-go:
	go build ./cmd/reftrace

proto:
	protoc -I=. --go_out=. --python_out=python/reftrace --pyi_out=python/reftrace proto/*
	# fix imports in generated python files
	find python/reftrace/proto -name "*_pb2.py" -exec sed -i.bak 's/from proto/from reftrace.proto/g' {} \;
	find python/reftrace/proto -name "*.bak" -delete
	
# Build the Python package
build-python: lib
	python3 -m venv build-venv
	. build-venv/bin/activate && python3 -m pip install --upgrade pip build
	. build-venv/bin/activate && python3 -m build

dev: lib
	python3 -m venv venv
	. venv/bin/activate && \
		python3 -m pip install --upgrade pip && \
		python3 -m pip install -e ".[dev]"

test: test-go test-python

test-go:
	go test ./...

test-python: dev
	. venv/bin/activate && cd python && python3 -m pytest

# doesn't rebuild the venv
test-python-quick:
	. venv/bin/activate && cd python && python3 -m pytest $(ARGS)

venv: build
	python3 -m venv venv
	. venv/bin/activate && pip install python/dist/reftrace-0.4.0-py3-none-any.whl
	. venv/bin/activate && pip install pytest

clean: clean-venv clean-build-env

clean-venv:
	rm -rf venv
	rm -rf build dist *.egg-info
	rm -f python/reftrace/libreftrace.so
	go clean

clean-build-env:
	rm -rf build-venv

docs-md:
	rm -rf _build docs/source/_autosummary docs/build
	sphinx-build -b markdown docs/source docs/build
	# python3 scripts/transform_docs.py _build/markdown/index.md

docs-html:
	rm -rf _build docs/source/_autosummary docs/build
	sphinx-build -M html docs/source docs/build
