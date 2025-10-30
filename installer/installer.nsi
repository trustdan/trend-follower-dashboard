; TF-Engine Installer Script
; Trend Following Trading Discipline System
; Version 1.0.0

!define APP_NAME "TF-Engine"
!define COMP_NAME "TF-Engine Project"
!define VERSION "1.0.0.0"
!define COPYRIGHT "Copyright © 2025 TF-Engine Project"
!define DESCRIPTION "Trend-Following Trading Discipline System"
!define INSTALLER_NAME "TF-Engine-Setup-v1.0.0.exe"
!define MAIN_APP_EXE "tf-engine.exe"
!define INSTALL_DIR "$PROGRAMFILES64\${APP_NAME}"

; Includes
!include "MUI2.nsh"
!include "FileFunc.nsh"
!include "LogicLib.nsh"

; MUI Settings
!define MUI_ABORTWARNING
!define MUI_ICON "..\backend\assets\trend_following_icon_proper.ico"
!define MUI_UNICON "..\backend\assets\trend_following_icon_proper.ico"

; Welcome page
!insertmacro MUI_PAGE_WELCOME

; Directory page
!insertmacro MUI_PAGE_DIRECTORY

; Instfiles page
!insertmacro MUI_PAGE_INSTFILES

; Finish page
!define MUI_FINISHPAGE_TEXT "TF-Engine has been installed successfully.$\r$\n$\r$\nClick Finish to close this wizard."
!define MUI_FINISHPAGE_RUN "$INSTDIR\${MAIN_APP_EXE}"
!define MUI_FINISHPAGE_RUN_PARAMETERS "server"
!define MUI_FINISHPAGE_RUN_TEXT "Launch TF-Engine"
!insertmacro MUI_PAGE_FINISH

; Uninstaller pages
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

; Language
!insertmacro MUI_LANGUAGE "English"

; Installer attributes
Name "${APP_NAME}"
OutFile "${INSTALLER_NAME}"
InstallDir "${INSTALL_DIR}"
InstallDirRegKey HKLM "Software\${APP_NAME}" "InstallDir"
RequestExecutionLevel admin
ShowInstDetails show
ShowUnInstDetails show

; Version Info
VIProductVersion "${VERSION}"
VIAddVersionKey "ProductName" "${APP_NAME}"
VIAddVersionKey "CompanyName" "${COMP_NAME}"
VIAddVersionKey "FileDescription" "${DESCRIPTION}"
VIAddVersionKey "FileVersion" "${VERSION}"
VIAddVersionKey "LegalCopyright" "${COPYRIGHT}"

; Installer Sections
Section "MainSection" SEC01
    ; Set output path to install directory
    SetOutPath "$INSTDIR"

    ; Copy main executable
    File "..\backend\tf-engine.exe"

    ; Create AppData directory for database
    CreateDirectory "$APPDATA\${APP_NAME}"

    ; Initialize database (silent, no window)
    DetailPrint "Initializing database..."
    nsExec::ExecToLog '"$INSTDIR\${MAIN_APP_EXE}" init --db "$APPDATA\${APP_NAME}\trading.db"'
    Pop $0  ; Return value
    ${If} $0 != 0
        MessageBox MB_OK|MB_ICONEXCLAMATION "Database initialization failed. You may need to run 'tf-engine.exe init' manually."
    ${EndIf}

    ; Create Start Menu folder
    CreateDirectory "$SMPROGRAMS\${APP_NAME}"

    ; Create Start Menu shortcut
    CreateShortCut "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk" \
        "$INSTDIR\${MAIN_APP_EXE}" \
        'server --db "$APPDATA\${APP_NAME}\trading.db"' \
        "$INSTDIR\${MAIN_APP_EXE}" 0 \
        SW_SHOWNORMAL "" \
        "Launch TF-Engine Trading System"

    ; Create Desktop shortcut
    CreateShortCut "$DESKTOP\${APP_NAME}.lnk" \
        "$INSTDIR\${MAIN_APP_EXE}" \
        'server --db "$APPDATA\${APP_NAME}\trading.db"' \
        "$INSTDIR\${MAIN_APP_EXE}" 0 \
        SW_SHOWNORMAL "" \
        "Launch TF-Engine Trading System"

    ; Write installation path to registry
    WriteRegStr HKLM "Software\${APP_NAME}" "InstallDir" "$INSTDIR"
    WriteRegStr HKLM "Software\${APP_NAME}" "Version" "${VERSION}"

    ; Write uninstall information
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "DisplayName" "${APP_NAME}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "DisplayVersion" "${VERSION}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "Publisher" "${COMP_NAME}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "UninstallString" "$INSTDIR\uninstall.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "DisplayIcon" "$INSTDIR\${MAIN_APP_EXE}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "InstallLocation" "$INSTDIR"
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "NoModify" 1
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "NoRepair" 1

    ; Calculate installed size
    ${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
    IntFmt $0 "0x%08X" $0
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" \
        "EstimatedSize" "$0"

    ; Create uninstaller
    WriteUninstaller "$INSTDIR\uninstall.exe"
SectionEnd

; Uninstaller Section
Section "Uninstall"
    ; Remove application files
    Delete "$INSTDIR\tf-engine.exe"
    Delete "$INSTDIR\uninstall.exe"

    ; Remove installation directory (if empty)
    RMDir "$INSTDIR"

    ; Remove shortcuts
    Delete "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk"
    RMDir "$SMPROGRAMS\${APP_NAME}"
    Delete "$DESKTOP\${APP_NAME}.lnk"

    ; Ask user if they want to delete database (user data)
    MessageBox MB_YESNO|MB_ICONQUESTION \
        "Do you want to delete your trading database? This will remove all your positions, decisions, and settings. This action cannot be undone." \
        IDNO skip_database

    ; Delete database if user confirmed
    Delete "$APPDATA\${APP_NAME}\trading.db"
    Delete "$APPDATA\${APP_NAME}\trading.db-shm"
    Delete "$APPDATA\${APP_NAME}\trading.db-wal"
    RMDir "$APPDATA\${APP_NAME}"

    skip_database:

    ; Remove registry keys
    DeleteRegKey HKLM "Software\${APP_NAME}"
    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}"
SectionEnd
