Attribute VB_Name = "TF_Presets"
Option Explicit

' ========================================
' TF_Presets Module
' FINVIZ integration and candidate management
' ========================================

Sub OpenPreset()
    ' Opens the FINVIZ URL for selected preset in browser
    Dim presetName As String
    Dim queryString As String
    Dim url As String
    Dim tbl As ListObject
    Dim row As ListRow

    presetName = Worksheets("TradeEntry").Range("B5").Value

    If presetName = "" Then
        MsgBox "Please select a Preset first", vbExclamation
        Exit Sub
    End If

    ' Look up query string from Presets table
    Set tbl = Worksheets("Presets").ListObjects("tblPresets")
    For Each row In tbl.ListRows
        If row.Range.Columns(1).Value = presetName Then
            queryString = row.Range.Columns(2).Value
            Exit For
        End If
    Next row

    If queryString = "" Then
        MsgBox "Preset not found: " & presetName, vbExclamation
        Exit Sub
    End If

    ' Build and open URL
    url = "https://finviz.com/screener.ashx?" & queryString

    On Error Resume Next
    CreateObject("WScript.Shell").Run url
    On Error GoTo 0
End Sub

Sub ImportCandidatesPrompt()
    ' Smart import: Auto-scrape with Python if available, else manual paste
    Dim presetName As String
    Dim queryString As String
    Dim tbl As ListObject
    Dim row As ListRow
    Dim pythonAvailable As Boolean

    presetName = Worksheets("TradeEntry").Range("B5").Value

    If presetName = "" Then
        MsgBox "Please select a Preset first", vbExclamation
        Exit Sub
    End If

    ' Look up query string
    Set tbl = Worksheets("Presets").ListObjects("tblPresets")
    For Each row In tbl.ListRows
        If row.Range.Columns(1).Value = presetName Then
            queryString = row.Range.Columns(2).Value
            Exit For
        End If
    Next row

    If queryString = "" Then
        MsgBox "Preset not found: " & presetName, vbExclamation
        Exit Sub
    End If

    ' Check if Python is available
    pythonAvailable = TF_Python_Bridge.IsPythonAvailable()

    If pythonAvailable Then
        ' Auto-scrape with Python
        Call ImportWithPython(presetName, queryString)
    Else
        ' Fallback to manual import
        Call ImportManual(presetName)
    End If
End Sub

Sub ImportWithPython(presetName As String, queryString As String)
    ' Auto-scrape using Python
    Dim tickerArray As Variant
    Dim tickers As Collection
    Dim ticker As Variant
    Dim sector As String
    Dim bucket As String
    Dim addedCount As Integer

    sector = Worksheets("TradeEntry").Range("B7").Value
    bucket = Worksheets("TradeEntry").Range("B8").Value

    ' Show progress
    Application.StatusBar = "Scraping FINVIZ for " & presetName & "..."

    ' Call Python scraper
    tickerArray = TF_Python_Bridge.CallPythonFinvizScraper(queryString)

    Application.StatusBar = False

    If Not IsArray(tickerArray) Or UBound(tickerArray) < 0 Then
        MsgBox "Python scraping failed or returned no tickers" & vbCrLf & _
               "Falling back to manual import...", vbInformation
        Call ImportManual(presetName)
        Exit Sub
    End If

    ' Convert array to collection
    Set tickers = New Collection
    For Each ticker In tickerArray
        On Error Resume Next
        tickers.Add CStr(ticker), CStr(ticker)
        On Error GoTo 0
    Next ticker

    ' Add to Candidates table
    addedCount = 0
    For Each ticker In tickers
        If AddCandidate(CStr(ticker), presetName, sector, bucket) Then
            addedCount = addedCount + 1
        End If
    Next ticker

    MsgBox "Auto-scraped " & addedCount & " candidates from FINVIZ" & vbCrLf & _
           "Date: " & Format(Date, "yyyy-mm-dd") & vbCrLf & _
           "Preset: " & presetName, vbInformation, "Python Import Success"
End Sub

Sub ImportManual(presetName As String)
    ' Manual paste import (fallback when Python unavailable)
    Dim inputText As String
    Dim tickers As Collection
    Dim ticker As String
    Dim sector As String
    Dim bucket As String
    Dim addedCount As Integer

    sector = Worksheets("TradeEntry").Range("B7").Value
    bucket = Worksheets("TradeEntry").Range("B8").Value

    ' Prompt for tickers
    inputText = InputBox( _
        "Python auto-scraping not available" & vbCrLf & _
        "Paste ticker symbols manually (comma or line-separated):" & vbCrLf & vbCrLf & _
        "Examples:" & vbCrLf & _
        "  AAPL, MSFT, NVDA" & vbCrLf & _
        "  AAPL" & vbCrLf & _
        "  MSFT" & vbCrLf & _
        "  NVDA", _
        "Import Candidates - Manual", _
        "")

    If inputText = "" Then
        Exit Sub
    End If

    ' Parse and normalize tickers
    Set tickers = ParseTickers(inputText)

    If tickers.Count = 0 Then
        MsgBox "No valid tickers found", vbExclamation
        Exit Sub
    End If

    ' Add to Candidates table
    addedCount = 0
    For Each ticker In tickers
        If AddCandidate(ticker, presetName, sector, bucket) Then
            addedCount = addedCount + 1
        End If
    Next ticker

    MsgBox "Imported " & addedCount & " candidates" & vbCrLf & _
           "Date: " & Format(Date, "yyyy-mm-dd") & vbCrLf & _
           "Preset: " & presetName, vbInformation
End Sub

Function ParseTickers(inputText As String) As Collection
    ' Parse input text and return normalized tickers
    Dim result As New Collection
    Dim lines() As String
    Dim line As String
    Dim parts() As String
    Dim part As String
    Dim ticker As String
    Dim i As Integer

    ' Replace common separators with commas
    inputText = Replace(inputText, vbTab, ",")
    inputText = Replace(inputText, ";", ",")
    inputText = Replace(inputText, "|", ",")
    inputText = Replace(inputText, " ", ",")

    ' Split by newlines and commas
    lines = Split(inputText, vbCrLf)

    For Each line In lines
        If Trim(line) <> "" Then
            parts = Split(line, ",")
            For Each part In parts
                ticker = NormalizeTicker(Trim(part))
                If ticker <> "" Then
                    ' Avoid duplicates
                    On Error Resume Next
                    result.Add ticker, ticker  ' Key = ticker for uniqueness
                    On Error GoTo 0
                End If
            Next part
        End If
    Next line

    Set ParseTickers = result
End Function

Function AddCandidate(ticker As String, presetName As String, sector As String, bucket As String) As Boolean
    ' Add ticker to Candidates table (skip if duplicate)
    Dim tbl As ListObject
    Dim row As ListRow
    Dim isDuplicate As Boolean

    Set tbl = Worksheets("Candidates").ListObjects("tblCandidates")

    ' Check for duplicates (same ticker and date)
    isDuplicate = False
    For Each row In tbl.ListRows
        If row.Range.Columns(2).Value = ticker And _
           row.Range.Columns(1).Value = Date Then
            isDuplicate = True
            Exit For
        End If
    Next row

    If isDuplicate Then
        AddCandidate = False
        Exit Function
    End If

    ' Add new candidate
    Set row = tbl.ListRows.Add
    With row.Range
        .Columns(1).Value = Date  ' Date
        .Columns(2).Value = ticker  ' Ticker
        .Columns(3).Value = presetName  ' Preset
        .Columns(4).Value = sector  ' Sector
        .Columns(5).Value = bucket  ' Bucket
    End With

    AddCandidate = True
End Function

Sub ImportFromPython()
    ' Enhanced import using Python scraper (if available)
    ' Falls back to manual if Python not available

    Dim presetName As String
    Dim queryString As String
    Dim tbl As ListObject
    Dim row As ListRow
    Dim pythonAvailable As Boolean

    presetName = Worksheets("TradeEntry").Range("B5").Value

    If presetName = "" Then
        MsgBox "Please select a Preset first", vbExclamation
        Exit Sub
    End If

    ' Look up query string
    Set tbl = Worksheets("Presets").ListObjects("tblPresets")
    For Each row In tbl.ListRows
        If row.Range.Columns(1).Value = presetName Then
            queryString = row.Range.Columns(2).Value
            Exit For
        End If
    Next row

    If queryString = "" Then
        MsgBox "Preset not found: " & presetName, vbExclamation
        Exit Sub
    End If

    ' Check if Python integration is available
    ' (This will be enhanced when TF_Python_Bridge is added)
    pythonAvailable = False  ' Placeholder

    If pythonAvailable Then
        ' Use Python scraper
        MsgBox "Python scraping not yet implemented" & vbCrLf & _
               "Falling back to manual import", vbInformation
        Call ImportCandidatesPrompt
    Else
        ' Fallback to manual
        MsgBox "Auto-scraping requires Python integration" & vbCrLf & _
               "Opening FINVIZ manually...", vbInformation
        Call OpenPreset
        MsgBox "Copy tickers from browser, then click OK to paste", vbInformation
        Call ImportCandidatesPrompt
    End If
End Sub

Sub ClearTodayCandidates()
    ' Clear today's candidates (useful for re-importing)
    Dim tbl As ListObject
    Dim i As Long
    Dim deletedCount As Integer

    Set tbl = Worksheets("Candidates").ListObjects("tblCandidates")

    deletedCount = 0
    For i = tbl.ListRows.Count To 1 Step -1
        If tbl.ListRows(i).Range.Columns(1).Value = Date Then
            tbl.ListRows(i).Delete
            deletedCount = deletedCount + 1
        End If
    Next i

    MsgBox "Cleared " & deletedCount & " candidates from today", vbInformation
End Sub

Sub ExportCandidatesToClipboard()
    ' Export today's candidates to clipboard (useful for sharing)
    Dim tbl As ListObject
    Dim row As ListRow
    Dim output As String
    Dim dataObj As Object

    Set tbl = Worksheets("Candidates").ListObjects("tblCandidates")

    output = ""
    For Each row In tbl.ListRows
        If row.Range.Columns(1).Value = Date Then
            If output <> "" Then output = output & ", "
            output = output & row.Range.Columns(2).Value
        End If
    Next row

    If output = "" Then
        MsgBox "No candidates found for today", vbExclamation
        Exit Sub
    End If

    ' Copy to clipboard
    On Error Resume Next
    Set dataObj = CreateObject("New:{1C3B4210-F441-11CE-B9EA-00AA006B1A69}")
    dataObj.SetText output
    dataObj.PutInClipboard
    On Error GoTo 0

    MsgBox "Copied " & tbl.ListRows.Count & " tickers to clipboard:" & vbCrLf & _
           Left(output, 100) & IIf(Len(output) > 100, "...", ""), vbInformation
End Sub
