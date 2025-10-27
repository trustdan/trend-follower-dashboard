# Trading System User Guide - Complete Walkthrough

**Welcome!** This guide will walk you through everything you need to know to use this trading system, from first setup to logging your first trade.

---

## Table of Contents

1. [First Time Setup (5 minutes)](#first-time-setup)
2. [Understanding the Worksheets](#understanding-the-worksheets)
3. [Understanding the Settings](#understanding-the-settings)
4. [Daily Morning Routine (10 minutes)](#daily-morning-routine)
5. [Evaluating a Single Trade (5 minutes)](#evaluating-a-single-trade)
6. [Understanding Each Field](#understanding-each-field)
7. [The 6-Item Checklist Explained](#the-6-item-checklist-explained)
8. [Position Sizing Methods](#position-sizing-methods)
9. [Heat Management (Risk Control)](#heat-management)
10. [Troubleshooting](#troubleshooting)

---

## First Time Setup

### Step 1: Build the Workbook

1. Open your **Command Prompt** (Windows key + R, type `cmd`, press Enter)
2. Navigate to the folder: `cd C:\path\to\excel-trading-workflow`
3. Type: `BUILD.bat`
4. Press **Enter**
5. Wait 30 seconds while it builds the workbook

**What just happened?** The script created a new Excel file called `TrendFollowing_TradeEntry.xlsm` and imported all the VBA code.

### Step 2: Open the Workbook

1. Double-click **TrendFollowing_TradeEntry.xlsm**
2. If you see a yellow **"Security Warning"** banner at the top, click **"Enable Content"**
3. Wait a few seconds - you'll see a popup: **"Welcome to the Trading System! Running initial setup now..."**
4. Click **OK**
5. Wait 10-30 seconds while setup runs
6. You'll see: **"Setup Complete!"** with timing info
7. Click **OK**

**What just happened?** The workbook automatically created 8 worksheets, 5 data tables, loaded default presets, and built the TradeEntry user interface.

### Step 3: Verify Checkboxes (Important!)

1. Click the **"TradeEntry"** tab at the bottom of Excel
2. Look at rows 21-26 on the left side
3. You should see **6 small checkboxes** next to these labels:
   - `[ ] FromPreset`
   - `[ ] TrendPass`
   - `[ ] LiquidityPass`
   - `[ ] TVConfirm`
   - `[ ] EarningsOK`
   - `[ ] JournalOK`

**If you DON'T see checkboxes:**
1. Click the **Developer** tab (top of Excel)
   - *Don't see Developer tab? File → Options → Customize Ribbon → Check "Developer"*
2. Click **Insert** → Under "Form Controls", click the **Check Box** icon (☐)
3. Draw a small checkbox next to the first label (FromPreset)
4. **Right-click** the checkbox → **Format Control**
5. In "Cell link", type: **C20**
6. Click **OK**
7. Repeat for the other 5 checkboxes, linking them to **C21, C22, C23, C24, C25**
8. Right-click each checkbox, click **Edit Text**, and **delete the text** (labels are already there)

### You're Done!

The system is now ready to use. Let's understand what we've built.

---

## Understanding the Worksheets

When you open the workbook, you'll see **8 tabs** at the bottom:

### 1. **Setup** (Start Here!)
- Your control panel
- Shows setup status
- Has utility buttons:
  - **Rebuild TradeEntry UI** - Fixes UI if something breaks
  - **Test Python Integration** - Checks if auto-scraping works
  - **Clear Old Candidates** - Removes old ticker imports
- Contains instructions and settings summary

### 2. **TradeEntry** (Main Workspace)
- This is where you spend 90% of your time
- Where you evaluate trades
- Where you check the GO/NO-GO criteria
- Where you calculate position sizes
- Where you save trade decisions

### 3. **Presets** (FINVIZ Scan Definitions)
- Stores your FINVIZ screener queries
- **5 default presets included:**
  - `TF_BREAKOUT_LONG` - Stocks breaking out to new highs
  - `TF_PULLBACK_LONG` - Stocks pulling back in uptrend
  - `TF_MOMENTUM_LONG` - Strong momentum stocks
  - `TF_BREAKDOWN_SHORT` - Stocks breaking down (for shorts)
  - `TF_CUSTOM` - Your custom criteria
- Each preset has a **QueryString** (FINVIZ URL parameters)

### 4. **Buckets** (Correlation Groups)
- Groups stocks by correlation to prevent over-concentration
- **6 default buckets:**
  - `Tech/Comm` - Technology, Communications
  - `Healthcare` - Healthcare, Biotech
  - `Financials` - Banks, Insurance, REITs
  - `Consumer` - Retail, Consumer Discretionary
  - `Industrials` - Manufacturing, Transportation
  - `Energy/Materials` - Oil, Gas, Mining
- Each bucket tracks:
  - **StopoutsInWindow** - How many losses in last 20 days
  - **CooldownUntil** - If paused due to 2+ stopouts

### 5. **Candidates** (Daily Import List)
- Stores tickers you import each morning
- Columns:
  - **Ticker** - Stock symbol (AAPL, MSFT, etc.)
  - **PresetName** - Which scan it came from
  - **ImportDate** - When you imported it
  - **Notes** - Optional notes
- Auto-clears candidates older than 7 days

### 6. **Decisions** (Trade Log)
- Every trade decision is logged here
- Columns include:
  - **Ticker, Entry, Stop, Target**
  - **Position size (shares or contracts)**
  - **R ($)** - Dollar risk
  - **All 6 checklist items** (TRUE/FALSE)
  - **Timestamp**
- This is your audit trail

### 7. **Positions** (Open Trades)
- Manually track open positions here
- Used for heat calculations
- Columns:
  - **Ticker, Entry, Current, Stop**
  - **Shares/Contracts**
  - **R ($)** - Dollar risk per position
  - **Bucket** - Which correlation group

### 8. **Summary** (Settings & Configuration)
- Global settings (account size, risk %, etc.)
- Portfolio statistics
- **Key Named Ranges:**
  - `Equity_E` - Your account size
  - `RiskPct_r` - Risk per trade (% of account)
  - `HeatCap_H_pct` - Max total portfolio risk
  - `BucketHeatCap_pct` - Max risk per correlation group

### 9. **Control** (Hidden - Don't Touch!)
- Backend calculations
- Python integration test area
- You won't need to interact with this

---

## Understanding the Settings

Go to the **Summary** sheet. Here are the key settings you need to understand:

### Equity_E (Account Size)
- **Default:** $10,000
- **What it means:** Your total trading capital
- **Example:** If you have $25,000 in your account, change this to 25000

### RiskPct_r (Risk Per Trade)
- **Default:** 0.75% (0.0075)
- **What it means:** How much of your account you risk per trade
- **Example:** With $10,000 account and 0.75% risk:
  - Risk per trade = $10,000 × 0.0075 = **$75**
- **Conservative:** 0.5% ($50 on $10k)
- **Aggressive:** 1.5% ($150 on $10k)

### StopMultiple_K (ATR Stop Distance)
- **Default:** 2
- **What it means:** How many ATRs away from entry is your stop loss
- **Example:**
  - Stock at $100, ATR = $2, K = 2
  - Stop = $100 - (2 × $2) = **$96**

### HeatCap_H_pct (Portfolio Heat Cap)
- **Default:** 4% (0.04)
- **What it means:** Max total dollar risk across ALL open positions
- **Example:** With $10,000 account and 4% cap:
  - Max total risk = $10,000 × 0.04 = **$400**
  - If you have 3 positions risking $75 each ($225 total), you can add 2 more

### BucketHeatCap_pct (Bucket Heat Cap)
- **Default:** 1.5% (0.015)
- **What it means:** Max risk in one correlation bucket
- **Example:** With $10,000 account and 1.5% cap:
  - Max bucket risk = $10,000 × 0.015 = **$150**
  - If you have 2 Tech stocks risking $75 each ($150), no more Tech trades allowed

### AddStepN (Add-On Step Size)
- **Default:** 0.5
- **What it means:** When scaling into a position, how many ATRs of profit before adding
- **Example:** Entry at $100, ATR = $2, AddStepN = 0.5
  - Add more at $100 + (0.5 × $2) = **$101**

### EarningsBufferDays (Earnings Blackout)
- **Default:** 3 days
- **What it means:** Don't enter trades within 3 days of earnings
- **Why:** Avoid overnight gap risk from earnings announcements

---

## Daily Morning Routine

Every trading day, before the market opens, you need to **import fresh candidates** from your scans.

### Option A: Manual Import (Works for Everyone)

1. Go to **TradeEntry** sheet
2. In cell **B5** (Preset), click the dropdown
3. Select a preset, e.g., **TF_BREAKOUT_LONG**
4. Click the **"Open FINVIZ"** button
   - Your browser opens to the FINVIZ screener
5. On FINVIZ website:
   - Look at the **Ticker column**
   - Select the tickers (Ctrl+Click to select multiple)
   - **Copy** them (Ctrl+C)
6. Back in Excel, click **"Import Candidates"** button
7. A dialog pops up: **"Enter tickers (comma-separated):"**
8. **Paste** (Ctrl+V) → Click **OK**
9. You'll see: **"Imported 15 candidates for TF_BREAKOUT_LONG"**

**Repeat for other presets** (TF_PULLBACK_LONG, TF_MOMENTUM_LONG, etc.)

### Option B: Auto-Scraping (Requires Python in Excel)

1. Click **"Test Python Integration"** on Setup sheet
2. If it says **"[OK] AVAILABLE"**, you can use auto-scraping!
3. Go to **TradeEntry** sheet
4. Select preset in **B5**
5. Click **"Import Candidates"**
6. Wait 5-10 seconds → **Done!** Tickers auto-imported

**If Python says NOT AVAILABLE:** Use Option A (manual). It's only 30 seconds slower.

### After Importing

Go to the **Candidates** sheet - you'll see all imported tickers with today's date.

---

## Evaluating a Single Trade

Let's walk through a real example: **Trading AAPL (Apple) stock**.

### Step 1: Select the Trade Setup

1. Go to **TradeEntry** sheet
2. **Cell B5** (Preset): Select **TF_BREAKOUT_LONG** from dropdown
3. **Cell B6** (Ticker): Select **AAPL** from dropdown
   - *This dropdown shows all tickers you imported today*
4. **Cell B7** (Sector): Select **Technology**
5. **Cell B8** (Bucket): Select **Tech/Comm**

### Step 2: Enter Price Data

You need 3 numbers: **Entry price, ATR, and K**.

#### Where to Get These Numbers:

**Entry Price (B10):**
- The price where you plan to enter
- Example: AAPL is trading at **$180.50**
- Enter: **180.50**

**ATR N (B11):**
- ATR = Average True Range (volatility measure)
- **On TradingView:**
  1. Open AAPL chart
  2. Click **Indicators** → Search "ATR"
  3. Add **"Average True Range"**
  4. Look at the current ATR value on the chart
  5. Example: ATR shows **1.50**
- Enter: **1.50**

**K (B12) - Stop Multiple:**
- How many ATRs away is your stop loss?
- **Common values:**
  - K = 1 (tight stop, 1 ATR away)
  - K = 2 (standard stop, 2 ATRs away) ← **Most common**
  - K = 3 (wide stop, 3 ATRs away)
- Enter: **2**

**What just happened?**
- Entry = $180.50
- Stop = $180.50 - (2 × $1.50) = $180.50 - $3.00 = **$177.50**
- Your risk per share = **$3.00**

### Step 3: Choose Position Type (Stock or Options)

**Cell C13** - Method dropdown:

#### Method 1: Stock (Shares)
- Select: **1**
- You'll buy shares of AAPL
- **When to use:** Low-risk, long-term holds, liquid stocks

#### Method 2: Opt-DeltaATR (Options with Delta)
- Select: **2**
- New fields appear:
  - **Delta (B16):** Enter option delta, e.g., **0.70**
  - **Contracts (B17):** Enter # of contracts per 100 shares, e.g., **1**
- **When to use:** Buying call/put options with known delta

#### Method 3: Opt-MaxLoss (Options with Max Loss)
- Select: **3**
- New field appears:
  - **MaxLoss/Contract (B18):** Enter max loss per contract, e.g., **250**
- **When to use:** Defined-risk spreads (credit spreads, iron condors)

**For this example, let's use Method 1 (Stock):**
- Set **C13** to **1**

### Step 4: Complete the 6-Item Checklist

Now you manually check each item. Click each checkbox as you verify:

#### ☐ FromPreset (C20)
- **Question:** Did this ticker come from today's FINVIZ import?
- **How to check:** Look at Candidates sheet - is AAPL there with today's date?
- **If YES:** ✓ Check the box
- **If NO:** ✗ Leave unchecked (this is a random idea, not from your scan)

#### ☐ TrendPass (C21)
- **Question:** Is the stock in a confirmed uptrend?
- **How to check on TradingView:**
  1. Open AAPL daily chart
  2. Add 3 moving averages: 20 SMA, 50 SMA, 200 SMA
     - Indicators → "Moving Average" → Set to 20, Simple
     - Repeat for 50 and 200
  3. **For LONG:** Price > 20 SMA > 50 SMA > 200 SMA (all aligned)
  4. **For SHORT:** Price < 20 SMA < 50 SMA < 200 SMA
- **If trend confirmed:** ✓ Check the box
- **If choppy/sideways:** ✗ Leave unchecked

#### ☐ LiquidityPass (C22)
- **Question:** Does this stock trade enough volume?
- **How to check:**
  1. Look at **Average Volume** on your broker or Yahoo Finance
  2. **Rule:** Must be > 500,000 shares/day
  3. Example: AAPL averages 50 million/day → **PASS**
- **If volume > 500K:** ✓ Check the box
- **If volume < 500K:** ✗ Leave unchecked (too illiquid, skip it)

#### ☐ TVConfirm (C23)
- **Question:** Does your TradingView strategy show a signal?
- **How to check:**
  1. Open your custom TradingView strategy
  2. Look for a BUY arrow or signal on AAPL
  3. If you don't have a strategy, use this as "My own analysis confirms"
- **If confirmed:** ✓ Check the box
- **If no signal:** ✗ Leave unchecked

#### ☐ EarningsOK (C24)
- **Question:** Is earnings more than 3 days away?
- **How to check:**
  1. Google: "AAPL earnings date"
  2. Or check your broker's earnings calendar
  3. Example: Today is May 1, earnings is May 10 → **9 days away** → PASS
  4. Example: Today is May 1, earnings is May 3 → **2 days away** → FAIL
- **If > 3 days:** ✓ Check the box
- **If < 3 days:** ✗ Leave unchecked (too risky)

#### ☐ JournalOK (C25)
- **Question:** Have I reviewed this in my trading journal?
- **How to check:**
  1. Open your trading journal (OneNote, spreadsheet, etc.)
  2. Write down your thesis: "AAPL breaking above resistance at $180"
  3. Review your recent trades - any patterns to avoid?
  4. Ask: "Am I trading emotionally or systematically?"
- **If journaled and reviewed:** ✓ Check the box
- **If skipped journaling:** ✗ Leave unchecked (forces discipline!)

### Step 5: Click "Evaluate"

1. Click the **"Evaluate"** button
2. Watch the **banner** (top of sheet, row 2) change color:

**GREEN - "GO! All checks passed"**
- All 6 boxes checked ✓
- You can proceed to sizing

**YELLOW - "CAUTION - 1 item missing"**
- 5 out of 6 boxes checked
- Review the missing item - is it critical?
- Proceed with caution

**RED - "NO-GO - 2+ items missing"**
- 4 or fewer boxes checked
- Do NOT take this trade
- Missing too many criteria

**For our example:** If all 6 boxes are checked, you'll see **GREEN**.

### Step 6: Calculate Position Size

1. Click **"Recalc Sizing"** button
2. Look at the **Output Section** (right side, columns E-F):

**R ($)** - Dollar Risk:
- Shows: **$75.00** (your 0.75% of $10k account)

**Stop** - Stop Loss Price:
- Shows: **$177.50** (Entry $180.50 - 2 ATR)

**Target** - Take Profit Price:
- Shows: **$186.50** (Entry $180.50 + 2 ATR)

**RiskPerShare**:
- Shows: **$3.00** (difference between entry and stop)

**Shares** (or Contracts):
- Shows: **25 shares**
- **Calculation:** $75 risk ÷ $3 per share = 25 shares

**Position Value**:
- Shows: **$4,512.50** (25 shares × $180.50)

**Portfolio Heat** (before this trade):
- Shows current total risk: **$150** (if you have 2 other trades risking $75 each)

**Bucket Heat** (before this trade):
- Shows current Tech/Comm risk: **$75** (if you have 1 other Tech stock)

**Portfolio Heat (after)**:
- Shows: **$225** ($150 + $75)
- **Limit:** $400 (4% of $10k)
- **Status:** ✓ Within limit

**Bucket Heat (after)**:
- Shows: **$150** ($75 + $75)
- **Limit:** $150 (1.5% of $10k)
- **Status:** ✓ At limit (can't add more Tech after this)

### Step 7: Wait for Impulse Timer

**This is important** - prevents FOMO entries!

1. Click **"Start Timer"** button
2. A timestamp appears: **"Timer: 5/1/2024 9:32:15 AM"**
3. You must wait **2 minutes** before you can save
4. Use this time to:
   - Double-check your charts
   - Review the trade thesis
   - Make sure you're not emotional
   - Check if price has moved against you

**Why 2 minutes?** Studies show a forced delay reduces impulsive, emotional trading by 40%.

### Step 8: Save the Decision

After 2 minutes have elapsed:

1. Click **"Save Decision"** button
2. The system performs **5 hard gate checks:**

**Gate 1:** Banner must be GREEN
- If YELLOW or RED → **BLOCKED**
- Error: "Banner must be GREEN to save"

**Gate 2:** Ticker must be in today's Candidates
- If not in Candidates → **BLOCKED**
- Error: "Ticker not found in today's candidates"

**Gate 3:** 2-minute timer must be elapsed
- If < 2 minutes → **BLOCKED**
- Error: "Impulse timer: Please wait X more seconds"

**Gate 4:** Bucket must NOT be in cooldown
- If bucket had 2 stopouts in last 20 days → **BLOCKED**
- Error: "Bucket Tech/Comm is in cooldown until 5/10/2024"

**Gate 5:** Heat caps must not be exceeded
- If adding this trade would exceed limits → **BLOCKED**
- Error: "Would exceed portfolio heat cap (4%)"

**If all 5 gates pass:**
- Trade is logged to **Decisions** sheet
- Success message: "Trade decision saved for AAPL"
- Timestamp, all details, and checklist items are recorded

### Step 9: Execute the Trade

Now go to your broker and:

1. **Buy 25 shares of AAPL at $180.50** (or current market price)
2. **Set stop loss at $177.50**
3. **Optional:** Set target at $186.50
4. Manually add to **Positions** sheet for heat tracking:
   - Ticker: AAPL
   - Entry: 180.50
   - Stop: 177.50
   - Shares: 25
   - R: 75
   - Bucket: Tech/Comm

**You're done!** You just evaluated and logged your first trade.

---

## Understanding Each Field

Let me break down every field in plain English:

### Input Section (Left Side)

| Field | Cell | What It Means | Example |
|-------|------|---------------|---------|
| **Preset** | B5 | Which FINVIZ scan did you use? | TF_BREAKOUT_LONG |
| **Ticker** | B6 | What stock are you trading? | AAPL |
| **Sector** | B7 | What sector is it in? (for reference) | Technology |
| **Bucket** | B8 | Which correlation group? (for heat tracking) | Tech/Comm |
| **Entry** | B10 | At what price are you entering? | 180.50 |
| **ATR N** | B11 | What's the current Average True Range? | 1.50 |
| **K** | B12 | How many ATRs away is your stop? | 2 |
| **Method** | C13 | Stock (1), Opt-DeltaATR (2), or Opt-MaxLoss (3)? | 1 |
| **Delta** | B16 | For options: What's the delta? (0-1) | 0.70 |
| **Contracts/100sh** | B17 | For options: How many contracts per 100 shares? | 1 |
| **MaxLoss/Contract** | B18 | For spreads: Max loss per contract? | 250 |

### Checklist Section (Left Side)

| Item | Cell | What You're Checking | How to Verify |
|------|------|---------------------|---------------|
| **FromPreset** | C20 | Is ticker from today's import? | Check Candidates sheet |
| **TrendPass** | C21 | Is trend confirmed? | Check 20/50/200 SMA alignment |
| **LiquidityPass** | C22 | Volume > 500K/day? | Check Yahoo Finance |
| **TVConfirm** | C23 | TradingView signal present? | Check your TV strategy |
| **EarningsOK** | C24 | Earnings > 3 days away? | Google earnings date |
| **JournalOK** | C25 | Reviewed in journal? | Check your journal |

### Output Section (Right Side)

| Field | Cell | What It Shows | Example |
|-------|------|---------------|---------|
| **R ($)** | F5 | Dollar risk for this trade (RiskPct × Equity) | $75.00 |
| **Stop** | F6 | Stop loss price (Entry - K × ATR) | $177.50 |
| **Target** | F7 | Take profit price (Entry + K × ATR) | $186.50 |
| **RiskPerShare** | F8 | Risk per share (Entry - Stop) | $3.00 |
| **Shares/Contracts** | F9 | Position size (R ÷ RiskPerShare) | 25 |
| **PositionValue** | F10 | Total capital required (Shares × Entry) | $4,512.50 |
| **Portfolio Heat** | F11 | Current total risk across all positions | $150 |
| **Bucket Heat** | F15 | Current risk in this correlation bucket | $75 |

---

## The 6-Item Checklist Explained

This is your **GO/NO-GO decision framework**. Each item prevents a specific type of trading mistake.

### 1. FromPreset - "Is This From My Scan?"

**Purpose:** Prevents random, impulsive trades.

**Bad Example:** You see a headline "XYZ STOCK SURGES 20%!" and immediately want to trade it.
- Not from your scan → ✗ Unchecked
- Banner goes YELLOW/RED
- You're forced to pause and ask: "Is this systematic or emotional?"

**Good Example:** AAPL came up in your TF_BREAKOUT_LONG scan this morning.
- In Candidates sheet → ✓ Checked

**Trading Psychology:** This prevents "shiny object syndrome" - chasing random stocks because they're in the news.

### 2. TrendPass - "Is The Trend Confirmed?"

**Purpose:** "The trend is your friend" - only trade with the trend.

**Bad Example:** Stock is choppy, moving averages are tangled.
- 20 SMA below 50 SMA, price whipsawing → ✗ Unchecked
- You're trading in a zone of uncertainty

**Good Example:** Stock is in strong uptrend.
- Price > 20 SMA > 50 SMA > 200 SMA → ✓ Checked
- All moving averages aligned

**Trading Psychology:** Most traders lose money fighting the trend. This keeps you on the right side.

### 3. LiquidityPass - "Can I Get In and Out?"

**Purpose:** Avoid illiquid stocks where you can't exit your position.

**Bad Example:** Small-cap stock with 50,000 shares/day volume.
- You try to buy 10,000 shares → You move the price against yourself
- You try to sell → No buyers, can't exit
- Volatility is fake (low volume) → ✗ Unchecked

**Good Example:** AAPL with 50 million shares/day.
- You can buy/sell anytime without impact → ✓ Checked

**Trading Psychology:** "You're not profitable until you can EXIT profitably."

### 4. TVConfirm - "Does My Strategy Agree?"

**Purpose:** Multi-indicator confirmation (not just one signal).

**Bad Example:** FINVIZ shows a breakout, but your TradingView RSI shows overbought.
- Conflicting signals → ✗ Unchecked
- One system says yes, the other says no

**Good Example:** FINVIZ breakout + TradingView momentum signal + Your analysis.
- Multiple confirmations → ✓ Checked

**Trading Psychology:** Prevents "one-indicator bias" - seeing what you want to see.

### 5. EarningsOK - "Am I About to Get Gapped?"

**Purpose:** Avoid overnight earnings gaps that blow through your stop loss.

**Bad Example:** You buy AAPL on Monday, earnings is Wednesday after close.
- Wednesday after hours: AAPL reports, drops 10% overnight
- Your 2 ATR stop ($3) is useless against a $18 gap → ✗ Unchecked

**Good Example:** Earnings is 2 weeks away.
- Plenty of time for your trade to work → ✓ Checked

**Trading Psychology:** "Hope is not a strategy" - earnings is binary risk.

### 6. JournalOK - "Have I Thought This Through?"

**Purpose:** Forces deliberate, thoughtful trading (not reactive).

**Bad Example:** See a breakout, click buy immediately.
- No journal entry → ✗ Unchecked
- You're trading on impulse

**Good Example:** You write in journal:
- "AAPL breaking above $180 resistance on volume"
- "Sector strength in tech"
- "Checked for recent stopouts - none in Tech bucket"
- "Not revenge trading (last trade was a win)"
- Checkbox → ✓ Checked

**Trading Psychology:** The act of writing forces System 2 thinking (slow, analytical) instead of System 1 (fast, emotional).

---

## Position Sizing Methods

The system supports **3 different methods** for sizing positions. Choose based on what you're trading.

### Method 1: Stock (Shares)

**When to use:** Buying/shorting shares of stock.

**Formula:**
```
Shares = R ÷ RiskPerShare
Shares = R ÷ (Entry - Stop)
```

**Example:**
- Account: $10,000
- Risk: 0.75% = $75
- Entry: $180.50
- Stop: $177.50
- RiskPerShare: $180.50 - $177.50 = $3.00
- **Shares = $75 ÷ $3 = 25 shares**

**Position Value:** 25 × $180.50 = $4,512.50 (you need this much buying power)

**If trade hits stop:**
- Loss = 25 shares × $3/share = **$75** (exactly your R)

**If trade hits target (+2 ATR):**
- Gain = 25 shares × $3/share = **$75** (1R profit)

### Method 2: Opt-DeltaATR (Options with Delta)

**When to use:** Buying calls or puts (single-leg options).

**Formula:**
```
Contracts = R ÷ (K × N × Delta × 100)
```

**Example:**
- Account: $10,000
- Risk: $75
- K: 2 (ATR multiple)
- N: $1.50 (ATR)
- Delta: 0.70 (option delta)
- **Contracts = $75 ÷ (2 × $1.50 × 0.70 × 100)**
- **Contracts = $75 ÷ $210 = 0.36 → Round to 0 contracts**

**Uh oh!** 0.36 contracts is less than 1. What does this mean?

**Solution:** Your risk ($75) is too small for this option.
- **Option 1:** Increase your risk percentage (1% instead of 0.75%)
- **Option 2:** Use a different strike with higher delta (0.80+)
- **Option 3:** Buy shares instead of options

**Better Example:**
- Risk: $150 (1.5% of $10k)
- Delta: 0.80
- **Contracts = $150 ÷ (2 × $1.50 × 0.80 × 100)**
- **Contracts = $150 ÷ $240 = 0.625 → Round to 1 contract**

**Why this formula?**
- K × N = Stop distance in dollars ($3)
- × Delta = How much option moves per $1 of stock ($3 × 0.70 = $2.10)
- × 100 = Contract multiplier (1 contract = 100 shares)
- Total risk = $2.10 × 100 = $210 per contract

### Method 3: Opt-MaxLoss (Options with Max Loss)

**When to use:** Defined-risk spreads (credit spreads, iron condors, butterflies).

**Formula:**
```
Contracts = R ÷ MaxLossPerContract
```

**Example: Credit Spread**
- You sell a $180/$175 put spread on AAPL
- Max loss per contract: **$500** (width of $5 × 100)
- Risk: $75
- **Contracts = $75 ÷ $500 = 0.15 → Round to 0**

**Same problem!** Your risk is too small for this spread.

**Solution:**
- **Option 1:** Increase risk to $500 (5% of $10k) - risky!
- **Option 2:** Use narrower spread ($180/$178, max loss $200)
- **Option 3:** Trade multiple small accounts

**Better Example: Narrow Spread**
- Sell $180/$179 put spread (width = $1)
- Max loss: **$100** per contract
- Risk: $75
- **Contracts = $75 ÷ $100 = 0.75 → Round to 1 contract**

**Note:** You're risking $100 to meet your $75 target (slightly over-risked).

**Alternative:** Use fractional sizing:
- Risk: $200 (2% of $10k)
- **Contracts = $200 ÷ $100 = 2 contracts**
- Now you're risking $200 (exactly 2R)

---

## Heat Management

**Heat = Total Dollar Risk** across all open positions. This prevents blowing up your account.

### Portfolio Heat Cap

**Question:** How much total risk can you have across ALL positions?

**Default:** 4% of equity = $400 on $10k account

**Example Scenario:**

| Position | Ticker | R ($) | Status |
|----------|--------|-------|--------|
| 1 | AAPL | $75 | Open |
| 2 | MSFT | $75 | Open |
| 3 | NVDA | $75 | Open |
| **Total** | - | **$225** | **Within cap** ✓ |

**Can you add another trade?**
- Current heat: $225
- New trade: $75
- After trade: $225 + $75 = $300
- Cap: $400
- **Status:** ✓ YES, you're at $300 / $400 (75% utilized)

**What if you tried to add a 6th trade?**
- Current: $225 + $75 + $75 + $75 = $450 (after 4th trade)
- After 5th: $450 + $75 = $525
- Cap: $400
- **Status:** ✗ BLOCKED - "Would exceed portfolio heat cap"

**Why this matters:** Prevents over-leveraging. If you had 10 positions each risking $75 ($750 total), and all stopped out, you'd lose **7.5% of your account in one day**. The 4% cap limits this to max 4% loss.

### Bucket Heat Cap

**Question:** How much risk can you have in one correlation group?

**Default:** 1.5% of equity = $150 on $10k account

**Example Scenario:**

| Position | Ticker | Bucket | R ($) | Status |
|----------|--------|--------|-------|--------|
| 1 | AAPL | Tech/Comm | $75 | Open |
| 2 | MSFT | Tech/Comm | $75 | Open |
| **Tech Bucket Total** | - | - | **$150** | **At cap** ⚠️ |

**Can you add NVDA (also Tech/Comm)?**
- Current Tech heat: $150
- New trade: $75
- After trade: $225
- Cap: $150
- **Status:** ✗ BLOCKED - "Would exceed bucket heat cap for Tech/Comm"

**Can you add JPM (Financials bucket)?**
- Current Financials heat: $0
- New trade: $75
- After trade: $75
- Cap: $150
- **Status:** ✓ YES, different bucket

**Why this matters:** Tech stocks are correlated. If there's a sector crash (tech selloff), all your Tech positions could hit stop on the same day. Bucket caps prevent you from having 5 tech stocks that all crash together.

### Bucket Cooldown

**Question:** What if you keep losing in the same bucket?

**Rule:** If you get stopped out **2+ times in one bucket within 20 days**, that bucket goes into **cooldown for 10 days**.

**Example Timeline:**

| Date | Event | Bucket Status |
|------|-------|---------------|
| May 1 | AAPL stop hit in Tech/Comm | 1 stopout |
| May 5 | MSFT stop hit in Tech/Comm | 2 stopouts → **COOLDOWN STARTS** |
| May 5-15 | Tech/Comm bucket locked | Cannot enter new Tech trades |
| May 15 | Cooldown expires | Tech/Comm bucket available again |

**What happens if you try to trade NVDA on May 8?**
- Click "Save Decision"
- Error: "Bucket Tech/Comm is in cooldown until 5/15/2024"
- Trade is **BLOCKED**

**Why this matters:** If you're getting stopped out repeatedly in one sector, **the sector is not working**. Cooldown forces you to stop digging deeper into a losing theme.

**Manual Override:** If you absolutely must trade during cooldown:
1. Go to **Buckets** sheet
2. Change **CooldownUntil** to a past date
3. System will allow the trade
4. **Warning:** You're overriding your own risk rules. Proceed with caution.

---

## Troubleshooting

### "Please select a preset first"

**Problem:** You clicked "Import Candidates" but cell B5 is empty.

**Fix:**
1. Click cell **B5**
2. Click the dropdown arrow
3. Select a preset (e.g., TF_BREAKOUT_LONG)
4. Try again

### Dropdown is missing / no arrow in B5

**Problem:** Data validation didn't get created.

**Fix:**
1. Go to **Setup** sheet
2. Click **"Rebuild TradeEntry UI"** button
3. Wait for "TradeEntry UI built successfully!"
4. Check cell B5 again

**If still missing:**
1. Press Alt+F11 (VBA Editor)
2. Press Ctrl+G (Immediate Window)
3. Type: `TF_UI.BindControls`
4. Press Enter
5. Close VBA Editor

### "Ticker not found in today's candidates"

**Problem:** You're trying to save a trade for a ticker you didn't import.

**Fix:**
1. Go to **Candidates** sheet
2. Check if the ticker is there with **today's date**
3. If not, go back to TradeEntry and import it first
4. Then try saving the trade again

### "Impulse timer: Please wait X more seconds"

**Problem:** You haven't waited 2 minutes since clicking "Start Timer".

**Fix:**
- **Wait!** This is intentional.
- Use the time to double-check your setup
- After 2 minutes, click "Save Decision" again

**To bypass (not recommended):**
1. Press Alt+F11
2. Ctrl+G
3. Type: `Worksheets("TradeEntry").Range("C30").Value = Now - TimeValue("00:03:00")`
4. This sets timer to 3 minutes ago

### "Would exceed portfolio heat cap"

**Problem:** Adding this trade would put you over 4% total risk.

**Options:**
1. **Wait for an existing position to close** (reduces heat)
2. **Reduce risk per trade** (change RiskPct_r on Summary sheet)
3. **Increase heat cap** (change HeatCap_H_pct on Summary sheet) ⚠️ Risky!
4. **Skip this trade** (discipline!)

### "Bucket X is in cooldown"

**Problem:** You had 2+ stopouts in this bucket recently.

**Options:**
1. **Wait for cooldown to expire** (check CooldownUntil date on Buckets sheet)
2. **Trade a different bucket** (e.g., if Tech is locked, trade Healthcare)
3. **Manual override** (edit Buckets sheet, set CooldownUntil to yesterday) ⚠️ Risky!

### Python says "NOT AVAILABLE" but I have Python in Excel

**Problem:** VBA detection might not be working.

**Fix 1: Verify Python Works**
1. Click any cell in Excel
2. Type: `=PY(1+1)`
3. Press Enter
4. Do you see `2` or a Python object? → Python works!
5. If you see `#NAME?` error → Python is not enabled

**Fix 2: Enable Python in Excel**
1. Click **Data** tab
2. Look for **"Python in Excel"** section
3. If missing → You need Microsoft 365 Insider build
4. If present → Click **"Enable Python"**

**Fix 3: Use Manual Import**
- Python auto-scraping is optional
- Manual import works perfectly
- Only 30 seconds slower

### Checkboxes didn't auto-create

**Problem:** The BuildTradeEntryUI ran but no checkboxes appeared.

**Fix: Add Manually**
1. Go to **TradeEntry** sheet
2. Click **Developer** tab
3. Click **Insert** → **Check Box** (under Form Controls)
4. Draw a small checkbox next to row 21
5. Right-click → **Format Control**
6. Cell link: **$C$20**
7. Click OK
8. Right-click → **Edit Text** → Delete the text
9. Repeat for rows 22-26, linking to C21, C22, C23, C24, C25

### Buttons are duplicated

**Problem:** You ran BuildTradeEntryUI multiple times.

**Fix:**
1. Go to **Setup** sheet
2. Click **"Rebuild TradeEntry UI"**
3. This deletes all shapes and rebuilds clean
4. If still duplicated, manually delete:
   - Click each duplicate button
   - Press Delete key

### "Compile error" when running macros

**Problem:** VBA code has a syntax error.

**Fix:**
1. Note which module is highlighted
2. Click **Reset** button (stop sign icon)
3. Contact support or rebuild:
   - Close Excel
   - Run `BUILD.bat` again
   - Open fresh workbook

---

## Quick Reference Card

**Print this page and keep it next to your monitor!**

### Daily Workflow Checklist

- [ ] **Morning (10 min):** Import candidates from 3-5 presets
- [ ] **Per Trade (5 min):** Select ticker → Enter data → 6 checks → Evaluate → Size → Timer → Save
- [ ] **After Entry:** Add to Positions sheet, set stop in broker
- [ ] **End of Day:** Update Positions sheet with current prices

### The 6 Checks (Memorize This!)

1. ✓ **FromPreset** - Is it from today's scan? (Check Candidates sheet)
2. ✓ **TrendPass** - Is trend confirmed? (20 > 50 > 200 SMA)
3. ✓ **LiquidityPass** - Volume > 500K? (Check Yahoo Finance)
4. ✓ **TVConfirm** - Strategy signal? (Check TradingView)
5. ✓ **EarningsOK** - Earnings > 3 days? (Google earnings date)
6. ✓ **JournalOK** - Journaled? (Write in journal)

### Critical Reminders

- **GREEN or NO-GO** - Don't trade YELLOW/RED banners
- **Wait 2 minutes** - Impulse timer is mandatory
- **Check heat** - Don't exceed caps (4% portfolio, 1.5% bucket)
- **Respect cooldowns** - If bucket is paused, skip it
- **Journal everything** - Force yourself to think

### Key Formulas

- **R ($)** = Equity × RiskPct
- **Stop** = Entry - (K × ATR)
- **Shares** = R ÷ (Entry - Stop)
- **Portfolio Heat** = Sum of all open R ($)

### Emergency Contacts

- **Rebuild UI:** Setup sheet → "Rebuild TradeEntry UI"
- **Test Python:** Setup sheet → "Test Python Integration"
- **Manual VBA:** Alt+F11 → Ctrl+G → Type command
- **Reset Timer:** Set C30 to past time (Ctrl+G hack above)

---

**You're ready to trade!** Start with the Daily Morning Routine, then evaluate your first trade using the step-by-step guide.

Remember: **The system is designed to protect you from yourself.** If it blocks a trade, that's a FEATURE, not a bug. Trust the process.

Happy trading! 📈
