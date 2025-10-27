Attribute VB_Name = "TF_UI"
Option Explicit

' ========================================
' TF_UI Module
' UI controls, checklist evaluation, sizing, save decision
' ========================================

Sub BindControls()
    ' Set up Data Validation dropdowns
    Call TF_Logger.WriteLogSection("BindControls() - Start")

    Dim errorMessages As String
    errorMessages = ""

    ' Check if tblPresets exists
    Dim tblExists As Boolean
    On Error Resume Next
    tblExists = Not Worksheets("Presets").ListObjects("tblPresets") Is Nothing
    If Err.Number <> 0 Then
        Call TF_Logger.WriteLog("tblPresets does not exist: " & Err.Description)
        tblExists = False
    Else
        Call TF_Logger.WriteLog("tblPresets exists with " & Worksheets("Presets").ListObjects("tblPresets").ListRows.Count & " rows")
    End If
    On Error GoTo 0

    With Worksheets("TradeEntry")
        ' Preset dropdown
        Call TF_Logger.WriteLog("Setting up Preset dropdown in B5...")
        .Range("B5").Validation.Delete
        On Error Resume Next
        Err.Clear
        .Range("B5").Validation.Add Type:=xlValidateList, _
            Formula1:="=tblPresets[Name]", _
            AlertStyle:=xlValidAlertStop
        If Err.Number <> 0 Then
            errorMessages = errorMessages & "Preset dropdown: " & Err.Description & vbCrLf
            Call TF_Logger.WriteLogError("BindControls", Err.Number, "Preset dropdown failed: " & Err.Description)
            Err.Clear
        Else
            Call TF_Logger.WriteLog("Preset dropdown created successfully")
        End If
        On Error GoTo 0

        ' Ticker dropdown (from Candidates)
        Call TF_Logger.WriteLog("Setting up Ticker dropdown in B6...")
        .Range("B6").Validation.Delete
        On Error Resume Next
        Err.Clear
        .Range("B6").Validation.Add Type:=xlValidateList, _
            Formula1:="=tblCandidates[Ticker]", _
            AlertStyle:=xlValidAlertStop
        If Err.Number <> 0 Then
            Call TF_Logger.WriteLogError("BindControls", Err.Number, "Ticker dropdown failed: " & Err.Description)
            Err.Clear
        Else
            Call TF_Logger.WriteLog("Ticker dropdown created successfully")
        End If
        On Error GoTo 0

        ' Sector dropdown (hardcoded list - always works)
        Call TF_Logger.WriteLog("Setting up Sector dropdown in B7...")
        .Range("B7").Validation.Delete
        On Error Resume Next
        Err.Clear
        .Range("B7").Validation.Add Type:=xlValidateList, _
            Formula1:="Technology,Healthcare,Financials,Consumer,Industrials,Energy", _
            AlertStyle:=xlValidAlertStop
        If Err.Number <> 0 Then
            Call TF_Logger.WriteLogError("BindControls", Err.Number, "Sector dropdown failed: " & Err.Description)
            Err.Clear
        Else
            Call TF_Logger.WriteLog("Sector dropdown created successfully")
        End If
        On Error GoTo 0

        ' Bucket dropdown
        Call TF_Logger.WriteLog("Setting up Bucket dropdown in B8...")
        .Range("B8").Validation.Delete
        On Error Resume Next
        Err.Clear
        .Range("B8").Validation.Add Type:=xlValidateList, _
            Formula1:="=tblBuckets[Bucket]", _
            AlertStyle:=xlValidAlertStop
        If Err.Number <> 0 Then
            Call TF_Logger.WriteLogError("BindControls", Err.Number, "Bucket dropdown failed: " & Err.Description)
            Err.Clear
        Else
            Call TF_Logger.WriteLog("Bucket dropdown created successfully")
        End If
        On Error GoTo 0
    End With

    ' Set default Method to Stock (1)
    On Error Resume Next
    Worksheets("TradeEntry").Range("C13").Value = 1
    Call TF_Logger.WriteLog("Set default method to Stock (1)")
    On Error GoTo 0

    Call ToggleMethodFields

    ' Show error message if any dropdowns failed
    If Len(errorMessages) > 0 Then
        Call TF_Logger.WriteLog("Dropdown errors detected, showing message to user")
        MsgBox "Some dropdowns could not be created:" & vbCrLf & vbCrLf & _
               errorMessages & vbCrLf & _
               "This usually means the data tables haven't been created yet." & vbCrLf & _
               "Run Setup.RunInitialSetup first.", vbExclamation, "Dropdown Setup Warning"
    End If

    Call TF_Logger.WriteLog("BindControls() - Complete")
End Sub

Sub ToggleMethodFields()
    ' Show/hide method-specific fields based on C13 value
    Dim methodChoice As Integer

    methodChoice = NzD(Worksheets("TradeEntry").Range("C13").Value, 1)

    With Worksheets("TradeEntry")
        If methodChoice = 1 Then  ' Stock
            .Rows("16:18").Hidden = True
        ElseIf methodChoice = 2 Then  ' Opt-DeltaATR
            .Rows("16:17").Hidden = False
            .Rows("18:18").Hidden = True
        ElseIf methodChoice = 3 Then  ' Opt-MaxLoss
            .Rows("17:18").Hidden = False
            .Rows("16:16").Hidden = True
        End If
    End With
End Sub

Sub EvaluateChecklist()
    ' GO/NO-GO engine - evaluates all checklist items
    Dim checks(1 To 6) As Boolean
    Dim missingCount As Integer
    Dim reasons As String
    Dim banner As String
    Dim bannerColor As Long

    With Worksheets("TradeEntry")
        ' Read checklist values
        checks(1) = .Range("C20").Value  ' FromPreset
        checks(2) = .Range("C21").Value  ' TrendPass
        checks(3) = .Range("C22").Value  ' LiquidityPass
        checks(4) = .Range("C23").Value  ' TVConfirm
        checks(5) = .Range("C24").Value  ' EarningsOK
        checks(6) = .Range("C25").Value  ' JournalOK
    End With

    ' Count missing checks
    missingCount = 0
    reasons = ""

    If Not checks(1) Then: missingCount = missingCount + 1: reasons = reasons & "FromPreset, "
    If Not checks(2) Then: missingCount = missingCount + 1: reasons = reasons & "TrendPass, "
    If Not checks(3) Then: missingCount = missingCount + 1: reasons = reasons & "LiquidityPass, "
    If Not checks(4) Then: missingCount = missingCount + 1: reasons = reasons & "TVConfirm, "
    If Not checks(5) Then: missingCount = missingCount + 1: reasons = reasons & "EarningsOK, "
    If Not checks(6) Then: missingCount = missingCount + 1: reasons = reasons & "JournalOK, "

    ' Determine banner
    If missingCount = 0 Then
        banner = "GREEN - GO"
        bannerColor = RGB(0, 200, 0)
        ' Start impulse timer
        Worksheets("Control").Range("A1").Value = Now
    ElseIf missingCount = 1 Then
        banner = "YELLOW - CAUTION"
        bannerColor = RGB(255, 200, 0)
    Else
        banner = "RED - NO-GO"
        bannerColor = RGB(255, 0, 0)
    End If

    ' Write banner and reasons
    With Worksheets("TradeEntry")
        .Range("A2:F2").Merge
        .Range("A2").Value = banner
        .Range("A2").Interior.Color = bannerColor
        .Range("A2").Font.Bold = True
        .Range("A2").Font.Size = 14
        .Range("A2").HorizontalAlignment = xlCenter

        If Len(reasons) > 0 Then
            .Range("A3").Value = "Missing: " & Left(reasons, Len(reasons) - 2)
        Else
            .Range("A3").Value = "All checks passed!"
        End If
    End With
End Sub

Sub RecalcSizing()
    ' Calculate position sizing for all methods
    Dim E As Double, r As Double, R As Double
    Dim entry As Double, n As Double, K As Double
    Dim methodChoice As Integer, delta As Double, maxLoss As Double
    Dim shares As Long, contracts As Long
    Dim stopDist As Double, initialStop As Double

    ' Read settings
    E = Range("Equity_E").Value
    r = Range("RiskPct_r").Value
    R = E * r

    With Worksheets("TradeEntry")
        entry = NzD(.Range("B9").Value, 0)
        n = NzD(.Range("B10").Value, 0)
        K = NzD(.Range("B11").Value, 2)
        methodChoice = NzD(.Range("C13").Value, 1)
        delta = NzD(.Range("B16").Value, 0.3)
        maxLoss = NzD(.Range("B18").Value, 100)
    End With

    ' Validate inputs
    If entry = 0 Or n = 0 Then
        MsgBox "Please enter Entry Price and ATR N values", vbExclamation
        Exit Sub
    End If

    ' Common calculations
    stopDist = K * n
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
            If delta > 0 And n > 0 Then
                contracts = WorksheetFunction.Floor_Precise(R / (K * n * delta * 100))
            Else
                contracts = 0
            End If
            shares = 0

        Case 3  ' Opt-MaxLoss
            If maxLoss > 0 Then
                contracts = WorksheetFunction.Floor_Precise(R / maxLoss)
            Else
                contracts = 0
            End If
            shares = 0
    End Select

    ' Write outputs
    With Worksheets("TradeEntry")
        .Range("F5").Value = R
        .Range("F5").NumberFormat = "$#,##0.00"

        .Range("F6").Value = stopDist
        .Range("F6").NumberFormat = "0.00"

        .Range("F7").Value = initialStop
        .Range("F7").NumberFormat = "0.00"

        .Range("F8").Value = shares
        .Range("F8").NumberFormat = "0"

        .Range("F9").Value = contracts
        .Range("F9").NumberFormat = "0"

        ' Add levels (formulas)
        .Range("F10").Formula = "=B9+(AddStepN*B10)"
        .Range("F11").Formula = "=B9+(2*AddStepN*B10)"
        .Range("F12").Formula = "=B9+(3*AddStepN*B10)"
    End With

    ' Calculate and display heat preview
    Call UpdateHeatPreview
End Sub

Sub UpdateHeatPreview()
    ' Update heat preview bars
    Dim addR As Double
    Dim bucket As String
    Dim portHeat As Double, buckHeat As Double
    Dim portCap As Double, buckCap As Double
    Dim portPct As Double, buckPct As Double

    With Worksheets("TradeEntry")
        addR = NzD(.Range("F5").Value, 0)
        bucket = NzS(.Range("B8").Value, "")
    End With

    ' Calculate heat
    portHeat = PortfolioHeatAfter(addR)
    buckHeat = BucketHeatAfter(bucket, addR)

    portCap = Range("HeatCap_H_pct").Value * Range("Equity_E").Value
    buckCap = Range("BucketHeatCap_pct").Value * Range("Equity_E").Value

    portPct = (portHeat / portCap) * 100
    buckPct = (buckHeat / buckCap) * 100

    ' Write heat values
    With Worksheets("TradeEntry")
        .Range("F14").Value = portHeat & " / " & portCap
        .Range("F15").Value = buckHeat & " / " & buckCap

        ' Color code based on percentage
        If portPct <= 70 Then
            .Range("F14").Interior.Color = RGB(0, 200, 0)  ' Green
        ElseIf portPct <= 100 Then
            .Range("F14").Interior.Color = RGB(255, 200, 0)  ' Amber
        Else
            .Range("F14").Interior.Color = RGB(255, 0, 0)  ' Red
        End If

        If buckPct <= 70 Then
            .Range("F15").Interior.Color = RGB(0, 200, 0)
        ElseIf buckPct <= 100 Then
            .Range("F15").Interior.Color = RGB(255, 200, 0)
        Else
            .Range("F15").Interior.Color = RGB(255, 0, 0)
        End If
    End With
End Sub

Sub SaveDecision()
    ' Hard-gate validation and save to Decisions table
    Dim banner As String, ticker As String, bucket As String
    Dim addR As Double, portHeat As Double, buckHeat As Double
    Dim timerStart As Variant, elapsed As Double
    Dim portCap As Double, buckCap As Double

    With Worksheets("TradeEntry")
        banner = .Range("A2").Value
        ticker = .Range("B6").Value
        bucket = .Range("B8").Value
        addR = NzD(.Range("F5").Value, 0)
    End With

    ' Gate 1: Banner must be GREEN
    If InStr(banner, "GREEN") = 0 Then
        MsgBox "Cannot save: Banner is not GREEN" & vbCrLf & _
               "Current: " & banner, vbCritical, "Hard Gate Failed"
        Exit Sub
    End If

    ' Gate 2: Ticker must be in today's Candidates
    If Not IsTickerInCandidates(ticker, Date) Then
        MsgBox "Cannot save: Ticker '" & ticker & "' not in today's Candidates" & vbCrLf & _
               "Import candidates first or verify ticker symbol", vbCritical, "Hard Gate Failed"
        Exit Sub
    End If

    ' Gate 3: Check impulse timer (2 minutes)
    timerStart = Worksheets("Control").Range("A1").Value
    If IsDate(timerStart) Then
        elapsed = (Now - timerStart) * 24 * 60  ' Convert to minutes
        If elapsed < 2 Then
            MsgBox "Cannot save: 2-minute cool-off not elapsed" & vbCrLf & _
                   "Elapsed: " & Format(elapsed, "0.0") & " minutes" & vbCrLf & _
                   "Remaining: " & Format(2 - elapsed, "0.0") & " minutes", _
                   vbCritical, "Impulse Brake Active"
            Exit Sub
        End If
    End If

    ' Gate 4: Check bucket cooldown
    If IsBucketInCooldown(bucket) Then
        MsgBox "Cannot save: Bucket '" & bucket & "' is in cooldown" & vbCrLf & _
               "This bucket has too many recent stop-outs", vbCritical, "Hard Gate Failed"
        Exit Sub
    End If

    ' Gate 5: Check heat caps
    portHeat = PortfolioHeatAfter(addR)
    buckHeat = BucketHeatAfter(bucket, addR)

    portCap = Range("HeatCap_H_pct").Value * Range("Equity_E").Value
    buckCap = Range("BucketHeatCap_pct").Value * Range("Equity_E").Value

    If portHeat > portCap Then
        MsgBox "Cannot save: Portfolio heat would exceed cap" & vbCrLf & _
               "Heat after trade: $" & Format(portHeat, "0.00") & vbCrLf & _
               "Portfolio cap: $" & Format(portCap, "0.00"), _
               vbCritical, "Heat Cap Exceeded"
        Exit Sub
    End If

    If buckHeat > buckCap Then
        MsgBox "Cannot save: Bucket heat would exceed cap" & vbCrLf & _
               "Bucket: " & bucket & vbCrLf & _
               "Heat after trade: $" & Format(buckHeat, "0.00") & vbCrLf & _
               "Bucket cap: $" & Format(buckCap, "0.00"), _
               vbCritical, "Heat Cap Exceeded"
        Exit Sub
    End If

    ' All gates passed - save decision
    Call AppendDecisionRow
    Call UpdatePositions

    ' Reset banner and timer
    With Worksheets("TradeEntry")
        .Range("A2").Value = ""
        .Range("A2").Interior.ColorIndex = xlNone
        .Range("A3").Value = ""
    End With
    Worksheets("Control").Range("A1").Value = ""

    MsgBox "Decision saved successfully!" & vbCrLf & _
           "Ticker: " & ticker & vbCrLf & _
           "R: $" & Format(addR, "0.00"), vbInformation, "Trade Logged"
End Sub

Sub AppendDecisionRow()
    ' Append current trade to Decisions table
    Dim tbl As ListObject
    Dim newRow As ListRow
    Dim ws As Worksheet

    Set ws = Worksheets("TradeEntry")
    Set tbl = Worksheets("Decisions").ListObjects("tblDecisions")

    Set newRow = tbl.ListRows.Add

    With newRow.Range
        .Columns(1).Value = Now  ' DateTime
        .Columns(2).Value = ws.Range("B6").Value  ' Ticker
        .Columns(3).Value = ws.Range("B5").Value  ' Preset
        .Columns(4).Value = ws.Range("B8").Value  ' Bucket
        .Columns(5).Value = ws.Range("B10").Value  ' N_ATR
        .Columns(6).Value = ws.Range("B11").Value  ' K
        .Columns(7).Value = ws.Range("B9").Value  ' Entry
        .Columns(8).Value = Range("RiskPct_r").Value  ' RiskPct_r
        .Columns(9).Value = ws.Range("F5").Value  ' R_dollars
        .Columns(10).Value = ws.Range("F8").Value  ' Size_Shares
        .Columns(11).Value = ws.Range("F9").Value  ' Size_Contracts
        .Columns(12).Value = GetMethodName(ws.Range("C13").Value)  ' Method
        .Columns(13).Value = ws.Range("B16").Value  ' Delta
        .Columns(14).Value = ws.Range("B17").Value  ' DTE
        .Columns(15).Value = ws.Range("F7").Value  ' InitialStop
        .Columns(16).Value = ws.Range("A2").Value  ' Banner
        .Columns(17).Value = PortfolioHeatAfter(0)  ' HeatAtEntry (before)
        .Columns(18).Value = BucketHeatAfter(ws.Range("B8").Value, 0)  ' BucketHeatPost (before)
        .Columns(19).Value = PortfolioHeatAfter(ws.Range("F5").Value)  ' PortHeatPost
        .Columns(20).Value = ""  ' Outcome (filled later)
        .Columns(21).Value = ""  ' Notes
    End With
End Sub

Sub UpdatePositions()
    ' Open new position or add to existing in Positions table
    Dim tbl As ListObject
    Dim ticker As String, bucket As String
    Dim addR As Double, units As Long
    Dim found As Boolean
    Dim row As ListRow
    Dim newRow As ListRow
    Dim ws As Worksheet

    Set ws = Worksheets("TradeEntry")
    Set tbl = Worksheets("Positions").ListObjects("tblPositions")

    ticker = ws.Range("B6").Value
    bucket = ws.Range("B8").Value
    addR = ws.Range("F5").Value

    ' Determine units (shares or contracts)
    If ws.Range("C13").Value = 1 Then
        units = ws.Range("F8").Value  ' Shares
    Else
        units = ws.Range("F9").Value  ' Contracts
    End If

    ' Check if position already exists
    found = False
    For Each row In tbl.ListRows
        If row.Range.Columns(1).Value = ticker And _
           row.Range.Columns(8).Value <> "Closed" Then
            ' Add to existing position
            row.Range.Columns(4).Value = row.Range.Columns(4).Value + units  ' UnitsOpen
            row.Range.Columns(6).Value = row.Range.Columns(6).Value + addR  ' TotalOpenR
            row.Range.Columns(9).Value = ws.Range("B9").Value  ' LastAddPrice
            row.Range.Columns(10).Value = ws.Range("F10").Value  ' NextAddPrice
            found = True
            Exit For
        End If
    Next row

    ' Create new position if not found
    If Not found Then
        Set newRow = tbl.ListRows.Add
        With newRow.Range
            .Columns(1).Value = ticker
            .Columns(2).Value = bucket
            .Columns(3).Value = Date  ' OpenDate
            .Columns(4).Value = units  ' UnitsOpen
            .Columns(5).Value = addR / units  ' RperUnit
            .Columns(6).Value = addR  ' TotalOpenR
            .Columns(7).Value = GetMethodName(ws.Range("C13").Value)  ' Method
            .Columns(8).Value = "Open"  ' Status
            .Columns(9).Value = ws.Range("B9").Value  ' LastAddPrice
            .Columns(10).Value = ws.Range("F10").Value  ' NextAddPrice
        End With
    End If
End Sub

Function GetMethodName(methodCode As Integer) As String
    Select Case methodCode
        Case 1: GetMethodName = "Stock"
        Case 2: GetMethodName = "Opt-DeltaATR"
        Case 3: GetMethodName = "Opt-MaxLoss"
        Case Else: GetMethodName = "Unknown"
    End Select
End Function

Sub StartImpulseTimer()
    ' Manual timer starter (for testing)
    Worksheets("Control").Range("A1").Value = Now
    MsgBox "Impulse timer started at " & Format(Now, "hh:mm:ss"), vbInformation
End Sub
