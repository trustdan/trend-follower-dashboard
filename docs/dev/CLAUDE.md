# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Excel-based trading platform enforcing disciplined trend-following through systematic constraints.**

This is not a flexible trading platform. It's a discipline enforcement system that makes impulsive trading impossible. The system's value comes from what it prevents (bad trades) not what it allows.

### Core Philosophy

- **Discipline over flexibility** - Constraints are features, not limitations
- **Behavior-driven development** - Gherkin scenarios define "done"
- **Fail loudly** - Silent failures are unacceptable
- **Simple over clever** - Boring code is good code
- **Process enforcement** - Trading is a protocol, not performance

## Project Status

**Phase: Planning & Foundation**

Currently in the design phase for v3 architecture. No code implementation exists yet. The plan calls for:

- **Backend Engine**: Go (single-binary CLI + optional HTTP server)
- **Transport**: CLI-first via JSON stdout/stderr, HTTP optional
- **Storage**: SQLite for all state
- **Frontend**: Excel as thin UI, VBA as bridge layer
- **Testing**: BDD (Gherkin) + unit + integration tests

## Architecture (Planned v3)

### Core Concepts

**Engine-first architecture**: All trading logic lives in a compiled backend. Excel is just a UI.

```
Excel (UI)
  ↓ VBA (thin bridge)
    ↓ CLI or HTTP
      ↓ Go Engine (business logic)
        ↓ SQLite (state)
```

### The 5 Hard Gates

Every trade must pass ALL gates before execution:

1. **Banner GREEN** - All 6 checklist items confirmed
2. **Ticker in today's candidates** - Must be from FINVIZ preset
3. **2-minute impulse brake** - Mandatory pause after evaluation
4. **Bucket not in cooldown** - Sector restrictions after losses
5. **Heat caps not exceeded** - Portfolio (4%) and bucket (1.5%) limits

These gates CANNOT be bypassed. They're enforced in the engine, not VBA.

### Key Business Rules

**Position Sizing (Van Tharp method)**:
```
1. Calculate risk dollars: R = Equity × RiskPct (0.75%)
2. Calculate stop distance: StopDist = K × ATR (K=2)
3. Calculate initial stop: InitStop = Entry - StopDist
4. Calculate shares: Shares = floor(R ÷ StopDist)
5. Verify: ActualRisk = Shares × StopDist ≤ R
```

**Checklist → Banner Logic**:
- 0 missing → GREEN (go)
- 1 missing → YELLOW (caution)
- 2+ missing → RED (no-go)

Only GREEN starts the impulse timer and allows eventual save.

**Heat Management**:
- Portfolio heat = sum of risk across all open positions
- Portfolio cap = Equity × 4%
- Bucket heat = sum of risk within one sector bucket
- Bucket cap = Equity × 1.5%

Any trade exceeding either cap is rejected.

## Development Commands

### Not Yet Implemented

The system is in planning phase. When implemented:

```bash
# Initialize database and config
tf-engine.exe init

# Run position sizing calculation
tf-engine.exe size --entry 180 --atr 1.5 --k 2 --method stock

# Evaluate checklist
tf-engine.exe checklist --ticker AAPL --checks [...]

# Check heat
tf-engine.exe heat --add-r 75 --bucket "Tech/Comm"

# Save decision (enforces 5 hard gates)
tf-engine.exe save-decision --ticker AAPL --entry 180 [...]

# Import candidates
tf-engine.exe import-candidates --preset TF_BREAKOUT_LONG

# Scrape FINVIZ
tf-engine.exe scrape-finviz --query "..."

# Start HTTP server (optional)
tf-engine.exe server --listen 127.0.0.1:18888

# Run tests
go test ./...
godog run features/
```

## Critical Development Rules

### 1. Gherkin First, Always

**NO CODE without Gherkin scenario first.**

Every feature starts with behavior specification:

```gherkin
Scenario: Portfolio heat exceeds cap
  Given portfolio heat is $420
  And portfolio cap is $400
  When I try to add a $75 trade
  Then I should see error "Portfolio heat exceeds cap by $45"
  And the trade should be rejected
```

Write scenario → Get agreement → Implement code → Verify behavior.

### 2. Question Every Feature Request

Ask these questions BEFORE implementing:

- **Does this support discipline or undermine it?**
- **Would Ed Seykota approve?**
- **Does it make impulsivity easier or harder?**
- **Is this solving a real problem or adding complexity?**

If it makes impulsivity easier, **push back**. That's your job.

### 3. Anti-Patterns to Reject Immediately

❌ "Let's make it configurable" → NO. Hard-code the rules.
❌ "Let's add a bypass for edge cases" → NO. Document it, don't build backdoors.
❌ "Let's make it more flexible" → NO. Flexibility = opportunity for impulsivity.
❌ "Let's add this convenience feature" → Only if it reduces technical complexity without reducing discipline.
❌ "This is how other systems do it" → We're not building other systems.
❌ "But what if the user wants to..." → Read ../project/WHY.md. This system constrains users. That's the point.

### 4. Code Quality Standards

**Every feature must have:**
1. Gherkin scenario
2. Unit tests
3. Integration test
4. Error handling
5. Logging with correlation IDs
6. Documentation

**Code style:**
- Simple over clever
- Explicit over implicit
- Verbose over terse
- Named constants (no magic numbers)
- Clear error messages with actual values

**Bad:**
```go
if heat > 400 { return errors.New("invalid") }
```

**Good:**
```go
const PORTFOLIO_HEAT_CAP = 400 // 4% of $10,000 account
if portfolioHeat > PORTFOLIO_HEAT_CAP {
    return fmt.Errorf("portfolio heat $%.2f exceeds cap $%.2f (overage: $%.2f)",
        portfolioHeat, PORTFOLIO_HEAT_CAP, portfolioHeat - PORTFOLIO_HEAT_CAP)
}
```

### 5. Error Messages Must Teach

Every error should:
- State what's wrong
- Show the actual values
- Show the limit/expectation
- Suggest how to fix it

**Bad:** "Invalid input"
**Good:** "Portfolio heat ($425) exceeds cap ($400). Reduce position size or close existing positions."

### 6. Testing Philosophy

**Unit tests** - Test one thing, fast (<10ms), independent
**Integration tests** - Test component interaction with real DB
**BDD tests** - Match Gherkin scenarios exactly
**Parity tests** - CLI and HTTP return identical JSON

Test behavior, not implementation. Tests should describe what the system does, not how.

## Repository Structure

```
excel-trading-platform/
├─ ../project/WHY.md                           # Core philosophy - READ THIS FIRST
├─ DEVELOPMENT_PHILOSOPHY.md        # How we build
├─ BDD_GUIDE.md                     # How we test behavior
├─ CLAUDE_RULES.md                  # Detailed development rules
├─ ../project/PLAN.md  # Implementation roadmap
├─ README.md                        # Project overview
└─ (future)
   ├─ engine/                       # Go backend
   │  ├─ cmd/tf-engine/            # CLI + server
   │  ├─ internal/domain/          # Business logic
   │  ├─ internal/storage/         # SQLite layer
   │  └─ features/                 # BDD tests
   ├─ excel/                        # Excel workbook + VBA
   └─ docs/                         # User documentation
```

## Key Files to Read Before Coding

**MANDATORY READING ORDER:**

1. **../project/WHY.md** (5 min) - Understand the psychology and purpose
2. **DEVELOPMENT_PHILOSOPHY.md** (10 min) - Understand the approach
3. **BDD_GUIDE.md** (15 min) - Understand how we test
4. **../project/PLAN.md** (20 min) - Understand the architecture

**Do not skip these.** They define what success looks like.

## Context Switching Checklist

When returning to this project after time away:

```
[ ] Re-read ../project/WHY.md (5 min) - Remember the purpose
[ ] Review ../project/PLAN.md - Current architecture
[ ] Check recent commits - Understand what's been done
[ ] Run tests - Verify everything works
[ ] Read relevant Gherkin scenarios - Understand expected behavior

NOW you can code.
```

## Standard Patterns

### Position Sizing Always Follows This Pattern:
```
1. Validate inputs (entry, ATR, K must be positive)
2. Calculate stop distance (K × ATR)
3. Calculate initial stop (entry - stop distance)
4. Calculate shares (risk ÷ stop distance, rounded down)
5. Verify actual risk ≤ specified risk
6. Return result with all components
```

Never deviate. This is the Van Tharp method.

### Heat Management Always Follows This Pattern:
```
1. Sum risk across all open positions
2. Add proposed new position risk
3. Compare to portfolio cap (equity × heat_pct)
4. Compare to bucket cap (equity × bucket_pct)
5. Reject if either exceeded
6. Return detailed breakdown
```

Never allow trades that exceed caps. No exceptions.

### Checklist Validation Always Follows This Pattern:
```
1. Count missing checks
2. If 0 missing → GREEN (go)
3. If 1 missing → YELLOW (caution)
4. If 2+ missing → RED (no-go)
5. Return banner, missing count, missing items
6. Start impulse timer only on GREEN
```

Never skip checklist. Never allow save without GREEN.

### Data Flow Always Follows This Pattern:
```
Excel (input)
  → VBA (call backend)
    → Backend (calculate)
      → Database (persist)
    ← Backend (return JSON)
  ← VBA (parse JSON)
← Excel (display result)
```

Keep VBA thin. Keep Excel dumb. Keep backend smart.

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

## Decision-Making Framework

When uncertain about a feature or approach:

**1. Does it support discipline?**
   - If yes → Consider it
   - If no → Reject it
   - If unclear → Read ../project/WHY.md

**2. Can it be tested with Gherkin?**
   - If yes → Good sign
   - If no → Probably too vague

**3. Is it as simple as possible?**
   - If yes → Good
   - If no → Simplify first

**4. Will you understand it in 6 months?**
   - If yes → Ship it
   - If no → Rewrite it

**5. Would you trust it with your money?**
   - If yes → Ship it
   - If no → Fix it

## Git Commit Format

```
[Type] Brief description

- Detail 1
- Detail 2
- Reasoning for change
- Reference to Gherkin scenario if applicable
```

Types: `[feat]` `[fix]` `[refactor]` `[test]` `[docs]` `[chore]`

Example:
```
[feat] Add portfolio heat validation

- Implement heat cap checking
- Reject trades exceeding 4% cap
- Return detailed error with overage amount
- Implements: features/heat-management.feature:15
```

## Communication Style

**With users (error messages):**
- Speak human, not computer
- Show actual values and limits
- Suggest fixes

**With future self (comments):**
- Comment the why, not the what
- Explain assumptions
- Document edge cases

**With other developers (PRs):**
- Provide context (what, why, how)
- Reference Gherkin scenarios
- Ask questions about uncertainties

## The Meta-Rule

**When in doubt, read ../project/WHY.md.**

If the answer isn't there, you don't understand the question well enough yet.

Think more. Code less.

Every line of code is a liability. Every feature is a potential failure point. Every configuration is a decision users must make.

**Minimize all three.**

## Remember

This is not a software project that happens to be about trading.

This is a discipline system that happens to be implemented in software.

**Code serves discipline. Discipline does not serve code.**

Act accordingly.
