Attribute VB_Name = "TF_UI_Builder"
Option Explicit

' ========================================
' TF_UI_Builder Module
' Automated TradeEntry sheet UI creation
' ========================================

Sub BuildTradeEntryUI()
    ' Complete UI builder - creates all labels, controls, and formatting
    Application.ScreenUpdating = False

    Dim ws As Worksheet
    Dim shp As Shape

    Set ws = Worksheets("TradeEntry")

    ' Delete all existing buttons and shapes first
    On Error Resume Next
    For Each shp In ws.Shapes
        shp.Delete
    Next shp
    On Error GoTo 0

    ' Clear existing content
    ws.Cells.Clear
    ws.Cells.ClearFormats
    ws.Cells.ClearComments

    ' Build UI sections
    Call CreateHeader(ws)
    Call CreateInputSection(ws)
    Call CreateChecklistSection(ws)
    Call CreateOutputSection(ws)
    Call CreateButtons(ws)
    Call CreateCheckboxes(ws)  ' Attempt to create checkboxes programmatically
    Call ApplyFormatting(ws)
    Call SetColumnWidths(ws)

    ' Bind controls
    Call TF_UI.BindControls

    Application.ScreenUpdating = True

    MsgBox "TradeEntry UI built successfully!" & vbCrLf & vbCrLf & _
           "NOTE: If checkboxes did not appear, add them manually:" & vbCrLf & _
           "Developer -> Insert -> Check Box (Form Control)" & vbCrLf & _
           "Link to cells C20, C21, C22, C23, C24, C25", vbInformation
End Sub

Private Sub CreateHeader(ws As Worksheet)
    ' Create title and banner
    With ws
        ' Title
        .Range("A1").Value = "TRADE ENTRY SYSTEM"
        .Range("A1").Font.Size = 16
        .Range("A1").Font.Bold = True
        .Range("A1:F1").Merge
        .Range("A1").HorizontalAlignment = xlCenter

        ' Banner (GO/NO-GO indicator)
        .Range("A2:F2").Merge
        .Range("A2").Value = "Click Evaluate to check entry criteria"
        .Range("A2").Interior.Color = RGB(240, 240, 240)
        .Range("A2").Font.Size = 14
        .Range("A2").Font.Bold = True
        .Range("A2").HorizontalAlignment = xlCenter
        .Range("A2").VerticalAlignment = xlCenter
        .Range("A2").RowHeight = 30

        ' Reason string
        .Range("A3:F3").Merge
        .Range("A3").Value = ""
        .Range("A3").Font.Italic = True
    End With
End Sub

Private Sub CreateInputSection(ws As Worksheet)
    ' Create input labels and cells
    With ws
        ' Row 5: Preset
        .Range("A5").Value = "Preset:"
        .Range("A5").Font.Bold = True
        .Range("B5").Interior.Color = RGB(255, 255, 200)  ' Light yellow

        ' Row 6: Ticker
        .Range("A6").Value = "Ticker:"
        .Range("A6").Font.Bold = True
        .Range("B6").Interior.Color = RGB(255, 255, 200)

        ' Row 7: Sector
        .Range("A7").Value = "Sector:"
        .Range("A7").Font.Bold = True
        .Range("B7").Interior.Color = RGB(255, 255, 200)

        ' Row 8: Bucket
        .Range("A8").Value = "Bucket:"
        .Range("A8").Font.Bold = True
        .Range("B8").Interior.Color = RGB(255, 255, 200)

        ' Row 9: Entry Price
        .Range("A9").Value = "Entry Price:"
        .Range("A9").Font.Bold = True
        .Range("B9").Interior.Color = RGB(255, 255, 200)
        .Range("B9").NumberFormat = "0.00"

        ' Row 10: ATR N
        .Range("A10").Value = "ATR N:"
        .Range("A10").Font.Bold = True
        .Range("B10").Interior.Color = RGB(255, 255, 200)
        .Range("B10").NumberFormat = "0.000"

        ' Row 11: K (Stop Multiple)
        .Range("A11").Value = "K (Stop Multiple):"
        .Range("A11").Font.Bold = True
        .Range("B11").Interior.Color = RGB(255, 255, 200)
        .Range("B11").Value = 2
        .Range("B11").NumberFormat = "0"

        ' Row 13: Method label
        .Range("A13").Value = "Method:"
        .Range("A13").Font.Bold = True

        ' Method option labels (actual buttons created separately)
        .Range("B13").Value = "1=Stock, 2=Opt-DeltaATR, 3=Opt-MaxLoss"
        .Range("C13").Value = 1  ' Default to Stock
        .Range("C13").Interior.Color = RGB(220, 220, 220)

        ' Row 16: Delta (for options)
        .Range("A16").Value = "Delta:"
        .Range("B16").Interior.Color = RGB(255, 255, 200)
        .Range("B16").Value = 0.3
        .Range("B16").NumberFormat = "0.00"

        ' Row 17: DTE (for options)
        .Range("A17").Value = "DTE:"
        .Range("B17").Interior.Color = RGB(255, 255, 200)
        .Range("B17").NumberFormat = "0"

        ' Row 18: Max Loss (for options)
        .Range("A18").Value = "Max Loss/Contract:"
        .Range("B18").Interior.Color = RGB(255, 255, 200)
        .Range("B18").NumberFormat = "$#,##0.00"
    End With
End Sub

Private Sub CreateChecklistSection(ws As Worksheet)
    ' Create checklist labels (checkboxes added separately)
    With ws
        .Range("A20").Value = "ENTRY CHECKLIST"
        .Range("A20").Font.Bold = True
        .Range("A20").Font.Size = 12

        .Range("A21").Value = "[ ] FromPreset"
        .Range("A22").Value = "[ ] TrendPass"
        .Range("A23").Value = "[ ] LiquidityPass"
        .Range("A24").Value = "[ ] TVConfirm"
        .Range("A25").Value = "[ ] EarningsOK"
        .Range("A26").Value = "[ ] JournalOK"

        ' Hidden cells for checkbox links (C20:C25)
        .Range("C20:C26").Interior.Color = RGB(220, 220, 220)
        .Range("C20:C26").Value = False

        ' Add comment note (safely)
        On Error Resume Next
        .Range("C20").AddComment "Link checkboxes to these cells"
        On Error GoTo 0
    End With
End Sub

Private Sub CreateOutputSection(ws As Worksheet)
    ' Create output labels and cells
    With ws
        ' Row 5: R ($)
        .Range("E5").Value = "R ($):"
        .Range("E5").Font.Bold = True
        .Range("F5").Interior.Color = RGB(220, 240, 255)  ' Light blue
        .Range("F5").NumberFormat = "$#,##0.00"

        ' Row 6: Stop Distance
        .Range("E6").Value = "Stop Distance:"
        .Range("E6").Font.Bold = True
        .Range("F6").Interior.Color = RGB(220, 240, 255)
        .Range("F6").NumberFormat = "0.00"

        ' Row 7: Initial Stop
        .Range("E7").Value = "Initial Stop:"
        .Range("E7").Font.Bold = True
        .Range("F7").Interior.Color = RGB(220, 240, 255)
        .Range("F7").NumberFormat = "0.00"

        ' Row 8: Shares
        .Range("E8").Value = "Shares:"
        .Range("E8").Font.Bold = True
        .Range("F8").Interior.Color = RGB(220, 240, 255)
        .Range("F8").NumberFormat = "0"

        ' Row 9: Contracts
        .Range("E9").Value = "Contracts:"
        .Range("E9").Font.Bold = True
        .Range("F9").Interior.Color = RGB(220, 240, 255)
        .Range("F9").NumberFormat = "0"

        ' Row 10-12: Add levels
        .Range("E10").Value = "Add 1:"
        .Range("F10").NumberFormat = "0.00"
        .Range("F10").Formula = "=IF(B9>0,B9+(AddStepN*B10),"""")"

        .Range("E11").Value = "Add 2:"
        .Range("F11").NumberFormat = "0.00"
        .Range("F11").Formula = "=IF(B9>0,B9+(2*AddStepN*B10),"""")"

        .Range("E12").Value = "Add 3:"
        .Range("F12").NumberFormat = "0.00"
        .Range("F12").Formula = "=IF(B9>0,B9+(3*AddStepN*B10),"""")"

        ' Row 14-15: Heat preview
        .Range("E14").Value = "Portfolio Heat:"
        .Range("E14").Font.Bold = True
        .Range("F14").Interior.Color = RGB(200, 255, 200)  ' Light green
        .Range("F14").Value = "Click Recalc"

        .Range("E15").Value = "Bucket Heat:"
        .Range("E15").Font.Bold = True
        .Range("F15").Interior.Color = RGB(200, 255, 200)
        .Range("F15").Value = "Click Recalc"
    End With
End Sub

Private Sub CreateButtons(ws As Worksheet)
    ' Create command buttons
    Dim btn As Button

    ' Button 1: Evaluate
    Set btn = ws.Buttons.Add(30, 430, 100, 25)
    btn.Text = "Evaluate"
    btn.OnAction = "TF_UI.EvaluateChecklist"

    ' Button 2: Recalc Sizing
    Set btn = ws.Buttons.Add(140, 430, 100, 25)
    btn.Text = "Recalc Sizing"
    btn.OnAction = "TF_UI.RecalcSizing"

    ' Button 3: Save Decision
    Set btn = ws.Buttons.Add(250, 430, 100, 25)
    btn.Text = "Save Decision"
    btn.OnAction = "TF_UI.SaveDecision"

    ' Button 4: Import Candidates
    Set btn = ws.Buttons.Add(30, 460, 100, 25)
    btn.Text = "Import Candidates"
    btn.OnAction = "TF_Presets.ImportCandidatesPrompt"

    ' Button 5: Open FINVIZ
    Set btn = ws.Buttons.Add(140, 460, 100, 25)
    btn.Text = "Open FINVIZ"
    btn.OnAction = "TF_Presets.OpenPreset"

    ' Button 6: Start Timer (optional)
    Set btn = ws.Buttons.Add(250, 460, 100, 25)
    btn.Text = "Start Timer"
    btn.OnAction = "TF_UI.StartImpulseTimer"
End Sub

Private Sub CreateCheckboxes(ws As Worksheet)
    ' Create 6 Form Control checkboxes linked to C20:C25
    Call TF_Logger.WriteLogSection("CreateCheckboxes() - Start")

    Dim chk As CheckBox
    Dim i As Integer
    Dim topPos As Double
    Dim leftPos As Double
    Dim cellRow As Integer
    Dim checkboxCount As Integer

    On Error Resume Next

    ' Starting position - place checkboxes in column B (before the text labels in column A)
    ' Column B starts around pixel 160 (after column A which is 20 characters wide)
    leftPos = ws.Columns("B").Left + 5  ' 5 pixels into column B

    Call TF_Logger.WriteLog("Column B Left position: " & ws.Columns("B").Left)
    Call TF_Logger.WriteLog("Checkbox leftPos: " & leftPos)

    checkboxCount = 0

    ' Create 6 checkboxes
    For i = 0 To 5
        cellRow = 21 + i  ' Rows 21-26 for the checklist items
        topPos = ws.Rows(cellRow).Top + 2  ' Align with row

        Call TF_Logger.WriteLog("Creating checkbox " & (i + 1) & " at row " & cellRow & ", topPos=" & topPos & ", leftPos=" & leftPos)

        Err.Clear

        ' Add checkbox (Form Control)
        Set chk = ws.CheckBoxes.Add(leftPos, topPos, 15, 15)

        If Err.Number <> 0 Then
            Call TF_Logger.WriteLogError("CreateCheckboxes", Err.Number, "Failed to add checkbox " & (i + 1) & ": " & Err.Description)
        Else
            checkboxCount = checkboxCount + 1
            Call TF_Logger.WriteLog("Checkbox " & (i + 1) & " created successfully")

            ' Link to corresponding cell in column C
            chk.LinkedCell = "$C$" & (20 + i)
            Call TF_Logger.WriteLog("Linked to cell: $C$" & (20 + i))

            ' Remove default text (labels are already in column A)
            chk.Text = ""

            ' Set initial value to FALSE
            ws.Range("C" & (20 + i)).Value = False
        End If
    Next i

    Call TF_Logger.WriteLog("Total checkboxes created: " & checkboxCount & " of 6")

    On Error GoTo 0
End Sub

Private Sub ApplyFormatting(ws As Worksheet)
    ' Apply borders and cell styles
    Dim rng As Range

    ' Input section border
    Set rng = ws.Range("A5:B18")
    With rng.Borders
        .LineStyle = xlContinuous
        .Weight = xlThin
    End With

    ' Output section border
    Set rng = ws.Range("E5:F15")
    With rng.Borders
        .LineStyle = xlContinuous
        .Weight = xlThin
    End With

    ' Checklist section border
    Set rng = ws.Range("A20:C26")
    With rng.Borders
        .LineStyle = xlContinuous
        .Weight = xlThin
    End With
End Sub

Private Sub SetColumnWidths(ws As Worksheet)
    ' Set optimal column widths
    ws.Columns("A:A").ColumnWidth = 20
    ws.Columns("B:B").ColumnWidth = 15
    ws.Columns("C:C").ColumnWidth = 8
    ws.Columns("D:D").ColumnWidth = 2
    ws.Columns("E:E").ColumnWidth = 18
    ws.Columns("F:F").ColumnWidth = 15
End Sub

Sub AddCheckboxes()
    ' Manually add Form Control checkboxes
    ' (This is harder to automate, so providing manual instructions)
    Dim msg As String

    msg = "To add checkboxes:" & vbCrLf & vbCrLf
    msg = msg & "1. Go to Developer tab → Insert → Check Box (Form Control)" & vbCrLf
    msg = msg & "2. Draw 6 checkboxes next to cells A21:A26" & vbCrLf
    msg = msg & "3. Right-click each checkbox → Format Control" & vbCrLf
    msg = msg & "4. Set Cell link:" & vbCrLf
    msg = msg & "   - Checkbox 1 → $C$20" & vbCrLf
    msg = msg & "   - Checkbox 2 → $C$21" & vbCrLf
    msg = msg & "   - Checkbox 3 → $C$22" & vbCrLf
    msg = msg & "   - Checkbox 4 → $C$23" & vbCrLf
    msg = msg & "   - Checkbox 5 → $C$24" & vbCrLf
    msg = msg & "   - Checkbox 6 → $C$25" & vbCrLf & vbCrLf
    msg = msg & "5. Delete the checkbox labels (text already in column A)"

    MsgBox msg, vbInformation, "Add Checkboxes Manually"
End Sub
