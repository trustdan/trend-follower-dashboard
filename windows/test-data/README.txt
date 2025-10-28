Test Data - Sample JSON Responses
===================================

This directory contains sample JSON responses from tf-engine.exe for testing VBA parsing functions.

These files were generated during M17-M18 (JSON contract validation) and represent actual engine outputs with validated schemas.

Usage in VBA Tests:
-------------------
1. Load JSON from these files in TFTests.bas test functions
2. Pass to TFHelpers.ParseXXXJSON() functions
3. Verify parsing correctness

Files:
------
- size-stock-success.json              - Stock position sizing
- size-opt-delta-atr-success.json      - Option sizing (delta-ATR method)
- size-opt-maxloss-success.json        - Option sizing (max loss method)
- checklist-green-success.json         - All 6 checks passed (GREEN)
- checklist-yellow-success.json        - 1 check failed (YELLOW)
- checklist-red-success.json           - 2+ checks failed (RED)
- heat-check-success.json              - Heat check with open positions
- heat-check-empty-success.json        - Heat check with no positions
- timer-check-active-success.json      - Impulse timer active
- candidate-check-yes-success.json     - Ticker found in candidates
- candidate-check-no-success.json      - Ticker not in candidates
- candidates-list-success.json         - List of candidate tickers
- candidates-import-success.json       - Import result
- cooldown-check-active-success.json   - Bucket on cooldown
- cooldown-check-inactive-success.json - Bucket not on cooldown
- cooldowns-list-empty-success.json    - No active cooldowns
- cooldowns-list-with-data-success.json - Active cooldowns present
- positions-list-empty-success.json    - No open positions
- settings-get-all-success.json        - Application settings

Note: These are SUCCESS responses only. Error responses are tested via stderr output.

Created: 2025-10-27 (M20 - Windows Integration Package)
