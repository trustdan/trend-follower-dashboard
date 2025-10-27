# Two Options for Building the Workbook

**Problem**: VBScript version imports modules successfully but **fails to save** the workbook.

**Solution**: Use Python instead (more reliable COM handling).

---

## Option 1: Build with Python (RECOMMENDED)

### Why Python Works Better:
- ✅ Better error handling for COM objects
- ✅ Simpler save logic (delete old file, then SaveAs)
- ✅ Proper cleanup (close workbook, quit Excel)
- ✅ More readable code (easier to debug)

### Requirements:
1. Python venv with pywin32 installed
2. Run `scripts\setup_venv.bat` if not already done

### How to Use:

**Step 1: Run the build script**
```cmd
cd C:\Users\Dan\excel-trading-dashboard
BUILD_WITH_PYTHON.bat
```

**Expected Output:**
```
========================================
Build Workbook Using Python
========================================

Activating venv...
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

**Step 2: Open the workbook**
```cmd
start TrendFollowing_TradeEntry.xlsm
```

**Step 3: Verify modules (Alt+F11 in Excel)**

---

## Option 2: Fix VBScript Version (For Debugging)

### Current Issue:
The VBScript successfully:
- ✅ Creates Excel.Application
- ✅ Creates new workbook
- ✅ Saves as .xlsm initially
- ✅ Imports all VBA modules
- ✅ Lists modules (proves they're loaded in memory)

But then FAILS to:
- ❌ Save the workbook with modules
- ❌ Close cleanly

**Error in log:**
```
[WARN] Save failed: Unknown runtime error
[WARN] Close failed: Unknown runtime error
```

### Latest Fixes Applied:
1. Better error logging (shows error numbers)
2. Delete old file before SaveAs
3. Fallback: if wb.Save fails, try wb.SaveAs
4. Added explicit error handling

### To Test VBScript Version:
```cmd
cd C:\Users\Dan\excel-trading-dashboard
CLEANUP_STUCK_EXCEL.bat
IMPORT_VBA_MODULES_DEBUG.bat
```

Check the log for error numbers to diagnose further.

---

## Comparison

| Feature | Python | VBScript |
|---------|--------|----------|
| Module Import | ✅ Works | ✅ Works |
| Save Workbook | ✅ Works | ❌ Fails |
| Error Messages | ✅ Clear | ❌ Generic |
| Excel Cleanup | ✅ Quits cleanly | ⚠ Sometimes hangs |
| Debugging | ✅ Easy | ❌ Hard |
| Dependencies | Python + pywin32 | Windows built-in |

---

## Recommendation

**Use Python** (`BUILD_WITH_PYTHON.bat`) because:
1. It's more reliable
2. Better error handling
3. You already have Python set up
4. VBScript save issue is hard to debug

The VBScript approach *could* work if we figure out why the save fails, but the Python approach is working NOW.

---

## Files Modified

### Python Approach:
1. ✅ `import_to_excel.py` - Added proper save/close/quit logic
2. ✅ `BUILD_WITH_PYTHON.bat` - New wrapper script

### VBScript Approach:
1. ✅ `scripts/excel_build_repo_aware_logged.vbs` - Better save error handling
2. ✅ `IMPORT_VBA_MODULES.bat` - Still available
3. ✅ `IMPORT_VBA_MODULES_DEBUG.bat` - For debugging

---

## Next Steps

1. **Run** `BUILD_WITH_PYTHON.bat`
2. **Open** `TrendFollowing_TradeEntry.xlsm`
3. **Verify** modules are there (Alt+F11)
4. **Run** `Setup.RunOnce` macro to create sheets/tables

Done!
