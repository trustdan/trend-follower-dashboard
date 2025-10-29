'=============================================================================
' check-vba-version.vbs - Check if VBA modules have the signature fix
'=============================================================================
' Purpose: Verify if TFEngine.bas has the corrected Parse function calls
'
' Usage: cscript check-vba-version.vbs
'
' What it checks:
'   - Whether TFEngine module exists
'   - Whether it has the OLD (broken) or NEW (fixed) Parse function calls
'
' Created: 2025-10-28
'=============================================================================

Option Explicit

Dim objExcel, objWorkbook, objFSO
Dim strScriptDir, strWorkbookPath
Dim comp, moduleCode
Dim hasOldSignature, hasNewSignature

' Create Excel instance
WScript.Echo "=========================================="
WScript.Echo " VBA Signature Version Check"
WScript.Echo "=========================================="
WScript.Echo ""

On Error Resume Next
Set objExcel = CreateObject("Excel.Application")
If Err.Number <> 0 Then
    WScript.Echo "[FAIL] Cannot create Excel.Application"
    WScript.Quit 1
End If
On Error Goto 0

objExcel.Visible = False
objExcel.DisplayAlerts = False

' Get script directory
Set objFSO = CreateObject("Scripting.FileSystemObject")
strScriptDir = objFSO.GetParentFolderName(WScript.ScriptFullName)

' Check if workbook exists
strWorkbookPath = strScriptDir & "\TradingPlatform.xlsm"
WScript.Echo "Checking: " & strWorkbookPath
WScript.Echo ""

If Not objFSO.FileExists(strWorkbookPath) Then
    WScript.Echo "[FAIL] TradingPlatform.xlsm not found"
    objExcel.Quit
    WScript.Quit 1
End If

' Open workbook
On Error Resume Next
Set objWorkbook = objExcel.Workbooks.Open(strWorkbookPath)
If Err.Number <> 0 Then
    WScript.Echo "[FAIL] Cannot open workbook: " & Err.Description
    objExcel.Quit
    WScript.Quit 1
End If
On Error Goto 0

' Find TFEngine module
On Error Resume Next
Set comp = objWorkbook.VBProject.VBComponents("TFEngine")
If Err.Number <> 0 Then
    WScript.Echo "[FAIL] TFEngine module not found in workbook"
    WScript.Echo ""
    WScript.Echo "Run fix-vba-modules.bat to import the VBA modules"
    objWorkbook.Close False
    objExcel.Quit
    WScript.Quit 1
End If
On Error Goto 0

WScript.Echo "[OK]   TFEngine module found"
WScript.Echo ""

' Get module code
moduleCode = comp.CodeModule.Lines(1, comp.CodeModule.CountOfLines)

' Check for old signature pattern
hasOldSignature = False
hasNewSignature = False

If InStr(moduleCode, "sizeResult = TFHelpers.ParseSizingJSON(result.JsonOutput)") > 0 Then
    hasOldSignature = True
End If

If InStr(moduleCode, "TFHelpers.ParseSizingJSON result.JsonOutput, sizeResult") > 0 Then
    hasNewSignature = True
End If

' Report results
WScript.Echo "Version Check:"
WScript.Echo ""

If hasNewSignature And Not hasOldSignature Then
    WScript.Echo "[OK]   ✅ YOU HAVE THE FIXED VERSION!"
    WScript.Echo ""
    WScript.Echo "Your VBA modules have the corrected function signatures."
    WScript.Echo "All macros should work correctly now."
    WScript.Echo ""
    WScript.Echo "Next steps:"
    WScript.Echo "1. Open TradingPlatform.xlsm"
    WScript.Echo "2. Enable macros"
    WScript.Echo "3. Test Position Sizing sheet"
    WScript.Echo "4. All buttons should work!"
    WScript.Echo ""
    objWorkbook.Close False
    objExcel.Quit
    WScript.Quit 0

ElseIf hasOldSignature Then
    WScript.Echo "[FAIL] ❌ YOU HAVE THE OLD (BROKEN) VERSION"
    WScript.Echo ""
    WScript.Echo "Your VBA modules still have the old function signatures."
    WScript.Echo "This will cause 'argument not optional' errors."
    WScript.Echo ""
    WScript.Echo "TO FIX:"
    WScript.Echo "1. Close this Excel workbook"
    WScript.Echo "2. Run: fix-vba-modules.bat"
    WScript.Echo "3. This will import the corrected VBA modules"
    WScript.Echo ""
    objWorkbook.Close False
    objExcel.Quit
    WScript.Quit 1

Else
    WScript.Echo "[WARN] ⚠️  CANNOT DETERMINE VERSION"
    WScript.Echo ""
    WScript.Echo "Could not find the signature pattern in TFEngine module."
    WScript.Echo "The module might be corrupted or incomplete."
    WScript.Echo ""
    WScript.Echo "TO FIX:"
    WScript.Echo "1. Close this Excel workbook"
    WScript.Echo "2. Run: fix-vba-modules.bat"
    WScript.Echo "3. This will re-import all VBA modules"
    WScript.Echo ""
    objWorkbook.Close False
    objExcel.Quit
    WScript.Quit 1
End If
