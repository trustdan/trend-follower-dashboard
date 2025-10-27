# Quick Start - Build Excel Workbook

**Goal**: Create `TrendFollowing_TradeEntry.xlsm` with all VBA modules imported.

---

## One-Line Build

```cmd
cd C:\Users\Dan\excel-trading-dashboard
BUILD_WITH_PYTHON.bat
```

That's it! The script will:
1. ✅ Check for Python venv
2. ✅ Auto-install pywin32 if missing
3. ✅ Configure Excel Trust Center
4. ✅ Kill any stuck Excel processes
5. ✅ Delete old workbook if exists
6. ✅ Import all VBA modules
7. ✅ Save and close properly

---

## What to Expect

### Successful Output:
```
========================================
Build Workbook Using Python
========================================

Activating venv...

Checking dependencies...
pywin32 not found - installing...
pywin32 installed successfully

Current directory: C:\Users\Dan\excel-trading-dashboard
Target workbook: TrendFollowing_TradeEntry.xlsm

Configuring Excel Trust Center...
Closing any existing Excel instances...

========================================
Running Python import script...
========================================

======================================================================
VBA Module Import Automation
======================================================================

📁 Modules found: 7 .bas, 2 .cls
Starting Excel…
Creating new workbook…

📥 Importing standard modules…
  ✅ PQ_Setup.bas
  ✅ Python_Run.bas
  ✅ Setup.bas
  ✅ TF_Data.bas
  ✅ TF_Presets.bas
  ✅ TF_UI.bas
  ✅ TF_Utils.bas

📥 Importing class modules…
  ✅ Sheet_TradeEntry.cls
  ✅ ThisWorkbook.cls (replaced)

💾 Saving to: C:\Users\Dan\excel-trading-dashboard\TrendFollowing_TradeEntry.xlsm
  (Deleted existing file)

✅ Import complete! 9 modules imported.
📁 File saved: C:\Users\Dan\excel-trading-dashboard\TrendFollowing_TradeEntry.xlsm
✅ Workbook closed
✅ Excel quit successfully

========================================
SUCCESS!
========================================

Workbook created: TrendFollowing_TradeEntry.xlsm
```

---

## Next Steps

### 1. Open the Workbook
```cmd
start TrendFollowing_TradeEntry.xlsm
```

Or double-click `TrendFollowing_TradeEntry.xlsm` in File Explorer.

### 2. Verify Modules (Alt+F11)
Press **Alt+F11** to open the VBA Editor.

You should see:
- **Modules**:
  - PQ_Setup
  - Python_Run
  - Setup
  - TF_Data
  - TF_Presets
  - TF_UI
  - TF_Utils
- **ThisWorkbook** (with code)
- **Sheet1** (or TradeEntry)

### 3. Run Initial Setup Macro
In the VBA Editor:
1. Press **F5** (Run)
2. Select **Setup.RunOnce**
3. Click **Run**

This will create all the sheets, tables, and named ranges.

---

## Troubleshooting

### "Python venv not found"
Run the venv setup first:
```cmd
scripts\setup_venv.bat
```

### "pywin32 installation failed"
Manually install:
```cmd
venv\Scripts\activate
pip install pywin32
```

### "Excel blocked VBProject access"
Enable in Excel:
1. File → Options
2. Trust Center → Trust Center Settings
3. Macro Settings
4. Check "Trust access to the VBA project object model"

### Workbook created but no modules
This was the old VBScript issue - **Python version should work**.

If modules are still missing:
1. Check the console output for errors
2. Try running again (script deletes old file first)
3. Check Task Manager - ensure no Excel.exe is running

---

## Old Scripts (Not Needed)

These are kept for reference but not actively used:
- `old_IMPORT_VBA_MODULES.bat` - VBScript version (save issues)
- `old_IMPORT_VBA_MODULES_DEBUG.bat` - VBScript with verbose logging
- `old_CLEANUP_STUCK_EXCEL.bat` - Manual cleanup (Python version handles this)
- `old_VERIFY_MODULES.bat` - Module checker (just use Alt+F11 instead)

---

## Summary

✅ **Use**: `BUILD_WITH_PYTHON.bat`
✅ **Auto-installs**: pywin32
✅ **Auto-cleans**: Old files, stuck processes
✅ **Works reliably**: Python COM > VBScript COM

**One command. Done.**
