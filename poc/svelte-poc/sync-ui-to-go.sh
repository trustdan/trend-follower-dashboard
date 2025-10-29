#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

echo "[1/3] Building Svelte UI..."
npm ci
npm run build

echo "[2/3] Staging UI for Go embed..."
rm -rf go-server/webui/dist
mkdir -p go-server/webui/dist
cp -R build/* go-server/webui/dist/

echo "[3/3] Done. Staged at go-server/webui/dist/"
