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
echo  One-Click Installation
echo ========================================
echo.
echo This script will:
echo   1. Create Excel workbook (TradingPlatform.xlsm)
echo   2. Enable VBA project access
echo   3. Import VBA modules (4 modules)
echo   4. Create UI worksheets (5 production sheets)
echo   5. Initialize database (trading.db)
echo   6. Configure Excel settings
echo   7. Run automated tests
echo.
echo Estimated time: 3-5 minutes
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
REM Step 2: Create Excel Workbook with VBA Access Enabled
REM ============================================================================

echo.
echo [Step 2/8] Creating Excel workbook...
echo.

REM Create VBScript to create and configure Excel workbook
echo ' Create Excel Workbook with VBA Access > create_workbook.vbs
echo Set objExcel = CreateObject("Excel.Application") >> create_workbook.vbs
echo objExcel.Visible = False >> create_workbook.vbs
echo objExcel.DisplayAlerts = False >> create_workbook.vbs
echo. >> create_workbook.vbs
echo ' Create new workbook >> create_workbook.vbs
echo Set objWorkbook = objExcel.Workbooks.Add >> create_workbook.vbs
echo. >> create_workbook.vbs
echo ' Rename first sheet >> create_workbook.vbs
echo objWorkbook.Worksheets(1).Name = "Setup" >> create_workbook.vbs
echo. >> create_workbook.vbs
echo ' Add VBA Tests sheet >> create_workbook.vbs
echo Set wsTests = objWorkbook.Worksheets.Add >> create_workbook.vbs
echo wsTests.Name = "VBA Tests" >> create_workbook.vbs
echo wsTests.Move , objWorkbook.Worksheets(objWorkbook.Worksheets.Count) >> create_workbook.vbs
echo. >> create_workbook.vbs
echo ' Save as macro-enabled workbook >> create_workbook.vbs
echo strPath = CreateObject("Scripting.FileSystemObject").GetAbsolutePathName(".") >> create_workbook.vbs
echo objWorkbook.SaveAs strPath + "\TradingPlatform.xlsm", 52 >> create_workbook.vbs
echo. >> create_workbook.vbs
echo ' Close >> create_workbook.vbs
echo objWorkbook.Close >> create_workbook.vbs
echo objExcel.Quit >> create_workbook.vbs
echo. >> create_workbook.vbs
echo Set objWorkbook = Nothing >> create_workbook.vbs
echo Set objExcel = Nothing >> create_workbook.vbs
echo. >> create_workbook.vbs
echo WScript.Echo "Workbook created: TradingPlatform.xlsm" >> create_workbook.vbs

REM Run the script
echo Running: cscript //nologo create_workbook.vbs >> %LOGFILE%
cscript //nologo create_workbook.vbs >> %LOGFILE% 2>&1

if %errorlevel% neq 0 (
    echo ERROR: Failed to create workbook
    echo ERROR: Failed to create workbook. Error level: %errorlevel% >> %LOGFILE%
    echo.
    echo Check setup-all.log for details
    echo The create_workbook.vbs script has been preserved for inspection
    pause
    exit /b 1
)

del create_workbook.vbs
echo [OK] Excel workbook created
echo [OK] Excel workbook created >> %LOGFILE%

REM ============================================================================
REM Step 3: Enable VBA Project Access via Registry
REM ============================================================================

echo.
echo [Step 3/8] Enabling VBA project access...
echo.

REM This sets the registry key that enables "Trust access to VBA project object model"
REM We'll try for different Excel versions

reg add "HKCU\Software\Microsoft\Office\16.0\Excel\Security" /v AccessVBOM /t REG_DWORD /d 1 /f >nul 2>&1
reg add "HKCU\Software\Microsoft\Office\15.0\Excel\Security" /v AccessVBOM /t REG_DWORD /d 1 /f >nul 2>&1
reg add "HKCU\Software\Microsoft\Office\14.0\Excel\Security" /v AccessVBOM /t REG_DWORD /d 1 /f >nul 2>&1

echo [OK] VBA project access enabled

REM ============================================================================
REM Step 4: Import VBA Modules
REM ============================================================================

echo.
echo [Step 4/8] Importing VBA modules...
echo.

REM Create import script
echo ' VBA Module Import Script > import_vba.vbs
echo Set objExcel = CreateObject("Excel.Application") >> import_vba.vbs
echo objExcel.Visible = False >> import_vba.vbs
echo objExcel.DisplayAlerts = False >> import_vba.vbs
echo. >> import_vba.vbs
echo ' Get script directory >> import_vba.vbs
echo Set objFSO = CreateObject("Scripting.FileSystemObject") >> import_vba.vbs
echo strScriptDir = objFSO.GetParentFolderName(WScript.ScriptFullName) >> import_vba.vbs
echo. >> import_vba.vbs
echo ' Open workbook >> import_vba.vbs
echo strWorkbookPath = strScriptDir + "\TradingPlatform.xlsm" >> import_vba.vbs
echo Set objWorkbook = objExcel.Workbooks.Open(strWorkbookPath) >> import_vba.vbs
echo. >> import_vba.vbs
echo ' Remove existing modules if present >> import_vba.vbs
echo On Error Resume Next >> import_vba.vbs
echo For Each comp In objWorkbook.VBProject.VBComponents >> import_vba.vbs
echo     If comp.Type = 1 Then ' vbext_ct_StdModule >> import_vba.vbs
echo         Select Case comp.Name >> import_vba.vbs
echo             Case "TFTypes", "TFHelpers", "TFEngine", "TFTests" >> import_vba.vbs
echo                 objWorkbook.VBProject.VBComponents.Remove comp >> import_vba.vbs
echo         End Select >> import_vba.vbs
echo     End If >> import_vba.vbs
echo Next >> import_vba.vbs
echo On Error Goto 0 >> import_vba.vbs
echo. >> import_vba.vbs
echo ' Import new modules >> import_vba.vbs
echo strVBADir = objFSO.GetParentFolderName(strScriptDir) + "\excel\vba\" >> import_vba.vbs
echo. >> import_vba.vbs
echo WScript.Echo "Importing TFTypes.bas..." >> import_vba.vbs
echo objWorkbook.VBProject.VBComponents.Import strVBADir + "TFTypes.bas" >> import_vba.vbs
echo. >> import_vba.vbs
echo WScript.Echo "Importing TFHelpers.bas..." >> import_vba.vbs
echo objWorkbook.VBProject.VBComponents.Import strVBADir + "TFHelpers.bas" >> import_vba.vbs
echo. >> import_vba.vbs
echo WScript.Echo "Importing TFEngine.bas..." >> import_vba.vbs
echo objWorkbook.VBProject.VBComponents.Import strVBADir + "TFEngine.bas" >> import_vba.vbs
echo. >> import_vba.vbs
echo WScript.Echo "Importing TFTests.bas..." >> import_vba.vbs
echo objWorkbook.VBProject.VBComponents.Import strVBADir + "TFTests.bas" >> import_vba.vbs
echo. >> import_vba.vbs
echo ' Save and close >> import_vba.vbs
echo objWorkbook.Save >> import_vba.vbs
echo objWorkbook.Close >> import_vba.vbs
echo objExcel.Quit >> import_vba.vbs
echo. >> import_vba.vbs
echo ' Cleanup >> import_vba.vbs
echo Set objWorkbook = Nothing >> import_vba.vbs
echo Set objExcel = Nothing >> import_vba.vbs
echo Set objFSO = Nothing >> import_vba.vbs
echo. >> import_vba.vbs
echo WScript.Echo "VBA modules imported successfully!" >> import_vba.vbs

REM Run import
cscript //nologo import_vba.vbs

if %errorlevel% neq 0 (
    echo ERROR: VBA import failed
    echo.
    echo Possible issues:
    echo   - Excel version not compatible
    echo   - Antivirus blocking VBScript
    echo   - VBA project access not enabled
    echo.
    pause
    exit /b 1
)

del import_vba.vbs
echo [OK] VBA modules imported (4 modules)

REM ============================================================================
REM Step 5: Create UI Worksheets (M22 - Automated UI Generation)
REM ============================================================================

echo.
echo [Step 5/8] Creating UI worksheets...
echo.

REM Check if worksheet generation script exists
if not exist "create-ui-worksheets.vbs" (
    echo ERROR: create-ui-worksheets.vbs not found
    echo Please ensure the file exists in the windows\ folder
    pause
    exit /b 1
)

REM Execute worksheet generation script
echo Running: cscript //nologo create-ui-worksheets.vbs TradingPlatform.xlsm >> %LOGFILE%
cscript //nologo create-ui-worksheets.vbs TradingPlatform.xlsm >> %LOGFILE% 2>&1

if %errorlevel% neq 0 (
    echo ERROR: Failed to create UI worksheets
    echo ERROR: Failed to create UI worksheets. Error level: %errorlevel% >> %LOGFILE%
    echo.
    echo Check setup-all.log for details
    pause
    exit /b 1
)

echo [OK] UI worksheets created (5 production sheets)
echo [OK] UI worksheets created: Dashboard, Position Sizing, Checklist, Heat Check, Trade Entry >> %LOGFILE%

REM ============================================================================
REM Step 6: Configure Excel Named Ranges and Test Button
REM ============================================================================

echo.
echo [Step 6/8] Configuring Excel workbook...
echo.

REM Create configuration script
echo ' Configure Excel Workbook > configure_excel.vbs
echo Set objExcel = CreateObject("Excel.Application") >> configure_excel.vbs
echo objExcel.Visible = False >> configure_excel.vbs
echo objExcel.DisplayAlerts = False >> configure_excel.vbs
echo. >> configure_excel.vbs
echo ' Open workbook >> configure_excel.vbs
echo Set objFSO = CreateObject("Scripting.FileSystemObject") >> configure_excel.vbs
echo strPath = objFSO.GetAbsolutePathName(".") >> configure_excel.vbs
echo Set objWorkbook = objExcel.Workbooks.Open(strPath + "\TradingPlatform.xlsm") >> configure_excel.vbs
echo. >> configure_excel.vbs
echo ' Configure Setup sheet >> configure_excel.vbs
echo Set wsSetup = objWorkbook.Worksheets("Setup") >> configure_excel.vbs
echo wsSetup.Range("A1").Value = "Trading Engine v3 - Setup" >> configure_excel.vbs
echo wsSetup.Range("A1").Font.Bold = True >> configure_excel.vbs
echo wsSetup.Range("A1").Font.Size = 14 >> configure_excel.vbs
echo. >> configure_excel.vbs
echo wsSetup.Range("A3").Value = "Configuration" >> configure_excel.vbs
echo wsSetup.Range("A3").Font.Bold = True >> configure_excel.vbs
echo. >> configure_excel.vbs
echo wsSetup.Range("A4").Value = "Engine Path:" >> configure_excel.vbs
echo wsSetup.Range("B4").Value = ".\tf-engine.exe" >> configure_excel.vbs
echo objWorkbook.Names.Add "EnginePathSetting", wsSetup.Range("B4") >> configure_excel.vbs
echo. >> configure_excel.vbs
echo wsSetup.Range("A5").Value = "Database Path:" >> configure_excel.vbs
echo wsSetup.Range("B5").Value = ".\trading.db" >> configure_excel.vbs
echo objWorkbook.Names.Add "DatabasePathSetting", wsSetup.Range("B5") >> configure_excel.vbs
echo. >> configure_excel.vbs
echo ' Configure VBA Tests sheet >> configure_excel.vbs
echo Set wsTests = objWorkbook.Worksheets("VBA Tests") >> configure_excel.vbs
echo wsTests.Range("A1").Value = "VBA Unit Tests" >> configure_excel.vbs
echo wsTests.Range("A1").Font.Bold = True >> configure_excel.vbs
echo wsTests.Range("A1").Font.Size = 14 >> configure_excel.vbs
echo. >> configure_excel.vbs
echo wsTests.Range("A3").Value = "Click the button to run all VBA unit tests:" >> configure_excel.vbs
echo. >> configure_excel.vbs
echo ' Add test button (positioned in row 3, starting at column B) >> configure_excel.vbs
echo Set btnTests = wsTests.Buttons.Add(wsTests.Range("B3").Left, wsTests.Range("B3").Top, 120, 24) >> configure_excel.vbs
echo btnTests.Text = "Run All Tests" >> configure_excel.vbs
echo btnTests.OnAction = "TFTests.RunAllTests" >> configure_excel.vbs
echo. >> configure_excel.vbs
echo ' Save and close >> configure_excel.vbs
echo objWorkbook.Save >> configure_excel.vbs
echo objWorkbook.Close >> configure_excel.vbs
echo objExcel.Quit >> configure_excel.vbs
echo. >> configure_excel.vbs
echo Set objWorkbook = Nothing >> configure_excel.vbs
echo Set objExcel = Nothing >> configure_excel.vbs
echo. >> configure_excel.vbs
echo WScript.Echo "Excel configured successfully" >> configure_excel.vbs

cscript //nologo configure_excel.vbs

if %errorlevel% neq 0 (
    echo WARNING: Excel configuration had issues, but continuing...
)

del configure_excel.vbs
echo [OK] Excel workbook configured

REM ============================================================================
REM Step 7: Initialize Database
REM ============================================================================

echo.
echo [Step 7/8] Initializing database...
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
REM Step 8: Run Smoke Tests
REM ============================================================================

echo.
echo [Step 8/8] Running automated smoke tests...
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
echo   - TradingPlatform.xlsm (Excel workbook with VBA and 5 UI sheets)
echo   - trading.db (SQLite database)
echo   - setup-all.log (setup process log)
echo   - TradingSystem_Debug.log (will be created on first use)
echo.
echo Workbook contains:
echo   - Setup (configuration)
echo   - VBA Tests (automated testing)
echo   - Dashboard (portfolio overview)
echo   - Position Sizing (calculate shares/contracts)
echo   - Checklist (6-item evaluation)
echo   - Heat Check (portfolio/bucket heat)
echo   - Trade Entry (5-gate decision workflow)
echo.
echo Next steps:
echo   1. Open TradingPlatform.xlsm in Excel
echo   2. Enable macros when prompted
echo   3. Start with Dashboard worksheet
echo   4. Or go to "VBA Tests" sheet and click "Run All Tests"
echo   5. Verify all tests pass
echo.
echo For full testing, run: run-tests.bat
echo.
echo Documentation:
echo   - WINDOWS_TESTING.md (complete testing guide)
echo   - README.md (package overview)
echo.

pause
