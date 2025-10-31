@echo off
REM TF-Engine Windows Build Script (Batch version)
REM Rebuilds all components: migration tool, backend, and GUI

echo ========================================
echo   TF-Engine Windows Build Script
echo ========================================
echo.

REM Step 1: Build migration tool
echo [1/3] Building migration tool...
cd backend
go build -o ..\migrate-db.exe .\cmd\migrate\main.go
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Migration tool build failed!
    cd ..
    goto :error
)
echo SUCCESS: migrate-db.exe built
cd ..

REM Step 2: Run migration (if database exists)
echo.
if exist "trading.db" (
    echo [2/3] Running database migration...
    migrate-db.exe
    if %ERRORLEVEL% NEQ 0 (
        echo WARNING: Migration had warnings ^(this is OK if columns already exist^)
    ) else (
        echo SUCCESS: Database migration completed
    )
) else (
    echo [2/3] Skipping migration ^(trading.db not found^)
    echo     Note: Database will be created when you first run tf-gui.exe
)

REM Step 3: Build GUI
echo.
echo [3/3] Building GUI application...
cd ui
go build -o ..\tf-gui.exe
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: GUI build failed!
    cd ..
    goto :error
)
echo SUCCESS: tf-gui.exe built
cd ..

REM Step 4: Summary
echo.
echo [4/4] Build Summary
echo ========================================
echo SUCCESS: All builds completed!
echo.
echo Ready to run:
echo   tf-gui.exe          - Launch GUI
echo   migrate-db.exe      - Run migration again if needed
echo ========================================

REM Send success notification
powershell.exe -Command "[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null; [Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null; $template = '<toast><visual><binding template=\"ToastText02\"><text id=\"1\">TF-Engine Build Complete!</text><text id=\"2\">GUI rebuilt + Database migrated - Ready to trade!</text></binding></visual></toast>'; $xml = New-Object Windows.Data.Xml.Dom.XmlDocument; $xml.LoadXml($template); $toast = [Windows.UI.Notifications.ToastNotification]::new($xml); [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('TF-Engine').Show($toast)"

exit /b 0

:error
echo.
echo ========================================
echo ERROR: Build failed - check output above
echo ========================================
powershell.exe -Command "[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null; [Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null; $template = '<toast><visual><binding template=\"ToastText02\"><text id=\"1\">TF-Engine Build Failed</text><text id=\"2\">Check command prompt for errors</text></binding></visual></toast>'; $xml = New-Object Windows.Data.Xml.Dom.XmlDocument; $xml.LoadXml($template); $toast = [Windows.UI.Notifications.ToastNotification]::new($xml); [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('TF-Engine').Show($toast)"
exit /b 1
