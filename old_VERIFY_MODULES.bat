@echo off
echo ========================================
echo Verify: Check VBA Modules in Workbook
echo ========================================

if not exist "TrendFollowing_TradeEntry.xlsm" (
    echo ERROR: Workbook not found
    echo Please run IMPORT_VBA_MODULES_DEBUG.bat first
    pause
    exit /b 1
)

echo.
echo Found: TrendFollowing_TradeEntry.xlsm
echo Opening workbook to check modules...
echo.

REM Create a temporary VBScript to list modules
echo Option Explicit > temp_verify.vbs
echo Dim xl, wb, vbcomp, count >> temp_verify.vbs
echo Set xl = CreateObject("Excel.Application") >> temp_verify.vbs
echo xl.Visible = False >> temp_verify.vbs
echo xl.DisplayAlerts = False >> temp_verify.vbs
echo Set wb = xl.Workbooks.Open("%CD%\TrendFollowing_TradeEntry.xlsm") >> temp_verify.vbs
echo. >> temp_verify.vbs
echo WScript.Echo "VBA Modules in workbook:" >> temp_verify.vbs
echo WScript.Echo "=========================" >> temp_verify.vbs
echo count = 0 >> temp_verify.vbs
echo For Each vbcomp In wb.VBProject.VBComponents >> temp_verify.vbs
echo     WScript.Echo vbcomp.Name ^& " (Type=" ^& vbcomp.Type ^& ")" >> temp_verify.vbs
echo     count = count + 1 >> temp_verify.vbs
echo Next >> temp_verify.vbs
echo. >> temp_verify.vbs
echo WScript.Echo "=========================" >> temp_verify.vbs
echo WScript.Echo "Total: " ^& count ^& " components" >> temp_verify.vbs
echo. >> temp_verify.vbs
echo wb.Close False >> temp_verify.vbs
echo xl.Quit >> temp_verify.vbs
echo Set wb = Nothing >> temp_verify.vbs
echo Set xl = Nothing >> temp_verify.vbs

cscript //nologo temp_verify.vbs

del temp_verify.vbs

echo.
echo ========================================
echo Expected modules:
echo   - ThisWorkbook (Type=100)
echo   - Sheet1 (Type=100)
echo   - PQ_Setup (Type=1)
echo   - Python_Run (Type=1)
echo   - Setup (Type=1)
echo   - TF_Data (Type=1)
echo   - TF_Presets (Type=1)
echo   - TF_UI (Type=1)
echo   - TF_Utils (Type=1)
echo ========================================
echo.
pause
