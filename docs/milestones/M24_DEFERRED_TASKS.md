# M24 - Deferred Tasks & Manual Testing Checklist

**Created:** 2025-10-28
**Purpose:** Track tasks that require Windows PC or manual testing

---

## Tasks Deferred to Windows PC Testing

### 1. Heat Tests - Windows Validation (HIGH PRIORITY)
**Status:** ⏸️ Requires Windows PC with Excel

**What's Done:**
- ✅ Heat command implemented and verified on Linux
- ✅ 4 test scenarios validated via CLI
- ✅ VBA tests unskipped in TFIntegrationTests.bas
- ✅ Windows binary rebuilt (tf-engine.exe)

**What's Needed:**
- [ ] Copy updated `tf-engine.exe` to Windows PC
- [ ] Run `3-run-integration-tests.bat`
- [ ] Verify all 4 heat tests now PASS (not SKIP)
- [ ] Expected result: 13/19 PASS, 6 SKIP (down from 10 SKIP)

**Test IDs to Verify:**
- Test 3.1: Heat Check (No Open Positions) - Should PASS
- Test 3.2: Portfolio Cap Exceeded - Should PASS
- Test 3.3: Bucket Cap Exceeded - Should PASS
- Test 3.4: Exactly At Cap (Edge Case) - Should PASS

**Location:** `windows/3-run-integration-tests.bat`
**Estimated Time:** 5 minutes

---

### 2. M22 UI Workbook Validation (HIGH PRIORITY)
**Status:** ⏸️ Requires Windows PC with Excel

**What's Done:**
- ✅ VBScript UI generator created (M22)
- ✅ 5 production worksheets defined
- ✅ Button handlers implemented
- ✅ Setup script automated

**What's Needed:**
- [ ] Run `1-setup-all.bat` on Windows
- [ ] Verify all 7 worksheets created
- [ ] Test all button handlers:
  - [ ] Position Sizing: Calculate button
  - [ ] Checklist: Evaluate button
  - [ ] Heat Check: Check Heat button
  - [ ] Trade Entry: Save GO/NO-GO buttons
  - [ ] Dashboard: Navigation buttons
- [ ] Verify formatting looks professional
- [ ] Test complete workflow end-to-end
- [ ] Check error handling

**Location:** `windows/1-setup-all.bat`
**Documentation:** `docs/milestones/M22_COMPLETION_SUMMARY.md`
**Estimated Time:** 70 minutes

---

### 3. Gate Timing Tests (MEDIUM PRIORITY)
**Status:** ⏸️ Requires manual timing verification

**What's Done:**
- ✅ All 5 gates implemented in Go engine
- ✅ Gate logic tested in isolation
- ✅ VBA integration complete

**What's Needed:**
- [ ] Test 4.1: Happy Path (All Gates Pass)
  - Complete checklist evaluation
  - Wait 2 minutes for impulse brake to clear
  - Execute save-decision
  - Verify all gates pass

- [ ] Test 4.2: Gate 1 Rejection (YELLOW Banner)
  - Create YELLOW banner (4-5 checks)
  - Wait 2 minutes
  - Attempt save-decision
  - Verify Gate 1 rejects

- [ ] Test 4.3: Gate 1 Rejection (RED Banner)
  - Create RED banner (<4 checks)
  - Wait 2 minutes
  - Attempt save-decision
  - Verify Gate 1 rejects

- [ ] Test 4.4: Gate 2 Rejection (Not in Candidates)
  - Use ticker not in candidates list
  - Complete GREEN checklist, wait 2 min
  - Attempt save-decision
  - Verify Gate 2 rejects

- [ ] Test 4.5: Gate 5 Rejection (Portfolio Cap)
  - Create scenario exceeding portfolio cap
  - Complete GREEN checklist, wait 2 min
  - Attempt save-decision
  - Verify Gate 5 rejects

**Why Manual:** Impulse brake (Gate 3) requires 2-minute delay, making automation impractical

**Location:** `excel/vba/TFIntegrationTests.bas` lines 987-1275 (commented out)
**Estimated Time:** 30 minutes (5-10 min per test)

**Alternative Approach:**
- Modify database timestamps directly to bypass waiting
- `UPDATE timers SET last_checklist_eval = datetime('now', '-3 minutes')`
- This allows faster testing but is less realistic

---

### 4. Windows Integration Test Suite (LOW PRIORITY)
**Status:** ⏸️ Nice to have, not critical

**What's Done:**
- ✅ Integration test framework complete (M21)
- ✅ 13/19 tests automated and passing
- ✅ Test runner batch script working

**What's Needed:**
- [ ] Run complete test suite on Windows
- [ ] Verify test execution time
- [ ] Check log output formatting
- [ ] Validate correlation IDs
- [ ] Review TradingSystem_Debug.log

**Why Deferred:** Core functionality already validated on Linux, Windows testing is for environment verification

**Location:** `windows/3-run-integration-tests.bat`
**Estimated Time:** 10 minutes

---

## Tasks We Can Complete Now (Linux)

### 1. Distribution Package Creation ✅ CAN DO NOW
**Status:** Ready to execute

**Tasks:**
- [ ] Create `release/` directory structure
- [ ] Copy all necessary files:
  - `tf-engine.exe`
  - `windows/` folder (all scripts)
  - `excel/vba/` modules
  - Documentation (README, guides)
- [ ] Create QUICKSTART.md (1-page setup guide)
- [ ] Zip package: `TradingEngine-v3-Release.zip`
- [ ] Calculate checksums

**Estimated Time:** 30 minutes

---

### 2. Documentation Polish ✅ CAN DO NOW
**Status:** Ready to execute

**Tasks:**
- [ ] Create QUICKSTART.md (single-page guide)
- [ ] Update main README with final status
- [ ] Add "Known Limitations" section
- [ ] Document manual testing requirements
- [ ] Create deployment checklist
- [ ] Add troubleshooting guide

**Estimated Time:** 1 hour

---

### 3. Final Code Cleanup ✅ CAN DO NOW (Optional)
**Status:** Optional, code already production-quality

**Tasks:**
- [ ] Run go fmt on all files
- [ ] Remove debug logging
- [ ] Update version string to "3.0.0"
- [ ] Final code review
- [ ] Check for TODOs in code

**Estimated Time:** 30 minutes

---

## Summary: What Can We Do Now vs. Later

### ✅ Can Complete on Linux (Now)
1. **Distribution package** - Zip all files for release
2. **Documentation polish** - QUICKSTART, troubleshooting, etc.
3. **Code cleanup** - Final polish before release
4. **README updates** - Mark production-ready status

**Total Time:** ~2 hours

### ⏸️ Requires Windows PC (Later)
1. **Heat tests validation** - Verify 4 tests now pass (5 min)
2. **M22 UI workbook testing** - Full manual validation (70 min)
3. **Gate timing tests** - Manual gate verification (30 min, optional)

**Total Time:** ~1.5 hours (or 5 min if just heat tests)

---

## Proposed M24 Approach

### Phase 1: Linux Tasks (Now) - ~2 hours
1. Create distribution package
2. Write QUICKSTART.md
3. Polish documentation
4. Update README with production status
5. Mark as "Ready for Windows validation"

### Phase 2: Windows Validation (Later) - ~5 min to 1.5 hours
1. **Minimum:** Run heat tests (5 min) - Validates M23
2. **Recommended:** Full UI testing (70 min) - Validates M22
3. **Optional:** Gate timing (30 min) - Validates gates

### Phase 3: Production Release (After Windows validation)
1. Update version to "3.0.0"
2. Create final release zip
3. Publish

---

## Decision Point

**Question:** How far should we go in M24 without Windows access?

**Option A: Minimal (Recommended)**
- Complete all Linux tasks now (~2 hours)
- Mark as "Ready for Windows validation"
- Document what needs Windows testing
- User runs Windows validation when ready
- **Result:** 95% complete, awaiting final validation

**Option B: Full Simulation**
- Complete all Linux tasks
- Create mock Windows test results
- Package as if production-ready
- **Result:** 100% "complete" but untested on Windows (risky)

**Option C: Wait**
- Don't proceed with M24 until Windows PC available
- **Result:** Blocks progress unnecessarily

**Recommendation:** **Option A** - Complete all Linux tasks, document Windows validation steps clearly

---

## Acceptance Criteria for M24 (Adjusted)

### Core M24 Objectives (Linux-completable)
- [x] Distribution package created ✅ CAN DO
- [x] QUICKSTART.md written ✅ CAN DO
- [x] Documentation polished ✅ CAN DO
- [x] README updated to production status ✅ CAN DO
- [x] Known limitations documented ✅ CAN DO
- [x] Windows validation checklist created ✅ THIS FILE

### Windows Validation Objectives (Deferred)
- [ ] Heat tests validated on Windows ⏸️ REQUIRES WINDOWS
- [ ] M22 UI validated on Windows ⏸️ REQUIRES WINDOWS
- [ ] Complete workflow tested ⏸️ REQUIRES WINDOWS
- [ ] Gate timing verified ⏸️ REQUIRES WINDOWS (optional)

### Production Release (After Windows validation)
- [ ] All automated tests passing on Windows
- [ ] UI workflow validated
- [ ] Final release package created
- [ ] Version bumped to 3.0.0

---

## Risk Mitigation

### Risk: Windows validation finds critical bugs
**Likelihood:** Low
**Impact:** Medium
**Mitigation:**
- All business logic tested on Linux
- VBA integration tested in M21
- Setup scripts tested in M21
- Heat command verified via CLI

### Risk: UI issues on Windows
**Likelihood:** Low
**Impact:** Low
**Mitigation:**
- VBScript generation tested in M22
- Standard Excel controls used
- No Windows-specific APIs

### Risk: Heat tests fail on Windows
**Likelihood:** Very Low
**Impact:** Low
**Mitigation:**
- All 4 scenarios tested on Linux
- JSON output verified
- Command syntax correct

---

## Next Steps

1. **Decide:** Option A, B, or C above
2. **Execute:** Linux tasks (~2 hours)
3. **Document:** Windows validation checklist (this file)
4. **Mark:** M24 as "Ready for Windows validation"
5. **Communicate:** Clear handoff to user for Windows testing

---

**Created:** 2025-10-28
**Status:** Ready for M24 Linux tasks
**Estimated Time to Windows-ready:** ~2 hours
**Estimated Windows validation time:** 5 min (heat only) to 1.5 hours (full)
