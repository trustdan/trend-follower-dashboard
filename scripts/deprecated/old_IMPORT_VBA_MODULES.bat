@echo off
setlocal EnableExtensions EnableDelayedExpansion

REM --- Config ---
set "WB=TrendFollowing_TradeEntry.xlsm"
set "VBAFOLDER=VBA"
set "SCRIPT=scripts\excel_build_repo_aware_logged.vbs"
set "LOGDIR=logs"
if not exist "%LOGDIR%" mkdir "%LOGDIR%"

REM Robust timestamp
for /f %%i in ('powershell -NoProfile -Command "(Get-Date).ToString('yyyyMMdd_HHmmss')"') do set "TS=%%i"
set "LOGFILE=%LOGDIR%\build_%TS%.log"

echo [INFO] Repo: %CD% > "%LOGFILE%"
echo [INFO] Log:  %LOGFILE% >> "%LOGFILE%"
echo [INFO] Script: %SCRIPT% >> "%LOGFILE%"

REM Trust Center toggles (best-effort)
reg add "HKCU\Software\Microsoft\Office\16.0\Excel\Security" /v AccessVBOM /t REG_DWORD /d 1 /f >> "%LOGFILE%" 2>&1
reg add "HKCU\Software\Microsoft\Office\16.0\Excel\Security" /v VBAWarnings /t REG_DWORD /d 1 /f >> "%LOGFILE%" 2>&1

REM Kill any stray Excel
taskkill /IM excel.exe /F >> "%LOGFILE%" 2>&1

REM Signal “automation” to any Workbook_Open logic
set "XL_SILENT_SETUP=1"
echo [INFO] XL_SILENT_SETUP=%XL_SILENT_SETUP% >> "%LOGFILE%"

REM Build / import (VBScript handles its own logging to LOGFILE)
cscript //nologo "%SCRIPT%" "%CD%\%WB%" "%CD%\%VBAFOLDER%" "%CD%\%LOGFILE%"
set "RC=%ERRORLEVEL%"
set "XL_SILENT_SETUP="

echo [INFO] ExitCode: %RC% >> "%LOGFILE%"
if not "%RC%"=="0" (
  echo [ERROR] Import failed. See "%LOGFILE%"
  exit /b %RC%
)

echo [OK] %WB% ready. See "%LOGFILE%"
REM Optional: auto-open to verify
REM start "" "%CD%\%WB%"
exit /b 0
