# Quick Troubleshooting Checklist

## When Things Don't Work - Follow These Steps

### Step 1: Open the Debug Log
1. Go to **Setup sheet**
2. Click **"Open Debug Log"** button
3. Log file opens in Notepad

If button doesn't work, manually open:
`TradingSystem_Debug.log` (same folder as Excel file)

---

## Problem: Checkboxes Overlapping Text

### What to Check in Log:

Search for: `CreateCheckboxes()`

**Look for this section:**
```
--- CreateCheckboxes() - Start ---
Column B Left position: [NUMBER]
Checkbox leftPos: [NUMBER]
Creating checkbox 1 at row 21, topPos=420, leftPos=[NUMBER]
Checkbox 1 created successfully
Total checkboxes created: 6 of 6
```

### Diagnosis:

| What You See | What It Means | Solution |
|--------------|---------------|----------|
| `leftPos: 10` or very small number | Positioning code not working | VBA didn't update - rebuild workbook |
| `leftPos: 165` or ~160-200 | Position correct | Checkboxes should be properly placed |
| `Failed to add checkbox` error | Checkbox creation failed | COM automation issue - add manually |
| `Total checkboxes created: 0 of 6` | All checkboxes failed | Security settings blocking automation |
| `Total checkboxes created: X of 6` (X < 6) | Partial failure | Some succeeded, manually add missing ones |

---

## Problem: Preset Dropdown Not Working

### What to Check in Log:

Search for: `BindControls()`

**Look for this section:**
```
--- BindControls() - Start ---
tblPresets exists with 5 rows
Setting up Preset dropdown in B5...
Preset dropdown created successfully
```

### Diagnosis:

| What You See | What It Means | Solution |
|--------------|---------------|----------|
| `tblPresets does not exist` | Tables not created | Run `Setup.RunInitialSetup` |
| `tblPresets exists with 0 rows` | Table empty | Run `Setup.RunInitialSetup` or add presets manually |
| `tblPresets exists with 5 rows` + `created successfully` | Should work! | Test again, if still broken - check cell B5 validation manually |
| `ERROR in BindControls: Preset dropdown failed` | Validation setup failed | Check error details - usually means table reference broken |

---

## Problem: Python Detection Says "Not Available"

### What to Check in Log:

Search for: `IsPythonAvailable()`

**Look for this section:**
```
--- IsPythonAvailable() - Start ---
Excel Version: [VERSION]
Formula2 property succeeded
Cell value type: [TYPE]
Cell value: [VALUE]
IsPythonAvailable() - Result: [True/False]
```

### Diagnosis:

| What You See | What It Means | Solution |
|--------------|---------------|----------|
| `Formula2 assignment error: [424]` | Excel too old for Formula2 | Normal - will fallback to Formula |
| `Cell value: 2` + `Result: True` | **PYTHON WORKS!** | Detection should show Available |
| `Cell value type: Error` | Python function not recognized | Python in Excel not available in your Office version |
| `Cell is empty` | Formula didn't execute | Increase wait time or Python disabled |
| `Cell value: 0` or wrong value | Formula ran but incorrectly | Excel version or Python engine issue |
| `Error 2015` | #NAME? error | `PY()` function doesn't exist in this Excel |

### Manual Test:

To confirm Python works:
1. Go to any cell
2. Type: `=PY(1+1)`
3. Press Enter
4. If you get `2` → Python works, detection should work
5. If you get `#NAME?` error → Python not available in your Excel

---

## Problem: Setup Doesn't Complete

### What to Check in Log:

Search for: `RunInitialSetup`

**Look for errors:**
```
ERROR in [FunctionName]: [Error Number] [Error Description]
```

### Common Errors:

| Error | What It Means | Solution |
|-------|---------------|----------|
| Error 9: Subscript out of range | Sheet not found | Sheet deleted or renamed |
| Error 1004: Method failed | Various VBA issues | Check specific error message |
| Error 70: Permission denied | File/folder locked | Close other Excel instances |

---

## General Troubleshooting Steps

### If Something Doesn't Work:

1. **Check the log file**
   - Look for ERROR messages
   - Look for the specific function (CreateCheckboxes, BindControls, IsPythonAvailable)
   - Note error numbers and descriptions

2. **Verify prerequisites**
   - [ ] Excel file built from latest VBA code?
   - [ ] All sheets created (Setup, TradeEntry, Presets, etc.)?
   - [ ] Setup.RunInitialSetup completed successfully?

3. **Try rebuilding**
   - Copy project to Windows
   - Run `BUILD.bat`
   - Open new Excel file
   - Check if issue persists

4. **Manual workarounds**
   - **Checkboxes:** Add manually via Developer → Insert → Check Box
   - **Dropdowns:** Manually set Data Validation on cells
   - **Tables:** Check Presets sheet, ensure tblPresets table exists

---

## How to Share Logs When Reporting Issues

If you need help:

1. Reproduce the issue
2. Open Debug Log (button or file)
3. **Copy entire log file** contents
4. Share with:
   - Your Excel version (File → Account → About Excel)
   - Screenshot of the problem
   - What you were trying to do

---

## Quick Reference: Where Everything Is

| Feature | Location | How to Test |
|---------|----------|-------------|
| Checkboxes | TradeEntry sheet, rows 21-26, column B | Look for ☐ symbols before text |
| Preset Dropdown | TradeEntry sheet, cell B5 | Click cell, see dropdown arrow? |
| Ticker Dropdown | TradeEntry sheet, cell B6 | Click cell, see dropdown arrow? |
| Sector Dropdown | TradeEntry sheet, cell B7 | Click cell, should always work |
| Bucket Dropdown | TradeEntry sheet, cell B8 | Click cell, see dropdown arrow? |
| Python Test | Setup sheet, "Test Python Integration" button | Should say "Available" if working |
| Debug Log | Setup sheet, "Open Debug Log" button | Opens Notepad with log |

---

## Still Having Issues?

If nothing works after checking logs and following troubleshooting:

1. **Share your log file** - Copy entire `TradingSystem_Debug.log`
2. **Share Excel version** - File → Account → About Excel
3. **Share screenshot** - Show what's not working
4. **Describe exact steps** - What you clicked, what you expected, what happened

The log file will reveal exactly what's going wrong!
