@echo off
REM ===========================================================================
REM run-tests.bat - Automated Test Runner for Windows Integration
REM ===========================================================================
REM Purpose: Run smoke tests and report results
REM
REM Note: VBA unit tests must be run manually from Excel
REM       This script tests CLI functionality and environment setup
REM
REM Created: 2025-10-27 (M20 - Windows Integration Package)
REM ===========================================================================

setlocal enabledelayedexpansion

echo ========================================
echo  Windows Integration Test Runner
echo  Trading Engine v3
echo ========================================
echo.
echo Date: %date% %time%
echo Environment: Windows, Excel
echo.

REM Initialize counters
set PASS_COUNT=0
set FAIL_COUNT=0
set TOTAL_COUNT=0

REM Create results file
set RESULTS_FILE=test-results.txt
echo ======================================== > %RESULTS_FILE%
echo  Windows Integration Test Report >> %RESULTS_FILE%
echo ======================================== >> %RESULTS_FILE%
echo Date: %date% %time% >> %RESULTS_FILE%
echo Environment: Windows >> %RESULTS_FILE%
echo. >> %RESULTS_FILE%

REM ============================================================================
REM SMOKE TESTS
REM ============================================================================

echo [SMOKE TESTS]
echo [SMOKE TESTS] >> %RESULTS_FILE%
echo.

REM ---------------------------------------------------------------------------
REM Test 1: Engine Version
REM ---------------------------------------------------------------------------

set /a TOTAL_COUNT+=1
echo Test 1: Engine version check...

if not exist "tf-engine.exe" (
    echo [X] FAIL: tf-engine.exe not found
    echo [X] FAIL: Engine version check - tf-engine.exe not found >> %RESULTS_FILE%
    set /a FAIL_COUNT+=1
    goto :test2
)

tf-engine.exe --version > nul 2>&1
if %errorlevel% equ 0 (
    echo [✓] PASS: Engine version check
    echo [✓] PASS: Engine version check >> %RESULTS_FILE%
    set /a PASS_COUNT+=1
) else (
    echo [X] FAIL: Engine version check - exit code %errorlevel%
    echo [X] FAIL: Engine version check - exit code %errorlevel% >> %RESULTS_FILE%
    set /a FAIL_COUNT+=1
)

:test2

REM ---------------------------------------------------------------------------
REM Test 2: Database Exists
REM ---------------------------------------------------------------------------

set /a TOTAL_COUNT+=1
echo Test 2: Database existence check...

if exist "trading.db" (
    echo [✓] PASS: Database file exists
    echo [✓] PASS: Database file exists >> %RESULTS_FILE%
    set /a PASS_COUNT+=1
) else (
    echo [X] FAIL: Database file not found
    echo [X] FAIL: Database file not found >> %RESULTS_FILE%
    set /a FAIL_COUNT+=1
)

REM ---------------------------------------------------------------------------
REM Test 3: Get Settings
REM ---------------------------------------------------------------------------

set /a TOTAL_COUNT+=1
echo Test 3: Get settings from database...

if not exist "trading.db" (
    echo [X] FAIL: Get settings - database missing
    echo [X] FAIL: Get settings - database missing >> %RESULTS_FILE%
    set /a FAIL_COUNT+=1
    goto :test4
)

tf-engine.exe get-settings --format json > nul 2>&1
if %errorlevel% equ 0 (
    echo [✓] PASS: Get settings
    echo [✓] PASS: Get settings >> %RESULTS_FILE%
    set /a PASS_COUNT+=1
) else (
    echo [X] FAIL: Get settings - exit code %errorlevel%
    echo [X] FAIL: Get settings - exit code %errorlevel% >> %RESULTS_FILE%
    set /a FAIL_COUNT+=1
)

:test4

REM ---------------------------------------------------------------------------
REM Test 4: Position Sizing (Stock)
REM ---------------------------------------------------------------------------

set /a TOTAL_COUNT+=1
echo Test 4: Position sizing calculation...

tf-engine.exe size --entry 180 --atr 1.5 --k 2 --method stock --format json > size_output.tmp 2>&1
if %errorlevel% equ 0 (
    REM Verify JSON contains expected values (basic check)
    findstr /C:"shares" size_output.tmp > nul
    if %errorlevel% equ 0 (
        findstr /C:"25" size_output.tmp > nul
        if %errorlevel% equ 0 (
            echo [✓] PASS: Position sizing calculation
            echo [✓] PASS: Position sizing calculation >> %RESULTS_FILE%
            set /a PASS_COUNT+=1
        ) else (
            echo [X] FAIL: Position sizing - incorrect shares value
            echo [X] FAIL: Position sizing - incorrect shares value >> %RESULTS_FILE%
            set /a FAIL_COUNT+=1
        )
    ) else (
        echo [X] FAIL: Position sizing - invalid JSON output
        echo [X] FAIL: Position sizing - invalid JSON output >> %RESULTS_FILE%
        set /a FAIL_COUNT+=1
    )
) else (
    echo [X] FAIL: Position sizing - command failed
    echo [X] FAIL: Position sizing - command failed >> %RESULTS_FILE%
    set /a FAIL_COUNT+=1
)

if exist size_output.tmp del size_output.tmp

REM ---------------------------------------------------------------------------
REM Test 5: Checklist Command
REM ---------------------------------------------------------------------------

set /a TOTAL_COUNT+=1
echo Test 5: Checklist evaluation...

tf-engine.exe checklist --ticker AAPL --checks true,true,true,true,true,true --format json > checklist_output.tmp 2>&1
if %errorlevel% equ 0 (
    findstr /C:"GREEN" checklist_output.tmp > nul
    if %errorlevel% equ 0 (
        echo [✓] PASS: Checklist evaluation
        echo [✓] PASS: Checklist evaluation >> %RESULTS_FILE%
        set /a PASS_COUNT+=1
    ) else (
        echo [X] FAIL: Checklist - expected GREEN banner
        echo [X] FAIL: Checklist - expected GREEN banner >> %RESULTS_FILE%
        set /a FAIL_COUNT+=1
    )
) else (
    echo [X] FAIL: Checklist - command failed
    echo [X] FAIL: Checklist - command failed >> %RESULTS_FILE%
    set /a FAIL_COUNT+=1
)

if exist checklist_output.tmp del checklist_output.tmp

REM ---------------------------------------------------------------------------
REM Test 6: Heat Check
REM ---------------------------------------------------------------------------

set /a TOTAL_COUNT+=1
echo Test 6: Heat check...

tf-engine.exe heat --add-r 75 --bucket "Tech/Comm" --format json > nul 2>&1
if %errorlevel% equ 0 (
    echo [✓] PASS: Heat check
    echo [✓] PASS: Heat check >> %RESULTS_FILE%
    set /a PASS_COUNT+=1
) else (
    echo [X] FAIL: Heat check - command failed
    echo [X] FAIL: Heat check - command failed >> %RESULTS_FILE%
    set /a FAIL_COUNT+=1
)

REM ---------------------------------------------------------------------------
REM Test 7: Import Candidates
REM ---------------------------------------------------------------------------

set /a TOTAL_COUNT+=1
echo Test 7: Import candidates...

tf-engine.exe import-candidates --tickers "TEST1,TEST2,TEST3" --preset TEST --format json > nul 2>&1
if %errorlevel% equ 0 (
    echo [✓] PASS: Import candidates
    echo [✓] PASS: Import candidates >> %RESULTS_FILE%
    set /a PASS_COUNT+=1
) else (
    echo [X] FAIL: Import candidates - command failed
    echo [X] FAIL: Import candidates - command failed >> %RESULTS_FILE%
    set /a FAIL_COUNT+=1
)

REM ---------------------------------------------------------------------------
REM Test 8: List Candidates
REM ---------------------------------------------------------------------------

set /a TOTAL_COUNT+=1
echo Test 8: List candidates...

tf-engine.exe list-candidates --format json > nul 2>&1
if %errorlevel% equ 0 (
    echo [✓] PASS: List candidates
    echo [✓] PASS: List candidates >> %RESULTS_FILE%
    set /a PASS_COUNT+=1
) else (
    echo [X] FAIL: List candidates - command failed
    echo [X] FAIL: List candidates - command failed >> %RESULTS_FILE%
    set /a FAIL_COUNT+=1
)

REM ============================================================================
REM FILE CHECKS
REM ============================================================================

echo.
echo [FILE CHECKS]
echo [FILE CHECKS] >> %RESULTS_FILE%
echo.

REM ---------------------------------------------------------------------------
REM Test 9: VBA Modules Exist
REM ---------------------------------------------------------------------------

set /a TOTAL_COUNT+=1
echo Test 9: VBA module files exist...

set VBA_MODULES_FOUND=1
if not exist "..\excel\vba\TFTypes.bas" set VBA_MODULES_FOUND=0
if not exist "..\excel\vba\TFHelpers.bas" set VBA_MODULES_FOUND=0
if not exist "..\excel\vba\TFEngine.bas" set VBA_MODULES_FOUND=0
if not exist "..\excel\vba\TFTests.bas" set VBA_MODULES_FOUND=0

if %VBA_MODULES_FOUND% equ 1 (
    echo [✓] PASS: All VBA module files exist
    echo [✓] PASS: All VBA module files exist >> %RESULTS_FILE%
    set /a PASS_COUNT+=1
) else (
    echo [X] FAIL: VBA module files missing
    echo [X] FAIL: VBA module files missing >> %RESULTS_FILE%
    set /a FAIL_COUNT+=1
)

REM ---------------------------------------------------------------------------
REM Test 10: Test Data Files Exist
REM ---------------------------------------------------------------------------

set /a TOTAL_COUNT+=1
echo Test 10: Test data files exist...

if exist "test-data" (
    dir /b test-data\*.json > nul 2>&1
    if %errorlevel% equ 0 (
        echo [✓] PASS: Test data files exist
        echo [✓] PASS: Test data files exist >> %RESULTS_FILE%
        set /a PASS_COUNT+=1
    ) else (
        echo [X] FAIL: No JSON files in test-data folder
        echo [X] FAIL: No JSON files in test-data folder >> %RESULTS_FILE%
        set /a FAIL_COUNT+=1
    )
) else (
    echo [X] FAIL: test-data folder not found
    echo [X] FAIL: test-data folder not found >> %RESULTS_FILE%
    set /a FAIL_COUNT+=1
)

REM ---------------------------------------------------------------------------
REM Test 11: Log Files Writable
REM ---------------------------------------------------------------------------

set /a TOTAL_COUNT+=1
echo Test 11: Log files writable...

echo Test > test_write.tmp 2>&1
if %errorlevel% equ 0 (
    del test_write.tmp
    echo [✓] PASS: Directory writable for logs
    echo [✓] PASS: Directory writable for logs >> %RESULTS_FILE%
    set /a PASS_COUNT+=1
) else (
    echo [X] FAIL: Cannot write to directory
    echo [X] FAIL: Cannot write to directory >> %RESULTS_FILE%
    set /a FAIL_COUNT+=1
)

REM ============================================================================
REM SUMMARY
REM ============================================================================

echo.
echo ========================================
echo  TEST SUMMARY
echo ========================================
echo Total Tests:  %TOTAL_COUNT%
echo Passed:       %PASS_COUNT%
echo Failed:       %FAIL_COUNT%
echo.

echo. >> %RESULTS_FILE%
echo ======================================== >> %RESULTS_FILE%
echo  TEST SUMMARY >> %RESULTS_FILE%
echo ======================================== >> %RESULTS_FILE%
echo Total Tests:  %TOTAL_COUNT% >> %RESULTS_FILE%
echo Passed:       %PASS_COUNT% >> %RESULTS_FILE%
echo Failed:       %FAIL_COUNT% >> %RESULTS_FILE%
echo. >> %RESULTS_FILE%

if %FAIL_COUNT% equ 0 (
    echo RESULT: ALL TESTS PASSED ✅
    echo RESULT: ALL TESTS PASSED >> %RESULTS_FILE%
    echo.
    echo Next steps:
    echo   1. Open TradingPlatform.xlsm in Excel
    echo   2. Run VBA unit tests from "VBA Tests" worksheet
    echo   3. Complete integration tests per WINDOWS_TESTING.md Phase 4
) else (
    echo RESULT: %FAIL_COUNT% TESTS FAILED ❌
    echo RESULT: %FAIL_COUNT% TESTS FAILED >> %RESULTS_FILE%
    echo.
    echo Please review failures above and check:
    echo   - tf-engine.exe is correct Windows binary
    echo   - Database initialized correctly
    echo   - All prerequisite files present
    echo   - Check tf-engine.log for errors
)

echo.
echo Test results saved to: %RESULTS_FILE%
echo.

REM ============================================================================
REM VBA TESTING REMINDER
REM ============================================================================

echo ========================================
echo  MANUAL VBA TESTS REQUIRED
echo ========================================
echo.
echo The automated tests above verify CLI functionality.
echo VBA integration tests must be run manually:
echo.
echo   1. Open TradingPlatform.xlsm in Excel
echo   2. Enable macros when prompted
echo   3. Go to "VBA Tests" worksheet
echo   4. Click "Run All Tests" button
echo   5. Verify all 14 tests pass
echo.
echo See WINDOWS_TESTING.md for complete testing procedures.
echo.

pause
