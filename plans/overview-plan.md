# TF-Engine GUI: Overview Plan
## The Forest Over the Trees

**TF = Trend Following** - Systematic trading using Donchian breakouts, ATR-based sizing, and mechanical exits

**Created:** 2025-10-29
**Purpose:** Strategic foundation for the entire GUI trading platform
**Philosophy:** Anti-impulsivity through systematic discipline enforcement
**Status:** 🎯 Planning Phase

---

## Table of Contents

1. [Vision & Philosophy](#vision--philosophy)
2. [System Architecture](#system-architecture)
3. [User Workflow](#user-workflow)
4. [Technical Components](#technical-components)
5. [Implementation Strategy](#implementation-strategy)
6. [Proof-of-Concept Approach](#proof-of-concept-approach)
7. [Behavioral Specifications (Gherkin)](#behavioral-specifications-gherkin)
8. [Success Criteria](#success-criteria)
9. [Risk Management](#risk-management)
10. [Next Steps](#next-steps)

---

## Vision & Philosophy

### The Purpose

This is **not** a trading platform. This is a **discipline enforcement system** that happens to be implemented as software.

**What is Trend Following?**
- Follow existing price trends (the "tide") rather than predict reversals
- Enter on breakouts (new highs/lows over N-bar period)
- Exit mechanically (opposite breakout or ATR-based stop)
- Add to winners (pyramid), cut losers quickly
- No opinion, no prediction - just systematic execution

From `docs/anti-impulsivity.md`:

> **Trade the tide, not the splash.** A breakout and a mechanical exit are the core.

**The TF-Engine (Trend Following Engine)** implements this philosophy through:
- 55-bar Donchian breakouts (Ed Seykota's System-2)
- ATR-based position sizing (Van Tharp method)
- 10-bar opposite Donchian exits (mechanical)
- 0.5×N pyramiding (add to winners)
- 2×N protective stops (cut losers)

The system exists to make impulsive trading **impossible** while making systematic, trend-following trading **effortless**.

### Core Principles (Non-Negotiable)

1. **Friction where it matters** - Hard gates for signal, risk, liquidity, exit, behavior
2. **Nudge for better trades** - Optional quality items affect score, not permission
3. **Immediate feedback** - Large 3-state banner (RED/YELLOW/GREEN) updates live
4. **Journal while you decide** - One-click logging of full trade plan
5. **Calendar awareness** - Rolling 10-week sector view prevents basket crowding
6. **No backdoors** - Cannot bypass gates, skip cooldowns, or override caps

### The 5 Hard Gates (Enforced by Backend)

Every trade must pass ALL gates before execution. These are **immutable**.

```gherkin
Given a trader wants to execute a trade
Then ALL of the following gates must be GREEN:

Gate 1: Signal Confirmed
  - 55-bar Donchian breakout (long: close > donHi[1], short: close < donLo[1])
  - Verified against TradingView chart using Ed-Seykota.pine script
  - No subjective pattern analysis allowed

Gate 2: Risk/Size Calculated
  - Per-unit risk = equity × risk% (typically 0.75-1.0%)
  - Stop distance = 2.0 × N (where N = ATR(20))
  - Initial stop = entry ± (2.0 × N)
  - Shares/contracts = floor(unitRisk$ ÷ stopDistance)
  - Add-on levels = entry + (k × 0.5 × N) for k=1,2,3

Gate 3: Options Requirements Met (if applicable)
  - DTE: 60-90 days at entry
  - Roll/close trigger: ~21 DTE
  - Liquidity: bid-ask spread < 10% of mid
  - Open interest: > 100 contracts
  - If stock trade: liquidity check via avg volume

Gate 4: Exits Defined
  - Primary: 10-bar opposite Donchian (donLo[1] for longs, donHi[1] for shorts)
  - Secondary: Initial stop = entry ± (2.0 × N)
  - Exit trigger: CLOSER of Donchian or 2×N stop
  - No discretionary exits allowed

Gate 5: Behavior Constraints Honored
  - 2-minute cool-off elapsed since checklist evaluation
  - Ticker not on sector cooldown (from recent loss)
  - Portfolio heat < 4.0% of equity
  - Sector bucket heat < 1.5% of equity
  - No intraday overrides permitted
```

### Banner System (Visual Discipline Enforcement)

**Size:** Minimum 20% of screen height. **Cannot be missed.**

**Visual Design:** Sleek, modern gradients with smooth transitions. No flat colors.

```
┌─────────────────────────────────────────────────┐
│  ╔════════════════════════════════════════╗    │
│  ║                                        ║    │
│  ║    🛑  DO NOT TRADE  🛑                ║    │
│  ║    One or more REQUIRED gates failed   ║    │
│  ║                                        ║    │
│  ╚════════════════════════════════════════╝    │
│  [Gradient: Deep Red → Dark Crimson]           │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│  ╔════════════════════════════════════════╗    │
│  ║                                        ║    │
│  ║    ⚠️  CAUTION  ⚠️                     ║    │
│  ║    Quality score below threshold       ║    │
│  ║                                        ║    │
│  ╚════════════════════════════════════════╝    │
│  [Gradient: Amber → Golden Yellow]             │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│  ╔════════════════════════════════════════╗    │
│  ║                                        ║    │
│  ║    ✓  OK TO TRADE  ✓                  ║    │
│  ║    All gates pass • Quality met        ║    │
│  ║                                        ║    │
│  ╚════════════════════════════════════════╝    │
│  [Gradient: Emerald → Forest Green]            │
└─────────────────────────────────────────────────┘
```

**State Logic:**
- RED: Any required gate fails → **Absolutely no trade**
- YELLOW: All required pass, quality score < threshold (default 3.0) → **Caution**
- GREEN: All required pass, quality score ≥ threshold → **Proceed**

**Gradient Specifications:**
- **RED Banner:** Linear gradient from `#DC2626` (deep red) to `#991B1B` (dark crimson)
- **YELLOW Banner:** Linear gradient from `#F59E0B` (amber) to `#FBBF24` (golden)
- **GREEN Banner:** Linear gradient from `#10B981` (emerald) to `#059669` (forest green)
- **Transitions:** Smooth 0.3s ease-in-out animation when state changes
- **Text:** Large, bold, white text with subtle shadow for readability
- **Border:** Subtle rounded corners with soft glow effect in banner color

---

## System Architecture

### High-Level Architecture

```
┌───────────────────────────────────────────────────────────────┐
│                    DESKTOP APPLICATION                         │
│                                                                │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │              SVELTE FRONTEND (Browser/WebView)           │ │
│  │                                                          │ │
│  │  Components:                                             │ │
│  │  - Dashboard (Positions, Candidates, Heat)               │ │
│  │  - FINVIZ Scanner (One-click daily scan)                 │ │
│  │  - Candidate Import (Review & approve tickers)           │ │
│  │  - Checklist (Required gates + quality items)            │ │
│  │  - Position Sizer (ATR-based Van Tharp calculations)     │ │
│  │  - Heat Monitor (Portfolio 4% / Bucket 1.5% caps)        │ │
│  │  - Trade Entry (Final 5-gate validation)                 │ │
│  │  - Calendar View (10-week sector diversification)        │ │
│  │  - TradingView Integration (Open charts w/ script)       │ │
│  │                                                          │ │
│  │  State Management: Svelte stores + localStorage          │ │
│  │  UI Framework: SvelteKit (static adapter)                │ │
│  │  Styling: TailwindCSS or similar (large, obvious UI)     │ │
│  └──────────────────────────────────────────────────────────┘ │
│                           │                                    │
│                           │ HTTP REST API (JSON)               │
│                           │                                    │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │                 GO BACKEND (tf-engine)                   │ │
│  │                                                          │ │
│  │  ┌────────────────┐  ┌────────────────┐                │ │
│  │  │  HTTP Server   │  │  Static Assets │                │ │
│  │  │  (Gin/Chi)     │  │  (Embedded)    │                │ │
│  │  └────────────────┘  └────────────────┘                │ │
│  │                                                          │ │
│  │  ┌────────────────────────────────────────────────────┐ │ │
│  │  │              DOMAIN LOGIC (✅ Complete)            │ │ │
│  │  │                                                    │ │ │
│  │  │  - Position Sizing (stock, opt-delta, opt-contracts) │ │
│  │  │  - Checklist Evaluation (6 required + 4 optional)  │ │ │
│  │  │  - Heat Management (portfolio & sector caps)       │ │ │
│  │  │  - Gate Enforcement (all 5 gates)                  │ │ │
│  │  │  - Settings Management (equity, risk%, caps)       │ │ │
│  │  │  - Candidate Management (FINVIZ imports)           │ │ │
│  │  │  - Cooldown Tracking (sector & ticker)             │ │ │
│  │  └────────────────────────────────────────────────────┘ │ │
│  │                                                          │ │
│  │  ┌────────────────────────────────────────────────────┐ │ │
│  │  │          STORAGE (SQLite - ✅ Complete)            │ │ │
│  │  │                                                    │ │ │
│  │  │  Tables: settings, positions, decisions,          │ │ │
│  │  │          candidates, cooldowns, evaluations       │ │ │
│  │  └────────────────────────────────────────────────────┘ │ │
│  │                                                          │ │
│  │  ┌────────────────────────────────────────────────────┐ │ │
│  │  │        SCRAPE (FINVIZ - ✅ Complete)               │ │ │
│  │  │                                                    │ │ │
│  │  │  - Parse FINVIZ screener pages                     │ │ │
│  │  │  - Extract ticker symbols                          │ │ │
│  │  │  - Support presets (TF_BREAKOUT_LONG, etc.)        │ │ │
│  │  └────────────────────────────────────────────────────┘ │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                │
│  Database: trading.db (SQLite, single file)                   │
│  Deployment: Single binary (app.exe on Windows)               │
│  Opens: Default browser to localhost:8080 OR embedded webview │
└───────────────────────────────────────────────────────────────┘
```

### Technology Stack

**Backend (✅ Complete - tf-engine):**
- Language: Go 1.24+
- HTTP Server: Standard library or Gin/Chi
- Database: SQLite via mattn/go-sqlite3
- Web Scraping: golang.org/x/net
- Logging: Structured logging with levels (DEBUG, INFO, WARN, ERROR)
  - Log to file: `logs/tf-engine.log` (rotated daily)
  - Log to console: Configurable verbosity
  - Correlation IDs for request tracking
  - Performance metrics (request duration, DB query time)
  - Feature usage tracking (which features are used, how often)
  - Error tracking (full stack traces, context)
- Testing: Go standard testing + testify

**Frontend (🚧 To Build):**
- Language: TypeScript + Svelte
- Framework: SvelteKit with static adapter
- UI Library: Custom components (anti-impulsivity focused)
- Styling: TailwindCSS + custom CSS for gradients and animations
- Color System: Day/night mode with smooth theme toggle
- Design Language: Modern, sleek, gradient-heavy (no flat colors)
- State: Svelte stores + localStorage for theme preference
- Logging: Browser console logging with levels
  - User actions tracked (button clicks, form submissions, navigation)
  - API calls logged (request/response, timing, errors)
  - Component lifecycle events (mount, unmount, errors)
  - Performance metrics (render time, interaction delays)
  - Feature usage analytics (which screens visited, how often)
- Build: Vite (via SvelteKit)
- Output: Static files in `ui/build/`

**Integration:**
- Frontend → Backend: HTTP REST API (JSON)
- Backend serves frontend: Go `embed` package for static files
- Deployment: Single binary (includes both frontend and backend)
- Desktop: Opens default browser OR uses go-webview for native feel

**Development Environment:**
- OS: Linux (WSL2/Kali) for development
- Build Target: Windows .exe (cross-compiled)
- Git: Windows only (per `1._RULES.md—Operating_Rules_for_This_Project-(Claude_Code).md`)
- Handoff: Export zip → Windows Git repo
- Logging: Comprehensive logging for all features (debug, feature evaluation, performance)

---

## Visual Design Philosophy

### Modern, Sleek, and Appealing

**Core Principle:** The interface should be a pleasure to use. Beautiful design encourages regular engagement with the systematic process.

**Design Language:**
- **Gradients over flat colors** - Every major UI element uses smooth gradients
- **Smooth transitions** - 0.3s ease-in-out for all state changes
- **Subtle shadows** - Depth and dimension without clutter
- **Rounded corners** - Soft, approachable feel (border-radius: 8-12px)
- **Whitespace** - Generous padding and margins for breathing room
- **Typography** - Clear hierarchy with modern sans-serif fonts (Inter, SF Pro, Segoe UI)

### Color System

**Day Mode (Light Theme):**
```css
/* Background */
--bg-primary: #FFFFFF         /* Main background */
--bg-secondary: #F9FAFB       /* Cards, panels */
--bg-tertiary: #F3F4F6        /* Hover states */

/* Text */
--text-primary: #111827       /* Main text */
--text-secondary: #6B7280     /* Subtext */
--text-tertiary: #9CA3AF      /* Muted text */

/* Gradients */
--gradient-red: linear-gradient(135deg, #DC2626 0%, #991B1B 100%)
--gradient-yellow: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%)
--gradient-green: linear-gradient(135deg, #10B981 0%, #059669 100%)
--gradient-blue: linear-gradient(135deg, #3B82F6 0%, #1D4ED8 100%)
--gradient-purple: linear-gradient(135deg, #8B5CF6 0%, #6D28D9 100%)

/* Borders */
--border-color: #E5E7EB
--border-focus: #3B82F6
```

**Night Mode (Dark Theme):**
```css
/* Background */
--bg-primary: #0F172A         /* Main background (slate-900) */
--bg-secondary: #1E293B       /* Cards, panels (slate-800) */
--bg-tertiary: #334155        /* Hover states (slate-700) */

/* Text */
--text-primary: #F1F5F9       /* Main text (slate-100) */
--text-secondary: #CBD5E1     /* Subtext (slate-300) */
--text-tertiary: #94A3B8      /* Muted text (slate-400) */

/* Gradients (slightly muted for dark mode) */
--gradient-red: linear-gradient(135deg, #DC2626 0%, #7F1D1D 100%)
--gradient-yellow: linear-gradient(135deg, #F59E0B 0%, #B45309 100%)
--gradient-green: linear-gradient(135deg, #10B981 0%, #047857 100%)
--gradient-blue: linear-gradient(135deg, #3B82F6 0%, #1E40AF 100%)
--gradient-purple: linear-gradient(135deg, #8B5CF6 0%, #5B21B6 100%)

/* Borders */
--border-color: #334155
--border-focus: #60A5FA
```

**Theme Toggle:**
- Icon: Sun ☀️ (day mode) / Moon 🌙 (night mode)
- Location: Top-right corner of header
- Transition: Smooth 0.3s fade on all color changes
- Persistence: Save preference to localStorage
- Default: Detect system preference (`prefers-color-scheme`)

### Component Design Guidelines

**Buttons:**
```css
/* Primary Action (e.g., "SAVE GO DECISION") */
- Background: Gradient (green for GO, red for NO-GO, blue for neutral)
- Padding: 16px 32px
- Font: Bold, 18px
- Border-radius: 8px
- Shadow: 0 4px 6px rgba(0,0,0,0.1)
- Hover: Lift effect (translateY(-2px)) + deeper shadow
- Disabled: Grayscale filter + reduced opacity

/* Secondary Action */
- Background: Transparent
- Border: 2px solid with gradient border (via pseudo-element)
- Same size and hover effects
```

**Cards/Panels:**
```css
- Background: --bg-secondary
- Border: 1px solid --border-color
- Border-radius: 12px
- Shadow: 0 1px 3px rgba(0,0,0,0.1)
- Padding: 24px
- Hover: Subtle lift + shadow increase (for clickable cards)
```

**Forms:**
```css
/* Input Fields */
- Background: --bg-primary
- Border: 1px solid --border-color
- Border-radius: 6px
- Padding: 12px 16px
- Focus: Border color → --border-focus + subtle glow

/* Checkboxes */
- Custom styled with gradient when checked
- Size: 24px × 24px
- Check icon with smooth animation

/* Dropdowns */
- Same styling as inputs
- Chevron icon with smooth rotation on open
```

**Banner (The Star of the Show):**
```css
- Height: 20% of viewport height (minimum 150px)
- Width: 90% of container (centered)
- Margin: 32px auto
- Border-radius: 16px
- Background: Gradient (changes based on state)
- Shadow: 0 8px 16px with banner color (alpha 0.3)
- Text: 36px bold, white, centered
- Animation: Pulse effect on state change
- Glow: Subtle outer glow in banner color
```

**Tables:**
```css
- Header: Gradient background (subtle blue)
- Rows: Alternating --bg-primary and --bg-secondary
- Hover: Highlight row with --bg-tertiary
- Borders: Minimal, using --border-color
- Cell padding: 16px
```

### Animation Guidelines

**State Transitions:**
- Banner color change: 0.3s ease-in-out
- Page navigation: 0.2s slide-in from right
- Modal appearance: 0.2s fade-in + scale(0.95 → 1.0)
- Loading spinner: Continuous rotation with gradient trail

**Interactions:**
- Button hover: 0.15s ease-out
- Input focus: 0.1s ease-in
- Checkbox check: 0.2s spring animation

**Micro-interactions:**
- Checkmark appears with checkmark draw animation
- Success notification slides in from top
- Error notification shakes slightly on appear

### Icon Usage

**Icon Library:** Lucide Icons or Heroicons (clean, modern, consistent)

**Sizes:**
- Small: 16px (inline with text)
- Medium: 24px (buttons, navigation)
- Large: 32px (headers, feature icons)
- XL: 48px (empty states, banners)

**Colors:**
- Match text color in context
- Use gradient fill for important icons (e.g., banner icons)

### Spacing System (8px base)

```
--space-1: 4px
--space-2: 8px
--space-3: 12px
--space-4: 16px
--space-5: 24px
--space-6: 32px
--space-8: 48px
--space-10: 64px
```

### Responsive Breakpoints

```
--mobile: 640px
--tablet: 768px
--desktop: 1024px
--wide: 1280px
```

Desktop-first design (optimized for trading desk), but gracefully degrades to tablet.

### Typography Scale

```
--text-xs: 12px
--text-sm: 14px
--text-base: 16px
--text-lg: 18px
--text-xl: 20px
--text-2xl: 24px
--text-3xl: 30px
--text-4xl: 36px (banner text)
```

---

## User Workflow

### Daily Trading Routine (The Happy Path)

This is the **intended workflow** that the system enforces:

```gherkin
Feature: Daily trend-following workflow
  As a systematic trader
  I want to execute a disciplined daily routine
  So that I trade the tide, not the splash

  Background:
    Given I am using the TF-Engine GUI application
    And the backend tf-engine is running
    And I have my TradingView charts ready
    And my Ed-Seykota.pine script is loaded on TradingView

  Scenario: Morning scan and candidate identification
    When I open the Dashboard at market open
    Then I see my current open positions with heat % and stops
    And I see yesterday's candidates (if any remain)
    And I see the total portfolio heat vs 4% cap
    And I see any sector buckets on cooldown

    When I click "Run Daily FINVIZ Scan"
    Then the backend scrapes my preset FINVIZ screener
    And returns a list of tickers meeting technical criteria
    And displays count: "Found 23 candidates from TF_BREAKOUT_LONG"

    When I click "Review Candidates"
    Then I see a table of tickers with basic info
    And I can check/uncheck tickers to import
    And I can see which sectors are represented
    And I can see which are on cooldown (grayed out)

    When I select 12 tickers and click "Import Selected"
    Then those tickers are saved to the candidates table
    And they appear in my "Today's Candidates" list
    And I see notification: "12 candidates imported for 2025-10-29"

  Scenario: Chart analysis in TradingView
    Given I have imported 12 candidates
    When I click "Open in TradingView" next to ticker AAPL
    Then a new browser tab opens to AAPL chart on TradingView
    And the Ed-Seykota.pine script is already applied (if preset in TV)

    # Manual step (outside system):
    And I manually verify the 55-bar Donchian breakout
    And I note the current N (ATR) value from the script
    And I note the entry price (current close or limit)
    And I identify the sector bucket (e.g., "Tech/Comm")

    When I return to the TF-Engine GUI
    Then I am ready to fill out the checklist

  Scenario: Checklist evaluation (gates + quality)
    Given I have analyzed AAPL in TradingView
    When I navigate to the Checklist screen

    And I enter ticker: "AAPL"
    And I enter sector: "Tech/Comm" (dropdown)
    And I enter entry price: 180.50
    And I enter N (ATR): 2.35 (from TradingView)

    Then the system pre-calculates:
      - Stop distance: 2.0 × 2.35 = 4.70
      - Initial stop: 180.50 - 4.70 = 175.80 (for long)
      - Add-on levels: [182.68, 184.85, 187.03] (entry + k×0.5×N)

    # Required Gates (must check all):
    When I check "✓ Signal: 55-bar breakout confirmed (close > donHi[1])"
    And I check "✓ Risk/Size: Will use 2×N stop, add every 0.5×N, max 4 units"
    And I check "✓ Liquidity: Stock avg volume > 1M shares (or options OI > 100)"
    And I check "✓ Exits: Will exit on 10-bar opposite Donchian OR 2×N stop"
    And I check "✓ Behavior: Not on cooldown, heat caps OK, will honor 2-min timer"

    # Optional Quality Items (improve score):
    And I check "✓ Regime OK: SPY > 200 SMA (market favorable for longs)"
    And I check "✓ No Chase: Entry within 2N of 20-EMA (not extended)"
    And I check "✓ Earnings OK: No earnings within next 2 weeks (for long options)"
    And I enter journal note: "Clean breakout, strong volume, sector underweight"

    Then the quality score updates to 4 / 4 optional items
    And the banner changes from RED → YELLOW → GREEN
    And the banner displays: "🟢 OK TO TRADE 🟢"

    When I click "Save Evaluation"
    Then the backend records the evaluation timestamp
    And the 2-minute impulse brake timer starts
    And I see: "Evaluation saved. Cool-off period: 2 minutes remaining."

  Scenario: Position sizing calculation
    Given my checklist evaluation is GREEN
    And 2 minutes have NOT yet elapsed (still in cool-off)

    When I navigate to the Position Sizing screen
    And the form is pre-filled with:
      - Ticker: AAPL
      - Entry: 180.50
      - N (ATR): 2.35
      - Method: "stock" (dropdown: stock, opt-delta-atr, opt-contracts)
      - Risk % per unit: 0.75% (from settings)
      - Equity: $100,000 (from settings)
      - Max units: 4 (from settings)

    When I click "Calculate Position Size"
    Then the backend calculates using Van Tharp method:
      - Risk per unit: $100,000 × 0.0075 = $750
      - Stop distance: 2.0 × 2.35 = 4.70
      - Shares per unit: floor(750 ÷ 4.70) = 159 shares
      - Actual risk: 159 × 4.70 = $747.30 (≤ $750 ✓)
      - Initial stop: $175.80
      - Add-on schedule:
        * Unit 1: 159 shares @ 180.50 (now)
        * Unit 2: 159 shares @ 182.68 (+0.5N)
        * Unit 3: 159 shares @ 184.85 (+1.0N)
        * Unit 4: 159 shares @ 187.03 (+1.5N)

    And the results are displayed prominently
    And I see a warning if shares × price > 25% of equity (concentration risk)
    And I can adjust entry or method and recalculate

    When I click "Save Position Plan"
    Then the plan is saved to the database
    And it's linked to my checklist evaluation

  Scenario: Heat check before final entry
    Given my position sizing is calculated (159 shares AAPL, $747 risk)
    And my 2-minute cool-off is still active

    When I navigate to the Heat Check screen
    And I see my current portfolio heat: $2,890 / $4,000 cap (72.25%)
    And I see my Tech/Comm bucket heat: $1,125 / $1,500 cap (75%)

    When I click "Check Heat for AAPL Trade"
    Then the backend calculates:
      - New portfolio heat: $2,890 + $747 = $3,637 (90.93%)
      - Portfolio cap: $4,000 (4% of equity)
      - Remaining capacity: $363 (9.07%)
      - Result: ✓ WITHIN CAP

      - New bucket heat: $1,125 + $747 = $1,872 (124.8%)
      - Bucket cap: $1,500 (1.5% of equity)
      - Overage: $372 (24.8%)
      - Result: ✗ EXCEEDS BUCKET CAP

    Then I see a RED warning:
      "⚠️ BUCKET CAP EXCEEDED
       Tech/Comm bucket would be $1,872 (124.8% of $1,500 cap)
       Overage: $372

       Options:
       1. Reduce position size to max 96 shares (risk = $451)
       2. Close an existing Tech/Comm position first
       3. Choose a different ticker from another sector"

    And the "Proceed to Trade Entry" button is DISABLED
    And the banner remains GREEN but shows: "⚠️ Heat cap issue - see Heat Check"

  Scenario: Adjust position size to respect heat caps
    Given the bucket cap is exceeded by $372
    When I return to Position Sizing screen
    And I click "Calculate Max Shares for Heat Cap"

    Then the backend calculates:
      - Remaining bucket capacity: $1,500 - $1,125 = $375
      - Max shares: floor(375 ÷ 4.70) = 79 shares
      - Actual risk: 79 × 4.70 = $371.30 (≤ $375 ✓)

    And the form updates to show: 79 shares per unit
    And I see: "Position reduced to respect heat cap (79 shares instead of 159)"

    When I click "Save Adjusted Position Plan"
    Then the plan is updated with reduced size

    When I return to Heat Check and click "Recheck"
    Then I see:
      - New bucket heat: $1,125 + $371 = $1,496 (99.73%)
      - Result: ✓ WITHIN CAP

    And the "Proceed to Trade Entry" button becomes ENABLED

  Scenario: Final trade entry with 5-gate validation
    Given my checklist is GREEN
    And my position sizing is calculated (79 shares)
    And my heat check passed
    And 2 minutes have elapsed since checklist evaluation

    When I navigate to Trade Entry screen
    Then I see a summary of my trade plan:
      - Ticker: AAPL
      - Direction: LONG
      - Entry: $180.50
      - Shares: 79 per unit (max 4 units)
      - Initial stop: $175.80 (2×N)
      - Risk per unit: $371.30
      - Add-on levels: [182.68, 184.85, 187.03]
      - Sector: Tech/Comm
      - Exit plan: 10-bar Donchian OR 2×N stop, whichever closer

    When I click "Run Final Gate Check"
    Then the backend validates ALL 5 gates:

      Gate 1: Banner Status
        - Current: GREEN ✓
        - All required checks confirmed ✓
        - Quality score ≥ threshold ✓
        - Result: PASS ✓

      Gate 2: Impulse Brake (2-minute cool-off)
        - Last evaluation: 2025-10-29 09:32:15
        - Current time: 2025-10-29 09:34:27
        - Elapsed: 2 min 12 sec ✓
        - Result: PASS ✓

      Gate 3: Cooldown Status
        - Ticker AAPL on cooldown? NO ✓
        - Bucket Tech/Comm on cooldown? NO ✓
        - Result: PASS ✓

      Gate 4: Heat Caps
        - Portfolio heat: 90.93% of cap ✓
        - Bucket heat: 99.73% of cap ✓
        - Result: PASS ✓

      Gate 5: Sizing Completed
        - Position plan saved? YES ✓
        - Risk calculated? YES ($371.30) ✓
        - Result: PASS ✓

    And all gates show GREEN ✓
    And the "SAVE GO DECISION" button is ENABLED (green, large)
    And the "SAVE NO-GO DECISION" button is also visible (for journaling rejections)

    When I click "SAVE GO DECISION"
    Then the backend saves a decision record with:
      - Timestamp
      - Ticker, sector, entry, stop, shares
      - All gate results (all PASS)
      - Banner status (GREEN)
      - Quality score (4/4)
      - Journal note
      - Risk amount

    And I see confirmation: "✓ GO decision saved for AAPL"
    And the ticker appears in my "Ready to Execute" list
    And I can now manually execute the trade in my broker

    # System does NOT execute trades - this is intentional
    # Trader must manually enter in broker (final human verification)

  Scenario: Logging a NO-GO decision
    Given I complete analysis of ticker XYZ
    But the heat check shows portfolio at 95% of cap

    When I navigate to Trade Entry screen
    And I click "Run Final Gate Check"
    Then Gate 4 (Heat Caps) shows FAIL ✗
    And the "SAVE GO DECISION" button is DISABLED

    When I click "SAVE NO-GO DECISION"
    And I enter reason: "Portfolio heat at 95%, no capacity for new position"

    Then the backend saves a NO-GO decision record
    And I see: "✓ NO-GO decision logged for XYZ"
    And this is journaled for future review
    And the ticker is removed from candidates

  Scenario: Position appears in calendar view
    Given I saved a GO decision for AAPL (Tech/Comm)
    When I navigate to Calendar view
    Then I see a 10-week grid (2 weeks back + 8 weeks forward)
    And each row is a sector bucket
    And each column is a week (Mon-Sun date range)

    And I see AAPL appear in:
      - Tech/Comm row
      - Week of 2025-10-27 column (entry week)
      - Plus subsequent weeks until expected exit (~60 days if options)

    And I can see sector diversification at a glance
    And I can identify if too many trades in one sector/week
    And I can spot gaps (weeks with no trades)
```

### Additional Workflows

**Portfolio Management:**
- View all open positions with current heat
- Update stops as positions move (trailing 10-bar Donchian)
- Close positions (log exit decision)
- Add units (if favorable move ≥ 0.5N from last add)

**Settings Management:**
- Update equity (monthly or after significant P&L change)
- Adjust risk% per unit (0.5-1.0% typical range)
- Configure heat caps (portfolio 4%, bucket 1.5% defaults)
- Manage sector bucket definitions

**FINVIZ Screener Presets:**
- Create custom FINVIZ screener URLs
- Save as named presets (e.g., "TF_BREAKOUT_LONG", "TF_PULLBACK_LONG")
- One-click scan using saved preset

---

## Technical Components

### Backend API Endpoints (to implement)

The backend already has the domain logic. We need to expose it via HTTP:

**Settings:**
- `GET /api/settings` - Get current settings
- `PUT /api/settings` - Update settings

**Candidates:**
- `POST /api/candidates/scan` - Trigger FINVIZ scan
  - Body: `{ "preset": "TF_BREAKOUT_LONG" }`
  - Returns: `{ "count": 23, "tickers": [...] }`
- `POST /api/candidates/import` - Import selected candidates
  - Body: `{ "tickers": ["AAPL", "MSFT"], "date": "2025-10-29" }`
- `GET /api/candidates?date=2025-10-29` - Get candidates for date
- `DELETE /api/candidates/:ticker` - Remove candidate

**Checklist:**
- `POST /api/checklist/evaluate` - Evaluate checklist
  - Body: Full checklist data (ticker, checks, quality items)
  - Returns: `{ "banner": "GREEN", "missing_count": 0, "score": 4, ... }`
- `GET /api/checklist/:ticker` - Get last evaluation for ticker

**Position Sizing:**
- `POST /api/size/calculate` - Calculate position size
  - Body: `{ "ticker", "entry", "atr", "method", "k" }`
  - Returns: Full sizing result (shares, risk, stops, add-ons)

**Heat:**
- `POST /api/heat/check` - Check heat caps
  - Body: `{ "ticker", "risk_amount", "bucket" }`
  - Returns: Heat analysis (current, new, caps, result)

**Gates:**
- `POST /api/gates/check` - Run all 5 gates
  - Body: Full trade data
  - Returns: All gate results + final pass/fail

**Decisions:**
- `POST /api/decisions` - Save GO or NO-GO decision
  - Body: Full decision record
- `GET /api/decisions?date=2025-10-29` - Get decisions for date
- `GET /api/decisions/:id` - Get specific decision

**Positions:**
- `GET /api/positions` - Get all open positions
- `POST /api/positions` - Open new position (after GO decision)
- `PUT /api/positions/:id` - Update position (add unit, move stop)
- `DELETE /api/positions/:id` - Close position (log exit)

**Calendar:**
- `GET /api/calendar?weeks=10` - Get calendar data
  - Returns: Positions grouped by sector × week

### Frontend Components (to build)

**Layout:**
- `App.svelte` - Main layout with navigation and theme provider
- `Header.svelte` - App header with title, settings icon, **theme toggle button**
- `Navigation.svelte` - Sidebar or top nav with gradient accent
- `Banner.svelte` - Large 3-state banner with gradient backgrounds (RED/YELLOW/GREEN)
- `ThemeToggle.svelte` - Day/night mode switcher with smooth transition

**Dashboard:**
- `Dashboard.svelte` - Main landing page
- `PositionList.svelte` - Table of open positions
- `HeatGauge.svelte` - Visual heat indicators
- `CandidatesSummary.svelte` - Today's candidates count

**FINVIZ Scan:**
- `FINVIZScanner.svelte` - Scan trigger and results
- `CandidateImport.svelte` - Review and select candidates
- `PresetManager.svelte` - Manage FINVIZ presets

**Checklist:**
- `Checklist.svelte` - Main checklist form
- `RequiredGates.svelte` - 5 required checkboxes (must check all)
- `QualityItems.svelte` - 4 optional checkboxes (score boosters)
- `JournalNote.svelte` - Free-text journal entry

**Position Sizing:**
- `PositionSizer.svelte` - Main sizing form
- `SizingResults.svelte` - Display calculated shares, stops, add-ons
- `AddOnSchedule.svelte` - Visual add-on levels

**Heat Check:**
- `HeatCheck.svelte` - Main heat analysis
- `HeatWarning.svelte` - RED warning if caps exceeded
- `HeatSuggestions.svelte` - Suggestions to resolve cap issues

**Trade Entry:**
- `TradeEntry.svelte` - Main trade entry screen
- `GateResults.svelte` - Display all 5 gate results
- `TradeSummary.svelte` - Summary of trade plan
- `DecisionButtons.svelte` - SAVE GO / SAVE NO-GO buttons

**Calendar:**
- `Calendar.svelte` - 10-week sector × week grid
- `CalendarCell.svelte` - Individual cell (sector + week)
- `CalendarLegend.svelte` - Color coding, heat indicators

**TradingView Integration:**
- `TradingViewLink.svelte` - Button to open TV chart
- Constructs URL: `https://www.tradingview.com/chart/?symbol={ticker}`

**Utility Components:**
- `Modal.svelte` - Generic modal with backdrop blur and gradient border
- `LoadingSpinner.svelte` - Smooth animated spinner with gradient
- `ErrorMessage.svelte` - Display API errors with red gradient accent
- `SuccessMessage.svelte` - Display success with green gradient accent
- `Button.svelte` - Reusable button with gradient backgrounds and hover effects
- `Card.svelte` - Container with subtle shadow and gradient border
- `Badge.svelte` - Small status indicator with gradient fill

### State Management (Svelte)

**Stores (`src/lib/stores/`):**

```typescript
// settings.ts
export const settings = writable({
  equity: 100000,
  riskPct: 0.75,
  portfolioCap: 4.0,
  bucketCap: 1.5,
  maxUnits: 4,
  // ... other settings
});

// candidates.ts
export const candidates = writable([]);
export const candidatesDate = writable(null);

// checklist.ts
export const currentChecklist = writable({
  ticker: '',
  banner: 'RED',
  requiredGates: {
    signal: false,
    riskSize: false,
    liquidity: false,
    exits: false,
    behavior: false
  },
  qualityItems: {
    regime: false,
    noChase: false,
    earnings: false,
    journal: ''
  },
  score: 0,
  evaluatedAt: null
});

// positions.ts
export const positions = writable([]);
export const portfolioHeat = writable(0);
export const bucketHeat = writable({});

// ui.ts
export const banner = writable({
  state: 'RED', // 'RED' | 'YELLOW' | 'GREEN'
  message: 'Complete checklist to begin'
});

export const notifications = writable([]);
export const isLoading = writable(false);
```

**API Client (`src/lib/api/`):**

```typescript
// client.ts
const API_BASE = '/api';

export async function get(endpoint: string) {
  const res = await fetch(`${API_BASE}${endpoint}`);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function post(endpoint: string, data: any) {
  const res = await fetch(`${API_BASE}${endpoint}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

// Similar for PUT, DELETE...

// Specific API functions
export const api = {
  settings: {
    get: () => get('/settings'),
    update: (data) => put('/settings', data)
  },
  candidates: {
    scan: (preset) => post('/candidates/scan', { preset }),
    import: (tickers, date) => post('/candidates/import', { tickers, date }),
    list: (date) => get(`/candidates?date=${date}`)
  },
  checklist: {
    evaluate: (data) => post('/checklist/evaluate', data)
  },
  // ... etc
};
```

---

## Implementation Strategy

### Phase 0: Foundation & Proof-of-Concept (Week 1-2)

**Goals:**
- Validate Go + Svelte architecture
- Prove frontend ↔ backend communication works
- Test cross-compilation to Windows .exe
- Establish build and deployment workflow

**Deliverables:**
1. **Fyne POC** (Quick validation)
   - Simple Fyne app that calls backend functions
   - Proves Go GUI can work
   - Tests embedding and packaging
   - Duration: 2-3 days

2. **Svelte POC** (Real implementation start)
   - SvelteKit project with static adapter
   - Single component: Settings form
   - API call: `GET /api/settings`, `PUT /api/settings`
   - Go backend serves Svelte static files
   - Duration: 3-4 days

3. **Build Pipeline**
   - `scripts/sync-ui-to-go.sh` (copies Svelte build → Go embed)
   - `scripts/build-go-windows.sh` (cross-compile .exe)
   - Test: Run .exe on Windows, access via browser
   - Duration: 1-2 days

**Success Criteria:**
- ✓ Fyne POC displays data from backend
- ✓ Svelte POC renders in browser
- ✓ Backend serves Svelte static files
- ✓ API call successfully retrieves settings
- ✓ Windows .exe runs and opens browser to app
- ✓ Documentation updated (LLM-update.md, PROGRESS.md)

---

### Phase 1: Dashboard & FINVIZ Scanner (Week 3-4)

**Goals:**
- Build core navigation and layout
- Implement FINVIZ scanning workflow
- Display candidates and positions

**Deliverables:**
1. **Layout Components**
   - App shell with navigation
   - Header, sidebar/nav
   - Routing (SvelteKit routes)

2. **Dashboard Screen**
   - Display open positions (from `GET /api/positions`)
   - Show portfolio heat (from `GET /api/heat/current`)
   - Show today's candidates count
   - Quick stats: equity, available capacity

3. **FINVIZ Scanner**
   - One-click "Run Daily Scan" button
   - Calls `POST /api/candidates/scan`
   - Displays results table
   - Checkbox selection for import
   - "Import Selected" button

4. **Backend API Endpoints**
   - Implement HTTP server wrapper for existing domain logic
   - Endpoints: `/api/positions`, `/api/candidates/scan`, `/api/candidates/import`

**Success Criteria:**
- ✓ Dashboard displays real data from database
- ✓ FINVIZ scan button fetches candidates
- ✓ Candidate import saves to database
- ✓ Navigation works between screens

---

### Phase 2: Checklist & Position Sizing (Week 5-6)

**Goals:**
- Implement 5 gates + quality items checklist
- Display 3-state banner
- Calculate position sizing
- Enforce 2-minute cool-off

**Deliverables:**
1. **Checklist Screen**
   - Form: ticker, entry, ATR, sector
   - 5 required checkboxes (gates)
   - 4 optional checkboxes (quality)
   - Journal text area
   - Large banner at top (RED/YELLOW/GREEN)
   - "Save Evaluation" button

2. **Banner Component**
   - Minimum 20% screen height
   - RED: Any required unchecked
   - YELLOW: All required, score < 3
   - GREEN: All required, score ≥ 3
   - Updates live as checkboxes change

3. **Position Sizing Screen**
   - Form: pre-filled from checklist
   - Method dropdown (stock, opt-delta, opt-contracts)
   - "Calculate" button
   - Results: shares, risk, stops, add-ons
   - "Save Position Plan" button

4. **Backend Endpoints**
   - `POST /api/checklist/evaluate`
   - `POST /api/size/calculate`
   - Timestamp storage for 2-min timer

**Success Criteria:**
- ✓ Banner changes color based on checklist state
- ✓ Checklist evaluation saves to database
- ✓ 2-minute timer starts
- ✓ Position sizing calculates correctly (matches Van Tharp)
- ✓ Results match existing CLI output

---

### Phase 3: Heat Check & Trade Entry (Week 7-8)

**Goals:**
- Implement heat cap validation
- Build 5-gate final check
- Save GO/NO-GO decisions

**Deliverables:**
1. **Heat Check Screen**
   - Display current portfolio heat
   - Display current bucket heat
   - "Check Heat" button for proposed trade
   - RED warning if cap exceeded
   - Suggestions to resolve

2. **Trade Entry Screen**
   - Trade summary display
   - "Run Final Gate Check" button
   - Display all 5 gate results
   - Large "SAVE GO DECISION" button (disabled until all pass)
   - "SAVE NO-GO DECISION" button (always enabled)

3. **Gate Validation**
   - Gate 1: Banner GREEN
   - Gate 2: 2-min elapsed
   - Gate 3: Not on cooldown
   - Gate 4: Heat caps OK
   - Gate 5: Sizing done

4. **Backend Endpoints**
   - `POST /api/heat/check`
   - `POST /api/gates/check`
   - `POST /api/decisions`

**Success Criteria:**
- ✓ Heat check prevents cap violations
- ✓ All 5 gates must pass for GO decision
- ✓ 2-minute timer is enforced
- ✓ NO-GO decisions are logged
- ✓ No backdoors exist

---

### Phase 4: Calendar & TradingView Integration (Week 9-10)

**Goals:**
- Visualize sector diversification
- Integrate with TradingView charts
- Polish UI/UX

**Deliverables:**
1. **Calendar View**
   - 10-week grid (2 back + 8 forward)
   - Rows: sector buckets
   - Columns: weeks (Mon-Sun)
   - Cells: tickers active in that sector/week
   - Color coding: heat levels

2. **TradingView Integration**
   - "Open in TradingView" buttons on candidates
   - Construct URL with ticker symbol
   - Open in new tab
   - Documentation for applying Ed-Seykota.pine script

3. **UI Polish**
   - Keyboard shortcuts
   - Better loading states
   - Error handling improvements
   - Mobile responsive (bonus)

4. **Backend Endpoint**
   - `GET /api/calendar?weeks=10`

**Success Criteria:**
- ✓ Calendar displays positions correctly
- ✓ TradingView links work
- ✓ UI feels fast and responsive
- ✓ No obvious bugs or UX issues

---

### Phase 5: Testing & Packaging (Week 11-12)

**Goals:**
- Comprehensive testing
- Windows installer
- User documentation

**Deliverables:**
1. **Testing**
   - Frontend unit tests (Vitest)
   - Backend API integration tests
   - End-to-end workflow testing
   - Test all 5 gates enforcement
   - Test heat cap edge cases

2. **Windows Packaging**
   - Build .msi installer (WiX on Windows)
   - Or .exe installer (NSIS)
   - Desktop shortcut
   - Start menu entry
   - Uninstaller

3. **Documentation**
   - User guide (screenshots)
   - Quick start (5-minute setup)
   - Trading workflow guide
   - Troubleshooting
   - FAQ

**Success Criteria:**
- ✓ All critical paths tested
- ✓ Windows installer works on clean PC
- ✓ User can complete full workflow in < 10 minutes
- ✓ Documentation is clear and complete

---

## Proof-of-Concept Approach

### Why Start with Fyne POC?

**Rationale:**
1. **Fast validation** - Prove Go GUI integration works (1-2 days)
2. **Desktop-native feel** - Test if we want pure desktop vs browser-based
3. **Single binary** - Verify packaging and deployment
4. **Fallback option** - If Svelte doesn't work out, we have Fyne

**Fyne POC Scope (Minimal):**
```
1. Create simple Fyne window
2. Display: "TF-Engine Settings"
3. Show equity, risk%, caps from database
4. Button: "Refresh" (re-query database)
5. Button: "Update" (save changes)
6. Status: "Saved successfully" or error
```

**Fyne POC Success:**
- ✓ Window opens
- ✓ Data loads from SQLite
- ✓ Updates save to SQLite
- ✓ Compiles to Windows .exe
- ✓ Runs on Windows without dependencies

**Then Pivot to Svelte:**

**Why Svelte over Fyne for production?**
1. **Richer UI** - Easier to build complex, beautiful interfaces with gradients, animations, and modern design
2. **Visual Appeal** - CSS gradients, smooth transitions, sophisticated layouts (hard to achieve in Fyne)
3. **Theme Support** - Built-in day/night mode with CSS variables and seamless switching
4. **Web Skills** - More developers familiar with HTML/CSS/JS
5. **Rapid Iteration** - Hot reload, browser dev tools, instant visual feedback
6. **Flexibility** - Can add advanced features (charts, tables, custom animations)
7. **Ecosystem** - Huge library of components and tools (icon libraries, gradient generators, etc.)
8. **Polish** - Professional, modern appearance that makes the tool a pleasure to use

**Svelte POC Scope (Minimal):**
```
1. SvelteKit project with static adapter
2. Single page: Settings form
3. Fetch settings from GET /api/settings
4. Display in form inputs
5. Button: "Save" → PUT /api/settings
6. Success/error notification
```

**Svelte POC Success:**
- ✓ Svelte app builds to static files
- ✓ Go serves static files via embed
- ✓ API call works (GET /api/settings)
- ✓ Form submission works (PUT /api/settings)
- ✓ Browser opens to app when .exe runs
- ✓ Cross-compiles to Windows .exe

**Decision Point:**

After both POCs, we choose:
- **Fyne** if: Desktop-native feel is critical, Svelte integration is too complex
- **Svelte** if: POC works well, UI development is faster, team prefers web tech

**Expected choice: Svelte** (based on flexibility and development speed)

---

## Behavioral Specifications (Gherkin)

### Feature: Anti-Impulsivity Enforcement

```gherkin
Feature: System prevents impulsive trading
  As the trading system
  I must make it impossible to trade impulsively
  So that only systematic, trend-following trades are executed

  Background:
    Given the trader has analyzed a potential trade
    And the backend tf-engine is operational

  Scenario: Cannot bypass checklist evaluation
    When the trader tries to go directly to Trade Entry
    And the checklist has not been evaluated
    Then the system shows: "Complete checklist first"
    And the Trade Entry screen is READ-ONLY
    And the "SAVE GO DECISION" button is DISABLED

  Scenario: Cannot skip 2-minute cool-off
    Given the checklist was evaluated at 09:30:00
    And the current time is 09:31:45 (1 min 45 sec later)
    When the trader tries to save a GO decision
    Then the system shows: "Cool-off period: 15 seconds remaining"
    And the "SAVE GO DECISION" button is DISABLED
    And the countdown timer displays: 00:15

    When the current time becomes 09:32:01 (2 min 1 sec later)
    Then the countdown timer shows: 00:00
    And the "SAVE GO DECISION" button becomes ENABLED

  Scenario: Cannot exceed portfolio heat cap
    Given the portfolio heat is $3,800
    And the portfolio cap is $4,000 (4% of $100k)
    And the trader proposes a trade with $300 risk
    When the heat check is performed
    Then the new portfolio heat would be $4,100
    And this exceeds the cap by $100
    And the system shows RED warning: "Portfolio cap exceeded by $100"
    And the "SAVE GO DECISION" button is DISABLED
    And the system suggests: "Reduce position size or close existing position"

  Scenario: Cannot exceed bucket heat cap
    Given the Tech/Comm bucket heat is $1,400
    And the bucket cap is $1,500 (1.5% of $100k)
    And the trader proposes an AAPL trade (Tech/Comm) with $200 risk
    When the heat check is performed
    Then the new bucket heat would be $1,600
    And this exceeds the cap by $100
    And the system shows RED warning: "Tech/Comm bucket cap exceeded by $100"
    And the "SAVE GO DECISION" button is DISABLED

  Scenario: Cannot trade ticker on cooldown
    Given ticker AAPL is on cooldown until 2025-11-05
    And today is 2025-10-29
    When the trader selects AAPL for analysis
    Then the system shows: "⚠️ AAPL on cooldown until 2025-11-05 (7 days)"
    And the checklist form is DISABLED
    And the system suggests: "Choose a different ticker"

  Scenario: Cannot save GO decision with RED banner
    Given the checklist evaluation resulted in RED banner
    Because the "Signal" gate is unchecked
    When the trader navigates to Trade Entry
    Then Gate 1 (Banner Status) shows: FAIL ✗
    And the overall gate check is: FAIL
    And the "SAVE GO DECISION" button is DISABLED
    And only "SAVE NO-GO DECISION" is enabled

  Scenario: Cannot save GO decision with incomplete position sizing
    Given the checklist is GREEN
    And the 2-minute timer has elapsed
    And heat caps are OK
    And no cooldowns
    But position sizing has NOT been calculated
    When the trader tries to save a GO decision
    Then Gate 5 (Sizing Completed) shows: FAIL ✗
    And the "SAVE GO DECISION" button is DISABLED
    And the system shows: "Complete position sizing first"
```

### Feature: Systematic Workflow Enforcement

```gherkin
Feature: Daily trading workflow
  As a systematic trader
  I must follow the complete workflow
  So that I don't skip critical analysis steps

  Scenario: Morning scan workflow
    Given I open the app at market open
    When I navigate to Dashboard
    Then I see the "Run Daily FINVIZ Scan" button prominently

    When I click "Run Daily FINVIZ Scan"
    Then the backend scrapes my TF_BREAKOUT_LONG preset
    And returns candidates: ["AAPL", "MSFT", "NVDA", ...]
    And displays: "Found 23 candidates"

    When I click "Review Candidates"
    Then I see a table with columns: [Ticker, Sector, Last Close, Volume]
    And I can check/uncheck tickers
    And tickers on cooldown are grayed out

    When I select 12 tickers and click "Import Selected"
    Then those 12 tickers are saved to candidates table
    And they appear in "Today's Candidates" list

  Scenario: Chart analysis integration
    Given I have imported candidates
    When I click "Open in TradingView" next to AAPL
    Then a new browser tab opens to:
      https://www.tradingview.com/chart/?symbol=AAPL
    And I manually apply my Ed-Seykota.pine script (if not default)
    And I verify the 55-bar Donchian breakout visually
    And I note the current N (ATR) value: 2.35
    And I note the entry price: 180.50

  Scenario: Complete trade evaluation
    Given I have analyzed AAPL in TradingView
    When I fill out the checklist with:
      | Field         | Value      |
      | Ticker        | AAPL       |
      | Entry         | 180.50     |
      | ATR (N)       | 2.35       |
      | Sector        | Tech/Comm  |
    And I check all 5 required gates
    And I check 4 quality items
    And I enter journal note
    Then the banner turns GREEN

    When I save the evaluation
    Then the 2-minute timer starts
    And I see: "Cool-off: 2:00 remaining"

    When I calculate position sizing
    Then I get: 159 shares per unit, $747 risk

    When I check heat
    And all caps are OK
    And 2 minutes have elapsed
    Then all 5 gates pass
    And "SAVE GO DECISION" is enabled

    When I save the GO decision
    Then it's logged to the database
    And appears in "Ready to Execute" list
```

### Feature: Banner State Transitions

```gherkin
Feature: Banner reflects trade readiness
  As a visual discipline tool
  The banner must clearly show whether I can trade
  So that I never miss a failed gate

  Scenario: Initial state is RED
    Given I open a new checklist
    Then the banner shows: "🔴 DO NOT TRADE 🔴"
    And the banner is red background
    And the message says: "Complete required checklist items"

  Scenario: Transition to YELLOW
    Given the banner is RED
    When I check all 5 required gates
    But I check only 1 quality item (score = 1)
    And the quality threshold is 3
    Then the banner shows: "🟡 CAUTION 🟡"
    And the banner is yellow background
    And the message says: "Quality score below threshold (1/3)"

  Scenario: Transition to GREEN
    Given the banner is YELLOW
    When I check 2 more quality items (score now 3)
    Then the banner shows: "🟢 OK TO TRADE 🟢"
    And the banner is green background
    And the message says: "All gates pass, quality score met"

  Scenario: Regression from GREEN to YELLOW
    Given the banner is GREEN (all required, score 3/3)
    When I uncheck 1 quality item (score now 2/3)
    Then the banner immediately shows: "🟡 CAUTION 🟡"

  Scenario: Regression from YELLOW to RED
    Given the banner is YELLOW
    When I uncheck 1 required gate
    Then the banner immediately shows: "🔴 DO NOT TRADE 🔴"
```

### Feature: TradingView Integration

```gherkin
Feature: Seamless TradingView chart access
  As a trader
  I want to quickly open candidates in TradingView
  So that I can verify signals with my Pine script

  Scenario: Open ticker in TradingView
    Given I am viewing my candidates list
    When I click "Open in TV" for ticker AAPL
    Then a new browser tab opens
    And the URL is: https://www.tradingview.com/chart/?symbol=AAPL
    And I can manually load my Ed-Seykota.pine script

  Scenario: Verify Donchian breakout in TradingView
    Given I have AAPL chart open with Ed-Seykota.pine script
    And the script shows entryLen = 55, exitLen = 10
    When I look at the current bar
    Then I see if close > donHi[1] (for long breakout)
    And I note the N value (ATR) displayed by the script
    And I note the plotted stop level (entry - 2×N)

  Scenario: Return to app with parameters
    Given I verified the breakout in TradingView
    And N = 2.35, entry = 180.50
    When I return to the TF-Engine GUI
    And I enter those values in the checklist
    Then the system pre-calculates:
      - Stop: 180.50 - (2.0 × 2.35) = 175.80
      - Add1: 180.50 + (0.5 × 2.35) = 181.68
      - Add2: 180.50 + (1.0 × 2.35) = 182.85
      - Add3: 180.50 + (1.5 × 2.35) = 184.03
```

---

## Success Criteria

### Functional Requirements ✅

**Must Have:**
- [ ] FINVIZ scanner imports candidates via one click
- [ ] Checklist enforces all 5 required gates
- [ ] Banner visually indicates RED/YELLOW/GREEN state with **gradients** (large, obvious, beautiful)
- [ ] **Day/night mode toggle** with smooth transition (button in header)
- [ ] Theme preference persists to localStorage
- [ ] All UI elements use **gradients** (no flat colors)
- [ ] Smooth animations on all state changes (0.3s ease-in-out)
- [ ] 2-minute cool-off timer is enforced (cannot bypass)
- [ ] Position sizing calculates using Van Tharp method (matches backend CLI)
- [ ] Heat caps prevent trades exceeding 4% portfolio / 1.5% bucket
- [ ] All 5 gates must pass before GO decision saves
- [ ] NO-GO decisions are logged for journaling
- [ ] Calendar shows 10-week sector × week grid
- [ ] TradingView links open charts for analysis
- [ ] Settings persist in SQLite database
- [ ] Positions, decisions, candidates, cooldowns all persist
- [ ] **Comprehensive logging** throughout application
  - Backend: All API calls, database operations, calculations logged
  - Frontend: User actions, API calls, component lifecycle, errors logged
  - Log files accessible for debugging
  - Feature usage tracked to identify useful vs problematic features
  - Performance metrics captured for optimization

**Should Have:**
- [ ] Keyboard shortcuts for common actions (with visual hints)
- [ ] Loading states with animated gradient spinners
- [ ] Error messages with helpful suggestions and gradient accents
- [ ] Success confirmations with slide-in animations
- [ ] Auto-refresh after data changes
- [ ] Tooltips with smooth fade-in on hover
- [ ] Breadcrumb navigation with gradient dividers

**Nice to Have:**
- [ ] Mobile responsive (for viewing, not trading)
- [ ] Export decisions to CSV with styled download button
- [ ] Trade log search and filter with animated results
- [ ] Performance charts (P&L over time) with gradient fill
- [ ] Browser notifications (optional)
- [ ] Custom gradient presets (user can choose accent colors)
- [ ] Confetti animation on successful GO decision save 🎉

### Non-Functional Requirements ✅

**Performance:**
- [ ] Backend API responses < 100ms for calculations
- [ ] Database queries < 50ms
- [ ] FINVIZ scan < 5 seconds (network dependent)
- [ ] UI feels instant (no perceived lag)
- [ ] Page load < 2 seconds

**Reliability:**
- [ ] No crashes during normal operation
- [ ] Database corruption protection (SQLite WAL mode)
- [ ] Graceful error handling (API failures)
- [ ] Input validation prevents bad data

**Security:**
- [ ] Runs locally only (no external API calls except FINVIZ)
- [ ] No sensitive data in logs
- [ ] Database file permissions restricted (user only)
- [ ] No SQL injection vulnerabilities

**Usability:**
- [ ] Intuitive navigation (can complete workflow without docs)
- [ ] Large click targets (especially GO decision button)
- [ ] Clear visual hierarchy
- [ ] Consistent terminology (matches trading literature)
- [ ] Helpful error messages with suggestions

**Maintainability:**
- [ ] Code is well-commented
- [ ] API endpoints follow REST conventions
- [ ] Components are modular and reusable
- [ ] Documentation is up to date
- [ ] Build process is automated

**Deployability:**
- [ ] Single .exe for Windows (includes frontend + backend)
- [ ] No installation required (just run .exe)
- [ ] Opens browser automatically OR embedded webview
- [ ] Database created on first run
- [ ] Minimal dependencies (Go + SQLite only)

---

## Risk Management

### Technical Risks

**Risk: Svelte ↔ Go integration complexity**
- **Likelihood:** Medium
- **Impact:** High
- **Mitigation:** Start with Fyne POC first; if Svelte fails, fall back to Fyne
- **Contingency:** Build entire UI in Fyne (slower, but proven to work)

**Risk: Cross-compilation to Windows fails**
- **Likelihood:** Low (Go is mature)
- **Impact:** Critical
- **Mitigation:** Test early and often; use pure Go (no cgo)
- **Contingency:** Develop on Windows directly (abandon Linux-first workflow)

**Risk: FINVIZ changes page structure (scraper breaks)**
- **Likelihood:** Medium (over time)
- **Impact:** Medium (can still use manual import)
- **Mitigation:** Document scraper logic; make selectors configurable
- **Contingency:** Fallback to manual CSV import; consider FINVIZ API (paid)

**Risk: Browser embedding (webview) doesn't work well**
- **Likelihood:** Medium
- **Impact:** Low (can use default browser)
- **Mitigation:** Test webview early; have browser fallback ready
- **Contingency:** Always open in default browser (still works)

### Process Risks

**Risk: Scope creep (adding non-essential features)**
- **Likelihood:** High
- **Impact:** Medium (delays delivery)
- **Mitigation:** Strict adherence to anti-impulsivity principles; reject flexibility
- **Contingency:** Timebox phases; cut nice-to-haves

**Risk: Linux-first workflow causes Windows issues**
- **Likelihood:** Medium
- **Impact:** Medium (deployment friction)
- **Mitigation:** Test on Windows early and often; automate builds
- **Contingency:** Switch to Windows development if too painful

**Risk: Documentation falls behind code**
- **Likelihood:** High
- **Impact:** Low (but violates RULES.md)
- **Mitigation:** Daily Loop mandates updating LLM-update.md and PROGRESS.md
- **Contingency:** Dedicated doc sprint before each phase completion

### User Risks

**Risk: User finds workaround to bypass gates**
- **Likelihood:** Low (gates enforced in backend)
- **Impact:** Critical (defeats purpose)
- **Mitigation:** Backend validates all gates; UI cannot override
- **Contingency:** Log all decision attempts; review for patterns

**Risk: User doesn't understand the workflow**
- **Likelihood:** Medium
- **Impact:** Medium (frustration, abandonment)
- **Mitigation:** Clear documentation; in-app tooltips; guided tutorial
- **Contingency:** Video walkthrough; one-on-one training

**Risk: User trades manually outside system**
- **Likelihood:** High (cannot prevent)
- **Impact:** Low (system still useful for analysis)
- **Mitigation:** Make system so easy that manual tracking is more work
- **Contingency:** Accept this; system is a tool, not a jail

---

## Next Steps

### Immediate Actions (This Week)

1. **Review and Approve This Plan**
   - Read thoroughly
   - Ask questions
   - Identify gaps
   - Confirm priorities

2. **Create Step-by-Step Sub-Plans**
   - `plans/phase0-foundation-poc.md` - Fyne + Svelte POCs
   - `plans/phase1-dashboard-finviz.md` - Dashboard and scanning
   - `plans/phase2-checklist-sizing.md` - Checklist and sizing
   - `plans/phase3-heat-gates.md` - Heat check and trade entry
   - `plans/phase4-calendar-polish.md` - Calendar and final polish
   - `plans/phase5-testing-packaging.md` - Testing and Windows installer

3. **Set Up Frontend Project Structure**
   - Create `ui/` directory
   - Initialize SvelteKit project
   - Install dependencies (SvelteKit, TailwindCSS, etc.)
   - Configure static adapter
   - Create initial directory structure

4. **Begin Fyne POC**
   - Create simple Fyne app
   - Integrate with backend domain logic
   - Test database access
   - Test cross-compilation
   - Document results

### Week 1-2: Proof-of-Concept

- [ ] Fyne POC complete (2-3 days)
- [ ] Svelte POC complete (3-4 days)
- [ ] Build pipeline working (1-2 days)
- [ ] Decision: Fyne vs Svelte (1 day)
- [ ] Document results in `docs/LLM-update.md` and `docs/PROGRESS.md`

### Week 3-4: Dashboard & Scanner

- [ ] Layout and routing (2 days)
- [ ] Dashboard screen (3 days)
- [ ] FINVIZ scanner (3 days)
- [ ] Backend API endpoints (2 days)

### Week 5-6: Checklist & Sizing

- [ ] Checklist form (3 days)
- [ ] Banner component (2 days)
- [ ] Position sizing (3 days)
- [ ] 2-minute timer (2 days)

### Week 7-8: Heat & Gates

- [ ] Heat check (3 days)
- [ ] Trade entry screen (3 days)
- [ ] 5-gate validation (2 days)
- [ ] Decision saving (2 days)

### Week 9-10: Calendar & Polish

- [ ] Calendar view (3 days)
- [ ] TradingView integration (1 day)
- [ ] UI polish (3 days)
- [ ] Bug fixes (3 days)

### Week 11-12: Testing & Packaging

- [ ] Comprehensive testing (5 days)
- [ ] Windows installer (3 days)
- [ ] Documentation (2 days)

---

## Appendices

### A. Ed-Seykota.pine Script Parameters

From `reference/Ed-Seykota.pine`:

**Core Parameters:**
- `entryLen = 55` - Donchian entry lookback (System-2)
- `exitLen = 10` - Donchian exit lookback
- `nLen = 20` - N = ATR length
- `stopN = 2.0` - Initial stop in N multiples
- `addStepN = 0.5` - Add every X × N
- `maxUnits = 4` - Max units (including initial)
- `riskPct = 1.0` - Risk % of equity per unit

**Optional Filters:**
- `useMarket = false` - Market regime filter (SPY > 200 SMA)
- `timeExitMode = "None"` - Time exit (None/Close/Roll)
- `minVol = 0` - Minimum 20-bar avg volume

**Calculation:**
```pine
N = ta.atr(nLen)  // ATR(20)
donHi = ta.highest(high, entryLen)  // 55-bar high
donLo = ta.lowest(low, entryLen)    // 55-bar low
stopL = strategy.position_avg_price - stopN * N_entry  // Long stop
stopS = strategy.position_avg_price + stopN * N_entry  // Short stop
```

**Signal:**
```pine
longBreak = allowLong and liqOK and longRegOK and (close > donHiPrev)
shortBreak = allowShort and liqOK and shortRegOK and (close < donLoPrev)
```

**Position Sizing:**
```pine
sharesForUnit(_equity, _Nentry) =>
    riskDollars = _equity * (riskPct/100.0)
    perShareRisk = math.max(stopN * _Nentry, syminfo.mintick)
    math.max(1, math.floor(riskDollars / perShareRisk))
```

### B. Backend API Reference (Existing Domain Logic)

**Position Sizing:**
- `domain.CalculateSizeStock(ticker, entry, atr, k, riskPct, equity)`
- `domain.CalculateSizeOptDeltaATR(ticker, entry, atr, delta, k, riskPct, equity)`
- `domain.CalculateSizeOptContracts(ticker, entry, atr, contracts, debit, k, riskPct, equity)`

**Checklist:**
- `domain.EvaluateChecklist(fromPreset, trend, liquidity, timeframe, earnings, journal, ...)`
- Returns: `{ banner: "GREEN"|"YELLOW"|"RED", missing_count, missing_items, score, ... }`

**Heat:**
- `domain.CheckHeat(db, riskAmount, bucket, equity, portfolioCap, bucketCap)`
- Returns: `{ current_portfolio_heat, new_portfolio_heat, portfolio_cap_exceeded, bucket_cap_exceeded, ... }`

**Gates:**
- `domain.CheckGates(db, ticker, banner, riskDollars, bucket, equity, portfolioCap, bucketCap)`
- Returns: `{ gate1_pass, gate2_pass, gate3_pass, gate4_pass, gate5_pass, all_gates_pass }`

**Database:**
- `storage.NewDB(path)` - Open/create SQLite database
- `storage.GetSettings()`, `storage.UpdateSettings()`
- `storage.GetPositions()`, `storage.AddPosition()`, `storage.UpdatePosition()`, `storage.ClosePosition()`
- `storage.GetCandidates(date)`, `storage.AddCandidates(tickers, date)`
- `storage.SaveDecision()`, `storage.GetDecisions()`
- `storage.GetCooldowns()`, `storage.AddCooldown()`, `storage.RemoveCooldown()`

**FINVIZ Scraper:**
- `scrape.ScrapeFinviz(url)` - Returns list of tickers
- Presets: `TF_BREAKOUT_LONG`, `TF_BREAKOUT_SHORT`, etc.

### C. Reference Documents

**Must Read:**
1. `docs/project/WHY.md` - Core philosophy (why this system exists)
2. `docs/anti-impulsivity.md` - Design principles
3. `docs/dev/DEVELOPMENT_PHILOSOPHY.md` - How we build
4. `docs/dev/CLAUDE_RULES.md` - Development standards
5. `CLAUDE.md` - Guidance for future AI sessions
6. `1._RULES.md—Operating_Rules_for_This_Project-(Claude_Code).md` - Operating rules (Linux-first, no Git in Linux)

**Reference:**
- `reference/Ed-Seykota.pine` - TradingView Pine Script
- `FRESH_START_PLAN.md` - Original GUI plan (superseded by this)
- `PROJECT_HISTORY.md` - Why we abandoned Excel/VBA
- `docs/PROJECT_STATUS.md` - M24 completion status

**Ongoing:**
- `docs/LLM-update.md` - Session-by-session log (always current)
- `docs/PROGRESS.md` - Narrative progress and decisions

---

**End of Overview Plan**

**Status:** 📋 Ready for Review
**Next:** Approve plan → Create phase-specific sub-plans → Begin Fyne POC
**Updated:** 2025-10-29
