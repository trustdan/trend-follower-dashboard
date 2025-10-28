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

        ' Brief pause to avoid CPU spinning
        Application.Wait Now + TimeValue("0:00:00.1")
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

    ' Try to get from Setup sheet
    Set ws = ThisWorkbook.Worksheets("Setup")
    If Not ws Is Nothing Then
        Set pathRange = ws.Range("EnginePathSetting")  ' Named range
        If Not pathRange Is Nothing Then
            If pathRange.Value <> "" Then
                GetEnginePath = pathRange.Value
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

    ' Try to get from Setup sheet
    Set ws = ThisWorkbook.Worksheets("Setup")
    If Not ws Is Nothing Then
        Set pathRange = ws.Range("DatabasePathSetting")  ' Named range
        If Not pathRange Is Nothing Then
            If pathRange.Value <> "" Then
                GetDatabasePath = pathRange.Value
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
