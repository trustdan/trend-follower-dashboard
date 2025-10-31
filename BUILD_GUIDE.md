# TF-Engine Build Guide for Windows

## Quick Start

### First Time Setup
```powershell
# Run full build (includes database migration)
.\build-windows.bat
```

### Daily Development
```powershell
# Quick GUI rebuild (fast iteration)
.\quick-rebuild.bat
```

## Build Scripts

### 1. `build-windows.bat` - Full Build
**Use when:** First setup, after pulling changes, fixing database errors

**What it does:**
1. Builds `migrate-db.exe` (database migration tool)
2. Runs database migration (adds missing columns)
3. Builds `tf-gui.exe` (main GUI application)
4. Sends toast notification when complete

**Runtime:** ~30-60 seconds

```batch
.\build-windows.bat
```

### 2. `build-windows.ps1` - PowerShell Version
Same as above, but with better error handling and colored output.

**Use when:** You prefer PowerShell or need detailed build output

```powershell
.\build-windows.ps1
```

### 3. `quick-rebuild.bat` - Fast Iteration
**Use when:** Making UI changes and need fast rebuild cycles

**What it does:**
1. Builds `tf-gui.exe` only
2. Skips migration (assumes database is already up-to-date)
3. Sends toast notification when complete

**Runtime:** ~10-20 seconds

```batch
.\quick-rebuild.bat
```

## Fixing the Database Error

If you see this error:
```
Failed to create session with options: SQL logic error:
table trade_sessions has no column named instrument_type (1)
```

**Solution:**
```batch
# Option 1: Run full build (includes migration)
.\build-windows.bat

# Option 2: Run migration only
.\migrate-db.exe
```

The migration adds 27 columns needed for options trading:
- `instrument_type` - 'STOCK' or 'OPTION'
- `options_strategy` - Strategy name
- `entry_date`, `dte`, `roll_threshold_dte` - Time tracking
- `legs_json` - Multi-leg options
- `add_price_1/2/3` - Pyramid add-on prices
- `entry_lookback`, `exit_lookback` - Breakout system
- ... and 17 more columns

**Safe to run multiple times** - Migration skips columns that already exist.

## Build Requirements

### Required Software
- **Go 1.21+** - Download from https://go.dev/dl/
- **Git** (optional) - For version control
- **Windows 10/11** - For toast notifications

### Verify Go Installation
```powershell
go version
# Should show: go version go1.21.x windows/amd64 (or higher)
```

## Directory Structure After Build

```
trend-follower-dashboard/
├── tf-gui.exe              ← Main GUI application
├── migrate-db.exe          ← Database migration tool
├── trading.db              ← SQLite database (created on first run)
├── build-windows.bat       ← Full build script
├── build-windows.ps1       ← PowerShell build script
├── quick-rebuild.bat       ← Fast rebuild script
└── backend/
    └── migrations/
        └── 002_add_options_columns.sql  ← Migration SQL
```

## Troubleshooting

### Build Fails: "go: command not found"
**Problem:** Go is not installed or not in PATH

**Solution:**
1. Download Go from https://go.dev/dl/
2. Run installer (adds to PATH automatically)
3. Open **new** command prompt
4. Verify: `go version`

### Build Fails: "package X is not in GOROOT"
**Problem:** Missing dependencies

**Solution:**
```powershell
cd backend
go mod download
go mod tidy
cd ..
.\build-windows.bat
```

### Migration Fails: "no such table: trade_sessions"
**Problem:** Database not initialized

**Solution:**
```powershell
# Initialize database first
cd backend
go build -o ..\tf-engine.exe .\cmd\tf-engine\main.go
cd ..
.\tf-engine.exe init

# Now run migration
.\migrate-db.exe

# Now build GUI
.\quick-rebuild.bat
```

### GUI Won't Start: "trading.db locked"
**Problem:** Another program has database open

**Solution:**
1. Close any open database browsers (DB Browser for SQLite, etc.)
2. Close any running `tf-gui.exe` instances
3. Try again

### Toast Notifications Not Showing
**Problem:** Windows notifications disabled or blocked

**Solution:**
1. Open Windows Settings
2. System → Notifications
3. Enable notifications for "Claude Code" or your terminal app
4. Builds still work - you just won't see toast notifications

## Performance Tips

### Fast Development Cycle
```powershell
# 1. Make UI changes in ui/*.go files
# 2. Quick rebuild (10-20 seconds)
.\quick-rebuild.bat

# 3. Test immediately
.\tf-gui.exe

# 4. Repeat
```

### Backend Changes
If you modify backend code (`backend/internal/*`):
```powershell
# Full rebuild needed
.\build-windows.bat
```

### Migration Changes
If you modify `backend/migrations/*.sql`:
```powershell
# 1. Rebuild migration tool
cd backend
go build -o ..\migrate-db.exe .\cmd\migrate\main.go
cd ..

# 2. Run new migration
.\migrate-db.exe

# 3. Rebuild GUI
.\quick-rebuild.bat
```

## Advanced Usage

### Build for Different Platforms

**From Windows to Linux:**
```powershell
cd ui
$env:GOOS="linux"; $env:GOARCH="amd64"
go build -o tf-gui-linux
```

**From Windows to macOS:**
```powershell
cd ui
$env:GOOS="darwin"; $env:GOARCH="amd64"
go build -o tf-gui-mac
```

### Clean Build
```powershell
# Remove all built executables
Remove-Item tf-gui.exe, migrate-db.exe, tf-engine.exe -ErrorAction SilentlyContinue

# Rebuild from scratch
.\build-windows.bat
```

### Run Tests
```powershell
cd backend
go test ./... -v
cd ..
```

## What Each Build Produces

### migrate-db.exe
- **Size:** ~15 MB
- **Purpose:** One-time database schema upgrade
- **When to run:** After upgrading project, if you see column errors
- **Safe to run:** Multiple times (idempotent)

### tf-gui.exe
- **Size:** ~25-30 MB
- **Purpose:** Main trading application GUI
- **When to run:** Daily trading use
- **Dependencies:** Needs `trading.db` in same directory

## Getting Help

**Build issues?**
1. Check this guide's Troubleshooting section
2. Verify Go installation: `go version`
3. Check dependencies: `go mod verify` in backend/
4. Run full build: `.\build-windows.bat`

**Database issues?**
1. Run migration: `.\migrate-db.exe`
2. Check `MIGRATION_INSTRUCTIONS.md`
3. Verify database exists: `dir trading.db`

**GUI issues?**
1. Check database is migrated
2. Verify no other instances running
3. Check Windows Event Viewer for crashes

---

**Last Updated:** 2025-10-30
**For:** TF-Engine v2.0 (Options Trading Enhancement)
