Attribute VB_Name = "TF_Presets"
' ============================================================================
' Module: TF_Presets
' Purpose: FINVIZ preset management and candidate import
' ============================================================================

Option Explicit

' ----------------------------------------------------------------------------
' Sub: OpenPreset
' Opens the selected FINVIZ preset in default browser
' ----------------------------------------------------------------------------
Sub OpenPreset()
    Dim presetName As String
    Dim queryString As String
    Dim tbl As ListObject
    Dim row As ListRow
    Dim url As String
    Dim found As Boolean

    On Error GoTo ErrorHandler

    ' Get selected preset from TradeEntry
    presetName = Sheets("TradeEntry").Range("B5").Value

    If presetName = "" Then
        MsgBox "Please select a Preset first.", vbExclamation
        Exit Sub
    End If

    ' Look up query string in tblPresets
    Set tbl = Sheets("Presets").ListObjects("tblPresets")
    found = False

    For Each row In tbl.ListRows
        If row.Range.Cells(1, 1).Value = presetName Then
            queryString = row.Range.Cells(1, 2).Value
            found = True
            Exit For
        End If
    Next row

    If Not found Then
        MsgBox "Preset '" & presetName & "' not found in Presets table.", vbExclamation
        Exit Sub
    End If

    ' Build full URL
    url = "https://finviz.com/screener.ashx?" & queryString

    ' Open in browser
    CreateObject("Shell.Application").ShellExecute url

    MsgBox "Opening FINVIZ preset: " & presetName & vbCrLf & vbCrLf & _
           "Copy tickers from the page, then click Import Candidates.", _
           vbInformation

    Exit Sub

ErrorHandler:
    MsgBox "Error opening preset: " & Err.Description, vbCritical
End Sub

' ----------------------------------------------------------------------------
' Sub: ImportCandidatesPrompt
' Prompts user to paste tickers, normalizes them, adds to tblCandidates
' ----------------------------------------------------------------------------
Sub ImportCandidatesPrompt()
    Dim rawInput As String
    Dim tickers() As String
    Dim ticker As String
    Dim i As Integer
    Dim tickerList As Collection
    Dim tbl As ListObject
    Dim newRow As ListRow
    Dim presetName As String, sector As String, bucket As String
    Dim importCount As Integer

    On Error GoTo ErrorHandler

    ' Get preset context
    presetName = NzS(Sheets("TradeEntry").Range("B5").Value, "Manual")
    sector = NzS(Sheets("TradeEntry").Range("B7").Value, "")
    bucket = NzS(Sheets("TradeEntry").Range("B8").Value, "")

    ' Prompt for ticker input
    rawInput = InputBox( _
        "Paste tickers from FINVIZ (comma or line-separated):" & vbCrLf & vbCrLf & _
        "Example: AAPL, MSFT, GOOGL" & vbCrLf & _
        "or" & vbCrLf & _
        "AAPL" & vbCrLf & _
        "MSFT" & vbCrLf & _
        "GOOGL", _
        "Import Candidates - " & presetName)

    If rawInput = "" Then
        Exit Sub  ' User cancelled
    End If

    ' Normalize input: replace newlines with commas
    rawInput = Replace(rawInput, vbCrLf, ",")
    rawInput = Replace(rawInput, vbLf, ",")
    rawInput = Replace(rawInput, vbTab, ",")

    ' Split by comma
    tickers = Split(rawInput, ",")

    ' Dedupe using Collection
    Set tickerList = New Collection
    On Error Resume Next  ' Ignore duplicate key errors in Collection.Add
    For i = LBound(tickers) To UBound(tickers)
        ticker = NormalizeTicker(tickers(i))
        If ticker <> "" And Len(ticker) <= 5 Then
            tickerList.Add ticker, ticker  ' Key = ticker ensures uniqueness
        End If
    Next i
    On Error GoTo ErrorHandler

    If tickerList.Count = 0 Then
        MsgBox "No valid tickers found in input.", vbExclamation
        Exit Sub
    End If

    ' Add to tblCandidates
    Set tbl = Sheets("Candidates").ListObjects("tblCandidates")
    Application.ScreenUpdating = False
    importCount = 0

    For i = 1 To tickerList.Count
        ticker = tickerList(i)

        ' Check if already exists for today
        If Not IsCandidateExists(ticker, Date) Then
            Set newRow = tbl.ListRows.Add
            newRow.Range.Cells(1, 1).Value = Date  ' Date
            newRow.Range.Cells(1, 2).Value = ticker  ' Ticker
            newRow.Range.Cells(1, 3).Value = presetName  ' Preset
            newRow.Range.Cells(1, 4).Value = sector  ' Sector (optional)
            newRow.Range.Cells(1, 5).Value = bucket  ' Bucket (optional)
            importCount = importCount + 1
        End If
    Next i

    Application.ScreenUpdating = True

    MsgBox "Import complete:" & vbCrLf & vbCrLf & _
           "Total tickers: " & tickerList.Count & vbCrLf & _
           "New candidates added: " & importCount & vbCrLf & _
           "Duplicates skipped: " & (tickerList.Count - importCount), _
           vbInformation

    ' Refresh dropdown validation
    Call BindControls

    Exit Sub

ErrorHandler:
    Application.ScreenUpdating = True
    MsgBox "Error importing candidates: " & Err.Description, vbCritical
End Sub

' ----------------------------------------------------------------------------
' Function: IsCandidateExists
' Returns: True if ticker already exists in Candidates for given date
' ----------------------------------------------------------------------------
Function IsCandidateExists(ticker As String, tradeDate As Date) As Boolean
    Dim tbl As ListObject
    Dim row As ListRow

    On Error Resume Next
    Set tbl = Sheets("Candidates").ListObjects("tblCandidates")
    On Error GoTo 0

    If tbl Is Nothing Then
        IsCandidateExists = False
        Exit Function
    End If

    If tbl.ListRows.Count = 0 Then
        IsCandidateExists = False
        Exit Function
    End If

    For Each row In tbl.ListRows
        If UCase(Trim(row.Range.Cells(1, 2).Value)) = UCase(Trim(ticker)) And _
           Int(row.Range.Cells(1, 1).Value) = Int(tradeDate) Then
            IsCandidateExists = True
            Exit Function
        End If
    Next row

    IsCandidateExists = False
End Function

' ----------------------------------------------------------------------------
' Sub: ClearOldCandidates
' Removes candidates older than N days (default 7)
' ----------------------------------------------------------------------------
Sub ClearOldCandidates(Optional daysToKeep As Integer = 7)
    Dim tbl As ListObject
    Dim row As ListRow
    Dim cutoffDate As Date
    Dim deletedCount As Integer

    On Error GoTo ErrorHandler

    Set tbl = Sheets("Candidates").ListObjects("tblCandidates")
    cutoffDate = Date - daysToKeep
    deletedCount = 0

    Application.ScreenUpdating = False

    ' Loop backwards to safely delete rows
    Dim i As Long
    For i = tbl.ListRows.Count To 1 Step -1
        If tbl.ListRows(i).Range.Cells(1, 1).Value < cutoffDate Then
            tbl.ListRows(i).Delete
            deletedCount = deletedCount + 1
        End If
    Next i

    Application.ScreenUpdating = True

    MsgBox "Cleared " & deletedCount & " candidates older than " & daysToKeep & " days.", vbInformation

    Exit Sub

ErrorHandler:
    Application.ScreenUpdating = True
    MsgBox "Error clearing old candidates: " & Err.Description, vbCritical
End Sub
