# Phase 4 - Step 21: TradingView Integration

**TF = Trend Following** - Systematic trading discipline enforcement system

**Phase:** 4 - Calendar & Polish
**Step:** 21 of 28
**Duration:** 1 day
**Dependencies:** Step 20 complete (Calendar View functional)

---

## Objectives

Implement seamless TradingView integration for quick chart access. Add "Open in TradingView" buttons next to tickers throughout the app, allowing traders to verify Donchian breakout signals using the Ed-Seykota.pine script.

**Purpose:** Bridge the gap between systematic screening (FINVIZ) and signal verification (TradingView), making the workflow effortless.

---

## Success Criteria

- [ ] "Open in TradingView" button/link component created and reusable
- [ ] Links added to Candidates list
- [ ] Links added to Dashboard positions
- [ ] Links added to Checklist form
- [ ] Links added to Calendar cells (optional)
- [ ] Clicking link opens new browser tab to TradingView chart
- [ ] URL format: `https://www.tradingview.com/chart/?symbol={ticker}`
- [ ] URL template is customizable in settings (for advanced users)
- [ ] Component styled with gradient accent and TradingView icon
- [ ] Documentation includes setup guide for Ed-Seykota.pine script

---

## Prerequisites

**Frontend Requirements:**
- Navigation functional
- Candidates list displays tickers
- Dashboard displays open positions
- Checklist form accepts ticker input
- Theme system working (day/night mode)

**Optional:**
- Settings screen for URL template customization

---

## Implementation Plan

### Step 1: TradingViewLink Component (1 hour)

**File:** `ui/src/lib/components/TradingViewLink.svelte`

Create a reusable component that constructs and opens TradingView chart URLs.

```svelte
<script lang="ts">
    import { settings } from '$lib/stores/settings';

    export let ticker: string;
    export let variant: 'button' | 'icon' | 'text' = 'button';

    // Default TradingView URL template
    // Users can customize in settings to add chart layout, intervals, etc.
    let urlTemplate = $settings.tradingViewUrlTemplate || 'https://www.tradingview.com/chart/?symbol={ticker}';

    function openChart() {
        const url = urlTemplate.replace('{ticker}', ticker.toUpperCase());
        window.open(url, '_blank', 'noopener,noreferrer');
    }
</script>

{#if variant === 'button'}
    <button class="tv-button" on:click={openChart} title="Open {ticker} in TradingView">
        <svg class="tv-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <path d="M3 3h18v18H3z" stroke-width="2"/>
            <path d="M9 9l6 6M15 9l-6 6" stroke-width="2"/>
        </svg>
        <span>Open in TradingView</span>
    </button>
{:else if variant === 'icon'}
    <button class="tv-icon-only" on:click={openChart} title="Open {ticker} in TradingView">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <path d="M3 3h18v18H3z" stroke-width="2"/>
            <path d="M9 9l6 6M15 9l-6 6" stroke-width="2"/>
        </svg>
    </button>
{:else}
    <a href="#" class="tv-text-link" on:click|preventDefault={openChart}>
        {ticker} ↗
    </a>
{/if}

<style>
    .tv-button {
        display: inline-flex;
        align-items: center;
        gap: var(--space-2);
        padding: 8px 16px;
        background: var(--gradient-blue);
        color: white;
        border: none;
        border-radius: 6px;
        font-size: var(--text-sm);
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s ease;
        box-shadow: 0 2px 4px rgba(59, 130, 246, 0.2);
    }

    .tv-button:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 8px rgba(59, 130, 246, 0.3);
    }

    .tv-icon {
        width: 16px;
        height: 16px;
    }

    .tv-icon-only {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 32px;
        height: 32px;
        padding: 6px;
        background: var(--bg-secondary);
        border: 1px solid var(--border-color);
        border-radius: 6px;
        cursor: pointer;
        transition: all 0.15s ease;
    }

    .tv-icon-only:hover {
        background: var(--gradient-blue);
        color: white;
        border-color: transparent;
    }

    .tv-icon-only svg {
        width: 20px;
        height: 20px;
    }

    .tv-text-link {
        color: var(--border-focus);
        text-decoration: none;
        font-weight: 500;
        cursor: pointer;
        transition: color 0.15s ease;
    }

    .tv-text-link:hover {
        color: #2563EB;
        text-decoration: underline;
    }
</style>
```

---

### Step 2: Add to Candidates List (30 min)

**File:** `ui/src/lib/components/CandidateList.svelte` (or `ui/src/routes/scanner/+page.svelte`)

Add TradingView link to each candidate row:

```svelte
<script lang="ts">
    import TradingViewLink from '$lib/components/TradingViewLink.svelte';
    // ... existing imports

    export let candidates: any[];
</script>

<table class="candidates-table">
    <thead>
        <tr>
            <th>Select</th>
            <th>Ticker</th>
            <th>Sector</th>
            <th>Last Close</th>
            <th>Volume</th>
            <th>Chart</th>
        </tr>
    </thead>
    <tbody>
        {#each candidates as candidate}
            <tr>
                <td>
                    <input type="checkbox" bind:checked={candidate.selected} />
                </td>
                <td>{candidate.ticker}</td>
                <td>{candidate.sector}</td>
                <td>${candidate.last_close.toFixed(2)}</td>
                <td>{formatVolume(candidate.volume)}</td>
                <td>
                    <TradingViewLink ticker={candidate.ticker} variant="icon" />
                </td>
            </tr>
        {/each}
    </tbody>
</table>

<style>
    /* ... existing styles */
</style>
```

---

### Step 3: Add to Dashboard Positions (30 min)

**File:** `ui/src/routes/dashboard/+page.svelte` or `ui/src/lib/components/PositionList.svelte`

Add TradingView link to each open position:

```svelte
<script lang="ts">
    import TradingViewLink from '$lib/components/TradingViewLink.svelte';
    // ... existing imports

    export let positions: any[];
</script>

<div class="positions-section">
    <h2>Open Positions</h2>
    <table class="positions-table">
        <thead>
            <tr>
                <th>Ticker</th>
                <th>Entry</th>
                <th>Current Stop</th>
                <th>Risk</th>
                <th>Days Held</th>
                <th>Chart</th>
            </tr>
        </thead>
        <tbody>
            {#each positions as pos}
                <tr>
                    <td>{pos.ticker}</td>
                    <td>${pos.entry.toFixed(2)}</td>
                    <td>${pos.current_stop.toFixed(2)}</td>
                    <td>${pos.risk.toFixed(2)}</td>
                    <td>{pos.days_held}</td>
                    <td>
                        <TradingViewLink ticker={pos.ticker} variant="icon" />
                    </td>
                </tr>
            {/each}
        </tbody>
    </table>
</div>

<style>
    /* ... existing styles */
</style>
```

---

### Step 4: Add to Checklist Form (30 min)

**File:** `ui/src/routes/checklist/+page.svelte`

Add TradingView button next to ticker input:

```svelte
<script lang="ts">
    import TradingViewLink from '$lib/components/TradingViewLink.svelte';
    // ... existing imports

    let ticker = '';
</script>

<div class="checklist-page">
    <!-- Banner component -->
    <Banner state={bannerState} message={bannerMessage} />

    <form class="checklist-form">
        <div class="form-row">
            <div class="form-group">
                <label for="ticker">Ticker</label>
                <input
                    type="text"
                    id="ticker"
                    bind:value={ticker}
                    placeholder="AAPL"
                    on:input={handleTickerChange}
                />
            </div>
            {#if ticker}
                <div class="tv-button-wrapper">
                    <TradingViewLink {ticker} variant="button" />
                </div>
            {/if}
        </div>

        <!-- Rest of form fields -->
    </form>
</div>

<style>
    .form-row {
        display: flex;
        gap: var(--space-4);
        align-items: flex-end;
    }

    .form-group {
        flex: 1;
    }

    .tv-button-wrapper {
        margin-bottom: var(--space-1);
    }
</style>
```

---

### Step 5: Add to Calendar Cells (Optional, 30 min)

**File:** `ui/src/lib/components/CalendarCell.svelte`

Make ticker tags clickable to open TradingView:

```svelte
<script lang="ts">
    import TradingViewLink from '$lib/components/TradingViewLink.svelte';
    export let cells: any[] = [];

    function openChart(ticker: string) {
        const url = `https://www.tradingview.com/chart/?symbol=${ticker}`;
        window.open(url, '_blank', 'noopener,noreferrer');
    }
</script>

<div class="calendar-cell">
    {#if cells.length > 0}
        <div class="cell-content">
            {#each cells as cell}
                <button
                    class="ticker-tag"
                    on:click={() => openChart(cell.ticker)}
                    title="Click to open {cell.ticker} in TradingView"
                >
                    {cell.ticker} ↗
                </button>
            {/each}
        </div>
    {/if}
</div>

<style>
    .ticker-tag {
        font-size: var(--text-sm);
        font-weight: 500;
        padding: 4px 8px;
        background: var(--bg-primary);
        border: 1px solid var(--border-color);
        border-radius: 4px;
        cursor: pointer;
        transition: all 0.15s ease;
    }

    .ticker-tag:hover {
        background: var(--gradient-blue);
        color: white;
        border-color: transparent;
    }
</style>
```

---

### Step 6: Settings Integration (Optional, 1 hour)

**File:** `ui/src/routes/settings/+page.svelte`

Allow users to customize the TradingView URL template:

```svelte
<script lang="ts">
    import { settings } from '$lib/stores/settings';

    let tvUrlTemplate = $settings.tradingViewUrlTemplate || 'https://www.tradingview.com/chart/?symbol={ticker}';

    function saveSettings() {
        $settings.tradingViewUrlTemplate = tvUrlTemplate;
        // Save to backend
        api.settings.update({ tradingViewUrlTemplate: tvUrlTemplate });
    }
</script>

<div class="settings-page">
    <h2>TradingView Integration</h2>

    <div class="form-group">
        <label for="tv-url">TradingView URL Template</label>
        <input
            type="text"
            id="tv-url"
            bind:value={tvUrlTemplate}
            placeholder="https://www.tradingview.com/chart/?symbol={ticker}"
        />
        <p class="help-text">
            Use <code>{ticker}</code> as placeholder. Example with chart layout:
            <code>https://www.tradingview.com/chart/XXX/?symbol={ticker}&interval=D</code>
        </p>
    </div>

    <button class="save-button" on:click={saveSettings}>
        Save Settings
    </button>
</div>

<style>
    .help-text {
        font-size: var(--text-sm);
        color: var(--text-secondary);
        margin-top: var(--space-2);
    }

    code {
        background: var(--bg-tertiary);
        padding: 2px 6px;
        border-radius: 4px;
        font-family: monospace;
    }
</style>
```

**Backend:** Add `tradingview_url_template` column to `settings` table (optional).

---

### Step 7: Documentation (1 hour)

**File:** `docs/TRADINGVIEW_SETUP.md` (new)

Create comprehensive setup guide:

```markdown
# TradingView Integration Setup Guide

## Overview

The TF-Engine GUI integrates with TradingView to allow quick chart access for signal verification using the Ed-Seykota.pine script.

## How It Works

1. Click "Open in TradingView" next to any ticker
2. TradingView chart opens in new browser tab
3. Manually verify 55-bar Donchian breakout using Ed-Seykota.pine script
4. Note entry price and N (ATR) value from script
5. Return to TF-Engine and complete checklist

## Ed-Seykota.pine Script Setup

### Step 1: Install Pine Script

1. Open TradingView and navigate to any chart
2. Click "Pine Editor" at the bottom of the screen
3. Copy the contents of `reference/Ed-Seykota.pine`
4. Paste into Pine Editor
5. Click "Add to Chart"

### Step 2: Configure Script Parameters

**Core Parameters (matching TF-Engine):**
- `entryLen = 55` - Donchian entry lookback
- `exitLen = 10` - Donchian exit lookback
- `nLen = 20` - ATR length for N
- `stopN = 2.0` - Initial stop distance (2×N)
- `addStepN = 0.5` - Add every 0.5×N
- `maxUnits = 4` - Max units per position
- `riskPct = 1.0` - Risk % per unit (adjust to match your account settings)

**Optional Filters:**
- `useMarket = false` - Enable SPY > 200 SMA filter if desired
- `minVol = 0` - Minimum volume filter (0 = disabled)

### Step 3: Save as Default (Optional)

To automatically load the script on every chart:

1. Click "Indicators" → "My Scripts" → "Ed-Seykota"
2. Click the gear icon (settings)
3. Check "Add to favorites"
4. TradingView doesn't support auto-load, so add manually to each chart

### Step 4: Verify Signal

When viewing a candidate ticker:

1. **Visual Check:** Look at the chart
   - For longs: Did price close above the 55-bar high (blue line)?
   - For shorts: Did price close below the 55-bar low (red line)?

2. **Read N Value:** Script displays "N = X.XX" on chart
   - This is the current ATR(20) value
   - Enter this into TF-Engine checklist

3. **Note Entry Price:**
   - Use current close or your planned limit order price
   - Enter this into TF-Engine checklist

4. **Check Stop Level:** Script plots stop as dotted line
   - For longs: Entry - 2×N
   - For shorts: Entry + 2×N

## Custom URL Templates

You can customize the TradingView URL in Settings:

**Default:**
```
https://www.tradingview.com/chart/?symbol={ticker}
```

**With Daily Timeframe:**
```
https://www.tradingview.com/chart/?symbol={ticker}&interval=D
```

**With Saved Chart Layout:**
1. Set up your chart with Ed-Seykota.pine script and preferred layout
2. Click "Share" → "Copy link"
3. Replace the ticker symbol with `{ticker}` placeholder
4. Paste into TF-Engine settings

Example:
```
https://www.tradingview.com/chart/XXX/?symbol={ticker}&interval=D&theme=dark
```

## Workflow Example

1. Run FINVIZ scan → Import 12 candidates
2. Click "Open in TradingView" for AAPL
3. TradingView opens with AAPL chart
4. Verify: Close > 55-bar high ✓
5. Read: N = 2.35
6. Entry: $180.50 (current close)
7. Stop: $175.80 (shown on chart)
8. Return to TF-Engine
9. Fill checklist with AAPL, 180.50, 2.35
10. Continue with position sizing

## Troubleshooting

**Problem:** TradingView chart doesn't open
**Solution:** Check popup blocker settings; allow popups from TF-Engine

**Problem:** Script doesn't display on chart
**Solution:** Click "Indicators" → "My Scripts" → "Ed-Seykota" to add it

**Problem:** N value doesn't match TF-Engine calculation
**Solution:** Verify ATR length is 20 in both script and TF-Engine

**Problem:** Custom URL template doesn't work
**Solution:** Ensure `{ticker}` placeholder is present and spelled correctly
```

---

## Testing Checklist

### Component Tests

- [ ] TradingViewLink component renders correctly
- [ ] Button variant displays icon and text
- [ ] Icon variant displays icon only
- [ ] Text variant displays as link
- [ ] Click opens new tab (not current tab)
- [ ] URL is correctly formatted with ticker symbol
- [ ] Component works with both uppercase and lowercase tickers

### Integration Tests

- [ ] Candidates list shows TradingView icon for each ticker
- [ ] Dashboard positions show TradingView icon
- [ ] Checklist form shows TradingView button when ticker entered
- [ ] Calendar cells open TradingView on ticker click (if implemented)
- [ ] All links open in new tab with `noopener,noreferrer`
- [ ] Custom URL template (if implemented) replaces `{ticker}` correctly

### Browser Tests

- [ ] Links work in Chrome
- [ ] Links work in Firefox
- [ ] Links work in Edge
- [ ] Popup blocker doesn't prevent opening (or shows clear message)

### Documentation Tests

- [ ] TradingView setup guide is clear and complete
- [ ] Screenshots show script installation (if added)
- [ ] Workflow example matches actual UI
- [ ] Troubleshooting section addresses common issues

---

## Troubleshooting

**Problem:** Popup blocker prevents TradingView from opening
**Solution:** Show notification: "Allow popups to open TradingView charts"

**Problem:** TradingView URL doesn't include chart layout
**Solution:** Document custom URL template setup in settings

**Problem:** Icon doesn't display correctly
**Solution:** Check SVG path and ensure stroke color inherits from parent

**Problem:** Click handler doesn't fire
**Solution:** Ensure `on:click` is properly bound and not prevented by parent

---

## Documentation Requirements

- [ ] Create `docs/TRADINGVIEW_SETUP.md` with comprehensive guide
- [ ] Add screenshots of Ed-Seykota.pine script on TradingView
- [ ] Update User Guide with TradingView integration section
- [ ] Document custom URL template feature (if implemented)
- [ ] Update `docs/PROGRESS.md` with completion status

---

## Next Steps

After completing Step 21:
1. Proceed to **Step 22: UI Polish & Refinements**
2. Test TradingView integration with real workflow
3. Gather user feedback on ease of use
4. Consider adding TradingView widget embedding (future enhancement)

---

**Estimated Completion Time:** 1 day
**Phase 4 Progress:** 2 of 4 steps complete
**Overall Progress:** 21 of 28 steps complete (75%)

---

**End of Step 21**
