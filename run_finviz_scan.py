"""
Standalone FINVIZ Scanner
Runs outside Excel, outputs tickers to copy/paste
"""

import sys
from pathlib import Path

# Add Python folder to path
sys.path.insert(0, str(Path(__file__).parent / "Python"))

from finviz_scraper import fetch_finviz_tickers

# Preset query strings (from your Presets table)
PRESETS = {
    "1": ("TF_BREAKOUT_LONG", "v=211&f=ta_highlow52w_nh&ft=4"),
    "2": ("TF_MOMENTUM_UPTREND", "v=211&f=ta_sma200_pa,ta_sma50_pa&ft=4"),
    "3": ("TF_UNUSUAL_VOLUME", "v=211&f=sh_relvol_o2&ft=4"),
    "4": ("TF_GAP_UP", "v=211&f=ta_gap_u5&ft=4"),
    "5": ("TF_STRONG_TREND", "v=211&f=ta_changeopen_u5&ft=4"),
}

def main():
    print("=" * 70)
    print("FINVIZ Scanner - Standalone")
    print("=" * 70)
    print("\nAvailable presets:")
    for key, (name, _) in PRESETS.items():
        print(f"  {key}. {name}")

    choice = input("\nEnter preset number (1-5): ").strip()

    if choice not in PRESETS:
        print("Invalid choice!")
        return

    preset_name, query_string = PRESETS[choice]

    print(f"\nScanning: {preset_name}")
    print(f"Query: {query_string}")
    print("\nFetching tickers from FINVIZ...")

    tickers = fetch_finviz_tickers(query_string, max_pages=3)

    if tickers:
        print(f"\n✓ Found {len(tickers)} tickers:")
        print("=" * 70)
        ticker_str = ", ".join(tickers)
        print(ticker_str)
        print("=" * 70)
        print("\nCOPY THE TICKERS ABOVE")
        print("Then in Excel:")
        print("  1. Click 'Import Candidates' button")
        print("  2. Paste tickers")
        print("  3. Click OK")
    else:
        print("\n✗ No tickers found")
        print("Possible reasons:")
        print("  - No internet connection")
        print("  - FINVIZ changed HTML structure")
        print("  - Query returned no results")

if __name__ == "__main__":
    main()
