# TF-Engine GUI v4 - Release Notes

**Build Date:** October 30, 2025
**File:** ui/tf-gui-v4.exe (49MB)

---

## What's New in v4

### ✅ Complete Bucket Grouping System

Implemented comprehensive sector-to-bucket mapping based on market correlation:

**12 Individual Sectors (FINVIZ order):**
1. Basic Materials
2. Communication Services
3. Consumer Cyclical
4. Consumer Defensive
5. Energy
6. Financial
7. Healthcare
8. Industrials
9. Real Estate
10. Technology
11. Utilities
12. ETFs

**8 Heat Tracking Buckets:**
1. **Materials/Industrials** - Basic Materials + Industrials (commodity/cycle driven)
2. **Tech/Comm** - Communication Services + Technology (growth/innovation)
3. **Financial/Cyclical** - Financial + Consumer Cyclical (economic sensitivity)
4. **Defensive/Utilities** - Consumer Defensive + Utilities (defensive/dividends)
5. **Energy** - Standalone (commodity/burst sector)
6. **Healthcare** - Standalone (demographics/innovation)
7. **Real Estate** - Standalone (interest rate sensitive)
8. **ETFs** - Standalone (multi-sector tracking)

### Key Features

- **Auto-mapping:** Select a sector, bucket auto-populates based on correlation
- **Manual override:** Can manually change bucket if needed
- **Calendar integration:** 10-week view shows all 8 buckets
- **Diversification enforcement:** 1.5% cap per bucket prevents concentration in correlated sectors

---

## What's Carried Forward from v3

### ✅ Larger Scroll Containers
- Dashboard candidates list: 300px height
- Scanner results: 400x200px
- Much better visibility

### ✅ Dual Dropdown System
- Sector selection (individual classification)
- Bucket selection (heat tracking)
- Clear separation of concerns

### ✅ Theme & Colors
- British Racing Green primary color
- Dark/light mode toggle
- Banner colors (GREEN/YELLOW/RED)

### ✅ FINVIZ Scanner
- All 5 presets
- Clickable TradingView links
- Auto-import to database

---

## Changes from Previous Versions

### v4 vs v3
- **New:** 12-sector list (FINVIZ order + ETFs)
- **New:** 8-bucket grouping system (4 grouped + 4 standalone)
- **New:** Materials/Industrials bucket
- **New:** Financial/Cyclical bucket
- **New:** Defensive/Utilities bucket
- **New:** ETFs bucket
- **Changed:** Tech/Comm now includes Communication Services (not just Tech)
- **Removed:** Old sector names (Financial Services → Financial, etc.)

### v3 vs v2
- **New:** Larger scroll containers (300px and 400x200px)
- **Removed:** Old GUI executables (saved ~200MB)

### v2 vs v1
- **New:** Dual dropdown system (Sector + Bucket)
- **New:** Button contrast improvements (ColorNameForegroundOnPrimary)
- **Changed:** Bucket grouping (preliminary version)

---

## Known Issues

### Button Text Contrast in Light Mode
**Status:** Partial fix applied, may still have visibility issues

**Issue:** Some buttons may show black text on dark green in light mode

**Workaround:** Use dark mode for best visibility

**Technical:** Fyne theme system doesn't provide easy way to distinguish button text vs general text. `ColorNameForegroundOnPrimary` is set to white but may not apply to all button types.

**Future Fix Options:**
1. Custom button widgets
2. Force button importance levels
3. Accept current state and use dark mode

---

## Upgrade Guide

### From v3 to v4

**No database changes required** - bucket system is UI-only

**Testing Steps:**
1. Run ui/tf-gui-v4.exe
2. Navigate to Trade Entry
3. Test sector selections:
   - Technology → should auto-select Tech/Comm
   - Basic Materials → should auto-select Materials/Industrials
   - Financial → should auto-select Financial/Cyclical
   - Consumer Defensive → should auto-select Defensive/Utilities
4. Navigate to Calendar
5. Verify 8 bucket rows display
6. Test creating a trade with new bucket system

**Data Migration:**
- Existing positions with old bucket names will still work
- New positions will use new bucket names
- No data loss or corruption

---

## Testing Checklist

### Critical Tests
- [ ] All 12 sectors appear in dropdown (FINVIZ order)
- [ ] Each sector maps to correct bucket
- [ ] Bucket auto-populates when sector selected
- [ ] Manual bucket override works
- [ ] Calendar shows 8 buckets
- [ ] Heat check uses bucket grouping correctly

### Correlation Tests
- [ ] Technology + Communication Services → same bucket (Tech/Comm)
- [ ] Basic Materials + Industrials → same bucket (Materials/Industrials)
- [ ] Financial + Consumer Cyclical → same bucket (Financial/Cyclical)
- [ ] Consumer Defensive + Utilities → same bucket (Defensive/Utilities)
- [ ] Energy, Healthcare, Real Estate, ETFs → separate buckets

### Integration Tests
- [ ] Scanner import still works
- [ ] Dashboard displays candidates
- [ ] Position sizing calculator works
- [ ] Heat check calculator works
- [ ] Checklist evaluation works
- [ ] Trade entry gate check works
- [ ] Calendar displays positions

---

## Documentation

**Comprehensive Bucket System Guide:**
- [BUCKET_SYSTEM_V4.md](BUCKET_SYSTEM_V4.md) - Complete implementation details
  - Design philosophy
  - Correlation rationale
  - Auto-mapping table
  - Example scenarios
  - Heat cap strategy
  - Testing checklist

**Previous Documentation:**
- [BUCKET_SYSTEM_FIX.md](BUCKET_SYSTEM_FIX.md) - v2 bucket fix
- [UI_IMPROVEMENTS_V3.md](UI_IMPROVEMENTS_V3.md) - v3 scroll containers

---

## File Sizes

```
ui/tf-gui-v2.exe  49MB  (dual dropdowns + contrast fix)
ui/tf-gui-v3.exe  49MB  (larger scroll containers)
ui/tf-gui-v4.exe  49MB  (complete bucket system)
```

**Recommended:** Delete v2 and v3 after testing v4 successfully

---

## What's Next

### Immediate
1. **Test v4** with new bucket groupings
2. **Verify** sector-to-bucket mappings match your trading strategy
3. **Feedback** on any bucket grouping adjustments needed

### Future Enhancements
1. **Button contrast fix** (if still problematic in light mode)
2. **Backend integration** (connect position sizing, heat check to actual backend)
3. **Database schema** update to store sector + bucket separately
4. **Heat visualization** (gauges, progress bars for each bucket)
5. **Cooldown tracking** by bucket (not just ticker)

---

## Quick Start

```bash
# Run the GUI
cd ui
./tf-gui-v4.exe

# Navigate to Trade Entry
# Select Technology → verify Tech/Comm auto-populates
# Select Basic Materials → verify Materials/Industrials auto-populates
# Select Financial → verify Financial/Cyclical auto-populates

# Navigate to Calendar
# Verify 8 bucket rows display correctly
```

---

## Support

**Issues?**
- Check [BUCKET_SYSTEM_V4.md](BUCKET_SYSTEM_V4.md) for detailed documentation
- Review testing checklist
- Verify FINVIZ sector names match your imported candidates

**Feedback Welcome:**
- Bucket grouping too aggressive or too loose?
- Should any sectors be regrouped?
- Any sectors missing from the list?

---

**Status:** ✅ v4 ready for testing with complete 8-bucket diversification system!
