Feature: Position sizing for options using Max-Loss method
  As an options trader
  I want to calculate contracts based on maximum loss per contract
  So that my total position risk stays within my account risk limits

  Background:
    Given I have initialized the database with default settings
    And Equity_E is 10000
    And RiskPct_r is 0.0075

  Scenario: Basic max-loss options sizing
    Given I want to trade options using max-loss method
    And my maximum loss per contract is 50
    When I calculate position size
    Then risk_dollars should be 75.00
    And contracts should be 1
    And actual_risk should be 50.00
    And shares should be 0

  Scenario: Higher max-loss = fewer contracts
    Given I want to trade options using max-loss method
    And my maximum loss per contract is 100
    When I calculate position size
    Then risk_dollars should be 75.00
    And contracts should be 0
    And actual_risk should be 0.00

  Scenario: Lower max-loss = more contracts
    Given I want to trade options using max-loss method
    And my maximum loss per contract is 25
    When I calculate position size
    Then risk_dollars should be 75.00
    And contracts should be 3
    And actual_risk should be 75.00

  Scenario: Fractional contracts round down
    Given I want to trade options using max-loss method
    And my maximum loss per contract is 30
    When I calculate position size
    Then risk_dollars should be 75.00
    And contracts should be 2
    And actual_risk should be 60.00

  Scenario: Very small max-loss
    Given I want to trade options using max-loss method
    And my maximum loss per contract is 10
    When I calculate position size
    Then risk_dollars should be 75.00
    And contracts should be 7
    And actual_risk should be 70.00

  Scenario: Exact divisor
    Given I want to trade options using max-loss method
    And my maximum loss per contract is 15
    When I calculate position size
    Then risk_dollars should be 75.00
    And contracts should be 5
    And actual_risk should be 75.00

  Scenario: Reject zero max-loss
    Given I want to trade options using max-loss method
    And my maximum loss per contract is 0
    When I calculate position size
    Then I should receive an error "max loss must be positive"

  Scenario: Reject negative max-loss
    Given I want to trade options using max-loss method
    And my maximum loss per contract is -50
    When I calculate position size
    Then I should receive an error "max loss must be positive"
