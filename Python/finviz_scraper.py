"""
FINVIZ Screener Web Scraper
Fetches ticker symbols from FINVIZ screener pages with pagination support.

Usage:
    # From Excel VBA via =PY() formula
    tickers = finviz_scraper.fetch_finviz_tickers("v=211&p=d&s=ta_newhigh")

    # Standalone testing
    python finviz_scraper.py
"""

import requests
from bs4 import BeautifulSoup
import time
import re
from typing import List, Optional

# Configuration
BASE_URL = "https://finviz.com/screener.ashx"
USER_AGENT = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
REQUEST_TIMEOUT = 10
MAX_RETRIES = 3
RATE_LIMIT_DELAY = 1  # seconds between requests


def fetch_finviz_tickers(query_string: str, max_pages: int = 10) -> List[str]:
    """
    Scrapes FINVIZ screener and returns list of tickers.

    Args:
        query_string: FINVIZ query parameters (e.g., "v=211&s=ta_newhigh")
        max_pages: Maximum number of pages to scrape (default: 10)

    Returns:
        List of ticker symbols (uppercase, deduped, sorted)

    Example:
        >>> tickers = fetch_finviz_tickers("v=211&s=ta_newhigh")
        >>> print(tickers)
        ['AAPL', 'MSFT', 'TSLA', ...]
    """
    all_tickers = []
    page_num = 1

    while page_num <= max_pages:
        # Calculate offset for pagination (FINVIZ shows 20 results per page)
        offset = (page_num - 1) * 20
        url = f"{BASE_URL}?{query_string}&r={offset + 1}" if offset > 0 else f"{BASE_URL}?{query_string}"

        # Fetch page with retries
        tickers = _fetch_page_tickers(url)

        if not tickers:
            # No more results, stop pagination
            break

        all_tickers.extend(tickers)

        # Check if we've reached the last page
        if len(tickers) < 20:
            break

        # Rate limiting
        if page_num < max_pages:
            time.sleep(RATE_LIMIT_DELAY)

        page_num += 1

    # Dedupe and sort
    unique_tickers = list(set(all_tickers))
    unique_tickers.sort()

    return unique_tickers


def _fetch_page_tickers(url: str) -> List[str]:
    """
    Fetches tickers from a single FINVIZ page with retry logic.

    Args:
        url: Complete FINVIZ URL

    Returns:
        List of tickers found on page (may be empty if error)
    """
    headers = {'User-Agent': USER_AGENT}

    for attempt in range(MAX_RETRIES):
        try:
            response = requests.get(url, headers=headers, timeout=REQUEST_TIMEOUT)

            if response.status_code == 200:
                return _parse_tickers_from_html(response.content)
            elif response.status_code == 404:
                # Page doesn't exist, stop trying
                return []
            else:
                # Server error, retry
                if attempt < MAX_RETRIES - 1:
                    time.sleep(2 ** attempt)  # Exponential backoff
                    continue
                else:
                    return []

        except requests.Timeout:
            if attempt < MAX_RETRIES - 1:
                time.sleep(2 ** attempt)
                continue
            else:
                return []
        except requests.RequestException:
            return []

    return []


def _parse_tickers_from_html(html_content: bytes) -> List[str]:
    """
    Parses ticker symbols from FINVIZ HTML content.

    Args:
        html_content: Raw HTML bytes from FINVIZ

    Returns:
        List of ticker symbols found in the screener table
    """
    soup = BeautifulSoup(html_content, 'html.parser')
    tickers = []

    # FINVIZ uses different table structures, try multiple selectors

    # Method 1: Look for screener-body-table-nw class (current FINVIZ structure)
    ticker_links = soup.find_all('a', class_='screener-link-primary')
    if ticker_links:
        for link in ticker_links:
            ticker = link.get_text().strip()
            if ticker and len(ticker) <= 5:  # Basic validation
                tickers.append(ticker.upper())
        return tickers

    # Method 2: Fallback - look for links to quote.ashx (ticker detail pages)
    ticker_links = soup.find_all('a', href=re.compile(r'quote\.ashx\?t='))
    if ticker_links:
        for link in ticker_links:
            ticker = link.get_text().strip()
            if ticker and len(ticker) <= 5:
                tickers.append(ticker.upper())
        return tickers

    # Method 3: Fallback - look for table with class 'table-light' or specific structure
    tables = soup.find_all('table')
    for table in tables:
        rows = table.find_all('tr')
        for row in rows[1:]:  # Skip header row
            cells = row.find_all('td')
            if len(cells) > 1:
                # Ticker is typically in the 2nd cell (index 1)
                ticker_cell = cells[1]
                ticker_text = ticker_cell.get_text().strip()

                # Clean ticker symbol
                ticker = re.sub(r'[^A-Z\-]', '', ticker_text.upper())
                if ticker and len(ticker) <= 5:
                    tickers.append(ticker)

    return list(set(tickers))  # Dedupe before returning


def normalize_tickers(raw_list: List[str]) -> List[str]:
    """
    Cleans and normalizes ticker symbols.

    Args:
        raw_list: Raw list of ticker strings (may contain special chars, duplicates)

    Returns:
        List of clean uppercase tickers (deduped, validated)

    Example:
        >>> normalize_tickers(['AAPL', 'msft', 'AAPL', 'T.S.L.A', ''])
        ['AAPL', 'MSFT', 'TSLA']
    """
    normalized = []

    for ticker in raw_list:
        if not ticker:
            continue

        # Convert to uppercase and strip whitespace
        clean = ticker.strip().upper()

        # Replace dots with dashes (e.g., BRK.B → BRK-B)
        clean = clean.replace('.', '-')

        # Remove any other special characters except dashes
        clean = re.sub(r'[^A-Z0-9\-]', '', clean)

        # Basic validation: 1-5 characters, at least one letter
        if clean and len(clean) <= 5 and re.search(r'[A-Z]', clean):
            normalized.append(clean)

    # Dedupe and sort
    return sorted(list(set(normalized)))


def get_ticker_count(query_string: str) -> Optional[int]:
    """
    Returns the total number of tickers for a given query (without fetching all pages).

    Args:
        query_string: FINVIZ query parameters

    Returns:
        Total ticker count, or None if error

    Example:
        >>> count = get_ticker_count("v=211&s=ta_newhigh")
        >>> print(f"Found {count} tickers")
    """
    url = f"{BASE_URL}?{query_string}"
    headers = {'User-Agent': USER_AGENT}

    try:
        response = requests.get(url, headers=headers, timeout=REQUEST_TIMEOUT)
        if response.status_code != 200:
            return None

        soup = BeautifulSoup(response.content, 'html.parser')

        # Look for total count text (e.g., "Total: 47")
        count_text = soup.find(text=re.compile(r'Total:\s*\d+'))
        if count_text:
            match = re.search(r'Total:\s*(\d+)', count_text)
            if match:
                return int(match.group(1))

        # Fallback: count rows in current page
        tickers = _parse_tickers_from_html(response.content)
        return len(tickers) if tickers else None

    except Exception:
        return None


def fetch_multiple_presets(preset_queries: dict) -> dict:
    """
    Batch fetch tickers from multiple FINVIZ presets.

    Args:
        preset_queries: Dict of {preset_name: query_string}

    Returns:
        Dict of {preset_name: [tickers]}

    Example:
        >>> presets = {
        ...     "Breakout": "v=211&s=ta_newhigh",
        ...     "Momentum": "v=211&s=ta_topgainers"
        ... }
        >>> results = fetch_multiple_presets(presets)
        >>> print(results)
        {'Breakout': ['AAPL', 'MSFT'], 'Momentum': ['TSLA', 'NVDA']}
    """
    results = {}

    for preset_name, query_string in preset_queries.items():
        tickers = fetch_finviz_tickers(query_string)
        results[preset_name] = tickers

        # Rate limiting between presets
        time.sleep(RATE_LIMIT_DELAY * 2)

    return results


def _test_scraper():
    """
    Standalone test function for development/debugging.
    Run: python finviz_scraper.py
    """
    print("=" * 70)
    print("FINVIZ Scraper Test")
    print("=" * 70)

    # Test preset: New Highs
    test_query = "v=211&f=ta_highlow52w_nh&ft=4"

    print(f"\nTesting query: {test_query}")
    print(f"URL: {BASE_URL}?{test_query}")
    print("\nFetching tickers...")

    tickers = fetch_finviz_tickers(test_query, max_pages=2)

    if tickers:
        print(f"\n✅ Success! Found {len(tickers)} tickers:")
        print(", ".join(tickers[:20]))  # Show first 20
        if len(tickers) > 20:
            print(f"... and {len(tickers) - 20} more")
    else:
        print("\n❌ No tickers found. Possible reasons:")
        print("  - No internet connection")
        print("  - FINVIZ changed HTML structure")
        print("  - Query returned no results")
        print("  - Rate limited by FINVIZ")

    print("\n" + "=" * 70)


if __name__ == "__main__":
    _test_scraper()
