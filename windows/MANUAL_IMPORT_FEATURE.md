# Manual Data Import Feature - Guide

**Version:** M24 - Manual Import Enhancement
**Created:** 2025-10-28
**Purpose:** Manually import database data into Excel Dashboard

---

## Overview

The manual import feature allows you to pull data from the `tf-engine` database directly into your Excel workbook's Dashboard sheet. This is useful for:

- ✅ **Troubleshooting** - Verify backend is working and returning data
- ✅ **Manual Refresh** - Update Dashboard with latest database contents
- ✅ **Data Verification** - Check positions, candidates, settings, cooldowns
- ✅ **Initial Setup** - Populate Dashboard after first import of candidates

---

## Available Import Functions

### 1. RefreshDashboardData()
**What it does:** Imports ALL data in one click

**Data imported:**
- Settings (equity, risk%, heat caps, K multiple)
- Open positions (ticker, bucket, units, risk $, entry, stop)
- Today's candidates (tickers from FINVIZ import)
- Active cooldowns (buckets with restrictions)

**How to run:**
1. Open Excel
2. Press Alt+F8
3. Select "RefreshDashboardData"
4. Click Run

**Expected result:**
- Dashboard sheet populates with current database data
- Status message: "Dashboard data refreshed successfully!"
- All sections show current data

---

### 2. ImportSettings()
**What it does:** Imports settings only

**Data imported:**
- Account Equity (row 3)
- Risk % per Trade (row 4)
- Portfolio Heat Cap (row 5)
- Bucket Heat Cap (row 6)
- Stop Multiple K (row 7)

**Dashboard layout (rows 3-7):**
```
A3: Account Equity:        B3: $100,000
A4: Risk % per Trade:      B4: 0.75%
A5: Portfolio Heat Cap:    B5: 4.00%
A6: Bucket Heat Cap:       B6: 1.50%
A7: Stop Multiple (K):     B7: 2
```

---

### 3. ImportPositions()
**What it does:** Imports open positions

**Data imported:**
- Ticker, Bucket, Units, Risk $, Entry, Stop, Open Date, Status
- Total open positions count
- Total portfolio risk

**Dashboard layout (rows 9-31):**
```
Row 9:  Open Positions:
Row 10: Headers (Ticker, Bucket, Units, Risk $, Entry, Stop, Open Date, Status)
Row 11-30: Position data (up to 20 positions)
Row 31: Total Open Positions: X  |  Total Risk: $Y
```

**Example:**
```
Ticker  Bucket      Units  Risk $  Entry   Stop    Open Date   Status
AAPL    Tech/Comm   250    750.00  180.00  177.00  2025-10-28  open
NVDA    Tech/Comm   150    900.00  450.00  444.00  2025-10-27  open
```

---

### 4. ImportTodaysCandidates()
**What it does:** Imports today's candidate tickers

**Data imported:**
- Ticker, Preset, Sector, Bucket
- Total candidates count
- Date

**Dashboard layout (columns K-N, rows 9-31):**
```
Row 9:  Today's Candidates:
Row 10: Headers (Ticker, Preset, Sector, Bucket)
Row 11-30: Candidate data (up to 20 candidates)
Row 31: Total Candidates: X  |  Date: YYYY-MM-DD
```

**Example:**
```
Ticker  Preset              Sector     Bucket
CCJ     TF-Breakout-Long    Energy     Energy
W       TF-Breakout-Long    Retail     Consumer
SOFI    TF-Breakout-Long    FinTech    Finance
```

---

### 5. ImportCooldowns()
**What it does:** Imports active cooldowns

**Data imported:**
- Bucket, Triggered At, Clears At, Reason
- Active cooldowns count

**Dashboard layout (columns K-N, rows 34-51):**
```
Row 34: Active Cooldowns:
Row 35: Headers (Bucket, Triggered At, Clears At, Reason)
Row 36-50: Cooldown data
Row 51: Active Cooldowns: X
```

**Example:**
```
Bucket      Triggered At          Clears At            Reason
Tech/Comm   2025-10-28 10:30:00   2025-10-28 10:32:00  2 losses in a row
```

---

## How to Use

### Method 1: Run from Alt+F8 (Macro List)

1. **Open** TradingPlatform.xlsm
2. **Press** Alt+F8 (opens macro list)
3. **Select** the function you want:
   - `RefreshDashboardData` - Import all data
   - `ImportSettings` - Settings only
   - `ImportPositions` - Positions only
   - `ImportTodaysCandidates` - Candidates only
   - `ImportCooldowns` - Cooldowns only
4. **Click** Run
5. **Check** Dashboard sheet for updated data

### Method 2: Add Buttons to Dashboard Sheet

You can add clickable buttons to the Dashboard sheet:

**For RefreshDashboardData button:**
1. Open TradingPlatform.xlsm
2. Go to Dashboard sheet
3. Developer tab → Insert → Button (Form Control)
4. Draw button on sheet (e.g., cell P3)
5. Assign macro: `RefreshDashboardData`
6. Right-click button → Edit Text → "Refresh All Data"
7. Click button to refresh all data

**Repeat for individual import buttons:**
- Button for `ImportSettings` → "Refresh Settings"
- Button for `ImportPositions` → "Refresh Positions"
- Button for `ImportTodaysCandidates` → "Refresh Candidates"
- Button for `ImportCooldowns` → "Refresh Cooldowns"

**Suggested Dashboard layout:**
```
[Refresh All Data]  (calls RefreshDashboardData)

Settings              [Refresh Settings]
... settings data ...

Open Positions        [Refresh Positions]
... positions data ...

Today's Candidates    [Refresh Candidates]
... candidates data ...

Active Cooldowns      [Refresh Cooldowns]
... cooldowns data ...
```

### Method 3: Run from VBA Editor

1. **Press** Alt+F11 (opens VBA Editor)
2. **Find** TFEngine module in left panel
3. **Put cursor** in the function you want (e.g., RefreshDashboardData)
4. **Press** F5 to run
5. **Alt+Tab** back to Excel to see results

---

## Troubleshooting

### Error: "Subscript out of range"

**Cause:** Dashboard sheet doesn't exist or is named differently

**Fix:**
- Make sure you have a sheet named "Dashboard" (exact spelling, case-sensitive)
- Or rename your sheet to "Dashboard"

### Error: "Failed to import [data type]"

**Cause:** Backend command failed or returned error

**Check:**
1. Is `tf-engine.exe` in the same folder as the workbook?
2. Is `trading.db` initialized? Run: `tf-engine.exe init`
3. Check TradingSystem_Debug.log for correlation ID and error details

**Test backend directly:**
```cmd
# Test settings
tf-engine.exe get-settings --format json

# Test positions
tf-engine.exe list-positions --open-only --format json

# Test candidates
tf-engine.exe list-candidates --format json

# Test cooldowns
tf-engine.exe list-cooldowns --active-only --format json
```

### No Data Showing Up

**Cause:** Database is empty or data wasn't imported yet

**Fix:**
1. **Import candidates** first:
   ```cmd
   import-candidates.bat
   ```
2. **Initialize settings** if needed:
   ```cmd
   tf-engine.exe init
   ```
3. **Add test positions** (if needed for testing):
   ```cmd
   tf-engine.exe open-position --ticker AAPL --bucket Tech/Comm --units 250 --risk 750 --entry 180 --stop 177
   ```

### Data is Stale

**Solution:** Click "Refresh All Data" button or run `RefreshDashboardData` macro

**Note:** The import is NOT automatic - you must manually trigger it by:
- Clicking a button
- Running the macro (Alt+F8)
- Or add VBA code to auto-refresh on workbook open (advanced)

---

## Dashboard Layout Reference

**Recommended cell layout for Dashboard sheet:**

```
A1: "Dashboard"

Settings Section (A3:B7):
A3: "Account Equity:"          B3: $100,000
A4: "Risk % per Trade:"        B4: 0.75%
A5: "Portfolio Heat Cap:"      B5: 4.00%
A6: "Bucket Heat Cap:"         B6: 1.50%
A7: "Stop Multiple (K):"       B7: 2

Open Positions Section (A9:H31):
A9: "Open Positions:"
A10-H10: Column headers
A11-H30: Position data (20 rows)
A31: "Total Open Positions: X  |  Total Risk: $Y"

Today's Candidates Section (K9:N31):
K9: "Today's Candidates:"
K10-N10: Column headers
K11-N30: Candidate data (20 rows)
K31: "Total Candidates: X  |  Date: YYYY-MM-DD"

Active Cooldowns Section (K34:N51):
K34: "Active Cooldowns:"
K35-N35: Column headers
K36-N50: Cooldown data (15 rows)
K51: "Active Cooldowns: X"
```

---

## Auto-Refresh on Workbook Open (Advanced)

To automatically refresh data when you open the workbook:

1. **Press** Alt+F11 (VBA Editor)
2. **Double-click** "ThisWorkbook" in left panel
3. **Add this code:**

```vba
Private Sub Workbook_Open()
    ' Auto-refresh dashboard data on workbook open
    On Error Resume Next
    Call TFEngine.RefreshDashboardData
End Sub
```

4. **Save** the workbook
5. **Close** and re-open workbook
6. **Dashboard** will auto-populate!

**Warning:** This runs every time you open the workbook, which calls `tf-engine.exe` multiple times. Only enable if you want automatic refresh.

---

## Technical Details

**VBA Modules Modified:**
- `TFEngine.bas` - Added 5 new import functions (lines 1254-1557)
- `TFHelpers.bas` - Added 3 new Parse list functions (lines 357-464)

**New Functions Added:**

| Function | Purpose | Returns |
|----------|---------|---------|
| `RefreshDashboardData()` | Import all data | Void (populates Dashboard) |
| `ImportSettings()` | Import settings | Void |
| `ImportPositions()` | Import open positions | Void |
| `ImportTodaysCandidates()` | Import candidates | Void |
| `ImportCooldowns()` | Import cooldowns | Void |

**New Parse Functions:**

| Function | Purpose | Output Type |
|----------|---------|-------------|
| `ParsePositionsListJSON()` | Parse positions array | TFPositionsList |
| `ParseCandidatesListJSON()` | Parse candidates array | TFCandidatesList |
| `ParseCooldownsListJSON()` | Parse cooldowns array | TFCooldownsList |

**Backend Commands Called:**
- `tf-engine.exe get-settings --format json`
- `tf-engine.exe list-positions --open-only --format json`
- `tf-engine.exe list-candidates --format json`
- `tf-engine.exe list-cooldowns --active-only --format json`

---

## Logs and Debugging

Every import operation logs to `TradingSystem_Debug.log`:

**Example log entries:**
```
[2025-10-28 14:30:52] [INFO] [20251028-143052133-7A3F] Imported 3 positions
[2025-10-28 14:30:52] [INFO] [20251028-143052133-7A3F] Imported 15 candidates
[2025-10-28 14:30:52] [INFO] [20251028-143052133-7A3F] Imported 1 active cooldowns
[2025-10-28 14:30:52] [INFO] [20251028-143052133-7A3F] Settings imported successfully
```

**If import fails:**
```
[2025-10-28 14:30:52] [ERROR] [20251028-143052133-7A3F] ImportPositions failed: exit code 1: database not found
```

Check both logs:
- `TradingSystem_Debug.log` - VBA frontend log
- `tf-engine.log` - Go backend log (if exists)

Both use the same correlation ID for cross-referencing.

---

## Summary

✅ **5 new import functions** - Manually refresh Dashboard data
✅ **Clickable buttons** - Add to Dashboard for easy access
✅ **Troubleshooting** - Direct backend testing
✅ **Auto-refresh option** - For advanced users
✅ **Full logging** - Every import logged with correlation ID

**After running `fix-vba-modules.bat`, all import functions will be available!**

---

## Next Steps

1. **Run** `fix-vba-modules.bat` to import updated VBA modules
2. **Open** TradingPlatform.xlsm
3. **Test** imports using Alt+F8 → RefreshDashboardData
4. **Add buttons** to Dashboard sheet (optional, for convenience)
5. **Use** manual import whenever you need fresh data from database

---

**Last Updated:** 2025-10-28
**Feature Version:** M24 (Manual Import Enhancement)
**Status:** ✅ Ready to use after running fix-vba-modules.bat
