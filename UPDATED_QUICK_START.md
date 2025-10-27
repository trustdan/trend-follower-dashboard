# 🚀 QUICK START - Updated Workflow (2 Steps!)

## What's New: Fully Automated Setup!

The system now has **automatic initialization** that runs when you first open the workbook. No manual macro execution needed!

## Complete Setup (2 Steps, 3 Minutes Total)

### Step 1: Build Workbook (30 seconds)

```cmd
BUILD.bat
```

**What it does**:
- ✅ Installs pywin32 if needed
- ✅ Clears COM cache
- ✅ Kills stuck Excel processes
- ✅ Creates TrendFollowing_TradeEntry.xlsm
- ✅ Imports all 11 VBA modules
- ✅ Imports ThisWorkbook with auto-setup code

**Result**: `TrendFollowing_TradeEntry.xlsm` created

---

### Step 2: Open Workbook - Setup Runs Automatically! (2 minutes)

1. **Double-click** `TrendFollowing_TradeEntry.xlsm`
2. **Click** "Enable Content" (if prompted)
3. **Wait** for automatic setup:
   - Message: "Welcome to the Trading System! Running initial setup now..."
   - Setup creates all sheets, tables, and UI automatically
   - **Setup sheet opens** with complete instructions

4. **Follow Setup sheet instructions** to add 6 checkboxes (2 minutes)

**That's it!** You're done!

---

## What Happens Automatically

When you open the workbook for the first time:

```
Workbook_Open event fires
    ↓
Checks: Does "Setup" sheet exist?
    ↓ NO (first time)
Shows welcome message
    ↓
Runs Setup.RunInitialSetup()
    ↓
Creates all 8 sheets
    ↓
Creates all 5 tables
    ↓
Seeds default presets & buckets
    ↓
Builds TradeEntry UI (labels, buttons, dropdowns)
    ↓
Creates Setup sheet with instructions
    ↓
Opens Setup sheet
    ↓
DONE! Just add 6 checkboxes
```

On subsequent opens:
- Skips setup (already done)
- Goes directly to TradeEntry sheet
- Binds controls
- Ready to trade!

---

## The Only Manual Step: Add 6 Checkboxes

**Why manual?** Excel COM automation cannot reliably create Form Control checkboxes programmatically.

**Instructions** (shown on Setup sheet):

1. Go to **TradeEntry** sheet
2. **Developer tab** → Insert → **Check Box (Form Control)**
3. Draw 6 checkboxes next to cells B21:B26
4. For each checkbox:
   - Right-click → **Format Control**
   - Cell link: `$C$20`, `$C$21`, `$C$22`, `$C$23`, `$C$24`, `$C$25`
   - Click OK
5. Delete checkbox text (labels already in column A)

**Time**: 2 minutes

---

## Quick Test Workflow

After adding checkboxes:

1. Go to **TradeEntry** sheet
2. Select **Preset**: TF_BREAKOUT_LONG
3. Click **"Import Candidates"** → Paste: `AAPL, MSFT, NVDA`
4. Select **Ticker**: AAPL
5. Enter: **Entry**=180, **ATR N**=1.50, **K**=2
6. **Check all 6 boxes** → Click **"Evaluate"** → See **GREEN**! ✅
7. Click **"Recalc Sizing"** → See calculated shares
8. Wait 2 minutes → Click **"Save Decision"** → Trade logged!

---

## Setup Sheet Features

The auto-created Setup sheet includes:

### Status Checklist
- ✓ Workbook created
- ✓ VBA modules imported
- ✓ Data structure created
- ✓ TradeEntry UI built
- → Add 6 checkboxes (your only task!)

### Utility Buttons
- **Rebuild TradeEntry UI** - Recreates UI if needed
- **Test Python Integration** - Checks if auto-scraping available
- **Clear Old Candidates** - Removes old ticker imports

### Complete Instructions
- Step-by-step checkbox setup
- Quick test workflow
- Key settings reference
- Documentation links

---

## Comparison: Old vs New Workflow

### Old Way (Manual):
1. Run BUILD.bat
2. Open workbook
3. Press Alt+F11
4. Press Ctrl+G
5. Type: `TF_Data.EnsureStructure`
6. Press Enter
7. Type: `TF_UI_Builder.BuildTradeEntryUI`
8. Press Enter
9. Close VBA Editor
10. Add checkboxes

**Total**: 10 steps, 5 minutes

### New Way (Automated):
1. Run BUILD.bat
2. Open workbook (setup runs automatically)
3. Add checkboxes

**Total**: 3 steps, 3 minutes

**Time saved**: 40% faster!

---

## Troubleshooting

### "Setup didn't run automatically"
**Fix**: Run manually:
1. Alt+F11 (VBA Editor)
2. Ctrl+G (Immediate Window)
3. Type: `Setup.RunInitialSetup`
4. Press Enter

### "I need to rebuild the UI"
**Fix**: Click **"Rebuild TradeEntry UI"** button on Setup sheet

### "Can't find Setup sheet"
**Fix**: It was deleted - recreate it:
1. Alt+F11
2. Ctrl+G
3. Type: `Setup.CreateSetupSheet`

### "Import button asks for manual paste"
**This is normal!** Python auto-scraping requires:
- Microsoft 365 Insider
- Python in Excel enabled

Manual import works perfectly - just paste tickers from FINVIZ.

---

## File Structure

```
excel-trading-workflow/
│
├── BUILD.bat                    ← Run this!
├── build_workbook_simple.py
│
├── VBA/                         ← 11 modules
│   ├── TF_Utils.bas
│   ├── TF_Data.bas
│   ├── TF_UI.bas
│   ├── TF_Presets.bas           ← Smart import (auto/manual)
│   ├── TF_Python_Bridge.bas
│   ├── TF_UI_Builder.bas
│   ├── Setup.bas                ← NEW! Auto-setup
│   ├── ThisWorkbook.cls         ← NEW! Auto-runs setup
│   └── Sheet_TradeEntry.cls
│
├── Python/                      ← 3 files
│   ├── finviz_scraper.py
│   ├── heat_calculator.py
│   └── requirements.txt
│
└── TrendFollowing_TradeEntry.xlsm  ← Generated
```

---

## What Changed

### New Features:
- ✅ **Auto-setup on first open** - No manual macro execution
- ✅ **Setup sheet with instructions** - Always know what to do next
- ✅ **Utility buttons** - Rebuild UI, test Python, clear data
- ✅ **Smart import** - Auto-detects Python, falls back to manual
- ✅ **Fixed encoding issues** - Checkbox labels now display correctly

### Updated Components:
- **Setup.bas** - NEW! One-click initialization
- **ThisWorkbook.cls** - Now runs setup automatically
- **TF_Presets.bas** - Smart import with Python detection
- **build_workbook_simple.py** - Updated instructions
- **All documentation** - Reflects new automated workflow

---

## Summary

**Old workflow**: Build → Open → Manual VBA commands → Add checkboxes
**New workflow**: Build → Open (auto-setup!) → Add checkboxes

Everything between "Build" and "Add checkboxes" is now **fully automated**!

---

## Next Steps

1. ✅ Run `BUILD.bat`
2. ✅ Open workbook (setup runs automatically)
3. ✅ Add 6 checkboxes (follow Setup sheet instructions)
4. ✅ Start trading!

**Total time**: 3 minutes

**Questions?** Check the Setup sheet - it has everything you need!

---

**Ready?** → Run `BUILD.bat` now!
