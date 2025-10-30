# TF-Engine Quick Start

**Get started in 10 minutes**

**TF = Trend Following** - Ed Seykota style trend-following with systematic discipline enforcement

---

## 1. Install (2 minutes)

### Windows

1. Download `tf-engine.exe` or `TF-Engine-Setup.exe`
2. **Option A:** Run setup installer (if using setup)
   - Follow wizard
   - Install to default location
   - Desktop shortcut created automatically
3. **Option B:** Extract standalone binary (if using tf-engine.exe)
   - Place in desired folder (e.g., `C:\TF-Engine\`)
   - Create shortcut on Desktop if desired

### First Launch

1. Double-click `tf-engine.exe` or use Desktop shortcut
2. Command window opens (keep this open - this is the server)
3. Browser opens automatically to http://localhost:8080
4. If browser doesn't open: Manually navigate to http://localhost:8080

**Expected:** Dashboard screen with "Welcome" message

---

## 2. Configure Settings (2 minutes)

### Initial Configuration

1. Click **Settings** (gear icon in header, top-right)
2. Enter your account information:

**Required Fields:**
```
Equity:              100000     (your account size in dollars)
Risk % per unit:     0.75       (0.50-1.00% typical)
Portfolio heat cap:  4.0        (% of equity - total risk across all positions)
Bucket heat cap:     1.5        (% of equity - risk per sector)
Max units:           4          (initial + 3 add-ons)
```

3. Click **Save Settings**
4. Success message: "Settings saved"

**Why these numbers?**
- **0.75% per unit:** With 4 units max = 3% total risk per position
- **4% portfolio cap:** Never risk more than $4,000 total (on $100k account)
- **1.5% bucket cap:** Max $1,500 in any one sector (limits concentration)

---

## 3. Set Up TradingView (3 minutes)

### Install Pine Script

1. **Open TradingView:** https://www.tradingview.com
2. **Navigate to Chart** → Any ticker (e.g., SPY)
3. **Open Pine Editor** (bottom panel, "Pine Editor" tab)
4. **Copy Pine Script:**
   - Navigate to TF-Engine installation folder
   - Open `reference/Ed-Seykota.pine` in text editor
   - Copy entire contents (Ctrl+A, Ctrl+C)
5. **Paste in Pine Editor** (Ctrl+V)
6. **Click "Add to Chart"** (top-right of Pine Editor)
7. **Result:** Chart shows Donchian channels (blue/red/green lines) and ATR indicator below

**What it shows:**
- **Blue lines:** 55-bar high (long entry signal)
- **Red lines:** 55-bar low (short entry signal)
- **Green lines:** 10-bar exit levels
- **Indicator window:** N (ATR value) - you'll need this for sizing

**Pin to favorites:**
- Pine Editor → Favorites icon → Save as "Ed-Seykota"
- Now accessible on any chart instantly

---

## 4. Run Daily Workflow (3 minutes)

### Morning Workflow

**Step 1: Check Dashboard** (30 seconds)
- Open TF-Engine (should open to Dashboard)
- Review open positions (if any)
- Note portfolio heat (should be < 90% of cap)

**Step 2: Scan for Candidates** (if FINVIZ configured - optional)
- Click "Scanner" in navigation
- Click "Run Daily FINVIZ Scan"
- Select 10-15 tickers (checkboxes)
- Click "Import Selected"

**Step 3: Verify Signal in TradingView** (1 minute per ticker)
- Click "Open in TradingView" next to ticker
- Check: Is close > 55-bar high (blue line)? ✓
- Note: N (ATR) value (e.g., 2.35)
- Note: Entry price (today's close, e.g., 180.50)
- Note: Sector (e.g., Tech/Comm)

**Step 4: Complete Checklist** (2 minutes)
- Navigate to "Checklist"
- Fill in trade data:
  - Ticker: AAPL
  - Entry: 180.50
  - N (ATR): 2.35
  - Sector: Tech/Comm
  - Structure: Stock
- Check all 5 required gates:
  - ✓ Signal: 55-bar breakout confirmed
  - ✓ Risk/Size: 2×N stop, pyramid every 0.5×N
  - ✓ Liquidity: Volume > 1M shares/day
  - ✓ Exits: 10-bar Donchian OR 2×N stop
  - ✓ Behavior: Will honor 2-min timer
- Check optional quality items (recommended):
  - ✓ Regime OK: SPY > 200 SMA
  - ✓ No Chase: Entry within 2N of 20-EMA
  - ✓ Earnings OK: No earnings next 2 weeks
  - ✓ Journal: "Clean breakout, high volume, sector underweight"
- **Banner turns GREEN** ✓
- Click "Save Evaluation"
- **2-minute timer starts:** 2:00 → 1:59 → ...

**Step 5: Calculate Position Sizing** (1 minute)
- Navigate to "Position Sizing"
- Data pre-filled from checklist
- Click "Calculate Position Size"
- Review results:
  - Shares per unit: 159
  - Risk per unit: $747
  - Stop: $175.80
  - Add-ons: Every 0.5×N
- Click "Save Position Plan"

**Step 6: Check Heat** (30 seconds)
- Navigate to "Heat Check"
- Click "Check Heat for This Trade"
- Result:
  - ✓ Portfolio heat: Within cap
  - ✓ Bucket heat: Within cap
  - Verdict: TRADE APPROVED

**Step 7: Final Trade Entry** (30 seconds)
- Navigate to "Trade Entry"
- Review trade summary
- Click "Run Final Gate Check"
- Result: ALL GATES PASS ✓
- Click **"SAVE GO DECISION"** (green button now enabled)
- Success: "GO decision saved for AAPL"

**Step 8: Execute in Broker** (manual)
- Open your broker (Interactive Brokers, TD Ameritrade, etc.)
- Enter order: BUY 159 shares AAPL @ $180.50 (limit or market)
- Set stop: SELL 159 @ $175.80 STOP
- Set alerts: Add-on levels at $181.68, $182.85, $184.03

**Daily ongoing:**
- Check TradingView for 10-bar exit level
- Exit if close < 10-bar low OR stop hit at $175.80

---

## Next Steps

### Learn More

Read the comprehensive [USER_GUIDE.md](USER_GUIDE.md) for:
- Detailed explanation of 5 gates
- Understanding the 3-state banner (RED/YELLOW/GREEN)
- Position management (add-ons, exits, trailing stops)
- Heat management (portfolio & sector caps)
- Calendar view (sector diversification)
- Tips & best practices
- Troubleshooting

### Optional Configuration

**TradingView URL Template:**
- Settings → TradingView URL Template
- Enter: `https://tradingview.com/chart/?symbol={ticker}`
- Enables one-click chart opening

**FINVIZ Presets:**
- Settings → FINVIZ Presets → Add New
- Name: "TF Breakout Long"
- URL: Your FINVIZ screener URL
- Enables automated daily scanning

---

## Key Concepts (5-Minute Summary)

### The 5 Gates (Cannot Be Bypassed)

1. **Gate 1: Banner Status** → Must be GREEN
2. **Gate 2: Impulse Brake** → 2-minute timer must elapse
3. **Gate 3: Cooldowns** → Ticker/sector not on cooldown
4. **Gate 4: Heat Caps** → Portfolio & sector risk within limits
5. **Gate 5: Sizing** → Position sizing calculated and saved

All 5 must pass to save GO decision. **No exceptions.**

### The 3-State Banner

- **RED:** Missing required checklist items → DO NOT TRADE
- **YELLOW:** Required items OK but quality score low → CAUTION
- **GREEN:** All requirements met, quality score good → OK TO TRADE

### Position Sizing (Van Tharp Method)

```
Risk per unit = Equity × Risk% = $100,000 × 0.0075 = $750
Stop distance = 2 × N (ATR) = 2 × $2.35 = $4.70
Shares = Risk ÷ Stop = $750 ÷ $4.70 = 159 shares
Initial stop = Entry - 2×N = $180.50 - $4.70 = $175.80
```

### Exits (Two Rules - Whichever is Closer)

1. **10-bar Donchian opposite breakout** (check daily on TradingView)
2. **2×N stop** (set in broker: $175.80 for this example)

Exit immediately if either is hit.

### Add-Ons (Pyramiding)

Add to winners every 0.5×N:
- Unit 1: 159 @ $180.50 (initial entry)
- Unit 2: 159 @ $181.68 (+0.5N = +$1.18)
- Unit 3: 159 @ $182.85 (+1.0N = +$2.35)
- Unit 4: 159 @ $184.03 (+1.5N = +$3.53)

Set price alerts at these levels. Add manually when hit.

---

## Common Mistakes (Avoid These!)

### 1. Skipping the 2-Minute Timer
**Wrong:** Rushing through to save GO at 1:30 elapsed
**Right:** Wait the full 2:00 - use time to review analysis

### 2. Not Verifying Signal on TradingView
**Wrong:** Assuming FINVIZ scan results are valid breakouts
**Right:** Open every candidate in TradingView and verify close > 55-bar high

### 3. Ignoring Heat Caps
**Wrong:** "I really like this trade, I'll exceed the cap just this once"
**Right:** Honor caps. Reduce size, close existing positions, or skip trade.

### 4. Widening Stops
**Wrong:** Price approaching 2×N stop → "I'll give it more room"
**Right:** Honor your stop. Exit at 2×N (or 10-bar, whichever closer)

### 5. Not Journaling
**Wrong:** Checking journal box without writing anything
**Right:** Write 1-2 sentences: "Clean breakout, volume 2x avg, sector underweight"

### 6. Overtrading
**Wrong:** Forcing trades when no valid setups exist
**Right:** Sometimes best decision is NO-GO. Cash is a position.

---

## System Requirements

**Minimum:**
- Windows 10/11 (64-bit)
- 512 MB RAM
- 100 MB disk space
- Internet connection (for FINVIZ scanning)

**Recommended:**
- Windows 10/11 (64-bit)
- 1 GB RAM
- Monitor with 1920×1080 resolution (for comfortable multi-screen workflow)
- Dual monitors (TradingView on one, TF-Engine on other)

---

## Getting Help

**Resources (in order):**
1. **This Quick Start** (you just read it!)
2. **[USER_GUIDE.md](USER_GUIDE.md)** - Comprehensive guide
3. **[FAQ.md](FAQ.md)** - Frequently asked questions
4. **[TROUBLESHOOTING.md](TROUBLESHOOTING.md)** - Common issues
5. **Contact Support:** [your-email or GitHub issues]

---

## Philosophy in 30 Seconds

**TF-Engine is NOT a trading platform.**

It's a **discipline enforcement system**.

Every "inconvenience" (2-min timer, 5 gates, heat caps) is intentional. They exist to protect you from impulsive, emotional decisions.

> **Trade the tide, not the splash.**

The system works if you work the system.

**Good luck, and trade systematically!**

---

**Version:** 1.0.0
**Last Updated:** 2025-10-29
**Next:** Read [USER_GUIDE.md](USER_GUIDE.md) for comprehensive documentation
