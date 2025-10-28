# Trading Platform v3 - Quick Start Guide

**Last Updated:** 2025-10-28 (M23)
**Latest:** FINVIZ scraper working + Interactive mode with ASCII art

---

## The Problem We Fixed

The original `1-setup-all.bat` script was failing when trying to create checkboxes using Excel automation. This caused the Checklist, Heat Check, and Trade Entry worksheets to fail or be incomplete.

## The Solution

We created a **simplified setup** that uses **TRUE/FALSE dropdowns** instead of checkboxes. This avoids Excel's fragile OLE automation while keeping the same functionality.

---

## Quick Setup (3 Minutes)

### Prerequisites

1. Excel installed (2013 or later)
2. All files from `/home/kali/excel-trading-platform` copied to Windows

### Step 1: Run the Simple Setup

```cmd
cd C:\Users\Dan\excel-trading-dashboard\windows
setup-simple.bat
```

This will:
- ✅ Delete old `TradingPlatform.xlsm` (if exists)
- ✅ Create new workbook with simplified UI
- ✅ Import all VBA modules (TFTypes, TFHelpers, TFEngine, TFTests)
- ✅ Initialize database (`trading.db`)
- ✅ Run smoke tests

### Step 2: Open & Test the Workbook

1. **Open** `TradingPlatform.xlsm`
2. **Enable macros** when prompted (click "Enable Content")
3. **Go to "VBA Tests" sheet**
4. **Click "Run All Tests" button**
5. **Verify** all tests pass (should see green PASS cells)

---

## What You Get

The workbook contains **7 sheets**:

### 1. **Setup** (Gray)
- Configuration settings
- Engine path: `.\tf-engine.exe`
- Database path: `.\trading.db`

### 2. **Dashboard** (Blue)
- Portfolio status overview
- Current equity, heat caps
- Today's candidates (from FINVIZ)

### 3. **Position Sizing** (Green)
- Calculate shares/contracts
- Inputs: Ticker, Entry, ATR, K, Method
- Supports stocks and options (delta-ATR, max loss)
- Click "Calculate" button to run

### 4. **Checklist** (Orange)
- **6-item checklist** evaluation
- **Uses TRUE/FALSE dropdowns** (not checkboxes)
- Select TRUE or FALSE for each item from dropdown
- Click "Evaluate" to get GREEN/YELLOW/RED banner

**The 6 Checklist Items:**
1. Ticker from today's FINVIZ preset
2. Trend alignment confirmed
3. Adequate volume and spread
4. TradingView setup confirmation
5. No earnings in next 7 days
6. Trade thesis documented

### 5. **Heat Check** (Red)
- Portfolio heat management
- Check if trade would exceed portfolio cap (4%)
- Check if trade would exceed bucket cap (1.5%)
- Inputs: Ticker, Risk Amount, Bucket

### 6. **Trade Entry** (Purple)
- **5 Hard Gates** enforcement
- Save GO or NO-GO decisions
- Gates check:
  1. Banner GREEN (from checklist)
  2. In today's candidates
  3. Impulse brake (2-minute timer)
  4. Bucket cooldown
  5. Heat caps (portfolio + bucket)

### 7. **VBA Tests** (Gray)
- Automated unit tests
- Click "Run All Tests" to verify everything works
- Tests JSON parsing, engine communication, etc.

---

## Key Differences from Original

### Original (Broken)
- Used OLE checkboxes (Forms.CheckBox.1)
- Excel automation would crash at line 264
- Checklist, Heat Check, Trade Entry incomplete

### New Simplified Version
- **Uses TRUE/FALSE dropdowns** instead
- No fragile OLE automation
- All sheets complete and functional
- Same business logic, just different UI

---

## Using the Checklist Sheet

**Old way (broken):** Click checkboxes
**New way (works):** Select TRUE/FALSE from dropdown

**Example:**

| Checklist Item | Your Selection |
|----------------|----------------|
| 1. Ticker from today's FINVIZ preset | **TRUE** ← select from dropdown |
| 2. Trend alignment confirmed | **TRUE** |
| 3. Adequate volume and spread | **FALSE** |
| 4. TradingView setup confirmation | **TRUE** |
| 5. No earnings in next 7 days | **TRUE** |
| 6. Trade thesis documented | **FALSE** |

Then click **"Evaluate"** button → You'll get a **YELLOW** banner (2 missing)

---

## Testing the System

### Step 1: Test VBA Modules

1. Open `TradingPlatform.xlsm`
2. Go to "VBA Tests" sheet
3. Click "Run All Tests"
4. Should see PASS for all tests

### Step 2: Test Position Sizing

1. Go to "Position Sizing" sheet
2. Enter:
   - Ticker: `AAPL`
   - Entry Price: `180`
   - ATR (N): `1.5`
   - K Multiple: `2`
   - Method: `stock`
3. Click "Calculate"
4. Should see results: Shares, Stop, Risk Dollars

### Step 3: Test Checklist

1. Go to "Checklist" sheet
2. Enter Ticker: `AAPL`
3. Select TRUE for all 6 items (use dropdowns)
4. Click "Evaluate"
5. Should see **GREEN** banner

### Step 4: Run Integration Tests

```cmd
cd C:\Users\Dan\excel-trading-dashboard\windows
4-run-integration-tests.bat
```

Should see all tests pass (including heat tests now!)

---

## Troubleshooting

### "Macros are disabled"
- Click "Enable Content" button in Excel
- Or: File → Options → Trust Center → Trust Center Settings → Macro Settings → Enable all macros

### "VBA modules not found"
Make sure folder structure is:
```
C:\Users\Dan\excel-trading-dashboard\
├── excel\
│   └── vba\
│       ├── TFTypes.bas
│       ├── TFHelpers.bas
│       ├── TFEngine.bas
│       └── TFTests.bas
└── windows\
    ├── tf-engine.exe
    ├── setup-simple.bat
    └── create-workbook-manual-ui.vbs
```

### "tf-engine.exe not found"
- Make sure `tf-engine.exe` is in the `windows\` folder
- Check that you copied the **NEW version** (built Oct 28, 26MB)
- Old version doesn't have the `heat` command

### Tests fail with "unknown command heat"
- You're using old `tf-engine.exe`
- Copy the new one from `/home/kali/excel-trading-platform/windows/tf-engine.exe`
- Should be 26MB, built Oct 28 11:01

---

## What's Next?

Once all tests pass:

1. **Import candidates** from FINVIZ (M23 - New!):
   ```cmd
   # Interactive mode - Beautiful ASCII art & menus
   import-candidates.bat

   # Or auto mode - No prompts, uses defaults
   import-candidates-auto.bat
   ```

   **What you'll see:**
   - 🎨 Epic ASCII banner
   - 📊 Preset selection (TF-Breakout-Long, TF-Momentum-Uptrend, etc.)
   - ⚡ Progress bars & animations
   - 💰 Ticker preview (CCJ, W, SOFI, FTAI, etc.)
   - ✅ Success confirmation

2. **Start using the system:**
   - Dashboard → See portfolio status & today's candidates
   - Position Sizing → Calculate trade size
   - Checklist → Evaluate setup (must be GREEN)
   - Heat Check → Verify caps not exceeded
   - Trade Entry → Save GO/NO-GO decision (5 gates)

3. **The 5 Hard Gates enforce discipline:**
   - Gate 1: Checklist must be GREEN
   - Gate 2: Ticker must be in today's candidates
   - Gate 3: Must wait 2 minutes (impulse brake)
   - Gate 4: Bucket must not be on cooldown
   - Gate 5: Heat caps must not be exceeded

**If any gate fails → Trade is rejected. No exceptions.**

---

## Files You Need

Make sure these files are present in `C:\Users\Dan\excel-trading-dashboard\windows\`:

**Core Files:**
- ✅ `tf-engine.exe` (32MB, M23 with FINVIZ scraper + cookie jar)
- ✅ `1-setup-all.bat` (one-click complete setup)
- ✅ `create-workbook-manual-ui.vbs` (workbook generator with dropdowns)

**Candidate Import (M23 - New!):**
- ✅ `import-candidates.bat` (interactive mode launcher)
- ✅ `import-candidates-auto.bat` (auto mode launcher)

**Testing:**
- ✅ `3-run-integration-tests.bat` (integration test suite)
- ✅ `4-run-tests.bat` (all automated tests)

---

## Support

- **Full testing guide:** `WINDOWS_TESTING.md`
- **Workbook template details:** `EXCEL_WORKBOOK_TEMPLATE.md`
- **VBA module docs:** `..\excel\VBA_MODULES_README.md`
- **Project docs:** `..\docs\README.md`

---

**Ready to trade with discipline!** 🎯
