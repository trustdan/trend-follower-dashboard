#!/usr/bin/env bash
# build-go-windows.sh
# Cross-compiles Go backend (with embedded UI) to Windows .exe

set -euo pipefail
cd "$(dirname "$0")/.."

PROJECT_ROOT="$(pwd)"

echo "========================================="
echo "Build Windows Binary"
echo "========================================="

# Step 1: Sync UI first
echo "[1/3] Syncing UI to Go..."
./scripts/sync-ui-to-go.sh

# Step 2: Build Windows executable
echo "[2/3] Cross-compiling to Windows..."
cd "${PROJECT_ROOT}/backend"

GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
  go build -trimpath -ldflags "-s -w" \
  -o tf-engine.exe \
  cmd/tf-engine/main.go

if [ ! -f "tf-engine.exe" ]; then
    echo "ERROR: Build failed, tf-engine.exe not found"
    exit 1
fi

SIZE=$(du -h tf-engine.exe | cut -f1)
echo "✓ Windows binary built: tf-engine.exe (${SIZE})"

# Step 3: Verify binary type
echo "[3/3] Verifying binary..."
file tf-engine.exe | grep -q "PE32+" && echo "✓ Verified: PE32+ Windows executable" || echo "⚠ Warning: Binary type unexpected"

echo "========================================="
echo "Build complete!"
echo "========================================="
echo ""
echo "Windows binary: ${PROJECT_ROOT}/backend/tf-engine.exe"
echo ""
echo "To test on Windows:"
echo "  1. Copy tf-engine.exe to Windows machine"
echo "  2. Ensure trading.db is in same directory (or will be created)"
echo "  3. Run: tf-engine.exe server"
echo "  4. Open browser to http://localhost:8080"
echo ""
