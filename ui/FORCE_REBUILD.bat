@echo off
echo Forcing clean rebuild of TF-Engine GUI...
echo.

REM Kill any running instances
taskkill /F /IM tf-gui.exe 2>nul
timeout /t 2 /nobreak >nul

REM Delete old executable
echo Deleting old executable...
del tf-gui.exe 2>nul

REM Clean build cache
echo Cleaning build cache...
go clean -cache
go clean -modcache

REM Rebuild
echo Building new executable...
go build -o tf-gui-new.exe .

if %ERRORLEVEL% EQU 0 (
    echo.
    echo Build successful!
    echo Renaming tf-gui-new.exe to tf-gui.exe...
    timeout /t 1 /nobreak >nul
    move /Y tf-gui-new.exe tf-gui.exe
    echo.
    echo Done! You can now run tf-gui.exe
) else (
    echo.
    echo Build failed! Check errors above.
)

pause
