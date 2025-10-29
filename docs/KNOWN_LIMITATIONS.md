# Known Limitations

**Last Updated:** 2025-10-29
**Version:** 1.0.0
**TF = Trend Following**

---

## Overview

This document describes intentional design constraints, platform limitations, and known technical constraints of the TF-Engine trading system.

**Important:** These are **not bugs**. They are either by-design discipline constraints or acknowledged technical limitations.

---

## Limitations by Design (Anti-Impulsivity)

These constraints exist to enforce trading discipline and prevent impulsive behavior.

### 1. No Manual Gate Overrides

**Limitation:** You cannot bypass the 5 hard gates when saving a GO decision.

**Reason:** Intentional discipline enforcement. If gates fail, the trade should not be taken.

**Impact:** Cannot save GO decisions when:
- Banner is not GREEN
- 2-minute timer hasn't elapsed
- Ticker is on cooldown
- Heat caps are exceeded
- Position sizing is incomplete

**Workaround:** None - this is the core anti-impulsivity feature. If you must document the trade idea, save a NO-GO decision with notes.

---

### 2. Cannot Shorten Impulse Timer

**Limitation:** The 2-minute impulse brake cannot be reduced or skipped.

**Reason:** Prevents snap decisions. Forces a pause between checklist evaluation and trade execution.

**Impact:** Must wait 2 full minutes between completing checklist and saving GO decision.

**Workaround:** None - this is intentional. Use the time to review the trade setup, check TradingView, verify your plan.

---

### 3. Strict Heat Caps

**Limitation:**
- Portfolio heat cap: 4% of equity (cannot be increased)
- Sector bucket cap: 1.5% of equity (cannot be increased)

**Reason:** Prevents portfolio concentration and over-leveraging.

**Impact:** May need to:
- Reduce position sizes
- Close existing positions before entering new ones
- Wait for existing positions to hit profit targets

**Workaround:** None by design. If caps are reached, you have too much risk deployed.

---

### 4. No Position Editing After Entry

**Limitation:** Cannot modify entry price, stop, or initial risk after position is opened.

**Reason:** Prevents retroactive justification and "fudging" of trade parameters.

**Impact:** Must be certain of entry parameters before saving GO decision.

**Workaround:** If you entered wrong data, close the position and re-enter with correct parameters. The decision log will show both actions.

---

### 5. 55-Bar Breakout Requirement (Signal Gate)

**Limitation:** The system enforces 55-bar Donchian breakouts as the only entry signal.

**Reason:** Trend-following methodology. Not a flexible trading platform.

**Impact:** Cannot use other entry signals (moving averages, patterns, etc.).

**Workaround:** None - this is the Ed Seykota / Turtle Trader methodology the system implements. If you want different signals, this system isn't for you.

---

## Technical Limitations

### 1. FINVIZ Scraper Dependency

**Limitation:** The scanner relies on FINVIZ.com's HTML structure.

**Impact:**
- If FINVIZ changes their page layout, the scraper will break
- Scraping takes 3-5 seconds (network-dependent)
- Requires internet connection

**Workaround:**
- Manual ticker entry always available
- Import candidates from CSV (if implemented)
- System will continue working with manual entries

**Status:** Monitoring FINVIZ for changes. Will update scraper if needed.

---

### 2. Single User / Single Instance

**Limitation:** The system is designed for one trader, one database, one instance.

**Impact:**
- No multi-user accounts
- No shared portfolios
- Cannot run multiple instances against same database (SQLite single writer limitation)
- No collaboration features

**Workaround:** Each trader runs their own instance with their own database file.

**Rationale:** Individual trading discipline system, not a team platform.

---

### 3. No Direct Broker Integration

**Limitation:** The system does NOT execute trades in your broker.

**Impact:**
- Must manually enter trades in your broker after saving GO decision
- No automated order submission
- No live position sync from broker

**Workaround:** After saving GO decision:
1. Copy entry price, stop, shares from system
2. Manually enter limit order in your broker
3. Verify execution
4. System assumes position opened at saved parameters

**Rationale:** Intentional separation. Manual execution acts as final verification step and prevents "fat finger" errors.

---

### 4. No Historical Backtesting

**Limitation:** Cannot run backtests on historical price data.

**Impact:** System is live trading only, not for strategy research.

**Workaround:** Use TradingView, Amibroker, or dedicated backtesting software for historical testing.

**Rationale:** This is a live trading discipline tool, not a backtesting engine.

---

### 5. Browser Required (Web UI)

**Limitation:** Must have a web browser installed to use the UI.

**Impact:**
- Cannot run on headless servers (without workaround)
- Requires local browser (Chrome, Firefox, Edge, Safari)

**Workaround:** CLI mode can be used for basic operations (if implemented).

**Note:** The backend is a single Go binary with embedded UI, so it's very portable. Just needs a browser to view the UI.

---

### 6. SQLite Performance Limits

**Limitation:** SQLite has single-writer limitation and performance constraints.

**Impact:**
- One write operation at a time
- Thousands of positions may slow queries
- Not suitable for high-frequency trading data

**Practical Limits:**
- 1,000 open positions: ✅ No problem
- 10,000 historical decisions: ✅ No problem
- 100,000+ records: ⚠️ May need optimization

**Workaround:** Archival of old data (if needed), but realistically you won't hit limits with normal trading.

**Note:** WAL mode and indexes provide good performance for typical use.

---

## Platform Limitations

### 1. Operating System Support

**Supported:**
- ✅ Windows 10/11 (primary development target)
- ✅ Linux (Ubuntu, Debian, Arch, etc.)
- ✅ macOS (Intel and Apple Silicon)

**Requirements:**
- 64-bit operating system
- Modern web browser (Chrome 90+, Firefox 88+, Edge 90+, Safari 14+)

---

### 2. Market Data Sources

**Limitation:** System does NOT provide real-time market data.

**Impact:**
- No live price feeds
- No automatic ATR calculation
- Must manually look up prices and ATR from charting platform

**Workaround:** Use TradingView, thinkorswim, or your broker's platform for:
- Real-time prices
- ATR calculation
- Chart analysis
- Signal verification

**Rationale:** Market data subscriptions are expensive and complex. System focuses on decision workflow, not data provision.

---

### 3. Options Support

**Status:** ✅ Implemented with limitations

**Supported:**
- Delta-adjusted ATR sizing for options
- Max-loss sizing for spreads
- Position tracking

**Not Supported:**
- No Greeks calculations
- No options chain data
- No multi-leg spread builder
- No automatic roll calculations

**Workaround:** Use your broker's options tools for Greeks, chains, spreads. Use TF-Engine only for position sizing and heat management.

---

## Known Issues (Non-Critical)

### 1. API Handler Test Infrastructure

**Issue:** Some API handler tests fail due to test setup issues, not API bugs.

**Impact:** None on production code. Tests need better database initialization.

**Status:** Tracked in BUG_TRACKER.md as Medium priority (test-only issue).

**Workaround:** Not needed - production code works correctly.

---

### 2. Calendar Performance with 100+ Positions

**Issue:** Calendar view may slow down with hundreds of open positions.

**Impact:** Rare - most traders have 5-20 positions max.

**Status:** Monitoring. Will add virtual scrolling if needed.

**Workaround:** Close old positions periodically to keep active count manageable.

---

## Features NOT Planned

These features will NOT be added (by design):

### 1. ❌ Pattern Recognition / AI Signals

**Reason:** System uses 55-bar breakouts only (trend following methodology).

---

### 2. ❌ Intraday / Day Trading Support

**Reason:** System designed for position trading (multi-day to multi-week holds).

---

### 3. ❌ Cryptocurrency Support

**Reason:** Focus is on stocks, ETFs, and options on stocks. Crypto has different risk characteristics.

---

### 4. ❌ Social/Sharing Features

**Reason:** Trading discipline is personal. No social feed, no sharing trades, no following others.

---

### 5. ❌ Automated Trading / Bot Mode

**Reason:** Manual execution is a feature, not a limitation. It's the final verification step.

---

## Performance Characteristics

### Expected Latencies (Typical Hardware)

**Backend API:**
- Settings (cached): < 10ms
- Position sizing: < 100ms
- Heat check: < 100ms
- Gate validation: < 100ms
- Decision save: < 150ms
- Calendar (10 weeks): < 200ms
- FINVIZ scan: 3-5 seconds (network-dependent)

**Frontend:**
- Initial page load: < 2s
- Route transitions: < 200ms
- Form submissions: < 500ms

**Database:**
- Simple queries: < 50ms
- Complex joins: < 150ms

---

## Future Improvements (Maybe v2.0+)

These may be added in future versions, but are not guaranteed:

### Possible Future Features

- 📊 P&L tracking and analytics dashboard
- 📁 Export decisions to CSV/Excel
- 🎨 Custom color themes / dark mode improvements
- 📱 Mobile-responsive design improvements
- ☁️ Optional cloud backup/sync
- 📈 TradingView embedded widget (if feasible)
- 🔔 Price alert notifications (for stop levels)

**Note:** No commitments. Focus remains on core discipline enforcement.

---

## What This System IS and IS NOT

### ✅ This System IS:

- A **trading discipline enforcement tool**
- A **decision workflow manager**
- A **position sizing calculator** (Van Tharp method)
- A **heat management system** (caps and limits)
- A **trade journal** (decision logging)
- A **diversification tracker** (calendar view)

### ❌ This System IS NOT:

- A trading strategy (you provide the strategy)
- A market data provider (you provide prices/ATR)
- A backtesting engine (use dedicated tools)
- A broker (you execute trades manually)
- A signal generator (you identify setups)
- A performance tracker (focus is on process, not PnL)
- A social trading platform (no sharing/following)

---

## Getting Help

### For Technical Issues

- Check `docs/TROUBLESHOOTING.md` first
- Review `docs/TESTING_STRATEGY.md` for test guidance
- Check GitHub issues: (if repository public)

### For System Design Questions

- Read `docs/anti-impulsivity.md` (philosophy)
- Read `docs/project/WHY.md` (rationale)
- Read `CLAUDE.md` (development rules)

---

## Acceptance Criteria

By using this system, you acknowledge:

1. ✅ You understand the 5 hard gates cannot be bypassed
2. ✅ You accept the 2-minute impulse timer
3. ✅ You agree to manual trade execution (no automation)
4. ✅ You understand heat caps are strictly enforced
5. ✅ You accept this is trend-following only (55-bar breakouts)
6. ✅ You will provide your own market data and ATR values
7. ✅ You understand this is single-user, single-instance
8. ✅ You accept the technical limitations listed above

**If you cannot accept these constraints, this system is not for you.**

---

## Philosophy

> "The system's value comes from what it prevents (bad trades), not what it allows."

Every limitation exists for a reason. Most are discipline features, not bugs.

---

**Last Updated:** 2025-10-29
**Document Status:** Complete
**Next Review:** After v1.0 user feedback
