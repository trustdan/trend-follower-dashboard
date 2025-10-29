# Phase 3 - Step 15: Heat Check Screen

**Phase:** 3 - Heat Check & Trade Entry
**Step:** 15 of 28, 1 of 5 (Phase 3)
**Duration:** 2-3 days
**Dependencies:** Phase 2 complete (Checklist & Position Sizing)

---

## Objectives

Build Heat Check screen to validate proposed trade won't exceed portfolio or sector bucket caps. This is a **critical gate** - prevents over-concentration.

---

## Key Components

1. **Current Heat Display:**
   - Portfolio heat gauge (visual progress bar)
   - Shows: Current heat / Cap (percentage)
   - Example: "$2,890 / $4,000 (72.25%)"

2. **Sector Bucket Heat:**
   - Table showing all buckets with current heat
   - Highlight bucket for proposed trade

3. **Heat Check Button:**
   - "Check Heat for This Trade"
   - Calls `POST /api/heat/check` with ticker, risk amount, bucket

4. **Results Display:**
   - New portfolio heat (current + proposed)
   - Portfolio cap (equity × 4%)
   - Result: ✓ WITHIN CAP or ✗ EXCEEDS CAP
   - New bucket heat (current + proposed)
   - Bucket cap (equity × 1.5%)
   - Result: ✓ WITHIN CAP or ✗ EXCEEDS CAP

5. **RED Warning (if caps exceeded):**
   - Large gradient warning card
   - "⚠️ PORTFOLIO CAP EXCEEDED" or "⚠️ BUCKET CAP EXCEEDED"
   - Show exact overage amount
   - Suggestions:
     - "Reduce position size to X shares"
     - "Close existing position first"
     - "Choose different sector"

6. **Calculate Max Shares:**
   - Helper button: "Calculate Max Shares for Caps"
   - Returns largest position size that fits within caps

7. **Proceed Button:**
   - "Proceed to Trade Entry" button
   - Only enabled if heat check passed
   - Disabled with clear message if caps exceeded

---

## Implementation

**Create** `ui/src/routes/heat/+page.svelte`:
- Display current portfolio heat with visual gauge
- Display sector buckets table with heat levels
- Form with ticker, risk amount, bucket (pre-filled from checklist/sizing)
- "Check Heat" button
- Results display with conditional RED warning
- "Calculate Max Shares" helper
- "Proceed to Trade Entry" button (conditional)

**Create** `ui/src/lib/components/heat/HeatGauge.svelte`:
- Visual progress bar component
- Color-coded: green (<70%), yellow (70-90%), red (>90%)
- Shows percentage and dollar amounts

**Create** `ui/src/lib/components/heat/HeatWarning.svelte`:
- RED gradient warning card
- Display overage amount
- List suggestions to resolve

**Create** `ui/src/lib/api/heat.ts`:
- `checkHeat(params)` - Wraps `POST /api/heat/check`
- `calculateMaxShares(params)` - Helper for max position

**Backend endpoints** (already exist from domain logic):
- `POST /api/heat/check` - Validates portfolio and bucket caps
- Response: current heat, new heat, caps, pass/fail, overage amounts

---

## Expected Outcome

User enters AAPL trade with $747 risk in Tech/Comm bucket. Current portfolio: $2,890/$4,000. Current Tech/Comm: $1,400/$1,500. Check heat:
- New portfolio: $3,637/$4,000 ✓ WITHIN CAP
- New bucket: $2,147/$1,500 ✗ EXCEEDS CAP by $647

Display RED warning. User clicks "Calculate Max Shares" → system suggests reducing to 79 shares (risk = $371). User returns to sizing screen, adjusts, rechecks. New bucket: $1,496/$1,500 ✓ WITHIN CAP. "Proceed to Trade Entry" enabled.

---

## Comprehensive Logging

Log all heat checks:
- Current heat levels (portfolio, all buckets)
- Proposed trade details (ticker, risk, bucket)
- Calculation results (new heat, caps, pass/fail)
- Overage amounts (if failed)
- User actions (reduced position, closed existing, changed ticker)
- Time spent resolving heat issues

**Feature evaluation:** Track heat cap violations to identify if caps are too restrictive or causing workflow friction.

---

## Time Estimate

~8-12 hours (2-3 days)

---

## References

- [overview-plan.md - Heat Management](../plans/overview-plan.md#heat-management)
- [roadmap.md - Step 15](../plans/roadmap.md#step-15-heat-check-screen)

---

## Next Step

Proceed to: **[Phase 3 - Step 16: Trade Entry Screen & Summary](phase3-step16-trade-entry.md)**

---

**Status:** 📋 Ready for Execution
**Created:** 2025-10-29
