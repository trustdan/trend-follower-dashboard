# Options Trend Dashboard — Build & Implementation Summary
_Last updated: 2025-10-26 13:58:56_

## Goals
- Create a **bird’s‑eye options calendar** by sector, showing positions across a **rolling 10‑week window** (2 weeks back + 8 weeks forward).
- Add an **interactive pre‑trade checklist** to curb impulsive entries while staying faithful to **Seykota / Turtle** trend methods.
- Introduce **contracts sizing helpers**, **preset toggles** (System‑1: 20/10 and System‑2: 55/10), and a **Decisions** log.
- Make the workbook resilient to Excel quirks (dynamic arrays repairs, table naming), with macros that can rebuild views.

## Deliverables & Where Things Live
- **Trades** (sheet): one row per position (Symbol, Sector, StartDate, EndDate, Active, Strategy, Notes).
- **Calendar** (sheet): auto‑populates sector/week cells for all active trades in the 2‑back/8‑forward window. Weeks start **Monday**.
- **Summary** (sheet): counts/coverage per week and simple bar charts.
- **Checklist** (sheet): interactive; buttons in column **L**; three‑state banner; contracts helper; preset toggles; save decisions.
- **Decisions** (sheet): created on first save via the **Save Decision** button; logs all ticket inputs and helper outputs.
- **Lists** (sheet, hidden): sector & week values for validation; safe to edit if you add/remove sectors.
- **Modules**
  - `CalendarModule`: `ForceRepairAndRefresh`, `RefreshCalendar`, `StartOfWeek` (rebuilds the calendar & summary as values).
  - `ChecklistModule`: full UX logic for the checklist; buttons; banner coloring; presets; decisions; contract helper.
  - `Worksheet (Checklist)`: event handlers to update banner/colors on any checkbox or key input change.

## Calendar Details (date‑based; no week numbers)
- **Header** `B1:K1` = Mondays; `B2:K2` = “mmm d – mmm d” ranges.
- Cell shows a symbol if **date range overlaps** that week: `StartDate <= WeekEnd AND EndDate >= WeekStart AND Active="Yes"`.
- **Totals row** counts active trades per week; **Summary** sheet charts duplicate those counts and sector coverage.
- If you ever see _“Removed Records: Formula from /xl/worksheets/sheet4.xml”_, use the macros in `CalendarModule` (they avoid structured refs).

## Checklist Logic (Seykota / Turtle)
- **Required groups** (minimal roadblocks):
  - **SIG_REQ**: 55‑bar breakout (long > 55‑high / short < 55‑low).
  - **RISK_REQ**: risk per unit = % of equity using **2×N** (ATR) stop; pyramids **every 0.5×N** up to **Max Units**.
  - **OPT_REQ**: 60–90 DTE; **roll ~21 DTE**; **liquidity required** (bid/ask <10% of mid; OI >100).
  - **EXIT_REQ**: exit by **10‑bar opposite Donchian OR closer of 2×N**.
  - **BEHAV_REQ**: 2‑minute cool‑off; no intraday overrides.
- **Optional “QUAL” items** give score; **default threshold = 3.0** (you can change cell `I6`).

### Three‑state banner (cell J5)
- **RED: DO NOT TRADE** — any required group fails (or score < 6 in the “soft band” variant).
- **YELLOW: CAUTION** — all required pass but score < threshold **or** theta rule fails.
- **GREEN: OK TO TRADE** — all required pass, theta OK, and score ≥ threshold.
- The color is set by code (not CF) so it updates reliably on every change.

### Buttons (column L)
- **Reset** — unchecks every box.
- **Add To Trades** — appends to the Trades table with StartDate = Today and EndDate = Today + DTE.
- **Recalculate** — forces a full refresh if needed.
- **Apply Preset** — journals the selected preset (System‑1 20/10 or System‑2 55/10) without altering formulas.
- **Save Decision** — logs date/time, symbol, status, score, N/stopN, risk%, max units, DTE, contracts helper outputs, etc.

## Contracts Helper
Inputs (left panel): **Account size**, **N (ATR)**, **Stop multiple (N)**, **Add step (N)**, **Risk % per unit**, **Max units**, **Underlier price**, **Delta**, **Debit per spread**, **Structure**.
- Per‑unit risk $ = `AccountSize × Risk%`.
- **Single-call contracts** ≈ `unitRisk$ / (delta × (2×N) × 100)`.
- **Debit vertical spreads** ≈ `unitRisk$ / debit`.
- **Add levels**: `Entry + (k × AddStepN × N)` for k=1..(MaxUnits-1).

## Presets
- **System‑1 (20/10)** and **System‑2 (55/10)** are recorded for journaling (“Preset” and “Entry/Exit lens”). The checklist still uses the simplified required groups above; you’re not locked to one preset’s math.

## Typical Run Order
1. **Enable macros** (yellow banner → “Enable Content”).  
2. **Build** checklist: run `BuildInteractiveChecklist` (once), then `ApplyChecklistFixes` (once).  
3. **Calendar**: if formulas were stripped, run `ForceRepairAndRefresh` to rebuild.  
4. Use **Checklist** to qualify a trade → **Add To Trades** → **Save Decision**.

## Troubleshooting
- Banner color doesn’t change → run **Recalculate** (re‑enables events) or `ApplyChecklistFixes` (wipes CF and re‑styles J5).
- Buttons overlap → run **FixChecklistUI** (or simply `ApplyChecklistFixes`; buttons live in column **L**).
- Calendar empty → check `Trades` dates + `Active="Yes"`; run `ForceRepairAndRefresh`.
- “Formula removed” on open → your Excel build dislikes cross‑sheet dynamic arrays; use the `CalendarModule` rebuild macros.

## Customization knobs
- **Sectors**: edit hidden **Lists** sheet.
- **Week start**: change Monday logic in `CalendarModule` if you prefer Sundays.
- **Quality threshold**: set `Checklist!I6` (e.g., 4–6 for stricter CAUTION).
- **Risk per unit**: set `Checklist!C11` (e.g., 0.5–1.0%).

## Requirements
- Excel 365 / 2021 for dynamic arrays (calendar); macros required for checklist UI.
- No external references; VBA is late‑bound (no library setup needed).





---



# CalendarModule

Option Explicit

' === CONFIG ===
Const WEEKS_BACK As Long = 2        ' 2 weeks back
Const WEEKS_TOTAL As Long = 10      ' total (2 back + 8 forward)
Const STARTS_ON_MONDAY As Boolean = True

Sub ForceRepairAndRefresh()
    ' Clears old formulas that cause #NAME? and rebuilds all values
    Dim wsCal As Worksheet, wsSum As Worksheet
    Set wsCal = Worksheets("Calendar")
    Set wsSum = Worksheets("Summary")

    Application.ScreenUpdating = False
    Application.EnableEvents = False
    
    ' Clear header labels and grid (safe to run repeatedly)
    wsCal.Range("B1:K2").ClearContents
    wsCal.Range("B3:K100").ClearContents
    wsCal.Range("B100:K200").ClearContents  ' harmless extra wipe for totals row area
    wsSum.Range("A2:C200").ClearContents
    
    RefreshCalendar   ' rebuilds Calendar + Summary with values
    
    Application.EnableEvents = True
    Application.ScreenUpdating = True
End Sub

Sub RefreshCalendar()
    Dim wsCal As Worksheet, wsTr As Worksheet, wsSum As Worksheet
    Set wsCal = Worksheets("Calendar")
    Set wsTr = Worksheets("Trades")
    Set wsSum = Worksheets("Summary")

    Application.ScreenUpdating = False
    Application.EnableEvents = False
    
    ' --- 1) Rolling week starts (row 1) + labels (row 2) ---
    Dim wkStart(1 To WEEKS_TOTAL) As Date
    Dim base As Date: base = StartOfWeek(Date, STARTS_ON_MONDAY) - 7 * WEEKS_BACK
    
    Dim i As Long
    For i = 1 To WEEKS_TOTAL
        wkStart(i) = base + 7 * (i - 1)
        With wsCal.Cells(1, i + 1)
            .Value = wkStart(i)
            .NumberFormat = "mm/dd/yyyy"
        End With
        wsCal.Cells(2, i + 1).Value = Format(wkStart(i), "mmm d") & " – " & Format(wkStart(i) + 6, "mmm d")
        wsCal.Columns(i + 1).ColumnWidth = 22
    Next i
    wsCal.Columns(1).ColumnWidth = 24 ' Sector column
    
    ' --- 2) Determine sector rows from column A ---
    Dim firstSecRow As Long, lastSecRow As Long
    firstSecRow = 3
    lastSecRow = wsCal.Cells(wsCal.Rows.count, 1).End(xlUp).Row
    If lastSecRow < firstSecRow Then GoTo WrapUp
    
    ' --- 3) Get last trade row ---
    Dim lastTrRow As Long
    lastTrRow = wsTr.Cells(wsTr.Rows.count, 1).End(xlUp).Row
    
    ' --- 4) Build grid (values only) ---
    Dim r As Long, tr As Long, cnt As Long
    Dim sec As String, out As String
    Dim sd As Variant, ed As Variant
    
    For r = firstSecRow To lastSecRow
        sec = UCase$(Trim$(wsCal.Cells(r, 1).Value))
        wsCal.Rows(r).RowHeight = 60
        For i = 1 To WEEKS_TOTAL
            out = "": cnt = 0
            For tr = 2 To lastTrRow
                If UCase$(Trim$(wsTr.Cells(tr, 5).Value)) = "YES" Then       ' Active
                    If UCase$(Trim$(wsTr.Cells(tr, 2).Value)) = sec Then     ' Sector match
                        sd = wsTr.Cells(tr, 3).Value                          ' StartDate
                        ed = wsTr.Cells(tr, 4).Value                          ' EndDate
                        If IsDate(sd) Then
                            If Not IsDate(ed) Or ed = "" Then ed = sd
                            If CDate(sd) <= wkStart(i) + 6 And CDate(ed) >= wkStart(i) Then
                                If Len(out) > 0 Then out = out & vbLf
                                out = out & UCase$(Trim$(wsTr.Cells(tr, 1).Value)) ' Symbol
                                cnt = cnt + 1
                            End If
                        End If
                    End If
                End If
            Next tr
    
            With wsCal.Cells(r, 1 + i)
                .Value = out
                .WrapText = True
                .VerticalAlignment = xlTop
                ' light green (1–2), light yellow (3+), none if 0
                If cnt = 0 Then
                    .Interior.Pattern = xlNone
                ElseIf cnt >= 3 Then
                    .Interior.Color = RGB(255, 235, 156)
                Else
                    .Interior.Color = RGB(198, 239, 206)
                End If
            End With
        Next i
    Next r
    
    ' --- 5) Totals row ---
    Dim totalsRow As Long: totalsRow = lastSecRow + 1
    wsCal.Cells(totalsRow, 1).Value = "Trades / Week"
    For i = 1 To WEEKS_TOTAL
        cnt = 0
        For tr = 2 To lastTrRow
            If UCase$(Trim$(wsTr.Cells(tr, 5).Value)) = "YES" Then
                sd = wsTr.Cells(tr, 3).Value
                ed = wsTr.Cells(tr, 4).Value
                If IsDate(sd) Then
                    If Not IsDate(ed) Or ed = "" Then ed = sd
                    If CDate(sd) <= wkStart(i) + 6 And CDate(ed) >= wkStart(i) Then cnt = cnt + 1
                End If
            End If
        Next tr
        With wsCal.Cells(totalsRow, 1 + i)
            .Value = cnt
            .HorizontalAlignment = xlCenter
        End With
    Next i
    
    ' --- 6) Summary (values only; charts will update) ---
    For i = 1 To WEEKS_TOTAL
        wsSum.Cells(1 + i, 1).Value = wkStart(i)                    ' Week start
        wsSum.Cells(1 + i, 1).NumberFormat = "mm/dd/yyyy"
        wsSum.Cells(1 + i, 2).Value = wsCal.Cells(totalsRow, 1 + i) ' Active trades
        ' Sectors covered (non-empty sector cells in that week)
        cnt = 0
        For r = firstSecRow To lastSecRow
            If Len(wsCal.Cells(r, 1 + i).Value) > 0 Then cnt = cnt + 1
        Next r
        wsSum.Cells(1 + i, 3).Value = cnt
    Next i

WrapUp:
    Application.EnableEvents = True
    Application.ScreenUpdating = True
End Sub

Private Function StartOfWeek(ByVal d As Date, ByVal monday As Boolean) As Date
    If monday Then
        StartOfWeek = d - Weekday(d, vbMonday) + 1  ' Monday
    Else
        StartOfWeek = d - Weekday(d, vbSunday) + 1  ' Sunday
    End If
End Function



---



# ChecklistModule

Option Explicit

' ======================================================================
' ChecklistModule (ASCII-only)
' - Builds "Checklist" sheet with checkboxes
' - 3-state banner: DO NOT TRADE / CAUTION / OK TO TRADE
' - Contracts helper (single ITM call & debit vertical)
' - Preset toggles (System-1 20/10, System-2 55/10) + Apply Preset button
' - Buttons in UI column L (no overlap)
' - Decisions logger to "Decisions"
' ======================================================================

Private Const UI_COL As Long = 12    ' Column L anchors all buttons

' ========================= PUBLIC ENTRY POINTS =========================

Public Sub BuildInteractiveChecklist()
    Dim ws As Worksheet
    On Error Resume Next
    Application.DisplayAlerts = False
    Set ws = ThisWorkbook.Worksheets("Checklist")
    If Not ws Is Nothing Then ws.Delete
    Application.DisplayAlerts = True
    On Error GoTo 0

    Set ws = ThisWorkbook.Worksheets.Add(After:=ThisWorkbook.Sheets(ThisWorkbook.Sheets.count))
    ws.Name = "Checklist"
    ws.Cells.Font.Name = "Calibri"
    ws.Cells.Font.Size = 11
    
    SetupTicket ws
    
    Dim items As Variant
    items = ItemsData()
    
    Dim firstRow As Long, lastRow As Long
    RenderChecklist ws, items, firstRow, lastRow
    RenderCalcs ws, firstRow, lastRow
    AddButtons ws
    
    ' tidy columns/rows
    ws.Columns("A").ColumnWidth = 20
    ws.Columns("B").ColumnWidth = 60
    ws.Columns("C").ColumnWidth = 9
    ws.Columns("D").ColumnWidth = 7
    ws.Columns("E").ColumnWidth = 7
    ws.Columns("F").ColumnWidth = 2          ' spacer
    ws.Columns("L").ColumnWidth = 20         ' UI column
    ws.Rows("1:8").RowHeight = 20
    ws.Range("A1").RowHeight = 28
    
    ws.Columns("F:G").Hidden = True          ' helper columns
    ws.Columns("A:D").EntireColumn.AutoFit
    
    PlaceButtonsNeatly ws
    EnsureEventsOn
    RefreshBanner
    ws.Activate
End Sub

Public Sub FixChecklistUI()
    Dim ws As Worksheet: Set ws = Worksheets("Checklist")
    ws.Rows("1:8").RowHeight = 20
    ws.Range("A1").RowHeight = 28
    ws.Columns("L").ColumnWidth = 20

    ' lenient theta (OK if inputs blank/zero)
    ws.Range("H24").Value = "Theta budget OK"
    ws.Range("I24").Formula = "=IF(OR(C4="""",I3="""",C4=0,I3=0),TRUE, I2/C4<=I3)"
    
    ' banner text (3-state)
    ws.Range("J5").Formula = "=IF(NOT(I22),""DO NOT TRADE"", IF(AND(I24, I5>=I6), ""OK TO TRADE"", ""CAUTION""))"
    
    ' remove any CF from banner; color is code-driven
    On Error Resume Next
    Dim i As Long
    For i = ws.Range("J5").FormatConditions.count To 1 Step -1
        ws.Range("J5").FormatConditions.Item(i).Delete
    Next i
    On Error GoTo 0
    
    PlaceButtonsNeatly ws
    EnsureEventsOn
    RefreshBanner
End Sub

' Reapply banner formula, nuke CF, tidy buttons, refresh
Public Sub ApplyChecklistFixes()
    Dim ws As Worksheet: Set ws = Worksheets("Checklist")

    ws.Range("J5").Formula = "=IF(NOT(I22),""DO NOT TRADE""," & _
                              "IF(AND(I24, I5>=I6),""OK TO TRADE""," & _
                              "IF(OR(NOT(I24), I5>=6),""CAUTION"",""DO NOT TRADE"")))"
    
    ws.Range("I24").Formula = "=IF(OR(C4="""",I3="""",C4=0,I3=0),TRUE, I2/C4<=I3)"
    
    With ws.Range("J5")
        On Error Resume Next
        Dim k As Long
        For k = .FormatConditions.count To 1 Step -1
            .FormatConditions.Item(k).Delete
        Next k
        On Error GoTo 0
        .Interior.Pattern = xlSolid
        .Interior.ColorIndex = xlColorIndexNone
        .Font.ColorIndex = xlColorIndexAutomatic
        .Font.Bold = True
        .HorizontalAlignment = xlCenter
        .VerticalAlignment = xlCenter
        .RowHeight = 24
        .ColumnWidth = 18
    End With
    
    PlaceButtonsNeatly ws
    EnsureEventsOn
    ws.Calculate
    RefreshBanner
End Sub

Public Sub ResetChecklist()
    Dim ws As Worksheet: Set ws = Worksheets("Checklist")
    Dim shp As Shape
    For Each shp In ws.Shapes
        If shp.Type = msoFormControl Then
            On Error Resume Next
            ws.Range(shp.ControlFormat.LinkedCell).Value = False
            On Error GoTo 0
        End If
    Next shp
    RecalcChecklist
End Sub

Public Sub AddToTrades()
    Dim wsC As Worksheet: Set wsC = Worksheets("Checklist")
    Dim wsT As Worksheet
    On Error Resume Next
    Set wsT = Worksheets("Trades")
    On Error GoTo 0
    If wsT Is Nothing Then
        MsgBox "Trades sheet not found.", vbExclamation
        Exit Sub
    End If

    If UCase$(wsC.Range("J5").Value) <> "OK TO TRADE" Then
        If MsgBox("Status is not OK TO TRADE. Add anyway?", vbQuestion + vbYesNo) = vbNo Then Exit Sub
    End If
    
    Dim sym As String, sec As String, strat As String
    Dim dte As Long, sd As Date, ed As Date
    sym = UCase$(Trim$(wsC.Range("C2").Value))
    sec = Trim$(wsC.Range("C3").Value))
    strat = Trim$(wsC.Range("C6").Value)
    If sym = "" Or sec = "" Then
        MsgBox "Enter Symbol and Sector first.", vbExclamation
        Exit Sub
    End If
    
    sd = Date
    If Len(wsC.Range("C5").Value) = 0 Then dte = 60 Else dte = CLng(wsC.Range("C5").Value)
    ed = sd + dte
    
    Dim lo As ListObject
    On Error Resume Next
    Set lo = wsT.ListObjects(1)
    On Error GoTo 0
    
    If Not lo Is Nothing Then
        Dim newRow As ListRow
        Set newRow = lo.ListRows.Add
        With newRow.Range
            .Cells(1, 1).Value = sym
            .Cells(1, 2).Value = sec
            .Cells(1, 3).Value = sd
            .Cells(1, 4).Value = ed
            .Cells(1, 5).Value = "Yes"
            .Cells(1, 6).Value = strat
            .Cells(1, 7).Value = wsC.Range("C7").Value
        End With
    Else
        Dim r As Long
        r = wsT.Cells(wsT.Rows.count, 1).End(xlUp).Row + 1
        wsT.Cells(r, 1).Resize(1, 7).Value = Array(sym, sec, sd, ed, "Yes", strat, wsC.Range("C7").Value)
    End If
    
    On Error Resume Next
    Application.Run "'" & ThisWorkbook.Name & "'!RefreshCalendar"
    On Error GoTo 0
    
    MsgBox "Added to Trades: " & sym & "  DTE=" & dte, vbInformation
End Sub

Public Sub RecalcChecklist()
    EnsureEventsOn
    Worksheets("Checklist").Calculate
    RefreshBanner
End Sub

' Contracts helper uses per-unit risk dollars (I27) to compute counts
Public Sub ApplyPreset()
    Dim ws As Worksheet: Set ws = Worksheets("Checklist")
    Dim p As String: p = Trim$(CStr(ws.Range("C18").Value))
    If p = "" Then
        MsgBox "Pick a preset in C18 (System-1 or System-2).", vbInformation
        Exit Sub
    End If
    ' These are journaling helpers — risk model stays the same; we note entry/exit lens.
    If p = "System-1 (20/10)" Then
        ws.Range("I36").Value = "System-1"
        ws.Range("I37").Value = "20/10"
    Else
        ws.Range("I36").Value = "System-2"
        ws.Range("I37").Value = "55/10"
    End If
    RecalcChecklist
End Sub

Public Sub SaveDecision()
    Dim wsC As Worksheet: Set wsC = Worksheets("Checklist")
    Dim wsD As Worksheet
    On Error Resume Next
    Set wsD = Worksheets("Decisions")
    On Error GoTo 0
    If wsD Is Nothing Then
        Set wsD = Worksheets.Add(After:=Worksheets(Worksheets.count))
        wsD.Name = "Decisions"
        wsD.Range("A1:Q1").Value = Array( _
            "Timestamp", "Symbol", "Sector", "Strategy", _
            "Score", "Threshold", "ThetaOK", "Status", "Notes", _
            "N_at_entry", "StopN", "AddStepN", "RiskPct_per_unit", "MaxUnits", "Planned_DTE", _
            "Struct", "Delta", "Debit", "PxEntry", "CallContracts", "VerticalSpreads", "Preset")
        wsD.Rows(1).Font.Bold = True
        wsD.Columns("A:Q").EntireColumn.AutoFit
    End If

    Dim nextRow As Long
    nextRow = wsD.Cells(wsD.Rows.count, 1).End(xlUp).Row + 1
    
    ' Inputs & helper outputs
    wsD.Cells(nextRow, 1).Value = Now
    wsD.Cells(nextRow, 2).Value = UCase$(Trim$(wsC.Range("C2").Value))
    wsD.Cells(nextRow, 3).Value = Trim$(wsC.Range("C3").Value)
    wsD.Cells(nextRow, 4).Value = Trim$(wsC.Range("C6").Value)
    wsD.Cells(nextRow, 5).Value = wsC.Range("I5").Value                  ' score
    wsD.Cells(nextRow, 6).Value = wsC.Range("I6").Value                  ' threshold
    wsD.Cells(nextRow, 7).Value = wsC.Range("I24").Value                 ' theta OK
    wsD.Cells(nextRow, 8).Value = wsC.Range("J5").Value                  ' status
    wsD.Cells(nextRow, 9).Value = wsC.Range("C7").Value                  ' notes
    wsD.Cells(nextRow, 10).Value = wsC.Range("C8").Value                 ' N
    wsD.Cells(nextRow, 11).Value = wsC.Range("C9").Value                 ' stopN
    wsD.Cells(nextRow, 12).Value = wsC.Range("C10").Value                ' add step N
    wsD.Cells(nextRow, 13).Value = wsC.Range("C11").Value                ' risk pct
    wsD.Cells(nextRow, 14).Value = wsC.Range("C12").Value                ' max units
    wsD.Cells(nextRow, 15).Value = wsC.Range("C5").Value                 ' DTE
    wsD.Cells(nextRow, 16).Value = wsC.Range("C16").Value                ' structure
    wsD.Cells(nextRow, 17).Value = wsC.Range("C14").Value                ' delta
    wsD.Cells(nextRow, 18).Value = wsC.Range("C15").Value                ' debit
    wsD.Cells(nextRow, 19).Value = wsC.Range("C13").Value                ' price entry
    wsD.Cells(nextRow, 20).Value = wsC.Range("I30").Value                ' call contracts
    wsD.Cells(nextRow, 21).Value = wsC.Range("I31").Value                ' vertical spreads
    wsD.Cells(nextRow, 22).Value = wsC.Range("C18").Value                ' preset
    
    wsD.Columns("A:Q").EntireColumn.AutoFit
    MsgBox "Decision saved.", vbInformation
End Sub

' Called by each checkbox via OnAction so colors update immediately
Public Sub Checklist_CheckBoxChanged()
    RecalcChecklist
End Sub

' ============================== INTERNALS ===============================

Private Sub SetupTicket(ws As Worksheet)
    ws.Range("A1").Value = "Seykota / Turtle Pre-Trade Checklist"
    With ws.Range("A1")
        .Font.Bold = True
        .Font.Size = 16
    End With

    ' Core ticket
    ws.Range("A2").Value = "Symbol:"
    ws.Range("A3").Value = "Sector:"
    ws.Range("A4").Value = "Account Size ($):"
    ws.Range("A5").Value = "Planned DTE:"
    ws.Range("A6").Value = "Strategy:"
    ws.Range("A7").Value = "Notes:"
    
    ' Sizing / contracts inputs
    ws.Range("A8").Value = "N (ATR) at entry:"
    ws.Range("A9").Value = "Stop multiple (N):"
    ws.Range("A10").Value = "Add step (N):"
    ws.Range("A11").Value = "Risk % per unit:"
    ws.Range("A12").Value = "Max units:"
    ws.Range("A13").Value = "Underlier @ entry ($):"
    ws.Range("A14").Value = "Option delta (single call):"
    ws.Range("A15").Value = "Debit per vertical ($):"
    ws.Range("A16").Value = "Structure:"
    ws.Range("A18").Value = "Preset:"
    
    ws.Range("C2:C16").Interior.Color = RGB(242, 242, 242)
    ws.Range("C2:C16").Borders.LineStyle = xlContinuous
    
    ws.Range("C4").NumberFormat = "$#,##0"
    ws.Range("C5").NumberFormat = "0"
    ws.Range("C8").NumberFormat = "0.00"
    ws.Range("C9").NumberFormat = "0.00"
    ws.Range("C10").NumberFormat = "0.00"
    ws.Range("C11").NumberFormat = "0.00%"
    ws.Range("C12").NumberFormat = "0"
    ws.Range("C13").NumberFormat = "$#,##0.00"
    ws.Range("C14").NumberFormat = "0.00"
    ws.Range("C15").NumberFormat = "$#,##0.00"
    
    ' Defaults per system playbook
    ws.Range("C9").Value = 2        ' stopN
    ws.Range("C10").Value = 0.5     ' add step N
    ws.Range("C11").Value = 0.0075  ' risk pct / unit
    ws.Range("C12").Value = 4
    ws.Range("C14").Value = 0.65
    ws.Range("C15").Value = 5
    
    ' Structure & preset dropdowns
    On Error Resume Next
    ws.Range("C16").Validation.Delete
    ws.Range("C16").Validation.Add Type:=xlValidateList, AlertStyle:=xlValidAlertStop, _
        Operator:=xlBetween, Formula1:="Single Call,Debit Vertical"
    
    ws.Range("C18").Validation.Delete
    ws.Range("C18").Validation.Add Type:=xlValidateList, AlertStyle:=xlValidAlertStop, _
        Operator:=xlBetween, Formula1:="System-1 (20/10),System-2 (55/10)"
    ws.Range("C18").Value = "System-2 (55/10)"
    On Error GoTo 0
    
    ' Theta/score/status
    ws.Range("G2").Value = "Portfolio Daily Theta ($):"
    ws.Range("G3").Value = "Max Theta Budget (% acct):"
    ws.Range("I2").NumberFormat = "$#,##0"
    ws.Range("I2").Value = 0
    ws.Range("I3").NumberFormat = "0.00%"
    ws.Range("I3").Value = 0.0025
    
    ws.Range("H5").Value = "Quality Score"
    ws.Range("H6").Value = "Threshold"
    ws.Range("H7").Value = "Required Rules"
    ws.Range("I5").NumberFormat = "0.0"
    ws.Range("I6").Value = 3                    ' fewer road-blocks
    ws.Range("I6").NumberFormat = "0.0"
    ws.Range("I7").Value = "-"
    
    ws.Range("J5").Value = "DO NOT TRADE"
    With ws.Range("J5")
        .Font.Bold = True
        .HorizontalAlignment = xlCenter
        .VerticalAlignment = xlCenter
        .RowHeight = 24
        .ColumnWidth = 18
    End With
End Sub

Private Function RowItem(sec As String, txt As String, req As Boolean, wt As Double, tag As String) As Variant
    RowItem = Array(sec, txt, req, wt, tag)
End Function

Private Function ItemsData() As Variant
    Dim col As Object: Set col = CreateObject("System.Collections.ArrayList")
    Dim v As Variant, i As Long

    ' SIGNAL (one required)
    col.Add RowItem("Signal", "55-bar Donchian breakout confirmed (long > 55-high / short < 55-low)", True, 0, "SIG_REQ")
    col.Add RowItem("Signal", "Optional: regime OK (SPY vs 200SMA) if used", False, 1, "QUAL")
    col.Add RowItem("Signal", "Optional: not >2N above 20-EMA at entry", False, 1, "QUAL")
    
    ' RISK / SIZE
    col.Add RowItem("Risk/Size", "Per-unit risk = % of equity from 2×N stop", True, 0, "RISK_REQ")
    col.Add RowItem("Risk/Size", "Pyramids planned: add every 0.5×N to max units", True, 0, "RISK_REQ")
    col.Add RowItem("Risk/Size", "Journal: units/max and why now", False, 1, "QUAL")
    
    ' OPTIONS (expression)
    col.Add RowItem("Options", "Use 60–90 DTE; roll/close ~21 DTE", True, 0, "OPT_REQ")
    col.Add RowItem("Options", "Liquidity OK (bid-ask <10% mid; OI >100)", True, 0, "OPT_REQ")
    col.Add RowItem("Options", "Optional: vertical profit take 50–75%", False, 1, "QUAL")
    
    ' EXITS
    col.Add RowItem("Exits", "Exit: 10-bar opposite Donchian OR closer of 2×N", True, 0, "EXIT_REQ")
    
    ' BEHAVIOR
    col.Add RowItem("Behavior", "No intraday overrides; use resting orders", True, 0, "BEHAV_REQ")
    col.Add RowItem("Behavior", "2-minute cool-off before entry", True, 0, "BEHAV_REQ")
    col.Add RowItem("Behavior", "Optional: earnings blackout when long premium", False, 1, "QUAL")
    
    ReDim v(0 To col.count - 1)
    For i = 0 To col.count - 1
        v(i) = col(i)
    Next i
    ItemsData = v
End Function

Private Sub RenderChecklist(ws As Worksheet, items As Variant, _
                            ByRef firstRow As Long, ByRef lastRow As Long)
    Dim r As Long, i As Long
    firstRow = 10
    r = firstRow

    ws.Range("A9:D9").Value = Array("Section", "Checklist Item", "Req?", "Weight")
    ws.Range("A9:D9").Font.Bold = True
    ws.Range("E9").Value = "Done"
    ws.Range("A9:E9").Interior.Color = RGB(221, 235, 247)
    
    For i = LBound(items) To UBound(items)
        ws.Cells(r, 1).Value = items(i)(0)
        ws.Cells(r, 2).Value = items(i)(1)
        ws.Cells(r, 3).Value = IIf(items(i)(2), "Req", "Opt")
        ws.Cells(r, 4).Value = items(i)(3)
        ws.Cells(r, 7).Value = items(i)(4)
    
        Dim cb As Object
        Set cb = ws.CheckBoxes.Add(ws.Cells(r, 5).Left + 2, ws.Cells(r, 5).Top + 2, 12, 12)
        cb.caption = ""
        cb.LinkedCell = ws.Cells(r, 6).Address(False, False)  ' column F TRUE/FALSE
        cb.Value = xlOff
        cb.Placement = xlMoveAndSize
        cb.OnAction = "Checklist_CheckBoxChanged"
    
        r = r + 1
    Next i
    lastRow = r - 1
    
    ws.Range(ws.Cells(firstRow, 1), ws.Cells(lastRow, 5)).Borders.LineStyle = xlContinuous
End Sub

Private Sub RenderCalcs(ws As Worksheet, firstRow As Long, lastRow As Long)
    ws.Range("H9").Value = "Calc Helpers"
    ws.Range("H9").Font.Bold = True

    ' Group gates
    ws.Range("H10").Value = "Signal OK"
    ws.Range("I10").Formula = "=COUNTIF($G$" & firstRow & ":$G$" & lastRow & ",""SIG_REQ"")=" & _
                               "SUMPRODUCT(--($F$" & firstRow & ":$F$" & lastRow & "=TRUE),--($G$" & firstRow & ":$G$" & lastRow & "=""SIG_REQ""))"
    
    ws.Range("H11").Value = "RISK_REQ OK"
    ws.Range("I11").Formula = "=COUNTIF($G$" & firstRow & ":$G$" & lastRow & ",""RISK_REQ"")=" & _
                               "SUMPRODUCT(--($F$" & firstRow & ":$F$" & lastRow & "=TRUE),--($G$" & firstRow & ":$G$" & lastRow & "=""RISK_REQ""))"
    
    ws.Range("H12").Value = "OPT_REQ OK"
    ws.Range("I12").Formula = "=COUNTIF($G$" & firstRow & ":$G$" & lastRow & ",""OPT_REQ"")=" & _
                               "SUMPRODUCT(--($F$" & firstRow & ":$F$" & lastRow & "=TRUE),--($G$" & firstRow & ":$G$" & lastRow & "=""OPT_REQ""))"
    
    ws.Range("H13").Value = "EXIT_REQ OK"
    ws.Range("I13").Formula = "=COUNTIF($G$" & firstRow & ":$G$" & lastRow & ",""EXIT_REQ"")=" & _
                               "SUMPRODUCT(--($F$" & firstRow & ":$F$" & lastRow & "=TRUE),--($G$" & firstRow & ":$G$" & lastRow & "=""EXIT_REQ""))"
    
    ws.Range("H14").Value = "BEHAV_REQ OK"
    ws.Range("I14").Formula = "=COUNTIF($G$" & firstRow & ":$G$" & lastRow & ",""BEHAV_REQ"")=" & _
                               "SUMPRODUCT(--($F$" & firstRow & ":$F$" & lastRow & "=TRUE),--($G$" & firstRow & ":$G$" & lastRow & "=""BEHAV_REQ""))"
    
    ' ALL required OK
    ws.Range("H22").Value = "ALL required OK"
    ws.Range("I22").Formula = "=AND(I10,I11,I12,I13,I14)"
    
    ' Quality score (optional rows have weight 1)
    ws.Range("H23").Value = "Quality score"
    ws.Range("I23").Formula = "=SUMPRODUCT($D$" & firstRow & ":$D$" & lastRow & ",--($F$" & firstRow & ":$F$" & lastRow & "=TRUE),--($C$" & firstRow & ":$C$" & lastRow & "=""Opt""))"
    ws.Range("I23").NumberFormat = "0.0"
    ws.Range("I5").Formula = "=I23"
    
    ' Theta budget OK
    ws.Range("H24").Value = "Theta budget OK"
    ws.Range("I24").Formula = "=IF(OR(C4="""",I3="""",C4=0,I3=0),TRUE, I2/C4<=I3)"
    
    ' Final OK (not used by color directly)
    ws.Range("H25").Value = "Final OK"
    ws.Range("I25").Formula = "=AND(I22,I24,I5>=I6)"
    
    ' Banner text
    ws.Range("I7").Formula = "=IF(I22,""All req passed"",""Missing required items"")"
    ws.Range("J5").Formula = "=IF(NOT(I22),""DO NOT TRADE"", IF(AND(I24, I5>=I6), ""OK TO TRADE"", ""CAUTION""))"
    
    ' ---- Contracts helper & add levels ----
    ws.Range("H27").Value = "Per-unit risk $"
    ws.Range("I27").Formula = "=IF(C4>0, C4*IF(C11>0,C11,0.0075), 0)"
    
    ws.Range("H28").Value = "2N $ move"
    ws.Range("I28").Formula = "=IF(C8>0, 2*C8, 0)"
    
    ws.Range("H29").Value = "Structure"
    ws.Range("I29").Formula = "=C16"
    
    ws.Range("H30").Value = "Single-call contracts"
    ws.Range("I30").Formula = "=IF(AND(I27>0,C14>0,C8>0),MAX(1,INT(I27/(100*C14*(2*C8)))), """")"
    
    ws.Range("H31").Value = "Debit-vertical spreads"
    ws.Range("I31").Formula = "=IF(AND(I27>0,C15>0),MAX(1,INT(I27/C15)), """")"
    
    ws.Range("H32").Value = "Add1 price"
    ws.Range("I32").Formula = "=IF(AND(C13>0,C8>0,C10>0,C12>=2),C13+C10*C8, """")"
    
    ws.Range("H33").Value = "Add2 price"
    ws.Range("I33").Formula = "=IF(C12>=3, C13+2*C10*C8, """")"
    
    ws.Range("H34").Value = "Add3 price"
    ws.Range("I34").Formula = "=IF(C12>=4, C13+3*C10*C8, """")"
    
    ws.Range("H36").Value = "Preset"
    ws.Range("I36").Formula = "=C18"
    
    ws.Range("H37").Value = "Preset Entry/Exit"
    ws.Range("I37").Formula = "=IF(C18=""System-1 (20/10)"",""20/10"",""55/10"")"
End Sub

Private Sub GroupOK(ws As Worksheet, tag As String, outRow As Long, firstRow As Long, lastRow As Long)
    ws.Range("H" & outRow).Value = tag & " OK"
    ws.Range("I" & outRow).Formula = _
        "=COUNTIF($G$" & firstRow & ":$G$" & lastRow & ",""" & tag & """)=" & _
        "SUMPRODUCT(--($F$" & firstRow & ":$F$" & lastRow & "=TRUE),--($G$" & firstRow & ":$G$" & lastRow & "=""" & tag & """))"
End Sub

Private Sub AddButtons(ws As Worksheet)
    Dim r As Long: r = 3
    AddOneButton ws, "Reset", r, 28, RGB(64, 64, 64), RGB(32, 32, 32), vbWhite, "ResetChecklist":                r = r + 2
    AddOneButton ws, "Add To Trades", r, 28, RGB(0, 120, 215), RGB(0, 84, 150), vbWhite, "AddToTrades":          r = r + 2
    AddOneButton ws, "Recalculate", r, 24, RGB(100, 100, 100), RGB(60, 60, 60), vbWhite, "RecalcChecklist":      r = r + 2
    AddOneButton ws, "Apply Preset", r, 24, RGB(60, 160, 80), RGB(40, 120, 60), vbWhite, "ApplyPreset":          r = r + 2
    AddOneButton ws, "Save Decision", r, 28, RGB(110, 70, 170), RGB(80, 40, 140), vbWhite, "SaveDecision"
End Sub

Private Sub AddOneButton(ws As Worksheet, caption As String, rowTop As Long, _
                         h As Single, fillRGB As Long, lineRGB As Long, fontRGB As Long, _
                         macroName As String)
    Dim anchor As Range: Set anchor = ws.Cells(rowTop, UI_COL)
    Dim shp As Shape
    Set shp = ws.Shapes.AddShape(msoShapeRoundedRectangle, anchor.Left, anchor.Top, 140, h)
    shp.TextFrame.Characters.Text = caption
    StyleButton shp, anchor, fillRGB, lineRGB, fontRGB
    shp.OnAction = macroName
End Sub

Private Sub PlaceButtonsNeatly(ws As Worksheet)
    Dim y As Single, shp As Shape
    y = ws.Cells(3, UI_COL).Top
    For Each shp In ws.Shapes
        If shp.Type = msoAutoShape Then
            Select Case shp.TextFrame.Characters.Text
                Case "Reset", "Add To Trades", "Recalculate", "Apply Preset", "Save Decision"
                    shp.Left = ws.Cells(3, UI_COL).Left
                    shp.Width = 140
                    If shp.TextFrame.Characters.Text = "Recalculate" Or shp.TextFrame.Characters.Text = "Apply Preset" Then
                        shp.Height = 24
                    Else
                        shp.Height = 28
                    End If
                    shp.Top = y
                    y = y + shp.Height + 8
            End Select
        End If
    Next shp
End Sub

Private Sub StyleButton(shp As Shape, anchor As Range, fillRGB As Long, lineRGB As Long, fontRGB As Long)
    shp.Fill.ForeColor.RGB = fillRGB
    shp.Line.ForeColor.RGB = lineRGB
    shp.TextFrame.Characters.Font.Color = fontRGB
    shp.TextFrame.Characters.Font.Bold = True
End Sub

Public Sub RefreshBanner()
    Dim ws As Worksheet: Set ws = Worksheets("Checklist")
    Dim s As String: s = UCase$(Trim$(CStr(ws.Range("J5").Value)))
    With ws.Range("J5")
        .Interior.Pattern = xlSolid
        .Interior.TintAndShade = 0
        .Font.TintAndShade = 0
        Select Case s
            Case "OK TO TRADE": .Interior.Color = RGB(0, 176, 80): .Font.Color = vbWhite
            Case "CAUTION":     .Interior.Color = RGB(255, 192, 0): .Font.Color = vbBlack
            Case Else:          .Interior.Color = RGB(192, 0, 0): .Font.Color = vbWhite
        End Select
        .Font.Bold = True
    End With
End Sub

Private Sub EnsureEventsOn()
    On Error Resume Next
    If Application.EnableEvents = False Then Application.EnableEvents = True
    On Error GoTo 0
End Sub



---



# READMEformatter

Sub FormatReadmeSheet()
    Dim ws As Worksheet
    On Error Resume Next
    Set ws = Worksheets("Readme")
    If ws Is Nothing Then
        Set ws = Worksheets.Add(Before:=Worksheets(1))
        ws.Name = "Readme"
    End If
    On Error GoTo 0

    With ws
        .Columns("A").ColumnWidth = 120
        .Rows.RowHeight = 18
        .Range("A:A").WrapText = True
        .Range("A1").Font.Bold = True
        .Range("A1").Font.Size = 14
        .Activate
        .Range("A3").Select
        ActiveWindow.FreezePanes = True
    End With
End Sub

---



Option Explicit

' =========
' Utilities
' =========
Private Function SheetExists(ByVal nm As String) As Boolean
    On Error Resume Next
    SheetExists = Not Worksheets(nm) Is Nothing
    On Error GoTo 0
End Function

Private Function GetOrCreateSheet(ByVal nm As String) As Worksheet
    If SheetExists(nm) Then
        Set GetOrCreateSheet = Worksheets(nm)
    Else
        Set GetOrCreateSheet = Worksheets.Add(After:=Worksheets(Worksheets.count))
        GetOrCreateSheet.Name = nm
    End If
End Function

Private Function GetOrCreateListObject(ws As Worksheet, _
    ByVal tableName As String, headers As Variant, _
    Optional topLeft As String = "A1") As ListObject

    Dim lo As ListObject
    For Each lo In ws.ListObjects
        If lo.Name = tableName Then Set GetOrCreateListObject = lo: Exit Function
    Next lo
    
    Dim rng As Range
    Set rng = ws.Range(topLeft).Resize(1, UBound(headers) - LBound(headers) + 1)
    rng.Value = headers
    Set lo = ws.ListObjects.Add(xlSrcRange, rng, , xlYes)
    lo.Name = tableName
    lo.Range.Font.Bold = True
    Set GetOrCreateListObject = lo
End Function

Private Function FindCol(lo As ListObject, ByVal headerText As String) As Long
    Dim i As Long
    For i = 1 To lo.ListColumns.count
        If LCase$(Trim$(lo.ListColumns(i).Name)) = LCase$(Trim$(headerText)) Then
            FindCol = i: Exit Function
        End If
    Next i
    FindCol = 0
End Function

Private Function EnsureName(ByVal nm As String, ByVal refersTo As String, Optional defaultValue As Variant) As Name
    Dim n As Name
    For Each n In ThisWorkbook.Names
        If n.Name = nm Then Set EnsureName = n: Exit Function
    Next n
    Set EnsureName = ThisWorkbook.Names.Add(Name:=nm, refersTo:=refersTo)
    If Not IsMissing(defaultValue) Then Range(nm).Value = defaultValue
End Function

Private Function NzD(v, Optional d As Double = 0#) As Double
    If IsError(v) Or Len(Trim$(v & "")) = 0 Then NzD = d Else NzD = CDbl(v)
End Function

Private Function NzS(v, Optional d As String = "") As String
    If IsError(v) Or Len(Trim$(v & "")) = 0 Then NzS = d Else NzS = CStr(v)
End Function

Private Function NormalizeTicker(ByVal raw As String) As String
    Dim s As String: s = UCase$(Trim$(raw))
    s = Replace(s, ",", " ")
    s = Replace(s, ";", " ")
    s = Replace(s, vbTab, " ")
    s = Replace(s, vbCr, " ")
    s = Replace(s, vbLf, " ")
    s = WorksheetFunction.Trim(s)
    NormalizeTicker = s
End Function

Private Function GetBucketForSector(ByVal sector As String, bucketsLo As ListObject) As String
    Dim R As ListRow
    For Each R In bucketsLo.ListRows
        If LCase$(Trim$(R.Range.Cells(1, 1).Value)) = LCase$(Trim$(sector)) Then
            GetBucketForSector = NzS(R.Range.Cells(1, 2).Value)
            Exit Function
        End If
    Next R
    GetBucketForSector = ""
End Function

' ========================
' Setup / Structure seeding
' ========================
Public Sub EnsureWorkbookStructure()
    Dim ws As Worksheet, lo As ListObject

    ' Summary with named settings
    Set ws = GetOrCreateSheet("Summary")
    With ws
        .Cells.ClearFormats
        .Range("A1").Resize(10, 2).Clear
        .Range("A1").Resize(10, 1).Font.Bold = True
        .Range("A1").Value = "Setting"
        .Range("B1").Value = "Value"
        .Range("A2").Resize(8, 1).Value = Application.Transpose(Array( _
            "Equity_E", "RiskPct_r", "StopMultiple_K", "HeatCap_H_pct", _
            "BucketHeatCap_pct", "AddStepN", "EarningsBufferDays", "Notes"))
        .Range("B2").Resize(1, 1).Value = 100000
        .Range("B3").Value = 0.005
        .Range("B4").Value = 2
        .Range("B5").Value = 0.04
        .Range("B6").Value = 0.015
        .Range("B7").Value = 0.5
        .Range("B8").Value = 3
        EnsureName "Equity_E", "=Summary!$B$2"
        EnsureName "RiskPct_r", "=Summary!$B$3"
        EnsureName "StopMultiple_K", "=Summary!$B$4"
        EnsureName "HeatCap_H_pct", "=Summary!$B$5"
        EnsureName "BucketHeatCap_pct", "=Summary!$B$6"
        EnsureName "AddStepN", "=Summary!$B$7"
        EnsureName "EarningsBufferDays", "=Summary!$B$8"
    End With
    
    ' Presets
    Set ws = GetOrCreateSheet("Presets")
    Set lo = GetOrCreateListObject(ws, "tblPresets", Array("Name", "QueryString"), "A1")
    If lo.ListRows.count = 0 Then SeedPresets lo
    
    ' Buckets
    Set ws = GetOrCreateSheet("Buckets")
    Set lo = GetOrCreateListObject(ws, "tblBuckets", _
        Array("Sector", "Bucket", "BucketHeatCapPct", "StopoutsToCooldown", "StopoutsWindowBars", "CooldownBars", "CooldownActive", "CooldownEndsOn"), "A1")
    If lo.ListRows.count = 0 Then
        Dim data, i As Long
        data = Array( _
          Array("Technology", "Tech/Comm", 0.015, 2, 30, 10, False, ""), _
          Array("Communication Services", "Tech/Comm", 0.015, 2, 30, 10, False, ""), _
          Array("Consumer Cyclical", "Consumer", 0.015, 2, 30, 10, False, ""), _
          Array("Consumer Defensive", "Consumer", 0.015, 2, 30, 10, False, ""), _
          Array("Financial", "Financials", 0.015, 2, 30, 10, False, ""), _
          Array("Industrials", "Industrials", 0.015, 2, 30, 10, False, ""), _
          Array("Energy", "Energy/Materials", 0.015, 2, 30, 10, False, ""), _
          Array("Basic Materials", "Energy/Materials", 0.015, 2, 30, 10, False, ""), _
          Array("Utilities", "Defensives/REITs", 0.015, 2, 30, 10, False, ""), _
          Array("Real Estate", "Defensives/REITs", 0.015, 2, 30, 10, False, ""), _
          Array("Healthcare", "Defensives/REITs", 0.015, 2, 30, 10, False, ""))
        For i = LBound(data) To UBound(data)
            lo.ListRows.Add
            lo.ListRows(lo.ListRows.count).Range.Value = data(i)
        Next i
    End If
    
    ' Candidates
    Set ws = GetOrCreateSheet("Candidates")
    Set lo = GetOrCreateListObject(ws, "tblCandidates", _
        Array("Date", "Ticker", "Preset", "Sector", "Bucket", "TVConfirm", "N_ATR", "K", _
              "EntryPrice", "RiskPct_r", "Method", "Delta", "DTE", "MaxLossPerContract", "Notes"), "A1")
    
    ' Checklist inputs + names
    Set ws = GetOrCreateSheet("Checklist")
    With ws
        .Cells.Clear
        .Range("A1").Resize(20, 2).Font.Bold = True
        .Range("A1").Value = "Field": .Range("B1").Value = "Value"
        Dim labels
        labels = Array( _
         "chkTicker", "chkPreset", "chkSector", "chkBucket", _
         "chkN", "chkK", "chkEntry", "chkRiskPct", "chkMethod", _
         "chkDelta", "chkDTE", "chkMaxLoss", _
         "chkFromPreset", "chkTrendPass", "chkLiquidityPass", "chkTVConfirm", "chkEarningsOK", "chkJournalOK", "chkBanner")
        Dim R As Long: R = 2
        Dim i2 As Long
        For i2 = LBound(labels) To UBound(labels)
            .Cells(R, 1).Value = labels(i2)
            EnsureName labels(i2), "=" & ws.Name & "!$B$" & R
            R = R + 1
        Next i2
        Range("chkK").Value = Range("StopMultiple_K").Value
        Range("chkRiskPct").Value = Range("RiskPct_r").Value
        Range("chkMethod").Value = "Opt-DeltaATR" ' default
    End With
    
    ' Decisions
    Set ws = GetOrCreateSheet("Decisions")
    Set lo = GetOrCreateListObject(ws, "tblDecisions", _
        Array("DateTime", "Ticker", "Preset", "Bucket", "N_ATR", "K", "Entry", "RiskPct_r", "R_dollars", _
              "Size_Shares", "Size_Contracts", "Method", "Delta", "DTE", "InitialStop", "Banner", _
              "HeatAtEntry", "BucketHeatPost", "PortHeatPost", "Outcome", "Notes"), "A1")
    
    ' Positions
    Set ws = GetOrCreateSheet("Positions")
    Set lo = GetOrCreateListObject(ws, "tblPositions", _
        Array("Ticker", "Bucket", "OpenDate", "UnitsOpen", "RperUnit", "TotalOpenR", "Method", "Status", "LastAddPrice", "NextAddPrice"), "A1")
    
    ' Review & Control
    Set ws = GetOrCreateSheet("Review"): ws.Cells.Clear
    Set ws = GetOrCreateSheet("Control"): ws.Cells.Clear
    
    MsgBox "Structure ready. Set your Summary values and assign buttons.", vbInformation
End Sub

Private Sub SeedPresets(lo As ListObject)
    Dim presets
    presets = Array( _
      Array("TF_BREAKOUT_LONG", "v=211&p=d&s=ta_newhigh&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pa,ta_sma200_pa&o=-relativevolume"), _
      Array("TF_MOMENTUM_UPTREND", "v=211&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pa,ta_sma200_pa&dr=y1&o=-marketcap"), _
      Array("TF_UNUSUAL_VOLUME", "v=211&p=d&s=ta_unusualvolume&f=cap_largeover,sh_price_o20,ta_sma50_pa,ta_sma200_pa&o=-relativevolume"), _
      Array("TF_BREAKDOWN_SHORT", "v=211&p=d&s=ta_newlow&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pb,ta_sma200_pb&o=-relativevolume"), _
      Array("TF_MOMENTUM_DOWNTREND", "v=211&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pb,ta_sma200_pb&dr=y1&o=-marketcap"))
    Dim i As Long
    For i = LBound(presets) To UBound(presets)
        lo.ListRows.Add
        lo.ListRows(lo.ListRows.count).Range.Cells(1, 1).Value = presets(i)(0)
        lo.ListRows(lo.ListRows.count).Range.Cells(1, 2).Value = presets(i)(1)
    Next i
End Sub

' =============
' Preset helpers
' =============
Public Sub OpenPreset()
    Dim lo As ListObject: Set lo = Worksheets("Presets").ListObjects("tblPresets")
    Dim nm As String
    nm = Application.InputBox("Preset name (e.g., TF_BREAKOUT_LONG)", Type:=2)
    If nm = "False" Or Len(nm) = 0 Then Exit Sub
    Dim R As ListRow, url As String
    For Each R In lo.ListRows
        If UCase$(R.Range.Cells(1, 1).Value) = UCase$(nm) Then
            url = "https://finviz.com/screener.ashx?" & R.Range.Cells(1, 2).Value
            ThisWorkbook.FollowHyperlink url
            Exit Sub
        End If
    Next R
    MsgBox "Preset not found.", vbExclamation
End Sub

Public Sub ImportCandidatesPrompt()
    Dim txt As String
    txt = Application.InputBox("Paste tickers from FINVIZ (one per line, or comma-separated):", Type:=2)
    If txt = "False" Or Len(Trim$(txt)) = 0 Then Exit Sub
    Dim presetName As String
    presetName = Application.InputBox("Which preset did these come from?", Type:=2)
    If presetName = "False" Then presetName = "(unspecified)"
    ImportCandidatesFromText txt, presetName
End Sub

Public Sub ImportCandidatesFromText(ByVal txt As String, ByVal presetName As String)
    Dim ws As Worksheet: Set ws = Worksheets("Candidates")
    Dim lo As ListObject: Set lo = ws.ListObjects("tblCandidates")

    Dim tokens() As String, t As Variant
    txt = Replace(txt, vbCrLf, vbLf)
    txt = Replace(txt, vbCr, vbLf)
    txt = Replace(txt, ",", vbLf)
    txt = Replace(txt, ";", vbLf)
    tokens = Split(txt, vbLf)
    
    Dim dict As Object: Set dict = CreateObject("Scripting.Dictionary")
    Dim s As String
    For Each t In tokens
        s = NormalizeTicker(t)
        If Len(s) > 0 Then If Not dict.Exists(s) Then dict.Add s, 1
    Next t
    
    If dict.count = 0 Then MsgBox "No tickers recognized.", vbExclamation: Exit Sub
    
    Dim bucketsLo As ListObject
    Set bucketsLo = Worksheets("Buckets").ListObjects("tblBuckets")
    
    Dim K As Variant, newRow As ListRow
    For Each K In dict.Keys
        ' de-dupe existing
        If Not CandidateExists(lo, CStr(K)) Then
            Set newRow = lo.ListRows.Add
            With newRow.Range
                .Cells(1, 1).Value = Date
                .Cells(1, 2).Value = K
                .Cells(1, 3).Value = presetName
                .Cells(1, 6).Value = False ' TVConfirm default
                .Cells(1, 8).Value = Range("StopMultiple_K").Value ' K default
                .Cells(1, 10).Value = Range("RiskPct_r").Value ' Risk default
            End With
        End If
    Next K
    
    ws.Columns.AutoFit
    MsgBox dict.count & " candidate(s) imported.", vbInformation
End Sub

Private Function CandidateExists(lo As ListObject, ByVal ticker As String) As Boolean
    Dim R As ListRow
    For Each R In lo.ListRows
        If UCase$(NzS(R.Range.Cells(1, 2).Value)) = UCase$(ticker) Then CandidateExists = True: Exit Function
    Next R
    CandidateExists = False
End Function

' ======================
' Checklist & Sizing API
' ======================
Public Sub EvaluateChecklist()
    ' Compute Banner and bucket from sector
    Dim sector As String: sector = NzS(Range("chkSector").Value)
    Dim bucketsLo As ListObject: Set bucketsLo = Worksheets("Buckets").ListObjects("tblBuckets")
    If Len(NzS(Range("chkBucket").Value)) = 0 And Len(sector) > 0 Then
        Range("chkBucket").Value = GetBucketForSector(sector, bucketsLo)
    End If

    Dim score As Long, req As Long
    score = 0: req = 0
    req = req + 1: If CBool(NzD(Range("chkFromPreset").Value, 0)) Then score = score + 1
    req = req + 1: If CBool(NzD(Range("chkTrendPass").Value, 0)) Then score = score + 1
    req = req + 1: If CBool(NzD(Range("chkLiquidityPass").Value, 0)) Then score = score + 1
    req = req + 1: If CBool(NzD(Range("chkTVConfirm").Value, 0)) Then score = score + 1
    req = req + 1: If CBool(NzD(Range("chkEarningsOK").Value, 0)) Then score = score + 1
    req = req + 1: If CBool(NzD(Range("chkJournalOK").Value, 0)) Then score = score + 1
    
    Dim banner As String
    If score = req Then
        banner = "GREEN"
    ElseIf score >= req - 1 Then
        banner = "YELLOW"
    Else
        banner = "RED"
    End If
    Range("chkBanner").Value = banner
    Worksheets("Checklist").Range("B1").Interior.Color = IIf(banner = "GREEN", RGB(198, 239, 206), IIf(banner = "YELLOW", RGB(255, 235, 156), RGB(255, 199, 206)))
End Sub

Public Sub RecalcSizing()
    Dim E As Double: E = NzD(Range("Equity_E").Value, 0)
    Dim rPct As Double: rPct = NzD(Range("chkRiskPct").Value, NzD(Range("RiskPct_r").Value, 0.005))
    Dim Nval As Double: Nval = NzD(Range("chkN").Value, 0)
    Dim K As Double: K = NzD(Range("chkK").Value, NzD(Range("StopMultiple_K").Value, 2))
    Dim entry As Double: entry = NzD(Range("chkEntry").Value, 0)
    Dim method As String: method = UCase$(NzS(Range("chkMethod").Value))

    Dim R As Double: R = E * rPct
    Dim stopDist As Double: stopDist = K * Nval
    Dim qtyShares As Long: If stopDist > 0 Then qtyShares = Int(R / stopDist)
    
    Dim qtyContracts As Long
    If method = UCase$("Opt-DeltaATR") Then
        Dim delta As Double: delta = NzD(Range("chkDelta").Value, 0.6)
        If K * Nval * delta * 100 > 0 Then qtyContracts = Int(R / (K * Nval * delta * 100))
    ElseIf method = UCase$("Opt-MaxLoss") Then
        Dim maxLoss As Double: maxLoss = NzD(Range("chkMaxLoss").Value, 0)
        If maxLoss > 0 Then qtyContracts = Int(R / (maxLoss * 100))
    End If
    
    Dim initStop As Double: If entry > 0 Then initStop = entry - stopDist
    
    With Worksheets("Checklist")
        .Range("D2").Value = "Computed R ($)": .Range("E2").Value = R
        .Range("D3").Value = "StopDist (K×N)": .Range("E3").Value = stopDist
        .Range("D4").Value = "Shares": .Range("E4").Value = qtyShares
        .Range("D5").Value = "Contracts": .Range("E5").Value = qtyContracts
        .Range("D6").Value = "Initial Stop": .Range("E6").Value = initStop
        .Range("D7").Value = "Add1": .Range("E7").Value = IIf(entry > 0, entry + NzD(Range("AddStepN").Value, 0.5) * Nval, "")
        .Range("D8").Value = "Add2": .Range("E8").Value = IIf(entry > 0, entry + 2 * NzD(Range("AddStepN").Value, 0.5) * Nval, "")
        .Range("D9").Value = "Add3": .Range("E9").Value = IIf(entry > 0, entry + 3 * NzD(Range("AddStepN").Value, 0.5) * Nval, "")
    End With
End Sub

' ==================
' Heat & Cooldown API
' ==================
Private Function PortfolioHeatAfter(ByVal addR As Double) As Double
    Dim lo As ListObject: Set lo = Worksheets("Positions").ListObjects("tblPositions")
    Dim R As ListRow, sumR As Double
    For Each R In lo.ListRows
        If LCase$(NzS(R.Range.Cells(1, 8).Value)) <> "closed" Then
            sumR = sumR + NzD(R.Range.Cells(1, 6).Value)
        End If
    Next R
    PortfolioHeatAfter = sumR + addR
End Function

Private Function BucketHeatAfter(ByVal bucket As String, ByVal addR As Double) As Double
    Dim lo As ListObject: Set lo = Worksheets("Positions").ListObjects("tblPositions")
    Dim R As ListRow, sumR As Double
    For Each R In lo.ListRows
        If LCase$(NzS(R.Range.Cells(1, 8).Value)) <> "closed" Then
            If UCase$(NzS(R.Range.Cells(1, 2).Value)) = UCase$(bucket) Then
                sumR = sumR + NzD(R.Range.Cells(1, 6).Value)
            End If
        End If
    Next R
    BucketHeatAfter = sumR + addR
End Function

Private Function BucketCapPct(ByVal bucket As String) As Double
    Dim lo As ListObject: Set lo = Worksheets("Buckets").ListObjects("tblBuckets")
    Dim R As ListRow
    For Each R In lo.ListRows
        If UCase$(NzS(R.Range.Cells(1, 2).Value)) = UCase$(bucket) Then
            BucketCapPct = NzD(R.Range.Cells(1, 3).Value, NzD(Range("BucketHeatCap_pct").Value, 0.015))
            Exit Function
        End If
    Next R
    BucketCapPct = NzD(Range("BucketHeatCap_pct").Value, 0.015)
End Function

Private Function IsBucketInCooldown(ByVal bucket As String) As Boolean
    Dim lo As ListObject: Set lo = Worksheets("Buckets").ListObjects("tblBuckets")
    Dim R As ListRow
    For Each R In lo.ListRows
        If UCase$(NzS(R.Range.Cells(1, 2).Value)) = UCase$(bucket) Then
            If CBool(NzD(R.Range.Cells(1, 7).Value, 0)) Then
                If NzS(R.Range.Cells(1, 8).Value) <> "" Then
                    IsBucketInCooldown = (Date <= CDate(R.Range.Cells(1, 8).Value))
                Else
                    IsBucketInCooldown = True
                End If
            End If
            Exit Function
        End If
    Next R
    IsBucketInCooldown = False
End Function

Public Sub UpdateCooldowns()
    ' If a bucket has >=StopoutsToCooldown within StopoutsWindowBars days ? activate Cooldown for CooldownBars days
    Dim dec As ListObject: Set dec = Worksheets("Decisions").ListObjects("tblDecisions")
    Dim b As ListObject: Set b = Worksheets("Buckets").ListObjects("tblBuckets")

    Dim R As ListRow, dict As Object: Set dict = CreateObject("Scripting.Dictionary")
    Dim cutoff As Date: cutoff = Date - 30
    
    ' build counts of recent stop-outs by bucket
    For Each R In dec.ListRows
        If NzS(R.Range.Cells(1, 20).Value) = "StopOut" Then
            If CDate(R.Range.Cells(1, 1).Value) >= cutoff Then
                Dim bk As String: bk = NzS(R.Range.Cells(1, 4).Value)
                If Not dict.Exists(bk) Then dict.Add bk, 1 Else dict(bk) = dict(bk) + 1
            End If
        End If
    Next R
    
    ' write cooldown flags/dates
    Dim row As ListRow, stopReq As Long, bars As Long
    For Each row In b.ListRows
        stopReq = CLng(NzD(row.Range.Cells(1, 4).Value, 2))
        bars = CLng(NzD(row.Range.Cells(1, 6).Value, 10))
        Dim cnt As Long: cnt = 0
        If dict.Exists(NzS(row.Range.Cells(1, 2).Value)) Then cnt = dict(NzS(row.Range.Cells(1, 2).Value))
        If cnt >= stopReq Then
            row.Range.Cells(1, 7).Value = True
            row.Range.Cells(1, 8).Value = Date + bars
        End If
    Next row
    
    MsgBox "Cooldowns updated.", vbInformation
End Sub

' ===============
' Impulse brake (2m)
' ===============
Private Sub StartImpulseTimer()
    Worksheets("Control").Range("A1").Value = Now
End Sub

Private Function ImpulseElapsed() As Boolean
    Dim t As Variant: t = Worksheets("Control").Range("A1").Value
    If IsDate(t) Then
        ImpulseElapsed = (Now >= DateAdd("n", 2, CDate(t)))
    Else
        ImpulseElapsed = True ' first time
    End If
End Function

' ===============
' Save a decision
' ===============
Public Sub SaveDecision()
    ' Hard gate: must be GREEN, in Candidates, within heat caps, bucket not in cooldown, and 2-minute delay post-Evaluate
    Dim banner As String: banner = UCase$(NzS(Range("chkBanner").Value))
    If banner <> "GREEN" Then MsgBox "Banner must be GREEN.", vbExclamation: Exit Sub

    Dim ticker As String: ticker = NzS(Range("chkTicker").Value)
    Dim preset As String: preset = NzS(Range("chkPreset").Value)
    Dim bucket As String: bucket = NzS(Range("chkBucket").Value)
    If Len(ticker) = 0 Or Len(bucket) = 0 Then MsgBox "Ticker & Bucket required.", vbExclamation: Exit Sub
    
    If Not ImpulseElapsed() Then MsgBox "2-minute cool-off not elapsed.", vbExclamation: Exit Sub
    
    If IsBucketInCooldown(bucket) Then MsgBox "Bucket is in cooldown.", vbExclamation: Exit Sub
    
    ' Sizing inputs
    Dim E As Double: E = NzD(Range("Equity_E").Value, 0)
    Dim rPct As Double: rPct = NzD(Range("chkRiskPct").Value, NzD(Range("RiskPct_r").Value, 0.005))
    Dim Nval As Double: Nval = NzD(Range("chkN").Value, 0)
    Dim K As Double: K = NzD(Range("chkK").Value, NzD(Range("StopMultiple_K").Value, 2))
    Dim entry As Double: entry = NzD(Range("chkEntry").Value, 0)
    Dim method As String: method = UCase$(NzS(Range("chkMethod").Value))
    
    Dim R As Double: R = E * rPct
    Dim stopDist As Double: stopDist = K * Nval
    Dim qtyShares As Long: If stopDist > 0 Then qtyShares = Int(R / stopDist)
    Dim qtyContracts As Long
    If method = "OPT-DELTAATR" Then
        Dim delta As Double: delta = NzD(Range("chkDelta").Value, 0.6)
        If K * Nval * delta * 100 > 0 Then qtyContracts = Int(R / (K * Nval * delta * 100))
    ElseIf method = "OPT-MAXLOSS" Then
        Dim maxLoss As Double: maxLoss = NzD(Range("chkMaxLoss").Value, 0)
        If maxLoss > 0 Then qtyContracts = Int(R / (maxLoss * 100))
    End If
    
    ' Heat checks
    Dim portHeatPost As Double: portHeatPost = PortfolioHeatAfter(R)
    Dim portCap As Double: portCap = NzD(Range("HeatCap_H_pct").Value, 0.04) * E
    If portHeatPost > portCap Then
        MsgBox "Portfolio heat would exceed cap.", vbExclamation: Exit Sub
    End If
    
    Dim bucketHeatPost As Double: bucketHeatPost = BucketHeatAfter(bucket, R)
    Dim bucketCap As Double: bucketCap = BucketCapPct(bucket) * E
    If bucketHeatPost > bucketCap Then
        MsgBox "Bucket heat would exceed cap.", vbExclamation: Exit Sub
    End If
    
    ' Append to Decisions
    Dim loD As ListObject: Set loD = Worksheets("Decisions").ListObjects("tblDecisions")
    Dim rw As ListRow: Set rw = loD.ListRows.Add
    With rw.Range
        .Cells(1, 1).Value = Now
        .Cells(1, 2).Value = ticker
        .Cells(1, 3).Value = preset
        .Cells(1, 4).Value = bucket
        .Cells(1, 5).Value = Nval
        .Cells(1, 6).Value = K
        .Cells(1, 7).Value = entry
        .Cells(1, 8).Value = rPct
        .Cells(1, 9).Value = R
        .Cells(1, 10).Value = qtyShares
        .Cells(1, 11).Value = qtyContracts
        .Cells(1, 12).Value = method
        .Cells(1, 13).Value = IIf(method = "OPT-DELTAATR", NzD(Range("chkDelta").Value, 0.6), "")
        .Cells(1, 14).Value = NzD(Range("chkDTE").Value, "")
        .Cells(1, 15).Value = IIf(entry > 0, entry - stopDist, "")
        .Cells(1, 16).Value = banner
        .Cells(1, 17).Value = PortfolioHeatAfter(0)
        .Cells(1, 18).Value = bucketHeatPost
        .Cells(1, 19).Value = portHeatPost
        .Cells(1, 20).Value = "" ' Outcome blank until closed
        .Cells(1, 21).Value = NzS(Range("chkTicker").Parent.Range("E1").Value) ' optional note placeholder
    End With
    
    ' Update/open position
    Dim loP As ListObject: Set loP = Worksheets("Positions").ListObjects("tblPositions")
    Dim pr As ListRow, found As Boolean
    For Each pr In loP.ListRows
        If UCase$(NzS(pr.Range.Cells(1, 1).Value)) = UCase$(ticker) Then
            pr.Range.Cells(1, 4).Value = NzD(pr.Range.Cells(1, 4).Value, 0) + 1
            pr.Range.Cells(1, 5).Value = R ' R per unit
            pr.Range.Cells(1, 6).Value = NzD(pr.Range.Cells(1, 6).Value, 0) + R
            pr.Range.Cells(1, 10).Value = entry + NzD(Range("AddStepN").Value, 0.5) * Nval
            found = True: Exit For
        End If
    Next pr
    If Not found Then
        Set pr = loP.ListRows.Add
        With pr.Range
            .Cells(1, 1).Value = ticker
            .Cells(1, 2).Value = bucket
            .Cells(1, 3).Value = Date
            .Cells(1, 4).Value = 1
            .Cells(1, 5).Value = R
            .Cells(1, 6).Value = R
            .Cells(1, 7).Value = method
            .Cells(1, 8).Value = "Open"
            .Cells(1, 9).Value = entry
            .Cells(1, 10).Value = entry + NzD(Range("AddStepN").Value, 0.5) * Nval
        End With
    End If
    
    StartImpulseTimer ' reset the cool-off clock for next action
    
    MsgBox "Decision saved + position updated.", vbInformation
End Sub

' ============
' Review (light)
' ============
Public Sub RefreshReview()
    Dim ws As Worksheet: Set ws = Worksheets("Review")
    ws.Cells.Clear
    ws.Range("A1").Value = "KPI": ws.Range("B1").Value = "Value"

    Dim dec As ListObject: Set dec = Worksheets("Decisions").ListObjects("tblDecisions")
    Dim last30 As Long, i As Long
    last30 = Application.WorksheetFunction.Max(1, dec.ListRows.count - 29)
    
    Dim impulseCnt As Long, totalCnt As Long
    For i = last30 To dec.ListRows.count
        totalCnt = totalCnt + 1
        If UCase$(NzS(dec.ListRows(i).Range.Cells(1, 16).Value)) <> "GREEN" Then impulseCnt = impulseCnt + 1
    Next i
    
    ws.Range("A2").Value = "% taken non-GREEN (last 30)"
    If totalCnt > 0 Then ws.Range("B2").Value = impulseCnt / totalCnt Else ws.Range("B2").Value = 0
    
    ws.Columns.AutoFit
    MsgBox "Review updated.", vbInformation
End Sub

' -------------
' Convenience stub (optional)
' -------------
Public Sub Setup_Workbook()
    EnsureWorkbookStructure
End Sub
