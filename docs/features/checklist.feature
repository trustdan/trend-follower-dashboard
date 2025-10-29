Feature: Entry checklist evaluation
  As a trader
  I want to evaluate my 6-item checklist
  So that I know if I have a valid trading setup

  Background:
    Given the checklist has 6 required items:
      | FromPreset    |
      | TrendPass     |
      | LiquidityPass |
      | TVConfirm     |
      | EarningsOK    |
      | JournalOK     |

  Scenario: Perfect setup - All 6 items checked
    Given I check all 6 items as true
    When I evaluate the checklist
    Then the banner should be "GREEN"
    And missing_count should be 0
    And missing_items should be empty
    And allow_save should be true

  Scenario: One item missing - Caution
    Given I check 5 items as true
    And "EarningsOK" is false
    When I evaluate the checklist
    Then the banner should be "YELLOW"
    And missing_count should be 1
    And missing_items should contain "EarningsOK"
    And allow_save should be false

  Scenario: Two items missing - No-go
    Given I check 4 items as true
    And "TrendPass" is false
    And "TVConfirm" is false
    When I evaluate the checklist
    Then the banner should be "RED"
    And missing_count should be 2
    And missing_items should contain "TrendPass"
    And missing_items should contain "TVConfirm"
    And allow_save should be false

  Scenario: Three items missing - No-go
    Given I check 3 items as true
    And "FromPreset" is false
    And "TrendPass" is false
    And "EarningsOK" is false
    When I evaluate the checklist
    Then the banner should be "RED"
    And missing_count should be 3
    And missing_items should contain "FromPreset"
    And missing_items should contain "TrendPass"
    And missing_items should contain "EarningsOK"
    And allow_save should be false

  Scenario: All items missing - No-go
    Given I check 0 items as true
    When I evaluate the checklist
    Then the banner should be "RED"
    And missing_count should be 6
    And missing_items should contain all 6 items
    And allow_save should be false

  Scenario Outline: Banner determination by missing count
    Given I have <missing> items unchecked
    When I evaluate the checklist
    Then the banner should be "<banner>"
    And allow_save should be <allow_save>

    Examples:
      | missing | banner | allow_save |
      | 0       | GREEN  | true       |
      | 1       | YELLOW | false      |
      | 2       | RED    | false      |
      | 3       | RED    | false      |
      | 4       | RED    | false      |
      | 5       | RED    | false      |
      | 6       | RED    | false      |

  Scenario: Reject empty ticker
    Given I check all 6 items as true
    And the ticker is empty
    When I evaluate the checklist
    Then I should receive an error "ticker is required"
