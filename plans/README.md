# Trade Session Planning - Overview

**Created:** 2025-10-30
**Status:** Planning Complete, Ready for Review

---

## Problem Identified

You identified a critical UX gap: **the tabs don't marry up with each other**.

Currently:
- Checklist tab has its own ticker input
- Position Sizing tab has its own ticker input
- Heat Check tab has its own ticker input
- Trade Entry tab has its own ticker input

**Result:** No cohesion. No way to know if all tabs are analyzing the same trade. No strategy context. No session concept.

---

## Solution Designed

Introduce a **Trade Session** system:

1. User clicks "Start New Trade"
2. Selects strategy (Long Breakout, Short Breakout, Custom)
3. Session automatically increments (Session #47)
4. FINVIZ provenance (preset, candidate id, scan date) captured with the session
5. All tabs now work on Session #47
6. Progress tracked: ✅ Checklist | ⏳ Sizing | ○ Heat | ○ Entry
7. Full audit trail in database

---

## Planning Documents Created

### 1. [TRADE_SESSION_ARCHITECTURE.md](TRADE_SESSION_ARCHITECTURE.md)
**28 pages** of comprehensive architecture planning:

- **Problem Statement:** Why sessions are needed
- **Design Principles:** How sessions align with anti-impulsivity philosophy
- **Database Schema:** Complete `trade_sessions` table design
- **Session Lifecycle:** 6-step workflow (NEW → CHECKLIST → SIZING → HEAT → ENTRY → COMPLETED)
- **UI Changes:** Top bar, session dropdown, navigation
- **Data Flow Examples:** 3 detailed scenarios
- **Implementation Plan:** 4 phases, 10-15 hours estimated
- **Open Questions:** 4 design decisions to make
- **Success Metrics:** How to measure if it works

**Key Insight:** Sessions are not just a UX improvement - they're a fundamental alignment with the discipline philosophy. Every gate must be completed sequentially, no skipping, full audit trail.

### 2. [SESSION_UI_MOCKUPS.md](SESSION_UI_MOCKUPS.md)
**40+ pages** of detailed UI mockups:

- **10 screen mockups** (ASCII art, fully specified)
  - Dashboard with no session
  - "Start New Trade" dialog
  - Checklist tab (GREEN banner)
  - Checklist tab (RED banner)
  - Position Sizing tab
  - Heat Check tab
  - Trade Entry tab (final gates)
  - Session completed (read-only)
  - "Resume Session" dropdown
  - Session history view

- **Color schemes** (banner colors, session bar colors, progress icons)
- **Keyboard shortcuts** (Ctrl+N, Ctrl+R, Ctrl+H, Ctrl+1-5)
- **Responsive behavior** (how UI adapts to window size)
- **Accessibility** (screen reader, high contrast, keyboard nav)
- **Animation/transitions** (timing, easing functions)
- **Error states** (3 error dialogs with clear messaging)

**Key Insight:** Every pixel specified. No guesswork. Visual hierarchy: Session bar > Banner > Content > Navigation.

### 3. [IMPLEMENTATION_CHECKLIST.md](IMPLEMENTATION_CHECKLIST.md)
**50+ pages** of step-by-step implementation guide:

- **6 phases with 100+ checkboxes:**
  - Phase 1: Database & Backend (2-3 hours)
  - Phase 2: UI Scaffolding (2-3 hours)
  - Phase 3: Tab Integration (3-4 hours)
  - Phase 4: Polish & Testing (2-3 hours)
  - Phase 5: Documentation (1 hour)
  - Phase 6: Release (30 min)

- **E2E test scenarios** (5 complete workflows to verify)
- **Edge case testing** (collision handling, invalid state, rapid switching)
- **Performance testing** (< 100ms create, < 200ms load)
- **Rollback plan** (if implementation fails)
- **Success criteria** (7 criteria to verify completion)

**Key Insight:** Nothing left to chance. Every function, every test, every edge case planned. Timeline: 11-15 hours total.

### 4. Database Migration Script
**[backend/internal/storage/migrations/001_add_trade_sessions.sql](../backend/internal/storage/migrations/001_add_trade_sessions.sql)**

Complete SQL migration:
- Creates `trade_sessions` table
- 7 indexes for performance
- 1 trigger for auto-update timestamps
- Rollback script included
- Verification query

---

## Key Design Decisions to Make

Before implementation, you need to decide:

### 1. Session Number Generation
- **Option A:** Random 1-99 (proposed)
  - Human-memorable ("Session #47")
  - Risk of collision if 99+ sessions
- **Option B:** Sequential (1, 2, 3, ...)
  - No collisions
  - Less memorable ("Session #4732")

**Recommendation:** Sequential auto-increment (based on session id) for clarity and uniqueness.

### 2. Session Persistence Duration
- **Option A:** Keep DRAFT sessions 30 days, auto-archive
- **Option B:** Keep forever, manual cleanup
- **Option C:** Keep COMPLETED forever, delete DRAFT after 7 days

**Recommendation:** 30-day auto-archive for DRAFT, keep COMPLETED forever (audit trail).

### 3. Multi-Session Workflow
- **Option A:** Allow multiple DRAFT sessions (proposed)
  - User can switch between sessions
  - More flexible
- **Option B:** One session at a time
  - Must complete/abandon before starting new
  - Simpler, enforces focus

**Recommendation:** Allow multiple DRAFT, warn if >5 active (prevent clutter). Use session status and updated_at ordering to keep recents on top.

### 4. Session Locking
- **Option A:** Read-only after COMPLETED (proposed)
  - Immutable audit trail
  - Must clone to re-analyze
- **Option B:** Allow edits
  - More flexible
  - Violates audit principle

**Recommendation:** Read-only after COMPLETED. Cloning available for corrections.

### 5. FINVIZ Provenance
- **Option A:** Require preset/candidate linkage on session creation (proposed)
  - Fast jump back to source scan
  - Enables analytics on preset effectiveness
- **Option B:** Treat provenance as optional metadata
  - Simpler dialogs
  - Loses ability to audit scan origin

**Recommendation:** Capture `candidate_id`, `preset_id`, and scan date when available; allow manual override for custom sessions.

---

## What Happens Next?

1. **You review the planning docs**
   - Read [TRADE_SESSION_ARCHITECTURE.md](TRADE_SESSION_ARCHITECTURE.md) first
   - Look at [SESSION_UI_MOCKUPS.md](SESSION_UI_MOCKUPS.md) for visual clarity
   - Check [IMPLEMENTATION_CHECKLIST.md](IMPLEMENTATION_CHECKLIST.md) for feasibility

2. **You make design decisions**
   - Session number approach?
   - Persistence duration?
   - Multi-session workflow?
   - Session locking?

3. **Implementation begins**
   - Follow [IMPLEMENTATION_CHECKLIST.md](IMPLEMENTATION_CHECKLIST.md)
   - Phase 1: Database (2-3 hours)
   - Phase 2: UI Scaffolding (2-3 hours)
   - Phase 3: Tab Integration (3-4 hours)
   - Phase 4: Polish & Testing (2-3 hours)
   - Phase 5: Documentation (1 hour)
   - Phase 6: Release (30 min)

4. **Testing & iteration**
   - Run 5 E2E scenarios
   - Test edge cases
   - Gather feedback
   - Refine as needed

---

## Why This Is Important

From the anti-impulsivity design philosophy:

> **Trade the tide, not the splash**
> Sessions force strategy selection FIRST, then find setups that match.

> **Friction where it matters**
> Must create session before evaluating (intentional friction).

> **Immediate feedback**
> Session bar shows status in real-time (✅ ⏳ ○).

> **Journal while deciding**
> Full session history = automatic journal of evaluation process.

**This is not just a UX improvement. It's a fundamental alignment with the discipline enforcement philosophy.**

Without sessions: isolated calculations, no context, easy to impulsive trade.

With sessions: sequential gates, full audit trail, impossible to skip steps.

---

## Files Created

```
plans/
├── README.md                           ← You are here
├── TRADE_SESSION_ARCHITECTURE.md       ← 28 pages, architecture design
├── SESSION_UI_MOCKUPS.md               ← 40+ pages, UI mockups
└── IMPLEMENTATION_CHECKLIST.md         ← 50+ pages, implementation guide

backend/internal/storage/migrations/
└── 001_add_trade_sessions.sql          ← Database migration script
```

**Total:** 4 files, 120+ pages of planning documentation.

---

## Questions for You

1. **Do you agree with the session concept?**
   - Does it solve the "tabs don't marry up" problem?
   - Does it align with the discipline philosophy?

2. **Do you like the UI mockups?**
   - Session bar at top clear?
   - Progress indicators (✅ ⏳ ○) intuitive?
   - "Start New Trade" dialog makes sense?

3. **Are the design decisions reasonable?**
   - Random 1-99 session numbers OK?
   - 30-day auto-archive for DRAFT sessions OK?
   - Allow multiple DRAFT sessions OK?
   - Read-only after COMPLETED OK?

4. **Is the implementation plan feasible?**
   - 11-15 hours reasonable?
   - Phased approach makes sense?
   - Any concerns about complexity?

5. **Ready to proceed?**
   - Should we start Phase 1 (Database)?
   - Any changes needed first?

---

## Next Steps

**If you approve:**
1. Say "Let's implement Phase 1" or "Start with the database"
2. I'll follow [IMPLEMENTATION_CHECKLIST.md](IMPLEMENTATION_CHECKLIST.md)
3. We'll work through phases sequentially
4. Testing at each step
5. Full deployment in 11-15 hours

**If you have questions:**
1. Ask about any design decision
2. Request changes to mockups
3. Suggest alternative approaches
4. I'll update planning docs accordingly

**If you need time to review:**
1. Read the 3 main docs at your own pace
2. Come back when ready
3. We'll discuss and refine
4. Implement when you're confident

---

## Alignment Check

From CLAUDE.md:

> **When in doubt, read `docs/project/WHY.md` and `docs/anti-impulsivity.md`.**

✅ **Trade sessions support discipline:**
- Force strategy selection first
- Prevent impulsive trades (must go through all gates)
- Create audit trail (know what you were thinking)
- Add intentional friction (2-minute cooloff, sequential gates)

✅ **Trade sessions are as simple as possible:**
- One concept: "A session represents one trade being evaluated"
- Clear lifecycle: NEW → CHECKLIST → SIZING → HEAT → ENTRY → COMPLETED
- No configuration: Hard-coded rules, no bypasses

✅ **Trade sessions fail loudly:**
- Can't skip ahead (UI enforces prerequisites)
- Clear error messages ("Session #47 has RED banner - resolve checklist first")
- Full visibility (session bar always shows progress)

**Recommendation: This design aligns perfectly with the project philosophy. Proceed with confidence.**

---

**Document Version:** 1.0
**Author:** Claude Code Planning Agent
**Status:** Awaiting Your Review & Decision
