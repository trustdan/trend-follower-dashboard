#!/usr/bin/env bash
set -euo pipefail

ts="$(date +%Y%m%d-%H%M%S)"
out="dist/EXPORT-${ts}.zip"
mkdir -p dist

# Exclude typical build/virtualenv/temp dirs; add or remove as needed.
zip -r "${out}" . \
  -x "*.git*" "dist/*" "target/*" "bin/*" "node_modules/*" \
     "__pycache__/*" ".venv/*" ".tox/*" "publish/*"

echo "Created ${out}"
echo "To copy from Windows: Explorer -> \\\\wsl$\\<YourDistro>\\<path>\\${out}"