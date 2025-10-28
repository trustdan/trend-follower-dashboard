Feature: Settings management
  As a trader
  I want to view and update system settings
  So that I can configure my account parameters

  Background:
    Given I have initialized the database with default settings

  Scenario: Get all settings
    When I run "get-settings"
    Then I should see setting "Equity_E" with value "10000"
    And I should see setting "RiskPct_r" with value "0.0075"
    And I should see setting "HeatCap_H_pct" with value "0.04"
    And I should see setting "BucketHeatCap_pct" with value "0.015"
    And I should see setting "StopMultiple_K" with value "2"

  Scenario: Get single setting
    When I run "get-settings --key Equity_E"
    Then I should see value "10000"

  Scenario: Update equity setting
    Given Equity_E is currently "10000"
    When I run "set-setting --key Equity_E --value 20000"
    Then setting "Equity_E" should be "20000"
    And the change should be persisted

  Scenario: Update risk percent
    Given RiskPct_r is currently "0.0075"
    When I run "set-setting --key RiskPct_r --value 0.01"
    Then setting "RiskPct_r" should be "0.01"
    And subsequent calculations should use the new value

  Scenario: Reject invalid setting key
    When I run "set-setting --key InvalidKey --value 100"
    Then I should receive an error "setting not found: InvalidKey"

  Scenario: Validate equity is positive
    When I run "set-setting --key Equity_E --value -1000"
    Then I should receive an error "Equity_E must be positive"

  Scenario: Validate risk percent range
    When I run "set-setting --key RiskPct_r --value 1.5"
    Then I should receive an error "RiskPct_r must be between 0 and 1"

  Scenario: Validate heat cap range
    When I run "set-setting --key HeatCap_H_pct --value 1.5"
    Then I should receive an error "HeatCap_H_pct must be between 0 and 1"

  Scenario: Validate bucket heat cap range
    When I run "set-setting --key BucketHeatCap_pct --value 1.5"
    Then I should receive an error "BucketHeatCap_pct must be between 0 and 1"

  Scenario: Validate K multiple is positive
    When I run "set-setting --key StopMultiple_K --value 0"
    Then I should receive an error "StopMultiple_K must be positive"

  Scenario: Reject non-numeric value
    When I run "set-setting --key Equity_E --value abc"
    Then I should receive an error "value must be a number"
