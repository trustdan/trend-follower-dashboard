# Trading Engine v3 - JSON API Specification

**Purpose:** Document all JSON request/response schemas for CLI and HTTP endpoints.

**For:** M17-M18 (JSON Contracts & Validation) and M19 (VBA Implementation)

**Date:** 2025-10-27

---

## Overview

This document defines the exact JSON structures for all trading engine commands. Both CLI and HTTP transports MUST return identical JSON for the same inputs.

### General Patterns

**Success Responses:**
- CLI: Writes JSON to stdout (may include non-JSON messages - see known issues)
- HTTP: Returns JSON with `200 OK` status

**Error Responses:**
- CLI: Writes error message to stderr, exits with non-zero code
- HTTP: Returns JSON error with appropriate HTTP status code (400, 404, 500)

**Common Fields:**
- All timestamps use RFC3339Nano format: `"2025-10-27T19:59:27.950139119-05:00"`
- All monetary values are float64 (e.g., `75.0`, `177.5`)
- All percentages are stored as decimals (e.g., `0.0075` for 0.75%)

---

## 1. Position Sizing

### 1.1 Stock Method

**CLI Command:**
```bash
tf-engine size --entry 180 --atr 1.5 --k 2 --method stock
```

**HTTP Endpoint:**
```
POST /api/size
Content-Type: application/json
```

**Request (HTTP):**
```json
{
  "equity": 10000,
  "risk_pct": 0.0075,
  "entry": 180.0,
  "atr_n": 1.5,
  "k": 2,
  "method": "stock"
}
```

**Response (Both):**
```json
{
  "risk_dollars": 75,
  "stop_distance": 3,
  "initial_stop": 177,
  "shares": 25,
  "contracts": 0,
  "actual_risk": 75,
  "method": "stock"
}
```

**Fields:**
- `risk_dollars` (float): Maximum risk in dollars (Equity × RiskPct)
- `stop_distance` (float): Distance from entry to stop (K × ATR)
- `initial_stop` (float): Actual stop price (Entry - StopDistance)
- `shares` (int): Number of shares to trade (floor(RiskDollars / StopDistance))
- `contracts` (int): Always 0 for stocks
- `actual_risk` (float): Actual risk based on calculated shares
- `method` (string): "stock"

### 1.2 Options - Delta-ATR Method

**CLI Command:**
```bash
tf-engine size --entry 5.00 --atr 0.50 --k 2 --delta 0.50 --method opt-delta-atr
```

**Response:**
```json
{
  "risk_dollars": 75,
  "stop_distance": 1,
  "initial_stop": 4,
  "shares": 0,
  "contracts": 0,
  "actual_risk": 0,
  "method": "opt-delta-atr"
}
```

**Fields:** Same as stock method, but:
- `shares`: Always 0 for options
- `contracts`: Number of option contracts
- `method`: "opt-delta-atr"

### 1.3 Options - Max Loss Method

**CLI Command:**
```bash
tf-engine size --entry 5.00 --atr 0.50 --k 2 --maxloss 50 --method opt-maxloss
```

**Response:**
```json
{
  "risk_dollars": 75,
  "stop_distance": 0,
  "initial_stop": 0,
  "shares": 0,
  "contracts": 1,
  "actual_risk": 50,
  "method": "opt-maxloss"
}
```

---

## 2. Checklist Evaluation

### 2.1 GREEN Banner (0 missing)

**CLI Command:**
```bash
tf-engine checklist --ticker AAPL \
  --from-preset --trend-pass --liquidity-pass \
  --tv-confirm --earnings-ok --journal-ok
```

**HTTP Endpoint:**
```
POST /api/checklist
```

**Request (HTTP):**
```json
{
  "ticker": "AAPL",
  "from_preset": true,
  "trend_pass": true,
  "liquidity_pass": true,
  "tv_confirm": true,
  "earnings_ok": true,
  "journal_ok": true
}
```

**Response (Both):**
```json
{
  "banner": "GREEN",
  "missing_count": 0,
  "missing_items": [],
  "evaluation_timestamp": "2025-10-27T19:59:27.950139119-05:00",
  "allow_save": true
}
```

**Fields:**
- `banner` (string): "GREEN", "YELLOW", or "RED"
- `missing_count` (int): Number of unchecked items (0-6)
- `missing_items` (array of strings): Names of missing items
- `evaluation_timestamp` (string): RFC3339Nano timestamp when evaluated
- `allow_save` (bool): true only if banner is GREEN

**Side Effect:** GREEN banner starts a 2-minute impulse brake timer for this ticker.

### 2.2 YELLOW Banner (1 missing)

**Response:**
```json
{
  "banner": "YELLOW",
  "missing_count": 1,
  "missing_items": ["JournalOK"],
  "evaluation_timestamp": "0001-01-01T00:00:00Z",
  "allow_save": false
}
```

### 2.3 RED Banner (2+ missing)

**Response:**
```json
{
  "banner": "RED",
  "missing_count": 2,
  "missing_items": ["EarningsOK", "JournalOK"],
  "evaluation_timestamp": "0001-01-01T00:00:00Z",
  "allow_save": false
}
```

**Note:** `evaluation_timestamp` is zero value when `allow_save` is false.

---

## 3. Heat Management

### 3.1 Check Heat

**CLI Command:**
```bash
tf-engine check-heat --add-r 75
```

**HTTP Endpoint:**
```
GET /api/heat?add_r=75&bucket=Tech/Comm
```

**Response:**
```json
{
  "current_portfolio_heat": 0,
  "new_portfolio_heat": 75,
  "portfolio_cap": 400,
  "portfolio_status": "OK",
  "current_bucket_heat": 0,
  "new_bucket_heat": 75,
  "bucket_cap": 150,
  "bucket_status": "OK",
  "allow_trade": true
}
```

**Fields:**
- `current_portfolio_heat` (float): Sum of all open position risks
- `new_portfolio_heat` (float): Current + proposed trade risk
- `portfolio_cap` (float): Equity × HeatCap% (default 4%)
- `portfolio_status` (string): "OK" or "EXCEEDED"
- `current_bucket_heat` (float): Sum of risks in this sector bucket
- `new_bucket_heat` (float): Current bucket + proposed trade risk
- `bucket_cap` (float): Equity × BucketHeatCap% (default 1.5%)
- `bucket_status` (string): "OK" or "EXCEEDED"
- `allow_trade` (bool): true only if both statuses are "OK"

**Note:** This endpoint has inconsistent output format in CLI - see known issues.

---

## 4. Candidates

### 4.1 List Candidates

**CLI Command:**
```bash
tf-engine list-candidates --date 2025-10-27
```

**HTTP Endpoint:**
```
GET /api/candidates?date=2025-10-27
```

**Response:**
```json
{
  "candidates": [
    {
      "id": 1,
      "ticker": "AAPL",
      "date": "2025-10-27",
      "preset": "TEST",
      "preset_id": 1,
      "sector": "",
      "bucket": ""
    },
    {
      "id": 2,
      "ticker": "MSFT",
      "date": "2025-10-27",
      "preset": "TEST",
      "preset_id": 1,
      "sector": "",
      "bucket": ""
    }
  ],
  "count": 2,
  "date": "2025-10-27"
}
```

### 4.2 Check Candidate (Found)

**CLI Command:**
```bash
tf-engine check-candidate --ticker AAPL
```

**Response:**
```json
{
  "ticker": "AAPL",
  "date": "2025-10-27",
  "found": true
}
```

### 4.3 Check Candidate (Not Found)

**Response:**
```json
{
  "ticker": "TSLA",
  "date": "2025-10-27",
  "found": false
}
```

---

## 5. Settings

### 5.1 Get All Settings

**CLI Command:**
```bash
tf-engine get-settings
```

**HTTP Endpoint:**
```
GET /api/settings
```

**Response:**
```json
{
  "Equity_E": "10000",
  "RiskPct_r": "0.0075",
  "StopMultiple_K": "2",
  "HeatCap_H_pct": "0.04",
  "BucketHeatCap_pct": "0.015"
}
```

**Note:** Values are returned as strings, not numbers, to preserve precision.

---

## 6. Impulse Brake Timer

### 6.1 Check Timer (No Timer)

**CLI Command:**
```bash
tf-engine check-timer --ticker AAPL
```

**Response:**
```json
{
  "ticker": "AAPL",
  "timer_active": false,
  "elapsed_seconds": 0,
  "remaining_seconds": 0,
  "brake_cleared": false
}
```

### 6.2 Check Timer (Active)

**Response:**
```json
{
  "ticker": "AAPL",
  "timer_active": true,
  "started_at": "2025-10-27T19:59:27.950139119-05:00",
  "elapsed_seconds": 45,
  "remaining_seconds": 75,
  "brake_cleared": false
}
```

### 6.3 Check Timer (Cleared)

**Response:**
```json
{
  "ticker": "AAPL",
  "timer_active": true,
  "started_at": "2025-10-27T19:57:27.950139119-05:00",
  "elapsed_seconds": 125,
  "remaining_seconds": 0,
  "brake_cleared": true
}
```

**Fields:**
- `timer_active` (bool): Timer exists for this ticker
- `started_at` (string): RFC3339Nano timestamp when timer started
- `elapsed_seconds` (int): Seconds since timer started
- `remaining_seconds` (int): Seconds until brake clears (0 if cleared)
- `brake_cleared` (bool): true if elapsed >= 120 seconds

---

## 7. Bucket Cooldowns

### 7.1 Check Cooldown (Not Active)

**CLI Command:**
```bash
tf-engine check-cooldown --bucket Energy
```

**Response:**
```json
{
  "bucket": "Energy",
  "in_cooldown": false
}
```

### 7.2 Check Cooldown (Active)

**Response:**
```json
{
  "bucket": "Tech/Comm",
  "in_cooldown": true,
  "started_at": "2025-10-27T19:59:30.000000000-05:00",
  "expires_at": "2025-10-28T19:59:30.000000000-05:00",
  "remaining_hours": 24.0,
  "reason": "Manual trigger"
}
```

### 7.3 List Cooldowns (Empty)

**CLI Command:**
```bash
tf-engine list-cooldowns
```

**Response:**
```json
{
  "cooldowns": [],
  "count": 0
}
```

### 7.4 List Cooldowns (With Data)

**Response:**
```json
{
  "cooldowns": [
    {
      "bucket": "Tech/Comm",
      "started_at": "2025-10-27T19:59:30.000000000-05:00",
      "expires_at": "2025-10-28T19:59:30.000000000-05:00",
      "remaining_hours": 24.0,
      "reason": "Manual trigger"
    }
  ],
  "count": 1
}
```

---

## 8. Positions

### 8.1 List Positions (Empty)

**CLI Command:**
```bash
tf-engine list-positions
```

**Response:**
```json
{
  "positions": [],
  "count": 0
}
```

### 8.2 List Positions (With Data)

**Response:**
```json
{
  "positions": [
    {
      "id": 1,
      "ticker": "AAPL",
      "bucket": "Tech/Comm",
      "open_date": "2025-10-27",
      "entry_price": 180.0,
      "current_stop": 177.0,
      "units_open": 25,
      "total_open_r": 75.0,
      "status": "OPEN"
    }
  ],
  "count": 1,
  "total_heat": 75.0
}
```

---

## 9. Save Decision

### 9.1 Save Decision (Success)

**CLI Command:**
```bash
tf-engine save-decision \
  --ticker AAPL \
  --entry 180 \
  --atr 1.5 \
  --k 2 \
  --method stock \
  --shares 25 \
  --risk-dollars 75 \
  --bucket "Tech/Comm"
```

**HTTP Endpoint:**
```
POST /api/decision
```

**Request (HTTP):**
```json
{
  "ticker": "AAPL",
  "entry": 180.0,
  "atr_n": 1.5,
  "k": 2,
  "method": "stock",
  "shares": 25,
  "contracts": 0,
  "risk_dollars": 75.0,
  "bucket": "Tech/Comm",
  "checklist": {
    "from_preset": true,
    "trend_pass": true,
    "liquidity_pass": true,
    "tv_confirm": true,
    "earnings_ok": true,
    "journal_ok": true
  }
}
```

**Response (Success):**
```json
{
  "accepted": true,
  "decision_id": 1,
  "timestamp": "2025-10-27T20:01:45.123456789-05:00",
  "gates_passed": {
    "banner_green": true,
    "in_candidates": true,
    "impulse_brake_cleared": true,
    "bucket_not_in_cooldown": true,
    "heat_caps_ok": true
  }
}
```

### 9.2 Save Decision (Rejected - Banner Not Green)

**Response (Error):**
```json
{
  "accepted": false,
  "reason": "Banner must be GREEN (current: YELLOW)",
  "gates_passed": {
    "banner_green": false,
    "in_candidates": true,
    "impulse_brake_cleared": false,
    "bucket_not_in_cooldown": true,
    "heat_caps_ok": true
  }
}
```

### 9.3 Save Decision (Rejected - Impulse Brake)

**Response (Error):**
```json
{
  "accepted": false,
  "reason": "2-minute impulse brake not cleared (45 seconds remaining)",
  "gates_passed": {
    "banner_green": true,
    "in_candidates": true,
    "impulse_brake_cleared": false,
    "bucket_not_in_cooldown": true,
    "heat_caps_ok": true
  }
}
```

### 9.4 Save Decision (Rejected - Not in Candidates)

**Response (Error):**
```json
{
  "accepted": false,
  "reason": "Ticker AAPL not in today's candidates",
  "gates_passed": {
    "banner_green": true,
    "in_candidates": false,
    "impulse_brake_cleared": true,
    "bucket_not_in_cooldown": true,
    "heat_caps_ok": true
  }
}
```

### 9.5 Save Decision (Rejected - Heat Cap)

**Response (Error):**
```json
{
  "accepted": false,
  "reason": "Portfolio heat would exceed cap by $25.00 (new: $425, cap: $400)",
  "gates_passed": {
    "banner_green": true,
    "in_candidates": true,
    "impulse_brake_cleared": true,
    "bucket_not_in_cooldown": true,
    "heat_caps_ok": false
  }
}
```

### 9.6 Save Decision (Rejected - Bucket Cooldown)

**Response (Error):**
```json
{
  "accepted": false,
  "reason": "Bucket Tech/Comm is in cooldown (expires 2025-10-28 19:59)",
  "gates_passed": {
    "banner_green": true,
    "in_candidates": true,
    "impulse_brake_cleared": true,
    "bucket_not_in_cooldown": false,
    "heat_caps_ok": true
  }
}
```

---

## 10. Error Responses

### 10.1 Validation Error (Invalid Entry Price)

**CLI Command:**
```bash
tf-engine size --entry -180 --atr 1.5 --k 2 --method stock
```

**CLI Output (stderr):**
```
Error: entry price must be positive (got -180.000000)
```

**CLI Exit Code:** Non-zero

**HTTP Response:**
```
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "error": "entry price must be positive (got -180.000000)"
}
```

### 10.2 Missing Required Parameters

**CLI Output (stderr):**
```
Error: required flag(s) "entry", "atr" not set
```

**HTTP Response:**
```
HTTP/1.1 400 Bad Request

{
  "error": "missing required parameter: entry"
}
```

### 10.3 Invalid Method

**CLI Output (stderr):**
```
Error: invalid sizing method: invalid (must be stock, opt-delta-atr, or opt-maxloss)
```

**HTTP Response:**
```
HTTP/1.1 400 Bad Request

{
  "error": "invalid sizing method: invalid (must be stock, opt-delta-atr, or opt-maxloss)"
}
```

---

## Known Issues (To Fix Before M19)

### Issue 1: Mixed JSON and Text Output in CLI

**Problem:** Several CLI commands output human-readable text along with JSON to stdout.

**Examples:**
- `checklist` command outputs: `⏱️ Impulse brake timer started\n  Wait 2 minutes before saving decision\n{JSON}`
- `cooldown` commands output: `✓ Bucket Energy is NOT in cooldown\n{JSON}`
- `heat-check` command outputs text without JSON

**Impact:** VBA cannot reliably parse stdout when it contains mixed content.

**Solution Options:**
1. **Recommended:** Add `--format json` flag to force JSON-only output (default remains human-friendly)
2. Alternative: Always output pure JSON to stdout, human text to stderr
3. Alternative: Separate commands for programmatic use (e.g., `tf-engine api size` vs `tf-engine size`)

**Priority:** HIGH - Blocks M19 VBA implementation

### Issue 2: Logging to Stdout

**Problem:** `logx` package writes logs to both stdout and file (line 45: `io.MultiWriter(os.Stdout, file)`).

**Impact:** VBA receives log lines mixed with JSON output, making parsing unreliable.

**Solution:** Change logger to write ONLY to file, not stdout. Stderr can be used for user-facing errors.

**Priority:** HIGH - Blocks M19 VBA implementation

### Issue 3: Missing JSON Output for Some Commands

**Problem:** Commands like `check-heat` produce no JSON output in some cases.

**Impact:** VBA cannot parse empty responses reliably.

**Solution:** Always return valid JSON, even if it's just `{"status": "ok"}`.

**Priority:** MEDIUM

---

## VBA Parsing Requirements

For M19 (VBA implementation), VBA needs to:

1. **Call CLI command** via `WScript.Shell.Exec`
2. **Capture stdout** as a string
3. **Parse JSON string** into VBA objects/dictionaries
4. **Handle errors** from stderr or exit codes

**Critical Requirements:**
- ✅ Stdout MUST contain ONLY valid JSON (no mixed text)
- ✅ JSON MUST be on a single line OR VBA must handle multi-line
- ✅ Errors SHOULD go to stderr, not stdout
- ✅ Exit code SHOULD be 0 for success, non-zero for errors

**Current Status:**
- ❌ Mixed text/JSON output in several commands
- ❌ Logs polluting stdout
- ⚠️ Some commands have no JSON output

**Action Required Before M19:**
- Fix output format issues (Issue 1 & 2)
- Test clean JSON parsing from VBA

---

## Next Steps for M18

1. ✅ Capture all CLI JSON examples (DONE)
2. ✅ Document JSON schemas (DONE)
3. ⏳ Test HTTP endpoints for parity
4. ⏳ Create parity test suite
5. ⏳ Fix output format issues
6. ⏳ Re-capture clean JSON examples

---

## Appendix: HTTP Endpoint Summary

| Command | HTTP Method | Endpoint | Request Body | Response |
|---------|-------------|----------|--------------|----------|
| size | POST | /api/size | SizingRequest | SizingResponse |
| checklist | POST | /api/checklist | ChecklistRequest | ChecklistResponse |
| check-heat | GET | /api/heat?add_r=X&bucket=Y | - | HeatResponse |
| save-decision | POST | /api/decision | DecisionRequest | DecisionResponse |
| list-candidates | GET | /api/candidates?date=YYYY-MM-DD | - | CandidatesListResponse |
| get-settings | GET | /api/settings | - | SettingsResponse |
| check-timer | GET | /api/timer?ticker=AAPL | - | TimerResponse |
| check-cooldown | GET | /api/cooldown?bucket=Tech/Comm | - | CooldownResponse |
| list-cooldowns | GET | /api/cooldowns | - | CooldownsListResponse |
| list-positions | GET | /api/positions | - | PositionsListResponse |

---

**Document Version:** 1.0
**Last Updated:** 2025-10-27
**Status:** In Progress (M17-M18)
