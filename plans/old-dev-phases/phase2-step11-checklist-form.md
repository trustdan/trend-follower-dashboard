# Phase 2 - Step 11: Checklist Form & Required Gates

**Phase:** 2 - Checklist & Position Sizing
**Step:** 11 of 28, 2 of 5 (Phase 2)
**Duration:** 2-3 days
**Dependencies:** Step 10 (Banner Component)

---

## Objectives

Create the Checklist screen with form inputs and the 5 required gates as checkboxes. Banner updates live as user checks/unchecks items.

---

## Key Components

1. **Checklist form** - ticker, entry price, ATR, sector dropdown, structure dropdown
2. **5 Required Gates** - Large custom checkboxes with gradient when checked
3. **Banner integration** - Display at top, updates live
4. **Pre-calculations** - Show stop distance, initial stop, add-on levels as user enters data
5. **Comprehensive logging** - Log every checkbox change, banner transitions, form field changes

---

## Implementation Summary

**Create** `ui/src/routes/checklist/+page.svelte`:
- Form inputs for ticker, entry, ATR, sector (dropdown), structure (dropdown)
- Display calculated values: stop distance = K × ATR, initial stop, add-on schedule
- 5 required gate checkboxes with custom styling
- Banner component at top driven by checklist store
- Log all interactions (which gate toggled, banner state transitions)
- "Save Evaluation" button (enabled only when banner is GREEN)

**Backend API endpoint** (already exists from domain logic):
- `POST /api/checklist/evaluate` - Save checklist evaluation with timestamp

**Logging requirements:**
- Log each checkbox change (gate name, checked/unchecked, timestamp)
- Log banner state transitions (RED→YELLOW→GREEN, with reasons)
- Log form field changes (ticker, entry, ATR changed)
- Track which gate is typically the last one checked (bottleneck analysis)

---

## Expected Outcome

User fills out form, checks all 5 required gates, banner turns from RED → YELLOW or GREEN, pre-calculated values display automatically, can save evaluation when GREEN.

---

## Time Estimate

~6-10 hours (1-2 days)

---

**Next:** [Step 12: Quality Items & Scoring](phase2-step12-quality-scoring.md)
