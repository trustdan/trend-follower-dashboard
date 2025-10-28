# M21 Phase 4 - Integration Test Scripts

**Created:** 2025-10-27
**Purpose:** Detailed step-by-step test scripts for Phase 4 manual integration testing
**Duration:** 45-120 minutes (depending on issues encountered)
**Prerequisites:** Phases 1-3 complete (setup, smoke tests, VBA unit tests pass)

---

## Pre-Test Checklist

Before starting Phase 4, verify:

- [ ] Windows environment available
- [ ] TradingPlatform.xlsm exists in windows/ directory
- [ ] trading.db exists in windows/ directory
- [ ] tf-engine.exe exists in windows/ directory (26 MB, CGO-enabled)
- [ ] All VBA unit tests pass (14/14) - Run "VBA Tests" sheet
- [ ] Excel macros enabled
- [ ] Developer tab visible in Excel

**If any checks fail:** Return to Phase 3 and resolve issues first.

---

## Test Environment Setup

### Quick Environment Verification

Open Command Prompt in windows/ directory:

```cmd
# Verify files present
dir TradingPlatform.xlsm
dir trading.db
dir tf-engine.exe

# Verify engine works
tf-engine.exe --version
# Expected: tf-engine version 3.0.0-dev

# Verify database accessible
tf-engine.exe --db trading.db get-settings --format json
# Expected: JSON with 5 settings (Equity_E, RiskPct_r, etc.)
```

**All commands should succeed before proceeding.**

---

## Test Workflow Overview

Phase 4 tests 4 complete workflows through Excel UI:

1. **Position Sizing** - Calculate shares/contracts for trade entry
2. **Checklist Evaluation** - Evaluate 6-item checklist, get GREEN/YELLOW/RED banner
3. **Heat Management** - Verify portfolio/bucket heat caps enforced
4. **Save Decision** - Validate all 5 hard gates prevent bad trades

Each workflow includes:
- Setup instructions (worksheet creation)
- Test data
- Expected results
- Validation criteria
- Common issues and fixes

---

# Workflow 1: Position Sizing

**Purpose:** Test engine-first position sizing calculations through Excel UI
**Duration:** 10-15 minutes
**Sheet:** "Position Sizing"

---

## Setup Instructions

### Step 1.1: Create Position Sizing Worksheet

1. Open TradingPlatform.xlsm
2. Insert new worksheet (right-click tabs → Insert → Worksheet)
3. Rename sheet to: **Position Sizing**
4. Set column widths:
   - Column A: 20 characters
   - Column B: 15 characters

### Step 1.2: Add Headers and Labels

```
Cell A1: "Position Sizing Calculator"
  Format: Bold, 14pt, Dark Blue

Cell A3: "Inputs"
  Format: Bold, 12pt

Cell A4: "Ticker:"
Cell A5: "Entry Price:"
Cell A6: "ATR (N):"
Cell A7: "K Multiple:"
Cell A8: "Method:"

Cell A10: "Optional Overrides"
  Format: Bold, Italic

Cell A11: "Equity (E):"
Cell A12: "Risk % (r):"

Cell A14: "Option Parameters (if applicable)"
  Format: Bold, Italic

Cell A15: "Delta:"
Cell A16: "Max Loss/Contract:"

Cell A20: "Results"
  Format: Bold, 12pt

Cell A21: "Risk Dollars (R):"
Cell A22: "Stop Distance:"
Cell A23: "Initial Stop:"
Cell A24: "Shares:"
Cell A25: "Contracts:"
Cell A26: "Actual Risk:"

Cell A28: "Status"
  Format: Bold, 12pt

Cell A29: "Message:"
```

### Step 1.3: Add Input Validation

Select cell B8 (Method), add data validation:
- Allow: List
- Source: `stock,opt-delta-atr,opt-maxloss`

### Step 1.4: Add Calculate Button

1. Developer tab → Insert → ActiveX Controls → Command Button
2. Position button at cell C18 (between inputs and results)
3. Right-click button → Properties:
   - Name: `btnCalcSize`
   - Caption: `Calculate Position Size`
   - Font: Bold, 11pt
4. Double-click button, add code:

```vba
Private Sub btnCalcSize_Click()
    Dim ws As Worksheet
    Dim cmdResult As TFCommandResult
    Dim sizeResult As TFSizingResult
    Dim corrID As String
    Dim entry As Double
    Dim atr As Double
    Dim k As Long
    Dim method As String

    Set ws = ThisWorkbook.Sheets("Position Sizing")

    ' Clear previous results
    ws.Range("B21:B26").ClearContents
    ws.Range("B29").ClearContents

    ' Validate inputs
    If ws.Range("B4").Value = "" Then
        ws.Range("B29").Value = "❌ Error: Ticker required"
        Exit Sub
    End If

    If Not IsNumeric(ws.Range("B5").Value) Then
        ws.Range("B29").Value = "❌ Error: Entry price must be numeric"
        Exit Sub
    End If

    If Not IsNumeric(ws.Range("B6").Value) Then
        ws.Range("B29").Value = "❌ Error: ATR must be numeric"
        Exit Sub
    End If

    ' Get inputs
    entry = ws.Range("B5").Value
    atr = ws.Range("B6").Value
    k = 2  ' Default K
    If IsNumeric(ws.Range("B7").Value) Then k = ws.Range("B7").Value
    method = ws.Range("B8").Value
    If method = "" Then method = "stock"

    ' Generate correlation ID
    corrID = TFHelpers.GenerateCorrelationID()

    ' Build command
    Dim cmd As String
    cmd = "size --entry " & entry & " --atr " & atr & " --k " & k & " --method " & method

    ' Add optional overrides if provided
    If IsNumeric(ws.Range("B11").Value) Then
        cmd = cmd & " --equity " & ws.Range("B11").Value
    End If
    If IsNumeric(ws.Range("B12").Value) Then
        cmd = cmd & " --risk-pct " & ws.Range("B12").Value
    End If

    ' Add option parameters if method is option-based
    If InStr(method, "opt-") > 0 Then
        If IsNumeric(ws.Range("B15").Value) Then
            cmd = cmd & " --delta " & ws.Range("B15").Value
        End If
        If IsNumeric(ws.Range("B16").Value) Then
            cmd = cmd & " --max-loss " & ws.Range("B16").Value
        End If
    End If

    ' Execute command
    cmdResult = TFEngine.ExecuteCommand(cmd, corrID)

    If cmdResult.Success Then
        ' Parse JSON result
        Call TFHelpers.ParseSizingJSON(cmdResult.JsonOutput, sizeResult)

        If sizeResult.Success Then
            ' Display results
            ws.Range("B21").Value = sizeResult.RiskDollars
            ws.Range("B21").NumberFormat = "$#,##0.00"

            ws.Range("B22").Value = sizeResult.StopDistance
            ws.Range("B22").NumberFormat = "0.00"

            ws.Range("B23").Value = sizeResult.InitialStop
            ws.Range("B23").NumberFormat = "$#,##0.00"

            ws.Range("B24").Value = sizeResult.Shares
            ws.Range("B24").NumberFormat = "0"

            ws.Range("B25").Value = sizeResult.Contracts
            ws.Range("B25").NumberFormat = "0"

            ws.Range("B26").Value = sizeResult.ActualRisk
            ws.Range("B26").NumberFormat = "$#,##0.00"

            ws.Range("B29").Value = "✅ Success (corr_id: " & corrID & ")"
            ws.Range("B29").Interior.Color = RGB(198, 239, 206)  ' Light green
        Else
            ws.Range("B29").Value = "❌ Parse Error: " & sizeResult.ErrorMessage
            ws.Range("B29").Interior.Color = RGB(255, 199, 206)  ' Light red
        End If
    Else
        ws.Range("B29").Value = "❌ Engine Error: " & cmdResult.ErrorOutput
        ws.Range("B29").Interior.Color = RGB(255, 199, 206)  ' Light red
    End If
End Sub
```

5. Exit Design Mode (Developer tab → Design Mode button - toggle off)

---

## Test 1.1: Stock Position Sizing

**Test Data:**
```
Ticker:      AAPL
Entry:       180
ATR:         1.5
K:           2
Method:      stock
```

**Action:**
1. Enter test data in cells B4-B8
2. Leave optional overrides (B11-B12) empty (will use DB defaults)
3. Click "Calculate Position Size" button

**Expected Results:**
```
Risk Dollars (R):      $75.00
Stop Distance:         3.00
Initial Stop:          $177.00
Shares:                25
Contracts:             0
Actual Risk:           $75.00
Status:                ✅ Success (corr_id: YYYYMMDD-HHMMSSFFF-XXXX)
```

**Validation Criteria:**
- [ ] Risk Dollars = $75.00 (0.75% of $10,000 equity)
- [ ] Stop Distance = 3.00 (ATR 1.5 × K 2)
- [ ] Initial Stop = $177.00 (Entry 180 - Stop Distance 3)
- [ ] Shares = 25 (Risk $75 / Stop Distance $3)
- [ ] Contracts = 0 (stock method)
- [ ] Actual Risk = $75.00 (Shares 25 × Stop Distance $3)
- [ ] Status shows success with correlation ID

**Manual Verification:**
```
R = E × r = $10,000 × 0.0075 = $75
Stop Distance = N × K = 1.5 × 2 = 3.00
Initial Stop = Entry - Stop Distance = 180 - 3 = $177
Shares = R / Stop Distance = $75 / $3 = 25
Actual Risk = Shares × Stop Distance = 25 × $3 = $75
```

**Common Issues:**

| Issue | Symptom | Fix |
|-------|---------|-----|
| Button doesn't respond | No results appear | Check Design Mode is OFF |
| "Engine not found" | Error in B29 | Verify EnginePathSetting named range |
| "Database not found" | Error in B29 | Verify DatabasePathSetting named range |
| Exit code -999 | VBA error in B29 | Check TradingSystem_Debug.log for details |
| Wrong risk amount | R ≠ $75 | Check settings: `tf-engine.exe get-settings` |

---

## Test 1.2: Stock Position Sizing with Overrides

**Test Data:**
```
Ticker:      MSFT
Entry:       400
ATR:         3.0
K:           2
Method:      stock
Equity:      20000  (override)
Risk %:      1.0    (override)
```

**Action:**
1. Enter test data including overrides (B11-B12)
2. Click "Calculate Position Size"

**Expected Results:**
```
Risk Dollars (R):      $200.00
Stop Distance:         6.00
Initial Stop:          $394.00
Shares:                33
Contracts:             0
Actual Risk:           $198.00
```

**Validation Criteria:**
- [ ] Risk Dollars = $200.00 (1.0% of $20,000 equity override)
- [ ] Stop Distance = 6.00 (ATR 3.0 × K 2)
- [ ] Shares = 33 (Risk $200 / Stop Distance $6)
- [ ] Actual Risk ≤ Risk Dollars ($198 ≤ $200) ✅

**Key Insight:** Actual Risk ≤ Specified Risk (can't buy fractional shares)

---

## Test 1.3: Option Position Sizing (Delta-ATR Method)

**Test Data:**
```
Ticker:      NVDA
Entry:       500
ATR:         5.0
K:           2
Method:      opt-delta-atr
Delta:       0.30
```

**Action:**
1. Enter test data
2. Click "Calculate Position Size"

**Expected Results:**
```
Risk Dollars (R):      $75.00
Stop Distance:         10.00 (ATR 5.0 × K 2)
Initial Stop:          $490.00
Shares:                0
Contracts:             25  (calculated via delta)
Actual Risk:           $75.00
```

**Validation Criteria:**
- [ ] Contracts > 0 (option method)
- [ ] Shares = 0 (option method)
- [ ] Actual Risk ≤ Risk Dollars

**Option Sizing Math:**
```
Stop Distance = N × K = 5.0 × 2 = 10.00
Equivalent Shares = R / Stop Distance = $75 / $10 = 7.5
Contracts = Equivalent Shares / (Delta × 100) = 7.5 / (0.30 × 100) = 0.25 → rounds to 0

Note: With low delta, may get 0 contracts. This is CORRECT behavior.
```

---

## Test 1.4: Option Position Sizing (Max Loss Method)

**Test Data:**
```
Ticker:      SPY
Entry:       450
Method:      opt-maxloss
Max Loss:    250  ($ per contract)
```

**Action:**
1. Enter test data
2. Click "Calculate Position Size"

**Expected Results:**
```
Risk Dollars (R):      $75.00
Contracts:             0  (Max Loss $250 > Risk $75, can't afford 1 contract)
Actual Risk:           $0.00
```

**Validation Criteria:**
- [ ] Contracts = 0 (can't afford even 1 contract with $75 risk budget)
- [ ] This is CORRECT - system prevents oversized risk

**Adjust Test:** Use Max Loss = $70
```
Expected:
Contracts:             1  (Risk $75 / Max Loss $70 = 1.07 → 1)
Actual Risk:           $70.00
```

---

## Workflow 1 Complete! ✅

**Checklist:**
- [ ] Test 1.1: Stock sizing (default settings)
- [ ] Test 1.2: Stock sizing (with overrides)
- [ ] Test 1.3: Option sizing (delta-ATR method)
- [ ] Test 1.4: Option sizing (max-loss method)
- [ ] All results match expected values
- [ ] Status messages show success with correlation IDs
- [ ] Actual Risk ≤ Specified Risk in all cases

**Time Taken:** _______ minutes

---

# Workflow 2: Checklist Evaluation

**Purpose:** Test 6-item checklist evaluation and banner logic (GREEN/YELLOW/RED)
**Duration:** 10-15 minutes
**Sheet:** "Checklist"

---

## Setup Instructions

### Step 2.1: Create Checklist Worksheet

1. Insert new worksheet, rename to: **Checklist**
2. Set column widths:
   - Column A: 25 characters
   - Column B: 20 characters

### Step 2.2: Add Headers and Ticker Input

```
Cell A1: "Checklist Evaluation"
  Format: Bold, 14pt, Dark Blue

Cell A3: "Ticker:"
Cell B3: (input cell)

Cell A5: "6-Item Checklist"
  Format: Bold, 12pt
```

### Step 2.3: Add ActiveX Checkboxes

Developer tab → Insert → ActiveX Controls → Check Box

Create 6 checkboxes with these properties:

**Checkbox 1 (Row 6):**
- Name: `Check1`
- Caption: `Higher high (price action)`
- Position: A6

**Checkbox 2 (Row 7):**
- Name: `Check2`
- Caption: `Wider range (volatility expansion)`
- Position: A7

**Checkbox 3 (Row 8):**
- Name: `Check3`
- Caption: `Close off low (strength)`
- Position: A8

**Checkbox 4 (Row 9):**
- Name: `Check4`
- Caption: `Liquidity OK (volume adequate)`
- Position: A9

**Checkbox 5 (Row 10):**
- Name: `Check5`
- Caption: `Not overbought (technical check)`
- Position: A10

**Checkbox 6 (Row 11):**
- Name: `Check6`
- Caption: `Bucket OK (sector not overexposed)`
- Position: A11

### Step 2.4: Add Results Section

```
Cell A15: "Results"
  Format: Bold, 12pt

Cell A16: "Banner:"
Cell B16: (result cell - will be color-coded)
  Format: Bold, 16pt, Center-aligned
  Height: 30 points

Cell A17: "Missing Count:"
Cell A18: "Missing Items:"
Cell A19: "Allow Save:"
Cell A20: "Evaluation Time:"

Cell A27: "Status"
  Format: Bold, 12pt
Cell A28: "Message:"
```

### Step 2.5: Add Evaluate Button

1. Insert ActiveX Command Button at C13
2. Properties:
   - Name: `btnEvaluate`
   - Caption: `Evaluate Checklist`
   - Font: Bold, 11pt
3. Double-click button, add code:

```vba
Private Sub btnEvaluate_Click()
    Dim ws As Worksheet
    Dim cmdResult As TFCommandResult
    Dim checkResult As TFChecklistResult
    Dim corrID As String
    Dim ticker As String
    Dim checksStr As String

    Set ws = ThisWorkbook.Sheets("Checklist")

    ' Clear previous results
    ws.Range("B16:B20").ClearContents
    ws.Range("B16").Interior.ColorIndex = xlNone
    ws.Range("B28").ClearContents

    ' Validate ticker
    ticker = Trim(ws.Range("B3").Value)
    If ticker = "" Then
        ws.Range("B28").Value = "❌ Error: Ticker required"
        Exit Sub
    End If

    ' Collect checkbox states (1 = checked, 0 = unchecked)
    checksStr = ""
    checksStr = checksStr & IIf(ws.Check1.Value, "1", "0") & ","
    checksStr = checksStr & IIf(ws.Check2.Value, "1", "0") & ","
    checksStr = checksStr & IIf(ws.Check3.Value, "1", "0") & ","
    checksStr = checksStr & IIf(ws.Check4.Value, "1", "0") & ","
    checksStr = checksStr & IIf(ws.Check5.Value, "1", "0") & ","
    checksStr = checksStr & IIf(ws.Check6.Value, "1", "0")

    ' Generate correlation ID
    corrID = TFHelpers.GenerateCorrelationID()

    ' Build command
    Dim cmd As String
    cmd = "checklist --ticker " & ticker & " --checks " & checksStr

    ' Execute command
    cmdResult = TFEngine.ExecuteCommand(cmd, corrID)

    If cmdResult.Success Then
        ' Parse JSON result
        Call TFHelpers.ParseChecklistJSON(cmdResult.JsonOutput, checkResult)

        If checkResult.Success Then
            ' Display results
            ws.Range("B16").Value = checkResult.Banner
            ws.Range("B17").Value = checkResult.MissingCount
            ws.Range("B18").Value = checkResult.MissingItems
            ws.Range("B19").Value = checkResult.AllowSave
            ws.Range("B20").Value = Now  ' Evaluation timestamp

            ' Color-code banner
            Select Case checkResult.Banner
                Case "GREEN"
                    ws.Range("B16").Interior.Color = RGB(198, 239, 206)  ' Light green
                Case "YELLOW"
                    ws.Range("B16").Interior.Color = RGB(255, 235, 156)  ' Light yellow
                Case "RED"
                    ws.Range("B16").Interior.Color = RGB(255, 199, 206)  ' Light red
            End Select

            ws.Range("B28").Value = "✅ Evaluated (corr_id: " & corrID & ")"
            ws.Range("B28").Interior.Color = RGB(198, 239, 206)
        Else
            ws.Range("B28").Value = "❌ Parse Error: " & checkResult.ErrorMessage
            ws.Range("B28").Interior.Color = RGB(255, 199, 206)
        End If
    Else
        ws.Range("B28").Value = "❌ Engine Error: " & cmdResult.ErrorOutput
        ws.Range("B28").Interior.Color = RGB(255, 199, 206)
    End If
End Sub
```

4. Exit Design Mode

---

## Test 2.1: GREEN Banner (All Checks Pass)

**Test Data:**
```
Ticker: AAPL
Checkboxes: ALL 6 CHECKED
```

**Action:**
1. Enter ticker: AAPL
2. Check ALL 6 checkboxes
3. Click "Evaluate Checklist"

**Expected Results:**
```
Banner:              GREEN (light green background)
Missing Count:       0
Missing Items:       (empty)
Allow Save:          TRUE
Evaluation Time:     (current timestamp)
Status:              ✅ Evaluated (corr_id: XXXXX)
```

**Validation Criteria:**
- [ ] Banner = "GREEN"
- [ ] Banner background = light green (RGB 198, 239, 206)
- [ ] Missing Count = 0
- [ ] Missing Items = empty
- [ ] Allow Save = TRUE
- [ ] Status shows success

**Key Rule:** Only GREEN banner allows save (allow_save = TRUE)

---

## Test 2.2: YELLOW Banner (1-2 Items Missing)

**Test Data:**
```
Ticker: MSFT
Checkboxes: 4 checked, 2 unchecked
  ✓ Higher high
  ✓ Wider range
  ✓ Close off low
  ✓ Liquidity OK
  ✗ Not overbought
  ✗ Bucket OK
```

**Action:**
1. Enter ticker: MSFT
2. Check boxes 1-4, leave 5-6 unchecked
3. Click "Evaluate Checklist"

**Expected Results:**
```
Banner:              YELLOW (light yellow background)
Missing Count:       2
Missing Items:       "Not overbought,Bucket OK"
Allow Save:          FALSE
Status:              ✅ Evaluated (corr_id: XXXXX)
```

**Validation Criteria:**
- [ ] Banner = "YELLOW"
- [ ] Banner background = light yellow (RGB 255, 235, 156)
- [ ] Missing Count = 2
- [ ] Missing Items lists both unchecked items
- [ ] Allow Save = FALSE

**Key Rule:** YELLOW banner prevents save (allow_save = FALSE)

---

## Test 2.3: YELLOW Banner (1 Item Missing - Edge Case)

**Test Data:**
```
Ticker: NVDA
Checkboxes: 5 checked, 1 unchecked
  ✓ Higher high
  ✗ Wider range
  ✓ Close off low
  ✓ Liquidity OK
  ✓ Not overbought
  ✓ Bucket OK
```

**Action:**
1. Enter ticker: NVDA
2. Uncheck only "Wider range"
3. Click "Evaluate Checklist"

**Expected Results:**
```
Banner:              YELLOW (light yellow background)
Missing Count:       1
Missing Items:       "Wider range"
Allow Save:          FALSE
```

**Validation Criteria:**
- [ ] Banner = "YELLOW" (even with only 1 missing)
- [ ] Missing Items shows correct item name
- [ ] Allow Save = FALSE

**Key Insight:** Even ONE missing item prevents save

---

## Test 2.4: RED Banner (3+ Items Missing)

**Test Data:**
```
Ticker: SPY
Checkboxes: 3 checked, 3 unchecked
  ✓ Higher high
  ✓ Wider range
  ✓ Close off low
  ✗ Liquidity OK
  ✗ Not overbought
  ✗ Bucket OK
```

**Action:**
1. Enter ticker: SPY
2. Check boxes 1-3, leave 4-6 unchecked
3. Click "Evaluate Checklist"

**Expected Results:**
```
Banner:              RED (light red background)
Missing Count:       3
Missing Items:       "Liquidity OK,Not overbought,Bucket OK"
Allow Save:          FALSE
```

**Validation Criteria:**
- [ ] Banner = "RED"
- [ ] Banner background = light red (RGB 255, 199, 206)
- [ ] Missing Count = 3
- [ ] Missing Items lists all 3 unchecked items
- [ ] Allow Save = FALSE

**Key Rule:** 3+ missing = RED = VERY BAD, do not trade

---

## Test 2.5: Banner Persistence

**Purpose:** Verify banner evaluations persist (don't auto-clear)

**Action:**
1. Evaluate AAPL with GREEN banner (all checked)
2. Evaluate MSFT with YELLOW banner (2 unchecked)
3. Return to AAPL row - verify banner still shows

**Expected Behavior:**
- Each ticker's evaluation persists independently
- Banner doesn't auto-clear when evaluating different ticker

**Validation:**
- [ ] Previous evaluations remain visible
- [ ] Can evaluate multiple tickers in session

---

## Workflow 2 Complete! ✅

**Checklist:**
- [ ] Test 2.1: GREEN banner (all 6 checks)
- [ ] Test 2.2: YELLOW banner (2 missing)
- [ ] Test 2.3: YELLOW banner (1 missing - edge case)
- [ ] Test 2.4: RED banner (3+ missing)
- [ ] Test 2.5: Banner persistence verified
- [ ] Color coding correct for all banners
- [ ] Allow Save = TRUE only for GREEN
- [ ] Status messages show success

**Time Taken:** _______ minutes

---

# Workflow 3: Heat Management

**Purpose:** Test portfolio and bucket heat cap enforcement (4% portfolio, 1.5% bucket)
**Duration:** 10-15 minutes
**Sheet:** "Heat Check"

---

## Setup Instructions

### Step 3.1: Create Heat Check Worksheet

1. Insert new worksheet, rename to: **Heat Check**
2. Set column widths:
   - Column A: 25 characters
   - Column B: 15 characters

### Step 3.2: Add Headers and Input Section

```
Cell A1: "Heat Management"
  Format: Bold, 14pt, Dark Blue

Cell A3: "New Trade"
  Format: Bold, 12pt

Cell A4: "Risk Dollars (R):"
Cell B4: (input cell)

Cell A5: "Bucket:"
Cell B5: (input cell)

Cell A7: (Button: "Check Heat")
```

Add dropdown to B5:
- Data Validation → List
- Source: `Tech/Comm,Finance,Healthcare,Industrial,Consumer,Energy,Utilities,Real Estate`

### Step 3.3: Add Results Section

```
Cell A9: "Portfolio Heat"
  Format: Bold, 12pt

Cell A10: "Current Heat:"
Cell A11: "New Heat:"
Cell A12: "Heat %:"
Cell A13: "Cap:"
Cell A14: "Exceeded:"
Cell A15: "Overage:"

Cell A17: "Bucket Heat"
  Format: Bold, 12pt

Cell A18: "Current Heat:"
Cell A19: "New Heat:"
Cell A20: "Heat %:"
Cell A21: "Cap:"
Cell A22: "Exceeded:"
Cell A23: "Overage:"

Cell A25: "Result"
  Format: Bold, 12pt

Cell A26: "Allowed:"
Cell B26: (large cell, will be color-coded)
  Format: Bold, 16pt, Center-aligned

Cell A28: "Status"
  Format: Bold, 12pt
Cell A29: "Message:"
```

### Step 3.4: Add Check Heat Button

1. Insert ActiveX Command Button at C7
2. Properties:
   - Name: `btnCheckHeat`
   - Caption: `Check Heat`
   - Font: Bold, 11pt
3. Double-click button, add code:

```vba
Private Sub btnCheckHeat_Click()
    Dim ws As Worksheet
    Dim cmdResult As TFCommandResult
    Dim heatResult As TFHeatResult
    Dim corrID As String
    Dim riskDollars As Double
    Dim bucket As String

    Set ws = ThisWorkbook.Sheets("Heat Check")

    ' Clear previous results
    ws.Range("B10:B15,B18:B23,B26,B29").ClearContents
    ws.Range("B26").Interior.ColorIndex = xlNone

    ' Validate inputs
    If Not IsNumeric(ws.Range("B4").Value) Then
        ws.Range("B29").Value = "❌ Error: Risk Dollars must be numeric"
        Exit Sub
    End If

    bucket = Trim(ws.Range("B5").Value)
    If bucket = "" Then
        ws.Range("B29").Value = "❌ Error: Bucket required"
        Exit Sub
    End If

    riskDollars = ws.Range("B4").Value

    ' Generate correlation ID
    corrID = TFHelpers.GenerateCorrelationID()

    ' Build command
    Dim cmd As String
    cmd = "heat --risk " & riskDollars & " --bucket """ & bucket & """"

    ' Execute command
    cmdResult = TFEngine.ExecuteCommand(cmd, corrID)

    If cmdResult.Success Then
        ' Parse JSON result
        Call TFHelpers.ParseHeatJSON(cmdResult.JsonOutput, heatResult)

        If heatResult.Success Then
            ' Display portfolio heat results
            ws.Range("B10").Value = heatResult.CurrentPortfolioHeat
            ws.Range("B10").NumberFormat = "$#,##0.00"

            ws.Range("B11").Value = heatResult.NewPortfolioHeat
            ws.Range("B11").NumberFormat = "$#,##0.00"

            ws.Range("B12").Value = heatResult.PortfolioHeatPct
            ws.Range("B12").NumberFormat = "0.00%"

            ws.Range("B13").Value = heatResult.PortfolioCap
            ws.Range("B13").NumberFormat = "$#,##0.00"

            ws.Range("B14").Value = heatResult.PortfolioExceeded

            If heatResult.PortfolioExceeded Then
                ws.Range("B15").Value = heatResult.PortfolioOverage
                ws.Range("B15").NumberFormat = "$#,##0.00"
            End If

            ' Display bucket heat results
            ws.Range("B18").Value = heatResult.CurrentBucketHeat
            ws.Range("B18").NumberFormat = "$#,##0.00"

            ws.Range("B19").Value = heatResult.NewBucketHeat
            ws.Range("B19").NumberFormat = "$#,##0.00"

            ws.Range("B20").Value = heatResult.BucketHeatPct
            ws.Range("B20").NumberFormat = "0.00%"

            ws.Range("B21").Value = heatResult.BucketCap
            ws.Range("B21").NumberFormat = "$#,##0.00"

            ws.Range("B22").Value = heatResult.BucketExceeded

            If heatResult.BucketExceeded Then
                ws.Range("B23").Value = heatResult.BucketOverage
                ws.Range("B23").NumberFormat = "$#,##0.00"
            End If

            ' Display final allowed result
            ws.Range("B26").Value = heatResult.Allowed

            If heatResult.Allowed Then
                ws.Range("B26").Interior.Color = RGB(198, 239, 206)  ' Green
            Else
                ws.Range("B26").Interior.Color = RGB(255, 199, 206)  ' Red
            End If

            ws.Range("B29").Value = "✅ Heat checked (corr_id: " & corrID & ")"
            ws.Range("B29").Interior.Color = RGB(198, 239, 206)
        Else
            ws.Range("B29").Value = "❌ Parse Error: " & heatResult.ErrorMessage
            ws.Range("B29").Interior.Color = RGB(255, 199, 206)
        End If
    Else
        ws.Range("B29").Value = "❌ Engine Error: " & cmdResult.ErrorOutput
        ws.Range("B29").Interior.Color = RGB(255, 199, 206)
    End If
End Sub
```

4. Exit Design Mode

---

## Test 3.1: Heat Check (No Open Positions - Portfolio)

**Test Data:**
```
Risk Dollars: 75
Bucket:       Tech/Comm
```

**Action:**
1. Verify no open positions exist (fresh DB):
   ```cmd
   tf-engine.exe --db trading.db --format json -c "SELECT COUNT(*) FROM positions WHERE status='open'"
   # Expected: 0
   ```
2. Enter test data
3. Click "Check Heat"

**Expected Results:**
```
Portfolio Heat:
  Current Heat:        $0.00
  New Heat:            $75.00
  Heat %:              18.75%  (of $400 cap)
  Cap:                 $400.00  (4% of $10,000 equity)
  Exceeded:            FALSE
  Overage:             (empty)

Bucket Heat:
  Current Heat:        $0.00
  New Heat:            $75.00
  Heat %:              50.00%  (of $150 cap)
  Cap:                 $150.00  (1.5% of $10,000 equity)
  Exceeded:            FALSE
  Overage:             (empty)

Allowed:               TRUE (green background)
Status:                ✅ Heat checked (corr_id: XXXXX)
```

**Validation Criteria:**
- [ ] Portfolio Cap = $400 (4% × $10,000)
- [ ] Bucket Cap = $150 (1.5% × $10,000)
- [ ] Heat % calculated correctly
- [ ] Exceeded = FALSE (both)
- [ ] Allowed = TRUE (green)

**Manual Verification:**
```
Portfolio Cap = Equity × HeatCap_H_pct = $10,000 × 0.04 = $400
Bucket Cap = Equity × BucketHeatCap_pct = $10,000 × 0.015 = $150

Portfolio Heat % = New Heat / Cap = $75 / $400 = 18.75%
Bucket Heat % = New Heat / Cap = $75 / $150 = 50.00%

Both under caps → Allowed = TRUE
```

---

## Test 3.2: Heat Check (Portfolio Cap Exceeded)

**Test Data:**
```
Risk Dollars: 450
Bucket:       Tech/Comm
```

**Action:**
1. Enter Risk = $450 (exceeds $400 portfolio cap)
2. Click "Check Heat"

**Expected Results:**
```
Portfolio Heat:
  Current Heat:        $0.00
  New Heat:            $450.00
  Heat %:              112.50%  (OVER 100%!)
  Cap:                 $400.00
  Exceeded:            TRUE
  Overage:             $50.00  ($450 - $400)

Bucket Heat:
  Current Heat:        $0.00
  New Heat:            $450.00
  Heat %:              300.00%  (WAY OVER!)
  Cap:                 $150.00
  Exceeded:            TRUE
  Overage:             $300.00  ($450 - $150)

Allowed:               FALSE (red background)
Status:                ✅ Heat checked (corr_id: XXXXX)
```

**Validation Criteria:**
- [ ] Portfolio Exceeded = TRUE
- [ ] Portfolio Overage = $50 ($450 - $400)
- [ ] Bucket Exceeded = TRUE
- [ ] Bucket Overage = $300 ($450 - $150)
- [ ] Allowed = FALSE (red background)

**Key Rule:** If EITHER cap exceeded → Allowed = FALSE

---

## Test 3.3: Heat Check (Bucket Cap Exceeded, Portfolio OK)

**Test Data:**
```
Risk Dollars: 200
Bucket:       Tech/Comm
```

**Action:**
1. Enter Risk = $200
2. Click "Check Heat"

**Expected Results:**
```
Portfolio Heat:
  Current Heat:        $0.00
  New Heat:            $200.00
  Heat %:              50.00%
  Cap:                 $400.00
  Exceeded:            FALSE

Bucket Heat:
  Current Heat:        $0.00
  New Heat:            $200.00
  Heat %:              133.33%  (OVER 100%)
  Cap:                 $150.00
  Exceeded:            TRUE
  Overage:             $50.00  ($200 - $150)

Allowed:               FALSE (red background)
```

**Validation Criteria:**
- [ ] Portfolio Exceeded = FALSE (under $400)
- [ ] Bucket Exceeded = TRUE (over $150)
- [ ] Allowed = FALSE (red background)

**Key Insight:** Bucket cap is TIGHTER than portfolio cap (1.5% vs 4%)

---

## Test 3.4: Heat Check (Edge Case - Exactly at Cap)

**Test Data:**
```
Risk Dollars: 150
Bucket:       Finance
```

**Action:**
1. Enter Risk = $150 (exactly at bucket cap)
2. Click "Check Heat"

**Expected Results:**
```
Bucket Heat:
  New Heat:            $150.00
  Heat %:              100.00%  (exactly at cap)
  Cap:                 $150.00
  Exceeded:            FALSE  (at cap, not over)

Allowed:               TRUE (green background)
```

**Validation Criteria:**
- [ ] At cap (100%) = Allowed (not exceeded)
- [ ] Exceeded = FALSE

**Key Rule:** At cap = OK, over cap = REJECT

---

## Test 3.5: Heat Check (With Open Positions)

**Prerequisites:** Add an open position to database

```cmd
cd windows
tf-engine.exe --db trading.db --format json save-decision ^
  --ticker AAPL --entry 180 --atr 1.5 --method stock ^
  --banner GREEN --risk 75 --shares 25 --bucket "Tech/Comm" --preset TEST
# This adds $75 risk to Tech/Comm bucket
```

**Test Data:**
```
Risk Dollars: 80
Bucket:       Tech/Comm
```

**Action:**
1. After adding AAPL position (R=$75), check heat for new $80 trade
2. Click "Check Heat"

**Expected Results:**
```
Portfolio Heat:
  Current Heat:        $75.00  (from AAPL position)
  New Heat:            $155.00  ($75 + $80)
  Heat %:              38.75%
  Exceeded:            FALSE

Bucket Heat:
  Current Heat:        $75.00  (from AAPL in Tech/Comm)
  New Heat:            $155.00  ($75 + $80)
  Heat %:              103.33%  (OVER 100%!)
  Cap:                 $150.00
  Exceeded:            TRUE
  Overage:             $5.00  ($155 - $150)

Allowed:               FALSE (red background)
```

**Validation Criteria:**
- [ ] Current Heat includes open positions
- [ ] New Heat = Current + New Trade
- [ ] Bucket exceeded (103.33% > 100%)
- [ ] Allowed = FALSE

**Key Insight:** System tracks cumulative risk across all open positions

---

## Test 3.6: Heat Check (Different Buckets Don't Interfere)

**Prerequisites:** AAPL position exists (R=$75 in Tech/Comm)

**Test Data:**
```
Risk Dollars: 80
Bucket:       Healthcare  (different bucket)
```

**Action:**
1. Check heat for Healthcare bucket (AAPL is in Tech/Comm)
2. Click "Check Heat"

**Expected Results:**
```
Portfolio Heat:
  Current Heat:        $75.00  (from AAPL - portfolio-wide)
  New Heat:            $155.00
  Heat %:              38.75%
  Exceeded:            FALSE

Bucket Heat:
  Current Heat:        $0.00  (Healthcare has no positions)
  New Heat:            $80.00
  Heat %:              53.33%
  Cap:                 $150.00
  Exceeded:            FALSE

Allowed:               TRUE (green background)
```

**Validation Criteria:**
- [ ] Portfolio Heat includes all buckets
- [ ] Bucket Heat isolated per bucket
- [ ] Healthcare bucket starts at $0
- [ ] Allowed = TRUE (both caps OK)

**Key Insight:** Buckets are independent, portfolio is cumulative

---

## Workflow 3 Complete! ✅

**Checklist:**
- [ ] Test 3.1: No open positions (clean state)
- [ ] Test 3.2: Portfolio cap exceeded
- [ ] Test 3.3: Bucket cap exceeded (portfolio OK)
- [ ] Test 3.4: Edge case (exactly at cap)
- [ ] Test 3.5: With open positions (cumulative)
- [ ] Test 3.6: Different buckets (isolation)
- [ ] Heat % calculations correct
- [ ] Allowed = TRUE only when both caps OK
- [ ] Red/green color coding correct

**Time Taken:** _______ minutes

**Cleanup (Optional):**
```cmd
# Remove test position if desired
tf-engine.exe --db trading.db --format json -c "DELETE FROM positions WHERE ticker='AAPL'"
```

---

# Workflow 4: Save Decision (5 Hard Gates)

**Purpose:** Test all 5 hard gates that prevent bad trades
**Duration:** 15-30 minutes
**Sheet:** "Trade Entry" (comprehensive)

---

## The 5 Hard Gates

Save Decision workflow validates ALL 5 gates. Trade is REJECTED if ANY gate fails:

1. **Gate 1: Banner must be GREEN**
   - YELLOW/RED banners rejected
   - Missing checklist items prevent save

2. **Gate 2: Ticker must be in today's candidates**
   - Only tickers imported today allowed
   - Prevents trading tickers not on watchlist

3. **Gate 3: 2-minute impulse brake**
   - Must wait 2 minutes after banner evaluation
   - Prevents impulsive entries on excitement

4. **Gate 4: Bucket not in cooldown**
   - If bucket has recent trade, must wait cooldown period
   - Prevents overtrading same sector

5. **Gate 5: Heat caps not exceeded**
   - Portfolio heat ≤ 4% of equity
   - Bucket heat ≤ 1.5% of equity
   - Both must be OK

---

## Setup Instructions

### Step 4.1: Create Trade Entry Worksheet

1. Insert new worksheet, rename to: **Trade Entry**
2. Set column widths:
   - Column A: 25 characters
   - Column B: 20 characters
   - Column C: 15 characters

### Step 4.2: Add Complete Trade Entry Form

```
Cell A1: "Trade Entry - Save Decision"
  Format: Bold, 14pt, Dark Blue

Cell A3: "Trade Details"
  Format: Bold, 12pt

Cell A4: "Ticker:"
Cell A5: "Entry Price:"
Cell A6: "ATR:"
Cell A7: "Method:"
Cell A8: "Banner:"
Cell A9: "Risk Dollars:"
Cell A10: "Shares:"
Cell A11: "Contracts:"
Cell A12: "Bucket:"
Cell A13: "Preset:"

Cell A15: (Button: "Save Decision")

Cell A17: "Gate Status"
  Format: Bold, 12pt

Cell A18: "Gate 1 - Banner:"
Cell A19: "Gate 2 - In Candidates:"
Cell A20: "Gate 3 - Impulse Brake:"
Cell A21: "Gate 4 - No Cooldown:"
Cell A22: "Gate 5 - Heat Caps OK:"

Cell A24: "Final Result"
  Format: Bold, 12pt

Cell A25: "Decision ID:"
Cell A26: "Saved:"

Cell A28: "Status"
  Format: Bold, 12pt
Cell A29: "Message:"
```

Add dropdowns:
- B7 (Method): `stock,opt-delta-atr,opt-maxloss`
- B8 (Banner): `GREEN,YELLOW,RED`
- B12 (Bucket): `Tech/Comm,Finance,Healthcare,Industrial,Consumer,Energy,Utilities,Real Estate`

### Step 4.3: Add Save Decision Button

1. Insert ActiveX Command Button at C15
2. Properties:
   - Name: `btnSaveDecision`
   - Caption: `Save Decision`
   - Font: Bold, 12pt
   - BackColor: Light Yellow (caution color)
3. Double-click button, add code:

```vba
Private Sub btnSaveDecision_Click()
    Dim ws As Worksheet
    Dim cmdResult As TFCommandResult
    Dim saveResult As TFSaveDecisionResult
    Dim corrID As String
    Dim cmd As String

    Set ws = ThisWorkbook.Sheets("Trade Entry")

    ' Clear previous results
    ws.Range("B18:B22,B25:B26,B29").ClearContents
    ws.Range("B26").Interior.ColorIndex = xlNone

    ' Validate required fields
    If ws.Range("B4").Value = "" Then
        ws.Range("B29").Value = "❌ Error: Ticker required"
        Exit Sub
    End If

    If Not IsNumeric(ws.Range("B5").Value) Then
        ws.Range("B29").Value = "❌ Error: Entry price must be numeric"
        Exit Sub
    End If

    If Not IsNumeric(ws.Range("B6").Value) Then
        ws.Range("B29").Value = "❌ Error: ATR must be numeric"
        Exit Sub
    End If

    If ws.Range("B8").Value = "" Then
        ws.Range("B29").Value = "❌ Error: Banner required"
        Exit Sub
    End If

    If Not IsNumeric(ws.Range("B9").Value) Then
        ws.Range("B29").Value = "❌ Error: Risk Dollars must be numeric"
        Exit Sub
    End If

    If ws.Range("B12").Value = "" Then
        ws.Range("B29").Value = "❌ Error: Bucket required"
        Exit Sub
    End If

    If ws.Range("B13").Value = "" Then
        ws.Range("B29").Value = "❌ Error: Preset required"
        Exit Sub
    End If

    ' Generate correlation ID
    corrID = TFHelpers.GenerateCorrelationID()

    ' Build command
    cmd = "save-decision" & _
          " --ticker " & ws.Range("B4").Value & _
          " --entry " & ws.Range("B5").Value & _
          " --atr " & ws.Range("B6").Value & _
          " --method " & ws.Range("B7").Value & _
          " --banner " & ws.Range("B8").Value & _
          " --risk " & ws.Range("B9").Value & _
          " --bucket """ & ws.Range("B12").Value & """" & _
          " --preset " & ws.Range("B13").Value

    ' Add shares or contracts
    If IsNumeric(ws.Range("B10").Value) And ws.Range("B10").Value > 0 Then
        cmd = cmd & " --shares " & ws.Range("B10").Value
    End If

    If IsNumeric(ws.Range("B11").Value) And ws.Range("B11").Value > 0 Then
        cmd = cmd & " --contracts " & ws.Range("B11").Value
    End If

    ' Execute command
    cmdResult = TFEngine.ExecuteCommand(cmd, corrID)

    If cmdResult.Success Then
        ' Parse JSON result
        Call TFHelpers.ParseSaveDecisionJSON(cmdResult.JsonOutput, saveResult)

        If saveResult.Success Then
            ' Display gate results
            ws.Range("B18").Value = FormatGateResult(saveResult.Gate1_Banner)
            ws.Range("B19").Value = FormatGateResult(saveResult.Gate2_InCandidates)
            ws.Range("B20").Value = FormatGateResult(saveResult.Gate3_ImpulseBrake)
            ws.Range("B21").Value = FormatGateResult(saveResult.Gate4_NoCooldown)
            ws.Range("B22").Value = FormatGateResult(saveResult.Gate5_HeatCaps)

            ' Display final result
            If saveResult.Saved Then
                ws.Range("B25").Value = saveResult.DecisionID
                ws.Range("B26").Value = "✅ SAVED"
                ws.Range("B26").Interior.Color = RGB(198, 239, 206)  ' Green

                ws.Range("B29").Value = "✅ Decision saved (ID: " & saveResult.DecisionID & ", corr_id: " & corrID & ")"
                ws.Range("B29").Interior.Color = RGB(198, 239, 206)

                ' Clear form for next trade
                ws.Range("B4:B13").ClearContents
            Else
                ws.Range("B26").Value = "❌ REJECTED"
                ws.Range("B26").Interior.Color = RGB(255, 199, 206)  ' Red

                ws.Range("B29").Value = "❌ REJECTED: " & saveResult.RejectionReason & " (corr_id: " & corrID & ")"
                ws.Range("B29").Interior.Color = RGB(255, 199, 206)

                ' Keep form populated for review
            End If
        Else
            ws.Range("B29").Value = "❌ Parse Error: " & saveResult.ErrorMessage
            ws.Range("B29").Interior.Color = RGB(255, 199, 206)
        End If
    Else
        ws.Range("B29").Value = "❌ Engine Error: " & cmdResult.ErrorOutput
        ws.Range("B29").Interior.Color = RGB(255, 199, 206)
    End If
End Sub

' Helper function to format gate results
Private Function FormatGateResult(ByVal passed As Boolean) As String
    If passed Then
        FormatGateResult = "✅ PASS"
    Else
        FormatGateResult = "❌ FAIL"
    End If
End Function
```

4. Exit Design Mode

---

## Test Setup: Prepare Test Data

Before testing gates, prepare environment:

### Setup 1: Import Candidates

```cmd
cd windows
tf-engine.exe --db trading.db import-candidates --tickers AAPL,MSFT,NVDA --preset TEST
# Expected: Imported 3 candidates
```

### Setup 2: Clear Existing Positions (Optional)

```cmd
tf-engine.exe --db trading.db --format json -c "DELETE FROM positions"
# Ensures heat starts at $0
```

### Setup 3: Evaluate Checklist for AAPL (GREEN)

Use Checklist sheet:
1. Enter ticker: AAPL
2. Check all 6 boxes
3. Click "Evaluate Checklist"
4. Verify GREEN banner
5. **Wait 2 minutes** (for Gate 3 impulse brake to clear)

---

## Test 4.1: Happy Path (All Gates Pass)

**Prerequisites:**
- AAPL in candidates (Setup 1)
- AAPL has GREEN banner (Setup 3)
- 2 minutes passed since banner evaluation (Setup 3)

**Test Data:**
```
Ticker:        AAPL
Entry:         180
ATR:           1.5
Method:        stock
Banner:        GREEN
Risk Dollars:  75
Shares:        25
Contracts:     0
Bucket:        Tech/Comm
Preset:        TEST
```

**Action:**
1. Enter all test data
2. Click "Save Decision"

**Expected Results:**
```
Gate Status:
  Gate 1 - Banner:           ✅ PASS
  Gate 2 - In Candidates:    ✅ PASS
  Gate 3 - Impulse Brake:    ✅ PASS
  Gate 4 - No Cooldown:      ✅ PASS
  Gate 5 - Heat Caps OK:     ✅ PASS

Final Result:
  Decision ID:               1
  Saved:                     ✅ SAVED (green background)

Status:                      ✅ Decision saved (ID: 1, corr_id: XXXXX)
```

**Validation Criteria:**
- [ ] All 5 gates PASS
- [ ] Decision ID returned (1)
- [ ] Saved = TRUE
- [ ] Form clears for next trade
- [ ] Status shows success

**Database Verification:**
```cmd
tf-engine.exe --db trading.db --format json -c "SELECT * FROM decisions WHERE id=1"
# Should show saved decision with all fields
```

---

## Test 4.2: Gate 1 Rejection (YELLOW Banner)

**Test Data:**
```
(Same as Test 4.1, except:)
Banner:        YELLOW
```

**Action:**
1. Change banner to YELLOW
2. Click "Save Decision"

**Expected Results:**
```
Gate Status:
  Gate 1 - Banner:           ❌ FAIL
  Gate 2 - In Candidates:    ✅ PASS
  Gate 3 - Impulse Brake:    ✅ PASS
  Gate 4 - No Cooldown:      ✅ PASS
  Gate 5 - Heat Caps OK:     ✅ PASS

Final Result:
  Saved:                     ❌ REJECTED (red background)

Status:                      ❌ REJECTED: Banner must be GREEN (corr_id: XXXXX)
```

**Validation Criteria:**
- [ ] Gate 1 FAILS
- [ ] Other gates PASS
- [ ] Saved = FALSE
- [ ] Rejection reason: "Banner must be GREEN"
- [ ] Form retained (not cleared)

**Key Rule:** Only GREEN banner allows save

---

## Test 4.3: Gate 1 Rejection (RED Banner)

**Test Data:**
```
Banner:        RED
```

**Expected Results:**
```
Gate 1 - Banner:           ❌ FAIL
Saved:                     ❌ REJECTED
Status:                    ❌ REJECTED: Banner must be GREEN
```

**Validation:**
- [ ] RED banner also rejected (same as YELLOW)

---

## Test 4.4: Gate 2 Rejection (Not in Candidates)

**Test Data:**
```
Ticker:        ZZZZ  (not in candidate list)
Banner:        GREEN
(All other fields same as happy path)
```

**Action:**
1. Enter ZZZZ (not imported in Setup 1)
2. Click "Save Decision"

**Expected Results:**
```
Gate Status:
  Gate 1 - Banner:           ✅ PASS
  Gate 2 - In Candidates:    ❌ FAIL
  Gate 3 - Impulse Brake:    ❌ FAIL  (no evaluation for ZZZZ)
  Gate 4 - No Cooldown:      ✅ PASS
  Gate 5 - Heat Caps OK:     ✅ PASS

Final Result:
  Saved:                     ❌ REJECTED

Status:                      ❌ REJECTED: Ticker not in today's candidates
```

**Validation Criteria:**
- [ ] Gate 2 FAILS
- [ ] Rejection reason: "Ticker not in today's candidates"
- [ ] Form retained

**Key Rule:** Must be on watchlist (imported candidates)

---

## Test 4.5: Gate 3 Rejection (Impulse Brake Active)

**Prerequisites:**
- Import MSFT as candidate
- Evaluate MSFT checklist (GREEN banner)
- **Immediately** try to save (don't wait 2 minutes)

**Test Data:**
```
Ticker:        MSFT
Entry:         400
ATR:           3.0
Method:        stock
Banner:        GREEN
Risk Dollars:  75
Shares:        25
Bucket:        Tech/Comm
Preset:        TEST
```

**Action:**
1. Evaluate MSFT checklist (GREEN)
2. **Immediately** (within 2 minutes) click "Save Decision"

**Expected Results:**
```
Gate Status:
  Gate 1 - Banner:           ✅ PASS
  Gate 2 - In Candidates:    ✅ PASS
  Gate 3 - Impulse Brake:    ❌ FAIL
  Gate 4 - No Cooldown:      ✅ PASS
  Gate 5 - Heat Caps OK:     ✅ PASS

Final Result:
  Saved:                     ❌ REJECTED

Status:                      ❌ REJECTED: Must wait XX seconds (impulse brake)
```

**Validation Criteria:**
- [ ] Gate 3 FAILS
- [ ] Rejection shows remaining seconds
- [ ] Example: "Must wait 87 seconds"

**Follow-up Test:**
1. Wait for remaining time to elapse
2. Click "Save Decision" again
3. Verify Gate 3 now PASSES

**Key Rule:** 2-minute cooling period prevents impulsive trades

---

## Test 4.6: Gate 4 Rejection (Bucket in Cooldown)

**Prerequisites:**
- AAPL decision saved (Test 4.1) in Tech/Comm bucket
- Bucket now in cooldown (default: 24 hours)

**Test Data:**
```
Ticker:        NVDA  (also Tech/Comm bucket)
Entry:         500
ATR:           5.0
Method:        stock
Banner:        GREEN
Risk Dollars:  75
Shares:        15
Bucket:        Tech/Comm  (same bucket as AAPL)
Preset:        TEST
```

**Action:**
1. After saving AAPL (Tech/Comm), try to save NVDA (also Tech/Comm)
2. Click "Save Decision"

**Expected Results:**
```
Gate Status:
  Gate 1 - Banner:           ✅ PASS
  Gate 2 - In Candidates:    ✅ PASS
  Gate 3 - Impulse Brake:    ✅ PASS
  Gate 4 - No Cooldown:      ❌ FAIL
  Gate 5 - Heat Caps OK:     ✅ PASS

Final Result:
  Saved:                     ❌ REJECTED

Status:                      ❌ REJECTED: Bucket Tech/Comm in cooldown (XX hours remaining)
```

**Validation Criteria:**
- [ ] Gate 4 FAILS
- [ ] Rejection shows remaining cooldown time
- [ ] Form retained

**Key Rule:** Can't trade same sector twice in cooldown period

**Alternative Test (Different Bucket):**
```
Ticker:        JPM
Bucket:        Finance  (different bucket)
# Expected: Gate 4 PASSES (different bucket)
```

---

## Test 4.7: Gate 5 Rejection (Heat Caps Exceeded)

**Test Data:**
```
Ticker:        AAPL
Entry:         180
ATR:           1.5
Method:        stock
Banner:        GREEN
Risk Dollars:  450  (exceeds both caps)
Shares:        150
Bucket:        Tech/Comm
Preset:        TEST
```

**Action:**
1. Enter Risk = $450 (over $400 portfolio cap, $150 bucket cap)
2. Click "Save Decision"

**Expected Results:**
```
Gate Status:
  Gate 1 - Banner:           ✅ PASS
  Gate 2 - In Candidates:    ✅ PASS
  Gate 3 - Impulse Brake:    ✅ PASS
  Gate 4 - No Cooldown:      ✅ PASS
  Gate 5 - Heat Caps OK:     ❌ FAIL

Final Result:
  Saved:                     ❌ REJECTED

Status:                      ❌ REJECTED: Portfolio heat cap exceeded (overage: $XX)
```

**Validation Criteria:**
- [ ] Gate 5 FAILS
- [ ] Rejection specifies which cap exceeded (portfolio or bucket)
- [ ] Shows overage amount

**Alternative Test (Bucket Only):**
```
Risk Dollars:  200  (under $400 portfolio, over $150 bucket)
# Expected: Gate 5 FAILS (bucket cap exceeded)
```

---

## Test 4.8: Multiple Gate Failures

**Test Data:**
```
Ticker:        ZZZZ  (not in candidates)
Banner:        YELLOW  (not GREEN)
Risk Dollars:  450  (exceeds caps)
(All other fields valid)
```

**Action:**
1. Enter data with 3 gate violations
2. Click "Save Decision"

**Expected Results:**
```
Gate Status:
  Gate 1 - Banner:           ❌ FAIL
  Gate 2 - In Candidates:    ❌ FAIL
  Gate 3 - Impulse Brake:    ❌ FAIL
  Gate 4 - No Cooldown:      ✅ PASS
  Gate 5 - Heat Caps OK:     ❌ FAIL

Final Result:
  Saved:                     ❌ REJECTED

Status:                      ❌ REJECTED: [First failing gate's reason]
```

**Validation Criteria:**
- [ ] Multiple gates FAIL
- [ ] System shows first rejection reason
- [ ] All gate statuses visible

**Key Rule:** ANY gate failure = REJECT (short-circuit logic)

---

## Test 4.9: Form Behavior on Success

**Purpose:** Verify form clears after successful save

**Action:**
1. Save successful AAPL trade (happy path)
2. Observe form state

**Expected Behavior:**
- [ ] All input fields (B4-B13) cleared
- [ ] Gate status cleared
- [ ] Ready for next trade entry
- [ ] Success message persists in B29

**Purpose:** Prevents accidental duplicate saves

---

## Test 4.10: Form Behavior on Rejection

**Purpose:** Verify form persists on rejection

**Action:**
1. Try to save with YELLOW banner (rejection)
2. Observe form state

**Expected Behavior:**
- [ ] All input fields retained
- [ ] Gate status shows which gates failed
- [ ] Rejection message in B29
- [ ] Can fix issue and retry

**Purpose:** Allows correction and resubmission

---

## Workflow 4 Complete! ✅

**Checklist:**
- [ ] Test 4.1: Happy path (all gates pass)
- [ ] Test 4.2: Gate 1 rejection (YELLOW banner)
- [ ] Test 4.3: Gate 1 rejection (RED banner)
- [ ] Test 4.4: Gate 2 rejection (not in candidates)
- [ ] Test 4.5: Gate 3 rejection (impulse brake)
- [ ] Test 4.6: Gate 4 rejection (bucket cooldown)
- [ ] Test 4.7: Gate 5 rejection (heat caps)
- [ ] Test 4.8: Multiple gate failures
- [ ] Test 4.9: Form clears on success
- [ ] Test 4.10: Form persists on rejection
- [ ] Database verification (saved records exist)
- [ ] All rejection messages clear and actionable

**Time Taken:** _______ minutes

**Database Verification:**
```cmd
# Check decisions table
tf-engine.exe --db trading.db --format json -c "SELECT COUNT(*) FROM decisions"
# Should show number of successful saves

# Check positions table
tf-engine.exe --db trading.db --format json -c "SELECT * FROM positions WHERE status='open'"
# Should show open positions with risk tracking
```

---

# Phase 4 Complete! ✅

## Final Checklist

**All 4 Workflows Tested:**
- [ ] Workflow 1: Position Sizing (4 tests)
- [ ] Workflow 2: Checklist Evaluation (5 tests)
- [ ] Workflow 3: Heat Management (6 tests)
- [ ] Workflow 4: Save Decision (10 tests)

**Total Tests:** 25 integration tests

**Critical Validations:**
- [ ] Position sizing: Actual Risk ≤ Specified Risk
- [ ] Checklist: Only GREEN allows save
- [ ] Heat: Both caps enforced (portfolio 4%, bucket 1.5%)
- [ ] Save Decision: All 5 gates enforced
- [ ] Error messages clear and actionable
- [ ] Correlation IDs present for debugging
- [ ] Logs capture all operations

**System Behavior:**
- [ ] VBA thin bridge works (no business logic in VBA)
- [ ] Go engine handles all calculations
- [ ] JSON parsing reliable
- [ ] Shell execution stable
- [ ] Named ranges resolve correctly
- [ ] Database operations succeed

---

## Post-Test Actions

### 1. Review Logs

Check logs for any unexpected errors:

```cmd
# VBA log
notepad TradingSystem_Debug.log

# Engine log (if exists)
notepad tf-engine.log
```

Look for:
- [ ] All correlation IDs present
- [ ] No ERROR entries (except expected rejections)
- [ ] Execution times reasonable (<1 second per command)

### 2. Database Integrity Check

```cmd
# Verify schema intact
sqlite3 trading.db ".schema"

# Check record counts
tf-engine.exe --db trading.db --format json -c "SELECT COUNT(*) FROM candidates"
tf-engine.exe --db trading.db --format json -c "SELECT COUNT(*) FROM decisions"
tf-engine.exe --db trading.db --format json -c "SELECT COUNT(*) FROM positions"
```

### 3. Performance Metrics

Calculate average execution time:
1. Open TradingSystem_Debug.log
2. Find 10 command executions
3. Calculate average time from start to "Command succeeded"

**Expected:** < 500ms per command

---

## Known Issues / Edge Cases

Document any issues encountered during testing:

### Issue Template

```
**Issue:** [Brief description]
**Test:** [Which test exposed it]
**Symptom:** [What happened]
**Expected:** [What should happen]
**Workaround:** [Temporary fix if any]
**Status:** [Fixed / Open / Won't Fix]
```

---

## Next Steps

After Phase 4 completion:

1. **Document Results:**
   - Update M21_PROGRESS.md with Phase 4 status
   - Note any issues found
   - Record actual time taken vs. estimates

2. **Create M21 Completion Summary:**
   - All phases complete
   - Total time spent
   - Issues fixed
   - System ready for use

3. **Begin M22 (if applicable):**
   - Review project plan
   - Determine next milestone
   - Start planning M22 work

---

**Test Scripts Created:** 2025-10-27
**Ready for Execution:** Yes
**Estimated Duration:** 45-120 minutes (full Phase 4)
**Prerequisites:** Phases 1-3 complete

**Good luck with testing! The system is designed to enforce discipline - if any gate rejects a trade, it's working as intended.** 🎯
