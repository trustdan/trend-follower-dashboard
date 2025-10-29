Feature: Save trading decision
  As a trader
  I want to save my trade decision
  So that only valid, disciplined trades are recorded

  Background:
    Given I have initialized the database with default settings
    And I have a candidate "AAPL" for today in bucket "Tech/Comm"

  # GATE 1: Banner GREEN
  Scenario: Reject save when banner is YELLOW
    Given banner is YELLOW for "AAPL"
    When I run "save-decision --ticker AAPL --entry 180 --atr 1.5 --action GO"
    Then I should receive error "Banner must be GREEN"
    And the decision should NOT be saved

  Scenario: Reject save when banner is RED
    Given banner is RED for "AAPL"
    When I run "save-decision --ticker AAPL --entry 180 --atr 1.5 --action GO"
    Then I should receive error "Banner must be GREEN"
    And the decision should NOT be saved

  # GATE 2: Ticker in candidates
  Scenario: Reject save for ticker not in today's candidates
    Given "XYZ" is NOT in today's candidates
    When I run "save-decision --ticker XYZ --entry 50 --atr 1.0 --action GO"
    Then I should receive error "not in today's candidates"
    And the decision should NOT be saved

  # GATE 3: Impulse brake
  Scenario: Reject save when timer not expired
    Given banner is GREEN for "AAPL"
    And impulse timer has 30 seconds remaining for "AAPL"
    When I run "save-decision --ticker AAPL --entry 180 --atr 1.5 --action GO"
    Then I should receive error "impulse brake active"
    And the decision should NOT be saved

  # GATE 4: Bucket cooldown (placeholder - M14)
  Scenario: Reject save when bucket in cooldown
    Given bucket "Tech/Comm" is in cooldown
    When I run "save-decision --ticker AAPL --entry 180 --atr 1.5 --action GO"
    Then I should receive error "bucket.*cooldown"
    And the decision should NOT be saved

  # GATE 5: Heat caps
  Scenario: Reject when portfolio heat exceeds cap
    Given portfolio heat is $350
    And portfolio cap is $400
    When I run "save-decision --ticker AAPL --entry 180 --atr 1.5 --action GO"
    Then position risk should be $75
    And projected portfolio heat should be $425
    And I should receive error "portfolio heat.*exceeds cap"
    And the decision should NOT be saved

  Scenario: Reject when bucket heat exceeds cap
    Given bucket "Tech/Comm" heat is $125
    And bucket cap is $150
    When I run "save-decision --ticker AAPL --entry 180 --atr 1.5 --action GO"
    Then position risk should be $75
    And projected bucket heat should be $200
    And I should receive error "bucket heat.*exceeds cap"
    And the decision should NOT be saved

  # Success path
  Scenario: Save GO decision passing all gates
    Given banner is GREEN for "AAPL"
    And impulse timer expired for "AAPL"
    And bucket "Tech/Comm" is NOT in cooldown
    And portfolio heat is below cap
    And bucket heat is below cap
    When I run "save-decision --ticker AAPL --entry 180 --atr 1.5 --action GO"
    Then the decision should be saved with action "GO"
    And position size should be 25 shares
    And initial stop should be $177.00
    And risk should be $75.00
    And I should see confirmation "Decision saved: AAPL GO"

  Scenario: Save NO-GO decision (gates not checked)
    When I run "save-decision --ticker AAPL --action NO-GO --reason 'Bad setup'"
    Then the decision should be saved with action "NO-GO"
    And position size should be 0
    And risk should be $0
    And gates should NOT be validated
    And I should see confirmation "Decision saved: AAPL NO-GO"

  Scenario: Save decision includes all calculated values
    Given banner is GREEN for "AAPL"
    And impulse timer expired for "AAPL"
    And all heat caps pass
    When I run "save-decision --ticker AAPL --entry 180 --atr 1.5 --action GO"
    Then saved decision should include:
      | field        | value   |
      | ticker       | AAPL    |
      | action       | GO      |
      | entry        | 180.00  |
      | atr          | 1.50    |
      | stop_dist    | 3.00    |
      | init_stop    | 177.00  |
      | shares       | 25      |
      | risk_dollars | 75.00   |
      | banner       | GREEN   |

  Scenario: Save decision with options (delta-ATR method)
    Given banner is GREEN for "AAPL"
    And impulse timer expired for "AAPL"
    When I run "save-decision --ticker AAPL --entry 5.50 --atr 1.50 --delta 0.30 --method opt-delta-atr --action GO"
    Then the decision should be saved with action "GO"
    And method should be "opt-delta-atr"
    And contracts should be calculated correctly

  Scenario: Save decision with options (max-loss method)
    Given banner is GREEN for "AAPL"
    And impulse timer expired for "AAPL"
    When I run "save-decision --ticker AAPL --max-loss 0.75 --method opt-maxloss --action GO"
    Then the decision should be saved with action "GO"
    And method should be "opt-maxloss"
    And contracts should be calculated correctly

  Scenario: Duplicate decision for same ticker and date rejected
    Given I have already saved a decision for "AAPL" today
    When I run "save-decision --ticker AAPL --entry 180 --atr 1.5 --action GO"
    Then I should receive error "already have a decision.*today"
    And the decision should NOT be saved

  Scenario: Include correlation ID in saved decision
    Given correlation ID is "test-123"
    When I run "save-decision --ticker AAPL --entry 180 --atr 1.5 --action GO --corr-id test-123"
    Then saved decision should include correlation ID "test-123"
