# Phase 3 - Step 16: Trade Entry Screen & Summary

**Phase:** 3 - Heat Check & Trade Entry
**Step:** 16 of 28, 2 of 5 (Phase 3)
**Duration:** 2 days
**Dependencies:** Step 15 (Heat Check)

---

## Objectives

Create Trade Entry screen - the culmination of the workflow where final GO/NO-GO decision is made.

---

## Key Components

1. **Trade Summary Card** (gradient border, prominent display):
   - Ticker & Direction (LONG/SHORT)
   - Entry price
   - Shares/contracts per unit
   - Max units (1-4)
   - Initial stop price (2×N)
   - Risk $ per unit
   - Total max risk (if all units added)
   - Add-on schedule (visual: 4 levels with prices)
   - Sector bucket
   - Exit plan (10-bar Donchian OR 2×N, whichever closer)
   - Quality score (X/4)

2. **Final Gate Check Section:**
   - Large prominent section titled "Final Gate Check"
   - "RUN FINAL GATE CHECK" button (gradient, large)
   - Space reserved for gate results (initially empty)

3. **Action Buttons** (bottom of screen):
   - "SAVE GO DECISION" (green gradient, large, bold) - **Initially disabled**
   - "SAVE NO-GO DECISION" (red gradient, large) - **Always enabled** (for journaling)

---

## Implementation

**Create** `ui/src/routes/entry/+page.svelte`:
- Trade summary card at top (read from stores: checklist, sizing)
- Final gate check section with button
- Placeholder for gate results (Step 17 will populate)
- Action buttons at bottom
- "SAVE GO DECISION" disabled until all gates pass
- "SAVE NO-GO DECISION" always enabled

**Create** `ui/src/lib/components/entry/TradeSummary.svelte`:
- Displays complete trade plan
- All fields from checklist + sizing
- Visual add-on schedule
- Gradient border styling

**Create** `ui/src/lib/components/entry/DecisionButtons.svelte`:
- GO button (green gradient, disabled state)
- NO-GO button (red gradient)
- Clear labeling and large size

---

## Expected Outcome

User navigates to Trade Entry screen. Sees comprehensive summary of AAPL trade: 79 shares/unit, $180.50 entry, $175.80 stop, $371 risk/unit, Tech/Comm bucket. "Run Final Gate Check" button prominent. GO button disabled (grayed out). NO-GO button enabled.

---

## Time Estimate

~6-8 hours (2 days)

---

**Next:** [Step 17: 5 Gates Validation](phase3-step17-gates-validation.md)

---

**Status:** 📋 Ready for Execution
**Created:** 2025-10-29
