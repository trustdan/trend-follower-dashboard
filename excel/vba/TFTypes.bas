Attribute VB_Name = "TFTypes"
'=============================================================================
' TFTypes.bas - Type Definitions for Trading Engine v3
'=============================================================================
' Purpose: Define VBA types matching JSON response structures from tf-engine
'
' Architecture: Engine-first (Go backend) with Excel as thin UI
' Transport: CLI via shell execution (stdout = JSON, stderr = errors)
' Philosophy: Keep VBA thin - just parse JSON and display results
'
' Created: 2025-10-27 (M19 - VBA Implementation)
'=============================================================================

Option Explicit

'-----------------------------------------------------------------------------
' POSITION SIZING TYPES
'-----------------------------------------------------------------------------

' Position sizing result (from "size" command)
Public Type TFSizingResult
    RiskDollars As Double       ' R = Equity × RiskPct
    StopDistance As Double      ' K × ATR
    InitialStop As Double       ' Entry - StopDistance
    Shares As Long              ' For stocks: floor(R ÷ StopDistance)
    Contracts As Long           ' For options: calculated based on method
    ActualRisk As Double        ' Actual risk = Shares × StopDistance
    Method As String            ' "stock", "opt-delta-atr", "opt-maxloss"
End Type

'-----------------------------------------------------------------------------
' CHECKLIST TYPES
'-----------------------------------------------------------------------------

' Checklist evaluation result (from "checklist" command)
Public Type TFChecklistResult
    Banner As String                    ' "GREEN", "YELLOW", "RED"
    MissingCount As Long                ' Number of failed checks
    MissingItems As String              ' Comma-separated list of missing items
    EvaluationTimestamp As String       ' ISO8601 timestamp
    AllowSave As Boolean                ' True if GREEN (all checks pass)
End Type

'-----------------------------------------------------------------------------
' HEAT MANAGEMENT TYPES
'-----------------------------------------------------------------------------

' Heat check result (from "heat" command)
Public Type TFHeatResult
    CurrentPortfolioHeat As Double      ' Sum of risk across all open positions
    NewPortfolioHeat As Double          ' Current + proposed trade risk
    PortfolioHeatPct As Double          ' (New / Cap) × 100
    PortfolioCap As Double              ' Equity × HeatCap_H_pct (e.g., 4%)
    PortfolioCapExceeded As Boolean     ' True if new heat > cap
    PortfolioOverage As Double          ' Amount over cap (0 if not exceeded)

    CurrentBucketHeat As Double         ' Sum of risk in target bucket
    NewBucketHeat As Double             ' Current + proposed risk in bucket
    BucketHeatPct As Double             ' (New / Cap) × 100
    BucketCap As Double                 ' Equity × BucketHeatCap_pct (e.g., 1.5%)
    BucketCapExceeded As Boolean        ' True if new bucket heat > cap
    BucketOverage As Double             ' Amount over bucket cap

    Allowed As Boolean                  ' True if both caps are OK
End Type

'-----------------------------------------------------------------------------
' IMPULSE TIMER TYPES
'-----------------------------------------------------------------------------

' Impulse timer check result (from "timer" command)
Public Type TFTimerResult
    TimerActive As Boolean              ' True if timer is running
    BrakeCleared As Boolean             ' True if 2 minutes elapsed
    ElapsedSeconds As Long              ' Time since evaluation
    RemainingSeconds As Long            ' Time until brake clears
    StartedAt As String                 ' ISO8601 timestamp when timer started
    Ticker As String                    ' Ticker being evaluated
End Type

'-----------------------------------------------------------------------------
' CANDIDATE TYPES
'-----------------------------------------------------------------------------

' Single candidate ticker
Public Type TFCandidate
    ID As Long                          ' Database primary key
    DateStr As String                   ' Date string (YYYY-MM-DD)
    Ticker As String                    ' Stock ticker symbol
    PresetName As String                ' Preset name (e.g., "TF_BREAKOUT_LONG")
    PresetID As Long                    ' Foreign key to presets table
    Sector As String                    ' Industry sector (optional)
    Bucket As String                    ' Sector bucket (e.g., "Tech/Comm")
End Type

' List of candidates (from "list-candidates" command)
Public Type TFCandidatesList
    Candidates() As TFCandidate         ' Dynamic array of candidates
    Count As Long                       ' Number of candidates
    DateStr As String                   ' Date for this candidate list
End Type

' Candidate check result (from "check-candidate" command)
Public Type TFCandidateCheck
    Found As Boolean                    ' True if ticker in today's candidates
    Ticker As String                    ' Ticker that was checked
    DateStr As String                   ' Date checked
    Message As String                   ' Human-readable message
End Type

'-----------------------------------------------------------------------------
' COOLDOWN TYPES
'-----------------------------------------------------------------------------

' Single cooldown record
Public Type TFCooldown
    ID As Long                          ' Database primary key
    Bucket As String                    ' Sector bucket on cooldown
    TriggeredAt As String               ' ISO8601 timestamp
    ClearsAt As String                  ' ISO8601 timestamp
    Active As Boolean                   ' True if still in cooldown
    Reason As String                    ' Why cooldown triggered
End Type

' Cooldown check result (from "check-cooldown" command)
Public Type TFCooldownCheck
    OnCooldown As Boolean               ' True if bucket is on cooldown
    Bucket As String                    ' Bucket that was checked
    ClearsAt As String                  ' When cooldown ends (if active)
    Message As String                   ' Human-readable message
End Type

' List of cooldowns (from "list-cooldowns" command)
Public Type TFCooldownsList
    Cooldowns() As TFCooldown           ' Dynamic array of cooldowns
    Count As Long                       ' Number of cooldowns
    ActiveCount As Long                 ' Number of active cooldowns
End Type

'-----------------------------------------------------------------------------
' POSITION TYPES
'-----------------------------------------------------------------------------

' Single position record
Public Type TFPosition
    ID As Long                          ' Database primary key
    Ticker As String                    ' Stock ticker
    Bucket As String                    ' Sector bucket
    OpenDate As String                  ' Date position opened (YYYY-MM-DD)
    UnitsOpen As Long                   ' Number of shares/contracts open
    TotalOpenR As Double                ' Total risk dollars (R) for position
    Status As String                    ' "open", "closed"
    EntryPrice As Double                ' Original entry price
    CurrentStop As Double               ' Current stop loss level
    HighestStop As Double               ' Highest stop ever set (for tracking)
End Type

' List of positions (from "list-positions" command)
Public Type TFPositionsList
    Positions() As TFPosition           ' Dynamic array of positions
    Count As Long                       ' Number of positions
    OpenCount As Long                   ' Number of open positions
    TotalRisk As Double                 ' Sum of TotalOpenR for all open positions
End Type

'-----------------------------------------------------------------------------
' SETTINGS TYPES
'-----------------------------------------------------------------------------

' Application settings (from "get-settings" command)
Public Type TFSettings
    Equity_E As Double                  ' Account equity ($)
    RiskPct_r As Double                 ' Risk per trade (as decimal, e.g., 0.0075 = 0.75%)
    HeatCap_H_pct As Double             ' Portfolio heat cap (as decimal, e.g., 0.04 = 4%)
    BucketHeatCap_pct As Double         ' Bucket heat cap (as decimal, e.g., 0.015 = 1.5%)
    StopMultiple_K As Long              ' K multiple for stop distance (e.g., 2)
End Type

'-----------------------------------------------------------------------------
' DECISION TYPES
'-----------------------------------------------------------------------------

' Save decision result (from "save-decision" command)
Public Type TFSaveDecisionResult
    Accepted As Boolean                 ' True if all 5 hard gates passed
    DecisionID As Long                  ' Database ID if accepted (0 if rejected)
    Timestamp As String                 ' ISO8601 timestamp
    Reason As String                    ' Rejection reason (empty if accepted)
    GatesFailed As String               ' Comma-separated list of failed gates
End Type

'-----------------------------------------------------------------------------
' ERROR TYPES
'-----------------------------------------------------------------------------

' Engine error (when command fails)
Public Type TFEngineError
    HasError As Boolean                 ' True if error occurred
    ErrorMessage As String              ' Human-readable error from stderr
    CorrelationID As String             ' Correlation ID for log lookup
    ExitCode As Long                    ' Process exit code (0 = success)
End Type

'-----------------------------------------------------------------------------
' COMMAND EXECUTION RESULT (Generic wrapper)
'-----------------------------------------------------------------------------

' Generic result wrapper for any engine command
Public Type TFCommandResult
    Success As Boolean                  ' True if command executed successfully
    JsonOutput As String                ' Raw JSON from stdout
    ErrorOutput As String               ' Raw error text from stderr
    ExitCode As Long                    ' Process exit code
    CorrelationID As String             ' Correlation ID for this call
End Type

'=============================================================================
' NOTES ON TYPE USAGE
'=============================================================================
'
' JSON PARSING STRATEGY:
' VBA does not have native JSON support. We use simple string parsing for
' key-value extraction. This is sufficient because:
' 1. JSON schemas are validated in the Go engine (M17-M18)
' 2. VBA is a thin bridge - no business logic here
' 3. Complex parsing would add unnecessary dependencies
'
' ERROR HANDLING PATTERN:
' All TFEngine functions return TFCommandResult first, then typed result
' Example:
'   Dim cmdResult As TFCommandResult
'   Dim sizeResult As TFSizingResult
'   cmdResult = TFEngine.ExecuteSize(...)
'   If cmdResult.Success Then
'       sizeResult = TFHelpers.ParseSizingJSON(cmdResult.JsonOutput)
'   Else
'       ' Handle error using cmdResult.ErrorOutput
'   End If
'
' CORRELATION ID PATTERN:
' Every call generates or accepts a correlation ID that appears in:
' 1. Excel status cells
' 2. tf-engine.log (Go backend)
' 3. TradingSystem_Debug.log (VBA frontend)
' This enables cross-referencing issues between frontend and backend logs.
'
' TYPE SAFETY NOTES:
' - String dates use ISO8601 format: "2025-10-27T14:30:00-05:00"
' - Percentages in TFSettings are decimals: 0.0075 = 0.75%
' - Currency values (RiskDollars, etc.) use Double (VBA Currency type has limits)
' - Arrays are 0-based in these types (VBA default)
'
'=============================================================================
