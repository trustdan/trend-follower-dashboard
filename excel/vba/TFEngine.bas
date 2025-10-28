Attribute VB_Name = "TFEngine"
'=============================================================================
' TFEngine.bas - Trading Engine Communication Layer
'=============================================================================
' Purpose: Bridge between Excel and tf-engine.exe (Go backend)
'
' Architecture:
' - Excel (UI) → VBA (this module) → CLI (tf-engine.exe) → Go Engine
' - Commands via WScript.Shell.Exec (synchronous execution)
' - JSON output via stdout, errors via stderr
' - All business logic in Go - VBA is a thin bridge
'
' Key Functions:
' - ExecuteCommand() - Low-level command execution
' - Engine_Size() - Position sizing
' - Engine_Checklist() - Checklist evaluation
' - Engine_Heat() - Heat management
' - Engine_SaveDecision() - Save trade decision (5 hard gates)
' - Engine_ImportCandidates() - Import candidate tickers
' - Engine_GetSettings() - Retrieve settings
'
' Philosophy: Keep VBA thin - just shell exec + JSON parsing
'
' Created: 2025-10-27 (M19 - VBA Implementation)
'=============================================================================

Option Explicit

'-----------------------------------------------------------------------------
' CONFIGURATION
'-----------------------------------------------------------------------------

' Default path to tf-engine.exe (can be overridden from Setup sheet)
Private Const DEFAULT_ENGINE_PATH As String = "tf-engine.exe"
Private Const DEFAULT_DB_PATH As String = "trading.db"
Private Const COMMAND_TIMEOUT_SECONDS As Long = 30

'-----------------------------------------------------------------------------
' CORE COMMAND EXECUTION
'-----------------------------------------------------------------------------

' ExecuteCommand runs tf-engine.exe with specified arguments
' This is the core function - all other Engine_XXX functions use this
'
' Arguments:
'   command - Command string (e.g., "size --entry 180 --atr 1.5 --method stock")
'   corrID - Correlation ID for logging (auto-generated if empty)
'
' Returns: TFCommandResult with success status, JSON output, and errors
'
Public Function ExecuteCommand(ByVal command As String, Optional ByVal corrID As String = "") As TFCommandResult
    Dim result As TFCommandResult
    Dim shell As Object
    Dim exec As Object
    Dim enginePath As String
    Dim dbPath As String
    Dim fullCommand As String
    Dim stdOut As String
    Dim stdErr As String
    Dim startTime As Double
    Dim timeout As Double

    ' Generate correlation ID if not provided
    If corrID = "" Then corrID = TFHelpers.GenerateCorrelationID()
    result.CorrelationID = corrID

    ' Get engine path (from Setup sheet or default)
    enginePath = GetEnginePath()
    dbPath = GetDatabasePath()

    ' Log paths for debugging
    TFHelpers.LogMessage corrID, "DEBUG", "Engine path: " & enginePath
    TFHelpers.LogMessage corrID, "DEBUG", "Database path: " & dbPath

    ' Build full command with global flags
    ' Format: tf-engine.exe --db trading.db --corr-id <id> --format json <command>
    fullCommand = """" & enginePath & """ --db """ & dbPath & """ --corr-id " & corrID & " --format json " & command

    ' Log command execution
    TFHelpers.LogMessage corrID, "INFO", "Executing: " & fullCommand

    On Error GoTo ErrorHandler

    ' Create shell and execute command
    Set shell = CreateObject("WScript.Shell")
    Set exec = shell.exec(fullCommand)

    ' Wait for command to complete (with timeout)
    startTime = Timer
    timeout = COMMAND_TIMEOUT_SECONDS

    Do While exec.Status = 0  ' 0 = Still Running
        DoEvents  ' Allow Excel to process events

        ' Check for timeout
        If Timer - startTime > timeout Then
            TFHelpers.LogMessage corrID, "ERROR", "Command timed out after " & timeout & " seconds"
            result.Success = False
            result.ErrorOutput = "Command timed out after " & timeout & " seconds"
            result.ExitCode = -1
            GoTo Cleanup
        End If

        ' Brief pause to avoid CPU spinning (100ms)
        Application.Wait Now + (0.1 / 86400)  ' Add 0.1 seconds as fraction of day
    Loop

    ' Read stdout and stderr
    If Not exec.StdOut.AtEndOfStream Then
        stdOut = exec.StdOut.ReadAll()
    End If

    If Not exec.StdErr.AtEndOfStream Then
        stdErr = exec.StdErr.ReadAll()
    End If

    ' Get exit code
    result.ExitCode = exec.ExitCode

    ' Determine success
    If result.ExitCode = 0 And stdErr = "" Then
        result.Success = True
        result.JsonOutput = stdOut
        result.ErrorOutput = ""
        TFHelpers.LogMessage corrID, "INFO", "Command succeeded (" & Len(stdOut) & " bytes JSON)"
    Else
        result.Success = False
        result.JsonOutput = stdOut  ' May contain partial output
        result.ErrorOutput = stdErr
        TFHelpers.LogMessage corrID, "ERROR", "Command failed (exit " & result.ExitCode & "): " & stdErr
    End If

Cleanup:
    Set exec = Nothing
    Set shell = Nothing
    ExecuteCommand = result
    Exit Function

ErrorHandler:
    TFHelpers.LogMessage corrID, "ERROR", "Exception during command execution: " & Err.Description
    result.Success = False
    result.ErrorOutput = "VBA Error: " & Err.Description & " (Error " & Err.Number & ")"
    result.ExitCode = -999
    Resume Cleanup
End Function

'-----------------------------------------------------------------------------
' CONFIGURATION HELPERS
'-----------------------------------------------------------------------------

' GetEnginePath retrieves the path to tf-engine.exe
' Checks: Setup sheet setting, then default
Private Function GetEnginePath() As String
    On Error Resume Next

    Dim ws As Worksheet
    Dim pathRange As Range
    Dim configPath As String

    ' Try to get from Setup sheet
    Set ws = ThisWorkbook.Worksheets("Setup")
    If Not ws Is Nothing Then
        Set pathRange = ws.Range("EnginePathSetting")  ' Named range
        If Not pathRange Is Nothing Then
            If pathRange.Value <> "" Then
                configPath = pathRange.Value

                ' If relative path (starts with .), make it absolute
                If Left(configPath, 2) = ".\" Or Left(configPath, 2) = "./" Then
                    GetEnginePath = ThisWorkbook.Path & Mid(configPath, 2)
                    Exit Function
                End If

                ' If absolute path or just filename, use as-is
                GetEnginePath = configPath
                Exit Function
            End If
        End If
    End If

    ' Default to current directory
    GetEnginePath = ThisWorkbook.Path & "\" & DEFAULT_ENGINE_PATH
End Function

' GetDatabasePath retrieves the path to trading.db
' Checks: Setup sheet setting, then default
Private Function GetDatabasePath() As String
    On Error Resume Next

    Dim ws As Worksheet
    Dim pathRange As Range
    Dim configPath As String

    ' Try to get from Setup sheet
    Set ws = ThisWorkbook.Worksheets("Setup")
    If Not ws Is Nothing Then
        Set pathRange = ws.Range("DatabasePathSetting")  ' Named range
        If Not pathRange Is Nothing Then
            If pathRange.Value <> "" Then
                configPath = pathRange.Value

                ' If relative path (starts with .), make it absolute
                If Left(configPath, 2) = ".\" Or Left(configPath, 2) = "./" Then
                    GetDatabasePath = ThisWorkbook.Path & Mid(configPath, 2)
                    Exit Function
                End If

                ' If absolute path or just filename, use as-is
                GetDatabasePath = configPath
                Exit Function
            End If
        End If
    End If

    ' Default to current directory
    GetDatabasePath = ThisWorkbook.Path & "\" & DEFAULT_DB_PATH
End Function

'-----------------------------------------------------------------------------
' POSITION SIZING
'-----------------------------------------------------------------------------

' Engine_Size calculates position size
'
' Parameters:
'   entry - Entry price
'   atr - Average True Range (ATR)
'   method - "stock", "opt-delta-atr", "opt-maxloss"
'   Optional parameters (uses settings from DB if not provided):
'     equity - Account equity
'     riskPct - Risk per trade (as decimal, e.g., 0.0075)
'     k - Stop multiple
'     delta - Option delta (required for opt-delta-atr)
'     maxLoss - Max loss per contract (required for opt-maxloss)
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
Public Function Engine_Size( _
    ByVal entry As Double, _
    ByVal atr As Double, _
    ByVal method As String, _
    Optional ByVal equity As Double = 0, _
    Optional ByVal riskPct As Double = 0, _
    Optional ByVal k As Long = 0, _
    Optional ByVal delta As Double = 0, _
    Optional ByVal maxLoss As Double = 0, _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String

    ' Build command
    command = "size --entry " & entry & " --atr " & atr & " --method " & method

    ' Add optional parameters
    If equity > 0 Then command = command & " --equity " & equity
    If riskPct > 0 Then command = command & " --risk " & riskPct
    If k > 0 Then command = command & " --k " & k
    If delta > 0 Then command = command & " --delta " & delta
    If maxLoss > 0 Then command = command & " --maxloss " & maxLoss

    ' Execute
    Engine_Size = ExecuteCommand(command, corrID)
End Function

'-----------------------------------------------------------------------------
' CHECKLIST EVALUATION
'-----------------------------------------------------------------------------

' Engine_Checklist evaluates 6-item checklist and returns banner
'
' Parameters:
'   ticker - Stock ticker
'   checks - Collection of 6 boolean values (checklist items)
'            Order: Higher high, Wider range, Close off low,
'                   Liquidity OK, Not overbought, Bucket OK
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
Public Function Engine_Checklist( _
    ByVal ticker As String, _
    ByVal checks As Collection, _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String
    Dim i As Long
    Dim checkStr As String

    ' Validate checks collection
    If checks.Count <> 6 Then
        Dim errResult As TFCommandResult
        errResult.Success = False
        errResult.ErrorOutput = "Checklist must have exactly 6 items (got " & checks.Count & ")"
        errResult.ExitCode = -1
        Engine_Checklist = errResult
        Exit Function
    End If

    ' Build checks string: --checks true,true,false,true,true,false
    checkStr = ""
    For i = 1 To checks.Count
        If i > 1 Then checkStr = checkStr & ","
        If checks(i) Then
            checkStr = checkStr & "true"
        Else
            checkStr = checkStr & "false"
        End If
    Next i

    ' Build command
    command = "checklist --ticker " & ticker & " --checks " & checkStr

    ' Execute
    Engine_Checklist = ExecuteCommand(command, corrID)
End Function

'-----------------------------------------------------------------------------
' HEAT MANAGEMENT
'-----------------------------------------------------------------------------

' Engine_Heat checks portfolio and bucket heat with a proposed new trade
'
' Parameters:
'   addR - Risk dollars for proposed trade
'   bucket - Sector bucket (e.g., "Tech/Comm")
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
Public Function Engine_Heat( _
    ByVal addR As Double, _
    ByVal bucket As String, _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String

    ' Build command
    command = "heat --add-r " & addR & " --bucket """ & bucket & """"

    ' Execute
    Engine_Heat = ExecuteCommand(command, corrID)
End Function

'-----------------------------------------------------------------------------
' IMPULSE TIMER
'-----------------------------------------------------------------------------

' Engine_CheckTimer checks the status of the 2-minute impulse brake
'
' Parameters:
'   ticker - Stock ticker to check
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
Public Function Engine_CheckTimer( _
    ByVal ticker As String, _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String
    command = "timer --ticker " & ticker
    Engine_CheckTimer = ExecuteCommand(command, corrID)
End Function

'-----------------------------------------------------------------------------
' CANDIDATE MANAGEMENT
'-----------------------------------------------------------------------------

' Engine_ImportCandidates imports a list of candidate tickers
'
' Parameters:
'   tickers - Comma-separated list of tickers (e.g., "AAPL,MSFT,NVDA")
'   preset - Preset name (e.g., "TF_BREAKOUT_LONG")
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
Public Function Engine_ImportCandidates( _
    ByVal tickers As String, _
    ByVal preset As String, _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String
    command = "import-candidates --tickers """ & tickers & """ --preset " & preset
    Engine_ImportCandidates = ExecuteCommand(command, corrID)
End Function

' Engine_ListCandidates retrieves today's candidate list
'
' Parameters:
'   dateStr - Date string (YYYY-MM-DD format, defaults to today if empty)
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
Public Function Engine_ListCandidates( _
    Optional ByVal dateStr As String = "", _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String
    command = "list-candidates"

    If dateStr <> "" Then
        command = command & " --date " & dateStr
    End If

    Engine_ListCandidates = ExecuteCommand(command, corrID)
End Function

' Engine_CheckCandidate checks if a ticker is in today's candidates
'
' Parameters:
'   ticker - Stock ticker to check
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
Public Function Engine_CheckCandidate( _
    ByVal ticker As String, _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String
    command = "check-candidate --ticker " & ticker
    Engine_CheckCandidate = ExecuteCommand(command, corrID)
End Function

'-----------------------------------------------------------------------------
' COOLDOWN MANAGEMENT
'-----------------------------------------------------------------------------

' Engine_CheckCooldown checks if a bucket is on cooldown
'
' Parameters:
'   bucket - Sector bucket to check
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
Public Function Engine_CheckCooldown( _
    ByVal bucket As String, _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String
    command = "check-cooldown --bucket """ & bucket & """"
    Engine_CheckCooldown = ExecuteCommand(command, corrID)
End Function

' Engine_ListCooldowns retrieves all cooldowns
'
' Parameters:
'   activeOnly - If True, only return active cooldowns
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
Public Function Engine_ListCooldowns( _
    Optional ByVal activeOnly As Boolean = False, _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String
    command = "list-cooldowns"

    If activeOnly Then
        command = command & " --active-only"
    End If

    Engine_ListCooldowns = ExecuteCommand(command, corrID)
End Function

'-----------------------------------------------------------------------------
' SETTINGS MANAGEMENT
'-----------------------------------------------------------------------------

' Engine_GetSettings retrieves all settings from database
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
Public Function Engine_GetSettings( _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String
    command = "get-settings"
    Engine_GetSettings = ExecuteCommand(command, corrID)
End Function

' Engine_SetSetting updates a single setting
'
' Parameters:
'   key - Setting key (e.g., "Equity_E", "RiskPct_r")
'   value - New value (as string)
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
Public Function Engine_SetSetting( _
    ByVal key As String, _
    ByVal value As String, _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String
    command = "set-setting --key " & key & " --value " & value
    Engine_SetSetting = ExecuteCommand(command, corrID)
End Function

'-----------------------------------------------------------------------------
' SAVE DECISION (5 HARD GATES)
'-----------------------------------------------------------------------------

' Engine_SaveDecision saves a trade decision (enforces 5 hard gates)
'
' The 5 hard gates:
'   1. Banner must be GREEN
'   2. Ticker must be in today's candidates
'   3. 2-minute impulse brake must have elapsed
'   4. Bucket must not be on cooldown
'   5. Heat caps must not be exceeded
'
' Parameters:
'   ticker - Stock ticker
'   entry - Entry price
'   atr - Average True Range
'   k - Stop multiple
'   method - Sizing method
'   riskDollars - Risk dollars (R)
'   shares - Number of shares (0 for options)
'   contracts - Number of contracts (0 for stocks)
'   banner - Checklist banner ("GREEN", "YELLOW", "RED")
'   bucket - Sector bucket
'   preset - Preset name
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
Public Function Engine_SaveDecision( _
    ByVal ticker As String, _
    ByVal entry As Double, _
    ByVal atr As Double, _
    ByVal k As Long, _
    ByVal method As String, _
    ByVal riskDollars As Double, _
    ByVal shares As Long, _
    ByVal contracts As Long, _
    ByVal banner As String, _
    ByVal bucket As String, _
    ByVal preset As String, _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String

    ' Build command with all required parameters
    command = "save-decision" & _
              " --ticker " & ticker & _
              " --entry " & entry & _
              " --atr " & atr & _
              " --k " & k & _
              " --method " & method & _
              " --risk " & riskDollars & _
              " --shares " & shares & _
              " --contracts " & contracts & _
              " --banner " & banner & _
              " --bucket """ & bucket & """" & _
              " --preset " & preset

    ' Execute (engine enforces all 5 gates)
    Engine_SaveDecision = ExecuteCommand(command, corrID)
End Function

'-----------------------------------------------------------------------------
' FINVIZ SCRAPER
'-----------------------------------------------------------------------------

' Engine_ScrapeFinviz scrapes FINVIZ with a query string
'
' Parameters:
'   query - FINVIZ query string (URL parameters)
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
' NOTE: Manual paste fallback should remain available in Excel UI
'
Public Function Engine_ScrapeFinviz( _
    ByVal query As String, _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String
    command = "scrape-finviz --query """ & query & """"
    Engine_ScrapeFinviz = ExecuteCommand(command, corrID)
End Function

'-----------------------------------------------------------------------------
' DATABASE INITIALIZATION
'-----------------------------------------------------------------------------

' Engine_Init initializes the database (creates schema if not exists)
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
Public Function Engine_Init( _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String
    command = "init"
    Engine_Init = ExecuteCommand(command, corrID)
End Function

'-----------------------------------------------------------------------------
' POSITION MANAGEMENT
'-----------------------------------------------------------------------------

' Engine_OpenPosition opens a new position after decision is saved
'
' Parameters:
'   ticker - Stock ticker
'   bucket - Sector bucket
'   units - Number of shares/contracts
'   riskDollars - Risk dollars (R) for this position
'   entryPrice - Entry price
'   initialStop - Initial stop loss
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
Public Function Engine_OpenPosition( _
    ByVal ticker As String, _
    ByVal bucket As String, _
    ByVal units As Long, _
    ByVal riskDollars As Double, _
    ByVal entryPrice As Double, _
    ByVal initialStop As Double, _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String
    command = "open-position" & _
              " --ticker " & ticker & _
              " --bucket """ & bucket & """" & _
              " --units " & units & _
              " --risk " & riskDollars & _
              " --entry " & entryPrice & _
              " --stop " & initialStop

    Engine_OpenPosition = ExecuteCommand(command, corrID)
End Function

' Engine_ListPositions retrieves all positions
'
' Parameters:
'   openOnly - If True, only return open positions
'
' Returns: TFCommandResult (check .Success before parsing JSON)
'
Public Function Engine_ListPositions( _
    Optional ByVal openOnly As Boolean = True, _
    Optional ByVal corrID As String = "" _
) As TFCommandResult

    Dim command As String
    command = "list-positions"

    If openOnly Then
        command = command & " --open-only"
    End If

    Engine_ListPositions = ExecuteCommand(command, corrID)
End Function

'=============================================================================
' USAGE PATTERNS
'=============================================================================
'
' STANDARD CALL PATTERN:
'
' Sub Example_PositionSizing()
'     Dim cmdResult As TFCommandResult
'     Dim sizeResult As TFSizingResult
'     Dim corrID As String
'
'     ' Generate correlation ID
'     corrID = TFHelpers.GenerateCorrelationID()
'
'     ' Call engine
'     cmdResult = TFEngine.Engine_Size(180, 1.5, "stock", corrID:=corrID)
'
'     ' Check success
'     If cmdResult.Success Then
'         ' Parse JSON
'         sizeResult = TFHelpers.ParseSizingJSON(cmdResult.JsonOutput)
'
'         ' Use result
'         Range("SharesResult").Value = sizeResult.Shares
'         Range("InitialStopResult").Value = sizeResult.InitialStop
'         Range("StatusCell").Value = "✅ Success (corr_id: " & corrID & ")"
'     Else
'         ' Handle error
'         Range("StatusCell").Value = "❌ " & cmdResult.ErrorOutput
'         MsgBox "Error: " & cmdResult.ErrorOutput & vbCrLf & _
'                "Correlation ID: " & corrID, vbCritical
'     End If
' End Sub
'
' ERROR DISPLAY PATTERN:
'
' - Always show correlation ID in status cells
' - Always log errors with TFHelpers.LogMessage
' - Never swallow errors silently
' - Provide actionable error messages to user
'
' TESTING PATTERN:
'
' 1. Unit test each Engine_XXX function with known inputs
' 2. Verify JSON output matches expected schema
' 3. Verify error handling (invalid inputs, missing exe, etc.)
' 4. Verify correlation IDs appear in both VBA and Go logs
'
'=============================================================================

'-----------------------------------------------------------------------------
' UI BUTTON HANDLERS (M22 - Automated UI Generation)
'-----------------------------------------------------------------------------

' CalculatePositionSize - Position Sizing worksheet button handler
Public Sub CalculatePositionSize()
    On Error GoTo ErrorHandler

    Dim ws As Worksheet
    Set ws = ThisWorkbook.Worksheets("Position Sizing")

    ' Generate correlation ID
    Dim corrID As String
    corrID = TFHelpers.GenerateCorrelationID()

    ' Read inputs
    Dim ticker As String, entryPrice As Double, atr As Double
    Dim kMultiple As Double, method As String
    Dim equityOverride As Variant, riskPctOverride As Variant
    Dim delta As Variant, maxLoss As Variant

    ticker = ws.Range("B4").Value
    entryPrice = ws.Range("B5").Value
    atr = ws.Range("B6").Value
    kMultiple = ws.Range("B7").Value
    method = ws.Range("B8").Value

    ' Optional fields
    equityOverride = ws.Range("B9").Value
    riskPctOverride = ws.Range("B10").Value
    delta = ws.Range("B11").Value
    maxLoss = ws.Range("B12").Value

    ' Validate required inputs
    If ticker = "" Or entryPrice <= 0 Or atr <= 0 Or kMultiple <= 0 Or method = "" Then
        ws.Range("B22").Value = "❌ Missing required inputs"
        MsgBox "Please fill in all required fields: Ticker, Entry Price, ATR, K Multiple, Method", vbExclamation, "Missing Inputs"
        Exit Sub
    End If

    ' Build command
    Dim cmd As String
    cmd = "size --entry " & entryPrice & " --atr " & atr & " --k " & kMultiple & " --method " & method

    ' Add optional parameters
    If Not IsEmpty(equityOverride) And equityOverride <> "" Then
        cmd = cmd & " --equity " & equityOverride
    End If
    If Not IsEmpty(riskPctOverride) And riskPctOverride <> "" Then
        cmd = cmd & " --risk-pct " & (riskPctOverride * 100)  ' Convert to percentage
    End If
    If Not IsEmpty(delta) And delta <> "" Then
        cmd = cmd & " --delta " & delta
    End If
    If Not IsEmpty(maxLoss) And maxLoss <> "" Then
        cmd = cmd & " --max-loss " & maxLoss
    End If

    ' Execute command
    Dim result As TFCommandResult
    result = ExecuteCommand(cmd, corrID)

    If result.Success Then
        ' Parse results
        Dim sizeResult As TFSizingResult
        sizeResult = TFHelpers.ParseSizingJSON(result.JsonOutput)

        ' Display results
        ws.Range("B16").Value = sizeResult.RiskDollars
        ws.Range("B17").Value = sizeResult.StopDistance
        ws.Range("B18").Value = sizeResult.InitialStop
        ws.Range("B19").Value = sizeResult.Shares
        ws.Range("B20").Value = sizeResult.Contracts
        ws.Range("B21").Value = sizeResult.ActualRisk
        ws.Range("B22").Value = "✅ Success (" & corrID & ")"
        ws.Range("B22").Font.Color = RGB(0, 128, 0)
    Else
        ws.Range("B22").Value = "❌ " & result.ErrorOutput
        ws.Range("B22").Font.Color = RGB(255, 0, 0)
        MsgBox "Position sizing failed:" & vbCrLf & result.ErrorOutput & vbCrLf & vbCrLf & _
               "Correlation ID: " & corrID, vbCritical, "Error"
    End If

    Exit Sub

ErrorHandler:
    ws.Range("B22").Value = "❌ Error: " & Err.Description
    ws.Range("B22").Font.Color = RGB(255, 0, 0)
    MsgBox "Error in position sizing: " & Err.Description & vbCrLf & _
           "Correlation ID: " & corrID, vbCritical, "Error"
    TFHelpers.LogMessage corrID, "ERROR", "CalculatePositionSize failed: " & Err.Description
End Sub

' ClearPositionSizing - Clear Position Sizing worksheet inputs/results
Public Sub ClearPositionSizing()
    On Error Resume Next
    Dim ws As Worksheet
    Set ws = ThisWorkbook.Worksheets("Position Sizing")

    ' Clear inputs
    ws.Range("B4:B12").ClearContents

    ' Clear results
    ws.Range("B16:B22").ClearContents
    ws.Range("B22").Font.Color = RGB(0, 0, 0)
End Sub

' EvaluateChecklist - Checklist worksheet button handler
Public Sub EvaluateChecklist()
    On Error GoTo ErrorHandler

    Dim ws As Worksheet
    Set ws = ThisWorkbook.Worksheets("Checklist")

    ' Generate correlation ID
    Dim corrID As String
    corrID = TFHelpers.GenerateCorrelationID()

    ' Read ticker
    Dim ticker As String
    ticker = ws.Range("B3").Value

    If ticker = "" Then
        ws.Range("B23").Value = "❌ Ticker required"
        MsgBox "Please enter a ticker symbol", vbExclamation, "Missing Ticker"
        Exit Sub
    End If

    ' Read checkbox states
    Dim fromPreset As Boolean, trendPass As Boolean, liquidityPass As Boolean
    Dim tvConfirm As Boolean, earningsOk As Boolean, journalOk As Boolean

    fromPreset = ws.OLEObjects("chk_from_preset").Object.Value
    trendPass = ws.OLEObjects("chk_trend_pass").Object.Value
    liquidityPass = ws.OLEObjects("chk_liquidity_pass").Object.Value
    tvConfirm = ws.OLEObjects("chk_tv_confirm").Object.Value
    earningsOk = ws.OLEObjects("chk_earnings_ok").Object.Value
    journalOk = ws.OLEObjects("chk_journal_ok").Object.Value

    ' Build command
    Dim cmd As String
    cmd = "checklist --ticker " & ticker

    If fromPreset Then cmd = cmd & " --from-preset"
    If trendPass Then cmd = cmd & " --trend-pass"
    If liquidityPass Then cmd = cmd & " --liquidity-pass"
    If tvConfirm Then cmd = cmd & " --tv-confirm"
    If earningsOk Then cmd = cmd & " --earnings-ok"
    If journalOk Then cmd = cmd & " --journal-ok"

    ' Execute command
    Dim result As TFCommandResult
    result = ExecuteCommand(cmd, corrID)

    If result.Success Then
        ' Parse results
        Dim checkResult As TFChecklistResult
        checkResult = TFHelpers.ParseChecklistJSON(result.JsonOutput)

        ' Display results
        ws.Range("B16").Value = checkResult.Banner
        ws.Range("B17").Value = checkResult.MissingCount

        ' Color-code banner
        If checkResult.Banner = "GREEN" Then
            ws.Range("B16").Interior.Color = RGB(0, 255, 0)
            ws.Range("B16").Font.Color = RGB(0, 0, 0)
        ElseIf checkResult.Banner = "YELLOW" Then
            ws.Range("B16").Interior.Color = RGB(255, 255, 0)
            ws.Range("B16").Font.Color = RGB(0, 0, 0)
        Else ' RED
            ws.Range("B16").Interior.Color = RGB(255, 0, 0)
            ws.Range("B16").Font.Color = RGB(255, 255, 255)
        End If

        ' Show missing items
        Dim missing As String
        missing = Join(checkResult.MissingItems, vbCrLf)
        ws.Range("B18").Value = missing

        ws.Range("B21").Value = checkResult.AllowSave
        ws.Range("B22").Value = checkResult.EvalTime
        ws.Range("B23").Value = "✅ Success (" & corrID & ")"
        ws.Range("B23").Font.Color = RGB(0, 128, 0)
    Else
        ws.Range("B23").Value = "❌ " & result.ErrorOutput
        ws.Range("B23").Font.Color = RGB(255, 0, 0)
        MsgBox "Checklist evaluation failed:" & vbCrLf & result.ErrorOutput & vbCrLf & vbCrLf & _
               "Correlation ID: " & corrID, vbCritical, "Error"
    End If

    Exit Sub

ErrorHandler:
    ws.Range("B23").Value = "❌ Error: " & Err.Description
    ws.Range("B23").Font.Color = RGB(255, 0, 0)
    MsgBox "Error in checklist evaluation: " & Err.Description & vbCrLf & _
           "Correlation ID: " & corrID, vbCritical, "Error"
    TFHelpers.LogMessage corrID, "ERROR", "EvaluateChecklist failed: " & Err.Description
End Sub

' ClearChecklist - Clear Checklist worksheet
Public Sub ClearChecklist()
    On Error Resume Next
    Dim ws As Worksheet
    Set ws = ThisWorkbook.Worksheets("Checklist")

    ' Clear ticker
    ws.Range("B3").ClearContents

    ' Uncheck all checkboxes
    ws.OLEObjects("chk_from_preset").Object.Value = False
    ws.OLEObjects("chk_trend_pass").Object.Value = False
    ws.OLEObjects("chk_liquidity_pass").Object.Value = False
    ws.OLEObjects("chk_tv_confirm").Object.Value = False
    ws.OLEObjects("chk_earnings_ok").Object.Value = False
    ws.OLEObjects("chk_journal_ok").Object.Value = False

    ' Clear results
    ws.Range("B16:B23").ClearContents
    ws.Range("B16").Interior.ColorIndex = xlNone
    ws.Range("B16").Font.Color = RGB(0, 0, 0)
    ws.Range("B23").Font.Color = RGB(0, 0, 0)
End Sub

' CheckHeat - Heat Check worksheet button handler
Public Sub CheckHeat()
    On Error GoTo ErrorHandler

    Dim ws As Worksheet
    Set ws = ThisWorkbook.Worksheets("Heat Check")

    ' Generate correlation ID
    Dim corrID As String
    corrID = TFHelpers.GenerateCorrelationID()

    ' Read inputs
    Dim ticker As String, riskAmount As Double, bucket As String

    ticker = ws.Range("B3").Value
    riskAmount = ws.Range("B4").Value
    bucket = ws.Range("B5").Value

    ' Validate inputs
    If ticker = "" Or riskAmount <= 0 Or bucket = "" Then
        MsgBox "Please fill in all required fields: Ticker, Risk Amount, Bucket", vbExclamation, "Missing Inputs"
        Exit Sub
    End If

    ' Build command
    Dim cmd As String
    cmd = "heat --ticker " & ticker & " --risk " & riskAmount & " --bucket """ & bucket & """"

    ' Execute command
    Dim result As TFCommandResult
    result = ExecuteCommand(cmd, corrID)

    If result.Success Then
        ' Parse results
        Dim heatResult As TFHeatResult
        heatResult = TFHelpers.ParseHeatJSON(result.JsonOutput)

        ' Display portfolio heat results
        ws.Range("B10").Value = heatResult.PortfolioCurrentHeat
        ws.Range("B11").Value = heatResult.PortfolioNewHeat
        ws.Range("B12").Value = heatResult.PortfolioHeatPct
        ws.Range("B13").Value = heatResult.PortfolioCap
        ws.Range("B14").Value = IIf(heatResult.PortfolioExceeded, "YES", "NO")
        ws.Range("B15").Value = heatResult.PortfolioOverage

        ' Highlight exceeded caps in red
        If heatResult.PortfolioExceeded Then
            ws.Range("B14:B15").Font.Color = RGB(255, 0, 0)
            ws.Range("B14:B15").Font.Bold = True
        Else
            ws.Range("B14:B15").Font.Color = RGB(0, 128, 0)
            ws.Range("B14:B15").Font.Bold = False
        End If

        ' Display bucket heat results
        ws.Range("B18").Value = heatResult.BucketCurrentHeat
        ws.Range("B19").Value = heatResult.BucketNewHeat
        ws.Range("B20").Value = heatResult.BucketHeatPct
        ws.Range("B21").Value = heatResult.BucketCap
        ws.Range("B22").Value = IIf(heatResult.BucketExceeded, "YES", "NO")
        ws.Range("B23").Value = heatResult.BucketOverage

        ' Highlight exceeded caps in red
        If heatResult.BucketExceeded Then
            ws.Range("B22:B23").Font.Color = RGB(255, 0, 0)
            ws.Range("B22:B23").Font.Bold = True
        Else
            ws.Range("B22:B23").Font.Color = RGB(0, 128, 0)
            ws.Range("B22:B23").Font.Bold = False
        End If

        MsgBox "Heat check complete!" & vbCrLf & "Correlation ID: " & corrID, vbInformation, "Heat Check"
    Else
        MsgBox "Heat check failed:" & vbCrLf & result.ErrorOutput & vbCrLf & vbCrLf & _
               "Correlation ID: " & corrID, vbCritical, "Error"
    End If

    Exit Sub

ErrorHandler:
    MsgBox "Error in heat check: " & Err.Description & vbCrLf & _
           "Correlation ID: " & corrID, vbCritical, "Error"
    TFHelpers.LogMessage corrID, "ERROR", "CheckHeat failed: " & Err.Description
End Sub

' ClearHeatCheck - Clear Heat Check worksheet
Public Sub ClearHeatCheck()
    On Error Resume Next
    Dim ws As Worksheet
    Set ws = ThisWorkbook.Worksheets("Heat Check")

    ' Clear inputs
    ws.Range("B3:B5").ClearContents

    ' Clear results
    ws.Range("B10:B15").ClearContents
    ws.Range("B18:B23").ClearContents

    ' Reset formatting
    ws.Range("B14:B15").Font.Color = RGB(0, 0, 0)
    ws.Range("B14:B15").Font.Bold = False
    ws.Range("B22:B23").Font.Color = RGB(0, 0, 0)
    ws.Range("B22:B23").Font.Bold = False
End Sub

' SaveDecisionGO - Trade Entry worksheet - Save GO decision
Public Sub SaveDecisionGO()
    On Error GoTo ErrorHandler

    Dim ws As Worksheet
    Set ws = ThisWorkbook.Worksheets("Trade Entry")

    ' Generate correlation ID
    Dim corrID As String
    corrID = TFHelpers.GenerateCorrelationID()

    ' Read inputs
    Dim ticker As String, entryPrice As Double, atr As Double, method As String
    Dim bannerStatus As String, delta As Variant, maxLoss As Variant
    Dim bucket As String, preset As String

    ticker = ws.Range("B4").Value
    entryPrice = ws.Range("B5").Value
    atr = ws.Range("B6").Value
    method = ws.Range("B7").Value
    bannerStatus = ws.Range("B8").Value
    delta = ws.Range("B9").Value
    maxLoss = ws.Range("B10").Value
    bucket = ws.Range("B11").Value
    preset = ws.Range("B12").Value

    ' Validate required inputs
    If ticker = "" Or entryPrice <= 0 Or atr <= 0 Or method = "" Or bucket = "" Then
        ws.Range("B30").Value = "❌ Missing required inputs"
        MsgBox "Please fill in all required fields", vbExclamation, "Missing Inputs"
        Exit Sub
    End If

    ' Build command
    Dim cmd As String
    cmd = "save-decision --ticker " & ticker & " --entry " & entryPrice & _
          " --atr " & atr & " --method " & method & _
          " --bucket """ & bucket & """ --action GO"

    If bannerStatus <> "" Then cmd = cmd & " --banner " & bannerStatus
    If preset <> "" Then cmd = cmd & " --preset """ & preset & """"
    If Not IsEmpty(delta) And delta <> "" Then cmd = cmd & " --delta " & delta
    If Not IsEmpty(maxLoss) And maxLoss <> "" Then cmd = cmd & " --max-loss " & maxLoss

    ' Execute command
    Dim result As TFCommandResult
    result = ExecuteCommand(cmd, corrID)

    If result.Success Then
        ' Parse results
        Dim decisionResult As TFDecisionResult
        decisionResult = TFHelpers.ParseDecisionJSON(result.JsonOutput)

        ' Display gate results
        ws.Range("B18").Value = FormatGateStatus(decisionResult.Gate1Pass)
        ws.Range("B19").Value = FormatGateStatus(decisionResult.Gate2Pass)
        ws.Range("B20").Value = FormatGateStatus(decisionResult.Gate3Pass)
        ws.Range("B21").Value = FormatGateStatus(decisionResult.Gate4Pass)
        ws.Range("B22").Value = FormatGateStatus(decisionResult.Gate5Pass)

        ' Display decision results
        ws.Range("B25").Value = IIf(decisionResult.DecisionSaved, "YES", "NO")
        ws.Range("B26").Value = decisionResult.DecisionID
        ws.Range("B27").Value = decisionResult.RejectionReason
        ws.Range("B30").Value = "✅ Decision saved (" & corrID & ")"
        ws.Range("B30").Font.Color = RGB(0, 128, 0)

        ' Show success/rejection message
        If decisionResult.DecisionSaved Then
            MsgBox "Trade decision saved!" & vbCrLf & "Decision ID: " & decisionResult.DecisionID & vbCrLf & _
                   "Correlation ID: " & corrID, vbInformation, "Success"
            ' Clear form on success
            Call ClearTradeEntry
        Else
            MsgBox "Trade rejected:" & vbCrLf & decisionResult.RejectionReason & vbCrLf & vbCrLf & _
                   "Correlation ID: " & corrID, vbExclamation, "Trade Rejected"
        End If
    Else
        ws.Range("B30").Value = "❌ " & result.ErrorOutput
        ws.Range("B30").Font.Color = RGB(255, 0, 0)
        MsgBox "Save decision failed:" & vbCrLf & result.ErrorOutput & vbCrLf & vbCrLf & _
               "Correlation ID: " & corrID, vbCritical, "Error"
    End If

    Exit Sub

ErrorHandler:
    ws.Range("B30").Value = "❌ Error: " & Err.Description
    ws.Range("B30").Font.Color = RGB(255, 0, 0)
    MsgBox "Error saving decision: " & Err.Description & vbCrLf & _
           "Correlation ID: " & corrID, vbCritical, "Error"
    TFHelpers.LogMessage corrID, "ERROR", "SaveDecisionGO failed: " & Err.Description
End Sub

' SaveDecisionNOGO - Trade Entry worksheet - Save NO-GO decision
Public Sub SaveDecisionNOGO()
    On Error GoTo ErrorHandler

    Dim ws As Worksheet
    Set ws = ThisWorkbook.Worksheets("Trade Entry")

    ' Generate correlation ID
    Dim corrID As String
    corrID = TFHelpers.GenerateCorrelationID()

    ' Read inputs
    Dim ticker As String, bucket As String, preset As String, reason As String

    ticker = ws.Range("B4").Value
    bucket = ws.Range("B11").Value
    preset = ws.Range("B12").Value

    If ticker = "" Or bucket = "" Then
        ws.Range("B30").Value = "❌ Ticker and Bucket required"
        MsgBox "Please enter Ticker and Bucket", vbExclamation, "Missing Inputs"
        Exit Sub
    End If

    ' Prompt for rejection reason
    reason = InputBox("Enter reason for NO-GO decision:", "Rejection Reason")
    If reason = "" Then
        ws.Range("B30").Value = "❌ Cancelled"
        Exit Sub
    End If

    ' Build command
    Dim cmd As String
    cmd = "save-decision --ticker " & ticker & _
          " --bucket """ & bucket & """ --action NO-GO" & _
          " --reason """ & reason & """"

    If preset <> "" Then cmd = cmd & " --preset """ & preset & """"

    ' Execute command
    Dim result As TFCommandResult
    result = ExecuteCommand(cmd, corrID)

    If result.Success Then
        ' Parse results
        Dim decisionResult As TFDecisionResult
        decisionResult = TFHelpers.ParseDecisionJSON(result.JsonOutput)

        ' Display results
        ws.Range("B25").Value = "YES"
        ws.Range("B26").Value = decisionResult.DecisionID
        ws.Range("B27").Value = reason
        ws.Range("B30").Value = "✅ NO-GO saved (" & corrID & ")"
        ws.Range("B30").Font.Color = RGB(0, 128, 0)

        MsgBox "NO-GO decision saved!" & vbCrLf & "Decision ID: " & decisionResult.DecisionID & vbCrLf & _
               "Correlation ID: " & corrID, vbInformation, "Success"

        ' Clear form on success
        Call ClearTradeEntry
    Else
        ws.Range("B30").Value = "❌ " & result.ErrorOutput
        ws.Range("B30").Font.Color = RGB(255, 0, 0)
        MsgBox "Save decision failed:" & vbCrLf & result.ErrorOutput & vbCrLf & vbCrLf & _
               "Correlation ID: " & corrID, vbCritical, "Error"
    End If

    Exit Sub

ErrorHandler:
    ws.Range("B30").Value = "❌ Error: " & Err.Description
    ws.Range("B30").Font.Color = RGB(255, 0, 0)
    MsgBox "Error saving NO-GO decision: " & Err.Description & vbCrLf & _
           "Correlation ID: " & corrID, vbCritical, "Error"
    TFHelpers.LogMessage corrID, "ERROR", "SaveDecisionNOGO failed: " & Err.Description
End Sub

' ClearTradeEntry - Clear Trade Entry worksheet
Public Sub ClearTradeEntry()
    On Error Resume Next
    Dim ws As Worksheet
    Set ws = ThisWorkbook.Worksheets("Trade Entry")

    ' Clear inputs
    ws.Range("B4:B12").ClearContents

    ' Clear gate status
    ws.Range("B18:B22").ClearContents

    ' Clear results
    ws.Range("B25:B30").ClearContents
    ws.Range("B30").Font.Color = RGB(0, 0, 0)
End Sub

' FormatGateStatus - Helper function to format gate pass/fail status
Private Function FormatGateStatus(ByVal passed As Boolean) As String
    If passed Then
        FormatGateStatus = "✅ PASS"
    Else
        FormatGateStatus = "❌ FAIL"
    End If
End Function

'=============================================================================
