Feature: Position sizing for options using Delta-ATR method
  As an options trader
  I want to calculate contracts based on delta and ATR
  So that my options position risks the same as an equivalent stock position

  Background:
    Given I have initialized the database with default settings
    And Equity_E is 10000
    And RiskPct_r is 0.0075
    And StopMultiple_K is 2

  Scenario: Basic delta-ATR options sizing (low delta yields 0 contracts)
    Given I want to trade options using delta-ATR method
    And the underlying entry price is 180
    And the ATR is 1.5
    And I want 30-delta options
    When I calculate position size
    Then risk_dollars should be 75.00
    And stop_distance should be 3.00
    And contracts should be 0
    And shares should be 0

  Scenario: Higher delta yields more contracts
    Given I want to trade options using delta-ATR method
    And the underlying entry price is 50
    And the ATR is 1.0
    And I want 70-delta options
    When I calculate position size
    Then risk_dollars should be 75.00
    And stop_distance should be 2.00
    And contracts should be 0
    And shares should be 0

  Scenario: 80-delta deep ITM options
    Given I want to trade options using delta-ATR method
    And the underlying entry price is 100
    And the ATR is 2.0
    And I want 80-delta options
    When I calculate position size
    Then risk_dollars should be 75.00
    And stop_distance should be 4.00
    And contracts should be 2
    And actual_risk should be approximately 64.00

  Scenario: 90-delta very deep ITM options
    Given I want to trade options using delta-ATR method
    And the underlying entry price is 100
    And the ATR is 1.0
    And I want 90-delta options
    When I calculate position size
    Then risk_dollars should be 75.00
    And stop_distance should be 2.00
    And contracts should be 3
    And actual_risk should be approximately 54.00

  Scenario: Reject invalid delta (< 0)
    Given I want to trade options using delta-ATR method
    And the underlying entry price is 180
    And the ATR is 1.5
    And I specify delta as -0.3
    When I calculate position size
    Then I should receive an error "delta must be between 0 and 1"

  Scenario: Reject invalid delta (> 1)
    Given I want to trade options using delta-ATR method
    And the underlying entry price is 180
    And the ATR is 1.5
    And I specify delta as 1.5
    When I calculate position size
    Then I should receive an error "delta must be between 0 and 1"

  Scenario: Reject zero delta
    Given I want to trade options using delta-ATR method
    And the underlying entry price is 180
    And the ATR is 1.5
    And I specify delta as 0
    When I calculate position size
    Then I should receive an error "delta must be between 0 and 1"

  Scenario: Reject missing delta
    Given I want to trade options using delta-ATR method
    And the underlying entry price is 180
    And the ATR is 1.5
    And I do not specify delta
    When I calculate position size
    Then I should receive an error "delta must be between 0 and 1"
