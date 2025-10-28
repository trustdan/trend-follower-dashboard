@echo off
REM =============================================================================
REM run-integration-tests.bat - Automated Phase 4 Integration Test Runner
REM =============================================================================
REM Purpose: One-click automated execution of all 25 Phase 4 integration tests
REM
REM What this script does:
REM   1. Verifies test environment (files exist, engine works)
REM   2. Imports TFIntegrationTests.bas module into workbook
REM   3. Runs automated tests via Excel macro
REM   4. Opens results worksheet
REM   5. Opens log file
REM
REM Usage:
REM   run-integration-tests.bat
REM
REM Requirements:
REM   - TradingPlatform.xlsm exists
REM   - tf-engine.exe exists and accessible
REM   - trading.db exists
REM   - Excel installed
REM
REM Output:
REM   - Integration Tests worksheet in workbook
REM   - logs/integration-tests-YYYYMMDD-HHMMSS.log
REM
REM Created: 2025-10-27 (M21 - Phase 4 Automation)
REM =============================================================================

setlocal enabledelayedexpansion

echo ========================================
echo M21 PHASE 4 INTEGRATION TESTS
echo Automated Test Runner
echo ========================================
echo.

REM -----------------------------------------------------------------------------
REM ENVIRONMENT CHECK
REM -----------------------------------------------------------------------------

echo [1/5] Checking environment...

if not exist "TradingPlatform.xlsm" (
    echo ERROR: TradingPlatform.xlsm not found
    echo        Run setup-all.bat first
    pause
    exit /b 1
)

if not exist "tf-engine.exe" (
    echo ERROR: tf-engine.exe not found
    pause
    exit /b 1
)

if not exist "trading.db" (
    echo ERROR: trading.db not found
    echo        Run setup-all.bat first
    pause
    exit /b 1
)

echo   - TradingPlatform.xlsm found
echo   - tf-engine.exe found
echo   - trading.db found
echo.

REM -----------------------------------------------------------------------------
REM VERIFY ENGINE
REM -----------------------------------------------------------------------------

echo [2/5] Verifying engine...

tf-engine.exe --version >nul 2>&1
if errorlevel 1 (
    echo ERROR: Engine not working
    echo        Try rebuilding with: go build
    pause
    exit /b 1
)

echo   - Engine accessible
echo.

REM -----------------------------------------------------------------------------
REM IMPORT TEST MODULE
REM -----------------------------------------------------------------------------

echo [3/5] Importing TFIntegrationTests module...

REM Check if module file exists
if not exist "..\excel\vba\TFIntegrationTests.bas" (
    echo ERROR: TFIntegrationTests.bas not found
    echo        Expected at: ..\excel\vba\TFIntegrationTests.bas
    pause
    exit /b 1
)

REM Create VBScript to import module
echo ' Import TFIntegrationTests module > import-test-module.vbs
echo Set objExcel = CreateObject("Excel.Application") >> import-test-module.vbs
echo objExcel.Visible = False >> import-test-module.vbs
echo objExcel.DisplayAlerts = False >> import-test-module.vbs
echo. >> import-test-module.vbs
echo ' Open workbook >> import-test-module.vbs
echo Set objWorkbook = objExcel.Workbooks.Open("%CD%\TradingPlatform.xlsm") >> import-test-module.vbs
echo. >> import-test-module.vbs
echo ' Access VBA project >> import-test-module.vbs
echo Set objVBProj = objWorkbook.VBProject >> import-test-module.vbs
echo. >> import-test-module.vbs
echo ' Remove existing module if present >> import-test-module.vbs
echo On Error Resume Next >> import-test-module.vbs
echo objVBProj.VBComponents.Remove objVBProj.VBComponents("TFIntegrationTests") >> import-test-module.vbs
echo On Error Goto 0 >> import-test-module.vbs
echo. >> import-test-module.vbs
echo ' Import new module >> import-test-module.vbs
echo objVBProj.VBComponents.Import "%CD%\..\excel\vba\TFIntegrationTests.bas" >> import-test-module.vbs
echo. >> import-test-module.vbs
echo ' Save and close >> import-test-module.vbs
echo objWorkbook.Save >> import-test-module.vbs
echo objWorkbook.Close >> import-test-module.vbs
echo objExcel.Quit >> import-test-module.vbs
echo. >> import-test-module.vbs
echo WScript.Echo "Module imported successfully" >> import-test-module.vbs

REM Run import script
cscript //nologo import-test-module.vbs
if errorlevel 1 (
    echo ERROR: Failed to import test module
    echo        Check that "Trust access to VBA project" is enabled
    echo        File ^> Options ^> Trust Center ^> Trust Center Settings
    echo        ^> Macro Settings ^> Trust access to VBA project object model
    pause
    del import-test-module.vbs
    exit /b 1
)

del import-test-module.vbs
echo   - TFIntegrationTests module imported
echo.

REM -----------------------------------------------------------------------------
REM CREATE LOGS FOLDER
REM -----------------------------------------------------------------------------

echo [4/5] Preparing logs folder...

if not exist "logs" (
    mkdir logs
    echo   - logs\ folder created
) else (
    echo   - logs\ folder exists
)
echo.

REM -----------------------------------------------------------------------------
REM RUN TESTS
REM -----------------------------------------------------------------------------

echo [5/5] Running integration tests...
echo.
echo This will take 1-2 minutes...
echo Tests are running in Excel (background)...
echo.

REM Create VBScript to run tests
echo ' Run integration tests > run-tests.vbs
echo Set objExcel = CreateObject("Excel.Application") >> run-tests.vbs
echo objExcel.Visible = True >> run-tests.vbs
echo objExcel.DisplayAlerts = False >> run-tests.vbs
echo. >> run-tests.vbs
echo ' Open workbook >> run-tests.vbs
echo Set objWorkbook = objExcel.Workbooks.Open("%CD%\TradingPlatform.xlsm") >> run-tests.vbs
echo. >> run-tests.vbs
echo ' Run tests >> run-tests.vbs
echo objExcel.Run "TFIntegrationTests.RunAllIntegrationTests" >> run-tests.vbs
echo. >> run-tests.vbs
echo ' Keep workbook open for review >> run-tests.vbs
echo WScript.Echo "Tests complete - check Integration Tests worksheet" >> run-tests.vbs

REM Run tests
cscript //nologo run-tests.vbs
if errorlevel 1 (
    echo ERROR: Test execution failed
    echo        Check TradingSystem_Debug.log for details
    pause
    del run-tests.vbs
    exit /b 1
)

del run-tests.vbs
echo.

REM -----------------------------------------------------------------------------
REM OPEN RESULTS
REM -----------------------------------------------------------------------------

echo ========================================
echo TESTS COMPLETE
echo ========================================
echo.
echo Results:
echo   - Check "Integration Tests" worksheet in Excel
echo   - Check logs\integration-tests-*.log file
echo.

REM Find most recent log file
for /f "delims=" %%f in ('dir /b /o-d logs\integration-tests-*.log 2^>nul') do (
    set "latestlog=%%f"
    goto :foundlog
)
:foundlog

if defined latestlog (
    echo Opening log file: logs\!latestlog!
    echo.
    start notepad "logs\!latestlog!"
) else (
    echo WARNING: No log file found in logs\
)

echo Excel is open with test results
echo.
echo Press any key to exit...
pause >nul

exit /b 0
