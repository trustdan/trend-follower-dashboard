# Options Trading Enhancement Plan

## Problem Statement

The current TF-Engine system lacks critical options trading metadata needed for proper trend-following with options. Specifically:

1. **No expiration tracking** - Can't track DTE or roll thresholds (21 DTE)
2. **No options strategy types** - Can't distinguish Long Call from Call Spread
3. **No pyramid add-on tracking** - Can't track Entry + 0.5N, Entry + 1N levels
4. **No calendar integration** - Can't display trades in 10-week sector calendar
5. **Limited strategy options** - Only Long/Short Breakout, missing System-1, System-2 distinction
6. **No time exit modes** - Missing "Roll" vs "Close" exit strategies

## Design Principles (from anti-impulsivity.md)

Per the anti-impulsivity document, this is **Seykota/Turtle-style options trading**:

- **55-bar breakout** (System-2) or **20-bar breakout** (System-1)
- **Options: 60-90 DTE** at entry, **roll/close at ~21 DTE**
- **Pyramid every 0.5N** up to **max 4 units**
- **Exit by 10-bar opposite Donchian OR 2×N stop**
- **Time exit modes**: None, Close, or Roll

## Required Strategy Types

### Breakout Strategies (from Ed-Seykota.pine)

1. **System-1 (20/10)** - 20-bar breakout, 10-bar exit (faster, more frequent)
2. **System-2 (55/10)** - 55-bar breakout, 10-bar exit (slower, higher quality) **[DEFAULT]**
3. **Custom** - Manual parameters

### Options Structures (Comprehensive List)

#### Directional (Simple)
1. **Long Call** - Bullish directional play
2. **Long Put** - Bearish directional play

#### Income Generation
3. **Covered Call** - Own stock, sell call (income on holdings)
4. **Cash-Secured Put** - Sell put with cash backing (acquire stock at discount)

#### Vertical Credit Spreads
5. **Bull Put Credit Spread** - Bullish assumption, collect premium
6. **Bear Call Credit Spread** - Bearish assumption, collect premium

#### Butterflies & Condors (Neutral)
7. **Iron Butterfly** - Sell ATM straddle, buy wings (low IV, range-bound)
8. **Iron Condor** - Sell OTM strangle, buy wings (wider range)
9. **Long Put Butterfly** - Debit butterfly with puts
10. **Long Call Butterfly** - Debit butterfly with calls
11. **Inverse Iron Butterfly** - Buy ATM straddle, sell wings (high IV expected)
12. **Inverse Iron Condor** - Buy OTM strangle, sell wings
13. **Short Put Butterfly** - Credit butterfly with puts
14. **Short Call Butterfly** - Credit butterfly with calls

#### Time Spreads (Calendar & Diagonal)
15. **Calendar Call Spread** - Buy longer-dated, sell shorter-dated call (same strike)
16. **Calendar Put Spread** - Buy longer-dated, sell shorter-dated put (same strike)
17. **Diagonal Call Spread** - Calendar + vertical (different strikes)
18. **Diagonal Put Spread** - Calendar + vertical (different strikes)

#### Vertical Debit Spreads
19. **Bull Call Spread** - Buy lower strike call, sell higher strike call
20. **Bear Put Spread** - Buy higher strike put, sell lower strike put

#### Volatility Plays
21. **Straddle** - Buy ATM call + ATM put (big move expected)
22. **Strangle** - Buy OTM call + OTM put (cheaper than straddle)

#### Ratio & Broken Wing
23. **Call Ratio Backspread** - Sell ITM calls, buy more OTM calls (bullish, limited risk)
24. **Put Broken Wing** - Butterfly with strikes skewed to one side
25. **Put Ratio Backspread** - Sell ITM puts, buy more OTM puts (bearish, limited risk)
26. **Call Broken Wing** - Butterfly with strikes skewed to one side

## Database Schema Changes

### Phase 1: Extend `trade_sessions` Table

Add these columns to `trade_sessions`:

```sql
-- Options metadata
options_strategy TEXT,           -- See strategy constants below
entry_date TEXT,                 -- ISO date when trade entered (YYYY-MM-DD)
primary_expiration_date TEXT,    -- Primary expiration (or nearest for calendars/diagonals)
dte INTEGER,                     -- Days to expiration at entry (for primary leg)
roll_threshold_dte INTEGER DEFAULT 21,  -- DTE to roll/close (default 21)
time_exit_mode TEXT DEFAULT 'Close',    -- 'None', 'Close', 'Roll'

-- Multi-leg structure (stored as JSON for flexibility)
legs_json TEXT,                  -- JSON array of legs: [{type, strike, expiration, contracts, action}]
                                 -- Example: [{"type":"CALL","strike":180,"exp":"2025-12-19","qty":1,"action":"BUY"}]

-- Aggregate pricing
net_debit REAL,                  -- Total debit paid (negative = credit received)
max_profit REAL,                 -- Maximum theoretical profit
max_loss REAL,                   -- Maximum theoretical loss (at expiration)
breakeven_lower REAL,            -- Lower breakeven price (NULL if none)
breakeven_upper REAL,            -- Upper breakeven price (NULL if none)
underlying_at_entry REAL,        -- Stock price at entry

-- Pyramiding (Van Tharp method)
max_units INTEGER DEFAULT 4,     -- Maximum pyramid units (default 4)
add_step_n REAL DEFAULT 0.5,     -- Add every X * N (default 0.5N)
current_units INTEGER DEFAULT 0, -- Current units in position (0-4)
add_price_1 REAL,                -- Entry + 0.5N
add_price_2 REAL,                -- Entry + 1.0N
add_price_3 REAL,                -- Entry + 1.5N

-- Breakout system parameters (for documentation)
entry_lookback INTEGER,          -- 20 or 55 for System-1/System-2
exit_lookback INTEGER DEFAULT 10, -- 10-bar exit
```

### Strategy Constants (Backend)

```go
// Options strategy types
const (
    // Directional
    StrategyLongCall           = "LONG_CALL"
    StrategyLongPut            = "LONG_PUT"

    // Income
    StrategyCoveredCall        = "COVERED_CALL"
    StrategyCashSecuredPut     = "CASH_SECURED_PUT"

    // Vertical Credit Spreads
    StrategyBullPutSpread      = "BULL_PUT_SPREAD"
    StrategyBearCallSpread     = "BEAR_CALL_SPREAD"

    // Butterflies & Condors
    StrategyIronButterfly      = "IRON_BUTTERFLY"
    StrategyIronCondor         = "IRON_CONDOR"
    StrategyLongPutButterfly   = "LONG_PUT_BUTTERFLY"
    StrategyLongCallButterfly  = "LONG_CALL_BUTTERFLY"
    StrategyInverseIronButterfly = "INVERSE_IRON_BUTTERFLY"
    StrategyInverseIronCondor  = "INVERSE_IRON_CONDOR"
    StrategyShortPutButterfly  = "SHORT_PUT_BUTTERFLY"
    StrategyShortCallButterfly = "SHORT_CALL_BUTTERFLY"

    // Time Spreads
    StrategyCalendarCallSpread = "CALENDAR_CALL_SPREAD"
    StrategyCalendarPutSpread  = "CALENDAR_PUT_SPREAD"
    StrategyDiagonalCallSpread = "DIAGONAL_CALL_SPREAD"
    StrategyDiagonalPutSpread  = "DIAGONAL_PUT_SPREAD"

    // Vertical Debit Spreads
    StrategyBullCallSpread     = "BULL_CALL_SPREAD"
    StrategyBearPutSpread      = "BEAR_PUT_SPREAD"

    // Volatility
    StrategyStraddle           = "STRADDLE"
    StrategyStrangle           = "STRANGLE"

    // Ratio & Broken Wing
    StrategyCallRatioBackspread = "CALL_RATIO_BACKSPREAD"
    StrategyPutBrokenWing      = "PUT_BROKEN_WING"
    StrategyPutRatioBackspread = "PUT_RATIO_BACKSPREAD"
    StrategyCallBrokenWing     = "CALL_BROKEN_WING"
)
```

### Phase 2: Extend `positions` Table

Add similar fields to `positions` for tracking active trades:

```sql
ALTER TABLE positions ADD COLUMN options_strategy TEXT;
ALTER TABLE positions ADD COLUMN entry_date TEXT;
ALTER TABLE positions ADD COLUMN primary_expiration_date TEXT;
ALTER TABLE positions ADD COLUMN dte INTEGER;
ALTER TABLE positions ADD COLUMN legs_json TEXT;           -- JSON array of legs
ALTER TABLE positions ADD COLUMN net_debit REAL;           -- Debit paid (or credit received)
ALTER TABLE positions ADD COLUMN max_profit REAL;
ALTER TABLE positions ADD COLUMN max_loss REAL;
ALTER TABLE positions ADD COLUMN breakeven_lower REAL;
ALTER TABLE positions ADD COLUMN breakeven_upper REAL;
ALTER TABLE positions ADD COLUMN underlying_at_entry REAL;
ALTER TABLE positions ADD COLUMN max_units INTEGER DEFAULT 4;
ALTER TABLE positions ADD COLUMN current_units INTEGER DEFAULT 1;
ALTER TABLE positions ADD COLUMN add_step_n REAL DEFAULT 0.5;
```

### Phase 3: Create `trade_history` Table (Calendar View)

New table for calendar display:

```sql
CREATE TABLE IF NOT EXISTS trade_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER,
    ticker TEXT NOT NULL,
    strategy TEXT NOT NULL,              -- LONG_BREAKOUT, SHORT_BREAKOUT
    options_strategy TEXT,               -- LONG_CALL, CALL_SPREAD, etc.
    sector TEXT,                         -- Tech/Comm, Finance, etc.
    entry_date TEXT NOT NULL,            -- YYYY-MM-DD
    expiration_date TEXT,                -- YYYY-MM-DD (NULL for stocks)
    exit_date TEXT,                      -- YYYY-MM-DD (NULL if still open)
    status TEXT NOT NULL DEFAULT 'OPEN', -- OPEN, CLOSED, ROLLED
    dte INTEGER,
    contracts INTEGER,
    shares INTEGER,
    risk_dollars REAL,
    pnl REAL,
    outcome TEXT,                        -- WIN, LOSS, SCRATCH
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES trade_sessions(id)
);

CREATE INDEX idx_trade_history_entry_date ON trade_history(entry_date);
CREATE INDEX idx_trade_history_sector ON trade_history(sector, entry_date);
CREATE INDEX idx_trade_history_status ON trade_history(status);
```

## UI Changes

### 1. Enhanced "Start New Trade Session" Dialog

**Current:**
```
Strategy:
○ Long Breakout (55-bar high breakout)
○ Short Breakout (55-bar low breakdown)
○ Custom (manual setup)
```

**Proposed:**
```
Ticker: [AAPL]

Breakout System:
○ System-2 (55/10) - 55-bar breakout, 10-bar exit [DEFAULT]
○ System-1 (20/10) - 20-bar breakout, 10-bar exit
○ Custom - Manual parameters

Direction:
○ Long (bullish)
○ Short (bearish)

Instrument Type:
○ Stock/ETF (no options)
○ Options (select strategy below)

[If Options selected, show dropdown:]

Options Strategy: [Dropdown organized by category]
  ├─ Directional
  │  ├─ Long Call
  │  └─ Long Put
  ├─ Income
  │  ├─ Covered Call
  │  └─ Cash-Secured Put
  ├─ Vertical Credit Spreads
  │  ├─ Bull Put Credit Spread
  │  └─ Bear Call Credit Spread
  ├─ Butterflies & Condors
  │  ├─ Iron Butterfly
  │  ├─ Iron Condor
  │  ├─ Long Put Butterfly
  │  ├─ Long Call Butterfly
  │  ├─ Inverse Iron Butterfly
  │  ├─ Inverse Iron Condor
  │  ├─ Short Put Butterfly
  │  └─ Short Call Butterfly
  ├─ Time Spreads
  │  ├─ Calendar Call Spread
  │  ├─ Calendar Put Spread
  │  ├─ Diagonal Call Spread
  │  └─ Diagonal Put Spread
  ├─ Vertical Debit Spreads
  │  ├─ Bull Call Spread
  │  └─ Bear Put Spread
  ├─ Volatility
  │  ├─ Straddle
  │  └─ Strangle
  └─ Ratio & Broken Wing
     ├─ Call Ratio Backspread
     ├─ Put Broken Wing
     ├─ Put Ratio Backspread
     └─ Call Broken Wing

[After strategy selected, show dynamic leg entry form:]

=== Strategy: Iron Condor ===
Underlying: AAPL @ $175.00
Entry Date: [2025-10-30]
Expiration: [2025-12-19] (50 DTE)

Leg 1: Buy Put  $160 @ $0.50 × [1] contracts
Leg 2: Sell Put  $165 @ $1.20 × [1] contracts
Leg 3: Sell Call $185 @ $1.30 × [1] contracts
Leg 4: Buy Call  $190 @ $0.60 × [1] contracts

Net Credit: $140 per spread
Max Profit: $140 (100% of credit)
Max Loss: $360 (spread width - credit)
Breakeven: $163.60 / $186.40

Roll at DTE: [21]
Time Exit Mode: [Close ▼]

[Start Session] [Cancel]
```

### 2. Position Sizing Screen Enhancements

Add fields after existing sizing calculation:

```
=== Pyramid Planning ===
Max Units: [Input, default: 4]
Add Every: [Input, default: 0.5] × N

📊 Add-On Prices:
  Add 1: $XXX.XX (Entry + 0.5N)
  Add 2: $XXX.XX (Entry + 1.0N)
  Add 3: $XXX.XX (Entry + 1.5N)

Current Units: 0 / 4
```

### 3. Trade Entry Screen Enhancements

Add session to positions when GO is clicked:

```
[Existing 5 gates display]

=== Trade Summary ===
Ticker: AAPL
Strategy: System-2 (55/10) Long Breakout
Options: Long Call $180 strike
Entry: $175.00
Expiration: 2025-12-19 (60 DTE)
Contracts: 5
Risk: $750 (0.75% of equity)

[SAVE GO] [SAVE NO-GO] [Cancel]
```

### 4. NEW: Calendar View (📅 Tab)

```
=== Trade Calendar (10-Week View) ===

Filter: [All Sectors ▼] [All Strategies ▼] [Open Only ☑]

        Week 1    Week 2    Week 3    Week 4    Week 5    Week 6    Week 7    Week 8    Week 9    Week 10
        (Oct 14)  (Oct 21)  (Oct 28)  (Nov 4)   (Nov 11)  (Nov 18)  (Nov 25)  (Dec 2)   (Dec 9)   (Dec 16)

Tech     AAPL      AAPL      AAPL      AAPL
         MSFT      MSFT      MSFT

Finance            JPM       JPM       JPM       JPM

Energy   XLE       XLE       XLE       XLE       XLE       XLE

Materials                    FCX       FCX       FCX

Legend:
🟢 Open position
🔴 Closed (loss)
🟡 Closed (win)
📊 Rolled

Click any cell to see trade details.
```

### 5. Dashboard Enhancements

Show active trades with expiration dates:

```
📊 Active Trades (Open Positions)

┌─────────────────────────────────────────────────────────┐
│ AAPL - System-2 Long Breakout                          │
│ Long Call $180 strike                                   │
│ Entry: Oct 15 | Exp: Dec 19 (65 DTE) | Roll at: 21 DTE │
│ Units: 2/4 | Risk: $1,500                               │
│ Status: Up +5.2% | Next add: $178.50                    │
└─────────────────────────────────────────────────────────┘
```

## Strategy Templates & Builders

To simplify entering complex multi-leg strategies, we'll create **strategy builders** that auto-populate leg structures:

### Builder Functions

```go
// BuildIronCondor creates a 4-leg iron condor structure
func BuildIronCondor(underlying float64, putSpread, callSpread float64, expiration string) []OptionLeg {
    return []OptionLeg{
        {Type: "PUT", Strike: underlying - putSpread - 5, Action: "BUY", Qty: 1, Exp: expiration},
        {Type: "PUT", Strike: underlying - putSpread, Action: "SELL", Qty: 1, Exp: expiration},
        {Type: "CALL", Strike: underlying + callSpread, Action: "SELL", Qty: 1, Exp: expiration},
        {Type: "CALL", Strike: underlying + callSpread + 5, Action: "BUY", Qty: 1, Exp: expiration},
    }
}

// BuildVerticalSpread creates a 2-leg vertical spread
func BuildVerticalSpread(optionType, direction string, lowerStrike, upperStrike float64, expiration string) []OptionLeg

// BuildCalendarSpread creates a 2-leg calendar spread with different expirations
func BuildCalendarSpread(optionType string, strike float64, nearExp, farExp string) []OptionLeg

// BuildButterfly creates a 3-leg butterfly
func BuildButterfly(optionType string, lowerStrike, middleStrike, upperStrike float64, expiration string) []OptionLeg
```

### UI Strategy Selection Flow

1. **User selects strategy from dropdown** (e.g., "Iron Condor")
2. **System shows simplified inputs**:
   - Underlying price (auto-fetched)
   - Wing width (e.g., $5 from center)
   - Spread width (e.g., $10 for put side, $10 for call side)
   - Expiration date
3. **Builder generates legs automatically**
4. **User can review/edit each leg** before submitting
5. **System calculates max profit/loss/breakevens** from leg structure

### Example: Iron Condor Builder

```
=== Iron Condor Builder ===
Underlying: AAPL @ $175.00
Expiration: [2025-12-19] (50 DTE)

Put Spread Width: [$5] (e.g., $160/$165)
Call Spread Width: [$5] (e.g., $185/$190)
Wing Distance: [$10] from underlying

[Generate Legs]

Generated Structure:
  Buy Put  $160 @ [auto-fill price]
  Sell Put  $165 @ [auto-fill price]
  Sell Call $185 @ [auto-fill price]
  Buy Call  $190 @ [auto-fill price]

[Fetch Prices] [Edit Manually] [Continue]
```

## Implementation Phases (UPDATED)

### Phase 1: Database Foundation (3-4 hours)
- [ ] Create migration script `002_add_options_metadata.sql`
- [ ] Add new columns to `trade_sessions` table (including legs_json)
- [ ] Add new columns to `positions` table
- [ ] Create `trade_history` table
- [ ] Update `storage/sessions.go` with new fields + JSON marshaling
- [ ] Update `storage/positions.go` with new fields
- [ ] Create `OptionLeg` struct for JSON serialization
- [ ] Write unit tests for new storage functions

### Phase 2: Backend Enhancement (4-5 hours)
- [ ] Add strategy type constants (System1, System2, Custom)
- [ ] Add ALL 26 options strategy constants
- [ ] Add time exit mode constants (None, Close, Roll)
- [ ] Create `OptionLeg` struct with JSON tags
- [ ] Create strategy builder functions (Iron Condor, Spreads, Butterflies, etc.)
- [ ] Create max profit/loss calculator functions per strategy type
- [ ] Create breakeven calculator functions per strategy type
- [ ] Update `CreateSession()` to accept legs_json parameter
- [ ] Update `UpdateSessionSizing()` with pyramid fields
- [ ] Create `AddTradeToHistory()` function
- [ ] Create `GetCalendarView()` function for 10-week display
- [ ] Write integration tests for each strategy type

### Phase 3: UI Strategy Builders (5-6 hours)
- [ ] Create strategy builder dialog widgets (one per category)
- [ ] Build Iron Condor builder UI
- [ ] Build Butterfly builder UI
- [ ] Build Vertical Spread builder UI
- [ ] Build Calendar/Diagonal spread builder UI
- [ ] Build Straddle/Strangle builder UI
- [ ] Build Ratio/Broken Wing builder UI
- [ ] Add "Generate Legs" button functionality
- [ ] Add "Fetch Prices" integration (if API available)
- [ ] Add manual leg editing capability
- [ ] Add validation for leg structure (must be balanced, valid strikes, etc.)

### Phase 4: UI Dialog Integration (3-4 hours)
- [ ] Redesign "Start New Trade Session" dialog with all fields
- [ ] Add breakout system radio buttons (System-1, System-2, Custom)
- [ ] Add instrument type radio buttons (Stock/ETF vs Options)
- [ ] Add options strategy dropdown (categorized menu)
- [ ] Wire dropdown selection to appropriate builder dialog
- [ ] Add date pickers for entry and expiration
- [ ] Add DTE auto-calculation
- [ ] Add roll threshold input
- [ ] Add time exit mode dropdown
- [ ] Display calculated max profit/loss/breakevens
- [ ] Update dialog validation logic
- [ ] Wire up to backend CreateSession() with legs_json

### Phase 5: Position Sizing Enhancements (2 hours)
- [ ] Add "Pyramid Planning" section
- [ ] Add max units input
- [ ] Add "add step N" input
- [ ] Calculate and display add-on prices
- [ ] Show current units counter
- [ ] Update UpdateSessionSizing() calls with new data
- [ ] Add validation (units 1-10, add step > 0)

### Phase 6: Trade Entry Enhancements (2-3 hours)
- [ ] Add trade summary display with options details
- [ ] Show expiration date and DTE in summary
- [ ] Show pyramid levels in summary
- [ ] Update SAVE GO to create position + trade_history entry
- [ ] Test full workflow end-to-end

### Phase 7: Calendar View (4-5 hours)
- [ ] Create `buildCalendarScreen()` function
- [ ] Design 10-week grid layout (2 back + 8 forward)
- [ ] Query trade_history for date range
- [ ] Group trades by sector and week
- [ ] Display ticker symbols in cells (color-coded by status)
- [ ] Add click handler to show trade details
- [ ] Add filter dropdowns (sector, strategy, status)
- [ ] Add legend for status colors

### Phase 8: Dashboard Enhancements (2-3 hours)
- [ ] Update active trades widget with expiration info
- [ ] Show current units / max units
- [ ] Show next add-on price
- [ ] Show days until roll threshold
- [ ] Add visual alerts for trades nearing expiration
- [ ] Color-code by urgency (green >30 DTE, yellow 21-30, red <21)

### Phase 9: Testing & Documentation (3-4 hours)
- [ ] End-to-end test: Create System-2 Long Call session
- [ ] Test pyramid add-on calculations
- [ ] Test calendar view with multiple trades
- [ ] Test roll scenario (close at 21 DTE, reopen)
- [ ] Update USER_GUIDE.md with new workflows
- [ ] Update README.md with options trading features
- [ ] Create CHANGELOG entry
- [ ] Record demo video showing full workflow

## Total Estimated Time: 30-40 hours

**Breakdown:**
- Backend (Database + Logic): 7-9 hours
- Strategy Builders: 5-6 hours
- UI Integration: 5-7 hours
- Position Sizing + Entry: 4-6 hours
- Calendar View: 4-5 hours
- Dashboard: 2-3 hours
- Testing + Docs: 3-4 hours

## Success Criteria

1. ✅ Can create trade session with options strategy + dates
2. ✅ Can calculate pyramid add-on prices automatically
3. ✅ Can view all trades in 10-week calendar by sector
4. ✅ Can track DTE and get alerts at roll threshold (21 DTE)
5. ✅ Can distinguish System-1 vs System-2 vs Custom strategies
6. ✅ Dashboard shows expiration dates for all active trades
7. ✅ Full session history includes all options metadata
8. ✅ All existing tests pass
9. ✅ Documentation updated

## Migration Path

For existing sessions in database:
- `options_strategy` defaults to NULL (interpreted as stock/ETF)
- `entry_date` defaults to `created_at` date
- `max_units` defaults to 4
- `add_step_n` defaults to 0.5
- `time_exit_mode` defaults to 'Close'

No breaking changes to existing data.

## Anti-Impulsivity Compliance

**CRITICAL:** Despite the increased complexity of 26 options strategies, the system must maintain its core discipline enforcement:

### 1. All 5 Gates Still Apply

Regardless of options strategy selected:
- **Gate 1:** Checklist must be GREEN (all required items complete)
- **Gate 2:** 2-minute cooloff after checklist evaluation
- **Gate 3:** Ticker not on cooldown (no recent losses)
- **Gate 4:** Heat caps respected (portfolio + sector)
- **Gate 5:** Position sizing complete with ATR-based risk

**No exceptions. No backdoors. No overrides.**

### 2. Systematic Entry Only

- **Only System-1 or System-2 breakouts** trigger trades (or explicitly documented Custom)
- **No discretionary entries** - must follow mechanical breakout rules
- **No "feeling" or "hunch" entries** - signal must be present

### 3. Risk-Based Sizing Applies to Options

Options position sizing **must** use same ATR-based risk calculation:
- For directional (Long Call/Put): Use delta-adjusted shares equivalent
- For spreads: Risk = max loss × contracts
- For credit strategies: Risk = max loss, NOT just margin requirement
- **Same 0.75% risk per trade** regardless of strategy complexity

### 4. Heat Caps Cannot Be Gamed

Complex options strategies might tempt users to hide risk. Prevent this:
- **Iron Condors count as TWO positions** (put side + call side) for heat calculation
- **Butterflies count as full max loss**, not just debit paid
- **Calendar spreads** - both legs count for heat until near-dated expires
- **No netting** - if you have opposing positions, both count for heat

### 5. Journal Entry Required

**Higher complexity = More required documentation:**
- Why this specific options strategy vs stock/ETF?
- Why these strike prices?
- What is the thesis (volatility contraction? direction?)
- What is the exit plan at 21 DTE?
- How does this fit the trend-following system?

If user can't articulate clear answers → **Banner stays YELLOW or RED**

### 6. No Intraday Adjustments

Once trade is entered:
- **No leg adjustments** without creating new session
- **No "rolling early"** without going through full 5-gate checklist again
- **No "just this once"** overrides for exit strategy

### 7. Complexity Budget

To prevent analysis paralysis and over-optimization:
- **Maximum 3 active options strategies per week** across all positions
- If user has 3+ active iron condors, cannot start a butterfly
- Forces focus on high-conviction setups only
- Prevents "strategy of the day" behavior

### 8. Strategy Approval List

Not all 26 strategies should be available immediately:
- **Phase 1 Release:** Long Call, Long Put, Bull Call Spread, Bear Put Spread only
- **Phase 2 Release:** Add Iron Condor, Butterflies after 20+ successful Phase 1 trades
- **Phase 3 Release:** Ratio/Broken Wing after 50+ successful trades total
- User must "unlock" complex strategies through demonstrated discipline

### 9. Time Exit Mode Enforcement

If user selects "Roll" as time exit mode:
- **System forces new session creation** at 21 DTE - cannot just "roll in place"
- Must re-evaluate checklist (conditions may have changed)
- Must re-check heat caps (portfolio may be different)
- **Rolling is not automatic** - it's a new trade decision

### 10. No Naked Options (Ever)

Even though system supports "Short Call Butterfly", it does NOT support:
- Naked calls (undefined risk)
- Naked puts without cash backing
- Any structure with unlimited loss potential

**If max_loss = NULL or unlimited → REJECT at validation**

## Open Questions

1. **Calendar date range** - Should it be configurable? (Currently: 2 weeks back, 8 forward)
2. **Multiple positions per ticker** - Should we allow? (e.g., AAPL Dec calls + AAPL Feb calls)
3. **Roll history** - Should we track each roll as a separate record or update in place?
4. **Spreads** - Do we need to track P&L on each leg separately?
5. **Greeks tracking** - Do we want to track delta, theta, vega over time?
6. **Strategy unlock system** - Implement now or Phase 2?
7. **Complexity budget** - Enforce 3-strategy limit immediately?

## Dependencies

- Fyne date picker widget (need to add)
- SQLite date functions for calendar queries
- Possible Fyne grid widget for calendar layout

## References

- `docs/anti-impulsivity.md` - Core design principles
- `reference/Ed-Seykota.pine` - PineScript implementation
- `backend/internal/domain/sizing.go` - Van Tharp position sizing
- `backend/internal/storage/sessions.go` - Current session structure
