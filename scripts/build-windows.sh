#!/bin/bash
# Build script for Windows executable
# Builds to /windows folder, then syncs to /release for distribution

set -e

echo "Building Windows executable..."
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o windows/tf-engine.exe ./cmd/tf-engine

echo "Copying to release folder..."
cp windows/tf-engine.exe release/TradingEngine-v3/
cp windows/*.bat release/TradingEngine-v3/ 2>/dev/null || true
cp windows/*.md release/TradingEngine-v3/ 2>/dev/null || true
cp windows/*.vbs release/TradingEngine-v3/ 2>/dev/null || true
cp -r windows/test-data release/TradingEngine-v3/ 2>/dev/null || true

SIZE=$(ls -lh windows/tf-engine.exe | awk '{print $5}')
echo "✅ Build complete! tf-engine.exe ($SIZE)"
echo ""
echo "📁 Active development folder: /windows/"
echo "📦 Release package:           /release/TradingEngine-v3/"
