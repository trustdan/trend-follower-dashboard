Feature: Database initialization and migrations
  As a trading system
  I need persistent storage for settings, candidates, decisions, and positions
  So that state is preserved across sessions

  Scenario: Initialize database on first run
    Given no database file exists
    When I run "tf-engine init"
    Then a database file "trading.db" should be created
    And the database should contain table "settings"
    And the database should contain table "presets"
    And the database should contain table "candidates"
    And the database should contain table "decisions"
    And the database should contain table "positions"

  Scenario: Bootstrap default settings
    Given I have initialized the database
    When I query the settings table
    Then setting "Equity_E" should be "10000"
    And setting "RiskPct_r" should be "0.0075"
    And setting "HeatCap_H_pct" should be "0.04"
    And setting "BucketHeatCap_pct" should be "0.015"
    And setting "StopMultiple_K" should be "2"

  Scenario: Idempotent initialization
    Given I have already initialized the database
    When I run "tf-engine init" again
    Then the database should not be recreated
    And existing settings should be preserved
    And the command should succeed
