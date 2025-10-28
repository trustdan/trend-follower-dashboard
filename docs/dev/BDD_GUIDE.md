# Behavior-Driven Development (BDD) Guide

## What Is BDD?

**Behavior-Driven Development is writing tests in human language before writing code.**

Instead of:
```go
func TestCalculateShares(t *testing.T) {
    result := calculateShares(75.0, 3.0)
    assert.Equal(t, 25, result)
}
```

We write:
```gherkin
Scenario: Calculate shares for stock position
  Given I have $75 risk
  And my stop distance is $3
  When I calculate position size
  Then I should get 25 shares
```

**Why?** Because humans understand stories, not assertions.

---

## Gherkin Syntax

### The Three Keywords

#### Given - Setup/Context
**The world before the action happens.**

```gherkin
Given I have an account size of $10000
And my risk per trade is 0.75%
And I have 3 open positions with $50 risk each
```

Sets up the initial state. No actions yet.

#### When - Action
**The thing you're testing.**

```gherkin
When I calculate position size for a $180 stock
And the ATR is 1.50
And the K multiple is 2
```

The trigger. The thing under test.

#### Then - Assertion
**What should happen as a result.**

```gherkin
Then the position size should be 25 shares
And the risk should be $75
And the stop should be at $177
And the trade should be allowed
```

Expected outcomes. If these don't happen, test fails.

### Optional: And, But
**Chain multiple conditions without repetition.**

```gherkin
Given I have an account size of $10000
And my risk per trade is 0.75%  # And = another Given
And I have no open positions

When I calculate position size
And I check portfolio heat       # And = another When

Then the shares should be 25
And the heat should be 18.75%    # And = another Then
But the heat should not exceed 100%  # But = Then with contrast
```

### Background
**Common setup for all scenarios in a feature.**

```gherkin
Feature: Position Sizing

Background:
  Given I have an account size of $10000
  And my risk per trade is 0.75%
  # This runs before EVERY scenario below

Scenario: Stock position
  When I calculate size for stock
  Then shares should be calculated

Scenario: Option position
  When I calculate size for options
  Then contracts should be calculated
```

Avoids repeating same Given steps.

### Scenario Outline / Examples
**Run same scenario with different data.**

```gherkin
Scenario Outline: Position sizing for various ATRs
  Given entry price is <entry>
  And ATR is <atr>
  And K multiple is 2
  When I calculate position size
  Then shares should be <shares>
  And risk should be approximately $75

  Examples:
    | entry | atr  | shares |
    | 180   | 1.50 | 25     |
    | 90    | 0.75 | 50     |
    | 360   | 3.00 | 12     |
```

One scenario, multiple test cases.

---

## Writing Good Scenarios

### ✅ Do: Focus on Behavior

```gherkin
# Good: Describes what happens
Scenario: Reject trade exceeding heat cap
  Given portfolio heat is $380
  When I try to add a $75 trade
  Then I should see error "Exceeds portfolio cap"
  And the trade should not be saved
```

**Why good:** Describes outcome, not implementation.

### ❌ Don't: Describe Implementation

```gherkin
# Bad: Describes how code works
Scenario: Heat cap check
  Given portfolioHeat variable is 380
  When addTrade function is called with 75
  Then calculateHeat method should return 455
  And shouldAllowTrade method should return false
```

**Why bad:** Tightly coupled to code structure. Breaks when refactoring.

---

### ✅ Do: Use Business Language

```gherkin
# Good: Trader/business terms
Scenario: Options trader uses Delta-ATR method
  Given I'm trading options
  And the underlying is at $180
  And I want 30-delta options
  When I calculate contracts
  Then I should risk $75 per contract
```

**Why good:** Anyone can read this - trader, PM, developer.

### ❌ Don't: Use Technical Jargon

```gherkin
# Bad: Developer terms
Scenario: Execute option_delta_atr algorithm
  Given params struct with method=2
  When sizing_calculator.calculate() is invoked
  Then result.contracts should be populated
```

**Why bad:** Only developers understand. Not truly BDD.

---

### ✅ Do: One Scenario, One Behavior

```gherkin
# Good: Tests one thing
Scenario: Portfolio heat exceeds cap
  Given portfolio heat is $420
  When I check if trade is allowed
  Then I should see error "Portfolio heat exceeds cap"
```

**Why good:** Clear what's being tested. Easy to debug when fails.

### ❌ Don't: Test Multiple Behaviors

```gherkin
# Bad: Tests everything at once
Scenario: Complete trade workflow
  Given I have account setup
  When I calculate sizing
  Then shares should be calculated
  And heat should be checked
  And checklist should be validated
  And database should be updated
  And email should be sent
  And logs should be written
```

**Why bad:** If it fails, which behavior broke? Hard to debug.

---

### ✅ Do: Make Expected Outcomes Explicit

```gherkin
# Good: Clear expectations
Scenario: Calculate initial stop
  Given entry price is $180
  And ATR is 1.50
  And K multiple is 2
  When I calculate stop
  Then stop distance should be $3.00
  And initial stop should be $177.00
  And stop percentage should be 1.67% below entry
```

**Why good:** No ambiguity. Pass/fail is obvious.

### ❌ Don't: Leave Outcomes Vague

```gherkin
# Bad: What does "correctly" mean?
Scenario: Calculate stop
  Given entry and ATR
  When I calculate stop
  Then stop should be calculated correctly
```

**Why bad:** "Correctly" means nothing. Not verifiable.

---

## Organizing Features

### Feature Files Structure

```
features/
├── position-sizing.feature
├── heat-management.feature
├── checklist-validation.feature
├── finviz-scraping.feature
└── trade-logging.feature
```

One feature per file. Related scenarios grouped together.

### Feature File Template

```gherkin
Feature: [Feature Name]
  As a [role]
  I want to [do something]
  So that [business value]

  Background:
    Given [common setup for all scenarios]

  Scenario: [Happy path - main use case]
    Given [context]
    When [action]
    Then [expected outcome]

  Scenario: [Edge case 1]
    Given [edge case context]
    When [action]
    Then [edge case outcome]

  Scenario: [Error case 1]
    Given [error condition]
    When [action]
    Then [error handling]

  Scenario Outline: [Multiple similar cases]
    Given [parameterized context]
    When [action with <parameter>]
    Then [outcome with <parameter>]

    Examples:
      | parameter | expected |
      | value1    | result1  |
      | value2    | result2  |
```

---

## Real-World Examples

### Example 1: Position Sizing

```gherkin
Feature: Position Sizing Calculation
  As a trader
  I want to calculate position size based on ATR
  So that every trade risks the same dollar amount

  Background:
    Given I have a $10,000 account
    And I risk 0.75% per trade ($75)

  Scenario: Calculate stock position
    Given I want to trade AAPL
    And the entry price is $180
    And the 20-day ATR is 1.50
    And I use a 2-ATR stop
    When I calculate position size
    Then my stop distance should be $3.00
    And my initial stop should be $177.00
    And I should buy 25 shares
    And my total risk should be $75.00

  Scenario: Handle zero ATR
    Given the ATR is 0
    When I calculate position size
    Then I should see error "ATR must be greater than zero"
    And no position size should be calculated

  Scenario: Handle negative entry price
    Given the entry price is -180
    When I calculate position size
    Then I should see error "Entry price must be positive"
    And no position size should be calculated

  Scenario Outline: Position sizing across price ranges
    Given the entry price is <entry>
    And the ATR is <atr>
    And I use a 2-ATR stop
    When I calculate position size
    Then I should buy <shares> shares
    And my risk should be approximately $75

    Examples:
      | entry | atr  | shares |
      | 10    | 0.25 | 150    |
      | 50    | 0.50 | 75     |
      | 100   | 1.00 | 37     |
      | 180   | 1.50 | 25     |
      | 500   | 5.00 | 7      |
```

### Example 2: Heat Management

```gherkin
Feature: Portfolio Heat Management
  As a trader
  I want to limit total portfolio risk exposure
  So that I don't blow up my account in a drawdown

  Background:
    Given I have a $10,000 account
    And my portfolio heat cap is 4% ($400)
    And my bucket heat cap is 1.5% ($150)

  Scenario: First trade of the day
    Given I have no open positions
    When I check heat for a $75 trade in "Tech/Comm" bucket
    Then portfolio heat should be $75
    And portfolio heat percentage should be 18.75% of cap
    And bucket heat should be $75
    And bucket heat percentage should be 50% of cap
    And the trade should be allowed

  Scenario: Portfolio heat approaching cap
    Given I have 5 open positions
    And total portfolio heat is $350
    When I check heat for a $75 trade
    Then portfolio heat would be $425
    And portfolio heat percentage would be 106.25% of cap
    And I should see error "Portfolio heat exceeds cap by $25"
    And the trade should be rejected

  Scenario: Bucket heat at limit
    Given I have 2 open positions in "Tech/Comm" bucket
    And "Tech/Comm" bucket heat is $130
    When I check heat for a $75 trade in "Tech/Comm"
    Then bucket heat would be $205
    And bucket heat percentage would be 136.67% of cap
    And I should see error "Bucket heat exceeds cap"
    And the trade should be rejected

  Scenario: Portfolio OK, bucket at risk
    Given portfolio heat is $200 (50% of cap)
    And "Healthcare" bucket heat is $145
    When I check heat for a $75 trade in "Healthcare"
    Then portfolio heat check should pass
    But bucket heat check should fail
    And I should see warning "Bucket would exceed cap"
    And the trade should be rejected
```

### Example 3: Checklist Validation

```gherkin
Feature: Entry Checklist Validation
  As a trader
  I want to validate all entry criteria
  So that I only take high-probability setups

  Scenario: Perfect setup - All checks pass
    Given I check all 6 checklist items:
      | FromPreset    | true |
      | TrendPass     | true |
      | LiquidityPass | true |
      | TVConfirm     | true |
      | EarningsOK    | true |
      | JournalOK     | true |
    When I evaluate the checklist
    Then the banner should be "GREEN - GO"
    And missing items should be 0
    And the impulse timer should start

  Scenario: One check missing - Caution
    Given I check 5 of 6 checklist items
    And "EarningsOK" is not checked
    When I evaluate the checklist
    Then the banner should be "YELLOW - CAUTION"
    And missing items should be 1
    And the missing list should show "EarningsOK"
    And the impulse timer should not start

  Scenario: Multiple checks missing - No-Go
    Given I check only 3 of 6 checklist items
    And "FromPreset" is not checked
    And "TrendPass" is not checked
    And "EarningsOK" is not checked
    When I evaluate the checklist
    Then the banner should be "RED - NO-GO"
    And missing items should be 3
    And the missing list should show "FromPreset, TrendPass, EarningsOK"
    And the impulse timer should not start

  Scenario: Attempting to save without GREEN banner
    Given the banner is "YELLOW - CAUTION"
    When I try to save the trade decision
    Then I should see error "Cannot save: Banner must be GREEN"
    And the decision should not be saved to database
```

---

## From Gherkin to Code

### Step 1: Write Gherkin
```gherkin
Scenario: Calculate shares
  Given I risk $75
  And stop distance is $3
  When I calculate shares
  Then I should get 25 shares
```

### Step 2: Implement Step Definitions (Go example)

```go
// features/steps/sizing_steps.go
package steps

import (
    "github.com/cucumber/godog"
    "trading-backend/internal/domain"
)

type sizingContext struct {
    risk         float64
    stopDistance float64
    result       domain.SizingResult
}

func (ctx *sizingContext) iRisk(amount float64) error {
    ctx.risk = amount
    return nil
}

func (ctx *sizingContext) stopDistanceIs(distance float64) error {
    ctx.stopDistance = distance
    return nil
}

func (ctx *sizingContext) iCalculateShares() error {
    calculator := domain.NewSizingCalculator()
    result, err := calculator.CalculateShares(ctx.risk, ctx.stopDistance)
    if err != nil {
        return err
    }
    ctx.result = result
    return nil
}

func (ctx *sizingContext) iShouldGetShares(expected int) error {
    if ctx.result.Shares != expected {
        return fmt.Errorf("expected %d shares, got %d", expected, ctx.result.Shares)
    }
    return nil
}

func InitializeSizingSteps(ctx *godog.ScenarioContext) {
    sizing := &sizingContext{}

    ctx.Step(`^I risk \$(\d+)$`, sizing.iRisk)
    ctx.Step(`^stop distance is \$(\d+)$`, sizing.stopDistanceIs)
    ctx.Step(`^I calculate shares$`, sizing.iCalculateShares)
    ctx.Step(`^I should get (\d+) shares$`, sizing.iShouldGetShares)
}
```

### Step 3: Implement Business Logic

```go
// internal/domain/sizing.go
package domain

type SizingCalculator struct{}

func NewSizingCalculator() *SizingCalculator {
    return &SizingCalculator{}
}

func (c *SizingCalculator) CalculateShares(risk float64, stopDistance float64) (SizingResult, error) {
    if risk <= 0 {
        return SizingResult{}, errors.New("risk must be positive")
    }
    if stopDistance <= 0 {
        return SizingResult{}, errors.New("stop distance must be positive")
    }

    shares := int(math.Floor(risk / stopDistance))

    return SizingResult{
        Shares:       shares,
        RiskDollars:  risk,
        StopDistance: stopDistance,
    }, nil
}
```

### Step 4: Run Tests

```bash
$ go test ./features
Feature: Position Sizing
  Scenario: Calculate shares                    # features/sizing.feature:5
    Given I risk $75                            # passed
    And stop distance is $3                     # passed
    When I calculate shares                     # passed
    Then I should get 25 shares                 # passed

1 scenario (1 passed)
4 steps (4 passed)
```

✅ **Gherkin → Code → Passing Tests**

---

## Common Mistakes

### ❌ Mistake: UI-focused scenarios

```gherkin
# Bad: Testing Excel UI details
Scenario: Update cell
  Given I'm on the TradeEntry sheet
  When I click cell B9
  And I type "180"
  And I press Enter
  Then cell B9 should show "180"
```

**Why bad:** Tests UI implementation, not business behavior.

**Fix:** Test the business action:
```gherkin
Scenario: Enter trade price
  Given I'm preparing a trade
  When I set entry price to $180
  Then the entry price should be recorded as $180
```

### ❌ Mistake: Too granular

```gherkin
# Bad: Testing individual lines of code
Scenario: Multiply equity by risk percent
  Given equity is 10000
  And riskPct is 0.0075
  When I multiply equity by riskPct
  Then result is 75
```

**Why bad:** This is a unit test, not a behavior test.

**Fix:** Test the business outcome:
```gherkin
Scenario: Calculate risk amount
  Given I have a $10,000 account
  And I risk 0.75% per trade
  When I calculate my risk per trade
  Then I should risk $75 per trade
```

### ❌ Mistake: Coupling to implementation

```gherkin
# Bad: References code structure
Scenario: Call sizing service
  Given SizingService is initialized
  When CalculateShares method is called
  Then result object is returned
```

**Why bad:** Changes if you rename classes/methods.

**Fix:** Describe behavior:
```gherkin
Scenario: Calculate position size
  Given I have trade parameters
  When I calculate how many shares to buy
  Then I get a position size
```

---

## Integration with Development

### Workflow

```
1. Write Gherkin scenario (failing)
   ↓
2. Write step definitions (failing)
   ↓
3. Implement business logic (passing)
   ↓
4. Refactor (still passing)
   ↓
5. Commit working feature
```

### Red → Green → Refactor

Just like TDD, but with human-readable tests.

---

## Tools

### Go: godog
```bash
go get github.com/cucumber/godog
godog run features/
```

### JavaScript: cucumber-js
```bash
npm install @cucumber/cucumber
npx cucumber-js
```

### Python: behave
```bash
pip install behave
behave features/
```

### Ruby: cucumber
```bash
gem install cucumber
cucumber features/
```

---

## Summary

### BDD in One Sentence
**Write how the system should behave in human language, then make the code match.**

### Why It Matters Here
This trading system's value is in its behavior (preventing bad trades), not its code. BDD tests behavior.

### The Golden Rule
**If you can't describe it in Gherkin, you don't understand the requirement.**

Write the Gherkin first. Always.
