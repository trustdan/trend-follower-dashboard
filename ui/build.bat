@echo off
REM Build script for TF-Engine GUI
REM Usage: build.bat [version]
REM Example: build.bat v10
REM Default: build.bat (builds tf-gui.exe)

SET VERSION=%1
IF "%VERSION%"=="" SET VERSION=

SET FILES=main.go dashboard.go checklist.go position_sizing.go heat_check.go trade_entry.go scanner.go calendar.go theme.go widgets.go utils.go keybindings.go keybindings_v2.go

IF "%VERSION%"=="" (
    echo Building tf-gui.exe...
    go build -o tf-gui.exe %FILES%
) ELSE (
    echo Building tf-gui-%VERSION%.exe...
    go build -o tf-gui-%VERSION%.exe %FILES%
)

IF %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo BUILD SUCCESSFUL!
    echo ========================================
    IF "%VERSION%"=="" (
        echo Binary: tf-gui.exe
    ) ELSE (
        echo Binary: tf-gui-%VERSION%.exe
    )
    echo Size:
    IF "%VERSION%"=="" (
        dir tf-gui.exe | find "tf-gui.exe"
    ) ELSE (
        dir tf-gui-%VERSION%.exe | find "tf-gui-%VERSION%.exe"
    )
    echo ========================================
) ELSE (
    echo.
    echo ========================================
    echo BUILD FAILED!
    echo ========================================
    echo Check the error messages above
    echo ========================================
    exit /b 1
)
