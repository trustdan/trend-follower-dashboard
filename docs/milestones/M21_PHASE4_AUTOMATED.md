# M21 Phase 4 - Automated Integration Tests

**QUICK START:** One-click automated testing - no manual steps required!

---

## TL;DR - Just Run This

```cmd
cd C:\Users\<YourUser>\excel-trading-dashboard\windows
run-integration-tests.bat
```

That's it! Tests run automatically in 1-2 minutes.

---

## What It Does

The automated test runner:

1. ✅ Verifies environment (files, engine, database)
2. ✅ Imports test module into Excel workbook
3. ✅ Clears previous test data
4. ✅ Imports test candidates (AAPL, MSFT, NVDA, SPY, JPM)
5. ✅ Runs all 25 integration tests automatically
6. ✅ Logs results to worksheet + file
7. ✅ Opens results for review

**No manual UI interaction required!**

---

## Test Coverage

### Automated Tests (19 tests)

**Workflow 1: Position Sizing (4 tests)**
- 1.1: Stock sizing (default settings)
- 1.2: Stock sizing (with overrides)
- 1.3: Option sizing (delta-ATR method)
- 1.4: Option sizing (max-loss method)

**Workflow 2: Checklist Evaluation (5 tests)**
- 2.1: GREEN banner (all 6 checks pass)
- 2.2: YELLOW banner (2 missing items)
- 2.3: YELLOW banner (1 missing item - edge case)
- 2.4: RED banner (3+ missing items)
- 2.5: Banner persistence (database verification)

**Workflow 3: Heat Management (4 tests)**
- 3.1: No open positions (clean state)
- 3.2: Portfolio cap exceeded
- 3.3: Bucket cap exceeded
- 3.4: Exactly at cap (edge case)

**Workflow 4: Save Decision - 5 Gates (6 tests)**
- 4.1: Happy path (all gates pass)
- 4.2: Gate 1 rejection (YELLOW banner)
- 4.3: Gate 1 rejection (RED banner)
- 4.4: Gate 2 rejection (not in candidates)
- 4.5: Gate 5 rejection (portfolio cap)
- 4.6: Gate 5 rejection (bucket cap)

### Manual Tests (6 tests - timing-dependent)

These require manual execution due to time-based gates:

**Workflow 3: Heat Management (2 tests)**
- 3.5: With open positions (cumulative heat)
- 3.6: Different buckets (isolation verification)

**Workflow 4: Save Decision (4 tests)**
- 4.7: Gate 3 rejection (2-minute impulse brake)
- 4.8: Gate 4 rejection (bucket cooldown)
- 4.9: Multiple gate failures
- 4.10: Form behavior (clear vs persist)

**Manual test procedures:** See `M21_PHASE4_TEST_SCRIPTS.md`

---

## System Requirements

**Prerequisites:**
- Windows 10/11
- Excel with macro support
- VBA project access enabled (see Setup section)
- setup-all.bat completed (Phases 1-3)

**Disk Space:**
- ~5 MB for test execution
- Log files: ~100 KB per run

---

## Setup (One-Time)

### Enable VBA Project Access

Required for automated module import:

1. Open Excel
2. File → Options → Trust Center → Trust Center Settings
3. Macro Settings tab
4. Check: **"Trust access to the VBA project object model"**
5. Click OK

**Why needed:** Allows batch file to automatically import test module

---

## Usage

### Option 1: Automated Tests (Recommended)

```cmd
cd windows
run-integration-tests.bat
```

**What happens:**
1. Script verifies environment (10 seconds)
2. Imports TFIntegrationTests.bas module
3. Runs tests in Excel (background)
4. Shows results in worksheet + log file
5. Opens Excel and log automatically

**Duration:** 1-2 minutes

**Output:**
- **Integration Tests worksheet:** Summary table with color-coded results
- **logs/integration-tests-YYYYMMDD-HHMMSS.log:** Detailed execution log

### Option 2: Manual Execution from Excel

If you prefer to run from Excel directly:

1. Open TradingPlatform.xlsm
2. Alt+F11 (open VBA editor)
3. Verify TFIntegrationTests module loaded
4. Alt+F8 (Macros dialog)
5. Select: `TFIntegrationTests.RunAllIntegrationTests`
6. Click Run

**Duration:** 1-2 minutes

---

## Reading Results

### Integration Tests Worksheet

Created automatically in TradingPlatform.xlsm:

```
Test ID | Workflow         | Test Name               | Status | Expected          | Actual            | Duration
--------|------------------|-------------------------|--------|-------------------|-------------------|----------
1.1     | Position Sizing  | Stock Sizing Default    | PASS   | R=$75, Shares=25  | R=75, Shares=25   | 0.234s
1.2     | Position Sizing  | Stock Sizing Overrides  | PASS   | R=$200, Shares=33 | R=200, Shares=33  | 0.187s
2.1     | Checklist        | GREEN Banner            | PASS   | Banner=GREEN      | Banner=GREEN      | 0.156s
...
```

**Color Coding:**
- 🟢 **Green (PASS):** Test passed
- 🟡 **Yellow (FAIL):** Test failed (expected behavior not met)
- 🔴 **Red (ERROR):** Test error (exception, parse failure, etc.)

**Summary:**
- Bottom of worksheet shows:
  - Total tests run
  - Pass/Fail/Error counts
  - Percentage breakdown
  - Total duration

### Log File

Location: `windows/logs/integration-tests-YYYYMMDD-HHMMSS.log`

Example:
```
========================================
M21 PHASE 4 INTEGRATION TESTS
========================================
Started: 2025-10-27 14:30:15

PRE-TEST SETUP
----------------------------------------
Checking engine accessibility...
  Engine accessible: {"version":"3.0.0-dev"}
Clearing previous test data...
  Test data cleared
Importing test candidates...
  Candidates imported: AAPL, MSFT, NVDA, SPY, JPM
Pre-test setup complete

========================================
WORKFLOW 1: POSITION SIZING
========================================

Test 1.1: Stock Sizing (Default Settings)
  PASS

Test 1.2: Stock Sizing (With Overrides)
  PASS

...

========================================
TEST SUMMARY
========================================
Total Tests:  19
PASS:         19 (100.0%)
FAIL:         0 (0.0%)
ERROR:        0 (0.0%)
Duration:     45.23 seconds
Completed:    2025-10-27 14:31:00

Log file: C:\...\windows\logs\integration-tests-20251027-143015.log
========================================
```

---

## Troubleshooting

### Error: "Trust access to VBA project" Required

**Symptom:** Script fails at step [3/5] importing module

**Fix:**
1. Excel → File → Options → Trust Center
2. Trust Center Settings → Macro Settings
3. Check "Trust access to VBA project object model"
4. Re-run script

### Error: "TradingPlatform.xlsm not found"

**Symptom:** Script fails at step [1/5]

**Fix:**
1. Verify you're in windows/ directory
2. Run `setup-all.bat` first to create workbook
3. Re-run integration tests

### Error: "Engine not working"

**Symptom:** Script fails at step [2/5]

**Fix:**
1. Verify tf-engine.exe exists: `dir tf-engine.exe`
2. Test manually: `tf-engine.exe --version`
3. If missing, rebuild: `cd .. && go build -o windows/tf-engine.exe ./cmd/tf-engine`
4. Ensure CGO enabled for SQLite support

### Tests Fail with "Failed to parse JSON"

**Symptom:** Multiple tests show ERROR status with "Failed to parse JSON"

**Fix:**
1. Check TradingSystem_Debug.log for engine errors
2. Verify database intact: `tf-engine.exe get-settings`
3. Check engine JSON output format matches VBA parser
4. Re-run setup-all.bat to reset environment

### Tests Fail with Gate Rejections

**Symptom:** Gate 1/2/5 tests fail unexpectedly

**Fix:**
1. Verify test data imported: `tf-engine.exe list-candidates`
2. Check settings: `tf-engine.exe get-settings`
3. Clear stale data: `tf-engine.exe -c "DELETE FROM candidates WHERE preset='AUTOTEST'"`
4. Re-run tests

### Excel Opens But Tests Don't Run

**Symptom:** Excel opens but Integration Tests worksheet not created

**Fix:**
1. Check if TFIntegrationTests module imported:
   - Alt+F11 in Excel
   - Look for TFIntegrationTests in Modules list
2. If missing, manually import:
   - File → Import File → Select TFIntegrationTests.bas
3. Run manually: Alt+F8 → RunAllIntegrationTests

---

## Manual Tests (Timing-Dependent)

After automated tests complete, optionally run manual tests for timing-based gates:

### Test 3.5: Heat with Open Positions

1. Save a decision (creates open position):
   ```cmd
   tf-engine.exe save-decision --ticker AAPL --entry 180 --atr 1.5 --method stock --banner GREEN --risk 75 --shares 25 --bucket "Tech/Comm" --preset MANUAL
   ```

2. Check heat with new trade in same bucket:
   ```cmd
   tf-engine.exe heat --risk 80 --bucket "Tech/Comm"
   ```

3. Verify:
   - Current Heat = $75 (from AAPL)
   - New Heat = $155 ($75 + $80)
   - Bucket Exceeded = TRUE (over $150 cap)

### Test 4.7: Gate 3 (Impulse Brake)

1. Evaluate checklist for MSFT (GREEN):
   ```cmd
   tf-engine.exe checklist --ticker MSFT --checks 1,1,1,1,1,1
   ```

2. **Immediately** try to save decision (within 2 minutes):
   ```cmd
   tf-engine.exe save-decision --ticker MSFT --entry 400 --atr 3.0 --method stock --banner GREEN --risk 75 --shares 25 --bucket "Tech/Comm" --preset MANUAL
   ```

3. Verify: Rejected with "Wait XX seconds" (Gate 3 FAIL)

4. Wait 2 minutes, retry - should succeed (Gate 3 PASS)

### Test 4.8: Gate 4 (Bucket Cooldown)

1. Save first decision in Tech/Comm bucket (Test 4.7 success)

2. Try to save second decision in same bucket:
   ```cmd
   tf-engine.exe save-decision --ticker NVDA --entry 500 --atr 5.0 --method stock --banner GREEN --risk 75 --shares 15 --bucket "Tech/Comm" --preset MANUAL
   ```

3. Verify: Rejected with "Bucket in cooldown" (Gate 4 FAIL)

4. Try different bucket (should succeed):
   ```cmd
   tf-engine.exe save-decision --ticker JPM --entry 150 --atr 1.5 --method stock --banner GREEN --risk 75 --shares 50 --bucket Finance --preset MANUAL
   ```

---

## Architecture

### How It Works

```
run-integration-tests.bat
  ↓
[1] Verify environment (files, engine, db)
  ↓
[2] Import TFIntegrationTests.bas via VBScript
  ↓
[3] Run TFIntegrationTests.RunAllIntegrationTests
  ↓
  ├─ PreTestSetup()
  │    ├─ Clear test data
  │    └─ Import candidates
  │
  ├─ RunWorkflow1_PositionSizing()
  │    ├─ Test_1_1_StockSizingDefault()
  │    ├─ Test_1_2_StockSizingOverrides()
  │    ├─ Test_1_3_OptionDeltaATR()
  │    └─ Test_1_4_OptionMaxLoss()
  │
  ├─ RunWorkflow2_ChecklistEvaluation()
  │    ├─ Test_2_1_GreenBanner()
  │    ├─ Test_2_2_YellowBanner2Missing()
  │    ├─ Test_2_3_YellowBanner1Missing()
  │    ├─ Test_2_4_RedBanner3Missing()
  │    └─ Test_2_5_BannerPersistence()
  │
  ├─ RunWorkflow3_HeatManagement()
  │    ├─ Test_3_1_NoOpenPositions()
  │    ├─ Test_3_2_PortfolioCapExceeded()
  │    ├─ Test_3_3_BucketCapExceeded()
  │    └─ Test_3_4_ExactlyAtCap()
  │
  └─ RunWorkflow4_SaveDecision()
       ├─ Test_4_1_HappyPathAllGatesPass()
       ├─ Test_4_2_Gate1RejectionYellow()
       ├─ Test_4_3_Gate1RejectionRed()
       ├─ Test_4_4_Gate2RejectionNotInCandidates()
       ├─ Test_4_5_Gate5RejectionPortfolioCap()
       └─ Test_4_6_Gate5RejectionBucketCap()
  ↓
[4] Write results to worksheet
  ↓
[5] Write detailed log to file
  ↓
[6] Display summary message box
```

### VBA Module Structure

**TFIntegrationTests.bas** (1,100+ lines)
- Main runner: `RunAllIntegrationTests()`
- Pre-test setup: `PreTestSetup()`
- Workflow runners: `RunWorkflow1()`, `RunWorkflow2()`, etc.
- Individual tests: `Test_1_1()`, `Test_1_2()`, etc.
- Logging: `LogMessage()`, `WriteResultsToWorksheet()`
- Result type: `IntegrationTestResult`

**Each test:**
1. Logs start
2. Builds engine command
3. Executes via `TFEngine.ExecuteCommand()`
4. Parses JSON result
5. Validates expected vs actual
6. Sets PASS/FAIL/ERROR status
7. Logs result
8. Adds to results collection

**All tests use existing VBA modules:**
- TFEngine.bas - Command execution
- TFHelpers.bas - JSON parsing
- TFTypes.bas - Type definitions

---

## Files Created

### New Files

```
excel/vba/TFIntegrationTests.bas          (25 KB) - Automated test runner module
windows/run-integration-tests.bat         (9 KB)  - One-click execution script
docs/milestones/M21_PHASE4_AUTOMATED.md   (this file)
```

### Generated at Runtime

```
windows/logs/integration-tests-YYYYMMDD-HHMMSS.log  - Detailed execution log
TradingPlatform.xlsm → Integration Tests worksheet   - Results table
```

---

## Comparison: Manual vs Automated

| Aspect | Manual Testing | Automated Testing |
|--------|----------------|-------------------|
| **Setup** | 10-15 min per workflow | One-time VBA access enable |
| **Execution** | 45-120 minutes | 1-2 minutes |
| **Worksheets** | Create 4 sheets manually | No UI work needed |
| **VBA Code** | Copy-paste 4 button handlers | Auto-imported |
| **Test Data** | Enter manually | Auto-loaded |
| **Results** | Manual observation | Auto-logged |
| **Repeatability** | Error-prone | 100% consistent |
| **Documentation** | Manual notes | Auto-generated logs |
| **Coverage** | All 25 tests | 19 automated + 6 manual |

**Recommendation:** Use automated tests for 19 tests, manual for 6 timing-dependent tests

---

## Success Criteria

**Phase 4 Complete When:**
- ✅ Automated tests: 19/19 PASS
- ✅ Manual tests: 6/6 PASS
- ✅ No ERROR status in results
- ✅ All workflows validated
- ✅ Logs confirm expected behavior

**Expected Results:**
- Position Sizing: All calculations correct (actual risk ≤ specified)
- Checklist: Only GREEN allows save
- Heat Management: Both caps enforced (4% portfolio, 1.5% bucket)
- Save Decision: All 5 gates enforced correctly

---

## Next Steps

After automated tests pass:

1. **Review Results:**
   - Check Integration Tests worksheet
   - Review log file for any warnings
   - Verify all tests PASS

2. **Manual Tests (Optional):**
   - Test Gate 3 (impulse brake)
   - Test Gate 4 (cooldown)
   - Test cumulative heat (3.5, 3.6)

3. **Complete M21:**
   - Document any issues found
   - Create M21 completion summary
   - Begin M22 planning

---

## FAQ

**Q: Do I need to create any worksheets manually?**
A: No! The automated tests work directly with the engine (no UI).

**Q: Can I run tests multiple times?**
A: Yes! Tests clean up after themselves (preset='AUTOTEST').

**Q: What if some tests fail?**
A: Check the log file for details. Common issues:
- Engine not working → rebuild
- Database corrupt → re-run setup-all.bat
- VBA module not imported → check VBA project access

**Q: How long do logs persist?**
A: Forever (unless manually deleted). Each run creates new timestamped log.

**Q: Can I run tests on Linux/Mac?**
A: No, requires Excel + VBA. But engine tests can run via CLI on any platform.

**Q: What about the 6 manual tests?**
A: Optional for Phase 4 completion. Main validation is automated tests.

---

**Created:** 2025-10-27
**For:** M21 Phase 4 Integration Tests
**Method:** Automated execution via VBA + Batch script
**Duration:** 1-2 minutes (19 tests)
**Maintainer:** See docs/dev/CLAUDE_RULES.md for contribution guidelines

---

**Ready to test? Just run:**
```cmd
cd windows
run-integration-tests.bat
```

Good luck! 🚀
