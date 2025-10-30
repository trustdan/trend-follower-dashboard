@echo off
echo Starting TF-Engine GUI with logging...
echo.
echo Log file will be created at: tf-gui.log
echo.

REM Run the application and capture all output to a log file
tf-gui.exe > tf-gui.log 2>&1

REM If the app crashed, show the log
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo Application crashed! Error code: %ERRORLEVEL%
    echo.
    echo Last 50 lines of log:
    echo ================================
    type tf-gui.log | more
    echo ================================
    echo.
    echo Full log saved to tf-gui.log
    pause
) else (
    echo Application exited normally.
)
