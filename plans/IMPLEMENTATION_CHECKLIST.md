# Trade Session Implementation Checklist

**Purpose:** Step-by-step guide for implementing trade sessions
**Related:**
- [TRADE_SESSION_ARCHITECTURE.md](TRADE_SESSION_ARCHITECTURE.md) - Architecture details
- [SESSION_UI_MOCKUPS.md](SESSION_UI_MOCKUPS.md) - Visual mockups

---

## Pre-Implementation

- [ ] **Review planning documents**
  - [ ] Read `TRADE_SESSION_ARCHITECTURE.md` fully
  - [ ] Review `SESSION_UI_MOCKUPS.md` for UI clarity
  - [ ] Understand the "why" (UX cohesion, audit trail, discipline)

- [ ] **Gather stakeholder feedback**
  - [ ] Session number formatting (sequential ids, zero-padding in UI)
  - [ ] Strategy types (Long/Short/Custom sufficient?)
  - [ ] Session persistence duration (30 days auto-archive?)
  - [ ] Multi-session workflow (allow multiple drafts?)

- [ ] **Set up test environment**
  - [ ] Backup current `trading.db` database
  - [ ] Create test database with sample data
  - [ ] Verify backend tests pass: `cd backend && go test ./...`

---

## Phase 1: Database & Backend (Est: 2-3 hours)

### 1.1 Database Schema

- [ ] **Create migration file**
  - [ ] Create `backend/internal/storage/migrations/001_add_trade_sessions.sql`
  - [ ] Add `trade_sessions` table (see architecture doc for schema)
  - [ ] Add indexes (session_num, ticker, status, created_at)
  - [ ] Test migration on fresh database

- [ ] **Update schema.go**
  - [ ] Add `trade_sessions` table to `schema` constant
  - [ ] Test: `go run cmd/tf-engine/main.go init` should create table
  - [ ] Verify with: `sqlite3 trading.db ".schema trade_sessions"`

### 1.2 Backend Storage Layer

- [ ] **Create sessions.go storage file**
  - [ ] `backend/internal/storage/sessions.go`
  - [ ] Implement `CreateSession(ticker, strategy string) (*TradeSession, error)`
    - [ ] Session number comes from autoincrement id (no collision logic required)
    - [ ] Persist FINVIZ provenance (`source`, `candidate_id`, `preset_id`, `preset_name`, `scan_date`)
    - [ ] INSERT into trade_sessions table
    - [ ] Return TradeSession struct
  - [ ] Implement `GetSession(sessionNum int) (*TradeSession, error)`
    - [ ] SELECT by session_num
    - [ ] Handle not found error
  - [ ] Implement `UpdateSessionChecklist(sessionNum int, banner string, completed bool) error`
    - [ ] UPDATE checklist fields
    - [ ] Set updated_at timestamp
  - [ ] Implement `UpdateSessionSizing(sessionNum int, method string, shares int, risk float64, completed bool) error`
  - [ ] Implement `UpdateSessionHeat(sessionNum int, status string, completed bool) error`
  - [ ] Implement `UpdateSessionEntry(sessionNum int, decision string, decisionID int) error`
    - [ ] Set status='COMPLETED'
    - [ ] Set completed_at timestamp
  - [ ] Implement `ListActiveSessions() ([]*TradeSession, error)`
    - [ ] SELECT WHERE status='DRAFT'
    - [ ] ORDER BY updated_at DESC
  - [ ] Implement `ListSessionHistory(limit int) ([]*TradeSession, error)`
    - [ ] SELECT all sessions
    - [ ] ORDER BY created_at DESC
    - [ ] LIMIT to N results

### 1.3 API Types

- [ ] **Add TradeSession struct**
  - [ ] `backend/internal/api/types.go`
  - [ ] Define `TradeSession` struct with all fields
  - [ ] Add JSON tags for serialization
  - [ ] Add strategy constants (LONG_BREAKOUT, SHORT_BREAKOUT, CUSTOM)

### 1.4 Backend Tests

- [ ] **Create sessions_test.go**
  - [ ] `backend/internal/storage/sessions_test.go`
  - [ ] Test CreateSession
    - [ ] Happy path: creates session with unique session_num
    - [ ] Collision handling: generates new num if collision
    - [ ] Error: all 99 slots full
  - [ ] Test GetSession
    - [ ] Happy path: retrieves session by num
    - [ ] Error: session not found
  - [ ] Test UpdateSessionChecklist
    - [ ] Updates checklist fields correctly
    - [ ] Sets timestamp
  - [ ] Test UpdateSessionSizing
  - [ ] Test UpdateSessionHeat
  - [ ] Test UpdateSessionEntry
    - [ ] Sets status to COMPLETED
    - [ ] Sets completed_at timestamp
  - [ ] Test ListActiveSessions
    - [ ] Returns only DRAFT sessions
    - [ ] Sorted by updated_at DESC
  - [ ] Test ListSessionHistory
  - [ ] **Run tests:** `go test ./internal/storage/sessions_test.go -v`
  - [ ] **All tests must pass** before proceeding

---

## Phase 2: UI Scaffolding (Est: 2-3 hours)

### 2.1 AppState Extension

- [ ] **Update AppState in main.go**
  - [ ] Add `currentSession *api.TradeSession` field
  - [ ] Add `sessionChangeCallbacks []func(*api.TradeSession)` for observers
  - [ ] Add method `SetCurrentSession(session *api.TradeSession)`
    - [ ] Updates currentSession
    - [ ] Calls all registered callbacks
    - [ ] Triggers UI refresh

### 2.2 Session Bar Widget

- [ ] **Create session_bar.go**
  - [ ] `ui/session_bar.go`
  - [ ] Create `sessionBarWidget` struct
    - [ ] Fields: sessionLabel, progressLabel, background
  - [ ] Implement `NewSessionBar(state *AppState) *sessionBarWidget`
  - [ ] Method `updateFromSession(session *api.TradeSession)`
    - [ ] If session nil: show "No Active Session"
    - [ ] If session: show "#47 • LONG_BREAKOUT • AAPL"
    - [ ] Update progress: "✅ Checklist | ⏳ Sizing | ○ Heat | ○ Entry"
  - [ ] Method `updateColors()`
    - [ ] Active session: blue border, light blue background
    - [ ] No session: gray border, gray background
  - [ ] Register callback with AppState to auto-update

### 2.3 New Trade Dialog

- [ ] **Create new_trade_dialog.go**
  - [ ] `ui/new_trade_dialog.go`
  - [ ] Function `showNewTradeDialog(state *AppState)`
  - [ ] UI elements:
    - [ ] Strategy radio buttons (Long Breakout, Short Breakout, Custom)
    - [ ] Ticker entry (optional, with placeholder)
    - [ ] Info text explaining strategy
    - [ ] Create button
    - [ ] Cancel button
  - [ ] On Create button:
    - [ ] Get selected strategy
    - [ ] Call `state.db.CreateSession(ticker, strategy)`
    - [ ] Call `state.SetCurrentSession(newSession)`
    - [ ] Navigate to Checklist tab
    - [ ] Close dialog
  - [ ] On Cancel: close dialog
  - [ ] Test: clicking "Start New Trade" shows dialog

### 2.4 Resume Session Dropdown

- [ ] **Create session_selector.go**
  - [ ] `ui/session_selector.go`
  - [ ] Function `createSessionSelector(state *AppState) *widget.Button`
  - [ ] Returns button with dropdown popup
  - [ ] On click:
    - [ ] Call `state.db.ListActiveSessions()`
    - [ ] Show menu with sessions:
      - [ ] Format: "#47 (AAPL - Long) | ✅ Check ⏳ Size | 2h ago"
    - [ ] On session select:
      - [ ] Call `state.db.GetSession(sessionNum)`
      - [ ] Call `state.SetCurrentSession(session)`
      - [ ] Navigate to current_step tab
  - [ ] Test: dropdown shows active DRAFT sessions

### 2.5 Top Bar Integration

- [ ] **Update main.go buildMainUI**
  - [ ] Add session bar at very top (before theme/help buttons)
  - [ ] Add "Start New Trade" button
  - [ ] Add "Resume Session" dropdown
  - [ ] Add "Session History" button (placeholder for now)
  - [ ] Layout:
    ```
    [Session Bar: full width]
    [Start New | Resume ▼ | History]  [Theme | Help | VIM | Welcome]
    [Tab navigation]
    [Tab content]
    ```
  - [ ] Test: session bar visible, buttons work

---

## Phase 3: Tab Integration (Est: 3-4 hours)

### 3.1 Checklist Tab Integration

- [ ] **Modify checklist.go**
  - [ ] Add session check at top of `buildChecklistScreen`:
    ```go
    if state.currentSession == nil {
        return showNoSessionPrompt(state, "Checklist")
    }
    ```
  - [ ] Auto-fill ticker from session: `tickerEntry.SetText(state.currentSession.Ticker)`
  - [ ] Disable ticker entry if session locked (COMPLETED)
  - [ ] On "Evaluate" button click:
    - [ ] Call existing `domain.EvaluateChecklist()`
    - [ ] Call `state.db.UpdateSessionChecklist(sessionNum, banner, completed)`
    - [ ] Call `state.SetCurrentSession(updatedSession)` to trigger refresh
  - [ ] Show session-specific info in dialogs:
    - [ ] "From Preset" dialog mentions session's strategy
  - [ ] Add "Next: Position Sizing" button:
    - [ ] Only enabled if checklist_completed=1
    - [ ] Click navigates to Position Sizing tab
  - [ ] Test:
    - [ ] Without session: shows prompt
    - [ ] With session: auto-fills ticker
    - [ ] Evaluation updates session in database
    - [ ] Next button appears when complete

### 3.2 Position Sizing Tab Integration

- [ ] **Modify position_sizing.go**
  - [ ] Add session check (show prompt if no session)
  - [ ] Add prerequisite check:
    ```go
    if !state.currentSession.ChecklistCompleted {
        return showPrerequisiteError(state, "Checklist", "Position Sizing")
    }
    ```
  - [ ] Auto-fill ticker from session
  - [ ] On "Calculate" button click:
    - [ ] Call existing `domain.CalculateSizeStock()`
    - [ ] Call `state.db.UpdateSessionSizing(sessionNum, method, shares, risk, true)`
    - [ ] Call `state.SetCurrentSession(updatedSession)`
  - [ ] Add "Next: Heat Check" button (enabled after calculation)
  - [ ] Test:
    - [ ] RED banner in Checklist → Position Sizing disabled with error
    - [ ] GREEN banner → Position Sizing enabled
    - [ ] Calculation updates session

### 3.3 Heat Check Tab Integration

- [ ] **Modify heat_check.go**
  - [ ] Add session check
  - [ ] Add prerequisite check (sizing_completed must be true)
  - [ ] Auto-load risk_dollars from session
  - [ ] Auto-load ticker from session
  - [ ] On "Check Heat" button click:
    - [ ] Call existing `domain.CheckHeat()`
    - [ ] Call `state.db.UpdateSessionHeat(sessionNum, heatStatus, true)`
    - [ ] Call `state.SetCurrentSession(updatedSession)`
  - [ ] Add "Next: Trade Entry" button (enabled after check)
  - [ ] Test:
    - [ ] Without sizing → Heat Check disabled
    - [ ] With sizing → auto-loads risk amount
    - [ ] Check updates session

### 3.4 Trade Entry Tab Integration

- [ ] **Modify trade_entry.go**
  - [ ] Add session check
  - [ ] Add prerequisite check (heat_completed must be true)
  - [ ] Display full session summary:
    - [ ] Strategy, ticker, entry, stop, shares, risk, bucket, banner
  - [ ] Display 5-gate status:
    - [ ] Gate 1: Banner GREEN (from session)
    - [ ] Gate 2: 2-min cooloff (call existing logic)
    - [ ] Gate 3: Ticker cooldown (call existing logic)
    - [ ] Gate 4: Heat caps (from Heat Check results)
    - [ ] Gate 5: Sizing complete (from session)
  - [ ] On "Save GO" or "Save NO-GO" button click:
    - [ ] Call existing `domain.CheckGates()`
    - [ ] Call existing decision insertion (decisions table)
    - [ ] Call `state.db.UpdateSessionEntry(sessionNum, decision, decisionID)`
    - [ ] Call `state.SetCurrentSession(updatedSession)`
    - [ ] Show success dialog
    - [ ] Navigate to Dashboard (or show read-only view)
  - [ ] Test:
    - [ ] Without heat check → Trade Entry disabled
    - [ ] With all gates → displays full summary
    - [ ] Save GO → marks session COMPLETED
    - [ ] Session becomes read-only

### 3.5 Navigation Guards

- [ ] **Create navigation_guard.go**
  - [ ] `ui/navigation_guard.go`
  - [ ] Function `checkUnsavedChanges(state *AppState) bool`
    - [ ] Compares UI state to session state
    - [ ] Returns true if changes detected
  - [ ] Function `showUnsavedChangesDialog(state *AppState, onSave, onDiscard, onCancel func())`
    - [ ] Shows dialog: "Session #X has unsaved changes"
    - [ ] Buttons: Save & Continue, Discard, Cancel
  - [ ] Integrate with tab switching:
    - [ ] Before navigating, check for unsaved changes
    - [ ] If changes, show dialog
    - [ ] On Save: evaluate current tab, then navigate
    - [ ] On Discard: reload session from DB, then navigate
    - [ ] On Cancel: stay on current tab

---

## Phase 4: Polish & Testing (Est: 2-3 hours)

### 4.1 Session History View

- [ ] **Create session_history.go**
  - [ ] `ui/session_history.go`
  - [ ] Function `buildSessionHistoryScreen(state *AppState) fyne.CanvasObject`
  - [ ] UI:
    - [ ] Filter dropdown (All, COMPLETED, DRAFT, ABANDONED)
    - [ ] Search entry (by ticker)
    - [ ] Table: Session #, Ticker, Strategy, Status, Decision, Date
  - [ ] On row select: load session (read-only if COMPLETED)
  - [ ] Buttons:
    - [ ] View (opens session in read-only mode)
    - [ ] Clone (creates new draft with same ticker/strategy)
    - [ ] Export (CSV export - future feature)
  - [ ] Add to tab navigation
  - [ ] Test: shows all sessions, filtering works

### 4.2 Clone Session Feature

- [ ] **Add CloneSession method to storage**
  - [ ] `backend/internal/storage/sessions.go`
  - [ ] `CloneSession(sourceSessionNum int) (*TradeSession, error)`
    - [ ] Get source session
    - [ ] Create new session with same ticker/strategy
    - [ ] Copy relevant fields (not completion states)
    - [ ] Return new session
  - [ ] Test: cloned session is new DRAFT

### 4.3 Read-Only Session View

- [ ] **Create read_only indicator**
  - [ ] When session status is COMPLETED:
    - [ ] Show 🔒 icon in session bar
    - [ ] Disable all form inputs
    - [ ] Replace action buttons with "Clone Session" button
  - [ ] Test: COMPLETED session cannot be edited

### 4.4 Keyboard Shortcuts

- [ ] **Add keyboard shortcuts**
  - [ ] Ctrl+N → Start New Trade
  - [ ] Ctrl+R → Resume Session
  - [ ] Ctrl+H → Session History
  - [ ] Ctrl+1-5 → Navigate to tabs (if session active)
  - [ ] Register shortcuts in main.go
  - [ ] Test: shortcuts work

### 4.5 End-to-End Testing

- [ ] **E2E Test Scenario 1: Happy Path**
  - [ ] Start app (clean database)
  - [ ] Click "Start New Trade"
  - [ ] Select Long Breakout, ticker AAPL
  - [ ] Create session → verify session #1 created
  - [ ] Checklist tab: check all boxes, evaluate → verify GREEN
  - [ ] Click "Next: Position Sizing"
  - [ ] Enter entry=180, ATR=1.5, K=2, calculate → verify 25 shares
  - [ ] Click "Next: Heat Check"
  - [ ] Check heat → verify OK
  - [ ] Click "Next: Trade Entry"
  - [ ] Verify all 5 gates PASS
  - [ ] Click "Save GO" → verify decision logged
  - [ ] Verify session #1 status=COMPLETED
  - [ ] Verify session is read-only
  - [ ] **Database check:**
    ```sql
    SELECT * FROM trade_sessions WHERE session_num=1;
    SELECT * FROM decisions WHERE ticker='AAPL';
    ```

- [ ] **E2E Test Scenario 2: RED Banner**
  - [ ] Create new session (Long Breakout, TSLA)
  - [ ] Checklist: check only 3 boxes (leave 2 unchecked)
  - [ ] Evaluate → verify RED banner
  - [ ] Try to go to Position Sizing → verify error message
  - [ ] Fix checklist (check all boxes), evaluate → verify GREEN
  - [ ] Now Position Sizing works

- [ ] **E2E Test Scenario 3: Resume Session**
  - [ ] Create session, complete Checklist + Sizing
  - [ ] Close app (Ctrl+Q or close window)
  - [ ] Restart app
  - [ ] Click "Resume Session" → verify session listed
  - [ ] Select session → verify navigates to Heat Check (current_step)
  - [ ] Verify all prior data (ticker, shares, risk) still there

- [ ] **E2E Test Scenario 4: Multiple Sessions**
  - [ ] Create session #1 (AAPL, complete Checklist)
  - [ ] Create session #2 (TSLA, complete Checklist)
  - [ ] Switch between sessions → verify data doesn't mix
  - [ ] Resume Session dropdown → verify both listed
  - [ ] Complete session #1 (save GO)
  - [ ] Resume Session dropdown → verify #1 not listed (COMPLETED)
  - [ ] Session History → verify #1 appears as COMPLETED

- [ ] **E2E Test Scenario 5: Unsaved Changes**
  - [ ] Create session, check 2 boxes (don't evaluate)
  - [ ] Try to switch tabs → verify warning dialog
  - [ ] Click "Discard" → verify switches tabs, data not saved
  - [ ] Go back, check boxes again, evaluate → verify GREEN
  - [ ] Switch tabs → no warning (changes saved)

### 4.6 Edge Cases

- [ ] **Test session number collision**
  - [ ] Create 99 sessions (script or manual)
  - [ ] Try to create 100th → verify error "All slots full"

- [ ] **Test invalid session state**
  - [ ] Manually UPDATE database: SET checklist_completed=0 for COMPLETED session
  - [ ] Load session → verify app handles gracefully (shows warning?)

- [ ] **Test missing decision_id**
  - [ ] Complete session with GO decision
  - [ ] Manually DELETE from decisions table
  - [ ] Load session history → verify doesn't crash (shows "Decision missing")

- [ ] **Test rapid session switching**
  - [ ] Create 3 sessions
  - [ ] Rapidly switch between them
  - [ ] Verify no race conditions, no data corruption

### 4.7 Performance Testing

- [ ] **Measure session creation time**
  - [ ] Add logging: `log.Printf("Session created in %v", elapsed)`
  - [ ] Verify < 100ms

- [ ] **Measure session load time**
  - [ ] Load session with full history (all gates completed)
  - [ ] Verify < 200ms

- [ ] **Measure tab switch time**
  - [ ] Switch between tabs with session active
  - [ ] Verify no lag (< 50ms perceived)

---

## Phase 5: Documentation (Est: 1 hour)

### 5.1 User Guide

- [ ] **Update USER_GUIDE.md**
  - [ ] Add "Trade Sessions" section
  - [ ] Explain what sessions are (cohesive trade evaluation)
  - [ ] How to start new session
  - [ ] How to resume session
  - [ ] How to view history
  - [ ] Include screenshots (from mockups)

### 5.2 Developer Docs

- [ ] **Update FRESH_START_PLAN.md or create SESSIONS.md**
  - [ ] Document session lifecycle
  - [ ] Document database schema
  - [ ] Document API methods
  - [ ] Include examples

### 5.3 Changelog

- [ ] **Add to CHANGELOG.md** (or create if doesn't exist)
  - [ ] Version: v2.0.0 (major feature)
  - [ ] Added: Trade Session system
  - [ ] Added: Session bar with progress tracking
  - [ ] Added: Session history view
  - [ ] Changed: Tabs now require active session
  - [ ] Migration: Existing data compatible, sessions optional

---

## Phase 6: Release (Est: 30 min)

### 6.1 Pre-Release Checklist

- [ ] **All tests pass**
  - [ ] Backend: `go test ./... -v`
  - [ ] E2E: All scenarios pass
  - [ ] Edge cases: All handled

- [ ] **Code review** (self or peer)
  - [ ] No hardcoded values (use constants)
  - [ ] Error messages are clear and helpful
  - [ ] Logging added for debugging
  - [ ] No commented-out code

- [ ] **Database migration tested**
  - [ ] Fresh install: creates trade_sessions table
  - [ ] Upgrade: existing data intact, new table added
  - [ ] Rollback plan: backup script provided

- [ ] **Documentation complete**
  - [ ] User guide updated
  - [ ] Developer docs updated
  - [ ] README updated (if needed)

### 6.2 Build & Package

- [ ] **Build binaries**
  - [ ] Windows: `ui/build.bat` or `go build -o tf-gui.exe`
  - [ ] Linux: `GOOS=linux go build -o tf-gui`
  - [ ] macOS: `GOOS=darwin go build -o tf-gui-mac`

- [ ] **Test on clean environment**
  - [ ] Fresh VM or container
  - [ ] No trading.db (should initialize)
  - [ ] Create session → full workflow
  - [ ] Restart app → resume session

### 6.3 Release Notes

- [ ] **Write release notes**
  - [ ] Version: v2.0.0 - Trade Sessions
  - [ ] Summary: Major UX improvement for cohesive trade evaluation
  - [ ] Features:
    - [ ] Trade session system (track evaluation across tabs)
    - [ ] Session bar with progress indicators
    - [ ] Resume session functionality
    - [ ] Session history view
    - [ ] Keyboard shortcuts (Ctrl+N, Ctrl+R, Ctrl+H)
  - [ ] Breaking changes: None (backward compatible)
  - [ ] Migration: Automatic on first run
  - [ ] Known issues: (if any)

- [ ] **Tag release**
  - [ ] Git tag: `git tag -a v2.0.0 -m "Trade Sessions Release"`
  - [ ] Push tag: `git push origin v2.0.0`

### 6.4 Notify User

- [ ] **Send notification** (Windows Toast)
  ```powershell
  powershell.exe -Command "[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null; ... $xml.LoadXml('<toast>...</toast>'); ... Show($toast)"
  ```
  - [ ] Message: "TF-Engine v2.0 Ready - Trade Sessions Implemented!"

---

## Post-Release

### Monitor & Iterate

- [ ] **Gather user feedback**
  - [ ] Session concept clear?
  - [ ] UI intuitive?
  - [ ] Performance acceptable?
  - [ ] Bugs or confusion?

- [ ] **Track issues**
  - [ ] Create GitHub issues for bugs
  - [ ] Prioritize by severity
  - [ ] Fix critical bugs immediately

- [ ] **Plan enhancements**
  - [ ] Session comparison (compare two setups side-by-side)
  - [ ] Session templates (save/load common patterns)
  - [ ] Export to CSV/JSON
  - [ ] Analytics (success rate by strategy, etc.)

---

## Rollback Plan

If implementation fails or causes critical bugs:

1. **Database rollback**
   ```sql
   -- Remove trade_sessions table
   DROP TABLE IF EXISTS trade_sessions;
   ```

2. **Code rollback**
   ```bash
   git revert <commit-hash>
   git push
   ```

3. **Restore backup**
   ```bash
   cp trading.db.backup trading.db
   ```

4. **Rebuild without sessions**
   - [ ] Revert UI changes
   - [ ] Remove session checks from tabs
   - [ ] Restore direct tab access

---

## Success Criteria

Implementation is successful when:

1. ✅ **All E2E tests pass** (5 scenarios)
2. ✅ **User can create, evaluate, and complete a session** (full workflow)
3. ✅ **User can resume a session after restart** (persistence)
4. ✅ **Session data is accurate** (no data loss or corruption)
5. ✅ **Performance is acceptable** (< 100ms session creation, < 200ms load)
6. ✅ **Documentation is complete** (user guide, developer docs)
7. ✅ **Code is maintainable** (tests, comments, clear structure)

---

## Timeline Estimate

| Phase                  | Estimated Time | Dependencies          |
|------------------------|----------------|-----------------------|
| Phase 1: Database      | 2-3 hours      | None                  |
| Phase 2: UI Scaffold   | 2-3 hours      | Phase 1               |
| Phase 3: Tab Integration| 3-4 hours     | Phase 2               |
| Phase 4: Polish & Test | 2-3 hours      | Phase 3               |
| Phase 5: Documentation | 1 hour         | Phase 4               |
| Phase 6: Release       | 30 min         | Phase 5               |
| **TOTAL**              | **11-15 hours**| Sequential            |

**Recommended approach:**
- Day 1 (4 hours): Phase 1 + start Phase 2
- Day 2 (4 hours): Finish Phase 2 + Phase 3
- Day 3 (3-4 hours): Phase 4 + Phase 5 + Phase 6

---

**Document Version:** 1.0
**Author:** Claude Code Planning Agent
**Status:** Ready for Implementation
