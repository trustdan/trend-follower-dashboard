Attribute VB_Name = "TFTests"
'=============================================================================
' TFTests.bas - VBA Unit Tests for Trading Engine v3
'=============================================================================
' Purpose: Test VBA modules before Windows integration testing
'
' Test Categories:
' 1. JSON Parsing Tests (using sample JSON files)
' 2. Shell Execution Tests (verify can call engine)
' 3. Validation Tests (input validation functions)
' 4. Formatting Tests (display formatters)
'
' Usage: Run tests from "VBA Tests" worksheet or call directly
'
' Expected Setup:
' - test-data/ folder with JSON examples
' - tf-engine.exe in workbook directory (or Setup sheet path)
'
' Created: 2025-10-27 (M19 - VBA Implementation)
'=============================================================================

Option Explicit

'-----------------------------------------------------------------------------
' TEST RESULT TYPE
'-----------------------------------------------------------------------------

Public Type TestResult
    TestName As String
    Passed As Boolean
    Message As String
    Duration As Double
End Type

'-----------------------------------------------------------------------------
' TEST RUNNER
'-----------------------------------------------------------------------------

' RunAllTests executes all test functions and reports results
Public Sub RunAllTests()
    Dim results() As TestResult
    Dim resultCount As Long
    Dim i As Long
    Dim passCount As Long
    Dim failCount As Long

    ' Clear previous results
    ClearTestResults

    ' Run test suites
    WriteTestHeader "JSON PARSING TESTS"
    AddTestResult results, resultCount, Test_ParseSizingJSON()
    AddTestResult results, resultCount, Test_ParseChecklistJSON_Green()
    AddTestResult results, resultCount, Test_ParseChecklistJSON_Yellow()
    AddTestResult results, resultCount, Test_ParseHeatJSON()
    AddTestResult results, resultCount, Test_ParseTimerJSON()
    AddTestResult results, resultCount, Test_ParseSettingsJSON()

    WriteTestHeader "HELPER FUNCTION TESTS"
    AddTestResult results, resultCount, Test_ExtractJSONValue()
    AddTestResult results, resultCount, Test_ExtractJSONArray()
    AddTestResult results, resultCount, Test_GenerateCorrelationID()

    WriteTestHeader "VALIDATION TESTS"
    AddTestResult results, resultCount, Test_ValidateTicker()
    AddTestResult results, resultCount, Test_ValidatePositiveNumber()

    WriteTestHeader "FORMATTING TESTS"
    AddTestResult results, resultCount, Test_FormatCurrency()
    AddTestResult results, resultCount, Test_FormatPercent()

    WriteTestHeader "SHELL EXECUTION TESTS"
    AddTestResult results, resultCount, Test_ShellExecution()

    ' Count results
    passCount = 0
    failCount = 0
    For i = 0 To resultCount - 1
        If results(i).Passed Then
            passCount = passCount + 1
        Else
            failCount = failCount + 1
        End If
    Next i

    ' Display summary
    WriteTestSummary passCount, failCount, results, resultCount
End Sub

' AddTestResult adds a test result to the results array
Private Sub AddTestResult(ByRef results() As TestResult, ByRef count As Long, ByVal result As TestResult)
    ReDim Preserve results(0 To count)
    results(count) = result
    count = count + 1

    ' Write result to worksheet
    WriteTestResult result
End Sub

'-----------------------------------------------------------------------------
' JSON PARSING TESTS
'-----------------------------------------------------------------------------

' Test parsing position sizing JSON
Private Function Test_ParseSizingJSON() As TestResult
    Dim result As TestResult
    Dim startTime As Double
    Dim sizeResult As TFSizingResult
    Dim jsonStr As String

    result.TestName = "ParseSizingJSON"
    startTime = Timer

    On Error GoTo ErrorHandler

    ' Sample JSON (from test-data/json-examples/responses/size-stock-success.json)
    jsonStr = "{" & vbCrLf & _
              "  ""risk_dollars"": 75," & vbCrLf & _
              "  ""stop_distance"": 3," & vbCrLf & _
              "  ""initial_stop"": 177," & vbCrLf & _
              "  ""shares"": 25," & vbCrLf & _
              "  ""contracts"": 0," & vbCrLf & _
              "  ""actual_risk"": 75," & vbCrLf & _
              "  ""method"": ""stock""" & vbCrLf & _
              "}"

    ' Parse
    sizeResult = TFHelpers.ParseSizingJSON(jsonStr)

    ' Verify
    If sizeResult.RiskDollars <> 75 Then GoTo FailTest
    If sizeResult.StopDistance <> 3 Then GoTo FailTest
    If sizeResult.InitialStop <> 177 Then GoTo FailTest
    If sizeResult.Shares <> 25 Then GoTo FailTest
    If sizeResult.Method <> "stock" Then GoTo FailTest

    result.Passed = True
    result.Message = "All fields parsed correctly"
    result.Duration = Timer - startTime
    Test_ParseSizingJSON = result
    Exit Function

FailTest:
    result.Passed = False
    result.Message = "Field validation failed"
    result.Duration = Timer - startTime
    Test_ParseSizingJSON = result
    Exit Function

ErrorHandler:
    result.Passed = False
    result.Message = "Error: " & Err.Description
    result.Duration = Timer - startTime
    Test_ParseSizingJSON = result
End Function

' Test parsing checklist JSON (GREEN)
Private Function Test_ParseChecklistJSON_Green() As TestResult
    Dim result As TestResult
    Dim startTime As Double
    Dim checkResult As TFChecklistResult
    Dim jsonStr As String

    result.TestName = "ParseChecklistJSON_Green"
    startTime = Timer

    On Error GoTo ErrorHandler

    jsonStr = "{" & vbCrLf & _
              "  ""banner"": ""GREEN""," & vbCrLf & _
              "  ""missing_count"": 0," & vbCrLf & _
              "  ""missing_items"": []," & vbCrLf & _
              "  ""evaluation_timestamp"": ""2025-10-27T20:11:16.30438317-05:00""," & vbCrLf & _
              "  ""allow_save"": true" & vbCrLf & _
              "}"

    checkResult = TFHelpers.ParseChecklistJSON(jsonStr)

    If checkResult.Banner <> "GREEN" Then GoTo FailTest
    If checkResult.MissingCount <> 0 Then GoTo FailTest
    If Not checkResult.AllowSave Then GoTo FailTest

    result.Passed = True
    result.Message = "GREEN banner parsed correctly"
    result.Duration = Timer - startTime
    Test_ParseChecklistJSON_Green = result
    Exit Function

FailTest:
    result.Passed = False
    result.Message = "Field validation failed"
    result.Duration = Timer - startTime
    Test_ParseChecklistJSON_Green = result
    Exit Function

ErrorHandler:
    result.Passed = False
    result.Message = "Error: " & Err.Description
    result.Duration = Timer - startTime
    Test_ParseChecklistJSON_Green = result
End Function

' Test parsing checklist JSON (YELLOW)
Private Function Test_ParseChecklistJSON_Yellow() As TestResult
    Dim result As TestResult
    Dim startTime As Double
    Dim checkResult As TFChecklistResult
    Dim jsonStr As String

    result.TestName = "ParseChecklistJSON_Yellow"
    startTime = Timer

    On Error GoTo ErrorHandler

    jsonStr = "{" & vbCrLf & _
              "  ""banner"": ""YELLOW""," & vbCrLf & _
              "  ""missing_count"": 1," & vbCrLf & _
              "  ""missing_items"": [""Higher high""]," & vbCrLf & _
              "  ""evaluation_timestamp"": ""2025-10-27T20:11:16-05:00""," & vbCrLf & _
              "  ""allow_save"": false" & vbCrLf & _
              "}"

    checkResult = TFHelpers.ParseChecklistJSON(jsonStr)

    If checkResult.Banner <> "YELLOW" Then GoTo FailTest
    If checkResult.MissingCount <> 1 Then GoTo FailTest
    If checkResult.AllowSave Then GoTo FailTest
    If checkResult.MissingItems <> "Higher high" Then GoTo FailTest

    result.Passed = True
    result.Message = "YELLOW banner with missing items parsed correctly"
    result.Duration = Timer - startTime
    Test_ParseChecklistJSON_Yellow = result
    Exit Function

FailTest:
    result.Passed = False
    result.Message = "Field validation failed (got: " & checkResult.MissingItems & ")"
    result.Duration = Timer - startTime
    Test_ParseChecklistJSON_Yellow = result
    Exit Function

ErrorHandler:
    result.Passed = False
    result.Message = "Error: " & Err.Description
    result.Duration = Timer - startTime
    Test_ParseChecklistJSON_Yellow = result
End Function

' Test parsing heat JSON
Private Function Test_ParseHeatJSON() As TestResult
    Dim result As TestResult
    Dim startTime As Double
    Dim heatResult As TFHeatResult
    Dim jsonStr As String

    result.TestName = "ParseHeatJSON"
    startTime = Timer

    On Error GoTo ErrorHandler

    jsonStr = "{" & vbCrLf & _
              "  ""current_portfolio_heat"": 0," & vbCrLf & _
              "  ""new_portfolio_heat"": 75," & vbCrLf & _
              "  ""portfolio_heat_pct"": 18.75," & vbCrLf & _
              "  ""portfolio_cap"": 400," & vbCrLf & _
              "  ""portfolio_cap_exceeded"": false," & vbCrLf & _
              "  ""portfolio_overage"": 0," & vbCrLf & _
              "  ""current_bucket_heat"": 0," & vbCrLf & _
              "  ""new_bucket_heat"": 75," & vbCrLf & _
              "  ""bucket_heat_pct"": 50," & vbCrLf & _
              "  ""bucket_cap"": 150," & vbCrLf & _
              "  ""bucket_cap_exceeded"": false," & vbCrLf & _
              "  ""bucket_overage"": 0," & vbCrLf & _
              "  ""allowed"": true" & vbCrLf & _
              "}"

    heatResult = TFHelpers.ParseHeatJSON(jsonStr)

    If heatResult.NewPortfolioHeat <> 75 Then GoTo FailTest
    If heatResult.PortfolioCap <> 400 Then GoTo FailTest
    If heatResult.PortfolioCapExceeded Then GoTo FailTest
    If Not heatResult.Allowed Then GoTo FailTest

    result.Passed = True
    result.Message = "Heat check parsed correctly"
    result.Duration = Timer - startTime
    Test_ParseHeatJSON = result
    Exit Function

FailTest:
    result.Passed = False
    result.Message = "Field validation failed"
    result.Duration = Timer - startTime
    Test_ParseHeatJSON = result
    Exit Function

ErrorHandler:
    result.Passed = False
    result.Message = "Error: " & Err.Description
    result.Duration = Timer - startTime
    Test_ParseHeatJSON = result
End Function

' Test parsing timer JSON
Private Function Test_ParseTimerJSON() As TestResult
    Dim result As TestResult
    Dim startTime As Double
    Dim timerResult As TFTimerResult
    Dim jsonStr As String

    result.TestName = "ParseTimerJSON"
    startTime = Timer

    On Error GoTo ErrorHandler

    jsonStr = "{" & vbCrLf & _
              "  ""brake_cleared"": false," & vbCrLf & _
              "  ""elapsed_seconds"": 0," & vbCrLf & _
              "  ""remaining_seconds"": 119," & vbCrLf & _
              "  ""started_at"": ""2025-10-27T20:11:16-05:00""," & vbCrLf & _
              "  ""ticker"": ""AAPL""," & vbCrLf & _
              "  ""timer_active"": true" & vbCrLf & _
              "}"

    timerResult = TFHelpers.ParseTimerJSON(jsonStr)

    If Not timerResult.TimerActive Then GoTo FailTest
    If timerResult.BrakeCleared Then GoTo FailTest
    If timerResult.RemainingSeconds <> 119 Then GoTo FailTest
    If timerResult.Ticker <> "AAPL" Then GoTo FailTest

    result.Passed = True
    result.Message = "Timer check parsed correctly"
    result.Duration = Timer - startTime
    Test_ParseTimerJSON = result
    Exit Function

FailTest:
    result.Passed = False
    result.Message = "Field validation failed"
    result.Duration = Timer - startTime
    Test_ParseTimerJSON = result
    Exit Function

ErrorHandler:
    result.Passed = False
    result.Message = "Error: " & Err.Description
    result.Duration = Timer - startTime
    Test_ParseTimerJSON = result
End Function

' Test parsing settings JSON
Private Function Test_ParseSettingsJSON() As TestResult
    Dim result As TestResult
    Dim startTime As Double
    Dim settings As TFSettings
    Dim jsonStr As String

    result.TestName = "ParseSettingsJSON"
    startTime = Timer

    On Error GoTo ErrorHandler

    jsonStr = "{" & vbCrLf & _
              "  ""BucketHeatCap_pct"": ""0.015""," & vbCrLf & _
              "  ""Equity_E"": ""10000""," & vbCrLf & _
              "  ""HeatCap_H_pct"": ""0.04""," & vbCrLf & _
              "  ""RiskPct_r"": ""0.0075""," & vbCrLf & _
              "  ""StopMultiple_K"": ""2""" & vbCrLf & _
              "}"

    settings = TFHelpers.ParseSettingsJSON(jsonStr)

    If settings.Equity_E <> 10000 Then GoTo FailTest
    If settings.RiskPct_r <> 0.0075 Then GoTo FailTest
    If settings.StopMultiple_K <> 2 Then GoTo FailTest

    result.Passed = True
    result.Message = "Settings parsed correctly"
    result.Duration = Timer - startTime
    Test_ParseSettingsJSON = result
    Exit Function

FailTest:
    result.Passed = False
    result.Message = "Field validation failed"
    result.Duration = Timer - startTime
    Test_ParseSettingsJSON = result
    Exit Function

ErrorHandler:
    result.Passed = False
    result.Message = "Error: " & Err.Description
    result.Duration = Timer - startTime
    Test_ParseSettingsJSON = result
End Function

'-----------------------------------------------------------------------------
' HELPER FUNCTION TESTS
'-----------------------------------------------------------------------------

' Test ExtractJSONValue function
Private Function Test_ExtractJSONValue() As TestResult
    Dim result As TestResult
    Dim startTime As Double
    Dim jsonStr As String
    Dim value As String

    result.TestName = "ExtractJSONValue"
    startTime = Timer

    On Error GoTo ErrorHandler

    jsonStr = "{""name"": ""AAPL"", ""price"": 180.50, ""active"": true}"

    value = TFHelpers.ExtractJSONValue(jsonStr, "name")
    If value <> "AAPL" Then GoTo FailTest

    value = TFHelpers.ExtractJSONValue(jsonStr, "price")
    If value <> "180.50" Then GoTo FailTest

    value = TFHelpers.ExtractJSONValue(jsonStr, "active")
    If value <> "true" Then GoTo FailTest

    result.Passed = True
    result.Message = "JSON value extraction works"
    result.Duration = Timer - startTime
    Test_ExtractJSONValue = result
    Exit Function

FailTest:
    result.Passed = False
    result.Message = "Value extraction failed (got: " & value & ")"
    result.Duration = Timer - startTime
    Test_ExtractJSONValue = result
    Exit Function

ErrorHandler:
    result.Passed = False
    result.Message = "Error: " & Err.Description
    result.Duration = Timer - startTime
    Test_ExtractJSONValue = result
End Function

' Test ExtractJSONArray function
Private Function Test_ExtractJSONArray() As TestResult
    Dim result As TestResult
    Dim startTime As Double
    Dim arrayStr As String
    Dim items As Collection

    result.TestName = "ExtractJSONArray"
    startTime = Timer

    On Error GoTo ErrorHandler

    arrayStr = "[""item1"", ""item2"", ""item3""]"
    Set items = TFHelpers.ExtractJSONArray(arrayStr)

    If items.Count <> 3 Then GoTo FailTest
    If items(1) <> "item1" Then GoTo FailTest
    If items(2) <> "item2" Then GoTo FailTest
    If items(3) <> "item3" Then GoTo FailTest

    result.Passed = True
    result.Message = "JSON array extraction works"
    result.Duration = Timer - startTime
    Test_ExtractJSONArray = result
    Exit Function

FailTest:
    result.Passed = False
    result.Message = "Array extraction failed"
    result.Duration = Timer - startTime
    Test_ExtractJSONArray = result
    Exit Function

ErrorHandler:
    result.Passed = False
    result.Message = "Error: " & Err.Description
    result.Duration = Timer - startTime
    Test_ExtractJSONArray = result
End Function

' Test GenerateCorrelationID function
Private Function Test_GenerateCorrelationID() As TestResult
    Dim result As TestResult
    Dim startTime As Double
    Dim corrID As String

    result.TestName = "GenerateCorrelationID"
    startTime = Timer

    On Error GoTo ErrorHandler

    corrID = TFHelpers.GenerateCorrelationID()

    ' Should be format: YYYYMMDD-HHMMSS-XXXX (23 chars)
    If Len(corrID) <> 23 Then GoTo FailTest
    If Mid(corrID, 9, 1) <> "-" Then GoTo FailTest
    If Mid(corrID, 16, 1) <> "-" Then GoTo FailTest

    result.Passed = True
    result.Message = "Correlation ID format correct: " & corrID
    result.Duration = Timer - startTime
    Test_GenerateCorrelationID = result
    Exit Function

FailTest:
    result.Passed = False
    result.Message = "Invalid format: " & corrID
    result.Duration = Timer - startTime
    Test_GenerateCorrelationID = result
    Exit Function

ErrorHandler:
    result.Passed = False
    result.Message = "Error: " & Err.Description
    result.Duration = Timer - startTime
    Test_GenerateCorrelationID = result
End Function

'-----------------------------------------------------------------------------
' VALIDATION TESTS
'-----------------------------------------------------------------------------

' Test ValidateTicker function
Private Function Test_ValidateTicker() As TestResult
    Dim result As TestResult
    Dim startTime As Double

    result.TestName = "ValidateTicker"
    startTime = Timer

    On Error GoTo ErrorHandler

    If Not TFHelpers.ValidateTicker("AAPL") Then GoTo FailTest
    If Not TFHelpers.ValidateTicker("BRK-B") Then GoTo FailTest
    If Not TFHelpers.ValidateTicker("BRK.B") Then GoTo FailTest
    If TFHelpers.ValidateTicker("") Then GoTo FailTest
    If TFHelpers.ValidateTicker("ABC@123") Then GoTo FailTest

    result.Passed = True
    result.Message = "Ticker validation works"
    result.Duration = Timer - startTime
    Test_ValidateTicker = result
    Exit Function

FailTest:
    result.Passed = False
    result.Message = "Validation logic incorrect"
    result.Duration = Timer - startTime
    Test_ValidateTicker = result
    Exit Function

ErrorHandler:
    result.Passed = False
    result.Message = "Error: " & Err.Description
    result.Duration = Timer - startTime
    Test_ValidateTicker = result
End Function

' Test ValidatePositiveNumber function
Private Function Test_ValidatePositiveNumber() As TestResult
    Dim result As TestResult
    Dim startTime As Double

    result.TestName = "ValidatePositiveNumber"
    startTime = Timer

    On Error GoTo ErrorHandler

    If Not TFHelpers.ValidatePositiveNumber(123.45) Then GoTo FailTest
    If Not TFHelpers.ValidatePositiveNumber(0.001) Then GoTo FailTest
    If TFHelpers.ValidatePositiveNumber(0) Then GoTo FailTest
    If TFHelpers.ValidatePositiveNumber(-5) Then GoTo FailTest

    result.Passed = True
    result.Message = "Number validation works"
    result.Duration = Timer - startTime
    Test_ValidatePositiveNumber = result
    Exit Function

FailTest:
    result.Passed = False
    result.Message = "Validation logic incorrect"
    result.Duration = Timer - startTime
    Test_ValidatePositiveNumber = result
    Exit Function

ErrorHandler:
    result.Passed = False
    result.Message = "Error: " & Err.Description
    result.Duration = Timer - startTime
    Test_ValidatePositiveNumber = result
End Function

'-----------------------------------------------------------------------------
' FORMATTING TESTS
'-----------------------------------------------------------------------------

' Test FormatCurrency function
Private Function Test_FormatCurrency() As TestResult
    Dim result As TestResult
    Dim startTime As Double
    Dim formatted As String

    result.TestName = "FormatCurrency"
    startTime = Timer

    On Error GoTo ErrorHandler

    formatted = TFHelpers.FormatCurrency(1234.56)
    If formatted <> "$1,234.56" Then GoTo FailTest

    formatted = TFHelpers.FormatCurrency(75)
    If formatted <> "$75.00" Then GoTo FailTest

    result.Passed = True
    result.Message = "Currency formatting works"
    result.Duration = Timer - startTime
    Test_FormatCurrency = result
    Exit Function

FailTest:
    result.Passed = False
    result.Message = "Format incorrect (got: " & formatted & ")"
    result.Duration = Timer - startTime
    Test_FormatCurrency = result
    Exit Function

ErrorHandler:
    result.Passed = False
    result.Message = "Error: " & Err.Description
    result.Duration = Timer - startTime
    Test_FormatCurrency = result
End Function

' Test FormatPercent function
Private Function Test_FormatPercent() As TestResult
    Dim result As TestResult
    Dim startTime As Double
    Dim formatted As String

    result.TestName = "FormatPercent"
    startTime = Timer

    On Error GoTo ErrorHandler

    formatted = TFHelpers.FormatPercent(0.0075)
    If formatted <> "0.75%" Then GoTo FailTest

    formatted = TFHelpers.FormatPercent(0.04)
    If formatted <> "4.00%" Then GoTo FailTest

    result.Passed = True
    result.Message = "Percent formatting works"
    result.Duration = Timer - startTime
    Test_FormatPercent = result
    Exit Function

FailTest:
    result.Passed = False
    result.Message = "Format incorrect (got: " & formatted & ")"
    result.Duration = Timer - startTime
    Test_FormatPercent = result
    Exit Function

ErrorHandler:
    result.Passed = False
    result.Message = "Error: " & Err.Description
    result.Duration = Timer - startTime
    Test_FormatPercent = result
End Function

'-----------------------------------------------------------------------------
' SHELL EXECUTION TEST
'-----------------------------------------------------------------------------

' Test shell execution (calls tf-engine --version)
Private Function Test_ShellExecution() As TestResult
    Dim result As TestResult
    Dim startTime As Double
    Dim cmdResult As TFCommandResult

    result.TestName = "ShellExecution"
    startTime = Timer

    On Error GoTo ErrorHandler

    ' Try to execute --version command (simplest test)
    cmdResult = TFEngine.ExecuteCommand("--version")

    ' Check if execution worked
    If cmdResult.ExitCode <> 0 Then
        result.Passed = False
        result.Message = "Engine not found or failed (exit code " & cmdResult.ExitCode & ")"
        result.Duration = Timer - startTime
        Test_ShellExecution = result
        Exit Function
    End If

    result.Passed = True
    result.Message = "Engine executable found and runs"
    result.Duration = Timer - startTime
    Test_ShellExecution = result
    Exit Function

ErrorHandler:
    result.Passed = False
    result.Message = "Error: " & Err.Description
    result.Duration = Timer - startTime
    Test_ShellExecution = result
End Function

'-----------------------------------------------------------------------------
' TEST OUTPUT HELPERS
'-----------------------------------------------------------------------------

' ClearTestResults clears the test results worksheet
Private Sub ClearTestResults()
    On Error Resume Next
    Dim ws As Worksheet
    Set ws = ThisWorkbook.Worksheets("VBA Tests")
    If Not ws Is Nothing Then
        ws.Range("A5:Z1000").ClearContents
        ws.Range("A5:Z1000").Interior.Color = xlNone
    End If
End Sub

' WriteTestHeader writes a section header
Private Sub WriteTestHeader(ByVal header As String)
    On Error Resume Next
    Dim ws As Worksheet
    Dim nextRow As Long

    Set ws = ThisWorkbook.Worksheets("VBA Tests")
    If ws Is Nothing Then Exit Sub

    nextRow = ws.Cells(ws.Rows.Count, "A").End(xlUp).Row + 1
    If nextRow < 5 Then nextRow = 5

    ws.Cells(nextRow, 1).Value = header
    ws.Cells(nextRow, 1).Font.Bold = True
    ws.Cells(nextRow, 1).Font.Size = 12
End Sub

' WriteTestResult writes a single test result to worksheet
Private Sub WriteTestResult(ByVal result As TestResult)
    On Error Resume Next
    Dim ws As Worksheet
    Dim nextRow As Long

    Set ws = ThisWorkbook.Worksheets("VBA Tests")
    If ws Is Nothing Then Exit Sub

    nextRow = ws.Cells(ws.Rows.Count, "A").End(xlUp).Row + 1

    ' Write result
    ws.Cells(nextRow, 1).Value = result.TestName
    ws.Cells(nextRow, 2).Value = IIf(result.Passed, "PASS", "FAIL")
    ws.Cells(nextRow, 3).Value = result.Message
    ws.Cells(nextRow, 4).Value = Format(result.Duration, "0.000") & "s"

    ' Color code
    If result.Passed Then
        ws.Cells(nextRow, 2).Interior.Color = RGB(198, 239, 206)  ' Green
        ws.Cells(nextRow, 2).Value = "✅ PASS"
    Else
        ws.Cells(nextRow, 2).Interior.Color = RGB(255, 199, 206)  ' Red
        ws.Cells(nextRow, 2).Value = "❌ FAIL"
    End If
End Sub

' WriteTestSummary writes the final summary
Private Sub WriteTestSummary(ByVal passCount As Long, ByVal failCount As Long, _
                             ByRef results() As TestResult, ByVal count As Long)
    On Error Resume Next
    Dim ws As Worksheet
    Dim nextRow As Long
    Dim totalTime As Double
    Dim i As Long

    Set ws = ThisWorkbook.Worksheets("VBA Tests")
    If ws Is Nothing Then Exit Sub

    ' Calculate total time
    For i = 0 To count - 1
        totalTime = totalTime + results(i).Duration
    Next i

    nextRow = ws.Cells(ws.Rows.Count, "A").End(xlUp).Row + 2

    ' Write summary
    ws.Cells(nextRow, 1).Value = "SUMMARY"
    ws.Cells(nextRow, 1).Font.Bold = True
    ws.Cells(nextRow, 1).Font.Size = 14

    nextRow = nextRow + 1
    ws.Cells(nextRow, 1).Value = "Total Tests:"
    ws.Cells(nextRow, 2).Value = count

    nextRow = nextRow + 1
    ws.Cells(nextRow, 1).Value = "Passed:"
    ws.Cells(nextRow, 2).Value = passCount
    ws.Cells(nextRow, 2).Interior.Color = RGB(198, 239, 206)

    nextRow = nextRow + 1
    ws.Cells(nextRow, 1).Value = "Failed:"
    ws.Cells(nextRow, 2).Value = failCount
    If failCount > 0 Then
        ws.Cells(nextRow, 2).Interior.Color = RGB(255, 199, 206)
    End If

    nextRow = nextRow + 1
    ws.Cells(nextRow, 1).Value = "Total Time:"
    ws.Cells(nextRow, 2).Value = Format(totalTime, "0.000") & "s"

    nextRow = nextRow + 1
    ws.Cells(nextRow, 1).Value = "Result:"
    If failCount = 0 Then
        ws.Cells(nextRow, 2).Value = "✅ ALL TESTS PASSED"
        ws.Cells(nextRow, 2).Font.Bold = True
        ws.Cells(nextRow, 2).Interior.Color = RGB(198, 239, 206)
    Else
        ws.Cells(nextRow, 2).Value = "❌ " & failCount & " TESTS FAILED"
        ws.Cells(nextRow, 2).Font.Bold = True
        ws.Cells(nextRow, 2).Interior.Color = RGB(255, 199, 206)
    End If

    ' Auto-fit columns
    ws.Columns("A:D").AutoFit
End Sub

'=============================================================================
' NOTES ON TESTING
'=============================================================================
'
' TEST PHILOSOPHY:
' - These are VBA unit tests (not full integration tests)
' - Integration tests happen in M21 (Windows manual testing)
' - Focus: Verify JSON parsing and basic shell execution
'
' TEST COVERAGE:
' - JSON parsing for all response types
' - Helper functions (validation, formatting)
' - Shell execution (can we call the engine?)
' - NOT TESTED HERE: Full command workflows (that's M21)
'
' RUNNING TESTS:
' 1. Manual: Call RunAllTests() from VBA editor
' 2. Button: Add button on "VBA Tests" worksheet linked to RunAllTests
' 3. Automated: Use test runner batch script (M21)
'
' EXTENDING TESTS:
' To add a new test:
' 1. Create Private Function Test_YourTestName() As TestResult
' 2. Add call to AddTestResult() in RunAllTests()
' 3. Follow existing pattern (setup, execute, verify, report)
'
' ERROR HANDLING:
' - All test functions have On Error GoTo ErrorHandler
' - Errors are caught and reported as test failures
' - Tests should never crash Excel
'
'=============================================================================
