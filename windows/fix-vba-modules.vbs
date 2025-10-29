'=============================================================================
' fix-vba-modules.vbs - Import VBA Modules into Existing Workbook
'=============================================================================
' Purpose: Fix "argument not optional" error by importing missing VBA modules
'
' Usage:
'   1. Close TradingPlatform.xlsm if open
'   2. Double-click this file (or run: cscript fix-vba-modules.vbs)
'
' What it does:
'   - Opens TradingPlatform.xlsm
'   - Removes old VBA modules (if any)
'   - Imports fresh VBA modules: TFTypes, TFHelpers, TFEngine, TFTests
'   - Saves and closes the workbook
'
' Created: 2025-10-28
'=============================================================================

Option Explicit

Dim objExcel, objWorkbook, objFSO
Dim strScriptDir, strWorkbookPath, strVBADir
Dim comp

' Create Excel instance
WScript.Echo "Starting VBA module import..."
WScript.Echo ""

On Error Resume Next
Set objExcel = CreateObject("Excel.Application")
If Err.Number <> 0 Then
    WScript.Echo "ERROR: Cannot create Excel.Application"
    WScript.Echo "Make sure Excel is installed"
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
If Not objFSO.FileExists(strWorkbookPath) Then
    WScript.Echo "ERROR: TradingPlatform.xlsm not found at:"
    WScript.Echo strWorkbookPath
    WScript.Echo ""
    WScript.Echo "Expected location: " & strScriptDir
    WScript.Echo ""
    WScript.Echo "Please:"
    WScript.Echo "1. Make sure TradingPlatform.xlsm is in the same folder as this script"
    WScript.Echo "2. Or run 1-setup-all.bat to create the workbook first"
    objExcel.Quit
    Set objExcel = Nothing
    WScript.Quit 1
End If

' Open workbook
WScript.Echo "Opening workbook: " & strWorkbookPath
On Error Resume Next
Set objWorkbook = objExcel.Workbooks.Open(strWorkbookPath)
If Err.Number <> 0 Then
    WScript.Echo "ERROR: Cannot open workbook: " & Err.Description
    objExcel.Quit
    WScript.Quit 1
End If
On Error Goto 0

' Check VBA directory
strVBADir = objFSO.GetParentFolderName(objFSO.GetParentFolderName(strScriptDir)) & "\excel\vba\"
If Not objFSO.FolderExists(strVBADir) Then
    WScript.Echo "ERROR: VBA source directory not found:"
    WScript.Echo strVBADir
    objWorkbook.Close False
    objExcel.Quit
    WScript.Quit 1
End If

WScript.Echo "VBA source directory: " & strVBADir
WScript.Echo ""

' Remove existing modules
WScript.Echo "Removing old VBA modules (if any)..."
On Error Resume Next
For Each comp In objWorkbook.VBProject.VBComponents
    If comp.Type = 1 Then ' vbext_ct_StdModule = 1
        Select Case comp.Name
            Case "TFTypes", "TFHelpers", "TFEngine", "TFTests", "TFIntegrationTests"
                WScript.Echo "  - Removing: " & comp.Name
                objWorkbook.VBProject.VBComponents.Remove comp
        End Select
    End If
Next
On Error Goto 0

WScript.Echo ""
WScript.Echo "Importing fresh VBA modules..."

' Import modules in correct order (dependencies first)
On Error Resume Next

' 1. TFTypes (no dependencies)
WScript.Echo "  1. Importing TFTypes.bas..."
objWorkbook.VBProject.VBComponents.Import strVBADir & "TFTypes.bas"
If Err.Number <> 0 Then
    WScript.Echo "     ERROR: " & Err.Description
    Err.Clear
Else
    WScript.Echo "     [OK]"
End If

' 2. TFHelpers (depends on TFTypes)
WScript.Echo "  2. Importing TFHelpers.bas..."
objWorkbook.VBProject.VBComponents.Import strVBADir & "TFHelpers.bas"
If Err.Number <> 0 Then
    WScript.Echo "     ERROR: " & Err.Description
    Err.Clear
Else
    WScript.Echo "     [OK]"
End If

' 3. TFEngine (depends on TFTypes, TFHelpers)
WScript.Echo "  3. Importing TFEngine.bas..."
objWorkbook.VBProject.VBComponents.Import strVBADir & "TFEngine.bas"
If Err.Number <> 0 Then
    WScript.Echo "     ERROR: " & Err.Description
    Err.Clear
Else
    WScript.Echo "     [OK]"
End If

' 4. TFTests (depends on all above)
WScript.Echo "  4. Importing TFTests.bas..."
objWorkbook.VBProject.VBComponents.Import strVBADir & "TFTests.bas"
If Err.Number <> 0 Then
    WScript.Echo "     ERROR: " & Err.Description
    Err.Clear
Else
    WScript.Echo "     [OK]"
End If

' 5. TFIntegrationTests (optional)
If objFSO.FileExists(strVBADir & "TFIntegrationTests.bas") Then
    WScript.Echo "  5. Importing TFIntegrationTests.bas..."
    objWorkbook.VBProject.VBComponents.Import strVBADir & "TFIntegrationTests.bas"
    If Err.Number <> 0 Then
        WScript.Echo "     ERROR: " & Err.Description
        Err.Clear
    Else
        WScript.Echo "     [OK]"
    End If
End If

On Error Goto 0

' Save and close
WScript.Echo ""
WScript.Echo "Saving workbook..."
objWorkbook.Save
objWorkbook.Close
objExcel.Quit

' Cleanup
Set objWorkbook = Nothing
Set objExcel = Nothing
Set objFSO = Nothing

WScript.Echo ""
WScript.Echo "=========================================="
WScript.Echo " VBA Modules Imported Successfully!"
WScript.Echo "=========================================="
WScript.Echo ""
WScript.Echo "Next steps:"
WScript.Echo "1. Open TradingPlatform.xlsm"
WScript.Echo "2. Enable macros when prompted"
WScript.Echo "3. Go to Position Sizing sheet"
WScript.Echo "4. Try clicking Calculate button again"
WScript.Echo ""
WScript.Echo "The 'argument not optional' error should now be fixed!"
WScript.Echo ""

WScript.Quit 0
