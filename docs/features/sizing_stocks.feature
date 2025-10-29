Feature: Position sizing for stock trades
  As a trader
  I want to calculate position size based on ATR
  So that every trade risks exactly 0.75% of my account

  Background:
    Given I have initialized the database with default settings
    And Equity_E is 10000
    And RiskPct_r is 0.0075
    And StopMultiple_K is 2

  Scenario: Basic stock position sizing
    Given I want to trade a stock
    And the entry price is 180
    And the ATR is 1.5
    When I calculate position size
    Then risk_dollars should be 75.00
    And stop_distance should be 3.00
    And initial_stop should be 177.00
    And shares should be 25
    And contracts should be 0
    And actual_risk should be 75.00

  Scenario Outline: Position sizing across price ranges
    Given I want to trade a stock
    And the entry price is <entry>
    And the ATR is <atr>
    When I calculate position size
    Then shares should be <shares>
    And risk_dollars should be approximately 75.00

    Examples:
      | entry | atr  | shares |
      | 10    | 0.25 | 150    |
      | 50    | 0.50 | 75     |
      | 100   | 1.00 | 37     |
      | 180   | 1.50 | 25     |
      | 500   | 5.00 | 7      |

  Scenario: Reject zero ATR
    Given I want to trade a stock
    And the entry price is 180
    And the ATR is 0
    When I calculate position size
    Then I should receive an error "ATR must be greater than zero"

  Scenario: Reject negative entry price
    Given I want to trade a stock
    And the entry price is -180
    And the ATR is 1.5
    When I calculate position size
    Then I should receive an error "entry price must be positive"

  Scenario: Reject negative ATR
    Given I want to trade a stock
    And the entry price is 180
    And the ATR is -1.5
    When I calculate position size
    Then I should receive an error "ATR must be greater than zero"

  Scenario: Reject zero K multiple
    Given I want to trade a stock
    And the entry price is 180
    And the ATR is 1.5
    And StopMultiple_K is 0
    When I calculate position size
    Then I should receive an error "K multiple must be positive"

  Scenario: Reject zero equity
    Given Equity_E is 0
    And I want to trade a stock
    And the entry price is 180
    And the ATR is 1.5
    When I calculate position size
    Then I should receive an error "equity must be greater than zero"
