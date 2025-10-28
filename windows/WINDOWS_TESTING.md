# Windows Integration Testing Guide (M21)
# Trading Engine v3

**Purpose:** Step-by-step manual testing procedures for VBA ↔ Go engine integration
**Target:** Windows 10/11 with Excel desktop
**Duration:** ~45 minutes (best case) to ~4 hours (with issues)
**Prerequisites:** M17-M20 completed

---

## Pre-Test Setup

### Equipment Needed
- Windows PC (Windows 10 or 11)
- Excel desktop (Microsoft 365, Excel 2019, or Excel 2021)
- Administrator privileges (for script execution)
- Internet access (for documentation lookup if needed)

### Files Required (from windows/ folder)
- ✅ `tf-engine.exe` (12 MB Windows binary)
- ✅ `windows-import-vba.bat`
- ✅ `windows-init-database.bat`
- ✅ `run-tests.bat`
- ✅ `test-data/` folder with JSON examples
- ✅ `EXCEL_WORKBOOK_TEMPLATE.md` (workbook structure spec)
- ✅ `../excel/vba/` folder with .bas modules

---

## Phase 1: Pre-Test Setup (~10 min)

### Step 1.1: Copy Files to Windows Machine

**Action:**
1. Copy entire `windows/` folder from Linux/WSL to Windows PC
2. Suggested location: `C:\trading-engine\`
3. Verify all files present:
   ```
   C:\trading-engine\
   ├── tf-engine.exe
   ├── windows-import-vba.bat
   ├── windows-init-database.bat
   ├── run-tests.bat
   ├── WINDOWS_TESTING.md (this file)
   ├── EXCEL_WORKBOOK_TEMPLATE.md
   └── test-data\
       └── (21 JSON files)
   ```

**Verification:**
- [ ] All files copied successfully
- [ ] No missing files or folders
- [ ] tf-engine.exe is 12 MB (approximately)

---

### Step 1.2: Create Excel Workbook

**Action:**
1. Open Excel on Windows
2. Create new blank workbook
3. Save as: `C:\trading-engine\TradingPlatform.xlsm`
4. File type: "Excel Macro-Enabled Workbook (*.xlsm)"

**Enable VBA Project Access:**
5. File > Options > Trust Center > Trust Center Settings
6. Macro Settings tab
7. ✅ Check "Trust access to the VBA project object model"
8. Click OK

**Create Basic Worksheets:**
9. Rename Sheet1 to "Setup"
10. Create new sheet: "VBA Tests"
11. Save workbook

**Verification:**
- [ ] Workbook saved as .xlsm in correct location
- [ ] VBA project access enabled
- [ ] Two worksheets created: Setup, VBA Tests

---

### Step 1.3: Import VBA Modules

**Action:**
1. Open Command Prompt as Administrator
2. Navigate to: `cd C:\trading-engine`
3. Run: `windows-import-vba.bat`
4. Watch for success message

**Expected Output:**
```
========================================
 VBA Module Import Script
 Trading Engine v3
========================================

Excel workbook: TradingPlatform.xlsm
VBA modules source: ..\excel\vba\

Creating import script...
Script created: import_vba.vbs

Importing VBA modules...

Importing TFTypes.bas...
Importing TFHelpers.bas...
Importing TFEngine.bas...
Importing TFTests.bas...
VBA modules imported successfully!

========================================
 VBA Import Complete!
========================================
```

**Verify in Excel:**
5. Open TradingPlatform.xlsm
6. Press Alt+F11 to open VBA Editor
7. Verify modules loaded in Project Explorer:
   - TFTypes
   - TFHelpers
   - TFEngine
   - TFTests

**Verification:**
- [ ] Import script ran without errors
- [ ] All 4 modules visible in VBA Editor
- [ ] No compilation errors (click Debug > Compile VBAProject)

**If Import Fails:**
- Check "Trust access to VBA project object model" setting
- Ensure Excel is not already running with workbook open
- Check antivirus is not blocking VBScript execution
- Try running as Administrator

---

### Step 1.4: Initialize Database

**Action:**
1. In Command Prompt: `windows-init-database.bat`
2. Watch for success message

**Expected Output:**
```
========================================
 Trading Database Initialization
 Trading Engine v3
========================================

Engine: tf-engine.exe
Database: trading.db

Initializing database...

Database initialized successfully
Settings initialized with defaults

========================================
 Database Initialized Successfully!
========================================

Database file created: C:\trading-engine\trading.db

Verifying database schema...
{
  "BucketHeatCap_pct": "0.015",
  "Equity_E": "10000",
  "HeatCap_H_pct": "0.04",
  "RiskPct_r": "0.0075",
  "StopMultiple_K": "2"
}

Database schema verified OK

Default Settings:
  - Equity (E):           $10,000
  - Risk per trade (r):   0.75%
  - Portfolio heat cap:   4%
  - Bucket heat cap:      1.5%
  - Stop multiple (K):    2
```

**Verification:**
- [ ] Database file created: `trading.db`
- [ ] Settings JSON displayed correctly
- [ ] No errors in output

---

### Step 1.5: Setup Excel Configuration

**Action:**
1. Open TradingPlatform.xlsm
2. Go to "Setup" worksheet
3. In cell A4, type: "Engine Path:"
4. In cell B4, type: `.\tf-engine.exe`
5. Define named range: Select B4, Name Box: `EnginePathSetting`
6. In cell A5, type: "Database Path:"
7. In cell B5, type: `.\trading.db`
8. Define named range: Select B5, Name Box: `DatabasePathSetting`
9. Save workbook

**Verification:**
- [ ] Named ranges created correctly
- [ ] Paths point to correct locations

---

**Phase 1 Complete! ✅**

**Checklist:**
- [ ] Files copied to Windows
- [ ] Excel workbook created (.xlsm)
- [ ] VBA modules imported (4 modules visible)
- [ ] Database initialized (trading.db exists)
- [ ] Named ranges configured

**Time Estimate:** Actual: _______ minutes

---

## Phase 2: Smoke Tests (~5 min)

**Purpose:** Verify basic engine functionality before VBA testing

### Test 2.1: Engine Version

**Action:**
```cmd
tf-engine.exe --version
```

**Expected Output:**
```
tf-engine version 3.0.0-dev
```

**Verification:**
- [ ] Version displayed
- [ ] No errors
- [ ] Exit code 0

---

### Test 2.2: Database Verification

**Action:**
```cmd
tf-engine.exe --db trading.db get-settings --format json
```

**Expected Output:**
```json
{
  "BucketHeatCap_pct": "0.015",
  "Equity_E": "10000",
  "HeatCap_H_pct": "0.04",
  "RiskPct_r": "0.0075",
  "StopMultiple_K": "2"
}
```

**Verification:**
- [ ] JSON output well-formed
- [ ] All 5 settings present
- [ ] No errors

---

### Test 2.3: Position Sizing Command

**Action:**
```cmd
tf-engine.exe --db trading.db size --entry 180 --atr 1.5 --k 2 --method stock --format json
```

**Expected Output:**
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

**Verification:**
- [ ] JSON output matches expected
- [ ] shares = 25
- [ ] initial_stop = 177
- [ ] No errors

---

**Phase 2 Complete! ✅**

**Checklist:**
- [ ] Engine executes without errors
- [ ] Database accessible
- [ ] Position sizing calculation correct

**Time Estimate:** Actual: _______ minutes

---

## Phase 3: VBA Unit Tests (~10 min)

**Purpose:** Test VBA modules in isolation before full integration

### Step 3.1: Prepare VBA Tests Worksheet

**Action:**
1. Open TradingPlatform.xlsm
2. Go to "VBA Tests" worksheet
3. In cell A1, type: "VBA Unit Tests"
4. Format as header (Bold, 14pt)
5. In cell A3, insert ActiveX CommandButton
6. Set properties:
   - Name: `btnRunTests`
   - Caption: "Run All Tests"
7. Double-click button, add code:
   ```vba
   Private Sub btnRunTests_Click()
       TFTests.RunAllTests
   End Sub
   ```
8. Save workbook

---

### Step 3.2: Run VBA Unit Tests

**Action:**
1. Enable macros if prompted
2. Exit Design Mode (Developer tab)
3. Click "Run All Tests" button
4. Watch results populate worksheet

**Expected Output:**

```
VBA Unit Tests

[Run All Tests]    Status: Running...

JSON PARSING TESTS
ParseSizingJSON                  ✅ PASS    All fields parsed correctly           0.003s
ParseChecklistJSON_Green         ✅ PASS    GREEN banner parsed correctly         0.002s
ParseChecklistJSON_Yellow        ✅ PASS    YELLOW banner with missing items...   0.002s
ParseHeatJSON                    ✅ PASS    Heat check parsed correctly           0.003s
ParseTimerJSON                   ✅ PASS    Timer check parsed correctly          0.002s
ParseSettingsJSON                ✅ PASS    Settings parsed correctly             0.002s

HELPER FUNCTION TESTS
ExtractJSONValue                 ✅ PASS    JSON value extraction works           0.001s
ExtractJSONArray                 ✅ PASS    JSON array extraction works           0.002s
GenerateCorrelationID            ✅ PASS    Correlation ID format correct: ...    0.001s

VALIDATION TESTS
ValidateTicker                   ✅ PASS    Ticker validation works               0.001s
ValidatePositiveNumber           ✅ PASS    Number validation works               0.001s

FORMATTING TESTS
FormatCurrency                   ✅ PASS    Currency formatting works             0.001s
FormatPercent                    ✅ PASS    Percent formatting works              0.001s

SHELL EXECUTION TESTS
ShellExecution                   ✅ PASS    Engine executable found and runs      0.150s

SUMMARY
Total Tests: 14
Passed: 14
Failed: 0
Total Time: 0.173s
Result: ✅ ALL TESTS PASSED
```

**Verification:**
- [ ] All 14 tests show PASS
- [ ] No VBA runtime errors
- [ ] Summary shows 0 failures
- [ ] ShellExecution test passed (engine accessible)

---

### Step 3.3: Review Test Logs

**Action:**
1. Open: `C:\trading-engine\TradingSystem_Debug.log`
2. Verify log entries for test execution

**Expected Log Entries:**
```
[2025-10-27 14:30:52] [INFO] [20251027-143052-7A3F] TradingPlatform.xlsm opened
[2025-10-27 14:31:15] [INFO] [20251027-143115-2B4E] Executing: "tf-engine.exe" --db "trading.db" --corr-id 20251027-143115-2B4E --format json --version
[2025-10-27 14:31:15] [INFO] [20251027-143115-2B4E] Command succeeded (23 bytes JSON)
```

**Verification:**
- [ ] Log file exists
- [ ] Correlation IDs present
- [ ] No ERROR entries

---

**Phase 3 Complete! ✅**

**Checklist:**
- [ ] VBA Tests worksheet configured
- [ ] All 14 unit tests passed
- [ ] Logs show successful execution

**Time Estimate:** Actual: _______ minutes

---

## Phase 4: Integration Tests (~15 min)

**Purpose:** Test complete workflows through Excel UI

### Test 4.1: Position Sizing Workflow

**Setup:**
1. Create new worksheet: "Position Sizing"
2. Add inputs (see EXCEL_WORKBOOK_TEMPLATE.md for full layout):
   - B4: Ticker (e.g., "AAPL")
   - B5: Entry (e.g., 180)
   - B6: ATR (e.g., 1.5)
   - B7: K (e.g., 2)
   - B8: Method (e.g., "stock")
3. Add button: "Calculate"
4. Link button to code (see template for full code)

**Test Action:**
1. Enter test values:
   - Ticker: AAPL
   - Entry: 180
   - ATR: 1.5
   - K: 2
   - Method: stock
2. Click "Calculate"

**Expected Result:**
- B21 (Risk Dollars): $75.00
- B22 (Stop Distance): 3.00
- B23 (Initial Stop): $177.00
- B24 (Shares): 25
- B25 (Contracts): 0
- B26 (Actual Risk): $75.00
- B29 (Status): "✅ Success (corr_id: XXXXX)"

**Verification:**
- [ ] All results match expected
- [ ] Status shows success with correlation ID
- [ ] No errors displayed

---

### Test 4.2: Checklist Evaluation (GREEN)

**Setup:**
1. Create worksheet: "Checklist"
2. Add ticker input: B3
3. Add 6 ActiveX checkboxes (Check1 - Check6)
4. Add button: "Evaluate"

**Test Action:**
1. Enter ticker: AAPL
2. Check all 6 checkboxes
3. Click "Evaluate"

**Expected Result:**
- B16 (Banner): "GREEN" with green background
- B17 (Missing Count): 0
- B18 (Missing Items): (empty)
- B19 (Allow Save): TRUE
- B20 (Evaluation Time): (timestamp)
- B27 (Status): "✅ Evaluated (corr_id: XXXXX)"

**Verification:**
- [ ] Banner shows GREEN with green background
- [ ] Missing count is 0
- [ ] Allow Save is TRUE
- [ ] Status shows success

---

### Test 4.3: Checklist Evaluation (YELLOW)

**Test Action:**
1. Uncheck 1 checkbox (e.g., "Higher high")
2. Click "Evaluate"

**Expected Result:**
- B16 (Banner): "YELLOW" with yellow background
- B17 (Missing Count): 1
- B18 (Missing Items): "Higher high"
- B19 (Allow Save): FALSE
- B27 (Status): "✅ Evaluated (corr_id: XXXXX)"

**Verification:**
- [ ] Banner shows YELLOW with yellow background
- [ ] Missing count is 1
- [ ] Missing item listed
- [ ] Allow Save is FALSE

---

### Test 4.4: Heat Management

**Setup:**
1. Create worksheet: "Heat Check"
2. Add inputs:
   - B4: Risk Dollars (e.g., 75)
   - B5: Bucket (e.g., "Tech/Comm")
3. Add button: "Check Heat"

**Test Action:**
1. Enter Risk: 75
2. Enter Bucket: Tech/Comm
3. Click "Check Heat"

**Expected Result:** (assuming no open positions)
- B10 (Current Portfolio Heat): $0.00
- B11 (New Portfolio Heat): $75.00
- B12 (Portfolio Heat %): 18.75%
- B13 (Portfolio Cap): $400.00 (4% of $10,000)
- B14 (Exceeded): FALSE
- B18 (Current Bucket Heat): $0.00
- B19 (New Bucket Heat): $75.00
- B20 (Bucket Heat %): 50.00%
- B21 (Bucket Cap): $150.00 (1.5% of $10,000)
- B22 (Exceeded): FALSE
- B26 (Allowed): TRUE (green background)

**Verification:**
- [ ] Heat calculations correct
- [ ] Caps shown correctly
- [ ] Not exceeded
- [ ] Allowed shows TRUE

---

### Test 4.5: Import Candidates

**Setup:**
1. Create worksheet: "Candidates"
2. Add inputs:
   - B4: Tickers (e.g., "AAPL,MSFT,NVDA")
   - B5: Preset (e.g., "TEST")
3. Add button: "Import"

**Test Action:**
1. Enter tickers: AAPL,MSFT,NVDA
2. Enter preset: TEST
3. Click "Import"

**Expected Result:**
- Status: "✅ Imported 3 candidates"
- Database contains candidates

**Verification (via CLI):**
```cmd
tf-engine.exe list-candidates --format json
```

Should show 3 candidates: AAPL, MSFT, NVDA

**Verification:**
- [ ] Import successful
- [ ] Candidates appear in database
- [ ] Status shows count

---

### Test 4.6: Save Decision (Happy Path)

**Prerequisites:**
1. Import candidates: AAPL (from Test 4.5)
2. Complete checklist for AAPL: GREEN banner
3. Wait 2 minutes for impulse brake (or mock timer for testing)

**Test Action:**
1. Fill trade entry form:
   - Ticker: AAPL
   - Entry: 180
   - ATR: 1.5
   - K: 2
   - Banner: GREEN
   - Risk Dollars: 75
   - Shares: 25
   - Bucket: Tech/Comm
   - Preset: TEST
2. Click "Save Decision"

**Expected Result:**
- Status: "✅ Decision saved (ID: 1)"
- Form clears for next trade
- Confirmation message shown

**Verification (via CLI):**
```cmd
tf-engine.exe --db trading.db --format json -c "SELECT * FROM decisions ORDER BY id DESC LIMIT 1"
```
(Note: May need direct SQL query or list-decisions command if implemented)

**Verification:**
- [ ] Decision saved (ID returned)
- [ ] No error messages
- [ ] Database contains decision record

---

### Test 4.7: Save Decision (Gate Rejection - Banner)

**Test Action:**
1. Complete form with YELLOW banner (not GREEN)
2. Click "Save Decision"

**Expected Result:**
- Status: "❌ REJECTED: Banner must be GREEN"
- Form retained (not cleared)
- Error message explains gate failure

**Verification:**
- [ ] Rejection message clear
- [ ] Form not cleared
- [ ] Correlation ID shown for debugging

---

### Test 4.8: Save Decision (Gate Rejection - Not in Candidates)

**Test Action:**
1. Enter ticker NOT in candidate list (e.g., "ZZZZ")
2. Complete form with GREEN banner
3. Click "Save Decision"

**Expected Result:**
- Status: "❌ REJECTED: Ticker not in today's candidates"
- Form retained
- Error explains which gate failed

**Verification:**
- [ ] Correct gate identified as failed
- [ ] Clear error message
- [ ] Form retained

---

**Phase 4 Complete! ✅**

**Checklist:**
- [ ] Position sizing workflow works
- [ ] Checklist evaluation (GREEN) works
- [ ] Checklist evaluation (YELLOW) works
- [ ] Heat management calculation works
- [ ] Candidate import works
- [ ] Save decision (happy path) works
- [ ] Gate rejections work correctly

**Time Estimate:** Actual: _______ minutes

---

## Phase 5: Issue Reporting

**If any test fails, document using this template:**

### Issue Report Template

```
==================================================
ISSUE REPORT
==================================================
Date: [YYYY-MM-DD]
Tester: [Your name]
Test ID: [e.g., 4.2 - Checklist Evaluation GREEN]
Status: ❌ FAILED

EXPECTED BEHAVIOR:
--------------------
[Describe what should happen]

ACTUAL BEHAVIOR:
--------------------
[Describe what actually happened]

STEPS TO REPRODUCE:
--------------------
1. [Step 1]
2. [Step 2]
3. [Step 3]

ERROR DETAILS:
--------------------
- Correlation ID: [e.g., 20251027-143052-7A3F]
- Error Message: [Exact error text from Excel or stderr]
- VBA Error (if any): [Error number and description]
- Engine Output (if relevant):
  [Paste stdout/stderr from command]

SCREENSHOTS:
--------------------
[Attach screenshot of Excel showing issue]

LOG EXCERPTS:
--------------------
From TradingSystem_Debug.log:
[Paste relevant lines with correlation ID]

From tf-engine.log:
[Paste relevant lines with correlation ID]

ENVIRONMENT:
--------------------
- Windows Version: [e.g., Windows 11 Pro]
- Excel Version: [e.g., Microsoft 365 Excel]
- tf-engine.exe size: [e.g., 12 MB]
- Database exists: [Yes/No]

==================================================
```

### Developer Fix Process

1. Developer fixes issue in Linux (fast iteration)
2. Developer updates affected files
3. Copy updated files to Windows test environment
4. Re-run only affected test(s)
5. If pass, continue; if fail, repeat

---

## Phase 6: Final Validation (~5 min)

**Purpose:** Automated test suite to verify all components

### Step 6.1: Run Automated Test Runner

**Action:**
```cmd
run-tests.bat
```

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

[INTEGRATION TESTS - MANUAL]
Manual tests must be verified from worksheet results.
See Phase 4 checklist above.

========================================
RESULT: ALL AUTOMATED TESTS PASSED ✅
========================================

Test results saved to: test-results.txt
Logs available at: TradingSystem_Debug.log
Correlation IDs: [list of IDs]

Next: Verify manual integration tests (Phase 4)
```

**Verification:**
- [ ] All automated tests passed
- [ ] test-results.txt created
- [ ] Logs accessible

---

### Step 6.2: Final Checklist

**Before declaring M21 complete, verify:**

**Engine Tests:**
- [ ] tf-engine.exe executes without errors
- [ ] Database initializes correctly
- [ ] All CLI commands work (size, checklist, heat, etc.)

**VBA Tests:**
- [ ] All 4 modules imported successfully
- [ ] All 14 VBA unit tests pass
- [ ] No VBA compilation errors

**Integration Tests:**
- [ ] Position sizing workflow end-to-end
- [ ] Checklist evaluation (GREEN, YELLOW, RED)
- [ ] Heat management calculations
- [ ] Candidate import
- [ ] Save decision (happy path)
- [ ] Save decision (gate rejections)

**Logging & Debugging:**
- [ ] Correlation IDs appear in both logs
- [ ] TradingSystem_Debug.log created and readable
- [ ] tf-engine.log created and readable
- [ ] Correlation IDs cross-reference successfully

**Error Handling:**
- [ ] Errors display clearly with correlation IDs
- [ ] Form retains data on errors (doesn't clear)
- [ ] Actionable error messages shown

---

**M21 Complete! ✅**

**Final Sign-Off:**
- Date: _________________
- Tester: _________________
- Total Time: _______ hours _______ minutes
- Issues Found: _______ (see issue reports above)
- All Tests Passed: YES / NO

**Ready for Phase E (Hardening & Release)**

---

## Troubleshooting Guide

### Issue: "Trust access to VBA project object model" Error

**Solution:**
1. File > Options > Trust Center > Trust Center Settings
2. Macro Settings > Check "Trust access to the VBA project object model"
3. Restart Excel
4. Re-run windows-import-vba.bat

---

### Issue: VBA Import Script Fails

**Symptoms:**
- Error during windows-import-vba.bat
- Modules not appearing in VBA Editor

**Solutions:**
1. Close Excel completely before running import
2. Run Command Prompt as Administrator
3. Check antivirus is not blocking VBScript (.vbs files)
4. Manually import: VBA Editor > File > Import File > Select .bas files

---

### Issue: Engine Not Found

**Symptoms:**
- "tf-engine.exe not found" error
- Shell execution test fails

**Solutions:**
1. Verify tf-engine.exe in same folder as Excel workbook
2. Check named range `EnginePathSetting` points to correct path
3. Try absolute path: `C:\trading-engine\tf-engine.exe`
4. Check Windows Defender hasn't quarantined the .exe

---

### Issue: Database Initialization Fails

**Symptoms:**
- windows-init-database.bat fails
- "Failed to open database" errors

**Solutions:**
1. Check write permissions on folder
2. Delete existing trading.db and retry
3. Run as Administrator
4. Check disk space available

---

### Issue: JSON Parsing Errors

**Symptoms:**
- VBA unit tests fail on JSON parsing
- "Type mismatch" or "Invalid data" errors

**Solutions:**
1. Check engine outputs well-formed JSON (run CLI commands manually)
2. Verify test data JSON files are intact
3. Check for character encoding issues
4. Review ExtractJSONValue function logic

---

### Issue: Correlation IDs Not in Logs

**Symptoms:**
- TradingSystem_Debug.log empty or missing entries
- Correlation IDs not matching between logs

**Solutions:**
1. Check file write permissions
2. Verify LogMessage function in TFHelpers.bas
3. Check log file path in LogMessage
4. Manually create log file if missing

---

### Issue: Tests Timeout

**Symptoms:**
- "Command timed out" errors
- Excel freezes during engine calls

**Solutions:**
1. Increase timeout in TFEngine.bas (COMMAND_TIMEOUT_SECONDS)
2. Check engine is not hanging (run CLI manually)
3. Check antivirus real-time scanning (may slow execution)
4. Verify database is not locked by another process

---

## Appendix: Quick Reference

### Key File Locations
```
C:\trading-engine\
├── tf-engine.exe              - Go backend
├── trading.db                 - SQLite database
├── TradingPlatform.xlsm       - Excel workbook
├── TradingSystem_Debug.log    - VBA logs
├── tf-engine.log              - Go logs
├── windows-import-vba.bat     - VBA import script
├── windows-init-database.bat  - DB initialization
└── run-tests.bat              - Automated test runner
```

### Common Commands
```cmd
REM Version check
tf-engine.exe --version

REM Initialize database
tf-engine.exe init

REM Get settings
tf-engine.exe get-settings --format json

REM Position sizing
tf-engine.exe size --entry 180 --atr 1.5 --k 2 --method stock --format json

REM Check candidates
tf-engine.exe list-candidates --format json

REM Import candidates
tf-engine.exe import-candidates --tickers "AAPL,MSFT" --preset TEST
```

### VBA Immediate Window Commands
```vba
' Test correlation ID generation
?TFHelpers.GenerateCorrelationID()

' Test engine connection
?TFEngine.ExecuteCommand("--version").Success

' Test JSON parsing
?TFHelpers.ExtractJSONValue("{""test"": 123}", "test")
```

---

**This testing guide ensures thorough validation of the VBA ↔ Go engine integration before production use.**

**Remember:** This is a discipline enforcement system. Testing must be rigorous. No shortcuts.
