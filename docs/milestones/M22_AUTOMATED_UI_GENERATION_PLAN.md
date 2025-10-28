# M22: Automated UI Generation Plan

**Milestone:** M22 - Automated Trading UI Worksheet Generation
**Goal:** Make workbook "come alive" with full trading UI automatically via batch scripts
**Approach:** Extend `1-setup-all.bat` to create 5 production-ready worksheets with VBA code
**Predecessor:** M21 (Windows Integration Validation)

---

## Overview

Currently, `1-setup-all.bat` creates a minimal workbook with:
- VBA modules (TFTypes, TFHelpers, TFEngine, TFTests, TFIntegrationTests)
- "VBA Tests" worksheet with test runner button
- Named ranges for configuration

**M22 Goal:** Enhance setup to automatically create 5 production worksheets:
1. **Dashboard** - Portfolio overview and quick actions
2. **Position Sizing** - Calculate shares/contracts for entry
3. **Checklist** - Evaluate 6-item checklist, get GREEN/YELLOW/RED banner
4. **Heat Check** - Verify portfolio/bucket heat caps
5. **Trade Entry** - Full trade decision workflow with 5 hard gates

After setup, user should have a **fully functional trading workbook** ready to use immediately.

---

## Architecture Decision

### Approach: VBScript-Based Worksheet Generation

**Current setup-all.bat uses VBScript to:**
- Create Excel Application object
- Create workbook programmatically
- Import VBA modules
- Create named ranges
- Add "Run All Tests" button

**M22 extends this to:**
- Create 5 additional worksheets
- Add labels, headers, and formatting
- Create input cells with data validation
- Add ActiveX controls (buttons, checkboxes, dropdowns)
- Wire up VBA event handlers
- Set up result display areas

**Why VBScript?**
- ✅ Already used successfully in current setup
- ✅ Can create complex Excel objects programmatically
- ✅ No additional dependencies (built into Windows)
- ✅ Can be embedded in batch files via HEREDOC
- ✅ Runs without Excel UI (background automation)

**Alternatives considered:**
- PowerShell: Requires execution policy changes
- Python/openpyxl: Requires Python installation
- Manual template: Requires version control of binary .xlsm file

---

## Implementation Strategy

### Phase 1: Create Worksheet Creation Library (VBScript)

Create reusable VBScript functions for common operations:

```vbscript
' CreateWorksheet(wb, sheetName, tabColor)
Function CreateWorksheet(wb, sheetName, tabColor)
    Dim ws
    Set ws = wb.Worksheets.Add
    ws.Name = sheetName
    ws.Tab.Color = tabColor ' RGB color
    Set CreateWorksheet = ws
End Function

' AddLabel(ws, cell, text, bold, size, color)
Sub AddLabel(ws, cell, text, bold, size, color)
    ws.Range(cell).Value = text
    ws.Range(cell).Font.Bold = bold
    ws.Range(cell).Font.Size = size
    ws.Range(cell).Font.Color = color
End Sub

' AddButton(ws, cell, caption, macroName)
Sub AddButton(ws, cell, caption, macroName)
    Dim btn
    Set btn = ws.OLEObjects.Add(ClassType:="Forms.CommandButton.1", _
        Link:=False, DisplayAsIcon:=False, _
        Left:=ws.Range(cell).Left, _
        Top:=ws.Range(cell).Top, _
        Width:=120, Height:=30)
    btn.Object.Caption = caption
    btn.Object.ForeColor = RGB(255, 255, 255)
    btn.Object.BackColor = RGB(0, 102, 204)
    btn.Object.Font.Bold = True
    ws.OLEObjects(btn.Name).OnAction = macroName
End Sub

' AddDropdown(ws, cell, items)
Sub AddDropdown(ws, cell, items)
    Dim validation
    Set validation = ws.Range(cell).Validation
    validation.Delete
    validation.Add Type:=3, AlertStyle:=1, Operator:=1, Formula1:=items
    validation.IgnoreBlank = True
    validation.InCellDropdown = True
End Sub

' AddCheckbox(ws, position, caption)
Function AddCheckbox(ws, position, caption)
    Dim chk
    Set chk = ws.OLEObjects.Add(ClassType:="Forms.CheckBox.1", _
        Link:=False, DisplayAsIcon:=False, _
        Left:=position(0), Top:=position(1), _
        Width:=200, Height:=20)
    chk.Object.Caption = caption
    Set AddCheckbox = chk
End Function

' FormatResultArea(ws, startCell, endCell, bgColor)
Sub FormatResultArea(ws, startCell, endCell, bgColor)
    Dim rng
    Set rng = ws.Range(startCell & ":" & endCell)
    rng.Interior.Color = bgColor
    rng.Font.Bold = True
    rng.Borders.LineStyle = 1
End Sub
```

### Phase 2: Define Worksheet Specifications

Create declarative specifications for each worksheet (JSON-like structure in VBScript):

#### Worksheet 1: Dashboard

```vbscript
' Dashboard worksheet
sheetName = "Dashboard"
tabColor = RGB(0, 102, 204) ' Blue

' Layout
headers = Array( _
    Array("A1", "Trading Platform Dashboard", True, 16, RGB(0, 0, 139)), _
    Array("A3", "Portfolio Status", True, 12, RGB(0, 0, 0)), _
    Array("A10", "Today's Candidates", True, 12, RGB(0, 0, 0)), _
    Array("A16", "Quick Actions", True, 12, RGB(0, 0, 0)) _
)

labels = Array( _
    Array("A4", "Current Equity:", False, 10), _
    Array("A5", "Portfolio Heat:", False, 10), _
    Array("A6", "Portfolio Cap:", False, 10), _
    Array("A7", "Heat %:", False, 10) _
)

resultCells = Array( _
    Array("B4", "=GetSetting('Equity_E')"), _
    Array("B5", "=GetPortfolioHeat()"), _
    Array("B6", "=GetPortfolioCap()"), _
    Array("B7", "=B5/B6") _
)

buttons = Array( _
    Array("B17", "Refresh", "RefreshDashboard"), _
    Array("B18", "Position Sizing", "GotoPositionSizing"), _
    Array("B19", "Checklist", "GotoChecklist"), _
    Array("B20", "Heat Check", "GotoHeatCheck"), _
    Array("B21", "Trade Entry", "GotoTradeEntry") _
)
```

#### Worksheet 2: Position Sizing

```vbscript
sheetName = "Position Sizing"
tabColor = RGB(0, 153, 76) ' Green

headers = Array( _
    Array("A1", "Position Sizing Calculator", True, 14, RGB(0, 102, 0)), _
    Array("A3", "Inputs", True, 12, RGB(0, 0, 0)), _
    Array("A15", "Results", True, 12, RGB(0, 0, 0)) _
)

inputs = Array( _
    Array("A4", "Ticker:", "B4", "TEXT"), _
    Array("A5", "Entry Price:", "B5", "NUMBER"), _
    Array("A6", "ATR (N):", "B6", "NUMBER"), _
    Array("A7", "K Multiple:", "B7", "NUMBER"), _
    Array("A8", "Method:", "B8", "DROPDOWN", "stock,opt-delta-atr,opt-maxloss") _
)

optionalInputs = Array( _
    Array("A9", "Equity Override:", "B9", "NUMBER"), _
    Array("A10", "Risk % Override:", "B10", "PERCENT"), _
    Array("A11", "Delta (options):", "B11", "NUMBER"), _
    Array("A12", "Max Loss (options):", "B12", "CURRENCY") _
)

results = Array( _
    Array("A16", "Risk Dollars (R):", "B16", "CURRENCY"), _
    Array("A17", "Stop Distance:", "B17", "NUMBER"), _
    Array("A18", "Initial Stop:", "B18", "CURRENCY"), _
    Array("A19", "Shares:", "B19", "NUMBER"), _
    Array("A20", "Contracts:", "B20", "NUMBER"), _
    Array("A21", "Actual Risk:", "B21", "CURRENCY"), _
    Array("A22", "Status:", "B22", "TEXT") _
)

buttons = Array( _
    Array("B13", "Calculate", "CalculatePositionSize"), _
    Array("C13", "Clear", "ClearPositionSizing") _
)
```

#### Worksheet 3: Checklist

```vbscript
sheetName = "Checklist"
tabColor = RGB(255, 192, 0) ' Orange

headers = Array( _
    Array("A1", "Entry Checklist Evaluation", True, 14, RGB(204, 102, 0)) _
)

inputs = Array( _
    Array("A3", "Ticker:", "B3", "TEXT") _
)

' 6 checkboxes for checklist items
checkboxes = Array( _
    Array("A5", "from-preset", "Ticker from today's FINVIZ preset"), _
    Array("A6", "trend-pass", "Trend alignment confirmed"), _
    Array("A7", "liquidity-pass", "Adequate volume and spread"), _
    Array("A8", "tv-confirm", "TradingView setup confirmation"), _
    Array("A9", "earnings-ok", "No earnings in next 7 days"), _
    Array("A10", "journal-ok", "Trade thesis documented in journal") _
)

results = Array( _
    Array("A15", "Banner:", "B15", "BANNER_COLOR"), _
    Array("A16", "Missing Items:", "B16", "NUMBER"), _
    Array("A17", "Missing:", "B17:B20", "LIST"), _
    Array("A21", "Allow Save:", "B21", "BOOLEAN"), _
    Array("A22", "Evaluation Time:", "B22", "DATETIME"), _
    Array("A23", "Status:", "B23", "TEXT") _
)

buttons = Array( _
    Array("B12", "Evaluate", "EvaluateChecklist"), _
    Array("C12", "Clear", "ClearChecklist") _
)
```

#### Worksheet 4: Heat Check

```vbscript
sheetName = "Heat Check"
tabColor = RGB(255, 0, 0) ' Red

headers = Array( _
    Array("A1", "Portfolio Heat Management", True, 14, RGB(139, 0, 0)), _
    Array("A9", "Portfolio Heat", True, 12, RGB(0, 0, 0)), _
    Array("A17", "Bucket Heat", True, 12, RGB(0, 0, 0)) _
)

inputs = Array( _
    Array("A3", "Ticker:", "B3", "TEXT"), _
    Array("A4", "Risk Amount ($):", "B4", "CURRENCY"), _
    Array("A5", "Bucket:", "B5", "DROPDOWN", "Tech/Comm,Finance,Healthcare,Consumer,Energy,Industrials") _
)

portfolioResults = Array( _
    Array("A10", "Current Heat:", "B10", "CURRENCY"), _
    Array("A11", "New Heat:", "B11", "CURRENCY"), _
    Array("A12", "Heat %:", "B12", "PERCENT"), _
    Array("A13", "Cap:", "B13", "CURRENCY"), _
    Array("A14", "Exceeded:", "B14", "BOOLEAN"), _
    Array("A15", "Overage:", "B15", "CURRENCY") _
)

bucketResults = Array( _
    Array("A18", "Current Heat:", "B18", "CURRENCY"), _
    Array("A19", "New Heat:", "B19", "CURRENCY"), _
    Array("A20", "Heat %:", "B20", "PERCENT"), _
    Array("A21", "Cap:", "B21", "CURRENCY"), _
    Array("A22", "Exceeded:", "B22", "BOOLEAN"), _
    Array("A23", "Overage:", "B23", "CURRENCY") _
)

buttons = Array( _
    Array("B7", "Check Heat", "CheckHeat"), _
    Array("C7", "Clear", "ClearHeatCheck") _
)
```

#### Worksheet 5: Trade Entry

```vbscript
sheetName = "Trade Entry"
tabColor = RGB(128, 0, 128) ' Purple

headers = Array( _
    Array("A1", "Trade Decision Entry (5 Hard Gates)", True, 14, RGB(75, 0, 130)), _
    Array("A3", "Trade Details", True, 12, RGB(0, 0, 0)), _
    Array("A17", "Gate Status", True, 12, RGB(0, 0, 0)), _
    Array("A24", "Results", True, 12, RGB(0, 0, 0)) _
)

inputs = Array( _
    Array("A4", "Ticker:", "B4", "TEXT"), _
    Array("A5", "Entry Price:", "B5", "CURRENCY"), _
    Array("A6", "ATR:", "B6", "NUMBER"), _
    Array("A7", "Method:", "B7", "DROPDOWN", "stock,opt-delta-atr,opt-maxloss"), _
    Array("A8", "Banner Status:", "B8", "DROPDOWN", "GREEN,YELLOW,RED"), _
    Array("A9", "Delta (options):", "B9", "NUMBER"), _
    Array("A10", "Max Loss (options):", "B10", "CURRENCY"), _
    Array("A11", "Bucket:", "B11", "DROPDOWN", "Tech/Comm,Finance,Healthcare,Consumer,Energy,Industrials"), _
    Array("A12", "Preset:", "B12", "TEXT") _
)

gateStatus = Array( _
    Array("A18", "Gate 1 - Banner GREEN:", "B18", "GATE_STATUS"), _
    Array("A19", "Gate 2 - In Candidates:", "B19", "GATE_STATUS"), _
    Array("A20", "Gate 3 - Impulse Brake:", "B20", "GATE_STATUS"), _
    Array("A21", "Gate 4 - Cooldown:", "B21", "GATE_STATUS"), _
    Array("A22", "Gate 5 - Heat Caps:", "B22", "GATE_STATUS") _
)

results = Array( _
    Array("A25", "Decision Saved:", "B25", "BOOLEAN"), _
    Array("A26", "Decision ID:", "B26", "NUMBER"), _
    Array("A27", "Rejection Reason:", "B27:C29", "TEXT_AREA"), _
    Array("A30", "Status:", "B30", "TEXT") _
)

buttons = Array( _
    Array("B14", "Save Decision (GO)", "SaveDecisionGO"), _
    Array("C14", "Save Decision (NO-GO)", "SaveDecisionNOGO"), _
    Array("D14", "Clear", "ClearTradeEntry") _
)
```

### Phase 3: Create VBA Support Functions

Add new VBA functions to support UI operations:

**In TFHelpers.bas:**

```vba
' Refresh dashboard data from database
Public Sub RefreshDashboard()
    ' Query current portfolio state
    ' Update dashboard cells
End Sub

' Navigation helpers
Public Sub GotoPositionSizing()
    Worksheets("Position Sizing").Activate
End Sub

Public Sub GotoChecklist()
    Worksheets("Checklist").Activate
End Sub

' ... etc
```

**In TFEngine.bas:**

```vba
' Wrapper functions for each worksheet button action
Public Sub CalculatePositionSize()
    ' Read inputs from Position Sizing sheet
    ' Call ExecuteCommand("size ...")
    ' Parse results
    ' Update result cells
    ' Format cells based on success/failure
End Sub

Public Sub EvaluateChecklist()
    ' Read checkbox states
    ' Build checklist command
    ' Execute and display results
    ' Color-code banner cell
End Sub

Public Sub CheckHeat()
    ' Read risk amount and bucket
    ' Execute heat command
    ' Display portfolio and bucket results
    ' Highlight exceeded caps in red
End Sub

Public Sub SaveDecisionGO()
    ' Read all trade entry inputs
    ' Execute save-decision command with action=GO
    ' Display gate results
    ' Show acceptance or rejection
    ' Clear form on success
End Sub

Public Sub SaveDecisionNOGO()
    ' Execute save-decision with action=NO-GO
    ' Prompt for reason
    ' Save to database
End Sub
```

### Phase 4: Enhance setup-all.bat

Modify `1-setup-all.bat` to include worksheet generation:

```batch
REM After VBA module import...
REM Create UI worksheets

echo [6/8] Creating UI worksheets...

REM Generate VBScript for worksheet creation
call :CreateWorksheetGeneration

REM Execute worksheet generation
cscript //nologo create-worksheets.vbs
if errorlevel 1 (
    echo ERROR: Failed to create worksheets
    pause
    exit /b 1
)

del create-worksheets.vbs
echo   - 5 UI worksheets created

REM Continue with existing named range setup...
```

### Phase 5: Testing Strategy

After worksheet generation:

1. **Visual inspection:** Open workbook, verify all 5 worksheets present and formatted
2. **Button tests:** Click each button, verify VBA code executes
3. **Input validation:** Test dropdowns, checkboxes, number fields
4. **Integration tests:** Run a complete trade workflow end-to-end
5. **Regression tests:** Ensure VBA Tests sheet still works

---

## Implementation Phases

### Phase 1: Foundation (1-2 hours)
- [ ] Create VBScript helper function library
- [ ] Test helper functions in isolation
- [ ] Verify worksheet creation, formatting, controls work

### Phase 2: Dashboard + Position Sizing (2-3 hours)
- [ ] Implement Dashboard worksheet generation
- [ ] Implement Position Sizing worksheet generation
- [ ] Add corresponding VBA functions (CalculatePositionSize, etc.)
- [ ] Test both worksheets manually
- [ ] Document any issues/limitations

### Phase 3: Checklist + Heat Check (2-3 hours)
- [ ] Implement Checklist worksheet generation
- [ ] Implement Heat Check worksheet generation
- [ ] Add VBA functions (EvaluateChecklist, CheckHeat)
- [ ] Test both worksheets manually
- [ ] Verify checkbox behavior

### Phase 4: Trade Entry (3-4 hours)
- [ ] Implement Trade Entry worksheet generation
- [ ] Add VBA functions (SaveDecisionGO, SaveDecisionNOGO)
- [ ] Test full 5-gate workflow
- [ ] Test form clear on success/failure

### Phase 5: Integration and Polish (2-3 hours)
- [ ] Integrate all worksheets into setup-all.bat
- [ ] Test complete setup from scratch
- [ ] Add navigation buttons between sheets
- [ ] Add conditional formatting for results
- [ ] Create user guide documentation

**Total estimated time:** 10-15 hours

---

## File Changes Required

### New Files:
- `windows/vbscript-lib.vbs` - Reusable VBScript functions
- `docs/USER_GUIDE.md` - How to use the trading workbook

### Modified Files:
- `windows/1-setup-all.bat` - Add worksheet generation step
- `excel/vba/TFEngine.bas` - Add UI button handler functions
- `excel/vba/TFHelpers.bas` - Add dashboard and navigation functions

### Generated Files (by setup):
- `TradingPlatform.xlsm` - Now includes 5 production worksheets

---

## Benefits

### For Users:
- ✅ **One-click setup** - Full trading workbook ready in 3 minutes
- ✅ **No manual worksheet creation** - Everything automated
- ✅ **Production-ready UI** - Professional appearance, fully functional
- ✅ **Consistent layout** - Same structure every time

### For Development:
- ✅ **Version controlled UI** - Worksheet structure in VBScript (text files)
- ✅ **Reproducible builds** - Same workbook every setup
- ✅ **Easy updates** - Modify VBScript, re-run setup
- ✅ **No binary files** - .xlsm generated, not committed

### For Testing:
- ✅ **Automated UI testing** - Can generate test workbooks on demand
- ✅ **Clean slate** - Fresh workbook for each test run
- ✅ **CI-friendly** - Could run in CI/CD pipeline (with Excel installed)

---

## Risks and Mitigations

### Risk 1: VBScript complexity
**Issue:** Generating complex Excel objects in VBScript can be verbose
**Mitigation:** Create helper function library, use declarative specifications

### Risk 2: Excel version differences
**Issue:** ActiveX controls may behave differently across Excel versions
**Mitigation:** Test on Excel 2016, 2019, 365; document known issues

### Risk 3: Setup time increase
**Issue:** Worksheet generation may slow down setup-all.bat
**Mitigation:** Profile execution time, optimize if >5 minutes total

### Risk 4: VBA code duplication
**Issue:** Button handlers in TFEngine.bas may be repetitive
**Mitigation:** Extract common patterns to helper functions

---

## Success Criteria

M22 is complete when:

1. ✅ `1-setup-all.bat` creates workbook with 6 worksheets:
   - VBA Tests (existing)
   - Dashboard (new)
   - Position Sizing (new)
   - Checklist (new)
   - Heat Check (new)
   - Trade Entry (new)

2. ✅ All worksheets have:
   - Professional formatting
   - Working buttons
   - Input validation
   - Result display areas
   - Navigation

3. ✅ User can complete full trade workflow:
   - Calculate position size
   - Evaluate checklist
   - Check heat
   - Save decision
   - All without writing VBA code

4. ✅ Setup completes in <5 minutes

5. ✅ All Phase 3 VBA tests still pass

6. ✅ Documentation updated with screenshots

---

## Next Steps

1. Review this plan with user
2. Get approval to proceed
3. Start with Phase 1 (foundation)
4. Iterate through phases 2-5
5. Demo complete workbook
6. Document and close M22

---

**Created:** 2025-10-28
**Status:** Planning
**Estimated Duration:** 10-15 hours development
**Target Completion:** M22
