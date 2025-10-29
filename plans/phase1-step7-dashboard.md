# Phase 1 - Step 7: Dashboard Screen

**Phase:** 1 - Dashboard & FINVIZ Scanner
**Step:** 7 of 9 (overall), 3 of 5 (Phase 1)
**Duration:** 2 days
**Dependencies:** Step 6 (Layout & Navigation), Step 5 (Backend API)

---

## Objectives

Build the Dashboard as the main landing page with key information display.

---

## Key Components to Build

1. **Portfolio Summary Card** - Equity, total heat, available capacity
2. **Positions Table Component** - Open positions with risk, stops, days held
3. **Candidates Summary** - Count for today
4. **Heat Gauge Visual** - Progress bar showing portfolio heat
5. **Quick Actions** - Navigate to Scanner or Checklist

---

## Implementation

Create `ui/src/lib/components/dashboard/PositionTable.svelte`, `HeatGauge.svelte`, `PortfolioCard.svelte`.

Update `ui/src/routes/+page.svelte` to fetch from `/api/settings` and `/api/positions`, display real data.

Use API client from `$lib/api/client.ts` (create this to wrap fetch calls with error handling and logging).

---

## Expected Outcome

Dashboard displays real portfolio data, positions table, heat levels, and candidates count. All data fetched from backend API on page load.

---

**Status:** 📋 Ready for Execution
**Next:** [Step 8: FINVIZ Scanner](phase1-step8-finviz-scanner.md)
