# Interactive Mode Guide - FINVIZ Candidate Import

**Last Updated:** 2025-10-28
**Feature:** New interactive CLI for importing candidates

---

## What Is Interactive Mode?

Interactive mode is a **guided, menu-driven interface** for importing candidates from FINVIZ. No need to remember command-line flags or URLs - just run one command and follow the prompts!

---

## Quick Start

```cmd
cd C:\Users\Dan\excel-trading-dashboard\windows
tf-engine.exe interactive
```

That's it! The program will guide you through:
1. ✅ Selecting a screener preset
2. ✅ Configuring scraper options
3. ✅ Scraping FINVIZ
4. ✅ Reviewing results
5. ✅ Importing as today's candidates

---

## Step-by-Step Walkthrough

### Step 1: Select FINVIZ Screener

```
Step 1: Select FINVIZ Screener
----------------------------------------
Choose a screener preset?
👉 TF-Breakout-Long
   TF-Momentum-Uptrend
   TF-Unusual-Volume
   TF-Breakdown-Short
   TF-Momentum-Downtrend
   Enter Custom URL
```

**Use arrow keys** to navigate, **Enter** to select.

**Built-in presets:**
- **TF-Breakout-Long** - Large cap stocks making new highs, above SMA 50/200, sorted by relative volume
- **TF-Momentum-Uptrend** - Large cap stocks in strong uptrend, above SMA 50/200, sorted by market cap
- **TF-Unusual-Volume** - Large cap stocks with unusual volume, above SMA 50/200
- **TF-Breakdown-Short** - Large cap stocks making new lows, below SMA 50/200 (for short candidates)
- **TF-Momentum-Downtrend** - Large cap stocks in strong downtrend, below SMA 50/200 (for short candidates)
- **Enter Custom URL** - Paste your own FINVIZ screener URL

---

### Step 2: Configure Scraper Options

```
Step 2: Configure Scraper Options
----------------------------------------
Max pages to scrape (0 = unlimited): [10]
```

**Max pages:** How many screener pages to scrape (each page = ~20 tickers)
- Enter `10` for up to 200 tickers
- Enter `0` for unlimited (scrape all)
- Enter `5` for ~100 tickers

```
Rate limit between pages?
👉 0.5 seconds
   1 second (recommended)
   2 seconds
   3 seconds
```

**Rate limit:** Delay between page requests (to be polite to FINVIZ)
- `1 second` recommended - good balance
- `0.5 seconds` faster but more aggressive
- `2-3 seconds` more polite, slower

---

### Step 3: Scrape FINVIZ

```
Step 3: Scrape FINVIZ
----------------------------------------
Start scraping (y/N): y

Scraping FINVIZ...
```

Confirm and the scraper starts fetching data.

**What happens:**
- Fetches each page from FINVIZ
- Extracts ticker symbols
- Normalizes them (BRK.B → BRK-B)
- Removes duplicates
- Shows progress

---

### Step 4: Review Results

```
Step 4: Review Results
----------------------------------------
Tickers found: 47
Pages scraped: 3
More available: true
Date: 2025-10-28

Sample tickers:
  - AAPL
  - MSFT
  - NVDA
  - GOOGL
  - TSLA
  - META
  - AMZN
  - JPM
  - BAC
  - WMT
  ... and 37 more
```

**Review the results:**
- **Tickers found** - Total unique tickers scraped
- **Pages scraped** - How many pages were processed
- **More available** - Whether more tickers exist beyond max pages
- **Sample tickers** - First 10 tickers from the list

---

### Step 5: Import as Candidates

```
Step 5: Import as Candidates
----------------------------------------
Import 47 tickers as today's candidates (y/N): y

Enter preset name (for tracking) [Trend_Following_(SMA_20/50)]:
```

**Confirm import** and optionally customize the preset name.

**Success:**
```
========================================
  ✓ Import Complete!
========================================

Imported: 47 tickers
Date: 2025-10-28
Preset: Trend_Following_(SMA_20/50)

Next steps:
  1. Open TradingPlatform.xlsm
  2. Go to Dashboard to see candidates
  3. Evaluate trades using Trade Entry sheet

The 5 Hard Gates will now enforce that you only
trade tickers from today's imported candidates!
```

---

## Auto Mode (Skip All Prompts)

For daily automation, use `--auto` flag:

```cmd
tf-engine.exe interactive --auto
```

**What happens:**
- Uses default preset (TF-Breakout-Long)
- Default config (10 pages, 1 second rate limit)
- Automatically imports without confirmation
- Uses preset name from selection

**Perfect for:**
- Daily batch scripts
- Scheduled tasks
- Quick updates

---

## Advanced: Custom URLs

If you select **"Enter Custom URL"**, you can paste your own FINVIZ screener URL:

**Example custom screeners:**

### High Momentum Tech
```
https://finviz.com/screener.ashx?v=111&f=ind_technology,ta_perf_4wup20,ta_rsi_os50
```

### Dividend Growth
```
https://finviz.com/screener.ashx?v=111&f=fa_div_pos,fa_epsqoq_pos,fa_sales5years_pos
```

### Small Cap Breakouts
```
https://finviz.com/screener.ashx?v=111&f=cap_smallover,ta_highlow52w_nh,ta_sma20_pa
```

**How to create custom URLs:**
1. Go to https://finviz.com/screener.ashx
2. Apply your desired filters
3. Copy the URL from your browser
4. Paste it when prompted

---

## Comparison: Interactive vs Manual

### Old Way (Manual Commands)

```cmd
# Remember the command syntax
tf-engine.exe scrape-finviz --query "https://finviz.com/screener.ashx?v=111&f=ta_sma20_pa,ta_sma50_pa" --preset trend-following

# Then import (separate command)
tf-engine.exe import-candidates --preset trend-following
```

**Problems:**
- ❌ Have to remember command names
- ❌ Have to remember flag syntax
- ❌ Have to look up or save screener URLs
- ❌ Two separate commands
- ❌ Easy to make typos

### New Way (Interactive)

```cmd
# Just run one command
tf-engine.exe interactive
```

**Benefits:**
- ✅ Menu-driven (no memorization)
- ✅ Built-in presets
- ✅ Guided configuration
- ✅ One command does everything
- ✅ Hard to make mistakes

---

## Daily Workflow

### Morning Routine (Before Trading)

```cmd
cd C:\Users\Dan\excel-trading-dashboard\windows

# Import today's candidates
tf-engine.exe interactive --auto

# Open Excel
TradingPlatform.xlsm
```

**Or manually:**
```cmd
# Import candidates interactively
tf-engine.exe interactive
  → Select "TF-Breakout-Long"
  → Accept defaults
  → Confirm import

# Open Excel
TradingPlatform.xlsm
```

---

## Verifying Import

After importing, verify candidates were saved:

```cmd
# List today's candidates
tf-engine.exe list-candidates

# Check a specific ticker
tf-engine.exe check-candidate --ticker AAPL
```

**Expected output:**
```json
{
  "success": true,
  "in_candidates": true,
  "ticker": "AAPL",
  "date": "2025-10-28"
}
```

---

## Troubleshooting

### "Failed to scrape FINVIZ"

**Cause:** Internet connection or FINVIZ website issue

**Fix:**
- Check internet connection
- Try again in a few minutes
- Try a different preset
- Use `--auto` flag to skip confirmations

---

### "No tickers found"

**Cause:** Screener returned no results

**Fix:**
- Try a different preset
- Check if FINVIZ is accessible in browser
- Try lowering the filter criteria (custom URL)

---

### Interactive prompt not working

**Cause:** Terminal doesn't support interactive prompts

**Fix:** Use auto mode instead:
```cmd
tf-engine.exe interactive --auto
```

---

## Available Commands

| Command | Purpose |
|---------|---------|
| `tf-engine.exe interactive` | Launch interactive mode |
| `tf-engine.exe interactive --auto` | Auto mode (no prompts) |
| `tf-engine.exe list-candidates` | List imported candidates |
| `tf-engine.exe check-candidate --ticker AAPL` | Check if ticker is candidate |
| `tf-engine.exe scrape-finviz --help` | See manual scraper options |
| `tf-engine.exe import-candidates --help` | See manual import options |

---

## Built-in Presets (URLs)

For reference, here are the URLs used by built-in presets:

### TF-Breakout-Long
```
https://finviz.com/screener.ashx?v=111&p=d&s=ta_newhigh&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pa,ta_sma50_pa&o=-relativevolume
```
Large caps making new highs, above SMA 50/200, high volume, price >$20, sorted by relative volume

### TF-Momentum-Uptrend
```
https://finviz.com/screener.ashx?v=111&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pa,ta_sma50_pa&dr=y1&o=-marketcap
```
Large caps in strong uptrend, above SMA 50/200, high volume, price >$20, sorted by market cap

### TF-Unusual-Volume
```
https://finviz.com/screener.ashx?v=111&p=d&s=ta_unusualvolume&f=cap_largeover,sh_price_o20,ta_sma200_pa,ta_sma50_pa&o=-relativevolume
```
Large caps with unusual volume spikes, above SMA 50/200, sorted by relative volume

### TF-Breakdown-Short
```
https://finviz.com/screener.ashx?v=111&p=d&s=ta_newlow&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pb,ta_sma50_pb&o=-relativevolume
```
Large caps making new lows, below SMA 50/200 (short candidates), sorted by relative volume

### TF-Momentum-Downtrend
```
https://finviz.com/screener.ashx?v=111&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pb,ta_sma50_pb&dr=y1&o=-marketcap
```
Large caps in strong downtrend, below SMA 50/200 (short candidates), sorted by market cap

---

## Tips & Best Practices

### Daily Import
✅ **Import fresh candidates every trading day** (markets change!)

### Rate Limit
✅ **Use 1 second rate limit** (default) - respectful to FINVIZ

### Max Pages
✅ **10 pages is usually enough** (~200 tickers) - you don't need thousands

### Preset Names
✅ **Use descriptive preset names** - helps with tracking in database

### Verification
✅ **Always verify import** with `list-candidates` before trading

### Backup
✅ **If scraper fails**, you can manually import:
```cmd
echo AAPL > tickers.txt
echo MSFT >> tickers.txt
tf-engine.exe import-candidates --file tickers.txt --preset manual
```

---

## Integration with Excel

After importing candidates via interactive mode:

1. **Open** `TradingPlatform.xlsm`
2. **Go to Dashboard** - candidates should appear
3. **Go to Trade Entry** sheet
4. **Enter ticker** (e.g., AAPL)
5. **Click "Save GO"**
6. **Gate 2 validates** - ✅ "AAPL is in today's candidates"

If you enter a ticker NOT in candidates:
- Gate 2 fails - ❌ "TICKER not in today's candidates"
- Trade is rejected automatically

**The system enforces discipline!**

---

## Summary

**Old workflow:**
```
Look up FINVIZ URL → Copy URL → Remember scrape-finviz syntax →
Run scrape command → Remember import syntax → Run import command
```

**New workflow:**
```
tf-engine.exe interactive
```

**That's it!** ✨

---

**Questions?** See `QUICK_START.md` or `README.md` for more help.
