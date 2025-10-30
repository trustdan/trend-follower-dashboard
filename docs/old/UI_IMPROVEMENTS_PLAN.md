# TF-Engine GUI Improvements Plan

**Date:** October 29, 2025
**Status:** Identified - Ready to Implement

## Issues Identified

### 1. Theme & Visual Design
- ❌ No dark/light mode toggle
- ❌ Generic color scheme (no British Racing Green / pink theme)
- ❌ Inconsistent visual hierarchy
- ❌ Basic UI polish

### 2. Functionality Gaps
- ❌ "Edit Settings" button does nothing
- ❌ Scanner missing 4 of 5 presets
- ❌ No tooltips/help text on checklist items
- ❌ No explanations for position sizing methods
- ❌ Heat Check concept unclear

### 3. User Guidance
- ❌ No onboarding or first-run experience
- ❌ No inline help or explanations
- ❌ Missing context for trading rules

## Proposed Solutions

### Theme System (Priority: HIGH)

**Color Schemes:**

**Light Mode (Default):**
- Background: Warm off-white (#F5F5F0)
- Primary: British Racing Green (#004225)
- Accent: Forest Green (#228B22)
- Secondary: Dusty Rose (#D8A7B1) for warnings/highlights
- Cards: White (#FFFFFF)

**Dark Mode:**
- Background: Dark charcoal (#1E1E1E)
- Primary: British Racing Green (#004225)
- Accent: Light Green (#90EE90)
- Secondary: Soft Pink (#FFB3C1) for warnings/highlights
- Cards: Dark card (#282828)

**Banner Colors (Discipline Enforcement):**
- GREEN (#2E7D32) - GO, all gates passed
- YELLOW (#FFC107) - CAUTION, missing optional items
- RED (#D32F2F) - NO-GO, required gates failed

**Implementation:**
```go
// Add to AppState
type AppState struct {
    db          *storage.DB
    window      fyne.Window
    isDarkMode  bool  // NEW
}

// Theme toggle button in header/toolbar
toggleThemeBtn := widget.NewButton("🌙/☀️", func() {
    state.isDarkMode = !state.isDarkMode
    if state.isDarkMode {
        myApp.Settings().SetTheme(&tfTheme{variant: ThemeDark})
    } else {
        myApp.Settings().SetTheme(&tfTheme{variant: ThemeLight})
    }
})
```

---

### Scanner Improvements (Priority: HIGH)

**Add All 5 FINVIZ Presets:**

From `backend/internal/cli/interactive.go`:

```go
var finvizPresets = map[string]string{
    "TF-Breakout-Long":      "https://finviz.com/screener.ashx?v=111&p=d&s=ta_newhigh&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pa,ta_sma50_pa&o=-relativevolume",
    "TF-Momentum-Uptrend":   "https://finviz.com/screener.ashx?v=111&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pa,ta_sma50_pa&dr=y1&o=-marketcap",
    "TF-Unusual-Volume":     "https://finviz.com/screener.ashx?v=111&p=d&s=ta_unusualvolume&f=cap_largeover,sh_price_o20,ta_sma200_pa,ta_sma50_pa&o=-relativevolume",
    "TF-Breakdown-Short":    "https://finviz.com/screener.ashx?v=111&p=d&s=ta_newlow&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pb,ta_sma50_pb&o=-relativevolume",
    "TF-Momentum-Downtrend": "https://finviz.com/screener.ashx?v=111&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pb,ta_sma50_pb&dr=y1&o=-marketcap",
}
```

**Current GUI has only 3:**
- TF Breakout (Long)
- TF Breakout (Short)
- Volatile Longs

**Need to add:**
- TF-Momentum-Uptrend
- TF-Unusual-Volume
- TF-Momentum-Downtrend

**Update scanner.go preset buttons section.**

---

### Edit Settings Dialog (Priority: HIGH)

**Add functional settings editor:**

```go
func showSettingsDialog(state *AppState) {
    // Load current settings
    settings, _ := state.db.GetAllSettings()

    // Create form
    equityEntry := widget.NewEntry()
    equityEntry.SetText(settings["equity"])

    riskPctEntry := widget.NewEntry()
    riskPctEntry.SetText(settings["risk_pct"])

    portfolioCapEntry := widget.NewEntry()
    portfolioCapEntry.SetText(settings["portfolio_heat_cap"])

    bucketCapEntry := widget.NewEntry()
    bucketCapEntry.SetText(settings["bucket_heat_cap"])

    form := &widget.Form{
        Items: []*widget.FormItem{
            {Text: "Account Equity ($)", Widget: equityEntry},
            {Text: "Risk per Trade (%)", Widget: riskPctEntry},
            {Text: "Portfolio Heat Cap (%)", Widget: portfolioCapEntry},
            {Text: "Bucket Heat Cap (%)", Widget: bucketCapEntry},
        },
        OnSubmit: func() {
            // Save to database
            state.db.SetSetting("equity", equityEntry.Text)
            state.db.SetSetting("risk_pct", riskPctEntry.Text)
            state.db.SetSetting("portfolio_heat_cap", portfolioCapEntry.Text)
            state.db.SetSetting("bucket_heat_cap", bucketCapEntry.Text)

            // Refresh dashboard
            dialog.ShowInformation("Success", "Settings saved successfully", state.window)
        },
    }

    dialog.ShowForm("Edit Settings", "Save", "Cancel", form.Items, form.OnSubmit, func() {}, state.window)
}
```

---

### Checklist Help Text (Priority: MEDIUM)

**Add info icons/tooltips for each item:**

From `docs/anti-impulsivity.md`:

**Required Gates:**
1. **From Preset (SIG_REQ)** - "Ticker came from today's FINVIZ preset scan (55-bar breakout filter)"
2. **Trend Pass (RISK_REQ)** - "Trend confirmed: Long > 55-high OR Short < 55-low. Uses 2×N stop distance."
3. **Liquidity OK (OPT_REQ)** - "Options liquidity: bid-ask < 10% of mid, OI > 100, 60-90 DTE"
4. **TV Confirm (EXIT_REQ)** - "Exit plan confirmed: 10-bar opposite Donchian OR 2×N, whichever closer"
5. **Earnings OK (BEHAV_REQ)** - "2-minute cooloff passed, no intraday overrides, earnings blackout OK"

**Optional Quality Items:**
6. **Regime OK** - "Market regime favorable (e.g., SPY > 200SMA for longs)"
7. **No Chase** - "Entry not > 2N above 20-EMA (avoids chasing)"
8. **Journal OK** - "Trade plan documented: why now, profit targets, risk/reward"

**Implementation:**
```go
// Add info button next to each checkbox
infoBtn := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
    dialog.ShowInformation("From Preset",
        "Ticker came from today's FINVIZ preset scan...",
        state.window)
})

row := container.NewBorder(nil, nil, nil, infoBtn, fromPresetCheck)
```

---

### Position Sizing Explanations (Priority: MEDIUM)

**Add help text for each method:**

```go
methodExplanations := map[string]string{
    "Stock/ETF": `Van Tharp position sizing method:

1. Risk $ = Account × Risk%
2. Stop Distance = K × ATR (typically K=2)
3. Initial Stop = Entry - Stop Distance
4. Shares = floor(Risk$ / Stop Distance)
5. Actual Risk = Shares × Stop Distance

Example: $10,000 account, 0.75% risk, AAPL @ $180, ATR=$1.50, K=2
→ Risk=$75, Stop=$3, Shares=25

Pyramiding: Add every 0.5×N up to Max Units (typically 4)`,

    "Options (Delta-ATR)": `Options sizing using delta-adjusted ATR:

Per-unit risk adjusted for delta:
Contracts = Risk$ / (Delta × Stop Distance × 100)

Use for: Single calls/puts
Best when: Clear directional bias, high IV
Delta: Typically 0.60-0.80 for trending moves`,

    "Options (Contracts)": `Options sizing using max loss per contract:

Contracts = Risk$ / (Contract Price × 100)

Use for: Debit spreads, defined-risk strategies
Best when: Lower cost, limited downside
Max Loss = Debit paid per spread`
}

// Show explanation when method changes
methodSelect.OnChanged = func(method string) {
    explanationLabel.SetText(methodExplanations[method])
    // Show/hide relevant fields...
}
```

---

### Heat Check Education (Priority: HIGH)

**Add prominent explanation panel:**

```go
func buildHeatCheckExplanation() fyne.CanvasObject {
    return widget.NewCard(
        "What is Heat Check?",
        "Portfolio Risk Management",
        widget.NewRichTextFromMarkdown(`
**Heat = Total risk across all open positions**

### The Rules (from docs/anti-impulsivity.md):

1. **Portfolio Heat Cap:** 4% of equity
   - Example: $10,000 account → $400 max total risk

2. **Bucket Heat Cap:** 1.5% of equity per sector
   - Example: $10,000 account → $150 max in Tech/Comm

3. **Purpose:** Prevent concentration risk
   - Forces diversification across sectors
   - Prevents "all-in" on one trade/sector
   - Mechanical enforcement (no override)

### How It Works:

- Sum risk from all open positions
- Add risk from proposed new trade
- If total > cap → **REJECT** (no exceptions)
- If OK → **ALLOW** trade

**This is discipline enforcement, not flexibility.**
        `),
    )
}
```

---

### First-Run Experience (Priority: LOW)

**Show welcome dialog on first launch:**

```go
// In main.go, check for first run
if isFirstRun(db) {
    dialog.ShowCustom(
        "Welcome to TF-Engine!",
        "Get Started",
        buildWelcomeContent(),
        mainWindow,
    )
    setNotFirstRun(db)
}

func buildWelcomeContent() fyne.CanvasObject {
    return widget.NewRichTextFromMarkdown(`
# Welcome to TF-Engine

This is a **discipline enforcement system** for trend-following trading.

## Quick Start:

1. **Dashboard** - Set your account size and risk parameters
2. **Scanner** - Import FINVIZ candidates daily
3. **Checklist** - Evaluate trades (RED/YELLOW/GREEN banner)
4. **Position Sizing** - Calculate shares/contracts
5. **Heat Check** - Verify portfolio limits
6. **Trade Entry** - Final gates check before GO/NO-GO

## Philosophy:

This system **prevents impulsive trading** through:
- 5 Hard Gates (cannot bypass)
- Heat caps (no concentration)
- 2-minute cooloff
- Mechanical exits

**The value is in what it prevents, not what it allows.**

See docs/anti-impulsivity.md for full details.
    `)
}
```

---

## Implementation Priority

### Phase 1 (Critical UX)
1. ✅ Theme system with British Racing Green + Dark/Light toggle
2. ✅ Scanner: Add all 5 FINVIZ presets
3. ✅ Edit Settings dialog (functional)
4. ✅ Heat Check explanation panel

### Phase 2 (User Guidance)
5. ✅ Checklist item tooltips/help
6. ✅ Position Sizing method explanations
7. ✅ Welcome dialog (first run)

### Phase 3 (Polish)
8. ✅ Improve visual hierarchy
9. ✅ Better spacing/padding
10. ✅ Icons and visual cues

---

## Files to Modify

1. `ui/theme.go` - Complete theme system
2. `ui/main.go` - Add theme toggle, first-run check
3. `ui/dashboard.go` - Wire up Edit Settings
4. `ui/scanner.go` - Add missing presets
5. `ui/checklist.go` - Add tooltips
6. `ui/position_sizing.go` - Add explanations
7. `ui/heat_check.go` - Add education panel

---

## Expected Outcome

**Before:**
- Basic GUI, unclear purpose
- Missing features
- No guidance

**After:**
- Professional discipline enforcement system
- Clear purpose and guidance
- Beautiful British Racing Green + Pink theme
- Dark/Light modes
- All backend features exposed
- Inline help everywhere

**This transforms the GUI from "basic prototype" to "production-ready trading discipline system".**
