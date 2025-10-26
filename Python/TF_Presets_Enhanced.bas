Attribute VB_Name = "TF_Presets_Enhanced"
' ============================================================================
' Module: TF_Presets_Enhanced
' Purpose: Enhanced version with Python integration for auto-scraping
' Usage: Replace TF_Presets module with this version after Python setup
' ============================================================================

Option Explicit

' ----------------------------------------------------------------------------
' Sub: ImportCandidatesPrompt
' Enhanced version: Tries Python scraper first, falls back to manual paste
' ----------------------------------------------------------------------------
Sub ImportCandidatesPrompt()
    Dim presetName As String, sector As String, bucket As String
    Dim queryString As String
    Dim tickers As Variant
    Dim tickerList As Collection
    Dim rawInput As String
    Dim tbl As ListObject
    Dim newRow As ListRow
    Dim i As Integer
    Dim importCount As Integer
    Dim usePython As Boolean

    On Error GoTo ErrorHandler

    ' Get preset context
    presetName = NzS(Sheets("TradeEntry").Range("B5").Value, "Manual")
    sector = NzS(Sheets("TradeEntry").Range("B7").Value, "")
    bucket = NzS(Sheets("TradeEntry").Range("B8").Value, "")

    ' Look up query string if preset selected
    If presetName <> "Manual" Then
        queryString = GetQueryStringForPreset(presetName)
        If queryString = "" Then
            MsgBox "Could not find query string for preset: " & presetName, vbExclamation
            Exit Sub
        End If
    End If

    ' === TRY PYTHON SCRAPER FIRST ===
    usePython = False
    If queryString <> "" And IsPythonAvailable() Then
        Dim response As VbMsgBoxResult
        response = MsgBox( _
            "Python integration detected!" & vbCrLf & vbCrLf & _
            "Auto-scrape FINVIZ? (5-10 seconds)" & vbCrLf & _
            "or manually paste tickers?", _
            vbYesNoCancel + vbQuestion, "Import Method")

        If response = vbCancel Then Exit Sub

        If response = vbYes Then
            usePython = True

            ' Show progress message
            Application.StatusBar = "Scraping FINVIZ for " & presetName & "..."
            Application.ScreenUpdating = False

            tickers = CallPythonFinvizScraper(queryString)

            Application.StatusBar = False
            Application.ScreenUpdating = True

            ' Check if scraping succeeded
            If Not IsArray(tickers) Or UBound(tickers) < 0 Then
                MsgBox "Python scraping returned no tickers." & vbCrLf & _
                       "Falling back to manual paste.", vbExclamation
                usePython = False
            End If
        End If
    End If

    ' === FALLBACK: MANUAL PASTE ===
    If Not usePython Then
        rawInput = InputBox( _
            "Paste tickers from FINVIZ (comma or line-separated):" & vbCrLf & vbCrLf & _
            "Example: AAPL, MSFT, GOOGL" & vbCrLf & _
            "or" & vbCrLf & _
            "AAPL" & vbCrLf & _
            "MSFT" & vbCrLf & _
            "GOOGL", _
            "Import Candidates - " & presetName)

        If rawInput = "" Then Exit Sub

        ' Parse manual input
        rawInput = Replace(rawInput, vbCrLf, ",")
        rawInput = Replace(rawInput, vbLf, ",")
        rawInput = Replace(rawInput, vbTab, ",")

        tickers = Split(rawInput, ",")
    End If

    ' === NORMALIZE AND DEDUPE ===
    Set tickerList = New Collection
    On Error Resume Next

    For i = LBound(tickers) To UBound(tickers)
        Dim ticker As String
        ticker = NormalizeTicker(CStr(tickers(i)))
        If ticker <> "" And Len(ticker) <= 5 Then
            tickerList.Add ticker, ticker  ' Key ensures uniqueness
        End If
    Next i

    On Error GoTo ErrorHandler

    If tickerList.Count = 0 Then
        MsgBox "No valid tickers found in input.", vbExclamation
        Exit Sub
    End If

    ' === ADD TO CANDIDATES TABLE ===
    Set tbl = Sheets("Candidates").ListObjects("tblCandidates")
    Application.ScreenUpdating = False
    importCount = 0

    For i = 1 To tickerList.Count
        ticker = tickerList(i)

        If Not IsCandidateExists(ticker, Date) Then
            Set newRow = tbl.ListRows.Add
            newRow.Range.Cells(1, 1).Value = Date
            newRow.Range.Cells(1, 2).Value = ticker
            newRow.Range.Cells(1, 3).Value = presetName
            newRow.Range.Cells(1, 4).Value = sector
            newRow.Range.Cells(1, 5).Value = bucket
            importCount = importCount + 1
        End If
    Next i

    Application.ScreenUpdating = True

    ' === REPORT RESULTS ===
    Dim msg As String
    msg = "Import complete:" & vbCrLf & vbCrLf & _
          "Method: " & IIf(usePython, "🐍 Python Auto-Scrape", "📋 Manual Paste") & vbCrLf & _
          "Total tickers: " & tickerList.Count & vbCrLf & _
          "New candidates added: " & importCount & vbCrLf & _
          "Duplicates skipped: " & (tickerList.Count - importCount)

    If usePython Then
        msg = msg & vbCrLf & vbCrLf & "⏱ Auto-scraping saved ~30 seconds!"
    End If

    MsgBox msg, vbInformation

    ' Refresh dropdown validation
    Call BindControls

    Exit Sub

ErrorHandler:
    Application.ScreenUpdating = True
    Application.StatusBar = False
    MsgBox "Error importing candidates: " & Err.Description, vbCritical
End Sub

' ----------------------------------------------------------------------------
' Function: GetQueryStringForPreset
' Looks up query string in tblPresets
' ----------------------------------------------------------------------------
Function GetQueryStringForPreset(presetName As String) As String
    Dim tbl As ListObject
    Dim row As ListRow

    On Error Resume Next
    Set tbl = Sheets("Presets").ListObjects("tblPresets")
    On Error GoTo 0

    If tbl Is Nothing Then
        GetQueryStringForPreset = ""
        Exit Function
    End If

    For Each row In tbl.ListRows
        If row.Range.Cells(1, 1).Value = presetName Then
            GetQueryStringForPreset = row.Range.Cells(1, 2).Value
            Exit Function
        End If
    Next row

    GetQueryStringForPreset = ""
End Function

' ----------------------------------------------------------------------------
' Sub: OpenPreset (same as original)
' ----------------------------------------------------------------------------
Sub OpenPreset()
    Dim presetName As String
    Dim queryString As String
    Dim url As String

    On Error GoTo ErrorHandler

    presetName = Sheets("TradeEntry").Range("B5").Value

    If presetName = "" Then
        MsgBox "Please select a Preset first.", vbExclamation
        Exit Sub
    End If

    queryString = GetQueryStringForPreset(presetName)

    If queryString = "" Then
        MsgBox "Preset '" & presetName & "' not found in Presets table.", vbExclamation
        Exit Sub
    End If

    url = "https://finviz.com/screener.ashx?" & queryString

    CreateObject("Shell.Application").ShellExecute url

    If IsPythonAvailable() Then
        MsgBox "Opening FINVIZ preset: " & presetName & vbCrLf & vbCrLf & _
               "💡 Tip: Use Import Candidates with Python auto-scrape instead of manual copy!", _
               vbInformation
    Else
        MsgBox "Opening FINVIZ preset: " & presetName & vbCrLf & vbCrLf & _
               "Copy tickers from the page, then click Import Candidates.", _
               vbInformation
    End If

    Exit Sub

ErrorHandler:
    MsgBox "Error opening preset: " & Err.Description, vbCritical
End Sub

' ----------------------------------------------------------------------------
' Other functions (same as original TF_Presets)
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

    For Each row In tbl.ListRows
        If UCase(Trim(row.Range.Cells(1, 2).Value)) = UCase(Trim(ticker)) And _
           Int(row.Range.Cells(1, 1).Value) = Int(tradeDate) Then
            IsCandidateExists = True
            Exit Function
        End If
    Next row

    IsCandidateExists = False
End Function

Sub ClearOldCandidates(Optional daysToKeep As Integer = 7)
    Dim tbl As ListObject
    Dim cutoffDate As Date
    Dim deletedCount As Integer
    Dim i As Long

    On Error GoTo ErrorHandler

    Set tbl = Sheets("Candidates").ListObjects("tblCandidates")
    cutoffDate = Date - daysToKeep
    deletedCount = 0

    Application.ScreenUpdating = False

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
