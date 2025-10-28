# M24 - Windows Transfer & Validation Guide

**Purpose:** Transfer release package from Linux to Windows PC for validation
**Estimated Time:** 5-10 minutes (transfer) + 5-90 minutes (validation)

---

## 📦 What You're Transferring

**File:** `TradingEngine-v3.0.0-rc1.zip`
**Location:** `/home/kali/excel-trading-platform/release/`
**Size:** 16 MB
**SHA256:** `cf3d2e72fb77ec30ad15cd7bb9568a7c22a9777804d84d9c21e96210406363f4`

**Contains:**
- Windows binary (tf-engine.exe)
- Complete documentation
- Setup scripts
- VBA source code
- All project files

---

## 🚀 Transfer Methods

Choose the method that works for your setup:

### Method 1: USB Drive (Simplest)

**On Linux:**
```bash
# 1. Insert USB drive
# 2. Mount if not auto-mounted (check with 'lsblk')
sudo mount /dev/sdb1 /mnt/usb  # Adjust device name

# 3. Copy file
cp /home/kali/excel-trading-platform/release/TradingEngine-v3.0.0-rc1.zip /mnt/usb/
cp /home/kali/excel-trading-platform/release/TradingEngine-v3.0.0-rc1.sha256 /mnt/usb/

# 4. Verify copy
ls -lh /mnt/usb/Trading*

# 5. Safely unmount
sudo umount /mnt/usb
```

**On Windows:**
1. Eject USB from Linux
2. Insert into Windows PC
3. Copy files to `C:\Trading\`
4. Verify file size (should be ~16 MB)

---

### Method 2: Network Transfer (SCP/SFTP)

**If Windows has SSH/SFTP client (e.g., WinSCP, FileZilla):**

**On Linux - Start SSH server (if not running):**
```bash
sudo systemctl start ssh
sudo systemctl status ssh
ip addr show  # Note your IP address
```

**On Windows - Using WinSCP:**
1. Download WinSCP: https://winscp.net/
2. Install and open
3. New Site → Protocol: SFTP
4. Host name: [Linux IP from above]
5. Port: 22
6. User name: kali (or your username)
7. Password: [your password]
8. Login
9. Navigate to: `/home/kali/excel-trading-platform/release/`
10. Download: `TradingEngine-v3.0.0-rc1.zip`
11. Download: `TradingEngine-v3.0.0-rc1.sha256`
12. Save to: `C:\Trading\`

**On Windows - Using Command Line (if OpenSSH installed):**
```cmd
REM From Windows Command Prompt
scp kali@[LINUX-IP]:/home/kali/excel-trading-platform/release/TradingEngine-v3.0.0-rc1.zip C:\Trading\
scp kali@[LINUX-IP]:/home/kali/excel-trading-platform/release/TradingEngine-v3.0.0-rc1.sha256 C:\Trading\
```

---

### Method 3: Shared Folder (VM Setup)

**If Linux is in a VM (VMware/VirtualBox):**

**VMware Shared Folders:**
```bash
# On Linux
sudo mkdir -p /mnt/hgfs/Shared
sudo vmhgfs-fuse .host:/Shared /mnt/hgfs/Shared -o allow_other
cp /home/kali/excel-trading-platform/release/TradingEngine-v3.0.0-rc1.* /mnt/hgfs/Shared/
```

**VirtualBox Shared Folders:**
```bash
# On Linux (assuming shared folder named "Shared")
sudo mkdir -p /mnt/shared
sudo mount -t vboxsf Shared /mnt/shared
cp /home/kali/excel-trading-platform/release/TradingEngine-v3.0.0-rc1.* /mnt/shared/
```

**On Windows:**
- Files appear in the shared folder automatically
- Copy to `C:\Trading\` for testing

---

### Method 4: Cloud Storage (Dropbox, Google Drive, OneDrive)

**On Linux:**
```bash
# Install cloud client if needed (example for Dropbox)
# Or use web interface

# Copy to cloud folder
cp /home/kali/excel-trading-platform/release/TradingEngine-v3.0.0-rc1.zip ~/Dropbox/
```

**On Windows:**
1. Open cloud storage client or web interface
2. Download file to `C:\Trading\`
3. Verify download complete

---

### Method 5: Direct HTTP Server (Quick & Easy)

**On Linux:**
```bash
# Navigate to release directory
cd /home/kali/excel-trading-platform/release

# Start simple HTTP server
python3 -m http.server 8000

# Note: Server will show "Serving HTTP on 0.0.0.0 port 8000"
# Get your IP address:
ip addr show | grep "inet "
```

**On Windows:**
1. Open browser
2. Go to: `http://[LINUX-IP]:8000/`
3. Click `TradingEngine-v3.0.0-rc1.zip` to download
4. Save to `C:\Trading\`

**On Linux (when done):**
- Press `Ctrl+C` to stop server

---

## ✅ Verify Transfer

**On Windows - Check File Integrity:**

### Option A: Using PowerShell (Recommended)
```powershell
# Open PowerShell
cd C:\Trading

# Check file exists
Get-Item TradingEngine-v3.0.0-rc1.zip

# Verify file size (should be ~16 MB)
(Get-Item TradingEngine-v3.0.0-rc1.zip).Length / 1MB

# Calculate SHA256 checksum
Get-FileHash TradingEngine-v3.0.0-rc1.zip -Algorithm SHA256

# Compare with expected:
# cf3d2e72fb77ec30ad15cd7bb9568a7c22a9777804d84d9c21e96210406363f4
```

### Option B: Using certutil (Windows built-in)
```cmd
cd C:\Trading
certutil -hashfile TradingEngine-v3.0.0-rc1.zip SHA256

REM Should show:
REM cf3d2e72fb77ec30ad15cd7bb9568a7c22a9777804d84d9c21e96210406363f4
```

### Option C: Visual Check
```cmd
dir C:\Trading\TradingEngine-v3.0.0-rc1.zip

REM File size should be approximately:
REM 16,000,000 bytes (16 MB)
```

**If checksums don't match:** Transfer corrupted, try again

**If checksums match:** ✅ Transfer successful!

---

## 📂 Extract Package

**On Windows:**

### Using Windows Explorer (GUI)
```
1. Navigate to C:\Trading\
2. Right-click TradingEngine-v3.0.0-rc1.zip
3. Select "Extract All..."
4. Extract to: C:\Trading\TradingEngine-v3\
5. Click "Extract"
```

### Using Command Line
```cmd
cd C:\Trading

REM Using PowerShell
powershell Expand-Archive TradingEngine-v3.0.0-rc1.zip -DestinationPath .

REM Or using 7-Zip (if installed)
"C:\Program Files\7-Zip\7z.exe" x TradingEngine-v3.0.0-rc1.zip

REM Or using tar (Windows 10+)
tar -xf TradingEngine-v3.0.0-rc1.zip
```

**Verify Extraction:**
```cmd
cd C:\Trading\TradingEngine-v3
dir

REM Should see:
REM - tf-engine.exe (~12 MB)
REM - QUICKSTART.md
REM - README.md
REM - TROUBLESHOOTING.md
REM - KNOWN_LIMITATIONS.md
REM - WINDOWS_VALIDATION.md
REM - windows\ directory
REM - excel\ directory
REM - docs\ directory
```

---

## 🧪 Quick Validation (5 minutes)

**Before full validation, quick smoke test:**

```cmd
cd C:\Trading\TradingEngine-v3

REM Test 1: Binary runs
tf-engine.exe --version
REM Should show: tf-engine version 3.0.0-dev

REM Test 2: Help works
tf-engine.exe --help
REM Should show command list

REM Test 3: Init works
tf-engine.exe init
REM Should create trading.db

REM Test 4: Settings work
tf-engine.exe get-setting Equity_E
REM Should show: 10000
```

**If all 4 tests pass:** ✅ Binary works on Windows!

**If any fail:** See `TROUBLESHOOTING.md` in the package

---

## 📋 Next Steps

### 1. Read Documentation (5 minutes)
```
Open in this order:
1. QUICKSTART.md ⭐ (5-min read)
2. WINDOWS_VALIDATION.md (understand test levels)
```

### 2. Choose Validation Level

**Level 1: Quick (5 min)**
- Run integration tests
- Verify heat tests pass
- Good for: Confirming M23 works

**Level 2: Standard (25 min)**
- Full UI testing
- All workflows
- Good for: Production confidence

**Level 3: Comprehensive (90 min)**
- Manual gate tests
- Complete validation
- Good for: Full certification

### 3. Run Validation

**See:** `WINDOWS_VALIDATION.md` for detailed checklist

**Quick Start:**
```cmd
cd C:\Trading\TradingEngine-v3\windows
1-setup-all.bat
REM Wait for setup to complete (~3 minutes)

REM Then run tests
3-run-integration-tests.bat
```

**Expected Results:**
```
✅ PASS: 13/19 (68.4%)
⏭️ SKIP: 6/19 (31.6%)

Tests 3.1-3.4 (Heat) should show PASS (not SKIP)
```

### 4. Report Results

**If All Tests Pass:**
- Document validation complete
- Mark as production-ready
- Proceed with final release

**If Any Tests Fail:**
1. Check `TROUBLESHOOTING.md`
2. Review logs (`tf-engine.log`, `TradingSystem_Debug.log`)
3. Document specific failures
4. Report via GitHub issues

---

## 🔧 Troubleshooting Transfer

### File Won't Transfer
**Possible causes:**
- Network connection issues
- Permission problems
- Disk space on target

**Solutions:**
- Check network connectivity
- Ensure write permissions on C:\Trading\
- Verify 50+ MB free space on Windows

### Checksum Mismatch
**Possible causes:**
- Transfer interrupted
- File corrupted
- Wrong file

**Solutions:**
- Delete and retry transfer
- Use different transfer method
- Verify source file on Linux first

### Can't Extract ZIP
**Possible causes:**
- Corrupted download
- Windows Defender blocking
- Insufficient permissions

**Solutions:**
- Verify checksum before extracting
- Temporarily disable Windows Defender
- Run as Administrator
- Use different extraction tool (7-Zip)

### Binary Won't Run
**Possible causes:**
- Windows Defender/Antivirus blocking
- Missing dependencies
- Wrong architecture (unlikely with Go)

**Solutions:**
- Add tf-engine.exe to antivirus exceptions
- Run as Administrator
- Check Windows Event Viewer for errors

---

## 📞 Support

### Before Asking for Help

1. ✅ Verify checksum matches
2. ✅ Try extraction with different tool
3. ✅ Check `TROUBLESHOOTING.md`
4. ✅ Review Windows Event Viewer
5. ✅ Try on different Windows PC (if available)

### Reporting Issues

**Include:**
- Windows version (Settings → System → About)
- Excel version (File → Account → About Excel)
- Transfer method used
- Checksum comparison results
- Error messages (exact text)
- Screenshots if helpful

**Where:** https://github.com/anthropics/claude-code/issues

---

## ✅ Success Checklist

**Transfer Complete When:**
- [ ] File copied to Windows PC
- [ ] Checksum verified (matches expected)
- [ ] ZIP extracted successfully
- [ ] `tf-engine.exe` exists and is ~12 MB
- [ ] All documentation files present
- [ ] Quick smoke test passed (4 commands)
- [ ] Ready to read QUICKSTART.md

**Next:** Follow `QUICKSTART.md` for setup and validation

---

## 📊 Transfer Summary

**Source:**
- Location: `/home/kali/excel-trading-platform/release/`
- File: `TradingEngine-v3.0.0-rc1.zip`
- Size: 16 MB
- Checksum: `cf3d2e72fb77ec30ad15cd7bb9568a7c22a9777804d84d9c21e96210406363f4`

**Destination:**
- Location: `C:\Trading\TradingEngine-v3\`
- Extracted size: ~20 MB
- Files: ~100+ files
- Ready for: Setup and validation

**Validation Time:**
- Quick: 5 minutes
- Standard: 25 minutes
- Comprehensive: 90 minutes

**Choose your validation level and proceed!**

---

**Last Updated:** 2025-10-28
**Status:** Ready for transfer
**Next:** Transfer file and run `QUICKSTART.md`
