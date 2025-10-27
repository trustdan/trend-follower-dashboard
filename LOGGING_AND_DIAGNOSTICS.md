# Logging and Diagnostics System

## Overview

A comprehensive logging system has been added to help diagnose issues with checkboxes, dropdowns, and Python integration. All debug information is automatically written to a log file.

## Quick Start

### Viewing the Log File

1. **Open Excel workbook** - Logging starts automatically
2. **Go to Setup sheet**
3. **Click "Open Debug Log" button** - Opens the log in Notepad
4. **Or manually navigate to:** `TradingSystem_Debug.log` in the same folder as your workbook

### What Gets Logged

The system automatically logs:
- ✓ Workbook open/close events
- ✓ All function calls with section headers
- ✓ Python availability detection with detailed diagnostics
- ✓ Checkbox creation (position, linking, errors)
- ✓ Dropdown binding (table checks, validation setup, errors)
- ✓ All errors with error numbers and descriptions
- ✓ Excel version information
- ✓ Timestamps for every event

## Log File Format

```
2025-10-27 10:30:45 | ========== WORKBOOK OPENED ==========
2025-10-27 10:30:45 | Setup sheet exists - workbook already configured
2025-10-27 10:30:45 | Activating TradeEntry sheet
2025-10-27 10:30:45 | --- BindControls() - Start ---
2025-10-27 10:30:45 | tblPresets exists with 5 rows
2025-10-27 10:30:46 | Setting up Preset dropdown in B5...
2025-10-27 10:30:46 | Preset dropdown created successfully
...
```

## Testing Python Integration

When you click "Test Python Integration" button, the log will show detailed diagnostics:

### Expected Log Output (Working Python):
```
--- IsPythonAvailable() - Start ---
Excel Version: 16.0
Test cell: $Z$1 on sheet: Control
Cleared test cell
Attempting Formula2 with: =PY(1+1)
Formula2 property succeeded
Cell formula after assignment: =PY(1+1)
Forcing calculation...
Waiting 2 seconds for cloud execution...
Cell value type: Double
Cell value: 2
SUCCESS: Python returned correct value (2)
Cell HasFormula: True
Cell Formula: =PY(1+1)
PY function availability check: True
Test cell cleared
IsPythonAvailable() - Result: True
```

### Expected Log Output (Python Not Working):
```
--- IsPythonAvailable() - Start ---
Excel Version: 16.0
Test cell: $Z$1 on sheet: Control
Cleared test cell
Attempting Formula2 with: =PY(1+1)
Formula2 property succeeded
Cell formula after assignment: =PY(1+1)
Forcing calculation...
Waiting 2 seconds for cloud execution...
Cell value type: Error
Cell contains error: Error 2015
ERROR in IsPythonAvailable: [2015] Name error
IsPythonAvailable() - Result: False
```

## Common Error Codes

### Python Errors
- **Error 2015** (Name error): `PY()` function not recognized - Python in Excel not available
- **Error 2023** (#REF! error): Reference error in formula
- **Error 2042** (#N/A error): Formula returns #N/A

### Dropdown Errors
- **Error 1004**: "Application-defined or object-defined error" - Usually means table doesn't exist
- **Error 9**: "Subscript out of range" - Sheet or table not found

### Checkbox Errors
- **Error 70**: "Permission denied" - COM automation issue
- **Error 438**: "Object doesn't support this property or method" - CheckBoxes collection issue

## Troubleshooting Guide

### Problem: Checkboxes Overlapping Text

**Check the log for:**
```
--- CreateCheckboxes() - Start ---
Column B Left position: [value]
Checkbox leftPos: [value]
Creating checkbox 1 at row 21, topPos=[Y], leftPos=[X]
```

**What to look for:**
- `leftPos` should be around 160-200 pixels (not 10)
- If `leftPos` is 10, the positioning code isn't working
- Check for errors after each "Creating checkbox" line

### Problem: Preset Dropdown Not Working

**Check the log for:**
```
--- BindControls() - Start ---
tblPresets exists with 5 rows
Setting up Preset dropdown in B5...
Preset dropdown created successfully
```

**What to look for:**
- "tblPresets exists with X rows" - should show 5 rows
- If you see "tblPresets does not exist" - run `Setup.RunInitialSetup`
- If you see "ERROR in BindControls" - check the error details

### Problem: Python Detection Always Says "Not Available"

**Check the log for:**
```
--- IsPythonAvailable() - Start ---
Excel Version: [version]
...
Cell value: [what value?]
```

**Common Issues:**

1. **Formula assignment fails:**
   ```
   Formula2 assignment error: [424] Object required
   ```
   **Solution:** Excel version too old, doesn't support Formula2

2. **Cell stays empty:**
   ```
   Cell is empty - Python may still be calculating
   ```
   **Solution:** Increase wait time or Python not enabled

3. **Error value returned:**
   ```
   Cell contains error: Error 2015
   ```
   **Solution:** Python in Excel feature not available in your Excel version

4. **Wrong value returned:**
   ```
   UNEXPECTED: Python returned: 0 (expected 2)
   ```
   **Solution:** Formula was accepted but didn't execute properly

## Advanced Diagnostics

### Manually Enable/Disable Logging

```vba
' In VBA Immediate window (Ctrl+G):
Call TF_Logger.EnableLogging    ' Turn on logging
Call TF_Logger.DisableLogging   ' Turn off logging
Call TF_Logger.ClearLog         ' Clear and restart log file
```

### Check If Logging Is Working

```vba
' In VBA Immediate window:
? TF_Logger.IsLoggingEnabled()   ' Should return True
? TF_Logger.GetLogPath()         ' Shows log file path
```

### Add Custom Log Messages

```vba
' In your own code:
Call TF_Logger.WriteLog("My custom message")
Call TF_Logger.WriteLogSection("My Section")
Call TF_Logger.WriteLogError("MyFunction", Err.Number, Err.Description)
```

## What to Share When Reporting Issues

When reporting bugs, please share:

1. **The entire log file** (`TradingSystem_Debug.log`)
2. **Your Excel version** (File → Account → About Excel)
3. **What you were trying to do** when the error occurred
4. **Screenshot** of the error or unexpected behavior

## Log File Location

The log file is saved in the same folder as your Excel workbook:
- If workbook is at: `C:\Users\YourName\Documents\TrendFollowing_TradeEntry.xlsm`
- Log file is at: `C:\Users\YourName\Documents\TradingSystem_Debug.log`

## Performance Notes

- Logging has minimal performance impact
- Log file appends, it doesn't overwrite
- If log file gets too large (>10 MB), click "Open Debug Log" and delete old entries
- Or click the "Clear Old Candidates" button and we can add a "Clear Log" button

## Files Modified for Logging

- **NEW:** `VBA/TF_Logger.bas` - Complete logging module
- **MODIFIED:** `VBA/TF_Python_Bridge.bas` - Added detailed Python diagnostics
- **MODIFIED:** `VBA/TF_UI_Builder.bas` - Added checkbox creation logging
- **MODIFIED:** `VBA/TF_UI.bas` - Added dropdown binding logging
- **MODIFIED:** `VBA/ThisWorkbook.cls` - Auto-initialize logger on open
- **MODIFIED:** `VBA/Setup.bas` - Added "Open Debug Log" button

All logging is backward compatible - old workbooks will work fine if TF_Logger module is missing (functions just won't be called).
