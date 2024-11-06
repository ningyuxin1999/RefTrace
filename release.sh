#!/bin/bash

# Set the base name of the binary
BINARY_NAME="reft"

# Check if a version was provided as an argument
if [ -z "$1" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

VERSION=$1

# Create a directory for the releases
mkdir -p releases

# Function to build with version injection
build() {
    local os=$1
    local arch=$2
    echo "Building for $os ($arch) with version ${VERSION}..."
    GOOS=$os GOARCH=$arch go build -ldflags "-X main.version=${VERSION}" -o "releases/${BINARY_NAME}-${os}-${arch}"
}

# Build for Linux (x86_64)
build linux amd64

# Build for macOS (ARM64)
build darwin arm64

echo "Build complete. Binaries are in the 'releases' directory."
echo "Version ${VERSION} has been applied to the binaries."

