#!/bin/bash
# Build script for TF-Engine Native GUI

echo "Building TF-Engine Native GUI..."
go build -o tf-gui.exe .

if [ $? -eq 0 ]; then
    echo "✓ Build successful: tf-gui.exe"
    ls -lh tf-gui.exe
else
    echo "✗ Build failed"
    exit 1
fi
