# Fresh Start Organization Summary

**Created:** October 29, 2025
**Location:** `/root/fresh-start/`
**Purpose:** Clean slate for custom GUI development

---

## What's Included

This directory contains all reusable components from the original Excel-based trading platform, organized for starting fresh with a custom GUI frontend.

---

## Directory Structure

```
/root/fresh-start/
├── README.md                      # Main entry point - START HERE
├── PROJECT_HISTORY.md             # What we tried with Excel/VBA and why we're pivoting
├── FRESH_START_PLAN.md            # Detailed plan for custom GUI (screens, architecture, timeline)
├── ORIGINAL_README.md             # Original project README for reference
├── .gitignore                     # Git ignore rules
│
├── backend/                       # ✅ 100% FUNCTIONAL Go backend
│   ├── cmd/tf-engine/            # CLI entry point
│   ├── internal/
│   │   ├── cli/                  # Command handlers (size, checklist, heat, etc.)
│   │   ├── domain/               # Business logic (sizing, gates, heat calculations)
│   │   ├── storage/              # SQLite database CRUD operations
│   │   ├── scrape/               # FINVIZ screener scraper
│   │   ├── server/               # HTTP server (for Excel integration)
│   │   └── logx/                 # Logging utilities
│   ├── go.mod                    # Go module definition
│   └── go.sum                    # Go dependency checksums
│
├── scripts/                       # Build and import automation
│   ├── import-candidates-auto.bat    # Windows: Auto-import from FINVIZ
│   ├── import-candidates.bat         # Windows: Manual import from FINVIZ
│   ├── build-windows.sh              # Linux: Cross-compile for Windows
│   └── test-finviz-scrape.sh         # Linux: Test FINVIZ scraper
│
├── docs/                          # Comprehensive documentation
│   ├── anti-impulsivity.md       # ⭐ CORE DESIGN PHILOSOPHY - Read this first!
│   ├── PROJECT_STATUS.md         # Current project status
│   ├── M24_UI_IMPLEMENTATION_PLAN.md  # UI implementation plans
│   ├── UI_QUICK_REFERENCE.md     # UI reference
│   │
│   ├── project/                  # Project planning docs
│   │   ├── PLAN.md               # Original project plan
│   │   └── WHY.md                # Why this project exists
│   │
│   ├── milestones/               # Milestone completion summaries (M17-M24)
│   │   ├── M17-M18_*.md          # Milestones 17-18
│   │   ├── M19_COMPLETION_SUMMARY.md
│   │   ├── M20_COMPLETION_SUMMARY.md
│   │   ├── M21_*.md              # Milestone 21 + Phase 4
│   │   ├── M22_*.md              # Milestone 22
│   │   ├── M23_COMPLETION_SUMMARY.md
│   │   ├── M24_*.md              # Milestone 24
│   │   └── TODO_ENABLE_SKIPPED_TESTS.md
│   │
│   ├── dev/                      # Development practices
│   │   ├── BDD_GUIDE.md          # Behavior-Driven Development guide
│   │   ├── CLAUDE.md             # Working with Claude Code
│   │   ├── CLAUDE_RULES.md       # Claude interaction rules
│   │   └── DEVELOPMENT_PHILOSOPHY.md
│   │
│   ├── features/                 # Gherkin feature files (BDD specs)
│   │   ├── bucket_cooldown.feature
│   │   ├── checklist.feature
│   │   ├── finviz-scraper.feature
│   │   ├── heat.feature
│   │   ├── http_api.feature
│   │   ├── import-candidates.feature
│   │   ├── impulse-brake.feature
│   │   ├── position_management.feature
│   │   ├── save-decision.feature
│   │   ├── settings.feature
│   │   ├── sizing_options_delta_atr.feature
│   │   ├── sizing_options_maxloss.feature
│   │   ├── sizing_stocks.feature
│   │   └── storage.feature
│   │
│   └── json-schemas/             # JSON API specifications
│       └── JSON_API_SPECIFICATION.md
│
├── art/                           # ASCII art and assets
│   └── tf-engine_exe-ASCII.txt
│
└── test-data/                     # Test fixtures and examples
    ├── json-examples/            # Sample JSON requests/responses
    │   ├── requests/
    │   ├── responses/
    │   └── errors/
    ├── test-contracts.db         # Sample SQLite database
    ├── phase4-test-data.sql      # SQL test data
    ├── phase4-test-scenarios.csv # CSV test scenarios
    └── capture-*.sh              # Scripts to capture JSON examples
```

---

## File Inventory

### Core Documentation (4 files)

| File | Size | Description |
|------|------|-------------|
| `README.md` | ~16KB | Main entry point, quick start guide |
| `PROJECT_HISTORY.md` | ~35KB | Excel/VBA issues, fixes attempted, lessons learned |
| `FRESH_START_PLAN.md` | ~40KB | GUI plan: architecture, screens, timeline, decisions |
| `ORIGINAL_README.md` | ~11KB | Original project README for reference |

**Start with:** `README.md` → `PROJECT_HISTORY.md` → `docs/anti-impulsivity.md` → `FRESH_START_PLAN.md`

---

### Backend Code (45+ Go files)

**Status:** ✅ 100% Functional, fully tested

**Key modules:**
- `cmd/tf-engine/main.go` - Entry point
- `internal/cli/*.go` - Command handlers (size, checklist, heat, gates, etc.)
- `internal/domain/*.go` - Business logic (sizing algorithms, heat calculations, gate checks)
- `internal/storage/*.go` - SQLite CRUD operations
- `internal/scrape/*.go` - FINVIZ screener scraper
- `internal/server/*.go` - HTTP server (for Excel integration)

**Test coverage:** 45+ test files (`*_test.go`)

---

### Scripts (4 files)

| Script | Platform | Purpose |
|--------|----------|---------|
| `import-candidates-auto.bat` | Windows | Auto-import candidates from FINVIZ |
| `import-candidates.bat` | Windows | Manual import with prompts |
| `build-windows.sh` | Linux/WSL | Cross-compile Go binary for Windows |
| `test-finviz-scrape.sh` | Linux | Test FINVIZ scraper |

---

### Documentation (50+ files)

**Philosophy & Vision:**
- `docs/anti-impulsivity.md` - ⭐ **Core design principles** (READ THIS FIRST!)
- `docs/project/WHY.md` - Why this system exists
- `docs/project/PLAN.md` - Original plan

**Implementation Guides:**
- `docs/M24_UI_IMPLEMENTATION_PLAN.md` - UI design and implementation
- `docs/UI_QUICK_REFERENCE.md` - Quick UI reference
- `docs/HTTP_CLI_PARITY.md` - HTTP vs CLI interface parity

**Development Practices:**
- `docs/dev/BDD_GUIDE.md` - Behavior-Driven Development
- `docs/dev/DEVELOPMENT_PHILOSOPHY.md` - Development approach
- `docs/dev/CLAUDE.md` - Working with Claude Code

**Feature Specifications (Gherkin):**
- `docs/features/*.feature` - 14 feature files with BDD scenarios

**Milestones (M17-M24):**
- `docs/milestones/M*_*.md` - 16 milestone documents

**API Specs:**
- `docs/json-schemas/JSON_API_SPECIFICATION.md` - JSON API documentation

---

### Test Data

**JSON Examples:**
- `test-data/json-examples/requests/` - Sample requests
- `test-data/json-examples/responses/` - Sample responses
- `test-data/json-examples/errors/` - Error cases

**Databases:**
- `test-data/test-contracts.db` - Sample SQLite database

**SQL Fixtures:**
- `test-data/phase4-test-data.sql` - SQL test data
- `test-data/phase4-test-scenarios.csv` - CSV test scenarios

---

## What's NOT Included (Intentionally)

These were Excel/VBA-specific and are being replaced:

❌ Excel workbook (`TradingPlatform.xlsm`)
❌ VBA modules (`TFEngine.bas`, `TFHelpers.bas`, etc.)
❌ VBA fix scripts (`fix-vba-modules.bat`, `check-vba-version.vbs`, etc.)
❌ Excel-specific documentation (VBA_SIGNATURE_FIX_README.md, etc.)
❌ Windows release binaries (we'll build fresh)

**These are archived** in `/home/kali/excel-trading-platform/release/TradingEngine-v3/`

---

## Key Concepts

### The 5 Hard Gates

1. **Signal (SIG_REQ):** 55-bar breakout confirmed
2. **Risk/Size (RISK_REQ):** 2×N stop, 0.5×N adds, max units
3. **Options (OPT_REQ):** 60-90 DTE, liquidity OK
4. **Exits (EXIT_REQ):** 10-bar or 2×N exit plan
5. **Behavior (BEHAV_REQ):** 2-min cooloff, no overrides

**RED:** Any gate fails → **DO NOT TRADE**
**YELLOW:** All gates pass, quality score < threshold → **CAUTION**
**GREEN:** All gates pass, quality meets threshold → **OK TO TRADE**

### Backend Commands

The tf-engine backend supports these commands (all working):

- `init` - Initialize database
- `settings` - Set equity, risk %, caps
- `size` - Calculate position size (stock/opt-delta-atr/opt-contracts)
- `checklist` - Evaluate 6 quality criteria, get banner color
- `heat` - Check portfolio and sector heat
- `save-decision` - Log GO/NO-GO trade decision with all details
- `scrape` - Import candidates from FINVIZ
- `import` - Import candidates from CSV
- `positions` - Manage open positions
- `cooldown` - Manage cooldown list
- `server` - Start HTTP server (for Excel/external integrations)

All commands work via:
1. **CLI** - `./tf-engine <command> [flags]`
2. **HTTP** - `POST http://localhost:8080/<command>` with JSON body
3. **Direct Go calls** - Import packages and call functions (for GUI)

---

## Backend API Examples

### Position Sizing (Go)

```go
import "github.com/youruser/tf-engine/internal/domain"

result := domain.CalculateSize(
    ticker:   "AAPL",
    entry:    180.0,
    atr:      1.5,
    method:   "stock",
    k:        2.0,
    riskPct:  0.0075,
    equity:   100000.0,
)
// result.Shares = 250
// result.RiskDollars = 750
```

### Checklist Evaluation (Go)

```go
import "github.com/youruser/tf-engine/internal/domain"

result := domain.EvaluateChecklist(
    fromPreset: true,
    trendPass: true,
    liquidityOK: true,
    timeframeConfirm: true,
    earningsOK: true,
    journalOK: true,
)
// result.Banner = "GREEN"
// result.AllowSave = true
```

### Heat Check (Go)

```go
import (
    "github.com/youruser/tf-engine/internal/domain"
    "github.com/youruser/tf-engine/internal/storage"
)

db := storage.OpenDB("trading.db")
result := domain.CheckHeat(db, "NVDA", 750.0, "Tech/Comm")
// result.PortfolioCapExceeded = false
// result.BucketCapExceeded = false
```

---

## Quick Verification

### 1. Verify Backend Works

```bash
cd /root/fresh-start/backend/
go build -o tf-engine cmd/tf-engine/main.go
./tf-engine init
./tf-engine settings --equity 100000 --risk-pct 0.75
./tf-engine size --ticker AAPL --entry 180 --atr 1.5 --method stock
```

Expected output:
```json
{
  "success": true,
  "shares": 250,
  "risk_dollars": 750,
  "stop_price": 177
}
```

✅ **If this works, backend is ready!**

---

### 2. Run Tests

```bash
cd /root/fresh-start/backend/
go test ./internal/domain/... -v
go test ./internal/storage/... -v
```

Expected: All tests pass ✅

---

### 3. Test FINVIZ Scraper

```bash
cd /root/fresh-start/
bash scripts/test-finviz-scrape.sh
```

Expected: CSV file with candidate tickers ✅

---

## Next Steps

### 1. Read the Documentation

**Priority order:**
1. `README.md` - Overview and quick start
2. `PROJECT_HISTORY.md` - Context on Excel/VBA issues
3. `docs/anti-impulsivity.md` - **Core design philosophy**
4. `FRESH_START_PLAN.md` - Detailed GUI plan

### 2. Choose GUI Framework

**Recommendation: Fyne**

Install:
```bash
go get fyne.io/fyne/v2
```

Test:
```go
package main
import (
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/widget"
)
func main() {
    myApp := app.New()
    myWindow := myApp.NewWindow("TF-Engine")
    myWindow.SetContent(widget.NewLabel("Hello TF-Engine!"))
    myWindow.ShowAndRun()
}
```

### 3. Build Dashboard (Phase 1)

Create `frontend/` directory and build first screen:
```bash
mkdir -p /root/fresh-start/frontend/ui
cd /root/fresh-start/frontend/
go mod init github.com/youruser/tf-engine-gui
```

See `FRESH_START_PLAN.md` for detailed implementation plan.

---

## Timeline

**Aggressive:** 8 weeks
**Conservative:** 12 weeks

**Phase 1:** Foundation (Week 1-2) - Dashboard working
**Phase 2:** Core Functionality (Week 3-4) - Position sizing + Checklist
**Phase 3:** Heat & Gates (Week 5-6) - Heat check + Trade entry
**Phase 4:** Calendar & Polish (Week 7-8) - Calendar + UX improvements
**Phase 5:** Testing & Deployment (Week 9-10) - Package and ship

---

## Success Criteria

### Functional Requirements

- ✅ Position sizing calculates correctly (all 3 methods)
- ✅ Checklist evaluates with RED/YELLOW/GREEN banner
- ✅ Heat check validates portfolio and sector caps
- ✅ Trade entry enforces all 5 gates
- ✅ Dashboard displays positions, candidates, settings
- ✅ Calendar shows 10-week sector × week grid
- ✅ FINVIZ import populates candidates
- ✅ All data persists to SQLite

### Non-Functional Requirements

- ⚡ Fast UI updates (< 100ms)
- 🚀 Single binary deployment
- 🖥️  Cross-platform (Windows, Linux, macOS)
- 📦 Small binary size (< 50MB)
- 🎨 Clean, modern UI
- ⌨️  Keyboard navigation

---

## Support Resources

### Learning Materials

**Fyne (GUI):**
- Docs: https://fyne.io/
- Tutorial: https://developer.fyne.io/tutorial/
- Examples: https://github.com/fyne-io/examples

**Go:**
- Effective Go: https://go.dev/doc/effective_go
- Go by Example: https://gobyexample.com/

**Trading:**
- Turtle Traders: https://en.wikipedia.org/wiki/Turtle_trading
- Ed Seykota: https://www.seykota.com/

### Internal Documentation

- `docs/anti-impulsivity.md` - Design philosophy
- `docs/features/*.feature` - Behavior specs (Gherkin)
- `docs/json-schemas/JSON_API_SPECIFICATION.md` - API docs
- `docs/dev/BDD_GUIDE.md` - Development practices

---

## What We Learned

### From Excel/VBA Experience

1. **Don't use Excel for complex integrations** - Great for spreadsheets, bad for software
2. **Type safety matters** - VBA's weak typing caused constant errors
3. **Backend separation was wise** - tf-engine can be reused with any frontend
4. **Good documentation pays off** - Made pivoting possible
5. **Know when to pivot** - After 5+ issues, try a different approach

### Going Forward

1. **Start simple** - Get basic UI working first
2. **Iterate quickly** - Build one screen at a time
3. **Test continuously** - Don't wait until the end
4. **Focus on workflow** - Not spreadsheet mimicry
5. **Anti-impulsivity first** - Banner, gates, cooloff are core features

---

## Directory Statistics

```
Total files:     ~200+
Go source:       ~45 files
Test files:      ~45 files
Documentation:   ~50 files
Feature specs:   14 files
Scripts:         4 files
Test data:       ~50+ files
```

**Backend Status:** ✅ 100% Functional (45+ tests passing)
**Frontend Status:** 🚧 To be built (fresh start)
**Documentation:** ✅ Comprehensive (100+ KB)
**Architecture:** ✅ Proven, tested, reliable

---

## Final Checklist

Before starting GUI development:

- [ ] Read `README.md`
- [ ] Read `PROJECT_HISTORY.md`
- [ ] Read `docs/anti-impulsivity.md` ⭐
- [ ] Read `FRESH_START_PLAN.md`
- [ ] Verify backend works (`go test ./...`)
- [ ] Choose GUI framework (Fyne recommended)
- [ ] Set up frontend project structure
- [ ] Build Hello World GUI
- [ ] Integrate backend (call one function)
- [ ] Build Dashboard screen (Phase 1)

---

## Questions?

1. Check this file (`ORGANIZATION_SUMMARY.md`)
2. Read `README.md` for quick start
3. Read `PROJECT_HISTORY.md` for context
4. Read `FRESH_START_PLAN.md` for implementation details
5. Check relevant docs in `docs/`

---

**Status:** ✅ Organization complete
**Next Action:** Read documentation and choose GUI framework
**Timeline:** 8-12 weeks to fully functional GUI
**Confidence:** High - backend is rock solid

**Let's build something better!** 🚀

---

**Last Updated:** October 29, 2025
**Location:** `/root/fresh-start/`
**Purpose:** Clean slate for custom GUI development
