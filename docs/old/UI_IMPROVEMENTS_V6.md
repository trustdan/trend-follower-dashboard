# TF-Engine GUI v6 - Complete Button Text Fix

**Build:** `ui/tf-gui-v6.exe` (49MB)
**Date:** 2025-10-30
**Status:** ✅ COMPLETE

---

## Changes from v5

### Fixed ALL Remaining Black Text Buttons

Applied `button.Importance = widget.HighImportance` to every remaining button that had black text:

#### 1. **main.go** - Theme Toggle
- ✅ "🌙 Dark Mode" / "☀️ Light Mode" button

#### 2. **checklist.go** - Checklist Screen
- ✅ "Evaluate Checklist" button
- ✅ "Reset" button

#### 3. **position_sizing.go** - Position Sizing Screen
- ✅ "Calculate Position Size" button

#### 4. **heat_check.go** - Heat Check Screen
- ✅ "Check If Trade Allowed" button

#### 5. **trade_entry.go** - Trade Entry Screen
- ✅ "Check All 5 Gates" button
- ✅ "Save NO-GO Decision" button

#### 6. **calendar.go** - Calendar Screen
- ✅ "Refresh Calendar" button

---

## How It Works

Setting `Importance = widget.HighImportance` tells Fyne to use the primary color scheme, which applies:
- **Light Mode:** White text on British Racing Green background
- **Dark Mode:** White text on British Racing Green background

All buttons should now be clearly readable in both modes.

---

## Previously Fixed (v4 & v5)

These were already working from previous versions:

### Scanner Screen (v4)
- ✅ "TF-Breakout-Long" preset button
- ✅ "TF-Momentum-Uptrend" preset button
- ✅ "TF-Unusual-Volume" preset button
- ✅ "TF-Breakdown-Short" preset button
- ✅ "TF-Momentum-Downtrend" preset button
- ✅ "Scan FINVIZ & Import" button

### Dashboard Screen (v5)
- ✅ "Edit Settings" button
- ✅ "Refresh" button

---

## Test Instructions

1. **Run the GUI:**
   ```powershell
   cd ui
   .\tf-gui-v6.exe
   ```

2. **Test Light Mode:**
   - Open GUI (starts in light mode by default)
   - Navigate to each screen
   - Verify ALL buttons show **white text** on green background

3. **Test Dark Mode:**
   - Click "🌙 Dark Mode" in top-left corner
   - Navigate to each screen again
   - Verify ALL buttons show **white text** on green background

4. **Screens to Check:**
   - ✅ Dashboard → "Edit Settings", "Refresh"
   - ✅ Checklist → "Evaluate Checklist", "Reset"
   - ✅ Position Sizing → "Calculate Position Size"
   - ✅ Heat Check → "Check If Trade Allowed"
   - ✅ Trade Entry → "Check All 5 Gates", "Save NO-GO Decision"
   - ✅ Scanner → All 5 preset buttons + "Scan FINVIZ & Import"
   - ✅ Calendar → "Refresh Calendar"
   - ✅ Theme Toggle → "🌙 Dark Mode" / "☀️ Light Mode"

---

## What's Next?

Per your feature requests (still TODO):

### 1. Expandable Calendar Rows (Medium Effort)
"If I enter multiple trades in the tech/comm bucket, it should create another row for tech/comm in the calendar tab"

**Plan:**
- Modify `calendar.go` to create multiple rows per bucket when there are multiple positions
- Example:
  ```
  Tech/Comm (1)  [NVDA]  [MSFT]  -      -     ...
  Tech/Comm (2)  -       [META]  [GOOGL] -    ...
  ```

### 2. Startup Welcome Tooltip (Quick Win)
"Tooltip on startup with a checkbox not to show again"

**Plan:**
- Create welcome dialog that shows on first run
- Add checkbox: "Don't show this again"
- Store preference in database or settings file

### 3. Help Button & FAQ/README (Moderate Effort)
"Help button at the top (question mark), FAQs and README pages"

**Plan:**
- Add "?" help button to top menu bar
- Create dialog/screen showing:
  - README content
  - FAQ section
  - Quick start guide

### 4. VIM Keybindings (Complex - Save for Later?)
"Optional VIM keybindings (like the vimium browser extension), where you press f and it highlights the clickable links and assigns keyboard shortcuts to them"

**Plan (Complex in Fyne):**
- Global keyboard handler
- Link hint overlay system
- VIM navigation (j/k/h/l for movement)
- Preference toggle to enable/disable
- This is the most involved feature - might save for v2.0

---

## Recommended Implementation Order

1. **Startup welcome tooltip** (Quick win, immediate value)
2. **Help button + FAQ/README** (Moderate effort, high value for new users)
3. **Expandable calendar rows** (Moderate effort, needs design clarification)
4. **VIM keybindings** (Complex, consider for future major version)

---

## Summary

**v6 = Complete Button Text Fix**

Every single button in the application now has white text on the British Racing Green background in both light and dark modes. No more black-on-green text visibility issues.

Next step: Get your feedback on v6, then tackle the feature requests above.
