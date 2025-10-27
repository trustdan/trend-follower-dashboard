# Latest Fix - VBScript GoTo Issue

**Date**: 2025-10-26 16:08
**Previous Error**: Line 78 - "Expected statement"

## Problem

VBScript has **very limited support for GoTo** statements. Unlike VBA (which runs inside Excel), VBScript (Windows Script Host) does NOT support `GoTo` with text labels.

### What Doesn't Work in VBScript:
```vbscript
If rc <> 0 Then
  GoTo Finalize    ' ❌ ERROR: VBScript doesn't support this
End If

' ... more code ...

Finalize:
  ' cleanup code
```

### What VBScript Actually Supports:
- `GoTo` with line numbers (rarely used, deprecated)
- `On Error GoTo 0` (disable error handler)
- `On Error Resume Next` (continue on errors)

**BUT NOT**: `GoTo` with text labels like `Finalize:`

## Solution

Removed the `GoTo Finalize` pattern entirely and restructured code to use conditional blocks:

### Before (lines 75-79):
```vbscript
On Error GoTo 0

If rc <> 0 Then
  GoTo Finalize    ' ❌ FAILS
End If
```

### After (line 75-76):
```vbscript
On Error GoTo 0

' Code continues - everything below wrapped in If rc = 0 checks
```

### Changes Made:

1. **Line 175**: Changed `If fso.FolderExists(vbaFolder) Then` → `If rc = 0 And fso.FolderExists(vbaFolder) Then`

2. **Line 191**: Changed `If fso.FolderExists(vbaFolder) Then` → `If rc = 0 And fso.FolderExists(vbaFolder) Then`

3. **Line 218-226**: Wrapped component listing in `If rc = 0 Then ... End If`

4. **Line 229-250**: Wrapped bootstrap attempts in `If rc = 0 Then ... End If`

5. **Removed**: The `Finalize:` label (line 253 in old version)

## Result

Now the script structure is:

```vbscript
' Setup and initialization
rc = 0

' Try to open/create Excel
If error Then rc = 10
If error Then rc = 11

' Only proceed with import if rc = 0
If rc = 0 And FolderExists Then
  ' Pass 0: Create sheets
End If

If rc = 0 And FolderExists Then
  ' Pass 1: Import modules
End If

If rc = 0 Then
  ' List components
End If

If rc = 0 Then
  ' Run bootstrap
End If

' ALWAYS run cleanup (save/close/quit)
' Save & close
If Not wb Is Nothing Then
  wb.Save
  wb.Close
End If
Set wb = Nothing
If Not xl Is Nothing Then
  xl.Quit
End If
Set xl = Nothing

' Exit with return code
WScript.Quit rc
```

## Why This Works

- **No GoTo**: VBScript doesn't need to support text labels
- **Conditional execution**: Code only runs if `rc = 0` (no errors)
- **Cleanup always runs**: Save/close/quit happens regardless of `rc` value
- **Proper exit code**: `WScript.Quit rc` returns success (0) or failure (10, 11)

## Testing

Run this on Windows:
```cmd
cd C:\Users\Dan\excel-trading-dashboard
IMPORT_VBA_MODULES_DEBUG.bat
```

Expected result: **Exit Code 0** (success)
