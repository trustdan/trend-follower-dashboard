# Phase 3 - Step 18: Decision Saving & Journaling

**Phase:** 3 - Heat Check & Trade Entry
**Step:** 18 of 28, 4 of 5 (Phase 3)
**Duration:** 1-2 days
**Dependencies:** Step 17 (Gates Validation)

---

## Objectives

Implement GO/NO-GO decision saving. Both paths are important - tracking what you DON'T trade is as valuable as tracking what you DO.

---

## Key Components

### 1. GO Decision (All Gates Pass)

**When:** User clicks "SAVE GO DECISION" (only enabled if all gates pass)

**Action:**
- Call `POST /api/decisions` with full trade record:
  - Timestamp
  - Decision type: "GO"
  - Ticker, sector, direction (LONG/SHORT)
  - Entry price, shares/contracts, max units
  - Initial stop, add-on prices
  - Structure (Stock, Call, Put, etc.)
  - Risk $ per unit, total max risk
  - Quality score
  - Journal note
  - Banner status (GREEN)
  - All 5 gate results (all PASS)

**Response:**
- Success notification: "✓ GO decision saved for AAPL"
- Ticker added to "Ready to Execute" list (Dashboard)
- Reset checklist/sizing for next trade
- Navigate to Dashboard or show success screen

**Optional (nice to have):**
- Confetti animation on success 🎉

### 2. NO-GO Decision (Any Reason)

**When:** User clicks "SAVE NO-GO DECISION" (always enabled)

**Action:**
- Show modal asking for rejection reason
- Modal fields:
  - Reason (textarea, required): "Why rejecting this trade?"
  - Category (dropdown, optional): "Heat cap", "Gate failed", "Changed mind", "Better opportunity", "Other"

**Save:**
- Call `POST /api/decisions` with:
  - Timestamp
  - Decision type: "NO-GO"
  - Ticker (if filled out)
  - Reason (user's text)
  - Category (if selected)
  - All available data (checklist, sizing, gate results if checked)

**Response:**
- Success notification: "✓ NO-GO decision logged for AAPL"
- Important: This is journaling, not failure - tracking rejections is valuable
- Reset checklist/sizing
- Navigate to Dashboard

---

## Implementation

**Update** `ui/src/routes/entry/+page.svelte`:
- "SAVE GO DECISION" click handler (when enabled)
- "SAVE NO-GO DECISION" click handler (always)
- NO-GO modal component

**Create** `ui/src/lib/components/entry/NoGoModal.svelte`:
- Modal overlay with backdrop
- Reason textarea (required)
- Category dropdown (optional)
- "Cancel" and "Save NO-GO" buttons
- Form validation (reason not empty)

**Create** `ui/src/lib/api/decisions.ts`:
- `saveDecision(data)` - Wraps `POST /api/decisions`
- Handles both GO and NO-GO types

**Backend endpoint** (already exists):
- `POST /api/decisions` - Saves decision record to database
- Table: decisions (id, timestamp, type, ticker, all trade data)

---

## Logging & Feature Evaluation

**Log all decisions:**
- GO decisions: full trade details, timing (how long from scan to decision?)
- NO-GO decisions: reasons, which gate failed (if applicable)
- **GO/NO-GO ratio**: Track percentage of trades that reach GO vs NO-GO
- **Rejection patterns**: Most common reasons for NO-GO
- **Success path timing**: Average time from candidate import to GO decision

**Feature evaluation:**
- If NO-GO >> GO, system may be too restrictive
- If specific gate causes many NO-GOs, may need adjustment
- Track user frustration indicators (abandoned workflows, repeated attempts)

---

## Expected Outcome

**GO Path:**
User with all gates passing clicks "SAVE GO DECISION". Backend saves complete record. Green success notification: "✓ GO decision saved for AAPL". Confetti animation (optional). Dashboard updates showing AAPL in "Ready to Execute". Checklist resets for next trade.

**NO-GO Path:**
User clicks "SAVE NO-GO DECISION". Modal appears: "Why are you rejecting this trade?" User types: "Portfolio heat at 95%, no capacity for new position". Selects category: "Heat cap". Clicks "Save NO-GO". Backend saves. Notification: "✓ NO-GO decision logged for AAPL". Checklist resets.

---

## Time Estimate

~4-6 hours (1-2 days)

---

## References

- [overview-plan.md - Decision Saving](../plans/overview-plan.md#user-workflow)
- [roadmap.md - Step 18](../plans/roadmap.md#step-18-decision-saving--journaling)

---

## Next Step

Proceed to: **[Phase 3 - Step 19: Integration Testing for Core Workflow](phase3-step19-integration-testing.md)**

---

**Status:** 📋 Ready for Execution
**Created:** 2025-10-29
