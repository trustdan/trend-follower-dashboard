#!/bin/bash
# Build script for TF-Engine Native GUI
# Usage: ./build.sh [version]
# Example: ./build.sh v10
# Default: ./build.sh (builds tf-gui.exe)

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "Building tf-gui.exe..."
    go build -o tf-gui.exe *.go
    OUTPUT="tf-gui.exe"
else
    echo "Building tf-gui-$VERSION.exe..."
    go build -o "tf-gui-$VERSION.exe" *.go
    OUTPUT="tf-gui-$VERSION.exe"
fi

if [ $? -eq 0 ]; then
    echo ""
    echo "========================================"
    echo "BUILD SUCCESSFUL!"
    echo "========================================"
    echo "Binary: $OUTPUT"
    echo "Size:"
    ls -lh "$OUTPUT"
    echo "========================================"
else
    echo ""
    echo "========================================"
    echo "BUILD FAILED!"
    echo "========================================"
    echo "Check the error messages above"
    echo "========================================"
    exit 1
fi
