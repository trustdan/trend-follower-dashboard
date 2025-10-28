# Trading Engine v3 — Step‑by‑Step Development Plan

**Purpose:** Rebuild the Excel trading workflow around a fast, testable backend (Go), with Excel as a thin UI. This plan merges your alternate v3 concept (CLI-first) with the best parts of the current v2 system (docs, logging, and UX patterns).

> Architecture in one line: **Engine-first** (Go), **CLI by default**, **HTTP optional**, **SQLite for truth**, **Excel as a thin client**.

---

## Development Environment Strategy

**Critical Constraint:** Development happens in **Linux/WSL** (fast iteration with Claude Code), but execution testing happens in **Windows** (Excel + VBA, manual, slower).

**Solution:** Minimize Windows testing burden by:
1. **M1-M16:** Develop and test entire Go engine in Linux (automated tests, no VBA)
2. **M17-M18:** Define and validate JSON contracts in Linux (engine outputs correct JSON)
3. **M19:** Write VBA as text files (`.bas` exports) in Linux; code against validated JSON schemas
4. **M20:** Create Windows integration package in Linux (batch scripts, test procedures, sample data)
5. **M21:** Manual Windows testing only (VBA execution + Excel UI integration, ~45 min if prep is solid)

**Key Insight:** If JSON contracts are correct and VBA is thin (just shell exec + parse JSON), Windows testing should succeed on first try or need only minor tweaks.

---

## 0) Goals & Non‑Goals

### Goals
- Move **all trading logic** (sizing, checklist → banner, heat caps, hard gates, FINVIZ import) into a compiled backend.
- Keep Excel **familiar**: same sheet layout, dropdowns, 6 checkboxes, buttons.
- Make behavior **provable** with BDD (Gherkin), unit & integration tests.
- Ship a **single Windows .exe** (CLI) with an optional `server` mode (HTTP/JSON) using the **same handlers**.
- Centralize **state** in **SQLite**, not in cells/VBA variables.

### Non‑Goals (v3.0)
- No broker API trading, no live market data stream.
- No Office‑JS rewrite; VBA remains a thin bridge.
- No web dashboard (leave for v3.1+ via HTTP parity).

---

## 1) Technical Decisions (locked for v3.0)

- **Language:** Go (fast bring‑up, single exe, mature stdlib).
- **Transport:** **CLI by default** (`tf-engine.exe <command> ...`), **HTTP optional** (`tf-engine.exe server --listen 127.0.0.1:18888`). Both return **identical JSON**.
- **Storage:** SQLite (one file `trading.db`), with migrations.
- **Logging:** Structured JSON logs with **correlation IDs** propagated from Excel.
- **OS:** Windows 10/11 (Excel desktop).

> Rationale: CLI-first is the simplest, most reliable path; HTTP later unlocks a web UI without refactoring Excel.

---

## 2) Repository Layout

```
excel-trading-dashboard/
├─ engine/
│  ├─ cmd/tf-engine/        # main (CLI + server subcommand)
│  ├─ internal/api/         # HTTP handlers (share domain with CLI)
│  ├─ internal/cli/         # CLI commands (share domain with HTTP)
│  ├─ internal/domain/      # sizing, checklist, heat, decisions
│  ├─ internal/scrape/      # finviz scraper (+ manual normalize)
│  ├─ internal/storage/     # sqlite (migrations, queries)
│  └─ internal/logx/        # logging helpers (corr_id)
├─ excel/
│  ├─ vba/
│  │  ├─ TFEngine.bas       # Engine communication (text export, version controlled)
│  │  ├─ TFHelpers.bas      # JSON parsing, error handling
│  │  └─ TFTypes.bas        # Type definitions
│  └─ TrendFollowing_TradeEntry.xlsm  # Workbook template
├─ windows/
│  ├─ tf-engine.exe         # Windows binary (cross-compiled in M20)
│  ├─ windows-import-vba.bat
│  ├─ windows-init-database.bat
│  ├─ run-tests.bat
│  ├─ WINDOWS_TESTING.md
│  └─ test-data/            # Sample JSON responses for testing
├─ docs/
│  ├─ start-here.md
│  ├─ quick-start.md
│  ├─ architecture.md
│  ├─ troubleshooting.md
│  └─ bdd/ *.feature
├─ features/                # BDD Gherkin scenarios
└─ .github/workflows/ci.yml # Linux CI only (Windows testing manual)
```

**Development Workflow:**
- **Linux/WSL:** All Go code, tests, VBA text files, build scripts
- **Windows:** Manual execution testing only (via M21 testing procedure)
- **Version Control:** VBA as `.bas` text exports (not binary `.xlsm` with embedded VBA)

---

## 3) Data Model (SQLite)

**Tables (MVP):**

- `settings(key TEXT PRIMARY KEY, value TEXT)` — Equity_E, RiskPct_r, HeatCap_H_pct, BucketHeatCap_pct, StopMultiple_K, etc.
- `presets(id INTEGER PK, name TEXT UNIQUE, query_string TEXT, active INTEGER)`
- `candidates(id INTEGER PK, date DATE, ticker TEXT, preset_id INT, sector TEXT, bucket TEXT)`
- `decisions(id INTEGER PK, timestamp DATETIME, ticker TEXT, preset_id INT, bucket TEXT, entry_price REAL, atr_n REAL, k_multiple INT, method TEXT, risk_dollars REAL, shares INT, contracts INT, banner TEXT, checklist_json TEXT, hard_gates_json TEXT, accepted INT, reason TEXT)`
- `positions(id INTEGER PK, ticker TEXT, bucket TEXT, open_date DATE, units_open INT, total_open_r REAL, status TEXT)`

**Notes:**
- Keep history in `decisions`; open/close rolls into `positions` only if you choose to track that in v3.0.
- JSON columns for checklist/gates to preserve traceability.

---

## 4) Engine Surface (Shared JSON for CLI & HTTP)

### Core commands / endpoints
- `size` / `POST /sizing` — position sizing (stocks; options: delta‑ATR, max‑loss)
- `checklist` / `POST /evaluate` — 6‑item → GREEN/YELLOW/RED (+ missing list)
- `heat` / `GET /heat?add_r=...&bucket=...` — portfolio/bucket heat preview, caps
- `save-decision` / `POST /decision` — enforces **5 hard gates**; persists
- `scrape-finviz`, `import-candidates`, `list-presets`
- `get-settings`, `set-setting`

### JSON (illustrative)
```jsonc
// /sizing request
{
  "equity": 10000,
  "risk_pct": 0.0075,
  "entry": 180.0,
  "atr_n": 1.5,
  "k": 2,
  "method": "stock",    // "opt-delta-atr" | "opt-maxloss"
  "delta": 0.30,        // optional
  "max_loss": 50        // optional
}

// /sizing response
{
  "risk_dollars": 75.0,
  "stop_distance": 3.0,
  "initial_stop": 177.0,
  "shares": 25,
  "contracts": 0
}
```

---

## 5) Excel Integration (thin & resilient)

- Default bridge: **CLI** via `WScript.Shell.Exec`, reading **stdout** (JSON), capturing **stderr** for user‑friendly messages.
- Optional bridge: **HTTP** via `MSXML2.ServerXMLHTTP.6.0` (toggle).
- Keep current UI: dropdowns (B5‑B8), checkboxes (C20‑C25), buttons (Evaluate / Recalc / Save / Import / Open FINVIZ).
- Always pass a generated **corr_id** into each engine call; show it in status cells and write to `TradingSystem_Debug.log` for cross‑reference.

**VBA public surface (stubs):**
- `Engine_ImportCandidates(presetName) As Variant()`
- `Engine_Evaluate(ticker, entry, atr_n, k, checks As Collection) As EngineEvalResult`
- `Engine_RecalcSizing(params As EngineSizingParams) As EngineSizingResult`
- `Engine_SaveDecision(dto As EngineDecisionDto) As EngineSaveResult`

---

## 6) Step‑by‑Step Development Tasks

### Phase A — Bootstrap (Days 1‑3)
1. **Repo/init**
   - Initialize Go module; add basic folders.
   - Vendor dependencies (`go mod tidy`).

2. **Storage**
   - Add SQLite driver.
   - Implement migrations & bootstrap `settings` with defaults (Equity_E=10000, RiskPct_r=0.0075, HeatCap_H_pct=0.04, BucketHeatCap_pct=0.015, StopMultiple_K=2).

3. **Logging**
   - Implement structured logger with `corr_id` propagation.
   - Single rotating log file next to the exe.

4. **CLI scaffold**
   - `tf-engine.exe init` → creates DB if missing.
   - `tf-engine.exe size ...` → implements stock sizing first.

**Exit criteria:** `size` works with unit tests; DB initializes; logs write correctly.

---

### Phase B — Core Domain (Days 4‑7)
5. **Sizing variants**
   - Implement **opt‑delta‑atr**, **opt‑maxloss**; add unit tests for examples.

6. **Checklist → banner**
   - Input: 6 booleans; Output: banner + missing list.
   - Start an **evaluation timestamp** used by the 2‑minute impulse gate.

7. **Heat**
   - Portfolio & bucket preview given a new trade **R**; compute caps from settings.

8. **Save decision**
   - Enforce **5 hard gates**:
     - Banner must be GREEN
     - Ticker in **today’s** candidates
     - **2‑minute impulse** elapsed since evaluation
     - Bucket not in cooldown
     - Heat caps not exceeded
   - Persist to `decisions`.

**Exit criteria:** Unit tests pass for sizing, checklist, heat, and decisions; CLI returns correct JSON.

---

### Phase C — FINVIZ & Data Paths (Days 8‑10)
9. **Import candidates**
   - `import-candidates --tickers "AAPL,MSFT,NVDA" --preset TF_BREAKOUT_LONG`

10. **FINVIZ scraper**
   - `scrape-finviz --query "<screen-string>"` with pagination, normalization, rate‑limit.
   - Return `tickers[]` + count + date.
   - Manual fallback (Excel prompts paste) remains unchanged on the UI side.

**Exit criteria:** Candidates appear in DB; manual paste path tested; scraper robust to pagination and rate limit.

---

### Phase D — HTTP Parity & Excel Wiring (Days 11‑16)

**Environment Note:** Development happens in **Linux/WSL**; VBA execution testing happens in **Windows** (manual, slower). Minimize Windows test iterations by validating everything possible in Linux first.

**M17-M18: Engine JSON Contracts (Linux, Days 11-12)**
11. **Define JSON schemas**
   - Document exact request/response schemas for all CLI commands
   - Create example JSON files for each command (success & error cases)
   - Implement JSON schema validation in engine
   - Unit test: engine produces valid JSON for all scenarios

12. **HTTP server with parity**
   - `tf-engine.exe server --listen 127.0.0.1:18888`
   - Handlers wrap the same domain logic as CLI
   - **Transport parity tests:** CLI vs HTTP return identical JSON for same inputs
   - Test both transports thoroughly in Linux

**Exit criteria (M17-M18):** All JSON schemas documented; engine produces valid, tested JSON; CLI and HTTP parity verified in Linux.

---

**M19: VBA Implementation (Linux, Days 13-14)**
13. **Create VBA modules as text exports**
   - Write `TFEngine.bas` (shell execution, JSON parsing)
   - Write `TFHelpers.bas` (error handling, logging with corr_id)
   - Write `TFTypes.bas` (type definitions for JSON responses)
   - Code against validated JSON schemas from M17-M18
   - Keep VBA **thin**: just shell exec → parse JSON → return
   - Export as `.bas` files (version-controllable text)

**Exit criteria (M19):** VBA modules written as `.bas` files; syntax valid; logic matches JSON contracts. **Execution testing deferred to M21.**

---

**M20: Windows Integration Package (Linux, Days 15-16)**
14. **Create Windows deployment artifacts**
   - Cross-compile `tf-engine.exe` for Windows (GOOS=windows GOARCH=amd64)
   - Create `windows-import-vba.bat` (imports `.bas` files into Excel workbook)
   - Create `windows-init-database.bat` (initializes SQLite DB)
   - Create `WINDOWS_TESTING.md` (step-by-step manual test guide)
   - Create test data files (`test-data/*.json` with sample responses)
   - Create Excel workbook template with placeholder VBA modules
   - Create automated test runner (`run-tests.bat`) that calls VBA test functions

**Exit criteria (M20):** Complete Windows integration package ready; all artifacts created in Linux; ready for manual testing.

---

**M21: Windows Integration Validation (Windows, Days 17-18)**
15. **Manual Windows testing** (see detailed M21 plan in §6.1 below)
   - Phase 1: Pre-test setup (~10 min)
   - Phase 2: Smoke tests (~5 min)
   - Phase 3: VBA unit tests (~10 min)
   - Phase 4: Integration tests (~15 min)
   - Phase 5: Issue reporting (if needed)
   - Phase 6: Final validation (~5 min)

**Exit criteria (M21):** All Windows integration tests pass; VBA correctly calls engine; JSON parsing works; Excel UI functions end-to-end.

---

### Phase E — Hardening & Release (Days 19‑22)
14. **BDD features** (see §8) — full suite green in CI.
15. **Error surfaces** — map engine errors to friendly Excel messages.
16. **Packaging** — zip with exe, DB bootstrap, workbook, docs.
17. **Docs** — update start‑here, quick‑start, troubleshooting, architecture.

**Exit criteria:** Zipped deliverable works on a clean Windows box in < 5 min.

---

## 6.1) Detailed M21 Plan: Windows Integration Validation

**Goal:** Validate VBA ↔ Go engine integration in Windows with minimal manual testing iterations.

**Prerequisites:** M17-M20 completed; Windows package ready with:
- `tf-engine.exe` (Windows binary)
- `windows-import-vba.bat`, `windows-init-database.bat`, `run-tests.bat`
- `WINDOWS_TESTING.md` (step-by-step guide)
- VBA modules as `.bas` exports
- Test data files (`test-data/*.json`)
- Excel template workbook

### Phase 1: Pre-Test Setup (Windows, ~10 min)

**Setup Steps:**
1. Copy `windows/` folder to Windows machine
2. Run `windows-init-database.bat` → Creates `trading.db` with schema
3. Open `TradingPlatform.xlsm` in Excel
4. Run `windows-import-vba.bat` → Imports VBA modules from `.bas` files
5. Enable macros in Excel (Trust Center settings if needed)
6. Open VBA editor (Alt+F11) and verify modules loaded

**Success Criteria:**
- ✅ `trading.db` exists with correct schema
- ✅ VBA modules visible in VBA editor
- ✅ No import errors or warnings

### Phase 2: Smoke Tests (Windows, ~5 min)

**Test CLI and engine basics before VBA testing**

**Test 2.1: Engine version**
```batch
tf-engine.exe --version
```
Expected: Version string printed (e.g., "tf-engine v3.0.0")

**Test 2.2: Database initialization**
```batch
tf-engine.exe init
```
Expected: "Database initialized" or "Database already initialized"

**Test 2.3: Position sizing command**
```batch
tf-engine.exe size --entry 180 --atr 1.5 --k 2 --method stock
```
Expected: Valid JSON with `shares`, `initial_stop`, `risk_dollars`

**Success Criteria:**
- ✅ All 3 commands execute without errors
- ✅ JSON output is well-formed
- ✅ Values match expected ranges

### Phase 3: VBA Unit Tests (Windows, ~10 min)

**Test VBA modules in isolation before full integration**

Excel workbook contains a "VBA Tests" worksheet with test buttons.

**Test 3.1: Shell Execution**
- Click button: "Test: Shell Execution"
- VBA function: `TFEngine.ExecuteCommand("--version")`
- Validates: Can spawn process, capture stdout, handle exit code
- Sheet cell displays: ✅ PASS or ❌ FAIL with error

**Test 3.2: JSON Parsing (Success Response)**
- Click button: "Test: Parse Success JSON"
- VBA function: `TFHelpers.ParseSizingJSON(sampleJson)`
- Uses: `test-data/size-response-success.json`
- Validates: Extracts `shares`, `initial_stop`, `risk_dollars` correctly
- Sheet cell displays: ✅ PASS or ❌ FAIL

**Test 3.3: JSON Parsing (Error Response)**
- Click button: "Test: Parse Error JSON"
- VBA function: `TFHelpers.ParseErrorJSON(errorJson)`
- Uses: `test-data/error-response.json`
- Validates: Extracts error message, displays user-friendly text
- Sheet cell displays: ✅ PASS or ❌ FAIL

**Success Criteria:**
- ✅ All 3 VBA unit tests show PASS
- ✅ No VBA runtime errors (check Immediate window)
- ✅ Test results logged to `TradingSystem_Debug.log`

### Phase 4: Integration Tests (Windows, ~15 min)

**Test complete workflows through Excel UI**

**Test 4.1: Position Sizing**
1. Navigate to "Position Sizing" worksheet
2. Enter test values:
   - Entry = $180.00
   - ATR (N) = $1.50
   - K = 2
3. Click button: "Calculate"
4. VBA calls: `tf-engine.exe size --entry 180 --atr 1.5 --k 2 --method stock`
5. VBA parses JSON response
6. Sheet auto-fills results

**Expected Results:**
- Shares = 25
- Initial Stop = $177.00
- Risk $ = $75.00
- Status cell shows: "✅ Calculated successfully" with corr_id

**Test 4.2: Checklist Validation (GREEN)**
1. Navigate to "Checklist" worksheet
2. Check all 6 checklist items
3. Click button: "Evaluate"
4. VBA calls: `tf-engine.exe checklist --checks [...all true...]`
5. VBA parses JSON response
6. Banner cell updates, timer starts

**Expected Results:**
- Banner = GREEN (cell background color green)
- Missing count = 0
- Timer cell shows: "Impulse timer: 2:00 remaining"
- Status: "✅ Ready to save (after timer)"

**Test 4.3: Checklist Validation (YELLOW)**
1. Navigate to "Checklist" worksheet
2. Check 5/6 items (leave one unchecked)
3. Click button: "Evaluate"
4. VBA calls: `tf-engine.exe checklist --checks [...one false...]`
5. VBA parses JSON response
6. Banner updates, timer does NOT start

**Expected Results:**
- Banner = YELLOW (cell background color yellow)
- Missing count = 1
- Missing items listed
- Timer cell: (blank - not started)
- Status: "⚠️ Caution - 1 item missing"

**Test 4.4: Heat Management**
1. Navigate to "Heat Check" worksheet
2. Enter: New trade risk $ = $75
3. Click button: "Check Heat"
4. VBA calls: `tf-engine.exe heat --add-r 75`
5. VBA parses JSON response
6. Sheet displays heat breakdown

**Expected Results:**
- Current portfolio heat: (calculated from open positions)
- New portfolio heat: Current + $75
- Portfolio cap: $400 (4% of $10k)
- Bucket heat: (calculated for relevant bucket)
- Bucket cap: $150 (1.5% of $10k)
- Status: "✅ Heat OK" or "❌ Exceeds cap by $X"

**Test 4.5: Save Decision (Happy Path)**
1. Navigate to "Trade Entry" worksheet
2. Complete full trade entry (use prior test results):
   - Ticker: AAPL
   - Entry: $180
   - Shares: 25
   - Banner: GREEN (from earlier evaluation)
   - Timer: (wait 2 minutes or override for testing)
3. Click button: "Save Decision"
4. VBA calls: `tf-engine.exe save-decision --ticker AAPL [...all params...]`
5. VBA parses JSON response
6. Sheet confirms save

**Expected Results:**
- Status: "✅ Decision saved" with decision ID and timestamp
- Sheet clears input fields (ready for next trade)
- Confirmation logged with corr_id
- Database query: `SELECT * FROM decisions ORDER BY id DESC LIMIT 1` shows saved record

**Test 4.6: Save Decision (Gate Rejection - Banner)**
1. Navigate to "Trade Entry" worksheet
2. Complete trade entry with Banner = YELLOW (not GREEN)
3. Click button: "Save Decision"
4. VBA calls: `tf-engine.exe save-decision [...YELLOW banner...]`
5. Engine returns error JSON
6. VBA parses and displays error

**Expected Results:**
- Status: "❌ REJECTED: Banner must be GREEN"
- Sheet retains input data (not cleared)
- Error logged with corr_id
- Database: No new record in `decisions` table

**Test 4.7: Save Decision (Gate Rejection - Impulse Timer)**
1. Navigate to "Trade Entry" worksheet
2. Complete trade with GREEN banner
3. Immediately click "Save Decision" (before 2 min elapsed)
4. VBA calls: `tf-engine.exe save-decision [...]`
5. Engine returns error JSON
6. VBA displays error

**Expected Results:**
- Status: "❌ REJECTED: 2-minute impulse timer not elapsed (XYZ seconds remaining)"
- Sheet retains input data
- Error logged with corr_id

**Success Criteria:**
- ✅ All 7 integration tests pass
- ✅ JSON parsing works correctly for all response types
- ✅ Excel UI updates appropriately (colors, values, status messages)
- ✅ Errors displayed clearly with actionable messages
- ✅ Correlation IDs appear in both Excel status cells and log files
- ✅ Database state matches expectations

### Phase 5: Issue Reporting (If Tests Fail)

**If any test fails, document using this template:**

```
Test ID: 4.2 (Checklist Validation - YELLOW)
Status: ❌ FAILED

Expected:
- Banner cell background = YELLOW
- Missing count = 1

Actual:
- Banner cell background = RED
- Missing count = 2

Error Details:
- VBA error: "Type mismatch in ParseChecklistJSON, line 42"
- Engine JSON output:
  {
    "banner": "YELLOW",
    "missing_count": 1,
    "missing_items": ["Higher high"]
  }
- Correlation ID: abc-123-def-456

Screenshots: [attach Excel screenshot showing issue]

Relevant Log Excerpt:
[paste from TradingSystem_Debug.log]
```

**Developer Fix Process:**
1. Developer fixes issue in Linux (fast iteration)
2. Developer updates affected files in `windows/` folder
3. You copy updated files to Windows test environment
4. Re-run only the affected test(s) (targeted re-test)
5. If pass, continue; if fail, repeat

### Phase 6: Final Validation (Windows, ~5 min)

**Once all individual tests pass, run automated test suite:**

```batch
run-tests.bat
```

This batch script:
1. Runs all smoke tests
2. Calls VBA test functions programmatically (if Excel supports COM automation)
3. Generates test report

**Expected Output:**
```
========================================
 Windows Integration Test Report
========================================
Date: 2025-10-27 14:30:00
Environment: Windows 11, Excel 2021

[SMOKE TESTS]
[✓] Engine version
[✓] Database initialization
[✓] Position sizing CLI

[VBA UNIT TESTS]
[✓] Shell execution
[✓] JSON parsing (success)
[✓] JSON parsing (error)

[INTEGRATION TESTS]
[✓] Position sizing workflow
[✓] Checklist validation (GREEN)
[✓] Checklist validation (YELLOW)
[✓] Heat management
[✓] Save decision (happy path)
[✓] Save decision (gate: banner)
[✓] Save decision (gate: impulse timer)

========================================
RESULT: 13/13 PASSED ✅
========================================

Test results saved to: test-results.txt
Logs available at: TradingSystem_Debug.log
Correlation IDs: [list of all corr_ids used]
```

**Success Criteria:**
- ✅ All tests show PASSED
- ✅ Test report saved to `test-results.txt`
- ✅ No errors or warnings in logs
- ✅ Ready to proceed to Phase E (Hardening & Release)

### M21 Deliverables

**Created by Developer (in M20, used in M21):**
- `WINDOWS_TESTING.md` - Complete testing guide
- `run-tests.bat` - Automated test runner
- `test-data/` - Sample JSON files for VBA tests
- VBA test functions in workbook
- Expected results documentation

**Created by Tester (you, in M21):**
- `test-results.txt` - Output from `run-tests.bat`
- Issue reports (if any tests failed)
- Screenshots of failures (if applicable)
- Confirmation email/message: "M21 complete - all tests pass ✅"

### Estimated Timeline

**Best case (everything works first try):** ~45 minutes
**Typical case (1-2 issues found):** ~2 hours
**Worst case (multiple issues):** ~4 hours (but each fix cycle is faster in Linux)

**Key to Success:** Thorough M17-M20 preparation minimizes M21 iterations.

---

## 7) Acceptance Criteria (Definition of Done)

- Engine binary provides **init / size / checklist / heat / save‑decision / import / scrape / settings** commands with JSON I/O.
- **5 hard gates** enforced in engine (not VBA).
- **6‑item checklist** → banner logic matches spec.
- Sizing math matches examples: E=10k; r=0.75%; N=1.50; K=2 → R=$75; stop=3.00; 25 shares.
- Heat caps honored (portfolio 4%, bucket 1.5%).
- Excel buttons work through the bridge; workbook state is **not** the source of truth.
- Logging: both Excel and engine record the same **corr_id** per action.
- BDD suite + unit + integration tests green in CI.
- Release bundle tested on a clean machine.

---

## 8) BDD (Gherkin) — Key Features

### 8.1 Checklist → Banner
```gherkin
Feature: Checklist evaluation to banner
  Scenario: All six pass yields GREEN
    Given all six checklist items are true
    When I evaluate AAPL at entry 180 with N 1.5 and K 2
    Then banner is "GREEN"
```

### 8.2 Sizing
```gherkin
Feature: Position sizing
  Background:
    Given Equity_E is 10000 and RiskPct_r is 0.0075

  Scenario: Stock sizing example
    When I size with entry 180, N 1.5, K 2, method "stock"
    Then risk_dollars is 75 and shares is 25
```

### 8.3 Heat & Hard Gates
```gherkin
Feature: SaveDecision hard gates
  Scenario: Impulse timer blocks premature save
    Given GREEN for AAPL and evaluation just occurred
    When I attempt to save after 60 seconds
    Then accepted is false and reason contains "2-minute impulse"

  Scenario: Not in today's candidates blocks save
    Given GREEN for MSFT
    And MSFT is absent from today's candidates
    When I save decision
    Then accepted is false and reason contains "today's Candidates"
```

### 8.4 Transport Parity
```gherkin
Feature: CLI/HTTP parity
  Scenario: Sizing identical via CLI and HTTP
    Given server listens on 127.0.0.1:18888
    When I invoke "tf-engine size ..." and POST /sizing with same payload
    Then JSON bodies are identical
```

---

## 9) Testing Strategy

**Linux/WSL (automated, fast iteration):**
- **Unit tests** — pure functions (sizing, checklist, heat, gating)
- **Integration tests** — CLI JSON outputs; DB read/write
- **Parity tests** — CLI vs HTTP JSON equivalence
- **BDD tests** — Gherkin scenarios executed via godog
- **JSON contract tests** — Validate all responses against schemas

**Windows (manual, slower):**
- **VBA unit tests** — Shell execution, JSON parsing (via test worksheet)
- **Excel integration tests** — Full workflows through UI buttons
- **Smoke tests** — Basic engine functionality on Windows
- **E2E tests** — Complete trade flow from Excel to database

**Testing Sequence:**
1. **M1-M16:** Linux-based automated testing (Go engine, no VBA)
2. **M17-M20:** JSON contract validation (Linux) + VBA syntax checks
3. **M21:** Windows manual testing (VBA execution + Excel integration)
4. **Phase E:** Final cross-platform validation

---

## 10) Excel Bridge — Implementation Notes

- Use `WScript.Shell.Exec` for CLI; block until `.Status <> 0`; read `.StdOut` (JSON) and `.StdErr` (error text).
- Always wrap calls with try/catch; surface errors in a **Status** cell with `corr_id`.
- Provide a workbook **setting cell** for `tf-engine.exe` path and another toggle for HTTP mode.
- Keep current **Setup sheet**, **Open Debug Log** button; add **Open Engine Log** button next to it.

---

## 11) FINVIZ Import — Reliability Notes

- Respect pagination & rate‑limits (≥1s between pages).
- Normalize tickers (`BRK.B` → `BRK-B` style decisions documented).
- Manual paste fallback must remain always available from Excel.

---

## 12) Packaging & Delivery

**Bundle:**
```
/release
  tf-engine.exe
  trading.db                (created on first run if absent)
  excel/TrendFollowing_TradeEntry.xlsm
  docs/{start-here,quick-start,architecture,troubleshooting}.md
  LICENSE
```

**Install:**
1. Unzip to `C:\trading-engine\`
2. Double-click Excel workbook
3. In **Setup** sheet, set engine path and click **Test Connection**
4. Use **Import Candidates** & standard flow

---

## 13) Timeline (indicative)

- **Week 1:** Phase A/B (init + core domain + tests) — Linux
- **Week 2:** Finish domain + FINVIZ + DB + CLI UX — Linux
- **Week 3:** HTTP parity + JSON contracts + VBA implementation — Linux
- **Week 3.5:** Windows integration validation (M21) — Windows (manual testing)
- **Week 4:** BDD suite, docs, packaging, E2E verification — Linux + Windows

> Adjust as needed; milestones are functional (not calendar‑locked).

**Development Environment Split:**
- **Phases A-D (M1-M20):** Linux/WSL development (fast iteration)
- **M21:** Windows manual testing (slower, but minimized through thorough Linux prep)
- **Phase E:** Final hardening across both environments

---

## 14) Risks & Mitigations

- **FINVIZ HTML changes / rate‑limit** → strict error handling, retries; keep manual paste fallback.
- **Excel UI quirks** (checkbox/dropdown creation) → correctness lives in engine; troubleshooting keeps existing Setup/Log UX.
- **Path/quoting issues** in CLI → one shared `RunCommand()` helper; robust quoting; stderr surfaced.
- **Doc drift** → consolidate into 4 living docs in `/docs` and keep a Project Index current.
- **Linux/Windows environment split** → Validate JSON contracts thoroughly in Linux (M17-M20) to minimize Windows test iterations (M21); keep VBA thin to reduce cross-platform issues.
- **VBA testing bottleneck** → Automated testing in Linux covers 90% of functionality; Windows testing focuses only on VBA bridge and Excel UI integration; provide clear test procedures and issue templates to speed feedback loop.
- **Cross-compilation issues** → Test Windows binary on actual Windows environment early; use `GOOS=windows GOARCH=amd64` build tags consistently.

---

## 15) Appendices

### A) Sizing Math (sanity example)
- E = 10,000; r = 0.75% → **R=$75**
- N = 1.50; K = 2 → Stop distance = **3.00**
- Stock shares = floor(75/3) = **25**, initial stop = **$177.00**
- Options (delta‑ATR): contracts = floor(75/(2×1.5×0.3×100)) = **0**
- Options (max‑loss $50): contracts = floor(75/50) = **1**

### B) Minimal VBA (pseudocode)
```vba
Function Engine_Size(...) As EngineSizingResult
  Dim cmd As String: cmd = BuildSizeCommand(...)
  Dim json As String: json = RunCommand(cmd, corrId)
  Set Engine_Size = ParseSizing(json)
End Function
```

### C) HTTP Endpoints
- `GET /health` → `{status, version}`
- `GET /settings` / `PUT /settings`
- `POST /sizing`, `POST /evaluate`, `GET /heat`, `POST /decision`
- `POST /presets/import`, `POST /candidates`, `GET /candidates?date=...`

---

**Ready to build.** This plan delivers a reliable, testable engine with Excel as a thin shell, preserves your existing workflow, and lays a clean runway for a future web UI without rework.
