# Bug Tracker - Step 25

**Last Updated:** 2025-10-29
**Status:** Bug fixing in progress

---

## Overview

This document tracks all bugs discovered during Step 24 (Comprehensive Testing) and Step 25 (Bug Fixing Sprint).

**Summary:**
- Critical: 0
- High: 0
- Medium: 1 (Test assertion mismatch)
- Low: 0

---

## BUG-001: API Handler Test Assertions Mismatch

**Severity:** MEDIUM
**Status:** IDENTIFIED
**Found:** 2025-10-29 (Step 25)
**Component:** Backend Tests (`internal/api/handlers/*_test.go`)

### Description

The API handler tests expect a response format with a `status` field:

```json
{
  "status": "success",
  "data": {...}
}
```

However, the actual API returns:

```json
{
  "data": {...}
}
```

### Impact

- **Production Code:** ✅ Working correctly - frontend expects and receives correct format
- **Tests:** ❌ Failing - tests expect wrong response structure
- **Severity Rationale:** Medium - tests fail but production code works

### Affected Tests

All success-path tests in:
1. `calendar_test.go` - Calendar data retrieval
2. `candidates_test.go` - Candidate import/retrieval
3. `decisions_test.go` - Decision saving
4. `heat_test.go` - Heat calculations
5. `positions_test.go` - Position retrieval
6. `settings_test.go` - Settings retrieval
7. `sizing_test.go` - Position sizing calculations

### Root Cause

Tests were written with an incorrect assumption about API response format. The actual response format (defined in `internal/api/responses/response.go`) correctly matches what the frontend expects (defined in `ui/src/lib/api/client.ts`).

### Reproduction

```bash
cd backend/
go test ./internal/api/handlers/... -v
```

**Expected:** All tests pass
**Actual:** Success-path tests fail with "Expected status 'success', got '<nil>'"

### Fix Plan

Update all test assertions to match actual API contract:

**Before:**
```go
if response["status"] != "success" {
    t.Errorf("Expected status 'success', got '%v'", response["status"])
}
data, ok := response["data"].(map[string]interface{})
```

**After:**
```go
// Response is ApiResponse[T] which only has "data" field
data, ok := response["data"].(map[string]interface{})
if !ok {
    t.Fatalf("Expected data to be a map")
}
```

### Files to Fix

1. `backend/internal/api/handlers/calendar_test.go`
2. `backend/internal/api/handlers/candidates_test.go`
3. `backend/internal/api/handlers/decisions_test.go`
4. `backend/internal/api/handlers/heat_test.go`
5. `backend/internal/api/handlers/positions_test.go`
6. `backend/internal/api/handlers/settings_test.go`
7. `backend/internal/api/handlers/sizing_test.go`

### Estimated Fix Time

1 hour (systematic find-and-replace across all test files)

---

## Production Verification

### ✅ Working Correctly

1. **Frontend Build:** ✅ Builds successfully (9.51s)
2. **API Contract:** ✅ Backend returns `{data: T}` format
3. **Frontend Client:** ✅ Expects `{data: T}` format (client.ts:6-8)
4. **Domain Tests:** ✅ All passing (96.9% coverage)
5. **Storage Tests:** ✅ All passing (66.7% coverage)
6. **Middleware Tests:** ✅ All passing (100% coverage)

### ⚠️ Needs Fixing

1. **Handler Tests:** ❌ Test assertions don't match actual API

---

## No Critical/High Bugs Found! ✅

**Great news:** Comprehensive testing revealed NO critical or high severity bugs. The system is production-ready. Only test assertion mismatches need fixing.

---

## Known Issues (By Design)

These are intentional limitations, not bugs:

### 1. No Direct Trade Execution
**Reason:** Intentional separation - trader must manually verify and execute

### 2. FINVIZ Scraper Dependent on HTML Structure
**Reason:** Third-party website - we have no control
**Mitigation:** Manual candidate entry always available

### 3. SQLite Single Writer Limitation
**Reason:** SQLite architecture
**Mitigation:** Single-user system by design

---

## Bug History

| Bug ID | Severity | Status | Found | Fixed | Days Open |
|--------|----------|--------|-------|-------|-----------|
| BUG-001 | Medium | Identified | 2025-10-29 | - | 0 |

---

## Testing Summary

**Test Status Before Fixes:**
- Domain: ✅ ALL PASS (96.9% coverage)
- Storage: ✅ ALL PASS (66.7% coverage)
- Middleware: ✅ ALL PASS (100% coverage)
- API Handlers: ❌ PARTIAL (test assertions wrong)

**Expected After Fixes:**
- All tests: ✅ PASS (100%)

---

**Next:** See `BUG_FIX_PLAN.md` for prioritized fix plan
