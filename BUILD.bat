@echo off
REM ========================================
REM Excel Trading Workbook - One-Click Build
REM ========================================

echo.
echo ========================================
echo Building Excel Trading Workbook
echo ========================================
echo.

REM First, kill any stuck Excel processes
taskkill /F /IM EXCEL.EXE >nul 2>&1

REM Run simplified Python build script
python build_workbook_simple.py

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo SUCCESS! Workbook created.
    echo ========================================
    echo.
    echo The workbook has been created with all VBA modules.
    echo.
    echo NEXT: Open TrendFollowing_TradeEntry.xlsm and run setup macros
    echo See instructions above or in START_HERE.md
    echo.
) else (
    echo.
    echo ========================================
    echo BUILD FAILED - See errors above
    echo ========================================
    echo.
)

pause
