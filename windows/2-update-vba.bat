@echo off
REM ===========================================================================
REM update-vba.bat - Quick VBA Module Update Script
REM ===========================================================================
REM Purpose: Re-import VBA modules without recreating workbook/database
REM
REM Usage: Run this after VBA .bas files have been updated
REM ===========================================================================

echo ========================================
echo  Update VBA Modules
echo ========================================
echo.

REM Check if workbook exists
if not exist "TradingPlatform.xlsm" (
    echo ERROR: TradingPlatform.xlsm not found
    echo Please run setup-all.bat first
    pause
    exit /b 1
)

echo Backing up workbook...
copy TradingPlatform.xlsm TradingPlatform.xlsm.backup >nul
echo [OK] Backup created

echo.
echo Importing updated VBA modules...
echo.

REM Create import script
echo ' VBA Module Update Script > update_vba.vbs
echo Set objExcel = CreateObject("Excel.Application") >> update_vba.vbs
echo objExcel.Visible = False >> update_vba.vbs
echo objExcel.DisplayAlerts = False >> update_vba.vbs
echo. >> update_vba.vbs
echo ' Get script directory >> update_vba.vbs
echo Set objFSO = CreateObject("Scripting.FileSystemObject") >> update_vba.vbs
echo strScriptDir = objFSO.GetParentFolderName(WScript.ScriptFullName) >> update_vba.vbs
echo. >> update_vba.vbs
echo ' Open workbook >> update_vba.vbs
echo strWorkbookPath = strScriptDir + "\TradingPlatform.xlsm" >> update_vba.vbs
echo Set objWorkbook = objExcel.Workbooks.Open(strWorkbookPath) >> update_vba.vbs
echo. >> update_vba.vbs
echo ' Remove existing modules >> update_vba.vbs
echo On Error Resume Next >> update_vba.vbs
echo For Each comp In objWorkbook.VBProject.VBComponents >> update_vba.vbs
echo     If comp.Type = 1 Then ' vbext_ct_StdModule >> update_vba.vbs
echo         Select Case comp.Name >> update_vba.vbs
echo             Case "TFTypes", "TFHelpers", "TFEngine", "TFTests" >> update_vba.vbs
echo                 objWorkbook.VBProject.VBComponents.Remove comp >> update_vba.vbs
echo         End Select >> update_vba.vbs
echo     End If >> update_vba.vbs
echo Next >> update_vba.vbs
echo On Error Goto 0 >> update_vba.vbs
echo. >> update_vba.vbs
echo ' Import updated modules >> update_vba.vbs
echo strVBADir = objFSO.GetParentFolderName(strScriptDir) + "\excel\vba\" >> update_vba.vbs
echo. >> update_vba.vbs
echo WScript.Echo "Importing TFTypes.bas..." >> update_vba.vbs
echo objWorkbook.VBProject.VBComponents.Import strVBADir + "TFTypes.bas" >> update_vba.vbs
echo. >> update_vba.vbs
echo WScript.Echo "Importing TFHelpers.bas..." >> update_vba.vbs
echo objWorkbook.VBProject.VBComponents.Import strVBADir + "TFHelpers.bas" >> update_vba.vbs
echo. >> update_vba.vbs
echo WScript.Echo "Importing TFEngine.bas..." >> update_vba.vbs
echo objWorkbook.VBProject.VBComponents.Import strVBADir + "TFEngine.bas" >> update_vba.vbs
echo. >> update_vba.vbs
echo WScript.Echo "Importing TFTests.bas..." >> update_vba.vbs
echo objWorkbook.VBProject.VBComponents.Import strVBADir + "TFTests.bas" >> update_vba.vbs
echo. >> update_vba.vbs
echo ' Save and close >> update_vba.vbs
echo objWorkbook.Save >> update_vba.vbs
echo objWorkbook.Close >> update_vba.vbs
echo objExcel.Quit >> update_vba.vbs
echo. >> update_vba.vbs
echo ' Cleanup >> update_vba.vbs
echo Set objWorkbook = Nothing >> update_vba.vbs
echo Set objExcel = Nothing >> update_vba.vbs
echo Set objFSO = Nothing >> update_vba.vbs
echo. >> update_vba.vbs
echo WScript.Echo "VBA modules updated successfully!" >> update_vba.vbs

REM Run import
cscript //nologo update_vba.vbs

if %errorlevel% neq 0 (
    echo ERROR: VBA update failed
    echo Restoring backup...
    copy TradingPlatform.xlsm.backup TradingPlatform.xlsm >nul
    pause
    exit /b 1
)

del update_vba.vbs
echo.
echo ========================================
echo  Update Complete!
echo ========================================
echo.
echo VBA modules have been updated
echo Backup saved as: TradingPlatform.xlsm.backup
echo.

pause
