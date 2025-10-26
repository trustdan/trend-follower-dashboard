# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

---

## Overview

This is an **Excel VBA-based trend-following trading system** implementing Seykota/Turtle methodology with options execution. The workflow integrates FINVIZ screeners → TradingView strategy validation → Excel workbook for position sizing, decision tracking, and impulse control.

**Key Philosophy**: Minimize bias, enforce mechanical rules, track portfolio/bucket heat, and prevent impulsive entries through structured checklists and cooldown timers.

---

## System Architecture

### Three-Layer Workflow

1. **FINVIZ Screeners** (external) → Generate candidate tickers via presets
2. **TradingView Strategy** (Pine Script v6) → Validate signals using Donchian breakouts, ATR-based stops, and pyramiding rules
3. **Excel Workbook** (VBA) → Execute decision checklist, size positions, enforce heat caps, log trades

### Core Trading Rules (documented in workflow-plan.md)

- **Entry**: 52-week highs (or 55-bar Donchian breakout), price above 50/200 SMA
- **Exit**: 10-bar opposite Donchian OR 2×N ATR stop (whichever is closer)
- **Position Sizing**: Risk 0.25-0.75% per unit using 2×N stop distance
- **Pyramiding**: Add every 0.5×N up to max 4 units (only to winners)
- **Portfolio Heat Cap**: 2-4% total; per-bucket cap at 1.0-1.5%
- **Correlation Buckets**: Tech/Comm, Consumer, Financials, Industrials, Energy/Materials, Defensives/REITs

### Options Execution (when using options instead of stock)

- **DTE**: 60-90 days; roll/close around 21 DTE
- **Structure**: ITM calls (Δ≈0.6-0.7) when IVR<20; debit verticals when IVR≥20
- **Sizing**: Delta-ATR method or MaxLoss method (per-contract risk calculation)

---

## Excel Workbook Structure

### Primary Workbook: Interactive Trade Entry System

**File Reference**: `newest-Interactive_TF_Workbook_Plan.md` (18KB detailed spec)

#### Key Sheets

1. **Trade Entry** (main UI sheet)
   - Single-screen interface for entire workflow
   - Dropdowns: Preset, Ticker, Sector, Bucket, Method (Stock/Opt-DeltaATR/Opt-MaxLoss)
   - Inputs: Entry price, ATR N, K (stop multiple), Delta/DTE/MaxLoss
   - Checklist: 6 required items (FromPreset, TrendPass, LiquidityPass, TVConfirm, EarningsOK, JournalOK)
   - GO/NO-GO banner: GREEN/YELLOW/RED based on checklist + heat caps + cooldown + impulse timer
   - Buttons: Evaluate, Recalc Sizing, Save Decision, Import Candidates

2. **Presets** (`tblPresets`)
   - 5 FINVIZ screener query strings (TF_BREAKOUT_LONG, TF_MOMENTUM_UPTREND, etc.)

3. **Buckets** (`tblBuckets`)
   - Sector → Bucket mapping with cooldown rules
   - Tracks StopoutsToCooldown (default 2), CooldownBars (default 10)

4. **Candidates** (`tblCandidates`)
   - Today's tickers imported from FINVIZ

5. **Decisions** (`tblDecisions`)
   - Append-only log of all trade decisions (timestamp, ticker, sizing, heat metrics, banner status)

6. **Positions** (`tblPositions`)
   - Open positions tracker (UnitsOpen, RperUnit, TotalOpenR, NextAddPrice)

7. **Summary**
   - Named ranges: `Equity_E`, `RiskPct_r`, `StopMultiple_K`, `HeatCap_H_pct`, `BucketHeatCap_pct`, `AddStepN`, `EarningsBufferDays`

#### GO/NO-GO Logic (critical)

**GREEN** (OK to trade) requires ALL of:
- All 6 required checklist items TRUE
- Ticker in today's Candidates
- Portfolio heat post-entry ≤ HeatCap_H_pct × Equity_E
- Bucket heat post-entry ≤ BucketHeatCapPct × Equity_E
- Bucket NOT in cooldown
- 2-minute impulse timer elapsed since Evaluate button clicked

**Save Decision** button hard-gates on GREEN status.

### Secondary Workbook: Options Calendar Dashboard

**File Reference**: `older-Options_Trend_Dashboard_Summary.md` (61KB VBA code + spec)

#### Key Sheets

1. **Calendar**
   - 10-week rolling view (2 weeks back, 8 forward)
   - Sector rows × week columns showing active symbols
   - Refreshes via `CalendarModule.RefreshCalendar` (values-only to avoid #NAME? errors)

2. **Checklist**
   - Interactive pre-trade checklist with 3-state banner (DO NOT TRADE / CAUTION / OK TO TRADE)
   - Contracts helper: computes single-call contracts and debit vertical spreads from per-unit risk
   - Preset toggles: System-1 (20/10) vs System-2 (55/10)
   - Buttons: Reset, Add To Trades, Recalculate, Apply Preset, Save Decision

3. **Trades**
   - Simple log (Symbol, Sector, StartDate, EndDate, Active, Strategy, Notes)

#### VBA Modules

- **CalendarModule**: `ForceRepairAndRefresh`, `RefreshCalendar`, `StartOfWeek`
- **ChecklistModule**: All UI logic, banner coloring, sizing helpers, decision logging

---

## Common Development Commands

### No build/test commands (Excel-only)

This is a VBA-based system; no external build tools. To work on the code:

1. Open `.xlsm` file in Excel
2. `Alt+F11` to open VBA Editor
3. Modify modules in place
4. `F5` to run/test procedures
5. Save workbook (macros must be enabled)

### Key VBA Entry Points

**Setup/Structure** (run once on new workbook):
```vba
EnsureWorkbookStructure  ' Creates all sheets, tables, named ranges, seed data
```

**Calendar Operations**:
```vba
ForceRepairAndRefresh    ' Clears old formulas and rebuilds calendar as values
RefreshCalendar          ' Updates calendar grid from Trades table
```

**Checklist Operations**:
```vba
BuildInteractiveChecklist  ' Creates Checklist sheet from scratch
ApplyChecklistFixes        ' Repairs banner formula, removes CF, tidies buttons
EvaluateChecklist          ' Computes GO/NO-GO banner
RecalcSizing               ' Computes R, shares/contracts, stops, add levels
SaveDecision               ' Hard-gates and logs to Decisions + Positions
```

**Data Import**:
```vba
OpenPreset                 ' Opens FINVIZ URL for selected preset in browser
ImportCandidatesPrompt     ' Pastes tickers from clipboard → tblCandidates
```

**Maintenance**:
```vba
UpdateCooldowns            ' Checks recent StopOuts and flags buckets for cooldown
RefreshReview              ' Computes adherence KPIs (% non-GREEN trades)
```

---

## Position Sizing Calculations (critical logic)

### Common Variables
- **E** = `Equity_E` (account size)
- **r** = `RiskPct_r` (risk % per unit, typically 0.5-0.75%)
- **R** = E × r (risk dollars per unit)
- **N** = ATR(20) from TradingView (volatility measure)
- **K** = `StopMultiple_K` (default 2.0)
- **StopDist** = K × N
- **InitialStop** = Entry - StopDist (for longs)

### Stock Sizing
```
Shares = floor(R / StopDist)
```

### Options Sizing - Delta-ATR Method
```
Contracts = floor(R / (K × N × Delta × 100))
```
Use for single ITM calls; Delta typically 0.60-0.70.

### Options Sizing - MaxLoss Method
```
Contracts = floor(R / (MaxLossPerContract × 100))
```
Use for debit verticals; MaxLoss = net debit.

### Add Levels (Pyramiding)
```
Add1 = Entry + (AddStepN × N)
Add2 = Entry + (2 × AddStepN × N)
Add3 = Entry + (3 × AddStepN × N)
```
Default `AddStepN = 0.5`.

---

## TradingView Strategy

**File Reference**: `SeykotaTurtleTrend-FollowingOptionsExecution+TradingViewStrategyGuide.md`

**Strategy Name**: "Seykota / Turtle Core v2.1 + Date Range"

### Key Parameters
- `entryLen=55` (Donchian entry lookback, System-2)
- `exitLen=10` (Donchian exit lookback)
- `nLen=20` (ATR period)
- `stopN=2.0` (stop distance in N)
- `addStepN=0.5` (add-on step in N)
- `maxUnits=4` (max pyramid units)
- `riskPct=0.5-1.0` (% risk per unit)

### Alert Conditions (available in strategy)
- Long Entry, Short Entry
- Long Add, Short Add
- Long Exit, Short Exit
- Time Exit (optional)

Use **"Once per bar close"** for cleaner alerts.

### Date Range Filter
```pinescript
fromDate = input.time(timestamp("2022-01-01T00:00:00"), "Backtest FROM")
toDate   = input.time(timestamp("2099-12-31T23:59:59"), "Backtest TO")
inRange  = time >= fromDate and time <= toDate
flatAtFrom = input.bool(true, "Force FLAT at range start?")
```

---

## FINVIZ Presets (URL query strings)

Base URL: `https://finviz.com/screener.ashx?`

1. **TF_BREAKOUT_LONG**
   ```
   v=211&p=d&s=ta_newhigh&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pa,ta_sma200_pa&o=-relativevolume
   ```

2. **TF_MOMENTUM_UPTREND**
   ```
   v=211&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pa,ta_sma200_pa&dr=y1&o=-marketcap
   ```

3. **TF_UNUSUAL_VOLUME**
   ```
   v=211&p=d&s=ta_unusualvolume&f=cap_largeover,sh_price_o20,ta_sma50_pa,ta_sma200_pa&o=-relativevolume
   ```

4. **TF_BREAKDOWN_SHORT**
   ```
   v=211&p=d&s=ta_newlow&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pb,ta_sma200_pb&o=-relativevolume
   ```

5. **TF_MOMENTUM_DOWNTREND**
   ```
   v=211&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pb,ta_sma200_pb&dr=y1&o=-marketcap
   ```

---

## Daily Workflow (15-25 minutes)

1. **After market close** (primary run):
   - Run FINVIZ presets 1-3 in order
   - Export tickers, dedupe, tag by sector → import to Excel Candidates
   - For each candidate: open TradingView, validate signal with strategy
   - Use Trade Entry sheet: fill inputs, check boxes, click Evaluate
   - If GREEN + 2-minute timer elapsed: click Save Decision

2. **Weekly maintenance** (Friday after close):
   - Review bucket heat and rule adherence
   - Run `UpdateCooldowns` if any bucket had ≥2 stop-outs in last 20-30 bars
   - Run `RefreshReview` to check impulse adherence

---

## Data Validation & Conditional Formatting

### Data Validation (DV)
- Preset dropdown: from `tblPresets[Name]`
- Ticker dropdown: from `tblCandidates[Ticker]` filtered to today
- Sector dropdown: common sector list
- Bucket dropdown: from `tblBuckets[Bucket]`
- Method: "Stock", "Opt-DeltaATR", "Opt-MaxLoss"

### Conditional Formatting (CF)
- Banner cell: GREEN/YELLOW/RED background (set by VBA code, not CF rules)
- Heat preview bars: green ≤70% of cap, amber 70-100%, red >100%

---

## Known Excel Quirks & Fixes

1. **"Removed Records: Formula from /xl/worksheets/sheet4.xml"**
   - Excel dislikes structured refs in dynamic arrays
   - Solution: `CalendarModule.ForceRepairAndRefresh` (uses values only)

2. **Banner color doesn't update**
   - CF rules may conflict with VBA coloring
   - Solution: `ApplyChecklistFixes` (nukes CF from banner cell)

3. **Buttons overlap**
   - Solution: `PlaceButtonsNeatly` (column L anchors)

4. **Events disabled**
   - Solution: `EnsureEventsOn` (called by most VBA subs)

---

## Important Constraints & Rules

1. **Never amend other people's commits** (git safety protocol in docs, though this is not a git repo)
2. **GREEN-only saves**: SaveDecision will reject if banner ≠ GREEN
3. **Cooldown enforcement**: No new entries in a bucket if CooldownActive=TRUE
4. **Impulse brake**: 2-minute delay between Evaluate and Save Decision (stored in Control!A1)
5. **Heat caps are hard limits**: Save will reject if portfolio or bucket heat would exceed caps
6. **No discretionary overrides**: Checklist must be completed as-is (bias minimization)

---

## File Organization

```
/
├── workflow-plan.md                              # High-level workflow & invariants
├── newest-Interactive_TF_Workbook_Plan.md       # Detailed spec for main workbook
├── older-Options_Trend_Dashboard_Summary.md     # Calendar dashboard + VBA code
├── SeykotaTurtleTrend-FollowingOptionsExecution+TradingViewStrategyGuide.md
│                                                 # TradingView strategy guide + Pine code
├── diversification-across-sectors.md            # Sector/bucket risk framework
├── diversification-across-sectors.pdf           # Visual diagrams
├── Options_Trading_Calendar_dynamic.xlsm        # Actual Excel workbook (binary)
└── CLAUDE.md                                    # This file
```

---

## When Modifying VBA Code

### Module Structure (from older-Options_Trend_Dashboard_Summary.md lines 90-1482)

- **CalendarModule**: Date/week logic, calendar grid refresh
- **ChecklistModule**: UI interactions, banner logic, contracts helper
- **Utilities** (embedded in older doc): `SheetExists`, `GetOrCreateSheet`, `GetOrCreateListObject`, `NormalizeTicker`, etc.

### Key Helper Functions

```vba
PortfolioHeatAfter(addR As Double) As Double
  ' Sums TotalOpenR from tblPositions where Status <> "Closed"

BucketHeatAfter(bucket As String, addR As Double) As Double
  ' Sums TotalOpenR for specified bucket

IsBucketInCooldown(bucket As String) As Boolean
  ' Checks tblBuckets CooldownActive and CooldownEndsOn

ImpulseElapsed() As Boolean
  ' Checks if 2 minutes elapsed since Control!A1 timestamp
```

### Testing Approach (Gherkin scenarios in newest-Interactive_TF_Workbook_Plan.md lines 333-401)

Manual testing via scenarios:
- Banner logic (all checks pass → GREEN)
- Impulse timer (attempt save too early → rejected)
- Heat caps (portfolio/bucket exceeded → blocked)
- Candidate gating (ticker not in today's list → blocked)
- Bucket cooldown (active cooldown → blocked)
- Sizing math (Delta-ATR yields correct contracts)

---

## Architecture Rationale

**Why separate TradingView + Excel?**
- TradingView validates mechanical entry/exit signals with backtesting
- Excel enforces portfolio-level constraints (heat caps, cooldowns, impulse control) that TradingView cannot

**Why VBA instead of Python/web app?**
- Traders already use Excel for journaling
- No deployment complexity; works offline
- Forms/buttons provide tactile friction (anti-impulsivity)

**Why correlation buckets?**
- Prevents overconcentration in single sector during sector rotation
- Cooldowns pause entries after sector shows weakness (2+ stop-outs in rolling window)

**Why 2-minute impulse timer?**
- Behavioral circuit-breaker; forces pause between signal recognition and order placement
- Reduces FOMO-driven entries

---

## Migration Notes (if rebuilding from scratch)

1. Duplicate existing `.xlsm` as backup
2. Run `EnsureWorkbookStructure` (creates tables/sheets/names)
3. Run `BuildInteractiveChecklist` (creates Trade Entry UI)
4. Manually wire buttons to procedures (Developer tab → Insert → Button → assign macro)
5. Import historical data to Decisions/Positions if migrating
6. Test with Gherkin scenarios before live use
7. Update Summary named ranges to match real account size

---

## Performance Considerations

- **Calendar refresh**: O(positions × weeks × sectors) — keep Trades table trimmed
- **Heat calculations**: O(open positions) — typically <50 rows
- **Cooldown checks**: O(decisions × buckets) — uses 30-day rolling window
- **Excel file size**: ~150KB without historical data; ~5MB with 1000+ decisions

---

## External Dependencies

- **FINVIZ** (free screener): No API; manual export/paste workflow
- **TradingView** (free or paid plan): Alerts require paid plan; Pine scripts are free
- **Excel** (365 or 2021): Dynamic arrays used in Calendar sheet; VBA for all logic

No Python libraries, no external databases, no web services.

---

## Future Enhancement Ideas (documented in workflow-plan.md)

- Excel VBA macro to scrape FINVIZ directly (current: manual copy/paste)
- TradingView webhook → Excel Cloud for automatic candidate import
- System-1 (20/10) vs System-2 (55/10) toggle in Trade Entry sheet
- Earnings blackout integration (via API or manual calendar)
- Portfolio heat gauge chart on Summary sheet
