#!/bin/bash
# Script to capture JSON examples from all CLI commands
# For M17-M18: JSON Contracts & Validation

DB="test-data/test-contracts.db"
OUT_DIR="test-data/json-examples"
RESPONSES="$OUT_DIR/responses"
ERRORS="$OUT_DIR/errors"

# Helper function to capture clean JSON (filter out log lines)
capture_json() {
    local output_file="$1"
    shift
    ./tf-engine "$@" --db "$DB" 2>/dev/null | grep -v '"level"' | grep -v "^$" > "$output_file"
}

# Helper function to capture error JSON
capture_error() {
    local output_file="$1"
    shift
    ./tf-engine "$@" --db "$DB" 2>&1 | grep -v '"level"' | grep -v "^$" > "$output_file"
}

echo "Capturing JSON examples for M17-M18..."

# =============================================================================
# 1. POSITION SIZING
# =============================================================================
echo "1. Position Sizing..."

# Stock method
capture_json "$RESPONSES/size-stock-success.json" \
    size --entry 180 --atr 1.5 --k 2 --method stock

# Options - delta-ATR method
capture_json "$RESPONSES/size-opt-delta-atr-success.json" \
    size --entry 5.00 --atr 0.50 --k 2 --delta 0.50 --method opt-delta-atr

# Options - maxloss method
capture_json "$RESPONSES/size-opt-maxloss-success.json" \
    size --entry 5.00 --atr 0.50 --k 2 --maxloss 50 --method opt-maxloss

# =============================================================================
# 2. CHECKLIST EVALUATION
# =============================================================================
echo "2. Checklist Evaluation..."

# All checks passed (GREEN)
capture_json "$RESPONSES/checklist-green-success.json" \
    checklist --ticker AAPL \
    --higher-high \
    --down-to-up \
    --strong \
    --good-bars \
    --fresh \
    --uptrend

# One check missing (YELLOW)
capture_json "$RESPONSES/checklist-yellow-success.json" \
    checklist --ticker AAPL \
    --higher-high \
    --down-to-up \
    --strong \
    --good-bars \
    --uptrend

# Two+ checks missing (RED)
capture_json "$RESPONSES/checklist-red-success.json" \
    checklist --ticker AAPL \
    --higher-high \
    --strong \
    --uptrend

# =============================================================================
# 3. HEAT MANAGEMENT
# =============================================================================
echo "3. Heat Management..."

# Check heat with no positions
capture_json "$RESPONSES/heat-check-empty-success.json" \
    check-heat --add-r 75

# =============================================================================
# 4. SETTINGS
# =============================================================================
echo "4. Settings..."

# Get all settings
capture_json "$RESPONSES/settings-get-all-success.json" \
    get-settings

# =============================================================================
# 5. CANDIDATES
# =============================================================================
echo "5. Candidates..."

# Import candidates first
./tf-engine import-candidates --db "$DB" --preset TEST --tickers "AAPL,MSFT,NVDA" 2>/dev/null

# List candidates for today
capture_json "$RESPONSES/candidates-list-success.json" \
    list-candidates --date $(date +%Y-%m-%d)

# Check if ticker is in candidates (YES)
capture_json "$RESPONSES/candidate-check-yes-success.json" \
    check-candidate --ticker AAPL

# Check if ticker is in candidates (NO)
capture_json "$RESPONSES/candidate-check-no-success.json" \
    check-candidate --ticker TSLA

# =============================================================================
# 6. IMPULSE BRAKE TIMER
# =============================================================================
echo "6. Impulse Brake Timer..."

# Start a timer (happens during checklist with GREEN banner)
./tf-engine checklist --db "$DB" --ticker AAPL \
    --higher-high --down-to-up --strong --good-bars --fresh --uptrend 2>/dev/null >/dev/null

# Check timer status
capture_json "$RESPONSES/timer-check-active-success.json" \
    check-timer --ticker AAPL

# Check timer for ticker with no timer
capture_json "$RESPONSES/timer-check-none-success.json" \
    check-timer --ticker MSFT

# =============================================================================
# 7. COOLDOWNS
# =============================================================================
echo "7. Bucket Cooldowns..."

# List cooldowns (initially empty)
capture_json "$RESPONSES/cooldowns-list-empty-success.json" \
    list-cooldowns

# Trigger a cooldown for testing
./tf-engine trigger-cooldown --db "$DB" --bucket "Tech/Comm" --hours 24 2>/dev/null

# Check cooldown status (active)
capture_json "$RESPONSES/cooldown-check-active-success.json" \
    check-cooldown --bucket "Tech/Comm"

# Check cooldown status (not active)
capture_json "$RESPONSES/cooldown-check-inactive-success.json" \
    check-cooldown --bucket "Energy"

# List cooldowns (with active cooldown)
capture_json "$RESPONSES/cooldowns-list-with-data-success.json" \
    list-cooldowns

# =============================================================================
# 8. POSITIONS
# =============================================================================
echo "8. Position Management..."

# List positions (initially empty)
capture_json "$RESPONSES/positions-list-empty-success.json" \
    list-positions

# =============================================================================
# 9. SAVE DECISION
# =============================================================================
echo "9. Save Decision (will capture errors for gate rejections)..."

# This will likely fail due to gates, but we'll capture the structure
# We'll test gate enforcement separately

# =============================================================================
# ERROR CASES
# =============================================================================
echo "10. Error Cases..."

# Invalid entry price (negative)
capture_error "$ERRORS/size-invalid-entry.json" \
    size --entry -180 --atr 1.5 --k 2 --method stock

# Missing required parameters
capture_error "$ERRORS/size-missing-params.json" \
    size --method stock

# Invalid method
capture_error "$ERRORS/size-invalid-method.json" \
    size --entry 180 --atr 1.5 --k 2 --method invalid

echo ""
echo "✓ JSON examples captured to $OUT_DIR"
echo ""
echo "Next steps:"
echo "  1. Review captured JSON files"
echo "  2. Document schema for each response"
echo "  3. Test HTTP endpoints for parity"
