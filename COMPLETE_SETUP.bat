@echo off
echo ========================================
echo COMPLETE WORKBOOK SETUP
echo ========================================
echo.
echo This will build the workbook AND add all enhancements:
echo - Python FINVIZ auto-scraper
echo - Python heat calculator
echo - UI buttons and formatting
echo - Enhanced workflows
echo.
pause

REM Step 1: Build basic workbook
echo.
echo ========================================
echo Step 1: Building workbook with VBA modules
echo ========================================
call BUILD_WITH_PYTHON.bat
if errorlevel 1 (
    echo ERROR: Workbook build failed
    pause
    exit /b 1
)

echo.
echo ========================================
echo Step 2: Opening workbook to add enhancements
echo ========================================

REM Step 2: Run post-build enhancements
echo Running TF_UI.InitializeUI to build buttons and UI...
cscript //nologo scripts\run_postbuild.vbs "TrendFollowing_TradeEntry.xlsm"

if errorlevel 1 (
    echo WARNING: Post-build enhancements may have failed
    echo Workbook is functional but may be missing UI elements
) else (
    echo SUCCESS: UI enhancements applied
)

echo.
echo ========================================
echo Step 3: Python Integration Setup
echo ========================================
echo.
echo The Python modules (finviz_scraper.py, heat_calculator.py) are ready.
echo To use Python features in Excel:
echo.
echo 1. Open Excel (Microsoft 365 required)
echo 2. File → Options → Trust Center → Python Python Settings
echo 3. Enable "Python in Excel"
echo.
echo Then in the workbook:
echo 4. Run macro: TestPythonIntegration (to verify Python works)
echo 5. Use "Auto Import" buttons (Python-powered FINVIZ scraping)
echo.
echo See PYTHON_SETUP_GUIDE.md for detailed instructions.
echo.

echo ========================================
echo SETUP COMPLETE!
echo ========================================
echo.
echo Workbook: TrendFollowing_TradeEntry.xlsm
echo.
echo What you have:
echo - [x] 8 Sheets with data
echo - [x] 5 Tables (Presets, Buckets, Candidates, Decisions, Positions)
echo - [x] 7 Named Ranges
echo - [x] 9 VBA Modules
echo - [x] Python modules ready (optional)
echo.
echo Next steps:
echo 1. Open TrendFollowing_TradeEntry.xlsm
echo 2. Press Alt+F11 to view VBA modules
echo 3. Run TF_UI.InitializeUI if UI needs adjustment
echo 4. (Optional) Enable Python in Excel for auto-scraping
echo.
pause
