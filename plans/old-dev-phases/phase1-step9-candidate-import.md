# Phase 1 - Step 9: Candidate Import & Review

**Phase:** 1 - Dashboard & FINVIZ Scanner
**Step:** 9 of 9 (overall), 5 of 5 (Phase 1)
**Duration:** 2 days
**Dependencies:** Step 8 (FINVIZ Scanner)

---

## Objectives

Build candidate import workflow with review and selection capabilities.

---

## Key Components

1. **CandidateReviewTable.svelte** - Table with checkboxes for each ticker
2. **SectorDistribution.svelte** - Show sector breakdown (e.g., "3 Tech/Comm, 2 Energy")
3. **CooldownIndicator** - Gray out tickers on cooldown (not selectable)
4. **ImportButton** - Large gradient button, disabled until ≥1 ticker selected
5. **Select All / Deselect All** buttons

---

## Implementation

Display scan results from Step 8. User checks/unchecks tickers. Show count: "12 of 23 selected". Call `POST /api/candidates/import` with selected tickers. Show success notification. Update Dashboard to reflect new candidates.

---

## Expected Outcome

User reviews 23 scan results, selects 12, clicks "Import Selected". Backend saves to database. Dashboard shows "12 candidates for 2025-10-29". Success notification: "✓ 12 candidates imported".

---

## Phase 1 Complete!

After this step, Phase 1 (Dashboard & FINVIZ Scanner) is complete. Users can:
- View dashboard with portfolio summary
- Run FINVIZ scans
- Import candidates for analysis

**Next Phase:** [Phase 2: Checklist & Position Sizing](phase2-step10-banner-component.md)

---

**Status:** 📋 Ready for Execution
**Created:** 2025-10-29
