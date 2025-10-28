' ===========================================================================
' create-workbook-manual-ui.vbs - Create Workbook with Simple UI
' ===========================================================================
' Purpose: Creates workbook with text-based UI (no fragile OLE controls)
' Strategy: Use simple cells + buttons instead of checkboxes
' ===========================================================================

Option Explicit

' Get workbook path from command line or use default
Dim objFSO, objExcel, objWorkbook, strWorkbookPath, strVBADir

Set objFSO = CreateObject("Scripting.FileSystemObject")

' Determine paths
If WScript.Arguments.Count > 0 Then
    strWorkbookPath = WScript.Arguments(0)
Else
    strWorkbookPath = objFSO.GetAbsolutePathName("TradingPlatform.xlsm")
End If

' Get VBA directory (one level up from windows/)
Dim strScriptDir
strScriptDir = objFSO.GetParentFolderName(WScript.ScriptFullName)
strVBADir = objFSO.GetParentFolderName(strScriptDir) & "\excel\vba\"

' Verify VBA directory exists
If Not objFSO.FolderExists(strVBADir) Then
    WScript.Echo "ERROR: VBA directory not found: " & strVBADir
    WScript.Quit 1
End If

WScript.Echo "=========================================="
WScript.Echo "Trading Platform Workbook Generator"
WScript.Echo "=========================================="
WScript.Echo ""
WScript.Echo "Creating: " & strWorkbookPath
WScript.Echo "VBA Path: " & strVBADir
WScript.Echo ""

' Create Excel application
On Error Resume Next
Set objExcel = CreateObject("Excel.Application")
If Err.Number <> 0 Then
    WScript.Echo "ERROR: Could not create Excel.Application"
    WScript.Echo "Error: " & Err.Description
    WScript.Quit 1
End If
On Error Goto 0

objExcel.Visible = False
objExcel.DisplayAlerts = False

' Create new workbook
WScript.Echo "[1/4] Creating new workbook..."
Set objWorkbook = objExcel.Workbooks.Add

' Import VBA modules FIRST (before creating UI)
WScript.Echo "[2/4] Importing VBA modules..."
On Error Resume Next

WScript.Echo "  - Importing TFTypes.bas..."
objWorkbook.VBProject.VBComponents.Import strVBADir & "TFTypes.bas"
If Err.Number <> 0 Then
    WScript.Echo "    WARNING: Could not import TFTypes.bas: " & Err.Description
    WScript.Echo "    Make sure 'Trust access to VBA project' is enabled in Excel"
    Err.Clear
End If

WScript.Echo "  - Importing TFHelpers.bas..."
objWorkbook.VBProject.VBComponents.Import strVBADir & "TFHelpers.bas"
If Err.Number <> 0 Then
    WScript.Echo "    WARNING: Could not import TFHelpers.bas: " & Err.Description
    Err.Clear
End If

WScript.Echo "  - Importing TFEngine.bas..."
objWorkbook.VBProject.VBComponents.Import strVBADir & "TFEngine.bas"
If Err.Number <> 0 Then
    WScript.Echo "    WARNING: Could not import TFEngine.bas: " & Err.Description
    Err.Clear
End If

WScript.Echo "  - Importing TFTests.bas..."
objWorkbook.VBProject.VBComponents.Import strVBADir & "TFTests.bas"
If Err.Number <> 0 Then
    WScript.Echo "    WARNING: Could not import TFTests.bas: " & Err.Description
    Err.Clear
End If

On Error Goto 0

' Create worksheets with simple UI
WScript.Echo "[3/4] Creating worksheets..."

' Rename first sheet to Setup
objWorkbook.Worksheets(1).Name = "Setup"
Call CreateSetupSheet(objWorkbook.Worksheets("Setup"))

' Create other sheets
Call CreateDashboardSheet(objWorkbook)
Call CreatePositionSizingSheet(objWorkbook)
Call CreateChecklistSheet(objWorkbook)
Call CreateHeatCheckSheet(objWorkbook)
Call CreateTradeEntrySheet(objWorkbook)
Call CreateVBATestsSheet(objWorkbook)

' Save workbook
WScript.Echo "[4/4] Saving workbook..."
On Error Resume Next
objWorkbook.SaveAs strWorkbookPath, 52  ' xlOpenXMLWorkbookMacroEnabled
If Err.Number <> 0 Then
    WScript.Echo "ERROR: Could not save workbook: " & Err.Description
    objWorkbook.Close False
    objExcel.Quit
    WScript.Quit 1
End If
On Error Goto 0

objWorkbook.Close
objExcel.Quit

WScript.Echo ""
WScript.Echo "=========================================="
WScript.Echo "SUCCESS!"
WScript.Echo "=========================================="
WScript.Echo ""
WScript.Echo "Workbook created: " & strWorkbookPath
WScript.Echo ""
WScript.Echo "Next steps:"
WScript.Echo "  1. Open TradingPlatform.xlsm"
WScript.Echo "  2. Enable macros when prompted"
WScript.Echo "  3. Go to 'VBA Tests' sheet"
WScript.Echo "  4. Click 'Run All Tests' button"
WScript.Echo ""

WScript.Quit 0

' ===========================================================================
' WORKSHEET CREATION FUNCTIONS
' ===========================================================================

Sub CreateSetupSheet(ws)
    WScript.Echo "  - Creating Setup sheet..."

    ws.Tab.Color = RGB(192, 192, 192)

    ' Header
    ws.Range("A1").Value = "Trading Engine v3 - Setup & Configuration"
    ws.Range("A1").Font.Bold = True
    ws.Range("A1").Font.Size = 16

    ' Configuration
    ws.Range("A3").Value = "Configuration"
    ws.Range("A3").Font.Bold = True
    ws.Range("A3").Font.Size = 12

    ws.Range("A4").Value = "Engine Path:"
    ws.Range("B4").Value = ".\tf-engine.exe"

    ws.Range("A5").Value = "Database Path:"
    ws.Range("B5").Value = ".\trading.db"

    ' Create named ranges
    On Error Resume Next
    ws.Parent.Names.Add "EnginePathSetting", ws.Range("B4")
    ws.Parent.Names.Add "DatabasePathSetting", ws.Range("B5")
    On Error Goto 0

    ' Instructions
    ws.Range("A8").Value = "Quick Start:"
    ws.Range("A8").Font.Bold = True
    ws.Range("A9").Value = "1. Make sure tf-engine.exe is in the same folder as this workbook"
    ws.Range("A10").Value = "2. Go to 'VBA Tests' sheet and click 'Run All Tests'"
    ws.Range("A11").Value = "3. If tests pass, start with the Dashboard sheet"

    ws.Columns("A").ColumnWidth = 20
    ws.Columns("B").ColumnWidth = 30
End Sub

Sub CreateDashboardSheet(wb)
    WScript.Echo "  - Creating Dashboard sheet..."

    Dim ws
    Set ws = wb.Worksheets.Add(, wb.Worksheets(wb.Worksheets.Count))
    ws.Name = "Dashboard"
    ws.Tab.Color = RGB(0, 102, 204)

    ' Header
    ws.Range("A1").Value = "Trading Platform Dashboard"
    ws.Range("A1").Font.Bold = True
    ws.Range("A1").Font.Size = 16
    ws.Range("A1").Font.Color = RGB(0, 0, 139)

    ' Portfolio Status
    ws.Range("A3").Value = "Portfolio Status"
    ws.Range("A3").Font.Bold = True
    ws.Range("A3").Font.Size = 12

    ws.Range("A4").Value = "Current Equity:"
    ws.Range("B4").Value = 10000
    ws.Range("B4").NumberFormat = "$#,##0.00"

    ws.Range("A5").Value = "Portfolio Heat Cap:"
    ws.Range("B5").Value = "4.0%"

    ws.Range("A6").Value = "Bucket Heat Cap:"
    ws.Range("B6").Value = "1.5%"

    ' Today's Candidates
    ws.Range("A9").Value = "Today's Candidates"
    ws.Range("A9").Font.Bold = True
    ws.Range("A9").Font.Size = 12

    ws.Range("A10").Value = "[Import candidates from FINVIZ to populate this list]"

    ws.Columns("A").ColumnWidth = 22
    ws.Columns("B").ColumnWidth = 20
End Sub

Sub CreatePositionSizingSheet(wb)
    WScript.Echo "  - Creating Position Sizing sheet..."

    Dim ws
    Set ws = wb.Worksheets.Add(, wb.Worksheets(wb.Worksheets.Count))
    ws.Name = "Position Sizing"
    ws.Tab.Color = RGB(0, 153, 76)

    ' Header
    ws.Range("A1").Value = "Position Sizing Calculator"
    ws.Range("A1").Font.Bold = True
    ws.Range("A1").Font.Size = 14
    ws.Range("A1").Font.Color = RGB(0, 102, 0)

    ' Inputs
    ws.Range("A3").Value = "INPUTS"
    ws.Range("A3").Font.Bold = True
    ws.Range("A3").Font.Size = 12

    ws.Range("A4").Value = "Ticker:"
    ws.Range("A5").Value = "Entry Price:"
    ws.Range("A6").Value = "ATR (N):"
    ws.Range("A7").Value = "K Multiple:"
    ws.Range("A8").Value = "Method:"

    ws.Range("B7").Value = 2
    ws.Range("B8").Value = "stock"

    ' Add dropdown for method
    On Error Resume Next
    ws.Range("B8").Validation.Delete
    ws.Range("B8").Validation.Add 3, 1, 1, "stock,opt-delta-atr,opt-maxloss"
    ws.Range("B8").Validation.InCellDropdown = True
    On Error Goto 0

    ' Optional inputs
    ws.Range("A10").Value = "Optional Overrides:"
    ws.Range("A10").Font.Bold = True
    ws.Range("A11").Value = "Equity Override:"
    ws.Range("A12").Value = "Risk % Override:"
    ws.Range("A13").Value = "Delta (options):"
    ws.Range("A14").Value = "Max Loss (options):"

    ' Add Calculate button
    On Error Resume Next
    Dim btn
    Set btn = ws.Buttons.Add(ws.Range("B16").Left, ws.Range("B16").Top, 100, 25)
    btn.Text = "Calculate"
    btn.OnAction = "'Position Sizing'!CalculateSize"
    On Error Goto 0

    ' Results
    ws.Range("A19").Value = "RESULTS"
    ws.Range("A19").Font.Bold = True
    ws.Range("A19").Font.Size = 12

    ws.Range("A20").Value = "Risk Dollars (R):"
    ws.Range("A21").Value = "Stop Distance:"
    ws.Range("A22").Value = "Initial Stop:"
    ws.Range("A23").Value = "Shares:"
    ws.Range("A24").Value = "Contracts:"
    ws.Range("A25").Value = "Status:"

    ' Format cells
    ws.Range("B5:B6").NumberFormat = "0.00"
    ws.Range("B11").NumberFormat = "$#,##0.00"
    ws.Range("B12").NumberFormat = "0.00%"
    ws.Range("B20").NumberFormat = "$#,##0.00"
    ws.Range("B21").NumberFormat = "0.00"
    ws.Range("B22").NumberFormat = "$#,##0.00"
    ws.Range("B23:B24").NumberFormat = "#,##0"

    ws.Columns("A").ColumnWidth = 22
    ws.Columns("B").ColumnWidth = 18
End Sub

Sub CreateChecklistSheet(wb)
    WScript.Echo "  - Creating Checklist sheet (simplified)..."

    Dim ws
    Set ws = wb.Worksheets.Add(, wb.Worksheets(wb.Worksheets.Count))
    ws.Name = "Checklist"
    ws.Tab.Color = RGB(255, 192, 0)

    ' Header
    ws.Range("A1").Value = "Entry Checklist Evaluation"
    ws.Range("A1").Font.Bold = True
    ws.Range("A1").Font.Size = 14
    ws.Range("A1").Font.Color = RGB(204, 102, 0)

    ' Ticker
    ws.Range("A3").Value = "Ticker:"

    ' Checklist - using TRUE/FALSE cells instead of checkboxes
    ws.Range("A5").Value = "CHECKLIST ITEMS (Enter TRUE/FALSE in column B):"
    ws.Range("A5").Font.Bold = True
    ws.Range("A5").Font.Size = 11

    ws.Range("A6").Value = "1. Ticker from today's FINVIZ preset:"
    ws.Range("A7").Value = "2. Trend alignment confirmed:"
    ws.Range("A8").Value = "3. Adequate volume and spread:"
    ws.Range("A9").Value = "4. TradingView setup confirmation:"
    ws.Range("A10").Value = "5. No earnings in next 7 days:"
    ws.Range("A11").Value = "6. Trade thesis documented:"

    ' Pre-fill with FALSE
    ws.Range("B6:B11").Value = "FALSE"

    ' Add dropdowns for TRUE/FALSE
    Dim i
    For i = 6 To 11
        On Error Resume Next
        ws.Range("B" & i).Validation.Delete
        ws.Range("B" & i).Validation.Add 3, 1, 1, "TRUE,FALSE"
        ws.Range("B" & i).Validation.InCellDropdown = True
        On Error Goto 0
    Next

    ' Evaluate button
    On Error Resume Next
    Dim btn
    Set btn = ws.Buttons.Add(ws.Range("B13").Left, ws.Range("B13").Top, 100, 25)
    btn.Text = "Evaluate"
    btn.OnAction = "Checklist!EvaluateChecklist"
    On Error Goto 0

    ' Results
    ws.Range("A16").Value = "RESULTS"
    ws.Range("A16").Font.Bold = True
    ws.Range("A16").Font.Size = 12

    ws.Range("A17").Value = "Banner:"
    ws.Range("B17").Interior.Color = RGB(220, 220, 220)
    ws.Range("B17").Font.Bold = True
    ws.Range("B17").Font.Size = 14
    ws.Range("B17").HorizontalAlignment = -4108

    ws.Range("A18").Value = "Missing Items:"
    ws.Range("A19").Value = "Allow Save:"
    ws.Range("A20").Value = "Status:"

    ws.Columns("A").ColumnWidth = 35
    ws.Columns("B").ColumnWidth = 20
End Sub

Sub CreateHeatCheckSheet(wb)
    WScript.Echo "  - Creating Heat Check sheet..."

    Dim ws
    Set ws = wb.Worksheets.Add(, wb.Worksheets(wb.Worksheets.Count))
    ws.Name = "Heat Check"
    ws.Tab.Color = RGB(255, 0, 0)

    ' Header
    ws.Range("A1").Value = "Portfolio Heat Management"
    ws.Range("A1").Font.Bold = True
    ws.Range("A1").Font.Size = 14
    ws.Range("A1").Font.Color = RGB(139, 0, 0)

    ' Inputs
    ws.Range("A3").Value = "Ticker:"
    ws.Range("A4").Value = "Risk Amount ($):"
    ws.Range("A5").Value = "Bucket:"

    ' Bucket dropdown
    On Error Resume Next
    ws.Range("B5").Validation.Delete
    ws.Range("B5").Validation.Add 3, 1, 1, "Tech/Comm,Finance,Healthcare,Consumer,Energy,Industrials"
    ws.Range("B5").Validation.InCellDropdown = True
    On Error Goto 0

    ws.Range("B4").NumberFormat = "$#,##0.00"

    ' Check Heat button
    On Error Resume Next
    Dim btn
    Set btn = ws.Buttons.Add(ws.Range("B7").Left, ws.Range("B7").Top, 100, 25)
    btn.Text = "Check Heat"
    btn.OnAction = "'Heat Check'!CheckHeat"
    On Error Goto 0

    ' Portfolio Heat Results
    ws.Range("A10").Value = "PORTFOLIO HEAT"
    ws.Range("A10").Font.Bold = True
    ws.Range("A11").Value = "Current Heat:"
    ws.Range("A12").Value = "New Heat:"
    ws.Range("A13").Value = "Heat %:"
    ws.Range("A14").Value = "Cap Exceeded:"

    ' Bucket Heat Results
    ws.Range("A17").Value = "BUCKET HEAT"
    ws.Range("A17").Font.Bold = True
    ws.Range("A18").Value = "Current Heat:"
    ws.Range("A19").Value = "New Heat:"
    ws.Range("A20").Value = "Heat %:"
    ws.Range("A21").Value = "Cap Exceeded:"

    ws.Range("B11:B12").NumberFormat = "$#,##0.00"
    ws.Range("B13").NumberFormat = "0.0%"
    ws.Range("B18:B19").NumberFormat = "$#,##0.00"
    ws.Range("B20").NumberFormat = "0.0%"

    ws.Columns("A").ColumnWidth = 22
    ws.Columns("B").ColumnWidth = 18
End Sub

Sub CreateTradeEntrySheet(wb)
    WScript.Echo "  - Creating Trade Entry sheet..."

    Dim ws
    Set ws = wb.Worksheets.Add(, wb.Worksheets(wb.Worksheets.Count))
    ws.Name = "Trade Entry"
    ws.Tab.Color = RGB(128, 0, 128)

    ' Header
    ws.Range("A1").Value = "Trade Decision Entry (5 Hard Gates)"
    ws.Range("A1").Font.Bold = True
    ws.Range("A1").Font.Size = 14
    ws.Range("A1").Font.Color = RGB(75, 0, 130)

    ' Trade Details
    ws.Range("A3").Value = "TRADE DETAILS"
    ws.Range("A3").Font.Bold = True
    ws.Range("A4").Value = "Ticker:"
    ws.Range("A5").Value = "Entry Price:"
    ws.Range("A6").Value = "ATR:"
    ws.Range("A7").Value = "Method:"
    ws.Range("A8").Value = "Banner Status:"
    ws.Range("A9").Value = "Bucket:"
    ws.Range("A10").Value = "Preset:"

    ' Add dropdowns
    On Error Resume Next
    ws.Range("B7").Validation.Delete
    ws.Range("B7").Validation.Add 3, 1, 1, "stock,opt-delta-atr,opt-maxloss"
    ws.Range("B7").Validation.InCellDropdown = True

    ws.Range("B8").Validation.Delete
    ws.Range("B8").Validation.Add 3, 1, 1, "GREEN,YELLOW,RED"
    ws.Range("B8").Validation.InCellDropdown = True

    ws.Range("B9").Validation.Delete
    ws.Range("B9").Validation.Add 3, 1, 1, "Tech/Comm,Finance,Healthcare,Consumer,Energy,Industrials"
    ws.Range("B9").Validation.InCellDropdown = True
    On Error Goto 0

    ws.Range("B5:B6").NumberFormat = "0.00"

    ' Action buttons
    On Error Resume Next
    Dim btn
    Set btn = ws.Buttons.Add(ws.Range("B12").Left, ws.Range("B12").Top, 100, 25)
    btn.Text = "Save GO"
    btn.OnAction = "'Trade Entry'!SaveGO"

    Set btn = ws.Buttons.Add(ws.Range("C12").Left, ws.Range("C12").Top, 100, 25)
    btn.Text = "Save NO-GO"
    btn.OnAction = "'Trade Entry'!SaveNOGO"
    On Error Goto 0

    ' Gate Status
    ws.Range("A15").Value = "GATE STATUS"
    ws.Range("A15").Font.Bold = True
    ws.Range("A16").Value = "Gate 1 - Banner GREEN:"
    ws.Range("A17").Value = "Gate 2 - In Candidates:"
    ws.Range("A18").Value = "Gate 3 - Impulse Brake:"
    ws.Range("A19").Value = "Gate 4 - Cooldown:"
    ws.Range("A20").Value = "Gate 5 - Heat Caps:"

    ' Results
    ws.Range("A23").Value = "RESULTS"
    ws.Range("A23").Font.Bold = True
    ws.Range("A24").Value = "Decision Saved:"
    ws.Range("A25").Value = "Decision ID:"
    ws.Range("A26").Value = "Rejection Reason:"
    ws.Range("A27").Value = "Status:"

    ws.Columns("A").ColumnWidth = 26
    ws.Columns("B").ColumnWidth = 18
    ws.Columns("C").ColumnWidth = 15
End Sub

Sub CreateVBATestsSheet(wb)
    WScript.Echo "  - Creating VBA Tests sheet..."

    Dim ws
    Set ws = wb.Worksheets.Add(, wb.Worksheets(wb.Worksheets.Count))
    ws.Name = "VBA Tests"
    ws.Tab.Color = RGB(100, 100, 100)

    ' Header
    ws.Range("A1").Value = "VBA Unit Tests"
    ws.Range("A1").Font.Bold = True
    ws.Range("A1").Font.Size = 14

    ws.Range("A3").Value = "Click the button to run all VBA unit tests:"

    ' Add Run Tests button
    On Error Resume Next
    Dim btn
    Set btn = ws.Buttons.Add(ws.Range("B3").Left, ws.Range("B3").Top, 120, 25)
    btn.Text = "Run All Tests"
    btn.OnAction = "TFTests.RunAllTests"
    On Error Goto 0

    ' Test results header
    ws.Range("A6").Value = "Test Name"
    ws.Range("B6").Value = "Result"
    ws.Range("C6").Value = "Message"
    ws.Range("D6").Value = "Duration"

    ws.Range("A6:D6").Font.Bold = True
    ws.Range("A6:D6").Interior.Color = RGB(200, 200, 200)

    ws.Columns("A").ColumnWidth = 35
    ws.Columns("B").ColumnWidth = 10
    ws.Columns("C").ColumnWidth = 50
    ws.Columns("D").ColumnWidth = 12
End Sub
