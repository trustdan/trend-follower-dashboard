' ===========================================================================
' create-ui-worksheets.vbs - Creates All UI Worksheets for Trading Platform
' ===========================================================================
' Purpose: Generates 5 production worksheets with full UI
' Used by: 1-setup-all.bat during workbook generation
' Created: M22 - Automated UI Generation
' ===========================================================================

Option Explicit

' Get workbook path from command line or use default
Dim objFSO, objExcel, objWorkbook, strWorkbookPath

Set objFSO = CreateObject("Scripting.FileSystemObject")

' Determine workbook path
If WScript.Arguments.Count > 0 Then
    strWorkbookPath = WScript.Arguments(0)
Else
    ' Default: TradingPlatform.xlsm in current directory
    strWorkbookPath = objFSO.GetAbsolutePathName("TradingPlatform.xlsm")
End If

' Verify workbook exists
If Not objFSO.FileExists(strWorkbookPath) Then
    WScript.Echo "ERROR: Workbook not found: " & strWorkbookPath
    WScript.Quit 1
End If

' Open Excel
WScript.Echo "Opening workbook: " & strWorkbookPath
Set objExcel = CreateObject("Excel.Application")
objExcel.Visible = False
objExcel.DisplayAlerts = False

' Open workbook
Set objWorkbook = objExcel.Workbooks.Open(strWorkbookPath)

' Create all worksheets
WScript.Echo ""
WScript.Echo "Creating UI worksheets..."
WScript.Echo ""

Call CreateDashboard(objWorkbook)
Call CreatePositionSizing(objWorkbook)
Call CreateChecklist(objWorkbook)
Call CreateHeatCheck(objWorkbook)
Call CreateTradeEntry(objWorkbook)

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

WScript.Echo "UI worksheets created successfully!"
WScript.Quit 0

' ===========================================================================
' Worksheet Creation Functions
' ===========================================================================

' ---------------------------------------------------------------------------
' Dashboard Worksheet
' ---------------------------------------------------------------------------
Sub CreateDashboard(wb)
    WScript.Echo "  [1/5] Creating Dashboard worksheet..."

    Dim ws
    Set ws = wb.Worksheets.Add(, wb.Worksheets(wb.Worksheets.Count))
    ws.Name = "Dashboard"
    ws.Tab.Color = RGB(0, 102, 204)  ' Blue

    ' Header
    ws.Range("A1").Value = "Trading Platform Dashboard"
    ws.Range("A1").Font.Bold = True
    ws.Range("A1").Font.Size = 16
    ws.Range("A1").Font.Color = RGB(0, 0, 139)

    ' Portfolio Status Section
    ws.Range("A3").Value = "Portfolio Status"
    ws.Range("A3").Font.Bold = True
    ws.Range("A3").Font.Size = 12

    ws.Range("A4").Value = "Current Equity:"
    ws.Range("A5").Value = "Portfolio Heat:"
    ws.Range("A6").Value = "Portfolio Cap:"
    ws.Range("A7").Value = "Heat %:"

    ' Placeholder formulas (will be updated when VBA functions are ready)
    ws.Range("B4").Value = "[Formula: Equity]"
    ws.Range("B5").Value = "[Formula: Current Heat]"
    ws.Range("B6").Value = "[Formula: Heat Cap]"
    ws.Range("B7").Value = "[Formula: Heat %]"

    ' Format currency cells
    ws.Range("B4:B6").NumberFormat = "$#,##0.00"
    ws.Range("B7").NumberFormat = "0.0%"

    ' Today's Candidates Section
    ws.Range("A10").Value = "Today's Candidates"
    ws.Range("A10").Font.Bold = True
    ws.Range("A10").Font.Size = 12

    ws.Range("A11").Value = "[Will show candidates from FINVIZ preset]"

    ' Quick Actions Section
    ws.Range("A16").Value = "Quick Actions"
    ws.Range("A16").Font.Bold = True
    ws.Range("A16").Font.Size = 12

    ' Add navigation buttons
    Dim btn
    Set btn = ws.Buttons.Add(ws.Range("B17").Left, ws.Range("B17").Top, 150, 30)
    btn.Text = "Refresh Dashboard"
    btn.OnAction = "TFHelpers.RefreshDashboard"

    Set btn = ws.Buttons.Add(ws.Range("B18").Left, ws.Range("B18").Top, 150, 30)
    btn.Text = "Position Sizing"
    btn.OnAction = "TFHelpers.GotoPositionSizing"

    Set btn = ws.Buttons.Add(ws.Range("B19").Left, ws.Range("B19").Top, 150, 30)
    btn.Text = "Checklist"
    btn.OnAction = "TFHelpers.GotoChecklist"

    Set btn = ws.Buttons.Add(ws.Range("B20").Left, ws.Range("B20").Top, 150, 30)
    btn.Text = "Heat Check"
    btn.OnAction = "TFHelpers.GotoHeatCheck"

    Set btn = ws.Buttons.Add(ws.Range("B21").Left, ws.Range("B21").Top, 150, 30)
    btn.Text = "Trade Entry"
    btn.OnAction = "TFHelpers.GotoTradeEntry"

    ' Set column widths
    ws.Columns("A").ColumnWidth = 18
    ws.Columns("B").ColumnWidth = 20
    ws.Columns("C").ColumnWidth = 20

    WScript.Echo "      Dashboard created"
End Sub

' ---------------------------------------------------------------------------
' Position Sizing Worksheet
' ---------------------------------------------------------------------------
Sub CreatePositionSizing(wb)
    WScript.Echo "  [2/5] Creating Position Sizing worksheet..."

    Dim ws
    Set ws = wb.Worksheets.Add(, wb.Worksheets(wb.Worksheets.Count))
    ws.Name = "Position Sizing"
    ws.Tab.Color = RGB(0, 153, 76)  ' Green

    ' Header
    ws.Range("A1").Value = "Position Sizing Calculator"
    ws.Range("A1").Font.Bold = True
    ws.Range("A1").Font.Size = 14
    ws.Range("A1").Font.Color = RGB(0, 102, 0)

    ' Inputs Section
    ws.Range("A3").Value = "Inputs"
    ws.Range("A3").Font.Bold = True
    ws.Range("A3").Font.Size = 12

    ws.Range("A4").Value = "Ticker:"
    ws.Range("A5").Value = "Entry Price:"
    ws.Range("A6").Value = "ATR (N):"
    ws.Range("A7").Value = "K Multiple:"
    ws.Range("A8").Value = "Method:"

    ' Optional inputs
    ws.Range("A9").Value = "Equity Override:"
    ws.Range("A10").Value = "Risk % Override:"
    ws.Range("A11").Value = "Delta (options):"
    ws.Range("A12").Value = "Max Loss (options):"

    ' Add dropdown for method
    On Error Resume Next
    ws.Range("B8").Validation.Delete
    On Error Goto 0
    ws.Range("B8").Validation.Add 3, 1, 1, "stock,opt-delta-atr,opt-maxloss"
    ws.Range("B8").Validation.IgnoreBlank = True
    ws.Range("B8").Validation.InCellDropdown = True

    ' Format input cells
    ws.Range("B5:B7").NumberFormat = "0.00"
    ws.Range("B9").NumberFormat = "$#,##0.00"
    ws.Range("B10").NumberFormat = "0.0%"
    ws.Range("B11:B12").NumberFormat = "0.00"

    ' Add buttons
    Dim btn
    Set btn = ws.Buttons.Add(ws.Range("B13").Left, ws.Range("B13").Top, 100, 30)
    btn.Text = "Calculate"
    btn.OnAction = "TFEngine.CalculatePositionSize"

    Set btn = ws.Buttons.Add(ws.Range("C13").Left, ws.Range("C13").Top, 100, 30)
    btn.Text = "Clear"
    btn.OnAction = "TFEngine.ClearPositionSizing"

    ' Results Section
    ws.Range("A15").Value = "Results"
    ws.Range("A15").Font.Bold = True
    ws.Range("A15").Font.Size = 12

    ws.Range("A16").Value = "Risk Dollars (R):"
    ws.Range("A17").Value = "Stop Distance:"
    ws.Range("A18").Value = "Initial Stop:"
    ws.Range("A19").Value = "Shares:"
    ws.Range("A20").Value = "Contracts:"
    ws.Range("A21").Value = "Actual Risk:"
    ws.Range("A22").Value = "Status:"

    ' Format result cells
    ws.Range("B16").NumberFormat = "$#,##0.00"
    ws.Range("B17").NumberFormat = "0.00"
    ws.Range("B18").NumberFormat = "$#,##0.00"
    ws.Range("B19:B20").NumberFormat = "#,##0"
    ws.Range("B21").NumberFormat = "$#,##0.00"

    ' Highlight results area
    ws.Range("A15:B22").Borders.LineStyle = 1
    ws.Range("A15:B15").Interior.Color = RGB(220, 230, 241)

    ' Set column widths
    ws.Columns("A").ColumnWidth = 22
    ws.Columns("B").ColumnWidth = 18
    ws.Columns("C").ColumnWidth = 15

    WScript.Echo "      Position Sizing created"
End Sub

' ---------------------------------------------------------------------------
' Checklist Worksheet
' ---------------------------------------------------------------------------
Sub CreateChecklist(wb)
    WScript.Echo "  [3/5] Creating Checklist worksheet..."

    Dim ws
    Set ws = wb.Worksheets.Add(, wb.Worksheets(wb.Worksheets.Count))
    ws.Name = "Checklist"
    ws.Tab.Color = RGB(255, 192, 0)  ' Orange

    ' Header
    ws.Range("A1").Value = "Entry Checklist Evaluation"
    ws.Range("A1").Font.Bold = True
    ws.Range("A1").Font.Size = 14
    ws.Range("A1").Font.Color = RGB(204, 102, 0)

    ' Ticker input
    ws.Range("A3").Value = "Ticker:"

    ' Checklist items (using checkboxes)
    ws.Range("A5").Value = "Checklist Items:"
    ws.Range("A5").Font.Bold = True

    ' Create checkboxes (positioned manually)
    Dim chk
    Set chk = ws.OLEObjects.Add("Forms.CheckBox.1", False, False, _
        ws.Range("A6").Left, ws.Range("A6").Top, 350, 20)
    chk.Name = "chk_from_preset"
    chk.Object.Caption = "1. Ticker from today's FINVIZ preset"
    chk.Object.Value = False

    Set chk = ws.OLEObjects.Add("Forms.CheckBox.1", False, False, _
        ws.Range("A7").Left, ws.Range("A7").Top, 350, 20)
    chk.Name = "chk_trend_pass"
    chk.Object.Caption = "2. Trend alignment confirmed"
    chk.Object.Value = False

    Set chk = ws.OLEObjects.Add("Forms.CheckBox.1", False, False, _
        ws.Range("A8").Left, ws.Range("A8").Top, 350, 20)
    chk.Name = "chk_liquidity_pass"
    chk.Object.Caption = "3. Adequate volume and spread"
    chk.Object.Value = False

    Set chk = ws.OLEObjects.Add("Forms.CheckBox.1", False, False, _
        ws.Range("A9").Left, ws.Range("A9").Top, 350, 20)
    chk.Name = "chk_tv_confirm"
    chk.Object.Caption = "4. TradingView setup confirmation"
    chk.Object.Value = False

    Set chk = ws.OLEObjects.Add("Forms.CheckBox.1", False, False, _
        ws.Range("A10").Left, ws.Range("A10").Top, 350, 20)
    chk.Name = "chk_earnings_ok"
    chk.Object.Caption = "5. No earnings in next 7 days"
    chk.Object.Value = False

    Set chk = ws.OLEObjects.Add("Forms.CheckBox.1", False, False, _
        ws.Range("A11").Left, ws.Range("A11").Top, 350, 20)
    chk.Name = "chk_journal_ok"
    chk.Object.Caption = "6. Trade thesis documented in journal"
    chk.Object.Value = False

    ' Add buttons
    Dim btn
    Set btn = ws.Buttons.Add(ws.Range("B12").Left, ws.Range("B12").Top, 100, 30)
    btn.Text = "Evaluate"
    btn.OnAction = "TFEngine.EvaluateChecklist"

    Set btn = ws.Buttons.Add(ws.Range("C12").Left, ws.Range("C12").Top, 100, 30)
    btn.Text = "Clear"
    btn.OnAction = "TFEngine.ClearChecklist"

    ' Results Section
    ws.Range("A15").Value = "Results"
    ws.Range("A15").Font.Bold = True
    ws.Range("A15").Font.Size = 12

    ws.Range("A16").Value = "Banner:"
    ws.Range("A17").Value = "Missing Items:"
    ws.Range("A18").Value = "Missing:"
    ws.Range("A21").Value = "Allow Save:"
    ws.Range("A22").Value = "Evaluation Time:"
    ws.Range("A23").Value = "Status:"

    ' Format result area
    ws.Range("B16").Interior.Color = RGB(220, 220, 220)
    ws.Range("B16").Font.Bold = True
    ws.Range("B16").Font.Size = 14
    ws.Range("B16").HorizontalAlignment = -4108  ' xlCenter

    ' Highlight results area
    ws.Range("A15:B23").Borders.LineStyle = 1
    ws.Range("A15:B15").Interior.Color = RGB(255, 230, 153)

    ' Set column widths
    ws.Columns("A").ColumnWidth = 22
    ws.Columns("B").ColumnWidth = 25
    ws.Columns("C").ColumnWidth = 15

    WScript.Echo "      Checklist created"
End Sub

' ---------------------------------------------------------------------------
' Heat Check Worksheet
' ---------------------------------------------------------------------------
Sub CreateHeatCheck(wb)
    WScript.Echo "  [4/5] Creating Heat Check worksheet..."

    Dim ws
    Set ws = wb.Worksheets.Add(, wb.Worksheets(wb.Worksheets.Count))
    ws.Name = "Heat Check"
    ws.Tab.Color = RGB(255, 0, 0)  ' Red

    ' Header
    ws.Range("A1").Value = "Portfolio Heat Management"
    ws.Range("A1").Font.Bold = True
    ws.Range("A1").Font.Size = 14
    ws.Range("A1").Font.Color = RGB(139, 0, 0)

    ' Inputs Section
    ws.Range("A3").Value = "Ticker:"
    ws.Range("A4").Value = "Risk Amount ($):"
    ws.Range("A5").Value = "Bucket:"

    ' Add dropdown for bucket
    On Error Resume Next
    ws.Range("B5").Validation.Delete
    On Error Goto 0
    ws.Range("B5").Validation.Add 3, 1, 1, "Tech/Comm,Finance,Healthcare,Consumer,Energy,Industrials"
    ws.Range("B5").Validation.IgnoreBlank = True
    ws.Range("B5").Validation.InCellDropdown = True

    ' Format input cells
    ws.Range("B4").NumberFormat = "$#,##0.00"

    ' Add buttons
    Dim btn
    Set btn = ws.Buttons.Add(ws.Range("B7").Left, ws.Range("B7").Top, 100, 30)
    btn.Text = "Check Heat"
    btn.OnAction = "TFEngine.CheckHeat"

    Set btn = ws.Buttons.Add(ws.Range("C7").Left, ws.Range("C7").Top, 100, 30)
    btn.Text = "Clear"
    btn.OnAction = "TFEngine.ClearHeatCheck"

    ' Portfolio Heat Section
    ws.Range("A9").Value = "Portfolio Heat"
    ws.Range("A9").Font.Bold = True
    ws.Range("A9").Font.Size = 12

    ws.Range("A10").Value = "Current Heat:"
    ws.Range("A11").Value = "New Heat:"
    ws.Range("A12").Value = "Heat %:"
    ws.Range("A13").Value = "Cap:"
    ws.Range("A14").Value = "Exceeded:"
    ws.Range("A15").Value = "Overage:"

    ' Format portfolio results
    ws.Range("B10:B11").NumberFormat = "$#,##0.00"
    ws.Range("B12").NumberFormat = "0.0%"
    ws.Range("B13").NumberFormat = "$#,##0.00"
    ws.Range("B15").NumberFormat = "$#,##0.00"

    ' Bucket Heat Section
    ws.Range("A17").Value = "Bucket Heat"
    ws.Range("A17").Font.Bold = True
    ws.Range("A17").Font.Size = 12

    ws.Range("A18").Value = "Current Heat:"
    ws.Range("A19").Value = "New Heat:"
    ws.Range("A20").Value = "Heat %:"
    ws.Range("A21").Value = "Cap:"
    ws.Range("A22").Value = "Exceeded:"
    ws.Range("A23").Value = "Overage:"

    ' Format bucket results
    ws.Range("B18:B19").NumberFormat = "$#,##0.00"
    ws.Range("B20").NumberFormat = "0.0%"
    ws.Range("B21").NumberFormat = "$#,##0.00"
    ws.Range("B23").NumberFormat = "$#,##0.00"

    ' Highlight results areas
    ws.Range("A9:B15").Borders.LineStyle = 1
    ws.Range("A9:B9").Interior.Color = RGB(255, 230, 153)

    ws.Range("A17:B23").Borders.LineStyle = 1
    ws.Range("A17:B17").Interior.Color = RGB(255, 230, 153)

    ' Set column widths
    ws.Columns("A").ColumnWidth = 22
    ws.Columns("B").ColumnWidth = 18
    ws.Columns("C").ColumnWidth = 15

    WScript.Echo "      Heat Check created"
End Sub

' ---------------------------------------------------------------------------
' Trade Entry Worksheet
' ---------------------------------------------------------------------------
Sub CreateTradeEntry(wb)
    WScript.Echo "  [5/5] Creating Trade Entry worksheet..."

    Dim ws
    Set ws = wb.Worksheets.Add(, wb.Worksheets(wb.Worksheets.Count))
    ws.Name = "Trade Entry"
    ws.Tab.Color = RGB(128, 0, 128)  ' Purple

    ' Header
    ws.Range("A1").Value = "Trade Decision Entry (5 Hard Gates)"
    ws.Range("A1").Font.Bold = True
    ws.Range("A1").Font.Size = 14
    ws.Range("A1").Font.Color = RGB(75, 0, 130)

    ' Trade Details Section
    ws.Range("A3").Value = "Trade Details"
    ws.Range("A3").Font.Bold = True
    ws.Range("A3").Font.Size = 12

    ws.Range("A4").Value = "Ticker:"
    ws.Range("A5").Value = "Entry Price:"
    ws.Range("A6").Value = "ATR:"
    ws.Range("A7").Value = "Method:"
    ws.Range("A8").Value = "Banner Status:"
    ws.Range("A9").Value = "Delta (options):"
    ws.Range("A10").Value = "Max Loss (options):"
    ws.Range("A11").Value = "Bucket:"
    ws.Range("A12").Value = "Preset:"

    ' Add dropdowns
    On Error Resume Next
    ws.Range("B7").Validation.Delete
    ws.Range("B8").Validation.Delete
    ws.Range("B11").Validation.Delete
    On Error Goto 0

    ws.Range("B7").Validation.Add 3, 1, 1, "stock,opt-delta-atr,opt-maxloss"
    ws.Range("B7").Validation.IgnoreBlank = True
    ws.Range("B7").Validation.InCellDropdown = True

    ws.Range("B8").Validation.Add 3, 1, 1, "GREEN,YELLOW,RED"
    ws.Range("B8").Validation.IgnoreBlank = True
    ws.Range("B8").Validation.InCellDropdown = True

    ws.Range("B11").Validation.Add 3, 1, 1, "Tech/Comm,Finance,Healthcare,Consumer,Energy,Industrials"
    ws.Range("B11").Validation.IgnoreBlank = True
    ws.Range("B11").Validation.InCellDropdown = True

    ' Format input cells
    ws.Range("B5").NumberFormat = "$#,##0.00"
    ws.Range("B6").NumberFormat = "0.00"
    ws.Range("B9:B10").NumberFormat = "0.00"

    ' Add action buttons
    Dim btn
    Set btn = ws.Buttons.Add(ws.Range("B14").Left, ws.Range("B14").Top, 120, 30)
    btn.Text = "Save GO"
    btn.OnAction = "TFEngine.SaveDecisionGO"

    Set btn = ws.Buttons.Add(ws.Range("C14").Left, ws.Range("C14").Top, 120, 30)
    btn.Text = "Save NO-GO"
    btn.OnAction = "TFEngine.SaveDecisionNOGO"

    Set btn = ws.Buttons.Add(ws.Range("D14").Left, ws.Range("D14").Top, 100, 30)
    btn.Text = "Clear"
    btn.OnAction = "TFEngine.ClearTradeEntry"

    ' Gate Status Section
    ws.Range("A17").Value = "Gate Status"
    ws.Range("A17").Font.Bold = True
    ws.Range("A17").Font.Size = 12

    ws.Range("A18").Value = "Gate 1 - Banner GREEN:"
    ws.Range("A19").Value = "Gate 2 - In Candidates:"
    ws.Range("A20").Value = "Gate 3 - Impulse Brake:"
    ws.Range("A21").Value = "Gate 4 - Cooldown:"
    ws.Range("A22").Value = "Gate 5 - Heat Caps:"

    ' Results Section
    ws.Range("A24").Value = "Results"
    ws.Range("A24").Font.Bold = True
    ws.Range("A24").Font.Size = 12

    ws.Range("A25").Value = "Decision Saved:"
    ws.Range("A26").Value = "Decision ID:"
    ws.Range("A27").Value = "Rejection Reason:"
    ws.Range("A30").Value = "Status:"

    ' Highlight gate status area
    ws.Range("A17:B22").Borders.LineStyle = 1
    ws.Range("A17:B17").Interior.Color = RGB(230, 230, 250)

    ' Highlight results area
    ws.Range("A24:B30").Borders.LineStyle = 1
    ws.Range("A24:B24").Interior.Color = RGB(230, 230, 250)

    ' Set column widths
    ws.Columns("A").ColumnWidth = 24
    ws.Columns("B").ColumnWidth = 18
    ws.Columns("C").ColumnWidth = 15
    ws.Columns("D").ColumnWidth = 15

    WScript.Echo "      Trade Entry created"
End Sub

' ===========================================================================
' End of create-ui-worksheets.vbs
' ===========================================================================
