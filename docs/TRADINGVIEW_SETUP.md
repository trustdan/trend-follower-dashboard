# TradingView Setup Guide

**Version:** 1.0.0
**Last Updated:** 2025-10-29

**TF = Trend Following** - Ed Seykota / Turtle Trader style visualization

This guide covers installing and using the Ed-Seykota Pine Script on TradingView for signal verification with TF-Engine.

---

## Table of Contents

1. [Overview](#overview)
2. [Create TradingView Account](#create-tradingview-account)
3. [Install Ed-Seykota Pine Script](#install-ed-seykota-pine-script)
4. [Understanding the Script](#understanding-the-script)
5. [Daily Usage Workflow](#daily-usage-workflow)
6. [Customize Chart Layout](#customize-chart-layout)
7. [TF-Engine Integration](#tf-engine-integration)
8. [Tips & Shortcuts](#tips--shortcuts)
9. [Troubleshooting](#troubleshooting)

---

## Overview

### Why TradingView?

TF-Engine does **NOT** generate signals automatically. You verify breakouts manually using TradingView charts with the Ed-Seykota Pine Script.

**TradingView provides:**
- Real-time charts (free account sufficient)
- Ed-Seykota Pine Script visualization
- 55-bar Donchian channels (entry signals)
- 10-bar Donchian channels (exit signals)
- N (ATR) value for position sizing
- Historical data for verification

### What You'll See

![TradingView Chart with Ed-Seykota Script](screenshots/tradingview-script.png)

**On chart:**
- **Blue lines:** 55-bar Donchian high/low (entry triggers)
- **Green triangles:** Long breakout signals (close > 55-high)
- **Red triangles:** Short breakout signals (close < 55-low)
- **Red stops:** Current stop levels (2×N or 10-bar, whichever closer)

**In indicator window (below chart):**
- **N (ATR):** Current ATR value (e.g., 2.35)
- Used for position sizing in TF-Engine

---

## Create TradingView Account

### Step 1: Sign Up

1. Navigate to: https://www.tradingview.com

2. Click **"Get started"** (top-right)

3. Sign up with:
   - Email + password, OR
   - Google account, OR
   - Apple ID, OR
   - Facebook, Twitter, LinkedIn

4. Verify email (check inbox for verification link)

5. Login

### Step 2: Choose Plan

**Free Plan (Sufficient for TF-Engine use):**
- Real-time data for stocks (delayed for some exchanges)
- 3 indicators per chart
- 1 saved chart layout
- Basic alerts

**Pro Plans ($14.95/mo+):**
- More indicators per chart
- More saved layouts
- Priority data
- Advanced alerts
- Volume profile, etc.

**Recommendation:** Start with **Free**. Upgrade later if desired.

### Step 3: Navigate to Chart

1. Top menu: **"Chart"** (or press Alt+C)

2. Search for a ticker (e.g., SPY)

3. Chart opens in new layout

---

## Install Ed-Seykota Pine Script

### Step 1: Open Pine Editor

![Pine Editor Location](screenshots/tradingview-pine-editor.png)

1. Chart screen (any ticker)

2. Bottom panel: **"Pine Editor"** tab
   - If not visible: View menu → Pine Editor

3. Pine Editor panel opens at bottom

### Step 2: Copy Script Code

**Option A: From TF-Engine installation folder**

```
1. Navigate to TF-Engine installation folder
   - Example: C:\TF-Engine\reference\Ed-Seykota.pine
2. Open Ed-Seykota.pine in text editor (Notepad, VS Code, etc.)
3. Select all (Ctrl+A)
4. Copy (Ctrl+C)
```

**Option B: From documentation**

The full Ed-Seykota Pine Script is included in your TF-Engine installation at:
```
reference/Ed-Seykota.pine
```

**Script Features:**
- Donchian breakout entry (55-bar System-2)
- Donchian exit (10-bar opposite)
- ATR-based stops (2×N initial stop)
- Risk-based position sizing
- Pyramiding every 0.5×N up to max units
- Market regime filter (optional)
- Alerts for entries, add-ons, exits

### Step 3: Paste Script in Pine Editor

![Pine Editor with Code](screenshots/pine-editor-code.png)

1. Click inside Pine Editor text area

2. Delete any existing code (if present)

3. Paste (Ctrl+V)

4. Script code appears in editor

### Step 4: Save Script

1. Click **"Save"** (diskette icon in Pine Editor toolbar)

2. Script name dialog:
   ```
   Name: Ed-Seykota
   ```

3. Click "Save" or press Enter

4. Script saved to your Pine Editor library

### Step 5: Add to Chart

1. Click **"Add to Chart"** (in Pine Editor toolbar, top-right)

2. Script loads on chart

3. **Expected result:**
   - Blue lines (Donchian channels) appear on chart
   - Indicator panel below chart shows "N = [value]"
   - Green/red triangles mark breakouts (if any in visible range)

**Success!** Script is installed.

---

## Understanding the Script

### Chart Overlays

![Chart with Donchian Channels](screenshots/tradingview-donchian.png)

#### Blue Lines: 55-Bar Donchian Channels (Entry)

- **Blue upper band:** Highest high of last 55 bars
- **Blue lower band:** Lowest low of last 55 bars

**Interpretation:**
- **Long entry:** Close > upper blue line (breakout above 55-high)
- **Short entry:** Close < lower blue line (breakout below 55-low)

#### Red Stops: Current Stop Levels

- **Red line (for longs):** Max of (Entry - 2×N, 10-bar low)
- **Red line (for shorts):** Min of (Entry + 2×N, 10-bar high)

**Interpretation:**
- Exit if price touches red line
- Red line trails as 10-bar level moves favorably

#### Green/Red Triangles: Breakout Signals

- **Green triangle (below bar):** Long breakout (close > 55-high)
- **Red triangle (above bar):** Short breakout (close < 55-low)

**Interpretation:**
- Visual confirmation of breakout
- Only appears on bars where breakout occurs

### Indicator Window

![Indicator Window with N Value](screenshots/tradingview-indicator.png)

#### N (ATR) Value

- **N = Average True Range** (volatility measure)
- **Default period:** 20 bars
- **Example:** N = 2.35

**Interpretation:**
- Used for position sizing: Stop distance = 2×N = $4.70
- Used for add-ons: Every 0.5×N = $1.18
- Higher N = higher volatility = fewer shares per unit

### Script Settings (Inputs)

Click gear icon (⚙️) next to script name in indicator list to adjust settings:

**Entry/Exit:**
- Donchian ENTRY lookback: 55 (System-2 standard)
- Donchian EXIT lookback: 10 (standard)
- N (ATR) length: 20 (standard)

**Risk Management:**
- Initial stop (in N): 2.0 (standard)
- Add every X * N: 0.5 (pyramid every half-N)
- Max units: 4 (initial + 3 add-ons)
- Risk % per unit: 1.0% (Ed Seykota standard)

**Optional Filters:**
- Use market regime filter: false (optional)
- Market symbol: SPY (if using regime filter)
- Market MA length: 200 (if using regime filter)

**Display:**
- Plot signals & stops: true
- Plot Donchian bands: true

**Recommendation:** Keep defaults unless you know what you're changing.

---

## Daily Usage Workflow

### Morning Routine

#### 1. Scan for Candidates (in TF-Engine)

```
TF-Engine → Scanner → Run FINVIZ Scan
Import 10-15 candidates
```

#### 2. Verify Each Candidate in TradingView

**For each ticker (e.g., AAPL):**

1. **Open chart:**
   - Click "Open in TradingView" in TF-Engine, OR
   - Manually search: AAPL in TradingView search bar

2. **Check if Ed-Seykota script is loaded:**
   - Script name should appear in indicator list
   - If not: Click Pine Editor → Open "Ed-Seykota" → Add to Chart

3. **Verify breakout signal:**
   ```
   For LONG:
   - Is today's close ABOVE blue upper line? ✓
   - Or: Is there a green triangle below today's bar? ✓

   For SHORT:
   - Is today's close BELOW blue lower line? ✓
   - Or: Is there a red triangle above today's bar? ✓
   ```

4. **Note key values:**
   ```
   - N (ATR): Read from indicator window (e.g., 2.35)
   - Entry price: Today's close (or your planned limit order)
   - Sector: Check company info (Tech/Comm, Energy, etc.)
   ```

5. **Visual quality check:**
   ```
   - Is breakout clean or choppy? (Clean = strong signal)
   - Is volume above average? (Check volume bars)
   - Is price near 20-EMA or extended? (Closer = better)
   - Any upcoming earnings? (TradingView → Calendar tab)
   ```

6. **Return to TF-Engine:**
   ```
   Fill checklist:
   - Ticker: AAPL
   - Entry: 180.50
   - N: 2.35
   - Sector: Tech/Comm
   - Check all gates...
   ```

#### 3. Repeat for All Candidates

- Process 10-15 tickers
- Takes 1-2 minutes per ticker
- Total: 15-30 minutes for daily scan

---

### Daily Exit Management

**Each evening (after market close):**

1. **Open TradingView for each open position**

2. **Check 10-bar exit level:**
   ```
   For longs:
   - Find 10-bar low (green line or manual)
   - Current 10-bar low: $188

   For shorts:
   - Find 10-bar high
   ```

3. **Compare to 2×N stop:**
   ```
   Original entry: $180.50
   Original stop: $175.80 (2×N)
   Current 10-bar low: $188

   Effective stop: $188 (closer of $188 vs $175.80)
   ```

4. **Update stop in broker:**
   ```
   Broker: Move stop to $188 (or slightly below for buffer)
   ```

5. **Exit if triggered:**
   ```
   If close < $188 → Exit position immediately
   ```

---

## Customize Chart Layout

### Save Chart Layout

**Once you've configured ideal layout:**

1. **Arrange chart:**
   - Timeframe: Daily (D)
   - Indicators: Ed-Seykota script + any other desired
   - Chart type: Candlesticks (or Bars)
   - Color scheme: Light/Dark to taste

2. **Save layout:**
   ```
   Chart top-right → Click dropdown arrow next to chart title
   → "Save as" → Name: "TF-Engine Daily"
   ```

3. **Load layout anytime:**
   ```
   Chart dropdown → Select "TF-Engine Daily"
   All settings restore instantly
   ```

**Free account limitation:** 1 saved layout
**Pro accounts:** 5-10+ saved layouts

---

### Recommended Chart Settings

**Timeframe:** Daily (D)
- Trend-following works on daily bars
- Use higher timeframes (Weekly) for context
- Avoid intraday (noise)

**Chart Type:** Candlesticks or Bars
- Candlesticks show open/close visually
- Bars are cleaner (personal preference)

**Indicators:**
1. Ed-Seykota script (required)
2. Volume (default, keep visible)
3. Optional: 20-EMA (for "No Chase" quality check)
4. Optional: 200-SMA (for regime check)

**Color Scheme:**
- Light mode: Good for daytime analysis
- Dark mode: Easy on eyes for evening reviews

---

### Multi-Chart Layouts (Pro Feature)

**If you have Pro account:**

Create 4-chart layout:
- Top-left: SPY (market regime)
- Top-right: Candidate 1
- Bottom-left: Candidate 2
- Bottom-right: Candidate 3

Quickly scan multiple tickers simultaneously.

---

## TF-Engine Integration

### One-Click TradingView Opening

**Configure in TF-Engine Settings:**

1. **Navigate to Settings**

2. **TradingView URL Template:**
   ```
   https://tradingview.com/chart/?symbol={ticker}
   ```

3. **Save settings**

4. **Usage:**
   ```
   TF-Engine → Scanner (or any ticker list)
   → Click "Open in TradingView" next to AAPL
   → Browser tab opens: https://tradingview.com/chart/?symbol=AAPL
   → Chart loads with AAPL
   ```

**Advanced template (if you saved chart layout):**
```
https://tradingview.com/chart/[LAYOUT_ID]/?symbol={ticker}
```
- Replace `[LAYOUT_ID]` with your saved layout ID
- Find ID: Open saved layout → Copy from URL

---

### Alert Setup (Optional)

**Set TradingView alerts for add-on levels:**

1. **On chart:** Right-click on price level
   ```
   → Add alert
   ```

2. **Configure alert:**
   ```
   Condition: AAPL crossing $181.68
   Options: Once per bar close
   Alert actions: Popup, Email, SMS (if Pro)
   Message: "AAPL add-on level 2 hit (+0.5N)"
   ```

3. **Create alert**

4. **When alert fires:**
   ```
   → Check TF-Engine heat caps still OK
   → Execute add-on in broker: BUY 159 shares @ market
   ```

---

## Tips & Shortcuts

### Keyboard Shortcuts (TradingView)

| Action | Shortcut |
|--------|----------|
| Search ticker | / (slash) |
| Zoom in | + or mouse wheel up |
| Zoom out | - or mouse wheel down |
| Pan left/right | Arrow keys or drag chart |
| Reset zoom | Double-click chart |
| Indicators | / (slash) → type indicator name |
| Full screen | F |
| Chart style | 1 (bars) / 2 (candles) / 3 (line) |
| Timeframe | D (daily), W (weekly), M (monthly) |

### Workflow Shortcuts

**Quick ticker switching:**
```
TradingView search bar (top) → Type ticker → Enter
Chart updates to new ticker
Ed-Seykota script remains loaded (yay!)
```

**Compare multiple tickers:**
```
Chart → Compare symbol (+ icon near search)
→ Add MSFT, GOOG, etc.
→ See relative performance overlaid
```

**Screenshot for journal:**
```
Chart → Camera icon (top toolbar)
→ Takes screenshot
→ Save to file or copy to clipboard
→ Paste in journal notes
```

---

### Mobile App (Optional)

**TradingView mobile app:**
- iOS: App Store
- Android: Google Play

**Mobile workflow:**
- Morning: Full analysis on desktop
- Intraday: Check exit levels on mobile
- Evening: Update stops on mobile

**Pine Scripts work on mobile:** Ed-Seykota script loads on mobile charts (if saved to favorites)

---

## Troubleshooting

### Script doesn't appear on chart

**Solution 1: Re-add script**
```
Pine Editor → Open "Ed-Seykota" (from saved list)
→ Click "Add to Chart"
```

**Solution 2: Check indicators list**
```
Chart → Indicators button (f(x) icon)
→ Look for "Ed-Seykota" in list
→ If hidden (eye icon with slash), click to show
```

**Solution 3: Re-install script**
```
Copy script code again
Pine Editor → New → Paste → Save → Add to Chart
```

---

### N (ATR) value not showing

**Solution 1: Check indicator window**
```
Below chart → Should see "Seykota / Turtle Core v2.1 + Date Range"
→ Panel shows N = [value]

If panel missing:
→ Script might not be loaded (see above)
```

**Solution 2: Expand indicator window**
```
Drag separator line between chart and indicator window
Make indicator panel taller
```

**Solution 3: Check script settings**
```
Gear icon next to script name
→ Style tab → Ensure "Plot N value" enabled
```

---

### Donchian lines don't match breakout signals

**Explanation:**

Donchian breakout = Close > 55-bar high **of previous bar**

The blue line you see on current bar is the current 55-high (includes today).

**Correct interpretation:**
- Yesterday's 55-high: $180.00
- Today's close: $180.50
- Breakout: YES (close > yesterday's 55-high)
- Today's 55-high now: $180.50 (updated)

The script uses `donHiPrev` and `donLoPrev` (previous bar's Donchian level) to avoid lookahead bias.

---

### Script shows errors in Pine Editor

**Common errors:**

1. **Syntax error:**
   - Copy script again from reference/Ed-Seykota.pine
   - Don't modify code unless you know Pine Script

2. **Version mismatch:**
   - Script uses `//@version=6`
   - TradingView occasionally updates Pine Script version
   - Use latest script from TF-Engine installation

3. **Symbol not supported:**
   - Some tickers don't have sufficient history (< 55 bars)
   - Script requires at least 55 bars of data

---

### Free account limitations

**TradingView Free limitations:**

1. **3 indicators per chart:**
   - Ed-Seykota = 1 indicator
   - Leaves room for 2 more (e.g., volume, EMA)

2. **1 saved chart layout:**
   - Save your primary "TF-Engine Daily" layout
   - Sacrifice multiple layouts unless you upgrade

3. **Delayed data for some exchanges:**
   - US stocks: Real-time (free)
   - Forex, crypto: Real-time (free)
   - Some international exchanges: 15-minute delay

4. **Basic alerts:**
   - Limited alert types
   - Pro has advanced alert options

**Workaround:**
- Free account is sufficient for TF-Engine use
- Upgrade to Pro ($14.95/mo) if you want more features

---

## Next Steps

**After setting up TradingView:**

1. **Practice workflow:**
   - Open 5-10 tickers
   - Verify breakouts using Ed-Seykota script
   - Note N values
   - Return to TF-Engine and complete checklist

2. **Read:**
   - [USER_GUIDE.md](USER_GUIDE.md) - Daily workflow section
   - [QUICK_START.md](QUICK_START.md) - 10-minute overview

3. **Start trading systematically!**

---

## Additional Resources

**TradingView Help:**
- https://www.tradingview.com/support/
- Pine Script docs: https://www.tradingview.com/pine-script-docs/

**Ed Seykota / Turtle Traders:**
- Ed Seykota: "The Trading Tribe" book
- Turtle Traders: "The Complete TurtleTrader" by Michael Covel
- Van Tharp: "Trade Your Way to Financial Freedom"

**TF-Engine Docs:**
- [USER_GUIDE.md](USER_GUIDE.md) - Comprehensive guide
- [FAQ.md](FAQ.md) - Common questions
- [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Technical issues

---

**Version:** 1.0.0
**Last Updated:** 2025-10-29
**Remember:** Verify every signal manually. TradingView shows you the breakout, TF-Engine enforces discipline.
