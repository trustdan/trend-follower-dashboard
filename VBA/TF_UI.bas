Attribute VB_Name = "TF_UI"
' ============================================================================
' Module: TF_UI
' Purpose: UI controls, checklist evaluation, sizing, and save logic
' ============================================================================

Option Explicit

' ----------------------------------------------------------------------------
' Sub: BindControls
' Sets up data validation dropdowns and links controls
' Call this when TradeEntry sheet is activated
' ----------------------------------------------------------------------------
Sub BindControls()
    On Error Resume Next

    With Sheets("TradeEntry")
        ' Clear existing validation
        .Range("B5:B8").Validation.Delete

        ' Preset dropdown
        .Range("B5").Validation.Add Type:=xlValidateList, _
            Formula1:="=tblPresets[Name]"

        ' Ticker dropdown (from Candidates table)
        .Range("B6").Validation.Add Type:=xlValidateList, _
            Formula1:="=tblCandidates[Ticker]", _
            AlertStyle:=xlValidAlertStop

        ' Sector dropdown
        .Range("B7").Validation.Add Type:=xlValidateList, _
            Formula1:="Technology,Healthcare,Financials,Consumer,Industrials,Energy"

        ' Bucket dropdown
        .Range("B8").Validation.Add Type:=xlValidateList, _
            Formula1:="=tblBuckets[Bucket]"

        ' Set default Method to Stock (1)
        If IsEmpty(.Range("C13").Value) Then
            .Range("C13").Value = 1
        End If

        ' Show/hide method-specific fields
        Call ToggleMethodFields
    End With

    On Error GoTo 0
End Sub

' ----------------------------------------------------------------------------
' Sub: ToggleMethodFields
' Shows/hides Delta, DTE, MaxLoss rows based on Method choice
' ----------------------------------------------------------------------------
Sub ToggleMethodFields()
    Dim methodChoice As Integer

    On Error Resume Next
    methodChoice = Sheets("TradeEntry").Range("C13").Value
    On Error GoTo 0

    With Sheets("TradeEntry")
        Select Case methodChoice
            Case 1  ' Stock
                .Rows("16:18").Hidden = True

            Case 2  ' Opt-DeltaATR
                .Rows("16:17").Hidden = False
                .Rows("18:18").Hidden = True

            Case 3  ' Opt-MaxLoss
                .Rows("17:18").Hidden = False
                .Rows("16:16").Hidden = True

            Case Else
                .Rows("16:18").Hidden = True
        End Select
    End With
End Sub

' ----------------------------------------------------------------------------
' Sub: EvaluateChecklist
' Reads checklist, computes GO/NO-GO banner, starts impulse timer if GREEN
' ----------------------------------------------------------------------------
Sub EvaluateChecklist()
    Dim checks(1 To 6) As Boolean
    Dim missingCount As Integer
    Dim reasons As String
    Dim banner As String
    Dim bannerColor As Long

    Application.ScreenUpdating = False

    With Sheets("TradeEntry")
        ' Read checklist values from C20:C25
        checks(1) = .Range("C20").Value  ' FromPreset
        checks(2) = .Range("C21").Value  ' TrendPass
        checks(3) = .Range("C22").Value  ' LiquidityPass
        checks(4) = .Range("C23").Value  ' TVConfirm
        checks(5) = .Range("C24").Value  ' EarningsOK
        checks(6) = .Range("C25").Value  ' JournalOK
    End With

    ' Count missing checks and build reason string
    missingCount = 0
    reasons = ""

    If Not checks(1) Then: missingCount = missingCount + 1: reasons = reasons & "FromPreset, "
    If Not checks(2) Then: missingCount = missingCount + 1: reasons = reasons & "TrendPass, "
    If Not checks(3) Then: missingCount = missingCount + 1: reasons = reasons & "LiquidityPass, "
    If Not checks(4) Then: missingCount = missingCount + 1: reasons = reasons & "TVConfirm, "
    If Not checks(5) Then: missingCount = missingCount + 1: reasons = reasons & "EarningsOK, "
    If Not checks(6) Then: missingCount = missingCount + 1: reasons = reasons & "JournalOK, "

    ' Remove trailing comma
    If Len(reasons) > 0 Then
        reasons = Left(reasons, Len(reasons) - 2)
    End If

    ' Determine banner color and text
    If missingCount = 0 Then
        banner = "GREEN - GO"
        bannerColor = RGB(0, 200, 0)

        ' Start impulse timer
        Sheets("Control").Range("A1").Value = Now

    ElseIf missingCount = 1 Then
        banner = "YELLOW - CAUTION"
        bannerColor = RGB(255, 200, 0)

    Else
        banner = "RED - NO-GO"
        bannerColor = RGB(255, 0, 0)
    End If

    ' Write banner and reasons
    With Sheets("TradeEntry")
        .Range("A2:F2").Merge
        .Range("A2").Value = banner
        .Range("A2").Interior.Color = bannerColor
        .Range("A2").Font.Color = RGB(255, 255, 255)
        .Range("A2").Font.Bold = True
        .Range("A2").Font.Size = 14
        .Range("A2").HorizontalAlignment = xlCenter

        .Range("A3:F3").Merge
        If reasons <> "" Then
            .Range("A3").Value = "Missing: " & reasons
        Else
            .Range("A3").Value = "All checks passed!"
        End If
        .Range("A3").Font.Size = 10
        .Range("A3").HorizontalAlignment = xlCenter
    End With

    Application.ScreenUpdating = True

    MsgBox "Evaluation complete: " & banner, vbInformation
End Sub

' ----------------------------------------------------------------------------
' Sub: RecalcSizing
' Computes position size based on method (Stock/Opt-DeltaATR/Opt-MaxLoss)
' ----------------------------------------------------------------------------
Sub RecalcSizing()
    Dim E As Double, r As Double, R As Double
    Dim entry As Double, N As Double, K As Double
    Dim methodChoice As Integer
    Dim delta As Double, maxLoss As Double
    Dim shares As Long, contracts As Long
    Dim stopDist As Double, initialStop As Double

    On Error GoTo ErrorHandler

    ' Read named ranges
    E = Range("Equity_E").Value
    r = Range("RiskPct_r").Value
    R = E * r

    ' Read inputs from TradeEntry
    With Sheets("TradeEntry")
        entry = NzD(.Range("B9").Value, 0)
        N = NzD(.Range("B10").Value, 0)
        K = NzD(.Range("B11").Value, 2)
        methodChoice = NzD(.Range("C13").Value, 1)
        delta = NzD(.Range("B16").Value, 0.5)
        maxLoss = NzD(.Range("B18").Value, 0)
    End With

    ' Validate inputs
    If entry = 0 Or N = 0 Then
        MsgBox "Please enter Entry Price and ATR N before calculating sizing.", vbExclamation
        Exit Sub
    End If

    ' Common calculations
    stopDist = K * N
    initialStop = entry - stopDist

    ' Method-specific sizing
    Select Case methodChoice
        Case 1  ' Stock
            If stopDist > 0 Then
                shares = WorksheetFunction.Floor_Precise(R / stopDist)
            Else
                shares = 0
            End If
            contracts = 0

        Case 2  ' Opt-DeltaATR
            If delta > 0 And N > 0 Then
                contracts = WorksheetFunction.Floor_Precise(R / (K * N * delta * 100))
            Else
                contracts = 0
            End If
            shares = 0

        Case 3  ' Opt-MaxLoss
            If maxLoss > 0 Then
                contracts = WorksheetFunction.Floor_Precise(R / (maxLoss * 100))
            Else
                contracts = 0
            End If
            shares = 0

        Case Else
            shares = 0
            contracts = 0
    End Select

    ' Write outputs
    With Sheets("TradeEntry")
        .Range("F5").Value = R
        .Range("F6").Value = stopDist
        .Range("F7").Value = initialStop
        .Range("F8").Value = shares
        .Range("F9").Value = contracts

        ' Add levels computed via formulas (set formulas if not already present)
        If Not .Range("F10").HasFormula Then
            .Range("F10").Formula = "=B9+(AddStepN*B10)"
            .Range("F11").Formula = "=B9+(2*AddStepN*B10)"
            .Range("F12").Formula = "=B9+(3*AddStepN*B10)"
        End If
    End With

    MsgBox "Sizing calculated:" & vbCrLf & _
           "R: $" & Format(R, "0.00") & vbCrLf & _
           "Stop Dist: " & Format(stopDist, "0.00") & vbCrLf & _
           "Initial Stop: " & Format(initialStop, "0.00") & vbCrLf & _
           "Shares: " & shares & vbCrLf & _
           "Contracts: " & contracts, _
           vbInformation

    Exit Sub

ErrorHandler:
    MsgBox "Error calculating sizing: " & Err.Description, vbCritical
End Sub

' ----------------------------------------------------------------------------
' Sub: SaveDecision
' Hard-gates validation and saves decision to log + positions
' ----------------------------------------------------------------------------
Sub SaveDecision()
    Dim banner As String, ticker As String, bucket As String
    Dim addR As Double, portHeat As Double, buckHeat As Double
    Dim portCap As Double, buckCap As Double
    Dim timerStart As Date, elapsed As Double

    On Error GoTo ErrorHandler

    Application.ScreenUpdating = False

    ' Read key values
    banner = Sheets("TradeEntry").Range("A2").Value
    ticker = Sheets("TradeEntry").Range("B6").Value
    bucket = Sheets("TradeEntry").Range("B8").Value
    addR = NzD(Sheets("TradeEntry").Range("F5").Value, 0)

    ' === HARD GATE 1: Banner must be GREEN ===
    If InStr(banner, "GREEN") = 0 Then
        MsgBox "Cannot save: Banner is not GREEN." & vbCrLf & vbCrLf & _
               "Current banner: " & banner, vbCritical, "BLOCKED"
        GoTo CleanExit
    End If

    ' === HARD GATE 2: Ticker must be in today's Candidates ===
    If Not IsTickerInCandidates(ticker, Date) Then
        MsgBox "Cannot save: Ticker '" & ticker & "' not in today's Candidates." & vbCrLf & vbCrLf & _
               "Import candidates first using the Import Candidates button.", vbCritical, "BLOCKED"
        GoTo CleanExit
    End If

    ' === HARD GATE 3: Impulse timer (2 minutes) ===
    timerStart = Sheets("Control").Range("A1").Value
    If IsDate(timerStart) Then
        elapsed = (Now - timerStart) * 24 * 60  ' Convert to minutes
        If elapsed < 2 Then
            MsgBox "Cannot save: 2-minute cool-off not elapsed." & vbCrLf & vbCrLf & _
                   "Elapsed: " & Format(elapsed, "0.0") & " minutes" & vbCrLf & _
                   "Remaining: " & Format(2 - elapsed, "0.0") & " minutes", vbCritical, "BLOCKED"
            GoTo CleanExit
        End If
    Else
        MsgBox "Cannot save: Impulse timer not started." & vbCrLf & vbCrLf & _
               "Click Evaluate first to start the timer.", vbCritical, "BLOCKED"
        GoTo CleanExit
    End If

    ' === HARD GATE 4: Bucket cooldown ===
    If IsBucketInCooldown(bucket) Then
        MsgBox "Cannot save: Bucket '" & bucket & "' is in cooldown." & vbCrLf & vbCrLf & _
               "Check the Buckets sheet for cooldown end date.", vbCritical, "BLOCKED"
        GoTo CleanExit
    End If

    ' === HARD GATE 5: Portfolio heat cap ===
    portHeat = PortfolioHeatAfter(addR)
    portCap = Range("HeatCap_H_pct").Value * Range("Equity_E").Value

    If portHeat > portCap Then
        MsgBox "Cannot save: Portfolio heat would exceed cap." & vbCrLf & vbCrLf & _
               "Portfolio Heat: $" & Format(portHeat, "0.00") & vbCrLf & _
               "Portfolio Cap: $" & Format(portCap, "0.00") & vbCrLf & _
               "Over by: $" & Format(portHeat - portCap, "0.00"), vbCritical, "BLOCKED"
        GoTo CleanExit
    End If

    ' === HARD GATE 6: Bucket heat cap ===
    buckHeat = BucketHeatAfter(bucket, addR)
    buckCap = Range("BucketHeatCap_pct").Value * Range("Equity_E").Value

    If buckHeat > buckCap Then
        MsgBox "Cannot save: Bucket heat would exceed cap." & vbCrLf & vbCrLf & _
               "Bucket Heat: $" & Format(buckHeat, "0.00") & vbCrLf & _
               "Bucket Cap: $" & Format(buckCap, "0.00") & vbCrLf & _
               "Over by: $" & Format(buckHeat - buckCap, "0.00"), vbCritical, "BLOCKED"
        GoTo CleanExit
    End If

    ' === ALL GATES PASSED - SAVE DECISION ===
    Call AppendDecisionRow
    Call UpdatePositions

    ' Clear banner
    With Sheets("TradeEntry")
        .Range("A2").Value = ""
        .Range("A2").Interior.ColorIndex = xlNone
        .Range("A3").Value = ""
    End With

    ' Reset timer
    Sheets("Control").Range("A1").Value = ""

    MsgBox "Decision saved successfully!" & vbCrLf & vbCrLf & _
           "Ticker: " & ticker & vbCrLf & _
           "Bucket: " & bucket & vbCrLf & _
           "R: $" & Format(addR, "0.00") & vbCrLf & _
           "Portfolio Heat: $" & Format(portHeat, "0.00") & " / $" & Format(portCap, "0.00"), _
           vbInformation, "SUCCESS"

CleanExit:
    Application.ScreenUpdating = True
    Exit Sub

ErrorHandler:
    Application.ScreenUpdating = True
    MsgBox "Error saving decision: " & Err.Description, vbCritical
End Sub

' ----------------------------------------------------------------------------
' Function: IsTickerInCandidates
' Returns: True if ticker exists in Candidates table with given date
' ----------------------------------------------------------------------------
Function IsTickerInCandidates(ticker As String, tradeDate As Date) As Boolean
    Dim tbl As ListObject
    Dim row As ListRow

    On Error Resume Next
    Set tbl = Sheets("Candidates").ListObjects("tblCandidates")
    On Error GoTo 0

    If tbl Is Nothing Then
        IsTickerInCandidates = False
        Exit Function
    End If

    If tbl.ListRows.Count = 0 Then
        IsTickerInCandidates = False
        Exit Function
    End If

    For Each row In tbl.ListRows
        ' Column 2 = Ticker, Column 1 = Date
        If UCase(Trim(row.Range.Cells(1, 2).Value)) = UCase(Trim(ticker)) And _
           Int(row.Range.Cells(1, 1).Value) = Int(tradeDate) Then
            IsTickerInCandidates = True
            Exit Function
        End If
    Next row

    IsTickerInCandidates = False
End Function

' ----------------------------------------------------------------------------
' Sub: AppendDecisionRow
' Appends current trade to tblDecisions with all fields
' ----------------------------------------------------------------------------
Sub AppendDecisionRow()
    Dim tbl As ListObject
    Dim newRow As ListRow

    Set tbl = Sheets("Decisions").ListObjects("tblDecisions")
    Set newRow = tbl.ListRows.Add

    With Sheets("TradeEntry")
        ' Columns: DateTime, Ticker, Preset, Bucket, N_ATR, K, Entry, RiskPct_r, R_dollars,
        '          Size_Shares, Size_Contracts, Method, Delta, DTE, InitialStop, Banner,
        '          HeatAtEntry, BucketHeatPost, PortHeatPost, Outcome, Notes
        newRow.Range.Cells(1, 1).Value = Now  ' DateTime
        newRow.Range.Cells(1, 2).Value = .Range("B6").Value  ' Ticker
        newRow.Range.Cells(1, 3).Value = .Range("B5").Value  ' Preset
        newRow.Range.Cells(1, 4).Value = .Range("B8").Value  ' Bucket
        newRow.Range.Cells(1, 5).Value = .Range("B10").Value  ' N_ATR
        newRow.Range.Cells(1, 6).Value = .Range("B11").Value  ' K
        newRow.Range.Cells(1, 7).Value = .Range("B9").Value  ' Entry
        newRow.Range.Cells(1, 8).Value = Range("RiskPct_r").Value  ' RiskPct_r
        newRow.Range.Cells(1, 9).Value = .Range("F5").Value  ' R_dollars
        newRow.Range.Cells(1, 10).Value = .Range("F8").Value  ' Size_Shares
        newRow.Range.Cells(1, 11).Value = .Range("F9").Value  ' Size_Contracts

        Dim methodText As String
        Select Case .Range("C13").Value
            Case 1: methodText = "Stock"
            Case 2: methodText = "Opt-DeltaATR"
            Case 3: methodText = "Opt-MaxLoss"
            Case Else: methodText = "Unknown"
        End Select
        newRow.Range.Cells(1, 12).Value = methodText  ' Method

        newRow.Range.Cells(1, 13).Value = .Range("B16").Value  ' Delta
        newRow.Range.Cells(1, 14).Value = .Range("B17").Value  ' DTE
        newRow.Range.Cells(1, 15).Value = .Range("F7").Value  ' InitialStop
        newRow.Range.Cells(1, 16).Value = .Range("A2").Value  ' Banner

        Dim addR As Double
        addR = NzD(.Range("F5").Value, 0)
        newRow.Range.Cells(1, 17).Value = PortfolioHeatAfter(0)  ' HeatAtEntry (current, before this trade)
        newRow.Range.Cells(1, 18).Value = BucketHeatAfter(.Range("B8").Value, addR)  ' BucketHeatPost
        newRow.Range.Cells(1, 19).Value = PortfolioHeatAfter(addR)  ' PortHeatPost

        newRow.Range.Cells(1, 20).Value = ""  ' Outcome (filled later)
        newRow.Range.Cells(1, 21).Value = ""  ' Notes
    End With
End Sub

' ----------------------------------------------------------------------------
' Sub: UpdatePositions
' Opens new position or adds to existing in tblPositions
' ----------------------------------------------------------------------------
Sub UpdatePositions()
    Dim tbl As ListObject
    Dim row As ListRow
    Dim ticker As String, bucket As String
    Dim found As Boolean
    Dim newRow As ListRow

    Set tbl = Sheets("Positions").ListObjects("tblPositions")
    ticker = Sheets("TradeEntry").Range("B6").Value
    bucket = Sheets("TradeEntry").Range("B8").Value

    ' Check if position already exists
    found = False
    For Each row In tbl.ListRows
        If row.Range.Cells(1, 1).Value = ticker And _
           row.Range.Cells(1, 8).Value <> "Closed" Then
            found = True

            ' Add to existing position
            row.Range.Cells(1, 4).Value = row.Range.Cells(1, 4).Value + 1  ' UnitsOpen
            row.Range.Cells(1, 6).Value = row.Range.Cells(1, 6).Value + Sheets("TradeEntry").Range("F5").Value  ' TotalOpenR

            Exit For
        End If
    Next row

    ' If not found, create new position
    If Not found Then
        Set newRow = tbl.ListRows.Add

        With Sheets("TradeEntry")
            ' Columns: Ticker, Bucket, OpenDate, UnitsOpen, RperUnit, TotalOpenR, Method, Status, LastAddPrice, NextAddPrice
            newRow.Range.Cells(1, 1).Value = ticker
            newRow.Range.Cells(1, 2).Value = bucket
            newRow.Range.Cells(1, 3).Value = Date
            newRow.Range.Cells(1, 4).Value = 1  ' UnitsOpen
            newRow.Range.Cells(1, 5).Value = .Range("F5").Value  ' RperUnit
            newRow.Range.Cells(1, 6).Value = .Range("F5").Value  ' TotalOpenR
            newRow.Range.Cells(1, 7).Value = IIf(.Range("C13").Value = 1, "Stock", "Options")  ' Method
            newRow.Range.Cells(1, 8).Value = "Open"  ' Status
            newRow.Range.Cells(1, 9).Value = .Range("B9").Value  ' LastAddPrice
            newRow.Range.Cells(1, 10).Value = .Range("F10").Value  ' NextAddPrice (Add1)
        End With
    End If
End Sub

' ----------------------------------------------------------------------------
' Sub: StartImpulseTimer
' Manually starts the 2-minute impulse timer
' ----------------------------------------------------------------------------
Sub StartImpulseTimer()
    Sheets("Control").Range("A1").Value = Now
    MsgBox "2-minute impulse timer started at " & Format(Now, "hh:mm:ss"), vbInformation
End Sub
