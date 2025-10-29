# Demo Worksheet Guide - Testing Sample Data

**Version:** M24 - Demo Worksheet Feature
**Created:** 2025-10-28
**Purpose:** Test macros with realistic sample trade data

---

## Overview

The Demo worksheet feature provides one-click sample data population for testing all macros. This is perfect for:

- ✅ **Testing macros** after setup or updates
- ✅ **Demonstrating the system** to others
- ✅ **Learning the workflow** with realistic examples
- ✅ **Verifying fixes** after troubleshooting

---

## Available Functions

### 1. PopulateSampleData() ⭐
**What it does:** Fills ALL worksheets with sample data in one click

**Data populated:**
- **Position Sizing:** AAPL stock (180, ATR 1.5, K=2)
- **Checklist:** NVDA with YELLOW banner (5 TRUE, 1 FALSE)
- **Heat Check:** MSFT ($750 risk, Tech/Comm bucket)
- **Trade Entry:** TSLA GO decision (GREEN banner)

**How to run:**
1. Press Alt+F8
2. Select `PopulateSampleData`
3. Click Run
4. Message box confirms data populated
5. Go test each worksheet's macro buttons!

---

### 2. ClearSampleData()
**What it does:** Clears all sample data from worksheets

**Clears:**
- Position Sizing inputs and results
- Checklist inputs and results
- Heat Check inputs and results
- Trade Entry inputs and results

**How to run:**
1. Press Alt+F8
2. Select `ClearSampleData`
3. Click Run
4. Confirm "Yes" when prompted
5. All worksheets cleared!

---

### 3. LoadScenario1_SimpleStock() 📊
**What it does:** Loads Scenario 1 - Simple stock trade (all GREEN)

**Data:**
- **Ticker:** AAPL
- **Entry:** $180
- **ATR:** 1.5
- **Method:** stock
- **Checklist:** All 6 items TRUE (GREEN banner)
- **Expected:** ~250 shares, GREEN, PASS 5 gates

**Use case:** Test perfect setup - everything passes

---

### 4. LoadScenario2_YellowBanner() ⚠️
**What it does:** Loads Scenario 2 - YELLOW banner (1 failed check)

**Data:**
- **Ticker:** NVDA
- **Entry:** $450
- **ATR:** 12.0
- **Method:** stock
- **Checklist:** 5 TRUE, 1 FALSE (earnings check fails)
- **Expected:** ~31 shares, YELLOW, FAIL gate 1

**Use case:** Test rejection - banner not GREEN

---

### 5. LoadScenario3_OptionTrade() 📈
**What it does:** Loads Scenario 3 - Option trade with delta-ATR

**Data:**
- **Ticker:** TSLA
- **Entry:** $250
- **ATR:** 8.5
- **Method:** opt-delta-atr
- **Delta:** 0.65
- **Checklist:** All 6 items TRUE (GREEN banner)
- **Expected:** Option contracts calculated, GREEN, PASS 5 gates

**Use case:** Test option position sizing

---

## How to Use

### Method 1: Run from Macro List (Alt+F8)

**To populate sample data:**
```
1. Press Alt+F8
2. Select: PopulateSampleData
3. Click Run
4. See confirmation message
5. Test macros on each sheet!
```

**To clear sample data:**
```
1. Press Alt+F8
2. Select: ClearSampleData
3. Click Run
4. Confirm Yes
5. Worksheets cleared!
```

**To load a specific scenario:**
```
1. Press Alt+F8
2. Select: LoadScenario1_SimpleStock (or 2, or 3)
3. Click Run
4. See scenario description
5. Test the scenario!
```

---

### Method 2: Add Buttons to Demo Worksheet (Recommended!)

Create a Demo worksheet with buttons:

**Step 1: Create Demo Sheet**
1. Right-click sheet tabs → Insert → Worksheet
2. Rename to "Demo"
3. Add title in A1: "Demo & Testing"

**Step 2: Add Buttons**

**Button 1: Populate Sample Data**
1. Developer tab → Insert → Button (Form Control)
2. Draw button (e.g., B3:D4)
3. Assign macro: `PopulateSampleData`
4. Edit text: "Populate Sample Data"

**Button 2: Clear Sample Data**
1. Developer tab → Insert → Button (Form Control)
2. Draw button (e.g., B6:D7)
3. Assign macro: `ClearSampleData`
4. Edit text: "Clear Sample Data"

**Button 3: Scenario 1 - Simple Stock**
1. Developer tab → Insert → Button (Form Control)
2. Draw button (e.g., B10:D11)
3. Assign macro: `LoadScenario1_SimpleStock`
4. Edit text: "Scenario 1: Simple Stock (GREEN)"

**Button 4: Scenario 2 - Yellow Banner**
1. Developer tab → Insert → Button (Form Control)
2. Draw button (e.g., B13:D14)
3. Assign macro: `LoadScenario2_YellowBanner`
4. Edit text: "Scenario 2: Yellow Banner (FAIL)"

**Button 5: Scenario 3 - Option Trade**
1. Developer tab → Insert → Button (Form Control)
2. Draw button (e.g., B16:D17)
3. Assign macro: `LoadScenario3_OptionTrade`
4. Edit text: "Scenario 3: Option Trade"

**Step 3: Add Labels**

```
A1: "Demo & Testing"

A3: "General Sample Data:"
B3: [Populate Sample Data button]
B6: [Clear Sample Data button]

A9: "Test Scenarios:"
B10: [Scenario 1 button] - "AAPL - All GREEN"
B13: [Scenario 2 button] - "NVDA - YELLOW Banner"
B16: [Scenario 3 button] - "TSLA - Option Trade"

A19: "Instructions:"
A20: "1. Click a scenario or sample data button"
A21: "2. Go to Position Sizing, Checklist, Heat Check, Trade Entry"
A22: "3. Click the macro buttons to test"
A23: "4. Verify results match expected outcomes"
A24: "5. Click Clear Sample Data when done"
```

---

## Sample Data Details

### Position Sizing Sheet
```
Ticker:         AAPL
Entry Price:    $180
ATR (N):        1.5
K Multiple:     2
Method:         stock

Expected Results:
  Risk Dollars: $750
  Stop Distance: 3.00
  Initial Stop: $177
  Shares: 250
  Actual Risk: $750
```

### Checklist Sheet
```
Ticker: NVDA

From preset?        TRUE
Trend pass?         TRUE
Liquidity OK?       TRUE
TV confirm?         TRUE
Earnings OK?        FALSE  ← Fails!
Journal OK?         TRUE

Expected Results:
  Banner: YELLOW
  Missing Count: 1
  Missing Items: "Earnings risk"
  Allow Save: FALSE
```

### Heat Check Sheet
```
Ticker:         MSFT
Risk Amount:    $750
Bucket:         Tech/Comm

Expected Results:
  Portfolio heat within 4% cap
  Bucket heat within 1.5% cap
  Both: Within caps (assuming no existing positions)
```

### Trade Entry Sheet
```
Ticker:         TSLA
Entry Price:    $250
ATR:            8.5
Method:         stock
Banner Status:  GREEN
Bucket:         Auto/Transport
Preset:         TF_BREAKOUT_LONG

Expected Results:
  5 Gates checked
  Decision saved (if all gates pass)
```

---

## Testing Workflow

### Complete Test Cycle

**1. Populate Sample Data**
```
Alt+F8 → PopulateSampleData → Run
```

**2. Test Position Sizing**
```
Go to Position Sizing sheet
Click "Calculate" button
Verify: 250 shares, $750 risk, $177 stop
```

**3. Test Checklist**
```
Go to Checklist sheet
Click "Evaluate" button
Verify: YELLOW banner, 1 missing (Earnings)
```

**4. Test Heat Check**
```
Go to Heat Check sheet
Click "Check Heat" button
Verify: Within caps (or see current heat)
```

**5. Test Trade Entry**
```
Go to Trade Entry sheet
Click "Save GO" button
Verify: Gates checked, result shown
```

**6. Clear Sample Data**
```
Alt+F8 → ClearSampleData → Run → Yes
```

---

## Scenario Testing

### Scenario 1: Perfect Trade (GREEN)
```
1. Alt+F8 → LoadScenario1_SimpleStock → Run
2. Test Position Sizing: Should get ~250 shares
3. Test Checklist: Should get GREEN banner
4. Test Heat Check: Should be within caps
5. Test Trade Entry: Should PASS all 5 gates
```

**Expected:** Everything passes, trade accepted

---

### Scenario 2: Failed Checklist (YELLOW)
```
1. Alt+F8 → LoadScenario2_YellowBanner → Run
2. Test Position Sizing: Should get ~31 shares
3. Test Checklist: Should get YELLOW banner (1 failed)
4. Test Trade Entry: Should FAIL gate 1 (banner not GREEN)
```

**Expected:** Trade rejected due to YELLOW banner

---

### Scenario 3: Option Trade
```
1. Alt+F8 → LoadScenario3_OptionTrade → Run
2. Test Position Sizing: Should calculate option contracts
3. Test Checklist: Should get GREEN banner
4. Test Heat Check: Should be within caps
5. Test Trade Entry: Should PASS all 5 gates
```

**Expected:** Option trade sized correctly, all gates pass

---

## Troubleshooting

### "Subscript out of range"

**Cause:** Worksheet doesn't exist with exact name

**Fix:** Make sure you have sheets named:
- "Position Sizing"
- "Checklist"
- "Heat Check"
- "Trade Entry"

### Functions Not in Macro List

**Cause:** VBA modules not imported

**Fix:**
```cmd
cd C:\Users\Dan\excel-trading-platform\release\TradingEngine-v3
fix-vba-modules.bat
```

### Buttons Don't Work

**Cause:** Macros disabled or wrong macro assigned

**Fix:**
1. Click "Enable Content" in Excel
2. Right-click button → Assign Macro
3. Select the correct function
4. Click OK

---

## Creating Your Own Scenarios

You can create custom scenario functions:

```vba
Public Sub LoadScenario4_MyCustomScenario()
    On Error Resume Next
    Dim ws As Worksheet

    Application.ScreenUpdating = False

    ' Position Sizing
    Set ws = ThisWorkbook.Worksheets("Position Sizing")
    ws.Range("B4").Value = "AMZN"
    ws.Range("B5").Value = 150
    ws.Range("B6").Value = 5.0
    ws.Range("B7").Value = 2
    ws.Range("B8").Value = "stock"

    ' Checklist
    Set ws = ThisWorkbook.Worksheets("Checklist")
    ws.Range("B3").Value = "AMZN"
    ws.Range("B5:B10").Value = "TRUE"  ' All green

    ' Heat Check
    Set ws = ThisWorkbook.Worksheets("Heat Check")
    ws.Range("B3").Value = "AMZN"
    ws.Range("B4").Value = 750
    ws.Range("B5").Value = "Tech/Comm"

    ' Trade Entry
    Set ws = ThisWorkbook.Worksheets("Trade Entry")
    ws.Range("B4").Value = "AMZN"
    ws.Range("B5").Value = 150
    ws.Range("B6").Value = 5.0
    ws.Range("B7").Value = "stock"
    ws.Range("B8").Value = "GREEN"
    ws.Range("B11").Value = "Tech/Comm"
    ws.Range("B12").Value = "TF_BREAKOUT_LONG"

    Application.ScreenUpdating = True

    MsgBox "Custom Scenario Loaded: AMZN", vbInformation
End Sub
```

Add this to TFEngine.bas, then assign to a button!

---

## Summary

| Function | Purpose | Use Case |
|----------|---------|----------|
| `PopulateSampleData()` | Fill all sheets | Quick general testing |
| `ClearSampleData()` | Clear all sheets | Start fresh |
| `LoadScenario1_SimpleStock()` | GREEN scenario | Test perfect setup |
| `LoadScenario2_YellowBanner()` | YELLOW scenario | Test rejection |
| `LoadScenario3_OptionTrade()` | Option scenario | Test option sizing |

**All functions log to:** `TradingSystem_Debug.log`

---

## Benefits

✅ **Fast testing** - One click to populate all sheets
✅ **Realistic data** - Real tickers, realistic prices and ATR
✅ **Multiple scenarios** - Test different outcomes
✅ **Easy cleanup** - One click to clear all
✅ **Great for demos** - Show others how system works
✅ **Learning tool** - See examples of each trade type

---

## Next Steps

After running `fix-vba-modules.bat`:

1. ✅ Create Demo worksheet (optional, but recommended)
2. ✅ Add buttons for each function
3. ✅ Test with `PopulateSampleData()`
4. ✅ Try each scenario (1, 2, 3)
5. ✅ Test all macro buttons
6. ✅ Clear with `ClearSampleData()`

---

**Happy Testing!** 🚀

**Last Updated:** 2025-10-28
**Feature Version:** M24 (Demo Worksheet & Testing Features)
