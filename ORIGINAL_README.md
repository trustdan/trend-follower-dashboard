# Trading Engine v3

**Excel-based trading platform enforcing disciplined trend-following through systematic constraints.**

This is not a flexible trading platform. It's a discipline enforcement system that makes impulsive trading impossible.

---

## 🚀 Quick Start

**New to the project?** Read in this order:
1. [docs/project/WHY.md](docs/project/WHY.md) - Why this exists (5 min) ⭐ **START HERE**
2. [docs/project/PLAN.md](docs/project/PLAN.md) - Architecture and roadmap (20 min)
3. [docs/dev/DEVELOPMENT_PHILOSOPHY.md](docs/dev/DEVELOPMENT_PHILOSOPHY.md) - How we build (10 min)

**Working with Claude Code?** Read [docs/dev/CLAUDE.md](docs/dev/CLAUDE.md)

**All documentation:** See [docs/README.md](docs/README.md)

---

## 🎯 What This Is

A **discipline enforcement system** that makes bad trades impossible through:

- ✅ **6-item checklist** → GREEN/YELLOW/RED banner (only GREEN allows saves)
- ✅ **2-minute impulse brake** → Mandatory pause after evaluation
- ✅ **Heat caps** → Portfolio (4%) and bucket (1.5%) limits enforced
- ✅ **Candidate validation** → Only trade tickers from today's FINVIZ list
- ✅ **Bucket cooldowns** → Sector restrictions after losses

**The 5 Hard Gates:** Every trade must pass ALL gates before execution. No bypasses. No exceptions.

---

## 📁 Project Structure

```
excel-trading-platform/
├── README.md                  # This file
├── docs/                      # All documentation (organized)
│   ├── README.md              # Documentation index
│   ├── project/               # Core project docs (WHY, PLAN)
│   ├── dev/                   # Development guides
│   └── milestones/            # Completion reports
├── cmd/tf-engine/             # Go CLI entry point
├── internal/                  # Go backend (all business logic)
│   ├── domain/                # Position sizing, checklist, heat, gates
│   ├── storage/               # SQLite persistence
│   ├── api/                   # HTTP handlers
│   └── cli/                   # CLI commands
├── features/                  # BDD Gherkin scenarios
├── excel/                     # VBA modules and workbook
│   ├── vba/                   # VBA modules (.bas text files)
│   └── VBA_MODULES_README.md  # VBA documentation
├── windows/                   # Windows deployment package
│   ├── tf-engine.exe          # Windows binary
│   ├── WINDOWS_TESTING.md     # M21 testing guide
│   └── ...                    # Setup scripts and docs
└── test-data/                 # Test fixtures and JSON examples
```

---

## 🏗️ Architecture

```
Excel UI (TradingPlatform.xlsm)
    ↓
VBA Bridge (thin layer - just shell exec + JSON parsing)
    ↓
tf-engine.exe (Go backend - ALL business logic)
    ↓
SQLite Database (trading.db - single source of truth)
```

**Key Principles:**
- **Engine-first**: All trading logic in Go backend
- **CLI by default**: Simple, reliable shell execution
- **HTTP optional**: Same handlers, enables future web UI
- **Thin VBA**: No business logic in Excel
- **Fail loudly**: Errors are never silently ignored

---

## 🛠️ Development Commands

```bash
# Build
go build -o tf-engine ./cmd/tf-engine

# Run tests
go test ./...                    # Unit tests
godog run features/              # BDD tests

# Position sizing example
./tf-engine size --entry 180 --atr 1.5 --k 2 --method stock

# Checklist evaluation
./tf-engine checklist --ticker AAPL --checks true,true,true,true,true,true

# Heat management (check portfolio and bucket heat)
./tf-engine heat --risk 75 --bucket "Tech/Comm"

# Save decision (enforces 5 hard gates)
./tf-engine save-decision --ticker AAPL --entry 180 --atr 1.5 \
  --k 2 --method stock --risk 75 --shares 25 --contracts 0 \
  --banner GREEN --bucket "Tech/Comm" --preset TF_BREAKOUT_LONG

# Initialize database
./tf-engine init

# Get settings
./tf-engine get-settings

# HTTP server (optional)
./tf-engine server --listen 127.0.0.1:18888
```

All commands support `--format json` for programmatic use.

---

## 📊 Current Status

**Phase:** M20 Complete ✅
**Next:** M21 - Windows Integration Validation (manual testing)

**Completed Milestones:**
- ✅ M1-M16: Go engine with all business logic
- ✅ M17-M18: JSON contract validation (CLI + HTTP parity)
- ✅ M19: VBA implementation (4 modules, 14 unit tests)
- ✅ M20: Windows integration package (ready for deployment)

**Next Milestone:**
- ⏸️ M21: Windows manual testing (~45 min to 4 hours)

See [docs/milestones/](docs/milestones/) for detailed completion reports.

---

## 🧪 Testing Strategy

**BDD Tests (Gherkin):**
```bash
godog run features/
```
- 50+ scenarios covering all business rules
- Position sizing (Van Tharp method)
- Checklist → banner logic
- Heat management
- 5 hard gates enforcement

**Unit Tests:**
```bash
go test ./...
```
- Domain logic validation
- Edge cases
- Error handling

**VBA Tests:**
- 14 VBA unit tests in TFTests.bas
- JSON parsing verification
- Shell execution validation

**Integration Tests:**
- CLI smoke tests (automated)
- HTTP parity tests (CLI vs HTTP)
- Windows manual tests (M21)

---

## 📖 Key Documents

### Core Philosophy
- [docs/project/WHY.md](docs/project/WHY.md) - **Read this first** ⭐
  - Why discipline over flexibility
  - Psychology of impulsive trading
  - Ed Seykota's philosophy
  - System design rationale

### Architecture & Plan
- [docs/project/PLAN.md](docs/project/PLAN.md)
  - Complete step-by-step plan
  - Architecture decisions
  - Data model
  - Milestone breakdown (M1-M21)

### Development Guides
- [docs/dev/DEVELOPMENT_PHILOSOPHY.md](docs/dev/DEVELOPMENT_PHILOSOPHY.md)
  - How we build this system
  - Code quality standards
  - Anti-patterns to reject
- [docs/dev/BDD_GUIDE.md](docs/dev/BDD_GUIDE.md)
  - Gherkin scenario writing
  - Testing workflow
- [docs/dev/CLAUDE.md](docs/dev/CLAUDE.md)
  - Claude Code guidance for this project

### Windows Deployment
- [windows/README.md](windows/README.md) - Windows package overview
- [windows/WINDOWS_TESTING.md](windows/WINDOWS_TESTING.md) - M21 testing guide (23 KB)
- [windows/EXCEL_WORKBOOK_TEMPLATE.md](windows/EXCEL_WORKBOOK_TEMPLATE.md) - Workbook spec

### VBA Implementation
- [excel/VBA_MODULES_README.md](excel/VBA_MODULES_README.md) - Complete VBA documentation
- [docs/milestones/M19_COMPLETION_SUMMARY.md](docs/milestones/M19_COMPLETION_SUMMARY.md) - M19 report

---

## 🤝 Contributing

**Before contributing:**
1. Read [docs/project/WHY.md](docs/project/WHY.md) - Understand the purpose
2. Read [docs/dev/DEVELOPMENT_PHILOSOPHY.md](docs/dev/DEVELOPMENT_PHILOSOPHY.md) - Understand the approach
3. Check anti-patterns - Know what to reject

**Development workflow:**
1. **Write Gherkin scenario first** (see [docs/dev/BDD_GUIDE.md](docs/dev/BDD_GUIDE.md))
2. Get agreement on behavior
3. Implement code
4. Write tests matching Gherkin
5. Verify behavior

**Questions to ask before adding a feature:**
- Does this support discipline or undermine it?
- Would Ed Seykota approve?
- Does it make impulsivity easier or harder?
- Is this solving a real problem or adding complexity?

**If unsure:** Read [docs/project/WHY.md](docs/project/WHY.md) again.

---

## 🎓 Core Concepts

### The 5 Hard Gates
Every trade must pass ALL gates before execution:
1. **Banner GREEN** - All 6 checklist items confirmed
2. **Ticker in today's candidates** - Must be from FINVIZ preset
3. **2-minute impulse brake** - Mandatory pause after evaluation
4. **Bucket not in cooldown** - Sector restrictions after losses
5. **Heat caps not exceeded** - Portfolio (4%) and bucket (1.5%) limits

**These gates CANNOT be bypassed.** They're enforced in the Go engine, not VBA.

### Position Sizing (Van Tharp Method)
```
1. Risk dollars:    R = Equity × RiskPct (0.75%)
2. Stop distance:   StopDist = K × ATR (K=2)
3. Initial stop:    InitStop = Entry - StopDist
4. Shares:          Shares = floor(R ÷ StopDist)
5. Verify:          ActualRisk = Shares × StopDist ≤ R
```

### Checklist → Banner Logic
- 0 missing → **GREEN** (go)
- 1 missing → **YELLOW** (caution)
- 2+ missing → **RED** (no-go)

Only GREEN starts the impulse timer and allows eventual save.

### Heat Management
- **Portfolio heat** = sum of risk across all open positions
- **Portfolio cap** = Equity × 4%
- **Bucket heat** = sum of risk within one sector bucket
- **Bucket cap** = Equity × 1.5%

Any trade exceeding either cap is rejected.

---

## 🔧 Technical Stack

**Backend (Go):**
- CLI framework: cobra
- Database: SQLite (mattn/go-sqlite3)
- HTTP: net/http (stdlib)
- BDD: godog
- Logging: logrus

**Frontend (Excel):**
- Excel desktop (Windows 10/11)
- VBA (thin bridge - no business logic)
- WScript.Shell for CLI execution
- Simple JSON parsing (no external dependencies)

**Database:**
- SQLite (single file: trading.db)
- Tables: settings, candidates, decisions, positions, cooldowns
- Migrations managed by Go code

---

## 📞 Support & Troubleshooting

**Documentation:**
- See [docs/README.md](docs/README.md) for complete documentation index
- Windows testing: [windows/WINDOWS_TESTING.md](windows/WINDOWS_TESTING.md)
- VBA specifics: [excel/VBA_MODULES_README.md](excel/VBA_MODULES_README.md)

**Logs:**
- `tf-engine.log` - Go backend logs (JSON, with correlation IDs)
- `TradingSystem_Debug.log` - VBA logs (text, with correlation IDs)
- Cross-reference using correlation IDs

**Common Issues:**
- Check [windows/WINDOWS_TESTING.md](windows/WINDOWS_TESTING.md) Troubleshooting section
- Verify correlation IDs in both log files
- Ensure database initialized: `./tf-engine init`

---

## 📜 License

[Add license information here]

---

## 🎯 Project Values

1. **Discipline** - The system enforces rules
2. **Simplicity** - Less code, less can break
3. **Clarity** - Errors are obvious and actionable
4. **Reliability** - Works every time, the same way
5. **Maintainability** - Future you can understand it

**Note what's NOT on the list:**
- Flexibility (discipline requires constraints)
- Power features (more features = more ways to fail)
- Customization (the rules are the rules)
- Convenience (friction is intentional)

---

## 🧭 The Guiding Question

**"Does this feature help or hurt discipline?"**

- If it helps → Build it
- If it hurts → Don't build it
- If unclear → Read [docs/project/WHY.md](docs/project/WHY.md) again

---

**Remember:** This is a discipline enforcement system, not a flexible trading platform. The system's value comes from what it prevents (bad trades), not what it allows.

**Code serves discipline. Discipline does not serve code.**
