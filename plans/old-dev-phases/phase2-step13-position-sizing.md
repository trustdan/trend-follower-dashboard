# Phase 2 - Step 13: Position Sizing Calculator

**Phase:** 2 - Checklist & Position Sizing
**Step:** 13 of 28, 4 of 5 (Phase 2)
**Duration:** 2-3 days
**Dependencies:** Step 12 (Quality Scoring)

---

## Objectives

Build Position Sizing screen that calculates shares/contracts using Van Tharp method from backend domain logic.

---

## Key Components

1. **Pre-filled Form** - ticker, entry, ATR from checklist
2. **Method Selector** - Dropdown: "stock", "opt-delta-atr", "opt-contracts"
3. **Settings Display** - Show equity, risk% per unit, max units from backend
4. **Calculate Button** - Calls `POST /api/size/calculate`
5. **Results Display:**
   - Shares/contracts per unit
   - Risk $ per unit
   - Initial stop price
   - Add-on schedule (visual display of 4 levels)
   - Actual risk verification

6. **Warnings** - Concentration risk if position > 25% of equity
7. **Save Button** - Saves position plan to database

---

## Implementation

**Create** `ui/src/routes/sizing/+page.svelte`:
- Form pre-filled with checklist data (reactive from store)
- Method dropdown with conditional fields (delta for options, etc.)
- "Calculate Position Size" button
- Results card with gradient border showing sizing breakdown
- Add-on schedule visual (4 units with prices)
- "Save Position Plan" button

**Create** `ui/src/lib/api/sizing.ts`:
- `calculateSize(params)` - Wraps `POST /api/size/calculate`
- `savePositionPlan(plan)` - Saves to database

**Backend endpoints** (already exist from domain logic):
- `POST /api/size/calculate` - Van Tharp calculation
- Response includes: shares, risk$, stop, add-on levels, warnings

---

## Expected Outcome

User enters AAPL @ $180.50, ATR = $2.35. Clicks Calculate. Backend returns: 159 shares/unit, $747 risk, stop $175.80, add-ons at [182.68, 184.85, 187.03]. Results display prominently. User saves plan.

---

## Time Estimate

~8-12 hours (2-3 days)

---

**Next:** [Step 14: 2-Minute Cool-Off Timer](phase2-step14-cooloff-timer.md)
