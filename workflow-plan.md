# Trend-Following Breakout Workflow (Finviz → TradingView → Execution)

**Goal:** Run a bias-minimized, mechanical process that surfaces liquid, trend-ready stocks from a broad universe using FINVIZ; confirm entries/exits with your Turtle-style TradingView strategy; size consistently; and keep portfolio heat controlled.

---

## Invariants (don’t change day-to-day)
- **Universe:** US-listed, liquid, large caps (price ≥ $20–50, avg. volume ≥ 1M shares). Prefer optionable, tight spreads when using options.
- **Entry logic:** Price breakout (e.g., 52‑week / multi‑month high) with price above the 50‑ and 200‑day SMAs.
- **Exit logic:** Managed in TradingView by your Donchian/ATR system (e.g., 10‑day exit / 2×N stop).
- **Risk:** Per-position risk = 0.25–0.75% of equity; portfolio heat cap = 2–4% across all open units.
- **Correlation buckets:** Cap risk per bucket (Tech/Comm, Cons., Financials, Industrials, Energy/Materials, Defensives/REITs) at ~1.0–1.5%.
- **Bias guardrails:** Follow the exact screener presets and the checklist below; no discretionary overrides on entry signals.

---

## Daily Schedule (15–25 minutes)

### A) After close (primary run)
1. **Run FINVIZ presets (in order):**
   1) **TF_BREAKOUT_LONG** → candidates for new long entries (fresh 52‑week highs).
   2) **TF_MOMENTUM_UPTREND** → replenishes the watchlist with strong, persistent uptrends.
   3) **TF_UNUSUAL_VOLUME** → breakout‑adjacent names with volume confirmation.
   4) **TF_BREAKDOWN_SHORT** *(optional)* → for hedges / short book if strategy allows.
2. **Export tickers** from each preset, **dedupe**, and tag by correlation bucket (Tech/Comm, Cons., Financials, Industrials, Energy/Materials, Defensives/REITs).
3. **TradingView validation:**
   - Load the **Seykota/Turtle** strategy on **Daily** (or Weekly for defensive names).
   - Confirm a valid **entry** (breakout/stop rules, no earnings within 3–5 trading days if you avoid earnings risk).
4. **Unit sizing:** Compute shares/contracts from risk % and ATR/N. Log in journal.
5. **Orders/alerts:** Place stop orders/alerts (once‑per‑bar close).

### B) Midday quick scan (optional, 5 minutes)
- Run **TF_BREAKOUT_LONG** and **TF_UNUSUAL_VOLUME** only, to catch intraday signals for end‑of‑day entries. No discretionary intraday chasing; entries still occur on close or next open per your plan.

### C) Weekly maintenance (Friday after close)
- Review **bucket heat** and **rule adherence** (≥90% following). If a bucket showed **2 stop‑outs within 20–30 bars**, apply a **10‑bar cool‑down** for *new* entries in that bucket next week.

---

## Decision Checklist (printable)

1. **Preset & signal present?** (ticker came from one of the FINVIZ presets)
2. **Trend filter:** Price above 50‑ & 200‑DMA (unless using Weekly for defensives).
3. **Liquidity:** Avg vol ≥ 1M, price ≥ $20–50 (per your risk preference).
4. **Portfolio heat:** New risk fits **per‑position** and **per‑bucket** caps.
5. **TradingView confirms:** Entry/exit rules & position size computed from ATR/N.
6. **Earnings/corporate events checked:** Within 3–5 trading days? (skip or reduce size if you avoid earnings)
7. **Order placement:** Initial stop = 2×N (or system default); alerts set for adds/exits.
8. **Journal entry created** (reason, unit size, stop, adds plan, bucket).

---

## Watchlist Hygiene
- Keep **two lists**: **Ready** (passed presets but no signal yet) and **Triggered** (signal confirmed today).
- Remove symbols violating liquidity or trend filters; archive exited names with brief notes (what worked/failed).

---

## Notes on Bias Minimization
- Work **top‑down via presets**—don’t browse charts first.
- Use the **same sort order** each time (see presets file) and don’t reshuffle mid‑scan.
- Hide company headlines while screening; read news **after** the signal.
- If in doubt, let the **risk & heat rules** decide—never “like/don’t like” the company.

---

## Hand‑off to Options (if used)
- Prefer ~0.60–0.70 Δ calls or call verticals 60–90 DTE; roll around ~21 DTE.
- For choppier groups (Utilities/REITs), consider stock or smaller verticals.
- Pyramids: add each +0.5×N from last add; exit/roll when price exit prints.

---

## Files & Links
- **Presets:** See `finviz-presets.md` (drop‑in query strings).
- **Journal template:** simple CSV/Sheets with columns: Date, Ticker, Bucket, Preset, Risk%, N/ATR, Shares/Contracts, Stop, Adds, Exit rule, Reason notes, Adherence (Y/N).

---

# Position Sizing, Pyramiding, and Option Rolls (Bias-Minimized)

**Purpose.** Convert signals into consistent bet sizes, add only to winners, and keep portfolio heat contained.

### A. Definitions
- **E** = Account equity
- **r** = per-position risk % (daily system: 0.5% typical; weekly: 0.25%)
- **R** = dollars at risk per new unit = `E × r`
- **N** = ATR(20) or your system’s ATR proxy (points)
- **K** = stop multiple (default `2×N`)
- **H** = portfolio heat cap (total open risk), e.g., 2–4% of equity
- **Buckets** = correlation groups (Tech/Comm, Cons., Financials, Industrials, Energy/Materials, Defensives/REITs); cap bucket heat at ~1.0–1.5%

---

### B. Share/Stock Sizing (mechanical)
1) **Stop distance:** `StopDist = K × N`
2) **Shares:** `Qty_shares = floor(R / StopDist)`
3) **Initial stop:** entry − `StopDist` (long) or entry + `StopDist` (short)
4) **Heat check:** sum of all open `R` across positions ≤ **H**; per-bucket sum ≤ bucket cap.

*Example.* E=$100,000, r=0.5% ⇒ R=$500. If N=3.00 and K=2 ⇒ StopDist=6.00.  
`Qty_shares = floor(500 / 6) = 83` shares.

---

### C. Options Sizing (two consistent approaches—pick one and stick with it)

**Method 1 — Delta-ATR (closest to stock logic)**
- Target **delta exposure** equal to share sizing.
- `Qty_contracts = floor( R / (K × N × Δ × 100) )`
  - Δ = option delta (use 0.60–0.70 for new longs)
- Place a **price-based** exit (from your system), not a premium stop.

*Example.* R=$500, N=3, K=2, Δ=0.60 ⇒ Denominator = 2×3×0.6×100 = 360 → **1** contract.

**Method 2 — Max-loss (simpler for debit structures)**
- Long call: max loss = premium × 100  
- Debit call vertical: max loss = **net debit × 100**  
- `Qty_contracts = floor( R / MaxLoss_per_contract )`
- Still manage exits by **price** (system exit), not premium.

*Consistency rule:* Use the same method across the book so heat math is stable.

---

### D. Poker-Inspired Pyramiding (add only as the trade proves itself)
- **Add levels:** every **+0.5×N** from the last add (max 3–4 adds per symbol).
- **Weights (choose one):**
  - *Equal-chips (default):* `1R, 1R, 1R, 1R`
  - *Front-loaded:* `1.25R, 1.0R, 0.75R, 0.50R` (only if heat allows)
- **Never average down.** Adds occur *only* above prior cost (longs) or below (shorts).
- **Table-stakes discipline:** Before each add, confirm total heat ≤ **H** and bucket heat ≤ cap.
- **Optional risk-lock:** After the first profitable add, you may trail stops so that *net* worst-case open risk across the stack stays ≤ initial `1–1.5R` on that symbol.

---

### E. When to Roll an Option vs. Close & Re-Enter

**Roll (stay in the same trend, same signal) when ALL apply:**
1) **System still long/short** (no exit printed); trend intact.
2) **Time-based:** ≤ **21 DTE** (to avoid gamma/theta blow-up) **or**  
   **Delta-based:** Δ ≥ **0.75–0.80** (deep ITM; reduce gamma, refresh time).
3) **Liquidity OK:** tight spreads, ≥ 500–1,000 OI on target strikes.
4) **Action:**  
   - *Long call:* **roll up-and-out** to regain Δ ≈ 0.60–0.70 with +45–90 DTE.  
     (Sell current, buy further-dated strike near ATM so portfolio Δ stays in range.)  
   - *Debit call vertical:* if P/L ≥ **65–75% of max**, **close** and **open a new vertical** higher strikes and/or later DTE. (Vertical “rolls” are functionally close-and-reopen.)

**Close (don’t roll) when ANY apply:**
- **System exit** triggers (10-day exit or stop). *Flatten options first; re-enter only on a fresh system signal.*  
- **Reward:Risk on the contemplated roll < 1** (little room left compared to debit/time).  
- **Liquidity deteriorates** (widespread/low OI) or **event risk** you avoid (earnings within 3–5 sessions).  
- **Spread at ≥ 80–90% of max profit** (vertical): take it off; re-establish only if a new add/entry triggers.

**Re-Enter (instead of rolling) when:**
- You closed due to **system exit** or due to **event/liquidity constraints**, and later a **fresh breakout** signal appears. Start over with base Δ (0.60–0.70) and DTE (60–90).

---

### F. Quick Checklist (every order)
- Position came from the **preset → signal** pipeline; TradingView strategy confirms entry/exit rules.
- **R calculated** and **Qty** computed (shares or contracts) via chosen method.
- **Heat check**: portfolio ≤ **H**, bucket ≤ cap.
- **Options only:** Δ in range, DTE 60–90 for new entries; roll/close logic applied if applicable.
- **Orders & alerts** placed; **journal** updated (reason, R, N, Δ/DTE, stop, next add price).

---



# FINVIZ Presets (drop‑in)

Paste each query string at the end of: `https://finviz.com/screener.ashx?`  
All use `v=211&p=d` (screener view + daily charts). You can change the sort `o=` or add a sector filter if desired.

---

## 1) TF_BREAKOUT_LONG  — Fresh long breakouts with trend filter

```
v=211&p=d&s=ta_newhigh&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pa,ta_sma200_pa&o=-relativevolume
```

**Why:** New 52‑week highs + price above 50/200‑DMA in liquid, large caps; sorted by **Relative Volume (desc)** to surface conviction days.

---

## 2) TF_MOMENTUM_UPTREND — Strong uptrends for watchlist refresh

```
v=211&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pa,ta_sma200_pa&dr=y1&o=-marketcap
```

**Why:** Matches your prior “up” idea but keeps the 200‑DMA filter to avoid false positives; large caps first.

---

## 3) TF_UNUSUAL_VOLUME — Volume confirmation / early movers

```
v=211&p=d&s=ta_unusualvolume&f=cap_largeover,sh_price_o20,ta_sma50_pa,ta_sma200_pa&o=-relativevolume
```

**Why:** Catches names with **unusual volume** inside existing uptrends; useful feeder list for imminent breakouts.

---

## 4) TF_BREAKDOWN_SHORT — Fresh breakdowns (optional short/hedge)

```
v=211&p=d&s=ta_newlow&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pb,ta_sma200_pb&o=-relativevolume
```

**Why:** Symmetric to #1 for downside; use cautiously or as a hedge framework.

---

## 5) TF_MOMENTUM_DOWNTREND — Persistent downtrends (optional)

```
v=211&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pb,ta_sma200_pb&dr=y1&o=-marketcap
```

**Why:** Mirrors your “down” preset but adds the 200‑DMA filter for trend integrity.

---

### Optional tweaks

- Raise `sh_price_o20` to `sh_price_o50` if you want only higher‑priced names.
- Sort alternatives: `o=-marketcap` (size first), `o=-volume`, or `o=-change` (for strongest day moves).
- To scan a **single sector**, append `&f=... ,sector` from the Finviz UI (or pick sector in the dropdown) after loading the preset.

---

## How to save as FINVIZ presets

1. Open each string (after `screener.ashx?`), confirm filters, and **Save Screener** (top right) with the given preset name.
2. Repeat for all five presets. During daily runs, open “My Presets” and click them in the order shown in the workflow.
