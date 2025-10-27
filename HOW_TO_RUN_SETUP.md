# How to Run Initial Setup

## Problem: "Dropdown Error - Run Setup.RunInitialSetup"

If you see this error message:
```
Some dropdowns could not be created:
Preset dropdown: Application-defined or object-defined error

This usually means the data tables haven't been created yet.
Run Setup.RunInitialSetup first.
```

## Solution: Click the Button!

### Easy Way (Recommended):

1. **Go to the Setup sheet**
2. **Look for the big button at the top:** `▶ RUN INITIAL SETUP`
3. **Click it**
4. **Wait for completion** (5-10 seconds)
5. **Done!**

The button will:
- Create all data sheets (Presets, Buckets, Candidates, Decisions, Positions, etc.)
- Create all data tables (tblPresets, tblBuckets, etc.)
- Seed default data (5 presets, 6 buckets)
- Build the TradeEntry UI
- Set up all dropdowns

---

## Alternative: Run from VBA Editor

If the button doesn't work, you can run it manually:

### Step 1: Open VBA Editor
- Press **Alt + F11** (or Alt + Fn + F11 on some keyboards)
- Or: Developer tab → Visual Basic

### Step 2: Open Immediate Window
- Press **Ctrl + G**
- Or: View → Immediate Window
- You'll see a window at the bottom labeled "Immediate"

### Step 3: Type and Run Command
In the Immediate Window, type:
```vba
Setup.RunInitialSetup
```

Then press **Enter**

### Step 4: Wait for Completion
- You'll see a progress message
- Wait 5-10 seconds
- A "Setup Complete!" message will appear

### Step 5: Close VBA Editor
- Press **Alt + Q**
- Or: File → Close and Return to Microsoft Excel

---

## What Initial Setup Does

When you run Setup.RunInitialSetup, it:

### 1. Creates Sheets
- ✅ TradeEntry (main UI)
- ✅ Presets (FINVIZ queries)
- ✅ Buckets (sector groupings)
- ✅ Candidates (daily stock lists)
- ✅ Decisions (trade log)
- ✅ Positions (current positions)
- ✅ Summary (settings)
- ✅ Control (hidden - internal use)
- ✅ Setup (instructions)

### 2. Creates Tables
- ✅ tblPresets - 5 default FINVIZ presets
- ✅ tblBuckets - 6 sector buckets
- ✅ tblCandidates - Empty (you import tickers here)
- ✅ tblDecisions - Empty (trade history)
- ✅ tblPositions - Empty (open positions)

### 3. Sets Up UI
- ✅ Builds TradeEntry sheet layout
- ✅ Creates buttons (Evaluate, Recalc, Save, etc.)
- ✅ Creates checkboxes (6 checklist items)
- ✅ Binds dropdowns to tables

### 4. Defines Settings
- ✅ Equity_E = $10,000 (account size)
- ✅ RiskPct_r = 0.75% (risk per trade)
- ✅ HeatCap_H_pct = 4% (max portfolio heat)
- ✅ And more...

---

## After Running Setup

You should see:
1. **All sheets exist** (check the sheet tabs at bottom)
2. **TradeEntry sheet has UI** (labels, buttons, input fields)
3. **Dropdowns work:**
   - B5 (Preset) → 5 options
   - B7 (Sector) → 6 options
   - B8 (Bucket) → 6 options
4. **No error messages**

---

## Troubleshooting

### "Button doesn't exist"
**Solution:** The Setup sheet wasn't created properly.
1. Press Alt+F11 to open VBA
2. Press Ctrl+G for Immediate window
3. Type: `Setup.RunInitialSetup`
4. Press Enter

### "Run-time error"
**Solution:** Setup already ran or sheets already exist.
- This is usually OK
- Check if TradeEntry sheet exists
- Check if dropdowns work
- If still broken, try: `Setup.RunInitialSetup` again (it will overwrite)

### "Compile error"
**Solution:** VBA modules not imported correctly.
- Close Excel
- Run BUILD.bat again
- Make sure all modules import successfully

### "Nothing happens"
**Solution:** Check the debug log.
1. Click "Open Debug Log" button
2. Look for errors
3. Search for "RunInitialSetup"
4. Check what went wrong

---

## When to Run Setup Again

You might need to re-run setup if:
- 🔄 Dropdowns stop working
- 🔄 Tables get deleted accidentally
- 🔄 Sheets get corrupted
- 🔄 You want to reset to defaults
- 🔄 Error messages about missing tables

**Running setup again is safe** - it will:
- ✅ Overwrite existing sheets/tables
- ✅ Preserve your existing data (if any)
- ✅ Re-seed default presets and buckets
- ✅ Rebuild the UI

---

## Quick Reference

| Issue | Solution |
|-------|----------|
| Dropdown error on opening | Click ▶ RUN INITIAL SETUP button |
| Can't find button | Go to Setup sheet |
| Button doesn't work | Use VBA: `Setup.RunInitialSetup` |
| Don't know VBA | Alt+F11, Ctrl+G, type command, Enter |
| Still doesn't work | Check debug log for errors |

---

## Video Walkthrough (Text Version)

**Opening the file:**
1. Double-click Excel file
2. Click "Enable Content" (yellow bar)
3. If error appears → Go to Step 4

**Running setup:**
4. Look at sheet tabs at bottom
5. Click "Setup" tab
6. Scroll down to buttons section
7. Click the big button: **▶ RUN INITIAL SETUP**
8. Wait for "Setup Complete!" message
9. Click OK

**Verifying it worked:**
10. Go to TradeEntry sheet
11. Click cell B5 (Preset)
12. Should see dropdown with 5 presets
13. ✅ Success!

---

## Still Having Issues?

If setup still doesn't work:

1. **Open debug log** (button on Setup sheet)
2. **Search for** "RunInitialSetup"
3. **Look for ERROR** messages
4. **Share the log** with the error details

Common errors and fixes:
- **Error 1004**: Sheets already exist (OK to ignore)
- **Error 9**: Can't find module/sheet (VBA not imported correctly - rebuild)
- **Error 70**: Permission denied (Close other Excel windows, try again)

---

## Summary

**The easy way:**
1. Open workbook
2. Go to Setup sheet
3. Click **▶ RUN INITIAL SETUP**
4. Done!

**The VBA way:**
1. Alt+F11
2. Ctrl+G
3. Type: `Setup.RunInitialSetup`
4. Enter

Both do the same thing - pick whichever is easier!
