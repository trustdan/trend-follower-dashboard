# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Notifications
**CRITICAL:** When tasks complete OR when prompting me for user input (yes/no, confirmation, choices), notify me using:

**Primary method (Windows Toast Notification):**
```bash
/mnt/c/Windows/System32/WindowsPowerShell/v1.0/powershell.exe -Command "
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

\$template = @'
<toast>
    <visual>
        <binding template=\"ToastText02\">
            <text id=\"1\">TF-Engine</text>
            <text id=\"2\">Task complete - ready for your input</text>
        </binding>
    </visual>
</toast>
'@

\$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
\$xml.LoadXml(\$template)
\$toast = [Windows.UI.Notifications.ToastNotification]::new(\$xml)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('Claude Code').Show(\$toast)
"
```

**Use cases:**
- ✓ Task completion (build finished, tests passed, files updated)
- ✓ Prompting for user decisions (yes/no, confirmation, choose option A/B)
- ✓ Asking clarifying questions before proceeding
- ✗ NOT for simple status updates or informational messages

**Customize message:** Change the second `<text id=\"2\">` line to reflect the context:
- Task done: "Task complete - ready for your input"
- Question: "Question - your response needed"
- Error/blocker: "Action required - awaiting decision"

## Project Overview

**TF-Engine: Fresh Start Trading Platform** - A trend-following trading system that enforces discipline through systematic constraints.

**TF = Trend Following** - The engine implements Ed Seykota/Turtle Trader style trend-following with Donchian breakouts, ATR-based position sizing, and mechanical exits.

This is a "fresh start" project replacing a problematic Excel/VBA frontend while keeping the proven Go backend.

**Status:** Backend 100% functional, Frontend to be built (custom GUI application)

**Core Philosophy:** This is not a flexible trading platform. It's a discipline enforcement system that makes impulsive trading impossible. The system's value comes from what it prevents (bad trades), not what it allows.

## Anti-Impulsivity Design Principles

From `docs/anti-impulsivity.md`:

1. **Trade the tide, not the splash** - Donchian breakouts with mechanical exits
2. **Friction where it matters** - Hard gates for signal, risk, liquidity, exit, behavior
3. **Nudge for better trades** - Optional quality items affect score, not permission
4. **Immediate feedback** - Large 3-state banner (RED/YELLOW/GREEN) updates live
5. **Journal while deciding** - One-click logging of full trade plan
6. **Calendar awareness** - 10-week sector view for diversification

### The 5 Hard Gates (Cannot Be Bypassed)

1. **Signal:** 55-bar breakout (long > 55-high / short < 55-low)
2. **Risk/Size:** Per-unit risk = % of equity using 2×N stop; pyramids every 0.5×N to max units
3. **Options:** 60–90 DTE, roll/close ~21 DTE, liquidity required
4. **Exits:** Exit by 10-bar opposite Donchian OR closer of 2×N
5. **Behavior:** 2-minute cool-off + no intraday overrides

**Banner States:**
- **RED:** Any required gate fails → DO NOT TRADE
- **YELLOW:** All required pass, quality score < threshold → CAUTION
- **GREEN:** All required pass, quality score ≥ threshold → OK TO TRADE

## Architecture

```
Current State (Working):
  Backend (Go) - tf-engine
    ├─ cmd/tf-engine/        - CLI entry point
    └─ internal/
       ├─ domain/            - Business logic (sizing, checklist, heat, gates)
       ├─ storage/           - SQLite persistence
       ├─ scrape/            - FINVIZ scraper
       ├─ cli/               - Command handlers
       └─ server/            - HTTP server (legacy, for Excel)

Future State (To Build):
  Custom GUI Frontend (Go + Fyne/Gio)
    └─ Direct in-process calls to backend
```

**Key Decision:** GUI will make direct function calls to backend (no HTTP, no CLI spawning). Single binary deployment.

## Development Commands

### Backend Development

**Linux/macOS (bash):**
```bash
# Navigate to backend
cd backend/

# Build the binary
go build -o tf-engine cmd/tf-engine/main.go

# Run tests (comprehensive test suite)
go test ./internal/domain/... -v
go test ./internal/storage/... -v
go test ./... -v

# Initialize database
./tf-engine init

# Configure settings
./tf-engine settings --equity 100000 --risk-pct 0.75 --portfolio-cap 4.0

# Position sizing calculation
./tf-engine size --ticker AAPL --entry 180 --atr 1.5 --method stock --k 2

# Checklist evaluation
./tf-engine checklist --ticker AAPL --from-preset true --trend true --liquidity true --timeframe true --earnings true --journal true

# Heat check
./tf-engine heat --risk 750 --bucket "Tech/Comm"

# Import candidates from FINVIZ
./tf-engine import-candidates --preset TF_BREAKOUT_LONG

# Start HTTP server (legacy)
./tf-engine server --listen 127.0.0.1:18888
```

**Windows (PowerShell):**
```powershell
# Navigate to backend
cd backend/

# Build the binary
go build -o tf-engine.exe cmd/tf-engine/main.go

# Run tests (comprehensive test suite)
go test ./internal/domain/... -v
go test ./internal/storage/... -v
go test ./... -v

# Initialize database
.\tf-engine.exe init

# Configure settings
.\tf-engine.exe settings --equity 100000 --risk-pct 0.75 --portfolio-cap 4.0

# Position sizing calculation
.\tf-engine.exe size --ticker AAPL --entry 180 --atr 1.5 --method stock --k 2

# Checklist evaluation
.\tf-engine.exe checklist --ticker AAPL --from-preset true --trend true --liquidity true --timeframe true --earnings true --journal true

# Heat check
.\tf-engine.exe heat --risk 750 --bucket "Tech/Comm"

# Import candidates from FINVIZ
.\tf-engine.exe import-candidates --preset TF_BREAKOUT_LONG

# Start HTTP server (legacy)
.\tf-engine.exe server --listen 127.0.0.1:18888
```

### Cross-Platform Building

**From Linux/macOS (bash):**
```bash
cd backend/
# Build for Windows
GOOS=windows GOARCH=amd64 go build -o tf-engine.exe cmd/tf-engine/main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o tf-engine cmd/tf-engine/main.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o tf-engine-mac cmd/tf-engine/main.go
```

**From Windows (PowerShell):**
```powershell
cd backend/
# Build for Windows (native - no env vars needed)
go build -o tf-engine.exe cmd/tf-engine/main.go

# Build for Linux
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o tf-engine cmd/tf-engine/main.go

# Build for macOS
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o tf-engine-mac cmd/tf-engine/main.go
```

### Testing Strategy

```bash
# Run all backend tests
cd backend/
go test ./... -v

# Run specific domain tests
go test ./internal/domain/ -v -run TestCalculateSize

# Run storage tests
go test ./internal/storage/ -v

# Check test coverage
go test ./... -cover
```

## Core Business Logic

### Position Sizing (Van Tharp Method)

From `backend/internal/domain/sizing.go` and `sizing_stock.go`:

```
1. Calculate risk dollars: R = Equity × RiskPct (typically 0.75%)
2. Calculate stop distance: StopDist = K × ATR (K=2)
3. Calculate initial stop: InitStop = Entry - StopDist (for longs)
4. Calculate shares: Shares = floor(R ÷ StopDist)
5. Verify: ActualRisk = Shares × StopDist ≤ R
```

**Methods supported:**
- `stock` - Direct stock/ETF trading
- `opt-delta-atr` - Options using delta-adjusted ATR
- `opt-contracts` - Options using contract risk

### Checklist Validation

From `backend/internal/domain/checklist.go`:

```
Count missing required checks:
- 0 missing → GREEN (go)
- 1 missing → YELLOW (caution)
- 2+ missing → RED (no-go)
```

**Required items:** Signal, Risk/Size, Options/Liquidity, Exits, Behavior
**Optional items:** Regime, Chase, Earnings, Journal (affect quality score only)

### Heat Management

From `backend/internal/domain/heat.go`:

```
Portfolio heat = sum of risk across all open positions
Portfolio cap = Equity × 4.0%
Bucket heat = sum of risk within one sector bucket
Bucket cap = Equity × 1.5%

Reject any trade that exceeds either cap.
```

### 5 Gates Check

From `backend/internal/domain/gates.go`:

```
Gate 1: Banner must be GREEN
Gate 2: 2-minute cooloff elapsed since checklist evaluation
Gate 3: Ticker not on cooldown (from recent losses)
Gate 4: Heat caps not exceeded (portfolio and bucket)
Gate 5: Position sizing completed

All gates must pass to save GO decision.
```

## Directory Structure

```
fresh-start-trading-platform/
├── README.md                    # Project overview
├── FRESH_START_PLAN.md         # Detailed GUI implementation plan
├── PROJECT_HISTORY.md          # Why we abandoned Excel/VBA
├── ORGANIZATION_SUMMARY.md     # Project organization notes
│
├── backend/                    # ✅ 100% FUNCTIONAL Go backend
│   ├── cmd/tf-engine/         # CLI entry point
│   ├── internal/
│   │   ├── api/               # API types
│   │   ├── cli/               # Command handlers
│   │   ├── domain/            # Core business logic ⭐
│   │   │   ├── sizing.go          # Position sizing algorithms
│   │   │   ├── checklist.go       # Checklist validation
│   │   │   ├── heat.go            # Heat management
│   │   │   ├── gates.go           # 5 gates enforcement
│   │   │   ├── settings.go        # Account settings
│   │   │   └── candidates.go      # Ticker candidates
│   │   ├── storage/           # SQLite persistence ⭐
│   │   │   ├── db.go              # Database interface
│   │   │   ├── positions.go       # Position CRUD
│   │   │   ├── decisions.go       # Trade decisions
│   │   │   ├── candidates.go      # Candidate tickers
│   │   │   └── cooldowns.go       # Cooldown tracking
│   │   ├── scrape/            # FINVIZ web scraping
│   │   ├── server/            # HTTP server (legacy)
│   │   └── logx/              # Logging utilities
│   ├── go.mod
│   └── go.sum
│
├── scripts/                    # Utility scripts
│   ├── import-candidates-auto.bat   # Automated FINVIZ import
│   ├── import-candidates.bat        # Manual FINVIZ import
│   ├── build-windows.sh             # Cross-compile for Windows
│   └── test-finviz-scrape.sh        # Test scraper
│
├── docs/                       # Comprehensive documentation
│   ├── anti-impulsivity.md    # Core design philosophy ⭐⭐⭐
│   ├── PROJECT_STATUS.md      # Current status (M24 complete)
│   ├── UI_QUICK_REFERENCE.md  # UI reference guide
│   ├── HTTP_CLI_PARITY.md     # API documentation
│   ├── dev/
│   │   ├── DEVELOPMENT_PHILOSOPHY.md  # How we build ⭐⭐
│   │   ├── CLAUDE_RULES.md            # Development rules ⭐
│   │   └── BDD_GUIDE.md               # BDD testing approach
│   ├── project/
│   │   ├── WHY.md             # Why this system exists ⭐⭐⭐
│   │   └── PLAN.md            # Original implementation plan
│   ├── milestones/            # Milestone completion summaries
│   └── json-schemas/          # JSON API specifications
│
├── test-data/                  # Test fixtures and examples
│   ├── test-contracts.db       # Sample database
│   ├── json-examples/          # JSON test cases
│   └── phase4-test-data.sql    # SQL fixtures
│
└── art/                        # ASCII art and assets
```

## Key Files to Read Before Coding

**MANDATORY READING ORDER:**

1. **`docs/project/WHY.md`** (5 min) - Understand the psychology and purpose
2. **`docs/anti-impulsivity.md`** (10 min) - Core design principles
3. **`docs/dev/DEVELOPMENT_PHILOSOPHY.md`** (10 min) - How we build
4. **`docs/dev/CLAUDE_RULES.md`** (10 min) - Development standards
5. **`README.md`** (5 min) - Project overview and quick start
6. **`FRESH_START_PLAN.md`** (15 min) - GUI implementation plan

**Do not skip these.** They define what success looks like.

## Critical Development Rules

### 1. Discipline Over Flexibility

**Every feature request must be evaluated:**
- Does this support discipline or undermine it?
- Would Ed Seykota approve?
- Does it make impulsivity easier or harder?
- Is this solving a real problem or adding complexity?

**If it makes impulsivity easier, push back.** That's your job.

### 2. Anti-Patterns to Reject Immediately

❌ "Let's make it configurable" → NO. Hard-code the rules.
❌ "Let's add a bypass for edge cases" → NO. Document it, don't build backdoors.
❌ "Let's make it more flexible" → NO. Flexibility = opportunity for impulsivity.
❌ "Let's add this convenience feature" → Only if it reduces technical complexity without reducing discipline.
❌ "This is how other systems do it" → We're not building other systems.

### 3. Behavior-Driven Development (BDD)

**Gherkin scenarios define what "done" means.**

While the system doesn't currently use Gherkin extensively, the principle remains: specify behavior before implementation.

Example:
```gherkin
Scenario: Trader tries to exceed heat cap
  Given portfolio heat is at 90% of cap
  When trader calculates a new position
  Then the system should show "EXCEEDS CAP"
  And the Save button should be disabled
  And the trader should see exact overage amount
```

### 4. Code Quality Standards

**Simple over clever:**
```go
// Bad: Magic number
if heat > 400 { reject() }

// Good: Named constant with context
const PORTFOLIO_HEAT_CAP = 400 // 4% of $10,000 account
if portfolioHeat > PORTFOLIO_HEAT_CAP {
    return fmt.Errorf("portfolio heat $%.2f exceeds cap $%.2f (overage: $%.2f)",
        portfolioHeat, PORTFOLIO_HEAT_CAP, portfolioHeat - PORTFOLIO_HEAT_CAP)
}
```

**Error messages must teach:**
- State what's wrong
- Show actual values
- Show the limit/expectation
- Suggest how to fix it

Bad: "Invalid input"
Good: "Portfolio heat ($425) exceeds cap ($400). Reduce position size or close existing positions."

### 5. Fail Loudly

**Silent failures are unacceptable.**
- Invalid input? Reject with clear error immediately
- Calculation wrong? Show the math step-by-step
- Rule violated? Display exactly which rule and why
- System error? Log everything, show user-friendly message

## Standard Calculation Patterns

### Position Sizing Pattern

```
1. Validate inputs (entry, ATR, K must be positive)
2. Calculate stop distance (K × ATR)
3. Calculate initial stop (entry - stop distance for longs)
4. Calculate shares (risk ÷ stop distance, rounded down)
5. Verify actual risk ≤ specified risk
6. Return result with all components
```

Never deviate. This is the Van Tharp method.

### Heat Management Pattern

```
1. Sum risk across all open positions
2. Add proposed new position risk
3. Compare to portfolio cap (equity × 4.0%)
4. Compare to bucket cap (equity × 1.5%)
5. Reject if either exceeded
6. Return detailed breakdown
```

Never allow trades that exceed caps. No exceptions.

### Checklist Validation Pattern

```
1. Count missing required checks
2. If 0 missing → GREEN (go)
3. If 1 missing → YELLOW (caution)
4. If 2+ missing → RED (no-go)
5. Return banner, missing count, missing items
6. Start impulse timer only on GREEN
```

Never skip checklist. Never allow save without GREEN.

## Example Calculations

### Position Sizing Example

```
Given:
  Equity E = $10,000
  Risk per trade r = 0.75%
  Entry = $180
  ATR (N) = $1.50
  K multiple = 2

Calculate:
  R = $10,000 × 0.0075 = $75.00
  StopDistance = 2 × $1.50 = $3.00
  InitialStop = $180 - $3.00 = $177.00
  Shares = floor($75 ÷ $3.00) = 25 shares
  ActualRisk = 25 × $3.00 = $75.00 ✓
```

### Heat Management Example

```
Given:
  Equity = $10,000
  Portfolio cap = 4% = $400
  Bucket cap = 1.5% = $150
  Current portfolio heat = $350
  Current Tech/Comm bucket heat = $125
  New trade risk = $75 in Tech/Comm

Check:
  New portfolio heat = $350 + $75 = $425
  $425 > $400 → REJECT (exceeds portfolio cap by $25)
```

## Accessing Backend Functions (Future GUI)

When building the GUI, import and call backend functions directly:

```go
import (
    "github.com/yourusername/trading-engine/internal/domain"
    "github.com/yourusername/trading-engine/internal/storage"
)

// Initialize database
db, err := storage.NewDB("trading.db")

// Calculate position size
result, err := domain.CalculateSizeStock(
    ticker:   "AAPL",
    entry:    180.0,
    atr:      1.5,
    k:        2.0,
    riskPct:  0.0075,
    equity:   100000.0,
)

// Evaluate checklist
result, err := domain.EvaluateChecklist(
    fromPreset: true,
    trend: true,
    liquidity: true,
    timeframe: true,
    earnings: true,
    journal: true,
)

// Check heat
result, err := domain.CheckHeat(
    db,
    riskAmount: 750.0,
    bucket: "Tech/Comm",
    equity: 100000.0,
    portfolioCap: 4.0,
    bucketCap: 1.5,
)

// Check gates
result, err := domain.CheckGates(
    db,
    ticker: "AAPL",
    banner: "GREEN",
    riskDollars: 750.0,
    bucket: "Tech/Comm",
    equity: 100000.0,
    portfolioCap: 4.0,
    bucketCap: 1.5,
)
```

## GUI Implementation Guidance

### Technology Recommendation: Fyne

**Why Fyne:**
- Pure Go (same language as backend)
- Cross-platform (Windows, Linux, macOS)
- Material Design look
- Good documentation
- Easy single binary packaging
- Active community

**Quick Start:**
```bash
go get fyne.io/fyne/v2
go run hello.go  # See README.md for example
```

### 6 Main Screens to Build

1. **Dashboard** - Portfolio overview, positions, candidates, cooldowns
2. **Checklist** - 5 required gates + optional quality items + banner
3. **Position Sizing** - Calculate shares/contracts using ATR-based risk
4. **Heat Check** - Verify portfolio and sector heat within caps
5. **Trade Entry** - Final 5-gate check before saving GO/NO-GO decision
6. **Calendar** - Rolling 10-week sector × week grid for diversification

See `FRESH_START_PLAN.md` for detailed UI specifications.

### UI Design Principles

1. **Large, Obvious Banner** - RED/YELLOW/GREEN must be impossible to miss
2. **Friction for Important Actions** - "Save GO" should require multiple confirmations
3. **Immediate Feedback** - No waiting, no spinners for calculations
4. **Clear Error Messages** - Show what's wrong and how to fix it
5. **No Backdoors** - Cannot bypass gates, cannot skip cooldowns, cannot override caps

## Database Schema

SQLite database in `trading.db`:

**Key Tables:**
- `settings` - Account settings (equity, risk_pct, caps)
- `positions` - Open positions (risk, bucket, entry, stop)
- `decisions` - Trade decisions (GO/NO-GO with all gate results)
- `candidates` - Daily ticker candidates from FINVIZ
- `cooldowns` - Sector/ticker cooldowns after losses
- `evaluations` - Checklist evaluation timestamps (for 2-min timer)

See `backend/internal/storage/*.go` for table definitions.

## Git Workflow

**Commit Format:**
```
[Type] Brief description

- Detail 1
- Detail 2
- Reasoning for change
```

**Types:** `[feat]` `[fix]` `[refactor]` `[test]` `[docs]` `[chore]`

**Every commit should:**
- Pass all tests
- Be a deployable state
- Have a clear purpose
- Include updated documentation if needed

## Performance Standards

- Backend calculations: < 100ms
- Database operations: < 500ms
- Web scraping: < 5s
- GUI responsiveness: Must feel instant

## Context Switching Checklist

When returning to this project after time away:

```
[ ] Re-read docs/project/WHY.md (5 min)
[ ] Review docs/anti-impulsivity.md (5 min)
[ ] Check docs/PROJECT_STATUS.md (current: M24 complete, awaiting Windows validation)
[ ] Run backend tests: cd backend/ && go test ./...
[ ] Review recent commits to understand latest changes

NOW you can code.
```

## Decision-Making Framework

When uncertain about a feature or approach:

**1. Does it support discipline?**
   - If yes → Consider it
   - If no → Reject it
   - If unclear → Read `docs/project/WHY.md`

**2. Is it as simple as possible?**
   - If yes → Good
   - If no → Simplify first

**3. Will you understand it in 6 months?**
   - If yes → Ship it
   - If no → Rewrite it

**4. Would you trust it with your money?**
   - If yes → Ship it
   - If no → Fix it

## The Meta-Rule

**When in doubt, read `docs/project/WHY.md` and `docs/anti-impulsivity.md`.**

If the answer isn't there, you don't understand the question well enough yet.

Think more. Code less.

Every line of code is a liability. Every feature is a potential failure point. Every configuration is a decision users must make.

**Minimize all three.**

## Remember

This is not a software project that happens to be about trading.

This is a discipline system that happens to be implemented in software.

**Code serves discipline. Discipline does not serve code.**

Act accordingly.
