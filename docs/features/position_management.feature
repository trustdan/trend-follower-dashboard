Feature: Position management
  As a trader
  I want to manage open positions
  So that I can track my trades through their lifecycle

  Background:
    Given I have initialized the database

  Scenario: Open position from GO decision
    Given I have a GO decision for "AAPL" with entry=$180, stop=$177, shares=25
    When I run "open-position --ticker AAPL"
    Then position should be created with status OPEN
    And position should have entry=$180, stop=$177, shares=25, risk=$75

  Scenario: Cannot open position without GO decision
    Given I have NO decision for "XYZ"
    When I run "open-position --ticker XYZ"
    Then I should receive error "no decision found"

  Scenario: Cannot open position for NO-GO decision
    Given I have a NO-GO decision for "MSFT"
    When I run "open-position --ticker MSFT"
    Then I should receive error "cannot open position for NO-GO decision"

  Scenario: Update position stop (trailing up)
    Given I have an open LONG position in "AAPL" with entry=$180, stop=$177
    When I run "update-stop --ticker AAPL --new-stop 179"
    Then position stop should be $179
    And risk should be recalculated to $25 (25 shares * $1 stop distance)

  Scenario: Cannot move stop against position (down for long)
    Given I have a LONG position in "AAPL" with stop=$177
    When I run "update-stop --ticker AAPL --new-stop 175"
    Then I should receive error "cannot move stop lower for LONG position"

  Scenario: Close position with WIN
    Given I have an open position in "AAPL" at entry=$180, shares=25
    When I run "close-position --ticker AAPL --exit 185 --outcome WIN"
    Then position should be CLOSED
    And P&L should be $125 (= 25 * ($185 - $180))
    And outcome should be WIN
    And bucket should NOT enter cooldown

  Scenario: Close position with LOSS triggers bucket cooldown
    Given I have an open position in "AAPL" at entry=$180, shares=25, bucket="Tech/Comm"
    When I run "close-position --ticker AAPL --exit 176 --outcome LOSS"
    Then position should be CLOSED
    And P&L should be -$100
    And outcome should be LOSS
    And bucket "Tech/Comm" should enter 24-hour cooldown

  Scenario: Close position with SCRATCH
    Given I have an open position in "AAPL" at entry=$180, shares=25
    When I run "close-position --ticker AAPL --exit 180 --outcome SCRATCH"
    Then position should be CLOSED
    And P&L should be $0
    And outcome should be SCRATCH
    And bucket should NOT enter cooldown

  Scenario: List all open positions
    Given I have 3 open positions
    When I run "list-positions"
    Then I should see 3 positions
    And each should show ticker, entry, shares, risk

  Scenario: List positions by status filter
    Given I have 2 open positions and 3 closed positions
    When I run "list-positions --status OPEN"
    Then I should see only 2 open positions

  Scenario: Get position details for specific ticker
    Given I have an open position in "AAPL" with entry=$180, stop=$177
    When I run "get-position --ticker AAPL"
    Then I should see full position details

  Scenario: Portfolio heat includes open position risk
    Given I have open positions with total risk of $225
    When I check portfolio heat
    Then current heat should be $225

  Scenario: Cannot open duplicate position for same ticker
    Given I have an open position in "AAPL"
    When I run "open-position --ticker AAPL"
    Then I should receive error "already have open position"

  Scenario: P&L calculation for stocks
    Given I have open position: entry=$100, shares=50
    When I close at exit=$110
    Then P&L should be $500 (= 50 * ($110 - $100))

  Scenario: Stop update increases risk when moved down
    Given I have position with entry=$180, stop=$177, shares=25, risk=$75
    When I update stop to $175
    Then new risk should be $125 (= 25 * ($180 - $175))
