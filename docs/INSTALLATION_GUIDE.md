# TF-Engine Installation Guide

**Version:** 1.0.0
**Last Updated:** 2025-10-29

**TF = Trend Following** - Systematic trading discipline enforcement system

This guide covers installation, initial setup, and verification for Windows users.

---

## Table of Contents

1. [System Requirements](#system-requirements)
2. [Download](#download)
3. [Installation Methods](#installation-methods)
4. [Initial Setup](#initial-setup)
5. [Verification](#verification)
6. [Troubleshooting Installation](#troubleshooting-installation)
7. [Uninstallation](#uninstallation)
8. [Upgrading](#upgrading)

---

## System Requirements

### Minimum Requirements

- **Operating System:** Windows 10 (64-bit) or Windows 11
- **Processor:** Intel Core i3 or equivalent (2.0 GHz+)
- **RAM:** 512 MB available
- **Disk Space:** 100 MB free space
- **Internet:** Required for FINVIZ scanning and TradingView integration
- **Browser:** Chrome, Firefox, or Edge (modern versions)

### Recommended Requirements

- **Operating System:** Windows 10/11 (64-bit, latest updates)
- **Processor:** Intel Core i5 or equivalent (2.5 GHz+)
- **RAM:** 1 GB+ available
- **Disk Space:** 500 MB free space (for database growth)
- **Monitor:** 1920×1080 resolution or higher
- **Setup:** Dual monitors (TradingView on one, TF-Engine on other)
- **Browser:** Chrome (latest version) - best performance

### Software Dependencies

**None required!** TF-Engine is a standalone executable with no external dependencies.

- ✗ No Python required
- ✗ No Node.js required
- ✗ No .NET Framework required
- ✗ No Visual C++ redistributables required
- ✓ Pure Go binary (self-contained)

---

## Download

### Official Release

**Download Location:** [Specify download URL]

**Available Packages:**

1. **Standalone Binary** (Recommended for most users)
   - File: `tf-engine.exe`
   - Size: ~15-20 MB
   - No installer needed
   - Run from any folder

2. **Windows Installer** (For system-wide installation)
   - File: `TF-Engine-Setup-v1.0.0.msi` or `TF-Engine-Setup-v1.0.0.exe`
   - Size: ~20-25 MB
   - Installs to Program Files
   - Creates Start Menu shortcuts
   - Desktop shortcut option

### Verify Download (Optional but Recommended)

**Check file hash to ensure authentic download:**

```powershell
# Open PowerShell
# Navigate to download folder
cd C:\Users\YourName\Downloads

# Calculate SHA256 hash
Get-FileHash tf-engine.exe -Algorithm SHA256

# Compare to official hash (provided on download page)
# Should match: [hash-value-here]
```

**If hash doesn't match:**
- Download might be corrupted
- Re-download from official source
- Do NOT run the file

---

## Installation Methods

### Method 1: Standalone Binary (Recommended)

**Best for:**
- Users who want portability
- Users without admin rights
- Users who prefer manual control

**Steps:**

1. **Download** `tf-engine.exe`

2. **Choose installation folder:**
   ```
   Recommended locations:
   - C:\TF-Engine\
   - C:\Users\[YourName]\Documents\TF-Engine\
   - C:\Users\[YourName]\AppData\Local\TF-Engine\
   ```

3. **Create folder:**
   ```
   Open File Explorer
   Navigate to desired location
   Right-click → New → Folder
   Name: TF-Engine
   ```

4. **Move binary:**
   ```
   Move tf-engine.exe from Downloads to C:\TF-Engine\
   ```

5. **Create desktop shortcut (optional):**
   ```
   Right-click tf-engine.exe
   Send to → Desktop (create shortcut)
   ```

6. **First launch:**
   ```
   Double-click tf-engine.exe
   (Or double-click desktop shortcut)

   Expected:
   - Command window opens (server console)
   - Browser opens to http://localhost:8080
   - Dashboard screen appears
   ```

**Done!** Skip to [Initial Setup](#initial-setup).

---

### Method 2: Windows Installer (.msi or .exe)

**Best for:**
- Users who want traditional installation
- Users with admin rights
- Users who want Start Menu integration

**Steps:**

1. **Download** `TF-Engine-Setup-v1.0.0.msi` (or `.exe`)

2. **Run installer:**
   ```
   Double-click TF-Engine-Setup-v1.0.0.msi
   Windows may show "SmartScreen" warning (new app)
   Click "More info" → "Run anyway"
   ```

3. **Installation wizard:**

   **Step 1: Welcome screen**
   - Click "Next"

   **Step 2: License agreement**
   - Read license
   - Accept if agreeable
   - Click "Next"

   **Step 3: Installation location**
   ```
   Default: C:\Program Files\TF-Engine\
   (Recommended: Keep default)

   Or: Click "Browse" to change location
   ```
   - Click "Next"

   **Step 4: Shortcuts**
   ```
   Options:
   - [✓] Create Desktop shortcut
   - [✓] Create Start Menu folder
   ```
   - Click "Next"

   **Step 5: Install**
   - Click "Install"
   - UAC prompt: Click "Yes" (requires admin)
   - Installation progress bar (10-30 seconds)

   **Step 6: Complete**
   - [✓] Launch TF-Engine
   - Click "Finish"

4. **First launch:**
   ```
   Expected:
   - Command window opens (server console)
   - Browser opens to http://localhost:8080
   - Dashboard screen appears
   ```

**Done!** Proceed to [Initial Setup](#initial-setup).

---

## Initial Setup

After installation, complete these one-time setup steps.

### Step 1: Initialize Database

**On first launch, TF-Engine automatically creates database at:**
```
C:\Users\[YourName]\AppData\Roaming\TF-Engine\trading.db
```

**Verify database created:**

1. Open File Explorer
2. Navigate to: `%APPDATA%\TF-Engine`
   - (Paste `%APPDATA%\TF-Engine` in address bar)
3. Should see: `trading.db` file

**If database not created:**

```cmd
# Open Command Prompt
# Navigate to tf-engine.exe location
cd C:\TF-Engine

# Run init command
tf-engine.exe init

# Expected output:
# Database initialized at: C:\Users\...\trading.db
# Tables created successfully
```

---

### Step 2: Configure Account Settings

1. **Open TF-Engine**
   - Browser should be at http://localhost:8080

2. **Navigate to Settings**
   - Click gear icon (⚙️) in header (top-right)

3. **Fill Account Settings:**

   ```
   Equity: 100000
   (Your trading account size in dollars)

   Risk % per unit: 0.75
   (Percent of equity risked per position unit)
   (Range: 0.50% - 1.00% typical)
   (Ed Seykota uses 1.00%)

   Portfolio heat cap: 4.0
   (Maximum total risk across all positions, in % of equity)

   Bucket heat cap: 1.5
   (Maximum risk per sector, in % of equity)

   Max units: 4
   (Maximum add-ons per position: initial + 3)
   ```

4. **Click "Save Settings"**

   Expected: "Settings saved successfully" message

5. **Verify settings saved:**
   - Refresh page (F5)
   - Settings should persist (values still filled)

**Example configuration for $100,000 account:**
- Equity: $100,000
- Risk per unit: 0.75% = $750
- Max risk per position: 4 units × $750 = $3,000 (3%)
- Portfolio cap: 4% = $4,000 total
- Bucket cap: 1.5% = $1,500 per sector

---

### Step 3: Configure FINVIZ Presets (Optional)

**If you want automated daily scans:**

1. **Create FINVIZ screener:**
   - Go to: https://www.finviz.com/screener.ashx
   - Build screener with filters
   - Example for trend-following longs:
     - Price > SMA200
     - RSI > 55
     - Volume > 1M
     - Market Cap > $500M

2. **Get screener URL:**
   - After applying filters, copy URL from browser address bar
   - Should look like:
     ```
     https://finviz.com/screener.ashx?v=111&f=ta_price_a200sma,ta_rsi_os55,sh_avgvol_o1000,sh_marketcap_o500
     ```

3. **Add preset in TF-Engine:**
   - Settings → FINVIZ Presets section
   - Click "Add New Preset"
   - Name: `TF Breakout Long`
   - URL: [Paste FINVIZ URL]
   - Click "Save"

4. **Test preset:**
   - Navigate to Scanner
   - Select preset from dropdown
   - Click "Run Daily FINVIZ Scan"
   - Should return 50-300 candidates (if market conditions favorable)

**If skip this step:**
- Can manually enter tickers in Checklist
- Or add preset later when needed

---

### Step 4: Set Up TradingView (Recommended)

**See [TRADINGVIEW_SETUP.md](TRADINGVIEW_SETUP.md) for detailed guide.**

**Quick steps:**

1. Create TradingView account (free): https://www.tradingview.com
2. Install Ed-Seykota.pine script:
   - Chart → Pine Editor
   - Copy script from `reference/Ed-Seykota.pine`
   - Add to Chart
3. Configure TF-Engine URL template (optional):
   - Settings → TradingView URL Template
   - Enter: `https://tradingview.com/chart/?symbol={ticker}`

---

## Verification

### Verify Installation Success

**Checklist:**

- [ ] tf-engine.exe launches without errors
- [ ] Command window stays open (server running)
- [ ] Browser opens to http://localhost:8080
- [ ] Dashboard screen appears (not blank page)
- [ ] Database file exists at `%APPDATA%\TF-Engine\trading.db`
- [ ] Settings screen loads and saves data
- [ ] No error messages in browser console (F12)

**If all checked:**
✓ Installation successful! Proceed to [QUICK_START.md](QUICK_START.md) or [USER_GUIDE.md](USER_GUIDE.md).

**If any not checked:**
See [Troubleshooting Installation](#troubleshooting-installation) below.

---

### Test Basic Workflow

**Quick smoke test (3 minutes):**

1. **Dashboard:** Should load with empty positions table

2. **Settings:**
   - Enter equity: 100000
   - Risk%: 0.75
   - Save → Should succeed

3. **Checklist:**
   - Enter ticker: AAPL
   - Entry: 180
   - N: 2.5
   - Sector: Tech/Comm
   - Check all 5 required gates
   - Banner should turn GREEN

4. **Position Sizing:**
   - Should pre-fill from checklist
   - Click "Calculate Position Size"
   - Should show results (shares, risk, stop)

**If all work:** ✓ System is functional!

**If any fail:** See [TROUBLESHOOTING.md](TROUBLESHOOTING.md).

---

## Troubleshooting Installation

### Issue: "Windows protected your PC" SmartScreen warning

**When launching tf-engine.exe or installer:**

```
Windows Defender SmartScreen prevented an unrecognized app from starting.
```

**Solution:**

1. Click "More info" (small link)
2. Click "Run anyway"
3. App launches normally

**Why this happens:**
- New/unrecognized app (not code-signed)
- SmartScreen is being cautious
- Safe to run if downloaded from official source

**Optional: Add exclusion**
```
Windows Security → Virus & threat protection
→ Manage settings → Exclusions → Add exclusion
→ File → Select tf-engine.exe
```

---

### Issue: Installer requires admin rights but I don't have them

**Solution:**

Use **Method 1: Standalone Binary** instead.
- Doesn't require admin rights
- Can run from user folder (Documents, AppData)
- Fully functional

---

### Issue: Port 8080 already in use

**Error message:**
```
Failed to start server: listen tcp :8080: bind: address already in use
```

**Solution 1: Find what's using port 8080**

```powershell
# Open PowerShell as Administrator
netstat -ano | findstr :8080

# Output shows PID (Process ID)
# Example: TCP 0.0.0.0:8080 0.0.0.0:0 LISTENING 12345

# Find process name
tasklist | findstr 12345

# If it's another tf-engine.exe: End it
# If it's something else: Stop that service or use different port
```

**Solution 2: Use different port**

```cmd
tf-engine.exe server --listen :8081
```

Then open browser to: http://localhost:8081

---

### Issue: Antivirus blocks or quarantines tf-engine.exe

**Some antivirus software flags new executables.**

**Solution:**

1. **Add exception in antivirus:**
   - Check antivirus quarantine/logs
   - Add tf-engine.exe to whitelist/exceptions
   - Process varies by antivirus:
     - Windows Defender: Virus & threat protection → Exclusions
     - McAfee: Real-Time Scanning → Excluded Files
     - Norton: Settings → Antivirus → Exclusions

2. **Restore quarantined file:**
   - Antivirus → Quarantine → Restore tf-engine.exe

3. **Re-download if needed:**
   - If file deleted: Download again from official source
   - After adding exception

---

### Issue: Cannot write to AppData folder

**Error:**
```
Permission denied: C:\Users\...\AppData\Roaming\TF-Engine\trading.db
```

**Solution:**

**Option 1: Fix permissions**
```
1. Navigate to: %APPDATA%
2. Right-click "TF-Engine" folder (or parent folder)
3. Properties → Security tab
4. Edit → Add → Your username
5. Grant "Modify" and "Write" permissions
6. Apply → OK
```

**Option 2: Run as Administrator once**
```
Right-click tf-engine.exe → Run as Administrator
(Database creates with correct permissions)
(Subsequent runs don't need admin)
```

**Option 3: Use different database location**
```
(Advanced) Specify custom database location:
tf-engine.exe server --db "C:\TF-Engine\data\trading.db"
```

---

## Uninstallation

### Uninstall Standalone Binary

1. **Close TF-Engine:**
   - Close browser tabs (http://localhost:8080)
   - Close command window (or End Task in Task Manager)

2. **Backup database (optional but recommended):**
   ```
   Navigate to: %APPDATA%\TF-Engine
   Copy trading.db to safe location (Dropbox, external drive)
   ```

3. **Delete binary:**
   ```
   Delete: C:\TF-Engine\tf-engine.exe
   Delete: C:\TF-Engine\ (entire folder if desired)
   ```

4. **Delete data folder (if desired):**
   ```
   Navigate to: %APPDATA%\TF-Engine
   Delete entire folder
   (WARNING: This deletes all your trading data!)
   ```

5. **Delete shortcuts:**
   - Desktop shortcut (if created)
   - Taskbar pin (if pinned)

---

### Uninstall Windows Installer

1. **Close TF-Engine** (same as above)

2. **Backup database** (same as above - optional but recommended)

3. **Uninstall via Windows Settings:**

   **Windows 10:**
   ```
   Settings → Apps → Apps & features
   Search: "TF-Engine"
   Click → Uninstall → Confirm
   ```

   **Windows 11:**
   ```
   Settings → Apps → Installed apps
   Search: "TF-Engine"
   Click ⋮ (three dots) → Uninstall → Confirm
   ```

4. **Installer asks: Keep data?**
   ```
   [  ] Keep trading data in AppData
   [✓] Delete all data

   (Check "Delete all data" if sure)
   (Uncheck to preserve database for potential reinstall)
   ```

5. **Verify removed:**
   ```
   Check: C:\Program Files\TF-Engine\ (should be deleted)
   Check: Start Menu → TF-Engine (should be gone)
   Check: Desktop shortcut (should be deleted)
   ```

6. **Delete AppData manually (if not deleted by installer):**
   ```
   Navigate to: %APPDATA%\TF-Engine
   Delete folder
   ```

---

## Upgrading

### Upgrade Standalone Binary

1. **Backup current database:**
   ```
   Navigate to: %APPDATA%\TF-Engine
   Copy trading.db to: trading-backup-2025-10-29.db
   ```

2. **Close current TF-Engine:**
   - Close browser tabs
   - Close command window

3. **Download new version:**
   - Example: tf-engine-v1.1.0.exe

4. **Replace binary:**
   ```
   Delete: C:\TF-Engine\tf-engine.exe (old version)
   Move: tf-engine-v1.1.0.exe → C:\TF-Engine\tf-engine.exe

   Or: Rename new binary to tf-engine.exe
   ```

5. **Launch new version:**
   - Double-click tf-engine.exe
   - Database automatically migrates (if needed)
   - Check version: Settings → About (or in UI footer)

6. **Verify upgrade:**
   - Check data still present (positions, settings)
   - Test basic workflow
   - If issues: Restore backup and rollback

---

### Upgrade Windows Installer

1. **Backup database** (same as above)

2. **Download new installer:**
   - Example: TF-Engine-Setup-v1.1.0.msi

3. **Run new installer:**
   - Double-click installer
   - Installer detects old version
   - Options:
     - Upgrade (recommended)
     - Uninstall old, install new
     - Cancel

4. **Choose "Upgrade":**
   - Preserves data automatically
   - Replaces binary
   - Keeps settings

5. **Restart TF-Engine:**
   - Launch from Start Menu or Desktop shortcut
   - Database migrates automatically (if needed)

6. **Verify upgrade:**
   - Check version
   - Check data present
   - Test basic workflow

---

## Data Locations

**Key files and folders:**

| Item | Location |
|------|----------|
| Binary | `C:\TF-Engine\tf-engine.exe` (or Program Files) |
| Database | `%APPDATA%\TF-Engine\trading.db` |
| Logs | `%APPDATA%\TF-Engine\logs\` (if logging enabled) |
| Config | `%APPDATA%\TF-Engine\config.json` (if exists) |
| Backups | `%APPDATA%\TF-Engine\backups\` (if auto-backup enabled) |

**To access AppData quickly:**
```
Win+R → type: %APPDATA%\TF-Engine → Enter
```

---

## Firewall Configuration

**If Windows Firewall blocks TF-Engine:**

### Allow Through Firewall

1. **Windows Security:**
   ```
   Start → Settings → Privacy & Security → Windows Security
   → Firewall & network protection
   → Allow an app through firewall
   ```

2. **Click "Change settings"** (requires admin)

3. **Click "Allow another app..."**

4. **Browse:**
   ```
   Navigate to: C:\TF-Engine\tf-engine.exe (or Program Files)
   Select: tf-engine.exe
   Click "Add"
   ```

5. **Check both:**
   ```
   [✓] Private networks
   [✓] Public networks
   ```

6. **Click "OK"**

7. **Restart TF-Engine**

---

## Next Steps

**After successful installation:**

1. **Quick Start:** [QUICK_START.md](QUICK_START.md) - Get started in 10 minutes
2. **TradingView Setup:** [TRADINGVIEW_SETUP.md](TRADINGVIEW_SETUP.md) - Install Pine script
3. **User Guide:** [USER_GUIDE.md](USER_GUIDE.md) - Comprehensive documentation
4. **Daily Workflow:** Learn the checklist → sizing → heat → entry flow

**Start trading systematically!**

---

## Support

**If you encounter issues not covered here:**

- [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Common issues
- [FAQ.md](FAQ.md) - Frequently asked questions
- Contact: [your-support-email or GitHub issues]

---

**Version:** 1.0.0
**Last Updated:** 2025-10-29
**System:** TF-Engine - Trend Following Engine
