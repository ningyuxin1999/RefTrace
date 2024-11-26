#!/usr/bin/env python3

import os
import platform
import shutil
import subprocess
import sys
import tarfile
import urllib.request
import zipfile
from pathlib import Path

def get_go_tarball():
    GO_VERSION = "1.22.0"
    system = platform.system().lower()
    machine = platform.machine().lower()
    
    if system == "linux":
        return f"go{GO_VERSION}.linux-amd64.tar.gz"
    elif system == "darwin":
        arch = "arm64" if machine == "arm64" else "amd64"
        return f"go{GO_VERSION}.darwin-{arch}.tar.gz"
    else:
        raise RuntimeError(f"Unsupported system: {system}")

def get_protoc_zip():
    PROTOC_VERSION = "28.3"
    system = platform.system().lower()
    
    if system == "linux":
        return f"protoc-{PROTOC_VERSION}-linux-x86_64.zip"
    elif system == "darwin":
        return f"protoc-{PROTOC_VERSION}-osx-universal_binary.zip"
    else:
        raise RuntimeError(f"Unsupported system: {system}")

def download_file(url, dest):
    print(f"Downloading {url} to {dest}")
    urllib.request.urlretrieve(url, dest)

def install_go():
    if shutil.which("go"):
        print("Go already installed, skipping...")
        return

    install_dir = Path("go")
    tarball = get_go_tarball()
    download_file(f"https://golang.org/dl/{tarball}", tarball)
    
    install_dir.mkdir(exist_ok=True)
    subprocess.run(["tar", "-C", str(install_dir), "--strip-components=1", "-xzf", tarball], check=True)
    os.remove(tarball)
    
    # Add to PATH for the current process
    os.environ["PATH"] = f"{os.path.abspath(install_dir/'bin')}{os.pathsep}{os.environ['PATH']}"
    
    # Verify installation
    subprocess.run(["go", "version"], check=True)

def install_protoc():
    if shutil.which("protoc"):
        print("Protoc already installed, skipping...")
        return

    install_dir = Path("protoc")
    protoc_zip = get_protoc_zip()
    download_file(
        f"https://github.com/protocolbuffers/protobuf/releases/download/v28.3/{protoc_zip}", 
        protoc_zip
    )
    
    install_dir.mkdir(exist_ok=True)
    with zipfile.ZipFile(protoc_zip) as z:
        z.extractall(install_dir)
    
    os.remove(protoc_zip)
    
    # Make the protoc binary executable
    protoc_bin = install_dir / "bin" / "protoc"
    os.chmod(protoc_bin, 0o755)
    
    # Add to PATH for the current process
    os.environ["PATH"] = f"{os.path.abspath(install_dir/'bin')}{os.pathsep}{os.environ['PATH']}"
    
    # Verify installation
    subprocess.run(["protoc", "--version"], check=True)

def install_protoc_go():
    if shutil.which("protoc-gen-go"):
        print("protoc-gen-go already installed, skipping...")
        return
        
    subprocess.run(
        ["go", "install", "google.golang.org/protobuf/cmd/protoc-gen-go@latest"],
        check=True
    )

def main():
    install_go()
    install_protoc()
    install_protoc_go()
    
    # Print Go environment for debugging
    subprocess.run(["go", "env"], check=True)

if __name__ == "__main__":
    main()