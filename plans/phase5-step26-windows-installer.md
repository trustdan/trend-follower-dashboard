# Phase 5 - Step 26: Windows Installer Creation

**TF = Trend Following** - Systematic trading discipline enforcement system

**Phase:** 5 - Testing & Packaging
**Step:** 26 of 28
**Duration:** 2-3 days
**Dependencies:** Step 25 complete (All bugs fixed, app stable)

---

## Objectives

Create a professional Windows installer (.msi or .exe) that provides easy deployment for end users. The installer should include the compiled application binary, create desktop/Start Menu shortcuts, register an uninstaller, and optionally add an application icon. Test the installer on a clean Windows machine to ensure smooth installation and clean removal.

**Purpose:** Provide a one-click installation experience that meets Windows user expectations for professional software.

---

## Success Criteria

- [ ] Windows installer created (.msi preferred, .exe acceptable)
- [ ] Application icon designed and included
- [ ] Installer includes tf-engine.exe binary
- [ ] Desktop shortcut created (optional, user choice)
- [ ] Start Menu shortcut created
- [ ] Uninstaller registered in Windows
- [ ] Database created in correct location (AppData)
- [ ] Installation tested on clean Windows VM
- [ ] Uninstallation tested (complete removal)
- [ ] Installer is code-signed (optional, if certificate available)
- [ ] Installation guide created

---

## Prerequisites

**Completed:**
- All bugs fixed (Step 25)
- Application tested and stable
- Cross-compilation to Windows working

**Required Tools:**
- **Windows machine** for building installer (WiX requires Windows)
- **WiX Toolset v4** (for .msi) OR **NSIS** (for .exe)
- Optional: Code signing certificate (for production release)

**Alternative:** Use NSIS on Linux (cross-platform), but testing still requires Windows.

---

## Implementation Plan

### Part 1: Application Icon (1-2 hours)

#### Task 1.1: Design or Source Icon (1 hour)

**Requirements:**
- Professional appearance
- Represents trend-following or trading
- Multiple sizes: 16x16, 32x32, 48x48, 64x64, 128x128, 256x256
- .ico format for Windows

**Option A: Create Custom Icon**

Use a free icon editor like:
- Inkscape (vector graphics)
- GIMP (raster graphics)
- Online: favicon.io, icons8

**Design ideas:**
- Upward trending arrow
- Candlestick chart
- "TF" monogram
- Wave/tide symbol (representing "trade the tide")

**Option B: Use Stock Icon**

Sources:
- Icons8 (free with attribution)
- Flaticon (free with attribution)
- Noun Project

**Licensing:** Ensure commercial use is allowed or purchase license.

#### Task 1.2: Convert to .ico Format (15 min)

Use ImageMagick or online converter:

```bash
# If you have ImageMagick
convert icon.png -define icon:auto-resize=256,128,64,48,32,16 app-icon.ico
```

Or use online: https://convertio.co/png-ico/

**File:** Save as `assets/app-icon.ico`

#### Task 1.3: Embed Icon in Binary (30 min)

**File:** `backend/cmd/tf-engine/main.go`

For Windows, create a resource file:

**File:** `backend/cmd/tf-engine/rsrc.rc`

```rc
IDI_ICON1 ICON "../../assets/app-icon.ico"
```

Install `go-winres` tool:

```bash
go install github.com/tc-hib/go-winres@latest
```

Generate resource files:

```bash
cd backend/cmd/tf-engine/
go-winres make --in rsrc.rc --out rsrc.syso --arch amd64
```

Rebuild Windows binary (will now include icon):

**Linux/macOS (bash):**
```bash
GOOS=windows GOARCH=amd64 go build -o tf-engine.exe .
```

**Windows (PowerShell):**
```powershell
go build -o tf-engine.exe .
```

---

### Part 2: WiX Installer (Recommended) (4-6 hours)

**Why WiX?**
- Creates professional .msi installers
- Industry standard for Windows
- Integrates with Windows Installer
- Supports uninstall, repair, upgrades

#### Task 2.1: Install WiX Toolset v4 (30 min)

**On Windows machine:**

Download from: https://wixtoolset.org/

Or via Chocolatey:

```powershell
choco install wixtoolset
```

Verify installation:

```powershell
wix --version
```

#### Task 2.2: Create WiX Configuration (2-3 hours)

**File:** `installer/Product.wxs`

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Wix xmlns="http://wixtoolset.org/schemas/v4/wxs">
  <Package
    Name="TF-Engine"
    Manufacturer="Your Name"
    Version="1.0.0"
    UpgradeCode="PUT-GUID-HERE">

    <!-- Generate GUID: https://www.guidgenerator.com/ -->

    <MajorUpgrade
      DowngradeErrorMessage="A newer version is already installed." />

    <!-- Installation directory -->
    <StandardDirectory Id="ProgramFiles6432Folder">
      <Directory Id="INSTALLDIR" Name="TF-Engine">
        <Component Id="MainExecutable" Bitness="always64">
          <File
            Id="TFEngineExe"
            Source="../backend/tf-engine.exe"
            KeyPath="yes">
            <Shortcut
              Id="StartMenuShortcut"
              Directory="ProgramMenuFolder"
              Name="TF-Engine"
              Description="Trend-Following Trading Discipline System"
              WorkingDirectory="INSTALLDIR"
              Icon="AppIcon"
              IconIndex="0"
              Advertise="yes" />
            <Shortcut
              Id="DesktopShortcut"
              Directory="DesktopFolder"
              Name="TF-Engine"
              Description="Trend-Following Trading Discipline System"
              WorkingDirectory="INSTALLDIR"
              Icon="AppIcon"
              IconIndex="0"
              Advertise="yes" />
          </File>
        </Component>
      </Directory>
    </StandardDirectory>

    <!-- Desktop and Start Menu folders -->
    <StandardDirectory Id="DesktopFolder" />
    <StandardDirectory Id="ProgramMenuFolder" />

    <!-- Define icon -->
    <Icon Id="AppIcon" SourceFile="../assets/app-icon.ico" />

    <!-- Feature: Main Application -->
    <Feature Id="MainApplication" Title="TF-Engine" Level="1">
      <ComponentRef Id="MainExecutable" />
    </Feature>

    <!-- UI -->
    <UIRef Id="WixUI_InstallDir" />
    <Property Id="WIXUI_INSTALLDIR" Value="INSTALLDIR" />

    <!-- License agreement (optional) -->
    <WixVariable Id="WixUILicenseRtf" Value="License.rtf" />

  </Package>
</Wix>
```

**Important:** Generate a unique UpgradeCode GUID and use it consistently across versions.

#### Task 2.3: Create License File (Optional) (15 min)

**File:** `installer/License.rtf`

```rtf
{\rtf1\ansi\ansicpg1252\deff0\nouicompat\deflang1033
{\fonttbl{\f0\fnil\fcharset0 Calibri;}}
{\*\generator Your Name}
\viewkind4\uc1
\pard\sa200\sl276\slmult1\f0\fs22\lang9
TF-Engine - Trend Following Trading Discipline System\par
\par
Copyright (c) 2025 Your Name\par
\par
Permission is hereby granted to use this software...\par
\par
(Insert your chosen license here: MIT, Apache, Proprietary, etc.)\par
}
```

Or skip license by removing `WixUILicenseRtf` line.

#### Task 2.4: Build Installer (30 min)

**File (PowerShell):** `installer/build.ps1`

```powershell
Write-Host "Building TF-Engine Installer..."

# Ensure tf-engine.exe is built
cd ..\backend
go build -o tf-engine.exe cmd/tf-engine/main.go
cd ..\installer

# Build MSI
wix build Product.wxs -o TF-Engine-Setup-v1.0.0.msi

Write-Host "Done! Installer: TF-Engine-Setup-v1.0.0.msi"
pause
```

**File (Batch - alternative):** `installer/build.bat`

```batch
@echo off
echo Building TF-Engine Installer...

REM Ensure tf-engine.exe is built
cd ..\backend
go build -o tf-engine.exe cmd/tf-engine/main.go
cd ..\installer

REM Build MSI
wix build Product.wxs -o TF-Engine-Setup-v1.0.0.msi

echo Done! Installer: TF-Engine-Setup-v1.0.0.msi
pause
```

Run:

```powershell
cd installer
.\build.bat
```

**Output:** `TF-Engine-Setup-v1.0.0.msi`

---

### Part 3: NSIS Installer (Alternative) (3-4 hours)

**Why NSIS?**
- Creates .exe installers
- Can be built on Linux (cross-platform)
- Lighter weight than WiX
- Still professional

#### Task 3.1: Install NSIS (15 min)

**On Windows:**
Download from: https://nsis.sourceforge.io/

**On Linux:**

```bash
sudo apt-get install nsis
```

#### Task 3.2: Create NSIS Script (2 hours)

**File:** `installer/installer.nsi`

```nsis
; TF-Engine Installer Script
!define APP_NAME "TF-Engine"
!define COMP_NAME "Your Name"
!define VERSION "1.0.0.0"
!define COPYRIGHT "Copyright © 2025 Your Name"
!define DESCRIPTION "Trend-Following Trading Discipline System"
!define INSTALLER_NAME "TF-Engine-Setup-v1.0.0.exe"
!define MAIN_APP_EXE "tf-engine.exe"
!define INSTALL_DIR "$PROGRAMFILES64\${APP_NAME}"

; MUI Settings
!include "MUI2.nsh"

!define MUI_ABORTWARNING
!define MUI_ICON "../assets/app-icon.ico"
!define MUI_UNICON "../assets/app-icon.ico"

; Welcome page
!insertmacro MUI_PAGE_WELCOME

; License page (optional)
; !insertmacro MUI_PAGE_LICENSE "License.txt"

; Directory page
!insertmacro MUI_PAGE_DIRECTORY

; Install page
!insertmacro MUI_PAGE_INSTFILES

; Finish page
!define MUI_FINISHPAGE_RUN "$INSTDIR\${MAIN_APP_EXE}"
!insertmacro MUI_PAGE_FINISH

; Uninstaller pages
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

; Language
!insertmacro MUI_LANGUAGE "English"

; General
Name "${APP_NAME}"
OutFile "${INSTALLER_NAME}"
InstallDir "${INSTALL_DIR}"
InstallDirRegKey HKLM "Software\${APP_NAME}" "InstallDir"
RequestExecutionLevel admin

; Version Info
VIProductVersion "${VERSION}"
VIAddVersionKey "ProductName" "${APP_NAME}"
VIAddVersionKey "CompanyName" "${COMP_NAME}"
VIAddVersionKey "FileDescription" "${DESCRIPTION}"
VIAddVersionKey "FileVersion" "${VERSION}"

Section "MainSection" SEC01
  ; Files
  SetOutPath "$INSTDIR"
  File "..\backend\tf-engine.exe"

  ; Shortcuts
  CreateDirectory "$SMPROGRAMS\${APP_NAME}"
  CreateShortCut "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk" "$INSTDIR\${MAIN_APP_EXE}" "" "$INSTDIR\${MAIN_APP_EXE}" 0
  CreateShortCut "$DESKTOP\${APP_NAME}.lnk" "$INSTDIR\${MAIN_APP_EXE}" "" "$INSTDIR\${MAIN_APP_EXE}" 0

  ; Registry
  WriteRegStr HKLM "Software\${APP_NAME}" "InstallDir" "$INSTDIR"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "DisplayName" "${APP_NAME}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "UninstallString" "$INSTDIR\uninstall.exe"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "DisplayIcon" "$INSTDIR\${MAIN_APP_EXE}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "Publisher" "${COMP_NAME}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "DisplayVersion" "${VERSION}"

  ; Uninstaller
  WriteUninstaller "$INSTDIR\uninstall.exe"
SectionEnd

Section "Uninstall"
  ; Remove files
  Delete "$INSTDIR\tf-engine.exe"
  Delete "$INSTDIR\uninstall.exe"
  RMDir "$INSTDIR"

  ; Remove shortcuts
  Delete "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk"
  RMDir "$SMPROGRAMS\${APP_NAME}"
  Delete "$DESKTOP\${APP_NAME}.lnk"

  ; Remove registry
  DeleteRegKey HKLM "Software\${APP_NAME}"
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}"
SectionEnd
```

#### Task 3.3: Build Installer (15 min)

**Windows:**

```batch
cd installer
makensis installer.nsi
```

**Linux:**

```bash
cd installer
makensis installer.nsi
```

**Output:** `TF-Engine-Setup-v1.0.0.exe`

---

### Part 4: Code Signing (Optional) (1-2 hours)

**Why Sign?**
- Windows SmartScreen won't block signed apps
- Users trust signed software
- Professional appearance

**Requirements:**
- Code signing certificate (from DigiCert, Sectigo, etc.)
- Cost: ~$100-$500/year

**If you have a certificate:**

**Using SignTool (Windows):**

```powershell
# Sign the installer
signtool sign /f "YourCertificate.pfx" /p "password" /tr http://timestamp.digicert.com /td sha256 /fd sha256 "TF-Engine-Setup-v1.0.0.msi"
```

**Verify signature:**

```powershell
signtool verify /pa "TF-Engine-Setup-v1.0.0.msi"
```

**If no certificate:**
- Skip signing for internal use or beta releases
- Document in installation guide that Windows SmartScreen may warn
- Consider self-signing for development (users must add cert to trusted store)

---

### Part 5: Testing on Clean Windows (2-3 hours)

#### Task 5.1: Prepare Test Environment (30 min)

**Option A: Windows VM**
- Use VirtualBox or VMware
- Install Windows 10/11 fresh
- No dev tools installed

**Option B: Physical Windows Machine**
- Use a clean test machine
- Create a test user account

#### Task 5.2: Install and Test (1 hour)

**Installation Test:**

1. [ ] Copy installer to Windows machine
2. [ ] Double-click installer
3. [ ] Follow installation wizard
4. [ ] Choose installation directory (default or custom)
5. [ ] Complete installation
6. [ ] Verify shortcuts created (Desktop, Start Menu)
7. [ ] Launch app from Desktop shortcut
8. [ ] Verify browser opens to localhost:8080
9. [ ] Verify UI loads correctly
10. [ ] Run through complete workflow (scan → decision)
11. [ ] Check database location: `%APPDATA%\TF-Engine\trading.db`
12. [ ] Close app

**Expected:** Smooth installation, no errors, app works identically to development environment.

#### Task 5.3: Uninstall and Test (30 min)

**Uninstallation Test:**

1. [ ] Open "Add or Remove Programs"
2. [ ] Find "TF-Engine"
3. [ ] Click "Uninstall"
4. [ ] Follow uninstaller
5. [ ] Verify shortcuts removed (Desktop, Start Menu)
6. [ ] Verify installation directory removed
7. [ ] Check registry: No leftover keys
8. [ ] Check AppData: Database file remains (user data, expected)

**Expected:** Clean removal, no leftover files except user data.

#### Task 5.4: Upgrade Test (Optional) (30 min)

If creating a v1.0.1 installer:

1. [ ] Install v1.0.0
2. [ ] Create some data (add positions)
3. [ ] Install v1.0.1 over it
4. [ ] Verify data persists
5. [ ] Verify new version runs

**Expected:** Upgrade preserves user data.

---

### Part 6: Installation Guide (1 hour)

**File:** `docs/INSTALLATION_GUIDE.md`

```markdown
# TF-Engine Installation Guide

**Version:** 1.0.0
**Platform:** Windows 10/11 (64-bit)

---

## System Requirements

- **OS:** Windows 10 or Windows 11 (64-bit)
- **RAM:** 4 GB minimum, 8 GB recommended
- **Disk Space:** 500 MB for application + database
- **Browser:** Chrome, Edge, or Firefox (latest version)
- **Internet:** Required for FINVIZ scanning

---

## Installation Steps

### Step 1: Download Installer

Download `TF-Engine-Setup-v1.0.0.msi` (or `.exe`) from:
- [Your distribution method: website, GitHub releases, etc.]

**File Size:** ~50 MB
**SHA256 Checksum:** [Insert checksum here]

### Step 2: Run Installer

1. Double-click `TF-Engine-Setup-v1.0.0.msi`
2. Windows may show "User Account Control" prompt → Click **Yes**
3. If Windows SmartScreen appears (unsigned installer):
   - Click "More info"
   - Click "Run anyway"

### Step 3: Follow Installation Wizard

1. **Welcome Screen:** Click **Next**
2. **License Agreement:** Accept terms, click **Next** (if applicable)
3. **Installation Directory:**
   - Default: `C:\Program Files\TF-Engine\`
   - Or choose custom location
   - Click **Next**
4. **Ready to Install:** Click **Install**
5. **Installation Progress:** Wait for files to copy (~30 seconds)
6. **Completion:**
   - Check "Launch TF-Engine" (optional)
   - Click **Finish**

### Step 4: First Launch

1. Double-click Desktop shortcut "TF-Engine"
   - OR: Start Menu → TF-Engine
2. Your default browser opens to `http://localhost:8080`
3. The TF-Engine UI loads (first load may take 5-10 seconds)

### Step 5: Initialize Database

On first run:
1. Navigate to **Settings**
2. Enter your account details:
   - Equity: Your trading account size
   - Risk % per unit: Typically 0.75% - 1.0%
   - Portfolio heat cap: 4.0% (default)
   - Bucket heat cap: 1.5% (default)
3. Click **Save Settings**

**Congratulations!** TF-Engine is now ready to use.

---

## Troubleshooting

### Installer won't run

**Problem:** "Windows protected your PC" message
**Solution:**
1. Click "More info"
2. Click "Run anyway"
3. This is expected for unsigned applications

**Problem:** "This app can't run on your PC"
**Solution:** Ensure you have 64-bit Windows 10/11

---

### App won't start

**Problem:** Browser doesn't open
**Solution:**
1. Manually open browser
2. Navigate to `http://localhost:8080`

**Problem:** "Cannot connect" in browser
**Solution:**
1. Check Task Manager → tf-engine.exe should be running
2. If not, try running as Administrator
3. Check firewall settings (allow localhost connections)

---

### Database issues

**Problem:** "Database is locked"
**Solution:**
1. Close any SQLite browser tools
2. Restart TF-Engine
3. Only one instance can run at a time

**Location:** Database is stored at:
```
C:\Users\[YourUsername]\AppData\Roaming\TF-Engine\trading.db
```

**Backup:** Copy this file to backup your data.

---

## Uninstallation

### Method 1: Windows Settings

1. Open **Settings** → **Apps** → **Apps & features**
2. Find "TF-Engine"
3. Click → **Uninstall**
4. Confirm uninstallation
5. Follow uninstaller prompts

### Method 2: Control Panel

1. Open **Control Panel** → **Programs and Features**
2. Find "TF-Engine"
3. Right-click → **Uninstall**

**Note:** Uninstaller removes the application but preserves your database file (user data). To completely remove all data, manually delete:
```
C:\Users\[YourUsername]\AppData\Roaming\TF-Engine\
```

---

## Upgrading

To upgrade to a newer version:

1. Download new installer (e.g., v1.1.0)
2. Run new installer
3. It will automatically upgrade over the old version
4. Your database and settings will be preserved

**Important:** Always backup your database before upgrading.

---

## Getting Help

- Documentation: [Link to docs]
- Issues: [Link to issue tracker]
- Email: [Your support email]
```

---

## Testing Checklist

### Installer Creation
- [ ] Application icon created and embedded
- [ ] WiX or NSIS installer configured
- [ ] License file included (if applicable)
- [ ] Build script working
- [ ] Installer file generated successfully

### Installation Testing
- [ ] Tested on clean Windows 10
- [ ] Tested on clean Windows 11
- [ ] Desktop shortcut created
- [ ] Start Menu shortcut created
- [ ] App launches from shortcuts
- [ ] Browser opens to correct URL
- [ ] All features work on Windows
- [ ] Database created in AppData

### Uninstallation Testing
- [ ] Uninstaller removes shortcuts
- [ ] Uninstaller removes program files
- [ ] Uninstaller removes registry entries
- [ ] Database preserved (user data)
- [ ] No leftover files (except AppData)

### Code Signing (if applicable)
- [ ] Installer signed with certificate
- [ ] Signature verified with SignTool
- [ ] SmartScreen doesn't block signed installer

### Documentation
- [ ] Installation guide created
- [ ] Screenshots added (optional)
- [ ] Troubleshooting section complete
- [ ] Upgrade instructions documented

---

## Deliverables

- [ ] `TF-Engine-Setup-v1.0.0.msi` (or `.exe`)
- [ ] `docs/INSTALLATION_GUIDE.md`
- [ ] Installer build script (`installer/build.bat`)
- [ ] WiX/NSIS configuration files
- [ ] Application icon (`assets/app-icon.ico`)
- [ ] SHA256 checksum file

---

## Documentation Requirements

- [ ] Create `docs/INSTALLATION_GUIDE.md`
- [ ] Document installer build process in `docs/BUILD.md`
- [ ] Add checksums to release notes
- [ ] Update `docs/PROGRESS.md` with installer status

---

## Next Steps

After completing Step 26:
1. Create final installer package
2. Test thoroughly on clean Windows
3. Proceed to **Step 27: User Documentation**

---

**Estimated Completion Time:** 2-3 days
**Phase 5 Progress:** 3 of 5 steps complete
**Overall Progress:** 26 of 28 steps complete (93%)

---

**End of Step 26**
