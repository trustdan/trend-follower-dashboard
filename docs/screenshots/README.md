# TF-Engine Screenshots Requirements

**Version:** 1.0.0
**Last Updated:** 2025-10-29

This document specifies all screenshots needed for user documentation.

---

## Overview

User documentation references 30+ screenshots to illustrate:
- UI screens and workflows
- Banner states (RED/YELLOW/GREEN)
- TradingView integration
- Error messages and warnings

**Current Status:** Placeholder images only. Real screenshots needed before final release.

---

## Screenshot Capture Tools

**Recommended tools:**

1. **Windows Snipping Tool** (built-in)
   - Win+Shift+S → Select area
   - Simple, fast, built-in

2. **ShareX** (free, open source)
   - https://getsharex.com/
   - Advanced features (annotations, auto-upload)

3. **Greenshot** (free, open source)
   - https://getgreenshot.org/
   - Easy annotations

**Settings:**
- Format: PNG (lossless)
- Resolution: 1920×1080 or native screen resolution
- DPI: 96-144 (standard screen DPI)
- Color space: sRGB

---

## Screenshot Naming Convention

```
[screen]-[variant]-[state].png

Examples:
dashboard-empty.png
dashboard-with-positions.png
checklist-banner-red.png
checklist-banner-yellow.png
checklist-banner-green.png
checklist-timer-active.png
gates-all-pass.png
gates-gate4-fail.png
```

**Rules:**
- Lowercase, hyphen-separated
- Descriptive names (not screenshot01.png)
- Include state/variant if multiple versions

---

## Required Screenshots by Category

### 1. Dashboard (3 screenshots)

#### dashboard-empty.png
- **Purpose:** First-time user view
- **State:**
  - No open positions
  - Empty candidates list
  - Portfolio heat at 0%
  - "Welcome" message visible
- **Highlights:**
  - Clean, minimal interface
  - Navigation menu visible
  - Theme toggle visible

#### dashboard-with-positions.png
- **Purpose:** Active trading view
- **State:**
  - 3-5 open positions in table
  - Portfolio heat at 60-80% of cap
  - Candidates count: 12
  - No cooldowns
- **Highlights:**
  - Position table with ticker, entry, stop, risk, days
  - Portfolio heat gauge (green/yellow zone)
  - Quick action buttons visible

#### dashboard-cooldowns.png
- **Purpose:** Show cooldown enforcement
- **State:**
  - 1-2 positions open
  - Cooldowns section visible:
    - Ticker cooldown: AAPL until 2025-11-05
    - Sector cooldown: Tech/Comm until 2025-11-12
- **Highlights:**
  - Cooldowns prominently displayed
  - Expiration dates clear

---

### 2. Scanner (2 screenshots)

#### scanner-results.png
- **Purpose:** FINVIZ scan successful
- **State:**
  - Results table with 50-100 tickers
  - Various sectors represented
  - Some tickers selected (checkboxes)
- **Highlights:**
  - "Import Selected" button enabled
  - Sector distribution visible
  - "Open in TradingView" links visible

#### scanner-empty.png
- **Purpose:** No results from scan
- **State:**
  - "0 candidates found" message
  - Empty table
  - Suggestion to adjust FINVIZ filters
- **Highlights:**
  - User-friendly empty state

---

### 3. Checklist & Banner (7 screenshots)

#### checklist-banner-red.png
- **Purpose:** Show RED banner (missing required items)
- **State:**
  - Trade data filled (AAPL, 180.50, 2.35, Tech/Comm)
  - Only 3 of 5 required gates checked
  - 2 quality items checked
  - Banner: Large RED gradient at top
  - Message: "STOP - 2 required items missing"
- **Highlights:**
  - RED banner very prominent
  - Missing items clearly indicated
  - "Save Evaluation" button disabled (grayed out)

#### checklist-banner-yellow.png
- **Purpose:** Show YELLOW banner (low quality score)
- **State:**
  - Trade data filled
  - All 5 required gates checked ✓
  - Only 2 quality items checked
  - Quality score: 2.0 / 4.0
  - Threshold: 3.0
  - Banner: Large YELLOW/ORANGE gradient
  - Message: "CAUTION - Quality score 2.0 < 3.0 threshold"
- **Highlights:**
  - YELLOW banner distinct from RED
  - Quality score display clear
  - "Save Evaluation" button enabled

#### checklist-banner-green.png
- **Purpose:** Show GREEN banner (all requirements met)
- **State:**
  - Trade data filled
  - All 5 required gates checked ✓
  - 4 quality items checked ✓
  - Quality score: 4.0 / 4.0
  - Threshold: 3.0
  - Banner: Large GREEN gradient
  - Message: "GO - All gates pass, quality 4.0 ≥ 3.0"
- **Highlights:**
  - GREEN banner prominent
  - Clear success state
  - "Save Evaluation" button enabled

#### checklist-timer-active.png
- **Purpose:** Show 2-minute countdown timer
- **State:**
  - After clicking "Save Evaluation"
  - Timer: 1:23 (counting down)
  - Progress bar showing time elapsed
  - Message: "Impulse brake - wait 1:23 before saving GO decision"
- **Highlights:**
  - Timer large and visible
  - Cannot be skipped (visual reinforcement)

#### checklist-timer-complete.png
- **Purpose:** Show timer reached 0:00
- **State:**
  - Timer: 0:00
  - Message: "2-minute timer complete - proceed to sizing"
  - Link to Position Sizing screen
- **Highlights:**
  - Clear completion state

#### checklist-form-filled.png
- **Purpose:** Show complete checklist form
- **State:**
  - All fields filled
  - All gates checked
  - Journal note filled with example text
- **Highlights:**
  - Entire form visible (scroll to show all)

#### checklist-quality-items.png
- **Purpose:** Show optional quality items section
- **State:**
  - Optional items section expanded
  - Some items checked, some unchecked
  - Journal textarea with example note
- **Highlights:**
  - Distinction between required (top) and optional (bottom)

---

### 4. Position Sizing (3 screenshots)

#### sizing-form.png
- **Purpose:** Show pre-filled form from checklist
- **State:**
  - Method: Stock (selected)
  - Trade data filled from checklist
  - "Calculate Position Size" button enabled
- **Highlights:**
  - Clean form layout
  - Pre-population working

#### sizing-results.png
- **Purpose:** Show calculation results
- **State:**
  - Results displayed:
    - Shares per unit: 159
    - Risk per unit: $747.30
    - Initial stop: $175.80
    - Add-on schedule table (4 rows)
    - Total max position: 636 shares
  - "Save Position Plan" button enabled
- **Highlights:**
  - Results clearly formatted
  - Add-on schedule table readable

#### sizing-concentration-warning.png
- **Purpose:** Show concentration warning
- **State:**
  - Results displayed (high-priced stock)
  - Warning: "⚠️ Position is 28% of equity - consider reducing"
  - Warning banner in orange/yellow
- **Highlights:**
  - Warning prominent but not blocking
  - Can still save (warning, not error)

---

### 5. Heat Check (3 screenshots)

#### heat-check-within-caps.png
- **Purpose:** Show successful heat check (trade approved)
- **State:**
  - Current heat gauges displayed
  - "Check Heat for This Trade" button clicked
  - Results:
    - ✓ Portfolio heat: $3,200 / $4,000 (80%) - WITHIN CAP
    - ✓ Bucket heat: $1,125 / $1,500 (75%) - WITHIN CAP
  - Verdict: "✓ TRADE APPROVED" (green banner)
- **Highlights:**
  - Green success banner
  - Checkmarks clear

#### heat-check-exceeds-cap.png
- **Purpose:** Show heat cap exceeded (trade rejected)
- **State:**
  - Current heat gauges displayed
  - Results:
    - ✗ Portfolio heat: $4,250 / $4,000 (106%) - EXCEEDS CAP by $250
    - OR ✗ Bucket heat: $1,650 / $1,500 (110%) - EXCEEDS CAP by $150
  - Verdict: "✗ TRADE REJECTED" (red banner)
  - Suggestions displayed:
    - "Reduce position size to XX shares"
    - "Close existing position first"
    - "Choose different sector"
- **Highlights:**
  - Red rejection banner
  - Overage amount clear
  - Actionable suggestions

#### heat-gauges.png
- **Purpose:** Show current heat visualization
- **State:**
  - Portfolio heat gauge: Bar chart (60% full, green)
  - Bucket heat table:
    - Tech/Comm: $1,125 / $1,500 (75%) - Yellow
    - Energy: $850 / $1,500 (57%) - Green
    - Financials: $0 / $1,500 (0%) - Green
- **Highlights:**
  - Color-coding (green/yellow/red based on %)
  - Table format readable

---

### 6. Trade Entry & Gates (4 screenshots)

#### trade-entry-summary.png
- **Purpose:** Show trade summary review
- **State:**
  - Trade summary card displays all data:
    - Ticker, direction, entry, stop, risk
    - Shares, sector, quality score
  - "Run Final Gate Check" button visible
- **Highlights:**
  - Summary complete and organized
  - Read-only (review only)

#### gates-all-pass.png
- **Purpose:** Show all 5 gates passing
- **State:**
  - Gate check results displayed:
    - Gate 1: Banner Status → GREEN ✓
    - Gate 2: Impulse Brake → 2:15 elapsed ✓
    - Gate 3: Cooldown Status → Not on cooldown ✓
    - Gate 4: Heat Caps → Within caps ✓
    - Gate 5: Sizing Completed → Plan saved ✓
  - Overall result: "ALL GATES PASS ✓" (large green banner)
  - "SAVE GO DECISION" button enabled (green, prominent)
  - "SAVE NO-GO DECISION" button visible (red, smaller)
- **Highlights:**
  - Green checkmarks for each gate
  - Prominent success banner
  - GO button clearly enabled

#### gates-gate2-fail.png
- **Purpose:** Show Gate 2 (impulse brake) failure
- **State:**
  - Gate 1: ✓ Pass
  - Gate 2: → RED ✗ "1:45 elapsed (need 2:00)"
  - Gate 3-5: Not evaluated (grayed out)
  - Overall result: "GATE CHECK FAILED ✗" (red banner)
  - "SAVE GO DECISION" button disabled (grayed out)
  - "SAVE NO-GO DECISION" button enabled
- **Highlights:**
  - Red X on failing gate
  - Time remaining clear
  - GO button disabled (visual enforcement)

#### gates-gate4-fail.png
- **Purpose:** Show Gate 4 (heat caps) failure
- **State:**
  - Gate 1-3: ✓ Pass
  - Gate 4: → RED ✗ "Portfolio heat exceeds cap by $250"
  - Gate 5: Not evaluated
  - Overall result: "GATE CHECK FAILED ✗"
  - Suggestions: "Reduce size or close positions"
  - GO button disabled
- **Highlights:**
  - Specific failure reason
  - Overage amount
  - Actionable suggestions

---

### 7. Calendar (2 screenshots)

#### calendar-view.png
- **Purpose:** Show 10-week calendar grid
- **State:**
  - 10 columns (weeks): 2 past, current, 7 future
  - Rows: Sector buckets (Tech/Comm, Energy, Financials, etc.)
  - Cells contain tickers (e.g., "AAPL, MSFT" in Tech/Comm week 1)
  - Current week highlighted (bold border)
  - Color-coding:
    - Green: Low heat (< 50% bucket cap)
    - Yellow: Medium heat (50-80%)
    - Red: High heat (> 80%)
- **Highlights:**
  - Grid layout clear
  - Color-coding visible
  - Current week distinct

#### calendar-tooltip.png
- **Purpose:** Show cell hover tooltip
- **State:**
  - Hovering over cell with "AAPL, MSFT"
  - Tooltip displays:
    - AAPL: Entry $180.50, Risk $747, Day 12
    - MSFT: Entry $350.20, Risk $820, Day 8
  - Tooltip positioned near cursor
- **Highlights:**
  - Tooltip readable
  - Position details clear

---

### 8. Settings (2 screenshots)

#### settings-account.png
- **Purpose:** Show account settings form
- **State:**
  - All fields filled with example values:
    - Equity: 100000
    - Risk %: 0.75
    - Portfolio cap: 4.0
    - Bucket cap: 1.5
    - Max units: 4
  - "Save Settings" button enabled
- **Highlights:**
  - Form layout clear
  - Tooltips visible (hover states)

#### settings-finviz.png
- **Purpose:** Show FINVIZ presets section
- **State:**
  - Presets table with 2-3 entries:
    - Name: TF Breakout Long
    - URL: https://finviz.com/...
    - Actions: Edit, Delete
  - "Add New Preset" button visible
- **Highlights:**
  - Table format
  - Action buttons clear

---

### 9. Theme Toggle (2 screenshots)

#### theme-day-mode.png
- **Purpose:** Show light theme
- **State:**
  - Dashboard or Checklist screen in day mode
  - Light background, dark text
  - Banner gradients vibrant
- **Highlights:**
  - High contrast, readable
  - Clean appearance

#### theme-night-mode.png
- **Purpose:** Show dark theme
- **State:**
  - Same screen as day mode but in night mode
  - Dark background, light text
  - Banner gradients still distinct
- **Highlights:**
  - Easy on eyes
  - Banner colors preserved

---

### 10. TradingView Integration (5 screenshots)

#### tradingview-script.png
- **Purpose:** Show TradingView chart with Ed-Seykota script
- **State:**
  - Ticker: Any liquid stock (e.g., SPY)
  - Timeframe: Daily
  - Ed-Seykota script loaded
  - Donchian channels visible (blue lines)
  - Green triangle (long breakout signal)
  - Red stop line visible
  - Indicator window shows: N = 2.35
- **Highlights:**
  - Channels clear
  - Signals visible
  - N value readable

#### tradingview-donchian.png
- **Purpose:** Close-up of Donchian channels
- **State:**
  - Zoomed view of price action near Donchian bands
  - Price breaking above upper blue line
  - Green triangle marking breakout
- **Highlights:**
  - Breakout mechanism clear

#### tradingview-indicator.png
- **Purpose:** Show indicator window with N value
- **State:**
  - Bottom panel below chart
  - Indicator name: "Seykota / Turtle Core v2.1 + Date Range"
  - Display: N = 2.35 (or similar)
- **Highlights:**
  - N value large and readable

#### tradingview-pine-editor.png
- **Purpose:** Show Pine Editor with script code
- **State:**
  - Pine Editor panel at bottom
  - Ed-Seykota.pine script code visible (first 20-30 lines)
  - "Save" and "Add to Chart" buttons visible
- **Highlights:**
  - Editor location clear
  - Toolbar visible

#### tradingview-link.png
- **Purpose:** Show "Open in TradingView" link in TF-Engine
- **State:**
  - Scanner or Checklist screen
  - Ticker row (e.g., AAPL) with "Open in TradingView" link/button
  - Cursor hovering (if possible)
- **Highlights:**
  - Link clearly visible and clickable

---

### 11. Error States & Messages (3 screenshots)

#### error-database-locked.png
- **Purpose:** Show database error
- **State:**
  - Error message banner or modal:
    - "Database is locked"
    - "Close other instances of TF-Engine"
  - UI partially functional (read-only mode)
- **Highlights:**
  - Error prominent
  - Instructions clear

#### error-finviz-failed.png
- **Purpose:** Show FINVIZ scan failure
- **State:**
  - Scanner screen
  - Error message:
    - "FINVIZ scan failed: Request timeout"
    - "Check network connection and try again"
  - "Retry" button visible
- **Highlights:**
  - User-friendly error message

#### warning-heat-cap.png
- **Purpose:** Show heat cap warning (not error, just warning)
- **State:**
  - Heat Check screen showing:
    - Portfolio heat at 95% (yellow/orange zone)
    - Warning: "Approaching cap - add trades cautiously"
- **Highlights:**
  - Warning distinct from error

---

## Screenshot Capture Workflow

### Preparation

1. **Clean installation:**
   - Fresh install of TF-Engine
   - Initialize with clean database
   - Configure settings with example values

2. **Test data:**
   - Import test candidates
   - Create 3-5 test positions
   - Generate test decisions (GO/NO-GO)

3. **Browser setup:**
   - Chrome (latest version)
   - 1920×1080 resolution (or scale to this)
   - Zoom: 100%
   - Close unnecessary tabs/windows

### Capture Process

For each screenshot:

1. **Set up state:**
   - Navigate to screen
   - Fill forms / trigger states as specified
   - Verify display matches requirements

2. **Clean UI:**
   - Hide cursor (unless specified to show hover)
   - Remove personal data (if any leaked)
   - Check for typos in UI

3. **Capture:**
   - Win+Shift+S (Snipping Tool)
   - Select entire browser window or specific area
   - Save to: `docs/screenshots/[filename].png`

4. **Annotate (optional):**
   - Use ShareX or Greenshot to add:
     - Arrows pointing to key elements
     - Numbered callouts (1, 2, 3...)
     - Highlight boxes around important areas
   - Keep annotations minimal (only if helpful)

5. **Verify:**
   - Open screenshot
   - Check readability
   - Check file size (< 500 KB ideal, < 1 MB max)
   - If > 1 MB: Compress with TinyPNG or similar

---

## Post-Capture Tasks

### Organize Files

```
docs/screenshots/
├── dashboard-empty.png
├── dashboard-with-positions.png
├── dashboard-cooldowns.png
├── scanner-results.png
├── ...
└── README.md (this file)
```

### Update Documentation

1. **Verify all references:**
   ```bash
   cd docs/
   grep -r "screenshots/" *.md
   # Check all image paths exist
   ```

2. **Test links:**
   - Open each doc in Markdown preview
   - Check images render
   - Fix broken links

### Optimize Images

**If screenshots are large (> 500 KB):**

```bash
# Use TinyPNG, ImageOptim, or similar
# Or command-line tools:

# With ImageMagick:
convert input.png -quality 85 output.png

# With pngquant:
pngquant --quality=80-90 input.png -o output.png
```

**Target sizes:**
- Dashboard screenshots: 300-500 KB
- Detail screenshots (charts): 200-400 KB
- Icon/small screenshots: 50-150 KB

---

## Placeholder Status

**Current:** All screenshot references in documentation use placeholder paths.

**Before final release:**
- [ ] Capture all 40 screenshots
- [ ] Verify filenames match documentation
- [ ] Optimize file sizes
- [ ] Test all image links
- [ ] Add alt text for accessibility

---

## Accessibility Notes

When adding screenshots to documentation, include **alt text** for screen readers:

```markdown
![Dashboard with open positions](screenshots/dashboard-with-positions.png)

Alt text: Dashboard screen showing 5 open positions in a table, portfolio heat at 72% in green zone, and quick action buttons at the bottom.
```

**Alt text guidelines:**
- Describe content concisely (1-2 sentences)
- Mention key UI elements
- Don't start with "Screenshot of..." or "Image of..."

---

## Contact

**Questions about screenshot requirements:**
- See project maintainer
- Or file issue in GitHub: [your-repo]/issues

---

**Version:** 1.0.0
**Last Updated:** 2025-10-29
**Status:** Requirements documented, screenshots pending capture
