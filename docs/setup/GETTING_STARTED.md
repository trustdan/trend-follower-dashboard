# Getting Started - Quick Reference

**Choose your path and get up and running in 60-90 minutes.**

---

## 🎯 What You're Building

An Excel-based trading system that:
- ✅ Enforces mechanical entry rules (6-item checklist)
- ✅ Sizes positions automatically (3 methods)
- ✅ Manages risk via heat caps
- ✅ Prevents impulsive trades (2-minute brake)
- ✅ Logs full audit trail

---

## 📁 File Organization

Your project folder structure:

```
excel-trading-workflow/
│
├── 📂 VBA/                        ← VBA modules (import these)
│   ├── TF_Utils.bas
│   ├── TF_Data.bas
│   ├── TF_UI.bas
│   ├── TF_Presets.bas
│   ├── ThisWorkbook.cls
│   └── Sheet_TradeEntry.cls
│
├── 📂 Python/                     ← Python modules (optional)
│   ├── finviz_scraper.py
│   ├── heat_calculator.py
│   ├── TF_Python_Bridge.bas
│   ├── TF_Presets_Enhanced.bas
│   └── requirements.txt
│
├── 📜 IMPORT_VBA_MODULES.bat     ← Double-click to auto-import (Windows)
├── 📜 import_to_excel.py         ← Python automation script
│
├── 📖 VBA_SETUP_GUIDE.md         ← START HERE for VBA
├── 📖 PYTHON_SETUP_GUIDE.md      ← Add Python later (optional)
├── 📖 README.md                  ← System overview
│
└── (other documentation files)
```

---

## 🚀 Quick Start Paths

### **Path 1: VBA Only** ⭐ RECOMMENDED FOR BEGINNERS

**Time:** 60-90 minutes
**Works on:** Excel 2016+, Windows/Mac
**Best for:** Getting started, learning the system

**Steps:**
1. **Import VBA modules** (choose one):
   - **Windows:** Double-click `IMPORT_VBA_MODULES.bat` (automated)
   - **Manual:** Follow `VBA_SETUP_GUIDE.md` Part 1 (10 min)

2. **Run setup macro:**
   - Alt+F11 → Ctrl+G (Immediate Window)
   - Type: `EnsureStructure` → Press Enter
   - Verify: 8 sheets created

3. **Build TradeEntry UI:**
   - Follow `VBA_SETUP_GUIDE.md` Part 2 (30 min)
   - Add labels, buttons, checkboxes, dropdowns

4. **Test workflow:**
   - Follow `VBA_SETUP_GUIDE.md` Part 3 (10 min)
   - Run sample trade from start to finish

5. **Done!** Start trading

---

### **Path 2: VBA + Python** ⚡ ADVANCED

**Time:** 120-150 minutes (includes VBA setup)
**Works on:** Microsoft 365 Insider (Windows)
**Best for:** Auto-scraping FINVIZ, faster performance

**Steps:**
1. **Complete Path 1** (VBA setup above)

2. **Enable Python in Excel:**
   - Data tab → Python → Enable
   - Follow `PYTHON_SETUP_GUIDE.md` Part 1-2 (30 min)

3. **Import Python modules:**
   - Copy `finviz_scraper.py` and `heat_calculator.py` to workbook folder
   - Import `TF_Python_Bridge.bas` into VBA
   - Install dependencies: `pip install requests beautifulsoup4 pandas`

4. **Test integration:**
   - Run `TestPythonIntegration` macro
   - Verify scraper and heat calculator work

5. **(Optional) Replace TF_Presets:**
   - Import `TF_Presets_Enhanced.bas`
   - Delete old `TF_Presets` module
   - Enjoy auto-scraping!

---

## 📖 Documentation Map

**Start here based on your goal:**

| Goal | Read This | Time |
|------|-----------|------|
| Import VBA modules | `VBA_SETUP_GUIDE.md` Part 1 | 10 min |
| Build TradeEntry UI | `VBA_SETUP_GUIDE.md` Part 2 | 30 min |
| Test the system | `VBA_SETUP_GUIDE.md` Part 3 | 10 min |
| Understand architecture | `VBA_README.md` | 15 min |
| Add Python auto-scraping | `PYTHON_SETUP_GUIDE.md` | 60 min |
| Troubleshoot issues | `VBA_README.md` Troubleshooting | varies |
| Learn trading rules | `workflow-plan.md` | 20 min |
| System overview | `README.md` | 10 min |

---

## ⚙️ Automation Scripts

### **Windows Users:**

**Option 1: Batch Script (No Python required)**
```batch
IMPORT_VBA_MODULES.bat
```
- Double-click in Windows Explorer
- Shows manual instructions if Python not available

**Option 2: Python Script (Automated)**
```bash
pip install pywin32
python import_to_excel.py
```
- Opens Excel
- Creates workbook
- Imports all VBA modules automatically

### **Mac/Linux Users:**

**Manual import only** (VBA editor method)
- Follow `VBA_SETUP_GUIDE.md` Method B
- File → Import File for each module

---

## 🧪 Testing Checklist

After setup, verify these work:

### **VBA Tests:**
- [ ] 8 sheets created (TradeEntry, Presets, Buckets, Candidates, Decisions, Positions, Summary, Control)
- [ ] 5 tables created (tblPresets, tblBuckets, tblCandidates, tblDecisions, tblPositions)
- [ ] 7 named ranges defined (Equity_E, RiskPct_r, etc.)
- [ ] Dropdowns work (Preset, Ticker, Sector, Bucket)
- [ ] Checkboxes link to cells (C20:C25)
- [ ] Evaluate button shows GREEN banner (when all boxes ticked)
- [ ] Recalc Sizing computes outputs
- [ ] Save Decision logs to Decisions + Positions tables

### **Python Tests** (if using):
- [ ] `TestPythonIntegration` runs without errors
- [ ] FINVIZ scraper returns tickers (5-10 sec)
- [ ] Heat calculator returns valid numbers
- [ ] Import button offers "Auto-scrape or manual paste?"

---

## 🎓 Learning Sequence

### **Day 1: Setup (1-2 hours)**
1. Import VBA modules (automated or manual)
2. Run `EnsureStructure` macro
3. Build TradeEntry UI (follow cell reference map)
4. Test with sample trade

### **Day 2: Configuration (30 min)**
5. Update Summary settings (account size, risk %)
6. Add test candidates to Candidates table
7. Run 3-5 test trades
8. Review Decisions table

### **Week 1: Paper Trading**
9. Import real candidates from FINVIZ
10. Validate signals on TradingView
11. Run full workflow 10+ times
12. Track adherence to GREEN-only rule

### **Week 2: Optimization**
13. Adjust heat caps based on comfort level
14. Customize bucket cooldown parameters
15. (Optional) Add Python integration
16. Go live with small size

---

## 🔧 Customization Priorities

**First:** Adjust these in Summary sheet to match your account

| Setting | Location | Your Value |
|---------|----------|------------|
| Equity_E | Summary!B2 | **$______** (your account size) |
| RiskPct_r | Summary!B3 | **0.____** (0.5-1.0% typical) |
| HeatCap_H_pct | Summary!B5 | **0.____** (2-6% portfolio cap) |

**Later:** Fine-tune these based on results

| Setting | Default | Adjust If... |
|---------|---------|--------------|
| StopMultiple_K | 2.0 | Too many stop-outs → increase to 2.5-3.0 |
| BucketHeatCap_pct | 1.5% | Want more/less sector concentration |
| StopoutsToCooldown | 2 | Buckets cooling down too often → increase to 3 |

---

## 💡 Pro Tips

### **Before You Start:**
1. ✅ Read `workflow-plan.md` to understand the trading rules
2. ✅ Watch a few TradingView tutorials (Donchian channels, ATR)
3. ✅ Paper trade for 1-2 weeks before going live

### **During Setup:**
1. ✅ Use automated import if on Windows (saves 5-10 min)
2. ✅ Copy cell reference map for TradeEntry UI (don't guess placements)
3. ✅ Test each button immediately after creating it

### **After Setup:**
1. ✅ Run `UpdateCooldowns` weekly (or after stop-outs)
2. ✅ Back up workbook before major changes
3. ✅ Review Decisions table monthly for patterns
4. ✅ Keep Summary settings conservative until confident

### **Common Mistakes to Avoid:**
1. ❌ Skipping the 2-minute timer (defeats impulse control)
2. ❌ Overriding YELLOW/RED banners (defeats checklist)
3. ❌ Not importing candidates (Save Decision will block)
4. ❌ Setting risk % too high (start at 0.5%, not 2%)
5. ❌ Ignoring cooldowns (buckets exist for a reason)

---

## 🐛 Quick Troubleshooting

**"Compile error: Sub or Function not defined"**
→ Import all 4 VBA modules (TF_Utils, TF_Data, TF_UI, TF_Presets)

**"Cannot save: Ticker not in today's Candidates"**
→ Click "Import Candidates" first, add at least one ticker

**"Cannot save: 2-minute cool-off not elapsed"**
→ Wait 2 full minutes after clicking Evaluate (or test with earlier timestamp)

**"Buttons don't do anything"**
→ Right-click button → Assign Macro → select correct procedure

**"Dropdowns show #REF!"**
→ Run `EnsureStructure` again to recreate tables

**"Python not available"**
→ Microsoft 365 Insider required. See `PYTHON_SETUP_GUIDE.md`

---

## 📞 Where to Get Help

1. **Setup Issues:** `VBA_SETUP_GUIDE.md` Troubleshooting section
2. **Button/Macro Errors:** `VBA_README.md` → Module Layout
3. **Trading Questions:** `workflow-plan.md` → Position Sizing Logic
4. **Python Issues:** `PYTHON_SETUP_GUIDE.md` → Troubleshooting
5. **System Architecture:** `README.md` or `CLAUDE.md`

---

## ✅ Success Checklist

You're ready to trade when:
- [ ] VBA modules imported (4 standard + 2 class)
- [ ] `EnsureStructure` ran successfully (8 sheets created)
- [ ] TradeEntry UI built (20+ elements)
- [ ] Test trade completes: Import → Evaluate → Size → Save
- [ ] Decisions table has at least 1 row
- [ ] Positions table has at least 1 row
- [ ] Summary settings customized to your account
- [ ] Understand GO/NO-GO logic (GREEN only!)
- [ ] 2-minute timer tested (blocks early saves)
- [ ] Heat caps tested (blocks over-cap trades)

---

## 🎯 First Trade Walkthrough

**Follow this exactly to verify everything works:**

1. **Import candidates:**
   - Click "Import Candidates" button
   - Paste: `MSFT, AAPL, GOOGL`
   - Verify added to Candidates table with today's date

2. **Fill inputs:**
   - Preset: TF_BREAKOUT_LONG
   - Ticker: MSFT (from dropdown)
   - Sector: Technology
   - Bucket: Tech/Comm (auto-fills)
   - Entry: 420.00
   - ATR N: 1.20
   - K: 2
   - Method: Stock

3. **Complete checklist:**
   - Tick all 6 checkboxes

4. **Evaluate:**
   - Click "Evaluate" button
   - Verify: GREEN banner appears
   - Note: Timer started (shows in Control!A1)

5. **Size position:**
   - Click "Recalc Sizing" button
   - Verify outputs: R=$75, Shares=31, etc.

6. **Save decision:**
   - Wait 2 minutes (or cheat: edit Control!A1 to earlier time)
   - Click "Save Decision" button
   - Verify: Success message, banner clears

7. **Check logs:**
   - Decisions sheet: 1 new row with all details
   - Positions sheet: 1 new row (MSFT, Open)

8. **Done!** System works. Now customize and use for real.

---

**Estimated Time to Complete This Guide:** 15 minutes

**Estimated Time to Full System:** 60-150 minutes (depending on path)

**Start Here:** `VBA_SETUP_GUIDE.md` Part 1

**Good Luck! 🚀**
