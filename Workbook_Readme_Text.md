# Readme — Options Calendar + Seykota/Turtle Checklist

Welcome! This workbook gives you a **rolling 10‑week options calendar** and a **minimal‑friction trend‑following checklist** based on Seykota/Turtle principles. Use the **Trades** sheet to schedule positions by dates; use the **Checklist** to enforce entry discipline and compute simple contract counts.

---

## Sheets Overview

- **Readme** (this page): how everything works.
- **Trades**: enter one row per idea — `Symbol, Sector, StartDate, EndDate, Active, Strategy, Notes`.
- **Calendar**: auto‑populates symbols by **Sector × Week**, covering **2 weeks back + 8 forward** (weeks start Monday).
- **Summary**: active-count and sector‑coverage charts per week.
- **Checklist**: interactive pre‑trade flow with **3‑state banner**, **contracts helper**, **preset toggles**, and **Save Decision** button.
- **Decisions**: populated by “Save Decision”; exports your ticket and helper outputs.
- **Lists** (hidden): sector list & validations — edit if you need more/fewer sectors.

---

## Quick Start

1. **Enable macros** if Excel asks.  
2. (If the Checklist sheet is missing) run `BuildInteractiveChecklist` then `ApplyChecklistFixes`.  
3. (If the Calendar is blank/“repaired”) run `ForceRepairAndRefresh` in `CalendarModule`.  
4. Enter a few rows on **Trades**; confirm they appear on **Calendar**.  
5. Use **Checklist** to qualify a trade, then **Add To Trades** and **Save Decision**.

---

## Trades → Calendar (how the rolling window works)

- **Header** `Calendar!B1:K1` holds **Mondays**; `B2:K2` shows **“mmm d – mmm d”**.
- A sector/week cell lists all symbols whose date range **overlaps** that week and `Active="Yes"`.
- The totals row counts active trades per week; **Summary** mirrors it with charts.
- To change the 2‑back/8‑forward window, edit `CalendarModule` (WEEKS_BACK / WEEKS_TOTAL).

---

## Checklist (Seykota / Turtle core)

### Required groups
- **Signal (SIG_REQ)**: 55‑bar breakout (long > 55‑high / short < 55‑low).  
- **Risk/Size (RISK_REQ)**: risk **% of equity per unit** via **2×N** stop; add **every 0.5×N** up to **Max Units**.  
- **Options (OPT_REQ)**: **60–90 DTE**, **roll/close ~21 DTE**, and **liquidity required** (bid‑ask <10% mid; OI >100).  
- **Exits (EXIT_REQ)**: exit by **10‑bar opposite Donchian** **or** **closer of 2×N**.  
- **Behavior (BEHAV_REQ)**: **2‑minute cool‑off** and **no intraday overrides**.

### Optional “QUAL” items (score)
These curb impulsivity without blocking trades (e.g., regime check, no chase >2N, earnings blackout). **Quality threshold** defaults to **3.0** (cell `I6`).

### Banner (J5)
- **DO NOT TRADE**: any required fails.  
- **CAUTION**: required pass but score < threshold **or** theta rule fails.  
- **OK TO TRADE**: required pass + theta OK + score ≥ threshold.  
_Banner color is set by code so it updates on every change._

### Buttons (column L)
- **Reset** — uncheck all boxes.  
- **Add To Trades** — append a row to **Trades** (StartDate = Today; EndDate = Today + DTE).  
- **Recalculate** — force refresh if needed.  
- **Apply Preset** — stamps **System‑1 (20/10)** or **System‑2 (55/10)** as a journal note.  
- **Save Decision** — logs ticket inputs + helper outputs to **Decisions**.

---

## Sizing & Contracts Helper

Fill in (left side):
- `Account Size`, `N (ATR)`, `Stop multiple (N)`, `Add step (N)`, `Risk % per unit`, `Max units`  
- `Underlier @ entry`, `Option delta`, `Debit per vertical`, `Structure (Single Call / Debit Vertical)`

Computed:
- **Per‑unit risk $** = `Account × Risk%`  
- **Single‑call contracts** ≈ `unitRisk$ / (delta × 2N × 100)`  
- **Debit‑vertical spreads** ≈ `unitRisk$ / debit`  
- **Add prices**: `Entry + (k × AddStepN × N)` for k = 1..(MaxUnits−1)

---

## Presets (journaling helpers)
- **System‑1 (20/10)** and **System‑2 (55/10)**: select in `C18`, then hit **Apply Preset**.  
This does **not** overwrite your signal/exit rules; it documents which lens you used.

---

## Decisions sheet (what gets logged)
Timestamp, Symbol, Sector, Strategy, Score, Threshold, ThetaOK, Status, Notes, N, StopN, AddStepN, Risk%, MaxUnits, DTE, Structure, Delta, Debit, Entry, CallContracts, VerticalSpreads, Preset.

---

## Troubleshooting

- **Banner color doesn’t change** → press **Recalculate** (re‑enables Excel events) or run **ApplyChecklistFixes** (wipes CF on J5).  
- **Buttons overlap** → run **FixChecklistUI** (or **ApplyChecklistFixes**; buttons live in column **L**).  
- **Calendar repaired/blank** → run **ForceRepairAndRefresh** (no structured refs).  
- **Dynamic arrays not supported** → calendar won’t compute; use the macro‑rebuilt version only.

---

## Customize
- Edit sectors on **Lists** (unhide it).
- Change quality threshold in `I6` (e.g., 4–6 for stricter “CAUTION”).
- Adjust risk% per unit at `C11` (e.g., 0.5–1.0%).

**Requirements:** Excel 365 / 2021 recommended; macros enabled.
