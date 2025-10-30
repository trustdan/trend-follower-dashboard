# TF-Engine Frequently Asked Questions (FAQ)

**Version:** 1.0.0
**Last Updated:** 2025-10-29

**TF = Trend Following** - Ed Seykota style systematic trading discipline

---

## Table of Contents

1. [General Questions](#general-questions)
2. [System Design & Philosophy](#system-design--philosophy)
3. [Installation & Setup](#installation--setup)
4. [The 5 Gates](#the-5-gates)
5. [Position Sizing & Risk Management](#position-sizing--risk-management)
6. [Entry & Exit Rules](#entry--exit-rules)
7. [Heat Management](#heat-management)
8. [Options Trading](#options-trading)
9. [Technical Issues](#technical-issues)
10. [Performance & Trading](#performance--trading)

---

## General Questions

### Q: What is TF-Engine?

**A:** TF-Engine (Trend-Following Engine) is a **discipline enforcement system** for systematic trend-following traders. It's NOT a signal generator, backtester, or auto-trader. It's a tool that makes impulsive trading impossible while making rule-based trading effortless.

The system implements Ed Seykota's System-2 approach: 55-bar Donchian breakouts, Van Tharp ATR-based sizing, pyramiding every 0.5×N, and 10-bar / 2×N exits.

---

### Q: Who is this for?

**A:** TF-Engine is designed for:
- **Systematic traders** who want to follow a proven trend-following system
- **Discretionary traders** struggling with discipline and impulsivity
- **Traders** who know the rules but have trouble following them consistently
- **Position traders** holding for days to months (NOT day traders)
- **Traders** committed to Ed Seykota style trend-following

**Not for:**
- Day traders (system designed for multi-day/week holds)
- Traders wanting flexibility and customization
- Traders looking for signal generation or backtesting
- Traders wanting auto-execution

---

### Q: How much does it cost?

**A:** [Specify pricing - free/open source, one-time purchase, subscription, etc.]

---

### Q: What does "TF" stand for?

**A:** **TF = Trend Following**. The system is purpose-built for trend-following strategies, not mean reversion, scalping, or other approaches.

---

## System Design & Philosophy

### Q: Can I bypass the 5 gates?

**A:** No. Absolutely not. That's the entire point.

The gates enforce discipline. If you want a flexible system where you can bypass rules "just this once," this is not the right tool for you.

**Remember:** Every bypass is an impulsive decision waiting to happen. The gates protect you from yourself.

---

### Q: Why is the system so rigid? Can't I have more flexibility?

**A:** The rigidity IS the feature, not a bug.

**Flexibility = Opportunity for impulsivity.**

Research shows that:
- Traders with discretion underperform systematic traders
- More rules = better compliance = better results
- Frictionless systems encourage overtrading

Every "inconvenience" in TF-Engine (2-min timer, 5 gates, heat caps, cooldowns) is an intentional design choice to prevent emotional, impulsive decisions.

If you want flexibility, this tool is not for you. If you want discipline, keep reading.

---

### Q: Why 2 minutes? Can I reduce the timer?

**A:** The 2-minute timer is backed by behavioral psychology research:
- Brief pauses reduce emotional decision-making
- 2 minutes is long enough to interrupt impulsive patterns
- Short enough to not disrupt workflow

**Cannot be reduced or disabled.** This is a core design principle.

**Use the 2 minutes productively:**
- Re-check TradingView chart
- Calculate position sizing
- Review sector diversification
- Take 3 deep breaths

If 2 minutes feels unbearable, that's a red flag. It means you're emotionally attached to the trade. That's exactly when you need the brake most.

---

### Q: Can I suggest features?

**A:** Yes! But understand the design philosophy:
- **Features that increase discipline:** Considered
- **Features that add flexibility:** Likely rejected
- **Features that add complexity:** Heavily scrutinized
- **Features that enable bypasses:** Never

Before suggesting, read:
1. `docs/anti-impulsivity.md` - Core design principles
2. `docs/project/WHY.md` - Why this system exists

Then ask: Does this feature support discipline or undermine it?

---

### Q: Why no auto-trading / broker integration?

**A:** Manual execution is intentional and crucial:

**Reasons:**
1. **Final human verification** before capital is committed
2. **Prevents catastrophic errors** (wrong symbol, wrong quantity, fat-finger mistakes)
3. **Keeps you engaged** in the process (not set-and-forget)
4. **You remain responsible** for your trades (not blaming software)
5. **Regulatory/legal simplicity** (no broker API, no accidental violations)

Auto-trading is seductive but dangerous. TF-Engine enforces the plan, you execute it.

---

### Q: What if I disagree with the sizing calculation?

**A:** The Van Tharp method is mathematically rigorous and time-tested by thousands of traders over decades.

If sizing seems wrong:
1. Verify your inputs (entry, ATR, equity, risk%)
2. Check math manually (see [USER_GUIDE.md](USER_GUIDE.md))
3. Understand the method (read Van Tharp's "Trade Your Way to Financial Freedom")

**The system does not allow manual overrides** because overrides lead to impulsivity ("I'll just add a few more shares...").

If you consistently want to override sizing, you either:
1. Have incorrect inputs, or
2. Are second-guessing the system (which means you need the discipline even more)

---

## Installation & Setup

### Q: What operating systems are supported?

**A:**
- **Windows 10/11 (64-bit):** Fully supported ✓
- **Linux:** Backend works, frontend may require manual setup
- **macOS:** Not officially tested but should work (Go is cross-platform)

Primary focus is Windows. Other OS users may need technical knowledge.

---

### Q: Do I need admin rights to install?

**A:** Depends on installation method:
- **Standalone binary:** No admin rights needed (run from any folder)
- **Installer (.msi or .exe):** May require admin for system-wide install
- **Recommended:** Use standalone binary in user folder (e.g., `Documents\TF-Engine\`)

---

### Q: Where is my data stored?

**A:** Database location:
```
C:\Users\[YourUsername]\AppData\Roaming\TF-Engine\trading.db
```

**Contains:**
- Settings (equity, risk%, caps)
- Positions (open positions, risk, stops)
- Decisions (GO/NO-GO history)
- Candidates (imported tickers)
- Cooldowns (sector/ticker cooldowns)
- Evaluations (checklist timestamps for 2-min timer)

**Backup regularly!** (See [USER_GUIDE.md](USER_GUIDE.md) for backup instructions)

---

### Q: Can I run multiple instances?

**A:** Not with the same database. SQLite locks the database file.

**Options:**
1. **One instance per computer:** Recommended
2. **Multiple instances with separate databases:** Possible but not recommended (why run multiple strategies in one account?)

**Attempting to run multiple instances on same database:**
- Second instance will fail with "database is locked" error
- Risk of data corruption
- Not supported

---

### Q: Can I sync across multiple computers?

**A:** Not natively supported.

**Workaround (advanced users only):**
1. Store `trading.db` in Dropbox/Google Drive/OneDrive
2. Create symlink from `%APPDATA%\TF-Engine\trading.db` to cloud folder
3. **CRITICAL:** Only run TF-Engine on ONE computer at a time
4. Risk of database corruption if both instances open simultaneously

**Not recommended for reliability.** Best practice: Use one primary trading computer.

---

## The 5 Gates

### Q: What are the 5 gates?

**A:** The 5 gates are the **final validation** before you can save a GO decision. All 5 must pass. They cannot be bypassed.

1. **Gate 1: Banner Status** → Checklist banner must be GREEN
2. **Gate 2: Impulse Brake** → 2-minute timer must elapse
3. **Gate 3: Cooldown Status** → Ticker/sector not on cooldown
4. **Gate 4: Heat Caps** → Portfolio & sector risk within limits
5. **Gate 5: Sizing Completed** → Position sizing calculated and saved

See [USER_GUIDE.md - The 5 Gates Explained](USER_GUIDE.md#the-5-gates-explained) for details.

---

### Q: Why do I need Gate 1 (banner) AND Gate 4 (heat)? Isn't that redundant?

**A:** No, they validate different things:

**Gate 1 (Banner):** Validates **checklist completion** at time of evaluation
- All required items checked
- Quality score ≥ threshold
- Ensures you've done your homework

**Gate 4 (Heat):** Validates **current portfolio state** at time of GO decision
- Portfolio heat ≤ cap (might have changed since checklist)
- Bucket heat ≤ cap (might have changed since checklist)
- Ensures real-time risk management

**Scenario where Gate 1 passes but Gate 4 fails:**
1. Complete checklist at 9:00 AM → Banner GREEN (Gate 1 pass)
2. During 2-min timer, add another position manually
3. Portfolio heat now at 95% of cap
4. Try to save GO decision at 9:05 AM → Gate 4 fails (heat cap exceeded)

Gates check different things at different times. Both necessary.

---

### Q: Can I skip the checklist and go straight to sizing?

**A:** No. The workflow is enforced:

1. Checklist → Save Evaluation (starts 2-min timer)
2. Position Sizing → Calculate & Save
3. Heat Check → Verify caps
4. Trade Entry → Run gates → Save GO decision

You cannot save GO decision without:
- Checklist evaluation saved (Gate 1)
- 2 minutes elapsed (Gate 2)
- Position sizing saved (Gate 5)

**Why?** Each step builds on the previous. Skipping steps = incomplete analysis = impulsive decision.

---

### Q: What triggers a cooldown?

**A:** Cooldowns are triggered by losing trades to prevent revenge trading.

**Ticker Cooldown:**
- Lose ≥1R on a trade → Ticker on cooldown for X days (configurable, default 7 days)
- Example: Lose -1.5R on AAPL → Cannot trade AAPL for 7 days

**Bucket Cooldown:**
- Multiple losses in same sector → Sector on cooldown for X days (configurable, default 14 days)
- Example: Lose 3 trades in Tech/Comm in 2 weeks → Cannot trade Tech/Comm for 14 days

**Purpose:** Forces you to take a break after losses. Prevents emotional revenge trading.

**Cannot be bypassed.** Take the break. Review what went wrong. Come back fresh.

---

## Position Sizing & Risk Management

### Q: What is "Van Tharp" position sizing?

**A:** Van Tharp's ATR-based position sizing method:

```
1. Define risk per unit: R = Equity × Risk% (e.g., $100,000 × 0.75% = $750)
2. Calculate stop distance: StopDist = K × ATR (e.g., 2 × $2.35 = $4.70)
3. Calculate shares: Shares = floor(R ÷ StopDist) = floor($750 ÷ $4.70) = 159
4. Verify: Actual risk = 159 × $4.70 = $747.30 ≤ $750 ✓
```

**Benefits:**
- Risk is constant across all trades (always $750 per unit, regardless of stock price)
- Stop is based on volatility (ATR), not arbitrary percentage
- Position size adjusts automatically (high volatility = fewer shares)

**Source:** Van Tharp's "Trade Your Way to Financial Freedom"

---

### Q: Why does the system round down shares?

**A:** To ensure actual risk never exceeds specified risk.

**Example:**
- Risk budget: $750
- Stop distance: $4.70
- Exact calculation: $750 ÷ $4.70 = 159.57 shares

**If we round up to 160 shares:**
- Actual risk: 160 × $4.70 = $752 (exceeds budget!)

**If we round down to 159 shares:**
- Actual risk: 159 × $4.70 = $747.30 (within budget ✓)

Rounding down is conservative. Ensures you never accidentally exceed risk budget.

---

### Q: Can I risk more than 1% per unit?

**A:** Technically yes (system allows 0.1% - 2.0%), but **not recommended**.

**Ed Seykota uses 1% per unit** with 4 units max = 4% per position.

**Why 0.75% - 1.0% is standard:**
- With 4 units, total risk per position is 3-4%
- With 4% portfolio cap, you can have ~5-6 positions
- Comfortable diversification and risk management

**If you use 2% per unit:**
- With 4 units, total risk per position is 8%
- Portfolio cap is 4% → You can only have 1 position at a time!
- Defeats the purpose of diversification

**Beginner recommendation:** Start with 0.50% per unit until comfortable.

---

### Q: What if I want to risk a fixed dollar amount (not % of equity)?

**A:** The system uses % of equity (not fixed dollars) because:

**% of equity is dynamic:**
- Account grows → Risk per trade grows proportionally
- Account shrinks → Risk per trade shrinks proportionally
- Prevents blowing up account on a losing streak

**Fixed dollar amount is static:**
- Account shrinks → Fixed risk becomes larger % of equity
- Riskier as you lose (opposite of what you want!)
- Example: $1000 risk on $100k account = 1%. Same $1000 on $50k account = 2%!

**If you strongly prefer fixed dollar:**
- Adjust equity in settings monthly to reflect actual account size
- Keep risk% constant (e.g., 0.75%)
- Result: Risk per trade stays approximately fixed

But % of equity is mathematically superior for long-term survival.

---

## Entry & Exit Rules

### Q: What is a 55-bar Donchian breakout?

**A:** A **55-bar Donchian breakout** occurs when:
- **Long entry:** Today's close > highest high of last 55 bars
- **Short entry:** Today's close < lowest low of last 55 bars

**Why 55 bars?**
- Captures intermediate-term trends (2-3 months on daily charts)
- Filters out short-term noise
- Proven by Turtle Traders and Ed Seykota

**How to verify:**
- Use Ed-Seykota.pine script on TradingView
- Blue line = 55-bar high (long trigger)
- Red line = 55-bar low (short trigger)
- Wait for close > blue line (for longs)

---

### Q: Can I enter on a limit order instead of market?

**A:** Yes, but be careful:

**Limit order pros:**
- Control your entry price
- Avoid slippage on low-volume stocks

**Limit order cons:**
- Might not fill (price runs away)
- Miss the trade entirely

**Recommendations:**
- Use limit at or slightly above breakout level (e.g., 55-high + $0.10)
- Set expiration (day order, not GTC)
- If doesn't fill by EOD, reconsider (momentum might be weak)

**Important:** TF-Engine calculates sizing based on your specified entry. If you use limit order at different price, recalculate sizing!

---

### Q: When exactly do I exit?

**A:** Two exit rules (**whichever is closer** to current price):

**Rule 1: 10-bar Donchian opposite breakout**
- For longs: Close < 10-bar low
- For shorts: Close > 10-bar high
- Check daily on TradingView chart
- Dynamic (trails as price moves in your favor)

**Rule 2: 2×N stop**
- For longs: Stop = Entry - 2×N (e.g., $180.50 - $4.70 = $175.80)
- For shorts: Stop = Entry + 2×N
- Static (set once in broker as stop-loss order)

**Exit immediately if either is hit.** No exceptions, no second-guessing.

**Example scenario (long AAPL @ $180.50, stop $175.80):**
- Day 1: 10-bar low = $178.00, 2×N stop = $175.80 → Stop is closer → Use $175.80
- Day 5: 10-bar low = $182.00, 2×N stop = $175.80 → 10-bar is closer → Exit if < $182
- Price hits $181.50 → Exit immediately (below 10-bar low of $182)

---

### Q: What if the 10-bar stop is way below my 2×N stop?

**A:** Use the **closer** of the two. If 10-bar is way below, that's actually good news!

**Scenario (long AAPL @ $180.50, stop $175.80):**
- Day 10: Price at $195, 10-bar low = $188
- 10-bar low ($188) is now closer than 2×N stop ($175.80)
- Your effective stop is $188 (locked in $7.50 profit!)

This is a **trailing stop**. As trade moves in your favor, 10-bar level moves up, protecting profits.

**Why this works:**
- Early in trade: 2×N stop protects from catastrophic loss
- Later in trade: 10-bar stop trails price, locking in gains
- Best of both worlds

---

### Q: How do I manage add-ons (pyramiding)?

**A:** Add to winners every 0.5×N (half an N):

**Initial position:**
- Unit 1: 159 shares @ $180.50 (entry)

**Add-ons (if price rises):**
- Unit 2: 159 shares @ $181.68 (entry + 0.5×N = $180.50 + $1.18)
- Unit 3: 159 shares @ $182.85 (entry + 1.0×N = $180.50 + $2.35)
- Unit 4: 159 shares @ $184.03 (entry + 1.5×N = $180.50 + $3.53)

**How to manage:**
1. Set price alerts at each add-on level ($181.68, $182.85, $184.03)
2. When alert triggers, verify position still open (haven't exited)
3. Execute add-on manually in broker: BUY 159 shares @ market
4. Update stop for all units to trailing 10-bar level

**Important:** Only add to winners (price moved 0.5×N in your favor). Never add to losers!

---

## Heat Management

### Q: What is "heat"?

**A:** **Heat** = Total risk across positions (sum of all $$ at risk if all stops hit)

**Portfolio heat:** Sum of risk across ALL positions
**Bucket heat:** Sum of risk within ONE sector

**Example ($100k account):**
- Position 1: AAPL (Tech) - risk $750
- Position 2: MSFT (Tech) - risk $820
- Position 3: XLE (Energy) - risk $690

**Portfolio heat:** $750 + $820 + $690 = $2,260 / $4,000 cap (56.5%)
**Tech bucket heat:** $750 + $820 = $1,570 / $1,500 cap (105% - EXCEEDS!)
**Energy bucket heat:** $690 / $1,500 cap (46%)

---

### Q: Why do we need BOTH portfolio cap AND bucket cap?

**A:** They serve different purposes:

**Portfolio heat cap (4%):** Limits **total risk** across all positions
- Prevents blowing up account on multiple simultaneous losses
- Example: With 4% cap on $100k, max loss if all stops hit = $4,000

**Bucket heat cap (1.5%):** Limits **sector concentration**
- Prevents correlated risk (all Tech stocks drop together)
- Forces diversification across sectors
- Example: Tech crash wipes out 3 Tech positions → Max loss $1,500 (not $4,000)

**Both necessary:** Portfolio cap limits total risk. Bucket cap limits correlation risk.

---

### Q: What if I want to concentrate in one sector (e.g., only trade Tech)?

**A:** The system won't let you. That's the point.

**Sector concentration = correlation risk.** When Tech crashes, all your positions crash together.

**Bucket cap (1.5%) forces diversification:**
- Max 1-2 positions per sector (depending on sizing)
- Spreads risk across Energy, Financials, Industrials, etc.
- Reduces portfolio volatility

**If you only want to trade one sector:**
1. This tool is not for you (system enforces diversification), OR
2. Adjust bucket cap in settings to match portfolio cap (e.g., both 4%)

But understand you're increasing correlation risk significantly.

---

### Q: Can I temporarily exceed heat caps "just this once"?

**A:** No. Gate 4 will fail. You cannot save GO decision.

**If heat caps exceeded:**
1. **Reduce position size:** Use "Calculate Max Shares" helper
2. **Close existing position:** Exit one position to free up capacity
3. **Choose different sector:** Find candidate in less crowded sector
4. **Wait:** Sometimes best decision is no trade

**Cannot be bypassed.** This is core risk management.

---

## Options Trading

### Q: Does TF-Engine support options?

**A:** Yes, with these methods:

1. **opt-delta-atr:** Delta-adjusted sizing for calls/puts
2. **opt-contracts:** Contract-based risk calculation

**Requirements:**
- Options must meet liquidity requirements:
  - Open Interest (OI) > 100 contracts
  - Bid-ask spread < 10% of mid price
- DTE (Days To Expiration): 60-90 days recommended
- Roll or close at ~21 DTE (don't hold to expiration)

---

### Q: How is options sizing different from stock sizing?

**A:** Options have unique risk characteristics:

**Stock sizing:**
- Risk = Shares × Stop distance
- Linear relationship (more shares = more risk)

**Options sizing (opt-delta-atr):**
- Risk = Contracts × 100 × (Premium × Delta × K × ATR)
- Accounts for leverage (delta)
- Accounts for volatility (ATR)
- More complex but mathematically sound

**Options sizing (opt-contracts):**
- Risk = Contracts × 100 × Contract price
- Simpler (total premium at risk)
- Conservative (assumes 100% loss possible)

See detailed examples in [USER_GUIDE.md - Options Section](USER_GUIDE.md#options-trading) (if exists).

---

### Q: Can I hold options through earnings?

**A:** Strongly discouraged.

**Risks:**
1. **Gap risk:** Price gaps past your stop (violates 2×N rule)
2. **Volatility crush:** IV drops post-earnings (hurts long options)
3. **Unpredictable:** Earnings surprises create large moves

**The optional "Earnings OK" checklist item** encourages you to avoid earnings.

**If you choose to hold anyway:**
- Understand you're accepting gap risk
- Might get stopped out with >2N loss
- System allows it but doesn't recommend it

**For stocks:** Less risky but still not ideal.
**For options:** Very risky. Avoid.

---

## Technical Issues

### Q: App won't start / browser doesn't open

**A:** See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for detailed solutions.

**Quick fixes:**
1. Check Task Manager → Is tf-engine.exe already running?
2. Manually open browser → Navigate to http://localhost:8080
3. Try different port: `tf-engine.exe server --listen :8081`
4. Run as Administrator

---

### Q: FINVIZ scan returns 0 candidates

**A:** Possible causes:

1. **FINVIZ URL wrong:** Test URL in browser first
2. **FINVIZ HTML changed:** Requires TF-Engine update
3. **Firewall blocking:** Add tf-engine.exe to whitelist
4. **Screener has 0 results:** Adjust FINVIZ filters

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for full diagnosis.

---

### Q: Database is locked

**A:** Causes:
1. **Multiple instances running:** Close all but one
2. **SQLite browser open:** Close database browser tools
3. **File permissions:** Check write access to `%APPDATA%\TF-Engine\`

**Solution:** End all tf-engine.exe processes, close database tools, restart once.

---

### Q: Can I use TF-Engine offline?

**A:** Partially.

**Offline functionality:**
- Checklist ✓
- Position sizing ✓
- Heat check ✓
- Gates check ✓
- Settings ✓

**Requires internet:**
- FINVIZ scanning ✗
- TradingView integration ✗ (opens web browser)

**Workflow for offline use:**
- Import candidates while online (morning)
- Analyze charts while online (TradingView)
- Complete checklist/sizing/gates offline (if needed)

---

### Q: How do I export trade history?

**A:** Currently no built-in export. Manual method:

1. Install DB Browser for SQLite: https://sqlitebrowser.org/
2. Open `%APPDATA%\TF-Engine\trading.db`
3. Navigate to "decisions" table
4. File → Export → Export to CSV
5. Open in Excel/Google Sheets

**Tables of interest:**
- `decisions`: GO/NO-GO decisions with timestamps
- `positions`: Open/closed positions
- `candidates`: Imported tickers
- `settings`: Account settings history

Future versions may include built-in export.

---

### Q: What data does TF-Engine collect?

**A:** **None.** Zero. Nada.

TF-Engine:
- Runs **entirely locally** on your computer
- **No telemetry** or analytics
- **No data** sent to remote servers
- **No tracking** of any kind
- FINVIZ requests go directly from your computer to FINVIZ (not through TF-Engine servers)

**Your trading data stays on your machine.**

Privacy is a core design principle.

---

## Performance & Trading

### Q: What win rate should I expect?

**A:** Trend-following systems typically have:
- **Win rate:** 30-45% (more losers than winners!)
- **Average winner:** 2-4R (multiples of risk)
- **Average loser:** -0.5 to -1R
- **Expectancy:** Positive (winners bigger than losers)

**Example:**
- 40% win rate
- Avg winner: +3R
- Avg loser: -1R
- Expectancy: (0.40 × 3R) + (0.60 × -1R) = 1.2R - 0.6R = **+0.6R per trade**

**Key insight:** You can be right less than half the time and still be profitable. Let winners run!

---

### Q: How many trades should I expect per month?

**A:** Depends on market conditions and number of sectors you trade:

**Active trends (2020, 2023):** 8-12 trades/month
**Choppy markets (2022):** 2-4 trades/month
**Average:** 4-8 trades/month

**Don't force trades.** Some months you'll have zero valid setups. Cash is a position.

---

### Q: How long should I hold positions?

**A:** As long as the trend lasts:

**Short-term trends:** 1-3 weeks (10-15 trading days)
**Medium-term trends:** 1-3 months
**Long-term trends:** 3-12+ months (rare but highly profitable)

**Average hold time:** 3-6 weeks

**The system lets the market decide:** Exit when 10-bar Donchian or 2×N stop is hit. Don't predict, react.

---

### Q: What if I keep getting stopped out?

**A:** Welcome to trend-following!

**Normal characteristics:**
- Many small losses (whipsaws)
- Occasional medium winners
- Rare but large winners (make up for all the losers)

**If excessive whipsaws (>60% loss rate):**
1. **Check entries:** Are you entering on true breakouts? (Verify on TradingView)
2. **Check exits:** Are you honoring stops? (Not widening them)
3. **Check markets:** Choppy/sideways markets generate more whipsaws (be patient)
4. **Check ATR:** Is N calculation correct? (Double-check on TradingView script)

**Losing streaks of 5-8 trades are normal.** Stay disciplined.

---

### Q: Should I trade both long and short?

**A:** Depends on your style and market regime:

**Long only (easier for most traders):**
- Simpler (most traders more comfortable)
- Market has upward bias long-term (SPY drift)
- Fewer false signals in bull markets

**Long and short (more opportunities):**
- More setups (can trade bear trends too)
- Better diversification
- Requires discipline to short (psychologically harder)

**TF-Engine supports both.** Start with longs, add shorts when comfortable.

---

### Q: What markets can I trade with TF-Engine?

**A:** Any liquid market with:
- ATR-based volatility (stocks, ETFs, futures)
- Clear trends (avoid mean-reverting markets)
- Sufficient liquidity (volume > 1M shares/day)

**Good candidates:**
- Large-cap stocks (S&P 500, Russell 1000)
- Sector ETFs (XLE, XLF, XLK, etc.)
- Commodity stocks (miners, energy)
- Index futures (if you trade futures)

**Poor candidates:**
- Penny stocks (low liquidity, high risk)
- Forex (mean-reverting, hard to trend-follow)
- Crypto (extreme volatility, gaps)
- Small-cap stocks (gaps, low volume)

---

### Q: Can I use this with a small account (<$10k)?

**A:** Technically yes, but challenging:

**With $10k account:**
- Risk per unit: 0.75% = $75
- Stop distance (typical): $3
- Shares: $75 / $3 = 25 shares
- Max stock price: ~$400 (for reasonable sizing)

**Challenges:**
- Limited universe (can't trade expensive stocks like TSLA, GOOG)
- Fewer positions (heat caps limit diversification)
- Commissions eat into returns

**Recommendations:**
- Start with $25k+ if possible
- Focus on lower-priced stocks ($50-$200 range)
- Consider sector ETFs (XLE, XLF) for diversification

**Account minimums:**
- Pattern Day Trader rule (US): $25k required for day trading
- TF-Engine is NOT day trading (hold days/weeks), so $10k+ is OK

---

## Still Have Questions?

**Resources:**
1. [USER_GUIDE.md](USER_GUIDE.md) - Comprehensive user guide
2. [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Technical issues
3. [QUICK_START.md](QUICK_START.md) - Get started in 10 minutes
4. **Contact Support:** [your-email or GitHub issues]

---

**Version:** 1.0.0
**Last Updated:** 2025-10-29
**Philosophy:** Trade the tide, not the splash.
