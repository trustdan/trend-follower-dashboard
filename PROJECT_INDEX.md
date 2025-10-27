# Excel Trading Workflow - Project Index

**Quick navigation to all project files and their purposes.**

---

## 🎯 START HERE (New Users)

1. **README.md** - Project overview and quick start
2. **USER_GUIDE.md** - 15,000-word beginner walkthrough ⭐
3. **UPDATED_QUICK_START.md** - 2-page quick reference

---

## 🤖 START HERE (AI Assistants)

1. **START_NEW_SESSION.md** - How to begin a new session
2. **.clinerules** - Critical project rules (auto-loaded)
3. **SESSION_2_SUMMARY.md** - Latest session recap
4. **DEVELOPMENT_LOG.md** - Technical context and issue tracker ⭐
5. **CHANGELOG.md** - Version history

---

## 📁 File Organization

### Build System
| File | Purpose | Run When |
|------|---------|----------|
| `BUILD.bat` | Main build command | First time + updates |
| `build_workbook_simple.py` | VBA import automation | Auto (via BUILD.bat) |
| `SCAN_FINVIZ.bat` | Standalone FINVIZ scanner | Optional (testing) |
| `run_finviz_scan.py` | CLI for FINVIZ scraper | Auto (via SCAN_FINVIZ.bat) |

### VBA Modules (Trading System Core)
| File | Lines | Purpose |
|------|-------|---------|
| `TF_Utils.bas` | 154 | Helper functions (sheets, tables, names) |
| `TF_Data.bas` | 320 | Data structure, heat, cooldowns |
| `TF_UI.bas` | 384 | Trading logic (evaluate, size, save) |
| `TF_Presets.bas` | 200 | FINVIZ integration + smart import |
| `TF_Python_Bridge.bas` | 280 | Python in Excel integration |
| `TF_UI_Builder.bas` | 300 | Automated UI generation + checkboxes |
| `Setup.bas` | 240 | One-click initialization + guide launcher |
| `ThisWorkbook.cls` | 100 | Auto-setup on workbook open |
| `Sheet_TradeEntry.cls` | 75 | TradeEntry sheet events |

**Total VBA:** 11 files, ~1,500 lines

### Python Modules (Optional Acceleration)
| File | Lines | Purpose |
|------|-------|---------|
| `finviz_scraper.py` | 280 | Active web scraping (NOT permalinks) |
| `heat_calculator.py` | 380 | Fast pandas calculations |
| `requirements.txt` | - | Python dependencies |

**Total Python:** 3 files, ~660 lines

### Documentation (For Users)
| File | Words | Purpose | Auto-Opens? |
|------|-------|---------|-------------|
| `USER_GUIDE.md` ⭐ | 15,000+ | Complete beginner walkthrough | Yes (first launch) |
| `README.md` | 2,000+ | Project overview | No |
| `UPDATED_QUICK_START.md` | 1,000+ | 2-page quick reference | No |
| `README_UPDATED.md` | 1,500+ | Feature summary | No |
| `START_HERE.md` | 3,000+ | Original detailed setup | No |
| `FINAL_SUMMARY.md` | 2,000+ | Project achievements | No |

### Documentation (For Developers/AI)
| File | Words | Purpose | Read When |
|------|-------|---------|-----------|
| `.clinerules` ⭐ | 5,000+ | Project rules + patterns | Every session (auto) |
| `DEVELOPMENT_LOG.md` ⭐ | 10,000+ | Technical issue tracker | Every session |
| `CHANGELOG.md` | 3,000+ | Version history | Every session |
| `SESSION_2_SUMMARY.md` | 2,000+ | Latest session recap | New session |
| `IMPLEMENTATION_STATUS.md` | 4,000+ | Technical architecture | As needed |
| `START_NEW_SESSION.md` | 500+ | Session start template | New session |
| `PROJECT_INDEX.md` | 1,000+ | This file (navigation) | Reference |

**Total Documentation:** 13 files, ~40,000+ words

---

## 🔍 Finding Things Quickly

### "Where is...?"

**...the user-facing documentation?**
→ USER_GUIDE.md (comprehensive) or UPDATED_QUICK_START.md (quick)

**...the technical context for AI assistants?**
→ DEVELOPMENT_LOG.md (issues, patterns) or .clinerules (rules)

**...the version history?**
→ CHANGELOG.md

**...the checkbox creation code?**
→ VBA/TF_UI_Builder.bas, CreateCheckboxes() function

**...the Python detection logic?**
→ VBA/TF_Python_Bridge.bas, IsPythonAvailable() function

**...the GO/NO-GO evaluation?**
→ VBA/TF_UI.bas, EvaluateChecklist() function

**...the position sizing calculations?**
→ VBA/TF_UI.bas, RecalcSizing() function

**...the FINVIZ scraping code?**
→ Python/finviz_scraper.py, fetch_finviz_tickers() function

**...the heat management logic?**
→ VBA/TF_Data.bas, PortfolioHeatAfter() and BucketHeatAfter()

**...the auto-setup code?**
→ VBA/ThisWorkbook.cls, Workbook_Open() event

**...the UI builder?**
→ VBA/TF_UI_Builder.bas, BuildTradeEntryUI() function

### "How do I...?"

**...set up the system for the first time?**
→ USER_GUIDE.md, "First Time Setup" section

**...understand what ATR, K, etc. mean?**
→ USER_GUIDE.md, "Understanding Each Field" section

**...trade my first position?**
→ USER_GUIDE.md, "Evaluating a Single Trade" section

**...fix Unicode encoding issues?**
→ DEVELOPMENT_LOG.md, Issue #1 or .clinerules, Rule #1

**...modernize Python in Excel code?**
→ DEVELOPMENT_LOG.md, Issue #2 or .clinerules, Rule #2

**...prevent duplicate buttons?**
→ DEVELOPMENT_LOG.md, Issue #3 or .clinerules, Rule #3

**...start a new session with AI?**
→ START_NEW_SESSION.md

**...check what changed in a version?**
→ CHANGELOG.md

**...troubleshoot common errors?**
→ USER_GUIDE.md, "Troubleshooting" section or DEVELOPMENT_LOG.md, "Known Issues"

---

## 📊 Project Statistics (v2.0.0)

### Code
- **Total Files:** 27
- **Total Lines:** ~6,000
- **VBA Code:** 11 modules, ~1,500 lines
- **Python Code:** 3 modules, ~660 lines
- **Documentation:** 13 files, ~40,000 words

### Functionality
- **Worksheets:** 8 (auto-created)
- **Tables:** 5 (auto-created)
- **Named Ranges:** 7
- **Buttons:** 6 (auto-created)
- **Dropdowns:** 4 (auto-created)
- **Checkboxes:** 6 (auto-created with fallback)
- **Presets:** 5 default FINVIZ queries
- **Buckets:** 6 correlation groups

### User Experience
- **Setup Time:** 3 minutes (from zero to trading)
- **Setup Steps:** 3 (Build → Open → Verify)
- **Manual Steps:** 0-1 (checkboxes only if automation fails)
- **Documentation Quality:** Comprehensive (15k word guide)

### Development
- **Sessions:** 2
- **Total Dev Time:** ~8 hours
- **Bugs Fixed:** 6 critical issues
- **Features:** 13 major features
- **Status:** ✅ Production Ready

---

## 🏗️ Architecture Overview

### Data Flow
```
User → TradeEntry UI → VBA Logic → Excel Tables → Trade Logged
         ↓
    Dropdowns (4) → Presets/Candidates/Buckets tables
         ↓
    Checkboxes (6) → Cells C20-C25 (TRUE/FALSE)
         ↓
    Evaluate Button → EvaluateChecklist() → Banner (GREEN/YELLOW/RED)
         ↓
    Recalc Sizing → RecalcSizing() → Position size calculated
         ↓
    Save Decision → SaveDecision() → 5 hard gates → Decisions table
```

### Python Integration (Optional)
```
User → Import Candidates Button → IsPythonAvailable()?
                                        ↓
                        Yes: CallPythonFinvizScraper() → Auto-scrape FINVIZ
                        No: ImportManual() → User pastes tickers
                                        ↓
                                 Candidates table populated
```

### Heat Management
```
User → Recalc Sizing → Calculate new position R
                             ↓
                 Read Positions table (current open trades)
                             ↓
                 Calculate PortfolioHeatAfter(newR)
                 Calculate BucketHeatAfter(newR, bucket)
                             ↓
                 Check against caps (4%, 1.5%)
                             ↓
                 Display on UI (before/after preview)
                             ↓
User → Save Decision → Check heat caps (hard gate #5)
                             ↓
                 If exceeds → BLOCK trade
                 If OK → Log to Decisions table
```

---

## 🎯 Current Status

**Version:** v2.0.0
**Release Date:** 2025-01-27
**Status:** ✅ Production Ready

**What Works:**
- ✅ Full trading system (6-item checklist)
- ✅ FINVIZ scraping (manual + Python)
- ✅ Position sizing (3 methods)
- ✅ Heat management (portfolio + bucket)
- ✅ Automated setup
- ✅ USER_GUIDE.md auto-opens
- ✅ Checkbox auto-creation (with fallback)

**Known Limitations:**
- ⚠️ Checkboxes: ~70% auto-creation rate (manual fallback available)
- ⚠️ Python: Requires Microsoft 365 Insider (optional feature)
- ⚠️ Platform: Windows Excel only (VBA limitation)

**Next Steps:**
- User testing
- Feedback incorporation
- Additional features (as requested)

---

## 📞 Support Resources

**For Users:**
1. USER_GUIDE.md - Comprehensive walkthrough
2. Setup sheet (in workbook) - Built-in help
3. README.md - Quick start

**For AI Assistants:**
1. .clinerules - Critical rules
2. DEVELOPMENT_LOG.md - Technical context
3. CHANGELOG.md - Version history

**For Issues:**
1. Check USER_GUIDE.md → Troubleshooting
2. Check DEVELOPMENT_LOG.md → Known Issues
3. Check CHANGELOG.md → Recent fixes

---

## 🗺️ Roadmap

### Completed (v2.0.0)
- ✅ Auto-setup on first open
- ✅ Checkbox auto-creation
- ✅ USER_GUIDE.md (comprehensive docs)
- ✅ Unicode fixes
- ✅ Python detection modernization

### Planned (v2.1.0+)
- [ ] ActiveX checkboxes (more reliable)
- [ ] Real-time price updates via API
- [ ] Automatic Positions table sync from broker
- [ ] Performance analytics dashboard
- [ ] Email/SMS alerts
- [ ] Multi-account support

---

**Last Updated:** 2025-01-27
**Maintained By:** AI Assistants + User
**License:** Private Use

**For navigation help, see this file (PROJECT_INDEX.md)**
**For quick start, see README.md or USER_GUIDE.md**
**For new sessions, see START_NEW_SESSION.md**
