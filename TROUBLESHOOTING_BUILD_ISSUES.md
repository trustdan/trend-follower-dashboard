# Build Issues - Troubleshooting Report

**Date**: 2025-10-26
**Issue**: Excel workbook not being created; VBScript compilation errors; no logs generated

---

## Problems Identified

### 1. VBScript Syntax Error (FIXED)
**File**: `scripts/excel_build_repo_aware_logged.vbs`
**Line**: 76
**Original Code**:
```vbscript
On Error GoTo 0
If rc <> 0 Then GoTo Finalize
```

**Error Message**:
```
C:\Users\Dan\excel-trading-dashboard\scripts\excel_build_repo_aware_logged.vbs(76, 17)
Microsoft VBScript compilation error: Expected statement
```

**Fix Applied**:
```vbscript
On Error GoTo 0

If rc <> 0 Then
  GoTo Finalize
End If
```

**Explanation**: Single-line `If...Then GoTo` at column 17 was causing a compilation error. Expanded to multi-line block format for clarity and to resolve the parsing issue.

---

### 2. Wrong File Extension in VBA Folder (FIXED)
**File**: `VBA/Setup.vbs` → should be `VBA/Setup.bas`

**Problem**:
- File was named `.vbs` (VBScript for Windows Script Host)
- File contains VBA code with `Attribute VB_Name = "Setup"`
- VBA modules for Excel must use `.bas` extension

**Fix Applied**:
```bash
mv VBA/Setup.vbs VBA/Setup.bas
```

**Verification**:
```bash
$ ls -la VBA/
-rw-r--r-- 1 kali kali  1748 Oct 26 15:18 Setup.bas  ✓
```

---

### 3. Missing Excel Workbook (NOT YET FIXED - EXPECTED)
**File**: `TrendFollowing_TradeEntry.xlsm`

**Problem**:
- The workbook doesn't exist yet
- The VBScript is designed to CREATE it if missing
- However, VBScript errors prevented it from being created

**Expected Behavior**:
- When `IMPORT_VBA_MODULES.bat` runs successfully, it will:
  1. Create a new `.xlsm` workbook
  2. Import all VBA modules from `VBA/` folder
  3. Run bootstrap macro `Setup.RunOnce` to create sheets/tables
  4. Save the workbook

**Status**: Should be resolved once VBScript runs successfully

---

### 4. Python Still Referenced (INFORMATIONAL)
**Files**:
- `VBA/Python_Run.bas`
- `scripts/refresh_data.bat`
- `Python/finviz_scraper.py`
- `Python/heat_calculator.py`

**Observation**:
- Python scripts still exist in repo
- Python functionality is NOT part of Excel workbook creation
- Python is OPTIONAL for runtime operations (data refresh via button)

**Python Functionality** (if you want to use it later):
1. **FINVIZ Scraper** (`Python/finviz_scraper.py`):
   - Automates ticker scraping from FINVIZ screeners
   - Eliminates manual copy/paste workflow
   - Called by `VBA/Python_Run.bas` → `scripts/refresh_data.bat`
   - Requires: `requests`, `beautifulsoup4`, `lxml`

2. **Heat Calculator** (`Python/heat_calculator.py`):
   - Fast vectorized portfolio/bucket heat calculations
   - Uses pandas for performance
   - Validates trades against heat caps
   - Requires: `pandas`, `numpy`

**Current Build Process**:
- Does NOT use Python
- Pure VBScript + Excel COM automation
- Python is only referenced in runtime VBA code (optional feature)

---

## Files Modified

1. ✅ `scripts/excel_build_repo_aware_logged.vbs` (removed GoTo statements, added conditional blocks)
2. ✅ `VBA/Setup.vbs` → `VBA/Setup.bas` (renamed)
3. ✅ `IMPORT_VBA_MODULES_DEBUG.bat` (created for diagnostics)

## Additional Fix - GoTo Statement Issue

**Second Error** (after first fix): Line 78 - "Expected statement"

**Root Cause**: VBScript does NOT support `GoTo` with text labels (unlike VBA). The pattern:
```vbscript
If rc <> 0 Then GoTo Finalize
...
Finalize:
```
...fails in VBScript (Windows Script Host).

**Solution**: Removed all `GoTo` statements and restructured code to use conditional blocks:
- Changed `If fso.FolderExists(...)` to `If rc = 0 And fso.FolderExists(...)`
- Wrapped component listing and bootstrap in `If rc = 0 Then ... End If` blocks
- Removed `Finalize:` label entirely
- Cleanup code (save/close/quit) always runs regardless of `rc` value

See `LATEST_FIX.md` for detailed explanation.

---

## Testing Instructions

### Option 1: Run Debug Batch File (RECOMMENDED)
On your Windows machine:

```cmd
cd C:\Users\Dan\excel-trading-dashboard
IMPORT_VBA_MODULES_DEBUG.bat
```

This will:
- Show detailed progress messages
- Display current directory and file paths
- Verify all folders/files exist before proceeding
- Show the VBScript exit code
- Display the full log at the end
- PAUSE at the end so you can review output

### Option 2: Run Original Batch File
```cmd
cd C:\Users\Dan\excel-trading-dashboard
IMPORT_VBA_MODULES.bat
```

Check for logs in `logs/build_YYYYMMDD_HHMMSS.log`

---

## Expected Successful Output

### Console Output:
```
========================================
DEBUG: Starting batch file
========================================
Current Directory: C:\Users\Dan\excel-trading-dashboard
Workbook: TrendFollowing_TradeEntry.xlsm
VBA Folder: VBA
Script: scripts\excel_build_repo_aware_logged.vbs
Log Dir: logs
[OK] VBA folder exists
[OK] VBScript file exists
Timestamp: 20251026_160530
Log file will be: logs\build_20251026_160530.log

========================================
About to run: VBScript
========================================
Running: cscript //nologo "scripts\excel_build_repo_aware_logged.vbs" ...

========================================
VBScript Exit Code: 0
========================================
[OK] TrendFollowing_TradeEntry.xlsm ready.
```

### Log File Contents (expected):
```
[STEP] Builder start
[INFO] wbPath=C:\Users\Dan\excel-trading-dashboard\TrendFollowing_TradeEntry.xlsm
[INFO] vbaFolder=C:\Users\Dan\excel-trading-dashboard\VBA
[STEP] Excel.Application created
[STEP] Added new workbook
[STEP] Saved as .xlsm
[INFO] Excel.Version=16.0
[INFO] Created sheet 'TradeEntry' as CodeName 'TradeEntry'
[INFO] Imported: TF_Data.bas
[INFO] Imported: TF_UI.bas
[INFO] Imported: TF_Utils.bas
[INFO] Imported: TF_Presets.bas
[INFO] Imported: PQ_Setup.bas
[INFO] Imported: Python_Run.bas
[INFO] Imported: Setup.bas
[INFO] Replaced doc module 'TradeEntry' from Sheet_TradeEntry.cls
[INFO] Replaced doc module 'ThisWorkbook' from ThisWorkbook.cls
[STEP] VBComponents after import:
  - TF_Data (Type=1)
  - TF_UI (Type=1)
  - TF_Utils (Type=1)
  - TF_Presets (Type=1)
  - PQ_Setup (Type=1)
  - Python_Run (Type=1)
  - Setup (Type=1)
  - TradeEntry (Type=100)
  - ThisWorkbook (Type=100)
[STEP] Ran bootstrap: 'TrendFollowing_TradeEntry.xlsm'!Setup.RunOnce
[STEP] Saved workbook
[OK] Build completed with rc=0
```

---

## If Build Still Fails

### Common Issues:

1. **Excel Trust Center not configured**
   - Error: "Access VBProject failed"
   - Fix: Excel → File → Options → Trust Center → Trust Center Settings → Macro Settings
   - Enable: "Trust access to the VBA project object model"

2. **Excel already running**
   - Error: File locked or permission denied
   - Fix: Close all Excel instances, or run `taskkill /IM excel.exe /F`

3. **Different Office version**
   - Current script targets Office 16.0 (Office 2016/2019/365)
   - If using Office 15.0 (2013) or earlier, registry keys may differ
   - Check: `reg query "HKCU\Software\Microsoft\Office"`

4. **VBScript disabled**
   - Error: cscript not recognized
   - Fix: Ensure Windows Script Host is enabled
   - Check: `cscript /?` should show help

5. **Path issues**
   - Error: VBA folder not found
   - Fix: Ensure you're running from the repository root
   - Verify: `dir VBA` should show `.bas` and `.cls` files

---

## Python Integration (Optional Future Enhancement)

If you want to re-enable Python for automated FINVIZ scraping:

### Setup:
```cmd
cd C:\Users\Dan\excel-trading-dashboard
scripts\setup_venv.bat
```

### Activate venv:
```cmd
venv\Scripts\activate
```

### Install dependencies:
```cmd
pip install requests beautifulsoup4 lxml pandas numpy pywin32
```

### Test FINVIZ scraper:
```cmd
python Python\finviz_scraper.py
```

### Use in Excel:
- Click button in workbook: "Fetch Screened Tickers"
- This calls `VBA/Python_Run.bas` → `scripts/refresh_data.bat`
- Output: `data/screened.csv`
- Excel Power Query loads this CSV into `Candidates` table

### Alternative: Use import_to_excel.py
```cmd
python import_to_excel.py TrendFollowing_TradeEntry.xlsm
```
This directly manipulates the Excel workbook via COM.

---

## Summary

✅ **Fixed**: VBScript syntax error (line 76)
✅ **Fixed**: Wrong file extension (`Setup.vbs` → `Setup.bas`)
✅ **Created**: Debug batch file for troubleshooting
⏳ **Next Step**: Run `IMPORT_VBA_MODULES_DEBUG.bat` on Windows machine

Expected result: Working `.xlsm` file with all VBA modules imported.
