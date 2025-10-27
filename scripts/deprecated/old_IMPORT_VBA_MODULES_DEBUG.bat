@echo off
setlocal EnableExtensions EnableDelayedExpansion

echo ========================================
echo DEBUG: Starting batch file
echo ========================================

REM --- Config ---
set "WB=TrendFollowing_TradeEntry.xlsm"
set "VBAFOLDER=VBA"
set "SCRIPT=scripts\excel_build_repo_aware_logged.vbs"
set "LOGDIR=logs"

echo Current Directory: %CD%
echo Workbook: %WB%
echo VBA Folder: %VBAFOLDER%
echo Script: %SCRIPT%
echo Log Dir: %LOGDIR%

REM Check if folders exist
if exist "%VBAFOLDER%" (
    echo [OK] VBA folder exists
) else (
    echo [ERROR] VBA folder NOT found
)

if exist "%SCRIPT%" (
    echo [OK] VBScript file exists
) else (
    echo [ERROR] VBScript file NOT found at: %SCRIPT%
)

if not exist "%LOGDIR%" (
    echo Creating logs directory...
    mkdir "%LOGDIR%"
)

REM Robust timestamp
echo Generating timestamp...
for /f %%i in ('powershell -NoProfile -Command "(Get-Date).ToString('yyyyMMdd_HHmmss')"') do set "TS=%%i"
echo Timestamp: %TS%
set "LOGFILE=%LOGDIR%\build_%TS%.log"
echo Log file will be: %LOGFILE%

echo [INFO] Repo: %CD% > "%LOGFILE%"
echo [INFO] Log:  %LOGFILE% >> "%LOGFILE%"
echo [INFO] Script: %SCRIPT% >> "%LOGFILE%"

echo.
echo ========================================
echo About to run: Registry edits
echo ========================================
REM Trust Center toggles (best-effort)
reg add "HKCU\Software\Microsoft\Office\16.0\Excel\Security" /v AccessVBOM /t REG_DWORD /d 1 /f >> "%LOGFILE%" 2>&1
reg add "HKCU\Software\Microsoft\Office\16.0\Excel\Security" /v VBAWarnings /t REG_DWORD /d 1 /f >> "%LOGFILE%" 2>&1

echo.
echo ========================================
echo About to run: Kill Excel
echo ========================================
REM Kill any stray Excel
taskkill /IM excel.exe /F >> "%LOGFILE%" 2>&1

echo.
echo ========================================
echo About to run: VBScript
echo ========================================
REM Signal "automation" to any Workbook_Open logic
set "XL_SILENT_SETUP=1"
echo [INFO] XL_SILENT_SETUP=%XL_SILENT_SETUP% >> "%LOGFILE%"

echo Running: cscript //nologo "%SCRIPT%" "%CD%\%WB%" "%CD%\%VBAFOLDER%" "%CD%\%LOGFILE%"
echo.
REM Don't redirect to logfile - VBScript handles its own logging
cscript //nologo "%SCRIPT%" "%CD%\%WB%" "%CD%\%VBAFOLDER%" "%CD%\%LOGFILE%"
set "RC=%ERRORLEVEL%"
set "XL_SILENT_SETUP="

echo.
echo ========================================
echo VBScript Exit Code: %RC%
echo ========================================

echo [INFO] ExitCode: %RC% >> "%LOGFILE%"
if not "%RC%"=="0" (
  echo [ERROR] Import failed. See "%LOGFILE%"
  type "%LOGFILE%"
  pause
  exit /b %RC%
)

echo [OK] %WB% ready. See "%LOGFILE%"
type "%LOGFILE%"
pause
exit /b 0
