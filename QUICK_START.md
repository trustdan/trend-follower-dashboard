# TF-Engine Quick Start

## Fix the Database Error NOW

You're seeing this error:
```
Failed to create session with options: SQL logic error:
table trade_sessions has no column named instrument_type (1)
```

### Solution (Takes 10 seconds)

**Option 1: Quick Fix (Recommended)**
```bash
.\migrate-db.exe
```

**Option 2: Full Rebuild**
```bash
.\build-windows.bat
```

This adds 27 missing columns to your database for options trading support.

## Daily Usage

### Build Scripts Available

1. **`build-windows.bat`** - Full rebuild (migration + GUI)
   ```bash
   .\build-windows.bat
   ```
   Use when: First time, after git pull, database errors

2. **`quick-rebuild.bat`** - Fast GUI rebuild only
   ```bash
   .\quick-rebuild.bat
   ```
   Use when: Making UI changes, quick iteration

3. **`migrate-db.exe`** - Fix database manually
   ```bash
   .\migrate-db.exe
   ```
   Use when: Just need to update database schema

### Run the Application

```bash
.\tf-gui.exe
```

## What Was Built

✅ **migrate-db.exe** (8.9 MB)
   - Fixes database schema
   - Adds options trading support
   - Safe to run multiple times

✅ **tf-gui.exe** (50 MB)
   - Main trading application
   - VIM mode enabled (press V)
   - Dark mode toggle
   - Trade sessions with options support

## Files Created

```
📦 Build Scripts (choose one):
├── build-windows.bat       ← Full build (Windows batch)
├── build-windows.ps1       ← Full build (PowerShell, colored output)
└── quick-rebuild.bat       ← Fast iteration (GUI only)

📚 Documentation:
├── BUILD_GUIDE.md          ← Comprehensive build instructions
├── MIGRATION_INSTRUCTIONS.md  ← Database migration details
└── QUICK_START.md          ← This file!

🔧 Executables (built):
├── migrate-db.exe          ← Database migration tool
└── tf-gui.exe              ← Main GUI application

📊 Database (created on first run):
└── trading.db              ← SQLite database
```

## Next Steps

1. **Fix the database error:**
   ```bash
   .\migrate-db.exe
   ```

2. **Launch the app:**
   ```bash
   .\tf-gui.exe
   ```

3. **Try VIM mode:**
   - Press `V` to toggle VIM mode
   - Press `F` to show button hints
   - Use `j/k` to scroll, `gg/G` for top/bottom
   - Type hint letters to click buttons

4. **Create a trade session:**
   - Click "Start New Trade" button
   - Should work without the `instrument_type` error!

## Troubleshooting

**Still seeing database error after migration?**
```bash
# Verify migration ran successfully
.\migrate-db.exe

# Check database file exists
dir trading.db

# Try full rebuild
.\build-windows.bat
```

**GUI won't start?**
- Close any running `tf-gui.exe` instances
- Check if `trading.db` is locked by another program
- Run migration: `.\migrate-db.exe`

**Build failed?**
- Verify Go is installed: `go version`
- Should show Go 1.21 or higher
- Download from: https://go.dev/dl/

## Development Workflow

Making UI changes:
```bash
# 1. Edit files in ui/*.go
# 2. Quick rebuild
.\quick-rebuild.bat
# 3. Test
.\tf-gui.exe
# 4. Repeat
```

Making backend changes:
```bash
# 1. Edit files in backend/internal/*
# 2. Full rebuild
.\build-windows.bat
# 3. Test
.\tf-gui.exe
```

## More Help

- **Comprehensive guide:** [BUILD_GUIDE.md](BUILD_GUIDE.md)
- **Migration details:** [MIGRATION_INSTRUCTIONS.md](MIGRATION_INSTRUCTIONS.md)
- **Project docs:** `docs/` folder

---

**Built:** 2025-10-30
**Status:** ✅ Ready to use!
