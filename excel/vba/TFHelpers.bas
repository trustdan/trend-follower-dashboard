Attribute VB_Name = "TFHelpers"
'=============================================================================
' TFHelpers.bas - Utility Functions for Trading Engine v3
'=============================================================================
' Purpose: JSON parsing, error handling, logging, and utility functions
'
' Key Functions:
' - ParseXXXJSON() - Parse engine JSON responses into typed results
' - GenerateCorrelationID() - Create unique IDs for log correlation
' - LogMessage() - Write to debug log with correlation IDs
' - ExtractJSONValue() - Simple JSON key-value extraction
'
' Philosophy: Simple string parsing (no external dependencies)
' The Go engine validates JSON schemas - VBA just extracts values
'
' Created: 2025-10-27 (M19 - VBA Implementation)
'=============================================================================

Option Explicit

'-----------------------------------------------------------------------------
' CONSTANTS
'-----------------------------------------------------------------------------

Private Const LOG_FILE_NAME As String = "TradingSystem_Debug.log"
Private Const MAX_LOG_SIZE_KB As Long = 5120  ' 5 MB before rotation

'-----------------------------------------------------------------------------
' CORRELATION ID GENERATION
'-----------------------------------------------------------------------------

' GenerateCorrelationID creates a unique ID for log correlation
' Format: YYYYMMDD-HHMMSS-RND (e.g., "20251027-143052-7A3F")
Public Function GenerateCorrelationID() As String
    Dim timestamp As String
    Dim milliseconds As String
    Dim randomPart As String

    ' Format: YYYYMMDD-HHMMSSFFF (with milliseconds)
    timestamp = Format(Now, "YYYYMMDD-HHNNSS")

    ' Add 3-digit milliseconds
    milliseconds = Right("000" & Format(Timer Mod 1 * 1000, "0"), 3)

    ' Add 4 random hex characters
    Randomize
    randomPart = Right("0000" & Hex(Int(Rnd * 65535)), 4)

    GenerateCorrelationID = timestamp & milliseconds & "-" & randomPart
End Function

'-----------------------------------------------------------------------------
' LOGGING
'-----------------------------------------------------------------------------

' LogMessage writes a timestamped message to the debug log
Public Sub LogMessage(ByVal corrID As String, ByVal level As String, ByVal message As String)
    On Error Resume Next  ' Don't fail if logging fails

    Dim logPath As String
    Dim logLine As String
    Dim fileNum As Integer

    ' Build log path (same directory as workbook)
    logPath = ThisWorkbook.Path & "\" & LOG_FILE_NAME

    ' Format: [TIMESTAMP] [LEVEL] [CORR_ID] Message
    logLine = "[" & Format(Now, "YYYY-MM-DD HH:NN:SS") & "] " & _
              "[" & UCase(level) & "] " & _
              "[" & corrID & "] " & _
              message

    ' Check log size and rotate if needed
    If Dir(logPath) <> "" Then
        If FileLen(logPath) > (MAX_LOG_SIZE_KB * 1024) Then
            RotateLogFile logPath
        End If
    End If

    ' Append to log file
    fileNum = FreeFile
    Open logPath For Append As #fileNum
    Print #fileNum, logLine
    Close #fileNum
End Sub

' RotateLogFile renames the current log and starts a new one
Private Sub RotateLogFile(ByVal logPath As String)
    On Error Resume Next

    Dim backupPath As String
    backupPath = Replace(logPath, ".log", "_" & Format(Now, "YYYYMMDD_HHNNSS") & ".log")

    ' Rename current log to backup
    Name logPath As backupPath
End Sub

'-----------------------------------------------------------------------------
' JSON PARSING - CORE UTILITIES
'-----------------------------------------------------------------------------

' ExtractJSONValue extracts a value for a given key from JSON string
' Simple parser - expects well-formed JSON from engine (validated in Go)
' Returns empty string if key not found
Public Function ExtractJSONValue(ByVal jsonStr As String, ByVal key As String) As String
    On Error GoTo ErrorHandler

    Dim startPos As Long
    Dim endPos As Long
    Dim searchKey As String
    Dim value As String

    ' Search for "key": or "key":
    searchKey = """" & key & """"
    startPos = InStr(jsonStr, searchKey)

    If startPos = 0 Then
        ExtractJSONValue = ""
        Exit Function
    End If

    ' Move past the key and colon
    startPos = InStr(startPos, jsonStr, ":")
    If startPos = 0 Then
        ExtractJSONValue = ""
        Exit Function
    End If
    startPos = startPos + 1

    ' Skip whitespace
    Do While Mid(jsonStr, startPos, 1) = " " Or Mid(jsonStr, startPos, 1) = vbTab
        startPos = startPos + 1
    Loop

    ' Check value type
    Dim firstChar As String
    firstChar = Mid(jsonStr, startPos, 1)

    If firstChar = """" Then
        ' String value - find closing quote
        startPos = startPos + 1
        endPos = InStr(startPos, jsonStr, """")
        value = Mid(jsonStr, startPos, endPos - startPos)
    ElseIf firstChar = "[" Then
        ' Array - find closing bracket
        endPos = InStr(startPos, jsonStr, "]")
        value = Mid(jsonStr, startPos, endPos - startPos + 1)
    ElseIf firstChar = "{" Then
        ' Object - find closing brace (simplified - doesn't handle nested objects)
        endPos = InStr(startPos, jsonStr, "}")
        value = Mid(jsonStr, startPos, endPos - startPos + 1)
    Else
        ' Number, boolean, or null - find next delimiter
        endPos = startPos
        Do While endPos <= Len(jsonStr)
            Dim ch As String
            ch = Mid(jsonStr, endPos, 1)
            If ch = "," Or ch = "}" Or ch = "]" Or ch = vbCr Or ch = vbLf Then
                Exit Do
            End If
            endPos = endPos + 1
        Loop
        value = Trim(Mid(jsonStr, startPos, endPos - startPos))
    End If

    ExtractJSONValue = value
    Exit Function

ErrorHandler:
    ExtractJSONValue = ""
End Function

' ExtractJSONArray extracts array items from a JSON array string
' Input: '["item1", "item2", "item3"]' or empty array '[]'
' Output: Collection of strings
Public Function ExtractJSONArray(ByVal arrayStr As String) As Collection
    Dim result As New Collection
    Dim items As String
    Dim parts() As String
    Dim i As Long

    On Error GoTo ErrorHandler

    ' Remove brackets and whitespace
    items = Trim(arrayStr)
    If Left(items, 1) = "[" Then items = Mid(items, 2)
    If Right(items, 1) = "]" Then items = Left(items, Len(items) - 1)
    items = Trim(items)

    ' Empty array
    If items = "" Then
        Set ExtractJSONArray = result
        Exit Function
    End If

    ' Split by comma (simple - doesn't handle nested arrays)
    parts = Split(items, ",")

    For i = LBound(parts) To UBound(parts)
        Dim item As String
        item = Trim(parts(i))

        ' Remove quotes if present
        If Left(item, 1) = """" And Right(item, 1) = """" Then
            item = Mid(item, 2, Len(item) - 2)
        End If

        result.Add item
    Next i

    Set ExtractJSONArray = result
    Exit Function

ErrorHandler:
    Set ExtractJSONArray = result
End Function

'-----------------------------------------------------------------------------
' JSON PARSING - TYPED PARSERS
'-----------------------------------------------------------------------------

' ParseSizingJSON parses position sizing JSON into TFSizingResult
Public Sub ParseSizingJSON(ByVal jsonStr As String, ByRef result As TFSizingResult)
    On Error GoTo ErrorHandler
    result.Success = False
    result.RiskDollars = CDbl(ExtractJSONValue(jsonStr, "risk_dollars"))
    result.StopDistance = CDbl(ExtractJSONValue(jsonStr, "stop_distance"))
    result.InitialStop = CDbl(ExtractJSONValue(jsonStr, "initial_stop"))
    result.Shares = CLng(ExtractJSONValue(jsonStr, "shares"))
    result.Contracts = CLng(ExtractJSONValue(jsonStr, "contracts"))
    result.ActualRisk = CDbl(ExtractJSONValue(jsonStr, "actual_risk"))
    result.Method = ExtractJSONValue(jsonStr, "method")
    result.Success = True
    Exit Sub
ErrorHandler:
    result.Success = False
End Sub

' ParseChecklistJSON parses checklist evaluation JSON into TFChecklistResult
Public Sub ParseChecklistJSON(ByVal jsonStr As String, ByRef result As TFChecklistResult)
    On Error GoTo ErrorHandler
    Dim missingArray As Collection
    Dim missingStr As String
    Dim item As Variant

    result.Success = False
    result.Banner = ExtractJSONValue(jsonStr, "banner")
    result.MissingCount = CLng(ExtractJSONValue(jsonStr, "missing_count"))
    result.EvaluationTimestamp = ExtractJSONValue(jsonStr, "evaluation_timestamp")

    ' Parse allow_save (boolean)
    Dim allowSaveStr As String
    allowSaveStr = LCase(ExtractJSONValue(jsonStr, "allow_save"))
    result.AllowSave = (allowSaveStr = "true")

    ' Parse missing_items array
    Dim arrayStr As String
    arrayStr = ExtractJSONValue(jsonStr, "missing_items")
    Set missingArray = ExtractJSONArray(arrayStr)

    ' Convert collection to comma-separated string
    missingStr = ""
    For Each item In missingArray
        If missingStr <> "" Then missingStr = missingStr & ", "
        missingStr = missingStr & item
    Next item
    result.MissingItems = missingStr
    result.Success = True
    Exit Sub
ErrorHandler:
    result.Success = False
End Sub

' ParseHeatJSON parses heat check JSON into TFHeatResult
Public Sub ParseHeatJSON(ByVal jsonStr As String, ByRef result As TFHeatResult)
    On Error GoTo ErrorHandler
    result.Success = False
    result.CurrentPortfolioHeat = CDbl(ExtractJSONValue(jsonStr, "current_portfolio_heat"))
    result.NewPortfolioHeat = CDbl(ExtractJSONValue(jsonStr, "new_portfolio_heat"))
    result.PortfolioHeatPct = CDbl(ExtractJSONValue(jsonStr, "portfolio_heat_pct"))
    result.PortfolioCap = CDbl(ExtractJSONValue(jsonStr, "portfolio_cap"))
    result.PortfolioCapExceeded = (LCase(ExtractJSONValue(jsonStr, "portfolio_cap_exceeded")) = "true")
    result.PortfolioOverage = CDbl(ExtractJSONValue(jsonStr, "portfolio_overage"))

    result.CurrentBucketHeat = CDbl(ExtractJSONValue(jsonStr, "current_bucket_heat"))
    result.NewBucketHeat = CDbl(ExtractJSONValue(jsonStr, "new_bucket_heat"))
    result.BucketHeatPct = CDbl(ExtractJSONValue(jsonStr, "bucket_heat_pct"))
    result.BucketCap = CDbl(ExtractJSONValue(jsonStr, "bucket_cap"))
    result.BucketCapExceeded = (LCase(ExtractJSONValue(jsonStr, "bucket_cap_exceeded")) = "true")
    result.BucketOverage = CDbl(ExtractJSONValue(jsonStr, "bucket_overage"))

    result.Allowed = (LCase(ExtractJSONValue(jsonStr, "allowed")) = "true")
    result.Success = True
    Exit Sub
ErrorHandler:
    result.Success = False
End Sub

' ParseTimerJSON parses impulse timer JSON into TFTimerResult
Public Sub ParseTimerJSON(ByVal jsonStr As String, ByRef result As TFTimerResult)
    result.TimerActive = (LCase(ExtractJSONValue(jsonStr, "timer_active")) = "true")
    result.BrakeCleared = (LCase(ExtractJSONValue(jsonStr, "brake_cleared")) = "true")
    result.ElapsedSeconds = CLng(ExtractJSONValue(jsonStr, "elapsed_seconds"))
    result.RemainingSeconds = CLng(ExtractJSONValue(jsonStr, "remaining_seconds"))
    result.StartedAt = ExtractJSONValue(jsonStr, "started_at")
    result.Ticker = ExtractJSONValue(jsonStr, "ticker")
End Sub

' ParseCandidateCheckJSON parses candidate check JSON into TFCandidateCheck
Public Sub ParseCandidateCheckJSON(ByVal jsonStr As String, ByRef result As TFCandidateCheck)
    result.Found = (LCase(ExtractJSONValue(jsonStr, "found")) = "true")
    result.Ticker = ExtractJSONValue(jsonStr, "ticker")
    result.DateStr = ExtractJSONValue(jsonStr, "date")
    result.Message = ExtractJSONValue(jsonStr, "message")
End Sub

' ParseCooldownCheckJSON parses cooldown check JSON into TFCooldownCheck
Public Sub ParseCooldownCheckJSON(ByVal jsonStr As String, ByRef result As TFCooldownCheck)
    result.OnCooldown = (LCase(ExtractJSONValue(jsonStr, "on_cooldown")) = "true")
    result.Bucket = ExtractJSONValue(jsonStr, "bucket")
    result.ClearsAt = ExtractJSONValue(jsonStr, "clears_at")
    result.Message = ExtractJSONValue(jsonStr, "message")
End Sub

' ParseSettingsJSON parses settings JSON into TFSettings
Public Sub ParseSettingsJSON(ByVal jsonStr As String, ByRef result As TFSettings)
    result.Equity_E = CDbl(ExtractJSONValue(jsonStr, "Equity_E"))
    result.RiskPct_r = CDbl(ExtractJSONValue(jsonStr, "RiskPct_r"))
    result.HeatCap_H_pct = CDbl(ExtractJSONValue(jsonStr, "HeatCap_H_pct"))
    result.BucketHeatCap_pct = CDbl(ExtractJSONValue(jsonStr, "BucketHeatCap_pct"))
    result.StopMultiple_K = CLng(ExtractJSONValue(jsonStr, "StopMultiple_K"))
End Sub

' ParseSaveDecisionJSON parses save-decision JSON into TFSaveDecisionResult
Public Sub ParseSaveDecisionJSON(ByVal jsonStr As String, ByRef result As TFSaveDecisionResult)
    On Error GoTo ErrorHandler
    result.Success = False
    result.Saved = (LCase(ExtractJSONValue(jsonStr, "accepted")) = "true")

    Dim decisionIDStr As String
    decisionIDStr = ExtractJSONValue(jsonStr, "decision_id")
    If IsNumeric(decisionIDStr) Then
        result.DecisionID = CLng(decisionIDStr)
    Else
        result.DecisionID = 0
    End If

    result.Timestamp = ExtractJSONValue(jsonStr, "timestamp")
    result.RejectionReason = ExtractJSONValue(jsonStr, "reason")
    result.GatesFailed = ExtractJSONValue(jsonStr, "gates_failed")
    result.Success = True
    Exit Sub
ErrorHandler:
    result.Success = False
End Sub

'-----------------------------------------------------------------------------
' ERROR HANDLING
'-----------------------------------------------------------------------------

' CreateError creates a TFEngineError from command result
Public Sub CreateError(ByVal corrID As String, ByVal errorMsg As String, ByVal exitCode As Long, ByRef result As TFEngineError)
    result.HasError = True
    result.ErrorMessage = errorMsg
    result.CorrelationID = corrID
    result.ExitCode = exitCode
End Sub

' FormatErrorMessage formats an error message for display in Excel
Public Function FormatErrorMessage(ByRef err As TFEngineError) As String
    If Not err.HasError Then
        FormatErrorMessage = ""
        Exit Function
    End If

    Dim msg As String
    msg = "ERROR: " & err.ErrorMessage & vbCrLf & _
          "Correlation ID: " & err.CorrelationID & vbCrLf & _
          "Exit Code: " & err.ExitCode

    FormatErrorMessage = msg
End Function

'-----------------------------------------------------------------------------
' STRING UTILITIES
'-----------------------------------------------------------------------------

' SafeString ensures a string value is not null or empty
Public Function SafeString(ByVal value As Variant, Optional ByVal defaultValue As String = "") As String
    If IsNull(value) Then
        SafeString = defaultValue
    ElseIf VarType(value) = vbString Then
        SafeString = CStr(value)
    Else
        SafeString = defaultValue
    End If
End Function

' SafeDouble converts a value to Double with default fallback
Public Function SafeDouble(ByVal value As Variant, Optional ByVal defaultValue As Double = 0#) As Double
    On Error Resume Next
    If IsNumeric(value) Then
        SafeDouble = CDbl(value)
    Else
        SafeDouble = defaultValue
    End If
End Function

' SafeLong converts a value to Long with default fallback
Public Function SafeLong(ByVal value As Variant, Optional ByVal defaultValue As Long = 0) As Long
    On Error Resume Next
    If IsNumeric(value) Then
        SafeLong = CLng(value)
    Else
        SafeLong = defaultValue
    End If
End Function

'-----------------------------------------------------------------------------
' VALIDATION HELPERS
'-----------------------------------------------------------------------------

' ValidateTicker checks if a ticker symbol is valid (basic check)
Public Function ValidateTicker(ByVal ticker As String) As Boolean
    If Len(Trim(ticker)) = 0 Then
        ValidateTicker = False
        Exit Function
    End If

    ' Ticker should be alphanumeric plus hyphen/period (e.g., BRK-B, BRK.B)
    Dim i As Long
    For i = 1 To Len(ticker)
        Dim ch As String
        ch = Mid(ticker, i, 1)
        If Not ((ch >= "A" And ch <= "Z") Or _
                (ch >= "a" And ch <= "z") Or _
                (ch >= "0" And ch <= "9") Or _
                ch = "-" Or ch = ".") Then
            ValidateTicker = False
            Exit Function
        End If
    Next i

    ValidateTicker = True
End Function

' ValidatePositiveNumber checks if a value is a positive number
Public Function ValidatePositiveNumber(ByVal value As Variant) As Boolean
    On Error Resume Next
    Dim num As Double
    num = CDbl(value)
    ValidatePositiveNumber = (num > 0)
End Function

'-----------------------------------------------------------------------------
' FORMATTING HELPERS
'-----------------------------------------------------------------------------

' FormatCurrency formats a double as currency string
Public Function FormatCurrency(ByVal value As Double) As String
    FormatCurrency = Format(value, "$#,##0.00")
End Function

' FormatPercent formats a decimal as percentage string
Public Function FormatPercent(ByVal value As Double) As String
    FormatPercent = Format(value * 100, "0.00") & "%"
End Function

' FormatTimestamp formats an ISO8601 timestamp for display
Public Function FormatTimestamp(ByVal isoTimestamp As String) As String
    On Error Resume Next

    ' Simple display format - just show date and time part
    ' ISO format: 2025-10-27T14:30:52-05:00
    ' Extract: YYYY-MM-DD HH:MM:SS

    If Len(isoTimestamp) < 19 Then
        FormatTimestamp = isoTimestamp
        Exit Function
    End If

    Dim datePart As String
    Dim timePart As String

    datePart = Left(isoTimestamp, 10)  ' YYYY-MM-DD
    timePart = Mid(isoTimestamp, 12, 8)  ' HH:MM:SS

    FormatTimestamp = datePart & " " & timePart
End Function

'=============================================================================
' NOTES ON IMPLEMENTATION
'=============================================================================
'
' JSON PARSING APPROACH:
' This module uses simple string parsing instead of external JSON libraries.
' Rationale:
' 1. No external dependencies (easier deployment to Windows)
' 2. JSON schemas are validated in Go engine (defensive programming there)
' 3. VBA is a thin bridge - extraction is sufficient
' 4. Performance is not critical (small JSON payloads)
'
' LIMITATIONS:
' - Does not handle nested objects/arrays well
' - Assumes well-formed JSON from engine
' - No JSON schema validation (done in Go)
'
' ACCEPTABLE BECAUSE:
' - Engine outputs are validated and tested (M17-M18)
' - VBA only extracts known keys from known schemas
' - Errors fail gracefully (return empty/default values)
'
' ERROR HANDLING STRATEGY:
' - Use "On Error Resume Next" for non-critical operations (logging, parsing)
' - Return sensible defaults for missing values
' - Log all errors with correlation IDs for debugging
'
' LOGGING PATTERN:
' All engine calls should log:
' 1. Before call: "Calling tf-engine [command] with corr_id [ID]"
' 2. After success: "Received JSON response [length] bytes"
' 3. After error: "Engine error: [message] (exit code [code])"
'
'=============================================================================

'-----------------------------------------------------------------------------
' UI NAVIGATION FUNCTIONS (M22 - Automated UI Generation)
'-----------------------------------------------------------------------------

' RefreshDashboard - Updates dashboard with current portfolio data
Public Sub RefreshDashboard()
    On Error GoTo ErrorHandler

    Dim ws As Worksheet
    Set ws = ThisWorkbook.Worksheets("Dashboard")

    ' Generate correlation ID for this operation
    Dim corrID As String
    corrID = GenerateCorrelationID()

    LogMessage corrID, "INFO", "Refreshing dashboard"

    ' TODO: Query portfolio state and update cells
    ' For now, just show refresh timestamp
    ws.Range("A8").Value = "Last Refresh:"
    ws.Range("B8").Value = Format(Now, "YYYY-MM-DD HH:NN:SS")

    MsgBox "Dashboard refreshed!" & vbCrLf & "Correlation ID: " & corrID, vbInformation, "Dashboard"

    Exit Sub

ErrorHandler:
    MsgBox "Error refreshing dashboard: " & Err.Description & vbCrLf & _
           "Correlation ID: " & corrID, vbCritical, "Error"
    LogMessage corrID, "ERROR", "Dashboard refresh failed: " & Err.Description
End Sub

' Navigation helper functions - Navigate between worksheets
Public Sub GotoPositionSizing()
    On Error Resume Next
    ThisWorkbook.Worksheets("Position Sizing").Activate
End Sub

Public Sub GotoChecklist()
    On Error Resume Next
    ThisWorkbook.Worksheets("Checklist").Activate
End Sub

Public Sub GotoHeatCheck()
    On Error Resume Next
    ThisWorkbook.Worksheets("Heat Check").Activate
End Sub

Public Sub GotoTradeEntry()
    On Error Resume Next
    ThisWorkbook.Worksheets("Trade Entry").Activate
End Sub

Public Sub GotoDashboard()
    On Error Resume Next
    ThisWorkbook.Worksheets("Dashboard").Activate
End Sub

'=============================================================================
