@echo off
REM ===========================================================================
REM windows-init-database.bat - Initialize Trading Database
REM ===========================================================================
REM Purpose: Initialize trading.db with schema and default settings
REM
REM Prerequisites:
REM   - tf-engine.exe in current directory
REM
REM Usage: windows-init-database.bat
REM
REM Created: 2025-10-27 (M20 - Windows Integration Package)
REM ===========================================================================

echo ========================================
echo  Trading Database Initialization
echo  Trading Engine v3
echo ========================================
echo.

REM Check if tf-engine.exe exists
if not exist "tf-engine.exe" (
    echo ERROR: tf-engine.exe not found in current directory
    echo Please ensure tf-engine.exe is in the same folder as this script.
    pause
    exit /b 1
)

echo Engine: tf-engine.exe
echo Database: trading.db
echo.

REM Check if database already exists
if exist "trading.db" (
    echo WARNING: trading.db already exists!
    echo.
    set /p confirm="Do you want to reinitialize? (y/n): "
    if /i not "%confirm%"=="y" (
        echo.
        echo Initialization cancelled.
        pause
        exit /b 0
    )
    echo.
    echo Backing up existing database...
    set timestamp=%date:~-4%%date:~-10,2%%date:~-7,2%_%time:~0,2%%time:~3,2%%time:~6,2%
    set timestamp=%timestamp: =0%
    copy trading.db trading.db.backup_%timestamp% >nul
    if %errorlevel% equ 0 (
        echo Backup created: trading.db.backup_%timestamp%
    )
    echo.
    echo Deleting existing database...
    del trading.db
)

echo Initializing database...
echo.

REM Run tf-engine init command
tf-engine.exe init

if %errorlevel% neq 0 (
    echo.
    echo ========================================
    echo  ERROR: Database initialization failed!
    echo ========================================
    echo.
    echo Check tf-engine.log for details.
    pause
    exit /b 1
)

echo.
echo ========================================
echo  Database Initialized Successfully!
echo ========================================
echo.
echo Database file: trading.db
echo.

REM Verify database exists
if exist "trading.db" (
    echo Database file created: %cd%\trading.db
    echo.

    REM Show database info
    echo Verifying database schema...
    tf-engine.exe get-settings --format json

    if %errorlevel% equ 0 (
        echo.
        echo Database schema verified OK
    )
) else (
    echo ERROR: Database file not created!
    pause
    exit /b 1
)

echo.
echo Default Settings:
echo   - Equity (E):           $10,000
echo   - Risk per trade (r):   0.75%%
echo   - Portfolio heat cap:   4%%
echo   - Bucket heat cap:      1.5%%
echo   - Stop multiple (K):    2
echo.
echo Next steps:
echo   1. Adjust settings if needed: tf-engine.exe set-setting --key Equity_E --value 20000
echo   2. Open TradingPlatform.xlsm in Excel
echo   3. Run VBA tests from "VBA Tests" worksheet
echo.

pause
