# HTTP/CLI Parity Analysis

## Test Date: 2025-10-27

### Test 1: Position Sizing (POST /api/size)

**HTTP Request:**
```json
{
  "equity": 10000,
  "risk_pct": 0.0075,
  "entry": 180,
  "atr": 1.5,
  "k": 2,
  "method": "stock"
}
```

**HTTP Response:**
```json
{
  "actual_risk": 75,
  "contracts": 0,
  "correlation_id": "670913b0-5085-437c-b84f-8b078aaab721",
  "initial_stop": 177,
  "method": "stock",
  "risk_dollars": 75,
  "shares": 25,
  "stop_distance": 3
}
```

**CLI Response:**
```json
{
  "risk_dollars": 75,
  "stop_distance": 3,
  "initial_stop": 177,
  "shares": 25,
  "contracts": 0,
  "actual_risk": 75,
  "method": "stock"
}
```

**Differences:**
- ✅ All data fields match
- ⚠️  HTTP includes `correlation_id` (extra field)
- ✅ Field order different but semantically identical

**Status:** PASS (with minor difference)

---

### Test 2: Settings (GET /api/settings)

**HTTP Response:**
```json
{
  "BucketHeatCap_pct": 0.015,
  "Equity_E": 10000,
  "HeatCap_H_pct": 0.04,
  "RiskPct_r": 0.0075,
  "StopMultiple_K": 2,
  "correlation_id": "..."
}
```

**CLI Response:**
```json
{
  "BucketHeatCap_pct": "0.015",
  "Equity_E": "10000",
  "HeatCap_H_pct": "0.04",
  "RiskPct_r": "0.0075",
  "StopMultiple_K": "2"
}
```

**Differences:**
- ❌ HTTP returns numbers, CLI returns strings
- ⚠️  HTTP includes `correlation_id`

**Status:** FAIL (type mismatch)
**Impact:** VBA parsing may need to handle both types

---

## Summary

| Endpoint | Parity Status | Notes |
|----------|---------------|-------|
| POST /api/size | ✅ PASS | Minor: correlation_id added |
| GET /api/settings | ❌ FAIL | Type mismatch: numbers vs strings |

## Issues Found

### Issue 1: correlation_id in HTTP responses

**Problem:** HTTP responses include `correlation_id` field, CLI doesn't
**Impact:** LOW - Extra field doesn't break parsing
**Recommendation:** Either add to CLI or document as HTTP-only

### Issue 2: Settings type mismatch

**Problem:** HTTP returns numbers, CLI returns strings
**Impact:** MEDIUM - VBA may need type conversion
**Recommendation:** Standardize on one format (suggest strings for precision)

## Recommendation

**For M19 (VBA):**
- Use HTTP API OR CLI with awareness of differences
- If using HTTP: expect `correlation_id` in responses
- If using CLI: settings are strings, parse as needed
- Both work, just handle the minor differences

**For M20 (Polish):**
- Standardize settings format (suggest strings)
- Decide on correlation_id strategy
- Add HTTP parity tests to CI

