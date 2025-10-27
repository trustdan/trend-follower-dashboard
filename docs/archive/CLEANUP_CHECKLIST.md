# File Cleanup Checklist

Quick action items to organize the project.

---

## IMMEDIATE ACTIONS (5 minutes)

### 1. Archive Old Build Scripts
```bash
mkdir -p scripts/deprecated
mv old_CLEANUP_STUCK_EXCEL.bat scripts/deprecated/
mv old_IMPORT_VBA_MODULES.bat scripts/deprecated/
mv old_IMPORT_VBA_MODULES_DEBUG.bat scripts/deprecated/
mv old_VERIFY_MODULES.bat scripts/deprecated/
```
- [ ] Done

### 2. Remove VBA Duplicates from Python/
```bash
# First, review if TF_Presets_Enhanced.bas should replace TF_Presets.bas
# Then:
rm Python/TF_Python_Bridge.bas
rm Python/TF_Presets_Enhanced.bas
```
- [ ] Reviewed TF_Presets_Enhanced.bas
- [ ] Removed duplicates

---

## MEDIUM PRIORITY (10-15 minutes)

### 3. Create Documentation Structure
```bash
mkdir -p docs/setup
mkdir -p docs/specifications
mkdir -p docs/reference
mkdir -p docs/archive
```
- [ ] Done

### 4. Move Documentation Files
Move these to `docs/setup/`:
- [ ] VBA_SETUP_GUIDE.md
- [ ] PYTHON_SETUP_GUIDE.md
- [ ] (Consolidate: VBA_IMPLEMENTATION_SUMMARY.md, PYTHON_IMPLEMENTATION_SUMMARY.md)
- [ ] (Consolidate: TROUBLESHOOTING_BUILD_ISSUES.md → TROUBLESHOOTING.md)

Move these to `docs/specifications/`:
- [ ] newest-Interactive_TF_Workbook_Plan.md
- [ ] workflow-plan.md
- [ ] SeykotaTurtleTrend-FollowingOptionsExecution+TradingViewStrategyGuide.md
- [ ] diversification-across-sectors.md
- [ ] diversification-across-sectors.pdf

Move these to `docs/reference/`:
- [ ] VBA_README.md
- [ ] Workbook_Readme_Text.md

Move these to `docs/archive/`:
- [ ] older-Options_Trend_Dashboard_Summary.md (rename to REFERENCE_OLDER_DASHBOARD.md)
- [ ] BUILD_NOW_COMPLETE.md
- [ ] FINAL_FIXES.md
- [ ] BUILD_COMPLETE.md
- [ ] FINAL_STATUS.md
- [ ] FIXES_APPLIED.md
- [ ] LATEST_FIX.md
- [ ] TWO_BUILD_OPTIONS.md
- [ ] FIX_BROKEN_VENV.md
- [ ] FIX_LOG_FILE_CONFLICT.md
- [ ] IGNORE_EXCEL_PYTHON_ERROR.md
- [ ] WHATS_MISSING.md
- [ ] TROUBLESHOOTING_BUILD_ISSUES.md
- [ ] VBA_IMPLEMENTATION_SUMMARY.md
- [ ] PYTHON_IMPLEMENTATION_SUMMARY.md
- [ ] README_BUILD.md
- [ ] VBA_IMPORT_SCRIPT.txt

### 5. Create docs/README.md (index)
- [ ] Create documentation index/guide

### 6. Archive Build Logs
```bash
mkdir -p logs/archive
mv logs/*.log logs/archive/
```
- [ ] Done

---

## LOW PRIORITY (Future)

- [ ] Add CHANGELOG.md
- [ ] Add CONTRIBUTING.md
- [ ] Consolidate implementation summaries into setup guides
- [ ] Update config/presets.json (only has 2 presets, should have 5)

---

## CLEANUP BENEFITS

| Metric | Before | After |
|--------|--------|-------|
| Root-level .md files | 28 | 4 |
| Old batch files in root | 4 | 0 |
| Python folder .bas duplicates | 2 | 0 |
| Organized documentation | No | Yes |
| Easier to navigate | No | Yes |
| Clear "what to read" order | No | Yes |

---

## VERIFICATION

After cleanup, run:
```bash
find . -name "*.md" -type f | wc -l
# Should be: 4 in root + ~24 in docs/ = 28 total (same, just organized)

find . -name "old_*.bat" -type f | wc -l
# Should be: 0 in root (moved to scripts/deprecated/)

find Python/ -name "*.bas" -type f | wc -l
# Should be: 0 (moved to VBA/)
```

