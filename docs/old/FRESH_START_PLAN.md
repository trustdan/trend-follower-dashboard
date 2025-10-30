# Fresh Start Plan - Custom GUI Trading Platform

**Created:** October 29, 2025
**Goal:** Replace Excel/VBA frontend with a custom GUI application
**Backend Status:** ✅ Fully functional (tf-engine in Go)
**Frontend Status:** 🚧 To be built from scratch

---

## Vision: Anti-Impulsivity Trading System

Based on [docs/anti-impulsivity.md](docs/anti-impulsivity.md), our system follows these principles:

### Core Philosophy
- **Trade the tide, not the splash** - Donchian breakouts with mechanical exits
- **Friction where it matters** - Hard gates for signal, risk, liquidity, exit, behavior
- **Nudge for better trades** - Optional quality items affect score, not permission
- **Immediate feedback** - Large 3-state banner (RED/YELLOW/GREEN) updates live
- **Journal while you decide** - One-click logging of full trade plan
- **Calendar awareness** - Rolling 10-week sector calendar for diversification visibility

### The 5 Hard Gates

1. **Signal (SIG_REQ):** 55-bar breakout (long > 55-high / short < 55-low)
2. **Risk/Size (RISK_REQ):** Per-unit risk = % of equity using 2×N stop; pyramids every 0.5×N to max units
3. **Options (OPT_REQ):** 60–90 DTE, roll/close ~21 DTE, liquidity required (bid/ask < 10% mid; OI > 100)
4. **Exits (EXIT_REQ):** Exit by 10-bar opposite Donchian OR closer of 2×N
5. **Behavior (BEHAV_REQ):** 2-minute cool-off + no intraday overrides

**RED (Do Not Trade):** Any required gate fails
**YELLOW (Caution):** All required pass, but quality score below threshold
**GREEN (OK to Trade):** All required pass + quality score meets threshold

### Optional Quality Nudges (affect score, not permission)

- Regime OK (e.g., SPY > 200SMA for longs)
- No chase (> 2N above 20-EMA at entry)
- Earnings blackout when long premium
- Journal note: why this unit/now, profit-take plan for debit verticals

Each optional box adds **1 point**. Default threshold: **3.0**

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     CUSTOM GUI FRONTEND                     │
│                      (Go + Fyne/Gio)                        │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  Dashboard   │  │ Checklist    │  │   Calendar   │     │
│  │              │  │  (5 Gates)   │  │ (10-week)    │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │Position Size │  │  Heat Check  │  │ Trade Entry  │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│                                                              │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       │ In-process function calls
                       │ (no HTTP, no CLI spawning)
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    TF-ENGINE BACKEND                        │
│                        (Go)                                  │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Domain     │  │   Storage    │  │   Scrape     │     │
│  │   (Logic)    │  │  (SQLite)    │  │  (FINVIZ)    │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

### Key Design Decisions

1. **Single Go binary** - GUI and backend in one executable
2. **No HTTP/CLI overhead** - Direct function calls within same process
3. **Cross-platform** - Works on Windows, Linux, macOS
4. **Native performance** - Fast, responsive, no browser overhead
5. **Reuse 100% of backend** - All domain logic, storage, scraping intact

---

## Technology Stack Options

### Option A: Fyne (Recommended)

**Pros:**
- Pure Go - same language as backend
- Cross-platform (Windows, Linux, macOS)
- Material Design look
- Good documentation
- Active community
- Easy to package as single binary

**Cons:**
- Less flexible layout options
- Somewhat limited styling

**Use Case:** Best for rapid development with clean, modern UI

**Example:**
```go
import "fyne.io/fyne/v2/app"
import "fyne.io/fyne/v2/widget"

func main() {
    myApp := app.New()
    myWindow := myApp.NewWindow("TF-Engine")

    // Build UI with widgets
    banner := widget.NewLabel("DO NOT TRADE")
    // ... more widgets

    myWindow.SetContent(banner)
    myWindow.ShowAndRun()
}
```

### Option B: Gio

**Pros:**
- Immediate mode GUI (like ImGui)
- Excellent performance
- Very flexible layouts
- Low-level control

**Cons:**
- Steeper learning curve
- More boilerplate
- Smaller community

**Use Case:** Best for performance-critical or highly custom UIs

### Option C: Wails (Go + Web Frontend)

**Pros:**
- Use web technologies (React, Vue, Svelte)
- Familiar to web devs
- Rich component libraries available
- Modern styling with Tailwind, etc.

**Cons:**
- Larger binary size
- More complex build process
- Still need to learn Go for backend

**Use Case:** Best if you prefer web UI development

### Recommendation: **Fyne**

For this project, **Fyne** is the best choice:
- Same language as backend (Go)
- Fast development iteration
- Native look and feel
- Easy deployment
- Good enough for our use case

---

## GUI Structure

### 1. Dashboard Screen

**Layout:**
```
┌────────────────────────────────────────────────────────────┐
│  TF-ENGINE                                    🟢 Connected  │
├────────────────────────────────────────────────────────────┤
│                                                              │
│  Account Equity: $100,000          Risk % per unit: 0.75%  │
│  Portfolio Heat: 2.3% / 4.0%       Active Positions: 3     │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Open Positions                                        │  │
│  ├──────────┬─────────┬──────────┬──────────┬──────────┤  │
│  │ Ticker   │ Entry   │ Risk $   │ Bucket   │ Days     │  │
│  ├──────────┼─────────┼──────────┼──────────┼──────────┤  │
│  │ AAPL     │ $180.50 │ $750     │ Tech     │ 12       │  │
│  │ XOM      │ $95.20  │ $750     │ Energy   │ 8        │  │
│  │ CAT      │ $250.00 │ $750     │ Industrl │ 3        │  │
│  └──────────┴─────────┴──────────┴──────────┴──────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Today's Candidates (from FINVIZ)                      │  │
│  ├──────────┬──────────┬──────────┬──────────┬─────────┤  │
│  │ Ticker   │ Sector   │ Price    │ Volume   │ ATR     │  │
│  ├──────────┼──────────┼──────────┼──────────┼─────────┤  │
│  │ NVDA     │ Tech     │ $450.00  │ 50M      │ 12.5    │  │
│  │ TSLA     │ Auto     │ $250.00  │ 120M     │ 8.3     │  │
│  └──────────┴──────────┴──────────┴──────────┴─────────┘  │
│                                                              │
│  [Refresh Data]  [Import from FINVIZ]                      │
│                                                              │
└────────────────────────────────────────────────────────────┘
```

**Components:**
- Account summary (equity, risk %, heat)
- Open positions table
- Today's candidates table
- Cooldowns list
- Refresh and import buttons

---

### 2. Checklist Screen (The Heart of Anti-Impulsivity)

**Layout:**
```
┌────────────────────────────────────────────────────────────┐
│  CHECKLIST: New Trade Evaluation                           │
├────────────────────────────────────────────────────────────┤
│                                                              │
│  Ticker: [AAPL    ]  Sector: [Tech/Comm ▼]                 │
│  Entry:  [$180.00 ]  ATR(N): [1.50      ]                  │
│  DTE:    [75      ]  Method: [stock     ▼]                 │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  REQUIRED GATES (All must pass)                      │  │
│  ├──────────────────────────────────────────────────────┤  │
│  │  ☑ Signal: 55-bar breakout confirmed                 │  │
│  │  ☑ Risk/Size: 2×N stop, 0.5×N adds, max 4 units     │  │
│  │  ☑ Options: 60-90 DTE, liquidity OK (bid/ask, OI)   │  │
│  │  ☑ Exits: 10-bar or 2×N exit plan in place          │  │
│  │  ☑ Behavior: 2-min cooloff, no intraday override    │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  OPTIONAL QUALITY (Score: 4 / Threshold: 3)         │  │
│  ├──────────────────────────────────────────────────────┤  │
│  │  ☑ Regime OK (SPY > 200SMA for longs)               │  │
│  │  ☑ No chase (< 2N above 20-EMA)                     │  │
│  │  ☑ Earnings blackout OK                             │  │
│  │  ☑ Journal note complete                            │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                  🟢 OK TO TRADE                       │  │
│  │                                                        │  │
│  │  All required gates passed                            │  │
│  │  Quality score: 4/3 (meets threshold)                │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  [Calculate Position]  [Add to Trades]  [Save Decision]    │
│                                                              │
└────────────────────────────────────────────────────────────┘
```

**Banner States:**
- **🔴 RED: DO NOT TRADE** - Background red, large text, any required gate fails
- **🟡 YELLOW: CAUTION** - Background yellow, required pass but quality score < threshold
- **🟢 GREEN: OK TO TRADE** - Background green, all gates pass, quality meets threshold

**Key Features:**
- Live banner updates as checkboxes change
- Clear visual feedback
- Required vs optional distinction
- Quality score calculation visible

---

### 3. Position Sizing Screen

**Layout:**
```
┌────────────────────────────────────────────────────────────┐
│  POSITION SIZING: Calculate Shares/Contracts               │
├────────────────────────────────────────────────────────────┤
│                                                              │
│  Account Size:  [$100,000]                                  │
│  Risk % per:    [0.75%   ]                                  │
│                                                              │
│  Ticker:        [AAPL    ]                                  │
│  Entry Price:   [$180.00 ]                                  │
│  ATR (N):       [1.50    ]                                  │
│  Stop Multiple: [2       ]  (Distance = 2 × 1.50 = $3.00)  │
│                                                              │
│  Method:        [● stock  ○ opt-delta-atr  ○ opt-contracts]│
│                                                              │
│  [Calculate Position]                                       │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  RESULTS                                              │  │
│  ├──────────────────────────────────────────────────────┤  │
│  │  Risk Dollars:      $750.00                          │  │
│  │  Stop Distance:     $3.00                            │  │
│  │  Initial Stop:      $177.00                          │  │
│  │  Shares:            250                              │  │
│  │  Position Value:    $45,000                          │  │
│  │  Actual Risk:       $750.00                          │  │
│  │                                                        │  │
│  │  Add Levels (0.5N increments):                       │  │
│  │    Add 1:  $181.50                                   │  │
│  │    Add 2:  $183.00                                   │  │
│  │    Add 3:  $184.50                                   │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  [Copy to Checklist]  [Save]                               │
│                                                              │
└────────────────────────────────────────────────────────────┘
```

**Features:**
- All 3 sizing methods (stock, opt-delta-atr, opt-contracts)
- Live calculation
- Add levels displayed
- Copy to clipboard functionality

---

### 4. Heat Check Screen

**Layout:**
```
┌────────────────────────────────────────────────────────────┐
│  HEAT CHECK: Risk Management                               │
├────────────────────────────────────────────────────────────┤
│                                                              │
│  Ticker:        [NVDA     ]                                 │
│  Risk Amount:   [$750     ]                                 │
│  Sector Bucket: [Tech/Comm ▼]                              │
│                                                              │
│  [Check Heat]                                               │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  PORTFOLIO HEAT                                       │  │
│  ├──────────────────────────────────────────────────────┤  │
│  │  Current:  2.25%  ████████░░░░░░░░░░  (2.25% / 4.0%)│  │
│  │  After:    3.00%  ████████████░░░░░░  (3.00% / 4.0%)│  │
│  │                                                        │  │
│  │  Status: ✅ Within 4% portfolio cap                  │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  SECTOR HEAT: Tech/Comm                              │  │
│  ├──────────────────────────────────────────────────────┤  │
│  │  Current:  0.75%  ████░░░░░░░░░░░░░░  (0.75% / 1.5%)│  │
│  │  After:    1.50%  ██████████░░░░░░░░  (1.50% / 1.5%)│  │
│  │                                                        │  │
│  │  Status: ⚠️  At 1.5% sector cap limit               │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  Overall: ⚠️  CAUTION - Sector at limit                   │
│                                                              │
│  [Copy to Checklist]                                       │
│                                                              │
└────────────────────────────────────────────────────────────┘
```

**Features:**
- Visual progress bars for heat levels
- Clear status indicators (✅/⚠️/❌)
- Before and after heat display
- Portfolio and sector heat breakdown

---

### 5. Trade Entry Screen (5 Hard Gates Enforcement)

**Layout:**
```
┌────────────────────────────────────────────────────────────┐
│  TRADE ENTRY: Final Gate Check                             │
├────────────────────────────────────────────────────────────┤
│                                                              │
│  Ticker:          [AAPL     ]                               │
│  Entry Price:     [$180.00  ]                               │
│  ATR (N):         [1.50     ]                               │
│  Method:          [stock    ▼]                             │
│  Banner Status:   [GREEN    ▼]                             │
│  Shares/Contracts:[250      ]                               │
│  Sector Bucket:   [Tech/Comm ▼]                            │
│  Strategy Preset: [TF_BREAKOUT_LONG ▼]                     │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  5 HARD GATES                                         │  │
│  ├──────────────────────────────────────────────────────┤  │
│  │  1. ✅ Banner is GREEN                               │  │
│  │  2. ✅ Cooloff timer elapsed (> 2 minutes)           │  │
│  │  3. ✅ Not on cooldown list                          │  │
│  │  4. ✅ Heat check passed (portfolio & sector)        │  │
│  │  5. ✅ Position sizing completed                     │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  Decision: ● GO  ○ NO-GO                                   │
│                                                              │
│  Notes: [Breakout of 55-day high. Trend confirmed.      ]  │
│         [4-week ATR stabilized. Earnings 6 weeks out.    ]  │
│                                                              │
│  [Save GO Decision]  [Save NO-GO Decision]                 │
│                                                              │
└────────────────────────────────────────────────────────────┘
```

**Features:**
- All 5 gates must pass for GO decision
- Clear visual feedback for each gate
- Notes field for journaling
- GO/NO-GO radio buttons
- Disabled save button if gates fail

---

### 6. Calendar Screen (10-Week Sector View)

**Layout:**
```
┌────────────────────────────────────────────────────────────┐
│  CALENDAR: Rolling 10-Week Sector View                     │
│  (2 weeks back + 8 weeks forward)                          │
├────────────────────────────────────────────────────────────┤
│                                                              │
│  Week Starting: [Oct 21] ◀ ▶  [Refresh]                    │
│                                                              │
│  ┌──────────┬─────────┬─────────┬─────────┬─────────┬────┐│
│  │ Sector   │ Oct 21  │ Oct 28  │ Nov 4   │ Nov 11  │... ││
│  ├──────────┼─────────┼─────────┼─────────┼─────────┼────┤│
│  │ Tech     │ AAPL    │ AAPL    │ AAPL    │ AAPL    │    ││
│  │          │ NVDA    │ NVDA    │ NVDA    │         │    ││
│  ├──────────┼─────────┼─────────┼─────────┼─────────┼────┤│
│  │ Energy   │ XOM     │ XOM     │ XOM     │ XOM     │    ││
│  │          │         │         │         │ CVX     │    ││
│  ├──────────┼─────────┼─────────┼─────────┼─────────┼────┤│
│  │ Industrl │         │ CAT     │ CAT     │ CAT     │    ││
│  │          │         │         │ BA      │ BA      │    ││
│  ├──────────┼─────────┼─────────┼─────────┼─────────┼────┤│
│  │ Finance  │         │         │         │ JPM     │    ││
│  │          │         │         │         │         │    ││
│  └──────────┴─────────┴─────────┴─────────┴─────────┴────┘│
│                                                              │
│  ● Active position  ○ Planned trade  × Cooldown            │
│                                                              │
└────────────────────────────────────────────────────────────┘
```

**Features:**
- Rolling 10-week window
- Sector × Week grid
- Visual density/crowding
- Color coding for status
- Quick navigation

---

## Implementation Plan

### Phase 1: Foundation (Week 1-2)

**Goal:** Basic GUI skeleton with backend integration

**Tasks:**
1. Set up Fyne project structure
2. Create main window with navigation
3. Integrate tf-engine backend (direct function calls)
4. Build Dashboard screen (read-only display)
5. Test SQLite database integration

**Deliverables:**
- Running GUI application
- Dashboard displaying positions, candidates, settings
- Refresh button working

---

### Phase 2: Core Functionality (Week 3-4)

**Goal:** Position sizing and checklist working

**Tasks:**
1. Build Position Sizing screen
2. Implement calculation logic (call domain.CalculateSize)
3. Build Checklist screen with checkboxes
4. Implement live banner updates (RED/YELLOW/GREEN)
5. Add "Copy to Checklist" functionality

**Deliverables:**
- Position sizing calculates correctly
- Checklist evaluates with banner colors
- Data flows between screens

---

### Phase 3: Heat & Gates (Week 5-6)

**Goal:** Heat check and trade entry working

**Tasks:**
1. Build Heat Check screen
2. Implement heat calculations and visual bars
3. Build Trade Entry screen
4. Implement 5 gates enforcement
5. Add decision logging (GO/NO-GO)

**Deliverables:**
- Heat check validates portfolio/sector caps
- Trade entry enforces all 5 gates
- Decisions saved to database

---

### Phase 4: Calendar & Polish (Week 7-8)

**Goal:** Calendar view and UX improvements

**Tasks:**
1. Build Calendar screen (10-week grid)
2. Implement sector × week display
3. Add navigation and filtering
4. Polish all screens (colors, fonts, spacing)
5. Add keyboard shortcuts
6. Implement validation and error messages

**Deliverables:**
- Calendar shows rolling 10-week view
- All screens polished and user-friendly
- Keyboard navigation working

---

### Phase 5: Testing & Deployment (Week 9-10)

**Goal:** Stable, deployable application

**Tasks:**
1. Integration testing (all screens)
2. End-to-end workflow testing
3. Performance optimization
4. Package for Windows/Linux/macOS
5. Create user documentation
6. Create installation guide

**Deliverables:**
- Fully tested application
- Single-binary executables for each platform
- User manual
- Installation instructions

---

## File Structure

```
fresh-start/
├── backend/               # tf-engine Go backend (existing)
│   ├── cmd/
│   ├── internal/
│   ├── go.mod
│   └── go.sum
│
├── frontend/              # New GUI frontend
│   ├── main.go           # Entry point
│   ├── ui/
│   │   ├── dashboard.go  # Dashboard screen
│   │   ├── checklist.go  # Checklist screen
│   │   ├── sizing.go     # Position sizing screen
│   │   ├── heat.go       # Heat check screen
│   │   ├── entry.go      # Trade entry screen
│   │   ├── calendar.go   # Calendar screen
│   │   └── components/   # Reusable UI components
│   │       ├── banner.go
│   │       ├── table.go
│   │       └── form.go
│   ├── bridge/            # Bridge to backend
│   │   ├── positions.go
│   │   ├── sizing.go
│   │   ├── checklist.go
│   │   ├── heat.go
│   │   └── gates.go
│   ├── go.mod
│   └── go.sum
│
├── docs/                  # Documentation
│   ├── anti-impulsivity.md
│   ├── PROJECT_HISTORY.md
│   ├── FRESH_START_PLAN.md
│   └── ...
│
├── scripts/               # Build and import scripts
├── art/                   # Assets
├── test-data/             # Test databases
└── README.md
```

---

## Key Technical Decisions

### 1. Direct Function Calls (No HTTP)

**Before (Excel/VBA):**
```vba
' VBA calls CLI
Set cmd = CreateObject("WScript.Shell")
cmd.Run "tf-engine.exe size --ticker AAPL --entry 180 --atr 1.5"
' Parse JSON output from stdout
```

**After (GUI):**
```go
// Direct function call
import "github.com/youruser/tf-engine/internal/domain"

func calculateSize(ticker string, entry, atr float64) domain.SizeResult {
    return domain.CalculateSize(ticker, entry, atr, method, k, riskPct)
}
```

**Benefits:**
- ✅ No process spawning overhead
- ✅ No JSON serialization/parsing
- ✅ Type safety
- ✅ Better error handling
- ✅ Faster execution

---

### 2. Single Binary Deployment

**Goal:** `tf-engine-gui.exe` contains everything

**Benefits:**
- ✅ No installation required
- ✅ No dependency management
- ✅ Portable (copy to any machine)
- ✅ Version control (one file = one version)

**How:**
```bash
# Build for Windows
GOOS=windows GOARCH=amd64 go build -o tf-engine-gui.exe

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o tf-engine-gui

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o tf-engine-gui-mac
```

---

### 3. Embedded SQLite Database

**Database file location:**
- Windows: `%APPDATA%\tf-engine\trading.db`
- Linux: `~/.config/tf-engine/trading.db`
- macOS: `~/Library/Application Support/tf-engine/trading.db`

**Benefits:**
- ✅ No database server required
- ✅ Automatic creation on first run
- ✅ Easy backup (just copy the file)
- ✅ Portable data

---

## Anti-Impulsivity Features in GUI

### 1. Large Banner (Impossible to Miss)

**Implementation:**
```go
// 3 states with distinct colors and large text
banner := widget.NewLabel("DO NOT TRADE")
banner.TextStyle = fyne.TextStyle{Bold: true}

switch status {
case "RED":
    banner.Text = "🔴 DO NOT TRADE"
    banner.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
case "YELLOW":
    banner.Text = "🟡 CAUTION"
    banner.Color = color.RGBA{R: 255, G: 200, B: 0, A: 255}
case "GREEN":
    banner.Text = "🟢 OK TO TRADE"
    banner.Color = color.RGBA{R: 0, G: 200, B: 0, A: 255}
}
```

### 2. Live Updates (No "Calculate" Button Lag)

**Implementation:**
```go
// Checkboxes trigger immediate banner refresh
checkbox1.OnChanged = func(checked bool) {
    evaluateGates()
    updateBanner()
}
```

### 3. 2-Minute Cooloff Timer (Visible Countdown)

**Implementation:**
```go
// Timer displayed in Trade Entry screen
timerLabel := widget.NewLabel("Cooloff: 1:45 remaining")

// Disable save button until timer expires
saveButton.Disable()

go func() {
    for remaining > 0 {
        time.Sleep(1 * time.Second)
        remaining--
        timerLabel.SetText(fmt.Sprintf("Cooloff: %d:%02d remaining", remaining/60, remaining%60))
    }
    saveButton.Enable()
}()
```

### 4. Required vs Optional Visual Distinction

**Implementation:**
```go
// Required gates in one panel with red border
requiredPanel := widget.NewCard("REQUIRED GATES (All must pass)", "", requiredCheckboxes)
requiredPanel.BorderColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}

// Optional quality in separate panel with blue border
optionalPanel := widget.NewCard("OPTIONAL QUALITY (Score: 4/3)", "", optionalCheckboxes)
optionalPanel.BorderColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
```

---

## Migration Path

### Week 1: Parallel Operation

**Both systems running:**
- Old: Excel/VBA (legacy, buggy)
- New: GUI (basic, read-only)

**User can compare:**
- Same database
- Same calculations
- Verify GUI matches Excel results

---

### Week 2-4: Feature Parity

**GUI reaches feature parity:**
- All 5 screens functional
- All calculations working
- All data displayed correctly

**User primarily uses GUI:**
- Excel available as backup
- Excel macros still work (via CLI)

---

### Week 5+: Full Migration

**GUI is primary system:**
- Excel retired
- All workflows in GUI
- No more VBA issues

---

## Success Metrics

### Functional

- ✅ Position sizing calculates correctly (all 3 methods)
- ✅ Checklist evaluates with correct banner colors
- ✅ Heat check validates portfolio/sector caps
- ✅ Trade entry enforces all 5 gates
- ✅ Dashboard displays all data
- ✅ Calendar shows 10-week sector view
- ✅ FINVIZ import populates candidates
- ✅ All data persists to SQLite

### Non-Functional

- ⚡ Instant UI updates (< 100ms)
- 🚀 Single binary deployment
- 🖥️  Cross-platform (Windows, Linux, macOS)
- 📦 Small binary size (< 50MB)
- 🔒 No macro security warnings
- 🎨 Clean, modern UI
- ⌨️  Keyboard navigation
- 📱 Responsive layout

### User Experience

- 😊 No VBA errors
- 😊 No manual module imports
- 😊 No "argument not optional" bugs
- 😊 No type mismatch errors
- 😊 Fast, responsive interface
- 😊 Clear visual feedback
- 😊 One-click actions

---

## Anti-Patterns to Avoid

### ❌ Don't: Replicate Excel's Layout

Excel has cells, rows, columns. GUI doesn't need to.
Focus on workflow, not mimicking spreadsheet structure.

### ❌ Don't: Use HTTP Server Mode

We're in the same process - use direct function calls.

### ❌ Don't: Parse JSON Strings

Use Go structs - type safety matters.

### ❌ Don't: Over-engineer

Start simple. Add features as needed.

### ❌ Don't: Ignore Anti-Impulsivity Design

The banner, gates, and cooloff are core features - not nice-to-haves.

---

## Next Steps

1. **Review this plan** - Verify it aligns with your vision
2. **Choose GUI framework** - Fyne recommended, but open to alternatives
3. **Set up project structure** - Create frontend/ directory
4. **Build Hello World** - Minimal Fyne app that displays "Hello TF-Engine"
5. **Integrate backend** - Import tf-engine packages, call one function
6. **Build Dashboard** - First real screen

---

## Questions to Answer

1. **GUI Framework:** Fyne, Gio, or Wails?
2. **Look & Feel:** Material Design, Windows native, or custom?
3. **Deployment:** Single binary, installer, or both?
4. **Platform Priority:** Windows first, or cross-platform from day 1?
5. **Timeline:** Aggressive (8 weeks) or conservative (12 weeks)?

---

## Conclusion

**This is a fresh start, but we're not starting from zero.**

We have:
- ✅ 100% functional backend
- ✅ All business logic tested
- ✅ Database schema proven
- ✅ Domain knowledge documented
- ✅ Clear anti-impulsivity design principles

We're replacing:
- ❌ Excel workbook
- ❌ VBA macros
- ❌ Form controls

With:
- ✅ Native GUI application
- ✅ Direct function calls
- ✅ Modern UI toolkit

**The hard part (business logic) is done. Now we build a better interface.**

---

**Last Updated:** October 29, 2025
**Status:** 📋 Planning complete, ready to start Phase 1
**Next Action:** Choose GUI framework and build Hello World
