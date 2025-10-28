Feature: Bucket cooldown after losses
  As a trader
  I want buckets to enter cooldown after losses
  So that I don't chase losses in the same sector

  Background:
    Given I have initialized the database

  Scenario: Block new trades in bucket during cooldown
    Given bucket "Tech/Comm" is in cooldown
    And "MSFT" is in bucket "Tech/Comm"
    And all other gates pass for "MSFT"
    When I try to save decision for "MSFT"
    Then I should receive error "bucket Tech/Comm is in cooldown"

  Scenario: Allow trades in different bucket during cooldown
    Given bucket "Tech/Comm" is in cooldown
    And "JNJ" is in bucket "Healthcare"
    And all other gates pass for "JNJ"
    When I try to save decision for "JNJ"
    Then bucket cooldown should NOT block the save
    And the decision should be saved successfully

  Scenario: Check active cooldown status
    Given bucket "Tech/Comm" entered cooldown 6 hours ago
    When I run "check-cooldown --bucket Tech/Comm"
    Then I should see cooldown is active
    And I should see approximately 18 hours remaining

  Scenario: Check expired cooldown
    Given bucket "Tech/Comm" entered cooldown 25 hours ago
    When I run "check-cooldown --bucket Tech/Comm"
    Then I should see cooldown is NOT active
    And I should see "Cooldown expired" or no active cooldown

  Scenario: Check bucket with no cooldown
    Given bucket "Energy" has never had a cooldown
    When I run "check-cooldown --bucket Energy"
    Then I should see "NOT in cooldown"

  Scenario: List all active cooldowns
    Given bucket "Tech/Comm" is in cooldown
    And bucket "Energy" is in cooldown
    And bucket "Finance" is in cooldown
    When I run "list-cooldowns"
    Then I should see 3 active cooldowns
    And each should show bucket name and remaining time

  Scenario: Trigger cooldown manually
    Given bucket "Tech/Comm" is NOT in cooldown
    When I trigger cooldown for "Tech/Comm" with reason "Manual test"
    Then bucket "Tech/Comm" should enter cooldown
    And cooldown duration should be 24 hours

  Scenario: Extend existing cooldown
    Given bucket "Tech/Comm" is in cooldown with 10 hours remaining
    When I trigger cooldown for "Tech/Comm" again
    Then cooldown should be extended to 24 hours from now

  Scenario: Save decision fails when bucket in cooldown
    Given I have set up all prerequisites
    And bucket "Tech/Comm" is in cooldown
    When I run "save-decision --ticker AAPL --entry 180 --atr 1.5 --bucket Tech/Comm --action GO"
    Then I should receive error "bucket Tech/Comm is in cooldown"
    And the decision should NOT be saved

  Scenario: Save decision succeeds when bucket not in cooldown
    Given I have set up all prerequisites
    And bucket "Tech/Comm" is NOT in cooldown
    When I run "save-decision --ticker AAPL --entry 180 --atr 1.5 --bucket Tech/Comm --action GO"
    Then all 5 gates should pass
    And the decision should be saved successfully

  Scenario: Empty bucket parameter bypasses cooldown check
    Given bucket "Tech/Comm" is in cooldown
    When I run "save-decision --ticker AAPL --entry 180 --atr 1.5 --action GO"
    Then bucket cooldown gate should be skipped
    And cooldown should NOT block the save
