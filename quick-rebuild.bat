@echo off
REM Quick rebuild script - GUI only (no migration)
REM Use this for fast iteration during development

echo Building tf-gui.exe...
cd ui
go build -o ..\tf-gui.exe
if %ERRORLEVEL% NEQ 0 (
    echo Build failed!
    cd ..
    exit /b 1
)
cd ..
echo.
echo SUCCESS! tf-gui.exe rebuilt
echo Run: .\tf-gui.exe

REM Send notification
powershell.exe -Command "[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null; [Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null; $template = '<toast><visual><binding template=\"ToastText02\"><text id=\"1\">TF-Engine GUI Rebuilt!</text><text id=\"2\">Ready to test - run tf-gui.exe</text></binding></visual></toast>'; $xml = New-Object Windows.Data.Xml.Dom.XmlDocument; $xml.LoadXml($template); $toast = [Windows.UI.Notifications.ToastNotification]::new($xml); [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('TF-Engine').Show($toast)"
