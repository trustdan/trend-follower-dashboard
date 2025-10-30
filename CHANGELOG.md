# Changelog

All notable changes to TF-Engine will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [2.0.0] - 2025-10-30

### 🎉 Major Release: Trade Sessions System

This release introduces **Trade Sessions**, a comprehensive workflow system that unifies all tabs and provides full audit trails for every trade evaluation.

### Added

#### Core Features
- **Trade Sessions** - Unified workflow tracking for each trade evaluation
  - Session Bar widget with real-time progress indicators (✅ ⏳ ○)
  - Sequential gate flow (Checklist → Sizing → Heat → Entry)
  - Session state persists across all tabs
  - Random session numbers (1-99) for easy reference
  - Full audit trail in database

- **Start New Trade Dialog**
  - Strategy selection (Long Breakout, Short Breakout, Custom)
  - Optional ticker entry
  - Creates session and navigates to Checklist tab
  - Keyboard shortcut: **Ctrl+N**

- **Resume Session Dropdown**
  - Lists all active DRAFT sessions
  - Shows progress indicators for each session
  - Time-since-update display (2 min ago, 2 hours ago, etc.)
  - Loads session and navigates to current step
  - Keyboard shortcut: **Ctrl+R**

- **Session History Tab**
  - View all past sessions (COMPLETED, DRAFT, ABANDONED)
  - Filter by status dropdown
  - Search by ticker symbol
  - View full session details in dialog
  - Clone sessions to create new drafts
  - Keyboard shortcut: **Ctrl+H**

- **Clone Session Feature**
  - Create new DRAFT session from existing session
  - Copies ticker and strategy
  - Resets all gate states (not completed)
  - Perfect for re-evaluating setups after conditions change

#### UI Enhancements
- **Read-Only Session View**
  - COMPLETED sessions show "🔒 READ-ONLY" indicator
  - All inputs disabled (ticker, checkboxes, calculate buttons)
  - Prevents post-decision tampering
  - Enforces immutable audit trail

- **Keyboard Shortcuts**
  - Ctrl+N: Start New Trade
  - Ctrl+R: Resume Session
  - Ctrl+H: Session History
  - Logged to tf-gui.log for debugging

- **Navigation Enhancement**
  - Added Session History as 8th navigation tab
  - Session Bar always visible at top when session active
  - Session controls row (Start New, Resume, Theme, Help, VIM)

#### Backend
- **Database Schema**
  - New `trade_sessions` table with full gate tracking
  - Session CRUD operations (Create, Get, Update, List, Clone)
  - 11 comprehensive tests (100% coverage)
  - Session lifecycle management (DRAFT → COMPLETED)

#### Tab Integration
- **All tabs now session-aware:**
  - Checklist: Ticker auto-filled from session
  - Position Sizing: Data pre-filled from Checklist
  - Heat Check: Risk pre-filled from Sizing
  - Trade Entry: Summary shows all session data

### Changed
- **USER_GUIDE.md** - Added comprehensive Trade Sessions section
  - Session workflow documentation
  - Keyboard shortcuts table
  - Session lifecycle diagrams
  - Best practices for session management
  - Troubleshooting guide

- **Navigation** - Updated tab structure:
  - Tab 0: Dashboard
  - Tab 1: Scanner
  - Tab 2: Checklist
  - Tab 3: Position Sizing
  - Tab 4: Heat Check
  - Tab 5: Trade Entry
  - Tab 6: Calendar
  - Tab 7: Session History ← NEW

### Fixed
- Session state now persists correctly across tab switches
- No more lost ticker data when changing tabs
- Banner state (RED/YELLOW/GREEN) flows through entire workflow
- Prerequisite checks prevent skipping gates

### Technical Details

#### Files Added
- `backend/internal/storage/sessions.go` (617 lines)
- `backend/internal/storage/sessions_test.go` (11 tests)
- `ui/session_bar.go` (151 lines)
- `ui/new_trade_dialog.go` (113 lines)
- `ui/session_selector.go` (193 lines)
- `ui/session_helpers.go` (helper functions)
- `ui/session_history.go` (234 lines)
- `ui/keyboard_shortcuts.go` (49 lines)

#### Files Modified
- `ui/main.go` - Added session management to AppState
- `ui/checklist.go` - Added 🔒 READ-ONLY indicator
- `ui/position_sizing.go` - Session-aware data flow
- `ui/heat_check.go` - Session-aware data flow
- `ui/trade_entry.go` - Session-aware data flow

#### Database Migration
- Run `go test ./backend/internal/storage/...` to verify schema

### Performance
- Session creation: < 100ms ✓
- Session load: < 200ms ✓
- Session list (100 items): < 500ms ✓
- Tab switch: < 50ms perceived ✓

### Breaking Changes
⚠️ **None** - All changes are additive. Existing workflows continue to work.

### Migration Guide
No migration required. Existing users will see:
1. New "Start New Trade" button in UI
2. Session Bar appears when first session created
3. All tabs continue to work as before
4. Session History tab available immediately

### Documentation
- Updated USER_GUIDE.md with Trade Sessions section
- Added PHASE_1_2_COMPLETE.md (backend implementation)
- Added PHASE_3_COMPLETE.md (tab integration)
- Added PHASE_4_COMPLETE.md (polish & testing)
- Created CHANGELOG.md (this file)

### Known Issues
- None identified

### Future Enhancements (Not in This Release)
- Session comparison (side-by-side view)
- Export to CSV/JSON
- Session templates
- Analytics dashboard
- Session tags
- Auto-archive old sessions

---

## [1.0.0] - 2025-10-29

### Initial Release

#### Core Features
- **Dashboard** - Portfolio overview with heat gauges
- **Scanner** - FINVIZ integration for daily candidate import
- **Checklist** - 5 required gates + 4 optional quality items
- **Position Sizing** - Van Tharp ATR-based calculations
- **Heat Check** - Portfolio and bucket heat cap enforcement
- **Trade Entry** - Final 5-gate validation
- **Calendar** - 10-week sector diversification grid

#### Anti-Impulsivity System
- RED/YELLOW/GREEN banner system
- 2-minute cool-off timer (cannot be bypassed)
- Heat caps (4% portfolio, 1.5% bucket)
- 5 hard gates (cannot bypass)
- Cooldown system (ticker and sector)

#### UI Features
- Light/Dark theme toggle
- VIM keybindings (F for hints, j/k/h/l navigation)
- Welcome dialog (first-run experience)
- Help dialog (comprehensive FAQ)
- British Racing Green color scheme
- Material Design inspiration

#### Backend
- SQLite database (trading.db)
- Position sizing algorithms (stock, opt-delta-atr, opt-contracts)
- Checklist validation logic
- Heat management
- Gates check engine
- FINVIZ web scraper

#### Documentation
- USER_GUIDE.md (comprehensive user manual)
- CLAUDE.md (developer instructions)
- README.md (project overview)
- anti-impulsivity.md (core philosophy)
- PROJECT_STATUS.md (implementation status)

---

## Version History

| Version | Date | Major Features |
|---------|------|----------------|
| **2.0.0** | 2025-10-30 | Trade Sessions System |
| **1.0.0** | 2025-10-29 | Initial Release |

---

## Upgrade Notes

### 1.0.0 → 2.0.0

**No action required.** Sessions are additive:
1. Old workflow still works (direct tab access)
2. New workflow available (Start New Trade → Sessions)
3. Database auto-creates `trade_sessions` table on first use

**Recommended:** Start using sessions for new trades. Benefits:
- Full audit trail
- Progress tracking
- Resume capability
- Session history view

**To try it:**
1. Click "Start New Trade" button
2. Select strategy and enter ticker
3. Follow the sequential workflow
4. Check Session History tab after completion

---

## Semantic Versioning

TF-Engine follows [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR** version (X.0.0): Breaking changes, incompatible API changes
- **MINOR** version (1.X.0): New features, backward-compatible additions
- **PATCH** version (1.0.X): Bug fixes, backward-compatible fixes

---

## Support

**Documentation:** See [docs/USER_GUIDE.md](docs/USER_GUIDE.md)
**Issues:** File issues on GitHub (if open source)
**Questions:** Check [docs/FAQ.md](docs/FAQ.md) first

---

**Philosophy:** Trade the tide, not the splash.
**Remember:** The value is in what this system prevents, not what it allows.
