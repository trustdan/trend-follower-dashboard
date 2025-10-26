Here’s a compact, portable guide you can paste into a fresh chat to pick up exactly where we left off.

------

# Seykota/Turtle Trend-Following — Options Execution + TradingView Strategy Guide

*Last updated: v2.1 + Date Range + “Flat at FROM” + Alerts*

## 1) Intent & Philosophy

- **Goal:** Ed Seykota/Turtles style trend-following that curbs impulsivity and lets a few big trends pay for many small losses.
- **Keep it simple:** buy strength/sell weakness; size by volatility (ATR “N”); add to winners, never to losers; cut losses mechanically.
- **Options expression (when trading the signals):**
  - 60–90 DTE; roll or close ~**21 DTE**.
  - IVR <20 → **single ITM** (Δ≈0.6–0.7).
     IVR ≥20 → **debit vertical** (buy Δ≈0.6 / sell Δ≈0.3–0.4).
  - Avoid holding long premium through earnings (unless defined-risk and intentional).

------

## 2) The Strategy We Built (TradingView / Pine v6)

**Name:** *Seykota / Turtle Core v2.1 + Date Range*

### Core rules (System-2 defaults)

- **Entry (breakout):** Close > **Donchian High** of prior **55** bars (long).
   Close < **Donchian Low** of prior **55** bars (short).
- **Exit (trend stop):** Opposite Donchian **10** (long exits on 10-low; short exits on 10-high).
- **Protective stop:** **2×N** from entry (N = ATR(20)) — choose the closer of 10-channel vs 2N stop.
- **Pyramiding:** add a unit every **0.5×N**, up to **max 4 units**.
- **Position sizing (per unit):** risk **riskPct%** of equity per unit, based on N at entry.

### Important inputs/toggles

- `entryLen` **55**, `exitLen` **10**, `nLen` **20**, `stopN` **2.0**, `addStepN` **0.5**, `maxUnits` **4**, `riskPct` **0.5–1.0** (daily) or **0.1–0.25** (intraday).
- `allowLong`, `allowShort` — enable both for symmetry.
- **Regime filter (optional):** `useMarket=true` with `marketSym=SPY`, `marketTF=D`, `marketLen=200`
   → Longs only when SPY > 200SMA, shorts only when SPY < 200SMA.
- **Time exit (optional):** `timeExitMode = None / Close / Roll`, `timeExitBars` default **60**.
- **Date range filter:**
   `fromDate`, `toDate` inputs + `flatAtFrom=true` to **start flat** at first in-range bar.
- **Alerts:** built-in conditions for **Long/Short Entry**, **Long/Short Add**, **Long/Short Exit**, **Time Exit**.

> **Tip:** Use **Once per bar close** for cleaner alerts; **Once per bar** if you want early pings.

------

## 3) Anti-Impulsivity Checklist (print this)

**Regime & Name (keep minimal to avoid missing trends)**

- ☐ (Optional) Market regime OK (SPY above/below 200SMA for long/short).
- ☐ Liquidity OK for your options (tight spreads, OI > 100 on target strikes).

**Entry**

- ☐ Breakout: price closes above 55-high (long) / below 55-low (short).
- ☐ Not chasing a giant gap? (Use limit/stop orders per plan.)

**Sizing & Adds**

- ☐ Risk per unit = **riskPct%** of equity from **2×N** stop distance.
- ☐ Adds every **0.5×N**, max **4** units.

**Exit & Roll**

- ☐ Exit on 10-day opposite Donchian **OR** 2×N stop (closer of the two).
- ☐ If using options: roll/close ~**21 DTE**; profit-take early on verticals at **50–75%** of max.

**Behavior**

- ☐ 2-minute cooldown before sending orders.
- ☐ If 2 losing days in a row → **no new trades until next close**.
- ☐ Journal one sentence: “Why this trade, now?”

------

## 4) How to Use It (step-by-step)

### A) Add the strategy to a chart

1. Paste the *v2.1 + Date Range* script into Pine Editor, **Add to chart**.
2. Set inputs (defaults are solid): `entryLen=55`, `exitLen=10`, `nLen=20`,
    `stopN=2.0`, `addStepN=0.5`, `maxUnits=4`, `riskPct=0.5–1.0` (daily).
3. **Date window:** set `fromDate` / `toDate`; leave `flatAtFrom=true`.

### B) Create alerts (entries, adds, exits)

- **Create Alert → Condition →** your script name → **Long Entry** (etc).
- Repeat for **Long Exit** on each **open long** (and **Short Exit** for shorts).
- Recommended: **Once per bar close**, push + email notifications.

### C) Backtest without Premium

- Use the **date-range inputs** (this constrains tests to the bars loaded by your plan).
- `flatAtFrom=true` so P&L & trades begin **inside** the window (no carry-in positions).

------

## 5) Presets that generally test well

### Daily swing/position (closest to classic Turtles)

- `entryLen=55`, `exitLen=10`, `nLen=20`, `stopN=2.0`, `addStepN=0.5`, `maxUnits=4`, `riskPct=0.5–1.0`.
- `timeExitMode=None`.
- `useMarket` optional; leave **off** to catch more trends.

### Intraday (liquid names only)

- `entryLen=20–55`, `exitLen=10–20`, `nLen=20`, `stopN=2.0`, `addStepN=0.5–1.0`, `riskPct=0.10–0.25`.
- Keep `useMarket=false`; `timeExitMode=None`.

------

## 6) Options Translation (from price signals)

- **When a long breakout prints:** prefer **ITM call** (Δ≈0.6–0.7) at **60–90 DTE**; or **debit vertical** if IVR ≥20.
- **Pyramids:** add contracts proportionally when price has moved **+0.5×N** from last add.
- **Exits:** close/roll when the **price exit** triggers; if single call loses **50%** before price exit, you may cut early (pick a rule and stick to it).
- **Always roll/close** around **21 DTE** to avoid late-cycle theta/gamma.

------

## 7) Troubleshooting (Pine)

- **“end of line without line continuation”** → you ended a function call line with a comma; remove it or continue args on the next line.
- **Undeclared variables inside functions** → declare locals (`float sh = …`) or assign outside.
- **No alert options showing** → ensure `alertcondition(...)` exists and the script is re-added to chart.

------

## 8) Roadmap / Nice-to-haves

- **System-1/System-2 toggle** (20/10 vs 55/10).
- **“Enter on close confirmation”** toggle (enter next bar after confirmed breakout).
- **Earnings blackout** input (skip entries X days before earnings).
- **Portfolio heat display** (sum of per-unit risks/open units).

------

## 9) Journal Ticket (copy/paste each trade)

```
Ticker / TF:
Entry type: 55-breakout  |  Units planned/max:
N at entry: ____ |  2×N stop distance: ____ |  Unit risk %: ____
Adds: every 0.5×N to max __ units
Exit: 10-day opposite Donchian OR closer of 2×N stop
Options plan: DTE __ | ITM Δ~0.6 (or debit vertical) | Roll at ~21 DTE
Cool-down observed? Y/N | Reason in one sentence:
```

------

## 10) One-screen defaults (good starting point)

- **Daily TF:** 55/10 channel, N=20, 2×N stop, add 0.5N, max 4 units, risk 0.75% per unit, time exit **None**, regime **Off**.
- **Alerts:** Long/Short Entry + Exit on watched symbols; Add-on alerts optional.

------

### Appendix: Minimal code bits you might reference

- **Date inputs (no confirmation):**

```pinescript
fromDate = input.time(timestamp("2022-01-01T00:00:00"), "Backtest FROM", confirm=false)
toDate   = input.time(timestamp("2024-12-31T23:59:59"), "Backtest TO",   confirm=false)
inRange  = time >= fromDate and time <= toDate
```

- **Start flat at window:**

```pinescript
flatAtFrom    = input.bool(true, "Force FLAT at range start?")
isRangeStart  = (time[1] < fromDate) and (time >= fromDate)
if flatAtFrom and isRangeStart and strategy.position_size != 0
    strategy.close_all(comment="Flat at FROM")
```

- **Alert conditions (so they appear in UI):**

```pinescript
alertcondition(longEntryCond,  title="Long Entry",  message="TurtleCore: LONG entry on {{ticker}} ({{interval}})")
alertcondition(shortEntryCond, title="Short Entry", message="TurtleCore: SHORT entry on {{ticker}} ({{interval}})")
alertcondition(longAddCond,    title="Long Add",    message="TurtleCore: LONG add-on on {{ticker}} ({{interval}})")
alertcondition(shortAddCond,   title="Short Add",   message="TurtleCore: SHORT add-on on {{ticker}} ({{interval}})")
alertcondition(longExitCond,   title="Long Exit",   message="TurtleCore: LONG EXIT on {{ticker}} ({{interval}})")
alertcondition(shortExitCond,  title="Short Exit",  message="TurtleCore: SHORT EXIT on {{ticker}} ({{interval}})")
alertcondition(timeExitCond,   title="Time Exit",   message="TurtleCore: TIME EXIT on {{ticker}} ({{interval}})")
```

------

If you paste this guide into a new chat, I can immediately:

- add the **System-1/System-2** toggle,
- wire in an **earnings blackout**, and
- generate a **one-page printable checklist** or a **Google Sheet** that converts “per-unit price risk” into **contracts** for your options.





// This Pine Script® code is subject to the terms of the Mozilla Public License 2.0 at https://mozilla.org/MPL/2.0/

// © danieltuckerrust

//**@version=**6

strategy("Seykota / Turtle Core v2.1 + Date Range",

   overlay=true,

   initial_capital=100000,

   commission_type=strategy.commission.percent, commission_value=0.005,

   default_qty_type=strategy.fixed, default_qty_value=0,

   pyramiding=10,

   max_bars_back=5000)

//================ Inputs

allowLong  = input.bool(true,  "Allow LONGs?")

allowShort  = input.bool(true,  "Allow SHORTs?")

entryLen   = input.int(55,   "Donchian ENTRY lookback (System-2)", minval=10)

exitLen   = input.int(10,   "Donchian EXIT lookback",       minval=5)

nLen     = input.int(20,   "N = ATR length",           minval=5)

stopN    = input.float(2.0,  "Initial stop (in N)",         minval=0.5, step=0.25)

addStepN   = input.float(0.5,  "Add every X * N",           minval=0.25, step=0.25)

maxUnits   = input.int(4,    "Max units (incl. initial)",      minval=1, maxval=10)

riskPct   = input.float(1.0,  "Risk % of equity PER UNIT",      minval=0.1, maxval=5, step=0.1)

useMarket  = input.bool(false, "Use market regime filter?")

marketSym  = input.symbol("SPY", "Market symbol for regime (if used)")

marketTF   = input.timeframe("D", "Regime timeframe (if used)")

marketLen  = input.int(200,   "Market MA length (if used)",     minval=50)

timeExitMode= input.string("None","Time exit", options=["None","Close","Roll"])

timeExitBars= input.int(60,   "Time exit bars (if used)",      minval=5)

minVol    = input.int(0,    "Min 20-bar avg volume (chart TF)",  minval=0)

showSignals = input.bool(true,  "Plot signals & stops?")

plotDon   = input.bool(true,  "Plot Donchian bands?")

// === Date-range filter (local chart timezone) — no trailing commas

fromDate = input.time(defval=timestamp("2022-01-01T00:00:00"), title="Backtest FROM (yyyy-mm-dd)", confirm=true)

toDate  = input.time(defval=timestamp("2099-12-31T23:59:59"), title="Backtest TO (yyyy-mm-dd)",  confirm=true)

//================ Helpers

sharesForUnit(_equity, _Nentry) =>

  riskDollars  = _equity * (riskPct/100.0)

  perShareRisk  = math.max(stopN * _Nentry, syminfo.mintick)

  math.max(1, math.floor(riskDollars / perShareRisk))

// Only trade inside the selected date window

inRange = (time >= fromDate) and (time <= toDate)

// Force-flat option at the first in-range bar

flatAtFrom = input.bool(true, "Force FLAT at range start?")

isRangeStart = (time[1] < fromDate) and (time >= fromDate)

// If a prior position exists when the window begins, close it so P&L starts clean

if flatAtFrom and isRangeStart and strategy.position_size != 0

  strategy.close_all(comment="Flat at FROM")

//================ Core calcs (chart timeframe)

N     = ta.atr(nLen)

volMA   = ta.sma(volume, 20)

liqOK   = volMA >= minVol

// Donchian levels

donHi   = ta.highest(high, entryLen)

donLo   = ta.lowest(low,  entryLen)

donHiPrev= donHi[1]

donLoPrev= donLo[1]

exitHiPrev = ta.highest(high, exitLen)[1]

exitLoPrev = ta.lowest(low,  exitLen)[1]

// Market regime (optional)

mClose = request.security(marketSym, marketTF, close)

mMA   = request.security(marketSym, marketTF, ta.sma(close, marketLen))

longRegOK  = not useMarket or (mClose > mMA)

shortRegOK = not useMarket or (mClose < mMA)

//================ Signals (no lookahead)

longBreak  = allowLong  and liqOK and longRegOK  and (close > donHiPrev)

shortBreak = allowShort and liqOK and shortRegOK and (close < donLoPrev)

//================ Position state

var **float** N_entry    = na

var **float** lastAddLong  = na

var **float** lastAddShort  = na

var **int**  units     = 0

var **int**  barsInPos   = 0

var **bool**  wantReenterL  = false

var **bool**  wantReenterS  = false

inPos   = strategy.position_size != 0

inLong   = strategy.position_size > 0

inShort  = strategy.position_size < 0

if inPos

  barsInPos += 1

else

  barsInPos   := 0

  units     := 0

  lastAddLong  := na

  lastAddShort := na

  N_entry    := na

  wantReenterL := false

  wantReenterS := false

//================ Entries (risk-based sizing, freeze N at entry) — gated by inRange

if inRange and not inPos and longBreak

  N_entry := N

  **float** sh = sharesForUnit(strategy.equity, N_entry)

  strategy.entry("L", strategy.long, qty=sh)

  units    := 1

  lastAddLong := close

if inRange and not inPos and shortBreak

  N_entry := N

  **float** sh = sharesForUnit(strategy.equity, N_entry)

  strategy.entry("S", strategy.short, qty=sh)

  units     := 1

  lastAddShort := close

//================ Add-on logic (every addStepN * N_entry)

canAddLong  = inRange and inLong  and units < maxUnits and close >= nz(lastAddLong)  + addStepN * N_entry

canAddShort = inRange and inShort and units < maxUnits and close <= nz(lastAddShort) - addStepN * N_entry

if canAddLong

  **float** sh = sharesForUnit(strategy.equity, N_entry)

  strategy.entry("L", strategy.long, qty=sh)

  units    += 1

  lastAddLong := close

if canAddShort

  **float** sh = sharesForUnit(strategy.equity, N_entry)

  strategy.entry("S", strategy.short, qty=sh)

  units     += 1

  lastAddShort := close

//================ Protective & Donchian exits — attach only inside range

var **float** stopL = na

var **float** stopS = na

if inRange and inLong

  initStopL = strategy.position_avg_price - stopN * N_entry

  stopL   := math.max(initStopL, exitLoPrev)

  strategy.exit("L-EXIT", from_entry="L", stop=stopL)

else

  stopL := na

if inRange and inShort

  initStopS = strategy.position_avg_price + stopN * N_entry

  stopS   := math.min(initStopS, exitHiPrev)

  strategy.exit("S-EXIT", from_entry="S", stop=stopS)

else

  stopS := na

//================ Optional time exit — evaluate only inside range

timeExitTrig = inRange and inPos and (timeExitMode != "None") and (barsInPos >= timeExitBars)

if timeExitTrig

  if inLong

​    if timeExitMode == "Close"

​      strategy.close("L", comment="Time Exit")

​    else if timeExitMode == "Roll"

​      wantReenterL := true

​      strategy.close("L", comment="Time Exit (Roll)")

  if inShort

​    if timeExitMode == "Close"

​      strategy.close("S", comment="Time Exit")

​    else if timeExitMode == "Roll"

​      wantReenterS := true

​      strategy.close("S", comment="Time Exit (Roll)")

// Re-enter after Roll only if still inside range and trend condition holds

if inRange and not inPos and wantReenterL and longRegOK and (close > donHiPrev)

  N_entry := N

  **float** sh = sharesForUnit(strategy.equity, N_entry)

  strategy.entry("L", strategy.long, qty=sh)

  units    := 1

  lastAddLong := close

  wantReenterL := false

if inRange and not inPos and wantReenterS and shortRegOK and (close < donLoPrev)

  N_entry := N

  **float** sh = sharesForUnit(strategy.equity, N_entry)

  strategy.entry("S", strategy.short, qty=sh)

  units     := 1

  lastAddShort := close

  wantReenterS := false

//================ Alerts (optional) — only fire inside date window

longEntryCond  = inRange and not inPos and longBreak

shortEntryCond  = inRange and not inPos and shortBreak

longAddCond   = inRange and canAddLong

shortAddCond   = inRange and canAddShort

longExitCond   = inRange and inLong  and not na(stopL) and close <= stopL

shortExitCond  = inRange and inShort and not na(stopS) and close >= stopS

timeExitCond   = inRange and timeExitTrig

alertcondition(longEntryCond,  title="Long Entry",  message="TurtleCore: LONG entry on {{ticker}} ({{interval}})")

alertcondition(shortEntryCond, title="Short Entry", message="TurtleCore: SHORT entry on {{ticker}} ({{interval}})")

alertcondition(longAddCond,   title="Long Add",   message="TurtleCore: LONG add-on on {{ticker}} ({{interval}})")

alertcondition(shortAddCond,  title="Short Add",  message="TurtleCore: SHORT add-on on {{ticker}} ({{interval}})")

alertcondition(longExitCond,  title="Long Exit",  message="TurtleCore: LONG EXIT on {{ticker}} ({{interval}})")

alertcondition(shortExitCond,  title="Short Exit",  message="TurtleCore: SHORT EXIT on {{ticker}} ({{interval}})")

alertcondition(timeExitCond,  title="Time Exit",  message="TurtleCore: TIME EXIT on {{ticker}} ({{interval}})")

//================ Plots

**float** dHiPlot = plotDon ? donHi : na

**float** dLoPlot = plotDon ? donLo : na

plot(dHiPlot, "Donchian High", color=color.new(color.blue, 40))

plot(dLoPlot, "Donchian Low",  color=color.new(color.blue, 40))

plotshape(showSignals and longBreak,  title="Long Breakout",  style=shape.triangleup,  location=location.belowbar, size=size.tiny,  color=color.new(color.green, 0), text="Long")

plotshape(showSignals and shortBreak, title="Short Breakout", style=shape.triangledown, location=location.abovebar, size=size.tiny,  color=color.new(color.red,  0), text="Short")

plot(showSignals and inLong  ? stopL : na, "Long Stop",  color=color.new(color.red,  0), style=plot.style_linebr)

plot(showSignals and inShort ? stopS : na, "Short Stop", color=color.new(color.red,  0), style=plot.style_linebr)