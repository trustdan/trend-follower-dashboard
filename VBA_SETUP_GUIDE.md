# VBA Setup Guide for Interactive Trade Entry Workbook

This guide walks you through importing the VBA modules and building the TradeEntry UI.

---

## Part 1: Import VBA Modules (10 minutes)

### **Method A: Automated Import (Windows Only, Recommended)**

**Requirements:**
- Windows OS
- Python 3.7+ installed
- `pywin32` package: `pip install pywin32`

**Steps:**
1. Double-click `IMPORT_VBA_MODULES.bat` in Windows Explorer
2. OR run from command line: `python import_to_excel.py`
3. Script will:
   - Open Excel
   - Create new workbook `TrendFollowing_TradeEntry.xlsm`
   - Import all 4 standard modules
   - Import class modules
   - Leave Excel open for you
4. Skip to Step 6 below (Run Initial Setup)

**If automated import fails**, see Method B below.

---

### **Method B: Manual Import (All Platforms)**

### Step 1: Create New Workbook
1. Open Excel
2. Create a new blank workbook
3. Save as `TrendFollowing_TradeEntry.xlsm` (macro-enabled format)
4. Close the workbook and reopen it to confirm macros are enabled

### Step 2: Enable Developer Tab
1. File → Options → Customize Ribbon
2. Check the box for "Developer" in the right column
3. Click OK

### Step 3: Open VBA Editor
1. Click Developer tab → Visual Basic (or press Alt+F11)
2. You should see the VBA Editor window

### Step 4: Import Standard Modules
1. In VBA Editor, select your workbook in the Project Explorer (left pane)
2. File → Import File
3. Navigate to `VBA/` folder (in same directory as this guide)
4. Import these files **in order**:
   - `TF_Utils.bas`
   - `TF_Data.bas`
   - `TF_UI.bas`
   - `TF_Presets.bas`

5. You should now see 4 modules in the Project Explorer under "Modules"

### Step 5: Import Class Modules (Event Handlers)
1. In VBA Editor, File → Import File
2. Import `ThisWorkbook.cls`
   - This will REPLACE your existing ThisWorkbook module (that's OK)

3. **SKIP** importing `Sheet_TradeEntry.cls` for now
   - We'll handle this after creating the TradeEntry sheet

---

### Step 6: Run Initial Setup (Both Methods)
1. In VBA Editor, View → Immediate Window (or Ctrl+G)
2. Type this command and press Enter:
   ```vba
   EnsureStructure
   ```
3. You should see a message box: "Workbook structure created successfully!"
4. Return to Excel (Alt+F11 or click Excel icon in taskbar)
5. Verify you now have 8 sheets:
   - TradeEntry (blank for now)
   - Presets (has 5 preset rows)
   - Buckets (has 6 bucket rows)
   - Candidates (empty table)
   - Decisions (empty table)
   - Positions (empty table)
   - Summary (has settings with default values)
   - Control (hidden)

---

## Part 2: Build TradeEntry UI (30 minutes)

### Step 7: Add Labels and Format Cells

Switch to the **TradeEntry** sheet and build the layout:

#### Row 1: Title
- Merge cells A1:F1
- Enter: "TRADE ENTRY SYSTEM"
- Font: Bold, 16pt, centered

#### Row 2-3: Banner (will be filled by VBA)
- Leave A2:F2 and A3:F3 blank (merged and colored by Evaluate button)

#### Rows 5-18: Input Section
Add these labels in column A:

| Cell | Label |
|------|-------|
| A5 | Preset: |
| A6 | Ticker: |
| A7 | Sector: |
| A8 | Bucket: |
| A9 | Entry Price: |
| A10 | ATR N: |
| A11 | K (Stop Multiple): |
| A13 | Method: |
| A16 | Delta: |
| A17 | DTE: |
| A18 | Max Loss/Contract: |

#### Rows 20-25: Checklist Labels
| Cell | Label |
|------|-------|
| A20 | From Preset |
| A21 | Trend Pass |
| A22 | Liquidity Pass |
| A23 | TV Confirm |
| A24 | Earnings OK |
| A25 | Journal OK |

#### Rows 5-12: Output Section (Column E-F)
| Cell | Label | Cell | Formula/Value |
|------|-------|------|---------------|
| E5 | R ($): | F5 | =(leave blank, filled by VBA) |
| E6 | Stop Distance: | F6 | (leave blank) |
| E7 | Initial Stop: | F7 | (leave blank) |
| E8 | Shares: | F8 | (leave blank) |
| E9 | Contracts: | F9 | (leave blank) |
| E10 | Add 1: | F10 | =B9+(AddStepN*B10) |
| E11 | Add 2: | F11 | =B9+(2*AddStepN*B10) |
| E12 | Add 3: | F12 | =B9+(3*AddStepN*B10) |

#### Rows 14-15: Heat Preview
| Cell | Label |
|------|-------|
| E14 | Portfolio Heat: |
| E15 | Bucket Heat: |

Format F14 and F15 with number format and green fill (will be updated by VBA).

### Step 8: Hide Column C
- Right-click column C header → Hide
- This column stores hidden values for option buttons and checkboxes

### Step 9: Add Form Controls

#### Add Option Buttons (Method Choice)
1. Developer tab → Insert → Option Button (Form Control)
2. Draw 3 option buttons in cells B13:B15:
   - B13: "Stock"
   - B14: "Opt-DeltaATR"
   - B15: "Opt-MaxLoss"
3. Right-click first option button → Format Control
4. Set "Cell link" to: `$C$13`
5. Click OK
6. Repeat for other 2 buttons (they will auto-group and share the same cell link)

#### Add Checkboxes (Checklist Items)
1. Developer tab → Insert → Check Box (Form Control)
2. Draw 6 checkboxes in column B, rows 20-25
3. For each checkbox:
   - Right-click → Edit Text → delete default text (label is in column A)
   - Right-click → Format Control → Cell link:
     - Checkbox in B20 → link to C20
     - Checkbox in B21 → link to C21
     - Checkbox in B22 → link to C22
     - Checkbox in B23 → link to C23
     - Checkbox in B24 → link to C24
     - Checkbox in B25 → link to C25

### Step 10: Add Command Buttons

#### Create 6 Buttons
1. Developer tab → Insert → Button (Form Control)
2. Draw buttons in these locations:

| Cell | Button Text | Assign to Macro |
|------|-------------|-----------------|
| A28 | Open Preset | OpenPreset |
| B28 | Import Candidates | ImportCandidatesPrompt |
| A29 | Evaluate | EvaluateChecklist |
| B29 | Recalc Sizing | RecalcSizing |
| A30 | Save Decision | SaveDecision |
| B30 | Start 2-min Timer | StartImpulseTimer |

3. For each button:
   - Right-click → Assign Macro
   - Select the macro name from the list
   - Click OK

### Step 11: Set Data Validation (Alternative Method)

Instead of running BindControls, you can manually set dropdowns:

1. **Cell B5 (Preset):**
   - Select B5 → Data tab → Data Validation
   - Allow: List
   - Source: `=tblPresets[Name]`

2. **Cell B6 (Ticker):**
   - Select B6 → Data tab → Data Validation
   - Allow: List
   - Source: `=tblCandidates[Ticker]`

3. **Cell B7 (Sector):**
   - Select B7 → Data tab → Data Validation
   - Allow: List
   - Source: `Technology,Healthcare,Financials,Consumer,Industrials,Energy`

4. **Cell B8 (Bucket):**
   - Select B8 → Data tab → Data Validation
   - Allow: List
   - Source: `=tblBuckets[Bucket]`

### Step 12: Add Number Formatting

Select these cells and apply formats:

| Cell(s) | Format |
|---------|--------|
| B9 (Entry) | Number, 2 decimals |
| B10 (N) | Number, 3 decimals |
| B11 (K) | Number, 1 decimal |
| B16 (Delta) | Number, 2 decimals |
| B17 (DTE) | Number, 0 decimals |
| B18 (MaxLoss) | Currency, 2 decimals |
| F5 (R) | Currency, 2 decimals |
| F6-F7 | Number, 2 decimals |
| F8-F9 | Number, 0 decimals |
| F10-F12 | Number, 2 decimals |

### Step 13: Import TradeEntry Sheet Events

Now that the sheet exists:

1. In Excel, right-click the **TradeEntry** sheet tab → View Code
2. This opens the VBA Editor for that specific sheet
3. Delete any existing code in that window
4. Open the file `/home/kali/excel-trading-workflow/VBA/Sheet_TradeEntry.cls` in a text editor
5. Copy all code AFTER the line `Attribute VB_Exposed = True`
6. Paste into the TradeEntry sheet code window in VBA Editor
7. Close VBA Editor (return to Excel)

### Step 14: Test the Setup

1. **Test Structure:** Click Developer → Visual Basic → Immediate Window → type:
   ```vba
   ? SheetExists("TradeEntry")
   ```
   Should return: `True`

2. **Test Dropdown:** Click cell B5 → dropdown arrow should show 5 presets

3. **Test Checkboxes:** Click a checkbox in B20 → cell C20 should show TRUE

4. **Test Buttons:**
   - Click "Evaluate" button → should show "RED - NO-GO" banner
   - Tick all 6 checkboxes → click "Evaluate" → should show "GREEN - GO"
   - Click "Recalc Sizing" button (without inputs) → should show error message

---

## Part 3: Verify and Customize (10 minutes)

### Step 15: Check Summary Settings

1. Go to **Summary** sheet
2. Verify named ranges have default values:

| Cell | Name | Default Value |
|------|------|---------------|
| B2 | Equity_E | 10,000 |
| B3 | RiskPct_r | 0.0075 (0.75%) |
| B4 | StopMultiple_K | 2 |
| B5 | HeatCap_H_pct | 0.04 (4%) |
| B6 | BucketHeatCap_pct | 0.015 (1.5%) |
| B7 | AddStepN | 0.5 |
| B8 | EarningsBufferDays | 3 |

3. **Customize these values** to match your account size and risk tolerance

### Step 16: Seed Test Data

1. Go to **Candidates** sheet
2. Manually add 2-3 test rows:

| Date | Ticker | Preset | Sector | Bucket |
|------|--------|--------|--------|--------|
| (today's date) | MSFT | TF_BREAKOUT_LONG | Technology | Tech/Comm |
| (today's date) | AAPL | TF_BREAKOUT_LONG | Technology | Tech/Comm |

3. Return to **TradeEntry** sheet
4. Test dropdown in B6 → should now show MSFT and AAPL

### Step 17: Full Workflow Test

Run a complete test trade:

1. **Preset:** Select "TF_BREAKOUT_LONG" (B5)
2. **Ticker:** Select "MSFT" (B6)
3. **Sector:** Select "Technology" (B7) → Bucket should auto-fill "Tech/Comm"
4. **Entry:** Enter 420.00 (B9)
5. **N:** Enter 1.20 (B10)
6. **K:** Enter 2 (B11)
7. **Method:** Click "Opt-DeltaATR" option button
   - Rows 16-17 should appear, row 18 should hide
8. **Delta:** Enter 0.30 (B16)
9. **DTE:** Enter 45 (B17)
10. **Checklist:** Tick all 6 checkboxes
11. Click **Evaluate** button → should show GREEN banner and start timer
12. Click **Recalc Sizing** → should show R=$75, Contracts=1
13. Wait 2 minutes (or manually edit Control!A1 to earlier time for testing)
14. Click **Save Decision** → should save successfully

15. **Verify:**
    - **Decisions** sheet has 1 new row
    - **Positions** sheet has 1 new row (MSFT, Open)
    - Banner is cleared on TradeEntry

---

## Part 4: Troubleshooting

### Common Issues

**Issue:** "Compile Error: Sub or Function not defined"
- **Fix:** Make sure all 4 .bas modules are imported

**Issue:** Buttons don't work
- **Fix:** Right-click button → Assign Macro → select correct macro name

**Issue:** Dropdowns show #REF!
- **Fix:** Run `EnsureStructure` again to recreate tables

**Issue:** Banner doesn't change color
- **Fix:** Check that TradeEntry sheet events are imported correctly

**Issue:** Method fields don't hide/show
- **Fix:** Verify C13 has a value (1, 2, or 3) and events are working

**Issue:** Save Decision always blocks
- **Fix:** Check these in order:
  1. Banner is GREEN (click Evaluate first)
  2. Ticker is in Candidates table with today's date
  3. 2 minutes have elapsed (check Control!A1 timestamp)
  4. Bucket not in cooldown (check Buckets sheet)
  5. Heat caps not exceeded (check Summary settings)

### Debug Mode

To see detailed error messages:
1. VBA Editor → Tools → Options → General tab
2. Error Trapping: Check "Break on All Errors"
3. Re-run the failing operation to see exact error location

---

## Part 5: Daily Workflow

Once setup is complete, your daily routine is:

### Morning (Market Open - 10 minutes)
1. Open workbook → goes straight to TradeEntry
2. Click **Open Preset** button → opens FINVIZ
3. Copy tickers from FINVIZ page
4. Click **Import Candidates** → paste tickers
5. Repeat for 2-3 different presets

### During Day (Per Trade - 2-3 minutes)
1. Select ticker from dropdown (B6)
2. Check TradingView to get Entry, N, and confirm signal
3. Fill Entry, N, K, Method, Delta/DTE
4. Tick checklist boxes as you verify each item
5. Click **Evaluate** → wait for GREEN
6. Click **Recalc Sizing** → verify size makes sense
7. Wait 2 minutes (browse TradingView, check chart)
8. Click **Save Decision** → execute trade in broker

### End of Week (Friday after close - 5 minutes)
1. Update any closed positions in Positions sheet (set Status to "Closed")
2. Update Outcome in Decisions sheet for stop-outs
3. Run `UpdateCooldowns` macro (VBA Immediate window):
   ```vba
   UpdateCooldowns
   ```
4. Review adherence to GREEN-only rule

---

## Next Steps

✅ **VBA Setup Complete!**

Optional enhancements:
- **Python Integration:** See Part 6 below (requires Microsoft 365 Insider)
- **Custom Presets:** Add more FINVIZ screeners to Presets sheet
- **Bucket Tuning:** Adjust cooldown parameters in Buckets sheet
- **Risk Settings:** Modify Summary settings based on actual trading results

---

## Part 6: Python Integration (Optional - Advanced)

If you have Microsoft 365 with Python in Excel enabled, you can add:
- Automatic FINVIZ scraping (no manual copy/paste)
- Faster heat calculations
- Real-time earnings calendar checks

See `newest-Interactive_TF_Workbook_Plan.md` section "Python Integration Strategy" for code samples.

---

## Support

For questions or issues:
1. Check the Gherkin test scenarios in the plan document
2. Review the pseudo-code in the module headers
3. Use VBA debugger (F8 to step through code)

**Definition of Done Checklist:**
- [ ] All 4 VBA modules imported and compile without errors
- [ ] TradeEntry UI built with all controls (buttons, checkboxes, dropdowns)
- [ ] Test trade completes full workflow (Import → Evaluate → Size → Save)
- [ ] Decisions and Positions tables update correctly
- [ ] All 6 Gherkin scenarios pass manual testing
- [ ] Summary settings customized to your account
- [ ] Workbook saved as .xlsm with macros enabled

**Estimated Total Setup Time:** 60-90 minutes

Enjoy your bias-free, mechanical trend-following system!
