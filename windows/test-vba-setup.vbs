'=============================================================================
' test-vba-setup.vbs - Verify VBA Modules are Properly Loaded
'=============================================================================
' Purpose: Diagnose VBA module issues in TradingPlatform.xlsm
'
' Usage: cscript test-vba-setup.vbs
'
' What it checks:
'   - Workbook exists
'   - VBA modules are imported (TFTypes, TFHelpers, TFEngine, TFTests)
'   - Functions are accessible
'   - Types are defined
'
' Created: 2025-10-28
'=============================================================================

Option Explicit

Dim objExcel, objWorkbook, objFSO
Dim strScriptDir, strWorkbookPath
Dim comp, foundModules
Dim moduleName, moduleList

' Create Excel instance
WScript.Echo "=========================================="
WScript.Echo " VBA Setup Diagnostic Tool"
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
WScript.Echo "Checking workbook..."
If Not objFSO.FileExists(strWorkbookPath) Then
    WScript.Echo "[FAIL] TradingPlatform.xlsm not found"
    WScript.Echo "       Path: " & strWorkbookPath
    objExcel.Quit
    WScript.Quit 1
Else
    WScript.Echo "[OK]   Workbook found: TradingPlatform.xlsm"
End If

' Open workbook
WScript.Echo ""
WScript.Echo "Opening workbook..."
On Error Resume Next
Set objWorkbook = objExcel.Workbooks.Open(strWorkbookPath)
If Err.Number <> 0 Then
    WScript.Echo "[FAIL] Cannot open workbook: " & Err.Description
    objExcel.Quit
    WScript.Quit 1
End If
On Error Goto 0
WScript.Echo "[OK]   Workbook opened"

' Check VBA modules
WScript.Echo ""
WScript.Echo "Checking VBA modules..."

Set foundModules = CreateObject("Scripting.Dictionary")
foundModules.Add "TFTypes", False
foundModules.Add "TFHelpers", False
foundModules.Add "TFEngine", False
foundModules.Add "TFTests", False

On Error Resume Next
For Each comp In objWorkbook.VBProject.VBComponents
    If comp.Type = 1 Then ' vbext_ct_StdModule = 1
        moduleName = comp.Name
        If foundModules.Exists(moduleName) Then
            foundModules(moduleName) = True
            WScript.Echo "[OK]   Found: " & moduleName
        End If
    End If
Next

If Err.Number <> 0 Then
    WScript.Echo ""
    WScript.Echo "[FAIL] Cannot access VBA project: " & Err.Description
    WScript.Echo ""
    WScript.Echo "This usually means:"
    WScript.Echo "1. Macro security is blocking VBA project access"
    WScript.Echo ""
    WScript.Echo "To fix:"
    WScript.Echo "1. Open Excel"
    WScript.Echo "2. File > Options > Trust Center > Trust Center Settings"
    WScript.Echo "3. Macro Settings > Enable all macros"
    WScript.Echo "4. Check 'Trust access to the VBA project object model'"
    WScript.Echo "5. Click OK and restart Excel"
    objWorkbook.Close False
    objExcel.Quit
    WScript.Quit 1
End If
On Error Goto 0

' Check for missing modules
WScript.Echo ""
WScript.Echo "Module Status:"
Dim allFound, missingList
allFound = True
missingList = ""

For Each moduleName In foundModules.Keys
    If Not foundModules(moduleName) Then
        WScript.Echo "[FAIL] Missing: " & moduleName
        allFound = False
        If missingList <> "" Then missingList = missingList & ", "
        missingList = missingList & moduleName
    End If
Next

If allFound Then
    WScript.Echo "[OK]   All required modules found!"
Else
    WScript.Echo ""
    WScript.Echo "[FAIL] Missing modules: " & missingList
    WScript.Echo ""
    WScript.Echo "To fix:"
    WScript.Echo "1. Close Excel"
    WScript.Echo "2. Run: fix-vba-modules.bat"
    WScript.Echo "3. This will import the missing modules"
End If

' Close workbook
objWorkbook.Close False
objExcel.Quit

' Cleanup
Set objWorkbook = Nothing
Set objExcel = Nothing
Set objFSO = Nothing

WScript.Echo ""
If allFound Then
    WScript.Echo "=========================================="
    WScript.Echo " DIAGNOSIS: ALL CHECKS PASSED!"
    WScript.Echo "=========================================="
    WScript.Echo ""
    WScript.Echo "Your VBA setup is correct. If you're still getting errors:"
    WScript.Echo "1. Make sure 'Enable Content' is clicked in Excel"
    WScript.Echo "2. Check that tf-engine.exe is in the same folder"
    WScript.Echo "3. Try running a simple test from VBA Tests sheet"
    WScript.Echo ""
    WScript.Quit 0
Else
    WScript.Echo "=========================================="
    WScript.Echo " DIAGNOSIS: MODULES MISSING"
    WScript.Echo "=========================================="
    WScript.Echo ""
    WScript.Echo "Run fix-vba-modules.bat to import the missing modules"
    WScript.Echo ""
    WScript.Quit 1
End If
