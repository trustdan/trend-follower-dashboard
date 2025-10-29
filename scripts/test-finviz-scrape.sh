#!/bin/bash
# Test script to debug FINVIZ scraping

URL="https://finviz.com/screener.ashx?v=111&f=ta_sma20_pa,ta_sma50_pa,ta_sma200_pa"

echo "Fetching FINVIZ page..."
echo "URL: $URL"
echo ""

# Fetch the page
curl -s -L \
  -H "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36" \
  "$URL" > /tmp/finviz-page.html

echo "Page fetched. Size: $(wc -c < /tmp/finviz-page.html) bytes"
echo ""

# Check for ticker links (old format)
echo "Looking for old format: /quote.ashx?t="
grep -o "/quote.ashx?t=[A-Z]*" /tmp/finviz-page.html | head -10
echo ""

# Check for any quote links
echo "Looking for any quote links:"
grep -o "quote[^\"]*" /tmp/finviz-page.html | head -10
echo ""

# Look for ticker-like patterns
echo "Looking for ticker patterns in table cells:"
grep -o "<td[^>]*>[A-Z][A-Z][A-Z]*</td>" /tmp/finviz-page.html | head -10
echo ""

# Check for JavaScript rendering
echo "Checking if content is JavaScript-rendered:"
grep -i "javascript" /tmp/finviz-page.html | head -5
echo ""

echo "Full HTML saved to: /tmp/finviz-page.html"
echo "Inspect with: cat /tmp/finviz-page.html | less"
