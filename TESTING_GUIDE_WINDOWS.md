# TF-Engine Windows Testing Guide

**Quick Reference for Testing Fixed Installer**

**Date:** 2025-10-29
**Version:** 1.0.0
**Fix:** Navigation bug (`e.subscribe is not a function`)

---

## 🔴 THE PROBLEM

**Symptom:** Clicking any navigation link (Scanner, Checklist, Sizing, Heat, Entry, Calendar) throws:
```
Uncaught TypeError: e.subscribe is not a function
```

**Impact:** Complete navigation failure - can't use the app

**Root Cause:** `ui/src/lib/components/Navigation.svelte:67` - Manual `page.subscribe()` call fails in prerendered Svelte app

**Fix Applied:** Changed to reactive `$page` store with null safety

---

## ✅ QUICK TEST (5 minutes)

### 1. Uninstall Old Version

```
Windows Settings → Apps → TF-Engine → Uninstall
When asked: "Delete database?" → Choose "No"
```

### 2. Verify New Installer

```powershell
cd C:\Users\Dan\trend-follower-dashboard\installer
Get-FileHash TF-Engine-Setup-v1.0.0.exe -Algorithm SHA256
```

**Expected:** `28A8026B882CF6F6C8FB68F8E279EF2FB74E7FE33400729E8C4A8185704B41EB`

### 3. Install Fixed Version

```
Right-click TF-Engine-Setup-v1.0.0.exe → Run as administrator
Follow installer wizard
Launch from desktop shortcut
```

### 4. Test Navigation ⭐

```
1. Browser opens to http://localhost:8080
2. Press F12 (open DevTools console)
3. Click EACH navigation link:
   - Dashboard
   - Scanner ← Should work now!
   - Checklist ← Should work now!
   - Sizing ← Should work now!
   - Heat ← Should work now!
   - Entry ← Should work now!
   - Calendar ← Should work now!
```

**Success Criteria:**
- ✅ No `e.subscribe is not a function` error
- ✅ All pages load
- ✅ Navigation works smoothly
- ✅ Console shows no red errors (warnings OK)

**If Still Broken:**
- Note exact error message
- Check console output
- Try different browser
- Report back

---

## 📋 FULL VALIDATION (After Navigation Works)

If navigation test passes, proceed with complete workflow validation:

### Phase 1: Configure Settings (5 min)

```
Navigate to Settings
Enter:
  Equity: 100000
  Risk %: 0.75
  Portfolio Cap %: 4.0
  Bucket Cap %: 1.5
  Max Units: 4
Click "Save Settings"
```

### Phase 2: Scanner Test (5 min)

```
Navigate to Scanner
Click "Run Daily FINVIZ Scan"
(Wait 10-20 seconds)
Verify candidates appear
```

### Phase 3: Checklist Test (10 min) ⭐ CRITICAL

```
Navigate to Checklist
Enter:
  Ticker: AAPL
  Entry: 180.50
  ATR (N): 2.35
  Sector: Tech/Comm
  Structure: Stock

Watch banner at top:
  1. Should start RED (no gates checked)
  2. Check all 5 required gates → Banner turns YELLOW
  3. Check 3+ quality items → Banner turns GREEN

Click "Save Evaluation"
→ 2-minute timer should start (2:00, 1:59, 1:58...)
```

### Phase 4: Position Sizing (5 min)

```
Navigate to Position Sizing
Form should be pre-filled from checklist
K multiple: 2.0
Click "Calculate Position Size"
→ Should show shares, risk, stop, add-on schedule
Click "Save Position Plan"
```

### Phase 5: Heat Check (5 min)

```
Navigate to Heat Check
Should show:
  Current Portfolio Heat: $0
  Portfolio Cap: $4,000
Click "Check Heat for This Trade"
→ Should show WITHIN CAP (green)
```

### Phase 6: Gates Validation (10 min) ⭐ THE CRITICAL TEST

```
WAIT for 2-minute timer to reach 0:00
Navigate to Trade Entry
Click "Run Final Gate Check"

Expected:
  ✅ Gate 1: Banner GREEN
  ✅ Gate 2: Impulse brake (2+ min)
  ✅ Gate 3: No cooldown
  ✅ Gate 4: Heat caps OK
  ✅ Gate 5: Sizing complete

"SAVE GO DECISION" button should be ENABLED (green)
Click "SAVE GO DECISION"
→ Success notification: "✓ GO decision saved for AAPL"
```

### Phase 7: Calendar View (2 min)

```
Navigate to Calendar
→ Should see AAPL in Tech/Comm row, current week
Hover over cell → Shows entry, risk, date
```

### Phase 8: Theme Toggle (2 min)

```
Click sun/moon icon (top right)
→ UI switches to dark mode
Navigate through pages → Theme persists
Toggle back to light mode
```

### Phase 9: Restart Test (5 min)

```
Stop server (Ctrl+C in console)
Close browser
Launch again from desktop shortcut
→ All data should persist (settings, candidates, GO decision)
```

---

## 🐛 TROUBLESHOOTING

### If Navigation Still Fails

**Check:**
1. Is checksum correct? (28a80...)
2. Did you uninstall old version first?
3. Clear browser cache (Ctrl+Shift+Delete)
4. Try different browser (Chrome, Edge, Firefox)

**Report:**
- Exact error message
- Browser console output (screenshot)
- Which page failed

### If Backend Issues

**Check Console Window:**
```
Should show:
  [TF-Engine] Server listening on http://127.0.0.1:8080
  [TF-Engine] Embedded UI loaded successfully
```

**If Database Errors:**
```powershell
# Check database exists
Test-Path "$env:APPDATA\TF-Engine\trading.db"

# If not, reinitialize
cd "C:\Program Files\TF-Engine"
.\tf-engine.exe init
```

---

## 📊 SUCCESS CHECKLIST

After testing, confirm:

- [ ] Installer runs without errors
- [ ] Navigation works (no `e.subscribe` errors)
- [ ] Settings save and persist
- [ ] Scanner finds candidates
- [ ] Checklist banner transitions (RED → YELLOW → GREEN)
- [ ] 2-minute timer works
- [ ] Position sizing calculates correctly
- [ ] Heat check validates
- [ ] All 5 gates pass
- [ ] GO decision saves
- [ ] Calendar displays trade
- [ ] Theme toggle works
- [ ] Data persists after restart

---

## 📁 REFERENCE DOCUMENTS

- **Bug Details:** `docs/FRONTEND_BUG_SUMMARY.md`
- **LLM History:** `docs/LLM-Update.md` (see Oct 29 21:00 session)
- **Step 28 Plan:** `plans/phase5-step28-final-validation.md`
- **Step 27 Complete:** `docs/STEP27_COMPLETE.md`

---

## 🎯 WHAT TO REPORT BACK

### If Navigation Works:

```
✅ Navigation fixed!
- All 7 pages load
- No console errors
- Proceeding with full validation
```

### If Navigation Still Broken:

```
❌ Still broken
- Error: [exact message]
- Browser: [Chrome/Edge/Firefox]
- Console screenshot: [attach]
```

### After Full Validation:

```
Report which tests passed/failed:
- Settings: Pass/Fail
- Scanner: Pass/Fail
- Checklist: Pass/Fail
- Sizing: Pass/Fail
- Heat: Pass/Fail
- Gates: Pass/Fail
- Calendar: Pass/Fail
- Theme: Pass/Fail
- Persistence: Pass/Fail
```

---

## 🚀 NEXT STEPS

### If All Tests Pass:

1. Create release documentation
2. Tag v1.0.0
3. Declare PRODUCTION READY
4. Celebrate! 🎉

### If Issues Found:

1. Document specific failure
2. Return to Linux for fixes
3. Rebuild → Re-test
4. Repeat until clean

---

**Good luck with testing!** 🍀

The fix should work - the root cause was clearly identified and properly addressed. Let me know the results!
