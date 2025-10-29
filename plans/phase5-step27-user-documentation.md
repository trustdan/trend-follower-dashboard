# Phase 5 - Step 27: User Documentation

**TF = Trend Following** - Systematic trading discipline enforcement system

**Phase:** 5 - Testing & Packaging
**Step:** 27 of 28
**Duration:** 2 days
**Dependencies:** Step 26 complete (Windows installer ready)

---

## Objectives

Create comprehensive, user-facing documentation that enables traders to install, configure, and use TF-Engine effectively. Documentation should cover installation, first-time setup, daily workflow, understanding the 5 gates, troubleshooting, and FAQs. Include screenshots for major screens. Make documentation accessible to non-technical traders.

**Purpose:** Ensure users can successfully adopt and use the system without requiring developer support.

---

## Success Criteria

- [ ] User Guide created (comprehensive, with screenshots)
- [ ] Quick Start guide created (one-page, get started in 10 minutes)
- [ ] FAQ document created (answers common questions)
- [ ] TROUBLESHOOTING.md created (common issues and solutions)
- [ ] In-app help text/tooltips added
- [ ] Screenshots captured for all major screens
- [ ] Documentation reviewed by non-technical user (if possible)
- [ ] All docs use clear, accessible language (no jargon)

---

## Prerequisites

**Completed:**
- All features implemented
- Windows installer created
- Application stable and tested

**Tools:**
- Screenshot tool (Windows Snipping Tool, ShareX, etc.)
- Image editor (optional, for annotations)
- Markdown editor

---

## Implementation Plan

### Part 1: User Guide (1 day)

#### Task 1.1: Capture Screenshots (2 hours)

**Required Screenshots:**

1. **Dashboard**
   - Empty state (first run)
   - With open positions
   - Portfolio heat gauge

2. **FINVIZ Scanner**
   - Scan button
   - Results table
   - Import candidates

3. **Checklist**
   - RED banner (gates unchecked)
   - YELLOW banner (low quality score)
   - GREEN banner (ready to trade)
   - Required gates section
   - Quality items section
   - 2-minute timer

4. **Position Sizing**
   - Input form
   - Calculation results
   - Add-on schedule display

5. **Heat Check**
   - Current heat display
   - Check result (within caps)
   - Check result (exceeds cap - RED warning)

6. **Trade Entry**
   - Trade summary
   - Gate check results (all pass)
   - GO/NO-GO buttons

7. **Calendar**
   - 10-week grid view
   - Positions in cells
   - Color coding legend

8. **Settings**
   - Account settings form

9. **Theme Toggle**
   - Day mode
   - Night mode

Save in: `docs/screenshots/`

#### Task 1.2: Write User Guide (4-6 hours)

**File:** `docs/USER_GUIDE.md`

```markdown
# TF-Engine User Guide

**Version:** 1.0.0
**Last Updated:** 2025-10-30

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
- **Entry:** 55-bar Donchian breakouts
- **Position Sizing:** Van Tharp ATR-based method
- **Pyramiding:** Add to winners every 0.5×N
- **Exits:** 10-bar opposite Donchian OR 2×N stop

### What It's NOT

- **NOT** a signal generator (you verify signals manually in TradingView)
- **NOT** an auto-trader (you execute trades manually in your broker)
- **NOT** a backtester (it's for live trading discipline)

### Core Philosophy

> **Trade the tide, not the splash.**

The system enforces discipline through **5 hard gates** that cannot be bypassed. If any gate fails, you cannot save a GO decision. This prevents impulsive trades driven by fear, greed, or impatience.

---

## Installation

See [INSTALLATION_GUIDE.md](INSTALLATION_GUIDE.md) for detailed installation instructions.

**Quick Summary:**
1. Download `TF-Engine-Setup-v1.0.0.msi`
2. Run installer
3. Follow wizard
4. Launch from Desktop shortcut
5. Browser opens to localhost:8080

---

## First-Time Setup

### Step 1: Configure Account Settings

![Settings Screenshot](screenshots/settings.png)

1. Navigate to **Settings** (gear icon in header)
2. Enter your account details:
   - **Equity:** Your trading account size (e.g., $100,000)
   - **Risk % per unit:** Typically 0.75% - 1.0% (Ed Seykota uses 1%)
   - **Portfolio heat cap:** 4.0% (total risk across all positions)
   - **Bucket heat cap:** 1.5% (max risk per sector)
   - **Max units:** 4 (initial + 3 add-ons)
3. Click **Save Settings**

**Why these numbers?**
- **0.75% risk per unit:** With 4 units max, you risk 3% per position
- **4% portfolio cap:** Total risk across all positions never exceeds 4% of equity
- **1.5% bucket cap:** Limits concentration in one sector (e.g., Tech/Comm)

### Step 2: Set Up TradingView

You'll use TradingView to verify breakout signals.

1. Create TradingView account (free or paid)
2. Install the **Ed-Seykota.pine** script:
   - Open Pine Editor in TradingView
   - Copy contents from `reference/Ed-Seykota.pine`
   - Click "Add to Chart"
3. Save script to favorites for easy access

**The script displays:**
- 55-bar Donchian channels (entry signals)
- 10-bar Donchian channels (exit signals)
- Current N (ATR value)
- Initial stop levels

See [TRADINGVIEW_SETUP.md](TRADINGVIEW_SETUP.md) for detailed guide.

### Step 3: Set Up FINVIZ Screener (Optional)

If you want to use the automated daily scan:

1. Create FINVIZ account (free or Elite)
2. Build your screener for trend-following candidates
3. Copy the screener URL
4. In TF-Engine Settings, add FINVIZ preset with your URL

**Suggested filters for long breakouts:**
- Technical: Price > SMA200
- Technical: RSI > 55
- Fundamental: Avg Volume > 1M
- Fundamental: Market Cap > $500M

---

## Daily Trading Workflow

### Morning Routine (30 minutes)

![Dashboard Screenshot](screenshots/dashboard.png)

#### 1. Check Dashboard

Start your day at the **Dashboard**:
- Review open positions (entry, current stop, risk, days held)
- Check portfolio heat (should be well below 4% cap)
- Note any cooldowns (sectors or tickers you cannot trade)

#### 2. Scan for Candidates

![FINVIZ Scanner Screenshot](screenshots/scanner.png)

Navigate to **Scanner** → Click **"Run Daily FINVIZ Scan"**

The system fetches tickers from your FINVIZ screener.

Review the candidates table:
- Ticker symbols
- Sector buckets
- Last close price
- Average volume

**Select 10-15 candidates** to analyze. Click **"Import Selected"**.

**Tip:** Focus on sectors you're underweight (check Dashboard).

#### 3. Verify Signals in TradingView

For each candidate:

![TradingView Link](screenshots/tradingview-link.png)

1. Click **"Open in TradingView"** next to ticker
2. Chart opens in new tab with your script
3. Verify the breakout:
   - **For longs:** Close > 55-bar high (blue line)
   - **For shorts:** Close < 55-bar low (red line)
4. Note the **N value** (ATR) from script
5. Note your **entry price** (current close or limit)
6. Note the **sector** (Tech/Comm, Energy, etc.)

**Return to TF-Engine** when ready to complete checklist.

---

### Checklist Evaluation (10 minutes)

![Checklist Screenshot](screenshots/checklist.png)

Navigate to **Checklist**:

#### Step 1: Enter Trade Data

- **Ticker:** AAPL
- **Entry Price:** 180.50
- **N (ATR):** 2.35 (from TradingView)
- **Sector:** Tech/Comm
- **Structure:** Stock (or Call, Put, etc. for options)

#### Step 2: Check All 5 Required Gates

These are **mandatory**. The banner is RED until all are checked.

1. **✓ Signal:** 55-bar Donchian breakout confirmed
   - You verified this in TradingView
2. **✓ Risk/Size:** Will use 2×N stop, add every 0.5×N, max 4 units
   - System calculates this automatically
3. **✓ Liquidity:** Stock avg volume > 1M (or options OI > 100)
   - Check on TradingView or Yahoo Finance
4. **✓ Exits:** Will exit on 10-bar opposite Donchian OR 2×N stop
   - Whichever is closer
5. **✓ Behavior:** Not on cooldown, heat OK, will honor 2-min timer
   - System validates this in final gates check

**Banner turns YELLOW** when all required are checked.

#### Step 3: Improve Quality Score (Optional)

These are **optional** but recommended. Each adds 1 point.

6. **✓ Regime OK:** SPY > 200 SMA (favorable market for longs)
   - Check SPY chart on TradingView
7. **✓ No Chase:** Entry within 2N of 20-EMA (not overextended)
   - Visual check on TradingView
8. **✓ Earnings OK:** No earnings within next 2 weeks (for long options)
   - Check earnings calendar (Yahoo Finance, TradingView)
9. **✓ Journal Note:** Why this trade now?
   - Write 1-2 sentences: "Clean breakout, strong volume, sector underweight"

**Quality Score:** 4 / 4
**Threshold:** 3.0 (configurable in settings)

**Banner turns GREEN** when score ≥ threshold.

#### Step 4: Save Evaluation

![Green Banner](screenshots/banner-green.png)

Click **"Save Evaluation"**

**The 2-minute cool-off timer starts:** 2:00... 1:59... 1:58...

![Timer](screenshots/timer.png)

This is the **impulse brake**. You cannot save a GO decision until the timer reaches 0:00.

**Use this time to:**
- Double-check your analysis
- Review the trade plan
- Calculate position sizing

---

### Position Sizing (5 minutes)

![Position Sizing Screenshot](screenshots/position-sizing.png)

Navigate to **Position Sizing**:

The form is pre-filled with your checklist data.

**Method:** Stock (default)
**K multiple:** 2.0 (stop distance)
**Max units:** 4 (from settings)

Click **"Calculate Position Size"**

**Results Display:**

```
Shares per unit: 159
Risk per unit: $747.30
Initial stop: $175.80

Add-on schedule:
Unit 1: 159 shares @ $180.50 (now)
Unit 2: 159 shares @ $181.68 (+0.5N)
Unit 3: 159 shares @ $182.85 (+1.0N)
Unit 4: 159 shares @ $184.03 (+1.5N)

Exit plan: 10-bar Donchian OR $175.80 stop, whichever is closer
```

**If concentration warning appears:**
- Position > 25% of equity → Consider reducing

Click **"Save Position Plan"**

---

### Heat Check (2 minutes)

![Heat Check Screenshot](screenshots/heat-check.png)

Navigate to **Heat Check**:

The system displays:
- **Current portfolio heat:** $2,890 / $4,000 cap (72%)
- **Current Tech/Comm heat:** $1,125 / $1,500 cap (75%)

Click **"Check Heat for This Trade"**

**Proposed trade:** $747 risk in Tech/Comm

**Result:**

```
✓ Portfolio heat: $3,637 / $4,000 (91%) - WITHIN CAP
✓ Bucket heat: $1,872 / $1,500 (125%) - EXCEEDS CAP by $372
```

![Heat Warning](screenshots/heat-warning.png)

**If cap exceeded:**
1. Reduce position size: Click "Calculate Max Shares for Heat Cap"
   - System suggests: 79 shares (instead of 159)
2. Close an existing Tech/Comm position first
3. Choose a different ticker from another sector

**If within caps:**
- Proceed to Trade Entry

---

### Final Trade Entry (2 minutes)

![Trade Entry Screenshot](screenshots/trade-entry.png)

Navigate to **Trade Entry**:

**Review Trade Summary:**
- Ticker: AAPL
- Direction: LONG
- Entry: $180.50
- Shares: 79 per unit (adjusted for heat cap)
- Initial stop: $175.80
- Risk: $371 per unit (max $1,484 if 4 units)
- Sector: Tech/Comm
- Quality score: 4/4

Click **"Run Final Gate Check"**

![Gate Check Results](screenshots/gates-all-pass.png)

**The 5 Gates Validate:**

```
Gate 1: Banner Status → GREEN ✓
Gate 2: Impulse Brake → 2:15 elapsed ✓
Gate 3: Cooldown Status → Not on cooldown ✓
Gate 4: Heat Caps → Within caps ✓
Gate 5: Sizing Completed → Plan saved ✓

Result: ALL GATES PASS ✓
```

**"SAVE GO DECISION" button is now ENABLED (green).**

Click **"SAVE GO DECISION"**

**Success notification:** "✓ GO decision saved for AAPL"

---

### Execute in Broker (Manual)

**TF-Engine does NOT execute trades.** You must manually enter the trade in your broker.

**Why manual execution?**
- Final human verification
- Prevents catastrophic auto-trading errors
- Keeps you engaged in the process

**Execute:**
1. Open your broker (Interactive Brokers, TD Ameritrade, etc.)
2. Enter order: BUY 79 shares AAPL @ $180.50 (limit or market)
3. Set initial stop: $175.80
4. Record entry in TF-Engine (optional: add position to Dashboard)

**Set alerts for add-ons:**
- Add2 @ $181.68
- Add3 @ $182.85
- Add4 @ $184.03

---

## Understanding the Banner

The banner is the **heart of the discipline enforcement system**.

### RED Banner: DO NOT TRADE

![Red Banner](screenshots/banner-red.png)

**Meaning:** One or more required gates are unchecked.

**What to do:**
- Complete all 5 required gates
- Do NOT proceed to trade entry

### YELLOW Banner: CAUTION

![Yellow Banner](screenshots/banner-yellow.png)

**Meaning:** All required gates pass, but quality score is below threshold (< 3.0).

**What to do:**
- Review quality items
- Add more quality checks to improve score
- Or accept lower quality trade (not recommended)

### GREEN Banner: OK TO TRADE

![Green Banner](screenshots/banner-green.png)

**Meaning:** All required gates pass, quality score ≥ threshold.

**What to do:**
- Proceed to position sizing, heat check, and trade entry
- You may save a GO decision (after 2-min timer and gates check)

**Note:** GREEN banner is **necessary** but not **sufficient**. The 5 gates check at Trade Entry is the final validation.

---

## The 5 Gates Explained

### Gate 1: Banner Status

**Rule:** Banner must be GREEN.

**Why:** Ensures all required checklist items and quality score are met.

**Failure:** If banner is RED or YELLOW, Gate 1 fails.

### Gate 2: Impulse Brake (2-Minute Cool-Off)

**Rule:** At least 2 minutes must elapse since checklist evaluation.

**Why:** Forces you to slow down and think. Prevents impulsive decisions.

**Failure:** If you try to save GO decision before 2:00, Gate 2 fails.

**Tip:** Use the 2 minutes to review analysis and calculate sizing.

### Gate 3: Cooldown Status

**Rule:** Ticker and sector must not be on cooldown.

**Why:** After a losing trade, take a break from that ticker or sector to avoid revenge trading.

**Cooldown duration:** Typically 1-4 weeks (configurable).

**Failure:** If ticker or bucket is on cooldown, Gate 3 fails.

### Gate 4: Heat Caps

**Rule:** Proposed trade must not exceed portfolio (4%) or bucket (1.5%) heat caps.

**Why:** Limits overall risk and prevents concentration in one sector.

**Portfolio cap:** Sum of risk across ALL positions ≤ 4% of equity
**Bucket cap:** Sum of risk in ONE sector ≤ 1.5% of equity

**Failure:** If adding this trade would exceed either cap, Gate 4 fails.

**Solution:** Reduce position size, close existing position, or choose different ticker.

### Gate 5: Sizing Completed

**Rule:** Position sizing must be calculated and saved.

**Why:** Cannot enter a trade without knowing how many shares and what the risk is.

**Failure:** If you skip position sizing, Gate 5 fails.

---

## Screen Reference

### Dashboard

**Purpose:** Overview of portfolio, positions, and candidates.

**Key Elements:**
- Portfolio summary (equity, heat, capacity)
- Open positions table (ticker, entry, stop, risk, days held)
- Today's candidates count
- Cooldowns (if any)
- Quick actions (Navigate to Scanner, Checklist)

### Scanner

**Purpose:** Daily FINVIZ scan to find trend-following candidates.

**Key Elements:**
- "Run Daily FINVIZ Scan" button
- Results table (ticker, sector, last close, volume)
- Checkbox selection for import
- Sector distribution summary
- "Import Selected" button

### Checklist

**Purpose:** Evaluate trade setup and calculate banner state.

**Key Elements:**
- Large banner (RED/YELLOW/GREEN)
- Trade data inputs (ticker, entry, ATR, sector)
- 5 required gates (checkboxes)
- 4 quality items (checkboxes + journal textarea)
- Quality score display
- "Save Evaluation" button
- 2-minute timer (after save)

### Position Sizing

**Purpose:** Calculate shares/contracts using Van Tharp ATR method.

**Key Elements:**
- Method dropdown (stock, opt-delta-atr, opt-contracts)
- Pre-filled trade data
- "Calculate Position Size" button
- Results display (shares, risk, stops, add-ons)
- Concentration warning (if > 25% equity)
- "Save Position Plan" button

### Heat Check

**Purpose:** Verify proposed trade doesn't exceed risk caps.

**Key Elements:**
- Current portfolio heat gauge
- Current bucket heat table
- "Check Heat for This Trade" button
- Result display (within caps or exceeded)
- Suggestions to resolve cap issues
- "Calculate Max Shares" helper button

### Trade Entry

**Purpose:** Final validation and GO/NO-GO decision.

**Key Elements:**
- Trade summary (all data from previous steps)
- "Run Final Gate Check" button
- Gate results display (5 gates, pass/fail)
- "SAVE GO DECISION" button (green, enabled only if all gates pass)
- "SAVE NO-GO DECISION" button (red, always enabled)

### Calendar

**Purpose:** Visualize sector diversification over time.

**Key Elements:**
- 10-week grid (2 back, 8 forward)
- Rows: sector buckets
- Columns: weeks (Monday-Sunday)
- Cells: tickers active in that sector/week
- Color coding: heat levels (green/yellow/red)
- Tooltips: position details
- Current week highlighted

### Settings

**Purpose:** Configure account settings and system parameters.

**Key Elements:**
- Equity (account size)
- Risk % per unit (default 0.75%)
- Portfolio heat cap (default 4.0%)
- Bucket heat cap (default 1.5%)
- Max units (default 4)
- TradingView URL template (optional)

---

## TradingView Integration

See [TRADINGVIEW_SETUP.md](TRADINGVIEW_SETUP.md) for detailed guide.

**Quick Tips:**
- Click "Open in TradingView" next to any ticker
- Verify breakout signal using Ed-Seykota.pine script
- Note the N (ATR) value displayed on chart
- Return to TF-Engine and fill checklist

**Custom URL template:**
If you have a saved TradingView chart layout, customize the URL in Settings to auto-load your layout.

---

## Theme Customization

**Day Mode (Light):**
- Bright, clean interface
- High contrast for readability

**Night Mode (Dark):**
- Easy on the eyes for extended use
- Gradient banners still vibrant

**Toggle:** Click sun/moon icon in header (top-right)

**Persistence:** Theme preference saves automatically (localStorage)

---

## Tips & Best Practices

### 1. Start with Small Position Sizes

If new to this system, reduce your risk % per unit:
- Beginner: 0.50% per unit
- Intermediate: 0.75% per unit
- Experienced: 1.00% per unit

### 2. Respect the 2-Minute Timer

Don't try to rush. Use the time productively:
- Double-check your TradingView chart
- Review sector diversification in Calendar
- Calculate potential profit at 1R, 2R, 3R

### 3. Journal Every Trade

Even if quality score doesn't require it, always add a journal note. Future you will thank present you.

**Good journal notes:**
- "Clean breakout above 6-month consolidation, volume 2x avg, sector underweight"
- "Pullback to 20-EMA in strong uptrend, SPY confirmed above 200 SMA"

**Bad journal notes:**
- "Looks good"
- "Strong"

### 4. Use the Calendar Weekly

Every Sunday evening:
- Review the Calendar view
- Identify gaps (weeks with no trades)
- Identify crowded sectors (too many positions in one bucket)
- Plan rebalancing for the week ahead

### 5. Save NO-GO Decisions

When a trade fails the gates, save it as NO-GO with a reason. This creates a valuable record of what you DON'T trade, which is just as important as what you DO trade.

### 6. Backup Your Database

**Location:** `C:\Users\[YourUsername]\AppData\Roaming\TF-Engine\trading.db`

**Frequency:** Weekly or after significant activity

**Method:** Copy `trading.db` to Dropbox, Google Drive, or external drive

---

## Troubleshooting

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for detailed solutions.

**Common Issues:**

### App won't start
- Check Task Manager → tf-engine.exe running?
- Manually open browser to localhost:8080
- Restart as Administrator

### FINVIZ scan returns 0 candidates
- Check FINVIZ URL is correct
- Test URL in browser first
- FINVIZ may have changed HTML structure (requires update)

### Banner stuck on RED
- Ensure all 5 required gates are checked
- Refresh page if UI is stale

### Gates check fails but all gates look OK
- Check cooldowns (might not be visible on Checklist)
- Verify 2-minute timer has elapsed (check timestamp)
- Close and reopen Trade Entry screen

### Database is locked
- Only one instance of TF-Engine can run
- Close SQLite browser tools
- Restart TF-Engine

---

## FAQ

See full [FAQ.md](FAQ.md) for more questions.

### Q: Can I bypass the 5 gates?

**A:** No. That's the entire point. The gates enforce discipline.

### Q: Why 2 minutes? Can I reduce it?

**A:** 2 minutes is the minimum to prevent impulsive decisions. This is a hardcoded rule and cannot be changed in the UI.

### Q: What if I disagree with the sizing calculation?

**A:** The Van Tharp method is mathematically rigorous. If the sizing seems wrong, verify your inputs (entry, ATR, equity, risk%). The system does not allow manual overrides.

### Q: Can I use this for day trading?

**A:** No. This is designed for swing/position trading with trend-following systems (holding days to months). The 2-minute timer and overnight reflection are intentional.

### Q: Does it work for options?

**A:** Yes. The system supports call/put spreads with delta-adjusted sizing and contract-based calculations. Options must meet liquidity requirements (OI > 100, bid-ask spread < 10%).

### Q: Can I run multiple instances?

**A:** Not with the same database. Each instance needs its own database file. However, you probably shouldn't be running multiple strategies in one account anyway.

---

## Conclusion

**Remember:**
- This system is a tool for discipline, not flexibility
- Every "inconvenience" is an intentional friction point
- The 5 gates exist to protect you from yourself
- Trade the tide, not the splash

**Good luck, and trade systematically!**

---

**Support:**
- Documentation: https://[your-docs-site]
- Issues: https://github.com/[your-repo]/issues
- Email: [your-support-email]

**Version:** 1.0.0
**Last Updated:** 2025-10-30
```

---

### Part 2: Quick Start Guide (2 hours)

**File:** `docs/QUICK_START.md`

```markdown
# TF-Engine Quick Start

**Get started in 10 minutes**

---

## 1. Install (2 minutes)

1. Download `TF-Engine-Setup-v1.0.0.msi`
2. Run installer → Follow wizard
3. Launch from Desktop shortcut

---

## 2. Configure Settings (2 minutes)

1. Navigate to **Settings**
2. Enter:
   - Equity: $100,000 (your account size)
   - Risk %: 0.75%
3. Click **Save**

---

## 3. Set Up TradingView (3 minutes)

1. Open TradingView
2. Pine Editor → Paste `Ed-Seykota.pine` script
3. Add to Chart

---

## 4. Run Daily Workflow (3 minutes)

### Morning:
1. **Scanner** → Run FINVIZ Scan → Import 10 candidates
2. Click "Open in TradingView" → Verify breakout
3. **Checklist** → Fill form, check 5 gates → GREEN banner
4. **Position Sizing** → Calculate shares
5. **Heat Check** → Verify within caps
6. **Trade Entry** → Run gates → Save GO decision
7. Execute in broker

---

## Next Steps

Read full [USER_GUIDE.md](USER_GUIDE.md) for detailed explanations.

**Trade systematically!**
```

---

### Part 3: FAQ & Troubleshooting (2 hours)

Already covered in User Guide sections above. Can be extracted as separate files:

- `docs/FAQ.md`
- `docs/TROUBLESHOOTING.md`

---

### Part 4: In-App Help Text (2 hours)

Add tooltips throughout the app:

**File:** `ui/src/lib/components/Tooltip.svelte` (already created in Step 22)

**Usage in Checklist:**

```svelte
<Tooltip text="This is the entry price for your position. Use current close or your planned limit order price.">
    <label for="entry">Entry Price</label>
</Tooltip>

<Tooltip text="N is the ATR (Average True Range). Get this value from your TradingView chart using the Ed-Seykota.pine script.">
    <label for="atr">N (ATR)</label>
</Tooltip>
```

Add tooltips to all complex fields across all screens.

---

## Testing Checklist

### Documentation Completeness
- [ ] User Guide covers all major features
- [ ] Screenshots included for all screens
- [ ] Quick Start guide is concise (one page)
- [ ] FAQ answers common questions
- [ ] Troubleshooting covers known issues
- [ ] In-app tooltips added

### Quality Check
- [ ] Language is clear and accessible (no jargon)
- [ ] Steps are numbered and easy to follow
- [ ] Examples are realistic
- [ ] No typos or grammatical errors
- [ ] Links work (internal references)

### User Testing (if possible)
- [ ] Non-technical user can follow Quick Start
- [ ] User can complete full workflow using User Guide
- [ ] User finds FAQ helpful
- [ ] User can resolve issues using Troubleshooting

---

## Deliverables

- [ ] `docs/USER_GUIDE.md` (comprehensive)
- [ ] `docs/QUICK_START.md` (one-page)
- [ ] `docs/FAQ.md`
- [ ] `docs/TROUBLESHOOTING.md`
- [ ] `docs/screenshots/` (folder with all images)
- [ ] In-app tooltips implemented

---

## Documentation Requirements

- [ ] All documentation uses Markdown format
- [ ] Screenshots have descriptive filenames
- [ ] Internal links tested (no broken references)
- [ ] Update `README.md` with links to all docs
- [ ] Update `docs/PROGRESS.md` with documentation status

---

## Next Steps

After completing Step 27:
1. Review all documentation for accuracy
2. Get feedback from beta user (if available)
3. Proceed to **Step 28: Final Validation & Release**

---

**Estimated Completion Time:** 2 days
**Phase 5 Progress:** 4 of 5 steps complete
**Overall Progress:** 27 of 28 steps complete (96%)

---

**End of Step 27**
