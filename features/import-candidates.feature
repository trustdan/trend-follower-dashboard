Feature: Import candidates
  As a trader
  I want to import my daily candidate tickers
  So that I can ensure only screened stocks are tradeable

  Background:
    Given I have initialized the database with default settings

  Scenario: Import candidates with preset
    When I run "import-candidates --tickers AAPL,MSFT,NVDA --preset TF_BREAKOUT_LONG"
    Then I should see 3 candidates imported
    And candidates should be stored for today's date
    And the preset should be "TF_BREAKOUT_LONG"

  Scenario: Import candidates without preset
    When I run "import-candidates --tickers TSLA,GOOGL"
    Then I should see 2 candidates imported
    And candidates should be stored for today's date
    And the preset should be empty

  Scenario: Import with sector and bucket information
    When I run "import-candidates --tickers AAPL --preset TF_BREAKOUT_LONG --sector Technology --bucket Tech/Comm"
    Then I should see 1 candidate imported
    And candidate "AAPL" should have sector "Technology"
    And candidate "AAPL" should have bucket "Tech/Comm"

  Scenario: Import candidates multiple times same day (idempotent)
    Given I have imported "AAPL,MSFT" with preset "TF_BREAKOUT_LONG" today
    When I run "import-candidates --tickers AAPL,NVDA --preset TF_BREAKOUT_LONG"
    Then I should see 2 candidates total for today
    And candidate list should be "AAPL,NVDA"

  Scenario: Import candidates replaces previous day's list
    Given I have imported "AAPL,MSFT" for yesterday
    When I run "import-candidates --tickers TSLA,GOOGL"
    Then I should see 2 candidates for today
    And yesterday's candidates should still exist
    And today's candidates should be "TSLA,GOOGL"

  Scenario: Empty ticker list rejected
    When I run "import-candidates --tickers ''"
    Then I should receive an error "at least one ticker required"

  Scenario: Tickers are normalized
    When I run "import-candidates --tickers 'aapl, msft , NVDA'"
    Then candidate list should be "AAPL,MSFT,NVDA"
    And all tickers should be uppercase
    And whitespace should be removed

  Scenario: Check if ticker is in today's candidates
    Given I have imported "AAPL,MSFT,NVDA" today
    When I check if "AAPL" is in today's candidates
    Then the result should be true
    When I check if "TSLA" is in today's candidates
    Then the result should be false

  Scenario: Get today's candidates count
    Given I have imported "AAPL,MSFT,NVDA,GOOGL,TSLA" today
    When I get today's candidates count
    Then the count should be 5

  Scenario: List all candidates for today
    Given I have imported "AAPL,MSFT,NVDA" with preset "TF_BREAKOUT_LONG" today
    When I get today's candidates list
    Then I should see 3 candidates
    And each candidate should have date, ticker, and preset information
