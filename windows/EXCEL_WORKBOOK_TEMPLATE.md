## Excel Workbook Template Structure
# TradingPlatform.xlsm

**Created:** 2025-10-27 (M20 - Windows Integration Package)
**Purpose:** Template structure for Excel workbook (to be created in Windows)

---

## Workbook Creation Instructions

Since .xlsm files cannot be created in Linux, follow these steps in Windows:

### Step 1: Create Blank Macro-Enabled Workbook
1. Open Excel on Windows
2. File > New > Blank Workbook
3. File > Save As
4. Choose file type: "Excel Macro-Enabled Workbook (*.xlsm)"
5. Name: `TradingPlatform.xlsm`
6. Save in the windows/ directory

### Step 2: Enable Developer Tab
1. File > Options > Customize Ribbon
2. Check "Developer" checkbox
3. Click OK

### Step 3: Enable Trust Access to VBA Project
1. File > Options > Trust Center > Trust Center Settings
2. Macro Settings tab
3. Check "Trust access to the VBA project object model"
4. Click OK

### Step 4: Import VBA Modules
1. Run `windows-import-vba.bat` from windows/ directory
2. Verify modules imported: Alt+F11 to open VBA Editor
3. Should see: TFTypes, TFHelpers, TFEngine, TFTests

---

## Worksheet Structure

### 1. Setup Sheet

**Purpose:** Configuration and connection testing

**Named Ranges:**
- `EnginePathSetting` (B2) - Path to tf-engine.exe (default: .\tf-engine.exe)
- `DatabasePathSetting` (B3) - Path to trading.db (default: .\trading.db)

**Layout:**

```
Row 1: [Header] Trading Engine v3 - Setup

Row 3: Configuration
Row 4: Engine Path:     [B4: EnginePathSetting]    [Button: Test Connection]
Row 5: Database Path:   [B5: DatabasePathSetting]   [Button: Initialize DB]

Row 7: Logs
Row 8: [Button: Open VBA Log]    [Button: Open Engine Log]

Row 10: Status
Row 11: Connection Status: [B11: Dynamic status message]
```

**Buttons:**
- **Test Connection** - Runs tf-engine.exe --version to verify engine works
- **Initialize DB** - Calls `Engine_Init()` to create database schema
- **Open VBA Log** - Opens TradingSystem_Debug.log in Notepad
- **Open Engine Log** - Opens tf-engine.log in Notepad

---

### 2. VBA Tests Sheet

**Purpose:** Run VBA unit tests before trading

**Layout:**

```
Row 1: [Header] VBA Unit Tests

Row 3: [Button: Run All Tests]    Status: [D3: Test status]

Row 5: Test Results Table
Headers: Test Name | Result | Message | Duration

Rows 6+: Test results populated by TFTests.RunAllTests()

Row 50+: Summary section (populated after test run)
```

**Auto-formatting:**
- ✅ PASS cells: Green background (RGB 198, 239, 206)
- ❌ FAIL cells: Red background (RGB 255, 199, 206)
- Summary cell: Bold, larger font

---

### 3. Position Sizing Sheet

**Purpose:** Calculate position size for a trade

**Layout:**

```
Row 1: [Header] Position Sizing Calculator

Row 3: Inputs
Row 4: Ticker:          [B4]
Row 5: Entry Price:     [B5: Currency format]
Row 6: ATR (N):         [B6: 2 decimal]
Row 7: K Multiple:      [B7: Whole number, default: 2]
Row 8: Method:          [B8: Dropdown: "stock", "opt-delta-atr", "opt-maxloss"]

Row 10: Optional (leave blank to use settings from DB)
Row 11: Equity:         [B11: Currency]
Row 12: Risk %:         [B12: Percentage, e.g., 0.75%]

Row 14: Option-specific (if applicable)
Row 15: Delta:          [B15: Decimal, e.g., 0.30]
Row 16: Max Loss:       [B16: Currency per contract]

Row 18: [Button: Calculate]    [Button: Clear]

Row 20: Results
Row 21: Risk Dollars (R):      [B21: Currency]
Row 22: Stop Distance:         [B22: 2 decimal]
Row 23: Initial Stop:          [B23: Currency]
Row 24: Shares:                [B24: Whole number]
Row 25: Contracts:             [B25: Whole number]
Row 26: Actual Risk:           [B26: Currency]

Row 28: Status
Row 29: [B29: Status message with correlation ID]
```

**Button Code (Calculate):**
```vba
Sub Button_CalcSize_Click()
    Dim cmdResult As TFCommandResult
    Dim sizeResult As TFSizingResult
    Dim corrID As String

    corrID = TFHelpers.GenerateCorrelationID()

    cmdResult = TFEngine.Engine_Size( _
        Range("B5").Value, _
        Range("B6").Value, _
        Range("B8").Value, _
        corrID:=corrID)

    If cmdResult.Success Then
        sizeResult = TFHelpers.ParseSizingJSON(cmdResult.JsonOutput)
        Range("B21").Value = sizeResult.RiskDollars
        Range("B22").Value = sizeResult.StopDistance
        Range("B23").Value = sizeResult.InitialStop
        Range("B24").Value = sizeResult.Shares
        Range("B25").Value = sizeResult.Contracts
        Range("B26").Value = sizeResult.ActualRisk
        Range("B29").Value = "✅ Success (corr_id: " & corrID & ")"
    Else
        Range("B29").Value = "❌ " & cmdResult.ErrorOutput
    End If
End Sub
```

---

### 4. Checklist Sheet

**Purpose:** Evaluate 6-item checklist and get banner

**Layout:**

```
Row 1: [Header] Checklist Evaluation

Row 3: Ticker:  [B3]

Row 5: 6-Item Checklist (checkboxes)
Row 6:  ☐ Higher high (price action)
Row 7:  ☐ Wider range (volatility expansion)
Row 8:  ☐ Close off low (strength)
Row 9:  ☐ Liquidity OK (volume adequate)
Row 10: ☐ Not overbought (technical check)
Row 11: ☐ Bucket OK (sector not overexposed)

Row 13: [Button: Evaluate]

Row 15: Results
Row 16: Banner:              [B16: Large cell, color-coded background]
Row 17: Missing Count:       [B17]
Row 18: Missing Items:       [B18: Comma-separated list]
Row 19: Allow Save:          [B19: TRUE/FALSE]
Row 20: Evaluation Time:     [B20: Timestamp]

Row 22: Impulse Timer
Row 23: Time Remaining:      [B23: MM:SS countdown if active]
Row 24: Brake Cleared:       [B24: TRUE/FALSE]

Row 26: Status
Row 27: [B27: Status message]
```

**Banner Formatting:**
- GREEN: Background RGB(198, 239, 206), Font: Bold, 16pt
- YELLOW: Background RGB(255, 235, 156), Font: Bold, 16pt
- RED: Background RGB(255, 199, 206), Font: Bold, 16pt

**Button Code (Evaluate):**
```vba
Sub Button_Evaluate_Click()
    Dim cmdResult As TFCommandResult
    Dim checkResult As TFChecklistResult
    Dim checks As New Collection
    Dim corrID As String

    ' Collect checkbox values (assume ActiveX checkboxes named Check1-Check6)
    checks.Add Sheet4.Check1.Value
    checks.Add Sheet4.Check2.Value
    checks.Add Sheet4.Check3.Value
    checks.Add Sheet4.Check4.Value
    checks.Add Sheet4.Check5.Value
    checks.Add Sheet4.Check6.Value

    corrID = TFHelpers.GenerateCorrelationID()

    cmdResult = TFEngine.Engine_Checklist(Range("B3").Value, checks, corrID)

    If cmdResult.Success Then
        checkResult = TFHelpers.ParseChecklistJSON(cmdResult.JsonOutput)

        Range("B16").Value = checkResult.Banner
        Range("B17").Value = checkResult.MissingCount
        Range("B18").Value = checkResult.MissingItems
        Range("B19").Value = checkResult.AllowSave
        Range("B20").Value = TFHelpers.FormatTimestamp(checkResult.EvaluationTimestamp)

        ' Color-code banner
        Select Case checkResult.Banner
            Case "GREEN"
                Range("B16").Interior.Color = RGB(198, 239, 206)
            Case "YELLOW"
                Range("B16").Interior.Color = RGB(255, 235, 156)
            Case "RED"
                Range("B16").Interior.Color = RGB(255, 199, 206)
        End Select

        Range("B27").Value = "✅ Evaluated (corr_id: " & corrID & ")"
    Else
        Range("B27").Value = "❌ " & cmdResult.ErrorOutput
    End If
End Sub
```

---

### 5. Heat Check Sheet

**Purpose:** Verify portfolio and bucket heat before trade

**Layout:**

```
Row 1: [Header] Heat Management

Row 3: New Trade
Row 4: Risk Dollars (R):  [B4: Currency]
Row 5: Bucket:            [B5: Dropdown of sector buckets]

Row 7: [Button: Check Heat]

Row 9: Portfolio Heat
Row 10: Current Heat:     [B10: Currency]
Row 11: New Heat:         [B11: Currency]
Row 12: Heat %:           [B12: Percentage of cap]
Row 13: Cap:              [B13: Currency]
Row 14: Exceeded:         [B14: TRUE/FALSE]
Row 15: Overage:          [B15: Currency, show if exceeded]

Row 17: Bucket Heat
Row 18: Current Heat:     [B18: Currency]
Row 19: New Heat:         [B19: Currency]
Row 20: Heat %:           [B20: Percentage of cap]
Row 21: Cap:              [B21: Currency]
Row 22: Exceeded:         [B22: TRUE/FALSE]
Row 23: Overage:          [B23: Currency, show if exceeded]

Row 25: Result
Row 26: Allowed:          [B26: Large cell, green if OK, red if exceeded]

Row 28: Status
Row 29: [B29: Status message]
```

---

### 6. Trade Entry Sheet

**Purpose:** Complete trade workflow with 5 hard gates

**Layout:**

```
Row 1: [Header] Trade Entry - Full Workflow

Row 3: Trade Details
Row 4: Ticker:            [B4]
Row 5: Entry Price:       [B5]
Row 6: ATR (N):           [B6]
Row 7: K Multiple:        [B7]
Row 8: Bucket:            [B8: Dropdown]
Row 9: Preset:            [B9: Dropdown, e.g., "TF_BREAKOUT_LONG"]

Row 11: [Button: Calculate Size]   [Button: Evaluate Checklist]   [Button: Check Heat]

Row 13: Calculated Results
Row 14: Risk Dollars:     [B14]
Row 15: Shares:           [B15]
Row 16: Initial Stop:     [B16]
Row 17: Banner:           [B17: Color-coded]

Row 19: Gate Status (auto-updated)
Row 20: 1. Banner GREEN:            [B20: ✅ or ❌]
Row 21: 2. Ticker in Candidates:    [B21: ✅ or ❌]
Row 22: 3. Impulse Brake Cleared:   [B22: ✅ or ❌ with timer]
Row 23: 4. Bucket Not on Cooldown:  [B23: ✅ or ❌]
Row 24: 5. Heat Caps OK:            [B24: ✅ or ❌]

Row 26: [Button: Save Decision] (enabled only if all 5 gates pass)

Row 28: Status
Row 29: [B29: Status message with correlation ID]
```

**Save Decision Flow:**
1. User fills trade details
2. Clicks Calculate Size → Populates B14-B16
3. Checks checklist items, clicks Evaluate → Updates B17, gate 1
4. System auto-checks gates 2-5 when Save clicked
5. If all gates pass → Decision saved, form clears
6. If any gate fails → Error shown, form retained

---

### 7. Candidates Sheet

**Purpose:** Manage daily candidate tickers

**Layout:**

```
Row 1: [Header] Candidate Tickers

Row 3: Import Candidates
Row 4: Tickers (comma-separated): [B4: Text box for "AAPL,MSFT,NVDA"]
Row 5: Preset:                    [B5: Dropdown, e.g., "TF_BREAKOUT_LONG"]
Row 6: [Button: Import]

Row 8: Or Open FINVIZ
Row 9: [Button: Open FINVIZ Screener]  (opens browser to FINVIZ)
Row 10: After screening, copy tickers and paste above

Row 12: Today's Candidates
Row 13: Date: [B13: Auto-filled with today]
Row 14: [Button: Refresh List]

Row 16: Candidates Table
Headers: ID | Ticker | Sector | Bucket | Preset | Date

Rows 17+: Populated from database

Row 50: Status
Row 51: [B51: Status message]
```

---

### 8. Positions Sheet

**Purpose:** View open positions and risk

**Layout:**

```
Row 1: [Header] Open Positions

Row 3: [Button: Refresh]   [Button: Show All] (include closed)

Row 5: Summary
Row 6: Open Positions:    [B6: Count]
Row 7: Total Risk:        [B7: Sum of all open R]
Row 8: Portfolio Heat %:  [B8: Percentage of cap]

Row 10: Positions Table
Headers: ID | Ticker | Bucket | Open Date | Units | Entry | Stop | Risk ($) | Status

Rows 11+: Populated from database

Row 50: Status
Row 51: [B51: Status message]
```

---

## Named Ranges Summary

**Setup Sheet:**
- EnginePathSetting: B4
- DatabasePathSetting: B5

**Position Sizing Sheet:**
- (No named ranges - direct cell references)

**Trade Entry Sheet:**
- (No named ranges - direct cell references)

---

## VBA Workbook Module (ThisWorkbook)

```vba
Private Sub Workbook_Open()
    ' Log workbook open
    Dim corrID As String
    corrID = TFHelpers.GenerateCorrelationID()
    TFHelpers.LogMessage corrID, "INFO", "TradingPlatform.xlsm opened"

    ' Check if engine exists
    Dim enginePath As String
    enginePath = ThisWorkbook.Path & "\tf-engine.exe"

    If Dir(enginePath) = "" Then
        MsgBox "WARNING: tf-engine.exe not found!" & vbCrLf & vbCrLf & _
               "Expected location: " & enginePath & vbCrLf & vbCrLf & _
               "Please run windows-init-database.bat to set up the system.", _
               vbExclamation, "Trading Engine Not Found"
    End If
End Sub
```

---

## Dropdown Lists (Data Validation)

### Method Dropdown (Position Sizing Sheet, B8):
- Source: "stock,opt-delta-atr,opt-maxloss"

### Bucket Dropdown (Multiple sheets):
- Source: Named range `BucketList` pointing to:
  - "Tech/Comm"
  - "Energy/Mat"
  - "Fin/REIT"
  - "Cons/Hlth"
  - "Ind/Util"

### Preset Dropdown (Multiple sheets):
- Source: Named range `PresetList` pointing to:
  - "TF_BREAKOUT_LONG"
  - "TF_PULLBACK_LONG"
  - "MANUAL"

---

## Formatting Standards

### Currency Cells:
- Format: `$#,##0.00`
- Examples: B5, B11, B21, B23, etc.

### Percentage Cells:
- Format: `0.00%`
- Examples: B12, B20, B8 (Heat sheet)

### Status Cells:
- Font: Consolas or Courier New, 10pt
- Success: Green background RGB(198, 239, 206), ✅ prefix
- Error: Red background RGB(255, 199, 206), ❌ prefix
- Warning: Yellow background RGB(255, 235, 156), ⚠️ prefix

### Header Rows:
- Font: Bold, 14pt
- Background: Light gray RGB(217, 217, 217)

---

## Initial Setup Checklist

After creating workbook in Windows:

1. ✅ Create 8 worksheets with names as specified
2. ✅ Import VBA modules via windows-import-vba.bat
3. ✅ Set up named ranges (EnginePathSetting, DatabasePathSetting)
4. ✅ Create dropdown lists (Method, Bucket, Preset)
5. ✅ Add buttons and link to VBA functions
6. ✅ Format cells (currency, percentage, colors)
7. ✅ Add ActiveX checkboxes to Checklist sheet (Check1-Check6)
8. ✅ Test connection to engine
9. ✅ Initialize database
10. ✅ Run VBA tests
11. ✅ Save as TradingPlatform.xlsm

---

## Notes

- **Checkbox Controls:** Use ActiveX CheckBox controls (not Form Controls) for programmatic access
- **Button Controls:** Can use either Form Controls or ActiveX CommandButton
- **Status Cells:** Always show correlation ID for log cross-referencing
- **Error Handling:** All buttons should wrap calls in Try/Catch and log errors
- **Workbook Protection:** Consider protecting structure to prevent accidental sheet deletion

---

**This template will be manually created in Windows during M21 setup phase.**

For now (M20), this document serves as the specification for the workbook structure.
