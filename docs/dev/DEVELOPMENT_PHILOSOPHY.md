# Development Philosophy

## Core Principles

### 1. Vibe Coding First
**Start with the experience, not the implementation.**

- How should it *feel* to use this system?
- What *should* happen when a trader opens it?
- Where *should* friction exist?
- What *should* be impossible to do?

Code the vibe before coding the features.

### 2. Behavior-Driven Development (BDD)
**Gherkin scenarios define what "done" means.**

```gherkin
Scenario: Trader tries to exceed heat cap
  Given portfolio heat is at 90% of cap
  When trader calculates a new position
  Then the system should show "EXCEEDS CAP"
  And the Save button should be disabled
  And the trader should see exact overage amount
```

If it's not in Gherkin, it's not a requirement. If it's in Gherkin, it's non-negotiable.

### 3. Constraints Are Features
**The system's power comes from what it prevents, not what it allows.**

- Cannot trade without checklist confirmation → FEATURE
- Cannot exceed heat caps → FEATURE
- Cannot bypass 2-minute timer → FEATURE
- Cannot skip position sizing → FEATURE

Limitations are the product.

### 4. Fail Loudly, Early, Obviously
**Silent failures are unacceptable.**

- Invalid input? Reject with clear error message immediately
- Calculation wrong? Show the math step-by-step
- Rule violated? Display exactly which rule and why
- System error? Log everything, show user-friendly message

Never fail silently. Never "just work around it."

### 5. Build for Actual Use, Not Theoretical Use
**Design for 2am impulsive trader, not rational 9am planner.**

The user is:
- Emotional
- Impulsive
- Tired
- Excited
- Wanting to bypass rules
- Looking for loopholes

Design for *that* person, not ideal-state disciplined trader.

---

## Technical Principles

### Code Is Communication
**Write for the next developer (probably future you).**

```go
// Bad: Magic number
if heat > 400 { reject() }

// Good: Named constant with context
const PORTFOLIO_HEAT_CAP = 400 // 4% of $10,000 account
if portfolioHeat > PORTFOLIO_HEAT_CAP {
    return fmt.Errorf("portfolio heat $%.2f exceeds cap $%.2f",
        portfolioHeat, PORTFOLIO_HEAT_CAP)
}
```

Every number has a name. Every name has a reason.

### Simple Over Clever
**Boring code is good code.**

```go
// Bad: Clever one-liner
r := func(e, p float64) float64 { return e * p }(equity, riskPct)

// Good: Obvious multi-line
func calculateRiskAmount(equity float64, riskPercent float64) float64 {
    return equity * riskPercent
}
riskAmount := calculateRiskAmount(equity, riskPct)
```

If you have to think about what code does, it's too clever.

### Explicit Over Implicit
**Make assumptions visible.**

```go
// Bad: Hidden assumption
func calculateShares(risk, stopDist float64) int {
    return int(risk / stopDist)
}

// Good: Assumption documented
func calculateShares(risk, stopDist float64) int {
    // Assumption: Always round down (floor) to avoid oversizing
    // Example: $75 risk / $3 stop = 25 shares (not 25.5)
    return int(math.Floor(risk / stopDist))
}
```

Future you will thank current you.

### Test Behavior, Not Implementation
**Tests should describe what the system does, not how.**

```go
// Bad: Testing implementation details
func TestSizingUsesFloorFunction(t *testing.T) {
    // Fragile - breaks if we change implementation
}

// Good: Testing behavior
func TestSizingNeverExceedsRiskAmount(t *testing.T) {
    result := calculateShares(75.0, 3.1)
    actual := result.Shares * result.StopDistance
    assert.LessOrEqual(t, actual, 75.0, "actual risk should never exceed specified risk")
}
```

Implementation can change. Behavior should not.

---

## Process Principles

### Iterate in Vertical Slices
**Complete one feature end-to-end before starting the next.**

Slice: Position sizing calculation
- [ ] Backend logic
- [ ] Unit tests
- [ ] CLI command
- [ ] Excel integration
- [ ] User acceptance test
- **DONE.** Ship it.

Then start next slice. No half-finished features.

### Commit Working Code Only
**Every commit should be a deployable state.**

```bash
# Bad: "WIP - broke everything"
git commit -m "half way through refactor"

# Good: Complete atomic change
git commit -m "Add position sizing for stock method with tests"
```

Main branch is always working. Always.

### Document Decisions, Not Just Code
**Why matters more than what.**

```markdown
# Architecture Decision Record: Use CLI Instead of HTTP

## Context
Need communication between Excel and backend.

## Decision
Use CLI with JSON output instead of HTTP REST API.

## Reasoning
- Simpler deployment (single .exe, no server)
- Easier debugging (can test with shell)
- No port conflicts
- Sufficient for single-user local use

## Consequences
- Each call spawns new process (slower)
- Cannot have persistent connections
- Would need rewrite for multi-user version

## Future
If we need multi-user or web UI, migrate to HTTP.
```

Future you needs to know *why* decisions were made.

### Question Everything, Assume Nothing
**"We always do it this way" is not a reason.**

- Why VBA at all? (Maybe we don't need it)
- Why Excel? (Maybe we don't need it)
- Why SQLite? (Maybe flat files work)
- Why Go? (Maybe Python is fine)

Question all assumptions. Keep only what has good answers.

---

## Anti-Patterns to Avoid

### ❌ "I'll fix it later"
**No you won't. Fix it now or delete it.**

Tech debt compounds. Every "I'll fix it later" is a lie you tell yourself.

### ❌ "It works on my machine"
**Then it doesn't work.**

If it requires specific setup, document it. If it can't be documented, fix it.

### ❌ "Just add a config option"
**Configuration is complexity in disguise.**

Every option is:
- A decision users must make
- A combination to test
- A support question waiting to happen

Default to zero configuration. Add options only when absolutely necessary.

### ❌ "We might need this someday"
**YAGNI - You Ain't Gonna Need It**

Don't build for hypothetical futures. Build for actual present needs.

When that future arrives, refactor. Code is meant to change.

### ❌ "This is how [framework/library] does it"
**Cargo cult programming.**

Don't copy patterns because they exist elsewhere. Understand why they exist, then decide if they apply here.

---

## Quality Standards

### Every Feature Must Have:
1. **Gherkin scenario** - What does success look like?
2. **Unit tests** - Do the pieces work in isolation?
3. **Integration test** - Do the pieces work together?
4. **Error handling** - What happens when it fails?
5. **Logging** - How do we debug issues?
6. **Documentation** - How do users use it?

Incomplete features don't ship. Period.

### Code Review Checklist:
- [ ] Does this solve the actual problem?
- [ ] Is it as simple as possible?
- [ ] Are edge cases handled?
- [ ] Is error handling clear?
- [ ] Are tests meaningful?
- [ ] Is documentation complete?
- [ ] Would you want to maintain this in 2 years?

If any answer is no, it's not ready.

---

## Learning Philosophy

### Spike, Then Implement
**Don't know how to do something? Experiment first.**

1. Create `spike/experiment/` directory
2. Try the approach, break things, learn
3. Once you understand it, delete the spike
4. Implement cleanly from scratch

Spikes are throwaway. Implementations are permanent. Keep them separate.

### Fail Fast, Learn Faster
**Mistakes are data.**

- Tried VBA? It has unicode issues. *Good to know.*
- Tried COM automation? It's unreliable. *Good to know.*
- Tried complex event handling? It creates loops. *Good to know.*

Each failure teaches what to avoid. Embrace failure as learning.

### Steal Smart
**Good artists copy, great artists steal.**

Don't reinvent position sizing math - Seykota defined it decades ago.
Don't reinvent risk management - Van Tharp wrote the book.
Don't reinvent trend following - the turtle traders proved it works.

Steal proven concepts. Invent only the glue code.

---

## Communication Philosophy

### With Users
**Speak human, not computer.**

```
# Bad: "Error 1004: Application-defined or object-defined error"
# Good: "Cannot save: Portfolio heat ($425) exceeds cap ($400)"

# Bad: "Invalid input"
# Good: "Entry price must be greater than 0. You entered: -5"

# Bad: "NULL reference exception"
# Good: "Cannot calculate heat: No open positions found in database"
```

Every error message is a teaching moment.

### With Future Self
**Comment the why, not the what.**

```go
// Bad:
// Loop through positions
for _, pos := range positions { ... }

// Good:
// Sum risk across all open positions to calculate current portfolio heat.
// Closed positions don't contribute to heat since capital is no longer at risk.
for _, pos := range positions {
    if pos.Status == "Open" { ... }
}
```

Code says what. Comments say why.

### With Other Developers
**Pull requests are conversations.**

```markdown
## What
Add heat cap validation to save-decision command.

## Why
Currently users can save trades that exceed heat caps. This violates core risk management rules.

## How
Check portfolio and bucket heat before allowing save. Reject with clear error if exceeded.

## Testing
Added test: TestSaveDecision_ExceedsHeatCap
Manual verification: Tried to save with heat at 110% - correctly rejected.

## Questions
Should we show how much heat would need to be reduced to allow trade?
```

Context helps reviewers understand. Understanding helps review quality.

---

## The Meta-Principle

### Build What You'd Want to Use
**Would you trust this system with your money?**

If the answer is no, fix it.

If the answer is "maybe," fix it.

Only "yes" ships.

This is not theoretical software. This is money. This is discipline. This is the difference between surviving and blowing up.

**Build accordingly.**

---

## Remember

> "Everyone has a plan until they get punched in the mouth." - Mike Tyson

The market will punch you. Your emotions will punch you. Your biases will punch you.

This system is your training regimen. Build it to withstand punches.

**Code like your account depends on it. Because it does.**
