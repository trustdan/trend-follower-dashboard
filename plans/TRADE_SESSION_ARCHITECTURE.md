# Trade Session Architecture Plan

**Created:** 2025-10-30
**Status:** Planning Phase
**Priority:** HIGH - Fundamental UX/Architecture Issue

---

## Problem Statement

### Current Architecture Flaw

Right now, the TF-Engine GUI has **disconnected tabs** that don't share state:

```
Checklist Tab       Position Sizing Tab     Trade Entry Tab
   ↓                      ↓                       ↓
Ticker: [____]      Ticker: [____]          Ticker: [____]
Strategy: ???       Strategy: ???           Strategy: ???
Banner: ???         Uses what data?         Uses what data?
```

**Issues:**
1. **No coherence** - User fills in AAPL on Checklist, then has to re-enter on Position Sizing
2. **No strategy context** - Which preset/strategy are we evaluating? Long breakout? Short breakout?
3. **No session concept** - The 5 gates refer to "a trade being evaluated" but there's no object representing that
4. **Impossible to track** - Can't tell if Position Sizing and Checklist are analyzing the same trade
5. **Data loss** - User enters ticker, switches tabs, data disappears
6. **No history** - Can't review "what was I thinking when I evaluated this?"

### What We Need

A **Trade Session** concept:

```
User clicks "New Trade"
   ↓
Selects Strategy (Long Breakout, Short Breakout, etc.)
   ↓
Session receives next sequential ID (Session #47)
   ↓
All tabs now work on Session #47 (Long Breakout, AAPL)
   ↓
Session persisted in database throughout workflow
   ↓
Can resume, review, or start fresh
```

**User experience:**
```
Top of window:
┌─────────────────────────────────────────────────────────┐
│ Trade Session: #47 (Long Breakout - AAPL)              │
│ Status: Checklist ✓ | Sizing ✓ | Heat ⏳ | Entry ○      │
└─────────────────────────────────────────────────────────┘

[Dashboard] [Checklist] [Position Sizing] [Heat Check] [Trade Entry] [Calendar]
```

All tabs see the same session. All data persisted. User always knows what they're working on.

---

## Design Principles

From CLAUDE.md:

> **1. Discipline Over Flexibility**
> This is not a flexible trading platform. It's a discipline enforcement system.

**Application to Trade Sessions:**
- Sessions MUST complete all 5 gates sequentially
- Cannot skip ahead (enforced by UI)
- Cannot delete sessions with GO decisions (audit trail)
- Cannot modify completed sessions (immutability)
- Friction is intentional - making you think through each step

> **2. Fail Loudly**
> Silent failures are unacceptable.

**Application:**
- Session state always visible (which gates passed/failed)
- Clear error if user tries to jump to Trade Entry before completing Checklist
- Warnings if switching sessions with unsaved work
- Validation on every field

> **3. Error Messages Must Teach**

**Application:**
- "Cannot evaluate Trade Entry - Checklist not completed for Session #47"
- "Session #47 has RED banner - resolve checklist issues before sizing"
- "Session #32 (TSLA) already has a GO decision dated 2025-10-29"

---

## Proposed Architecture

### 1. Database Schema

New table: `trade_sessions`

```sql
CREATE TABLE IF NOT EXISTS trade_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_num INTEGER GENERATED ALWAYS AS (id) STORED UNIQUE,
    ticker TEXT NOT NULL,
    strategy TEXT NOT NULL,
    source TEXT NOT NULL DEFAULT 'MANUAL' CHECK (source IN ('MANUAL', 'PRESET', 'CUSTOM')),
    candidate_id INTEGER,
    preset_id INTEGER,
    preset_name TEXT,
    scan_date TEXT,
    status TEXT NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'EVALUATING', 'COMPLETED', 'ABANDONED')),
    current_step TEXT NOT NULL DEFAULT 'CHECKLIST' CHECK (current_step IN ('CHECKLIST', 'SIZING', 'HEAT', 'ENTRY')),
    checklist_completed INTEGER NOT NULL DEFAULT 0 CHECK (checklist_completed IN (0,1)),
    checklist_banner TEXT,
    checklist_missing_count INTEGER DEFAULT 0,
    checklist_quality_score INTEGER DEFAULT 0,
    checklist_completed_at DATETIME,
    sizing_completed INTEGER NOT NULL DEFAULT 0 CHECK (sizing_completed IN (0,1)),
    sizing_method TEXT,
    sizing_shares INTEGER,
    sizing_risk_dollars REAL,
    sizing_completed_at DATETIME,
    heat_completed INTEGER NOT NULL DEFAULT 0 CHECK (heat_completed IN (0,1)),
    heat_status TEXT,
    heat_completed_at DATETIME,
    entry_completed INTEGER NOT NULL DEFAULT 0 CHECK (entry_completed IN (0,1)),
    entry_decision TEXT,
    entry_decision_id INTEGER,
    entry_completed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    FOREIGN KEY (candidate_id) REFERENCES candidates(id),
    FOREIGN KEY (preset_id) REFERENCES presets(id),
    FOREIGN KEY (entry_decision_id) REFERENCES decisions(id)
);

CREATE UNIQUE INDEX idx_sessions_session_num ON trade_sessions(session_num);
CREATE INDEX idx_sessions_status ON trade_sessions(status);
CREATE INDEX idx_sessions_candidate ON trade_sessions(candidate_id);
```

### 2. Session Lifecycle

```
1. NEW
   User: "New Trade" button
   System: Next sequential session number assigned (session_num = id)
   System: Prompt for Strategy selection and record preset provenance (candidate_id/preset_id)
   System: INSERT trade_sessions (status='DRAFT', current_step='CHECKLIST')

2. CHECKLIST
   User: Fill checkboxes, evaluate
   System: UPDATE trade_sessions SET checklist_banner='GREEN', checklist_completed=1
   System: INSERT checklist_evaluations (ticker, banner, ...) -- existing table
   System: Enable "Next: Position Sizing" button

3. SIZING
   User: Enter entry, ATR, calculate
   System: Call domain.CalculateSizeStock()
   System: UPDATE trade_sessions SET sizing_shares=X, sizing_risk_dollars=Y, sizing_completed=1
   System: Enable "Next: Heat Check" button

4. HEAT CHECK
   User: Verify portfolio/bucket heat
   System: Call domain.CheckHeat()
   System: UPDATE trade_sessions SET heat_status='OK', heat_completed=1
   System: Enable "Next: Trade Entry" button

5. TRADE ENTRY
   User: Final 5-gate review, click "Save GO" or "Save NO-GO"
   System: Call domain.CheckGates()
   System: INSERT decisions (action='GO' or 'NO-GO')
   System: UPDATE trade_sessions SET entry_decision='GO', entry_decision_id=X, status='COMPLETED'

6. COMPLETED
   Session locked for editing (read-only)
   Can view history
   Can clone to new session
```

### 3. Strategy Types

From FINVIZ presets and Donchian logic:

```go
const (
    StrategyLongBreakout  = "LONG_BREAKOUT"   // TF_BREAKOUT_LONG preset
    StrategyShortBreakout = "SHORT_BREAKOUT"  // TF_BREAKOUT_SHORT preset
    StrategyCustom        = "CUSTOM"          // Manual entry
)
```

**Strategy determines:**
- Checklist preset (from_preset check auto-passes if ticker in candidates table for that strategy)
- Info dialog content (what the strategy means)
- Exit rules reminder (long uses upper Donchian, short uses lower)

### 4. UI Changes

#### Top Bar (New)

```
┌─────────────────────────────────────────────────────────────────┐
│ TF-Engine Dashboard                                             │
│                                                                 │
│ ╔═══════════════════════════════════════════════════════════╗  │
│ ║ Trade Session: #47 • LONG_BREAKOUT • AAPL                 ║  │
│ ║ Progress: ✅ Checklist | ✅ Sizing | ⏳ Heat | ○ Entry      ║  │
│ ╚═══════════════════════════════════════════════════════════╝  │
│                                                                 │
│ [New Trade] [Resume Session ▼] [Session History]              │
└─────────────────────────────────────────────────────────────────┘
```

#### "New Trade" Dialog

```
┌─────────────────────────────────────┐
│ Start New Trade Session             │
├─────────────────────────────────────┤
│                                     │
│ Strategy:                           │
│   (*) Long Breakout                 │
│   ( ) Short Breakout                │
│   ( ) Custom                        │
│                                     │
│ Ticker (optional): [____]           │
│                                     │
│ [ Create Session ]  [ Cancel ]      │
└─────────────────────────────────────┘

On Create:
  → Generate session_num (e.g., 47)
  → INSERT trade_sessions
  → Navigate to Checklist tab
  → Show session bar at top
```

#### Modified Tab Headers

**Before (Current):**
```
Checklist Tab:
  Ticker: [____]
  (no context about what strategy)
```

**After (Proposed):**
```
Checklist Tab:
  Session: #47 (Long Breakout - AAPL)     [locked, shown at top]
  Ticker: AAPL                             [auto-filled from session, editable]

  [Info: Long Breakout means price > 55-high with 2×N stop...]
```

#### Session Dropdown (Top Bar)

```
Active Session: #47 (AAPL - Long) ▼
  ├─ #47 (AAPL - Long)       ← current
  ├─ #32 (TSLA - Short)      [COMPLETED 2025-10-29]
  ├─ #18 (NVDA - Long)       [DRAFT]
  └─ [View All Sessions...]
```

Switching sessions:
- Loads all data for that session into tabs
- Shows warning if current session has unsaved changes
- Updates top bar

---

## Data Flow Examples

### Scenario 1: Happy Path (All Gates Pass)

```
1. User clicks "New Trade"
   → Dialog: Select strategy = Long Breakout, Ticker = AAPL
   → Creates Session #47
   → DB: INSERT trade_sessions (
         ticker='AAPL',
         strategy='LONG_BREAKOUT',
         source='PRESET',
         candidate_id=981,
         preset_id=3,
         preset_name='TF_BREAKOUT_LONG',
         scan_date='2025-10-29',
         status='DRAFT'
       )

2. User goes to Checklist tab
   → UI auto-fills: Ticker = AAPL (from session)
   → User checks all boxes (from_preset=1, trend=1, liquidity=1, timeframe=1, earnings=1)
   → User clicks "Evaluate"
   → Backend: domain.EvaluateChecklist() → banner='GREEN'
   → DB: UPDATE trade_sessions SET checklist_banner='GREEN', checklist_completed=1
   → UI: Shows GREEN banner, enables "Next: Position Sizing" button

3. User clicks "Next: Position Sizing"
   → Navigate to Position Sizing tab
   → UI auto-fills: Ticker = AAPL, loads settings (equity, risk_pct)
   → User enters: Entry=180, ATR=1.5, K=2.0
   → User clicks "Calculate"
   → Backend: domain.CalculateSizeStock() → shares=25, risk=$75
   → DB: UPDATE trade_sessions SET sizing_shares=25, sizing_risk_dollars=75, sizing_completed=1
   → UI: Shows calculation, enables "Next: Heat Check" button

4. User clicks "Next: Heat Check"
   → Navigate to Heat Check tab
   → UI auto-loads: Risk=$75 (from session), Bucket="Tech/Comm" (from candidates or manual)
   → User clicks "Check Heat"
   → Backend: domain.CheckHeat() → heat_status='OK'
   → DB: UPDATE trade_sessions SET heat_status='OK', heat_completed=1
   → UI: Shows heat OK, enables "Next: Trade Entry" button

5. User clicks "Next: Trade Entry"
   → Navigate to Trade Entry tab
   → UI shows: Session summary (all gates, green checkmarks)
   → User clicks "Save GO"
   → Backend: domain.CheckGates() → all pass
   → DB: INSERT decisions (action='GO', ticker='AAPL', risk_dollars=75, banner='GREEN')
   → DB: UPDATE trade_sessions SET entry_decision='GO', entry_decision_id=123, status='COMPLETED'
   → UI: Success message, locks session

Result: Complete audit trail in `trade_sessions` table showing when each gate passed.
```

### Scenario 2: User Switches Tabs Mid-Session

```
1. User on Checklist tab, Session #47 (AAPL)
   → Checks 3 boxes (not evaluated yet)
   → Clicks "Position Sizing" tab

2. System warns:
   "Session #47 has unsaved changes. Continue?"
   [ Save & Continue ] [ Discard ] [ Cancel ]

3. User clicks "Save & Continue"
   → Auto-evaluate checklist (incomplete → RED banner)
   → DB: UPDATE trade_sessions SET checklist_banner='RED'
   → Navigate to Position Sizing tab
   → Position Sizing tab shows: "⚠️ Session #47 has RED banner - resolve checklist first"
   → Position Sizing form is disabled (grayed out)
```

### Scenario 3: User Resumes Old Session

```
1. User clicks "Resume Session" dropdown
   → Shows list of DRAFT sessions:
      #32 (TSLA - Short)    Updated 2 hours ago
      #18 (NVDA - Long)     Updated yesterday

2. User selects #32
   → DB: SELECT * FROM trade_sessions WHERE session_num=32
   → Load state:
      - checklist_completed=1, banner='GREEN'
      - sizing_completed=0
      - current_step='SIZING'
   → UI: Show session bar "#32 (TSLA - Short)"
   → Navigate to Position Sizing tab (current_step)
   → Position Sizing auto-fills ticker=TSLA

3. User completes sizing → updates session
```

---

## Implementation Plan

### Phase 1: Database & Backend (2-3 hours)

1. **Create migration script** `backend/internal/storage/migrations/001_trade_sessions.sql`
   - Add `trade_sessions` table
   - Add session_num generator function
   - Test with sample data

2. **Add session methods** to `backend/internal/storage/sessions.go`
   ```go
   func (db *DB) CreateSession(ticker, strategy string) (*TradeSession, error)
   func (db *DB) GetSession(sessionNum int) (*TradeSession, error)
   func (db *DB) UpdateSessionChecklist(sessionNum int, banner string, completed bool) error
   func (db *DB) UpdateSessionSizing(sessionNum int, shares int, risk float64, completed bool) error
   func (db *DB) UpdateSessionHeat(sessionNum int, status string, completed bool) error
   func (db *DB) UpdateSessionEntry(sessionNum int, decision string, decisionID int) error
   func (db *DB) ListActiveSessions() ([]*TradeSession, error)
   func (db *DB) ListSessionHistory(limit int) ([]*TradeSession, error)
   ```

3. **Add session struct** to `backend/internal/api/types.go`
   ```go
   type TradeSession struct {
       ID                  int       `json:"id"`
       SessionNum          int       `json:"session_num"`
       Ticker              string    `json:"ticker"`
       Strategy            string    `json:"strategy"`
       Status              string    `json:"status"`
       CurrentStep         string    `json:"current_step"`
       ChecklistCompleted  bool      `json:"checklist_completed"`
       ChecklistBanner     string    `json:"checklist_banner,omitempty"`
       SizingCompleted     bool      `json:"sizing_completed"`
       SizingShares        int       `json:"sizing_shares,omitempty"`
       SizingRiskDollars   float64   `json:"sizing_risk_dollars,omitempty"`
       HeatCompleted       bool      `json:"heat_completed"`
       HeatStatus          string    `json:"heat_status,omitempty"`
       EntryCompleted      bool      `json:"entry_completed"`
       EntryDecision       string    `json:"entry_decision,omitempty"`
       CreatedAt           time.Time `json:"created_at"`
       UpdatedAt           time.Time `json:"updated_at"`
   }
   ```

4. **Write tests** for session CRUD operations

### Phase 2: UI Scaffolding (2-3 hours)

1. **Add session state** to `ui/main.go` AppState:
   ```go
   type AppState struct {
       db              *storage.DB
       window          fyne.Window
       isDarkMode      bool
       myApp           fyne.App
       currentSession  *api.TradeSession  // NEW
   }
   ```

2. **Create session bar widget** `ui/session_bar.go`
   - Shows: "Session #47 • LONG_BREAKOUT • AAPL"
   - Shows progress: ✅ Checklist | ✅ Sizing | ⏳ Heat | ○ Entry
   - Updates on session change

3. **Create "New Trade" dialog** `ui/new_trade_dialog.go`
   - Strategy radio buttons (Long Breakout, Short Breakout, Custom)
   - Optional ticker entry
   - Creates session, updates AppState.currentSession

4. **Create "Resume Session" dropdown** `ui/session_selector.go`
   - Lists active DRAFT sessions
   - On select: loads session into AppState.currentSession
   - Triggers UI refresh

### Phase 3: Tab Integration (3-4 hours)

1. **Modify Checklist tab** `ui/checklist.go`
   - Check if session exists: `if state.currentSession == nil { show "Start New Trade" prompt }`
   - Auto-fill ticker from session
   - On "Evaluate": call `UpdateSessionChecklist()`
   - Show session-specific info dialogs

2. **Modify Position Sizing tab** `ui/position_sizing.go`
   - Check if checklist completed: `if !state.currentSession.ChecklistCompleted { disable form }`
   - Auto-fill ticker from session
   - On "Calculate": call `UpdateSessionSizing()`

3. **Modify Heat Check tab** `ui/heat_check.go`
   - Check if sizing completed
   - Auto-load risk_dollars from session
   - On "Check": call `UpdateSessionHeat()`

4. **Modify Trade Entry tab** `ui/trade_entry.go`
   - Check if heat completed
   - Display full session summary (all gates)
   - On "Save GO/NO-GO": call `UpdateSessionEntry()`

5. **Add navigation warnings**
   - Detect unsaved changes before tab switch
   - Show dialog: "Session #X has unsaved changes"

### Phase 4: Polish & Testing (2-3 hours)

1. **Session history view**
   - New tab or dialog showing all COMPLETED sessions
   - Filterable by ticker, strategy, date
   - Read-only view of session details

2. **Clone session feature**
   - "Copy Session #47" → creates new session with same ticker/strategy
   - Useful for re-evaluating a setup

3. **Keyboard shortcuts**
   - Ctrl+N: New Trade
   - Ctrl+R: Resume Session
   - Ctrl+S: Save current step

4. **End-to-end testing**
   - Create session → Checklist → Sizing → Heat → Entry → GO decision
   - Verify database state at each step
   - Test resume mid-session
   - Test switching sessions

5. **Edge case handling**
   - What if user creates 100 sessions? (session_num wraps, show error)
   - What if database write fails? (rollback UI state)
   - What if user force-quits mid-session? (recover on restart)

---

## Migration Strategy

### For Existing Users

Current system has no sessions. Some users may have:
- Saved decisions in `decisions` table
- Open positions in `positions` table
- Checklist evaluations in `checklist_evaluations` table

**Migration approach:**

1. **New table only** - `trade_sessions` table is new, doesn't affect existing tables
2. **Backward compatible** - Old workflow (direct tab access) still works, but shows "Create Session First" prompt
3. **Optional import** - Script to retroactively create sessions from `decisions` table entries
   ```sql
   -- For each decision, create a completed session
   INSERT INTO trade_sessions (session_num, ticker, strategy, status, entry_decision, entry_decision_id)
   SELECT
       ROW_NUMBER() OVER (ORDER BY created_at),
       ticker,
       'CUSTOM',  -- Unknown strategy for historical data
       'COMPLETED',
       action,
       id
   FROM decisions
   WHERE created_at > '2025-10-01';  -- Recent decisions only
   ```

4. **Gradual rollout**
   - Release 1: Add sessions table, optional usage
   - Release 2: Make sessions required, prompt users to create
   - Release 3: Full integration, remove old direct-tab workflow

---

## Open Questions

### 1. Session Number Generation

**Options:**
- **Sequential ID** (1, 2, 3, ...)
  - Pros: No collisions, trivially generated from primary key
  - Cons: May need zero-padding in UI if you want fixed width

- **Ticker + Timestamp** (AAPL-1030-1400)
  - Pros: Self-documenting
  - Cons: Long, harder to reference verbally

**Recommendation:** Use sequential ids (`session_num` generated from `id`). Format in UI (`Session #047`) if you want compact presentation.

### 2. Session Persistence Duration

**How long do DRAFT sessions live?**

Options:
- **Forever** (until user deletes)
- **Auto-delete after 7 days** (cleanup stale drafts)
- **Auto-delete when COMPLETED** (keep only active)

**Recommendation:** Keep DRAFT sessions for 30 days, auto-archive after that. Keep COMPLETED sessions forever (audit trail).

### 3. Multi-Session Workflows

**Can user work on multiple sessions simultaneously?**

- **Yes** (proposed): Dropdown to switch between sessions
  - Pros: Flexibility, can compare setups side-by-side
  - Cons: More complex UI, risk of confusion

- **No**: One session at a time, must complete/abandon before starting new
  - Pros: Simpler, enforces focus
  - Cons: Less flexible

**Recommendation:** Allow multiple DRAFT sessions, but encourage completion. Show warning if user has >5 active sessions.

### 4. Session Locking

**After GO/NO-GO decision, is session read-only?**

- **Yes** (proposed): Immutable audit trail
  - Cannot change ticker, strategy, gate results
  - Can view history

- **No**: Allow edits
  - Pros: Flexibility if user made mistake
  - Cons: Violates audit trail principle

**Recommendation:** Read-only after COMPLETED. If user needs to fix, they must "Clone Session" to create new draft.

---

## Success Metrics

After implementation, measure:

1. **User Comprehension**
   - Do users understand the session concept?
   - Feedback from user testing: "I can now track what I'm analyzing"

2. **Workflow Completion**
   - % of sessions that reach Trade Entry tab
   - Dropout rate at each gate (Checklist → Sizing → Heat → Entry)

3. **Data Integrity**
   - Zero orphaned data (sizing calculations without ticker context)
   - All GO decisions have complete session history

4. **Performance**
   - Session creation: < 100ms
   - Session load: < 200ms
   - No lag when switching sessions

---

## Next Steps

1. **Review this plan** with stakeholder (you!)
   - Discuss strategy types (Long/Short/Custom sufficient?)
   - Confirm session number approach (sequential auto-increment)
   - Decide on session persistence duration

2. **Prototype session bar** in UI
   - Mock the top bar widget
   - Test visual clarity

3. **Implement Phase 1** (database)
   - Create table, write methods, test

4. **Iterate on UI** (Phases 2-3)
   - Get feedback early on "New Trade" dialog
   - Test session switching UX

5. **Full integration** (Phase 4)
   - End-to-end testing
   - Document in user guide

---

## Alignment with Anti-Impulsivity Principles

From `docs/anti-impulsivity.md`:

✅ **Trade the tide, not the splash**
   → Sessions force you to define strategy FIRST, then find setups that match

✅ **Friction where it matters**
   → Must create session before evaluating (intentional friction)

✅ **Nudge for better trades**
   → Progress bar shows incomplete gates, nudging toward completion

✅ **Immediate feedback**
   → Session bar shows status in real-time (✅ ⏳ ○)

✅ **Journal while deciding**
   → Full session history = automatic journal of evaluation process

✅ **Calendar awareness**
   → Sessions can link to Calendar tab for sector diversification check

---

## Conclusion

**Trade Sessions solve the fundamental UX problem:**
> "What trade am I analyzing right now, and what data belongs to it?"

By introducing sessions, we:
- ✅ Create coherence across tabs
- ✅ Maintain audit trail of decision process
- ✅ Enable workflow resumption
- ✅ Prevent data loss
- ✅ Enforce discipline (can't skip gates)
- ✅ Build toward more advanced features (comparison, backtesting, etc.)

This is a **foundational change** that makes the system usable for real trading decisions, not just isolated calculations.

**Recommendation: Proceed with implementation.**

---

**Document Version:** 1.0
**Author:** Claude Code Planning Agent
**Next Review:** After Phase 1 completion
