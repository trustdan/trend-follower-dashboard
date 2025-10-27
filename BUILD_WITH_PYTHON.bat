@echo off
setlocal

echo ========================================
echo Build Workbook Using Python
echo ========================================

set "WB=TrendFollowing_TradeEntry.xlsm"
set "VENV=venv"
set "PYTHON=%VENV%\Scripts\python.exe"

REM Check if venv exists, create if missing
if not exist "%PYTHON%" (
    echo Python venv not found - creating...
    call scripts\setup_venv.bat
    if errorlevel 1 (
        echo ERROR: Failed to create Python venv
        echo.
        echo Make sure Python 3 is installed and available as 'py -3'
        pause
        exit /b 1
    )
    echo Venv created successfully
)

REM Check if pip is working in venv
%PYTHON% -m pip --version >nul 2>&1
if errorlevel 1 (
    echo Venv exists but pip is broken - recreating...
    if exist "%VENV%" (
        echo Deleting old venv...
        rd /s /q "%VENV%"
    )
    echo Creating fresh venv...
    call scripts\setup_venv.bat
    if errorlevel 1 (
        echo ERROR: Failed to recreate Python venv
        pause
        exit /b 1
    )
    echo Fresh venv created
)

echo.
echo Activating venv...
call %VENV%\Scripts\activate.bat

echo.
echo Checking dependencies...
%PYTHON% -c "import win32com.client" 2>nul
if errorlevel 1 (
    echo pywin32 not found - installing...
    %PYTHON% -m pip install --quiet pywin32
    if errorlevel 1 (
        echo ERROR: Failed to install pywin32
        pause
        exit /b 1
    )
    echo pywin32 installed successfully
)

echo.
echo Current directory: %CD%
echo Target workbook: %WB%

REM Trust Center toggles (best-effort)
echo.
echo Configuring Excel Trust Center...
reg add "HKCU\Software\Microsoft\Office\16.0\Excel\Security" /v AccessVBOM /t REG_DWORD /d 1 /f >nul 2>&1
reg add "HKCU\Software\Microsoft\Office\16.0\Excel\Security" /v VBAWarnings /t REG_DWORD /d 1 /f >nul 2>&1

REM Kill any stray Excel
echo Closing any existing Excel instances...
taskkill /IM excel.exe /F >nul 2>&1

REM Delete existing workbook if present
if exist "%WB%" (
    echo Deleting existing workbook...
    del /F "%WB%" 2>nul
    if exist "%WB%" (
        echo WARNING: Could not delete %WB% - it may be locked
        pause
    )
)

echo.
echo ========================================
echo Running Python import script...
echo ========================================
echo.

python import_to_excel.py "%CD%\%WB%"
set "RC=%ERRORLEVEL%"

if "%RC%"=="0" (
    echo.
    echo ========================================
    echo SUCCESS!
    echo ========================================
    echo.
    echo Workbook created: %WB%
    echo Excel is still running - check the workbook
    echo.
    echo Press Alt+F11 in Excel to view VBA modules
    echo.
) else (
    echo.
    echo ========================================
    echo FAILED - Exit code: %RC%
    echo ========================================
    echo.
)

pause
exit /b %RC%
