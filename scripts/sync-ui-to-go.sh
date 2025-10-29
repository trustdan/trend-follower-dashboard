#!/usr/bin/env bash
# sync-ui-to-go.sh
# Builds Svelte UI and copies to Go embed directory

set -euo pipefail
cd "$(dirname "$0")/.."

PROJECT_ROOT="$(pwd)"
UI_DIR="${PROJECT_ROOT}/ui"
EMBED_DIR="${PROJECT_ROOT}/backend/internal/webui/dist"

echo "========================================="
echo "Sync UI to Go Embed"
echo "========================================="

# Step 1: Build Svelte UI
echo "[1/3] Building Svelte UI..."
cd "${UI_DIR}"

if [ ! -f "package.json" ]; then
    echo "ERROR: package.json not found in ${UI_DIR}"
    exit 1
fi

npm ci --silent
npm run build

if [ ! -d "build" ]; then
    echo "ERROR: Svelte build failed, build/ directory not found"
    exit 1
fi

echo "✓ Svelte build complete"

# Step 2: Clear old embedded files
echo "[2/3] Clearing old embedded files..."
rm -rf "${EMBED_DIR}"
mkdir -p "${EMBED_DIR}"
echo "✓ Old files cleared"

# Step 3: Copy new build to embed directory
echo "[3/3] Copying build to Go embed directory..."
cp -R build/* "${EMBED_DIR}/"

# Verify files copied
FILE_COUNT=$(find "${EMBED_DIR}" -type f | wc -l)
echo "✓ Copied ${FILE_COUNT} files to ${EMBED_DIR}"

echo "========================================="
echo "Sync complete!"
echo "========================================="
echo ""
echo "Next steps:"
echo "  1. cd backend/"
echo "  2. go run cmd/tf-engine/main.go server"
echo "  3. Open http://localhost:8080 in browser"
echo ""
