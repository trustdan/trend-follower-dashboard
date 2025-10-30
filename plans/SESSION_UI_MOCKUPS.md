# Trade Session UI Mockups

**Purpose:** Visual reference for implementing trade session UI
**Status:** Design Phase
**Related:** [TRADE_SESSION_ARCHITECTURE.md](TRADE_SESSION_ARCHITECTURE.md)

---

## Screen 1: Dashboard with No Active Session

```
┌─────────────────────────────────────────────────────────────────────┐
│ TF-Engine - Trend Following Dashboard                        [-][□][×]│
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ╔═══════════════════════════════════════════════════════════════╗  │
│  ║  No Active Trade Session                                      ║  │
│  ║                                                                ║  │
│  ║  [📝 Start New Trade]  [📂 Resume Session ▼]  [📜 History]    ║  │
│  ╚═══════════════════════════════════════════════════════════════╝  │
│                                                                       │
│  🌙 Dark Mode  ❓ Help  VIM: On  👋 Welcome                          │
│                                                                       │
├───────────────────────────────────────────────────────────────────────┤
│ [Dashboard] [Checklist] [Position Sizing] [Heat Check] [Trade Entry] [Calendar] [Scanner]
├───────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  Dashboard - Portfolio Overview                                      │
│  ════════════════════════════                                        │
│                                                                       │
│  💼 Account Status                                                   │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ Equity:           $100,000                                  │    │
│  │ Risk per trade:   0.75% ($750)                              │    │
│  │ Portfolio heat:   $2,100 / $4,000 (52%)                     │    │
│  │ Open positions:   3                                         │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  📊 Open Positions                                                   │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ Ticker  Entry   Stop    Risk    Bucket       P&L           │    │
│  ├─────────────────────────────────────────────────────────────┤    │
│  │ AAPL    $180    $177    $750    Tech/Comm    +$450         │    │
│  │ MSFT    $380    $376    $800    Tech/Comm    -$120         │    │
│  │ XLE     $95     $93     $550    Energy       +$310         │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                       │
└───────────────────────────────────────────────────────────────────────┘
```

**Key UI Elements:**
- **Session bar** (top, gray when no session)
- **Three buttons**: Start New Trade, Resume Session (dropdown), History
- **Dashboard content** unchanged (existing functionality)

---

## Screen 2: "Start New Trade" Dialog

```
         ┌─────────────────────────────────────────────┐
         │ Start New Trade Session                     │
         ├─────────────────────────────────────────────┤
         │                                             │
         │  Select Trading Strategy:                   │
         │                                             │
         │  Strategy determines which FINVIZ preset    │
         │  is used for signal validation.             │
         │                                             │
         │  ┌───────────────────────────────────────┐  │
         │  │ ⚫ Long Breakout (55-bar high)        │  │
         │  │    Price > 55-high, mechanical long   │  │
         │  │                                       │  │
         │  │ ○ Short Breakout (55-bar low)        │  │
         │  │    Price < 55-low, mechanical short   │  │
         │  │                                       │  │
         │  │ ○ Custom Setup                        │  │
         │  │    Manual entry, no preset filter     │  │
         │  └───────────────────────────────────────┘  │
         │                                             │
         │  Ticker (optional):  [________]             │
         │                                             │
         │  ℹ️  Leave blank to fill in later           │
         │                                             │
         │  ┌─────────────────┐  ┌──────────────────┐ │
         │  │ Create Session  │  │     Cancel       │ │
         │  └─────────────────┘  └──────────────────┘ │
         │                                             │
         └─────────────────────────────────────────────┘
```

**User Flow:**
1. User clicks "Start New Trade" from dashboard
2. Dialog appears
3. User selects strategy (defaults to Long Breakout)
4. User optionally enters ticker
5. Clicks "Create Session"
   - System generates session #47
   - System navigates to Checklist tab
   - System shows session bar

---

## Screen 3: Active Session - Checklist Tab (GREEN)

```
┌─────────────────────────────────────────────────────────────────────┐
│ TF-Engine - Trend Following Dashboard                        [-][□][×]│
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ╔═══════════════════════════════════════════════════════════════╗  │
│  ║ Trade Session #47 • LONG_BREAKOUT • AAPL                     ║  │
│  ║ ✅ Checklist | ○ Sizing | ○ Heat | ○ Entry                    ║  │
│  ╚═══════════════════════════════════════════════════════════════╝  │
│                                                                       │
│  🌙 Dark Mode  ❓ Help  VIM: On  👋 Welcome                          │
│                                                                       │
├───────────────────────────────────────────────────────────────────────┤
│ [Dashboard] [Checklist] [Position Sizing] [Heat Check] [Trade Entry] [Calendar] [Scanner]
├───────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  Checklist - 5 Gates Evaluation                                      │
│  ════════════════════════════════                                    │
│                                                                       │
│  Session: #47 (Long Breakout)                                        │
│  Ticker: AAPL                                                        │
│                                                                       │
│  ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓  │
│  ┃                                                              ┃  │
│  ┃                        ✅ GREEN - GO                         ┃  │
│  ┃                  All Required Gates Pass                     ┃  │
│  ┃                                                              ┃  │
│  ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛  │
│                                                                       │
│  ✅ Required Gates (All Must Pass)                                   │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ ☑ From Preset (SIG_REQ)                              [ℹ️]   │    │
│  │   ✓ Ticker from FINVIZ Long Breakout preset                │    │
│  │                                                             │    │
│  │ ☑ Trend Confirmed (RISK_REQ)                         [ℹ️]   │    │
│  │   ✓ Price > 55-high, 2×N stop confirmed                    │    │
│  │                                                             │    │
│  │ ☑ Liquidity OK (OPT_REQ)                             [ℹ️]   │    │
│  │   ✓ Options bid-ask < 10%, OI > 100                        │    │
│  │                                                             │    │
│  │ ☑ TV Confirm (EXIT_REQ)                              [ℹ️]   │    │
│  │   ✓ Exit plan: 10-bar Donchian or 2×N stop                 │    │
│  │                                                             │    │
│  │ ☑ Earnings OK (BEHAV_REQ)                            [ℹ️]   │    │
│  │   ✓ 2-minute cooloff passed, no blackout               │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  📊 Optional Quality Items (Improve Score)                           │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ ☑ Regime Favorable                                   [ℹ️]   │    │
│  │ ☐ Not Chasing (entry near breakout)                 [ℹ️]   │    │
│  │ ☑ Journal Entry Logged                               [ℹ️]   │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  Quality Score: 8/10 (Excellent)                                     │
│                                                                       │
│  ┌──────────────┐  ┌────────────────────────────┐                   │
│  │  Evaluate    │  │  Next: Position Sizing →   │                   │
│  └──────────────┘  └────────────────────────────┘                   │
│                                                                       │
└───────────────────────────────────────────────────────────────────────┘
```

**Key Features:**
- Session bar shows: "#47 • LONG_BREAKOUT • AAPL"
- Progress: ✅ Checklist (completed) | ○ Sizing (pending)
- Large GREEN banner (impossible to miss)
- All checkboxes from session context
- "Next: Position Sizing" button enabled (after evaluation)

---

## Screen 4: Active Session - Checklist Tab (RED)

```
┌─────────────────────────────────────────────────────────────────────┐
│  ╔═══════════════════════════════════════════════════════════════╗  │
│  ║ Trade Session #47 • LONG_BREAKOUT • AAPL                     ║  │
│  ║ ⏳ Checklist | ○ Sizing | ○ Heat | ○ Entry                    ║  │
│  ╚═══════════════════════════════════════════════════════════════╝  │
│                                                                       │
│  Checklist - 5 Gates Evaluation                                      │
│  ════════════════════════════════                                    │
│                                                                       │
│  Session: #47 (Long Breakout)                                        │
│  Ticker: AAPL                                                        │
│                                                                       │
│  ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓  │
│  ┃                                                              ┃  │
│  ┃                    ❌ RED - DO NOT TRADE                     ┃  │
│  ┃                  2 Required Gates Failed                     ┃  │
│  ┃                                                              ┃  │
│  ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛  │
│                                                                       │
│  ✅ Required Gates (All Must Pass)                                   │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ ☐ From Preset (SIG_REQ)                              [ℹ️]   │    │
│  │   ❌ Ticker not in today's FINVIZ scan                      │    │
│  │                                                             │    │
│  │ ☐ Trend Confirmed (RISK_REQ)                         [ℹ️]   │    │
│  │   ❌ Price below 55-high, no breakout signal               │    │
│  │                                                             │    │
│  │ ☑ Liquidity OK (OPT_REQ)                             [ℹ️]   │    │
│  │ ☑ TV Confirm (EXIT_REQ)                              [ℹ️]   │    │
│  │ ☑ Earnings OK (BEHAV_REQ)                            [ℹ️]   │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  ⚠️ Missing Required Items:                                          │
│  • From Preset (SIG_REQ)                                             │
│  • Trend Confirmed (RISK_REQ)                                        │
│                                                                       │
│  📊 Optional Quality Items (Improve Score)                           │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ ☐ Regime Favorable                                   [ℹ️]   │    │
│  │ ☐ Not Chasing                                        [ℹ️]   │    │
│  │ ☐ Journal Entry Logged                               [ℹ️]   │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  Quality Score: 3/10 (Poor - Do Not Trade)                           │
│                                                                       │
│  ┌──────────────┐  ┌────────────────────────────┐                   │
│  │  Evaluate    │  │  Next: Position Sizing →   │  [DISABLED]       │
│  └──────────────┘  └────────────────────────────┘                   │
│                                                                       │
└───────────────────────────────────────────────────────────────────────┘
```

**Key Features:**
- RED banner (critical visual feedback)
- Progress: ⏳ Checklist (in progress, not completed)
- Clear list of missing items
- "Next: Position Sizing" button DISABLED (enforces discipline)
- User cannot proceed until gates pass

---

## Screen 5: Active Session - Position Sizing Tab

```
┌─────────────────────────────────────────────────────────────────────┐
│  ╔═══════════════════════════════════════════════════════════════╗  │
│  ║ Trade Session #47 • LONG_BREAKOUT • AAPL                     ║  │
│  ║ ✅ Checklist | ⏳ Sizing | ○ Heat | ○ Entry                    ║  │
│  ╚═══════════════════════════════════════════════════════════════╝  │
│                                                                       │
│  Position Sizing Calculator                                          │
│  ═════════════════════════                                           │
│                                                                       │
│  Session: #47 (Long Breakout)                                        │
│  Ticker: AAPL                    [from session, locked]              │
│  Banner: ✅ GREEN                 [from Checklist tab]               │
│                                                                       │
│  Account Settings:                                                   │
│  • Equity: $100,000                                                  │
│  • Risk per trade: 0.75% ($750)                                      │
│                                                                       │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ Sizing Method:  [Stock/ETF ▼]                               │    │
│  │                                                             │    │
│  │ ℹ️ Stock/ETF: Direct shares, ATR-based stop                 │    │
│  │                                                             │    │
│  │ Entry Price:      [$180.00___]                             │    │
│  │ ATR (N):          [$1.50______]                             │    │
│  │ K Multiple:       [$2.0_______]  (stop distance = 2×N)     │    │
│  │                                                             │    │
│  │ ┌──────────────────┐                                        │    │
│  │ │ Calculate Size   │                                        │    │
│  │ └──────────────────┘                                        │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  📊 Results:                                                         │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ Stop Distance:      $3.00  (2.0 × $1.50)                   │    │
│  │ Initial Stop:       $177.00                                 │    │
│  │ Shares:             25                                      │    │
│  │ Position Value:     $4,500  (25 × $180)                     │    │
│  │ Actual Risk:        $75.00  (25 × $3.00)                    │    │
│  │                                                             │    │
│  │ ✅ Risk within limit: $75 ≤ $750                            │    │
│  │                                                             │    │
│  │ 📈 Max Pyramids:    4 units (0.5×N spacing)                 │    │
│  │                                                             │    │
│  │ Next Add:           $180.75  (entry + 0.5×N)                │    │
│  │ Max Total Risk:     $300  (4 units × $75)                   │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  ┌────────────────────────┐                                          │
│  │  Next: Heat Check →    │                                          │
│  └────────────────────────┘                                          │
│                                                                       │
└───────────────────────────────────────────────────────────────────────┘
```

**Key Features:**
- Ticker auto-filled from session (locked)
- Banner status shown (from Checklist)
- Progress: ✅ Checklist | ⏳ Sizing (in progress)
- Calculation results saved to session
- "Next: Heat Check" button enabled after calculation

---

## Screen 6: Active Session - Heat Check Tab

```
┌─────────────────────────────────────────────────────────────────────┐
│  ╔═══════════════════════════════════════════════════════════════╗  │
│  ║ Trade Session #47 • LONG_BREAKOUT • AAPL                     ║  │
│  ║ ✅ Checklist | ✅ Sizing | ⏳ Heat | ○ Entry                    ║  │
│  ╚═══════════════════════════════════════════════════════════════╝  │
│                                                                       │
│  Heat Check - Portfolio Risk Management                              │
│  ═════════════════════════════════════                               │
│                                                                       │
│  Session: #47 (Long Breakout)                                        │
│  Ticker: AAPL                                                        │
│  Risk: $75                       [from Position Sizing]              │
│  Bucket: Tech/Comm                                                   │
│                                                                       │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ 🌡️ Portfolio Heat Check                                     │    │
│  │                                                             │    │
│  │ Current Portfolio Heat:  $2,100                             │    │
│  │ New Position Risk:       +$75                               │    │
│  │ ────────────────────────────────                            │    │
│  │ New Total Heat:          $2,175                             │    │
│  │                                                             │    │
│  │ Portfolio Cap:           $4,000  (4% of $100k)              │    │
│  │                                                             │    │
│  │ ✅ Within Cap: $2,175 < $4,000                              │    │
│  │ Remaining Capacity: $1,825  (45%)                           │    │
│  │                                                             │    │
│  │ ┌─────────────────────────────────────────────────┐         │    │
│  │ │████████████████████████████░░░░░░░░░░░░░░░░░░│ 54%       │    │
│  │ └─────────────────────────────────────────────────┘         │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ 🗂️ Bucket Heat Check (Tech/Comm)                            │    │
│  │                                                             │    │
│  │ Current Bucket Heat:     $1,400                             │    │
│  │ New Position Risk:       +$75                               │    │
│  │ ────────────────────────────────                            │    │
│  │ New Total Heat:          $1,475                             │    │
│  │                                                             │    │
│  │ Bucket Cap:              $1,500  (1.5% of $100k)            │    │
│  │                                                             │    │
│  │ ⚠️  Near Cap: $1,475 < $1,500 (98% utilized)                │    │
│  │ Remaining Capacity: $25  (2%)                               │    │
│  │                                                             │    │
│  │ ┌─────────────────────────────────────────────────┐         │    │
│  │ │███████████████████████████████████████████████│ 98%       │    │
│  │ └─────────────────────────────────────────────────┘         │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  💡 Recommendation: Consider diversifying to other sectors           │
│                                                                       │
│  ┌────────────────────────┐                                          │
│  │  Next: Trade Entry →   │                                          │
│  └────────────────────────┘                                          │
│                                                                       │
└───────────────────────────────────────────────────────────────────────┘
```

**Key Features:**
- Risk auto-loaded from Position Sizing ($75)
- Progress: ✅ Checklist | ✅ Sizing | ⏳ Heat (in progress)
- Visual progress bars for portfolio and bucket heat
- Clear capacity indicators
- Warning when near cap (yellow)
- "Next: Trade Entry" enabled after check

---

## Screen 7: Active Session - Trade Entry Tab (Final Gate)

```
┌─────────────────────────────────────────────────────────────────────┐
│  ╔═══════════════════════════════════════════════════════════════╗  │
│  ║ Trade Session #47 • LONG_BREAKOUT • AAPL                     ║  │
│  ║ ✅ Checklist | ✅ Sizing | ✅ Heat | ⏳ Entry                    ║  │
│  ╚═══════════════════════════════════════════════════════════════╝  │
│                                                                       │
│  Trade Entry - Final 5-Gate Check                                    │
│  ═══════════════════════════════                                     │
│                                                                       │
│  Session: #47 (Long Breakout - AAPL)                                 │
│                                                                       │
│  📋 Session Summary                                                  │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ Strategy:      Long Breakout (55-bar high)                  │    │
│  │ Ticker:        AAPL                                         │    │
│  │ Entry:         $180.00                                      │    │
│  │ Stop:          $177.00  (2×N = $3.00)                       │    │
│  │ Shares:        25                                           │    │
│  │ Risk:          $75                                          │    │
│  │ Bucket:        Tech/Comm                                    │    │
│  │ Banner:        ✅ GREEN                                      │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  🚦 5-Gate Status Check                                              │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ Gate 1: Banner GREEN                          ✅ PASS       │    │
│  │         All required checklist items complete               │    │
│  │                                                             │    │
│  │ Gate 2: 2-Minute Cooloff Elapsed              ✅ PASS       │    │
│  │         Evaluated 3 min 42 sec ago (> 2 min)                │    │
│  │                                                             │    │
│  │ Gate 3: Ticker Not on Cooldown                ✅ PASS       │    │
│  │         No recent losses in AAPL                            │    │
│  │                                                             │    │
│  │ Gate 4: Heat Caps Not Exceeded                ✅ PASS       │    │
│  │         Portfolio: $2,175 / $4,000 (54%)                    │    │
│  │         Bucket:    $1,475 / $1,500 (98%)                    │    │
│  │                                                             │    │
│  │ Gate 5: Position Sizing Completed             ✅ PASS       │    │
│  │         25 shares, $75 risk, $4,500 total                   │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  ✅ All 5 Gates PASS - Ready to Trade                                │
│                                                                       │
│  ⚠️  FINAL DECISION - This action will be logged                     │
│                                                                       │
│  ┌──────────────────┐  ┌──────────────────┐                         │
│  │   Save GO ✅     │  │  Save NO-GO ❌   │                         │
│  └──────────────────┘  └──────────────────┘                         │
│                                                                       │
│  💡 Saving GO will:                                                  │
│  • Log trade decision to database                                    │
│  • Complete Session #47 (read-only)                                  │
│  • Start 2-minute impulse timer                                      │
│  • Add AAPL to your active watchlist                                 │
│                                                                       │
└───────────────────────────────────────────────────────────────────────┘
```

**Key Features:**
- Progress: All gates ✅ (Checklist, Sizing, Heat) | ⏳ Entry (in progress)
- Full session summary (all data from prior tabs)
- Clear 5-gate status (all pass/fail)
- Two big buttons: Save GO / Save NO-GO
- Warning that decision will be logged
- Clear explanation of what happens next

---

## Screen 8: Session Completed (Read-Only View)

```
┌─────────────────────────────────────────────────────────────────────┐
│  ╔═══════════════════════════════════════════════════════════════╗  │
│  ║ Trade Session #47 • LONG_BREAKOUT • AAPL • [COMPLETED]       ║  │
│  ║ ✅ Checklist | ✅ Sizing | ✅ Heat | ✅ Entry (GO)              ║  │
│  ╚═══════════════════════════════════════════════════════════════╝  │
│                                                                       │
│  Trade Entry - Session Complete                                      │
│  ══════════════════════════════                                      │
│                                                                       │
│  🔒 Session #47 is READ-ONLY (Decision logged 2025-10-30 14:23)     │
│                                                                       │
│  📋 Session Summary                                                  │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ Strategy:      Long Breakout (55-bar high)                  │    │
│  │ Ticker:        AAPL                                         │    │
│  │ Entry:         $180.00                                      │    │
│  │ Stop:          $177.00  (2×N = $3.00)                       │    │
│  │ Shares:        25                                           │    │
│  │ Risk:          $75                                          │    │
│  │ Bucket:        Tech/Comm                                    │    │
│  │ Banner:        ✅ GREEN                                      │    │
│  │ Decision:      ✅ GO                                         │    │
│  │ Decided:       2025-10-30 14:23:17                          │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  🚦 Final 5-Gate Check (at time of decision)                         │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ Gate 1: Banner GREEN                          ✅ PASS       │    │
│  │ Gate 2: 2-Minute Cooloff Elapsed              ✅ PASS       │    │
│  │ Gate 3: Ticker Not on Cooldown                ✅ PASS       │    │
│  │ Gate 4: Heat Caps Not Exceeded                ✅ PASS       │    │
│  │ Gate 5: Position Sizing Completed             ✅ PASS       │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  ✅ Decision: GO (all gates passed)                                  │
│                                                                       │
│  ┌──────────────────┐  ┌──────────────────────┐                     │
│  │  Clone Session   │  │  View Decision Log   │                     │
│  └──────────────────┘  └──────────────────────┘                     │
│                                                                       │
│  💡 To analyze this setup again, clone to new session                │
│                                                                       │
└───────────────────────────────────────────────────────────────────────┘
```

**Key Features:**
- Session bar shows [COMPLETED] status
- All gates ✅ including Entry
- Read-only view (no editing)
- "Clone Session" button (creates new draft with same ticker/strategy)
- "View Decision Log" button (shows full audit trail)
- Immutable record

---

## Screen 9: "Resume Session" Dropdown

```
┌─────────────────────────────────────────────────────────────────────┐
│  ╔═══════════════════════════════════════════════════════════════╗  │
│  ║  No Active Trade Session                                      ║  │
│  ║                                                                ║  │
│  ║  [📝 Start New Trade]  [📂 Resume Session ▼]  [📜 History]    ║  │
│  ╚═══════════════════════════════════════════════════════════════╝  │
│                                       ▼                               │
│                    ┌─────────────────────────────────────────┐       │
│                    │ Resume Session                          │       │
│                    ├─────────────────────────────────────────┤       │
│                    │ Active Draft Sessions:                  │       │
│                    │                                         │       │
│                    │ #32  TSLA • Short Breakout              │       │
│                    │      ✅ Checklist | ⏳ Sizing            │       │
│                    │      Updated 2 hours ago                │       │
│                    │                                         │       │
│                    │ #18  NVDA • Long Breakout               │       │
│                    │      ✅ Checklist | ✅ Sizing | ○ Heat   │       │
│                    │      Updated yesterday                  │       │
│                    │                                         │       │
│                    │ #7   MSFT • Custom                      │       │
│                    │      ○ Checklist                        │       │
│                    │      Updated 3 days ago                 │       │
│                    │                                         │       │
│                    ├─────────────────────────────────────────┤       │
│                    │ [View All Sessions (12 total)...]      │       │
│                    └─────────────────────────────────────────┘       │
│                                                                       │
```

**User Flow:**
1. User clicks "Resume Session" dropdown
2. Shows list of DRAFT sessions (incomplete)
3. Sorted by most recently updated
4. Shows progress for each session
5. User clicks session → loads into AppState.currentSession
6. UI navigates to last incomplete tab

---

## Screen 10: Session History View

```
┌─────────────────────────────────────────────────────────────────────┐
│ TF-Engine - Trade Session History                            [-][□][×]│
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  Trade Session History                                               │
│  ════════════════════                                                │
│                                                                       │
│  Filter: [All ▼]  [COMPLETED ▼]  [DRAFT ▼]                          │
│  Search: [________________________________________]  [🔍]             │
│                                                                       │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ Session  Ticker  Strategy        Status      Decision  Date │    │
│  ├─────────────────────────────────────────────────────────────┤    │
│  │ #47      AAPL    Long Breakout   COMPLETED   ✅ GO    10/30 │    │
│  │ #32      TSLA    Short Breakout  DRAFT       -        10/30 │    │
│  │ #28      NVDA    Long Breakout   COMPLETED   ❌ NO-GO 10/29 │    │
│  │ #18      MSFT    Custom          DRAFT       -        10/28 │    │
│  │ #12      XLE     Long Breakout   COMPLETED   ✅ GO    10/27 │    │
│  │ #7       SPY     Short Breakout  ABANDONED   -        10/26 │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  Select session to view details, clone, or export                    │
│                                                                       │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                           │
│  │   View   │  │  Clone   │  │  Export  │                           │
│  └──────────┘  └──────────┘  └──────────┘                           │
│                                                                       │
└───────────────────────────────────────────────────────────────────────┘
```

**Key Features:**
- Filterable by status (DRAFT, COMPLETED, ABANDONED)
- Searchable by ticker
- Shows all historical sessions
- "View" → opens read-only session
- "Clone" → creates new draft with same parameters
- "Export" → CSV/JSON export for analysis

---

## Color Scheme

### Banner Colors (High Contrast)

```
GREEN (Go):
  Background: #00C853  (Material Green A700)
  Text:       #FFFFFF  (White)

YELLOW (Caution):
  Background: #FFD600  (Material Yellow A700)
  Text:       #000000  (Black)

RED (No-Go):
  Background: #D50000  (Material Red A700)
  Text:       #FFFFFF  (White)

GRAY (Not Evaluated):
  Background: #9E9E9E  (Material Gray 500)
  Text:       #FFFFFF  (White)
```

### Session Bar Colors

```
Active Session (light mode):
  Background: #E3F2FD  (Light Blue 50)
  Border:     #1976D2  (Blue 700)
  Text:       #0D47A1  (Blue 900)

Active Session (dark mode):
  Background: #263238  (Blue Gray 900)
  Border:     #42A5F5  (Blue 400)
  Text:       #BBDEFB  (Blue 100)

No Session (both modes):
  Background: #EEEEEE  (Gray 200) / #424242 (Gray 800)
  Border:     #9E9E9E  (Gray 500)
  Text:       #616161  (Gray 700) / #E0E0E0 (Gray 300)
```

### Progress Icons

```
✅  Completed (green checkmark)
⏳  In Progress (hourglass)
○   Pending (hollow circle)
❌  Failed (red X)
```

---

## Keyboard Shortcuts

```
Ctrl+N          Start New Trade
Ctrl+R          Resume Session (opens dropdown)
Ctrl+H          View Session History
Ctrl+S          Save current step (auto-evaluate)

Tab Navigation (while in session):
Ctrl+1          Checklist tab
Ctrl+2          Position Sizing tab
Ctrl+3          Heat Check tab
Ctrl+4          Trade Entry tab
Ctrl+5          Calendar tab

VIM Mode (when enabled):
g               Navigate to Checklist (first gate)
s               Navigate to Sizing
h               Navigate to Heat
e               Navigate to Entry
```

---

## Responsive Behavior

### Minimum Window Size
- Width: 1000px
- Height: 700px

### Session Bar Collapsing
```
Full Width (1200px+):
  ║ Trade Session #47 • LONG_BREAKOUT • AAPL                     ║
  ║ ✅ Checklist | ✅ Sizing | ⏳ Heat | ○ Entry                    ║

Medium Width (900-1199px):
  ║ Session #47 • LONG • AAPL                                    ║
  ║ ✅ Check | ✅ Size | ⏳ Heat | ○ Entry                         ║

Narrow (< 900px):
  ║ #47 • AAPL                                                   ║
  ║ ✅ ✅ ⏳ ○                                                      ║
```

---

## Accessibility

### Screen Reader Support
- Session bar announces: "Trade Session 47, Long Breakout strategy, ticker AAPL, Checklist completed, Sizing completed, Heat in progress, Entry pending"
- Banner announces: "Banner status GREEN, all required gates pass, proceed to trade"

### High Contrast Mode
- Banner colors meet WCAG AAA contrast ratio (7:1+)
- Focus indicators clearly visible
- Icons have text labels

### Keyboard Navigation
- All buttons tabbable
- Enter/Space to activate
- Esc to close dialogs
- Arrow keys in dropdowns

---

## Animation/Transitions

### Session Creation
```
1. User clicks "Start New Trade"
2. Dialog slides up (200ms ease-out)
3. User selects strategy, clicks Create
4. Dialog fades out (150ms)
5. Session bar slides down from top (300ms ease-out)
6. Navigate to Checklist tab (fade 200ms)
```

### Tab Switching (with session)
```
1. User clicks "Position Sizing" tab
2. Check for unsaved changes
3. If changes: show warning dialog (instant)
4. If no changes or user confirms:
   - Current tab fades out (100ms)
   - New tab fades in (100ms)
   - Session bar progress updates (instant)
```

### Gate Completion
```
1. User clicks "Evaluate" on Checklist
2. Button shows spinner (instant)
3. Backend evaluates (< 100ms)
4. Banner color changes (200ms transition)
5. Progress icon updates: ○ → ⏳ → ✅ (instant)
6. "Next" button enables (fade in 150ms)
```

---

## Error States

### Session Creation Failed
```
         ┌─────────────────────────────────────────────┐
         │ Error Creating Session                      │
         ├─────────────────────────────────────────────┤
         │                                             │
         │  ❌ Failed to create trade session          │
         │                                             │
         │  All 99 session slots are in use.           │
         │                                             │
         │  Please complete or abandon existing        │
         │  sessions before starting a new one.        │
         │                                             │
         │  [View Active Sessions]  [OK]               │
         │                                             │
         └─────────────────────────────────────────────┘
```

### Cannot Proceed (RED Banner)
```
         ┌─────────────────────────────────────────────┐
         │ Cannot Proceed                              │
         ├─────────────────────────────────────────────┤
         │                                             │
         │  ❌ Session #47 has RED banner               │
         │                                             │
         │  Position Sizing is disabled until all      │
         │  required checklist items pass.             │
         │                                             │
         │  Missing items:                             │
         │  • From Preset (SIG_REQ)                    │
         │  • Trend Confirmed (RISK_REQ)               │
         │                                             │
         │  [Go to Checklist]  [Cancel]                │
         │                                             │
         └─────────────────────────────────────────────┘
```

### Unsaved Changes Warning
```
         ┌─────────────────────────────────────────────┐
         │ Unsaved Changes                             │
         ├─────────────────────────────────────────────┤
         │                                             │
         │  ⚠️  Session #47 has unsaved changes         │
         │                                             │
         │  You have filled in some checklist items    │
         │  but not clicked "Evaluate".                │
         │                                             │
         │  What would you like to do?                 │
         │                                             │
         │  [Save & Continue]  [Discard]  [Cancel]     │
         │                                             │
         └─────────────────────────────────────────────┘
```

---

## Implementation Notes

1. **Session bar should be a reusable component** (`ui/session_bar.go`)
   - Updates when AppState.currentSession changes
   - Can be placed at top of any screen

2. **Progress tracking should be automatic**
   - Backend: `UpdateSessionChecklist()` sets `checklist_completed=1`
   - UI: Observes AppState.currentSession, updates icons

3. **Tab navigation should be gated**
   - Each tab checks prerequisites
   - If prerequisite fails, show error dialog and navigate to blocker

4. **Session persistence is critical**
   - Every user action updates database
   - No in-memory-only state
   - If app crashes, session recoverable

5. **Visual hierarchy matters**
   - Session bar: MOST PROMINENT (top, colored border)
   - Banner: VERY PROMINENT (large, bright colors)
   - Tab content: NORMAL (standard Fyne widgets)
   - Navigation buttons: SUBTLE (bottom, secondary importance)

---

**Document Version:** 1.0
**Author:** Claude Code Planning Agent
**Next:** Implement Phase 1 (Database)
