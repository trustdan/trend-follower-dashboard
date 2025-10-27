# Excel Trading Workbook - Build Instructions

**One-command build for `TrendFollowing_TradeEntry.xlsm`**

---

## Quick Start

```cmd
cd C:\Users\Dan\excel-trading-dashboard
BUILD_WITH_PYTHON.bat
```

That's it! The script handles everything automatically.

---

## What Gets Created

### Workbook: `TrendFollowing_TradeEntry.xlsm`

**8 Sheets**:
- TradeEntry (main UI)
- Presets (5 FINVIZ screeners)
- Buckets (6 sector buckets)
- Candidates (daily tickers)
- Decisions (trade log)
- Positions (open positions)
- Summary (settings)
- Control (hidden state)

**9 VBA Modules**:
- PQ_Setup, Python_Run, Setup, TF_Data, TF_Presets, TF_UI, TF_Utils
- ThisWorkbook (code)
- TradeEntry sheet (event handlers)

**5 Tables**:
- tblPresets, tblBuckets, tblCandidates, tblDecisions, tblPositions

**7 Named Ranges**:
- Equity_E, RiskPct_r, StopMultiple_K, HeatCap_H_pct, BucketHeatCap_pct, AddStepN, EarningsBufferDays

---

## Build Process (Automated)

1. ✅ Check/create Python venv
2. ✅ Check/install pywin32
3. ✅ Configure Excel Trust Center
4. ✅ Kill stuck Excel processes
5. ✅ Delete old workbook
6. ✅ Create new Excel workbook
7. ✅ Import VBA modules
8. ✅ Run TF_Data.EnsureStructure (creates sheets/tables/data)
9. ✅ Apply sheet code (TradeEntry events)
10. ✅ Save and close

**Time**: ~10-15 seconds

---

## Requirements

- **Python 3.7+** installed (check: `py -3 --version`)
- **Excel 2016+** with VBA enabled
- **Windows** (uses COM automation)

Everything else is auto-installed.

---

## Troubleshooting

### "Python not found"
Install Python from python.org, then retry.

### "Excel blocked VBProject access"
Excel → File → Options → Trust Center → Trust Center Settings → Macro Settings
→ Check "Trust access to the VBA project object model"

### "Proxy error" / "beautifulsoup not found"
**Ignore it!** That's Excel's Python (not used). See `IGNORE_EXCEL_PYTHON_ERROR.md`

### Build completes but sheets are empty
This was fixed. Delete old workbook and rebuild:
```cmd
del TrendFollowing_TradeEntry.xlsm
BUILD_WITH_PYTHON.bat
```

### "Sheet11" appears instead of TradeEntry
This was fixed. Rebuild with latest version.

---

## After Building

### Open the Workbook
```cmd
start TrendFollowing_TradeEntry.xlsm
```

### Verify Modules (Alt+F11)
Should see:
- 7 standard modules in "Modules" folder
- ThisWorkbook with code
- TradeEntry sheet with code (under "Microsoft Excel Objects")

### Check Data
- **Presets sheet**: 5 rows (TF_BREAKOUT_LONG, etc.)
- **Buckets sheet**: 6 rows (Tech/Comm, Consumer, etc.)
- **Summary sheet**: Labels + values (Equity_E = 10000)

### Compile Check
Alt+F11 → Debug → Compile VBAProject
Should compile with no errors.

---

## Customization

### Change Account Settings
Open workbook → Summary sheet → Edit values in column B

### Add More Presets
Presets sheet → Add rows to tblPresets table

### Modify VBA Code
1. Edit files in `VBA/` folder
2. Run `BUILD_WITH_PYTHON.bat` to rebuild

---

## Files You Can Ignore

- `old_*.bat` - Archived VBScript versions (kept for reference)
- `*.md` files - Documentation
- `logs/` - Build logs
- `venv/` - Python virtual environment (auto-created)

---

## Files You Need

- **BUILD_WITH_PYTHON.bat** ← Main build script
- **import_to_excel.py** ← Python import logic
- **VBA/*.bas** ← VBA modules
- **VBA/*.cls** ← Class modules
- **scripts/setup_venv.bat** ← Venv setup (auto-called)

---

## Optional Features

### Python FINVIZ Scraper
If you want the automated ticker scraping feature:

1. Run: `scripts\setup_venv.bat`
2. Installs: requests, beautifulsoup4, lxml
3. Use button in Excel: "Fetch Screened Tickers"

**Not required for manual FINVIZ workflow.**

---

## Support

See detailed troubleshooting:
- `QUICK_START.md` - Step-by-step guide
- `FINAL_FIXES.md` - Recent fixes
- `IGNORE_EXCEL_PYTHON_ERROR.md` - Excel Python vs. regular Python
- `BUILD_COMPLETE.md` - Full journey summary

---

## Summary

**Build command**: `BUILD_WITH_PYTHON.bat`
**Result**: Fully functional trading workbook
**Time**: ~15 seconds
**Manual steps**: Zero

✅ Ready to trade!
