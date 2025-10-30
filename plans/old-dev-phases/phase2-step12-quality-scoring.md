# Phase 2 - Step 12: Quality Items & Scoring

**Phase:** 2 - Checklist & Position Sizing
**Step:** 12 of 28, 3 of 5 (Phase 2)
**Duration:** 1-2 days
**Dependencies:** Step 11 (Checklist Form)

---

## Objectives

Add optional quality items section below required gates. Each item adds 1 point to quality score. Banner transitions from YELLOW to GREEN when score ≥ threshold (default 3).

---

## Key Components

1. **4 Optional Checkboxes:**
   - Regime OK (SPY > 200 SMA for longs)
   - No Chase (entry within 2N of 20-EMA)
   - Earnings OK (no earnings within next 2 weeks)
   - Journal Note (textarea for reasoning)

2. **Quality Score Display** - "Quality Score: 3 / 4" prominently shown
3. **Threshold Display** - "Threshold: 3.0" (configurable in settings)
4. **Visual Distinction** - Blue accent for optional vs red for required

---

## Implementation

Update `ui/src/routes/checklist/+page.svelte`:
- Add quality items section below required gates
- Journal note textarea (counts as 1 point if not empty)
- Display: "Quality Score: X / 4"
- Display: "Threshold: 3.0" (from settings)
- Banner updates when quality score crosses threshold

Update checklist store (already done in Step 10):
- `qualityScore` derived store calculates 0-4
- `bannerState` uses quality score for YELLOW/GREEN logic

---

## Expected Outcome

User checks 3+ quality items (or enters journal note), banner turns GREEN. Quality score visible. Optional items visually distinct from required.

---

## Time Estimate

~4-6 hours (1 day)

---

**Next:** [Step 13: Position Sizing Calculator](phase2-step13-position-sizing.md)
