.PHONY: build test clean lib venv

# Main targets
build: lib build-go build-python

# Build the shared library
lib:
	go build -buildmode=c-shared \
		-o python/reftrace/libreftrace.so \
		./pkg/capi

# Build the CLI tool
build-go:
	go build ./cmd/reftrace

# Build the Python package
build-python: lib
	python -m pip install --upgrade pip build
	cd python && python -m build

test: test-go test-python

test-go:
	go test ./...

test-python: venv
	. venv/bin/activate && cd python && python -m pytest

# doesn't rebuild the venv
test-python-quick:
	. venv/bin/activate && cd python && python -m pytest

venv: build
	python -m venv venv
	. venv/bin/activate && pip install python/dist/reftrace-0.4.0-py3-none-any.whl
	. venv/bin/activate && pip install pytest

clean: clean-venv

clean-venv:
	rm -rf venv
	rm -rf build dist *.egg-info
	rm -f python/reftrace/libreftrace.so
	go clean