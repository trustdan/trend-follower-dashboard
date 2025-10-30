# TF-Engine User Guide

**Version:** 1.0.0
**Last Updated:** 2025-10-29

---

## Table of Contents

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [First-Time Setup](#first-time-setup)
4. [Daily Trading Workflow](#daily-trading-workflow)
5. [Understanding the Banner](#understanding-the-banner)
6. [The 5 Gates Explained](#the-5-gates-explained)
7. [Screen Reference](#screen-reference)
8. [TradingView Integration](#tradingview-integration)
9. [Theme Customization](#theme-customization)
10. [Tips & Best Practices](#tips--best-practices)
11. [Troubleshooting](#troubleshooting)
12. [FAQ](#faq)

---

## Introduction

### What is TF-Engine?

TF-Engine (Trend-Following Engine) is a **discipline enforcement system** for systematic trend-following traders. It's designed to make impulsive trading impossible while making systematic, rule-based trading effortless.

**TF = Trend Following** - The system implements Ed Seykota's System-2 approach:
- **Entry:** 55-bar Donchian breakouts (long > 55-high / short < 55-low)
- **Position Sizing:** Van Tharp ATR-based method with 2×N stops
- **Pyramiding:** Add to winners every 0.5×N up to max units (typically 4)
- **Exits:** 10-bar opposite Donchian OR 2×N stop (whichever is closer)

### What It's NOT

- **NOT** a signal generator (you verify signals manually in TradingView)
- **NOT** an auto-trader (you execute trades manually in your broker)
- **NOT** a backtester (it's for live trading discipline)
- **NOT** a flexible platform (it's a discipline enforcement system)

### Core Philosophy

> **Trade the tide, not the splash.**

The system enforces discipline through **5 hard gates** that cannot be bypassed. If any gate fails, you cannot save a GO decision. This prevents impulsive trades driven by fear, greed, or impatience.

**Key Insight:** This system's value comes from what it **prevents** (bad trades), not just what it allows (good trades).

---

## Installation

See [INSTALLATION_GUIDE.md](INSTALLATION_GUIDE.md) for detailed installation instructions.

**Quick Summary:**
1. Download `TF-Engine-Setup.exe` or `tf-engine.exe` binary
2. Run installer (if using setup) or extract binary to desired location
3. Launch `tf-engine.exe` (double-click or run from command line)
4. Browser opens automatically to http://localhost:8080
5. If browser doesn't open, manually navigate to http://localhost:8080

**System Requirements:**
- Windows 10/11 (64-bit)
- 100 MB disk space
- 512 MB RAM minimum
- Internet connection (for FINVIZ scanning)

---

## First-Time Setup

### Step 1: Initialize Database

On first launch, TF-Engine will automatically create its database in:
```
C:\Users\[YourUsername]\AppData\Roaming\TF-Engine\trading.db
```

If you need to manually initialize:
```bash
tf-engine.exe init
```

### Step 2: Configure Account Settings

![Settings Screenshot](screenshots/settings.png)

1. Navigate to **Settings** (gear icon in header)
2. Enter your account details:
   - **Equity:** Your trading account size (e.g., $100,000)
   - **Risk % per unit:** Typically 0.75% - 1.0% (Ed Seykota uses 1%)
   - **Portfolio heat cap:** 4.0% (total risk across all positions)
   - **Bucket heat cap:** 1.5% (max risk per sector)
   - **Max units:** 4 (initial position + 3 add-ons)
3. Click **Save Settings**

**Why these numbers?**
- **0.75% risk per unit:** With 4 units max, you risk 3% per position (0.75% × 4)
- **4% portfolio cap:** Total risk across ALL positions never exceeds 4% of equity
- **1.5% bucket cap:** Limits concentration in one sector (e.g., Tech/Comm)

**Example with $100,000 account:**
- Risk per unit: $750 (0.75% × $100,000)
- Max risk per position: $3,000 (4 units × $750)
- Portfolio heat cap: $4,000 (4% × $100,000)
- Bucket heat cap: $1,500 (1.5% × $100,000)

### Step 3: Set Up TradingView

You'll use TradingView to verify breakout signals.

1. Create TradingView account (free or paid)
2. Install the **Ed-Seykota.pine** script:
   - Open Pine Editor in TradingView
   - Copy contents from `reference/Ed-Seykota.pine` in installation directory
   - Click "Add to Chart"
3. Save script to favorites for easy access

**The script displays:**
- 55-bar Donchian channels (blue/red lines - entry signals)
- 10-bar Donchian channels (exit signals)
- Current N (ATR value) in indicator window
- Initial stop levels

See [TRADINGVIEW_SETUP.md](TRADINGVIEW_SETUP.md) for detailed guide with screenshots.

### Step 4: Set Up FINVIZ Screener (Optional)

If you want to use the automated daily scan:

1. Create FINVIZ account (free or Elite)
2. Build your screener for trend-following candidates
3. Export your screener and save the URL
4. In TF-Engine Settings → FINVIZ Presets, add preset with your URL

**Suggested filters for long breakouts:**
- **Technical:** Price > SMA200 (in uptrend)
- **Technical:** RSI > 55 (momentum)
- **Fundamental:** Avg Volume > 1M (liquidity)
- **Fundamental:** Market Cap > $500M (avoid micro-caps)

**Example FINVIZ URL:**
```
https://finviz.com/screener.ashx?v=111&f=ta_price_a200sma,ta_rsi_os55,sh_avgvol_o1000,sh_marketcap_o500
```

---

## Daily Trading Workflow

### Morning Routine (30 minutes)

![Dashboard Screenshot](screenshots/dashboard.png)

#### 1. Check Dashboard

Start your day at the **Dashboard**:
- Review open positions (ticker, entry, current stop, risk, days held)
- Check portfolio heat gauge (should be well below 4% cap)
- Note sector distribution (avoid concentration)
- Check for cooldowns (sectors or tickers you cannot trade)

**Red Flags:**
- Portfolio heat > 90% of cap → Close positions before adding new ones
- One sector has 3+ positions → Diversify into other sectors
- Multiple cooldowns active → Review recent losing trades

#### 2. Scan for Candidates

![FINVIZ Scanner Screenshot](screenshots/scanner.png)

Navigate to **Scanner** → Click **"Run Daily FINVIZ Scan"**

The system fetches tickers from your FINVIZ screener and displays:
- Ticker symbols
- Sector buckets (Tech/Comm, Energy, Financials, etc.)
- Last close price
- Average daily volume

**Select 10-15 candidates** to analyze in detail. Click **"Import Selected"**.

**Tips:**
- Focus on sectors where you're underweight (check Calendar view)
- Avoid sectors on cooldown
- Prioritize high-volume stocks (> 5M shares/day)

#### 3. Verify Signals in TradingView

For each candidate:

![TradingView Link](screenshots/tradingview-link.png)

1. Click **"Open in TradingView"** next to ticker
2. Chart opens in new tab with your Ed-Seykota script loaded
3. **Verify the breakout:**
   - **For longs:** Today's close > 55-bar high (blue line)
   - **For shorts:** Today's close < 55-bar low (red line)
4. Note the **N value** (ATR) from script indicator window
5. Note your **entry price** (current close, or your planned limit order)
6. Note the **sector bucket** (Tech/Comm, Energy, etc.)

**Quality Check:**
- Is this a clean breakout or a choppy one?
- Is the stock near 20-EMA or overextended?
- Is volume above average?
- Are earnings coming up in next 2 weeks?

**Return to TF-Engine** when ready to complete checklist.

---

### Checklist Evaluation (10 minutes)

![Checklist Screenshot](screenshots/checklist.png)

Navigate to **Checklist**:

#### Step 1: Enter Trade Data

Fill in the form:
- **Ticker:** AAPL
- **Entry Price:** 180.50
- **N (ATR):** 2.35 (from TradingView script)
- **Sector:** Tech/Comm
- **Structure:** Stock (or Call/Put for options)

#### Step 2: Check All 5 Required Gates

These are **mandatory**. The banner stays RED until all are checked.

1. **✓ Signal:** 55-bar Donchian breakout confirmed in TradingView
   - You verified this visually on the chart

2. **✓ Risk/Size:** Will use 2×N stop, add every 0.5×N, max 4 units
   - This is the Van Tharp method (system calculates automatically)

3. **✓ Liquidity:** Stock avg volume > 1M (or options OI > 100)
   - Check on TradingView or Yahoo Finance
   - For options: Open Interest > 100 contracts, bid-ask spread < 10%

4. **✓ Exits:** Will exit on 10-bar opposite Donchian OR 2×N stop
   - Whichever is closer to current price
   - You promise to honor this rule

5. **✓ Behavior:** Not on cooldown, heat OK, will honor 2-min timer
   - System validates heat and cooldown in final gates check
   - You commit to waiting the full 2 minutes

**Banner turns YELLOW** when all required are checked (but quality score might still be low).

#### Step 3: Improve Quality Score (Optional but Recommended)

These are **optional** quality items. Each checked adds to your quality score.

6. **✓ Regime OK:** SPY > 200 SMA (favorable market for long positions)
   - Check SPY chart on TradingView
   - For shorts, check if SPY < 200 SMA

7. **✓ No Chase:** Entry within 2N of 20-EMA (not overextended)
   - Visual check on TradingView
   - Avoids buying into parabolic moves

8. **✓ Earnings OK:** No earnings within next 2 weeks (for options)
   - Check earnings calendar on Yahoo Finance or TradingView
   - Earnings can cause gap risk that violates 2×N stop

9. **✓ Journal Note:** Why this trade now? (1-2 sentences)
   - Example: "Clean breakout above 6-month consolidation, volume 2x average, Tech sector underweight in portfolio"
   - Forces you to articulate your thesis

**Quality Score Display:**
```
Quality Score: 4 / 4 checked
Threshold: 3.0 (from settings)
Result: ABOVE THRESHOLD ✓
```

**Banner turns GREEN** when score ≥ threshold (e.g., 3 out of 4).

#### Step 4: Save Evaluation

![Green Banner](screenshots/banner-green.png)

Click **"Save Evaluation"**

**The 2-minute cool-off timer starts:** 2:00... 1:59... 1:58...

![Timer](screenshots/timer.png)

This is the **impulse brake**. You cannot save a GO decision until the timer reaches 0:00.

**Purpose:** Prevents emotional, impulsive decisions. Forces you to slow down.

**Use this time productively:**
- Double-check your TradingView chart analysis
- Review the trade plan and exits
- Calculate expected profit at 1R, 2R, 3R targets
- Review sector diversification in Calendar view
- Calculate position sizing (next step)

**Cannot be bypassed.** This is a core anti-impulsivity feature.

---

### Position Sizing (5 minutes)

![Position Sizing Screenshot](screenshots/position-sizing.png)

Navigate to **Position Sizing**:

The form is pre-filled with your checklist data.

**Inputs:**
- **Method:** Stock (or opt-delta-atr / opt-contracts for options)
- **Entry:** $180.50 (from checklist)
- **ATR (N):** $2.35 (from checklist)
- **K multiple:** 2.0 (stop distance in multiples of N)
- **Max units:** 4 (from settings)

Click **"Calculate Position Size"**

**Results Display:**

```
Van Tharp Position Sizing Results:

Shares per unit: 159
Risk per unit: $747.30
Initial stop: $175.80 (entry - 2×N)

Add-on schedule (pyramid every 0.5×N):
  Unit 1: 159 shares @ $180.50 (NOW - initial entry)
  Unit 2: 159 shares @ $181.68 (+0.5N = +$1.18)
  Unit 3: 159 shares @ $182.85 (+1.0N = +$2.35)
  Unit 4: 159 shares @ $184.03 (+1.5N = +$3.53)

Total max position: 636 shares (if all 4 units filled)
Total max risk: $2,989.20 (4 × $747.30)

Exit plan:
  - 10-bar Donchian (check daily)
  - OR $175.80 stop (2×N below entry)
  - Whichever is CLOSER to current price
```

**Concentration Warning:**

If position would be > 25% of equity:
```
⚠️ WARNING: This position is 18.2% of your $100,000 account.
Consider reducing size or diversifying across more positions.
```

**Verify the calculation makes sense:**
- Risk per unit = Equity × Risk% = $100,000 × 0.0075 = $750 ✓
- Stop distance = 2 × $2.35 = $4.70
- Shares = floor($750 / $4.70) = 159 shares ✓
- Actual risk = 159 × $4.70 = $747.30 ≈ $750 ✓

Click **"Save Position Plan"**

---

### Heat Check (2 minutes)

![Heat Check Screenshot](screenshots/heat-check.png)

Navigate to **Heat Check**:

The system displays current risk exposure:

**Current Portfolio Heat:**
```
Total risk across all positions: $2,890 / $4,000 cap (72.3%)
Remaining capacity: $1,110
```

**Current Bucket Heat:**
```
Tech/Comm:  $1,125 / $1,500 cap (75.0%)
Energy:     $850 / $1,500 cap (56.7%)
Financials: $915 / $1,500 cap (61.0%)
```

Click **"Check Heat for This Trade"**

**Proposed trade:** $747 risk in Tech/Comm sector

**Result (Example 1 - Within Caps):**

```
✓ Portfolio heat: $3,637 / $4,000 (90.9%) - WITHIN CAP
✓ Bucket heat (Tech/Comm): $1,872 / $1,500 (125%) - EXCEEDS CAP by $372

VERDICT: TRADE REJECTED - Bucket heat cap exceeded
```

![Heat Warning](screenshots/heat-warning.png)

**If cap exceeded, you have 3 options:**

1. **Reduce position size:**
   - Click "Calculate Max Shares for Bucket Cap"
   - System suggests: 79 shares (instead of 159)
   - Keeps you within $1,500 bucket cap

2. **Close an existing Tech/Comm position first:**
   - Review Dashboard for positions to exit
   - Close one, then recheck heat

3. **Choose a different ticker from another sector:**
   - Go back to Scanner/Checklist
   - Pick from Energy or Financials (more capacity)

**Result (Example 2 - Within Caps):**

```
✓ Portfolio heat: $3,637 / $4,000 (90.9%) - WITHIN CAP
✓ Bucket heat (Energy): $1,597 / $1,500 (106%) - EXCEEDS CAP by $97

VERDICT: TRADE REJECTED - Bucket heat cap exceeded
```

**Result (Example 3 - All Clear):**

```
✓ Portfolio heat: $2,137 / $4,000 (53.4%) - WITHIN CAP
✓ Bucket heat (Energy): $747 / $1,500 (49.8%) - WITHIN CAP

VERDICT: TRADE APPROVED - Proceed to Trade Entry
```

**Heat management is non-negotiable.** If caps are exceeded, trade is blocked.

---

### Final Trade Entry (2 minutes)

![Trade Entry Screenshot](screenshots/trade-entry.png)

Navigate to **Trade Entry**:

**Review Trade Summary:**

```
Ticker:     AAPL
Direction:  LONG
Entry:      $180.50
Structure:  Stock
Shares:     159 per unit (max 636 total)
Stop:       $175.80 (2×N below entry)
Risk:       $747 per unit (max $2,989 total)
Sector:     Tech/Comm
ATR (N):    $2.35
Quality:    4/4 (GREEN banner)
```

Click **"Run Final Gate Check"**

![Gate Check Results](screenshots/gates-all-pass.png)

**The 5 Gates Validate:**

```
Gate 1: Banner Status        → GREEN ✓
  - All required items checked
  - Quality score 4/4 ≥ 3.0 threshold

Gate 2: Impulse Brake        → 2:15 elapsed ✓
  - Checklist saved at 09:23:45
  - Current time: 09:26:00
  - Required: 2:00 minimum

Gate 3: Cooldown Status      → Not on cooldown ✓
  - Ticker AAPL: No cooldown active
  - Sector Tech/Comm: No cooldown active

Gate 4: Heat Caps            → Within caps ✓
  - Portfolio: $2,137 / $4,000 (53%)
  - Bucket: $747 / $1,500 (50%)

Gate 5: Sizing Completed     → Plan saved ✓
  - 159 shares per unit
  - Risk $747 per unit
  - Stop $175.80

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Result: ALL GATES PASS ✓
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**"SAVE GO DECISION" button is now ENABLED (green).**

Click **"SAVE GO DECISION"**

**Success notification:**
```
✓ GO decision saved for AAPL at 09:26:00
  Trade ID: #147
  Proceed to broker to execute trade
```

**If any gate fails:**

```
Gate 4: Heat Caps            → EXCEEDS CAP ✗
  - Portfolio: $4,250 / $4,000 (106%)
  - Overage: $250

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Result: GATE CHECK FAILED ✗
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

"SAVE GO DECISION" button remains DISABLED.
```

You cannot save GO decision. **"SAVE NO-GO DECISION"** button is always enabled to record why trade was rejected.

---

### Execute in Broker (Manual)

**TF-Engine does NOT execute trades automatically.** You must manually enter the trade in your broker.

**Why manual execution?**
- Final human verification before committing capital
- Prevents catastrophic auto-trading errors
- Keeps you engaged and present in the process
- You remain responsible for your trades

**Execution Steps:**

1. Open your broker platform (Interactive Brokers, TD Ameritrade, Fidelity, etc.)
2. Enter order:
   ```
   BUY 159 shares AAPL @ $180.50
   Type: Limit order (or Market if acceptable)
   Time: Day order
   ```
3. Set initial stop-loss order:
   ```
   SELL 159 shares AAPL @ $175.80 STOP
   ```
4. Record fill in TF-Engine (optional):
   - Navigate to Dashboard → Add Position
   - Or manually track in spreadsheet

**Set alerts for add-on levels:**
- Alert @ $181.68: Add unit 2 (+0.5N)
- Alert @ $182.85: Add unit 3 (+1.0N)
- Alert @ $184.03: Add unit 4 (+1.5N)

**Daily management:**
- Check 10-bar Donchian exit level daily
- Update stop as 10-bar moves (trailing stop)
- Exit if price touches 10-bar OR 2×N stop (whichever closer)

---

## Understanding the Banner

The banner is the **heart of the discipline enforcement system**. It provides immediate visual feedback on trade readiness.

### RED Banner: DO NOT TRADE

![Red Banner](screenshots/banner-red.png)

**Visual:** Large red gradient banner at top of Checklist screen

**Meaning:** One or more required gates are unchecked.

**Message Examples:**
- "STOP - 2 required items missing"
- "STOP - Signal not verified"

**What to do:**
- Review the 5 required gates
- Check all that apply
- Do NOT proceed to position sizing or trade entry
- Do NOT try to bypass (system prevents it)

**Common causes:**
- Forgot to check "Signal" after TradingView verification
- Skipped "Liquidity" check
- Haven't verified exits plan

### YELLOW Banner: CAUTION

![Yellow Banner](screenshots/banner-yellow.png)

**Visual:** Large yellow/orange gradient banner

**Meaning:** All 5 required gates pass, but quality score is below threshold.

**Message Examples:**
- "CAUTION - Quality score 2.0 < 3.0 threshold"
- "CAUTION - 2 quality items recommended"

**What to do:**
- Review optional quality items (items 6-9)
- Add regime check (SPY vs 200 SMA)
- Add chase check (price vs 20-EMA)
- Add earnings check
- Add journal note (always recommended!)
- **Or** accept lower quality trade (not recommended but allowed)

**Trade is technically allowed but not optimal.**

### GREEN Banner: OK TO TRADE

![Green Banner](screenshots/banner-green.png)

**Visual:** Large green gradient banner

**Meaning:** All 5 required gates pass AND quality score ≥ threshold.

**Message Examples:**
- "GO - All gates pass, quality score 4.0 ≥ 3.0"
- "GO - Ready to proceed"

**What to do:**
- Proceed to position sizing
- Complete heat check
- Complete final trade entry
- Save GO decision (after 2-min timer and gates check)

**Important:** GREEN banner is **necessary** but not **sufficient**. The 5 gates check at Trade Entry is the final validation (checks heat caps, cooldowns, timer, etc.).

---

## The 5 Gates Explained

The 5 gates are the **final validation** before you can save a GO decision. They cannot be bypassed.

### Gate 1: Banner Status

**Rule:** Checklist banner must be GREEN.

**Why:** Ensures all required pre-trade checks are complete and quality score meets threshold.

**Validates:**
- All 5 required items checked (signal, risk, liquidity, exits, behavior)
- Quality score ≥ threshold (e.g., 3.0 out of 4.0)

**Failure modes:**
- Banner is RED: Missing required items → Go back and check them
- Banner is YELLOW: Quality score too low → Add optional quality checks

**Fix:** Complete checklist properly. No shortcuts.

### Gate 2: Impulse Brake (2-Minute Cool-Off)

**Rule:** At least 2 minutes must elapse between saving checklist evaluation and saving GO decision.

**Why:** Forces you to slow down. Prevents emotional, impulsive decisions made in the heat of the moment.

**Validates:**
- Timestamp of checklist evaluation (saved in database)
- Current timestamp when attempting to save GO
- Difference must be ≥ 120 seconds

**Failure modes:**
- Trying to save GO at 1:45 elapsed → Too soon, wait 15 more seconds
- System clock manipulated → Database timestamp is authoritative

**Fix:** Wait the full 2 minutes. Use the time productively (calculate sizing, review chart, check diversification).

**Cannot be reduced or disabled.** This is a core design principle.

### Gate 3: Cooldown Status

**Rule:** Neither the ticker nor its sector bucket can be on cooldown.

**Why:** After a losing trade, emotions are high. Taking a break prevents revenge trading and allows for reflection.

**Validates:**
- Ticker cooldown: Check if AAPL is on cooldown (after recent loss)
- Bucket cooldown: Check if Tech/Comm sector is on cooldown (after multiple losses)

**Cooldown triggers:**
- Lose -1R or more on a trade → Ticker goes on cooldown for 1-2 weeks
- Lose multiple trades in same sector → Bucket goes on cooldown for 2-4 weeks
- (Exact cooldown durations configurable in settings)

**Failure modes:**
- "Ticker AAPL on cooldown until Nov 5" → Choose different ticker
- "Sector Tech/Comm on cooldown until Nov 12" → Choose different sector

**Fix:**
- Wait for cooldown to expire (recommended)
- Trade a different ticker in a different sector
- **Cannot be bypassed** (by design)

### Gate 4: Heat Caps

**Rule:** Proposed trade must not exceed portfolio heat cap (4%) or bucket heat cap (1.5%).

**Why:** Limits overall risk and prevents concentration in one sector. Core risk management.

**Validates:**
- Portfolio heat: Sum of risk across ALL open positions + proposed trade ≤ 4% equity
- Bucket heat: Sum of risk in ONE sector + proposed trade ≤ 1.5% equity

**Example with $100,000 account:**
- Portfolio cap: $4,000
- Bucket cap: $1,500

**Failure modes:**
- "Portfolio heat $4,250 exceeds cap $4,000 (overage $250)"
- "Tech/Comm bucket heat $1,650 exceeds cap $1,500 (overage $150)"

**Fix:**
1. Reduce position size (use "Calculate Max Shares" helper)
2. Close an existing position first
3. Choose different ticker in less crowded sector

**Cannot be bypassed.** If caps exceeded, trade is blocked.

### Gate 5: Sizing Completed

**Rule:** Position sizing must be calculated and saved before saving GO decision.

**Why:** Cannot enter a trade without knowing how many shares and what the risk is.

**Validates:**
- Position sizing calculation completed
- Results saved to database
- Shares, risk, stop levels all determined

**Failure modes:**
- "Position sizing not completed" → Go to Position Sizing screen
- Sizing data missing from database → Recalculate and save

**Fix:** Navigate to Position Sizing screen, calculate, and save results.

**Why this matters:** Entering trades without pre-defined sizing is gambling, not systematic trading.

---

## Screen Reference

### Dashboard

**Purpose:** Portfolio overview and daily starting point.

**Location:** Home screen (default on launch)

**Key Elements:**

1. **Portfolio Summary Card:**
   - Total equity
   - Current portfolio heat ($$ and % of cap)
   - Portfolio heat gauge (visual bar)
   - Number of open positions
   - Available capacity

2. **Open Positions Table:**
   - Ticker symbol
   - Entry price
   - Current stop level
   - Risk per unit ($$$)
   - Days held
   - Sector bucket
   - Actions (view detail, close position)

3. **Today's Candidates Count:**
   - Number of imported candidates
   - Link to Scanner

4. **Cooldowns Section:**
   - Active ticker cooldowns (if any)
   - Active sector cooldowns (if any)
   - Expiration dates

5. **Quick Actions:**
   - "Run FINVIZ Scan" button
   - "New Checklist" button
   - "View Calendar" button

**Typical Daily Flow:**
- Start here every morning
- Review positions and heat
- Note cooldowns
- Click "Run FINVIZ Scan" to begin workflow

---

### Scanner

**Purpose:** Daily FINVIZ scan to find trend-following candidates.

**Location:** Main navigation → Scanner

**Key Elements:**

1. **Scan Controls:**
   - "Run Daily FINVIZ Scan" button (primary action)
   - FINVIZ preset selector (if multiple presets configured)
   - Scan status (running... / complete / error)

2. **Results Table:**
   - Ticker symbol (sortable)
   - Sector bucket (sortable, filterable)
   - Last close price
   - Average volume (30-day)
   - Checkbox for selection

3. **Filters:**
   - Filter by sector (Tech/Comm, Energy, etc.)
   - Filter by price range
   - Filter by volume threshold

4. **Sector Distribution Summary:**
   - Bar chart showing candidates per sector
   - Helps identify balanced opportunities

5. **Actions:**
   - "Import Selected" button
   - "Open in TradingView" link for each ticker

**Typical Workflow:**
1. Click "Run Daily FINVIZ Scan"
2. Wait 3-5 seconds for results
3. Filter/sort by sector if needed
4. Select 10-15 candidates (checkboxes)
5. Click "Import Selected"
6. Click "Open in TradingView" for each to verify signals

---

### Checklist

**Purpose:** Evaluate trade setup and calculate banner state.

**Location:** Main navigation → Checklist

**Key Elements:**

1. **Large 3-State Banner (most prominent element):**
   - RED: "STOP - X required items missing"
   - YELLOW: "CAUTION - Quality score below threshold"
   - GREEN: "GO - Ready to proceed"

2. **Trade Data Inputs:**
   - Ticker (text input)
   - Entry price (number)
   - N / ATR (number)
   - Sector bucket (dropdown: Tech/Comm, Energy, Financials, etc.)
   - Structure (dropdown: Stock, Call, Put, etc.)

3. **Required Gates Section (5 checkboxes):**
   - ☐ Signal: 55-bar Donchian breakout confirmed
   - ☐ Risk/Size: 2×N stop, add every 0.5×N, max 4 units
   - ☐ Liquidity: Volume > 1M (or OI > 100 for options)
   - ☐ Exits: 10-bar Donchian OR 2×N stop (closer)
   - ☐ Behavior: Not on cooldown, heat OK, will wait 2 min

4. **Optional Quality Items Section (4 checkboxes + textarea):**
   - ☐ Regime OK: SPY vs 200 SMA favorable
   - ☐ No Chase: Entry within 2N of 20-EMA
   - ☐ Earnings OK: No earnings next 2 weeks
   - ☐ Journal Note: (textarea for 1-2 sentence thesis)

5. **Quality Score Display:**
   - "Quality Score: 3 / 4 checked"
   - "Threshold: 3.0 (from settings)"
   - "Result: ABOVE THRESHOLD ✓" or "BELOW THRESHOLD ✗"

6. **Action Button:**
   - "Save Evaluation" (enabled when data filled)

7. **2-Minute Timer (appears after save):**
   - Large countdown: "2:00" → "1:59" → ... → "0:00"
   - Message: "Impulse brake - wait X:XX before saving GO decision"
   - Visual progress bar

**Typical Workflow:**
1. Fill trade data from TradingView
2. Check all 5 required gates
3. Check optional quality items
4. Watch banner turn GREEN
5. Click "Save Evaluation"
6. Use 2-minute timer to calculate sizing

---

### Position Sizing

**Purpose:** Calculate shares/contracts using Van Tharp ATR method.

**Location:** Main navigation → Position Sizing

**Key Elements:**

1. **Method Selector (dropdown):**
   - Stock (default)
   - opt-delta-atr (options with delta adjustment)
   - opt-contracts (options with contract risk)

2. **Pre-filled Trade Data:**
   - Ticker (from checklist)
   - Entry price (from checklist)
   - ATR / N (from checklist)
   - K multiple (default 2.0, usually not changed)
   - Max units (from settings, usually 4)

3. **Account Data (from settings):**
   - Equity (read-only)
   - Risk % per unit (read-only)

4. **Calculate Button:**
   - "Calculate Position Size" (primary action)

5. **Results Display (after calculation):**
   - **Shares per unit:** 159
   - **Risk per unit:** $747.30
   - **Initial stop:** $175.80
   - **Add-on schedule table:**
     - Unit 1: 159 @ $180.50 (NOW)
     - Unit 2: 159 @ $181.68 (+0.5N)
     - Unit 3: 159 @ $182.85 (+1.0N)
     - Unit 4: 159 @ $184.03 (+1.5N)
   - **Total max position:** 636 shares
   - **Total max risk:** $2,989.20
   - **Exit plan:** 10-bar Donchian OR $175.80 stop

6. **Concentration Warning (if applicable):**
   - "⚠️ Position is 18% of equity - consider reducing"

7. **Action Button:**
   - "Save Position Plan" (enabled after calculation)

**Typical Workflow:**
1. Data pre-filled from checklist
2. Click "Calculate Position Size"
3. Review results for reasonableness
4. Note add-on levels
5. Click "Save Position Plan"

---

### Heat Check

**Purpose:** Verify proposed trade doesn't exceed risk caps.

**Location:** Main navigation → Heat Check

**Key Elements:**

1. **Current Portfolio Heat Gauge:**
   - Visual bar chart: $2,890 / $4,000 cap (72%)
   - Color-coded: Green < 75%, Yellow 75-90%, Red > 90%
   - Remaining capacity: $1,110

2. **Current Bucket Heat Table:**
   - Rows: Each sector bucket
   - Columns: Heat ($), Cap ($), % Used, Remaining
   - Example:
     - Tech/Comm: $1,125 / $1,500 (75%) → $375 remaining
     - Energy: $850 / $1,500 (57%) → $650 remaining

3. **Proposed Trade Info:**
   - Ticker: AAPL
   - Risk per unit: $747
   - Sector: Tech/Comm

4. **Check Button:**
   - "Check Heat for This Trade" (primary action)

5. **Results Display (after check):**
   - **Portfolio Heat Check:**
     - New total: $3,637 / $4,000 (91%)
     - Result: ✓ WITHIN CAP or ✗ EXCEEDS CAP by $XXX
   - **Bucket Heat Check:**
     - New Tech/Comm: $1,872 / $1,500 (125%)
     - Result: ✓ WITHIN CAP or ✗ EXCEEDS CAP by $372
   - **Overall Verdict:**
     - "✓ TRADE APPROVED" (green) or "✗ TRADE REJECTED" (red)

6. **Helper Actions (if cap exceeded):**
   - "Calculate Max Shares for Portfolio Cap" → suggests reduced size
   - "Calculate Max Shares for Bucket Cap" → suggests reduced size
   - Link to Dashboard to review existing positions

**Typical Workflow:**
1. Navigate after saving position sizing
2. Click "Check Heat for This Trade"
3. If within caps: Proceed to Trade Entry
4. If exceeds caps: Adjust size or choose different trade

---

### Trade Entry

**Purpose:** Final 5-gate validation before saving GO/NO-GO decision.

**Location:** Main navigation → Trade Entry

**Key Elements:**

1. **Trade Summary Card:**
   - All trade details (ticker, entry, stop, risk, sector, quality)
   - Read-only review

2. **Gate Check Button:**
   - "Run Final Gate Check" (primary action)

3. **Gate Results Display (after check):**
   - **Gate 1: Banner Status**
     - GREEN ✓ or RED/YELLOW ✗
   - **Gate 2: Impulse Brake**
     - "2:15 elapsed ✓" or "1:45 elapsed ✗ (wait 15s)"
   - **Gate 3: Cooldown Status**
     - "Not on cooldown ✓" or "On cooldown ✗"
   - **Gate 4: Heat Caps**
     - "Within caps ✓" or "Exceeds caps ✗"
   - **Gate 5: Sizing Completed**
     - "Plan saved ✓" or "Not completed ✗"
   - **Overall Result:**
     - "ALL GATES PASS ✓" (green) or "GATE CHECK FAILED ✗" (red)

4. **Action Buttons:**
   - **"SAVE GO DECISION"** (green, large)
     - Enabled ONLY if all gates pass
     - Disabled (grayed out) if any gate fails
   - **"SAVE NO-GO DECISION"** (red, smaller)
     - Always enabled
     - Prompts for reason (dropdown + textarea)

5. **Decision History (bottom):**
   - List of recent GO/NO-GO decisions
   - Timestamp, ticker, result, reason (if NO-GO)

**Typical Workflow:**
1. Navigate after heat check passes
2. Review trade summary
3. Click "Run Final Gate Check"
4. If all gates pass: Click "SAVE GO DECISION"
5. If any gate fails: Fix issue or click "SAVE NO-GO DECISION"
6. Execute trade manually in broker

---

### Calendar

**Purpose:** Visualize sector diversification over rolling 10-week period.

**Location:** Main navigation → Calendar

**Key Elements:**

1. **10-Week Grid:**
   - Columns: Weeks (Mon-Sun date ranges)
   - Rows: Sector buckets (Tech/Comm, Energy, Financials, etc.)
   - Range: 2 weeks back, current week, 7 weeks forward

2. **Cell Contents:**
   - Tickers active in that sector/week
   - Example: "AAPL, MSFT" in Tech/Comm week of Nov 5-11
   - Color-coded by heat level:
     - Green: Low heat (< 50% of bucket cap)
     - Yellow: Medium heat (50-80% of bucket cap)
     - Red: High heat (> 80% of bucket cap)

3. **Current Week Indicator:**
   - Bold border or highlight on current week column

4. **Tooltips (hover over cell):**
   - Position details:
     - AAPL: Entry $180.50, Risk $747, Day 3
     - MSFT: Entry $350.20, Risk $820, Day 7

5. **Summary Stats (bottom):**
   - Diversification score (0-100)
   - Sectors covered this week: 4 / 7
   - Weeks with no trades: 2
   - Concentration warnings (if any)

**Typical Use Cases:**
- Review on Sunday evening before trading week
- Identify gaps (weeks with no trades or missing sectors)
- Avoid crowded sectors when selecting new candidates
- Visualize portfolio balance over time

---

### Settings

**Purpose:** Configure account settings and system parameters.

**Location:** Main navigation → Settings (gear icon)

**Key Elements:**

1. **Account Settings:**
   - Equity (number input, $)
   - Risk % per unit (number input, %, default 0.75)
   - Portfolio heat cap (number input, %, default 4.0)
   - Bucket heat cap (number input, %, default 1.5)
   - Max units (number input, default 4)

2. **Checklist Settings:**
   - Quality score threshold (number input, default 3.0)

3. **TradingView Integration:**
   - URL template (text input, optional)
   - Example: `https://tradingview.com/chart/?symbol={ticker}`

4. **FINVIZ Presets (table):**
   - Name (e.g., "TF Breakout Long")
   - URL (FINVIZ screener URL)
   - Actions: Edit, Delete, Add New

5. **Cooldown Settings:**
   - Ticker cooldown duration (days, default 7)
   - Bucket cooldown duration (days, default 14)

6. **Action Buttons:**
   - "Save Settings" (primary)
   - "Reset to Defaults" (secondary, confirmation required)

**Typical Use:**
- First-time setup: Enter equity and risk parameters
- Periodic updates: Adjust equity as account grows
- Rare changes: Modify caps or max units (not recommended)

---

## TradingView Integration

TF-Engine integrates with TradingView for signal verification. You verify breakouts manually using the Ed-Seykota Pine Script.

### Quick Setup

1. **Create TradingView account:** https://www.tradingview.com (free or paid)
2. **Install Pine Script:**
   - Navigate to Chart → Pine Editor (bottom panel)
   - Copy contents of `reference/Ed-Seykota.pine` from installation directory
   - Click "Add to Chart"
   - Pin script to favorites
3. **Configure URL template (optional):**
   - TF-Engine Settings → TradingView URL Template
   - Enter: `https://tradingview.com/chart/?symbol={ticker}`
   - System replaces `{ticker}` with actual ticker when you click "Open in TradingView"

### What the Pine Script Shows

![TradingView with Ed-Seykota Script](screenshots/tradingview-script.png)

**On Chart:**
- **Blue lines:** 55-bar Donchian channel high (long entry trigger)
- **Red lines:** 55-bar Donchian channel low (short entry trigger)
- **Green lines:** 10-bar Donchian channel (exit trigger)
- **Orange lines:** 2×N stops (calculated from entry)

**In Indicator Window:**
- **N (ATR):** Current ATR value (e.g., 2.35)
- **Initial Stop:** Entry - 2×N for longs (Entry + 2×N for shorts)

### Signal Verification Workflow

**For Long Entry:**

1. Click "Open in TradingView" for ticker AAPL
2. TradingView chart opens with Ed-Seykota script
3. **Check breakout:**
   - Is today's close > 55-bar high (blue line)? YES ✓
4. **Note values:**
   - N (ATR): 2.35
   - Entry: $180.50 (today's close or your limit order)
   - Initial stop: $175.80 (shown on chart)
5. Return to TF-Engine Checklist and fill in:
   - Entry: 180.50
   - N: 2.35
   - Check "Signal" gate ✓

**For Short Entry:**

1. Check if today's close < 55-bar low (red line)
2. Initial stop = Entry + 2×N (above entry for shorts)
3. Add-ons work in reverse (subtract 0.5×N)

### Daily Exit Management

**Two exit rules (whichever is closer):**

1. **10-bar Donchian opposite breakout:**
   - For longs: Close < 10-bar low (green line)
   - For shorts: Close > 10-bar high (green line)
   - Check daily on TradingView

2. **2×N stop:**
   - For longs: Stop = Entry - 2×N = $175.80
   - For shorts: Stop = Entry + 2×N
   - Set stop-loss order in broker

**Each day:**
- Open TradingView chart
- Check current 10-bar level
- Compare to 2×N stop
- Exit if price touches EITHER level

**Trailing stops:**
- As 10-bar level rises (for longs), move stop higher
- Stop becomes: max(10-bar low, original 2×N stop)
- Protects profits as trade moves in your favor

### Tips

- Save chart layout in TradingView for fast loading
- Add multiple tickers to watchlist for quick switching
- Use TradingView mobile app for on-the-go monitoring
- Set price alerts at add-on levels (every 0.5×N)

For detailed setup instructions with screenshots, see [TRADINGVIEW_SETUP.md](TRADINGVIEW_SETUP.md).

---

## Theme Customization

TF-Engine supports light and dark modes.

**Day Mode (Light):**
- Clean, bright interface
- High contrast for readability
- Best for well-lit environments

**Night Mode (Dark):**
- Easy on eyes for extended use
- Reduced blue light
- Gradient banners remain vibrant
- Best for evening trading sessions

**Toggle Theme:**
- Click sun/moon icon in header (top-right corner)
- Theme switches instantly
- Preference saved automatically to browser localStorage

**Persistence:**
- Theme preference persists across sessions
- Survives browser close/reopen
- Per-browser (not synced across devices)

**Accessibility:**
- Both themes meet WCAG AA contrast standards
- Banner colors remain distinct in both modes
- All text remains readable

---

## Tips & Best Practices

### 1. Start with Small Position Sizes

If new to this system, reduce your risk % per unit:
- **Beginner:** 0.50% per unit (max 2% per position)
- **Intermediate:** 0.75% per unit (max 3% per position)
- **Experienced:** 1.00% per unit (max 4% per position)

**Why?** Smaller positions reduce emotional attachment and help you learn the system without large drawdowns.

### 2. Respect the 2-Minute Timer

Don't try to rush through it. Use the time productively:
- Re-check your TradingView chart (still a valid breakout?)
- Review sector diversification in Calendar view
- Calculate potential profit at 1R, 2R, 3R targets
- Check news headlines (any breaking news?)
- Take 3 deep breaths (seriously - it helps)

### 3. Journal Every Trade

Even if quality score doesn't require it, **always add a journal note**. Future you will thank present you.

**Good journal notes (specific, actionable):**
- "Clean breakout above 6-month consolidation, volume 2x avg, Tech sector underweight in portfolio, SPY confirmed bullish"
- "Pullback to 20-EMA in strong uptrend, RSI reset to 55, no earnings for 6 weeks, sector rotation into Energy"

**Bad journal notes (vague, unhelpful):**
- "Looks good"
- "Strong momentum"
- "Like the chart"

**After 10-20 trades, review your journal:**
- Which setups worked best?
- Which journal notes correlated with winners?
- What patterns do you notice in losers?

### 4. Use the Calendar Weekly

**Sunday Evening Routine (15 minutes):**
- Open TF-Engine → Calendar view
- Review the 10-week grid
- Identify:
  - **Gaps:** Weeks with no trades (are you being too picky?)
  - **Crowding:** Weeks with 5+ positions (too concentrated?)
  - **Sector bias:** One sector with 3+ positions (diversify!)
- Plan for upcoming week:
  - Which sectors need more coverage?
  - Any cooldowns expiring this week?
  - How much heat capacity available?

### 5. Save NO-GO Decisions

When a trade fails the gates, **save it as NO-GO with a reason**. This creates a valuable record of what you DON'T trade.

**Why?**
- Prevents re-evaluating same ticker repeatedly
- Documents near-misses (setup almost worked but failed one gate)
- Shows discipline (you followed the rules, even when tempted)
- Provides learning data (why did it fail gates?)

**Example NO-GO reasons:**
- "Heat cap exceeded by $250 - Tech sector too crowded"
- "Impulse brake not elapsed - only 1:20 passed"
- "Ticker on cooldown until Nov 5 after -1.5R loss"

### 6. Backup Your Database

**Location:** `C:\Users\[YourUsername]\AppData\Roaming\TF-Engine\trading.db`

**Backup Frequency:**
- Weekly (Sunday evening after review)
- After significant activity (10+ trades)
- Before major system updates

**Backup Method:**
1. Close TF-Engine (stop tf-engine.exe process)
2. Navigate to `%APPDATA%\TF-Engine\`
3. Copy `trading.db` to:
   - External drive
   - Cloud storage (Dropbox, Google Drive)
   - Network drive
4. Rename with date: `trading-backup-2025-10-29.db`

**Restore from Backup:**
1. Close TF-Engine
2. Replace `trading.db` with backup file
3. Restart TF-Engine

### 7. Avoid Backseat Trading

**Don't second-guess closed trades.**

Once you've saved a GO decision and executed in broker:
- The decision is final
- Don't replay "what if" scenarios
- Don't adjust stops based on intraday moves
- Honor your 10-bar Donchian / 2×N exit rules

**Why?** Backseat trading leads to:
- Cutting winners early (impatience)
- Widening stops on losers (hope)
- Overtrading (boredom)
- Revenge trading (frustration)

The system worked when you made the decision. Trust past-you.

### 8. Review Monthly Performance

**First Sunday of Each Month (30 minutes):**

1. Export trade history (if TF-Engine supports export, or manually)
2. Calculate:
   - Total R-multiple (sum of all trades in R terms)
   - Win rate (% of trades with +R)
   - Average winner (in R terms)
   - Average loser (in R terms)
   - Expectancy: (Win% × AvgWin) - (Loss% × AvgLoss)
3. Review:
   - Which sectors performed best?
   - Which setups had highest win rate?
   - Did you honor all exits (check journal notes)?
   - Were there any impulsive trades (breaches of discipline)?
4. Adjust (if needed):
   - Equity (account grew/shrank)
   - Sector focus (rotate based on market regime)

### 9. Stick to Your System

**The system is designed to be boring.**

- Same rules every day
- Same checklist every trade
- Same sizing method every position
- Same exit rules every time

**Boring = Profitable** in systematic trading.

If you find yourself wanting to:
- "Adjust" the rules for this one trade → STOP
- "Skip" the 2-minute timer because you're sure → STOP
- "Increase" size because you really like this setup → STOP
- "Hold" past your exit because you think it'll bounce → STOP

**These urges are the enemy.** The system exists to protect you from them.

### 10. Trade the Tide, Not the Splash

**Remember the core philosophy:**

- Trends last weeks to months
- Daily noise doesn't matter
- Whipsaws are part of the game
- Big winners make up for many small losers

**Don't:**
- Check positions 10 times a day
- Panic on red days
- Celebrate too much on green days
- Get emotionally attached to any position

**Do:**
- Check once daily (before/after market)
- Honor your exits mechanically
- Add to winners per your plan (every 0.5×N)
- Keep your process consistent

---

## Troubleshooting

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for comprehensive solutions.

### App Won't Start

**Symptoms:**
- Double-click tf-engine.exe → nothing happens
- Browser doesn't open automatically
- Error message: "Port already in use"

**Solutions:**

1. **Check if already running:**
   - Open Task Manager (Ctrl+Shift+Esc)
   - Processes tab → Look for "tf-engine.exe"
   - If running: Right-click → End Task → Try again

2. **Manually open browser:**
   - Start tf-engine.exe
   - Open browser manually
   - Navigate to: http://localhost:8080

3. **Port conflict:**
   - Another app is using port 8080
   - Run: `tf-engine.exe server --listen :8081`
   - Open browser to: http://localhost:8081

4. **Run as Administrator:**
   - Right-click tf-engine.exe
   - "Run as Administrator"

### FINVIZ Scan Returns 0 Candidates

**Symptoms:**
- Click "Run FINVIZ Scan" → Success but 0 results
- Or: Error message "Failed to parse FINVIZ HTML"

**Solutions:**

1. **Test FINVIZ URL in browser:**
   - Copy FINVIZ URL from Settings → FINVIZ Presets
   - Paste in browser
   - Should show candidate tickers
   - If 0 results in browser: Adjust screener filters

2. **FINVIZ HTML changed:**
   - FINVIZ occasionally updates their HTML structure
   - Requires TF-Engine update
   - Check for software updates
   - Temporary workaround: Manual candidate entry

3. **Network/firewall issue:**
   - FINVIZ request blocked by firewall
   - Add tf-engine.exe to firewall whitelist
   - Check antivirus logs

### Banner Stuck on RED

**Symptoms:**
- All 5 required items checked
- Banner still shows RED
- Can't proceed to sizing

**Solutions:**

1. **Verify all 5 checkboxes:**
   - Signal ✓
   - Risk/Size ✓
   - Liquidity ✓
   - Exits ✓
   - Behavior ✓
   - Scroll up/down to ensure all visible

2. **Missing trade data:**
   - Ticker must be filled
   - Entry must be > 0
   - N (ATR) must be > 0
   - Fill all required fields

3. **Refresh page:**
   - UI state might be stale
   - Refresh browser (F5)
   - Re-check items if needed

### Gates Check Fails But All Gates Look OK

**Symptoms:**
- Click "Run Final Gate Check"
- One or more gates fail
- But you think they should pass

**Solutions:**

1. **Check cooldowns (hidden):**
   - Go to Dashboard → Cooldowns section
   - Ticker or sector might be on cooldown
   - Not visible on Checklist screen
   - Wait for cooldown to expire or choose different ticker

2. **Verify 2-minute timer:**
   - Check timestamp of checklist save
   - Check current time
   - Must be ≥ 2:00 elapsed
   - System clock accurate?

3. **Heat caps changed:**
   - Check current heat on Dashboard
   - Might have changed since Heat Check screen
   - Re-run Heat Check
   - Another position might have been added

4. **Close and reopen Trade Entry:**
   - UI state might be stale
   - Navigate away and back
   - Re-run gate check

### Database is Locked

**Symptoms:**
- Error: "Database is locked"
- Cannot save settings/positions
- Operations fail with timeout

**Solutions:**

1. **Only one instance allowed:**
   - Task Manager → End all tf-engine.exe processes
   - Restart tf-engine.exe once
   - Should resolve

2. **Close database browser tools:**
   - If you have SQLite browser open on trading.db
   - Close it
   - Try again

3. **File permissions:**
   - Check `%APPDATA%\TF-Engine\` folder permissions
   - Should have write access
   - Run as Administrator if needed

4. **Restart Windows:**
   - Extreme case: Reboot
   - Clears any locks

### Calculations Look Wrong

**Symptoms:**
- Position sizing results don't match expectations
- Heat values seem off
- Stop distances incorrect

**Solutions:**

1. **Verify inputs:**
   - Entry, ATR, K must all be positive
   - Equity in settings must be correct
   - Risk % typically 0.50 - 1.00%

2. **Check math manually:**
   - Risk = Equity × Risk% = $100,000 × 0.0075 = $750
   - Stop Dist = K × ATR = 2 × $2.35 = $4.70
   - Shares = floor($750 / $4.70) = 159
   - Actual Risk = 159 × $4.70 = $747.30 ✓

3. **Review Van Tharp method:**
   - System uses standard Van Tharp formulas
   - Rounding down shares (floor function)
   - See [Position Sizing section](#position-sizing-5-minutes)

4. **Report bug:**
   - If calculations consistently wrong
   - Provide: Inputs, expected result, actual result
   - See [FAQ](#faq) for contact info

---

## FAQ

### General

**Q: Can I bypass the 5 gates?**

**A:** No. That's the entire point. The gates enforce discipline. If you want a flexible system, this is not the right tool for you.

---

**Q: Why 2 minutes? Can I reduce it?**

**A:** 2 minutes is the minimum to prevent impulsive decisions. Research shows that a brief pause reduces emotional decision-making. This is a hardcoded rule and cannot be changed in the UI.

---

**Q: What if I disagree with the sizing calculation?**

**A:** The Van Tharp method is mathematically rigorous and time-tested. If the sizing seems wrong, verify your inputs (entry, ATR, equity, risk%). The system does not allow manual overrides because overrides lead to impulsivity.

---

**Q: Can I use this for day trading?**

**A:** No. This is designed for swing/position trading with trend-following systems (holding days to weeks or months). The 2-minute timer, overnight reflection, and Donchian breakouts are all designed for longer timeframes.

---

**Q: Does it work for options?**

**A:** Yes. The system supports call/put strategies with:
- Delta-adjusted sizing (opt-delta-atr method)
- Contract-based sizing (opt-contracts method)
- Liquidity requirements: OI > 100, bid-ask spread < 10%
- DTE requirements: 60-90 days, roll/close at ~21 DTE

---

**Q: Can I run multiple instances?**

**A:** Not with the same database. Each instance needs its own database file. However, you probably shouldn't be running multiple independent strategies in one account anyway.

---

**Q: Why can't I edit the 5 required gates checklist items?**

**A:** Because they define the trend-following system. Changing them would make it a different system. If you want different rules, this tool is not for you.

---

**Q: What happens if my computer crashes mid-trade?**

**A:** TF-Engine saves data to database immediately on each action (Save Evaluation, Save Position Plan, etc.). If crash occurs:
- Checklist evaluation: Saved (2-min timer may need recheck)
- Position sizing: Saved
- Gates check: Re-run when you restart
- GO decision: Only saved when you click "SAVE GO DECISION"

Your broker order is separate and unaffected (since execution is manual).

---

### Position Management

**Q: How do I add pyramid units (add-ons)?**

**A:** TF-Engine calculates add-on levels for you (every 0.5×N). You execute manually in your broker:

1. Set price alerts at add-on levels (e.g., $181.68, $182.85, $184.03)
2. When alert triggers, verify:
   - Position still open (haven't exited)
   - Still above entry (winner)
   - Heat caps allow additional risk
3. Execute: BUY 159 shares AAPL @ $181.68
4. Update stop for all units to new trailing level

**Important:** Add-ons are optional. Only add to winners (price moved 0.5×N in your favor).

---

**Q: When do I exit a position?**

**A:** Two exit rules (whichever is closer):

1. **10-bar Donchian opposite breakout:**
   - For longs: Close < 10-bar low
   - Check daily on TradingView

2. **2×N stop:**
   - Set stop-loss order in broker
   - For longs: Stop = Entry - 2×N

**Exit immediately if either is hit.** No exceptions.

---

**Q: Can I hold through earnings?**

**A:** Not recommended, especially for options. Earnings create gap risk that violates your 2×N stop. The optional "Earnings OK" checklist item encourages you to avoid earnings.

For stocks, if you choose to hold:
- Understand gap risk
- Might get stopped out with >2N loss
- Accept this as part of the system

For options:
- Strongly discouraged
- Volatility crush post-earnings hurts long options

---

**Q: What if the 10-bar stop is way below my 2×N stop?**

**A:** Use the **closer** of the two. If 10-bar low is at $178 and 2×N stop is at $175.80, your stop is $178. This trails your stop higher as the trade moves in your favor.

---

**Q: How do I handle dividends or splits?**

**A:** TF-Engine doesn't automatically adjust for corporate actions. You must manually:

**Splits:**
- 2-for-1 split: Double your shares, halve entry price and stops
- Example: 159 shares @ $180 → 318 shares @ $90
- ATR (N) also halves: $2.35 → $1.18

**Dividends:**
- No adjustment needed (cash paid out)
- Minor: Continue as normal
- Special dividend: Might affect Donchian levels on chart

**Rare edge case:** If confused, close position and re-enter fresh after split settles.

---

### Technical Issues

**Q: Can I use TF-Engine offline?**

**A:** Partially. The app runs locally, but:
- FINVIZ scanning requires internet
- TradingView integration requires internet
- Checklist, sizing, heat check, gates work offline

**Offline workflow:**
- Import candidates while online
- Analyze charts while online
- Complete checklist/sizing/gates offline if needed

---

**Q: How do I export my trade history?**

**A:** Currently, trade history is in SQLite database. To export:

1. Install SQLite browser: https://sqlitebrowser.org/
2. Open `%APPDATA%\TF-Engine\trading.db`
3. Navigate to "decisions" table
4. File → Export → CSV
5. Open in Excel/Google Sheets

Future versions may include built-in export.

---

**Q: Can I sync across multiple computers?**

**A:** Not natively. Workaround:

1. Store `trading.db` in Dropbox/Google Drive/OneDrive
2. Create symlink from `%APPDATA%\TF-Engine\trading.db` to cloud folder
3. Only run TF-Engine on one computer at a time (database locking)

**Not recommended for reliability.** Best practice: Use one primary trading computer.

---

**Q: What data does TF-Engine collect?**

**A:** None. TF-Engine:
- Runs entirely locally on your computer
- No telemetry or analytics
- No data sent to remote servers
- FINVIZ requests go directly from your computer to FINVIZ
- Your trading data stays on your machine

---

### System Design

**Q: Why is the system so rigid?**

**A:** Because flexibility = opportunity for impulsivity. The rigidity IS the feature. Every "inconvenience" is an intentional friction point designed to protect you from yourself.

---

**Q: Can I suggest features?**

**A:** Yes, but understand the design philosophy:
- Features that increase discipline: Considered
- Features that add flexibility: Likely rejected
- Features that add complexity: Scrutinized heavily

Read `docs/anti-impulsivity.md` and `docs/project/WHY.md` first.

---

**Q: Why no auto-trading / broker integration?**

**A:** Manual execution is intentional:
- Final human verification before capital is committed
- Prevents catastrophic auto-trading errors
- Keeps you engaged in the process
- You remain responsible for your trades

Auto-trading is seductive but dangerous.

---

**Q: What if I want to trade a different system (not Donchian)?**

**A:** TF-Engine is purpose-built for Ed Seykota style trend-following. For other systems:
- You'd need different signal checks (not 55-bar Donchian)
- Different sizing (not Van Tharp with 2×N stops)
- Different exits (not 10-bar opposite)

This tool is not designed to be a general-purpose trading platform. It enforces one specific system.

---

### Support

**Q: How do I get help?**

**A:**
1. Read this User Guide (you're doing it!)
2. Check [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
3. Check [FAQ.md](FAQ.md) (extended version)
4. Review documentation in `docs/` folder
5. Contact: [your-support-email or GitHub issues]

---

**Q: Can I contribute to the project?**

**A:** If open source: Yes! See CONTRIBUTING.md (if exists).
If proprietary: Feature requests and bug reports welcome.

---

**Q: What's the license?**

**A:** [Specify license - MIT, GPL, proprietary, etc.]

---

## Conclusion

**Remember the core principles:**

1. This system is a tool for **discipline, not flexibility**
2. Every "inconvenience" is an **intentional friction point**
3. The 5 gates exist to **protect you from yourself**
4. **Trade the tide, not the splash**
5. **Boring = Profitable** in systematic trading

**Success with TF-Engine requires:**
- Following the rules **every single time**
- Honoring your exits **mechanically**
- Accepting small losses **gracefully**
- Letting big winners run **patiently**
- Reviewing performance **objectively**

**The system works if you work the system.**

---

**Good luck, and trade systematically!**

---

## Support & Resources

**Documentation:**
- [Quick Start Guide](QUICK_START.md) - Get started in 10 minutes
- [Installation Guide](INSTALLATION_GUIDE.md) - Detailed install instructions
- [TradingView Setup](TRADINGVIEW_SETUP.md) - Pine Script integration
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues and solutions
- [FAQ](FAQ.md) - Extended frequently asked questions

**Technical:**
- GitHub: [your-repo-url] (if open source)
- Issues: [your-issues-url]
- Email: [your-support-email]

**Version:** 1.0.0
**Last Updated:** 2025-10-29
**System:** TF-Engine - Trend Following Engine
**Philosophy:** Trade the tide, not the splash.

---

**End of User Guide**
