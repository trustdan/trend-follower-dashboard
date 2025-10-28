# Trading Platform v3 - Setup Workflow

**Last Updated:** 2025-10-28
**Status:** Updated and streamlined

---

## Numbered Batch File Workflow

Run these scripts **in order** for a complete setup and validation:

### 1️⃣ **`1-setup-all.bat`** - Initial Setup

**Purpose:** Complete one-click setup from scratch

**What it does:**
- ✅ Checks prerequisites (tf-engine.exe, VBA source files)
- ✅ Enables VBA project access (registry)
- ✅ Creates TradingPlatform.xlsm with all 7 worksheets and 4 VBA modules
- ✅ Initializes trading.db database
- ✅ Runs smoke tests

**When to use:**
- First time setup
- After deleting TradingPlatform.xlsm to start fresh
- After major updates

**Runtime:** 2-3 minutes

**Output files:**
- `TradingPlatform.xlsm` (50-100 KB)
- `trading.db` (8-16 KB)
- `setup-all.log`

---

### 2️⃣ **`2-update-vba.bat`** - Update VBA Only

**Purpose:** Re-import VBA modules without recreating workbook

**What it does:**
- ✅ Backs up existing TradingPlatform.xlsm
- ✅ Removes old VBA modules
- ✅ Imports updated VBA modules from `../excel/vba/`
- ✅ Preserves all worksheets and data

**When to use:**
- After editing VBA .bas files
- After pulling VBA updates from Git
- VBA code changed but workbook structure didn't

**Runtime:** 10-20 seconds

**Prerequisites:**
- TradingPlatform.xlsm must exist

---

### 3️⃣ **`3-run-integration-tests.bat`** - Integration Tests

**Purpose:** Automated execution of 25 Phase 4 integration tests

**What it does:**
- ✅ Imports TFIntegrationTests.bas module
- ✅ Runs automated tests via Excel macro
- ✅ Tests all 4 workflows:
  - Position Sizing
  - Checklist Evaluation
  - Heat Management
  - Save Decision (5 gates)
- ✅ Creates log file
- ✅ Populates results in "Integration Tests" worksheet

**When to use:**
- After initial setup
- Before releasing
- After making changes to engine or VBA

**Runtime:** 2-5 seconds

**Output:**
- `logs/integration-tests-YYYYMMDD-HHMMSS.log`
- Integration Tests worksheet in workbook

---

### 4️⃣ **`4-run-tests.bat`** - CLI Smoke Tests

**Purpose:** Quick validation of CLI functionality

**What it does:**
- ✅ Tests engine version
- ✅ Tests database access
- ✅ Tests position sizing
- ✅ Tests checklist evaluation
- ✅ Tests heat management
- ✅ Tests all major CLI commands

**When to use:**
- Quick validation after setup
- Before starting a trading session
- After updating tf-engine.exe

**Runtime:** 5-10 seconds

**Output:**
- `test-results.txt`
- Console output (PASS/FAIL summary)

---

## Typical Usage Scenarios

### Scenario 1: First Time Setup

```cmd
cd C:\Users\Dan\excel-trading-dashboard\windows
1-setup-all.bat
3-run-integration-tests.bat
4-run-tests.bat
```

**Result:** Fully configured system with all tests passing

---

### Scenario 2: VBA Code Update

```cmd
cd C:\Users\Dan\excel-trading-dashboard\windows
2-update-vba.bat
3-run-integration-tests.bat
```

**Result:** Updated VBA modules with validation

---

### Scenario 3: Fresh Start (Reset Everything)

```cmd
cd C:\Users\Dan\excel-trading-dashboard\windows
del TradingPlatform.xlsm
del trading.db
1-setup-all.bat
3-run-integration-tests.bat
```

**Result:** Brand new setup from scratch

---

### Scenario 4: Daily Pre-Trading Check

```cmd
cd C:\Users\Dan\excel-trading-dashboard\windows
4-run-tests.bat
```

**Result:** Quick validation that everything works

---

## File Naming Convention

| Pattern | Purpose | Examples |
|---------|---------|----------|
| `N-*.bat` | Numbered workflow scripts (run in order) | `1-setup-all.bat`, `2-update-vba.bat` |
| `windows-*.bat` | Utility scripts (standalone) | `windows-import-vba.bat` |
| `*.vbs` | VBScript automation (called by .bat) | `create-workbook-manual-ui.vbs` |
| `*.md` | Documentation | `QUICK_START.md`, `README.md` |
| `*.log` | Log files (generated) | `setup-all.log`, `tf-engine.log` |

---

## Changes from Previous Version

### October 28, 2025 Update

**Problem:** Original `1-setup-all.bat` failed when creating Excel checkboxes via OLE automation

**Solution:** Consolidated Steps 2, 4, and 5 into one robust VBScript

**Changes:**

| Old Approach | New Approach |
|--------------|--------------|
| Step 2: Create basic workbook | Step 2: Enable VBA access |
| Step 3: Enable VBA access | Step 3: Create full workbook (consolidated) |
| Step 4: Import VBA modules | (merged into Step 3) |
| Step 5: Create UI worksheets (FAILED) | (merged into Step 3) |
| Step 6: Configure Excel | (merged into Step 3) |
| Step 7: Initialize database | Step 4: Initialize database |
| Step 8: Run smoke tests | Step 5: Run smoke tests |

**UI Changes:**
- **Old:** Checkboxes (OLE controls) - fragile, broke during automation
- **New:** TRUE/FALSE dropdowns (cell validation) - reliable, works consistently

**File Changes:**
- **New:** `create-workbook-manual-ui.vbs` (consolidated script)
- **Removed:** Inline VBScript generation in Steps 2/4/5/6
- **Kept:** All other batch files (2-4) unchanged

---

## Troubleshooting

### `1-setup-all.bat` fails at Step 3

**Error:** "create-workbook-manual-ui.vbs not found"

**Fix:** Ensure `create-workbook-manual-ui.vbs` exists in windows folder

---

### Tests fail with "unknown command heat"

**Error:** `Error: unknown command "heat" for "tf-engine"`

**Fix:** You're using old tf-engine.exe. Copy the new version (26MB, Oct 28 build) from `/home/kali/excel-trading-platform/windows/tf-engine.exe`

---

### VBA Tests button doesn't run

**Error:** Click "Run All Tests" but nothing happens

**Fix:**
1. Enable macros: Click "Enable Content" yellow bar
2. Check macro settings: File → Options → Trust Center → Enable all macros
3. Verify VBA modules imported: Press Alt+F11, check for TFTests module

---

### Workbook has no UI (buttons/dropdowns missing)

**Error:** Sheets exist but no interactive elements

**Fix:**
1. Delete TradingPlatform.xlsm
2. Ensure `create-workbook-manual-ui.vbs` is present
3. Run `1-setup-all.bat` again
4. Check setup-all.log for errors

---

## Quick Reference

| Task | Command |
|------|---------|
| First time setup | `1-setup-all.bat` |
| Update VBA after code changes | `2-update-vba.bat` |
| Run full integration tests | `3-run-integration-tests.bat` |
| Quick smoke test | `4-run-tests.bat` |
| Reset everything | Delete .xlsm and .db, run `1-setup-all.bat` |
| Import VBA only | `windows-import-vba.bat TradingPlatform.xlsm` |
| Initialize DB only | `windows-init-database.bat` |

---

## Documentation Index

| File | Purpose |
|------|---------|
| `QUICK_START.md` | Quick setup guide (start here!) |
| `README_UI_FIX.md` | Technical explanation of dropdown vs checkbox change |
| `SETUP_WORKFLOW.md` | This file - numbered batch workflow |
| `WINDOWS_TESTING.md` | Complete testing guide |
| `EXCEL_WORKBOOK_TEMPLATE.md` | Workbook structure details |
| `README.md` | Windows package overview |

---

**Ready to set up!** Start with `1-setup-all.bat` and follow the numbered sequence.
