#!/bin/bash
# Capture clean JSON examples with --format json flag
# All commands now output ONLY JSON (no logs, no mixed text)

DB="test-data/test-contracts.db"
OUT_DIR="test-data/json-examples"
RESPONSES="$OUT_DIR/responses"
ERRORS="$OUT_DIR/errors"

echo "Capturing clean JSON examples with --format json..."

# =============================================================================
# 1. POSITION SIZING
# =============================================================================
echo "1. Position Sizing..."

./tf-engine size --db "$DB" --entry 180 --atr 1.5 --k 2 --method stock --format json \
    > "$RESPONSES/size-stock-success.json"

./tf-engine size --db "$DB" --entry 5.00 --atr 0.50 --k 2 --delta 0.50 --method opt-delta-atr --format json \
    > "$RESPONSES/size-opt-delta-atr-success.json"

./tf-engine size --db "$DB" --entry 5.00 --atr 0.50 --k 2 --maxloss 50 --method opt-maxloss --format json \
    > "$RESPONSES/size-opt-maxloss-success.json"

# =============================================================================
# 2. CHECKLIST EVALUATION
# =============================================================================
echo "2. Checklist Evaluation..."

./tf-engine checklist --db "$DB" --ticker AAPL \
    --from-preset --trend-pass --liquidity-pass --tv-confirm --earnings-ok --journal-ok --format json \
    > "$RESPONSES/checklist-green-success.json"

./tf-engine checklist --db "$DB" --ticker MSFT \
    --from-preset --trend-pass --liquidity-pass --tv-confirm --earnings-ok --format json \
    > "$RESPONSES/checklist-yellow-success.json"

./tf-engine checklist --db "$DB" --ticker NVDA \
    --from-preset --trend-pass --liquidity-pass --tv-confirm --format json \
    > "$RESPONSES/checklist-red-success.json"

# =============================================================================
# 3. HEAT MANAGEMENT
# =============================================================================
echo "3. Heat Management..."

./tf-engine check-heat --db "$DB" --add-risk 75 --add-bucket "Tech/Comm" --format json \
    > "$RESPONSES/heat-check-success.json"

./tf-engine check-heat --db "$DB" --format json \
    > "$RESPONSES/heat-check-empty-success.json"

# =============================================================================
# 4. SETTINGS
# =============================================================================
echo "4. Settings..."

./tf-engine get-settings --db "$DB" --format json \
    > "$RESPONSES/settings-get-all-success.json"

# =============================================================================
# 5. CANDIDATES
# =============================================================================
echo "5. Candidates..."

# Import candidates first
./tf-engine import-candidates --db "$DB" --preset TEST --tickers "AAPL,MSFT,NVDA" --format json \
    > "$RESPONSES/candidates-import-success.json"

./tf-engine list-candidates --db "$DB" --date $(date +%Y-%m-%d) --format json \
    > "$RESPONSES/candidates-list-success.json"

./tf-engine check-candidate --db "$DB" --ticker AAPL --format json \
    > "$RESPONSES/candidate-check-yes-success.json"

./tf-engine check-candidate --db "$DB" --ticker TSLA --format json \
    > "$RESPONSES/candidate-check-no-success.json"

# =============================================================================
# 6. IMPULSE BRAKE TIMER
# =============================================================================
echo "6. Impulse Brake Timer..."

./tf-engine check-timer --db "$DB" --ticker AAPL --format json \
    > "$RESPONSES/timer-check-active-success.json"

./tf-engine check-timer --db "$DB" --ticker MSFT --format json \
    > "$RESPONSES/timer-check-none-success.json"

# =============================================================================
# 7. COOLDOWNS
# =============================================================================
echo "7. Bucket Cooldowns..."

./tf-engine list-cooldowns --db "$DB" --format json \
    > "$RESPONSES/cooldowns-list-empty-success.json"

./tf-engine trigger-cooldown --db "$DB" --bucket "Tech/Comm" --format json \
    > "$RESPONSES/cooldown-trigger-success.json"

./tf-engine check-cooldown --db "$DB" --bucket "Tech/Comm" --format json \
    > "$RESPONSES/cooldown-check-active-success.json"

./tf-engine check-cooldown --db "$DB" --bucket "Energy" --format json \
    > "$RESPONSES/cooldown-check-inactive-success.json"

./tf-engine list-cooldowns --db "$DB" --format json \
    > "$RESPONSES/cooldowns-list-with-data-success.json"

# =============================================================================
# 8. POSITIONS
# =============================================================================
echo "8. Position Management..."

./tf-engine list-positions --db "$DB" --format json \
    > "$RESPONSES/positions-list-empty-success.json"

# =============================================================================
# ERROR CASES
# =============================================================================
echo "9. Error Cases..."

# Invalid entry price (negative)
./tf-engine size --db "$DB" --entry -180 --atr 1.5 --k 2 --method stock --format json \
    2>&1 > "$ERRORS/size-invalid-entry.json" || true

# Missing required parameters
./tf-engine size --db "$DB" --method stock --format json \
    2>&1 > "$ERRORS/size-missing-params.json" || true

# Invalid method
./tf-engine size --db "$DB" --entry 180 --atr 1.5 --k 2 --method invalid --format json \
    2>&1 > "$ERRORS/size-invalid-method.json" || true

echo ""
echo "✓ Clean JSON examples captured to $OUT_DIR"
echo ""
echo "File count:"
ls -1 "$RESPONSES" | wc -l | xargs echo "  Success responses:"
ls -1 "$ERRORS" | wc -l | xargs echo "  Error responses:"
