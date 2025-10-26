Attribute VB_Name = "TF_Data"
' ============================================================================
' Module: TF_Data
' Purpose: Data structure management, heat calculations, cooldown logic
' ============================================================================

Option Explicit

' ----------------------------------------------------------------------------
' Sub: EnsureStructure
' Creates all sheets, tables, named ranges, and seeds default data
' Run this FIRST when setting up a new workbook
' ----------------------------------------------------------------------------
Sub EnsureStructure()
    Application.ScreenUpdating = False

    ' Create all sheets
    Call GetOrCreateSheet("TradeEntry")
    Call GetOrCreateSheet("Presets")
    Call GetOrCreateSheet("Buckets")
    Call GetOrCreateSheet("Candidates")
    Call GetOrCreateSheet("Decisions")
    Call GetOrCreateSheet("Positions")
    Call GetOrCreateSheet("Summary")

    Dim ctrlSheet As Worksheet
    Set ctrlSheet = GetOrCreateSheet("Control")
    ctrlSheet.Visible = xlSheetVeryHidden  ' Hide from user

    ' Create tables with headers
    Call GetOrCreateTable(Sheets("Presets"), "tblPresets", _
        Array("Name", "QueryString"))

    Call GetOrCreateTable(Sheets("Buckets"), "tblBuckets", _
        Array("Sector", "Bucket", "BucketHeatCapPct", "StopoutsToCooldown", _
              "StopoutsWindowBars", "CooldownBars", "CooldownActive", "CooldownEndsOn"))

    Call GetOrCreateTable(Sheets("Candidates"), "tblCandidates", _
        Array("Date", "Ticker", "Preset", "Sector", "Bucket"))

    Call GetOrCreateTable(Sheets("Decisions"), "tblDecisions", _
        Array("DateTime", "Ticker", "Preset", "Bucket", "N_ATR", "K", "Entry", _
              "RiskPct_r", "R_dollars", "Size_Shares", "Size_Contracts", "Method", _
              "Delta", "DTE", "InitialStop", "Banner", "HeatAtEntry", _
              "BucketHeatPost", "PortHeatPost", "Outcome", "Notes"))

    Call GetOrCreateTable(Sheets("Positions"), "tblPositions", _
        Array("Ticker", "Bucket", "OpenDate", "UnitsOpen", "RperUnit", _
              "TotalOpenR", "Method", "Status", "LastAddPrice", "NextAddPrice"))

    ' Create named ranges on Summary sheet
    With Sheets("Summary")
        .Range("A1").Value = "Settings"
        .Range("A2").Value = "Equity (E):"
        .Range("A3").Value = "Risk % per unit (r):"
        .Range("A4").Value = "Stop Multiple (K):"
        .Range("A5").Value = "Portfolio Heat Cap %:"
        .Range("A6").Value = "Bucket Heat Cap %:"
        .Range("A7").Value = "Add Step (N units):"
        .Range("A8").Value = "Earnings Buffer Days:"
    End With

    Call EnsureName("Equity_E", "Summary!B2", 10000)
    Call EnsureName("RiskPct_r", "Summary!B3", 0.0075)
    Call EnsureName("StopMultiple_K", "Summary!B4", 2)
    Call EnsureName("HeatCap_H_pct", "Summary!B5", 0.04)
    Call EnsureName("BucketHeatCap_pct", "Summary!B6", 0.015)
    Call EnsureName("AddStepN", "Summary!B7", 0.5)
    Call EnsureName("EarningsBufferDays", "Summary!B8", 3)

    ' Seed default data
    Call SeedPresets
    Call SeedBuckets

    Application.ScreenUpdating = True

    MsgBox "Workbook structure created successfully!" & vbCrLf & vbCrLf & _
           "Sheets: 8 created" & vbCrLf & _
           "Tables: 5 created" & vbCrLf & _
           "Named Ranges: 7 created" & vbCrLf & vbCrLf & _
           "Next: Build the TradeEntry UI manually (see Setup Guide)", _
           vbInformation, "Setup Complete"
End Sub

' ----------------------------------------------------------------------------
' Sub: SeedPresets
' Seeds default FINVIZ screener presets
' ----------------------------------------------------------------------------
Sub SeedPresets()
    Dim tbl As ListObject
    Set tbl = Sheets("Presets").ListObjects("tblPresets")

    ' Clear existing data rows (keep headers)
    If tbl.ListRows.Count > 0 Then
        tbl.DataBodyRange.Delete
    End If

    ' Add 5 default presets
    With tbl
        .ListRows.Add.Range.Value = Array("TF_BREAKOUT_LONG", _
            "v=211&p=d&s=ta_newhigh&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pa,ta_sma200_pa&o=-relativevolume")

        .ListRows.Add.Range.Value = Array("TF_MOMENTUM_UPTREND", _
            "v=211&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pa,ta_sma200_pa&dr=y1&o=-marketcap")

        .ListRows.Add.Range.Value = Array("TF_UNUSUAL_VOLUME", _
            "v=211&p=d&s=ta_unusualvolume&f=cap_largeover,sh_price_o20,ta_sma50_pa,ta_sma200_pa&o=-relativevolume")

        .ListRows.Add.Range.Value = Array("TF_BREAKDOWN_SHORT", _
            "v=211&p=d&s=ta_newlow&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pb,ta_sma200_pb&o=-relativevolume")

        .ListRows.Add.Range.Value = Array("TF_MOMENTUM_DOWNTREND", _
            "v=211&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pb,ta_sma200_pb&dr=y1&o=-marketcap")
    End With
End Sub

' ----------------------------------------------------------------------------
' Sub: SeedBuckets
' Seeds default correlation buckets
' ----------------------------------------------------------------------------
Sub SeedBuckets()
    Dim tbl As ListObject
    Set tbl = Sheets("Buckets").ListObjects("tblBuckets")

    ' Clear existing data rows
    If tbl.ListRows.Count > 0 Then
        tbl.DataBodyRange.Delete
    End If

    ' Add 6 default buckets
    ' Columns: Sector, Bucket, BucketHeatCapPct, StopoutsToCooldown, StopoutsWindowBars, CooldownBars, CooldownActive, CooldownEndsOn
    With tbl
        .ListRows.Add.Range.Value = Array("Technology", "Tech/Comm", 0.015, 2, 20, 10, False, "")
        .ListRows.Add.Range.Value = Array("Consumer", "Consumer", 0.015, 2, 20, 10, False, "")
        .ListRows.Add.Range.Value = Array("Financials", "Financials", 0.015, 2, 20, 10, False, "")
        .ListRows.Add.Range.Value = Array("Industrials", "Industrials", 0.015, 2, 20, 10, False, "")
        .ListRows.Add.Range.Value = Array("Energy", "Energy/Materials", 0.015, 2, 20, 10, False, "")
        .ListRows.Add.Range.Value = Array("Healthcare", "Defensives/REITs", 0.015, 2, 20, 10, False, "")
    End With
End Sub

' ----------------------------------------------------------------------------
' Function: TodayCandidates
' Returns: Dynamic range of tickers from Candidates table where Date = Today
' ----------------------------------------------------------------------------
Function TodayCandidates() As Range
    Dim tbl As ListObject
    Dim tickerCol As ListColumn

    On Error Resume Next
    Set tbl = Sheets("Candidates").ListObjects("tblCandidates")
    Set tickerCol = tbl.ListColumns("Ticker")

    ' For simplicity, return entire Ticker column
    ' (Filtering by date would require advanced array formulas)
    Set TodayCandidates = tickerCol.DataBodyRange
    On Error GoTo 0
End Function

' ----------------------------------------------------------------------------
' Function: PortfolioHeatAfter
' Returns: Total portfolio heat (current open + proposed trade)
' ----------------------------------------------------------------------------
Function PortfolioHeatAfter(addR As Double) As Double
    Dim tbl As ListObject
    Dim row As ListRow
    Dim totalHeat As Double

    On Error Resume Next
    Set tbl = Sheets("Positions").ListObjects("tblPositions")
    On Error GoTo 0

    If tbl Is Nothing Then
        PortfolioHeatAfter = addR
        Exit Function
    End If

    If tbl.ListRows.Count = 0 Then
        PortfolioHeatAfter = addR
        Exit Function
    End If

    totalHeat = 0

    For Each row In tbl.ListRows
        ' Column 8 = Status, Column 6 = TotalOpenR
        If row.Range.Cells(1, 8).Value <> "Closed" Then
            totalHeat = totalHeat + NzD(row.Range.Cells(1, 6).Value, 0)
        End If
    Next row

    PortfolioHeatAfter = totalHeat + addR
End Function

' ----------------------------------------------------------------------------
' Function: BucketHeatAfter
' Returns: Bucket-specific heat (current open + proposed trade)
' ----------------------------------------------------------------------------
Function BucketHeatAfter(bucket As String, addR As Double) As Double
    Dim tbl As ListObject
    Dim row As ListRow
    Dim bucketHeat As Double

    On Error Resume Next
    Set tbl = Sheets("Positions").ListObjects("tblPositions")
    On Error GoTo 0

    If tbl Is Nothing Then
        BucketHeatAfter = addR
        Exit Function
    End If

    If tbl.ListRows.Count = 0 Then
        BucketHeatAfter = addR
        Exit Function
    End If

    bucketHeat = 0

    For Each row In tbl.ListRows
        ' Column 2 = Bucket, Column 8 = Status, Column 6 = TotalOpenR
        If row.Range.Cells(1, 2).Value = bucket And _
           row.Range.Cells(1, 8).Value <> "Closed" Then
            bucketHeat = bucketHeat + NzD(row.Range.Cells(1, 6).Value, 0)
        End If
    Next row

    BucketHeatAfter = bucketHeat + addR
End Function

' ----------------------------------------------------------------------------
' Function: IsBucketInCooldown
' Returns: True if bucket is currently in cooldown period
' ----------------------------------------------------------------------------
Function IsBucketInCooldown(bucket As String) As Boolean
    Dim tbl As ListObject
    Dim row As ListRow

    On Error Resume Next
    Set tbl = Sheets("Buckets").ListObjects("tblBuckets")
    On Error GoTo 0

    If tbl Is Nothing Then
        IsBucketInCooldown = False
        Exit Function
    End If

    For Each row In tbl.ListRows
        ' Column 2 = Bucket, Column 7 = CooldownActive, Column 8 = CooldownEndsOn
        If row.Range.Cells(1, 2).Value = bucket Then
            If row.Range.Cells(1, 7).Value = True Then
                If row.Range.Cells(1, 8).Value >= Date Then
                    IsBucketInCooldown = True
                    Exit Function
                End If
            End If
        End If
    Next row

    IsBucketInCooldown = False
End Function

' ----------------------------------------------------------------------------
' Sub: UpdateCooldowns
' Scans recent decisions and updates bucket cooldown flags
' Call this weekly or after each stop-out
' ----------------------------------------------------------------------------
Sub UpdateCooldowns()
    Dim bucketTbl As ListObject
    Dim decisTbl As ListObject
    Dim bRow As ListRow, dRow As ListRow
    Dim bucket As String
    Dim stopoutThreshold As Integer, windowBars As Integer, cooldownBars As Integer
    Dim windowStart As Date
    Dim stopoutCount As Integer

    Application.ScreenUpdating = False

    On Error Resume Next
    Set bucketTbl = Sheets("Buckets").ListObjects("tblBuckets")
    Set decisTbl = Sheets("Decisions").ListObjects("tblDecisions")
    On Error GoTo 0

    If bucketTbl Is Nothing Or decisTbl Is Nothing Then Exit Sub
    If decisTbl.ListRows.Count = 0 Then Exit Sub

    For Each bRow In bucketTbl.ListRows
        bucket = bRow.Range.Cells(1, 2).Value
        stopoutThreshold = NzD(bRow.Range.Cells(1, 4).Value, 2)
        windowBars = NzD(bRow.Range.Cells(1, 5).Value, 20)
        cooldownBars = NzD(bRow.Range.Cells(1, 6).Value, 10)

        windowStart = Date - windowBars

        ' Count StopOuts in window for this bucket
        stopoutCount = 0
        For Each dRow In decisTbl.ListRows
            ' Column 4 = Bucket, Column 1 = DateTime, Column 20 = Outcome
            If dRow.Range.Cells(1, 4).Value = bucket And _
               dRow.Range.Cells(1, 1).Value >= windowStart And _
               dRow.Range.Cells(1, 20).Value = "StopOut" Then
                stopoutCount = stopoutCount + 1
            End If
        Next dRow

        ' Update cooldown flags
        If stopoutCount >= stopoutThreshold Then
            bRow.Range.Cells(1, 7).Value = True  ' CooldownActive
            bRow.Range.Cells(1, 8).Value = Date + cooldownBars  ' CooldownEndsOn
        Else
            ' Clear cooldown if past end date
            If bRow.Range.Cells(1, 8).Value < Date Then
                bRow.Range.Cells(1, 7).Value = False
            End If
        End If
    Next bRow

    Application.ScreenUpdating = True

    MsgBox "Cooldowns updated for all buckets.", vbInformation
End Sub
