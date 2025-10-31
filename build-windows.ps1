# TF-Engine Windows Build Script
# Rebuilds all components: migration tool, backend, and GUI

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  TF-Engine Windows Build Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Track overall success
$buildSuccess = $true
$migrateSuccess = $true

# Step 1: Build migration tool
Write-Host "[1/3] Building migration tool..." -ForegroundColor Yellow
Set-Location backend
go build -o ..\migrate-db.exe .\cmd\migrate\main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Migration tool build failed!" -ForegroundColor Red
    $buildSuccess = $false
} else {
    Write-Host "✅ migrate-db.exe built successfully" -ForegroundColor Green
}
Set-Location ..

# Step 2: Run migration (if database exists)
if ($buildSuccess -and (Test-Path "trading.db")) {
    Write-Host ""
    Write-Host "[2/3] Running database migration..." -ForegroundColor Yellow
    .\migrate-db.exe
    if ($LASTEXITCODE -ne 0) {
        Write-Host "⚠️  Migration had warnings (this is OK if columns already exist)" -ForegroundColor Yellow
        $migrateSuccess = $false
    } else {
        Write-Host "✅ Database migration completed" -ForegroundColor Green
    }
} elseif (-not (Test-Path "trading.db")) {
    Write-Host ""
    Write-Host "[2/3] Skipping migration (trading.db not found)" -ForegroundColor Yellow
    Write-Host "    Note: Database will be created when you first run tf-gui.exe" -ForegroundColor Gray
} else {
    Write-Host ""
    Write-Host "[2/3] Skipping migration (migration tool build failed)" -ForegroundColor Red
}

# Step 3: Build GUI
Write-Host ""
Write-Host "[3/3] Building GUI application..." -ForegroundColor Yellow
Set-Location ui
go build -o ..\tf-gui.exe
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ GUI build failed!" -ForegroundColor Red
    $buildSuccess = $false
} else {
    Write-Host "✅ tf-gui.exe built successfully" -ForegroundColor Green
}
Set-Location ..

# Step 4: Summary and notification
Write-Host ""
Write-Host "[4/4] Build Summary" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Cyan

if ($buildSuccess) {
    Write-Host "✅ All builds successful!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Ready to run:" -ForegroundColor White
    Write-Host "  .\tf-gui.exe          - Launch GUI" -ForegroundColor Cyan
    Write-Host "  .\migrate-db.exe      - Run migration again if needed" -ForegroundColor Cyan

    # Send success notification
    $notificationTitle = "TF-Engine Build Complete!"
    if ($migrateSuccess) {
        $notificationMessage = "GUI rebuilt + Database migrated - Ready to trade!"
    } else {
        $notificationMessage = "GUI rebuilt successfully - Ready to test!"
    }

} else {
    Write-Host "❌ Build had errors - check output above" -ForegroundColor Red
    $notificationTitle = "TF-Engine Build Failed"
    $notificationMessage = "Check PowerShell output for errors"
}

Write-Host "========================================" -ForegroundColor Cyan

# Send toast notification
Write-Host ""
Write-Host "Sending notification..." -ForegroundColor Gray
& /mnt/c/Windows/System32/WindowsPowerShell/v1.0/powershell.exe -Command "
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

`$template = @'
<toast>
    <visual>
        <binding template=`"ToastText02`">
            <text id=`"1`">$notificationTitle</text>
            <text id=`"2`">$notificationMessage</text>
        </binding>
    </visual>
</toast>
'@

`$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
`$xml.LoadXml(`$template)
`$toast = [Windows.UI.Notifications.ToastNotification]::new(`$xml)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('TF-Engine').Show(`$toast)
"

Write-Host ""
if ($buildSuccess) {
    exit 0
} else {
    exit 1
}
