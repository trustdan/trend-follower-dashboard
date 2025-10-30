# Phase 1 - Step 8: FINVIZ Scanner Implementation

**Phase:** 1 - Dashboard & FINVIZ Scanner
**Step:** 8 of 9 (overall), 4 of 5 (Phase 1)
**Duration:** 2-3 days
**Dependencies:** Step 7 (Dashboard), Step 5 (Backend API with scan endpoint)

---

## Objectives

Create FINVIZ Scanner screen for one-click daily scanning.

---

## Key Components

1. **FINVIZScanner.svelte** - Main scanner component with large "Run Daily Scan" button
2. **PresetSelector.svelte** - Dropdown to select different FINVIZ presets
3. **ScanResults.svelte** - Table displaying scan results (ticker, last close, volume, sector)
4. **LoadingState** - Animated spinner during 3-5 second scan

---

## Implementation

Call `POST /api/candidates/scan` with preset name. Display results in table. Show count of candidates found. Handle errors gracefully (network issues, FINVIZ changes).

Add "Review Candidates" button that navigates to Step 9 (import screen).

---

## Expected Outcome

User clicks "Run Daily Scan", backend scrapes FINVIZ, frontend displays 20+ candidates in table. Success message shows count.

---

**Status:** 📋 Ready for Execution
**Next:** [Step 9: Candidate Import](phase1-step9-candidate-import.md)
