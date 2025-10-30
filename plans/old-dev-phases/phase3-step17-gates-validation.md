# Phase 3 - Step 17: 5 Gates Validation

**Phase:** 3 - Heat Check & Trade Entry
**Step:** 17 of 28, 3 of 5 (Phase 3)
**Duration:** 2-3 days
**Dependencies:** Step 16 (Trade Entry Screen)

---

## Objectives

Implement the **heart of the discipline enforcement system**: the 5 gates check. All gates must pass for GO decision.

---

## The 5 Gates (from overview-plan.md)

**Gate 1: Banner Status**
- Current banner must be GREEN
- All required checklist items confirmed
- Quality score ≥ threshold

**Gate 2: Impulse Brake (2-minute cool-off)**
- Last evaluation timestamp must be ≥ 2 minutes ago
- Backend validates elapsed time

**Gate 3: Cooldown Status**
- Ticker not on cooldown (from recent loss)
- Bucket not on cooldown

**Gate 4: Heat Caps**
- Portfolio heat < 4% cap
- Bucket heat < 1.5% cap

**Gate 5: Sizing Completed**
- Position plan saved to database
- Risk calculated and stored

---

## Key Components

1. **Gate Check API Call:**
   - `POST /api/gates/check` with full trade data
   - Backend validates all 5 gates
   - Returns pass/fail for each gate

2. **Gate Results Display:**
   - Visual cards for each gate
   - Color-coded: ✓ GREEN (pass), ✗ RED (fail)
   - Show details for each gate:
     - Gate 1: "Banner: GREEN ✓" or "Banner: YELLOW ✗"
     - Gate 2: "Elapsed: 2m 15s ✓" or "Elapsed: 1m 45s ✗ (15s remaining)"
     - Gate 3: "No cooldowns ✓" or "AAPL on cooldown until 2025-11-05 ✗"
     - Gate 4: "Portfolio: 90.9% of cap ✓, Bucket: 99.7% of cap ✓" or "Bucket exceeds cap ✗"
     - Gate 5: "Position plan saved ✓" or "Sizing not completed ✗"

3. **Overall Result:**
   - Large display at bottom
   - "ALL GATES PASS ✓" (green gradient) or "GATES FAILED ✗" (red gradient)

4. **Button State Update:**
   - If all gates pass: Enable "SAVE GO DECISION"
   - If any gate fails: Keep "SAVE GO DECISION" disabled

---

## Implementation

**Update** `ui/src/routes/entry/+page.svelte`:
- "Run Final Gate Check" button click handler
- Call `POST /api/gates/check` with trade data from stores
- Display gate results in cards
- Update button states based on results

**Create** `ui/src/lib/components/entry/GateResults.svelte`:
- Display all 5 gates with pass/fail status
- Color-coded cards
- Details for each gate
- Overall result display

**Create** `ui/src/lib/api/gates.ts`:
- `checkGates(tradeData)` - Wraps `POST /api/gates/check`
- Returns gate results object

**Backend endpoint** (already exists):
- `POST /api/gates/check` - Validates all 5 gates
- Response: `{ gate1_pass, gate2_pass, gate3_pass, gate4_pass, gate5_pass, all_gates_pass, details }`

---

## Comprehensive Logging

**Log every gate check:**
- All gate results (pass/fail for each)
- Reasons for failures
- Timing data (elapsed since eval, heat levels, etc.)
- Whether user proceeds to save decision after check
- **Gate failure patterns** (which gates fail most often?)

**Feature evaluation:**
- Track which gates are frequent bottlenecks
- Identify if any gates are too strict (causing frustration)
- Measure time spent resolving gate failures
- Data-driven decisions to adjust or remove problematic gates

---

## Expected Outcome

User clicks "Run Final Gate Check". Backend validates:
- Gate 1: Banner GREEN ✓
- Gate 2: 2m 15s elapsed ✓
- Gate 3: No cooldowns ✓
- Gate 4: Heat caps OK ✓
- Gate 5: Sizing saved ✓

Display: "ALL GATES PASS ✓" (green gradient). "SAVE GO DECISION" button becomes enabled (green, clickable). User can now save GO decision.

**Failure example:** Gate 4 fails (bucket exceeds cap). Display: "GATES FAILED ✗" (red gradient). Show: "Gate 4: Heat Caps ✗ - Tech/Comm bucket would be $1,872 (124.8% of $1,500 cap)". "SAVE GO DECISION" remains disabled. User must return to heat check or sizing to resolve.

---

## Time Estimate

~8-12 hours (2-3 days)

---

## References

- [overview-plan.md - The 5 Hard Gates](../plans/overview-plan.md#the-5-hard-gates-enforced-by-backend)
- [overview-plan.md - User Workflow](../plans/overview-plan.md#user-workflow)

---

## Next Step

Proceed to: **[Phase 3 - Step 18: Decision Saving & Journaling](phase3-step18-decision-saving.md)**

---

**Status:** 📋 Ready for Execution
**Created:** 2025-10-29
