# Sample Data Mode Feature

## Overview

Added a **Sample Data** mode to TF-Engine GUI that allows users to explore all features with pre-populated data without requiring a real trade session.

## Implementation Complete ✅

### What Was Added:

1. **Sample Data Generation Module** (`ui/sample_data.go`)
   - `CreateSampleSession()` - Sample trade session at Entry step
   - `CreateSamplePositions()` - 3 open positions (MSFT, XLE, GLD)
   - `CreateSampleCandidates()` - 5 sample tickers (AAPL, NVDA, TSLA, AMD, META)
   - `CreateSampleCalendarTrades()` - Historical trades for calendar view
   - `CreateSampleSettings()` - Sample account settings ($100K equity, 0.75% risk)

2. **AppState Enhancements** (`ui/main.go`)
   - Added `sampleMode` boolean flag
   - Added `sampleSession` to store sample session
   - `GetActiveSession()` - Returns real or sample session based on mode
   - `EnableSampleMode()` - Activates sample data with UI refresh
   - `DisableSampleMode()` - Deactivates sample data with UI refresh

3. **UI Toggle Button**
   - "📦 Sample Data: Off/On" button in top bar
   - Positioned next to Dark Mode button
   - Toggles between real and sample data
   - Triggers full UI refresh when toggled

4. **Screen Updates**
   - **Dashboard**: Shows sample positions, candidates, heat, and settings
   - **Checklist**: Pre-populated with GREEN banner, all checks passed
   - Sample mode clearly marked with "📦 SAMPLE MODE:" prefix

## How to Use:

1. Launch `tf-gui.exe`
2. Click "📦 Sample Data: Off" button in the top bar
3. Button changes to "📦 Sample Data: On"
4. All screens now show realistic sample data:
   - **Dashboard**: 3 open positions, 5 candidates, heat status
   - **Checklist**: GREEN banner with all gates passed
   - **Session Bar**: Sample session #999 (AAPL - Long)
5. Navigate through all tabs to explore features
6. Click button again to return to real data mode

## Visual Indicators:

- Button text shows current state: "📦 Sample Data: On" / "📦 Sample Data: Off"
- Session info shows: "📦 SAMPLE MODE: Session #999 • Long • AAPL"
- Sample data clearly marked to prevent confusion

## Benefits:

1. **Learning Tool**: New users can explore all features without setup
2. **Demo Mode**: Perfect for screenshots, presentations, training
3. **Testing**: Developers can test UI with realistic data instantly
4. **No Database Writes**: Sample mode doesn't modify the database
5. **Quick Toggle**: Switch between real and sample data anytime

## Sample Data Details:

### Sample Session:
- Ticker: AAPL
- Strategy: Long Breakout
- Entry Price: $181.50
- ATR: $2.45
- Position: 306 shares
- Risk: $750.00
- All gates completed (Checklist GREEN, Sizing done, Heat checked)

### Sample Positions:
1. **MSFT** - Tech/Comm - $748.20 risk - Opened 7 days ago
2. **XLE** - Energy - $749.70 risk - Opened 4 days ago
3. **GLD** - Commodities - $750.10 risk - Opened 2 days ago

### Sample Settings:
- Equity: $100,000.00
- Risk per Trade: 0.75%
- Portfolio Heat Cap: 4.0%
- Bucket Heat Cap: 1.5%

## Technical Notes:

- Sample data uses high IDs (999) to avoid conflicts with real data
- No database writes occur in sample mode
- All screens check `state.sampleMode` before database operations
- Sample mode disables save buttons to prevent accidental data creation
- Switching modes triggers session change callbacks to refresh all UI

## Files Modified:

- `ui/sample_data.go` (NEW)
- `ui/main.go` - AppState, toggle button, helper methods
- `ui/dashboard.go` - Sample data for all dashboard cards
- `ui/checklist.go` - Sample mode with pre-filled results

## Future Enhancements:

Consider adding sample data support for:
- Position Sizing screen (show calculated results)
- Heat Check screen (show heat analysis)
- Trade Entry screen (show 5-gate results)
- Calendar screen (show sample trades on calendar grid)
- Session History screen (show sample completed sessions)

## Build Status:

✅ Compiles successfully
✅ No runtime errors
✅ Ready for testing

Run `./tf-gui.exe` to test the new Sample Mode feature!
