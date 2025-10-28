# Windows Integration Package
# Trading Engine v3

**Created:** 2025-10-27 (M20 - Windows Integration Package)
**Updated:** 2025-10-28 (M22 - Automated UI Generation)
**Purpose:** Complete Windows deployment and testing package
**Status:** Production Ready

---

## Quick Start (3 Minutes)

### One-Command Setup (Recommended - M22+)

1. **Copy to Windows**
   - Copy entire `windows/` folder to Windows PC
   - Suggested location: `C:\trading-engine\`

2. **Run Complete Setup**
   ```cmd
   cd C:\trading-engine
   1-setup-all.bat
   ```

3. **That's it!** The script will:
   - Create Excel workbook (TradingPlatform.xlsm)
   - Enable VBA project access
   - Import all VBA modules
   - Create 5 production UI worksheets
   - Initialize database (trading.db)
   - Run automated smoke tests

4. **Open and Use**
   - Open `TradingPlatform.xlsm`
   - Enable macros when prompted
   - Start with the Dashboard worksheet
   - Navigate to other worksheets using Dashboard buttons

### Manual Setup (Legacy - M20)

For step-by-step manual setup, see "Legacy Setup Process" section below.

---

## Package Contents

```
windows/
├── tf-engine.exe                   - Go backend (12 MB Windows binary)
├── 1-setup-all.bat                 - ⭐ ONE-CLICK COMPLETE SETUP (M22)
├── 2-update-vba.bat                - Update VBA modules only
├── 3-run-integration-tests.bat     - Run integration test suite
├── 4-run-tests.bat                 - Run all automated tests
├── create-ui-worksheets.vbs        - UI worksheet generator (M22)
├── vbscript-lib.vbs                - VBScript helper library (M22)
├── windows-import-vba.bat          - [DEPRECATED] VBA import (use 1-setup-all.bat)
├── windows-init-database.bat       - [DEPRECATED] DB init (use 1-setup-all.bat)
├── README.md                       - This file
├── WINDOWS_TESTING.md              - Comprehensive testing guide
├── EXCEL_WORKBOOK_TEMPLATE.md      - Workbook structure specification
└── test-data/
    ├── README.txt                  - Test data description
    └── *.json                      - 21 sample JSON response files
```

### Additional Files (in parent directories)
```
../excel/vba/
├── TFTypes.bas                     - Type definitions
├── TFHelpers.bas                   - JSON parsing & utilities
├── TFEngine.bas                    - Engine communication
└── TFTests.bas                     - VBA unit tests
```

---

## File Descriptions

### tf-engine.exe (12 MB)
**Purpose:** Go backend - all business logic and 5 hard gates enforcement

**Capabilities:**
- Position sizing (stock, options delta-ATR, options max-loss)
- Checklist evaluation → GREEN/YELLOW/RED banner
- Heat management (portfolio and bucket caps)
- 2-minute impulse brake timing
- Candidate management (import, list, check)
- Cooldown management
- Save decision with 5 hard gate enforcement
- Position tracking
- Settings management
- FINVIZ scraper integration

**Usage:**
```cmd
tf-engine.exe --help
tf-engine.exe --version
tf-engine.exe init
tf-engine.exe size --entry 180 --atr 1.5 --k 2 --method stock
tf-engine.exe checklist --ticker AAPL --checks true,true,true,true,true,true
tf-engine.exe heat --add-r 75 --bucket "Tech/Comm"
```

All commands support `--format json` for programmatic use (VBA bridge).

---

### 1-setup-all.bat ⭐ (M22)
**Purpose:** One-click complete setup - creates entire trading workbook

**What it does:**
1. ✅ Creates Excel workbook (TradingPlatform.xlsm)
2. ✅ Enables VBA project access (registry setting)
3. ✅ Imports all VBA modules (TFTypes, TFHelpers, TFEngine, TFTests)
4. ✅ **Creates 5 production UI worksheets** (M22 new feature)
5. ✅ Initializes database (trading.db)
6. ✅ Configures named ranges
7. ✅ Runs smoke tests

**Result:** Fully functional trading workbook with:
- **Dashboard** - Portfolio overview and navigation hub
- **Position Sizing** - Calculate shares/contracts for trades
- **Checklist** - 6-item entry validation with GREEN/YELLOW/RED banner
- **Heat Check** - Portfolio and bucket heat cap verification
- **Trade Entry** - Full 5-gate trade decision workflow

**Estimated Time:** 3-5 minutes

**Requirements:**
- Excel installed on Windows
- tf-engine.exe in current directory
- VBA modules in ../excel/vba/

---

### create-ui-worksheets.vbs (M22)
**Purpose:** Automated UI worksheet generator

**What it does:**
- Creates 5 production worksheets with complete UI
- Adds buttons, dropdowns, checkboxes (ActiveX controls)
- Formats cells, colors, borders
- Sets up result display areas
- Configures navigation between sheets

**Called by:** 1-setup-all.bat (Step 5/8)

**Can also run standalone:**
```cmd
cscript //nologo create-ui-worksheets.vbs TradingPlatform.xlsm
```

**Architecture:** Pure VBScript - no Go involvement during UI generation

---

### 2-update-vba.bat
**Purpose:** Update VBA modules only (without recreating workbook)

**Use when:** You've modified VBA .bas files and want to reimport them

---

### 3-run-integration-tests.bat
**Purpose:** Run Phase 3 integration test suite

**What it tests:** VBA → Go engine integration for all commands

---

### 4-run-tests.bat
**Purpose:** Run complete automated test suite

**What it tests:** CLI commands, VBA tests, integration tests

---

## Legacy Scripts (Deprecated)

The following scripts are **deprecated** as of M22 (2025-10-28). Use `1-setup-all.bat` instead.

### windows-import-vba.bat [DEPRECATED]
**Purpose:** Import VBA modules from `.bas` files into Excel workbook

**Status:** ⚠️ Deprecated - use `1-setup-all.bat` instead

**Requirements:**
- TradingPlatform.xlsm exists in current directory
- VBA modules exist in ../excel/vba/
- "Trust access to VBA project object model" enabled in Excel

**What it does:**
1. Creates temporary VBScript
2. Opens Excel workbook programmatically
3. Removes old modules if present (TFTypes, TFHelpers, TFEngine, TFTests)
4. Imports new modules from .bas files
5. Saves and closes workbook

**Troubleshooting:**
- If fails: Check Excel Trust Center settings
- If still fails: Manually import via VBA Editor > File > Import File

---

### windows-init-database.bat [DEPRECATED]
**Purpose:** Initialize trading.db with schema and default settings

**Status:** ⚠️ Deprecated - use `1-setup-all.bat` instead

**What it does:**
1. Checks if tf-engine.exe exists
2. Backs up existing database if present
3. Runs `tf-engine.exe init` to create schema
4. Verifies database with `get-settings` command

**Default Settings Created:**
- Equity (E): $10,000
- Risk per trade (r): 0.75%
- Portfolio heat cap: 4%
- Bucket heat cap: 1.5%
- Stop multiple (K): 2

**Outputs:**
- `trading.db` - SQLite database file
- `trading.db.backup_YYYYMMDD_HHMMSS` - Backup if reinitialized

---

### run-tests.bat
**Purpose:** Automated CLI testing (smoke tests + environment checks)

**Tests Performed:**
1. Engine version check
2. Database existence
3. Get settings from database
4. Position sizing calculation
5. Checklist evaluation
6. Heat check
7. Import candidates
8. List candidates
9. VBA module files exist
10. Test data files exist
11. Directory writable (for logs)

**Outputs:**
- `test-results.txt` - Test report
- Console output with pass/fail status

**Note:** VBA unit tests and integration tests must be run manually (see WINDOWS_TESTING.md)

---

### WINDOWS_TESTING.md (Comprehensive Testing Guide)
**Purpose:** Step-by-step manual testing procedures for M21

**Sections:**
- Phase 1: Pre-Test Setup (~10 min)
- Phase 2: Smoke Tests (~5 min)
- Phase 3: VBA Unit Tests (~10 min)
- Phase 4: Integration Tests (~15 min)
- Phase 5: Issue Reporting (if needed)
- Phase 6: Final Validation (~5 min)

**Total Estimated Time:** 45 minutes (if all passes on first try)

**This is the primary testing document for M21.**

---

### EXCEL_WORKBOOK_TEMPLATE.md (Workbook Structure)
**Purpose:** Specification for creating TradingPlatform.xlsm in Windows

**Worksheets Defined:**
1. Setup - Configuration and connection testing
2. VBA Tests - Run VBA unit tests
3. Position Sizing - Calculate position size
4. Checklist - Evaluate 6-item checklist
5. Heat Check - Verify portfolio/bucket heat
6. Trade Entry - Complete workflow with 5 gates
7. Candidates - Manage daily candidate tickers
8. Positions - View open positions and risk

**Includes:**
- Cell layouts
- Named ranges
- Button placements
- VBA code snippets for buttons
- Formatting standards
- Dropdown lists
- Initial setup checklist

**Note:** Workbook must be manually created in Windows (cannot create .xlsm in Linux)

---

### test-data/ (Sample JSON Responses)
**Purpose:** Validated JSON examples for VBA unit testing

**Contains 21 files:**
- Position sizing responses (stock, opt-delta-atr, opt-maxloss)
- Checklist responses (GREEN, YELLOW, RED)
- Heat check responses
- Timer check responses
- Candidate responses
- Cooldown responses
- Position responses
- Settings response

**Usage:** VBA unit tests load these files to verify JSON parsing correctness

---

## VBA Modules (in ../excel/vba/)

### TFTypes.bas (283 lines)
**Type definitions for all JSON response structures**

Key types:
- TFSizingResult
- TFChecklistResult
- TFHeatResult
- TFTimerResult
- TFCandidate types
- TFCooldown types
- TFPosition types
- TFSettings
- TFSaveDecisionResult
- TFCommandResult (generic wrapper)

### TFHelpers.bas (593 lines)
**JSON parsing, logging, validation, formatting**

Key functions:
- ExtractJSONValue() - Simple JSON parser
- ParseXXXJSON() - Typed parsers for all responses
- GenerateCorrelationID() - Create unique tracking IDs
- LogMessage() - Write to TradingSystem_Debug.log
- ValidateTicker(), ValidatePositiveNumber()
- FormatCurrency(), FormatPercent(), FormatTimestamp()

### TFEngine.bas (539 lines)
**Engine communication layer via shell execution**

Key functions:
- ExecuteCommand() - Core shell execution
- Engine_Size() - Position sizing
- Engine_Checklist() - Checklist evaluation
- Engine_Heat() - Heat management
- Engine_CheckTimer() - Impulse brake check
- Engine_ImportCandidates() - Import tickers
- Engine_SaveDecision() - Save with 5 gates
- Plus: candidates, cooldowns, positions, settings management

### TFTests.bas (689 lines)
**VBA unit tests**

14 tests:
- 6 JSON parsing tests
- 3 helper function tests
- 2 validation tests
- 2 formatting tests
- 1 shell execution test

Test runner: `RunAllTests()` - Outputs to "VBA Tests" worksheet

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────┐
│  Excel UI (TradingPlatform.xlsm)                    │
│  - Worksheets with inputs/outputs                   │
│  - Buttons trigger VBA functions                    │
│  - Named ranges for configuration                   │
└─────────────────┬───────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────┐
│  VBA Bridge (TFEngine.bas)                          │
│  - ExecuteCommand() via WScript.Shell.Exec          │
│  - Reads stdout (JSON), stderr (errors)             │
│  - Passes correlation IDs                           │
└─────────────────┬───────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────┐
│  tf-engine.exe (Go Backend)                         │
│  - All business logic                               │
│  - 5 hard gates enforcement                         │
│  - Outputs JSON to stdout                           │
│  - Logs to tf-engine.log                            │
└─────────────────┬───────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────┐
│  SQLite Database (trading.db)                       │
│  - Settings, candidates, decisions                  │
│  - Positions, cooldowns                             │
│  - Single source of truth                           │
└─────────────────────────────────────────────────────┘
```

**Key Principles:**
- **Thin VBA** - No business logic in VBA, just shell execution + JSON parsing
- **Engine-first** - All rules enforced in Go backend
- **Correlation IDs** - Every call tracked across VBA and Go logs
- **Fail loudly** - Errors never silently ignored

---

## Correlation ID Flow

```
User clicks button in Excel
    ↓
VBA generates: "20251027-143052-7A3F"
    ↓
VBA logs: "Calling tf-engine with corr_id 20251027-143052-7A3F"
    → TradingSystem_Debug.log
    ↓
VBA passes: --corr-id 20251027-143052-7A3F
    ↓
Go logs: {"corr_id": "20251027-143052-7A3F", "msg": "Processing..."}
    → tf-engine.log
    ↓
Go returns JSON to stdout
    ↓
VBA logs: "Received response for corr_id 20251027-143052-7A3F"
    → TradingSystem_Debug.log
    ↓
Excel displays: "✅ Success (corr_id: 20251027-143052-7A3F)"
```

**Cross-referencing issues:**
1. User reports issue with correlation ID
2. Search TradingSystem_Debug.log for ID → VBA side
3. Search tf-engine.log for ID → Go side
4. Compare timestamps and trace data flow

---

## Testing Strategy

### Automated Tests (run-tests.bat)
- ✅ CLI functionality
- ✅ Database operations
- ✅ File existence
- ✅ Environment setup

### VBA Unit Tests (Excel - TFTests.bas)
- ✅ JSON parsing for all response types
- ✅ Helper functions (validation, formatting)
- ✅ Shell execution
- ⏸️ Manual execution required

### Integration Tests (WINDOWS_TESTING.md Phase 4)
- ✅ Position sizing workflow
- ✅ Checklist evaluation
- ✅ Heat management
- ✅ Save decision (happy path + gate rejections)
- ⏸️ Manual execution required

**Test Coverage:**
- Automated: ~40% (CLI + environment)
- VBA Unit: ~30% (parsing + helpers)
- Integration: ~30% (end-to-end workflows)

---

## Deployment Checklist

**Before M21 Testing:**
- [ ] Copy windows/ folder to Windows PC
- [ ] Verify tf-engine.exe (12 MB, Windows binary)
- [ ] Verify all batch scripts present
- [ ] Verify test-data/ folder has 21 JSON files
- [ ] Verify ../excel/vba/ has 4 .bas files

**During M21 Setup:**
- [ ] Create TradingPlatform.xlsm workbook
- [ ] Enable "Trust access to VBA project object model"
- [ ] Run windows-import-vba.bat
- [ ] Run windows-init-database.bat
- [ ] Verify modules imported (Alt+F11 in Excel)
- [ ] Verify trading.db created

**M21 Testing:**
- [ ] Run run-tests.bat (automated tests)
- [ ] Run VBA unit tests in Excel
- [ ] Complete integration tests per WINDOWS_TESTING.md
- [ ] Verify all correlation IDs in logs
- [ ] Document any issues with correlation IDs

**M21 Sign-Off:**
- [ ] All automated tests passed
- [ ] All VBA unit tests passed (14/14)
- [ ] All integration tests passed
- [ ] Logs cross-reference correctly
- [ ] Ready for Phase E (Hardening & Release)

---

## Troubleshooting Quick Reference

### Issue: VBA Import Fails
**Fix:** File > Options > Trust Center > Enable "Trust access to VBA project object model"

### Issue: Engine Not Found
**Fix:** Verify tf-engine.exe in same folder as TradingPlatform.xlsm

### Issue: Database Init Fails
**Fix:** Run as Administrator, check write permissions

### Issue: Tests Timeout
**Fix:** Increase COMMAND_TIMEOUT_SECONDS in TFEngine.bas

### Issue: JSON Parsing Errors
**Fix:** Verify engine outputs well-formed JSON (test CLI manually)

### Issue: Correlation IDs Missing
**Fix:** Check file write permissions for TradingSystem_Debug.log

**Full Troubleshooting:** See WINDOWS_TESTING.md Appendix

---

## Log Files

**TradingSystem_Debug.log** (VBA side)
- Location: Same directory as TradingPlatform.xlsm
- Format: `[TIMESTAMP] [LEVEL] [CORR_ID] Message`
- Auto-rotates at 5 MB
- Use: VBA operations, UI events

**tf-engine.log** (Go side)
- Location: Same directory as tf-engine.exe
- Format: JSON structured logs with correlation IDs
- Rotation: Managed by Go logger
- Use: Business logic, database operations, gate validations

---

## Next Steps

### For M21 (Windows Integration Validation)
1. Follow WINDOWS_TESTING.md step-by-step
2. Complete all 6 testing phases
3. Document any issues with correlation IDs
4. Create issue reports for failures
5. Sign off when all tests pass

### After M21 (Phase E - Hardening & Release)
1. BDD test suite (full Gherkin scenarios green)
2. Error message refinement
3. Final packaging (ZIP with exe, DB, workbook, docs)
4. Documentation updates
5. Clean Windows box test (< 5 min setup)

---

## Support & Documentation

**Primary Documents:**
- `WINDOWS_TESTING.md` - M21 testing procedures
- `EXCEL_WORKBOOK_TEMPLATE.md` - Workbook structure
- `../excel/VBA_MODULES_README.md` - VBA module documentation
- `../WHY.md` - Core philosophy
- `../DEVELOPMENT_PHILOSOPHY.md` - Development approach
- `../Trading-Engine-v3_Step-by-Step-Plan.md` - Overall plan

**Issue Reporting:**
- Use template in WINDOWS_TESTING.md Phase 5
- Always include correlation ID
- Attach log excerpts from both VBA and Go logs
- Screenshots of Excel showing issue

---

## Version Information

**Package:** M20 - Windows Integration Package
**Created:** 2025-10-27
**Binary:** tf-engine.exe version 3.0.0-dev
**VBA Modules:** M19 (2025-10-27)
**Target:** Windows 10/11 with Excel desktop
**Status:** ✅ Ready for M21 testing

---

**This package contains everything needed to test VBA ↔ Go engine integration on Windows.**

**Remember:** This is a discipline enforcement system. Testing must be thorough. No shortcuts.

Code serves discipline. Discipline does not serve code.
