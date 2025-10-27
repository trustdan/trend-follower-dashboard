# Excel Trading Workflow - Files at a Glance

## Quick Navigation

### Start Here
1. **README.md** - System overview and philosophy
2. **QUICK_START.md** - Build the workbook in one command
3. **GETTING_STARTED.md** - First-time setup walkthrough
4. **CLAUDE.md** - Project guidelines for Claude AI

### Build & Run
- **BUILD_WITH_PYTHON.bat** - Execute this (primary build script)
- **import_to_excel.py** - Automation engine (called by BUILD_WITH_PYTHON.bat)
- **COMPLETE_SETUP.bat** - Alternative build flow

### Development
- **VBA/** folder - 11 VBA modules (all active, ready to import)
- **Python/** folder - 3 Python scripts (optional enhancement)
- **scripts/** folder - Helper scripts and utilities

---

## File Organization Status

### ACTIVE (Keep)
```
Root Level:
  README.md
  CLAUDE.md
  GETTING_STARTED.md
  QUICK_START.md
  BUILD_WITH_PYTHON.bat
  COMPLETE_SETUP.bat
  import_to_excel.py

Folders:
  VBA/                (11 modules)
  Python/             (3 scripts)
  scripts/            (5 helper scripts)
  config/             (presets.json)
  venv/               (Python environment)
  logs/               (build logs)
```

### NEEDS CLEANUP
```
Root Level (28 .md files - too many):
  - 4 essential (keep above)
  - 5 specification docs (move to docs/specifications/)
  - 2 setup guides (move to docs/setup/)
  - 3 reference docs (move to docs/reference/)
  - 24 archive docs (move to docs/archive/)

Old Batch Files (in root):
  - old_CLEANUP_STUCK_EXCEL.bat
  - old_IMPORT_VBA_MODULES.bat
  - old_IMPORT_VBA_MODULES_DEBUG.bat
  - old_VERIFY_MODULES.bat
  → Move to scripts/deprecated/

Python Folder (contains misplaced VBA):
  - Python/TF_Python_Bridge.bas (duplicate)
  - Python/TF_Presets_Enhanced.bas (enhanced version)
  → Remove (keep only VBA/ originals)

Build Logs (in logs/):
  - 3 old build logs from Oct 26
  → Move to logs/archive/
```

---

## Recommended Structure (After Cleanup)

```
excel-trading-workflow/
├── README.md                    ← START HERE
├── CLAUDE.md                    ← For Claude AI
├── GETTING_STARTED.md           ← Setup guide
├── QUICK_START.md               ← Build command
│
├── BUILD_WITH_PYTHON.bat        ← Run this
├── COMPLETE_SETUP.bat           ← Alt build
├── import_to_excel.py           ← Automation
│
├── VBA/                         ← 11 active modules
├── Python/                      ← 3 scripts (no .bas files)
├── scripts/                     ← Helper scripts
│   └── deprecated/              ← 4 old .bat files
├── config/                      ← Configuration
├── logs/                        ← Current build logs
│   └── archive/                 ← Old logs
├── docs/                        ← NEW: organized docs
│   ├── README.md                (Index)
│   ├── setup/                   (Setup guides)
│   ├── specifications/          (Design docs)
│   ├── reference/               (Architecture)
│   └── archive/                 (Development history)
│
├── venv/                        ← Python environment
└── .gitignore
```

---

## File Statistics

| Category | Count | Status |
|----------|-------|--------|
| Essential docs in root | 4 | Keep |
| Total markdown files | 28 | 24 need archiving |
| VBA modules | 11 | All active |
| Python scripts | 3 | All active |
| Build scripts | 11 | 4 deprecated |
| Old batch files | 4 | Archive |
| VBA duplicates | 2 | Remove from Python/ |
| Build logs | 3 | Archive |

---

## Quick Cleanup Command

```bash
# Move old batch files
mkdir -p scripts/deprecated
mv old_*.bat scripts/deprecated/

# Clean VBA duplicates
rm Python/TF_*.bas

# Organize docs
mkdir -p docs/{setup,specifications,reference,archive}
# (Then move files - see CLEANUP_CHECKLIST.md)
```

---

## What Each Doc Category Contains

### Core Docs (in root - keep)
- **README.md** - What is this system?
- **CLAUDE.md** - How to use it with Claude AI
- **GETTING_STARTED.md** - How to set it up
- **QUICK_START.md** - How to build it

### Specification Docs (move to docs/specifications/)
- newest-Interactive_TF_Workbook_Plan.md - Detailed workbook design
- workflow-plan.md - Trading rules and logic
- SeykotaTurtleTrend-*.md - TradingView strategy guide
- diversification-across-sectors.* - Sector strategy

### Setup Guides (move to docs/setup/)
- VBA_SETUP_GUIDE.md - Manual VBA import steps
- PYTHON_SETUP_GUIDE.md - Python integration guide
- TROUBLESHOOTING.md - Common issues and fixes

### Reference (move to docs/reference/)
- VBA_README.md - VBA architecture overview
- Workbook_Readme_Text.md - Workbook structure guide

### Archive (move to docs/archive/)
- BUILD_NOW_COMPLETE.md, FINAL_FIXES.md, etc. - Build history
- WHATS_MISSING.md, FIXES_APPLIED.md, etc. - Development logs
- VBA_IMPLEMENTATION_SUMMARY.md, etc. - Implementation status
- older-Options_Trend_Dashboard_Summary.md - Alternative design

---

## Active Code You Care About

### VBA (11 modules in VBA/)
1. **TF_Utils.bas** - Helper utilities
2. **TF_Data.bas** - Data structure & calculations
3. **TF_UI.bas** - User interface logic
4. **TF_Presets.bas** - FINVIZ integration
5. **TF_UI_Builder.bas** - Auto-builds UI
6. **Setup.bas** - Initialization
7. **Python_Run.bas** - Python bridge
8. **PQ_Setup.bas** - Power Query (optional)
9. **ThisWorkbook.cls** - Workbook events
10. **Sheet_TradeEntry.cls** - Sheet events

### Python (3 scripts)
1. **import_to_excel.py** - VBA import automation
2. **finviz_scraper.py** - FINVIZ data scraper
3. **heat_calculator.py** - Portfolio calculations

### Build Tools
- **BUILD_WITH_PYTHON.bat** - Run this to build everything
- **import_to_excel.py** - Automation script
- **COMPLETE_SETUP.bat** - Alternative approach

---

## Decision Tree: Which Doc to Read?

```
Q: I want to understand the system
A: Start with README.md

Q: I want to build the workbook now
A: Run BUILD_WITH_PYTHON.bat or read QUICK_START.md

Q: I'm setting up for the first time
A: Read GETTING_STARTED.md

Q: I'm working with Claude AI
A: Read CLAUDE.md

Q: I want to understand the workbook design
A: Read newest-Interactive_TF_Workbook_Plan.md

Q: I want to understand the trading rules
A: Read workflow-plan.md

Q: I need to manually import VBA
A: Read VBA_SETUP_GUIDE.md

Q: I want to see the current build status
A: Read BUILD_NOW_COMPLETE.md + FINAL_FIXES.md

Q: I want to understand VBA architecture
A: Read VBA_README.md

Q: I want TradingView strategy details
A: Read SeykotaTurtleTrend-FollowingOptionsExecution+TradingViewStrategyGuide.md
```

---

## One-Minute Summary

This is an **Excel trading system** with:
- **VBA backend** (11 modules handling logic)
- **Python integration** (optional; auto-scraping & calculations)
- **Build automation** (BUILD_WITH_PYTHON.bat does everything)
- **Heavy documentation** (28 .md files, needs organizing)

To use it:
1. Run `BUILD_WITH_PYTHON.bat`
2. Open the generated workbook
3. Follow QUICK_START.md if any issues

To improve it:
1. See FILE_INVENTORY.md for detailed analysis
2. See CLEANUP_CHECKLIST.md for action items
3. See ORGANIZATION_SUMMARY.txt for overview

