Feature: Impulse brake timer
  As a trader
  I want a mandatory 2-minute delay after evaluation
  So that I don't make impulsive trading decisions

  Background:
    Given I have initialized the database
    And I have a candidate "AAPL" for today

  Scenario: Start timer when banner is GREEN
    Given all checklist items are satisfied
    When I run "checklist --ticker AAPL" with all checks true
    Then banner should be "GREEN"
    And impulse timer should start
    And I should see message "Wait 2 minutes before saving decision"

  Scenario: Do not start timer when banner is YELLOW
    Given 1 checklist item is missing
    When I run "checklist --ticker AAPL" with 1 check false
    Then banner should be "YELLOW"
    And impulse timer should NOT start

  Scenario: Do not start timer when banner is RED
    Given 2+ checklist items are missing
    When I run "checklist --ticker AAPL" with 2+ checks false
    Then banner should be "RED"
    And impulse timer should NOT start

  Scenario: Reject save before timer expires
    Given timer started 60 seconds ago
    When I try to check impulse brake for "AAPL"
    Then I should receive error "impulse brake active, wait 60 more seconds"
    And the decision should NOT be allowed

  Scenario: Allow save after timer expires
    Given timer started 120 seconds ago
    When I try to check impulse brake for "AAPL"
    Then the brake check should pass
    And I should NOT see timer error

  Scenario: Check remaining time (in progress)
    Given timer started 90 seconds ago
    When I run "check-timer --ticker AAPL"
    Then I should see "30 seconds" remaining
    And timer status should be "active"

  Scenario: Check timer when expired
    Given timer started 180 seconds ago
    When I run "check-timer --ticker AAPL"
    Then I should see "Timer expired"
    And I should see message "you may proceed"

  Scenario: Check timer when no timer exists
    Given I have not evaluated "AAPL" yet
    When I run "check-timer --ticker AAPL"
    Then I should see "No active timer for AAPL"
    And no error should occur

  Scenario: Timer replaces previous timer for same ticker
    Given I evaluated "AAPL" 5 minutes ago with GREEN
    When I evaluate "AAPL" again with GREEN
    Then a new timer should start
    And the old timer should be deactivated
    And remaining time should be 120 seconds

  Scenario: Multiple tickers have independent timers
    Given I evaluated "AAPL" 60 seconds ago with GREEN
    And I evaluated "MSFT" 30 seconds ago with GREEN
    When I check timer for "AAPL"
    Then I should see 60 seconds remaining
    When I check timer for "MSFT"
    Then I should see 90 seconds remaining

  Scenario: Timer precision (fractional seconds)
    Given timer started 119.5 seconds ago
    When I check impulse brake
    Then I should receive error "impulse brake active"
    And error should show "1 second" remaining

  Scenario: Timer output includes timestamps
    Given timer started 45 seconds ago
    When I run "check-timer --ticker AAPL"
    Then I should see the start time
    And I should see the expiry time
    And I should see "75 seconds" remaining
