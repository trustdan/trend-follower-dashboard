#!/bin/bash
set -e

echo "=========================================="
echo "TF-Engine Installer Build Script"
echo "=========================================="
echo ""

# Step 1: Build Windows binary
echo "Step 1: Building Windows binary..."
cd ../backend
GOOS=windows GOARCH=amd64 go build -o tf-engine.exe ./cmd/tf-engine
echo "  ✓ Binary built: backend/tf-engine.exe"
echo ""

# Step 2: Verify icon embedded
echo "Step 2: Verifying icon..."
if [ -f "assets/trend_following_icon_proper.ico" ]; then
    echo "  ✓ Icon found: backend/assets/trend_following_icon_proper.ico"
else
    echo "  ✗ Warning: Icon not found at backend/assets/trend_following_icon_proper.ico"
fi
echo ""

# Step 3: Build installer
echo "Step 3: Building NSIS installer..."
cd ../installer
makensis installer.nsi
echo "  ✓ Installer built successfully"
echo ""

# Step 4: Calculate checksum
echo "Step 4: Calculating SHA256 checksum..."
sha256sum TF-Engine-Setup-v1.0.0.exe > TF-Engine-Setup-v1.0.0.exe.sha256
cat TF-Engine-Setup-v1.0.0.exe.sha256
echo ""

echo "=========================================="
echo "Build Complete!"
echo "=========================================="
echo "Installer: installer/TF-Engine-Setup-v1.0.0.exe"
echo "Checksum:  installer/TF-Engine-Setup-v1.0.0.exe.sha256"
echo ""
echo "Next steps:"
echo "1. Test installer on Windows VM"
echo "2. Verify all features work"
echo "3. Test uninstaller"
echo ""
