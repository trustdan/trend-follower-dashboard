# TF-Engine Troubleshooting Guide

**Version:** 1.0.0
**Last Updated:** 2025-10-29

**TF = Trend Following** - Systematic trading discipline enforcement

This guide covers common issues and their solutions. If you don't find your issue here, see [FAQ.md](FAQ.md) or contact support.

---

## Table of Contents

1. [Installation & Launch Issues](#installation--launch-issues)
2. [Browser & UI Issues](#browser--ui-issues)
3. [FINVIZ Scanning Issues](#finviz-scanning-issues)
4. [Checklist & Banner Issues](#checklist--banner-issues)
5. [Position Sizing Issues](#position-sizing-issues)
6. [Gates Check Issues](#gates-check-issues)
7. [Database Issues](#database-issues)
8. [TradingView Integration Issues](#tradingview-integration-issues)
9. [Performance Issues](#performance-issues)
10. [Data Loss & Recovery](#data-loss--recovery)

---

## Installation & Launch Issues

### Issue: App won't start when double-clicking tf-engine.exe

**Symptoms:**
- Double-click tf-engine.exe → Nothing happens
- No window opens
- No browser opens

**Diagnosis:**

**Step 1: Check if already running**
```
1. Open Task Manager (Ctrl+Shift+Esc)
2. Processes tab
3. Look for "tf-engine.exe"
```

**If running:**
- Right-click → End Task
- Try launching again

**If not running:**
- Proceed to Step 2

**Step 2: Try running from command line**
```cmd
1. Open Command Prompt (Win+R → type "cmd" → Enter)
2. Navigate to tf-engine.exe location:
   cd C:\Path\To\TF-Engine
3. Run:
   tf-engine.exe server
4. Look for error messages
```

**Common errors:**

| Error | Solution |
|-------|----------|
| "Port already in use" | Another app using port 8080. Try: `tf-engine.exe server --listen :8081` |
| "Permission denied" | Run as Administrator (right-click → Run as Administrator) |
| "Database is locked" | Close any SQLite browser tools. End all tf-engine.exe processes. |
| "Cannot find database" | Run: `tf-engine.exe init` to create database |

**Step 3: Check firewall/antivirus**
- Windows Firewall might be blocking tf-engine.exe
- Add tf-engine.exe to firewall whitelist:
  - Windows Security → Firewall & network protection
  - Allow an app through firewall
  - Add `tf-engine.exe`

**Step 4: Check system requirements**
- Windows 10/11 (64-bit) required
- 512 MB RAM minimum
- 100 MB disk space
- Try running as Administrator

---

### Issue: Browser doesn't open automatically

**Symptoms:**
- tf-engine.exe starts (command window visible)
- Browser doesn't open automatically
- Server is running (can see logs in command window)

**Solution:**

**Manually open browser:**
```
1. Leave tf-engine.exe running
2. Open your browser (Chrome, Firefox, Edge)
3. Navigate to: http://localhost:8080
```

**If page doesn't load:**
```
1. Check command window for actual port
   - Look for: "Server listening on :8080" or similar
2. Try: http://localhost:8080
3. If different port shown (e.g., :8081), use that
```

**Set default browser (Windows 10/11):**
```
Settings → Apps → Default apps → Web browser
Select: Chrome, Firefox, or Edge
```

**Prevent auto-open (if desired):**
```
tf-engine.exe server --no-browser
```

---

### Issue: Error "Database not initialized"

**Symptoms:**
- tf-engine.exe starts
- Error message: "Database not initialized" or "Cannot find database"
- UI shows "Database error"

**Solution:**

**Initialize database:**
```cmd
cd C:\Path\To\TF-Engine
tf-engine.exe init
```

**Expected output:**
```
Creating database at: C:\Users\[You]\AppData\Roaming\TF-Engine\trading.db
Database initialized successfully
Tables created: settings, positions, decisions, candidates, cooldowns, evaluations
```

**If error persists:**

**Check AppData folder permissions:**
```
1. Navigate to: %APPDATA%\TF-Engine
   - Win+R → type: %APPDATA%\TF-Engine → Enter
2. Check if folder exists
3. Check if you have write permissions:
   - Right-click folder → Properties → Security tab
   - Your user should have "Modify" and "Write" permissions
```

**If no permissions:**
- Run tf-engine.exe as Administrator once
- Or manually create folder with correct permissions

---

## Browser & UI Issues

### Issue: UI shows "Cannot connect to server"

**Symptoms:**
- Browser opens to http://localhost:8080
- Error: "Cannot connect" or "Connection refused"
- Page doesn't load

**Diagnosis:**

**Step 1: Is server running?**
```
Task Manager → Look for tf-engine.exe
```

**If not running:**
- Start tf-engine.exe
- Wait 2-3 seconds
- Refresh browser (F5)

**If running:**
- Proceed to Step 2

**Step 2: Is server listening on correct port?**
```
1. Check command window (tf-engine.exe)
2. Look for: "Server listening on :8080"
3. Match port in browser URL
```

**If different port:**
- Update browser URL to match
- Example: Server on :8081 → Navigate to http://localhost:8081

**Step 3: Try different localhost formats**
- http://localhost:8080
- http://127.0.0.1:8080
- http://[::1]:8080 (IPv6)

**Step 4: Check firewall**
- Windows Firewall might block localhost connections
- Temporarily disable firewall to test
- If works, add exception for tf-engine.exe

---

### Issue: UI elements not loading / blank screen

**Symptoms:**
- Browser connects to http://localhost:8080
- Blank white page
- Or partial UI (missing buttons, forms)

**Solutions:**

**Solution 1: Hard refresh**
```
Ctrl+Shift+R (or Ctrl+F5)
- Clears browser cache
- Reloads all assets
```

**Solution 2: Clear browser cache**
```
Chrome: Ctrl+Shift+Delete → Clear browsing data → Cached images and files
Firefox: Ctrl+Shift+Delete → Cache → Clear Now
Edge: Ctrl+Shift+Delete → Cached images and files → Clear
```

**Solution 3: Try different browser**
- Chrome (recommended)
- Firefox
- Edge

**Solution 4: Check browser console for errors**
```
1. Press F12 (opens DevTools)
2. Console tab
3. Look for red errors
4. Common errors:
   - "Failed to load resource" → Assets not found (hard refresh)
   - "CORS error" → Server config issue (restart tf-engine.exe)
   - "JavaScript error" → Browser compatibility (try Chrome)
```

**Solution 5: Rebuild/reinstall**
- If using embedded UI (tf-engine.exe includes UI):
  - Re-download tf-engine.exe
  - Replace old version
- If using separate UI files:
  - Check `ui/dist/` folder exists
  - Run build: `npm run build` (if you have source)

---

### Issue: UI is slow / laggy

**Symptoms:**
- Buttons take 2-3 seconds to respond
- Typing in forms is delayed
- Scrolling is choppy

**Solutions:**

**Solution 1: Check system resources**
```
Task Manager → Performance tab
- CPU usage: Should be <50% with tf-engine.exe running
- RAM usage: tf-engine.exe should use <200 MB
```

**If high CPU/RAM:**
- Close other apps
- Restart computer
- Check for background processes

**Solution 2: Check number of open positions**
- TF-Engine with 50+ open positions might be slow
- Archive old positions (if feature exists)
- Clean up database

**Solution 3: Browser extensions**
- Ad blockers sometimes interfere with local apps
- Try disabling extensions temporarily
- Or use Incognito/Private mode: Ctrl+Shift+N (Chrome/Edge)

**Solution 4: Database size**
- Large database (1000+ decisions) might slow queries
- Backup and compact database:
  ```
  1. Close tf-engine.exe
  2. Navigate to %APPDATA%\TF-Engine
  3. Copy trading.db (backup)
  4. Open in DB Browser for SQLite
  5. Database menu → Compact Database
  6. Save and close
  7. Restart tf-engine.exe
  ```

---

## FINVIZ Scanning Issues

### Issue: FINVIZ scan returns 0 candidates

**Symptoms:**
- Click "Run Daily FINVIZ Scan"
- Success message but 0 results
- Or error: "Failed to parse FINVIZ HTML"

**Diagnosis & Solutions:**

**Step 1: Test FINVIZ URL in browser**
```
1. Go to Settings → FINVIZ Presets
2. Copy FINVIZ URL
3. Paste in browser
4. Check results
```

**If 0 results in browser:**
- Your FINVIZ screener has no matches today
- Adjust screener filters (loosen criteria)
- Examples:
  - Price > 50% above 52-wk low (instead of 100%)
  - Volume > 500k (instead of 1M)
  - Remove some technical filters

**If results in browser but not in TF-Engine:**
- FINVIZ HTML structure changed (requires update)
- Proceed to Step 2

**Step 2: Check FINVIZ URL format**

**Correct format:**
```
https://finviz.com/screener.ashx?v=111&f=ta_price_a200sma,ta_rsi_os55
```

**Common errors:**
- Missing `v=111` (table view required)
- URL uses `v=152` or other view → Change to `v=111`
- URL too long (250+ characters) → Might get truncated

**Step 3: Test with default preset**

**Create simple test preset:**
```
Name: Test
URL: https://finviz.com/screener.ashx?v=111&f=cap_largeover,ta_averagevolume_o1000

(Large cap stocks with volume > 1M - should return 100+ results)
```

**If test works:**
- Your original URL/filters are the issue
- Build up filters incrementally

**If test doesn't work:**
- Network issue or FINVIZ blocked
- Proceed to Step 4

**Step 4: Check network/firewall**

**Test network connection:**
```cmd
ping finviz.com

Expected: Reply from 104.26.x.x (or similar)
Error: "Request timed out" → Network or firewall issue
```

**Check firewall:**
```
1. Windows Security → Firewall & network protection
2. Allow an app through firewall
3. Ensure tf-engine.exe has "Private" and "Public" checked
```

**Check antivirus:**
- Some antivirus blocks web scraping
- Temporarily disable antivirus to test
- If works, add exception for tf-engine.exe

**Step 5: Check for FINVIZ HTML changes**

**FINVIZ occasionally updates their HTML structure.**

**Workaround:**
- Manually copy candidates from FINVIZ
- Enter in TF-Engine manually (Checklist → Ticker input)

**Long-term fix:**
- Wait for TF-Engine update
- Or if open source: Update scraper code yourself

---

### Issue: FINVIZ scan fails with "Request timeout"

**Symptoms:**
- Click "Run Daily FINVIZ Scan"
- Loading... (hangs for 30+ seconds)
- Error: "Request timeout"

**Solutions:**

**Solution 1: Check internet connection**
```
1. Open browser
2. Navigate to https://finviz.com
3. If doesn't load: Internet connection issue
```

**Solution 2: Increase timeout (if configurable)**
- Check Settings → FINVIZ timeout
- Increase to 30 seconds (default might be 10s)

**Solution 3: Try again later**
- FINVIZ might be slow/overloaded
- Try different time of day
- Morning/evening EST (FINVIZ peak hours)

**Solution 4: Check VPN/proxy**
- If using VPN: Try disabling
- If using corporate proxy: Might block FINVIZ
- Try from home network

---

## Checklist & Banner Issues

### Issue: Banner stuck on RED

**Symptoms:**
- Checked all 5 required gates
- Banner still shows RED
- Message: "STOP - X required items missing"

**Solutions:**

**Solution 1: Verify ALL 5 required gates checked**
```
Scroll through entire checklist:
1. ☑ Signal: 55-bar Donchian breakout confirmed
2. ☑ Risk/Size: 2×N stop, pyramid every 0.5×N
3. ☑ Liquidity: Volume > 1M (or OI > 100)
4. ☑ Exits: 10-bar Donchian OR 2×N stop
5. ☑ Behavior: Will honor timer, no cooldown

ALL must have checkmarks (✓)
```

**Common mistake:**
- Scrolling past an unchecked item
- Checkbox state not saved (form bug)

**Solution 2: Fill ALL required trade data**
```
Ticker: [Must be filled]
Entry Price: [Must be > 0]
N (ATR): [Must be > 0]
Sector: [Must be selected]
Structure: [Must be selected]
```

**If any field empty → Banner stays RED**

**Solution 3: Refresh page**
```
F5 or Ctrl+R
- UI state might be stale
- Re-check all 5 gates after refresh
```

**Solution 4: Check browser console**
```
F12 → Console tab
Look for JavaScript errors:
- "State not updating" → Refresh page
- "Validation error" → Check required fields
```

---

### Issue: Banner stuck on YELLOW (want GREEN)

**Symptoms:**
- All 5 required gates checked
- Banner shows YELLOW: "CAUTION - Quality score below threshold"
- Want GREEN banner

**Explanation:**

YELLOW means:
- Required gates: ✓ Pass
- Quality score: ✗ Below threshold (e.g., 2.0 < 3.0)

**Solutions:**

**Solution 1: Check more quality items**
```
Optional quality items (check to increase score):
6. ☑ Regime OK: SPY > 200 SMA
7. ☑ No Chase: Entry within 2N of 20-EMA
8. ☑ Earnings OK: No earnings next 2 weeks
9. ☑ Journal Note: [Write 1-2 sentences]

Each checked adds to quality score.
```

**Example:**
- 2 quality items checked → Score: 2.0
- Need ≥ 3.0 for GREEN
- Check 1 more item → Score: 3.0 → GREEN ✓

**Solution 2: Lower quality threshold (not recommended)**
```
Settings → Quality Score Threshold
Change from 3.0 to 2.0
Save settings

Banner turns GREEN with only 2 quality items checked.
```

**Not recommended:** Quality items exist for a reason (better trade quality).

**Solution 3: Accept YELLOW and proceed**
- YELLOW banner allows you to proceed
- Trade is allowed (just lower quality)
- You can still save GO decision (if gates pass)

**Trade-off:**
- GREEN = High-quality setups → Better win rate
- YELLOW = Lower-quality → More trades but more losers

---

### Issue: 2-minute timer doesn't start

**Symptoms:**
- Click "Save Evaluation"
- No timer appears
- Expected countdown: 2:00 → 1:59 → ...

**Solutions:**

**Solution 1: Check banner state**
```
Timer only starts if banner is GREEN.

If RED or YELLOW:
- Timer doesn't start
- Complete required gates first
- Get banner to GREEN
- Then click "Save Evaluation"
```

**Solution 2: Look for timer in UI**
```
Timer location (might vary by UI version):
- Below "Save Evaluation" button
- In checklist summary section
- In separate "Timer" panel

Check entire screen for countdown display.
```

**Solution 3: Check evaluation saved**
```
1. Navigate away (go to Dashboard)
2. Return to Checklist
3. If evaluation saved: Data should persist
4. If not saved: Bug in save logic
```

**If evaluation not saving:**
- Check browser console (F12) for errors
- Refresh page and try again
- Check database (is evaluations table writable?)

**Solution 4: Timestamp-based timer (backend)**
```
Timer might be tracked on backend (not visible on frontend).

Test at Trade Entry:
1. Complete checklist
2. Wait 2 minutes manually (use phone timer)
3. Go to Trade Entry → Run Gate Check
4. Gate 2 should pass (2 min elapsed)

If passes: Timer works (just not visible on Checklist)
If fails: Timer not tracking (database issue)
```

---

## Position Sizing Issues

### Issue: Calculation results seem wrong

**Symptoms:**
- Position sizing calculation returns unexpected shares
- Example: Expected 100 shares, got 159
- Or: Risk amount doesn't match expectations

**Diagnosis:**

**Step 1: Verify inputs**
```
Check all inputs carefully:
- Entry: 180.50 (correct?)
- ATR (N): 2.35 (from TradingView? Double-check)
- K multiple: 2.0 (should be 2 for most cases)
- Equity: 100000 (from Settings → Check current equity)
- Risk %: 0.75% (from Settings)
```

**Step 2: Calculate manually**
```
Risk per unit:
R = Equity × Risk% = $100,000 × 0.0075 = $750

Stop distance:
StopDist = K × ATR = 2.0 × $2.35 = $4.70

Shares per unit:
Shares = floor(R ÷ StopDist) = floor($750 ÷ $4.70) = floor(159.57) = 159

Actual risk:
ActualRisk = Shares × StopDist = 159 × $4.70 = $747.30 ≤ $750 ✓
```

**If manual calculation matches TF-Engine:**
- Calculation is correct
- Your expectations might be off

**If manual calculation differs:**
- Bug in TF-Engine (report with details)

**Step 3: Check settings values**
```
Settings → Account Settings
- Equity: Should match current account size
- Risk %: Should be 0.50% - 1.00% (typical)

If equity is old (account grew/shrank):
- Update equity
- Recalculate sizing
```

---

### Issue: Concentration warning appears

**Symptoms:**
- Position sizing calculation succeeds
- Warning: "Position is X% of equity - consider reducing"
- Appears when position > 25% of equity

**Explanation:**

This is a **warning, not an error**. Trade is still allowed.

**Example:**
- Equity: $50,000
- Position: $15,000 (300 shares @ $50)
- Percentage: $15k / $50k = 30%
- 30% > 25% → Warning appears

**Causes:**
1. **High-priced stock:** Trading TSLA ($300), GOOG ($150), etc.
2. **Large position:** Max units (4) × shares × price
3. **Small account:** $25k - $50k equity

**Solutions:**

**Solution 1: Reduce max units**
```
Settings → Max units
Change from 4 to 3 or 2
Reduces total position size
```

**Solution 2: Accept the risk**
- 25% is a guideline, not a hard rule
- If comfortable with concentration, proceed
- Heat caps still protect overall portfolio

**Solution 3: Trade lower-priced stocks**
- Focus on stocks $50-$200 (instead of $300+)
- More shares, lower total position value

**Solution 4: Increase equity (if possible)**
- Add capital to account
- Update Settings → Equity
- Same shares, lower % concentration

---

## Gates Check Issues

### Issue: Gate 1 (Banner Status) fails but banner is GREEN

**Symptoms:**
- Trade Entry → Run Gate Check
- Gate 1: Banner Status → RED ✗
- But Checklist screen shows GREEN banner

**Diagnosis:**

**Step 1: Return to Checklist**
```
Navigate: Trade Entry → Checklist
Check current banner color:
- Is it GREEN on Checklist screen?
- Or did it revert to RED/YELLOW?
```

**If banner is RED/YELLOW on Checklist:**
- Gate check is correct
- Re-complete checklist
- Get banner back to GREEN

**If banner is still GREEN on Checklist:**
- Stale state in Trade Entry
- Proceed to Step 2

**Step 2: Refresh Trade Entry screen**
```
1. Navigate away (Dashboard)
2. Return to Trade Entry
3. Re-run Gate Check
```

**If still fails:**
- Database issue (evaluation not saved properly)
- Proceed to Step 3

**Step 3: Re-save checklist evaluation**
```
1. Go to Checklist
2. Verify banner GREEN
3. Click "Save Evaluation" again
4. Wait 2 minutes
5. Go to Trade Entry
6. Re-run Gate Check
```

---

### Issue: Gate 2 (Impulse Brake) fails - says not elapsed

**Symptoms:**
- Trade Entry → Run Gate Check
- Gate 2: Impulse Brake → RED ✗
- Message: "1:45 elapsed (need 2:00)"
- But you waited 2+ minutes

**Diagnosis:**

**Step 1: Check actual time elapsed**
```
Gate 2 result should show:
- Time evaluation saved: 09:23:45
- Current time: 09:26:00
- Elapsed: 2:15 ✓

If shows < 2:00:
- System clock is authoritative
- Check your computer clock (taskbar)
- Might be slow/fast
```

**Step 2: Check if evaluation was re-saved**
```
Did you:
1. Save evaluation at 09:20
2. Make changes to checklist
3. Save evaluation again at 09:24

If re-saved → Timer resets!
- Must wait 2 min from LAST save (09:24)
- Current time: 09:25 → Only 1:00 elapsed ✗
```

**Solution:**
- Don't re-save evaluation after initial save
- If you must make changes: Wait 2 min from last save

**Step 3: Database timestamp check**
```
(Advanced) Check database directly:
1. Close tf-engine.exe
2. Open trading.db in DB Browser for SQLite
3. Browse Data → evaluations table
4. Find latest row for your ticker
5. Check timestamp column
6. Compare to current time
```

---

### Issue: Gate 3 (Cooldown) fails - but don't see any cooldown

**Symptoms:**
- Trade Entry → Run Gate Check
- Gate 3: Cooldown Status → RED ✗
- Message: "Ticker AAPL on cooldown until 2025-11-05"
- Dashboard doesn't show any cooldowns

**Diagnosis:**

**Step 1: Check Dashboard → Cooldowns section**
```
Dashboard might not display cooldowns prominently.

Look for:
- "Active Cooldowns" section (might be collapsed)
- "Ticker Cooldowns" table
- "Sector Cooldowns" table

Find:
- Ticker: AAPL
- Expires: 2025-11-05
- Reason: Lost -1.2R on 2025-10-28
```

**Step 2: Check cooldown expiration date**
```
Today: 2025-11-01
Cooldown expires: 2025-11-05
Days remaining: 4 days

Cannot trade AAPL until expiration.
```

**Solution:**
- Wait for cooldown to expire (4 days)
- Or trade different ticker

**Step 3: Check sector cooldown**
```
Error might say:
"Sector Tech/Comm on cooldown until 2025-11-12"

If sector on cooldown:
- Cannot trade ANY ticker in Tech/Comm
- Must choose different sector:
  - Energy
  - Financials
  - Industrials
  - Materials
  - etc.
```

**Why cooldowns exist:**
- Prevent revenge trading after losses
- Force reflection and learning
- Core discipline feature (cannot be bypassed)

---

### Issue: Gate 4 (Heat Caps) fails - but Heat Check screen said OK

**Symptoms:**
- Heat Check screen: "TRADE APPROVED ✓"
- Trade Entry → Run Gate Check
- Gate 4: Heat Caps → RED ✗
- Message: "Portfolio heat exceeds cap by $150"

**Explanation:**

Heat can change between Heat Check and Trade Entry!

**Scenario:**
1. Heat Check at 09:20: Portfolio $3,800 / $4,000 (95%) → APPROVED ✓
2. Added another position at 09:22 (outside TF-Engine, or different evaluation)
3. Trade Entry at 09:25: Portfolio now $3,950 + $750 = $4,700 / $4,000 (117%) → REJECTED ✗

**Solution:**

**Re-run Heat Check:**
```
1. Go to Heat Check screen
2. Click "Check Heat for This Trade"
3. See updated result (might now exceed cap)
4. Adjust:
   - Reduce position size, or
   - Close existing position, or
   - Choose different ticker
```

**Best practice:**
- Complete entire workflow (Checklist → Sizing → Heat → Entry) quickly
- Don't add other positions mid-workflow
- Re-check heat if significant time elapsed

---

### Issue: Gate 5 (Sizing Completed) fails

**Symptoms:**
- Trade Entry → Run Gate Check
- Gate 5: Sizing Completed → RED ✗
- Message: "Position sizing not completed"
- But you calculated sizing on Position Sizing screen

**Diagnosis:**

**Step 1: Did you click "Save Position Plan"?**
```
Position Sizing screen:
1. Click "Calculate Position Size" → Shows results
2. Must also click "Save Position Plan" → Saves to database

If you only calculated (didn't save):
- Results not persisted
- Gate 5 fails
```

**Solution:**
- Go back to Position Sizing
- Click "Calculate Position Size"
- Click "Save Position Plan"
- Return to Trade Entry
- Re-run Gate Check

**Step 2: Check database (advanced)**
```
1. Close tf-engine.exe
2. Open trading.db in DB Browser
3. Browse Data → sizing_results table (or similar)
4. Look for row with your ticker and timestamp
5. If no row: Sizing not saved (bug or forgot to save)
```

---

## Database Issues

### Issue: Database is locked

**Symptoms:**
- Error: "Database is locked"
- Cannot save settings/positions/decisions
- Operations timeout

**Causes:**
1. Multiple tf-engine.exe instances running
2. SQLite browser tool open on trading.db
3. File permissions issue
4. Orphaned lock file

**Solutions:**

**Solution 1: Close all tf-engine.exe instances**
```
1. Task Manager (Ctrl+Shift+Esc)
2. Processes tab
3. Find all "tf-engine.exe"
4. Right-click each → End Task
5. Wait 5 seconds
6. Start tf-engine.exe ONCE
```

**Solution 2: Close database browser tools**
```
Close any of these if open:
- DB Browser for SQLite
- SQLite Studio
- DBeaver
- Any database tool viewing trading.db
```

**Solution 3: Check file permissions**
```
1. Navigate to: %APPDATA%\TF-Engine
2. Right-click trading.db → Properties
3. Security tab
4. Ensure your user has:
   - Read & execute ✓
   - Read ✓
   - Write ✓
```

**If missing permissions:**
- Click Edit → Add → Your username → Grant permissions
- Or: Run tf-engine.exe as Administrator once

**Solution 4: Delete lock file (if exists)**
```
1. Close tf-engine.exe
2. Navigate to: %APPDATA%\TF-Engine
3. Look for: trading.db-journal or trading.db-shm or trading.db-wal
4. Delete these files (safe to delete if tf-engine closed)
5. Restart tf-engine.exe
```

**Solution 5: Restart computer**
- Nuclear option
- Clears all file locks
- Try other solutions first

---

### Issue: Database corruption / cannot open database

**Symptoms:**
- Error: "Database disk image is malformed"
- Error: "Database file is corrupt"
- tf-engine.exe crashes on startup

**Diagnosis:**

**Check database file:**
```
1. Navigate to: %APPDATA%\TF-Engine
2. Check trading.db file:
   - Size: Should be > 0 bytes
   - If 0 bytes: File is empty (corrupted)
   - If missing: Never initialized
```

**Solutions:**

**Solution 1: Restore from backup**
```
If you have backup (recommended weekly backups):
1. Close tf-engine.exe
2. Navigate to: %APPDATA%\TF-Engine
3. Rename trading.db to trading-corrupted.db
4. Copy backup file: trading-backup-2025-10-28.db
5. Rename to: trading.db
6. Start tf-engine.exe
```

**Solution 2: Attempt repair (advanced)**
```
1. Install DB Browser for SQLite
2. Open trading.db
3. File → Export → Export to SQL
4. Save as: trading-export.sql
5. Close DB Browser
6. Delete trading.db
7. Run: tf-engine.exe init (creates new DB)
8. Import SQL:
   - Open new trading.db in DB Browser
   - File → Import → Import from SQL
   - Select trading-export.sql
   - Execute
9. Restart tf-engine.exe
```

**Solution 3: Start fresh (last resort)**
```
If no backup and can't repair:
1. Close tf-engine.exe
2. Rename trading.db to trading-old.db
3. Run: tf-engine.exe init
4. Reconfigure Settings
5. Manually re-enter any critical data
```

**Prevention:**
- Backup weekly: Copy trading.db to safe location
- Close tf-engine properly (don't kill process forcefully)
- Don't edit database while tf-engine is running

---

## TradingView Integration Issues

### Issue: "Open in TradingView" link doesn't work

**Symptoms:**
- Click "Open in TradingView" next to ticker
- Nothing happens
- Or: Opens TradingView but wrong chart

**Solutions:**

**Solution 1: Manual navigation**
```
1. Open TradingView manually
2. Search for ticker (e.g., AAPL)
3. Open chart
4. Add Ed-Seykota.pine script (if not already added)
```

**Solution 2: Check URL template in Settings**
```
Settings → TradingView URL Template

Should be:
https://tradingview.com/chart/?symbol={ticker}

If blank or wrong:
- Update to correct URL
- {ticker} is placeholder (replaced with actual ticker)
- Save settings
```

**Solution 3: Browser popup blocker**
```
TradingView opens in new tab/window.

Check popup blocker:
1. Look for icon in address bar (usually right side)
2. Click → Allow popups from localhost
3. Or: Browser settings → Popups → Allow localhost:8080
```

**Solution 4: TradingView account required**
```
Free TradingView account required:
1. Go to: https://www.tradingview.com
2. Sign up (free)
3. Login
4. Try link again from TF-Engine
```

---

### Issue: Ed-Seykota.pine script not showing on chart

**Symptoms:**
- Opened TradingView chart
- Script not visible
- No Donchian channels displayed

**Solutions:**

**Solution 1: Add script to chart**
```
1. Pine Editor (bottom panel)
2. Open saved script: "Ed-Seykota"
3. Click "Add to Chart" (top-right of Pine Editor)
4. Script appears on chart
```

**Solution 2: Re-install script**
```
If script not in Pine Editor:
1. Pine Editor → New
2. Copy contents from reference/Ed-Seykota.pine
3. Paste
4. Save → Name: "Ed-Seykota"
5. Click "Add to Chart"
```

**Solution 3: Check script visibility**
```
Script might be hidden:
1. Top of chart → Indicators
2. Look for "Ed-Seykota" in list
3. If eye icon has slash (hidden):
   - Click to show
```

**Solution 4: Script might have errors**
```
Pine Editor → Check for red error messages

Common errors:
- Syntax error (code changed)
- Re-copy from original reference/Ed-Seykota.pine
- Don't modify script
```

---

## Performance Issues

### Issue: UI is slow / laggy / unresponsive

**See [Browser & UI Issues → UI is slow / laggy](#issue-ui-is-slow--laggy) above.**

---

### Issue: FINVIZ scan takes 30+ seconds

**Symptoms:**
- Click "Run FINVIZ Scan"
- Loading... (hangs for 30+ seconds)
- Eventually returns results (or times out)

**Causes:**
1. Large screener (1000+ results)
2. FINVIZ server slow
3. Network latency
4. TF-Engine parsing 1000s of tickers

**Solutions:**

**Solution 1: Reduce screener results**
```
Adjust FINVIZ filters to return fewer candidates:
- Add price range: $50 - $500
- Increase volume filter: > 2M shares/day
- Add market cap filter: > $1B

Target: 100-300 candidates (not 1000+)
```

**Solution 2: Run scan during off-peak hours**
```
FINVIZ is faster:
- Early morning (6-8 AM EST)
- Late evening (8-10 PM EST)

Avoid:
- Market open (9:30-10 AM EST)
- Mid-day (12-2 PM EST)
```

**Solution 3: Use pagination (if supported)**
```
Some FINVIZ URLs paginate results.

Check URL for:
- &r=1 (results 1-20)
- &r=21 (results 21-40)

Scan multiple pages separately if needed.
```

**Solution 4: Accept the delay**
```
30 seconds once per day (morning routine) is acceptable.
- Run scan
- Get coffee
- Come back to results
```

---

## Data Loss & Recovery

### Issue: Lost data after computer crash

**Symptoms:**
- Computer crashed / power outage
- Restarted tf-engine.exe
- Some data missing (recent decisions, positions)

**Diagnosis:**

**Check database:**
```
1. Navigate to: %APPDATA%\TF-Engine
2. Check trading.db file:
   - Modified timestamp: Should be recent
   - Size: Should be > 0 bytes
```

**Recovery:**

**Solution 1: Check for auto-backups (if feature exists)**
```
%APPDATA%\TF-Engine\backups\
- Look for: trading-backup-DATE.db
- Restore most recent
```

**Solution 2: Restore from manual backup**
```
If you backed up weekly (recommended):
1. Locate backup: trading-backup-2025-10-28.db
2. Copy to: %APPDATA%\TF-Engine
3. Rename to: trading.db
4. Restart tf-engine.exe
```

**Solution 3: Recover from shadow copies (Windows)**
```
Windows keeps shadow copies (if enabled):
1. Navigate to: %APPDATA%\TF-Engine
2. Right-click trading.db
3. Properties → Previous Versions tab
4. Select recent version
5. Click Restore
```

**Solution 4: Accept data loss (if recent)**
```
If only lost last 1-2 hours of work:
- Re-enter recent decisions manually
- Check broker for executed trades
- Rebuild position list from broker account
```

**Prevention:**
- Backup weekly (copy trading.db to cloud storage)
- Close tf-engine.exe properly (don't force kill)
- Use UPS (battery backup) for computer

---

### Issue: Need to recover old evaluation/decision

**Symptoms:**
- Made decision yesterday
- Want to review: What was quality score? Journal note?
- Data seems gone from UI

**Solution:**

**Check Trade History (if feature exists):**
```
Dashboard → Trade History
Or: Decisions → History

Filter by:
- Ticker: AAPL
- Date: 2025-10-28

Should show:
- Decision: GO
- Quality score: 4/4
- Journal note: "Clean breakout..."
- Timestamp
```

**Query database directly (advanced):**
```
1. Close tf-engine.exe
2. Open trading.db in DB Browser for SQLite
3. Browse Data → decisions table
4. Find row:
   - ticker = "AAPL"
   - timestamp LIKE '2025-10-28%'
5. Read columns:
   - quality_score
   - journal_note
   - banner_state
   - etc.
```

---

## Still Having Issues?

**If your issue is not covered here:**

1. **Check other docs:**
   - [USER_GUIDE.md](USER_GUIDE.md) - Comprehensive guide
   - [FAQ.md](FAQ.md) - Frequently asked questions
   - [QUICK_START.md](QUICK_START.md) - Basic setup

2. **Gather information:**
   - Exact error message
   - Steps to reproduce
   - TF-Engine version
   - Windows version
   - Screenshots (if applicable)

3. **Contact support:**
   - Email: [your-support-email]
   - GitHub Issues: [your-repo/issues] (if open source)
   - Include gathered information above

4. **Community:**
   - Forum/Discord (if exists)
   - Other users may have solved same issue

---

**Version:** 1.0.0
**Last Updated:** 2025-10-29
**Remember:** Most issues are solved by: Restart, Refresh, Re-save, or Restore backup.
