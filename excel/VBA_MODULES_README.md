# VBA Modules Documentation - Trading Engine v3

**Created:** 2025-10-27 (M19 - VBA Implementation)
**Purpose:** Bridge between Excel UI and Go backend (tf-engine.exe)
**Philosophy:** Keep VBA thin - just shell execution and JSON parsing

---

## Architecture Overview

```
Excel UI (worksheets, buttons, cells)
    ↓
VBA Bridge (these modules)
    ↓
Shell Execution (WScript.Shell.Exec)
    ↓
tf-engine.exe (Go backend)
    ↓
SQLite Database (trading.db)
```

**Key Design Principles:**
1. **No business logic in VBA** - All trading rules enforced in Go engine
2. **JSON as contract** - Engine outputs validated JSON, VBA just parses it
3. **Correlation IDs** - Every call tracked for debugging across VBA and Go logs
4. **Fail loudly** - Never swallow errors silently
5. **Simple over clever** - No external dependencies, basic string parsing

---

## Module Structure

### 1. **TFTypes.bas** - Type Definitions

**Purpose:** Define VBA types matching JSON response structures from engine

**Key Types:**
- `TFSizingResult` - Position sizing calculation results
- `TFChecklistResult` - Checklist evaluation (GREEN/YELLOW/RED banner)
- `TFHeatResult` - Portfolio and bucket heat management
- `TFTimerResult` - 2-minute impulse brake status
- `TFCandidateCheck` - Candidate ticker validation
- `TFCooldownCheck` - Bucket cooldown status
- `TFSettings` - Application settings (Equity, RiskPct, etc.)
- `TFSaveDecisionResult` - Save decision outcome (5 hard gates)
- `TFCommandResult` - Generic wrapper for any engine command

**Usage Pattern:**
```vba
Dim cmdResult As TFCommandResult
Dim sizeResult As TFSizingResult

cmdResult = TFEngine.Engine_Size(180, 1.5, "stock")
If cmdResult.Success Then
    sizeResult = TFHelpers.ParseSizingJSON(cmdResult.JsonOutput)
    ' Use sizeResult.Shares, sizeResult.InitialStop, etc.
End If
```

---

### 2. **TFHelpers.bas** - Utility Functions

**Purpose:** JSON parsing, error handling, logging, validation, formatting

**Key Functions:**

#### JSON Parsing
- `ExtractJSONValue(jsonStr, key)` - Extract value for a key from JSON
- `ExtractJSONArray(arrayStr)` - Parse JSON array into Collection
- `ParseSizingJSON(jsonStr)` → `TFSizingResult`
- `ParseChecklistJSON(jsonStr)` → `TFChecklistResult`
- `ParseHeatJSON(jsonStr)` → `TFHeatResult`
- `ParseTimerJSON(jsonStr)` → `TFTimerResult`
- `ParseSettingsJSON(jsonStr)` → `TFSettings`
- `ParseSaveDecisionJSON(jsonStr)` → `TFSaveDecisionResult`

#### Logging
- `GenerateCorrelationID()` - Create unique ID (format: `YYYYMMDD-HHMMSS-XXXX`)
- `LogMessage(corrID, level, message)` - Write to `TradingSystem_Debug.log`

#### Validation
- `ValidateTicker(ticker)` - Check ticker format
- `ValidatePositiveNumber(value)` - Check positive number

#### Formatting
- `FormatCurrency(value)` - Format as `$1,234.56`
- `FormatPercent(value)` - Format as `12.34%`
- `FormatTimestamp(iso8601)` - Convert ISO timestamp to readable format

**JSON Parsing Approach:**
- Simple string parsing (no external libraries)
- Assumes well-formed JSON from engine (validated in Go)
- Returns sensible defaults for missing values
- Graceful error handling with `On Error Resume Next`

**Limitations:**
- Does not handle deeply nested objects/arrays
- No JSON schema validation (done in Go engine)
- **This is acceptable because** engine outputs are tested and validated (M17-M18)

---

### 3. **TFEngine.bas** - Engine Communication

**Purpose:** Bridge to tf-engine.exe via shell execution

**Core Function:**
```vba
Function ExecuteCommand(command As String, Optional corrID As String = "") As TFCommandResult
```

All other functions wrap this core function with specific commands.

#### Command Functions

**Position Sizing:**
```vba
Engine_Size(entry, atr, method, [equity], [riskPct], [k], [delta], [maxLoss], [corrID])
' Returns: TFCommandResult with JSON for parsing
' Example: Engine_Size(180, 1.5, "stock")
```

**Checklist Evaluation:**
```vba
Engine_Checklist(ticker, checks, [corrID])
' checks = Collection of 6 boolean values
' Returns: TFCommandResult with banner (GREEN/YELLOW/RED)
```

**Heat Management:**
```vba
Engine_Heat(addR, bucket, [corrID])
' Checks if adding trade would exceed portfolio or bucket caps
' Returns: TFCommandResult with heat breakdown
```

**Impulse Timer:**
```vba
Engine_CheckTimer(ticker, [corrID])
' Checks if 2-minute impulse brake has elapsed
' Returns: TFCommandResult with timer status
```

**Candidate Management:**
```vba
Engine_ImportCandidates(tickers, preset, [corrID])
Engine_ListCandidates([dateStr], [corrID])
Engine_CheckCandidate(ticker, [corrID])
```

**Cooldown Management:**
```vba
Engine_CheckCooldown(bucket, [corrID])
Engine_ListCooldowns([activeOnly], [corrID])
```

**Settings:**
```vba
Engine_GetSettings([corrID])
Engine_SetSetting(key, value, [corrID])
```

**Save Decision (5 Hard Gates):**
```vba
Engine_SaveDecision(ticker, entry, atr, k, method, riskDollars,
                    shares, contracts, banner, bucket, preset, [corrID])
' Enforces ALL 5 gates in Go engine:
'   1. Banner must be GREEN
'   2. Ticker in today's candidates
'   3. 2-minute impulse brake elapsed
'   4. Bucket not on cooldown
'   5. Heat caps not exceeded
```

**Database Initialization:**
```vba
Engine_Init([corrID])
' Creates trading.db schema if not exists
```

#### Configuration

Engine path and database path are read from named ranges on Setup sheet:
- `EnginePathSetting` - Path to tf-engine.exe
- `DatabasePathSetting` - Path to trading.db

Defaults if not set:
- Engine: `<workbook_path>\tf-engine.exe`
- Database: `<workbook_path>\trading.db`

---

### 4. **TFTests.bas** - VBA Unit Tests

**Purpose:** Test VBA modules before Windows integration (M21)

**Main Function:**
```vba
RunAllTests()
' Executes all test functions and reports results to "VBA Tests" worksheet
```

**Test Categories:**
1. **JSON Parsing Tests** - Verify parsing for all response types
2. **Helper Function Tests** - Validate utility functions
3. **Validation Tests** - Check input validation logic
4. **Formatting Tests** - Verify display formatters
5. **Shell Execution Test** - Confirm can execute tf-engine.exe

**Test Output:**
- Results written to "VBA Tests" worksheet
- Color-coded (✅ green for pass, ❌ red for fail)
- Summary with pass/fail counts and total time

**Adding New Tests:**
```vba
Private Function Test_YourTestName() As TestResult
    Dim result As TestResult
    result.TestName = "YourTestName"

    ' ... test logic ...

    result.Passed = True/False
    result.Message = "Description"
    Test_YourTestName = result
End Function

' Add to RunAllTests():
AddTestResult results, resultCount, Test_YourTestName()
```

---

## Usage Patterns

### Standard Call Pattern

```vba
Sub Example_PositionSizing()
    Dim cmdResult As TFCommandResult
    Dim sizeResult As TFSizingResult
    Dim corrID As String

    ' Generate correlation ID
    corrID = TFHelpers.GenerateCorrelationID()

    ' Log call
    TFHelpers.LogMessage corrID, "INFO", "Calculating position size for AAPL"

    ' Call engine
    cmdResult = TFEngine.Engine_Size(180, 1.5, "stock", corrID:=corrID)

    ' Check success
    If cmdResult.Success Then
        ' Parse JSON
        sizeResult = TFHelpers.ParseSizingJSON(cmdResult.JsonOutput)

        ' Update UI
        Range("SharesResult").Value = sizeResult.Shares
        Range("InitialStopResult").Value = sizeResult.InitialStop
        Range("RiskDollarsResult").Value = sizeResult.RiskDollars
        Range("StatusCell").Value = "✅ Success (corr_id: " & corrID & ")"

        TFHelpers.LogMessage corrID, "INFO", "Position sizing succeeded"
    Else
        ' Handle error
        Range("StatusCell").Value = "❌ " & cmdResult.ErrorOutput
        Range("StatusCell").Interior.Color = RGB(255, 199, 206)  ' Red

        TFHelpers.LogMessage corrID, "ERROR", "Position sizing failed: " & cmdResult.ErrorOutput

        MsgBox "Error: " & cmdResult.ErrorOutput & vbCrLf & _
               "Correlation ID: " & corrID & vbCrLf & _
               "Check logs for details.", vbCritical, "Engine Error"
    End If
End Sub
```

### Error Display Pattern

**Always include:**
1. Correlation ID in status cells
2. Error message from stderr
3. Log entry with correlation ID
4. User-friendly error dialog

**Example:**
```vba
If Not cmdResult.Success Then
    Dim errMsg As String
    errMsg = "ERROR: " & cmdResult.ErrorOutput & vbCrLf & _
             "Correlation ID: " & cmdResult.CorrelationID & vbCrLf & vbCrLf & _
             "Check TradingSystem_Debug.log and tf-engine.log for details."

    Range("StatusCell").Value = "❌ Failed"
    Range("ErrorDetailsCell").Value = errMsg

    TFHelpers.LogMessage cmdResult.CorrelationID, "ERROR", cmdResult.ErrorOutput

    MsgBox errMsg, vbCritical, "Trade Decision Failed"
End If
```

### Complete Workflow Example (Save Decision)

```vba
Sub Button_SaveDecision_Click()
    Dim cmdResult As TFCommandResult
    Dim saveResult As TFSaveDecisionResult
    Dim corrID As String

    ' Generate correlation ID
    corrID = TFHelpers.GenerateCorrelationID()

    ' Read inputs from worksheet
    Dim ticker As String
    Dim entry As Double
    Dim atr As Double
    Dim k As Long
    Dim riskDollars As Double
    Dim shares As Long
    Dim banner As String
    Dim bucket As String

    ticker = Range("TickerInput").Value
    entry = Range("EntryInput").Value
    atr = Range("ATRInput").Value
    k = Range("KMultipleInput").Value
    riskDollars = Range("RiskDollarsResult").Value
    shares = Range("SharesResult").Value
    banner = Range("BannerResult").Value
    bucket = Range("BucketInput").Value

    ' Validate inputs
    If Not TFHelpers.ValidateTicker(ticker) Then
        MsgBox "Invalid ticker symbol", vbExclamation
        Exit Sub
    End If

    If Not TFHelpers.ValidatePositiveNumber(entry) Or _
       Not TFHelpers.ValidatePositiveNumber(atr) Then
        MsgBox "Entry and ATR must be positive numbers", vbExclamation
        Exit Sub
    End If

    ' Log attempt
    TFHelpers.LogMessage corrID, "INFO", "Attempting to save decision for " & ticker

    ' Call engine (enforces 5 hard gates)
    cmdResult = TFEngine.Engine_SaveDecision( _
        ticker, entry, atr, k, "stock", riskDollars, _
        shares, 0, banner, bucket, "TF_BREAKOUT_LONG", corrID)

    ' Process result
    If cmdResult.Success Then
        saveResult = TFHelpers.ParseSaveDecisionJSON(cmdResult.JsonOutput)

        If saveResult.Accepted Then
            ' Success - clear form and show confirmation
            Range("StatusCell").Value = "✅ Decision saved (ID: " & saveResult.DecisionID & ")"
            Range("StatusCell").Interior.Color = RGB(198, 239, 206)  ' Green

            TFHelpers.LogMessage corrID, "INFO", "Decision saved with ID " & saveResult.DecisionID

            MsgBox "Trade decision saved!" & vbCrLf & _
                   "Decision ID: " & saveResult.DecisionID & vbCrLf & _
                   "Timestamp: " & TFHelpers.FormatTimestamp(saveResult.Timestamp), _
                   vbInformation, "Success"

            ' Clear form for next trade
            Call ClearTradeForm
        Else
            ' Rejected by hard gates
            Range("StatusCell").Value = "❌ REJECTED: " & saveResult.Reason
            Range("StatusCell").Interior.Color = RGB(255, 199, 206)  ' Red

            TFHelpers.LogMessage corrID, "WARN", "Decision rejected: " & saveResult.Reason

            MsgBox "Trade decision REJECTED" & vbCrLf & vbCrLf & _
                   "Reason: " & saveResult.Reason & vbCrLf & _
                   "Failed gates: " & saveResult.GatesFailed & vbCrLf & vbCrLf & _
                   "Correlation ID: " & corrID, _
                   vbExclamation, "Decision Rejected"
        End If
    Else
        ' Engine error
        Range("StatusCell").Value = "❌ Engine Error"
        Range("StatusCell").Interior.Color = RGB(255, 199, 206)  ' Red

        TFHelpers.LogMessage corrID, "ERROR", "Engine error: " & cmdResult.ErrorOutput

        MsgBox "Engine execution failed!" & vbCrLf & vbCrLf & _
               cmdResult.ErrorOutput & vbCrLf & vbCrLf & _
               "Correlation ID: " & corrID, vbCritical, "Engine Error"
    End If
End Sub
```

---

## Logging & Debugging

### Log Files

**TradingSystem_Debug.log** (VBA side)
- Location: Same directory as workbook
- Format: `[TIMESTAMP] [LEVEL] [CORR_ID] Message`
- Rotation: Auto-rotates at 5 MB
- Use: VBA operations, UI events, errors

**tf-engine.log** (Go side)
- Location: Same directory as tf-engine.exe
- Format: JSON structured logs with correlation IDs
- Rotation: Managed by Go logger
- Use: Business logic, database operations, validations

### Correlation ID Flow

```
User clicks button in Excel
    ↓
VBA generates corr_id: "20251027-143052-7A3F"
    ↓
VBA logs: "Calling tf-engine size with corr_id 20251027-143052-7A3F"
    ↓
VBA passes: --corr-id 20251027-143052-7A3F
    ↓
Go engine logs: {"corr_id": "20251027-143052-7A3F", "msg": "Position sizing..."}
    ↓
Go returns JSON to stdout
    ↓
VBA logs: "Received response for corr_id 20251027-143052-7A3F"
    ↓
Excel displays: "✅ Success (corr_id: 20251027-143052-7A3F)"
```

**Cross-referencing issues:**
1. User reports issue with correlation ID
2. Search TradingSystem_Debug.log for corr_id → See VBA side
3. Search tf-engine.log for corr_id → See Go side
4. Compare timestamps and data flow

---

## Error Handling Strategy

### Principle: Fail Loudly

**Never silently ignore errors. Always:**
1. Log with correlation ID
2. Display to user with actionable message
3. Show correlation ID for support
4. Preserve data (don't clear form on errors)

### Error Types

**1. VBA Errors (Err object)**
- Trapped with `On Error GoTo ErrorHandler`
- Logged as "VBA Error: [description]"
- Exit code: -999

**2. Engine Errors (stderr output)**
- Captured from exec.StdErr
- Exit code from exec.ExitCode (non-zero)
- Displayed verbatim to user

**3. Validation Errors (pre-flight checks)**
- Caught before calling engine
- Show friendly message
- Don't call engine if inputs invalid

**4. Business Logic Rejections (5 hard gates)**
- Returned as JSON with `accepted: false`
- Not errors - expected behavior
- Show reason and failed gates

---

## Testing Strategy

### VBA Unit Tests (TFTests.bas)
- Run before Windows integration
- Test JSON parsing, helpers, validation
- Verify shell execution works
- Output to "VBA Tests" worksheet

### Manual Integration Tests (M21)
- Full workflows through Excel UI
- All buttons and commands
- Error scenarios
- Correlation ID verification
- See: WINDOWS_TESTING.md (created in M20)

### Test Data
- `test-data/json-examples/responses/` - Sample JSON responses
- Use for VBA unit tests
- Matches actual engine outputs (validated in M17-M18)

---

## Deployment Instructions

### For Development (Linux/WSL - M19)
1. VBA modules created as `.bas` text files
2. Version controlled in git
3. Syntax-checked but NOT executed yet

### For Windows Testing (M21)
1. Copy `.bas` files to Windows machine
2. Run `windows-import-vba.bat` to import into Excel workbook
3. Enable macros in Excel
4. Run VBA tests from "VBA Tests" worksheet
5. Execute integration tests per WINDOWS_TESTING.md

### Module Import Order
1. TFTypes.bas (no dependencies)
2. TFHelpers.bas (uses TFTypes)
3. TFEngine.bas (uses TFTypes, TFHelpers)
4. TFTests.bas (uses all above)

---

## Troubleshooting

### "Engine not found" error
- Check Setup sheet `EnginePathSetting` value
- Verify tf-engine.exe exists at that path
- Try absolute path instead of relative

### "Command timed out" error
- Check if engine is hanging
- Increase `COMMAND_TIMEOUT_SECONDS` in TFEngine.bas
- Check tf-engine.log for issues

### JSON parsing errors
- Check cmdResult.JsonOutput in Immediate window
- Verify JSON is well-formed
- Check if engine returned error on stderr instead

### Correlation ID not appearing in logs
- Verify TradingSystem_Debug.log file is writable
- Check file permissions
- Look for VBA error in LogMessage function

### Tests failing on Windows
- Ensure tf-engine.exe is in correct location
- Check Windows execution permissions
- Verify trading.db is initialized (`tf-engine init`)
- See detailed M21 testing procedures

---

## Design Decisions & Rationale

### Why no external JSON library?
- **Reason:** Minimizes dependencies, easier deployment
- **Trade-off:** Limited parsing capabilities
- **Acceptable because:** Engine validates JSON; VBA just extracts known keys

### Why synchronous shell execution?
- **Reason:** Simple, reliable, matches Excel's synchronous UI
- **Trade-off:** Excel freezes during command execution
- **Acceptable because:** Commands complete in < 1 second typically

### Why correlation IDs?
- **Reason:** Essential for debugging distributed system (VBA + Go)
- **Trade-off:** Slight overhead in logging
- **Acceptable because:** Debugging benefits far outweigh cost

### Why thin VBA?
- **Reason:** Business logic in Go is testable, maintainable, cross-platform
- **Trade-off:** More complex architecture
- **Acceptable because:** Enables future web UI without rewriting logic

### Why no error recovery/retry?
- **Reason:** Fail loudly, require manual intervention
- **Trade-off:** Less "convenient" for users
- **Acceptable because:** Trading system - no silent failures ever

---

## Future Enhancements (Post-M19)

### Potential Improvements (v3.1+)
1. **Async execution** - Run commands in background with progress bar
2. **JSON library** - Use VBA-JSON for more robust parsing
3. **HTTP transport** - Switch to HTTP API instead of CLI (toggle in Setup)
4. **Caching** - Cache settings and candidates to reduce engine calls
5. **Bulk operations** - Import multiple tickers at once
6. **Error recovery** - Auto-retry transient errors

**Current Status:** NOT IMPLEMENTED - Keep it simple for v3.0

---

## M19 Exit Criteria

**Definition of Done:**
- ✅ Four VBA modules created as `.bas` text files
- ✅ All functions implement validated JSON contracts (M17-M18)
- ✅ Syntax is valid (no compilation errors)
- ✅ Code matches architecture philosophy (thin bridge)
- ✅ Comprehensive inline documentation
- ✅ Usage examples provided
- ✅ Test functions written (execution deferred to M21)

**Next Steps:** M20 (Create Windows integration package)

---

## References

- **Trading-Engine-v3_Step-by-Step-Plan.md** - Overall project plan
- **CLAUDE.md** - Project philosophy and development rules
- **WHY.md** - Core philosophy (discipline over flexibility)
- **BDD_GUIDE.md** - Testing approach
- **test-data/json-examples/** - Validated JSON response examples

---

**Remember:** This is a discipline enforcement system, not a flexible trading platform. VBA serves that mission by being a thin, reliable bridge to the Go engine where all rules are enforced.

Code serves discipline. Discipline does not serve code.
