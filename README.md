# Excel Trading Workbook - Complete Automated Solution

**Version:** v2.0.0 | **Status:** ✅ Production Ready | **Setup Time:** 3 minutes

## 🎯 New in v2.1.0

- 🎉 **Comprehensive Logging System** - Debug log tracks everything automatically
- 🎉 **Enhanced Python Detection** - Detailed diagnostics for troubleshooting
- 🎉 **Improved Error Handling** - All operations logged with error details
- 🔧 **Fixed Checkbox Positioning** - Properly placed in column B (not overlapping text)
- 🔧 **Fixed Dropdown Binding** - Better error messages when tables missing
- 📚 **LOGGING_AND_DIAGNOSTICS.md** - Complete logging guide
- 📚 **TROUBLESHOOTING_CHECKLIST.md** - Quick problem-solving reference

## What's New in v2.0.0

- 🎉 **USER_GUIDE.md** - Comprehensive 15,000-word beginner guide (auto-opens on first launch!)
- 🎉 **Automatic Checkboxes** - No more manual creation (with fallback if needed)
- 🔧 **Fixed Unicode Issues** - All text displays correctly now
- 🔧 **Fixed Python Detection** - Modernized for Python in Excel 2023+
- 🔧 **Fixed Duplicate Buttons** - UI rebuild now works perfectly
- 📚 **CHANGELOG.md** - Track all version changes
- 📚 **DEVELOPMENT_LOG.md** - Technical notes for future development

## Overview

A **fully automated** trend-following trade entry system with:
- ✅ **FINVIZ integration** - Web scraping of screener results (NOT just permalinks)
- ✅ **Automated GO/NO-GO decisions** - 6-item checklist with GREEN/YELLOW/RED banner
- ✅ **Position sizing** - Stocks and 2 option methods (Delta-ATR, MaxLoss)
- ✅ **Heat caps** - Portfolio and bucket level risk management
- ✅ **Cooldown logic** - Auto-pause buckets after stop-outs
- ✅ **2-minute impulse brake** - Prevents FOMO entries
- ✅ **Python acceleration** - 5-10x faster import and calculations
- ✅ **Comprehensive Documentation** - USER_GUIDE.md explains everything in plain English

## Quick Start (3 Minutes) ⚡

### Step 1: Build the Workbook (30 seconds)

```cmd
BUILD.bat
```

The script will:
- Install pywin32 if needed
- Kill any stuck Excel processes
- Create `TrendFollowing_TradeEntry.xlsm`
- Import all 11 VBA modules
- Ready to open!

### Step 2: Open & Auto-Setup (2 minutes)

1. **Double-click** `TrendFollowing_TradeEntry.xlsm`
2. **Click "Enable Content"** (security warning)
3. **Wait** for auto-setup to complete:
   - Creates 8 worksheets
   - Creates 5 data tables
   - Builds TradeEntry UI
   - Creates 6 checkboxes (automated!)
   - **Opens USER_GUIDE.md** (read while exploring)
4. **Done!** System is ready to trade

### Step 3: Verify (30 seconds)

1. Go to **TradeEntry** sheet
2. Check rows 21-26 for **6 checkboxes**
3. If missing, follow instructions on **Setup** sheet

**That's it!** Read the USER_GUIDE.md that just opened for detailed instructions.

### 3. Test the Workflow

1. **Import candidates**:
   - Select Preset: "TF_BREAKOUT_LONG"
   - Click "Import Candidates" button
   - Paste: `AAPL, MSFT, NVDA, TSLA`
   - Click OK

2. **Enter trade**:
   - Select Ticker: AAPL
   - Entry Price: 180.00
   - ATR N: 1.50
   - K: 2

3. **Evaluate**:
   - Check all 6 checkboxes
   - Click "Evaluate" → Should see GREEN banner

4. **Size position**:
   - Click "Recalc Sizing" → See calculated shares/stop

5. **Save decision**:
   - Wait 2 minutes (impulse brake)
   - Click "Save Decision" → Trade logged!

## File Structure

```
excel-trading-workflow/
│
├── BUILD.bat                      # ← RUN THIS to build workbook
├── build_workbook.py              # Python build automation
│
├── VBA/                           # VBA modules (10 files)
│   ├── TF_Utils.bas              # Helper functions
│   ├── TF_Data.bas               # Data structure setup
│   ├── TF_UI.bas                 # Trading logic (evaluate, size, save)
│   ├── TF_Presets.bas            # FINVIZ integration
│   ├── TF_Python_Bridge.bas      # Python integration
│   ├── TF_UI_Builder.bas         # Automated UI creation
│   ├── ThisWorkbook.cls          # Workbook events
│   └── Sheet_TradeEntry.cls      # Sheet events
│
├── Python/                        # Python modules (3 files)
│   ├── finviz_scraper.py         # Web scraping engine
│   ├── heat_calculator.py        # Fast heat calculations
│   └── requirements.txt          # Dependencies
│
└── TrendFollowing_TradeEntry.xlsm # ← Generated workbook
```

## How FINVIZ Integration Works

### NOT Just Permalinks - Active Web Scraping

**FINVIZ Screener URLs**:
```
https://finviz.com/screener.ashx?v=211&s=ta_newhigh&ft=4
```

**What the Python scraper does**:
1. Sends HTTP request to FINVIZ with query params
2. Parses HTML table using BeautifulSoup
3. Extracts ticker symbols from table cells
4. Handles pagination (20 tickers per page, up to 10 pages)
5. Normalizes and dedupes tickers
6. Returns: `['AAPL', 'MSFT', 'NVDA', ...]`

**VBA Integration** (when Python available):
- Button click → VBA calls `finviz_scraper.fetch_finviz_tickers()`
- Python scrapes FINVIZ (5-10 seconds)
- Returns ticker array
- VBA writes to Candidates table
- **Fallback**: Manual paste still works

### Testing FINVIZ Scraper

```bash
cd Python
python finviz_scraper.py
```

Should output:
```
✅ Success! Found 47 tickers:
AAPL, MSFT, NVDA, TSLA, META, ...
```

## Workbook Structure

### 8 Sheets

1. **TradeEntry** - Main UI (all trading happens here)
2. **Presets** - FINVIZ query strings (5 default presets)
3. **Buckets** - Correlation groups with cooldown settings
4. **Candidates** - Daily ticker imports
5. **Decisions** - Complete trade log (20 fields)
6. **Positions** - Open positions tracker
7. **Summary** - Settings and named ranges
8. **Control** - Hidden helper sheet (impulse timer)

### 5 Tables

- `tblPresets` - 5 FINVIZ screener presets
- `tblBuckets` - 6 correlation buckets
- `tblCandidates` - Daily imported tickers
- `tblDecisions` - Full trade history
- `tblPositions` - Open positions

### 7 Named Ranges (Settings)

| Name | Default | Description |
|---|---|---|
| Equity_E | 10,000 | Account equity for sizing |
| RiskPct_r | 0.0075 | Risk per unit (0.75%) |
| StopMultiple_K | 2 | ATR multiple for stop |
| HeatCap_H_pct | 0.04 | Portfolio heat cap (4%) |
| BucketHeatCap_pct | 0.015 | Bucket heat cap (1.5%) |
| AddStepN | 0.5 | Add-on step (0.5N) |
| EarningsBufferDays | 3 | Days around earnings |

## Trading Workflow

### Daily Routine

**Morning (5 minutes)**:
1. Click "Open FINVIZ" for each preset
2. Click "Import Candidates" → paste tickers
3. Review in TradingView

**During Market (per trade, 3 minutes)**:
1. Select Ticker from dropdown
2. Enter Entry Price, ATR N, K
3. Check 6 checklist items
4. Click "Evaluate" → Wait for GREEN
5. Click "Recalc Sizing"
6. Wait 2 minutes (impulse timer)
7. Click "Save Decision"

### The 6-Item Checklist

Every trade must pass all 6 checks for GREEN:

1. **FromPreset** - Ticker came from today's FINVIZ import
2. **TrendPass** - Meets trend criteria (20/50/200 SMA alignment)
3. **LiquidityPass** - Sufficient volume (> 500K shares/day)
4. **TVConfirm** - TradingView strategy signal fired
5. **EarningsOK** - No earnings within 3 days
6. **JournalOK** - Reviewed in trading journal, no conflicts

**Banner Logic**:
- All 6 pass → **GREEN** (can save)
- 1 missing → **YELLOW** (caution)
- 2+ missing → **RED** (no-go)

### The 5 Hard Gates (SaveDecision)

Even if GREEN, the system blocks if:
1. ❌ Banner not GREEN
2. ❌ Ticker not in today's Candidates
3. ❌ 2-minute impulse timer not elapsed
4. ❌ Bucket in cooldown
5. ❌ Heat caps exceeded (portfolio or bucket)

**All 5 must pass** to log the trade.

## Position Sizing

### Method 1: Stock

```
R = Equity × RiskPct
StopDist = K × N
Shares = floor(R / StopDist)
```

**Example**: E=$10,000, r=0.75%, N=1.50, K=2
- R = $75
- StopDist = 3.00
- Shares = 25

### Method 2: Options - Delta-ATR

```
Contracts = floor(R / (K × N × Delta × 100))
```

**Example**: R=$75, N=1.50, K=2, Delta=0.30
- Contracts = floor(75 / (2 × 1.50 × 0.30 × 100))
- Contracts = floor(75 / 90) = 0
- (Need higher R or lower Delta)

### Method 3: Options - MaxLoss

```
Contracts = floor(R / MaxLossPerContract)
```

**Example**: R=$75, MaxLoss=$50 (debit spread)
- Contracts = floor(75 / 50) = 1

## Heat Management

### Portfolio Heat Cap

**Definition**: Sum of all open R across all positions
**Cap**: 4% of equity (default = $400 on $10K account)

**Example**:
- Position 1: $75 R
- Position 2: $50 R
- Position 3: $100 R
- **Total**: $225 / $400 = 56% (OK)

If new trade is $100 R → $325 / $400 = 81% (still OK)
If new trade is $200 R → $425 / $400 = **BLOCKED**

### Bucket Heat Cap

**Definition**: Sum of open R within one correlation bucket
**Cap**: 1.5% of equity (default = $150 on $10K account)

**Example** (Tech/Comm bucket):
- AAPL: $75 R
- MSFT: $50 R
- **Total**: $125 / $150 = 83% (OK)

New NVDA trade: $75 R → $200 / $150 = **BLOCKED**

### Bucket Cooldown

**Trigger**: 2 stop-outs in 20 days → 10-day cooldown
**Effect**: Cannot enter new trades in that bucket
**Purpose**: Prevents over-trading weak sectors

## Python Integration (Optional)

### Requirements

- Windows with Excel 365
- Microsoft 365 Insider (for Python in Excel)
- Internet connection

### Setup (if using Python)

1. **Enable Python in Excel**:
   - Data tab → Python in Excel
   - Follow Microsoft's setup wizard

2. **Test Python**:
   ```cmd
   cd Python
   python finviz_scraper.py
   python heat_calculator.py
   ```

3. **Test in Excel**:
   - Alt+F11 → Immediate Window
   - Type: `TF_Python_Bridge.TestPythonIntegration`
   - Press Enter

### Python vs VBA

| Feature | VBA Only | With Python |
|---|---|---|
| Import time | 30-60 sec (manual) | 5-10 sec (auto) |
| Heat calc speed | 1-3 sec | <0.5 sec |
| FINVIZ scraping | Manual copy/paste | Automated |
| Multi-page results | Manual | Automatic |

**Recommendation**: Start with VBA-only, add Python later as enhancement.

## Troubleshooting

### "Build failed - VBA project access denied"

**Fix**:
1. File → Options → Trust Center → Trust Center Settings
2. Macro Settings
3. Check "Trust access to the VBA project object model"
4. Click OK
5. Re-run BUILD.bat

### "Buttons don't work / macro not found"

**Fix**:
1. Enable macros: Click "Enable Content" banner
2. Verify modules imported: Alt+F11 → Check modules list
3. Re-assign button: Right-click button → Assign Macro

### "Dropdown lists are empty"

**Fix**:
1. Go to Presets/Buckets/Candidates sheets
2. Verify tables have data
3. If empty, run: `TF_Data.EnsureStructure` (Alt+F11 → Immediate Window)

### "Cannot save decision - Ticker not in candidates"

**Fix**:
1. Import candidates first (Import Candidates button)
2. Verify ticker exists in Candidates sheet with today's date
3. Check spelling

### "Python integration doesn't work"

**Expected**: Python in Excel is only available in Microsoft 365 Insider
**Workaround**: Use manual import (works perfectly)

## Customization

### Add New FINVIZ Preset

1. Go to Presets sheet
2. Add row:
   - Name: "MY_CUSTOM_PRESET"
   - QueryString: Copy from FINVIZ URL (after `screener.ashx?`)
3. Preset now appears in TradeEntry dropdown

### Change Risk Settings

1. Go to Summary sheet
2. Modify values in column B:
   - Equity_E: Your account size
   - RiskPct_r: Risk per trade (0.0075 = 0.75%)
   - HeatCap_H_pct: Max total heat (0.04 = 4%)
3. Values update immediately

### Add New Bucket

1. Go to Buckets sheet
2. Add row with bucket settings
3. Bucket appears in TradeEntry dropdown

## Development

### File Modifications

After modifying VBA files:
```cmd
BUILD.bat
```

After modifying Python files:
```bash
cd Python
python finviz_scraper.py  # Test
```

### Export VBA Modules

To export current VBA code:
1. Alt+F11 → VBA Editor
2. Right-click module → Export File
3. Save to `/VBA/` folder

## Credits

Built with:
- Excel VBA for UI and automation
- Python for web scraping and calculations
- BeautifulSoup for HTML parsing
- pandas for data processing
- pywin32 for COM automation

## License

MIT License - Free to use and modify

## 📚 Documentation

### For Users (Start Here!)
1. **USER_GUIDE.md** ⭐ - Complete beginner walkthrough (15,000+ words)
   - First-time setup step-by-step
   - What every field means (ATR, K, Delta, etc.)
   - Real trading examples with AAPL
   - The 6-item checklist explained
   - Position sizing for stocks and options
   - Troubleshooting common issues
   - **Auto-opens on first launch!**

2. **UPDATED_QUICK_START.md** - 2-page quick reference
3. **README_UPDATED.md** - Feature summary

### For Developers
1. **CHANGELOG.md** - Version history and upgrade notes
2. **DEVELOPMENT_LOG.md** - Technical issue tracker and AI assistant context
3. **IMPLEMENTATION_STATUS.md** - Architecture and technical details
4. **START_HERE.md** - Original detailed setup guide

### Quick Access
- **Setup Sheet** (in workbook) - Status, utility buttons, instructions
- **"Open User Guide" Button** - On Setup sheet, reopens USER_GUIDE.md anytime

## 🔄 Changelog

### v2.0.0 (2025-01-27)
- ✅ Added automatic checkbox creation
- ✅ Fixed Unicode encoding issues
- ✅ Fixed Python detection for Python in Excel 2023+
- ✅ Fixed duplicate button creation
- ✅ Added USER_GUIDE.md (comprehensive beginner guide)
- ✅ Added auto-open guide on first launch
- ✅ Enhanced error handling throughout

See **CHANGELOG.md** for complete version history.

## Support

For issues or questions:
1. **USER_GUIDE.md** - Comprehensive troubleshooting section
2. **DEVELOPMENT_LOG.md** - Known issues and solutions
3. **Setup Sheet** - Built-in help and utility buttons
4. **CHANGELOG.md** - Check if issue was fixed in latest version

---

**Ready to build?** Run `BUILD.bat` and start trading in 3 minutes! 🚀
