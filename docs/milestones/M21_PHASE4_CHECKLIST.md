# M21 Phase 4 - Execution Checklist

**Quick Reference:** Step-by-step checklist for Phase 4 testing
**Duration:** 45-120 minutes
**See:** M21_PHASE4_TEST_SCRIPTS.md for detailed instructions

---

## Pre-Flight Check (5 min)

**Environment Ready:**
- [ ] Windows 10/11 available
- [ ] Excel installed with macro support
- [ ] Command Prompt available
- [ ] Admin rights (if needed for registry)

**Files Present (windows/ directory):**
- [ ] TradingPlatform.xlsm exists
- [ ] trading.db exists (not empty)
- [ ] tf-engine.exe exists (26 MB)
- [ ] Test data files copied from test-data/

**Quick Verification:**
```cmd
cd windows
dir TradingPlatform.xlsm
dir trading.db
dir tf-engine.exe

# Test engine
tf-engine.exe --version
# Expected: tf-engine version 3.0.0-dev

# Test database
tf-engine.exe --db trading.db get-settings --format json
# Expected: JSON with 5 settings
```

**Prerequisites Complete:**
- [ ] Phase 1: Setup complete ✅
- [ ] Phase 2: Smoke tests pass (4/4) ✅
- [ ] Phase 3: VBA unit tests pass (14/14) ✅

**If any pre-flight checks fail → STOP and fix before proceeding**

---

## Workflow 1: Position Sizing (10-15 min)

**Setup:**
- [ ] Create "Position Sizing" worksheet
- [ ] Add headers and labels (A1-A29)
- [ ] Add input validation (B8 dropdown: stock/opt-delta-atr/opt-maxloss)
- [ ] Add Calculate button with VBA code
- [ ] Exit Design Mode

**Tests:**
- [ ] Test 1.1: Stock sizing (default settings)
  - Input: AAPL, Entry=180, ATR=1.5, K=2, Method=stock
  - Expected: R=$75, Stop=3, Shares=25
- [ ] Test 1.2: Stock sizing (with overrides)
  - Input: MSFT, Entry=400, ATR=3.0, Equity=20000, Risk%=1.0
  - Expected: R=$200, Stop=6, Shares=33
- [ ] Test 1.3: Option sizing (delta-ATR)
  - Input: NVDA, Entry=500, ATR=5.0, Method=opt-delta-atr, Delta=0.30
  - Expected: Contracts calculated via delta
- [ ] Test 1.4: Option sizing (max-loss)
  - Input: SPY, Method=opt-maxloss, MaxLoss=70
  - Expected: Contracts=1

**Validation:**
- [ ] All calculations match expected values
- [ ] Actual Risk ≤ Specified Risk in all cases
- [ ] Status shows success with correlation ID
- [ ] No errors in TradingSystem_Debug.log

**Issues Found:** _____________

**Time Taken:** _______ min

---

## Workflow 2: Checklist Evaluation (10-15 min)

**Setup:**
- [ ] Create "Checklist" worksheet
- [ ] Add headers and ticker input (A1, A3)
- [ ] Add 6 ActiveX checkboxes (Check1-Check6) with labels
- [ ] Add results section (A15-A28)
- [ ] Add Evaluate button with VBA code
- [ ] Exit Design Mode

**Tests:**
- [ ] Test 2.1: GREEN banner (all 6 checked)
  - Input: AAPL, all boxes checked
  - Expected: Banner=GREEN, Missing=0, AllowSave=TRUE
- [ ] Test 2.2: YELLOW banner (2 missing)
  - Input: MSFT, 4 checked, 2 unchecked
  - Expected: Banner=YELLOW, Missing=2, AllowSave=FALSE
- [ ] Test 2.3: YELLOW banner (1 missing)
  - Input: NVDA, 5 checked, 1 unchecked
  - Expected: Banner=YELLOW, Missing=1, AllowSave=FALSE
- [ ] Test 2.4: RED banner (3+ missing)
  - Input: SPY, 3 checked, 3 unchecked
  - Expected: Banner=RED, Missing=3, AllowSave=FALSE
- [ ] Test 2.5: Banner persistence
  - Evaluate multiple tickers, verify each persists

**Validation:**
- [ ] Color coding correct (green/yellow/red backgrounds)
- [ ] Allow Save = TRUE only for GREEN
- [ ] Missing items listed correctly
- [ ] Status shows success

**Issues Found:** _____________

**Time Taken:** _______ min

---

## Workflow 3: Heat Management (10-15 min)

**Setup:**
- [ ] Create "Heat Check" worksheet
- [ ] Add headers and inputs (A1-A7)
- [ ] Add dropdown for buckets (B5)
- [ ] Add results sections (portfolio A9-A15, bucket A17-A23)
- [ ] Add Check Heat button with VBA code
- [ ] Exit Design Mode

**Tests:**
- [ ] Test 3.1: No open positions
  - Input: Risk=$75, Bucket=Tech/Comm
  - Expected: Current=$0, New=$75, Allowed=TRUE
- [ ] Test 3.2: Portfolio cap exceeded
  - Input: Risk=$450
  - Expected: Portfolio Exceeded=TRUE, Allowed=FALSE
- [ ] Test 3.3: Bucket cap exceeded
  - Input: Risk=$200
  - Expected: Bucket Exceeded=TRUE, Allowed=FALSE
- [ ] Test 3.4: Exactly at cap
  - Input: Risk=$150, Bucket=Finance
  - Expected: Heat=100%, Allowed=TRUE
- [ ] Test 3.5: With open positions
  - Setup: Add AAPL position (R=$75)
  - Input: Risk=$80, Bucket=Tech/Comm
  - Expected: Current=$75, New=$155, Bucket Exceeded=TRUE
- [ ] Test 3.6: Different buckets
  - Input: Risk=$80, Bucket=Healthcare (AAPL in Tech/Comm)
  - Expected: Bucket Current=$0, Allowed=TRUE

**Validation:**
- [ ] Portfolio Cap = $400 (4% × $10,000)
- [ ] Bucket Cap = $150 (1.5% × $10,000)
- [ ] Heat % calculated correctly
- [ ] Allowed = FALSE when either cap exceeded
- [ ] Green/red color coding correct

**Issues Found:** _____________

**Time Taken:** _______ min

---

## Workflow 4: Save Decision (15-30 min)

**Setup:**
- [ ] Create "Trade Entry" worksheet
- [ ] Add complete form (A1-A29)
- [ ] Add dropdowns (Method B7, Banner B8, Bucket B12)
- [ ] Add gate status section (A17-A22)
- [ ] Add Save Decision button with VBA code
- [ ] Exit Design Mode

**Pre-Test Setup:**
- [ ] Import candidates:
  ```cmd
  tf-engine.exe --db trading.db import-candidates --tickers AAPL,MSFT,NVDA,SPY,JPM --preset TEST
  ```
- [ ] Clear positions (optional):
  ```cmd
  tf-engine.exe --db trading.db -c "DELETE FROM positions"
  ```
- [ ] Evaluate AAPL checklist (GREEN, all 6 boxes)
- [ ] **Wait 2 minutes** (impulse brake)

**Tests:**
- [ ] Test 4.1: Happy path (all gates pass)
  - Input: AAPL, Entry=180, ATR=1.5, Banner=GREEN, Risk=75, Shares=25, Bucket=Tech/Comm
  - Expected: All gates PASS, Decision ID=1, Saved=TRUE
- [ ] Test 4.2: Gate 1 rejection (YELLOW banner)
  - Input: Same as 4.1, Banner=YELLOW
  - Expected: Gate 1 FAIL, Rejected: "Banner must be GREEN"
- [ ] Test 4.3: Gate 1 rejection (RED banner)
  - Input: Same as 4.1, Banner=RED
  - Expected: Gate 1 FAIL, Rejected
- [ ] Test 4.4: Gate 2 rejection (not in candidates)
  - Input: Ticker=ZZZZ (not imported)
  - Expected: Gate 2 FAIL, Rejected: "Not in candidates"
- [ ] Test 4.5: Gate 3 rejection (impulse brake)
  - Setup: Evaluate MSFT checklist (GREEN)
  - Input: Immediately try to save MSFT (don't wait)
  - Expected: Gate 3 FAIL, Rejected: "Wait XX seconds"
- [ ] Test 4.6: Gate 4 rejection (bucket cooldown)
  - Setup: AAPL saved (Test 4.1) in Tech/Comm
  - Input: NVDA (also Tech/Comm bucket)
  - Expected: Gate 4 FAIL, Rejected: "Bucket in cooldown"
- [ ] Test 4.7: Gate 5 rejection (portfolio cap)
  - Input: Risk=$450 (exceeds $400 cap)
  - Expected: Gate 5 FAIL, Rejected: "Portfolio heat exceeded"
- [ ] Test 4.8: Gate 5 rejection (bucket cap)
  - Input: Risk=$200 (exceeds $150 bucket cap)
  - Expected: Gate 5 FAIL, Rejected: "Bucket heat exceeded"
- [ ] Test 4.9: Multiple gate failures
  - Input: Ticker=ZZZZ, Banner=YELLOW, Risk=$450
  - Expected: Gates 1,2,5 FAIL, Rejected with first reason
- [ ] Test 4.10: Form behavior
  - Verify form clears on success
  - Verify form persists on rejection

**Validation:**
- [ ] All 5 gates enforced
- [ ] Rejection messages clear and specific
- [ ] Gate status displayed correctly
- [ ] Database records created on success
- [ ] Form behavior correct (clear vs. persist)

**Database Verification:**
```cmd
# Check saved decisions
tf-engine.exe --db trading.db -c "SELECT COUNT(*) FROM decisions"

# Check open positions
tf-engine.exe --db trading.db -c "SELECT * FROM positions WHERE status='open'"
```

**Issues Found:** _____________

**Time Taken:** _______ min

---

## Post-Test Review (5-10 min)

**Log Review:**
- [ ] Open TradingSystem_Debug.log
- [ ] Check for ERROR entries (none expected except rejections)
- [ ] Verify all correlation IDs present
- [ ] Note execution times (<500ms per command)

**Performance Metrics:**
```
Average execution time: _______ ms
Slowest operation: _______
Total Phase 4 duration: _______ minutes
```

**Database Integrity:**
```cmd
# Verify schema
sqlite3 trading.db ".schema"

# Check record counts
tf-engine.exe --db trading.db -c "SELECT COUNT(*) FROM candidates"
tf-engine.exe --db trading.db -c "SELECT COUNT(*) FROM decisions"
tf-engine.exe --db trading.db -c "SELECT COUNT(*) FROM positions"
```

**System Health:**
- [ ] No VBA compile errors
- [ ] No shell execution failures
- [ ] JSON parsing 100% successful
- [ ] Named ranges resolved correctly
- [ ] Database operations successful

---

## Issues Log

Document any issues encountered:

| Test | Issue | Severity | Status | Notes |
|------|-------|----------|--------|-------|
| ____ | _____ | Low/Med/High | Open/Fixed | _____ |
| ____ | _____ | Low/Med/High | Open/Fixed | _____ |
| ____ | _____ | Low/Med/High | Open/Fixed | _____ |

---

## Final Sign-Off

**Phase 4 Complete:**
- [ ] All 4 workflows tested ✅
- [ ] 25 integration tests executed ✅
- [ ] All tests pass or issues documented ✅
- [ ] Logs reviewed ✅
- [ ] Database verified ✅
- [ ] Ready for M21 completion summary ✅

**Total Time:** _______ minutes

**Overall Assessment:**
- [ ] System ready for production use
- [ ] Minor issues found (list above)
- [ ] Major issues found - requires fixes before M21 complete

**Tested By:** _____________
**Date:** _____________
**Session ID:** _____________

---

## Next Steps

After Phase 4 completion:

1. **Update M21_PROGRESS.md:**
   - Change status to "Phase 4 Complete"
   - Add Phase 4 results
   - Document issues found

2. **Create M21 Completion Summary:**
   - All 4 phases complete
   - Total issues fixed: 9 (Phases 1-3) + X (Phase 4)
   - System validated and ready

3. **Begin M22 (Next Milestone):**
   - Review PLAN.md for next milestone
   - Start planning M22 work

**Congratulations on completing M21 Phase 4! 🎉**

---

**Checklist Created:** 2025-10-27
**For Use With:** M21_PHASE4_TEST_SCRIPTS.md
**Estimated Duration:** 45-120 minutes
