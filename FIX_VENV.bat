@echo off
echo ========================================
echo Fix Broken Python Venv
echo ========================================

set "VENV=venv"

echo.
echo This will delete and recreate the Python virtual environment.
echo.
pause

if exist "%VENV%" (
    echo Deleting old venv folder...
    rd /s /q "%VENV%"
    if exist "%VENV%" (
        echo ERROR: Could not delete venv folder
        echo Please close any programs using files in the venv folder
        pause
        exit /b 1
    )
    echo Old venv deleted.
)

echo.
echo Creating fresh venv...
py -3 -m venv "%VENV%"
if errorlevel 1 (
    echo ERROR: Failed to create venv
    echo.
    echo Make sure Python 3 is installed.
    echo Try running: py -3 --version
    pause
    exit /b 1
)

echo.
echo Upgrading pip...
%VENV%\Scripts\python.exe -m pip install --upgrade pip
if errorlevel 1 (
    echo WARNING: pip upgrade failed, but continuing...
)

echo.
echo Installing pywin32...
%VENV%\Scripts\pip.exe install pywin32
if errorlevel 1 (
    echo ERROR: Failed to install pywin32
    pause
    exit /b 1
)

echo.
echo ========================================
echo SUCCESS!
echo ========================================
echo.
echo Python venv fixed and ready.
echo pywin32 installed.
echo.
echo You can now run: BUILD_WITH_PYTHON.bat
echo.
pause
