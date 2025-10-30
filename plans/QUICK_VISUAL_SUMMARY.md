# Trade Session - Quick Visual Summary

**For the impatient:** This is the 5-minute version. Read this first, then dive into the detailed docs.

---

## The Problem (Before)

```
┌─────────────────────────────────────────────┐
│ TF-Engine Dashboard                         │
├─────────────────────────────────────────────┤
│ [Checklist] [Position Sizing] [Heat Check]  │
├─────────────────────────────────────────────┤
│                                             │
│  Checklist Tab:                             │
│  Ticker: [____]  ← User types AAPL          │
│  Evaluate → GREEN banner                    │
│                                             │
└─────────────────────────────────────────────┘

User switches to Position Sizing tab...

┌─────────────────────────────────────────────┐
│ TF-Engine Dashboard                         │
├─────────────────────────────────────────────┤
│ [Checklist] [Position Sizing] [Heat Check]  │
├─────────────────────────────────────────────┤
│                                             │
│  Position Sizing Tab:                       │
│  Ticker: [____]  ← Empty! Where's AAPL?     │
│  Entry: [____]                              │
│  Calculate → ???                            │
│                                             │
└─────────────────────────────────────────────┘

❌ Problems:
- Ticker lost when switching tabs
- No connection between Checklist and Sizing
- Which strategy are we evaluating?
- Is this the same trade or a different one?
- No audit trail of the evaluation process
```

---

## The Solution (After)

```
┌──────────────────────────────────────────────────────────┐
│ TF-Engine Dashboard                                      │
├──────────────────────────────────────────────────────────┤
│ ╔════════════════════════════════════════════════════╗  │
│ ║ Session #47 • LONG_BREAKOUT • AAPL                 ║  │
│ ║ ✅ Checklist | ⏳ Sizing | ○ Heat | ○ Entry         ║  │
│ ╚════════════════════════════════════════════════════╝  │
│                                                          │
│ [Start New] [Resume ▼] [History]  [Theme] [Help] [VIM]  │
├──────────────────────────────────────────────────────────┤
│ [Checklist] [Position Sizing] [Heat Check] [Trade Entry]│
├──────────────────────────────────────────────────────────┤
│                                                          │
│  Position Sizing Tab:                                    │
│  Session: #47 (Long Breakout)                            │
│  Ticker: AAPL              ← Auto-filled from session!   │
│  Entry: [180.00]                                         │
│  ATR: [1.50]                                             │
│  Calculate → 25 shares, $75 risk                         │
│                                                          │
│  [Next: Heat Check →]                                    │
│                                                          │
└──────────────────────────────────────────────────────────┘

✅ Benefits:
- All tabs work on same session (#47)
- Ticker, strategy, banner all shared
- Progress visible at all times (session bar)
- Full audit trail in database
- Sequential workflow enforced
```

---

## How It Works (User Flow)

### Step 1: Start New Trade

```
User clicks: [Start New Trade]

         ┌──────────────────────────┐
         │ Start New Trade          │
         ├──────────────────────────┤
         │ Strategy:                │
         │  ⚫ Long Breakout         │
         │  ○ Short Breakout        │
         │  ○ Custom                │
         │                          │
         │ Ticker: [AAPL___]        │
         │                          │
         │ [Create] [Cancel]        │
         └──────────────────────────┘

Result: Session #47 created (provenance from preset/candidate stored), navigate to Checklist
```

### Step 2: Evaluate Checklist

```
┌──────────────────────────────────────────────┐
│ Session #47 • LONG_BREAKOUT • AAPL           │
│ ⏳ Checklist | ○ Sizing | ○ Heat | ○ Entry   │
├──────────────────────────────────────────────┤
│ Checklist Tab:                               │
│                                              │
│ Ticker: AAPL    [from session, auto-filled]  │
│                                              │
│ ☑ From Preset                                │
│ ☑ Trend Confirmed                            │
│ ☑ Liquidity OK                               │
│ ☑ TV Confirm                                 │
│ ☑ Earnings OK                                │
│                                              │
│ [Evaluate]                                   │
└──────────────────────────────────────────────┘

User clicks Evaluate →

┌──────────────────────────────────────────────┐
│ ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓  │
│ ┃         ✅ GREEN - GO                  ┃  │
│ ┃     All Required Gates Pass            ┃  │
│ ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛  │
│                                              │
│ [Next: Position Sizing →]                    │
└──────────────────────────────────────────────┘

Session updated:
- checklist_completed = 1
- checklist_banner = 'GREEN'
- Progress: ✅ Checklist (done!)
```

### Step 3: Calculate Position Size

```
┌──────────────────────────────────────────────┐
│ Session #47 • LONG_BREAKOUT • AAPL           │
│ ✅ Checklist | ⏳ Sizing | ○ Heat | ○ Entry  │
├──────────────────────────────────────────────┤
│ Position Sizing Tab:                         │
│                                              │
│ Ticker: AAPL    [from session]               │
│ Banner: ✅ GREEN [from Checklist]            │
│                                              │
│ Entry: [180.00]                              │
│ ATR:   [1.50]                                │
│ K:     [2.0]                                 │
│                                              │
│ [Calculate]                                  │
└──────────────────────────────────────────────┘

User clicks Calculate →

┌──────────────────────────────────────────────┐
│ Results:                                     │
│ - Stop Distance: $3.00                       │
│ - Initial Stop: $177.00                      │
│ - Shares: 25                                 │
│ - Risk: $75                                  │
│                                              │
│ [Next: Heat Check →]                         │
└──────────────────────────────────────────────┘

Session updated:
- sizing_completed = 1
- sizing_shares = 25
- sizing_risk_dollars = 75
- Progress: ✅ Checklist | ✅ Sizing (done!)
```

### Step 4: Check Heat

```
┌──────────────────────────────────────────────┐
│ Session #47 • LONG_BREAKOUT • AAPL           │
│ ✅ Checklist | ✅ Sizing | ⏳ Heat | ○ Entry │
├──────────────────────────────────────────────┤
│ Heat Check Tab:                              │
│                                              │
│ Ticker: AAPL           [from session]        │
│ Risk: $75              [from Sizing]         │
│ Bucket: Tech/Comm                            │
│                                              │
│ [Check Heat]                                 │
└──────────────────────────────────────────────┘

User clicks Check Heat →

┌──────────────────────────────────────────────┐
│ Portfolio Heat:                              │
│ - Current: $2,100                            │
│ - New: $2,175 (with this trade)              │
│ - Cap: $4,000                                │
│ - Status: ✅ OK (54% utilized)               │
│                                              │
│ Bucket Heat (Tech/Comm):                     │
│ - Current: $1,400                            │
│ - New: $1,475                                │
│ - Cap: $1,500                                │
│ - Status: ⚠️ Near Cap (98% utilized)         │
│                                              │
│ [Next: Trade Entry →]                        │
└──────────────────────────────────────────────┘

Session updated:
- heat_completed = 1
- heat_status = 'OK'
- Progress: ✅ Checklist | ✅ Sizing | ✅ Heat (done!)
```

### Step 5: Final Gate Check & Decision

```
┌──────────────────────────────────────────────┐
│ Session #47 • LONG_BREAKOUT • AAPL           │
│ ✅ Checklist | ✅ Sizing | ✅ Heat | ⏳ Entry │
├──────────────────────────────────────────────┤
│ Trade Entry Tab:                             │
│                                              │
│ Session Summary:                             │
│ - Strategy: Long Breakout                    │
│ - Ticker: AAPL                               │
│ - Entry: $180, Stop: $177                    │
│ - Shares: 25, Risk: $75                      │
│ - Bucket: Tech/Comm                          │
│ - Banner: ✅ GREEN                           │
│                                              │
│ 5-Gate Check:                                │
│ ✅ Gate 1: Banner GREEN                      │
│ ✅ Gate 2: 2-min cooloff passed              │
│ ✅ Gate 3: Ticker not on cooldown            │
│ ✅ Gate 4: Heat caps OK                      │
│ ✅ Gate 5: Sizing complete                   │
│                                              │
│ All Gates PASS - Ready to Trade              │
│                                              │
│ [Save GO ✅]  [Save NO-GO ❌]                │
└──────────────────────────────────────────────┘

User clicks Save GO →

Session updated:
- status = 'COMPLETED'
- entry_decision = 'GO'
- entry_decision_id = 123 (link to decisions table)
- completed_at = now()
- Progress: ✅ Checklist | ✅ Sizing | ✅ Heat | ✅ Entry (DONE!)

Session now READ-ONLY (cannot edit)
```

---

## Resume Session Feature

```
User has 3 active sessions:

[Resume Session ▼]
  ├─ #47 (AAPL - Long)     ✅ ✅ ✅ ⏳  [2 min ago]
  ├─ #32 (TSLA - Short)    ✅ ⏳ ○ ○   [2 hours ago]
  └─ #18 (NVDA - Long)     ✅ ✅ ○ ○   [yesterday]

User clicks #32 →
- Loads Session #32 (TSLA, Short Breakout)
- Navigates to Sizing tab (current_step)
- All Checklist data still there (GREEN banner)
- Can continue where left off
```

---

## Session History

```
┌────────────────────────────────────────────────────┐
│ Trade Session History                              │
├────────────────────────────────────────────────────┤
│ Session  Ticker  Strategy      Status    Decision  │
├────────────────────────────────────────────────────┤
│ #47      AAPL    Long Breakout COMPLETED ✅ GO     │
│ #32      TSLA    Short         DRAFT     -         │
│ #28      NVDA    Long Breakout COMPLETED ❌ NO-GO  │
│ #18      MSFT    Custom        DRAFT     -         │
│ #12      XLE     Long Breakout COMPLETED ✅ GO     │
└────────────────────────────────────────────────────┘

Click on any session to view details (read-only)
Click "Clone" to create new draft with same setup
```

---

## Database Schema (Simplified)

```sql
CREATE TABLE trade_sessions (
    id INTEGER PRIMARY KEY,
    session_num INTEGER UNIQUE,           -- 1-99
    ticker TEXT,                          -- AAPL
    strategy TEXT,                        -- LONG_BREAKOUT
    status TEXT,                          -- DRAFT, COMPLETED

    -- Checklist gate
    checklist_completed INTEGER,          -- 0 or 1
    checklist_banner TEXT,                -- GREEN, YELLOW, RED

    -- Sizing gate
    sizing_completed INTEGER,             -- 0 or 1
    sizing_shares INTEGER,                -- 25
    sizing_risk_dollars REAL,             -- 75.00

    -- Heat gate
    heat_completed INTEGER,               -- 0 or 1
    heat_status TEXT,                     -- OK, WARN, REJECT

    -- Entry gate
    entry_completed INTEGER,              -- 0 or 1
    entry_decision TEXT,                  -- GO, NO-GO
    entry_decision_id INTEGER,            -- FK to decisions table

    -- Audit trail
    created_at DATETIME,
    updated_at DATETIME,
    completed_at DATETIME
);
```

---

## Progress Indicators Explained

```
✅  Green checkmark  = Gate completed
⏳  Hourglass        = Gate in progress (current step)
○   Hollow circle   = Gate pending (not started)
❌  Red X            = Gate failed (for banner states)
```

---

## Workflow Enforcement

**Without sessions (current):**
```
User on Checklist → RED banner
User clicks "Position Sizing" tab
→ Position Sizing loads with no context
→ User can calculate position anyway (bypass RED banner!)
❌ NO DISCIPLINE ENFORCEMENT
```

**With sessions (proposed):**
```
User on Checklist → RED banner
Session #47 marked: checklist_completed = 0, banner = RED
User clicks "Position Sizing" tab
→ Tab checks: if !session.checklistCompleted { showError() }
→ Position Sizing form DISABLED
→ Error message: "Session #47 has RED banner - resolve Checklist first"
✅ DISCIPLINE ENFORCED
```

---

## Why This Matters

### Alignment with Anti-Impulsivity Philosophy

| Principle                | How Sessions Support It                                      |
|--------------------------|--------------------------------------------------------------|
| Trade the tide           | Strategy selected FIRST, then find setups that match         |
| Friction where it matters| Must create session before evaluating (intentional friction) |
| Nudge for better trades  | Progress bar shows incomplete gates, nudges completion       |
| Immediate feedback       | Session bar shows status in real-time (✅ ⏳ ○)              |
| Journal while deciding   | Full session history = automatic journal of evaluation       |

### What Sessions Prevent

❌ Impulsive trading: "I'll just quickly size this without checking the banner"
❌ Lost context: "Wait, was this AAPL or TSLA?"
❌ Skipping gates: "I'll check heat later" (no, you check it NOW)
❌ No audit trail: "What was my reasoning 2 weeks ago?"

### What Sessions Enable

✅ Sequential workflow: Checklist → Sizing → Heat → Entry (cannot skip)
✅ Cohesive analysis: All tabs work on same trade
✅ Full history: Every decision logged with all gates
✅ Resumable sessions: Start today, finish tomorrow
✅ Strategy context: Know if evaluating Long vs Short breakout

---

## Implementation Timeline

```
Day 1 (4 hours):
  Phase 1: Database & Backend
    ├─ Create trade_sessions table
    ├─ Write CRUD methods (Create, Get, Update)
    ├─ Write tests
    └─ Verify all tests pass ✅

  Phase 2: UI Scaffolding (start)
    ├─ Add currentSession to AppState
    └─ Create session bar widget

Day 2 (4 hours):
  Phase 2: UI Scaffolding (finish)
    ├─ Create "New Trade" dialog
    ├─ Create "Resume Session" dropdown
    └─ Integrate into main UI

  Phase 3: Tab Integration (start)
    ├─ Modify Checklist tab
    └─ Modify Position Sizing tab

Day 3 (3-4 hours):
  Phase 3: Tab Integration (finish)
    ├─ Modify Heat Check tab
    ├─ Modify Trade Entry tab
    └─ Add navigation guards

  Phase 4: Polish & Testing
    ├─ Session history view
    ├─ Clone session feature
    ├─ E2E testing (5 scenarios)
    └─ Edge case testing

  Phase 5: Documentation
    └─ Update user guide

  Phase 6: Release
    └─ Build, test, deploy

Total: 11-15 hours (3 half-days)
```

---

## Key Design Decisions

| Decision                  | Proposed Choice       | Rationale                                  |
|---------------------------|-----------------------|--------------------------------------------|
| Session number format     | Random 1-99           | Human-memorable, fits in small UI space    |
| Session persistence       | 30-day auto-archive   | Cleanup stale drafts, keep audit trail    |
| Multi-session workflow    | Allow multiple DRAFT  | Flexibility, compare setups side-by-side   |
| Session locking           | Read-only after COMPLETED | Immutable audit trail (clone to re-analyze) |

---

## What You Need to Do

1. **Read this document** (5 minutes) ✅ You're here!

2. **Review detailed docs** (30-60 minutes)
   - [TRADE_SESSION_ARCHITECTURE.md](TRADE_SESSION_ARCHITECTURE.md) - Full architecture
   - [SESSION_UI_MOCKUPS.md](SESSION_UI_MOCKUPS.md) - Visual mockups
   - [IMPLEMENTATION_CHECKLIST.md](IMPLEMENTATION_CHECKLIST.md) - Step-by-step guide

3. **Make design decisions**
   - Session numbers: Random 1-99 or sequential?
   - Persistence: 30-day auto-archive or keep forever?
   - Multi-session: Allow multiple or one at a time?
   - Locking: Read-only after COMPLETED or allow edits?

4. **Approve or request changes**
   - "Looks good, let's implement" → We start Phase 1
   - "I have questions about X" → We discuss and refine
   - "I want to change Y" → We update planning docs

5. **Implementation begins**
   - Follow [IMPLEMENTATION_CHECKLIST.md](IMPLEMENTATION_CHECKLIST.md)
   - 11-15 hours total, 3 phases over 3 days
   - Testing at each step
   - Full deployment when done

---

## Questions?

- **"Is this too complex?"**
  - No. The complexity is in the planning, not the implementation. The concept is simple: one session = one trade being evaluated. The implementation is straightforward CRUD operations.

- **"Will this slow down my workflow?"**
  - No. Creating a session takes < 100ms. The session bar is always visible. Resuming a session loads in < 200ms. It's actually FASTER than re-entering ticker on every tab.

- **"Can I still use the old workflow?"**
  - Not recommended. The old workflow (direct tab access) will show "Create Session First" prompt. Sessions are the new standard.

- **"What if I don't like it?"**
  - Rollback plan included. We can revert the database migration and code changes. But let's get it right in planning first.

- **"When can we start?"**
  - Right now. Just say "Let's implement Phase 1" and we'll begin with the database.

---

## Bottom Line

**Before:** Disconnected tabs, no cohesion, no context, easy to bypass gates.

**After:** Unified sessions, progress tracking, full audit trail, discipline enforced.

**Cost:** 11-15 hours implementation, no breaking changes.

**Benefit:** Fundamental UX improvement that aligns perfectly with anti-impulsivity philosophy.

**Recommendation:** Approve and implement. This is the right solution.

---

**Next Step:** Your decision. Ready to proceed?

---

**Document Version:** 1.0
**Author:** Claude Code Planning Agent
**Read Time:** 5 minutes
