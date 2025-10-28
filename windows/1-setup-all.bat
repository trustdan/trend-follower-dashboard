@echo off
REM ===========================================================================
REM setup-all.bat - Complete Automated Setup for Trading Engine v3
REM ===========================================================================
REM Purpose: One-click setup - creates workbook, imports VBA, initializes DB
REM
REM What this does:
REM   1. Creates TradingPlatform.xlsm workbook with VBA project access enabled
REM   2. Imports all VBA modules
REM   3. Initializes database
REM   4. Configures Excel named ranges
REM   5. Creates test button
REM   6. Runs smoke tests
REM
REM Usage: Just run this script!
REM   setup-all.bat
REM
REM Created: 2025-10-27 (M21 - Windows Integration Validation)
REM ===========================================================================

setlocal enabledelayedexpansion

REM Setup logging
set LOGFILE=setup-all.log
echo Setup started at %date% %time% > %LOGFILE%

echo ========================================
echo  Trading Engine v3 - Complete Setup
echo  One-Click Installation (UPDATED)
echo ========================================
echo.
echo This script will:
echo   1. Check prerequisites (tf-engine.exe, VBA files)
echo   2. Enable VBA project access (registry)
echo   3. Create workbook with UI and VBA (consolidated)
echo   4. Initialize database (trading.db)
echo   5. Run automated smoke tests
echo.
echo Changes from previous version:
echo   - Consolidated workbook/VBA/UI creation into one step
echo   - Uses TRUE/FALSE dropdowns instead of checkboxes
echo   - More reliable (avoids Excel OLE automation issues)
echo.
echo Estimated time: 2-3 minutes
echo.

pause

REM ============================================================================
REM Step 1: Check Prerequisites
REM ============================================================================

echo.
echo [Step 1/8] Checking prerequisites...
echo.

REM Check if tf-engine.exe exists
if not exist "tf-engine.exe" (
    echo ERROR: tf-engine.exe not found in current directory
    echo Please ensure you're running this from the windows\ folder
    pause
    exit /b 1
)

REM Check if VBA modules exist
if not exist "..\excel\vba\TFTypes.bas" (
    echo ERROR: VBA modules not found at ..\excel\vba\
    echo Please ensure the project structure is intact
    pause
    exit /b 1
)

echo [OK] Prerequisites verified

REM ============================================================================
REM Step 2: Enable VBA Project Access via Registry (MOVED UP)
REM ============================================================================
REM NOTE: This must happen BEFORE creating the workbook so VBA import works
REM ============================================================================

echo.
echo [Step 2/8] Enabling VBA project access...
echo.

reg add "HKCU\Software\Microsoft\Office\16.0\Excel\Security" /v AccessVBOM /t REG_DWORD /d 1 /f >nul 2>&1
reg add "HKCU\Software\Microsoft\Office\15.0\Excel\Security" /v AccessVBOM /t REG_DWORD /d 1 /f >nul 2>&1
reg add "HKCU\Software\Microsoft\Office\14.0\Excel\Security" /v AccessVBOM /t REG_DWORD /d 1 /f >nul 2>&1

echo [OK] VBA project access enabled

REM ============================================================================
REM Step 3: Create Excel Workbook with UI and VBA Modules (CONSOLIDATED)
REM ============================================================================
REM
REM NOTE: Old Steps 2, 4, and 5 have been consolidated into one VBScript
REM      to avoid Excel OLE automation issues with checkboxes
REM
REM What this does:
REM   - Creates TradingPlatform.xlsm
REM   - Imports all 4 VBA modules (TFTypes, TFHelpers, TFEngine, TFTests)
REM   - Creates all 7 worksheets (Setup, Dashboard, Position Sizing, etc.)
REM   - Uses TRUE/FALSE dropdowns instead of checkboxes (more reliable)
REM
REM ============================================================================

echo.
echo [Step 3/8] Creating workbook with UI and VBA...
echo.

REM Check if consolidated script exists
if not exist "create-workbook-manual-ui.vbs" (
    echo ERROR: create-workbook-manual-ui.vbs not found
    echo Please ensure the file exists in the windows\ folder
    pause
    exit /b 1
)

REM Run consolidated script
echo Running: cscript //nologo create-workbook-manual-ui.vbs >> %LOGFILE%
cscript //nologo create-workbook-manual-ui.vbs >> %LOGFILE% 2>&1

if %errorlevel% neq 0 (
    echo ERROR: Failed to create workbook
    echo ERROR: Failed to create workbook. Error level: %errorlevel% >> %LOGFILE%
    echo.
    echo Check setup-all.log for details
    pause
    exit /b 1
)

echo [OK] Workbook created with all UI and VBA modules
echo [OK] Created: 7 worksheets, 4 VBA modules, buttons and dropdowns >> %LOGFILE%

REM ============================================================================
REM Step 4: Initialize Database
REM ============================================================================

echo.
echo [Step 4/6] Initializing database...
echo.

REM Check if database already exists
if exist "trading.db" (
    echo WARNING: trading.db already exists
    set /p overwrite="Overwrite existing database? (y/n): "
    if /i not "!overwrite!"=="y" (
        echo Skipping database initialization
        goto :SkipDB
    )
    del trading.db
)

REM Initialize database
tf-engine.exe init

if %errorlevel% neq 0 (
    echo ERROR: Database initialization failed
    pause
    exit /b 1
)

echo [OK] Database initialized

:SkipDB

REM ============================================================================
REM Step 5: Run Smoke Tests
REM ============================================================================

echo.
echo [Step 5/6] Running automated smoke tests...
echo.

REM Run quick smoke tests (subset of run-tests.bat)
echo Testing engine version...
tf-engine.exe --version >nul 2>&1
if %errorlevel% neq 0 (
    echo [FAIL] Engine version check failed
) else (
    echo [PASS] Engine version check
)

echo Testing database access...
tf-engine.exe get-settings --format json >nul 2>&1
if %errorlevel% neq 0 (
    echo [FAIL] Database access failed
) else (
    echo [PASS] Database access
)

echo Testing position sizing...
tf-engine.exe size --entry 180 --atr 1.5 --k 2 --method stock --format json >nul 2>&1
if %errorlevel% neq 0 (
    echo [FAIL] Position sizing failed
) else (
    echo [PASS] Position sizing
)

REM ============================================================================
REM Complete!
REM ============================================================================

echo.
echo ========================================
echo  Setup Complete!
echo ========================================
echo.
echo Created files:
echo   - TradingPlatform.xlsm (Excel workbook with VBA and 7 sheets)
echo   - trading.db (SQLite database)
echo   - setup-all.log (setup process log)
echo   - TradingSystem_Debug.log (will be created on first use)
echo.
echo Workbook contains (7 worksheets):
echo   - Setup (configuration with named ranges)
echo   - Dashboard (portfolio overview)
echo   - Position Sizing (calculate shares/contracts)
echo   - Checklist (6-item evaluation with TRUE/FALSE dropdowns)
echo   - Heat Check (portfolio/bucket heat management)
echo   - Trade Entry (5-gate decision workflow)
echo   - VBA Tests (automated testing with Run All Tests button)
echo.
echo VBA Modules imported (4 modules):
echo   - TFTypes (data structures)
echo   - TFHelpers (JSON parsing, validation)
echo   - TFEngine (engine communication)
echo   - TFTests (unit tests)
echo.
echo Next steps:
echo   1. Open TradingPlatform.xlsm in Excel
echo   2. Enable macros when prompted
echo   3. Go to "VBA Tests" sheet
echo   4. Click "Run All Tests" button
echo   5. Verify all tests pass (GREEN)
echo   6. Start trading with Dashboard sheet
echo.
echo For integration testing, run: 3-run-integration-tests.bat
echo For CLI smoke tests, run: 4-run-tests.bat
echo.
echo Documentation:
echo   - QUICK_START.md (quick setup guide)
echo   - README_UI_FIX.md (explains dropdown vs checkbox change)
echo   - WINDOWS_TESTING.md (complete testing guide)
echo   - README.md (package overview)
echo.

pause
