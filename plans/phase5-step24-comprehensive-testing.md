# Phase 5 - Step 24: Comprehensive Testing Suite

**TF = Trend Following** - Systematic trading discipline enforcement system

**Phase:** 5 - Testing & Packaging
**Step:** 24 of 28
**Duration:** 3 days
**Dependencies:** Phase 4 complete (All features implemented and polished)

---

## Objectives

Create and execute a comprehensive testing suite covering unit tests, integration tests, end-to-end tests, and manual testing scenarios. Verify that all features work correctly, edge cases are handled, and the 5 gates enforcement is unbreakable. Test on both Linux (development) and Windows (deployment target).

**Purpose:** Ensure the application is production-ready with zero critical bugs before packaging and release.

---

## Success Criteria

- [ ] Frontend unit tests written (Vitest) for critical components
- [ ] Backend unit tests pass (`go test ./...`)
- [ ] Integration tests cover complete workflows
- [ ] Edge case testing complete (boundary conditions, empty states, errors)
- [ ] Manual test plan executed with all scenarios
- [ ] Day/night mode tested thoroughly
- [ ] Keyboard navigation and accessibility tested
- [ ] Windows functionality verified (cross-compilation test)
- [ ] Backend regression tests pass (no broken functionality)
- [ ] Bug tracking document created with all issues logged
- [ ] All critical bugs identified and documented
- [ ] High-severity bugs documented with reproduction steps

---

## Prerequisites

**All Features Complete:**
- Dashboard, Scanner, Checklist, Position Sizing, Heat Check, Trade Entry, Calendar, TradingView integration
- Theme toggle, keyboard shortcuts, debug panel
- Performance optimizations applied

**Testing Tools:**
- Vitest (frontend unit testing)
- Go testing package (backend)
- Manual testing checklists

---

## Implementation Plan

### Part 1: Frontend Unit Tests (1 day)

#### Task 1.1: Setup Testing Framework (30 min)

**File:** `ui/vite.config.ts`

Configure Vitest:

```typescript
import { defineConfig } from 'vite';
import { sveltekit } from '@sveltejs/kit/vite';

export default defineConfig({
    plugins: [sveltekit()],
    test: {
        include: ['src/**/*.{test,spec}.{js,ts}'],
        environment: 'jsdom',
        globals: true,
        setupFiles: ['./src/setupTests.ts']
    }
});
```

**File:** `ui/src/setupTests.ts`

```typescript
import { expect, afterEach } from 'vitest';
import { cleanup } from '@testing-library/svelte';
import '@testing-library/jest-dom';

// Auto-cleanup after each test
afterEach(() => {
    cleanup();
});
```

Install dependencies:

```bash
cd ui/
npm install -D vitest @testing-library/svelte @testing-library/jest-dom jsdom
```

#### Task 1.2: Banner Component Tests (1 hour)

**File:** `ui/src/lib/components/Banner.test.ts`

```typescript
import { render, screen } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import Banner from './Banner.svelte';

describe('Banner Component', () => {
    it('displays RED state correctly', () => {
        render(Banner, {
            props: {
                state: 'RED',
                message: 'DO NOT TRADE',
                details: 'Required gates failed'
            }
        });

        const banner = screen.getByText(/DO NOT TRADE/i);
        expect(banner).toBeInTheDocument();
        expect(banner).toHaveClass('banner-red');
    });

    it('displays YELLOW state correctly', () => {
        render(Banner, {
            props: {
                state: 'YELLOW',
                message: 'CAUTION',
                details: 'Quality score below threshold'
            }
        });

        const banner = screen.getByText(/CAUTION/i);
        expect(banner).toBeInTheDocument();
        expect(banner).toHaveClass('banner-yellow');
    });

    it('displays GREEN state correctly', () => {
        render(Banner, {
            props: {
                state: 'GREEN',
                message: 'OK TO TRADE',
                details: 'All gates pass'
            }
        });

        const banner = screen.getByText(/OK TO TRADE/i);
        expect(banner).toBeInTheDocument();
        expect(banner).toHaveClass('banner-green');
    });

    it('transitions smoothly between states', async () => {
        const { component } = render(Banner, {
            props: { state: 'RED', message: 'Test', details: '' }
        });

        // Change state
        await component.$set({ state: 'GREEN' });

        // Verify transition
        const banner = screen.getByText(/Test/i);
        expect(banner).toHaveClass('banner-green');
    });
});
```

#### Task 1.3: Checklist Logic Tests (1 hour)

**File:** `ui/src/lib/utils/checklist.test.ts`

```typescript
import { describe, it, expect } from 'vitest';
import { calculateBannerState, calculateQualityScore } from './checklist';

describe('Checklist Logic', () => {
    describe('calculateBannerState', () => {
        it('returns RED when any required gate is false', () => {
            const requiredGates = {
                signal: false,
                riskSize: true,
                liquidity: true,
                exits: true,
                behavior: true
            };

            const result = calculateBannerState(requiredGates, 4);
            expect(result).toBe('RED');
        });

        it('returns YELLOW when all required pass but quality score low', () => {
            const requiredGates = {
                signal: true,
                riskSize: true,
                liquidity: true,
                exits: true,
                behavior: true
            };

            const result = calculateBannerState(requiredGates, 2); // score < 3
            expect(result).toBe('YELLOW');
        });

        it('returns GREEN when all required pass and quality score >= 3', () => {
            const requiredGates = {
                signal: true,
                riskSize: true,
                liquidity: true,
                exits: true,
                behavior: true
            };

            const result = calculateBannerState(requiredGates, 3);
            expect(result).toBe('GREEN');
        });
    });

    describe('calculateQualityScore', () => {
        it('calculates score correctly', () => {
            const qualityItems = {
                regime: true,
                noChase: true,
                earnings: false,
                journal: 'Test note'
            };

            const score = calculateQualityScore(qualityItems);
            expect(score).toBe(3); // regime + noChase + journal = 3
        });

        it('returns 0 when all unchecked', () => {
            const qualityItems = {
                regime: false,
                noChase: false,
                earnings: false,
                journal: ''
            };

            const score = calculateQualityScore(qualityItems);
            expect(score).toBe(0);
        });
    });
});
```

#### Task 1.4: Position Sizing Tests (1 hour)

**File:** `ui/src/lib/utils/sizing.test.ts`

```typescript
import { describe, it, expect } from 'vitest';
import { calculateShares, calculateAddOnLevels } from './sizing';

describe('Position Sizing Calculations', () => {
    describe('calculateShares', () => {
        it('calculates shares correctly for stock', () => {
            const result = calculateShares({
                riskDollars: 750,
                stopDistance: 4.70,
                method: 'stock'
            });

            expect(result.shares).toBe(159); // floor(750 / 4.70)
            expect(result.actualRisk).toBe(747.30); // 159 * 4.70
            expect(result.actualRisk).toBeLessThanOrEqual(750);
        });

        it('handles zero risk correctly', () => {
            const result = calculateShares({
                riskDollars: 0,
                stopDistance: 4.70,
                method: 'stock'
            });

            expect(result.shares).toBe(0);
            expect(result.actualRisk).toBe(0);
        });

        it('handles very small stop distance', () => {
            const result = calculateShares({
                riskDollars: 750,
                stopDistance: 0.10,
                method: 'stock'
            });

            expect(result.shares).toBe(7500); // floor(750 / 0.10)
        });
    });

    describe('calculateAddOnLevels', () => {
        it('calculates add-on prices correctly', () => {
            const levels = calculateAddOnLevels({
                entry: 180.50,
                atr: 2.35,
                addStepN: 0.5,
                maxUnits: 4
            });

            expect(levels).toHaveLength(3); // maxUnits - 1
            expect(levels[0]).toBeCloseTo(181.68, 2); // 180.50 + (0.5 * 2.35)
            expect(levels[1]).toBeCloseTo(182.85, 2); // 180.50 + (1.0 * 2.35)
            expect(levels[2]).toBeCloseTo(184.03, 2); // 180.50 + (1.5 * 2.35)
        });
    });
});
```

#### Task 1.5: Run Frontend Tests (15 min)

```bash
cd ui/
npm run test

# With coverage
npm run test -- --coverage
```

**Target:** 80% code coverage for critical utilities and components.

---

### Part 2: Backend Tests (1 day)

#### Task 2.1: Verify Existing Tests Pass (30 min)

```bash
cd backend/

# Run all tests
go test ./... -v

# Run with coverage
go test ./... -cover

# Run specific package
go test ./internal/domain/ -v
go test ./internal/storage/ -v
```

**Expected:** All existing tests should pass. If any fail, document and fix.

#### Task 2.2: Add Integration Tests (2 hours)

**File:** `backend/internal/server/integration_test.go`

```go
package server_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "your-module/internal/server"
    "your-module/internal/storage"
)

func setupTestServer(t *testing.T) (*server.Server, *storage.DB) {
    // Create test database
    db, err := storage.NewDB(":memory:")
    require.NoError(t, err)

    // Initialize server
    srv := server.NewServer(db)

    return srv, db
}

func TestCompleteWorkflow(t *testing.T) {
    srv, db := setupTestServer(t)
    defer db.Close()

    // Step 1: Get settings
    t.Run("Get initial settings", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/api/settings", nil)
        w := httptest.NewRecorder()

        srv.GetSettings(w, req)

        assert.Equal(t, http.StatusOK, w.Code)

        var settings map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &settings)
        require.NoError(t, err)

        assert.Equal(t, float64(100000), settings["equity"])
    })

    // Step 2: Evaluate checklist
    t.Run("Evaluate checklist - GREEN", func(t *testing.T) {
        body := map[string]interface{}{
            "ticker":     "AAPL",
            "entry":      180.50,
            "atr":        2.35,
            "sector":     "Tech/Comm",
            "required": map[string]bool{
                "signal":    true,
                "riskSize":  true,
                "liquidity": true,
                "exits":     true,
                "behavior":  true,
            },
            "quality": map[string]interface{}{
                "regime":   true,
                "noChase":  true,
                "earnings": true,
                "journal":  "Test trade",
            },
        }

        jsonBody, _ := json.Marshal(body)
        req := httptest.NewRequest("POST", "/api/checklist/evaluate", bytes.NewBuffer(jsonBody))
        w := httptest.NewRecorder()

        srv.EvaluateChecklist(w, req)

        assert.Equal(t, http.StatusOK, w.Code)

        var result map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &result)
        require.NoError(t, err)

        assert.Equal(t, "GREEN", result["banner"])
        assert.Equal(t, float64(0), result["missing_count"])
        assert.Equal(t, float64(4), result["quality_score"])
    })

    // Step 3: Calculate position sizing
    t.Run("Calculate position size", func(t *testing.T) {
        body := map[string]interface{}{
            "ticker": "AAPL",
            "entry":  180.50,
            "atr":    2.35,
            "method": "stock",
            "k":      2.0,
        }

        jsonBody, _ := json.Marshal(body)
        req := httptest.NewRequest("POST", "/api/size/calculate", bytes.NewBuffer(jsonBody))
        w := httptest.NewRecorder()

        srv.CalculateSize(w, req)

        assert.Equal(t, http.StatusOK, w.Code)

        var result map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &result)
        require.NoError(t, err)

        assert.Equal(t, float64(159), result["shares"])
        assert.InDelta(t, 747.30, result["actual_risk"], 0.01)
    })

    // Step 4: Check heat caps
    t.Run("Check heat - within caps", func(t *testing.T) {
        body := map[string]interface{}{
            "ticker":      "AAPL",
            "risk_amount": 750,
            "bucket":      "Tech/Comm",
        }

        jsonBody, _ := json.Marshal(body)
        req := httptest.NewRequest("POST", "/api/heat/check", bytes.NewBuffer(jsonBody))
        w := httptest.NewRecorder()

        srv.CheckHeat(w, req)

        assert.Equal(t, http.StatusOK, w.Code)

        var result map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &result)
        require.NoError(t, err)

        assert.False(t, result["portfolio_cap_exceeded"].(bool))
        assert.False(t, result["bucket_cap_exceeded"].(bool))
    })

    // Step 5: Check all gates
    t.Run("Check gates - all pass", func(t *testing.T) {
        // Wait 2+ minutes (simulate timer)
        // In real test, we'd mock the timer

        body := map[string]interface{}{
            "ticker":       "AAPL",
            "banner":       "GREEN",
            "risk_dollars": 750,
            "bucket":       "Tech/Comm",
        }

        jsonBody, _ := json.Marshal(body)
        req := httptest.NewRequest("POST", "/api/gates/check", bytes.NewBuffer(jsonBody))
        w := httptest.NewRecorder()

        srv.CheckGates(w, req)

        assert.Equal(t, http.StatusOK, w.Code)

        var result map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &result)
        require.NoError(t, err)

        // Note: Gate 2 (impulse brake) will fail in this test
        // unless we mock time or actually wait
        // For unit test, we'd mock the time check
    })
}

func TestHeatCapEnforcement(t *testing.T) {
    srv, db := setupTestServer(t)
    defer db.Close()

    t.Run("Reject trade exceeding portfolio cap", func(t *testing.T) {
        // Add positions to reach 95% of cap
        // ... (add test data)

        // Attempt to add position that exceeds cap
        body := map[string]interface{}{
            "risk_amount": 500, // Would exceed cap
            "bucket":      "Tech/Comm",
        }

        jsonBody, _ := json.Marshal(body)
        req := httptest.NewRequest("POST", "/api/heat/check", bytes.NewBuffer(jsonBody))
        w := httptest.NewRecorder()

        srv.CheckHeat(w, req)

        var result map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &result)

        assert.True(t, result["portfolio_cap_exceeded"].(bool))
    })
}
```

Run integration tests:

```bash
go test ./internal/server/ -v -run TestCompleteWorkflow
go test ./internal/server/ -v -run TestHeatCapEnforcement
```

#### Task 2.3: Edge Case Tests (2 hours)

**File:** `backend/internal/domain/sizing_edge_test.go`

```go
package domain_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "your-module/internal/domain"
)

func TestPositionSizing_EdgeCases(t *testing.T) {
    t.Run("Zero equity", func(t *testing.T) {
        result, err := domain.CalculateSizeStock("AAPL", 180.0, 2.0, 2.0, 0.0075, 0)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "equity must be positive")
    })

    t.Run("Negative ATR", func(t *testing.T) {
        result, err := domain.CalculateSizeStock("AAPL", 180.0, -1.0, 2.0, 0.0075, 100000)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "ATR must be positive")
    })

    t.Run("Zero ATR", func(t *testing.T) {
        result, err := domain.CalculateSizeStock("AAPL", 180.0, 0, 2.0, 0.0075, 100000)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "ATR must be positive")
    })

    t.Run("Very large position", func(t *testing.T) {
        // Entry $10, ATR $0.01, large equity
        result, err := domain.CalculateSizeStock("PENNY", 10.0, 0.01, 2.0, 0.01, 10000000)
        assert.NoError(t, err)

        // Should calculate shares, but verify it's reasonable
        assert.Greater(t, result.Shares, 0)
        assert.LessOrEqual(t, result.Entry*float64(result.Shares), 10000000.0*0.5) // Not more than 50% of equity
    })

    t.Run("Entry price at zero", func(t *testing.T) {
        result, err := domain.CalculateSizeStock("TEST", 0, 2.0, 2.0, 0.0075, 100000)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "entry must be positive")
    })
}
```

---

### Part 3: Manual Testing (1 day)

#### Task 3.1: Create Manual Test Plan (1 hour)

**File:** `docs/MANUAL_TEST_PLAN.md`

```markdown
# Manual Test Plan

## Test Environment

- OS: Windows 10
- Browser: Chrome (latest)
- Database: Fresh trading.db (delete existing)
- Date: _______

## Test Scenarios

### Scenario 1: Complete Happy Path

**Objective:** Verify full workflow from scan to GO decision

**Steps:**
1. [ ] Start app → Opens browser to localhost:8080
2. [ ] Navigate to Dashboard → Displays empty state correctly
3. [ ] Click "Run Daily FINVIZ Scan" → Returns candidates
4. [ ] Review candidates → Select 10 tickers
5. [ ] Click "Import Selected" → Success notification shown
6. [ ] Click "Open in TradingView" for AAPL → Opens in new tab
7. [ ] Return to app → Navigate to Checklist
8. [ ] Enter: AAPL, 180.50, 2.35, Tech/Comm
9. [ ] Check all 5 required gates → Banner turns YELLOW
10. [ ] Check 3 quality items + add journal note → Banner turns GREEN
11. [ ] Click "Save Evaluation" → Timer starts (2:00)
12. [ ] Navigate to Position Sizing
13. [ ] Click "Calculate Position Size" → Shows 159 shares, $747 risk
14. [ ] Click "Save Position Plan" → Success notification
15. [ ] Navigate to Heat Check
16. [ ] Click "Check Heat" → Shows within caps, GREEN result
17. [ ] Wait for timer to reach 0:00
18. [ ] Navigate to Trade Entry
19. [ ] Review trade summary → All data correct
20. [ ] Click "Run Final Gate Check" → All 5 gates PASS
21. [ ] Click "SAVE GO DECISION" (enabled) → Success notification
22. [ ] Navigate to Dashboard → Position appears in "Ready to Execute"
23. [ ] Navigate to Calendar → AAPL appears in Tech/Comm row, current week

**Expected:** All steps complete without errors. GO decision saved.

---

### Scenario 2: Gate Failures

#### 2A: Banner RED (missing required gate)

**Steps:**
1. [ ] Open Checklist
2. [ ] Enter ticker data
3. [ ] Check only 4 of 5 required gates → Banner RED
4. [ ] Attempt to save evaluation → Button disabled
5. [ ] Navigate to Trade Entry → "Complete checklist first"

**Expected:** Cannot proceed with RED banner

#### 2B: Impulse Brake Not Elapsed

**Steps:**
1. [ ] Complete checklist (GREEN)
2. [ ] Save evaluation → Timer starts
3. [ ] Immediately navigate to Trade Entry
4. [ ] Click "Run Final Gate Check"
5. [ ] Observe Gate 2 status → FAIL (timer still active)
6. [ ] "SAVE GO DECISION" button → DISABLED

**Expected:** Gate 2 fails, cannot save GO decision

#### 2C: Portfolio Heat Cap Exceeded

**Steps:**
1. [ ] Add test data: 3 positions totaling $3,800 heat
2. [ ] Start new checklist for $500 risk trade
3. [ ] Complete checklist (GREEN), sizing, wait 2 min
4. [ ] Navigate to Heat Check
5. [ ] Click "Check Heat"
6. [ ] Observe result → Portfolio cap exceeded (RED warning)
7. [ ] Navigate to Trade Entry → Run gates → Gate 4 FAIL

**Expected:** Heat cap prevents GO decision

#### 2D: Ticker on Cooldown

**Steps:**
1. [ ] Add cooldown for AAPL (until next week)
2. [ ] Start checklist for AAPL
3. [ ] Enter ticker → Warning shown: "AAPL on cooldown until..."
4. [ ] Try to proceed → Gate 3 fails

**Expected:** Cooldown prevents trade on AAPL

---

### Scenario 3: NO-GO Decision Logging

**Steps:**
1. [ ] Complete analysis for ticker XYZ
2. [ ] Heat check shows portfolio at 98% of cap
3. [ ] Navigate to Trade Entry → Run gates → Gate 4 FAIL
4. [ ] Click "SAVE NO-GO DECISION"
5. [ ] Enter reason: "Portfolio heat at cap"
6. [ ] Submit → Success notification
7. [ ] Check database → NO-GO decision logged

**Expected:** NO-GO saved with reason

---

### Scenario 4: Theme Toggle

**Steps:**
1. [ ] Open app (day mode default)
2. [ ] Navigate through all screens → Visual check
3. [ ] Click theme toggle (sun icon) → Switches to night mode
4. [ ] Verify all screens update correctly (dark backgrounds, light text)
5. [ ] Check banner gradients → Still vibrant in night mode
6. [ ] Check buttons → Hover effects work
7. [ ] Refresh page → Theme persists (localStorage)
8. [ ] Toggle back to day mode → All correct

**Expected:** Smooth theme transitions, no broken styles

---

### Scenario 5: Keyboard Navigation

**Steps:**
1. [ ] Open Checklist
2. [ ] Press Tab → Focus moves through inputs in order
3. [ ] Type in ticker field
4. [ ] Press Ctrl+K → Ticker input focused
5. [ ] Fill form, press Ctrl+S → Saves evaluation
6. [ ] Open modal (any), press Escape → Modal closes
7. [ ] Navigate forms with Tab → All inputs accessible
8. [ ] Press Enter in form → Submits (where appropriate)

**Expected:** All keyboard shortcuts work

---

### Scenario 6: Error Handling

#### 6A: Network Error During Scan

**Steps:**
1. [ ] Disconnect network (or block localhost:8080 with firewall)
2. [ ] Click "Run Daily FINVIZ Scan"
3. [ ] Observe error notification → Clear message, red gradient
4. [ ] Reconnect network
5. [ ] Retry scan → Works correctly

**Expected:** Graceful error handling, helpful message

#### 6B: Invalid Input Validation

**Steps:**
1. [ ] Open Position Sizing
2. [ ] Enter negative ATR: -2.0
3. [ ] Click Calculate → Error shown: "ATR must be positive"
4. [ ] Enter ATR: 0
5. [ ] Click Calculate → Error shown: "ATR must be positive"
6. [ ] Enter valid ATR: 2.35
7. [ ] Click Calculate → Success

**Expected:** Input validation works, clear error messages

---

### Scenario 7: Debug Panel (Dev Mode)

**Steps:**
1. [ ] Press Ctrl+Shift+D → Debug panel slides in
2. [ ] Navigate through app → Logs appear in panel
3. [ ] Filter logs: "All" → "Errors" → "Warnings" → Works correctly
4. [ ] Click "Clear" → Logs cleared
5. [ ] Perform actions → New logs appear
6. [ ] Click "Export" → JSON file downloads
7. [ ] Open JSON → Logs are valid, readable
8. [ ] Close panel (X button) → Slides out

**Expected:** Debug panel fully functional

---

### Scenario 8: Calendar View

**Steps:**
1. [ ] Add 3 positions in different sectors
2. [ ] Navigate to Calendar
3. [ ] Verify 10 weeks shown (2 back, 8 forward)
4. [ ] Current week highlighted → Correct
5. [ ] Positions appear in correct cells (sector × week)
6. [ ] Hover over ticker → Tooltip shows entry, risk, expected exit
7. [ ] Check color coding → Matches heat levels
8. [ ] Click ticker → Opens TradingView (optional feature)

**Expected:** Calendar displays correctly, all interactions work

---

### Scenario 9: Long-Running Session (30 min)

**Objective:** Verify no memory leaks or performance degradation

**Steps:**
1. [ ] Start app, open DevTools (Performance tab)
2. [ ] Take heap snapshot (baseline)
3. [ ] Navigate through all screens 5 times
4. [ ] Perform 5 complete workflows (scan → decision)
5. [ ] Toggle theme 10 times
6. [ ] Open/close debug panel 5 times
7. [ ] After 30 min, take another heap snapshot
8. [ ] Compare heap sizes → Should be similar (< 10MB growth)
9. [ ] Check CPU usage → Should be low when idle

**Expected:** No significant memory leaks, stable performance

---

### Scenario 10: Data Persistence

**Steps:**
1. [ ] Complete full workflow (save GO decision)
2. [ ] Close browser tab
3. [ ] Close backend (Ctrl+C)
4. [ ] Restart backend
5. [ ] Open browser to localhost:8080
6. [ ] Navigate to Dashboard → Position still shown
7. [ ] Navigate to Calendar → Position still appears
8. [ ] Check theme → Theme preference persisted
9. [ ] Open Debug Panel → Old logs cleared (expected)

**Expected:** All database data persists, localStorage works

---

## Edge Cases Checklist

- [ ] Empty database (first run)
- [ ] No candidates found (FINVIZ returns 0)
- [ ] Exactly at heat cap (99.99% vs 100.00%)
- [ ] Maximum positions (50+ open positions)
- [ ] Very long ticker symbols (10+ chars)
- [ ] Special characters in journal notes
- [ ] Large numbers (equity = $10,000,000)
- [ ] Small numbers (penny stocks, entry = $0.50)
- [ ] Rapid clicking (double-submit prevention)
- [ ] Concurrent tabs (two browsers, same database)

---

## Regression Tests (Features from Earlier Phases)

- [ ] Dashboard displays positions correctly
- [ ] FINVIZ scanner still works
- [ ] Candidates import successfully
- [ ] Settings persist and update
- [ ] Position sizing calculations accurate
- [ ] Heat caps enforce correctly
- [ ] All 5 gates validate properly
- [ ] TradingView links open correctly

---

## Windows-Specific Tests

**After cross-compilation:**

1. [ ] .exe runs without errors
2. [ ] Browser opens automatically
3. [ ] Database created in correct location (AppData)
4. [ ] Theme toggle works
5. [ ] All features function identically to Linux
6. [ ] FINVIZ scan works (network access)
7. [ ] Close app cleanly (no zombie processes)

---

## Accessibility Tests

- [ ] Screen reader labels present on inputs
- [ ] Tab order logical
- [ ] Focus indicators visible
- [ ] Color contrast meets WCAG AA (4.5:1)
- [ ] All interactive elements keyboard-accessible
- [ ] Alt text on images/icons
- [ ] Form validation announces errors

---

## Test Results Summary

**Date Tested:** _______
**Tester:** _______

**Total Scenarios:** 10
**Passed:** _______
**Failed:** _______

**Critical Bugs Found:** _______
**High Bugs Found:** _______
**Medium Bugs Found:** _______
**Low Bugs Found:** _______

**Notes:**
```

#### Task 3.2: Execute Manual Test Plan (4-6 hours)

Execute all scenarios, document results, log bugs.

#### Task 3.3: Windows Cross-Compilation Test (1 hour)

```bash
cd backend/

# Cross-compile to Windows
GOOS=windows GOARCH=amd64 go build -o tf-engine.exe cmd/tf-engine/main.go

# Copy to Windows machine (or test in Wine on Linux)
# Run on Windows and verify all functionality
```

---

### Part 4: Bug Tracking (Ongoing)

**File:** `docs/BUG_TRACKER.md`

```markdown
# Bug Tracker

## Critical Bugs (App Crashes or Data Loss)

### BUG-001: [Example] App crashes when...
**Severity:** Critical
**Found:** 2025-10-29
**Steps to Reproduce:**
1. Step 1
2. Step 2
**Expected:** ...
**Actual:** ...
**Status:** Open / In Progress / Fixed

---

## High Bugs (Feature Broken, No Workaround)

### BUG-101: [Example] Cannot save decision when...
**Severity:** High
**Found:** ...

---

## Medium Bugs (Feature Broken, Has Workaround)

---

## Low Bugs (Cosmetic, Minor UX Issues)

---

## Resolved Bugs

```

---

## Testing Checklist

### Frontend Unit Tests
- [ ] Banner component tests pass
- [ ] Checklist logic tests pass
- [ ] Position sizing calculation tests pass
- [ ] Code coverage > 80% for utilities

### Backend Tests
- [ ] All existing tests pass (`go test ./...`)
- [ ] Integration tests pass (complete workflow)
- [ ] Edge case tests pass (zero values, negative numbers, etc.)
- [ ] Heat cap enforcement tested

### Manual Testing
- [ ] Happy path (scan → GO decision) works
- [ ] All gate failures tested and working
- [ ] NO-GO decision logging works
- [ ] Theme toggle works perfectly
- [ ] Keyboard shortcuts work
- [ ] Error handling graceful
- [ ] Debug panel functional
- [ ] Calendar displays correctly
- [ ] No memory leaks in 30-min session
- [ ] Data persists after restart

### Windows Testing
- [ ] .exe runs on Windows
- [ ] All features work on Windows
- [ ] Performance is acceptable

### Accessibility
- [ ] Keyboard navigation works
- [ ] Screen reader compatible
- [ ] Color contrast meets standards

---

## Documentation Requirements

- [ ] Create `docs/MANUAL_TEST_PLAN.md` with all scenarios
- [ ] Create `docs/BUG_TRACKER.md` for issue logging
- [ ] Document test results in `docs/TEST_RESULTS.md`
- [ ] Update `docs/PROGRESS.md` with testing status

---

## Next Steps

After completing Step 24:
1. Review all bugs found
2. Prioritize bugs by severity
3. Proceed to **Step 25: Bug Fixing Sprint**

---

**Estimated Completion Time:** 3 days
**Phase 5 Progress:** 1 of 5 steps complete
**Overall Progress:** 24 of 28 steps complete (86%)

---

**End of Step 24**
