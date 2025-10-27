# Excel Trading Workflow

**Complete Excel VBA + Python trading system implementing Seykota/Turtle trend-following methodology with options execution.**

---

## Table of Contents

1. [System Overview](#system-overview)
2. [Quick Start](#quick-start)
3. [Project Structure](#project-structure)
4. [Setup Instructions](#setup-instructions)
5. [Daily Workflow](#daily-workflow)
6. [System Capabilities](#system-capabilities)
7. [Testing & Validation](#testing--validation)
8. [Key Concepts](#key-concepts)
9. [Customization](#customization)
10. [Troubleshooting](#troubleshooting)
11. [Architecture & Design](#architecture--design)

---

## System Overview

### Philosophy

This trading system is designed to **eliminate discretionary bias** and **enforce mechanical rules** through:
- ✅ Structured checklists (all boxes must be checked)
- ✅ Automated position sizing (3 methods: Stock, Delta-ATR, MaxLoss)
- ✅ Portfolio/bucket heat caps (hard limits on risk exposure)
- ✅ Impulse control timer (2-minute delay between evaluation and execution)
- ✅ Bucket cooldown (pauses entries after consecutive stop-outs)
- ✅ Full audit trail (every decision logged with timestamp)

### Core Features

**Trade Entry System:**
- 🎨 Single-screen UI with all inputs/outputs
- 🚦 GO/NO-GO banner (GREEN/YELLOW/RED with reasons)
- 🧮 Position sizing for stocks and options
- 🔥 Real-time heat management
- ⏱️ 2-minute impulse brake
- ❄️ Bucket cooldown enforcement
- 📝 Append-only decision log

**Risk Controls:**
- Portfolio heat cap (default: 4% of equity)
- Per-bucket heat cap (default: 1.5% of equity)
- Cooldown trigger (2 stop-outs in 20 days → 10-day pause)
- Candidate gating (ticker must be in today's import)
- Impulse timer (2 minutes between evaluate and save)

**Data Tracking:**
- **Decisions:** Every trade decision logged (timestamp, inputs, heat, banner status)
- **Positions:** Open positions tracker (units, R per unit, add prices)
- **Candidates:** Daily ticker imports from FINVIZ
- **Buckets:** Sector cooldown status
- **Presets:** FINVIZ screener configurations

### Trading Methodology

Based on **Seykota/Turtle trend-following** principles:
- **Entry:** 52-week highs (or 55-bar Donchian breakout), price above 50/200 SMA
- **Exit:** 10-bar opposite Donchian OR 2×N ATR stop (whichever is closer)
- **Position Sizing:** Risk 0.5-0.75% per unit using 2×N stop distance
- **Pyramiding:** Add every 0.5×N up to max 4 units (only to winners)
- **Diversification:** 6 correlation buckets (Tech/Comm, Consumer, Financials, Industrials, Energy/Materials, Defensives/REITs)

---

## Quick Start

### Prerequisites

**For VBA-Only (Recommended):**
- Excel 2016+ (Windows or Mac)
- Macros enabled
- 10 minutes setup time

**For VBA + Python (Optional):**
- Microsoft 365 with Python enabled
- Python 3.8+
- `pywin32` package (for automation)
- Additional 30 minutes setup time

### Installation Options

#### Option 1: Automated Build (Windows Only)

```batch
# Double-click this file or run from command prompt:
BUILD_WITH_PYTHON.bat
```

**What it does:**
1. Creates Python virtual environment (if needed)
2. Installs dependencies (`pywin32`)
3. Runs `import_to_excel.py` to create workbook
4. Imports all VBA modules automatically
5. Creates all sheets/tables/named ranges
6. Leaves Excel open for you to review

**Result:** `TrendFollowing_TradeEntry.xlsm` with all modules imported and structure created.

#### Option 2: Manual Setup (All Platforms)

See [Setup Instructions](#setup-instructions) below for step-by-step guide.

### First Test Trade

Once workbook is built:

```gherkin
Feature: First Test Trade
  As a new user
  I want to execute a complete trade workflow
  So that I can verify the system works correctly

  Scenario: Successful trade entry with all gates passing
    Given I have opened "TrendFollowing_TradeEntry.xlsm"
    And I have run the "EnsureStructure" macro
    And the TradeEntry sheet exists with UI controls

    When I add a test ticker to Candidates table
      | Ticker | DateImported | Sector     | Bucket    |
      | AAPL   | 2025-10-26   | Technology | Tech/Comm |

    And I fill the Trade Entry inputs
      | Field       | Value         |
      | Ticker      | AAPL          |
      | Entry Price | 180.00        |
      | ATR (N)     | 3.50          |
      | K Multiple  | 2.0           |
      | Method      | Stock         |
      | Bucket      | Tech/Comm     |

    And I check all 6 checklist boxes
      | Box            | Status  |
      | FromPreset     | Checked |
      | TrendPass      | Checked |
      | LiquidityPass  | Checked |
      | TVConfirm      | Checked |
      | EarningsOK     | Checked |
      | JournalOK      | Checked |

    And I click the "Evaluate" button

    Then the banner should display "GREEN: OK TO TRADE"
    And the sizing outputs should populate
      | Output        | Expected Value |
      | R Risk        | ~75.00         |
      | Stop Price    | 173.00         |
      | Shares        | 10             |

    When I wait 2 minutes
    And I click the "Save Decision" button

    Then a new row should appear in tblDecisions
    And a new row should appear in tblPositions
    And the banner should display success message
```

---

## Project Structure

### Root Directory

```
/home/kali/excel-trading-workflow/
│
├── BUILD_WITH_PYTHON.bat          # Automated build script (Windows)
├── import_to_excel.py              # Python automation script
├── scripts/                        # Utility scripts
│   ├── setup_venv.bat             # Virtual environment setup
│   └── deprecated/                # Old scripts (archived)
│
├── VBA/                            # VBA modules for manual import
│   ├── TF_Utils.bas               # Helper functions (NormalizeTicker, etc.)
│   ├── TF_Data.bas                # Structure setup + heat calculations
│   ├── TF_UI.bas                  # Checklist, sizing, save logic
│   ├── TF_Presets.bas             # FINVIZ integration (VBA-only)
│   ├── TF_Config.bas              # Configuration constants
│   ├── TF_PortfolioHeat.bas       # Heat calculation logic
│   ├── TF_PositionSizing.bas      # Sizing calculations
│   ├── TF_Cooldown.bas            # Bucket cooldown logic
│   ├── TF_Validation.bas          # Input validation
│   ├── TF_Logging.bas             # Decision/position logging
│   ├── ThisWorkbook.cls           # Workbook event handlers
│   └── Sheet_TradeEntry.cls       # Sheet event handlers
│
├── Python/                         # Python enhancement modules (optional)
│   ├── finviz_scraper.py          # Auto-scrape FINVIZ (web scraping)
│   ├── heat_calculator.py         # Fast heat calculations (pandas)
│   ├── TF_Python_Bridge.bas       # VBA-Python integration layer
│   ├── TF_Presets_Enhanced.bas    # Enhanced preset import with Python
│   └── requirements.txt           # Python dependencies
│
├── docs/                           # Documentation (organized)
│   ├── setup/                     # Setup guides
│   │   ├── GETTING_STARTED.md    # Quick start guide
│   │   ├── VBA_SETUP_GUIDE.md    # Manual VBA setup (detailed)
│   │   ├── PYTHON_SETUP_GUIDE.md # Python integration guide
│   │   └── QUICK_START.md        # Condensed setup instructions
│   │
│   ├── specifications/            # System specifications
│   │   ├── newest-Interactive_TF_Workbook_Plan.md  # Master plan
│   │   ├── workflow-plan.md                        # Trading rules
│   │   ├── diversification-across-sectors.md       # Bucket framework
│   │   └── SeykotaTurtleTrend-FollowingOptionsExecution+TradingViewStrategyGuide.md
│   │
│   ├── reference/                 # Reference documentation
│   │   ├── VBA_README.md         # VBA architecture & reference
│   │   ├── TWO_BUILD_OPTIONS.md  # Build comparison
│   │   ├── TROUBLESHOOTING_BUILD_ISSUES.md
│   │   └── older-Options_Trend_Dashboard_Summary.md
│   │
│   └── archive/                   # Archived/outdated docs
│       ├── BUILD_COMPLETE.md
│       ├── BUILD_STATUS.md
│       └── (19+ other archived files)
│
├── logs/                           # Build and execution logs
│   └── archive/                   # Old log files
│
├── venv/                           # Python virtual environment (auto-created)
│
├── TrendFollowing_TradeEntry.xlsm # Generated workbook (after build)
│
├── README.md                       # This file
└── CLAUDE.md                       # System overview for AI assistance
```

### Key Files

#### Build & Automation

| File | Purpose | When to Use |
|------|---------|-------------|
| `BUILD_WITH_PYTHON.bat` | Automated workbook creation | **START HERE** (Windows users) |
| `import_to_excel.py` | Python script to build workbook | Called by batch file, or run directly |
| `scripts/setup_venv.bat` | Creates Python virtual environment | Called automatically by BUILD_WITH_PYTHON.bat |

#### VBA Modules (Production Code)

| Module | Lines | Purpose | Key Functions |
|--------|-------|---------|---------------|
| `TF_Utils.bas` | ~150 | Helper utilities | `NormalizeTicker`, `SheetExists`, `GetOrCreateSheet` |
| `TF_Data.bas` | ~300 | Data structure setup | `EnsureStructure`, `CreateTables`, `SeedBuckets` |
| `TF_UI.bas` | ~400 | UI interactions | `EvaluateChecklist`, `RecalcSizing`, `SaveDecision` |
| `TF_Presets.bas` | ~150 | FINVIZ integration | `OpenPreset`, `ImportCandidatesPrompt` |
| `TF_Config.bas` | ~50 | Configuration constants | Named range definitions |
| `TF_PortfolioHeat.bas` | ~100 | Heat calculations | `PortfolioHeatAfter`, `BucketHeatAfter` |
| `TF_PositionSizing.bas` | ~150 | Position sizing | `CalcStockShares`, `CalcOptContracts` |
| `TF_Cooldown.bas` | ~100 | Cooldown logic | `IsBucketInCooldown`, `UpdateCooldowns` |
| `TF_Validation.bas` | ~80 | Input validation | `ValidateInputs`, `CheckHeatCaps` |
| `TF_Logging.bas` | ~120 | Audit trail | `LogDecision`, `LogPosition` |
| `ThisWorkbook.cls` | ~50 | Workbook events | `Workbook_Open`, `EnsureEventsOn` |
| `Sheet_TradeEntry.cls` | ~100 | Sheet events | `Worksheet_Change`, auto-refresh banner |

**Total VBA:** ~1,750 lines of production code across 12 modules.

#### Python Modules (Optional Enhancement)

| File | Lines | Purpose | Performance Gain |
|------|-------|---------|------------------|
| `finviz_scraper.py` | ~250 | Web scraping FINVIZ | 5x faster than manual copy/paste |
| `heat_calculator.py` | ~200 | Pandas-based heat calcs | 15-30x faster for large portfolios |
| `TF_Python_Bridge.bas` | ~150 | VBA-Python bridge | N/A (integration layer) |
| `TF_Presets_Enhanced.bas` | ~180 | Python-enhanced preset import | Uses Python scraper automatically |

**Total Python:** ~780 lines (optional).

#### Documentation

| Document | Category | Audience | Length |
|----------|----------|----------|--------|
| `GETTING_STARTED.md` | Setup | New users | 15 min read |
| `VBA_SETUP_GUIDE.md` | Setup | Manual setup users | 30 min read |
| `PYTHON_SETUP_GUIDE.md` | Setup | Advanced users | 20 min read |
| `newest-Interactive_TF_Workbook_Plan.md` | Spec | Developers/AI | 60 min read (detailed) |
| `workflow-plan.md` | Spec | Traders | 15 min read |
| `VBA_README.md` | Reference | Developers | 20 min read |
| `CLAUDE.md` | Reference | AI assistants | 10 min read |

---

## Setup Instructions

### Path A: Automated Build (Windows, Recommended)

```gherkin
Feature: Automated Workbook Build
  As a Windows user
  I want to build the workbook automatically
  So that I can start trading quickly

  Scenario: Successful automated build
    Given I have Python 3.8+ installed
    And I am in the project root directory

    When I double-click "BUILD_WITH_PYTHON.bat"

    Then the script should:
      | Step | Action |
      | 1 | Create Python virtual environment in ./venv/ |
      | 2 | Install pywin32 package |
      | 3 | Configure Excel Trust Center settings |
      | 4 | Kill any existing Excel processes |
      | 5 | Delete old workbook (if exists) |
      | 6 | Create new "TrendFollowing_TradeEntry.xlsm" |
      | 7 | Import all 12 VBA modules |
      | 8 | Run EnsureStructure macro |
      | 9 | Create all sheets/tables/named ranges |
      | 10 | Leave Excel open for review |

    And I should see "SUCCESS!" message
    And Excel should be open with the workbook loaded

    When I press Alt+F11

    Then I should see all 12 modules in VBA Editor
      | Module | Type |
      | TF_Utils | Standard Module |
      | TF_Data | Standard Module |
      | TF_UI | Standard Module |
      | TF_Presets | Standard Module |
      | TF_Config | Standard Module |
      | TF_PortfolioHeat | Standard Module |
      | TF_PositionSizing | Standard Module |
      | TF_Cooldown | Standard Module |
      | TF_Validation | Standard Module |
      | TF_Logging | Standard Module |
      | ThisWorkbook | Class Module |
      | Sheet_TradeEntry | Class Module |
```

**Estimated Time:** 2-3 minutes

**If Build Fails:**
1. Check Python installation: `python --version` (should be 3.8+)
2. Check pywin32 installation: `pip show pywin32`
3. Ensure Excel is closed before running
4. Review error messages in console
5. See `docs/reference/TROUBLESHOOTING_BUILD_ISSUES.md`

### Path B: Manual Setup (All Platforms)

```gherkin
Feature: Manual Workbook Setup
  As a Mac user or user without Python
  I want to manually import VBA modules
  So that I can use the system on any platform

  Scenario: Manual VBA import and UI build
    Given I have created a new Excel workbook
    And I have saved it as "TrendFollowing_TradeEntry.xlsm"
    And macros are enabled

    When I press Alt+F11 to open VBA Editor
    And I import all VBA modules from ./VBA/ folder
      | File | Import Method |
      | TF_Utils.bas | File → Import File |
      | TF_Data.bas | File → Import File |
      | TF_UI.bas | File → Import File |
      | TF_Presets.bas | File → Import File |
      | TF_Config.bas | File → Import File |
      | TF_PortfolioHeat.bas | File → Import File |
      | TF_PositionSizing.bas | File → Import File |
      | TF_Cooldown.bas | File → Import File |
      | TF_Validation.bas | File → Import File |
      | TF_Logging.bas | File → Import File |
      | ThisWorkbook.cls | File → Import File |
      | Sheet_TradeEntry.cls | File → Import File |

    And I press F5 to open Run Macro dialog
    And I run the "EnsureStructure" macro

    Then I should see 7 sheets created
      | Sheet Name | Purpose |
      | TradeEntry | Main UI |
      | Summary | Settings and named ranges |
      | Presets | FINVIZ screener configs |
      | Buckets | Sector/bucket definitions |
      | Candidates | Daily ticker imports |
      | Decisions | Append-only trade log |
      | Positions | Open positions tracker |

    And I should see 5 tables created
      | Table Name | Sheet | Rows |
      | tblPresets | Presets | 5 |
      | tblBuckets | Buckets | 6 |
      | tblCandidates | Candidates | 0 |
      | tblDecisions | Decisions | 0 |
      | tblPositions | Positions | 0 |

    And I should see 10+ named ranges in Summary sheet
      | Named Range | Default Value |
      | Equity_E | 10000 |
      | RiskPct_r | 0.0075 |
      | StopMultiple_K | 2.0 |
      | HeatCap_H_pct | 0.04 |
      | BucketHeatCap_pct | 0.015 |
      | AddStepN | 0.5 |

    When I navigate to TradeEntry sheet

    Then I should see the UI layout
      | Section | Elements |
      | Inputs | Dropdowns (Preset, Ticker, Bucket, Method), Entry, N, K, Delta, DTE, MaxLoss |
      | Checklist | 6 checkboxes + descriptions |
      | Buttons | Evaluate, Recalc Sizing, Save Decision, Import Candidates |
      | Outputs | R, Stop Price, Shares/Contracts, Add prices (1-3) |
      | Banner | Large cell displaying GO/NO-GO status |

    # Note: Manual UI build requires adding dropdowns, buttons, and formatting
    # See docs/setup/VBA_SETUP_GUIDE.md Part 2 for detailed instructions
```

**Estimated Time:** 60-90 minutes (including UI build)

**Detailed Instructions:** See `docs/setup/VBA_SETUP_GUIDE.md`

---

## Daily Workflow

### Morning Routine (10-15 minutes)

```gherkin
Feature: Daily Candidate Import
  As a trader
  I want to import potential trade candidates
  So that I can evaluate them for entry

  Background:
    Given I have opened "TrendFollowing_TradeEntry.xlsm"
    And all macros are enabled
    And I am on the TradeEntry sheet

  Scenario: Manual import from FINVIZ (VBA-only)
    When I select "TF_BREAKOUT_LONG" from the Preset dropdown
    And I click the "Import Candidates" button

    Then FINVIZ should open in my default browser
    And I should see the screener results

    When I select all tickers (Ctrl+A on ticker column)
    And I copy them (Ctrl+C)
    And I return to Excel
    And I click OK on the import prompt

    Then I should see a paste dialog

    When I paste the tickers (Ctrl+V)
    And I click OK

    Then tblCandidates should populate with today's imports
    And each row should have DateImported = today
    And the Ticker dropdown should include all imported tickers

  Scenario: Automated import with Python (advanced)
    Given I have Python integration enabled
    And TF_Presets_Enhanced module is imported

    When I select "TF_BREAKOUT_LONG" from the Preset dropdown
    And I click the "Import Candidates" button

    Then Python should auto-scrape FINVIZ
    And tblCandidates should populate automatically (5-10 seconds)
    And I should see a success message

    # No manual copy/paste required!
```

### Per-Trade Workflow (2-3 minutes)

```gherkin
Feature: Trade Entry Decision
  As a trader
  I want to evaluate a potential trade
  So that I can execute only high-quality setups

  Background:
    Given I have imported today's candidates
    And I have validated signals on TradingView
    And I am ready to enter a trade

  Scenario: GREEN banner → Successful trade entry
    When I select "AAPL" from the Ticker dropdown
    And I fill the inputs
      | Field | Value |
      | Entry Price | 180.00 |
      | ATR (N) | 3.50 |
      | K Multiple | 2.0 |
      | Method | Stock |
      | Bucket | Tech/Comm |

    And I check all 6 checklist boxes
      | Box | Reason |
      | FromPreset | Ticker is in today's Candidates |
      | TrendPass | Price > 50/200 SMA on TradingView |
      | LiquidityPass | Avg volume > 1M shares |
      | TVConfirm | TradingView strategy shows long entry |
      | EarningsOK | No earnings within 7 days |
      | JournalOK | I reviewed my trading journal |

    And I click "Evaluate"

    Then the banner should display "GREEN: OK TO TRADE"
    And the timer should start (2-minute countdown)

    When I click "Recalc Sizing"

    Then I should see sizing outputs
      | Output | Formula | Example |
      | R Risk | Equity × RiskPct | $75.00 |
      | Stop Price | Entry - (K × N) | $173.00 |
      | Shares | floor(R / (K × N)) | 10 |
      | Add1 Price | Entry + (AddStepN × N) | $181.75 |
      | Add2 Price | Entry + (2 × AddStepN × N) | $183.50 |
      | Add3 Price | Entry + (3 × AddStepN × N) | $185.25 |

    And I should verify the outputs match my expectations

    When I wait 2 minutes (impulse brake)
    And I click "Save Decision"

    Then a new row should be added to tblDecisions
      | Field | Value |
      | Timestamp | 2025-10-26 14:32:15 |
      | Ticker | AAPL |
      | BannerStatus | GREEN |
      | Shares | 10 |
      | EntryPrice | 180.00 |
      | StopPrice | 173.00 |

    And a new row should be added to tblPositions
      | Field | Value |
      | Ticker | AAPL |
      | UnitsOpen | 1 |
      | RperUnit | 75.00 |
      | TotalOpenR | 75.00 |
      | NextAddPrice | 181.75 |

    And I should execute the trade in my broker
      | Action | Details |
      | Order Type | Buy Limit |
      | Ticker | AAPL |
      | Shares | 10 |
      | Limit Price | 180.00 |
      | Stop Loss | 173.00 |

  Scenario: YELLOW banner → Missing 1 checklist item
    When I select "TSLA" from the Ticker dropdown
    And I fill all inputs correctly
    And I check only 5 of 6 checklist boxes (missing TVConfirm)
    And I click "Evaluate"

    Then the banner should display "YELLOW: REVIEW - 1 item unchecked"
    And the "Save Decision" button should be disabled

    # User must check all boxes to proceed

  Scenario: RED banner → Multiple violations
    When I select "NVDA" from the Ticker dropdown
    And I fill all inputs correctly
    And I check all 6 checklist boxes
    And I click "Evaluate"

    But portfolio heat would exceed 4% after entry
    And the bucket "Tech/Comm" is in active cooldown

    Then the banner should display "RED: BLOCKED - Heat cap exceeded, Cooldown active"
    And the "Save Decision" button should be disabled

    # User cannot proceed until violations are resolved

  Scenario: Impulse timer blocks early save
    Given I have a GREEN banner
    And I clicked "Evaluate" 30 seconds ago

    When I click "Save Decision"

    Then I should see an error message
      | Message | "Impulse brake: Wait 2 minutes after Evaluate before saving" |

    And the decision should NOT be saved

    # User must wait full 2 minutes
```

### Weekly Maintenance (5 minutes)

```gherkin
Feature: Weekly Portfolio Maintenance
  As a trader
  I want to update closed positions and cooldown status
  So that my risk calculations remain accurate

  Scenario: Update closed positions and check cooldowns
    Given I have closed 2 positions this week (stop-outs)

    When I navigate to the Positions sheet
    And I update the Status column to "Closed" for stopped-out trades

    And I press Alt+F11 to open VBA Editor
    And I run the "UpdateCooldowns" macro

    Then tblBuckets should recalculate cooldown status
    And any bucket with ≥2 stop-outs in last 20 days should show
      | Field | Value |
      | CooldownActive | TRUE |
      | CooldownEndsOn | 10 days from last stop-out |

    And I should review adherence metrics
      | Metric | Target | Actual |
      | % GREEN-only saves | 100% | Check tblDecisions |
      | Avg time per trade | 2-3 min | Review timestamps |
      | Impulse violations | 0 | Check decision log |
```

---

## System Capabilities

### Position Sizing Methods

#### Stock Sizing

**Formula:**
```
Shares = floor(R / StopDist)
where:
  R = Equity × RiskPct             (e.g., $10,000 × 0.75% = $75)
  StopDist = K × N                 (e.g., 2.0 × 3.50 = $7.00)
  Shares = floor(75 / 7) = 10
```

**Example:**
```gherkin
Scenario: Calculate stock position size
  Given Equity = $10,000
  And RiskPct = 0.75%
  And Entry = $180.00
  And N (ATR) = $3.50
  And K (stop multiple) = 2.0

  When I calculate position size

  Then R = $75.00
  And StopDist = $7.00
  And Stop Price = $173.00
  And Shares = 10
  And Total Position Value = $1,800
  And Position Risk = $70 (10 shares × $7 stop)
```

#### Options Sizing - Delta-ATR Method

**Formula:**
```
Contracts = floor(R / (K × N × Delta × 100))
where:
  R = risk dollars per unit
  K × N = stop distance in underlying
  Delta = option delta (0.60-0.70 for ITM calls)
  100 = shares per contract
```

**Example:**
```gherkin
Scenario: Calculate options contracts (Delta-ATR method)
  Given Equity = $10,000
  And RiskPct = 0.75%
  And Entry = $180.00
  And N (ATR) = $3.50
  And K = 2.0
  And Delta = 0.65

  When I select Method = "Opt-DeltaATR"
  And I calculate position size

  Then R = $75.00
  And Denominator = 2.0 × 3.50 × 0.65 × 100 = 455
  And Contracts = floor(75 / 455) = 0

  # Note: If contracts = 0, position is too small for options
  # Increase account size or use stock instead
```

#### Options Sizing - MaxLoss Method

**Formula:**
```
Contracts = floor(R / (MaxLossPerContract × 100))
where:
  MaxLossPerContract = net debit for debit spread
```

**Example:**
```gherkin
Scenario: Calculate options contracts (MaxLoss method)
  Given Equity = $10,000
  And RiskPct = 0.75%
  And Method = "Opt-MaxLoss"
  And MaxLoss per contract = $2.50 (debit spread)

  When I calculate position size

  Then R = $75.00
  And Contracts = floor(75 / (2.50 × 100)) = floor(75 / 250) = 0

  # Note: MaxLoss method is for defined-risk spreads
```

### Heat Management

#### Portfolio Heat Cap

**Calculation:**
```
Portfolio Heat = Sum of TotalOpenR for all open positions
Portfolio Heat % = Portfolio Heat / Equity

Constraint: Portfolio Heat % ≤ HeatCap_H_pct (default 4%)
```

**Example:**
```gherkin
Scenario: Portfolio heat cap enforcement
  Given Equity = $10,000
  And HeatCap_H_pct = 4.0%
  And I have 4 open positions
    | Ticker | UnitsOpen | RperUnit | TotalOpenR |
    | AAPL   | 2         | 75       | 150        |
    | MSFT   | 1         | 75       | 75         |
    | GOOGL  | 1         | 75       | 75         |
    | AMZN   | 1         | 75       | 75         |

  When I calculate portfolio heat

  Then Portfolio Heat = 150 + 75 + 75 + 75 = $375
  And Portfolio Heat % = 375 / 10000 = 3.75%

  When I attempt to enter a new trade with R = $75

  Then Portfolio Heat After = 375 + 75 = $450
  And Portfolio Heat % After = 450 / 10000 = 4.5%
  And Heat Cap Exceeded = TRUE
  And Banner = RED: BLOCKED - Portfolio heat would exceed 4.0%
```

#### Bucket Heat Cap

**Calculation:**
```
Bucket Heat = Sum of TotalOpenR for all open positions in specified bucket
Bucket Heat % = Bucket Heat / Equity

Constraint: Bucket Heat % ≤ BucketHeatCap_pct (default 1.5%)
```

**Example:**
```gherkin
Scenario: Bucket heat cap enforcement
  Given Equity = $10,000
  And BucketHeatCap_pct = 1.5%
  And Bucket "Tech/Comm" has 2 open positions
    | Ticker | UnitsOpen | RperUnit | TotalOpenR |
    | AAPL   | 1         | 75       | 75         |
    | MSFT   | 1         | 75       | 75         |

  When I calculate bucket heat for "Tech/Comm"

  Then Bucket Heat = 75 + 75 = $150
  And Bucket Heat % = 150 / 10000 = 1.5%

  When I attempt to enter a new "Tech/Comm" trade with R = $75

  Then Bucket Heat After = 150 + 75 = $225
  And Bucket Heat % After = 225 / 10000 = 2.25%
  And Bucket Cap Exceeded = TRUE
  And Banner = RED: BLOCKED - Tech/Comm bucket heat would exceed 1.5%
```

### Bucket Cooldown

**Logic:**
```
IF Bucket has ≥ StopoutsToCooldown (default 2)
   within StopoutsWindowBars (default 20 days)
THEN CooldownActive = TRUE
     CooldownEndsOn = LastStopoutDate + CooldownBars (default 10 days)
ELSE CooldownActive = FALSE
```

**Example:**
```gherkin
Scenario: Bucket cooldown trigger
  Given Bucket "Tech/Comm" has parameters
    | Parameter | Value |
    | StopoutsToCooldown | 2 |
    | StopoutsWindowBars | 20 |
    | CooldownBars | 10 |

  And tblDecisions contains recent Tech/Comm stop-outs
    | Date | Ticker | Outcome |
    | 2025-10-15 | AAPL | StopOut |
    | 2025-10-20 | MSFT | StopOut |
    | 2025-10-22 | GOOGL | Still Open |

  When I run "UpdateCooldowns" macro

  Then Bucket "Tech/Comm" should show
    | Field | Value |
    | CooldownActive | TRUE |
    | CooldownEndsOn | 2025-10-30 (10 days after 2025-10-20) |

  When I attempt to enter a new "Tech/Comm" trade on 2025-10-26

  Then Banner = RED: BLOCKED - Tech/Comm bucket in cooldown until 2025-10-30
  And Save Decision button = disabled

  When I attempt to enter a new "Consumer" trade on 2025-10-26

  Then Cooldown check = PASS (different bucket)
  And Banner = GREEN (assuming all other checks pass)
```

### Impulse Brake

**Logic:**
```
Evaluate button → Store timestamp in Control!A1
Save button → Check if Now() - Control!A1 >= 2 minutes

IF elapsed < 2 minutes THEN
  Display error: "Wait 2 minutes after Evaluate"
  Block save
ELSE
  Proceed with save
END IF
```

**Example:**
```gherkin
Scenario: Impulse brake enforcement
  Given I am on TradeEntry sheet
  And all inputs are filled correctly
  And checklist is complete

  When I click "Evaluate" at 14:30:00

  Then Control!A1 = 14:30:00
  And Banner = GREEN: OK TO TRADE

  When I click "Save Decision" at 14:30:45

  Then Elapsed = 45 seconds
  And Error = "Impulse brake: Wait 2 minutes after Evaluate"
  And Decision NOT saved

  When I click "Save Decision" at 14:32:15

  Then Elapsed = 2 minutes 15 seconds
  And Save proceeds normally
  And Decision logged to tblDecisions
```

---

## Testing & Validation

### Acceptance Test Scenarios

```gherkin
Feature: Trade Entry System Validation
  As a developer or QA tester
  I want to validate all system constraints
  So that the system enforces trading rules correctly

Background:
  Given I have opened "TrendFollowing_TradeEntry.xlsm"
  And EnsureStructure has been run
  And Summary named ranges are set to defaults
  And TradeEntry UI exists

Scenario: All checks pass → GREEN banner
  Given I have added "AAPL" to tblCandidates with today's date
  And Portfolio heat is currently 2.0% ($200)
  And Bucket "Tech/Comm" heat is currently 0.5% ($50)
  And Bucket "Tech/Comm" cooldown is inactive

  When I fill inputs
    | Ticker | Entry | N | K | Method | Bucket |
    | AAPL | 180 | 3.5 | 2.0 | Stock | Tech/Comm |

  And I check all 6 checklist boxes
  And I click "Evaluate"

  Then ChecklistPassed = TRUE (6 of 6 checked)
  And TickerInCandidates = TRUE
  And PortfolioHeatOK = TRUE (2.75% after entry < 4.0%)
  And BucketHeatOK = TRUE (1.25% after entry < 1.5%)
  And CooldownOK = TRUE
  And Banner = "GREEN: OK TO TRADE"
  And Save button enabled after 2 minutes

Scenario: 1 checklist item missing → YELLOW banner
  Given setup from previous scenario

  When I check only 5 of 6 boxes (missing TVConfirm)
  And I click "Evaluate"

  Then ChecklistPassed = FALSE (5 of 6)
  And Banner = "YELLOW: REVIEW - 1 item unchecked"
  And Save button = disabled

Scenario: 2+ checklist items missing → RED banner
  Given setup from previous scenario

  When I check only 4 of 6 boxes
  And I click "Evaluate"

  Then ChecklistPassed = FALSE (4 of 6)
  And Banner = "RED: BLOCKED - 2 items unchecked"
  And Save button = disabled

Scenario: Ticker not in today's Candidates → RED banner
  Given tblCandidates contains only "MSFT" with today's date
  And I attempt to trade "AAPL"
  And all other checks would pass

  When I click "Evaluate"

  Then TickerInCandidates = FALSE
  And Banner = "RED: BLOCKED - Ticker not in today's Candidates"
  And Save button = disabled

Scenario: Portfolio heat cap exceeded → RED banner
  Given Equity = $10,000
  And HeatCap_H_pct = 4.0%
  And Portfolio heat is currently 3.8% ($380)
  And I attempt to enter a trade with R = $75

  When I click "Evaluate"

  Then PortfolioHeatAfter = 380 + 75 = $455 (4.55%)
  And PortfolioHeatOK = FALSE (4.55% > 4.0%)
  And Banner = "RED: BLOCKED - Portfolio heat would exceed 4.0%"
  And Save button = disabled

Scenario: Bucket heat cap exceeded → RED banner
  Given Equity = $10,000
  And BucketHeatCap_pct = 1.5%
  And Bucket "Tech/Comm" heat is currently 1.3% ($130)
  And I attempt to enter a "Tech/Comm" trade with R = $75

  When I click "Evaluate"

  Then BucketHeatAfter = 130 + 75 = $205 (2.05%)
  And BucketHeatOK = FALSE (2.05% > 1.5%)
  And Banner = "RED: BLOCKED - Tech/Comm bucket heat would exceed 1.5%"
  And Save button = disabled

Scenario: Bucket in cooldown → RED banner
  Given Bucket "Tech/Comm" has
    | Field | Value |
    | CooldownActive | TRUE |
    | CooldownEndsOn | 2025-10-30 |

  And Today is 2025-10-26
  And I attempt to enter a "Tech/Comm" trade
  And all other checks would pass

  When I click "Evaluate"

  Then CooldownOK = FALSE
  And Banner = "RED: BLOCKED - Tech/Comm bucket in cooldown until 2025-10-30"
  And Save button = disabled

Scenario: Impulse timer not elapsed → Save blocked
  Given I have GREEN banner
  And I clicked "Evaluate" 90 seconds ago

  When I click "Save Decision"

  Then Error = "Impulse brake: Wait 2 minutes after Evaluate"
  And Decision NOT saved

  When I wait another 30 seconds (total 120 seconds)
  And I click "Save Decision"

  Then Save proceeds normally
  And tblDecisions has new row
  And tblPositions has new row
```

### Unit Test Examples

Run these in VBA Immediate Window (Ctrl+G in VBA Editor):

```vba
' Test structure creation
Call EnsureStructure
Debug.Print "Sheets created: " & ThisWorkbook.Sheets.Count

' Test heat calculations
Debug.Print "Portfolio heat after $75 entry: " & PortfolioHeatAfter(75)
Debug.Print "Tech/Comm bucket heat after $75 entry: " & BucketHeatAfter("Tech/Comm", 75)

' Test cooldown check
Debug.Print "Is Tech/Comm in cooldown? " & IsBucketInCooldown("Tech/Comm")

' Test sizing calculations
Debug.Print "Stock shares: " & CalcStockShares(10000, 0.0075, 2.0, 3.5)
Debug.Print "Opt contracts (Delta-ATR): " & CalcOptContracts_DeltaATR(75, 2.0, 3.5, 0.65)

' Test validation
Debug.Print "Input validation: " & ValidateInputs()

' Test ticker normalization
Debug.Print "Normalized: " & NormalizeTicker("  aapl  ") ' Should return "AAPL"
```

---

## Key Concepts

### Position Sizing Philosophy

**Why risk 0.5-0.75% per unit?**
- Allows 4 units per position (max pyramid)
- Max exposure per position = 2-3% of equity
- Survives 30-40 consecutive stop-outs before account wipeout
- Matches Turtle trader risk parameters

**Why 2×N stop distance (K=2.0)?**
- Wider stops than 1N (less noise-based stop-outs)
- Narrower than 3N (limits per-trade risk)
- Balances stop-out frequency vs. risk per trade

**Why pyramid every 0.5×N?**
- Frequent enough to capture strong trends (4-6 adds in 20-30% move)
- Infrequent enough to avoid over-trading
- Each add is smaller than previous (shares = R / new stop distance)

### Heat Caps Philosophy

**Why 4% portfolio cap?**
- Conservative enough to survive drawdowns
- Aggressive enough to participate in trends
- Allows 5-6 simultaneous positions at 0.75% risk each
- Matches professional money management standards

**Why 1.5% bucket cap?**
- Limits sector concentration
- Allows 1-2 positions per bucket
- Prevents overexposure to single macro theme
- Forces diversification across uncorrelated buckets

### Cooldown Philosophy

**Why cooldown after 2 stop-outs?**
- Prevents "revenge trading" in weak sectors
- Acknowledges that sector may be in correction
- Forced pause to reassess market conditions
- Reduces emotional decision-making

**Why 10-day cooldown?**
- Long enough to change market regime (1-2 weeks)
- Short enough to re-enter if trend resumes
- Balances discipline vs. opportunity cost

### Impulse Brake Philosophy

**Why 2-minute delay?**
- Long enough to prevent FOMO-driven entries
- Short enough to avoid missing fills
- Forces conscious pause to review decision
- Reduces impulsive "chart-chasing" behavior

---

## Customization

### Settings Customization

All parameters are customizable in the Summary sheet via named ranges:

```gherkin
Feature: System Customization
  As an experienced trader
  I want to adjust risk parameters
  So that the system matches my risk tolerance

  Scenario: Increase position size for larger account
    Given Equity_E = $50,000 (increased from $10,000)
    And RiskPct_r = 1.0% (increased from 0.75%)

    When I enter a new trade

    Then R = 50000 × 0.01 = $500 (vs. $75 previously)
    And position sizes scale proportionally

  Scenario: Widen stops to reduce stop-out frequency
    Given StopMultiple_K = 3.0 (increased from 2.0)

    When I calculate stop price

    Then StopDist = 3.0 × N (vs. 2.0 × N previously)
    And stops are 50% wider
    And position size automatically reduces (same R, wider stop)

  Scenario: Tighten heat caps for conservative trading
    Given HeatCap_H_pct = 2.0% (decreased from 4.0%)
    And BucketHeatCap_pct = 0.75% (decreased from 1.5%)

    When I attempt to enter multiple positions

    Then system blocks entries more aggressively
    And max simultaneous positions = 2-3 (vs. 5-6 previously)

  Scenario: Adjust pyramid frequency
    Given AddStepN = 1.0 (increased from 0.5)

    When I calculate add prices

    Then adds occur less frequently (every 1×N vs. 0.5×N)
    And I have fewer add opportunities in trends
```

### Bucket Customization

Modify bucket parameters in tblBuckets:

| Bucket | StopoutsToCooldown | StopoutsWindowBars | CooldownBars | Notes |
|--------|--------------------|--------------------|--------------|-------|
| Tech/Comm | 2 | 20 | 10 | Default |
| Consumer | 3 | 30 | 15 | More lenient (consumer staples less volatile) |
| Financials | 2 | 15 | 20 | Stricter (financials prone to cascading failures) |
| Industrials | 2 | 20 | 10 | Default |
| Energy/Materials | 1 | 10 | 30 | Very strict (commodity sectors highly volatile) |
| Defensives/REITs | 3 | 40 | 5 | Lenient (low volatility, short cooldown) |

### FINVIZ Preset Customization

Add custom screeners to tblPresets:

```gherkin
Scenario: Add custom FINVIZ preset
  Given I want to screen for "Small-cap breakouts"

  When I create a FINVIZ screener on https://finviz.com/screener.ashx
  And I set filters
    | Filter | Value |
    | Market Cap | Small ($300M - $2B) |
    | Price | Over $10 |
    | Volume | Over 500K |
    | Technical | New High |
    | SMA50 | Price above SMA50 |

  And I copy the URL query string after "?"
  # Example: v=211&f=cap_small,sh_price_o10,sh_avgvol_o500,ta_newhigh,ta_sma50_pa

  And I add a new row to tblPresets
    | Name | QueryString |
    | TF_SMALLCAP_BREAKOUT | v=211&f=cap_small,sh_price_o10,sh_avgvol_o500,ta_newhigh,ta_sma50_pa |

  Then the new preset appears in the Preset dropdown
  And I can use it for daily candidate imports
```

---

## Troubleshooting

### Build Issues

```gherkin
Feature: Troubleshooting Build Failures
  As a user encountering build errors
  I want to diagnose and fix common issues
  So that I can successfully build the workbook

  Scenario: Python not found
    Given I run BUILD_WITH_PYTHON.bat
    And I see error "Python not found"

    When I run "python --version" in command prompt

    Then if error, install Python 3.8+ from python.org
    And ensure "Add Python to PATH" is checked during install
    And restart command prompt
    And run BUILD_WITH_PYTHON.bat again

  Scenario: pywin32 installation fails
    Given I see error "Failed to install pywin32"

    When I manually run "pip install pywin32"

    Then if error, check Python version (must be 3.8+)
    And check internet connection
    And try "python -m pip install --upgrade pip" first
    And retry pywin32 install

  Scenario: Excel Trust Center blocks VBA
    Given I see error "Programmatic access to VBA project is not trusted"

    When I open Excel manually
    And I navigate to File → Options → Trust Center → Trust Center Settings
    And I check "Trust access to the VBA project object model"
    And I click OK

    Then BUILD_WITH_PYTHON.bat should succeed on next run

  Scenario: Workbook locked by Excel process
    Given I see error "Could not delete existing workbook"

    When I run Task Manager (Ctrl+Shift+Esc)
    And I find Excel.exe in Processes tab
    And I right-click → End Task

    Then run BUILD_WITH_PYTHON.bat again

  Scenario: Build succeeds but workbook is empty
    Given BUILD_WITH_PYTHON.bat shows "SUCCESS!"
    But workbook has no VBA modules

    When I check logs in console output

    Then look for "Import failed" messages
    And verify VBA files exist in ./VBA/ folder
    And check file permissions (must be readable)
    And re-run script with admin privileges
```

### Runtime Issues

```gherkin
Feature: Troubleshooting Runtime Errors
  As a user encountering runtime errors
  I want to diagnose and fix common issues
  So that the system operates correctly

  Scenario: "Compile error: Sub or Function not defined"
    Given I click a button and see this error

    Then check that all VBA modules are imported
    And press Alt+F11 → View modules in left pane
    And verify all 12 modules exist
    And if missing, run BUILD_WITH_PYTHON.bat again
    Or manually import missing modules

  Scenario: Buttons don't respond
    Given I click button and nothing happens

    When I right-click button → Assign Macro

    Then verify correct procedure is assigned
      | Button | Macro |
      | Evaluate | EvaluateChecklist |
      | Recalc Sizing | RecalcSizing |
      | Save Decision | SaveDecision |
      | Import Candidates | ImportCandidatesPrompt |

    And if wrong/missing, select correct macro and click OK

  Scenario: Dropdowns show #REF! or empty
    Given dropdowns are broken

    When I run EnsureStructure macro

    Then tables should be recreated
    And dropdowns should reference correct ranges
    And data validation should work

  Scenario: Banner doesn't update automatically
    Given I check/uncheck boxes but banner stays same

    When I check that Sheet_TradeEntry.cls is imported
    And I verify event handlers exist
      | Event | Handler |
      | Worksheet_Change | Auto-calls EvaluateChecklist |

    Then if missing, import Sheet_TradeEntry.cls
    And close/reopen workbook

  Scenario: Save Decision always blocked
    Given banner is GREEN but Save fails

    When I check Decisions table for last entry

    Then verify:
      | Check | Debug Method |
      | Ticker in Candidates | tblCandidates DateImported = today? |
      | Portfolio heat OK | Run ? PortfolioHeatAfter(75) in Immediate Window |
      | Bucket heat OK | Run ? BucketHeatAfter("Tech/Comm", 75) |
      | Cooldown OK | Run ? IsBucketInCooldown("Tech/Comm") |
      | Timer elapsed | Check Control!A1 timestamp vs. Now() |

  Scenario: Python not available (for Python integration)
    Given I see "Python not available" error

    Then check Excel version (must be Microsoft 365)
    And navigate to Data tab → Python button
    And if missing, update to Microsoft 365 Insider channel
    And enable Python in Excel settings
```

### Data Issues

```gherkin
Feature: Troubleshooting Data Problems
  As a user with data inconsistencies
  I want to clean and repair data
  So that calculations are accurate

  Scenario: Heat calculations seem wrong
    Given portfolio heat shows unexpected value

    When I navigate to Positions sheet
    And I check Status column

    Then ensure closed positions have Status = "Closed"
    And only "Open" positions should count toward heat
    And update any stale Status values
    And recalculate

  Scenario: Cooldown stuck active
    Given cooldown shows active but should be expired

    When I check tblBuckets CooldownEndsOn date
    And current date is past CooldownEndsOn

    Then run UpdateCooldowns macro
    And verify CooldownActive updates to FALSE
    And if still stuck, manually set CooldownActive = FALSE

  Scenario: Candidates list is empty
    Given Ticker dropdown is empty

    When I check tblCandidates DateImported column

    Then ensure at least one row has today's date
    And if not, import candidates via Import Candidates button
    Or manually add test ticker with today's date

  Scenario: Decisions table grows too large (>1000 rows)
    Given workbook is slow due to large Decisions table

    When I want to archive old decisions

    Then copy old rows (> 6 months) to separate "Archive" sheet
    And delete from tblDecisions (keep last 6 months only)
    And save archive sheet to separate workbook for history
```

For more troubleshooting, see:
- `docs/reference/TROUBLESHOOTING_BUILD_ISSUES.md`
- `docs/reference/VBA_README.md` (Section: Common Issues)

---

## Architecture & Design

### Design Philosophy

**1. Separation of Concerns**

Each VBA module has a single responsibility:

| Module | Responsibility | Dependencies |
|--------|---------------|--------------|
| TF_Utils | Low-level helpers | None |
| TF_Data | Data structure setup | TF_Utils |
| TF_Config | Configuration constants | None |
| TF_PortfolioHeat | Heat calculations | TF_Utils |
| TF_PositionSizing | Sizing calculations | TF_Config |
| TF_Cooldown | Cooldown logic | TF_Utils |
| TF_Validation | Input validation | TF_PortfolioHeat, TF_Cooldown |
| TF_Logging | Audit trail | TF_Utils |
| TF_Presets | FINVIZ integration | TF_Utils, TF_Logging |
| TF_UI | User interface | All of the above |

**2. Idempotent Operations**

`EnsureStructure` can be run multiple times safely:
- Creates sheets/tables only if missing
- Never deletes existing data
- Updates structure to match latest spec

**3. Append-Only Logging**

`tblDecisions` is **append-only**:
- Never edit/delete rows (audit trail integrity)
- Every decision timestamped
- Immutable history for compliance/review

**4. Hard Gates (No Overrides)**

5 hard gates in `SaveDecision`:
1. Banner must be GREEN
2. Ticker must be in today's Candidates
3. 2-minute timer must have elapsed
4. Portfolio heat cap must not be exceeded
5. Bucket heat cap must not be exceeded
6. Bucket must not be in cooldown

**No backdoors, no admin overrides.** This is intentional.

**5. VBA + Optional Python**

Core system is **100% VBA** (works offline, no cloud).

Python is **optional enhancement**:
- Faster candidate import (5x speedup)
- Faster heat calculations (15-30x for large portfolios)
- But requires Microsoft 365 + cloud execution

### Three-Layer Architecture

```
┌─────────────────────────────────────────────────────────────┐
│ LAYER 1: SIGNAL GENERATION (External)                      │
├─────────────────────────────────────────────────────────────┤
│ • FINVIZ Screeners → Candidate tickers                     │
│ • TradingView Strategy → Entry/exit signals                │
│ • Manual validation → Earnings, liquidity, sector          │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ LAYER 2: DECISION SYSTEM (Excel VBA)                       │
├─────────────────────────────────────────────────────────────┤
│ • TradeEntry UI → Input capture                            │
│ • Checklist → Mechanical gates                             │
│ • Sizing → Position calculations                           │
│ • Heat → Portfolio/bucket risk                             │
│ • Cooldown → Sector pause logic                            │
│ • Impulse brake → 2-minute delay                           │
│ • Logging → Audit trail                                    │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ LAYER 3: EXECUTION (Manual)                                │
├─────────────────────────────────────────────────────────────┤
│ • Broker platform → Place orders                           │
│ • Trade management → Stop adjustment, adds, exits          │
│ • Position tracking → Update Positions table               │
└─────────────────────────────────────────────────────────────┘
```

**Why this separation?**

- **Layer 1 (Signal):** External tools are best-in-class for screening and backtesting
- **Layer 2 (Decision):** Excel enforces portfolio constraints that TradingView cannot
- **Layer 3 (Execution):** Manual execution maintains human oversight and prevents automation errors

### Data Flow

```
Daily Workflow:
  FINVIZ → Clipboard → ImportCandidatesPrompt → tblCandidates
  TradingView → Manual validation → TradeEntry inputs
  EvaluateChecklist → Banner status
  RecalcSizing → Sizing outputs
  SaveDecision → tblDecisions + tblPositions

Weekly Workflow:
  Broker → Manual position updates → tblPositions Status
  UpdateCooldowns → tblBuckets CooldownActive

All workflows append to tblDecisions (immutable history)
```

### VBA vs. Python Tradeoffs

| Feature | VBA | Python | Winner |
|---------|-----|--------|--------|
| Offline operation | ✅ Yes | ❌ No (Azure cloud) | VBA |
| Setup complexity | ✅ Simple | ❌ Complex | VBA |
| FINVIZ import | ❌ Manual (30-60s) | ✅ Auto (5-10s) | Python |
| Heat calc speed (10 pos) | ✅ <1s | ✅ <0.1s | Tie |
| Heat calc speed (100 pos) | ❌ ~3s | ✅ <0.2s | Python |
| Data privacy | ✅ Local | ⚠️ Cloud | VBA |
| Excel compatibility | ✅ 2016+ | ❌ 365 only | VBA |
| Cross-platform | ✅ Win/Mac | ❌ Win only (for build) | VBA |

**Recommendation:** Start with VBA-only. Add Python later if needed.

### Testing Strategy

**Manual acceptance tests** (Gherkin scenarios):
- See [Testing & Validation](#testing--validation) section above
- Run full workflow 5-10 times before live trading
- Verify all 5 hard gates block correctly

**Unit tests** (VBA Immediate Window):
- Test each helper function independently
- Verify calculations match expected formulas
- Check edge cases (zero positions, max heat, etc.)

**Integration tests** (real data):
- Import real FINVIZ tickers
- Use real TradingView signals
- Execute test trades with small size
- Verify logs match actual decisions

### Extension Points

**To add new position sizing method:**
1. Add column to TradeEntry sheet for new inputs
2. Add case to `RecalcSizing` in `TF_UI.bas`
3. Implement calculation in `TF_PositionSizing.bas`
4. Update dropdown validation (Method column)

**To add new heat constraint:**
1. Add calculation function to `TF_PortfolioHeat.bas`
2. Add check to `EvaluateChecklist` in `TF_UI.bas`
3. Update banner logic in `TF_Validation.bas`
4. Add gate to `SaveDecision` in `TF_UI.bas`

**To add new bucket:**
1. Add row to tblBuckets (Bucket name + parameters)
2. Map sectors to bucket (Sector column in tblBuckets)
3. Update dropdown validation (Bucket column in TradeEntry)
4. Cooldown logic auto-applies to new bucket

**To customize Python integration:**
1. Modify `finviz_scraper.py` (add new screener sources)
2. Modify `heat_calculator.py` (add custom risk metrics)
3. Update `TF_Python_Bridge.bas` (add new function calls)
4. Replace standard modules with Enhanced versions

### Performance Characteristics

**VBA Workbook Size:**
- Empty structure: ~150 KB
- After 100 decisions: ~500 KB
- After 1,000 decisions: ~2 MB
- After 10,000 decisions: ~15 MB

**Calculation Speed:**
- EnsureStructure: 1-2 seconds (one-time)
- EvaluateChecklist: <0.5 seconds (per trade)
- RecalcSizing: <0.2 seconds (per trade)
- SaveDecision: ~1 second (per trade)
- UpdateCooldowns: 1-3 seconds (weekly)

**Bottlenecks:**
- Large Decisions table (>5,000 rows): Archive old data
- Large Positions table (>100 rows): Clean closed positions
- Cooldown checks (>10 buckets × 1,000 decisions): Optimize query

**Optimization Tips:**
- Keep Decisions table < 1,000 rows (archive annually)
- Keep Positions table < 50 rows (close old positions)
- Run UpdateCooldowns weekly (not daily)

### Future Enhancements

**Potential additions (not yet implemented):**

1. **Earnings calendar integration** (API or manual)
   - Auto-populate EarningsDate column in Candidates
   - Auto-check EarningsOK box if >7 days away

2. **TradingView webhook integration** (requires web server)
   - Auto-import signals from TradingView alerts
   - Populate Entry, N, K from strategy parameters

3. **Broker integration** (via API)
   - Auto-update Positions table from broker account
   - Auto-close positions when exits triggered

4. **Dashboard charts** (Excel charting)
   - Portfolio heat over time
   - Win rate by bucket
   - Adherence metrics (% GREEN-only)

5. **System-1 vs. System-2 toggle** (20/10 vs. 55/10)
   - Dual-system mode like original Turtles
   - Separate heat tracking per system

See `docs/specifications/workflow-plan.md` for full wishlist.

---

## Documentation Map

### For New Users

**Start here:**
1. **README.md** (this file) → System overview + quick start
2. `docs/setup/GETTING_STARTED.md` → Step-by-step first-time setup
3. `docs/setup/VBA_SETUP_GUIDE.md` → Detailed manual setup (if not using automation)

**Then practice:**
4. Run first test trade (see [First Test Trade](#first-test-trade))
5. Import real candidates from FINVIZ
6. Execute 5-10 trades with small size

### For Developers

**Understand architecture:**
1. **CLAUDE.md** → High-level system overview (for AI assistants)
2. `docs/specifications/newest-Interactive_TF_Workbook_Plan.md` → Detailed specification
3. `docs/reference/VBA_README.md` → VBA module architecture

**Understand trading rules:**
4. `docs/specifications/workflow-plan.md` → Trading methodology
5. `docs/specifications/diversification-across-sectors.md` → Bucket framework
6. `docs/specifications/SeykotaTurtleTrend-FollowingOptionsExecution+TradingViewStrategyGuide.md` → TradingView strategy

**Modify code:**
7. Open VBA Editor (Alt+F11)
8. Read module comments (all functions documented)
9. Test changes with Immediate Window (Ctrl+G)

### For Advanced Users (Python)

**After VBA setup is working:**
1. `docs/setup/PYTHON_SETUP_GUIDE.md` → Python integration guide
2. `Python/finviz_scraper.py` → Review scraper code
3. `Python/TF_Python_Bridge.bas` → Review VBA-Python bridge
4. Test Python integration before using in production

### For Troubleshooting

**Build issues:**
- `docs/reference/TROUBLESHOOTING_BUILD_ISSUES.md`
- `docs/reference/TWO_BUILD_OPTIONS.md`

**Runtime issues:**
- [Troubleshooting](#troubleshooting) section above
- `docs/reference/VBA_README.md` (Common Issues section)

**Data issues:**
- Check table contents (tblCandidates, tblPositions, tblBuckets)
- Run unit tests in VBA Immediate Window

---

## Credits & License

### Methodology

Based on **Seykota/Turtle trend-following** principles:
- **Ed Seykota:** Risk management, position sizing, mechanical systems
- **Richard Dennis & William Eckhardt:** Turtle trading rules (Donchian breakouts, pyramiding)
- **Curtis Faith:** *Way of the Turtle* (detailed implementation)
- **Van Tharp:** *Trade Your Way to Financial Freedom* (position sizing, R-multiples)

### Key References

**Books:**
- *The New Market Wizards* by Jack Schwager (Seykota interview)
- *Way of the Turtle* by Curtis Faith
- *Trade Your Way to Financial Freedom* by Van Tharp
- *Following the Trend* by Andreas Clenow

**Tools:**
- **FINVIZ:** Stock screener (https://finviz.com)
- **TradingView:** Charting and strategy backtesting (https://tradingview.com)
- **Excel VBA:** Microsoft Office suite

### System Stats

**Code:**
- VBA: 1,750+ lines across 12 modules
- Python: 780+ lines (optional)
- Documentation: 4,500+ lines (this README + guides)

**Generated:** October 26, 2025

**Maintenance:**
- Quarterly review of trading rules
- Annual archive of old Decisions data
- Update Python dependencies as needed

---

## License

**This is a personal trading system.** Use at your own risk.

**No warranty:**
- Code provided as-is
- No guarantee of profitability
- Past performance ≠ future results
- Trading involves risk of loss

**Recommended usage:**
- Paper trade first (6-12 months)
- Start with small position sizes
- Review all decisions in journal
- Adjust parameters to match your risk tolerance

---

## Next Steps

### Immediate (Today)

```gherkin
Scenario: Get system running
  When I run BUILD_WITH_PYTHON.bat (Windows)
  Or I follow docs/setup/VBA_SETUP_GUIDE.md (Mac/manual)

  Then I should have working workbook
  And I should run first test trade
  And I should verify all 5 gates work correctly
```

### Short-term (This Week)

```gherkin
Scenario: Learn the workflow
  When I import real FINVIZ candidates (5-10 tickers)
  And I validate signals on TradingView
  And I execute 10 test trades (paper trading)

  Then I should understand:
    - How checklist gates work
    - How heat caps enforce limits
    - How impulse timer prevents FOMO
    - How cooldown pauses weak buckets
```

### Medium-term (This Month)

```gherkin
Scenario: Customize for my account
  When I update Summary named ranges
    | Setting | My Value |
    | Equity_E | [my account size] |
    | RiskPct_r | [my risk tolerance] |
    | HeatCap_H_pct | [my max exposure] |

  And I adjust bucket parameters if needed
  And I add custom FINVIZ presets

  Then system should match my trading style
  And I should start live trading (small size)
```

### Long-term (Next Quarter)

```gherkin
Scenario: Optimize and enhance
  When I have 50+ logged decisions

  Then I should review:
    - Win rate by bucket (any patterns?)
    - Heat cap violations (too tight/loose?)
    - Cooldown triggers (effective or not?)
    - Adherence metrics (am I following rules?)

  And I should consider:
    - Python integration (if FINVIZ import is tedious)
    - Custom enhancements (dashboard, charts, etc.)
    - System-2 parameters (55/10 vs. 20/10)
```

---

## Support

**For questions about:**
- **Setup:** See `docs/setup/` guides
- **Trading rules:** See `docs/specifications/workflow-plan.md`
- **VBA code:** See `docs/reference/VBA_README.md`
- **Python integration:** See `docs/setup/PYTHON_SETUP_GUIDE.md`

**For bug reports:**
- Check [Troubleshooting](#troubleshooting) first
- Review VBA code comments (all functions documented)
- Test with unit tests (Immediate Window)
- Use VBA debugger (F8 to step through code)

**For enhancements:**
- Fork the project
- Modify VBA/Python code
- Test thoroughly before live use
- Document changes in code comments

---

**Happy Trading!**

**Remember:**
- ✅ Always follow the GREEN banner
- ✅ Never override the 5 hard gates
- ✅ Let the system enforce discipline
- ✅ Review decisions weekly
- ✅ Adjust parameters as you learn

**The system is designed to protect you from yourself. Let it do its job.**

---

*Last updated: October 26, 2025*
*Version: 1.0*
*Build: Automated (BUILD_WITH_PYTHON.bat)*
