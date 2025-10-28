# M17-M18: Critical Issues to Fix Before M19

**Purpose:** Track issues discovered during JSON contract validation that MUST be fixed before VBA implementation.

**Priority:** BLOCKER for M19 (VBA Implementation)

---

## Issue 1: Mixed JSON and Text Output in CLI Commands

### Problem

Several CLI commands write human-readable text AND JSON to stdout, making it impossible for VBA to parse reliably.

### Examples

**checklist command:**
```bash
$ tf-engine checklist --ticker AAPL --from-preset --trend-pass --liquidity-pass --tv-confirm --earnings-ok --journal-ok
⏱️  Impulse brake timer started
   Wait 2 minutes before saving decision
{
  "banner": "GREEN",
  "missing_count": 0,
  ...
}
```

**cooldown commands:**
```bash
$ tf-engine check-cooldown --bucket Energy
✓ Bucket Energy is NOT in cooldown
```

**heat check command:**
```bash
$ tf-engine check-heat --add-r 75
(no JSON output at all - just text or empty)
```

### Impact

- ❌ VBA cannot parse mixed text/JSON from stdout
- ❌ Breaks automated parsing in Excel integration
- ❌ Makes CLI/HTTP parity impossible to achieve

### Solution Options

#### Option A: JSON-Only Flag (RECOMMENDED)

Add `--format` flag to all commands:
- `--format human` (default): Current human-friendly output
- `--format json`: Pure JSON output to stdout

**Pros:**
- Backward compatible (existing users get same output)
- Clean separation between human and machine interfaces
- Easy to implement
- Standard pattern (e.g., `kubectl`, `aws cli`)

**Cons:**
- Requires adding flag to every command call from VBA

**Implementation:**
```go
// Add to each command
cmd.Flags().String("format", "human", "Output format: human or json")

// In command handler
format, _ := cmd.Flags().GetString("format")
if format == "json" {
    // Output only JSON
    fmt.Println(jsonString)
} else {
    // Output human-friendly text
    fmt.Println("⏱️  Impulse brake timer started")
    fmt.Println(jsonString)
}
```

#### Option B: Separate Stderr/Stdout Streams

- Human text → stderr
- JSON → stdout

**Pros:**
- No flag needed
- Clean programmatic interface

**Cons:**
- Changes user experience (text on stderr might look like errors)
- Less flexible

#### Option C: Dedicated API Commands

Create separate commands for programmatic use:
- User commands: `tf-engine checklist` (human output)
- API commands: `tf-engine api checklist` (JSON only)

**Pros:**
- Very clear separation
- No breaking changes

**Cons:**
- More commands to maintain
- Duplication

### Recommendation

**Use Option A: `--format json` flag**

**Rationale:**
- Industry standard pattern
- Backward compatible
- Flexible (can add XML, YAML later if needed)
- Clear intent when reading VBA code

**Files to Modify:**
- `internal/cli/checklist.go`
- `internal/cli/check_cooldown.go`
- `internal/cli/list_cooldowns.go`
- `internal/cli/trigger_cooldown.go`
- `internal/cli/check_heat.go`
- `internal/cli/check_timer.go`
- `internal/cli/check_candidate.go`
- Any other commands with mixed output

**Acceptance Criteria:**
- ✅ All commands support `--format json` flag
- ✅ With `--format json`, stdout contains ONLY valid JSON (one object per command)
- ✅ Human-readable messages still available with `--format human` (default)
- ✅ VBA can parse JSON without string manipulation

---

## Issue 2: Logging to Stdout

### Problem

Logger writes to BOTH stdout and file, polluting stdout with log lines:

```go
// internal/logx/logger.go:45
logger.SetOutput(io.MultiWriter(os.Stdout, file))
```

This causes CLI output to be:
```
{"corr_id":"...","level":"info","message":"Trading Engine starting",...}
{"corr_id":"...","level":"info","message":"Starting position sizing",...}
{
  "risk_dollars": 75,
  ...
}
{"corr_id":"...","level":"info","message":"Trading Engine completed",...}
```

### Impact

- ❌ VBA receives log lines mixed with actual JSON response
- ❌ Must parse/filter stdout to extract actual data
- ❌ Fragile and error-prone

### Solution

**Remove stdout from logger output:**

```go
// internal/logx/logger.go:45
// OLD:
logger.SetOutput(io.MultiWriter(os.Stdout, file))

// NEW:
logger.SetOutput(file)  // Log ONLY to file
```

**Rationale:**
- Logs belong in log file, not stdout
- Stdout should be reserved for command output
- User can always `tail -f tf-engine.log` to watch logs
- Stderr can still be used for user-facing errors

**Alternative:** If we want logs visible during development:
```go
if os.Getenv("TF_DEBUG") == "1" {
    logger.SetOutput(io.MultiWriter(os.Stderr, file))  // Use stderr, not stdout
} else {
    logger.SetOutput(file)
}
```

**Files to Modify:**
- `internal/logx/logger.go`

**Acceptance Criteria:**
- ✅ Logs written ONLY to `tf-engine.log`
- ✅ Stdout contains ONLY command output (JSON or human text)
- ✅ Stderr can be used for errors/warnings if needed
- ✅ VBA receives clean stdout

---

## Issue 3: Inconsistent JSON Output

### Problem

Some commands produce no JSON output at all, or only text:

**check-heat command:**
```bash
$ tf-engine check-heat --add-r 75
(empty or text-only output)
```

### Solution

**Every command MUST return valid JSON** (when `--format json` is used):

At minimum:
```json
{
  "status": "ok"
}
```

Better:
```json
{
  "current_portfolio_heat": 0,
  "new_portfolio_heat": 75,
  "portfolio_cap": 400,
  "allow_trade": true
}
```

**Files to Check:**
- `internal/cli/check_heat.go`
- Any command with empty or text-only output

**Acceptance Criteria:**
- ✅ All commands return valid JSON (when `--format json`)
- ✅ JSON is parseable by standard parsers
- ✅ Schema documented in JSON_API_SPECIFICATION.md

---

## Issue 4: Inconsistent HTTP Response Structure

### Problem

Need to verify that HTTP endpoints return identical JSON to CLI commands.

### Solution

Test all HTTP endpoints and compare to CLI output.

**Process:**
1. Start HTTP server
2. Call each endpoint with identical params
3. Compare JSON responses
4. Document differences
5. Fix to achieve parity

**Acceptance Criteria:**
- ✅ CLI (with `--format json`) and HTTP return IDENTICAL JSON for same inputs
- ✅ Parity verified with automated tests
- ✅ Documented in JSON_API_SPECIFICATION.md

---

## Implementation Plan

### Step 1: Fix Logging (PRIORITY 1)
- Modify `internal/logx/logger.go` to write only to file
- Test all CLI commands
- Verify stdout is clean

**Estimated Time:** 15 minutes

### Step 2: Add `--format` Flag (PRIORITY 1)
- Add global flag to root command
- Modify each command to respect flag
- Ensure `--format json` outputs ONLY JSON
- Ensure `--format human` preserves current UX

**Estimated Time:** 1-2 hours

### Step 3: Fix Missing JSON Outputs (PRIORITY 2)
- Audit all commands for missing/incomplete JSON
- Add proper JSON responses
- Update documentation

**Estimated Time:** 30 minutes

### Step 4: Verify HTTP Parity (PRIORITY 2)
- Start server
- Test all endpoints
- Compare with CLI JSON
- Fix discrepancies

**Estimated Time:** 1 hour

### Step 5: Re-Capture JSON Examples (PRIORITY 3)
- Run capture script with `--format json`
- Update test-data/json-examples/
- Verify all examples are clean JSON

**Estimated Time:** 15 minutes

**Total Estimated Time:** 3-4 hours

---

## Testing Strategy

### Manual Testing

After each fix, test:
```bash
# Should output ONLY JSON (no logs, no text)
./tf-engine size --entry 180 --atr 1.5 --k 2 --method stock --format json

# Should output human-friendly text + JSON
./tf-engine size --entry 180 --atr 1.5 --k 2 --method stock --format human
```

### Automated Testing

Create test script:
```bash
#!/bin/bash
# Test that --format json produces clean, parseable JSON

# Test 1: No log lines in output
output=$(./tf-engine size --entry 180 --atr 1.5 --k 2 --method stock --format json)
if echo "$output" | grep -q '"level"'; then
    echo "FAIL: Log lines detected in JSON output"
    exit 1
fi

# Test 2: Valid JSON
if ! echo "$output" | python3 -m json.tool > /dev/null 2>&1; then
    echo "FAIL: Invalid JSON"
    exit 1
fi

echo "PASS: Clean JSON output"
```

### VBA Parsing Test

Create minimal VBA test:
```vba
Sub TestJSONParsing()
    Dim shell As Object
    Dim exec As Object
    Dim stdout As String

    Set shell = CreateObject("WScript.Shell")
    Set exec = shell.Exec("tf-engine.exe size --entry 180 --atr 1.5 --k 2 --method stock --format json")

    ' Wait for command to complete
    Do While exec.Status = 0
        DoEvents
    Loop

    ' Read stdout
    stdout = exec.StdOut.ReadAll()

    ' Should be parseable as JSON
    Dim json As Object
    Set json = ParseJSON(stdout)  ' Using VBA-JSON library

    Debug.Print "Shares: " & json("shares")
    Debug.Print "Risk: " & json("risk_dollars")
End Sub
```

---

## Success Criteria for M17-M18 Completion

Before proceeding to M19 (VBA Implementation):

- ✅ Logging writes ONLY to file (not stdout)
- ✅ All commands support `--format json` flag
- ✅ With `--format json`, stdout contains ONLY valid, parseable JSON
- ✅ HTTP endpoints return identical JSON to CLI
- ✅ All JSON schemas documented
- ✅ Test data captured with clean JSON examples
- ✅ VBA can parse stdout without string manipulation

---

## Action Items

- [ ] Fix `internal/logx/logger.go` (remove stdout from MultiWriter)
- [ ] Add `--format` flag to root command
- [ ] Modify checklist command to respect `--format`
- [ ] Modify cooldown commands to respect `--format`
- [ ] Modify heat command to respect `--format` and return JSON
- [ ] Modify timer command to respect `--format`
- [ ] Test all commands with `--format json`
- [ ] Start HTTP server and test all endpoints
- [ ] Create CLI/HTTP parity test suite
- [ ] Re-capture JSON examples with clean output
- [ ] Update JSON_API_SPECIFICATION.md with final schemas
- [ ] Create VBA parsing test to validate

**Blocked:** M19 (VBA Implementation) until all items complete

**Estimated Completion:** 3-4 hours of focused work

---

**Document Version:** 1.0
**Last Updated:** 2025-10-27
**Status:** Action Required
