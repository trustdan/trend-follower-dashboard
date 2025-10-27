# Excel Trading Workflow - Comprehensive File Inventory

**Generated**: 2025-10-26  
**Directory**: /home/kali/excel-trading-workflow

---

## EXECUTIVE SUMMARY

This project contains an **Excel VBA + Python trading system** with significant documentation redundancy. The codebase is mature and functional, but the documentation folder has accumulated multiple overlapping and superseding files during development iterations.

**Key Findings**:
- ✅ 10 active VBA modules (clean, organized)
- ✅ 3 active Python scripts (well-structured)
- ✅ Multiple build automation scripts
- ⚠ 28 markdown documentation files (12+ are outdated/redundant)
- ⚠ 4 build batch files marked "old_"
- ⚠ 3 build log files in logs/ directory

---

## 1. SCRIPT FILES INVENTORY

### Active Build Automation Scripts

| File | Location | Purpose | Status |
|------|----------|---------|--------|
| `BUILD_WITH_PYTHON.bat` | /home/kali/excel-trading-workflow/ | **PRIMARY BUILD SCRIPT** - Creates workbook, imports VBA, builds UI | ✅ ACTIVE |
| `COMPLETE_SETUP.bat` | /home/kali/excel-trading-workflow/ | Alternative build flow | ✅ ACTIVE |
| `FIX_VENV.bat` | /home/kali/excel-trading-workflow/ | Repair Python virtual environment | ✅ UTILITY |
| `setup_venv.bat` | /home/kali/excel-trading-workflow/scripts/ | Initialize Python venv | ✅ ACTIVE |
| `refresh_data.bat` | /home/kali/excel-trading-workflow/scripts/ | Refresh candidate data | ✅ UTILITY |

### Deprecated/Old Build Scripts (Marked "old_")

| File | Status | Reason |
|------|--------|--------|
| `old_CLEANUP_STUCK_EXCEL.bat` | ❌ DEPRECATED | Replaced by better cleanup in BUILD_WITH_PYTHON.bat |
| `old_IMPORT_VBA_MODULES.bat` | ❌ DEPRECATED | Replaced by Python-based import (import_to_excel.py) |
| `old_IMPORT_VBA_MODULES_DEBUG.bat` | ❌ DEPRECATED | Replaced by Python-based import with logging |
| `old_VERIFY_MODULES.bat` | ❌ DEPRECATED | Replaced by import_to_excel.py validation |

### VBA Script Files (Execution)

| File | Location | Purpose |
|------|----------|---------|
| `run_macro.vbs` | /home/kali/excel-trading-workflow/scripts/ | VBScript to invoke Excel macros |
| `excel_build_repo_aware.vbs` | /home/kali/excel-trading-workflow/scripts/ | Context-aware workbook builder |
| `excel_build_repo_aware_logged.vbs` | /home/kali/excel-trading-workflow/scripts/ | Builder with detailed logging |

### Python Scripts

| File | Location | Purpose | Status |
|------|----------|---------|--------|
| `import_to_excel.py` | /home/kali/excel-trading-workflow/ | **PRIMARY** - Automates VBA import into Excel | ✅ ACTIVE |
| `finviz_scraper.py` | /home/kali/excel-trading-workflow/Python/ | Scrapes FINVIZ screener tickers (optional) | ✅ READY |
| `heat_calculator.py` | /home/kali/excel-trading-workflow/Python/ | Fast portfolio heat calculations | ✅ READY |

---

## 2. VBA MODULES INVENTORY

All located in `/home/kali/excel-trading-workflow/VBA/`

### Standard Modules (.bas)

| Module | Lines | Purpose | Status |
|--------|-------|---------|--------|
| `TF_Utils.bas` | ? | Helper functions (sheet/table/name management) | ✅ ACTIVE |
| `TF_Data.bas` | ? | Structure setup, heat calculations, cooldown logic | ✅ ACTIVE |
| `TF_UI.bas` | ? | UI controls, checklist, sizing, save decision | ✅ ACTIVE |
| `TF_Presets.bas` | ? | FINVIZ integration, candidate import | ✅ ACTIVE |
| `TF_UI_Builder.bas` | ? | Auto-builds TradeEntry sheet UI (NEW) | ✅ ACTIVE |
| `Setup.bas` | ? | Workbook initialization | ✅ ACTIVE |
| `Python_Run.bas` | ? | Calls Python scripts from VBA | ✅ READY |
| `PQ_Setup.bas` | ? | Power Query integration (optional) | ✅ READY |

### Class Modules (.cls)

| Module | Purpose | Status |
|--------|---------|--------|
| `ThisWorkbook.cls` | Workbook event handlers | ✅ ACTIVE |
| `Sheet_TradeEntry.cls` | TradeEntry sheet event handlers | ✅ ACTIVE |

### Duplicate/Legacy VBA Files

Located in `/home/kali/excel-trading-workflow/Python/` (should be in /VBA/ only):

| File | Status | Note |
|------|--------|------|
| `TF_Python_Bridge.bas` | ⚠ DUPLICATE | Copy exists in VBA/, this one in Python/ |
| `TF_Presets_Enhanced.bas` | ⚠ DUPLICATE | Enhanced version; may supersede TF_Presets.bas |

**Recommendation**: Remove Python/ copies; keep only VBA/ originals.

---

## 3. EXCEL FILES

### Built Workbooks

| File | Status | Note |
|------|--------|------|
| `TrendFollowing_TradeEntry.xlsm` | ⚠ GIT-IGNORED | Built by BUILD_WITH_PYTHON.bat; not tracked |

**Note**: Per `.gitignore` (lines 13-16), all `.xlsm` files are excluded from version control.

---

## 4. DOCUMENTATION FILES - COMPREHENSIVE ANALYSIS

**Total**: 28 markdown files (extremely high for a single-project workbook)

### CORE DOCUMENTATION (Read First)

| File | Lines | Purpose | Status | Priority |
|------|-------|---------|--------|----------|
| **CLAUDE.md** | 435 | Project guidelines for Claude AI | ✅ CURRENT | 🔴 ESSENTIAL |
| **README.md** | 460 | Main system overview + quick reference | ✅ CURRENT | 🔴 ESSENTIAL |
| **GETTING_STARTED.md** | 351 | First-time setup walkthrough | ✅ CURRENT | 🟠 IMPORTANT |
| **QUICK_START.md** | 326 | One-line build instructions | ✅ CURRENT | 🟠 IMPORTANT |

### SPECIFICATION DOCUMENTS (Design Reference)

| File | Lines | Purpose | Status | Priority |
|------|-------|---------|--------|----------|
| **newest-Interactive_TF_Workbook_Plan.md** | 1,386 | Detailed workbook spec (18KB) | ✅ CURRENT | 🟠 IMPORTANT |
| **workflow-plan.md** | 239 | Trading rules + workflow (11KB) | ✅ CURRENT | 🟠 IMPORTANT |
| **SeykotaTurtleTrend-FollowingOptionsExecution+TradingViewStrategyGuide.md** | 563 | TradingView strategy + Pine script guide | ✅ CURRENT | 🟡 REFERENCE |
| **older-Options_Trend_Dashboard_Summary.md** | 1,481 | Alternative dashboard spec (61KB) | ⚠ OLDER | 🟡 ARCHIVE |

### SETUP/BUILD GUIDES

| File | Lines | Purpose | Status | Priority |
|------|-------|---------|--------|----------|
| **VBA_SETUP_GUIDE.md** | 441 | Step-by-step VBA import | ✅ CURRENT | 🟠 IMPORTANT |
| **PYTHON_SETUP_GUIDE.md** | 449 | Python venv setup (optional) | ✅ CURRENT | 🟡 OPTIONAL |
| **VBA_README.md** | 312 | VBA architecture reference | ✅ CURRENT | 🟡 REFERENCE |
| **README_BUILD.md** | 183 | Build process documentation | ⚠ PARTIAL | 🟡 REFERENCE |

### STATUS/PROGRESS FILES (Outdated Build Logs)

**All superseded by BUILD_NOW_COMPLETE.md (most recent)**

| File | Lines | Date | Status | Action |
|------|-------|------|--------|--------|
| **BUILD_NOW_COMPLETE.md** | 319 | Oct 26 22:01 | ✅ LATEST | KEEP |
| **FINAL_FIXES.md** | 177 | Oct 26 21:34 | ✅ CURRENT | KEEP |
| **BUILD_COMPLETE.md** | 253 | Oct 26 16:27 | ⚠ OLDER | CONSIDER REMOVING |
| **FINAL_STATUS.md** | 255 | Oct 26 16:19 | ⚠ OLDER | CONSIDER REMOVING |
| **FIXES_APPLIED.md** | 177 | Oct 26 16:35 | ⚠ OLDER | CONSIDER REMOVING |
| **LATEST_FIX.md** | 165 | Oct 26 16:10 | ⚠ OLDER | REMOVE |
| **TWO_BUILD_OPTIONS.md** | 155 | Oct 26 16:23 | ⚠ OLDER | REMOVE |

### IMPLEMENTATION SUMMARIES (Duplicative Content)

| File | Lines | Purpose | Status | Action |
|------|-------|---------|--------|--------|
| **VBA_IMPLEMENTATION_SUMMARY.md** | 325 | VBA module checklist | ⚠ PARTIAL | CONSOLIDATE |
| **PYTHON_IMPLEMENTATION_SUMMARY.md** | 482 | Python status report | ⚠ PARTIAL | CONSOLIDATE |

### TROUBLESHOOTING / EDGE CASES

| File | Lines | Purpose | Status | Action |
|------|-------|---------|--------|--------|
| **TROUBLESHOOTING_BUILD_ISSUES.md** | 306 | Build error diagnosis | ⚠ PARTIAL | CONSOLIDATE |
| **WHATS_MISSING.md** | 336 | What still needs implementation | ⚠ PARTIAL | CONSOLIDATE |
| **FIX_BROKEN_VENV.md** | 90 | Venv recovery steps | ⚠ MINOR | CONSIDER REMOVING |
| **FIX_LOG_FILE_CONFLICT.md** | 139 | Log file issue resolution | ⚠ MINOR | REMOVE |
| **IGNORE_EXCEL_PYTHON_ERROR.md** | 139 | Known harmless errors | ⚠ MINOR | CONSIDER REMOVING |

### REFERENCE / DESIGN

| File | Lines | Purpose | Status |
|------|-------|---------|--------|
| **diversification-across-sectors.md** | 80 | Sector bucket strategy | ✅ REFERENCE |
| **diversification-across-sectors.pdf** | N/A (3.7MB) | Visual diagrams | ✅ REFERENCE |
| **Workbook_Readme_Text.md** | 174 | Workbook content guide | ⚠ PARTIAL |
| **VBA_IMPORT_SCRIPT.txt** | N/A | Old import notes | ⚠ OUTDATED |

---

## 5. CONFIGURATION FILES

| File | Location | Purpose | Status |
|------|----------|---------|--------|
| `presets.json` | /home/kali/excel-trading-workflow/config/ | FINVIZ screener presets (partial) | ⚠ INCOMPLETE |
| `settings.local.json` | /home/kali/excel-trading-workflow/.claude/ | Claude AI permissions | ✅ ACTIVE |

---

## 6. PYTHON VIRTUAL ENVIRONMENT & DEPENDENCIES

| Item | Location | Status |
|------|----------|--------|
| `venv/` | /home/kali/excel-trading-workflow/venv/ | ✅ ACTIVE (Python 3.14) |
| `requirements.txt` | /home/kali/excel-trading-workflow/Python/ | ✅ ACTIVE |

---

## 7. BUILD LOGS

| File | Created | Status | Action |
|------|---------|--------|--------|
| `build_20251026_161147.log` | Oct 26 16:11 | ⚠ OLDER | ARCHIVE |
| `build_20251026_161544.log` | Oct 26 16:15 | ⚠ OLDER | ARCHIVE |
| `build_20251026_162035.log` | Oct 26 16:20 | ⚠ OLDER | ARCHIVE |

**Recommendation**: Archive logs to `logs/archive/` or delete if already reviewed.

---

## 8. GITIGNORE COMPLIANCE

Files that SHOULD be tracked but are git-ignored:

```
# IGNORED (per .gitignore):
✗ *.xlsm  (workbook binaries)
✗ *.log   (build logs)
✗ venv/   (Python virtual env)
```

These are intentionally excluded (correct behavior).

---

## DUPLICATE/OUTDATED FILE ANALYSIS

### Clearly Deprecated Files

#### Old Build Scripts (safe to remove)
```
old_CLEANUP_STUCK_EXCEL.bat          ← Use BUILD_WITH_PYTHON.bat instead
old_IMPORT_VBA_MODULES.bat           ← Use BUILD_WITH_PYTHON.bat + import_to_excel.py
old_IMPORT_VBA_MODULES_DEBUG.bat     ← Use BUILD_WITH_PYTHON.bat + import_to_excel.py
old_VERIFY_MODULES.bat               ← Use import_to_excel.py validation
```

**Action**: Move to `scripts/deprecated/` or delete.

#### Documentation: Build Status Files

Created during development iterations (May 2025 - Oct 26 22:01):

**Definitively Obsolete** (earlier timestamps):
- `LATEST_FIX.md` (Oct 26 16:10)
- `FIX_LOG_FILE_CONFLICT.md` (Oct 26 16:14)
- `FIX_BROKEN_VENV.md` (Oct 26 16:30)
- `TWO_BUILD_OPTIONS.md` (Oct 26 16:23)
- `BUILD_COMPLETE.md` (Oct 26 16:27)
- `FINAL_STATUS.md` (Oct 26 16:19)
- `FIXES_APPLIED.md` (Oct 26 16:35)

**Superseded By**: `BUILD_NOW_COMPLETE.md` (Oct 26 22:01) + `FINAL_FIXES.md` (Oct 26 21:34)

**Action**: Archive to `docs/archive/` or delete.

#### Documentation: Partial Implementations

- `TROUBLESHOOTING_BUILD_ISSUES.md` - Partially resolved; consolidate into QUICK_START.md
- `WHATS_MISSING.md` - Dated Oct 26 22:00; may be accurate but redundant with build logs
- `IGNORE_EXCEL_PYTHON_ERROR.md` - Specific workaround; keep but consolidate into README.md or troubleshooting section
- `PYTHON_IMPLEMENTATION_SUMMARY.md` - Duplicate content with PYTHON_SETUP_GUIDE.md; merge them
- `VBA_IMPLEMENTATION_SUMMARY.md` - Duplicate content with VBA_README.md; merge them

### Structural Duplicates

#### In Python/ Folder (Should Be in VBA/ Only)
```
Python/TF_Python_Bridge.bas          ← Duplicate (move out of Python/)
Python/TF_Presets_Enhanced.bas       ← Enhanced version; may supersede TF_Presets.bas
```

**Action**: Keep VBA/ originals; remove Python/ copies. Review TF_Presets_Enhanced.bas to determine if it should replace TF_Presets.bas.

#### Specification Documents: Older vs Newer

```
older-Options_Trend_Dashboard_Summary.md  (1,481 lines, 61KB)  ← Legacy/alternative design
newest-Interactive_TF_Workbook_Plan.md    (1,386 lines, 52KB)  ← Current design
```

**Status**: Both exist because they represent different worksheet designs (dashboard vs trade entry). The "newest" is the currently implemented version.

**Action**: Keep both but clearly label in a manifest.

---

## 9. RECOMMENDED FOLDER STRUCTURE (Cleaned Up)

```
/home/kali/excel-trading-workflow/
│
├── 📂 VBA/                                    (ACTIVE)
│   ├── TF_Utils.bas
│   ├── TF_Data.bas
│   ├── TF_UI.bas
│   ├── TF_UI_Builder.bas
│   ├── TF_Presets.bas
│   ├── TF_Python_Bridge.bas
│   ├── Setup.bas
│   ├── Python_Run.bas
│   ├── PQ_Setup.bas
│   ├── ThisWorkbook.cls
│   └── Sheet_TradeEntry.cls
│
├── 📂 Python/                                 (ACTIVE)
│   ├── finviz_scraper.py
│   ├── heat_calculator.py
│   └── requirements.txt
│   └── (NO .bas files - move to VBA/)
│
├── 📂 scripts/                                (ACTIVE)
│   ├── BUILD_WITH_PYTHON.bat                 (PRIMARY)
│   ├── COMPLETE_SETUP.bat
│   ├── FIX_VENV.bat
│   ├── setup_venv.bat
│   ├── refresh_data.bat
│   ├── run_macro.vbs
│   ├── excel_build_repo_aware.vbs
│   ├── excel_build_repo_aware_logged.vbs
│   └── 📂 deprecated/                        (ARCHIVED)
│       ├── old_CLEANUP_STUCK_EXCEL.bat
│       ├── old_IMPORT_VBA_MODULES.bat
│       ├── old_IMPORT_VBA_MODULES_DEBUG.bat
│       └── old_VERIFY_MODULES.bat
│
├── 📂 config/                                 (ACTIVE)
│   └── presets.json
│
├── 📂 logs/                                   (UTILITY)
│   └── (current build logs only)
│   └── 📂 archive/                           (OPTIONAL)
│       └── (old build logs)
│
├── 📂 docs/                                   (ACTIVE)
│   ├── 📖 README.md                          (START HERE)
│   ├── 📖 CLAUDE.md                          (For Claude AI)
│   ├── 📖 GETTING_STARTED.md                 (First-time setup)
│   ├── 📖 QUICK_START.md                     (Build in 1 command)
│   │
│   ├── 📂 setup/
│   │   ├── VBA_SETUP_GUIDE.md
│   │   ├── PYTHON_SETUP_GUIDE.md
│   │   └── TROUBLESHOOTING.md (consolidated)
│   │
│   ├── 📂 specifications/
│   │   ├── newest-Interactive_TF_Workbook_Plan.md
│   │   ├── workflow-plan.md
│   │   ├── SeykotaTurtleTrend-FollowingOptionsExecution+TradingViewStrategyGuide.md
│   │   ├── diversification-across-sectors.md
│   │   └── diversification-across-sectors.pdf
│   │
│   ├── 📂 reference/
│   │   ├── VBA_README.md
│   │   ├── REFERENCE_OLDER_DASHBOARD.md (renamed from older-Options...)
│   │   └── Workbook_Readme_Text.md
│   │
│   └── 📂 archive/
│       ├── older-Options_Trend_Dashboard_Summary.md
│       ├── BUILD_NOW_COMPLETE.md
│       ├── FINAL_FIXES.md
│       ├── BUILD_COMPLETE.md
│       ├── FINAL_STATUS.md
│       ├── FIXES_APPLIED.md
│       ├── LATEST_FIX.md
│       ├── TWO_BUILD_OPTIONS.md
│       ├── FIX_BROKEN_VENV.md
│       ├── FIX_LOG_FILE_CONFLICT.md
│       ├── IGNORE_EXCEL_PYTHON_ERROR.md
│       ├── WHATS_MISSING.md
│       ├── TROUBLESHOOTING_BUILD_ISSUES.md
│       ├── VBA_IMPLEMENTATION_SUMMARY.md
│       ├── PYTHON_IMPLEMENTATION_SUMMARY.md
│       └── README_BUILD.md
│
├── import_to_excel.py                        (PRIMARY BUILD TOOL)
├── BUILD_WITH_PYTHON.bat                     (KEEP IN ROOT - entry point)
├── COMPLETE_SETUP.bat                        (KEEP IN ROOT - alternative)
├── .gitignore
└── venv/                                      (Python 3.14)
```

---

## 10. SUMMARY TABLE: WHAT TO KEEP VS ARCHIVE

### KEEP (Active Development)

| Category | Files | Count |
|----------|-------|-------|
| VBA Modules | TF_Utils, TF_Data, TF_UI, TF_UI_Builder, TF_Presets, TF_Python_Bridge, Setup, Python_Run, PQ_Setup, ThisWorkbook, Sheet_TradeEntry | 11 |
| Python Scripts | finviz_scraper.py, heat_calculator.py, import_to_excel.py | 3 |
| Build Scripts | BUILD_WITH_PYTHON.bat, COMPLETE_SETUP.bat, scripts/*.bat, scripts/*.vbs | 8 |
| Core Documentation | README.md, CLAUDE.md, GETTING_STARTED.md, QUICK_START.md | 4 |
| Specification Docs | newest-Interactive_TF_Workbook_Plan.md, workflow-plan.md, SeykotaTurtleTrend_*, diversification* | 5 |
| Setup Guides | VBA_SETUP_GUIDE.md, PYTHON_SETUP_GUIDE.md | 2 |
| Reference | VBA_README.md, Workbook_Readme_Text.md | 2 |
| Config | config/presets.json, .gitignore, .claude/settings.local.json | 3 |
| **TOTAL KEEP** | | **38 files** |

### ARCHIVE (Documentation History / Development Logs)

| Category | Files | Count |
|----------|-------|-------|
| Old Build Scripts | old_CLEANUP_STUCK_EXCEL.bat, old_IMPORT_*.bat, old_VERIFY_MODULES.bat | 4 |
| Build Status Logs | BUILD_NOW_COMPLETE.md, FINAL_FIXES.md, BUILD_COMPLETE.md, FINAL_STATUS.md, FIXES_APPLIED.md | 5 |
| Fix/Workaround Docs | LATEST_FIX.md, FIX_BROKEN_VENV.md, FIX_LOG_FILE_CONFLICT.md, IGNORE_EXCEL_PYTHON_ERROR.md | 4 |
| Redundant Guides | TROUBLESHOOTING_BUILD_ISSUES.md, TWO_BUILD_OPTIONS.md, README_BUILD.md | 3 |
| Implementation Status | VBA_IMPLEMENTATION_SUMMARY.md, PYTHON_IMPLEMENTATION_SUMMARY.md, WHATS_MISSING.md | 3 |
| Alternative Design | older-Options_Trend_Dashboard_Summary.md (if keeping newest as current) | 1 |
| Legacy Notes | VBA_IMPORT_SCRIPT.txt, build logs | 4 |
| **TOTAL ARCHIVE** | | **24 files** |

---

## 11. ORPHANED / POTENTIALLY MISSING FILES

### No .xlsm Files Found (Expected)
- Per `.gitignore`, Excel workbooks are not tracked
- `TrendFollowing_TradeEntry.xlsm` is built by `BUILD_WITH_PYTHON.bat`

### Incomplete Configuration
- `config/presets.json` only has 2 presets; CLAUDE.md lists 5 presets (mismatch)
- Could be seeded dynamically by VBA

### Missing Documentation
- No CHANGELOG.md (would be helpful for version history)
- No DEVELOPMENT.md (contributor guidelines)
- No API/function reference (VBA functions are documented inline)

---

## RECOMMENDATIONS

### Immediate Actions (High Priority)

1. **Move old batch files**:
   ```bash
   mkdir -p /home/kali/excel-trading-workflow/scripts/deprecated/
   mv /home/kali/excel-trading-workflow/old_*.bat scripts/deprecated/
   ```

2. **Remove VBA duplicates from Python/**:
   ```bash
   rm /home/kali/excel-trading-workflow/Python/TF_*.bas
   ```
   (Review TF_Presets_Enhanced.bas first to determine if it should replace TF_Presets.bas in VBA/)

3. **Create docs/ folder and reorganize**:
   ```bash
   mkdir -p /home/kali/excel-trading-workflow/docs/{setup,specifications,reference,archive}
   # Move files accordingly (see folder structure above)
   ```

4. **Archive outdated documentation**:
   Move 24 files from root to `docs/archive/` (see Summary Table, section 10)

5. **Update .gitignore** (optional):
   Add explicit pattern for logs:
   ```
   logs/*.log
   !logs/.gitkeep
   ```

### Medium Priority (Quality of Life)

1. **Create docs/README.md** (index):
   Lists where to find each type of documentation

2. **Consolidate setup guides**:
   - Merge VBA_IMPLEMENTATION_SUMMARY.md into VBA_SETUP_GUIDE.md
   - Merge PYTHON_IMPLEMENTATION_SUMMARY.md into PYTHON_SETUP_GUIDE.md
   - Move TROUBLESHOOTING_BUILD_ISSUES.md content into TROUBLESHOOTING.md in setup/ folder

3. **Create REFERENCE_OLDER_DASHBOARD.md**:
   Rename `older-Options_Trend_Dashboard_Summary.md` for clarity

4. **Archive build logs** (optional):
   Move logs to `logs/archive/` to keep root clean

### Low Priority (Future Enhancement)

1. Add CHANGELOG.md (document version history)
2. Add CONTRIBUTING.md (if collaborating)
3. Add function reference / API documentation
4. Auto-generate VBA documentation from comments

---

## FILE STATISTICS

| Metric | Value |
|--------|-------|
| Total files in project | ~115 |
| Markdown documentation files | 28 |
| VBA modules | 11 |
| Python scripts | 3 |
| Batch/VBS scripts | 11 |
| Build logs | 3 |
| Config files | 2 |
| Venv files (auto-generated) | ~1000+ |
| Recommended files to archive | 24 |
| Recommended files to keep | 38 |

