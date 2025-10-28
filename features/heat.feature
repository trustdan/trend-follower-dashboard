Feature: Portfolio and bucket heat management
  As a trader
  I want to check my portfolio and bucket heat levels
  So that I don't overexpose my account to correlated risks

  Background:
    Given I have initialized the database with default settings
    And Equity_E is 10000
    And HeatCap_H_pct is 0.04 (4% = $400)
    And BucketHeatCap_pct is 0.015 (1.5% = $150)

  Scenario: First trade of the day - no existing positions
    Given I have no open positions
    When I check heat for a $75 trade in "Tech/Comm" bucket
    Then current_portfolio_heat should be 0
    And new_portfolio_heat should be 75
    And portfolio_cap_exceeded should be false
    And current_bucket_heat should be 0
    And new_bucket_heat should be 75
    And bucket_cap_exceeded should be false
    And allowed should be true

  Scenario: Portfolio heat approaching cap
    Given I have 4 open positions with $75 risk each
    And total portfolio heat is $300
    When I check heat for a $75 trade in "Tech/Comm" bucket
    Then current_portfolio_heat should be 300
    And new_portfolio_heat should be 375
    And portfolio_cap_exceeded should be false
    And allowed should be true

  Scenario: Portfolio heat exceeds cap
    Given I have 5 open positions with $75 risk each
    And total portfolio heat is $375
    When I check heat for a $75 trade in "Tech/Comm" bucket
    Then current_portfolio_heat should be 375
    And new_portfolio_heat should be 450
    And portfolio_cap_exceeded should be true
    And portfolio_overage should be 50
    And allowed should be false
    And rejection_reason should contain "Portfolio heat"

  Scenario: Bucket heat approaching cap
    Given I have 1 open position in "Tech/Comm" with $100 risk
    And "Tech/Comm" bucket heat is $100
    When I check heat for a $49 trade in "Tech/Comm" bucket
    Then current_bucket_heat should be 100
    And new_bucket_heat should be 149
    And bucket_cap_exceeded should be false
    And allowed should be true

  Scenario: Bucket heat exceeds cap
    Given I have 1 open position in "Tech/Comm" with $100 risk
    And "Tech/Comm" bucket heat is $100
    When I check heat for a $75 trade in "Tech/Comm" bucket
    Then current_bucket_heat should be 100
    And new_bucket_heat should be 175
    And bucket_cap_exceeded should be true
    And bucket_overage should be 25
    And allowed should be false
    And rejection_reason should contain "Bucket"

  Scenario: Portfolio OK but bucket exceeds
    Given portfolio heat is $200 (50% of cap)
    And "Healthcare" bucket heat is $130
    When I check heat for a $75 trade in "Healthcare" bucket
    Then portfolio_cap_exceeded should be false
    And bucket_cap_exceeded should be true
    And allowed should be false
    And rejection_reason should contain "Bucket"

  Scenario: Different bucket - no bucket conflict
    Given "Tech/Comm" bucket heat is $140
    And portfolio heat is $200
    When I check heat for a $75 trade in "Energy" bucket
    Then bucket_cap_exceeded should be false
    And current_bucket_heat should be 0
    And portfolio_cap_exceeded should be false
    And allowed should be true

  Scenario: Multiple positions across buckets
    Given I have positions across 3 buckets:
      | Bucket      | Heat |
      | Tech/Comm   | 100  |
      | Healthcare  | 75   |
      | Energy      | 50   |
    And total portfolio heat is $225
    When I check heat for a $75 trade in "Finance" bucket
    Then current_portfolio_heat should be 225
    And new_portfolio_heat should be 300
    And current_bucket_heat should be 0
    And allowed should be true

  Scenario: Zero risk trade (query only)
    Given portfolio heat is $300
    When I check heat for a $0 trade in "Tech/Comm" bucket
    Then new_portfolio_heat should be 300
    And allowed should be true

  Scenario: Reject negative risk
    When I check heat for a $-75 trade in "Tech/Comm" bucket
    Then I should receive an error "add_risk_dollars must be non-negative"
