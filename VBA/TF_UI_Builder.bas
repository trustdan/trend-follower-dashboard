Attribute VB_Name = "TF_UI_Builder"
' ============================================================================
' Module: TF_UI_Builder
' Purpose: Programmatically build the TradeEntry UI (labels, buttons, formatting)
' Run InitializeUI to create the full interactive interface
' ============================================================================

Option Explicit

' ----------------------------------------------------------------------------
' Sub: InitializeUI
' Main entry point - creates complete TradeEntry UI
' ----------------------------------------------------------------------------
Sub InitializeUI()
    Application.ScreenUpdating = False

    Call FormatTradeEntrySheet
    Call CreateButtons
    Call ApplyTheme
    Call SetupDataValidation

    Application.ScreenUpdating = True

    MsgBox "Trade Entry UI initialized successfully!", vbInformation, "UI Setup Complete"
End Sub

' ----------------------------------------------------------------------------
' Sub: FormatTradeEntrySheet
' Creates the layout: labels, input cells, checklist, banner
' ----------------------------------------------------------------------------
Sub FormatTradeEntrySheet()
    Dim ws As Worksheet
    Set ws = Sheets("TradeEntry")

    With ws
        .Cells.Clear  ' Start fresh

        ' === HEADER ===
        .Range("A1:F1").Merge
        .Range("A1").Value = "TRADE ENTRY DASHBOARD"
        .Range("A1").Font.Bold = True
        .Range("A1").Font.Size = 14
        .Range("A1").HorizontalAlignment = xlCenter

        ' === INPUT SECTION ===
        .Range("A3").Value = "Preset:"
        .Range("A4").Value = "Ticker:"
        .Range("A5").Value = "Sector:"
        .Range("A6").Value = "Bucket:"
        .Range("A8").Value = "Entry Price:"
        .Range("A9").Value = "ATR (N):"
        .Range("A10").Value = "Stop Multiple (K):"
        .Range("A11").Value = "Method:"

        ' Method options (option buttons would be better, but dropdowns work)
        .Range("A13").Value = "Delta (options):"
        .Range("A14").Value = "DTE (options):"
        .Range("A15").Value = "Max Loss (spreads):"

        ' Input cells (column B)
        .Range("B3").Name = "inp_Preset"
        .Range("B4").Name = "inp_Ticker"
        .Range("B5").Name = "inp_Sector"
        .Range("B6").Name = "inp_Bucket"
        .Range("B8").Name = "inp_Entry"
        .Range("B9").Name = "inp_N"
        .Range("B10").Name = "inp_K"
        .Range("B11").Name = "inp_Method"
        .Range("B13").Name = "inp_Delta"
        .Range("B14").Name = "inp_DTE"
        .Range("B15").Name = "inp_MaxLoss"

        ' Set default values
        .Range("B10").Value = 2  ' Default K = 2
        .Range("B11").Value = "Stock"

        ' === BANNER SECTION ===
        .Range("A17:F17").Merge
        .Range("A17").Name = "BannerCell"
        .Range("A17").Value = "Click EVALUATE to check trade"
        .Range("A17").HorizontalAlignment = xlCenter
        .Range("A17").VerticalAlignment = xlCenter
        .Range("A17").Font.Bold = True
        .Range("A17").Font.Size = 12
        .Range("A17").Interior.Color = RGB(200, 200, 200)  ' Gray default

        ' === CHECKLIST SECTION ===
        .Range("A19").Value = "PRE-TRADE CHECKLIST"
        .Range("A19").Font.Bold = True

        .Range("A20").Value = "From Preset?"
        .Range("A21").Value = "Trend Pass?"
        .Range("A22").Value = "Liquidity Pass?"
        .Range("A23").Value = "TradingView Confirm?"
        .Range("A24").Value = "Earnings OK?"
        .Range("A25").Value = "Journal OK?"

        ' Checkbox cells (column B)
        .Range("B20:B25").Value = False
        .Range("B20:B25").HorizontalAlignment = xlCenter

        ' === SIZING OUTPUT ===
        .Range("D3").Value = "Position Sizing:"
        .Range("D3").Font.Bold = True
        .Range("D4").Value = "R (dollars):"
        .Range("D5").Value = "Shares:"
        .Range("D6").Value = "Contracts:"
        .Range("D7").Value = "Initial Stop:"
        .Range("D8").Value = "Add Level 1:"
        .Range("D9").Value = "Add Level 2:"
        .Range("D10").Value = "Add Level 3:"

        .Range("E4").Name = "out_R"
        .Range("E5").Name = "out_Shares"
        .Range("E6").Name = "out_Contracts"
        .Range("E7").Name = "out_InitialStop"
        .Range("E8").Name = "out_Add1"
        .Range("E9").Name = "out_Add2"
        .Range("E10").Name = "out_Add3"

        ' === HEAT PREVIEW ===
        .Range("D12").Value = "Heat Preview:"
        .Range("D12").Font.Bold = True
        .Range("D13").Value = "Portfolio Heat:"
        .Range("D14").Value = "Bucket Heat:"

        .Range("E13").Name = "out_PortHeat"
        .Range("E14").Name = "out_BucketHeat"

        ' === COLUMN WIDTHS ===
        .Columns("A:A").ColumnWidth = 20
        .Columns("B:B").ColumnWidth = 15
        .Columns("C:C").ColumnWidth = 3  ' Spacer
        .Columns("D:D").ColumnWidth = 18
        .Columns("E:E").ColumnWidth = 12
        .Columns("F:F").ColumnWidth = 12

        ' === ROW HEIGHTS ===
        .Rows("1:1").RowHeight = 30
        .Rows("17:17").RowHeight = 40

    End With
End Sub

' ----------------------------------------------------------------------------
' Sub: CreateButtons
' Adds action buttons (Evaluate, Recalc, Save, Import)
' ----------------------------------------------------------------------------
Sub CreateButtons()
    Dim ws As Worksheet
    Dim btn As Button

    Set ws = Sheets("TradeEntry")

    ' Clear existing buttons
    On Error Resume Next
    ws.Buttons.Delete
    On Error GoTo 0

    ' Button 1: Evaluate
    Set btn = ws.Buttons.Add(300, 100, 80, 25)
    btn.OnAction = "TF_UI.EvaluateChecklist"
    btn.Text = "Evaluate"
    btn.Font.Bold = True

    ' Button 2: Recalc Sizing
    Set btn = ws.Buttons.Add(390, 100, 80, 25)
    btn.OnAction = "TF_UI.RecalcSizing"
    btn.Text = "Recalc"

    ' Button 3: Save Decision
    Set btn = ws.Buttons.Add(480, 100, 80, 25)
    btn.OnAction = "TF_UI.SaveDecision"
    btn.Text = "Save"
    btn.Font.Bold = True

    ' Button 4: Import Candidates
    Set btn = ws.Buttons.Add(300, 140, 120, 25)
    btn.OnAction = "TF_Presets.ImportCandidatesPrompt"
    btn.Text = "Import Candidates"

    ' Button 5: Open Preset URL
    Set btn = ws.Buttons.Add(430, 140, 100, 25)
    btn.OnAction = "TF_Presets.OpenPreset"
    btn.Text = "Open FINVIZ"
End Sub

' ----------------------------------------------------------------------------
' Sub: ApplyTheme
' Colors, fonts, borders
' ----------------------------------------------------------------------------
Sub ApplyTheme()
    Dim ws As Worksheet
    Set ws = Sheets("TradeEntry")

    With ws
        ' Header styling
        .Range("A1").Interior.Color = RGB(0, 112, 192)  ' Blue
        .Range("A1").Font.Color = RGB(255, 255, 255)    ' White text

        ' Input labels (column A)
        .Range("A3:A15").Font.Bold = True
        .Range("A3:A15").Interior.Color = RGB(217, 225, 242)  ' Light blue

        ' Checklist header
        .Range("A19").Interior.Color = RGB(255, 242, 204)  ' Light yellow

        ' Checklist items
        .Range("A20:B25").Interior.Color = RGB(255, 250, 240)  ' Cream
        .Range("B20:B25").Interior.Color = RGB(255, 255, 255)  ' White checkboxes

        ' Output section
        .Range("D3").Interior.Color = RGB(217, 225, 242)  ' Light blue
        .Range("D4:D10").Interior.Color = RGB(242, 242, 242)  ' Light gray
        .Range("E4:E10").Interior.Color = RGB(255, 255, 255)  ' White

        ' Heat section
        .Range("D12").Interior.Color = RGB(255, 242, 204)  ' Light yellow
        .Range("D13:D14").Interior.Color = RGB(242, 242, 242)  ' Light gray
        .Range("E13:E14").Interior.Color = RGB(255, 255, 255)  ' White

        ' Add borders
        .Range("A1:F25").Borders.LineStyle = xlContinuous
        .Range("A1:F25").Borders.Weight = xlThin
        .Range("A1:F25").Borders.Color = RGB(128, 128, 128)

        ' Thicker borders around sections
        .Range("A1:F1").BorderAround Weight:=xlMedium
        .Range("A17:F17").BorderAround Weight:=xlMedium
        .Range("A19:B25").BorderAround Weight:=xlMedium
        .Range("D3:E10").BorderAround Weight:=xlMedium

    End With
End Sub

' ----------------------------------------------------------------------------
' Sub: SetupDataValidation
' Add dropdowns to input cells
' ----------------------------------------------------------------------------
Sub SetupDataValidation()
    Dim ws As Worksheet
    Set ws = Sheets("TradeEntry")

    On Error Resume Next

    With ws
        ' Clear existing validation
        .Range("B3:B15").Validation.Delete

        ' Preset dropdown
        .Range("B3").Validation.Add Type:=xlValidateList, _
            Formula1:="=tblPresets[Name]"

        ' Ticker dropdown
        .Range("B4").Validation.Add Type:=xlValidateList, _
            Formula1:="=tblCandidates[Ticker]"

        ' Sector dropdown
        .Range("B5").Validation.Add Type:=xlValidateList, _
            Formula1:="Technology,Healthcare,Financials,Consumer,Industrials,Energy,Materials,Real Estate"

        ' Bucket dropdown
        .Range("B6").Validation.Add Type:=xlValidateList, _
            Formula1:="=tblBuckets[Bucket]"

        ' Method dropdown
        .Range("B11").Validation.Add Type:=xlValidateList, _
            Formula1:="Stock,Opt-DeltaATR,Opt-MaxLoss"

        ' Numeric validations
        .Range("B8:B10,B13:B15").Validation.Add Type:=xlValidateDecimal, _
            AlertStyle:=xlValidAlertStop, _
            Operator:=xlGreater, _
            Formula1:="0"

    End With

    On Error GoTo 0
End Sub
