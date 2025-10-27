# Quick Setup Guide - Get Trading System Working in 30 Minutes

## Current Situation

You have:
- ✅ Python modules created (finviz_scraper.py, heat_calculator.py)
- ✅ Basic VBA utility module (TF_Utils.bas)
- ✅ Complete specifications
- ❌ Missing: Core VBA logic + UI elements

## What's Missing from Your Workbook

Based on your description, your TradeEntry sheet likely looks like this:
```
Empty sheet OR minimal table headers, no buttons, no checkboxes
```

You need:
1. **6 checkboxes** for the checklist
2. **4 dropdown lists** (Preset, Ticker, Sector, Bucket)
3. **5 buttons** (Evaluate, Recalc, Save, Import, Open FINVIZ)
4. **VBA code** to make them work

## The Missing Pieces Explained

### 1. The 6 Checkboxes

These are **Form Control checkboxes** that users tick before clicking Evaluate:

| Checklist Item | Purpose |
|---|---|
| ☐ FromPreset | Ticker came from today's FINVIZ import |
| ☐ TrendPass | Meets trend criteria (verified in TradingView) |
| ☐ LiquidityPass | Sufficient volume/liquidity |
| ☐ TVConfirm | TradingView strategy gives entry signal |
| ☐ EarningsOK | No earnings within buffer days |
| ☐ JournalOK | Reviewed in trading journal, no conflicts |

**How to add them**:
1. Developer tab → Insert → Form Controls → Check Box
2. Place 6 checkboxes in cells B20:B25
3. Link each to cells C20:C25 (stores TRUE/FALSE)
4. Label them with the names above

### 2. The 4 Dropdowns

These are **Data Validation dropdowns**:

| Cell | Dropdown | Source |
|---|---|---|
| B5 | Preset | List from Presets table (tblPresets[Name]) |
| B6 | Ticker | List from today's Candidates (tblCandidates[Ticker]) |
| B7 | Sector | Manual list: Technology,Healthcare,Financials,Consumer,Industrials,Energy |
| B8 | Bucket | List from Buckets table (tblBuckets[Bucket]) |

**How to add them**:
1. Select cell B5
2. Data tab → Data Validation → List
3. Source: `=tblPresets[Name]`
4. Repeat for B6, B7, B8 with their respective sources

### 3. The 5 Buttons

These are **Command Buttons** that trigger VBA macros:

| Button | Calls Macro | What It Does |
|---|---|---|
| Evaluate | TF_UI.EvaluateChecklist | Checks all 6 boxes, shows GREEN/YELLOW/RED banner |
| Recalc Sizing | TF_UI.RecalcSizing | Calculates shares/contracts based on Entry/ATR/K |
| Save Decision | TF_UI.SaveDecision | Validates and logs trade to Decisions table |
| Import Candidates | TF_Presets.ImportCandidatesPrompt | Imports tickers (manual paste or Python scrape) |
| Open FINVIZ | TF_Presets.OpenPreset | Opens selected preset URL in browser |

**How to add them**:
1. Developer tab → Insert → Button (Form Control)
2. Draw button on sheet
3. Assign macro when prompted
4. Repeat 5 times

## Step-by-Step Setup (30 Minutes)

### Step 1: Test Python Modules (5 min)

```bash
cd /home/kali/excel-trading-workflow/Python

# This should show ~47 tickers from FINVIZ New Highs screener
python finviz_scraper.py

# This should show portfolio heat calculations
python heat_calculator.py
```

If both work, Python is ready. ✅

### Step 2: Create Remaining VBA Modules (I'll do this next)

I need to create:
- TF_Data.bas (structure setup)
- TF_UI.bas (core trading logic)
- TF_Presets.bas (FINVIZ integration)

These are ~800 lines total. Should I create them now?

### Step 3: Manual UI Setup (10 min - YOU do this)

Once VBA is ready:

1. **Open your Excel file** (TrendFollowing_TradeEntry.xlsm)
2. **Go to TradeEntry sheet**
3. **Add labels in column A**:
   - A5: "Preset:"
   - A6: "Ticker:"
   - A7: "Sector:"
   - A8: "Bucket:"
   - A9: "Entry Price:"
   - A10: "ATR N:"
   - A11: "K:"
   - A20-A25: Checklist labels

4. **Add 6 checkboxes**:
   - Developer → Insert → Check Box (Form Control)
   - Place next to A20:A25
   - Right-click each → Format Control → Cell link: C20:C25

5. **Add 4 dropdowns**:
   - Select B5 → Data → Data Validation → List → `=tblPresets[Name]`
   - Select B6 → Data → Data Validation → List → `=tblCandidates[Ticker]`
   - Select B7 → Data → Data Validation → List → `Technology,Healthcare,Financials,Consumer,Industrials,Energy`
   - Select B8 → Data → Data Validation → List → `=tblBuckets[Bucket]`

6. **Add 5 buttons**:
   - Developer → Insert → Button
   - Assign macros:
     - Button 1 → TF_UI.EvaluateChecklist
     - Button 2 → TF_UI.RecalcSizing
     - Button 3 → TF_UI.SaveDecision
     - Button 4 → TF_Presets.ImportCandidatesPrompt
     - Button 5 → TF_Presets.OpenPreset

### Step 4: Test Workflow (5 min)

1. **Add test candidate**:
   - Go to Candidates sheet
   - Add row: Today's date, AAPL, TestPreset, Technology, Tech/Comm

2. **Test Evaluate button**:
   - Go to TradeEntry
   - Select Ticker: AAPL
   - Check all 6 boxes
   - Click Evaluate
   - Should see GREEN banner in A2

3. **Test sizing**:
   - Enter: Entry=180, N=1.50, K=2
   - Click Recalc Sizing
   - Should see calculated shares/stop

## About FINVIZ Integration

### How It Works (Technical):

FINVIZ URLs look like:
```
https://finviz.com/screener.ashx?v=211&s=ta_newhigh&ft=4
```

Breaking down the parameters:
- `v=211` - View type (overview table)
- `s=ta_newhigh` - Signal: Technical Analysis - New High
- `ft=4` - Filter type

The Python scraper:
1. Sends HTTP request to this URL
2. Parses the HTML table with BeautifulSoup
3. Extracts ticker symbols from table cells
4. Handles pagination (20 tickers per page)
5. Returns clean list: `['AAPL', 'MSFT', 'NVDA', ...]`

### NOT Just a Permalink:

❌ **NOT**: Just opening a saved URL
✅ **IS**: Active web scraping of live FINVIZ results

**Benefits**:
- Always gets latest screener results
- No manual copy/paste
- Handles multi-page results automatically
- Normalizes ticker symbols

**Limitations**:
- Requires internet connection
- May break if FINVIZ changes HTML structure
- Rate limited (1-2 sec delay between requests)

### Manual Fallback:

If Python isn't available or fails:
1. Click "Open FINVIZ" button → opens URL in browser
2. Manually copy tickers from webpage
3. Click "Import Candidates" → paste tickers → auto-normalized

## Current Implementation Status

| Component | Status | Location |
|---|---|---|
| Python scraper | ✅ Done | `/Python/finviz_scraper.py` |
| Python heat calc | ✅ Done | `/Python/heat_calculator.py` |
| VBA utilities | ✅ Done | `/VBA/TF_Utils.bas` |
| VBA data layer | ❌ Need to create | `/VBA/TF_Data.bas` |
| VBA UI logic | ❌ Need to create | `/VBA/TF_UI.bas` |
| VBA FINVIZ integration | ❌ Need to create | `/VBA/TF_Presets.bas` |
| VBA-Python bridge | ❌ Optional | `/VBA/TF_Python_Bridge.bas` |
| Excel UI elements | ❌ Your task | Manual setup (10 min) |

## Next Decision Point

**What do you want me to do next?**

### Option A: "Create ALL remaining VBA code" (50 min for me)
- I'll write TF_Data.bas (~320 lines)
- I'll write TF_UI.bas (~384 lines)
- I'll write TF_Presets.bas (~150 lines)
- I'll write TF_Python_Bridge.bas (~280 lines)
- Result: Complete VBA codebase, you just import and add UI elements

### Option B: "Create MINIMAL working VBA" (15 min for me)
- I'll write simplified versions of the 3 core modules
- You get working Evaluate/Recalc/Save buttons
- Simplified import (manual paste only initially)
- Can enhance later

### Option C: "Help me set up UI first, code later"
- I'll guide you through adding checkboxes/buttons/dropdowns
- We'll add VBA code as we test each piece
- Incremental approach

**Which path should I take?**

Also:
- Do you have an existing .xlsm file with tables already set up?
- Are you on Windows or Linux (WSL)?
- Do you want Python integration now or can it wait?

Let me know and I'll continue accordingly!
