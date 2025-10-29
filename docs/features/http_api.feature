Feature: HTTP API for Excel integration
  As an Excel user
  I want to call engine functions via HTTP
  So that VBA can interact with the trading engine

  Background:
    Given HTTP server is running on port 18888
    And I have initialized the database

  Scenario: Health check endpoint
    When I send GET to "/health"
    Then response status should be 200
    And response should contain "status": "ok"

  Scenario: Calculate position size via HTTP
    When I send POST to "/api/size" with:
      """
      {
        "equity": 10000,
        "risk_pct": 0.0075,
        "entry": 180,
        "atr": 1.5,
        "k": 2,
        "method": "stock"
      }
      """
    Then response status should be 200
    And response should contain "shares": 25
    And response should contain "risk_dollars": 75
    And response should contain "initial_stop": 177

  Scenario: Evaluate checklist via HTTP
    When I send POST to "/api/checklist" with:
      """
      {
        "ticker": "AAPL",
        "checks": {
          "uptrend": true,
          "above_ema": true,
          "volume": true,
          "clean_chart": true,
          "risk_reward": true,
          "timing": true
        }
      }
      """
    Then response status should be 200
    And response should contain "banner": "GREEN"
    And response should contain "missing_count": 0

  Scenario: Save decision via HTTP
    Given I have candidate "AAPL" for today
    And checklist is GREEN for "AAPL"
    And impulse timer expired for "AAPL"
    When I send POST to "/api/decision" with:
      """
      {
        "ticker": "AAPL",
        "action": "GO",
        "entry": 180,
        "atr": 1.5
      }
      """
    Then response status should be 200
    And response should contain "decision_id"

  Scenario: Gate validation failure returns 400
    Given banner is RED for "AAPL"
    When I send POST to "/api/decision" with GO action
    Then response status should be 400
    And response should contain error message about banner

  Scenario: List candidates via HTTP
    Given I have imported 3 candidates for today
    When I send GET to "/api/candidates"
    Then response status should be 200
    And response should contain array with 3 items

  Scenario: Get heat status via HTTP
    Given I have 2 open positions
    When I send GET to "/api/heat"
    Then response status should be 200
    And response should contain "portfolio_heat"
    And response should contain "portfolio_cap"
    And response should contain "buckets" array

  Scenario: Check impulse timer via HTTP
    Given checklist is GREEN for "AAPL"
    When I send GET to "/api/timer?ticker=AAPL"
    Then response status should be 200
    And response should contain "elapsed_seconds"
    And response should contain "remaining_seconds"
    And response should contain "ready" boolean

  Scenario: Check cooldown status via HTTP
    Given bucket "Tech/Comm" is in cooldown
    When I send GET to "/api/cooldown?bucket=Tech/Comm"
    Then response status should be 200
    And response should contain "in_cooldown": true
    And response should contain "expires_at"

  Scenario: CORS headers for Excel
    When I send OPTIONS to "/api/size"
    Then response should include CORS headers
    And Access-Control-Allow-Origin should be "*"

  Scenario: Error responses include correlation ID
    When I send invalid request
    Then response should include "correlation_id" field
