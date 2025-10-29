# Phase 5 - Step 25: Bug Fixing Sprint

**TF = Trend Following** - Systematic trading discipline enforcement system

**Phase:** 5 - Testing & Packaging
**Step:** 25 of 28
**Duration:** 2-3 days
**Dependencies:** Step 24 complete (Testing complete, bugs identified)

---

## Objectives

Address all bugs and issues discovered during comprehensive testing. Prioritize by severity: Critical (app crashes, data loss, gates not enforced) → High (incorrect calculations, UI broken) → Medium (cosmetic issues, minor UX problems) → Low (nice-to-haves). Fix critical and high severity bugs first. Document any bugs that cannot be fixed immediately with workarounds.

**Purpose:** Ensure zero critical bugs and minimal high-severity bugs before proceeding to packaging and release.

---

## Success Criteria

- [ ] All critical bugs fixed
- [ ] All high-severity bugs fixed or documented with workarounds
- [ ] Medium bugs addressed where feasible
- [ ] Regression tests added for each bug fix
- [ ] Verification testing passed (re-test fixed bugs)
- [ ] `docs/KNOWN_LIMITATIONS.md` created for unfixable issues
- [ ] Bug tracker updated with resolution status
- [ ] No new bugs introduced by fixes (regression testing)

---

## Prerequisites

**Completed:**
- Step 24 (Comprehensive testing)
- Bug tracker populated with all found issues
- Test results documented

**Required:**
- Development environment fully set up
- Access to both Linux and Windows for testing

---

## Implementation Plan

### Part 1: Bug Prioritization & Planning (2 hours)

#### Task 1.1: Review and Categorize Bugs (1 hour)

Review `docs/BUG_TRACKER.md` and categorize all bugs:

**Severity Definitions:**

**CRITICAL:**
- App crashes or hangs
- Data loss or corruption
- Gates can be bypassed (security issue)
- Cannot save decisions
- Cannot complete core workflow

**HIGH:**
- Incorrect calculations (position sizing, heat, etc.)
- UI completely broken (cannot navigate)
- Feature non-functional with no workaround
- Performance severely degraded

**MEDIUM:**
- UI partially broken (layout issues)
- Feature non-functional but has workaround
- Minor calculation errors
- Confusing error messages

**LOW:**
- Cosmetic issues (alignment, colors)
- Minor UX annoyances
- Typos in UI text
- Missing tooltips

#### Task 1.2: Create Fix Plan (1 hour)

**File:** `docs/BUG_FIX_PLAN.md`

```markdown
# Bug Fix Plan

## Sprint Goal

Fix all CRITICAL and HIGH severity bugs. Address MEDIUM bugs if time permits. Log LOW bugs for future releases.

## Priority Order

### Phase 1: Critical Bugs (Day 1)
- [ ] BUG-001: [Description] - Estimated: 2 hours
- [ ] BUG-002: [Description] - Estimated: 1 hour
- [ ] BUG-003: [Description] - Estimated: 3 hours

**Total Estimated:** 6 hours

### Phase 2: High Severity Bugs (Day 2)
- [ ] BUG-101: [Description] - Estimated: 1 hour
- [ ] BUG-102: [Description] - Estimated: 2 hours
- [ ] BUG-103: [Description] - Estimated: 1 hour

**Total Estimated:** 4 hours

### Phase 3: Medium Severity Bugs (Day 3, if time permits)
- [ ] BUG-201: [Description] - Estimated: 30 min
- [ ] BUG-202: [Description] - Estimated: 1 hour

**Total Estimated:** 1.5 hours

### Phase 4: Low Severity (Deferred to v1.1)
- [ ] BUG-301: [Description] - Logged for future
- [ ] BUG-302: [Description] - Logged for future

## Notes
- Each fix requires regression test
- Re-test on Windows after all fixes
- Update documentation if behavior changes
```

---

### Part 2: Fix Critical Bugs (1 day)

#### Example Bug Fix Workflow

**BUG-001: App crashes when database file is locked**

**Reproduction:**
1. Open app (creates trading.db)
2. Open trading.db in SQLite browser (locks file)
3. Try to save decision in app
4. App crashes: "database is locked"

**Root Cause:**
SQLite default behavior when file is locked.

**Fix:**

**File:** `backend/internal/storage/db.go`

```go
func NewDB(path string) (*DB, error) {
    // Add connection options for better concurrency
    dsn := fmt.Sprintf("%s?_journal_mode=WAL&_timeout=5000&_busy_timeout=5000", path)

    db, err := sql.Open("sqlite3", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Set connection pool limits
    db.SetMaxOpenConns(1) // SQLite works best with single writer
    db.SetMaxIdleConns(1)

    // Test connection with retry
    for i := 0; i < 3; i++ {
        if err := db.Ping(); err == nil {
            break
        }
        if i == 2 {
            return nil, fmt.Errorf("failed to connect to database after 3 attempts: %w", err)
        }
        time.Sleep(time.Second)
    }

    return &DB{db: db}, nil
}
```

**Regression Test:**

**File:** `backend/internal/storage/db_test.go`

```go
func TestDB_ConcurrentAccess(t *testing.T) {
    db, err := NewDB("test_concurrent.db")
    require.NoError(t, err)
    defer db.Close()
    defer os.Remove("test_concurrent.db")

    // Simulate concurrent writes
    var wg sync.WaitGroup
    errors := make(chan error, 10)

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()

            // Try to write to database
            _, err := db.db.Exec("INSERT INTO settings (key, value) VALUES (?, ?)",
                fmt.Sprintf("test_%d", n), "value")
            if err != nil {
                errors <- err
            }
        }(i)
    }

    wg.Wait()
    close(errors)

    // Should not have any errors
    errCount := 0
    for err := range errors {
        t.Errorf("Concurrent write failed: %v", err)
        errCount++
    }

    assert.Equal(t, 0, errCount, "No concurrent access errors expected")
}
```

**Verification:**
1. Run regression test: `go test ./internal/storage/ -run TestDB_ConcurrentAccess`
2. Manually test: Open app + SQLite browser, try to save decision
3. Expected: No crash, operation retries with timeout

**Update Bug Tracker:**

```markdown
### BUG-001: App crashes when database file is locked
**Severity:** Critical
**Found:** 2025-10-29
**Status:** FIXED (2025-10-30)
**Fix:** Added WAL mode, busy timeout, and connection retry logic
**Commit:** abc123
**Tested:** ✓ Regression test added, manual test passed
```

---

#### Common Bug Patterns & Fixes

**Pattern 1: Race Conditions**

**Symptom:** Intermittent failures, data inconsistency

**Fix:** Add proper locking:

```go
type SafeCache struct {
    mu    sync.RWMutex
    data  map[string]interface{}
}

func (c *SafeCache) Get(key string) interface{} {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.data[key]
}

func (c *SafeCache) Set(key string, val interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = val
}
```

**Pattern 2: Null Pointer Dereference**

**Symptom:** Panic: "nil pointer dereference"

**Fix:** Add nil checks:

```go
func ProcessPosition(pos *Position) error {
    if pos == nil {
        return fmt.Errorf("position cannot be nil")
    }

    if pos.Ticker == "" {
        return fmt.Errorf("ticker is required")
    }

    // ... rest of logic
}
```

**Pattern 3: Off-by-One Errors**

**Symptom:** Array index out of bounds

**Fix:** Validate bounds:

```go
func GetWeek(weeks []Week, index int) (*Week, error) {
    if index < 0 || index >= len(weeks) {
        return nil, fmt.Errorf("index %d out of range [0, %d)", index, len(weeks))
    }
    return &weeks[index], nil
}
```

**Pattern 4: Frontend State Not Updating**

**Symptom:** UI doesn't reflect changes

**Fix:** Ensure reactivity:

```svelte
<script lang="ts">
    // Bad: mutation doesn't trigger update
    function addItem() {
        items.push(newItem); // Won't trigger reactivity
    }

    // Good: reassignment triggers update
    function addItem() {
        items = [...items, newItem]; // Triggers reactivity
    }
</script>
```

**Pattern 5: Memory Leaks**

**Symptom:** Memory grows over time

**Fix:** Clean up event listeners:

```svelte
<script lang="ts">
    import { onMount, onDestroy } from 'svelte';

    let handler;

    onMount(() => {
        handler = () => { /* ... */ };
        window.addEventListener('resize', handler);
    });

    onDestroy(() => {
        if (handler) {
            window.removeEventListener('resize', handler);
        }
    });
</script>
```

---

### Part 3: Fix High Severity Bugs (1 day)

#### Example: Incorrect Position Sizing Calculation

**BUG-101: Position sizing returns wrong shares for options**

**Reproduction:**
1. Select method: "opt-delta-atr"
2. Enter: entry=5.00, delta=0.60, atr=0.50, k=2.0
3. Expected shares calculation incorrect

**Root Cause:**
Delta not applied correctly in calculation.

**Fix:**

**File:** `backend/internal/domain/sizing_options.go`

```go
func CalculateSizeOptDeltaATR(ticker string, entry, atr, delta, k, riskPct, equity float64) (*SizingResult, error) {
    // Validate inputs
    if delta <= 0 || delta > 1 {
        return nil, fmt.Errorf("delta must be between 0 and 1, got %.2f", delta)
    }

    // Risk per unit
    riskDollars := equity * riskPct

    // Stop distance adjusted for delta
    stopDistance := k * atr / delta // FIX: Was missing "/ delta"

    // Contracts per unit
    contracts := int(math.Floor(riskDollars / (stopDistance * 100)))

    // Actual risk
    actualRisk := float64(contracts) * stopDistance * 100

    return &SizingResult{
        Ticker:       ticker,
        Method:       "opt-delta-atr",
        Entry:        entry,
        Contracts:    contracts,
        RiskDollars:  riskDollars,
        ActualRisk:   actualRisk,
        StopDistance: stopDistance,
        InitialStop:  entry - stopDistance,
    }, nil
}
```

**Regression Test:**

**File:** `backend/internal/domain/sizing_options_test.go`

```go
func TestCalculateSizeOptDeltaATR(t *testing.T) {
    result, err := CalculateSizeOptDeltaATR("SPY", 450.0, 5.0, 0.60, 2.0, 0.01, 100000)
    require.NoError(t, err)

    // With delta=0.60, atr=5.0, k=2.0
    // stopDistance = 2.0 * 5.0 / 0.60 = 16.67
    assert.InDelta(t, 16.67, result.StopDistance, 0.01)

    // Risk = $1000 (1% of $100k)
    // Contracts = floor(1000 / (16.67 * 100)) = floor(0.60) = 0
    // (This shows the position is too risky - correct behavior)
    assert.Equal(t, 0, result.Contracts)
}
```

**Verification:**
1. Run test: `go test ./internal/domain/ -run TestCalculateSizeOptDeltaATR`
2. Manual test: Enter same values in UI, verify correct calculation
3. Compare with Van Tharp formula (external verification)

---

### Part 4: Medium & Low Bugs (As Time Permits)

**Medium Example: Calendar cell tooltips not showing**

**Fix:**

**File:** `ui/src/lib/components/CalendarCell.svelte`

```svelte
<script lang="ts">
    import Tooltip from './Tooltip.svelte';
    export let cell: any;
</script>

<!-- Before: title attribute (doesn't work in some browsers) -->
<div title="Entry: {cell.entry_date}, Risk: ${cell.risk}">
    {cell.ticker}
</div>

<!-- After: use Tooltip component -->
<Tooltip text="Entry: {cell.entry_date}, Risk: ${cell.risk.toFixed(2)}">
    <div class="ticker-tag">
        {cell.ticker}
    </div>
</Tooltip>
```

**Low Example: Typo in error message**

**File:** `backend/internal/domain/gates.go`

```go
// Before
return fmt.Errorf("ticker is on cooldown untill %s", cooldown.Until)

// After
return fmt.Errorf("ticker is on cooldown until %s", cooldown.Until)
```

---

### Part 5: Verification Testing (4 hours)

#### Task 5.1: Re-run Automated Tests

```bash
# Backend
cd backend/
go test ./... -v

# Frontend
cd ui/
npm run test
```

**Expected:** All tests pass, including new regression tests.

#### Task 5.2: Re-test Fixed Bugs Manually

For each fixed bug:
1. Follow original reproduction steps
2. Verify bug no longer occurs
3. Test related functionality (regression check)
4. Mark as "VERIFIED" in bug tracker

#### Task 5.3: Full Workflow Smoke Test

Run complete happy path one more time:
1. FINVIZ scan → Import candidates
2. Complete checklist (GREEN)
3. Calculate position size
4. Check heat
5. Wait 2 minutes
6. Run gates → All pass
7. Save GO decision
8. Verify appears in Dashboard and Calendar

**Expected:** No errors, smooth workflow.

#### Task 5.4: Windows Verification

```bash
# Rebuild for Windows
cd backend/
GOOS=windows GOARCH=amd64 go build -o tf-engine.exe cmd/tf-engine/main.go

# Test on Windows machine
# Run through critical bug scenarios to verify fixes work on Windows
```

---

### Part 6: Document Known Limitations (1 hour)

**File:** `docs/KNOWN_LIMITATIONS.md`

```markdown
# Known Limitations

**Last Updated:** 2025-10-30
**Version:** 1.0.0

## Limitations by Design

These are intentional design decisions, not bugs:

### 1. No Direct Trade Execution

**Limitation:** The system does not execute trades in your broker.

**Reason:** Intentional separation. The trader must manually enter trades in their broker as a final verification step.

**Workaround:** After saving a GO decision, manually execute in your broker.

---

### 2. Single User Only

**Limitation:** The system supports only one user (no multi-user accounts).

**Reason:** Designed for individual traders, not teams.

**Workaround:** Each trader runs their own instance with their own database.

---

### 3. No Historical Backtesting

**Limitation:** Cannot run backtests on historical data.

**Reason:** This is a live trading discipline tool, not a backtesting engine.

**Workaround:** Use TradingView or dedicated backtesting software.

---

## Technical Limitations

### 1. FINVIZ Scraper Dependent on Page Structure

**Limitation:** If FINVIZ changes their HTML structure, the scraper may break.

**Impact:** Cannot scan for candidates automatically.

**Workaround:** Manually enter tickers, or wait for update.

**Status:** Monitoring FINVIZ for changes.

---

### 2. Browser Required for UI

**Limitation:** Must have a web browser installed.

**Impact:** Cannot run on headless servers.

**Workaround:** Use CLI mode (if implemented), or run on machine with browser.

---

### 3. SQLite Single Writer

**Limitation:** Only one process can write to database at a time.

**Impact:** Cannot run multiple instances against same database.

**Workaround:** Each instance uses its own database file.

---

## Unfixable Bugs (if any)

### None Currently

All critical and high severity bugs have been fixed as of v1.0.0.

---

## Planned Improvements (Future Versions)

### v1.1 (Next Release)
- Improved calendar cell tooltips
- Export decisions to CSV
- Custom color themes

### v1.2
- TradingView widget embedding (if feasible)
- Trade log search and filtering
- Performance charts (P&L over time)

### v2.0 (Long Term)
- Multi-monitor support
- Mobile companion app (read-only)
- Cloud sync for database (optional)
```

---

## Testing Checklist

### Critical Bugs Fixed
- [ ] All critical bugs from Step 24 addressed
- [ ] Regression tests added for each fix
- [ ] Manual verification passed
- [ ] No new bugs introduced

### High Severity Bugs Fixed
- [ ] All high bugs from Step 24 addressed
- [ ] Calculations verified correct
- [ ] UI functionality restored
- [ ] Workarounds documented (if any)

### Medium/Low Bugs
- [ ] Medium bugs fixed where time permits
- [ ] Low bugs logged for future releases
- [ ] Prioritization documented

### Verification
- [ ] All automated tests pass
- [ ] Fixed bugs re-tested manually
- [ ] Full workflow smoke test passed
- [ ] Windows verification complete

### Documentation
- [ ] Bug tracker updated with resolutions
- [ ] KNOWN_LIMITATIONS.md created
- [ ] API docs updated (if behavior changed)
- [ ] User guide updated (if needed)

---

## Troubleshooting

**Problem:** Fix works on Linux but fails on Windows
**Solution:** Check for OS-specific code (file paths, line endings, etc.)

**Problem:** Test passes in isolation but fails in suite
**Solution:** Check for test interdependence; ensure proper setup/teardown

**Problem:** Cannot reproduce bug
**Solution:** Review exact reproduction steps; check environment differences

**Problem:** Fix causes new bug (regression)
**Solution:** Roll back fix, add more comprehensive test, re-implement carefully

---

## Documentation Requirements

- [ ] Update `docs/BUG_TRACKER.md` with all resolutions
- [ ] Create `docs/KNOWN_LIMITATIONS.md`
- [ ] Create `docs/BUG_FIX_PLAN.md` with sprint plan
- [ ] Update `docs/PROGRESS.md` with bug fix summary
- [ ] Document any behavior changes in release notes

---

## Next Steps

After completing Step 25:
1. Verify all critical/high bugs resolved
2. Run final smoke test on Linux and Windows
3. Proceed to **Step 26: Windows Installer Creation**

---

**Estimated Completion Time:** 2-3 days
**Phase 5 Progress:** 2 of 5 steps complete
**Overall Progress:** 25 of 28 steps complete (89%)

---

**End of Step 25**
