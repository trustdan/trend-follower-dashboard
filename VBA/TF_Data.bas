Attribute VB_Name = "TF_Data"
Option Explicit

' ========================================
' TF_Data Module
' Structure setup, heat calculations, cooldown logic
' ========================================

Sub EnsureStructure()
    ' Creates all sheets, tables, and named ranges
    Application.ScreenUpdating = False

    ' Create all sheets
    Dim ws As Worksheet
    Set ws = GetOrCreateSheet("TradeEntry")
    Set ws = GetOrCreateSheet("Presets")
    Set ws = GetOrCreateSheet("Buckets")
    Set ws = GetOrCreateSheet("Candidates")
    Set ws = GetOrCreateSheet("Decisions")
    Set ws = GetOrCreateSheet("Positions")
    Set ws = GetOrCreateSheet("Summary")
    Set ws = GetOrCreateSheet("Control")

    ' Hide Control sheet
    Worksheets("Control").Visible = xlSheetVeryHidden

    ' Create tables
    Call CreateTables

    ' Define named ranges
    Call DefineNamedRanges

    ' Seed default data
    Call SeedPresets
    Call SeedBuckets

    Application.ScreenUpdating = True

    MsgBox "Structure created successfully!", vbInformation
End Sub

Sub CreateTables()
    Dim tbl As ListObject

    ' Presets table
    Set tbl = GetOrCreateTable(Worksheets("Presets"), "tblPresets", _
        Array("Name", "QueryString"))

    ' Buckets table
    Set tbl = GetOrCreateTable(Worksheets("Buckets"), "tblBuckets", _
        Array("Sector", "Bucket", "BucketHeatCapPct", "StopoutsToCooldown", _
              "StopoutsWindowBars", "CooldownBars", "CooldownActive", "CooldownEndsOn"))

    ' Candidates table
    Set tbl = GetOrCreateTable(Worksheets("Candidates"), "tblCandidates", _
        Array("Date", "Ticker", "Preset", "Sector", "Bucket"))

    ' Decisions table
    Set tbl = GetOrCreateTable(Worksheets("Decisions"), "tblDecisions", _
        Array("DateTime", "Ticker", "Preset", "Bucket", "N_ATR", "K", "Entry", _
              "RiskPct_r", "R_dollars", "Size_Shares", "Size_Contracts", "Method", _
              "Delta", "DTE", "InitialStop", "Banner", "HeatAtEntry", _
              "BucketHeatPost", "PortHeatPost", "Outcome", "Notes"))

    ' Positions table
    Set tbl = GetOrCreateTable(Worksheets("Positions"), "tblPositions", _
        Array("Ticker", "Bucket", "OpenDate", "UnitsOpen", "RperUnit", _
              "TotalOpenR", "Method", "Status", "LastAddPrice", "NextAddPrice"))
End Sub

Sub DefineNamedRanges()
    ' Define named ranges on Summary sheet with defaults
    With Worksheets("Summary")
        .Range("A1").Value = "Setting"
        .Range("B1").Value = "Value"

        .Range("A2").Value = "Equity_E"
        .Range("A3").Value = "RiskPct_r"
        .Range("A4").Value = "StopMultiple_K"
        .Range("A5").Value = "HeatCap_H_pct"
        .Range("A6").Value = "BucketHeatCap_pct"
        .Range("A7").Value = "AddStepN"
        .Range("A8").Value = "EarningsBufferDays"
    End With

    Call EnsureName("Equity_E", "Summary!B2", 10000)
    Call EnsureName("RiskPct_r", "Summary!B3", 0.0075)
    Call EnsureName("StopMultiple_K", "Summary!B4", 2)
    Call EnsureName("HeatCap_H_pct", "Summary!B5", 0.04)
    Call EnsureName("BucketHeatCap_pct", "Summary!B6", 0.015)
    Call EnsureName("AddStepN", "Summary!B7", 0.5)
    Call EnsureName("EarningsBufferDays", "Summary!B8", 3)
End Sub

Sub SeedPresets()
    Dim tbl As ListObject
    Dim ws As Worksheet

    Set ws = Worksheets("Presets")
    Set tbl = ws.ListObjects("tblPresets")

    ' Clear existing data (keep headers)
    If tbl.ListRows.Count > 0 Then
        On Error Resume Next
        tbl.DataBodyRange.Delete
        On Error GoTo 0
    End If

    ' Add 5 default FINVIZ presets
    With tbl
        .ListRows.Add.Range.Value = Array("TF_BREAKOUT_LONG", "v=211&f=ta_highlow52w_nh&ft=4")
        .ListRows.Add.Range.Value = Array("TF_MOMENTUM_UPTREND", "v=211&f=ta_sma200_pa,ta_sma50_pa&ft=4")
        .ListRows.Add.Range.Value = Array("TF_UNUSUAL_VOLUME", "v=211&f=sh_relvol_o2&ft=4")
        .ListRows.Add.Range.Value = Array("TF_GAP_UP", "v=211&f=ta_gap_u5&ft=4")
        .ListRows.Add.Range.Value = Array("TF_STRONG_TREND", "v=211&f=ta_changeopen_u5&ft=4")
    End With
End Sub

Sub SeedBuckets()
    Dim tbl As ListObject
    Dim ws As Worksheet

    Set ws = Worksheets("Buckets")
    Set tbl = ws.ListObjects("tblBuckets")

    ' Clear existing data
    If tbl.ListRows.Count > 0 Then
        On Error Resume Next
        tbl.DataBodyRange.Delete
        On Error GoTo 0
    End If

    ' Add 6 correlation buckets
    With tbl
        .ListRows.Add.Range.Value = Array("Technology", "Tech/Comm", 0.015, 2, 20, 10, False, Empty)
        .ListRows.Add.Range.Value = Array("Healthcare", "Healthcare", 0.015, 2, 20, 10, False, Empty)
        .ListRows.Add.Range.Value = Array("Financials", "Financials", 0.015, 2, 20, 10, False, Empty)
        .ListRows.Add.Range.Value = Array("Consumer", "Consumer", 0.015, 2, 20, 10, False, Empty)
        .ListRows.Add.Range.Value = Array("Industrials", "Industrials", 0.015, 2, 20, 10, False, Empty)
        .ListRows.Add.Range.Value = Array("Energy", "Energy/Materials", 0.015, 2, 20, 10, False, Empty)
    End With
End Sub

Function PortfolioHeatAfter(addR As Double) As Double
    ' Calculates total portfolio heat including proposed trade
    Dim tbl As ListObject
    Dim row As ListRow
    Dim total As Double

    On Error Resume Next
    Set tbl = Worksheets("Positions").ListObjects("tblPositions")
    On Error GoTo 0

    If tbl Is Nothing Then
        PortfolioHeatAfter = addR
        Exit Function
    End If

    total = 0
    For Each row In tbl.ListRows
        If row.Range.Columns(8).Value <> "Closed" Then  ' Status column
            total = total + NzD(row.Range.Columns(6).Value, 0)  ' TotalOpenR column
        End If
    Next row

    PortfolioHeatAfter = total + addR
End Function

Function BucketHeatAfter(bucket As String, addR As Double) As Double
    ' Calculates bucket-specific heat including proposed trade
    Dim tbl As ListObject
    Dim row As ListRow
    Dim total As Double

    On Error Resume Next
    Set tbl = Worksheets("Positions").ListObjects("tblPositions")
    On Error GoTo 0

    If tbl Is Nothing Then
        BucketHeatAfter = addR
        Exit Function
    End If

    total = 0
    For Each row In tbl.ListRows
        If row.Range.Columns(2).Value = bucket And _
           row.Range.Columns(8).Value <> "Closed" Then
            total = total + NzD(row.Range.Columns(6).Value, 0)
        End If
    Next row

    BucketHeatAfter = total + addR
End Function

Function IsBucketInCooldown(bucket As String) As Boolean
    ' Check if bucket is currently in cooldown
    Dim tbl As ListObject
    Dim row As ListRow

    On Error Resume Next
    Set tbl = Worksheets("Buckets").ListObjects("tblBuckets")
    On Error GoTo 0

    If tbl Is Nothing Then
        IsBucketInCooldown = False
        Exit Function
    End If

    For Each row In tbl.ListRows
        If row.Range.Columns(2).Value = bucket Then  ' Bucket column
            If row.Range.Columns(7).Value = True Then  ' CooldownActive
                If NzD(row.Range.Columns(8).Value, 0) >= Date Then  ' CooldownEndsOn
                    IsBucketInCooldown = True
                    Exit Function
                End If
            End If
        End If
    Next row

    IsBucketInCooldown = False
End Function

Sub UpdateCooldowns()
    ' Scans decisions and updates bucket cooldowns
    Dim bucketTbl As ListObject
    Dim decisTbl As ListObject
    Dim bRow As ListRow
    Dim dRow As ListRow
    Dim bucket As String
    Dim stopoutCount As Integer
    Dim windowStart As Date
    Dim stopoutThreshold As Integer
    Dim windowBars As Integer
    Dim cooldownBars As Integer

    Set bucketTbl = Worksheets("Buckets").ListObjects("tblBuckets")
    Set decisTbl = Worksheets("Decisions").ListObjects("tblDecisions")

    For Each bRow In bucketTbl.ListRows
        bucket = bRow.Range.Columns(2).Value
        stopoutThreshold = NzD(bRow.Range.Columns(4).Value, 2)
        windowBars = NzD(bRow.Range.Columns(5).Value, 20)
        cooldownBars = NzD(bRow.Range.Columns(6).Value, 10)

        windowStart = Date - windowBars

        ' Count StopOuts in window
        stopoutCount = 0
        For Each dRow In decisTbl.ListRows
            If dRow.Range.Columns(4).Value = bucket And _
               dRow.Range.Columns(1).Value >= windowStart And _
               dRow.Range.Columns(20).Value = "StopOut" Then
                stopoutCount = stopoutCount + 1
            End If
        Next dRow

        ' Update cooldown flags
        If stopoutCount >= stopoutThreshold Then
            bRow.Range.Columns(7).Value = True  ' CooldownActive
            bRow.Range.Columns(8).Value = Date + cooldownBars  ' CooldownEndsOn
        Else
            ' Clear cooldown if past end date
            If NzD(bRow.Range.Columns(8).Value, 0) < Date Then
                bRow.Range.Columns(7).Value = False
            End If
        End If
    Next bRow
End Sub

Function IsTickerInCandidates(ticker As String, tradeDate As Date) As Boolean
    ' Check if ticker exists in today's candidates
    Dim tbl As ListObject
    Dim row As ListRow

    On Error Resume Next
    Set tbl = Worksheets("Candidates").ListObjects("tblCandidates")
    On Error GoTo 0

    If tbl Is Nothing Then
        IsTickerInCandidates = False
        Exit Function
    End If

    For Each row In tbl.ListRows
        If row.Range.Columns(2).Value = ticker And _
           row.Range.Columns(1).Value = tradeDate Then
            IsTickerInCandidates = True
            Exit Function
        End If
    Next row

    IsTickerInCandidates = False
End Function

Sub ClearOldCandidates(daysOld As Integer)
    ' Remove candidates older than specified days
    Dim tbl As ListObject
    Dim row As ListRow
    Dim cutoffDate As Date
    Dim i As Long

    Set tbl = Worksheets("Candidates").ListObjects("tblCandidates")
    cutoffDate = Date - daysOld

    ' Loop backwards to safely delete rows
    For i = tbl.ListRows.Count To 1 Step -1
        If tbl.ListRows(i).Range.Columns(1).Value < cutoffDate Then
            tbl.ListRows(i).Delete
        End If
    Next i
End Sub
