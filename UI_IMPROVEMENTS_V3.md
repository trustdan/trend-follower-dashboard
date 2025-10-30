# TF-Engine GUI v3 - UI Improvements

**Date:** October 30, 2025
**Build:** ui/tf-gui-v3.exe
**Status:** ✅ Scroll container heights increased

---

## Changes in v3

### 1. Increased Scroll Container Heights ✅

**Problem:** Scroll containers for candidates list and scanner results were too small vertically (tiny boxes)

**Solution:** Set minimum sizes for scroll containers to provide better visibility

**Files Modified:**
- [ui/dashboard.go](ui/dashboard.go:239-240) - Candidates list scroll container
- [ui/scanner.go](ui/scanner.go:204-206) - Scanner results scroll container

**Dashboard Candidates List:**
```go
// Create scroll container with minimum height for better visibility
scroll := container.NewScroll(candidatesList)
scroll.SetMinSize(fyne.NewSize(200, 300))  // Min width 200, height 300 pixels
```

**Scanner Results:**
```go
// Create scroll container for results with minimum height
resultsScroll := container.NewScroll(resultsContainer)
resultsScroll.SetMinSize(fyne.NewSize(400, 200))  // Min width 400, height 200 pixels
```

**Result:** Candidates and scanner results now display in much larger, more usable containers

---

### 2. Old Executables Cleaned Up ✅

Removed old GUI builds to save space (each was ~50MB):
- ❌ Deleted: tf-gui.exe
- ❌ Deleted: tf-gui-debug.exe
- ❌ Deleted: tf-gui-final.exe
- ❌ Deleted: tf-gui-fixed.exe
- ✅ Kept: tf-gui-v2.exe (button contrast + dual dropdowns)
- ✅ **New:** tf-gui-v3.exe (larger scroll containers)

---

## Outstanding Issues

### 1. Button Text Contrast (Still Present)

**Issue:** Black text on dark green buttons in light mode is hard to read

**Status:** ⚠️ Partial fix attempted but not fully working

**Technical Challenge:**
Fyne's button widget uses `ColorNameForeground` for button text color, but we need `ColorNameForeground` to be black for general text (labels) on white backgrounds in light mode. We can't easily distinguish between "text on button" vs "text on white background" in the theme's `Color()` method.

**Attempted Fixes:**
1. ✅ Added `ColorNameForegroundOnPrimary` = white (should work but Fyne may not use it for all button types)
2. ✅ Added `ColorNameForegroundOnError`, `ForegoundOnSuccess` = white
3. ❌ Cannot set `ColorNameForeground` = white in light mode (would make all text invisible on white backgrounds)

**Possible Solutions:**
1. **Custom Button Widget** - Create a custom button that manually sets text color to white
2. **Button Importance** - Set button importance to `widget.HighImportance` which may trigger different color handling
3. **Custom Render** - Override button rendering to force white text
4. **Accept Compromise** - Keep buttons in dark mode only, or accept reduced contrast in light mode

**Need User Decision:** Which approach to take? Or is the current state acceptable?

---

### 2. Bucket Grouping Strategy (Awaiting User Input)

**Current Implementation:**
- Dual dropdown system: Sector + Bucket
- Auto-mapping from sector to bucket
- Two buckets currently grouped:
  - **Tech/Comm** = Technology + Communication Services
  - **Consumer** = Consumer Discretionary + Consumer Staples

**User Feedback:** "the buckets should be more joined than just tech and comm together"

**Question for User:** What is the intended bucket grouping strategy?

**Option A - Growth/Defensive/Cyclical (3 large buckets):**
1. **Growth** - Tech + Comm + Consumer Discretionary (risk-on, growth-oriented)
2. **Defensive** - Utilities + Real Estate + Consumer Staples (low volatility, defensive)
3. **Cyclical** - Industrials + Materials + Energy (commodity/cycle driven)
4. **Financial** - Financials (standalone)
5. **Healthcare** - Healthcare (standalone)

**Option B - More Granular (7 buckets):**
1. **Tech/Comm** - Technology + Communication
2. **Consumer** - Consumer Discretionary + Staples
3. **Energy/Materials** - Energy + Materials (commodities)
4. **Industrials** - Industrials
5. **Financials** - Financials
6. **Healthcare** - Healthcare
7. **Defensive** - Utilities + Real Estate

**Option C - Custom Grouping:**
Please specify which sectors should be grouped together for heat tracking purposes

**Rationale from Documentation:**
- Bucket caps (1.5% each) prevent concentration in CORRELATED sectors
- When Tech crashes, Communication crashes too → group them
- When commodities spike, Energy + Materials move together → group them?
- Utilities + Real Estate both defensive/dividend plays → group them?

---

## Testing Checklist

### Scroll Container Fixes
- [ ] Run tf-gui-v3.exe
- [ ] Navigate to Dashboard
- [ ] Check "Today's Candidates" card - list should be taller (300px min height)
- [ ] Navigate to Scanner
- [ ] Run a FINVIZ scan
- [ ] Check results display - should be wider and taller (400x200 min size)
- [ ] Verify clickable TradingView links still work

### Button Contrast Test
- [ ] Toggle to light mode
- [ ] Check Scanner preset buttons - text readable?
- [ ] Check "Scan FINVIZ & Import" button - text readable?
- [ ] Check "Edit Settings" button - text readable?
- [ ] Check all buttons throughout app
- [ ] Take screenshots if still problematic

### Bucket System Test
- [ ] Navigate to Trade Entry
- [ ] Verify both dropdowns exist: Sector and Bucket
- [ ] Select Technology → verify Tech/Comm auto-populates
- [ ] Select Communication Services → verify Tech/Comm auto-populates
- [ ] Select Consumer Discretionary → verify Consumer auto-populates
- [ ] **Provide feedback:** What other sectors should be grouped?

---

## Build Info

**File:** ui/tf-gui-v3.exe
**Size:** ~49MB
**Changes from v2:**
- Larger scroll containers for candidates and scanner results
- No functional changes to button colors or bucket grouping

**Unchanged from v2:**
- Dual dropdown system (Sector + Bucket)
- British Racing Green theme
- Dark/light mode toggle
- All 5 FINVIZ presets
- TradingView clickable links

---

## Next Steps

### Immediate (Pending User Feedback)
1. **Button Color Decision:**
   - Test tf-gui-v3.exe in light mode
   - Report if button text contrast is acceptable or needs more work
   - If needs work: choose solution approach

2. **Bucket Grouping Strategy:**
   - Review proposed bucket options (A, B, or C above)
   - Specify which sectors should be grouped together
   - Confirm heat cap strategy (1.5% per bucket)

### After User Feedback
3. **Implement Chosen Bucket Strategy:**
   - Update sector-to-bucket mapping in trade_entry.go
   - Update bucket list in trade_entry.go and calendar.go
   - Test heat calculations with new groupings
   - Build v4

4. **Final Button Color Fix (if needed):**
   - Implement chosen solution (custom widget, importance, or accept current)
   - Test thoroughly in both light and dark modes
   - Build v4 or v5

---

## Questions for User

1. **Button Contrast:** Test light mode - is it acceptable or needs more work?

2. **Bucket Grouping:** Which option (A, B, C, or custom)?
   - How many total buckets do you want?
   - Which sectors correlate and should share a bucket?
   - Do you want Energy+Materials grouped (commodities)?
   - Do you want Utilities+Real Estate grouped (defensive)?

3. **Scroll Heights:** Are the new container sizes good, or need adjustment?

4. **Any other UI issues** noticed during testing?

---

**Status:** ✅ v3 ready for testing - awaiting feedback on buttons and bucket strategy
