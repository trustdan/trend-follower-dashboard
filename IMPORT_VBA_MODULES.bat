@echo off
REM ============================================================================
REM VBA Module Import Helper Script (Windows Batch)
REM ============================================================================
REM
REM This script helps you import VBA modules into Excel.
REM
REM Usage:
REM   1. Double-click this file from Windows Explorer
REM   2. Follow the prompts
REM
REM OR run from command line:
REM   IMPORT_VBA_MODULES.bat [path_to_workbook.xlsm]
REM
REM ============================================================================

echo.
echo ========================================================================
echo VBA Module Import Helper
echo ========================================================================
echo.

REM Check if Python is available
where python >nul 2>nul
if %errorlevel% == 0 (
    echo Python found! Checking for required packages...
    echo.

    REM Try Python automation
    python --version
    echo.
    echo Attempting automated import via Python...
    echo.
    python import_to_excel.py %1

    if %errorlevel% == 0 (
        echo.
        echo ========================================================================
        echo SUCCESS! VBA modules imported.
        echo ========================================================================
        echo.
        echo Next steps:
        echo   1. Open the Excel workbook that was created/updated
        echo   2. Press Alt+F11 to open VBA Editor
        echo   3. Press Ctrl+G to open Immediate Window
        echo   4. Type: EnsureStructure
        echo   5. Press Enter to create all sheets and tables
        echo   6. See VBA_SETUP_GUIDE.md Part 2 to build the UI
        echo.
        pause
        exit /b 0
    ) else (
        echo.
        echo Python automation failed. See manual instructions below.
        echo.
    )
) else (
    echo Python not found. Using manual import method.
    echo.
)

REM ============================================================================
REM Manual Import Instructions
REM ============================================================================

echo ========================================================================
echo MANUAL IMPORT INSTRUCTIONS
echo ========================================================================
echo.
echo Python automation is not available. Please import modules manually:
echo.
echo 1. Open Excel and create a new workbook
echo 2. Save as: TrendFollowing_TradeEntry.xlsm (enable macros)
echo 3. Press Alt+F11 to open VBA Editor
echo 4. For EACH file in the VBA\ folder:
echo.
echo    Standard Modules (.bas files):
echo      - File -^> Import File
echo      - Select: VBA\TF_Utils.bas
echo      - Repeat for: TF_Data.bas, TF_UI.bas, TF_Presets.bas
echo.
echo    Class Modules (.cls files):
echo      - File -^> Import File
echo      - Select: VBA\ThisWorkbook.cls (replaces existing)
echo      - Select: VBA\Sheet_TradeEntry.cls (add new)
echo.
echo 5. After importing, in Immediate Window (Ctrl+G):
echo    Type: EnsureStructure
echo    Press: Enter
echo.
echo 6. Build the TradeEntry UI:
echo    - See VBA_SETUP_GUIDE.md Part 2 for detailed steps
echo    - Add labels, buttons, checkboxes, dropdowns
echo.
echo ========================================================================
echo.
echo Full manual guide: VBA_SETUP_GUIDE.md
echo.
echo If you have Python installed, try:
echo   pip install pywin32
echo Then run this script again for automated import.
echo.
echo ========================================================================
pause
