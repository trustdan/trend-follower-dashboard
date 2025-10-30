# Step 26 Implementation Plan: Windows Installer Creation

**Date:** 2025-10-29
**Phase:** 5 - Testing & Packaging
**Step:** 26 of 28
**Status:** Ready to Start (after blocker fixes)

---

## Overview

Create a professional Windows installer that provides one-click deployment of TF-Engine. Based on Windows testing results (see `step26-windows-testing-results.md`), the binary is 95% ready - we need to fix database initialization before building the installer.

**Total Estimated Time:** 10-13 hours (2-3 days)

---

## Prerequisites

### ✅ Completed
- [x] Pure-Go SQLite driver working on Windows
- [x] Binary runs correctly on Windows
- [x] Embedded UI loads successfully
- [x] All core functionality tested
- [x] FINVIZ scraper working

### ❌ Blockers (Must Fix First)
- [ ] **Blocker 1:** Implement `init` command (15 min)
- [ ] **Blocker 2:** Add AppData path support (30 min)
- [ ] Test fixes on Windows (30 min)

**Total Blocker Fix Time:** 1-2 hours

---

## Phase 0: Fix Blockers (REQUIRED)

### Task 0.1: Implement Init Command (15 minutes)

**File:** `backend/cmd/tf-engine/main.go`

**Current state:**
```go
case "init":
    fmt.Println("TODO: Implement init command")
    os.Exit(1)
```

**Required implementation:**

1. Create new function in `main.go` or new file `backend/cmd/tf-engine/init.go`:

```go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "runtime"

    "github.com/yourusername/trading-engine/internal/storage"
)

// InitCommand initializes the database with schema and default settings
func InitCommand() {
    dbPath := flag.String("db", getDefaultDBPath(), "Path to database file")
    flag.Parse()

    // Create directory if it doesn't exist
    dbDir := filepath.Dir(*dbPath)
    if err := os.MkdirAll(dbDir, 0755); err != nil {
        log.Fatalf("Failed to create database directory: %v", err)
    }

    // Open database
    db, err := storage.New(*dbPath)
    if err != nil {
        log.Fatalf("Failed to open database: %v", err)
    }
    defer db.Close()

    // Initialize schema
    if err := db.Initialize(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

    fmt.Printf("Database initialized successfully: %s\n", *dbPath)
}

// getDefaultDBPath returns the default database path based on OS
func getDefaultDBPath() string {
    if runtime.GOOS == "windows" {
        appData := os.Getenv("APPDATA")
        if appData != "" {
            return filepath.Join(appData, "TF-Engine", "trading.db")
        }
    }

    // Fallback for non-Windows or if APPDATA not set
    return "trading.db"
}
```

2. Update `main.go` switch statement:

```go
case "init":
    InitCommand()
```

3. Update `ServerCommand()` to use same path logic:

```go
func ServerCommand() {
    listen := flag.String("listen", "127.0.0.1:8080", "Address to listen on")
    dbPath := flag.String("db", getDefaultDBPath(), "Path to database file")
    flag.Parse()

    // ... rest of function
}
```

**Testing:**
```bash
# Linux/macOS (for cross-compilation testing)
cd backend/
go build -o tf-engine ./cmd/tf-engine
./tf-engine init
./tf-engine init --db /tmp/test.db

# Cross-compile for Windows and test
GOOS=windows GOARCH=amd64 go build -o tf-engine.exe ./cmd/tf-engine
```

```powershell
# Windows (actual testing)
cd backend/
go build -o tf-engine.exe .\cmd\tf-engine
.\tf-engine.exe init
.\tf-engine.exe init --db "C:\temp\test.db"
.\tf-engine.exe init --db "%APPDATA%\TF-Engine\trading.db"
.\tf-engine.exe server --db "%APPDATA%\TF-Engine\trading.db"
```

**Success Criteria:**
- [ ] `init` command creates database file
- [ ] `init` command creates all tables
- [ ] `init` command inserts default settings
- [ ] `init` command works with custom `--db` path
- [ ] `init` command creates directories if needed
- [ ] Default path on Windows is `%APPDATA%\TF-Engine\trading.db`
- [ ] Server can connect to initialized database
- [ ] No errors in API calls after init

### Task 0.2: Rebuild and Test Binary (30 minutes)

1. Rebuild Windows binary with fixes:

```bash
cd /home/kali/trend-follower-dashboard/backend
GOOS=windows GOARCH=amd64 go build -o tf-engine.exe ./cmd/tf-engine
rsync -av tf-engine.exe /mnt/c/Users/Dan/trend-follower-dashboard/backend/
```

2. Test on Windows:

```powershell
cd C:\Users\Dan\trend-follower-dashboard\backend

# Test init command
.\tf-engine.exe init

# Verify database created
Test-Path .\trading.db

# Start server
.\tf-engine.exe server

# Open browser to http://localhost:8080
# Verify no "no such table" errors
# Test all pages (Dashboard, Checklist, Settings, etc.)
```

3. Test AppData location:

```powershell
# Init in AppData
.\tf-engine.exe init --db "%APPDATA%\TF-Engine\trading.db"

# Verify created
Test-Path "$env:APPDATA\TF-Engine\trading.db"

# Start server with AppData db
.\tf-engine.exe server --db "%APPDATA%\TF-Engine\trading.db"

# Test functionality
```

**Success Criteria:**
- [ ] Init creates database successfully
- [ ] Server connects to database
- [ ] All API endpoints return data (not 500 errors)
- [ ] AppData location works correctly
- [ ] No console errors

---

## Phase 1: Application Icon (1-2 hours)

### Task 1.1: Design/Source Icon (1 hour)

**Requirements:**
- Professional appearance
- Represents trend-following or trading
- Multiple sizes: 16x16, 32x32, 48x48, 64x64, 128x128, 256x256
- .ico format for Windows

**Option A: Simple Text Icon**
Use "TF" monogram with upward arrow:
```
  ↗
 TF
```

**Option B: Trend Line Icon**
Upward trending line chart

**Option C: Stock Icons**
- Icons8: https://icons8.com/icons/set/stock-market (free with attribution)
- Flaticon: https://www.flaticon.com/search?word=trading
- Noun Project: https://thenounproject.com/search/?q=trend

**Recommended Quick Option:**
Use Icons8 or generate simple icon with favicon.io

**Tool:** https://favicon.io/favicon-generator/
- Text: "TF"
- Background: Blue (#1E88E5)
- Font: Bold
- Shape: Rounded

### Task 1.2: Convert to .ico (15 minutes)

**Method 1: Online Converter**
- Upload PNG to https://convertio.co/png-ico/
- Download .ico file with all sizes

**Method 2: ImageMagick (if installed)**
```bash
convert icon.png -define icon:auto-resize=256,128,64,48,32,16 app-icon.ico
```

**Save as:** `backend/assets/app-icon.ico`

Create directory:
```bash
mkdir -p backend/assets
# Place app-icon.ico here
```

### Task 1.3: Embed Icon in Binary (30 minutes)

**Install go-winres:**
```bash
go install github.com/tc-hib/go-winres@latest
```

**Create resource config:**

**File:** `backend/cmd/tf-engine/winres.json`
```json
{
  "RT_ICON": {
    "#1": "../../assets/app-icon.ico"
  },
  "RT_GROUP_ICON": {
    "#1": ["#1"]
  },
  "RT_VERSION": {
    "#1": {
      "fixed": {
        "file_version": "1.0.0.0",
        "product_version": "1.0.0.0"
      },
      "info": {
        "0409": {
          "CompanyName": "Your Name",
          "FileDescription": "TF-Engine - Trend Following Trading System",
          "FileVersion": "1.0.0",
          "InternalName": "tf-engine",
          "LegalCopyright": "Copyright © 2025",
          "OriginalFilename": "tf-engine.exe",
          "ProductName": "TF-Engine",
          "ProductVersion": "1.0.0"
        }
      }
    }
  }
}
```

**Generate resource file:**
```bash
cd backend/cmd/tf-engine
go-winres make
```

This creates `rsrc_windows_amd64.syso` (automatically included in builds)

**Rebuild with icon:**
```bash
cd backend
GOOS=windows GOARCH=amd64 go build -o tf-engine.exe ./cmd/tf-engine
```

**Verify icon embedded:**
- Right-click tf-engine.exe on Windows
- Properties → Details tab
- Should show icon, version info, company name

---

## Phase 2: Choose Installer Technology (Decision)

### Option A: WiX Toolset v4 (.msi)
**Pros:**
- Professional .msi installer
- Industry standard
- Windows Installer integration
- Supports upgrades, repair
- Can register uninstaller automatically

**Cons:**
- Requires Windows to build
- More complex XML configuration
- Steeper learning curve

**Best for:** Professional deployment, enterprise users

### Option B: NSIS (.exe)
**Pros:**
- Creates .exe installer
- Can build on Linux (makensis)
- Simpler scripting language
- Lighter weight
- Still professional

**Cons:**
- Not native Windows Installer
- Manual uninstaller registration
- Less enterprise-friendly

**Best for:** Quick deployment, individual users, cross-platform build

### Recommendation: NSIS

**Reasons:**
1. Can build on Linux (current dev environment)
2. Simpler configuration
3. Faster iteration during testing
4. Still professional appearance
5. Adequate for individual trader use case

**Decision:** Use NSIS for v1.0, consider WiX for v2.0+ if enterprise demand

---

## Phase 3: Create NSIS Installer (3-4 hours)

### Task 3.1: Install NSIS (15 minutes)

**On Linux (Kali):**
```bash
sudo apt-get update
sudo apt-get install nsis
makensis -VERSION
```

**Expected output:** `v3.x`

### Task 3.2: Create Installer Script (2 hours)

**File:** `installer/installer.nsi`

```nsis
; TF-Engine Installer Script
; Trend Following Trading Discipline System
; Version 1.0.0

!define APP_NAME "TF-Engine"
!define COMP_NAME "Your Name"
!define VERSION "1.0.0.0"
!define COPYRIGHT "Copyright © 2025 Your Name"
!define DESCRIPTION "Trend-Following Trading Discipline System"
!define INSTALLER_NAME "TF-Engine-Setup-v1.0.0.exe"
!define MAIN_APP_EXE "tf-engine.exe"
!define INSTALL_DIR "$PROGRAMFILES64\${APP_NAME}"

; Includes
!include "MUI2.nsh"
!include "FileFunc.nsh"

; MUI Settings
!define MUI_ABORTWARNING
!define MUI_ICON "..\backend\assets\app-icon.ico"
!define MUI_UNICON "..\backend\assets\app-icon.ico"

; Welcome page
!insertmacro MUI_PAGE_WELCOME

; License page (optional - uncomment if you have license)
; !insertmacro MUI_PAGE_LICENSE "LICENSE.txt"

; Directory page
!insertmacro MUI_PAGE_DIRECTORY

; Instfiles page
!insertmacro MUI_PAGE_INSTFILES

; Finish page
!define MUI_FINISHPAGE_TEXT "TF-Engine has been installed successfully. Click Finish to close this wizard."
!define MUI_FINISHPAGE_RUN "$INSTDIR\${MAIN_APP_EXE}"
!define MUI_FINISHPAGE_RUN_PARAMETERS "server"
!define MUI_FINISHPAGE_RUN_TEXT "Launch TF-Engine"
!insertmacro MUI_PAGE_FINISH

; Uninstaller pages
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

; Language
!insertmacro MUI_LANGUAGE "English"

; Installer attributes
Name "${APP_NAME}"
OutFile "${INSTALLER_NAME}"
InstallDir "${INSTALL_DIR}"
InstallDirRegKey HKLM "Software\${APP_NAME}" "InstallDir"
RequestExecutionLevel admin
ShowInstDetails show
ShowUnInstDetails show

; Version Info
VIProductVersion "${VERSION}"
VIAddVersionKey "ProductName" "${APP_NAME}"
VIAddVersionKey "CompanyName" "${COMP_NAME}"
VIAddVersionKey "FileDescription" "${DESCRIPTION}"
VIAddVersionKey "FileVersion" "${VERSION}"
VIAddVersionKey "LegalCopyright" "${COPYRIGHT}"

; Installer Sections
Section "MainSection" SEC01
    ; Set output path to install directory
    SetOutPath "$INSTDIR"

    ; Copy main executable
    File "..\backend\tf-engine.exe"

    ; Create AppData directory for database
    CreateDirectory "$APPDATA\${APP_NAME}"

    ; Initialize database (silent, no window)
    DetailPrint "Initializing database..."
    nsExec::ExecToLog '"$INSTDIR\${MAIN_APP_EXE}" init --db "$APPDATA\${APP_NAME}\trading.db"'
    Pop $0  ; Return value
    ${If} $0 != 0
        MessageBox MB_OK|MB_ICONEXCLAMATION "Database initialization failed. You may need to run 'tf-engine.exe init' manually."
    ${EndIf}

    ; Create Start Menu folder
    CreateDirectory "$SMPROGRAMS\${APP_NAME}"

    ; Create Start Menu shortcut
    CreateShortCut "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk" \
        "$INSTDIR\${MAIN_APP_EXE}" \
        'server --db "$APPDATA\${APP_NAME}\trading.db"' \
        "$INSTDIR\${MAIN_APP_EXE}" 0 \
        SW_SHOWNORMAL "" \
        "Launch TF-Engine Trading System"

    ; Create Desktop shortcut
    CreateShortCut "$DESKTOP\${APP_NAME}.lnk" \
        "$INSTDIR\${MAIN_APP_EXE}" \
        'server --db "$APPDATA\${APP_NAME}\trading.db"' \
        "$INSTDIR\${MAIN_APP_EXE}" 0 \
        SW_SHOWNORMAL "" \
        "Launch TF-Engine Trading System"

    ; Write installation path to registry
    WriteRegStr HKLM "Software\${APP_NAME}" "InstallDir" "$INSTDIR"
    WriteRegStr HKLM "Software\${APP_NAME}" "Version" "${VERSION}"

    ; Write uninstall information
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "DisplayName" "${APP_NAME}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "DisplayVersion" "${VERSION}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "Publisher" "${COMP_NAME}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "UninstallString" "$INSTDIR\uninstall.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "DisplayIcon" "$INSTDIR\${MAIN_APP_EXE}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "InstallLocation" "$INSTDIR"
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "NoModify" 1
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "NoRepair" 1

    ; Calculate installed size
    ${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
    IntFmt $0 "0x%08X" $0
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "EstimatedSize" "$0"

    ; Create uninstaller
    WriteUninstaller "$INSTDIR\uninstall.exe"
SectionEnd

; Uninstaller Section
Section "Uninstall"
    ; Remove application files
    Delete "$INSTDIR\tf-engine.exe"
    Delete "$INSTDIR\uninstall.exe"

    ; Remove installation directory (if empty)
    RMDir "$INSTDIR"

    ; Remove shortcuts
    Delete "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk"
    RMDir "$SMPROGRAMS\${APP_NAME}"
    Delete "$DESKTOP\${APP_NAME}.lnk"

    ; Ask user if they want to delete database (user data)
    MessageBox MB_YESNO|MB_ICONQUESTION \
        "Do you want to delete your trading database? This will remove all your positions, decisions, and settings. This action cannot be undone." \
        IDNO skip_database

    ; Delete database if user confirmed
    Delete "$APPDATA\${APP_NAME}\trading.db"
    Delete "$APPDATA\${APP_NAME}\trading.db-shm"
    Delete "$APPDATA\${APP_NAME}\trading.db-wal"
    RMDir "$APPDATA\${APP_NAME}"

    skip_database:

    ; Remove registry keys
    DeleteRegKey HKLM "Software\${APP_NAME}"
    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}"
SectionEnd
```

### Task 3.3: Create Build Script (15 minutes)

**File:** `installer/build.sh` (Linux)

```bash
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
if [ -f "assets/app-icon.ico" ]; then
    echo "  ✓ Icon found: backend/assets/app-icon.ico"
else
    echo "  ✗ Warning: Icon not found at backend/assets/app-icon.ico"
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
```

Make executable:
```bash
chmod +x installer/build.sh
```

**File:** `installer/build.bat` (Windows - optional)

```batch
@echo off
echo ==========================================
echo TF-Engine Installer Build Script
echo ==========================================
echo.

REM Step 1: Build Windows binary
echo Step 1: Building Windows binary...
cd ..\backend
go build -o tf-engine.exe .\cmd\tf-engine
echo   [OK] Binary built: backend\tf-engine.exe
echo.

REM Step 2: Build installer
echo Step 2: Building NSIS installer...
cd ..\installer
makensis installer.nsi
echo   [OK] Installer built successfully
echo.

REM Step 3: Calculate checksum
echo Step 3: Calculating SHA256 checksum...
certutil -hashfile TF-Engine-Setup-v1.0.0.exe SHA256 > TF-Engine-Setup-v1.0.0.exe.sha256
type TF-Engine-Setup-v1.0.0.exe.sha256
echo.

echo ==========================================
echo Build Complete!
echo ==========================================
echo Installer: installer\TF-Engine-Setup-v1.0.0.exe
echo Checksum:  installer\TF-Engine-Setup-v1.0.0.exe.sha256
echo.
echo Next steps:
echo 1. Test installer on Windows
echo 2. Verify all features work
echo 3. Test uninstaller
echo.
pause
```

---

## Phase 4: Testing on Windows (2-3 hours)

### Task 4.1: Installation Testing (1 hour)

**Test Environment:**
- Clean Windows 10/11 machine (VM or physical)
- No development tools installed
- Standard user account (will prompt for admin when needed)

**Installation Test Steps:**

1. [ ] Copy installer to Windows test machine
2. [ ] Double-click `TF-Engine-Setup-v1.0.0.exe`
3. [ ] Windows SmartScreen may appear (expected for unsigned):
   - Click "More info"
   - Click "Run anyway"
4. [ ] User Account Control prompt appears → Click "Yes"
5. [ ] Installer wizard opens
6. [ ] Click "Next" on Welcome screen
7. [ ] Choose installation directory (default: `C:\Program Files\TF-Engine`)
8. [ ] Click "Install"
9. [ ] Watch progress bar (should take 10-30 seconds)
10. [ ] Verify "Database initialized successfully" message
11. [ ] Check "Launch TF-Engine" box
12. [ ] Click "Finish"
13. [ ] Browser should open to `http://localhost:8080`
14. [ ] UI should load without errors
15. [ ] Check console (F12) - no JavaScript errors

**Verify Installation:**

```powershell
# Check files installed
Test-Path "C:\Program Files\TF-Engine\tf-engine.exe"

# Check database created
Test-Path "$env:APPDATA\TF-Engine\trading.db"

# Check shortcuts
Test-Path "$env:USERPROFILE\Desktop\TF-Engine.lnk"
Test-Path "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\TF-Engine\TF-Engine.lnk"

# Check registry
Get-ItemProperty "HKLM:\Software\TF-Engine"
Get-ItemProperty "HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\TF-Engine"
```

### Task 4.2: Functionality Testing (1 hour)

**From fresh installation, test complete workflow:**

1. **Settings Page**
   - [ ] Navigate to Settings
   - [ ] Enter equity: $100,000
   - [ ] Enter risk per unit: 0.75%
   - [ ] Enter portfolio cap: 4.0%
   - [ ] Enter bucket cap: 1.5%
   - [ ] Click Save
   - [ ] Verify success message
   - [ ] Refresh page - settings persist

2. **Candidates Scan**
   - [ ] Navigate to Candidates
   - [ ] Click "Scan FINVIZ"
   - [ ] Select preset: TF_BREAKOUT_LONG
   - [ ] Wait for scan (10-15 seconds)
   - [ ] Verify candidates appear (should show 50-100)

3. **Checklist**
   - [ ] Navigate to Checklist
   - [ ] Select a candidate ticker
   - [ ] Check all required boxes
   - [ ] Verify banner turns GREEN
   - [ ] Verify timer starts (2 minutes)

4. **Position Sizing**
   - [ ] Navigate to Size Calculator
   - [ ] Enter ticker: AAPL
   - [ ] Enter entry: $180
   - [ ] Enter ATR: $1.50
   - [ ] Click Calculate
   - [ ] Verify shares calculated
   - [ ] Verify risk dollars shown

5. **Database Persistence**
   - [ ] Close browser
   - [ ] Stop server (Ctrl+C in terminal if running manually)
   - [ ] Relaunch from Desktop shortcut
   - [ ] Verify settings still present
   - [ ] Verify candidates still present

### Task 4.3: Uninstall Testing (30 minutes)

1. **Uninstall Process:**
   - [ ] Open Settings → Apps → Apps & features
   - [ ] Find "TF-Engine"
   - [ ] Click Uninstall
   - [ ] UAC prompt appears → Click Yes
   - [ ] Uninstaller runs
   - [ ] Prompt: "Delete database?" → Click NO (to preserve data)
   - [ ] Uninstall completes

2. **Verify Removal:**
```powershell
# Should NOT exist
Test-Path "C:\Program Files\TF-Engine\tf-engine.exe"
Test-Path "$env:USERPROFILE\Desktop\TF-Engine.lnk"
Test-Path "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\TF-Engine"

# Should STILL exist (user data preserved)
Test-Path "$env:APPDATA\TF-Engine\trading.db"

# Registry should be clean
Get-ItemProperty "HKLM:\Software\TF-Engine"  # Should error
Get-ItemProperty "HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\TF-Engine"  # Should error
```

3. **Reinstall Test:**
   - [ ] Run installer again
   - [ ] Complete installation
   - [ ] Launch application
   - [ ] Verify previous data is preserved (settings, candidates)
   - [ ] Verify application works correctly

4. **Complete Uninstall Test:**
   - [ ] Uninstall again
   - [ ] This time, click YES to delete database
   - [ ] Verify database folder deleted
   - [ ] Verify `%APPDATA%\TF-Engine` removed

---

## Phase 5: Documentation (1 hour)

### Task 5.1: Create Installation Guide (45 minutes)

**File:** `docs/INSTALLATION_GUIDE.md`

See `plans/phase5-step26-windows-installer.md` lines 543-703 for template.

Key sections:
- System requirements
- Download instructions
- Installation steps (with screenshots if possible)
- First launch guide
- Troubleshooting common issues
- Uninstallation instructions
- Upgrade instructions
- Getting help

### Task 5.2: Create Build Documentation (15 minutes)

**File:** `docs/BUILD_INSTALLER.md`

```markdown
# Building the TF-Engine Installer

## Prerequisites

- Go 1.24+ installed
- NSIS 3.x installed (`sudo apt-get install nsis` on Linux)
- go-winres tool: `go install github.com/tc-hib/go-winres@latest`

## Build Steps

### Linux/macOS

```bash
cd installer
./build.sh
```

### Windows

```batch
cd installer
build.bat
```

## Output

- `installer/TF-Engine-Setup-v1.0.0.exe` - The installer
- `installer/TF-Engine-Setup-v1.0.0.exe.sha256` - SHA256 checksum

## Testing

1. Copy installer to Windows test machine
2. Run installer
3. Test all functionality
4. Test uninstaller

## Troubleshooting

**Build fails with "command not found: makensis"**
- Install NSIS: `sudo apt-get install nsis`

**Icon not embedded**
- Ensure `backend/assets/app-icon.ico` exists
- Re-run `go-winres make` in `backend/cmd/tf-engine/`

**Database init fails during install**
- Check that `init` command is implemented
- Test manually: `tf-engine.exe init`
```

---

## Success Criteria (Final Checklist)

### Blocker Fixes
- [ ] `init` command implemented and tested
- [ ] AppData path support added
- [ ] Binary rebuilt and tested on Windows
- [ ] No "no such table" errors

### Icon
- [ ] Icon created (app-icon.ico)
- [ ] Icon embedded in binary
- [ ] Icon visible in Windows Explorer
- [ ] Icon shows in Task Manager

### Installer
- [ ] NSIS installer script created
- [ ] Build script working
- [ ] Installer creates successfully
- [ ] Installer size reasonable (< 50 MB)
- [ ] SHA256 checksum generated

### Installation Testing
- [ ] Installer runs on clean Windows 10
- [ ] Installer runs on clean Windows 11
- [ ] UAC prompt appears (admin rights)
- [ ] Files copied to Program Files
- [ ] Database initialized in AppData
- [ ] Desktop shortcut created
- [ ] Start Menu shortcut created
- [ ] Shortcuts work correctly
- [ ] Application launches from shortcuts
- [ ] Browser opens to localhost:8080
- [ ] UI loads without errors
- [ ] All features functional
- [ ] Settings persist after restart
- [ ] No console errors

### Uninstallation Testing
- [ ] Uninstaller accessible from Apps & features
- [ ] Uninstaller removes executable
- [ ] Uninstaller removes shortcuts
- [ ] Uninstaller removes registry keys
- [ ] Uninstaller preserves database (optional)
- [ ] Prompt to delete database works
- [ ] Complete removal when database deleted
- [ ] No leftover files (except if database preserved)
- [ ] Reinstall works correctly
- [ ] Data preserved after reinstall

### Documentation
- [ ] Installation guide created
- [ ] Build documentation created
- [ ] Troubleshooting section complete
- [ ] Upgrade instructions documented
- [ ] SHA256 checksum documented

---

## Deliverables

1. **Binary:**
   - `backend/tf-engine.exe` (with icon embedded)

2. **Installer:**
   - `installer/TF-Engine-Setup-v1.0.0.exe`
   - `installer/TF-Engine-Setup-v1.0.0.exe.sha256`

3. **Configuration:**
   - `installer/installer.nsi`
   - `installer/build.sh`
   - `installer/build.bat` (optional)

4. **Assets:**
   - `backend/assets/app-icon.ico`
   - `backend/cmd/tf-engine/winres.json`

5. **Documentation:**
   - `docs/INSTALLATION_GUIDE.md`
   - `docs/BUILD_INSTALLER.md`

6. **Testing Report:**
   - `docs/milestones/step26-testing-report.md` (to be created after testing)

---

## Risk Mitigation

### Risk: SmartScreen Warning Scares Users
**Mitigation:**
- Document in installation guide
- Add "This is expected for unsigned software" message
- Consider code signing for v2.0 ($300-500/year)

### Risk: Firewall Blocks Localhost
**Mitigation:**
- Localhost should be exempt from firewall
- If issue occurs, document workaround in troubleshooting

### Risk: Port 8080 Already in Use
**Mitigation:**
- Phase 2: Add `--listen` flag support to change port
- Document how to change port in troubleshooting

### Risk: Database Lock Issues
**Mitigation:**
- SQLite WAL mode already enabled (from storage/db.go)
- Single-instance check not needed (localhost server)

### Risk: Upgrade Fails
**Mitigation:**
- NSIS handles upgrades automatically
- Database in AppData preserved
- Document backup procedure before upgrade

---

## Timeline

| Phase | Task | Time | Cumulative |
|-------|------|------|------------|
| 0 | Fix blockers | 1-2h | 1-2h |
| 1 | Create icon | 1-2h | 2-4h |
| 2 | NSIS installer | 3-4h | 5-8h |
| 3 | Testing | 2-3h | 7-11h |
| 4 | Documentation | 1h | 8-12h |

**Total: 8-12 hours (1-2 days)**

---

## Next Steps

After Step 26 completion:
1. **Step 27:** User Documentation (comprehensive guides, tutorials)
2. **Step 28:** Final Release Preparation (polish, final testing, release notes)

---

**Status:** Ready to Start
**Last Updated:** 2025-10-29
**Dependencies:** Windows testing completed successfully
