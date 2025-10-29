@echo off
REM ===========================================================================
REM fix-vba-modules.bat - Quick Fix for VBA Module Import Issues
REM ===========================================================================
REM Purpose: Fix "argument not optional" error by importing VBA modules
REM
REM Usage: Double-click this file or run from command prompt
REM ===========================================================================

echo ==========================================
echo  VBA Module Fix Tool
echo ==========================================
echo.
echo This will fix the "argument not optional" error
echo by importing the VBA modules into your workbook.
echo.

REM Check if workbook exists
if not exist "TradingPlatform.xlsm" (
    echo ERROR: TradingPlatform.xlsm not found
    echo.
    echo This script must be run from the same folder as TradingPlatform.xlsm
    echo.
    echo Current directory: %CD%
    echo.
    pause
    exit /b 1
)

REM Check if VBA files exist
if not exist "..\excel\vba\TFEngine.bas" (
    echo ERROR: VBA source files not found
    echo Expected location: ..\excel\vba\
    echo.
    pause
    exit /b 1
)

echo Workbook found: TradingPlatform.xlsm
echo VBA sources found: ..\excel\vba\
echo.
echo IMPORTANT: Please close TradingPlatform.xlsm if it's currently open!
echo.
pause

echo.
echo Running VBA import script...
echo.

REM Run the VBS script
cscript //nologo fix-vba-modules.vbs

if %errorlevel% neq 0 (
    echo.
    echo ERROR: VBA import failed
    echo.
    echo Common causes:
    echo 1. TradingPlatform.xlsm is still open - close it and try again
    echo 2. Macro security is blocking VBA project access
    echo 3. Excel Trust Center settings need adjustment
    echo.
    echo To fix macro security:
    echo 1. Open Excel
    echo 2. File ^> Options ^> Trust Center ^> Trust Center Settings
    echo 3. Macro Settings ^> Enable all macros
    echo 4. Check "Trust access to the VBA project object model"
    echo 5. Click OK and restart Excel
    echo.
    pause
    exit /b 1
)

echo.
echo ==========================================
echo  SUCCESS! VBA Modules Imported
echo ==========================================
echo.
echo You can now:
echo 1. Open TradingPlatform.xlsm
echo 2. Enable macros when prompted
echo 3. Test the macros - they should work now!
echo.

pause
