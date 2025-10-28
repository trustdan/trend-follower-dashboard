#!/bin/bash
# Test HTTP endpoints and compare with CLI

BASE_URL="http://127.0.0.1:18888"
DB="test-data/test-contracts.db"

echo "=== HTTP ENDPOINT PARITY TESTS ==="
echo ""

# Test 1: Position Sizing
echo "1. Position Sizing (POST /api/size)"
./tf-engine size --db "$DB" --entry 180 --atr 1.5 --k 2 --method stock --format json > /tmp/cli-size.json
curl -s -X POST "$BASE_URL/api/size" \
  -H "Content-Type: application/json" \
  -d '{"equity":10000,"risk_pct":0.0075,"entry":180,"atr_n":1.5,"k":2,"method":"stock"}' \
  > /tmp/http-size.json
echo "  CLI:  $(jq -c . /tmp/cli-size.json 2>/dev/null || cat /tmp/cli-size.json)"
echo "  HTTP: $(jq -c . /tmp/http-size.json 2>/dev/null || cat /tmp/http-size.json)"
echo ""

# Test 2: Settings  
echo "2. Settings (GET /api/settings)"
./tf-engine get-settings --db "$DB" --format json > /tmp/cli-settings.json
curl -s "$BASE_URL/api/settings" > /tmp/http-settings.json
echo "  CLI:  $(jq -c . /tmp/cli-settings.json 2>/dev/null || cat /tmp/cli-settings.json)"
echo "  HTTP: $(jq -c . /tmp/http-settings.json 2>/dev/null || cat /tmp/http-settings.json)"
echo ""

# Test 3: Heat Check
echo "3. Heat Check (GET /api/heat)"
./tf-engine check-heat --db "$DB" --add-risk 75 --format json > /tmp/cli-heat.json
curl -s "$BASE_URL/api/heat?add_r=75" > /tmp/http-heat.json
echo "  CLI fields:  $(jq -c 'keys' /tmp/cli-heat.json 2>/dev/null || echo 'error')"
echo "  HTTP fields: $(jq -c 'keys' /tmp/http-heat.json 2>/dev/null || echo 'error')"
echo ""

echo "✓ Tests complete. Files in /tmp/cli-*.json and /tmp/http-*.json"
