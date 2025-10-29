# Phase 3 - Step 19: Integration Testing for Core Workflow

**Phase:** 3 - Heat Check & Trade Entry
**Step:** 19 of 28, 5 of 5 (Phase 3) - **FINAL STEP OF PHASE 3**
**Duration:** 1-2 days
**Dependencies:** Steps 15-18 (Heat Check, Trade Entry, Gates, Decision Saving)

---

## Objectives

Comprehensive integration testing of the complete workflow from FINVIZ scan through GO/NO-GO decision. Verify all components work together correctly. Fix any bugs discovered.

**This is the culmination of Phase 1-3 work.** After this step, the core discipline enforcement workflow is complete and functional.

---

## What to Test

### 1. Happy Path - All Gates Pass

**Complete workflow from start to finish:**

1. **Dashboard → FINVIZ Scan**
   - Click "Run FINVIZ Scan"
   - Verify loading state
   - Verify results table populates
   - Check sector distribution display

2. **Import Candidates**
   - Select 3-5 tickers from different sectors
   - Click "Import Selected"
   - Verify success notification
   - Verify candidates appear in Dashboard

3. **Checklist Evaluation**
   - Fill out form completely (ticker, entry, ATR, sector, structure)
   - Check all 5 required gates
   - Verify banner transitions: RED → YELLOW → GREEN
   - Check all 4 optional quality items
   - Verify quality score: 4/4
   - Verify banner: GREEN
   - Verify "Continue to Position Sizing" enabled

4. **Position Sizing**
   - Verify pre-filled data from checklist
   - Select method (stock/opt-delta-atr/opt-contracts)
   - Click "Calculate Position Size"
   - Verify results display (shares, risk, add-ons)
   - Verify visual add-on schedule
   - Click "Save Position Plan"
   - Verify success notification
   - Verify "Continue to Heat Check" enabled

5. **Heat Check**
   - Verify current portfolio heat display
   - Verify sector bucket table
   - Click "Check Heat for This Trade"
   - Verify results: WITHIN CAP (both portfolio and bucket)
   - Verify green checkmarks
   - Verify "Proceed to Trade Entry" enabled

6. **Trade Entry**
   - Verify trade summary card (all details correct)
   - Click "Run Final Gate Check"
   - Verify all 5 gates display as PASS ✓
   - Verify "ALL GATES PASS ✓" message (green gradient)
   - Verify "SAVE GO DECISION" button enabled

7. **Save GO Decision**
   - Click "SAVE GO DECISION"
   - Verify success notification
   - Optional: Verify confetti animation (if implemented)
   - Verify navigation to Dashboard
   - Verify ticker appears in "Ready to Execute" list

**Expected Result:** Complete workflow with zero friction. All gates pass. GO decision saved. User sees confetti and success message.

---

### 2. Gate Failure Scenarios

Test each gate failing individually to verify proper error handling and UI feedback.

**Scenario A: Gate 1 Fails (Banner Not GREEN)**

1. Fill out checklist but check only 4 of 5 required gates
2. Banner should be RED
3. Navigate to Trade Entry
4. Run gate check
5. **Expected:** Gate 1 shows FAIL ✗ with message "Banner: RED ✗"
6. **Expected:** "SAVE GO DECISION" remains disabled
7. **Expected:** "GATES FAILED ✗" message (red gradient)

**Scenario B: Gate 2 Fails (2-Minute Cool-Off Not Elapsed)**

1. Complete checklist (banner GREEN)
2. Immediately navigate through sizing → heat → entry
3. Run gate check **within 2 minutes** of checklist save
4. **Expected:** Gate 2 shows FAIL ✗ with message "Elapsed: 1m 45s ✗ (15s remaining)"
5. **Expected:** User must wait for timer to complete
6. **Expected:** After 2 minutes, re-run gate check → Gate 2 PASS ✓

**Scenario C: Gate 3 Fails (Ticker on Cooldown)**

1. Setup: Add AAPL to cooldown table (recent loss)
2. Attempt to create trade for AAPL
3. Run gate check
4. **Expected:** Gate 3 shows FAIL ✗ with message "AAPL on cooldown until 2025-11-05 ✗"
5. **Expected:** User must choose different ticker

**Scenario D: Gate 4 Fails (Heat Caps Exceeded)**

1. Setup: Portfolio at 95% of cap
2. Create trade with risk that exceeds cap
3. Heat check should fail
4. **Expected:** RED warning: "⚠️ PORTFOLIO CAP EXCEEDED"
5. **Expected:** Show exact overage amount
6. **Expected:** "Proceed to Trade Entry" disabled
7. User clicks "Calculate Max Shares"
8. Adjust position size to fit within cap
9. Re-run heat check → PASS ✓
10. Proceed to entry, run gate check
11. **Expected:** Gate 4 PASS ✓

**Scenario E: Gate 5 Fails (Sizing Not Completed)**

1. Complete checklist
2. Skip position sizing step (navigate directly to entry)
3. Run gate check
4. **Expected:** Gate 5 shows FAIL ✗ with message "Sizing not completed ✗"
5. **Expected:** User must return to sizing step

---

### 3. NO-GO Decision Path

**Test the "other" path - rejecting trades is journaling, not failure:**

1. Complete checklist (banner GREEN)
2. Navigate to entry
3. Click "SAVE NO-GO DECISION" (always enabled)
4. **Expected:** Modal appears with:
   - "Why are you rejecting this trade?" textarea (required)
   - Category dropdown (optional)
5. User types reason: "Portfolio heat at 95%, no capacity for new position"
6. User selects category: "Heat cap"
7. Click "Save NO-GO"
8. **Expected:** Success notification: "✓ NO-GO decision logged for AAPL"
9. **Expected:** Navigation to Dashboard
10. **Expected:** Checklist resets for next trade

**Verify database:** NO-GO decision record saved with timestamp, ticker, reason, category.

---

### 4. UI Behavior Testing

**Banner Transitions:**
- RED (0 required) → YELLOW (5 required, quality <3) → GREEN (5 required, quality ≥3)
- Verify pulse animation (0.5s) on each transition
- Verify gradient colors correct
- Verify text changes ("DO NOT TRADE" / "CAUTION" / "OK TO TRADE")

**Cool-Off Timer:**
- Starts on checklist save (banner GREEN)
- Countdown from 120s → 0s
- Persists across navigation (visit Dashboard, return to checklist, timer still counting)
- After 0s, timer shows "Ready ✓"
- Backend validates elapsed time on gate check

**Heat Gauges:**
- Progress bars color-coded: green (<70%), yellow (70-90%), red (>90%)
- Percentages displayed correctly
- Dollar amounts match calculations

**Gate Results Display:**
- Each gate shows ✓ or ✗ with color coding
- Details displayed for each gate
- Overall result card: "ALL GATES PASS ✓" (green) or "GATES FAILED ✗" (red)

**Button States:**
- "Continue to X" buttons enabled/disabled correctly
- "SAVE GO DECISION" only enabled when all gates pass
- "SAVE NO-GO DECISION" always enabled

---

### 5. Data Persistence Testing

**Verify database writes:**

1. Complete workflow and save GO decision
2. Check database tables:
   - `decisions` table: New row with type="GO", all trade data
   - `positions` table: New position if user executes (Phase 4)
   - `evaluations` table: Checklist timestamp for impulse timer
   - `sizing_plans` table: Position plan saved

3. Complete workflow and save NO-GO decision
4. Check database tables:
   - `decisions` table: New row with type="NO-GO", reason, category

5. Restart application
6. Verify Dashboard shows correct data (persisted from database)

---

### 6. Edge Case Testing

**Zero Candidates:**
- Run FINVIZ scan when no tickers match preset
- **Expected:** Empty results table with message "No candidates found"
- **Expected:** No errors, graceful handling

**Network Errors:**
- Simulate backend down (stop server)
- Attempt FINVIZ scan
- **Expected:** Error notification with clear message
- **Expected:** No silent failures, no console errors

**At Heat Cap:**
- Portfolio heat at 99.9% of cap
- Attempt to add position
- **Expected:** RED warning, exact overage shown
- **Expected:** "Calculate Max Shares" shows zero or minimal position

**Invalid Inputs:**
- Checklist: Enter negative ATR
- **Expected:** Validation error, clear message
- Position Sizing: Invalid entry price
- **Expected:** Validation error, clear message

**Rapid Navigation:**
- Quickly navigate between pages during calculations
- **Expected:** No race conditions, no stale data
- **Expected:** Loading states prevent double-clicks

**Multiple Tabs (if web-based):**
- Open two browser tabs
- Complete workflow in Tab 1
- **Expected:** Tab 2 data does not update (requires refresh)
- **Future:** WebSocket for live updates (Phase 4)

---

### 7. Logging Verification

**Check logs for comprehensive tracking:**

1. **FINVIZ Scan:**
   - Scan requested, duration, result count, sectors found

2. **Checklist:**
   - Checkbox changes (which items checked/unchecked)
   - Banner transitions (RED → YELLOW → GREEN with timestamps)
   - Quality score changes
   - Timer started

3. **Position Sizing:**
   - Method selected, inputs, calculation results
   - Save successful

4. **Heat Check:**
   - Current heat, proposed trade, new heat, caps, pass/fail
   - Overage amounts if failed

5. **Gate Check:**
   - All 5 gate results
   - Reasons for failures
   - Overall pass/fail

6. **Decision Saving:**
   - Type (GO/NO-GO), ticker, timestamp
   - Full trade data (for GO)
   - Reason and category (for NO-GO)

**Verify correlation IDs:** Each request should have unique correlation ID for tracing across logs.

---

### 8. Performance Testing

**Measure response times:**

- FINVIZ scan: < 5 seconds (network-bound)
- Position sizing calculation: < 100ms
- Heat check: < 500ms (database query)
- Gate check: < 500ms (multiple validations)
- Decision save: < 500ms (database write)

**UI responsiveness:**
- Banner transitions: < 100ms (CSS animation)
- Button clicks: Immediate feedback (loading state)
- Navigation: < 200ms (SvelteKit routing)

**Database queries:**
- Dashboard load: < 500ms (multiple tables)
- Positions list: < 200ms (single query)

**Identify bottlenecks:**
- Log slow requests (>500ms) in backend
- Frontend: Log slow renders (>100ms) if any
- Optimize if needed

---

## Bug Fixing Protocol

**When bugs are found:**

1. **Document the bug:**
   - What were you doing?
   - What did you expect?
   - What actually happened?
   - Steps to reproduce

2. **Check logs:**
   - Backend logs (correlation ID)
   - Frontend console (errors)
   - Database state (inspect tables)

3. **Fix and verify:**
   - Make the fix
   - Re-test the specific scenario
   - Run regression tests (other scenarios)

4. **Log the fix:**
   - Update LLM-update.md with bug description and fix
   - Commit changes with clear message

---

## Testing Checklist

**Complete this checklist before marking Phase 3 as done:**

```
[ ] Happy path test (all gates pass) - Complete workflow from scan to GO decision
[ ] Gate 1 failure test (banner not GREEN)
[ ] Gate 2 failure test (2-minute timer not elapsed)
[ ] Gate 3 failure test (ticker on cooldown)
[ ] Gate 4 failure test (heat caps exceeded)
[ ] Gate 5 failure test (sizing not completed)
[ ] NO-GO decision path test (modal, reason, save)
[ ] Banner transitions test (RED → YELLOW → GREEN with pulse)
[ ] Cool-off timer test (countdown, persistence, validation)
[ ] Heat gauges test (colors, percentages, dollar amounts)
[ ] Gate results display test (visual cards, colors, details)
[ ] Button states test (enabled/disabled correctly)
[ ] GO decision database write test
[ ] NO-GO decision database write test
[ ] Dashboard persistence test (restart app, data intact)
[ ] Zero candidates edge case test
[ ] Network error edge case test
[ ] At heat cap edge case test
[ ] Invalid input validation test
[ ] Rapid navigation test (no race conditions)
[ ] Logging verification (all events logged with correlation IDs)
[ ] Performance testing (response times within targets)
[ ] Bug fixes documented and verified
```

---

## Expected Outcome

**All tests pass.** Core workflow (Phase 1-3) is functional and ready for user testing.

**User Experience:**
- Smooth flow from FINVIZ scan through GO/NO-GO decision
- Clear visual feedback at every step
- Gates enforce discipline (cannot bypass)
- Errors handled gracefully with helpful messages
- Performance feels instant

**Technical Quality:**
- All database writes successful
- All API endpoints working
- All UI components rendering correctly
- No console errors
- Comprehensive logging for debugging

**Confidence Level:**
You should feel confident saying: "The core discipline enforcement system works. Let's move to Phase 4 (Calendar & Polish)."

---

## Common Issues and Solutions

**Issue:** Timer doesn't persist across navigation
**Solution:** Timer store must use localStorage or sessionStorage

**Issue:** Gate check passes when it shouldn't
**Solution:** Backend validation logic - check gates.go implementation

**Issue:** Banner doesn't transition smoothly
**Solution:** Check CSS transition timing, verify reactive statements in Svelte

**Issue:** Heat check shows wrong percentages
**Solution:** Verify equity value in settings, check heat.go calculations

**Issue:** GO decision button enabled when gates failed
**Solution:** Check derived store logic for `allGatesPass` state

---

## Time Estimate

~6-10 hours (1-2 days)

- Initial testing: 3-4 hours
- Bug fixing: 2-4 hours
- Regression testing: 1-2 hours

---

## References

- [overview-plan.md - User Workflow](../plans/overview-plan.md#user-workflow)
- [overview-plan.md - The 5 Hard Gates](../plans/overview-plan.md#the-5-hard-gates-enforced-by-backend)
- [roadmap.md - Step 19](../plans/roadmap.md#step-19-integration-testing-core-workflow)

---

## Next Phase

After completing integration testing:

**Proceed to: [Phase 4: Calendar & Polish](../plans/roadmap.md#phase-4-calendar--polish-4-steps)**

Phase 4 includes:
- Step 20: Calendar Grid Component (10-week sector view)
- Step 21: Calendar Integration
- Step 22: Visual Polish & Animations
- Step 23: Theme Refinement

---

**Status:** 📋 Ready for Execution
**Created:** 2025-10-29
**Phase 3 Milestone:** ✅ Core workflow complete after this step

