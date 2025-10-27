"""
FINVIZ Web Scraper for Excel Python Integration

Purpose:
    Automatically scrapes tickers from FINVIZ screener pages.
    Eliminates manual copy/paste workflow.

Usage in Excel:
    =PY("finviz_scraper.fetch_finviz_tickers", Presets!B2)

Dependencies:
    - requests
    - beautifulsoup4
    - lxml (parser)

Author: Generated from newest-Interactive_TF_Workbook_Plan.md
"""

import requests
from bs4 import BeautifulSoup
from typing import List, Optional
import time


def fetch_finviz_tickers(query_string: str, max_retries: int = 3) -> List[str]:
    """
    Scrapes FINVIZ screener page and returns list of ticker symbols.

    Args:
        query_string: The FINVIZ query string from tblPresets[QueryString]
                     Example: "v=211&p=d&s=ta_newhigh&f=cap_largeover..."
        max_retries: Number of retry attempts for failed requests

    Returns:
        List of ticker symbols (uppercase, deduped)
        Returns empty list if scraping fails

    Example:
        >>> tickers = fetch_finviz_tickers("v=211&p=d&s=ta_newhigh")
        >>> print(tickers)
        ['AAPL', 'MSFT', 'GOOGL', ...]
    """
    base_url = "https://finviz.com/screener.ashx"
    url = f"{base_url}?{query_string}"

    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
        'Accept-Language': 'en-US,en;q=0.5',
        'Connection': 'keep-alive',
    }

    for attempt in range(max_retries):
        try:
            response = requests.get(url, headers=headers, timeout=10)
            response.raise_for_status()  # Raise exception for 4xx/5xx status codes

            # Parse HTML
            soup = BeautifulSoup(response.content, 'lxml')

            # Find the screener results table
            # FINVIZ uses different table classes - try multiple selectors
            table = None

            # Try primary selector
            table = soup.find('table', {'class': 'table-light'})

            # Fallback: Find table with ticker links
            if not table:
                table = soup.find('table', {'id': 'screener-table'})

            # Fallback: Find any table with ticker-like content
            if not table:
                # Look for table containing ticker cells
                all_tables = soup.find_all('table')
                for t in all_tables:
                    if t.find('a', href=lambda x: x and '/quote.ashx?t=' in x):
                        table = t
                        break

            if not table:
                print(f"Warning: Could not find screener table in HTML (attempt {attempt + 1}/{max_retries})")
                if attempt < max_retries - 1:
                    time.sleep(1)  # Wait before retry
                    continue
                return []

            # Extract tickers from table
            tickers = []

            # Method 1: Look for ticker links in cells
            ticker_links = table.find_all('a', href=lambda x: x and '/quote.ashx?t=' in x)
            for link in ticker_links:
                ticker = link.get_text(strip=True)
                if ticker and len(ticker) <= 5:  # Basic validation
                    tickers.append(ticker.upper())

            # Method 2: If no links found, look in second column of data rows
            if not tickers:
                rows = table.find_all('tr')
                for row in rows[1:]:  # Skip header row
                    cells = row.find_all('td')
                    if len(cells) >= 2:
                        # Ticker is typically in 2nd column (index 1)
                        ticker_cell = cells[1]
                        ticker = ticker_cell.get_text(strip=True)
                        if ticker and len(ticker) <= 5:
                            tickers.append(ticker.upper())

            # Dedupe while preserving order
            seen = set()
            unique_tickers = []
            for ticker in tickers:
                if ticker not in seen:
                    seen.add(ticker)
                    unique_tickers.append(ticker)

            if unique_tickers:
                return unique_tickers
            else:
                print(f"Warning: No tickers found in table (attempt {attempt + 1}/{max_retries})")
                if attempt < max_retries - 1:
                    time.sleep(1)
                    continue
                return []

        except requests.exceptions.Timeout:
            print(f"Timeout error (attempt {attempt + 1}/{max_retries})")
            if attempt < max_retries - 1:
                time.sleep(2)
                continue
            return []

        except requests.exceptions.RequestException as e:
            print(f"Request error: {e} (attempt {attempt + 1}/{max_retries})")
            if attempt < max_retries - 1:
                time.sleep(2)
                continue
            return []

        except Exception as e:
            print(f"Unexpected error: {e}")
            return []

    return []


def normalize_tickers(raw_list: List[str]) -> List[str]:
    """
    Cleans and normalizes ticker symbols.

    Removes:
        - Special characters (except hyphens)
        - Duplicates
        - Blanks
        - Tickers longer than 5 characters

    Args:
        raw_list: List of raw ticker strings

    Returns:
        List of clean, uppercase, deduped ticker symbols

    Example:
        >>> normalize_tickers(['AAPL', 'aapl', 'MSFT.A', '', 'BRK-B', 'TOOLONG'])
        ['AAPL', 'MSFT-A', 'BRK-B']
    """
    normalized = []
    seen = set()

    for ticker in raw_list:
        if not ticker:
            continue

        # Clean ticker
        clean = str(ticker).strip().upper()
        clean = clean.replace('.', '-')  # Convert periods to hyphens

        # Remove invalid characters (keep only letters, numbers, hyphens)
        clean = ''.join(c for c in clean if c.isalnum() or c == '-')

        # Validate length and uniqueness
        if clean and len(clean) <= 5 and clean not in seen:
            normalized.append(clean)
            seen.add(clean)

    return normalized


def fetch_multiple_presets(preset_queries: dict) -> dict:
    """
    Fetches tickers from multiple presets in one call.

    Args:
        preset_queries: Dict mapping preset names to query strings
                       Example: {"TF_BREAKOUT_LONG": "v=211&p=d&s=ta_newhigh...",
                                "TF_MOMENTUM": "v=211&p=d&f=cap_largeover..."}

    Returns:
        Dict mapping preset names to lists of tickers
        Example: {"TF_BREAKOUT_LONG": ["AAPL", "MSFT"],
                 "TF_MOMENTUM": ["GOOGL", "AMZN"]}

    Note:
        Adds 1-second delay between requests to avoid rate limiting
    """
    results = {}

    for preset_name, query_string in preset_queries.items():
        tickers = fetch_finviz_tickers(query_string)
        results[preset_name] = tickers

        # Rate limiting: wait 1 second between requests
        if len(preset_queries) > 1:
            time.sleep(1)

    return results


def get_ticker_count(query_string: str) -> Optional[int]:
    """
    Gets the total count of tickers matching a screener without fetching all tickers.

    Args:
        query_string: The FINVIZ query string

    Returns:
        Integer count of matching tickers, or None if count cannot be determined

    Example:
        >>> count = get_ticker_count("v=211&p=d&s=ta_newhigh")
        >>> print(f"Found {count} tickers")
        Found 47 tickers
    """
    base_url = "https://finviz.com/screener.ashx"
    url = f"{base_url}?{query_string}"

    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
    }

    try:
        response = requests.get(url, headers=headers, timeout=10)
        response.raise_for_status()

        soup = BeautifulSoup(response.content, 'lxml')

        # Look for count text (e.g., "Total: 47")
        count_element = soup.find('td', {'class': 'count-text'})
        if count_element:
            count_text = count_element.get_text(strip=True)
            # Extract number from text like "Total: 47"
            import re
            match = re.search(r'(\d+)', count_text)
            if match:
                return int(match.group(1))

        # Fallback: count rows in table
        table = soup.find('table', {'class': 'table-light'})
        if table:
            rows = table.find_all('tr')
            return len(rows) - 1  # Subtract header row

        return None

    except Exception as e:
        print(f"Error getting ticker count: {e}")
        return None


# Test function for development
def _test_scraper():
    """
    Test function - not called by Excel.
    Run this in Python to verify scraper works before using in Excel.
    """
    # Test query: 52-week highs, large cap, above SMAs
    test_query = "v=211&p=d&s=ta_newhigh&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pa,ta_sma200_pa&o=-relativevolume"

    print("Testing FINVIZ scraper...")
    print(f"Query: {test_query[:50]}...")

    tickers = fetch_finviz_tickers(test_query)

    if tickers:
        print(f"\n✅ Success! Found {len(tickers)} tickers:")
        print(f"First 10: {tickers[:10]}")
    else:
        print("\n❌ No tickers found. Check FINVIZ website or query string.")

    # Test normalization
    test_raw = ['AAPL', 'aapl', 'MSFT.A', '', 'BRK-B', 'TOOLONGticker', '123']
    normalized = normalize_tickers(test_raw)
    print(f"\nNormalization test:")
    print(f"Input:  {test_raw}")
    print(f"Output: {normalized}")


if __name__ == "__main__":
    # Run test if executed directly
    _test_scraper()

if __name__ == "__main__":
    import argparse, json, csv, os, sys
    from datetime import datetime

    parser = argparse.ArgumentParser(description="Fetch FINVIZ tickers and write a CSV")
    g = parser.add_mutually_exclusive_group(required=True)
    g.add_argument("--query", help="FINVIZ query string (e.g., v=211&p=d&f=...)")
    g.add_argument("--url", help="Full FINVIZ URL")
    g.add_argument("--preset-file", help="JSON file with list of {name, query|url}")
    parser.add_argument("--out", required=True, help="Output CSV path")
    args = parser.parse_args()

    rows = []  # each row: {"Ticker": "...", "Preset": "...", "AsOf": "YYYY-MM-DD"}
    asof = datetime.utcnow().date().isoformat()

    def write_rows(path, rows):
        os.makedirs(os.path.dirname(path), exist_ok=True)
        with open(path, "w", newline="", encoding="utf-8") as f:
            w = csv.DictWriter(f, fieldnames=["Ticker", "Preset", "AsOf"])
            w.writeheader()
            for r in rows:
                w.writerow(r)

    if args.preset_file:
        with open(args.preset_file, "r", encoding="utf-8") as f:
            presets = json.load(f)
        for p in presets:
            name = p.get("name", "Preset")
            q = p.get("query")
            u = p.get("url")
            tickers = fetch_finviz_tickers(q if q else u)
            for t in tickers:
                rows.append({"Ticker": t, "Preset": name, "AsOf": asof})

    else:
        q_or_u = args.query if args.query else args.url
        tickers = fetch_finviz_tickers(q_or_u)
        for t in tickers:
            rows.append({"Ticker": t, "Preset": "Ad-hoc", "AsOf": asof})

    write_rows(args.out, rows)
    print(f"[OK] Wrote {len(rows)} tickers → {args.out}")
