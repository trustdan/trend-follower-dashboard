@echo off
echo ========================================
echo Cleanup: Kill stuck Excel processes
echo ========================================

echo.
echo Attempting to close all Excel instances...
taskkill /F /IM excel.exe /T
timeout /t 2 /nobreak > nul

echo.
echo Checking if Excel is still running...
tasklist /FI "IMAGENAME eq excel.exe" 2>NUL | find /I /N "excel.exe" > nul
if "%ERRORLEVEL%"=="0" (
    echo WARNING: Excel is still running!
    echo Please close Excel manually and run this script again.
    pause
    exit /b 1
) else (
    echo OK: No Excel processes found.
)

echo.
echo Checking for locked workbook...
if exist "TrendFollowing_TradeEntry.xlsm" (
    echo Found: TrendFollowing_TradeEntry.xlsm
    echo Attempting to delete...
    del /F "TrendFollowing_TradeEntry.xlsm" 2>nul
    if exist "TrendFollowing_TradeEntry.xlsm" (
        echo ERROR: Could not delete workbook (may be locked)
        echo Please close Excel and delete manually.
        pause
        exit /b 1
    ) else (
        echo OK: Workbook deleted.
    )
) else (
    echo OK: No workbook found.
)

echo.
echo ========================================
echo Cleanup complete!
echo ========================================
echo You can now run IMPORT_VBA_MODULES_DEBUG.bat
pause
