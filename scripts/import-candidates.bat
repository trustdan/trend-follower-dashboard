@echo off
REM ===========================================================================
REM import-candidates.bat - Quick Launcher for Interactive Candidate Import
REM ===========================================================================
REM Purpose: One-click launcher for FINVIZ candidate import in interactive mode
REM
REM What this does:
REM   - Launches tf-engine.exe in interactive mode
REM   - Shows ASCII art and guided prompts
REM   - Walks you through screener selection and import
REM
REM Usage: Just double-click this file or run from command line
REM ===========================================================================

REM Set console to UTF-8 for emoji support
chcp 65001 >nul 2>&1

REM Set window title
title Trading Engine - Candidate Import

REM Check if tf-engine.exe exists
if not exist "tf-engine.exe" (
    echo ERROR: tf-engine.exe not found in current directory
    echo Please run this script from the windows\ folder
    pause
    exit /b 1
)

REM Check if database is initialized
if not exist "trading.db" (
    echo.
    echo ====================================================
    echo  WARNING: Database not initialized!
    echo ====================================================
    echo.
    echo The trading.db database doesn't exist yet.
    echo.
    echo Initializing database now...
    echo.
    tf-engine.exe init
    if %errorlevel% neq 0 (
        echo.
        echo ERROR: Database initialization failed!
        pause
        exit /b 1
    )
    echo.
    echo Database initialized successfully!
    echo.
)

REM Launch interactive mode
tf-engine.exe interactive

REM Pause to see results
echo.
pause
