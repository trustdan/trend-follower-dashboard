# TF-Engine: Fresh Start with Custom GUI

**TF = Trend Following** - Systematic Donchian breakout system following Ed Seykota/Turtle Trader principles

**Created:** October 29, 2025
**Status:** 🚀 Backend Ready | Documentation Complete (Step 27)
**Backend:** ✅ 100% Functional (tf-engine in Go)
**Frontend:** ✅ Embedded UI (Svelte) - Production Ready
**Installer:** 🚧 Phase 1 Complete (Step 26)

---

## 📚 User Documentation

**New to TF-Engine? Start here:**

### Getting Started
- **[Quick Start](docs/QUICK_START.md)** - Get started in 10 minutes ⚡
- **[Installation Guide](docs/INSTALLATION_GUIDE.md)** - Windows setup instructions
- **[TradingView Setup](docs/TRADINGVIEW_SETUP.md)** - Install Ed-Seykota Pine Script

### Complete Guide
- **[User Guide](docs/USER_GUIDE.md)** - Comprehensive documentation (primary reference)
  - Daily trading workflow
  - Understanding the 5 gates
  - Banner states (RED/YELLOW/GREEN)
  - Position sizing (Van Tharp method)
  - Heat management
  - Tips & best practices

### Help & Support
- **[FAQ](docs/FAQ.md)** - Frequently asked questions
- **[Troubleshooting](docs/TROUBLESHOOTING.md)** - Common issues and solutions

---

## What is This?

This is a **fresh start** for the **TF-Engine (Trend Following Engine)** trading platform. We're building a **custom GUI application** to replace the problematic Excel/VBA frontend, while keeping the proven, tested Go backend.

---

## Why the Fresh Start?

The Excel/VBA frontend had fundamental integration issues:
- Parse function signature mismatches
- Type name and property name errors
- OLE control compatibility problems
- Difficult testing and deployment
- Poor developer experience

**The Go backend works perfectly.** We're just replacing the interface.

See [PROJECT_HISTORY.md](PROJECT_HISTORY.md) for full details.

---

## What We Have

### ✅ Reusable Components

**Backend (Go)** - `backend/`
- All position sizing algorithms (stock, opt-delta-atr, opt-contracts)
- Checklist evaluation with GREEN/YELLOW/RED banners
- Heat check calculations (portfolio and sector caps)
- 5 hard gates enforcement
- SQLite database with full CRUD operations
- FINVIZ screener import
- Comprehensive test coverage

**Scripts** - `scripts/`
- `import-candidates-auto.bat` - Automated FINVIZ import
- `import-candidates.bat` - Manual FINVIZ import
- `build-windows.sh` - Cross-compile for Windows
- `test-finviz-scrape.sh` - Test scraper

**Documentation** - `docs/`
- `CLAUDE.md` - Guidance for Claude Code (future AI sessions) ⭐
- `1._RULES.md` - Operating rules for this project (Linux-first, no Git in Linux) ⭐
- `LLM-update.md` - Session-by-session tracking log (always current)
- `PROGRESS.md` - Narrative progress and decisions
- `anti-impulsivity.md` - Core design philosophy
- `PROJECT_STATUS.md` - Current project status (M24 complete)
- `M24_UI_IMPLEMENTATION_PLAN.md` - UI plans
- `UI_QUICK_REFERENCE.md` - UI reference
- Plus project docs, milestones, dev docs, JSON schemas

**Test Data** - `test-data/`
- Sample databases
- JSON examples
- Test scenarios
- SQL fixtures

---

## What We're Building

A **native GUI application** using Go + Fyne (or similar) that:

1. **Directly calls tf-engine backend** (no HTTP, no CLI spawning)
2. **Enforces anti-impulsivity design** (large banners, 2-min cooloff, 5 gates)
3. **Deploys as single binary** (no installation, no dependencies)
4. **Works cross-platform** (Windows, Linux, macOS)
5. **Has modern UX** (fast, responsive, clear feedback)

### 6 Main Screens

1. **Dashboard** - Overview of positions, candidates, settings, cooldowns
2. **Checklist** - 5 required gates + optional quality items + RED/YELLOW/GREEN banner
3. **Position Sizing** - Calculate shares/contracts using ATR-based risk
4. **Heat Check** - Verify portfolio and sector heat within caps
5. **Trade Entry** - Final 5-gate check before saving GO/NO-GO decision
6. **Calendar** - Rolling 10-week sector × week grid for diversification

See [FRESH_START_PLAN.md](FRESH_START_PLAN.md) for complete details.

---

## Anti-Impulsivity Design

Based on [docs/anti-impulsivity.md](docs/anti-impulsivity.md):

### Core Principles

- **Trade the tide, not the splash** - Donchian breakouts with mechanical exits
- **Friction where it matters** - Hard gates for signal, risk, liquidity, exit, behavior
- **Nudge for better trades** - Optional quality items affect score, not permission
- **Immediate feedback** - Large 3-state banner updates live
- **Journal while deciding** - One-click logging of full trade plan
- **Calendar awareness** - 10-week sector view for diversification

### The 5 Hard Gates

1. **Signal:** 55-bar breakout (long > 55-high / short < 55-low)
2. **Risk/Size:** Per-unit risk = % of equity using 2×N stop; pyramids every 0.5×N to max units
3. **Options:** 60–90 DTE, roll/close ~21 DTE, liquidity required
4. **Exits:** Exit by 10-bar opposite Donchian OR closer of 2×N
5. **Behavior:** 2-minute cool-off + no intraday overrides

**RED:** Any required gate fails → **DO NOT TRADE**
**YELLOW:** All required pass, quality score < threshold → **CAUTION**
**GREEN:** All required pass, quality meets threshold → **OK TO TRADE**

---

## Directory Structure

```
fresh-start/
├── README.md                  # This file
├── PROJECT_HISTORY.md         # What happened with Excel/VBA
├── FRESH_START_PLAN.md        # Detailed plan for GUI
├── ORIGINAL_README.md         # Original project README
├── .gitignore                 # Git ignore rules
│
├── backend/                   # tf-engine Go backend (WORKING)
│   ├── cmd/                   # CLI entry point
│   ├── internal/
│   │   ├── cli/               # Command handlers
│   │   ├── domain/            # Business logic
│   │   ├── storage/           # SQLite persistence
│   │   ├── scrape/            # FINVIZ scraper
│   │   ├── server/            # HTTP server (for Excel)
│   │   └── logx/              # Logging
│   ├── go.mod
│   └── go.sum
│
├── scripts/                   # Build and import scripts
│   ├── import-candidates-auto.bat
│   ├── import-candidates.bat
│   ├── build-windows.sh
│   └── test-finviz-scrape.sh
│
├── docs/                      # Documentation
│   ├── anti-impulsivity.md    # Core design philosophy ⭐
│   ├── PROJECT_STATUS.md
│   ├── M24_UI_IMPLEMENTATION_PLAN.md
│   ├── UI_QUICK_REFERENCE.md
│   ├── project/               # Project documentation
│   ├── milestones/            # Milestone docs
│   ├── dev/                   # Development docs
│   └── json-schemas/          # JSON schemas
│
├── art/                       # ASCII art and assets
│   └── tf-engine_exe-ASCII.txt
│
└── test-data/                 # Test databases and examples
    ├── json-examples/
    ├── test-contracts.db
    ├── phase4-test-data.sql
    └── ...
```

---

## Quick Start

### 1. Verify Backend Works

```bash
cd backend/
go build -o tf-engine cmd/tf-engine/main.go
./tf-engine init
./tf-engine settings --equity 100000 --risk-pct 0.75 --portfolio-cap 4.0
./tf-engine size --ticker AAPL --entry 180 --atr 1.5 --method stock --k 2
```

**Expected output:**
```json
{
  "success": true,
  "ticker": "AAPL",
  "entry_price": 180,
  "shares": 250,
  "risk_dollars": 750,
  "stop_price": 177,
  ...
}
```

If this works, backend is good to go! ✅

---

### 2. Read the Documentation

**Start here:**
1. [PROJECT_HISTORY.md](PROJECT_HISTORY.md) - Understand what we tried with Excel/VBA
2. [docs/anti-impulsivity.md](docs/anti-impulsivity.md) - Core design philosophy
3. [FRESH_START_PLAN.md](FRESH_START_PLAN.md) - Detailed GUI plan

---

### 3. Choose GUI Framework

**Recommendation: Fyne**

```bash
# Install Fyne
go get fyne.io/fyne/v2

# Hello World
cat > hello.go << 'EOF'
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
EOF

go run hello.go
```

If a window opens with "Hello TF-Engine!", you're ready to build! ✅

---

### 4. Start Building

**Phase 1: Dashboard (Week 1-2)**

Create `frontend/` directory:
```bash
mkdir -p frontend/ui
cd frontend/
go mod init github.com/youruser/tf-engine-gui
go get fyne.io/fyne/v2
```

Build first screen:
```go
// frontend/ui/dashboard.go
package ui

import (
    "fyne.io/fyne/v2/widget"
    "github.com/youruser/tf-engine/internal/storage"
)

func BuildDashboard(db *storage.DB) *widget.Label {
    // Get positions from database
    positions, _ := db.GetAllPositions()

    // Display count
    return widget.NewLabel(fmt.Sprintf("Active Positions: %d", len(positions)))
}
```

---

## Backend API Reference

### Position Sizing

```go
import "github.com/youruser/tf-engine/internal/domain"

// Calculate position size
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
// result.StopPrice = 177
```

### Checklist Evaluation

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

// result.Banner = "GREEN" | "YELLOW" | "RED"
// result.MissingCount = 0
// result.AllowSave = true
```

### Heat Check

```go
import "github.com/youruser/tf-engine/internal/domain"

result := domain.CheckHeat(
    db,
    ticker: "NVDA",
    riskAmount: 750.0,
    bucket: "Tech/Comm",
)

// result.CurrentPortfolioHeat = 2.25
// result.AfterPortfolioHeat = 3.0
// result.PortfolioCapExceeded = false
// result.BucketCapExceeded = false
```

### 5 Gates Check

```go
import "github.com/youruser/tf-engine/internal/domain"

result := domain.CheckGates(
    db,
    banner: "GREEN",
    ticker: "AAPL",
    riskDollars: 750.0,
    bucket: "Tech/Comm",
)

// result.Gate1BannerGreen = true
// result.Gate2CooloffElapsed = true
// result.Gate3NotOnCooldown = true
// result.Gate4HeatOK = true
// result.Gate5SizingDone = true
// result.AllGatesPass = true
```

---

## Testing Strategy

### Backend Testing (Already Done)

```bash
cd backend/
go test ./internal/domain/... -v
go test ./internal/storage/... -v
```

All tests pass ✅

### GUI Testing (To Be Added)

**Manual testing checklist:**
- [ ] Dashboard displays positions correctly
- [ ] Position sizing calculates accurately
- [ ] Checklist banner shows correct colors
- [ ] Heat check validates caps
- [ ] Trade entry enforces all 5 gates
- [ ] Calendar shows 10-week view

**Automated testing:**
- Unit tests for UI components
- Integration tests for backend calls
- End-to-end workflow tests

---

## Deployment

### Build for Windows

**From Linux/macOS (bash):**
```bash
cd frontend/
GOOS=windows GOARCH=amd64 go build -o tf-engine-gui.exe
```

**From Windows (PowerShell):**
```powershell
cd frontend/
# Build for Windows (native - no env vars needed)
go build -o tf-engine-gui.exe
```

### Build for Linux

**From Linux/macOS (bash):**
```bash
cd frontend/
GOOS=linux GOARCH=amd64 go build -o tf-engine-gui
```

**From Windows (PowerShell):**
```powershell
cd frontend/
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o tf-engine-gui
```

### Build for macOS

**From Linux/macOS (bash):**
```bash
cd frontend/
GOOS=darwin GOARCH=amd64 go build -o tf-engine-gui-mac
```

**From Windows (PowerShell):**
```powershell
cd frontend/
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o tf-engine-gui-mac
```

### Package with Fyne

```bash
# Install fyne CLI
go install fyne.io/fyne/v2/cmd/fyne@latest

# Package for Windows
fyne package -os windows -icon icon.png

# Package for Linux
fyne package -os linux -icon icon.png

# Package for macOS
fyne package -os darwin -icon icon.png
```

---

## Roadmap

### Phase 1: Foundation (Week 1-2)
- ✅ Set up Fyne project
- ✅ Integrate tf-engine backend
- ✅ Build Dashboard (read-only)

### Phase 2: Core Functionality (Week 3-4)
- ⬜ Build Position Sizing screen
- ⬜ Build Checklist screen with banner

### Phase 3: Heat & Gates (Week 5-6)
- ⬜ Build Heat Check screen
- ⬜ Build Trade Entry screen
- ⬜ Implement 5 gates enforcement

### Phase 4: Calendar & Polish (Week 7-8)
- ⬜ Build Calendar screen
- ⬜ Polish all screens
- ⬜ Add keyboard shortcuts

### Phase 5: Testing & Deployment (Week 9-10)
- ⬜ Integration testing
- ⬜ Package for all platforms
- ⬜ Create user documentation

---

## Contributing

This is a personal trading system, but contributions are welcome:

1. **Report bugs** - Open an issue on GitHub
2. **Suggest features** - Propose enhancements aligned with anti-impulsivity design
3. **Submit PRs** - Backend improvements, GUI components, documentation

**Anti-patterns to avoid:**
- Don't add subjective pattern checks to required gates
- Don't over-complicate the workflow
- Don't add features that encourage impulsivity

---

## FAQ

### Q: Why not keep Excel?

Excel/VBA had fundamental integration issues (see PROJECT_HISTORY.md). After 5+ major issues, it's time for a better approach.

### Q: Why not use a web UI?

Native GUI is faster, more responsive, and easier to deploy (single binary).

### Q: Can I still use the CLI?

Yes! The backend supports both CLI and programmatic usage.

```bash
./tf-engine size --ticker AAPL --entry 180 --atr 1.5 --method stock
```

### Q: What about the Excel workbook?

It's archived in the original project. The VBA code and documentation are preserved in `/home/kali/excel-trading-platform/release/TradingEngine-v3/`.

### Q: How long will this take?

**Aggressive timeline:** 8 weeks
**Conservative timeline:** 12 weeks

See [FRESH_START_PLAN.md](FRESH_START_PLAN.md) for detailed schedule.

### Q: What GUI framework are you using?

**Recommendation: Fyne** (pure Go, cross-platform, Material Design)

Alternatives: Gio (immediate mode), Wails (Go + web frontend)

### Q: Will it work on macOS/Linux?

Yes! Go and Fyne are cross-platform. Build once, deploy everywhere.

---

## Resources

### Learning Fyne

- **Official Docs:** https://fyne.io/
- **Tutorial:** https://developer.fyne.io/tutorial/
- **Widget Tour:** https://apps.fyne.io/
- **Examples:** https://github.com/fyne-io/examples

### Go Best Practices

- **Effective Go:** https://go.dev/doc/effective_go
- **Go by Example:** https://gobyexample.com/
- **Project Layout:** https://github.com/golang-standards/project-layout

### Anti-Impulsivity Trading

- **Turtle Traders:** https://en.wikipedia.org/wiki/Turtle_trading
- **Ed Seykota:** https://www.seykota.com/
- **ATR-Based Stops:** Technical analysis standard (Wilder's ATR)

---

## License

This is a personal trading system. Use at your own risk.

**Disclaimer:** Trading involves risk. This software is for educational purposes and personal use. Past performance does not guarantee future results.

---

## Support

For questions or issues:
1. Check [PROJECT_HISTORY.md](PROJECT_HISTORY.md) for context
2. Review [FRESH_START_PLAN.md](FRESH_START_PLAN.md) for implementation details
3. Read [docs/anti-impulsivity.md](docs/anti-impulsivity.md) for design philosophy
4. Open an issue on GitHub

---

## Status

**Backend:** ✅ 100% Functional
**Frontend:** 🚧 To be built
**Timeline:** 8-12 weeks
**Next Action:** Choose GUI framework and build Hello World

---

**Last Updated:** October 29, 2025
**Version:** v4.0 (Fresh Start)
**Status:** 📋 Planning complete, ready to build

**Let's build something better!** 🚀
