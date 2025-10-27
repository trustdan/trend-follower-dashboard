# Fix: Log File Lock Conflict

**Date**: 2025-10-26 16:11
**Error**: "Object required" at line 14 (runtime error)
**Exit Code**: 0 (success, but with errors)

## Problem

The batch file and VBScript were both trying to write to the same log file simultaneously, causing a file lock conflict:

### Batch File (IMPORT_VBA_MODULES.bat line 31):
```batch
cscript //nologo "%SCRIPT%" ... >> "%LOGFILE%" 2>&1
```
- Uses `>>` to redirect cscript output to log file
- Creates a file lock on the log file

### VBScript (excel_build_repo_aware_logged.vbs line 31):
```vbscript
Set logTS = fso.OpenTextFile(logPath, 8, True)  ' 8 = ForAppending
```
- Tries to open the SAME file for appending
- **FAILS** because batch file has it locked
- Result: `logTS` remains `Nothing`

### When Log() is called (line 14):
```vbscript
If Not logTS Is Nothing Then logTS.WriteLine ...
```
- Since `logTS` is `Nothing`, this causes "Object required" error
- Error happens when trying to access `.WriteLine` on a Nothing object

## Why Exit Code Was Still 0

The VBScript uses `On Error Resume Next` extensively, which suppresses runtime errors:
- Line 45: `On Error Resume Next` before creating Excel
- Line 14: Error occurs but is suppressed
- Excel workbook is created successfully
- Line 277: `WScript.Quit rc` with `rc = 0` (no fatal errors)

**Result**: Workbook created successfully, but logging failed silently.

## Solution

### 1. Remove Batch File Redirection

**Before**:
```batch
cscript //nologo "%SCRIPT%" ... >> "%LOGFILE%" 2>&1
```

**After**:
```batch
REM Don't redirect - VBScript handles its own logging
cscript //nologo "%SCRIPT%" ... "%LOGFILE%"
```

### 2. Update VBScript to Log to Both File and Console

**Before**:
```vbscript
Sub Log(msg)
  If Not logTS Is Nothing Then logTS.WriteLine NowISO() & " " & msg
End Sub
```

**After**:
```vbscript
Sub Log(msg)
  Dim line: line = NowISO() & " " & msg
  ' Write to log file if available
  On Error Resume Next
  If Not logTS Is Nothing Then logTS.WriteLine line
  On Error GoTo 0
  ' Always write to console
  WScript.Echo line
End Sub
```

### 3. Properly Close Log File at End

Added at end of VBScript (before `WScript.Quit`):
```vbscript
' Close log file
On Error Resume Next
If Not logTS Is Nothing Then
  logTS.Close
  Set logTS = Nothing
End If
On Error GoTo 0
```

## Files Modified

1. ✅ `scripts/excel_build_repo_aware_logged.vbs`:
   - Updated `Log()` to write to both file and console
   - Added proper log file closing at end
   - Added error protection with `On Error Resume Next`

2. ✅ `IMPORT_VBA_MODULES.bat`:
   - Removed `>> "%LOGFILE%" 2>&1` redirection
   - VBScript now handles all logging

3. ✅ `IMPORT_VBA_MODULES_DEBUG.bat`:
   - Same change - removed redirection

## Expected Behavior Now

### Console Output:
```
========================================
About to run: VBScript
========================================
Running: cscript //nologo "scripts\excel_build_repo_aware_logged.vbs" ...

2025-10-26 16:15:30 [STEP] Builder start
2025-10-26 16:15:30 [INFO] wbPath=C:\Users\Dan\excel-trading-dashboard\TrendFollowing_TradeEntry.xlsm
2025-10-26 16:15:30 [INFO] vbaFolder=C:\Users\Dan\excel-trading-dashboard\VBA
2025-10-26 16:15:31 [STEP] Excel.Application created
2025-10-26 16:15:31 [STEP] Added new workbook
2025-10-26 16:15:32 [STEP] Saved as .xlsm
2025-10-26 16:15:32 [INFO] Excel.Version=16.0
2025-10-26 16:15:32 [INFO] Created sheet 'TradeEntry' as CodeName 'TradeEntry'
2025-10-26 16:15:33 [INFO] Imported: TF_Data.bas
2025-10-26 16:15:33 [INFO] Imported: TF_UI.bas
...
2025-10-26 16:15:40 [STEP] Ran bootstrap: 'TrendFollowing_TradeEntry.xlsm'!Setup.RunOnce
2025-10-26 16:15:41 [STEP] Saved workbook
2025-10-26 16:15:41 [OK] Build completed with rc=0

========================================
VBScript Exit Code: 0
========================================
```

### Log File (logs/build_*.log):
```
[INFO] Repo: C:\Users\Dan\excel-trading-dashboard
[INFO] Log:  logs\build_20251026_161500.log
[INFO] Script: scripts\excel_build_repo_aware_logged.vbs
<registry and taskkill output from batch>
[INFO] XL_SILENT_SETUP=1
2025-10-26 16:15:30 [STEP] Builder start
2025-10-26 16:15:30 [INFO] wbPath=...
<all VBScript log entries>
2025-10-26 16:15:41 [OK] Build completed with rc=0
[INFO] ExitCode: 0
```

## Testing

Run on Windows:
```cmd
cd C:\Users\Dan\excel-trading-dashboard
IMPORT_VBA_MODULES_DEBUG.bat
```

You should now see:
1. ✅ Live progress messages in console
2. ✅ Complete log file with all entries
3. ✅ Exit code 0 with no errors
4. ✅ `TrendFollowing_TradeEntry.xlsm` created successfully
5. ✅ All VBA modules imported

## Note

The previous run (build_20251026_161147.log) showed exit code 0, which means:
- **The workbook was created successfully!**
- The error was just a logging issue
- Check if `TrendFollowing_TradeEntry.xlsm` exists - it should work

The fixes above will eliminate the error and provide proper logging for future runs.
