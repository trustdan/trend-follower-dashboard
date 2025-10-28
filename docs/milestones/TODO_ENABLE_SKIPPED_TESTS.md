# TODO: Enable Skipped Integration Tests

**Created:** 2025-10-28 (M21 Phase 4)
**Location:** `excel/vba/TFIntegrationTests.bas`

## Tests Currently Skipped (9 tests)

These tests are commented out and need to be re-enabled when their dependencies are implemented:

### 1. Heat Management Tests (4 tests) - ✅ **COMPLETED (M23)**

**Lines:** ~752-940 in TFIntegrationTests.bas

**Tests:**
- ✅ Test 3.1: Heat Check (No Open Positions)
- ✅ Test 3.2: Portfolio Cap Exceeded
- ✅ Test 3.3: Bucket Cap Exceeded
- ✅ Test 3.4: Exactly At Cap (Edge Case)

**Status:** Re-enabled on 2025-10-28

**Changes Made:**
1. ✅ Renamed `check-heat` command to `heat` in internal/cli/heat.go
2. ✅ Changed flags from `--add-risk`/`--add-bucket` to `--risk`/`--bucket`
3. ✅ Unskipped all 4 tests in TFIntegrationTests.bas (lines 758, 811, 864, 917)
4. ✅ Fixed percentage formatting in VBA tests (divide by 100 before Format())
5. ✅ Verified all 4 test scenarios produce expected output

**Command syntax:**
```bash
tf-engine heat --risk <amount> --bucket <bucket-name>
```

**JSON output (verified):**
```json
{
  "current_portfolio_heat": 0,
  "new_portfolio_heat": 75,
  "portfolio_heat_pct": 18.75,
  "portfolio_cap": 400,
  "portfolio_cap_exceeded": false,
  "portfolio_overage": 0,
  "current_bucket_heat": 0,
  "new_bucket_heat": 75,
  "bucket_heat_pct": 50,
  "bucket_cap": 150,
  "bucket_cap_exceeded": false,
  "bucket_overage": 0,
  "allowed": true
}
```

---

### 2. Save-Decision Gate Tests (5 tests) - **BLOCKED: Impulse brake (Gate 3) runs first**

**Lines:** ~987-1275 in TFIntegrationTests.bas

**Tests:**
- Test 4.1: Happy Path (All Gates Pass) - requires 2-minute wait
- Test 4.2: Gate 1 Rejection (YELLOW Banner) - blocked by Gate 3
- Test 4.3: Gate 1 Rejection (RED Banner) - blocked by Gate 3
- Test 4.4: Gate 2 Rejection (Not in Candidates) - blocked by Gate 3
- Test 4.5: Gate 5 Rejection (Portfolio Cap) - blocked by Gate 3

**Issue:**
The gates run in order:
1. Banner GREEN check
2. Ticker in candidates
3. **Impulse brake (2-minute timer)** ← BLOCKS all subsequent gates
4. Bucket cooldown
5. Heat caps

Tests 4.2-4.5 try to test gates 1, 2, and 5, but Gate 3 fails first because the checklist evaluation happened <2 minutes ago.

**Possible solutions:**

#### Option A: Add test mode to skip Gate 3
Add a `--skip-gate` flag to save-decision command:
```bash
tf-engine save-decision --ticker AAPL --action GO --skip-gate 3
```

**Pros:** Enables automated testing of all gates
**Cons:** Creates a bypass mechanism (against discipline philosophy)

#### Option B: Add 2-minute delay in pre-test setup
Modify `PreTestSetup()` in TFIntegrationTests.bas to:
```vba
' Wait 2 minutes for impulse brake to clear
LogMessage "Waiting 2 minutes for impulse brake to clear...", True
Application.Wait Now + TimeValue("0:02:00")
```

**Pros:** Tests real gate behavior
**Cons:** Adds 2 minutes to test runtime

#### Option C: Keep manual testing for these
**Current approach** - these tests verify discipline enforcement, which is the core value prop. Manual testing ensures they work as intended.

**Pros:** Tests remain rigorous, no shortcuts
**Cons:** Not automated

**Recommendation:** Keep Option C (manual testing only) for these 5 tests. The discipline gates are the core of the system and should not have test bypasses.

---

## Re-enabling Checklist

Heat command implementation (M23):

- [x] Implement `internal/cli/heat.go` with proper flags
- [x] Add command registration in `cmd/tf-engine/main.go`
- [x] Test heat command manually: `tf-engine heat --risk 75 --bucket "Tech/Comm"`
- [x] Search TFIntegrationTests.bas for "SKIP: heat command not yet implemented"
- [x] Uncomment 4 test command lines (Tests 3.1-3.4)
- [x] Remove 4 SKIP blocks
- [x] Verify all 4 test scenarios produce expected JSON output
- [x] Fix VBA percentage formatting (divide by 100)
- [x] Update this file to mark heat tests as re-enabled

---

## Current Test Status

**Total tests:** 19
- **PASS:** 13 (68.4%) - All automatable tests passing ✅
- **FAIL:** 0 (0.0%)
- **ERROR:** 0 (0.0%)
- **SKIP:** 6 (31.6%)

**Automatable tests:** 13/13 passing (100%)
- Position sizing: 4/4 ✅
- Checklist: 5/5 ✅
- Heat management: 4/4 ✅ (re-enabled in M23)

**Manual tests:** 6 require timing-based testing
- Save-decision gate tests: 5 (require 2-minute delays)
- UI validation: 1 (requires Windows testing)

---

**Last updated:** 2025-10-28 (M23 complete)
**Milestone:** M23 - Heat Command Implementation
