# Trading Platform UI - Quick Reference Guide

**Version:** 1.0 (M22)
**Last Updated:** 2025-10-28

---

## Overview

The Trading Platform workbook contains 7 worksheets providing a complete trading workflow from analysis to decision execution.

---

## Worksheets

### 1. Setup (Configuration)
**Purpose:** System configuration
**Contents:**
- Engine path setting
- Database path setting

**Usage:** Generally no changes needed after initial setup

---

### 2. VBA Tests (Testing)
**Purpose:** Automated testing of VBA integration
**Contents:**
- "Run All Tests" button

**Usage:**
1. Click "Run All Tests" button
2. Verify all tests pass (green checkmarks)
3. Check for any failures (red X)

---

### 3. Dashboard (Overview)
**Purpose:** Portfolio overview and navigation hub
**Contents:**
- Portfolio status (equity, heat, cap, %)
- Today's candidates section
- Quick action navigation buttons

**Buttons:**
- **Refresh Dashboard** - Update portfolio data
- **Position Sizing** - Go to Position Sizing sheet
- **Checklist** - Go to Checklist sheet
- **Heat Check** - Go to Heat Check sheet
- **Trade Entry** - Go to Trade Entry sheet

**Workflow:**
- Start here to see portfolio status
- Use navigation buttons to access other features

---

### 4. Position Sizing (Calculate Trade Size)
**Purpose:** Calculate shares/contracts for a trade entry
**Tab Color:** Green

#### Required Inputs
- **Ticker:** Stock symbol (e.g., AAPL)
- **Entry Price:** Planned entry price ($)
- **ATR (N):** Average True Range in dollars
- **K Multiple:** Stop distance multiplier (typically 2-3)
- **Method:** Dropdown selection:
  - `stock` - Regular stock position
  - `opt-delta-atr` - Options using delta and ATR
  - `opt-maxloss` - Options using max loss

#### Optional Inputs
- **Equity Override:** Custom equity amount (overrides system setting)
- **Risk % Override:** Custom risk percentage (overrides system setting)
- **Delta:** Options delta (required for opt-delta-atr method)
- **Max Loss:** Maximum loss per contract (required for opt-maxloss method)

#### Results Displayed
- **Risk Dollars (R):** Dollar amount at risk
- **Stop Distance:** Distance to stop in $
- **Initial Stop:** Stop price
- **Shares:** Number of shares to buy
- **Contracts:** Number of options contracts to buy
- **Actual Risk:** Actual risk after rounding
- **Status:** Success/error message

#### Buttons
- **Calculate** - Run position sizing calculation
- **Clear** - Clear all inputs and results

#### Example Workflow
1. Enter ticker: `AAPL`
2. Enter entry price: `180.00`
3. Enter ATR: `1.50`
4. Enter K Multiple: `2`
5. Select method: `stock`
6. Click **Calculate**
7. Review results (shares, stop, risk)
8. Use results in Trade Entry sheet

---

### 5. Checklist (Entry Validation)
**Purpose:** Evaluate 6-item pre-entry checklist
**Tab Color:** Orange

#### Input
- **Ticker:** Stock symbol

#### Checklist Items (Check all that apply)
1. **Ticker from today's FINVIZ preset** - Stock is in today's scan results
2. **Trend alignment confirmed** - Meets trend criteria (e.g., above 20 EMA)
3. **Adequate volume and spread** - Liquid enough to trade
4. **TradingView setup confirmation** - Chart setup verified
5. **No earnings in next 7 days** - Earnings date clear
6. **Trade thesis documented in journal** - Rationale written down

#### Results Displayed
- **Banner:** Color-coded result
  - **GREEN** - All 6 items checked (ready to trade)
  - **YELLOW** - 1-2 items missing (proceed with caution)
  - **RED** - 3+ items missing (do not trade)
- **Missing Items:** Count of unchecked items
- **Missing:** List of specific items missing
- **Allow Save:** Whether decision can be saved (TRUE/FALSE)
- **Evaluation Time:** Timestamp of evaluation
- **Status:** Success/error message

#### Buttons
- **Evaluate** - Run checklist evaluation
- **Clear** - Clear ticker and uncheck all boxes

#### Example Workflow
1. Enter ticker: `AAPL`
2. Check all applicable items (goal: all 6)
3. Click **Evaluate**
4. Review banner color:
   - GREEN → Proceed to Heat Check
   - YELLOW → Review missing items, decide if acceptable
   - RED → Do not trade, too many missing items

---

### 6. Heat Check (Risk Management)
**Purpose:** Verify portfolio and bucket heat caps
**Tab Color:** Red

#### Required Inputs
- **Ticker:** Stock symbol
- **Risk Amount ($):** Dollar amount at risk (from Position Sizing)
- **Bucket:** Sector dropdown:
  - Tech/Comm
  - Finance
  - Healthcare
  - Consumer
  - Energy
  - Industrials

#### Portfolio Heat Results
- **Current Heat:** Current portfolio heat ($)
- **New Heat:** Portfolio heat after adding this trade ($)
- **Heat %:** New heat as % of cap
- **Cap:** Portfolio heat cap ($)
- **Exceeded:** YES/NO (red if exceeded)
- **Overage:** Amount over cap if exceeded ($)

#### Bucket Heat Results
- **Current Heat:** Current bucket heat ($)
- **New Heat:** Bucket heat after adding this trade ($)
- **Heat %:** New heat as % of bucket cap
- **Cap:** Bucket heat cap ($)
- **Exceeded:** YES/NO (red if exceeded)
- **Overage:** Amount over bucket cap if exceeded ($)

#### Buttons
- **Check Heat** - Run heat cap validation
- **Clear** - Clear all inputs and results

#### Example Workflow
1. Enter ticker: `AAPL`
2. Enter risk amount: `500.00` (from Position Sizing)
3. Select bucket: `Tech/Comm`
4. Click **Check Heat**
5. Review results:
   - If both Portfolio and Bucket show "NO" for Exceeded → Proceed to Trade Entry
   - If either shows "YES" → Review overage, decide if acceptable or resize position

---

### 7. Trade Entry (5 Hard Gates)
**Purpose:** Final trade decision with 5-gate validation
**Tab Color:** Purple

#### Required Inputs
- **Ticker:** Stock symbol
- **Entry Price:** Planned entry price ($)
- **ATR:** Average True Range
- **Method:** Dropdown (stock, opt-delta-atr, opt-maxloss)
- **Bucket:** Sector dropdown (same as Heat Check)

#### Optional Inputs
- **Banner Status:** Dropdown (GREEN, YELLOW, RED) - from Checklist
- **Delta:** Options delta (for options trades)
- **Max Loss:** Max loss per contract (for options trades)
- **Preset:** FINVIZ preset name if applicable

#### Gate Status (Displayed after Save)
1. **Gate 1 - Banner GREEN:** Checklist must be GREEN
2. **Gate 2 - In Candidates:** Ticker must be in today's preset
3. **Gate 3 - Impulse Brake:** No impulsive trading (cooldown check)
4. **Gate 4 - Cooldown:** Not trading same ticker too frequently
5. **Gate 5 - Heat Caps:** Portfolio and bucket heat within limits

Each gate shows: ✅ PASS or ❌ FAIL

#### Results Displayed
- **Decision Saved:** YES/NO
- **Decision ID:** Database ID if saved
- **Rejection Reason:** Why trade was rejected (if any)
- **Status:** Success/error message

#### Buttons
- **Save GO** - Save GO decision (attempt to enter trade)
- **Save NO-GO** - Save NO-GO decision (decline trade)
- **Clear** - Clear all inputs and results

#### Example Workflow - GO Decision
1. Enter ticker: `AAPL`
2. Enter entry price: `180.00` (from Position Sizing)
3. Enter ATR: `1.50`
4. Select method: `stock`
5. Select banner status: `GREEN` (from Checklist)
6. Select bucket: `Tech/Comm` (from Heat Check)
7. (Optional) Enter preset: `tech_breakout`
8. Click **Save GO**
9. Review gate results:
   - All gates PASS → Trade accepted, Decision ID assigned
   - Any gate FAIL → Trade rejected, reason displayed
10. Form clears automatically on successful save

#### Example Workflow - NO-GO Decision
1. Enter ticker: `AAPL`
2. Select bucket: `Tech/Comm`
3. (Optional) Enter preset: `tech_breakout`
4. Click **Save NO-GO**
5. Enter rejection reason in prompt (e.g., "Spread too wide")
6. Decision saved with reason
7. Form clears automatically

---

## Typical Trading Workflow

### Full Workflow (New Trade Idea)
1. **Dashboard** - Check current portfolio status
2. **Position Sizing** - Calculate shares/contracts
   - Enter ticker, entry price, ATR, K, method
   - Note Risk Dollars (R) for Heat Check
3. **Checklist** - Evaluate 6-item checklist
   - Goal: Achieve GREEN banner
4. **Heat Check** - Verify heat caps
   - Use Risk Dollars from Position Sizing
   - Ensure no caps exceeded
5. **Trade Entry** - Save final decision
   - Enter all trade details
   - Click "Save GO" if proceeding
   - Click "Save NO-GO" if declining

### Quick Workflow (Pre-analyzed Trade)
If you've already done analysis outside Excel:
1. **Trade Entry** - Go directly to Trade Entry
2. Enter all required fields
3. Click "Save GO"
4. Gates will validate automatically

### Rejection Workflow
If trade fails any validation:
1. Review rejection reason
2. Decide if you can fix the issue (e.g., reduce position size for heat)
3. If fixable, adjust inputs and try again
4. If not fixable, save as NO-GO with reason

---

## Tips and Best Practices

### Position Sizing
- Use K=2 for conservative stops, K=3 for wider stops
- Options delta should be between 0 and 1 (e.g., 0.70 for 70 delta)
- Review Actual Risk vs. Risk Dollars to see rounding impact

### Checklist
- Aim for GREEN (all 6 items)
- YELLOW is acceptable if you understand the risk
- Never trade on RED (3+ missing items)
- Document your trade thesis before checking journal-ok

### Heat Check
- Always run before entering trade
- If exceeded, consider reducing position size
- Remember bucket caps are per-sector diversification limits

### Trade Entry
- Banner status must be GREEN to pass Gate 1
- Preset field is optional but recommended for tracking
- Save NO-GO decisions to track trades you declined (valuable data)
- Correlation IDs in status help with debugging

### Navigation
- Use Dashboard navigation buttons to move between sheets
- Each sheet has Clear button to reset form
- Status messages show correlation IDs for log tracing

### Error Handling
- Red status messages indicate errors
- Correlation IDs link to TradingSystem_Debug.log
- Check log file for detailed error information

---

## Keyboard Shortcuts

**Excel Standard:**
- `Ctrl+PgDn` - Next worksheet
- `Ctrl+PgUp` - Previous worksheet
- `Ctrl+Home` - Go to cell A1
- `Alt+F11` - Open VBA editor

**Custom:** (Future enhancement)
- TBD: Custom shortcuts for common operations

---

## Troubleshooting

### Button Doesn't Respond
- Ensure macros are enabled
- Check if Excel is calculating (wait for completion)
- Look for error messages in status cells

### "❌ Missing required inputs"
- Review all required fields (marked with red * would be future enhancement)
- Ensure numeric fields have valid numbers
- Check dropdown selections are made

### "❌ Error: ..." Messages
- Note the Correlation ID in the message
- Check TradingSystem_Debug.log for details
- Verify tf-engine.exe is running
- Ensure trading.db exists

### Results Don't Update
- Click the Calculate/Evaluate button again
- Check if previous operation completed
- Look for error in status cell

### Checkboxes Don't Appear
- Excel version may not support ActiveX controls
- Try running setup again
- Check Excel security settings

---

## Advanced Features

### Correlation IDs
- Every operation generates a unique correlation ID
- Format: `YYYYMMDD-HHMMSSmmm-XXXX`
- Links Excel operation to Go engine log entries
- Appears in status messages and error messages
- Use to trace issues across VBA and Go logs

### Logging
- All operations logged to `TradingSystem_Debug.log`
- Log includes correlation IDs, timestamps, and details
- Useful for troubleshooting and audit trail

### Custom Settings
- Equity and Risk % can be overridden in Position Sizing
- System defaults from Setup sheet
- Overrides apply only to current calculation

---

## Future Enhancements (Roadmap)

### Planned Features
- Dashboard portfolio queries (real data)
- Candidates display on Dashboard
- Chart visualizations
- Historical trade review
- Performance metrics
- Export to CSV/PDF
- Email notifications
- Real-time updates

### Requested Features
- Copy trade details between sheets
- Bulk position sizing
- What-if scenario analysis
- Heat visualization charts
- Trade journal integration

---

## Support and Feedback

### Getting Help
1. Check status messages for correlation IDs
2. Review TradingSystem_Debug.log
3. Check WINDOWS_TESTING.md for test procedures
4. Review M22_COMPLETION_SUMMARY.md for implementation details

### Reporting Issues
When reporting issues, include:
- Worksheet name
- Operation attempted
- Error message (with correlation ID)
- Relevant log entries
- Excel version

### Documentation
- `README.md` - Project overview
- `WINDOWS_TESTING.md` - Testing guide
- `M22_COMPLETION_SUMMARY.md` - Implementation details
- `M22_AUTOMATED_UI_GENERATION_PLAN.md` - Original plan
- `UI_QUICK_REFERENCE.md` - This guide

---

## Version History

### Version 1.0 (M22 - 2025-10-28)
- Initial release
- 5 production worksheets
- Complete trading workflow
- Position sizing, checklist, heat check, trade entry
- Dashboard navigation
- Full VBA integration

---

**End of Quick Reference Guide**

**For detailed implementation information, see:** `M22_COMPLETION_SUMMARY.md`
**For testing procedures, see:** `WINDOWS_TESTING.md`
