@echo off
REM ===========================================================================
REM import-candidates-auto.bat - Automated Daily Candidate Import
REM ===========================================================================
REM Purpose: Quick auto-import using default settings (no prompts)
REM
REM What this does:
REM   - Runs interactive mode with --auto flag
REM   - Uses default preset (Trend Following)
REM   - Skips all confirmation prompts
REM   - Shows cool ASCII art and animations
REM   - Perfect for daily morning routine
REM
REM Usage: Double-click this file or run from command line
REM        Add to Windows Task Scheduler for automation
REM ===========================================================================

REM Set console to UTF-8 for emoji support
chcp 65001 >nul 2>&1

REM Set window title
title Trading Engine - Auto Import

echo ╔═══════════════════════════════════════════════════════════╗
echo ║                                                           ║
echo ║           🚀 AUTOMATED CANDIDATE IMPORT 🚀                ║
echo ║                                                           ║
echo ║   Running with default settings (no prompts)...          ║
echo ║                                                           ║
echo ╚═══════════════════════════════════════════════════════════╝
echo.

REM Check if tf-engine.exe exists
if not exist "tf-engine.exe" (
    echo ❌ ERROR: tf-engine.exe not found
    pause
    exit /b 1
)

REM Check if database is initialized (silent auto-init)
if not exist "trading.db" (
    echo.
    echo ⚙️  Initializing database...
    tf-engine.exe init >nul 2>&1
    if %errorlevel% neq 0 (
        echo ❌ ERROR: Database initialization failed!
        pause
        exit /b 1
    )
    echo ✅ Database initialized
    echo.
)

REM Launch auto mode
tf-engine.exe interactive --auto

echo.
echo ═══════════════════════════════════════════════════════════
echo   ✅ Done! Press any key to close...
echo ═══════════════════════════════════════════════════════════
pause >nul
