@echo off
setlocal enabledelayedexpansion
REM --- repo root from scripts/ folder ---
pushd "%~dp0\.."
set "REPO=%CD%"
set "VENV=%REPO%\venv"

if not exist "%VENV%\Scripts\python.exe" (
  py -3 -m venv "%VENV%" || (echo [ERR] Failed to create venv & pause & exit /b 1)
)

REM --- pick requirements path: arg > root > python\requirements.txt ---
set "REQ=%~1"
if not defined REQ set "REQ=%REPO%\requirements.txt"
if not exist "%REQ%" if exist "%REPO%\python\requirements.txt" set "REQ=%REPO%\python\requirements.txt"

REM --- auto-create a default if still missing ---
if not exist "%REQ%" (
  echo [INFO] %REQ% not found; creating a minimal one.
  > "%REPO%\requirements.txt" (
    echo pywin32>=305
    echo requests>=2.31
    echo beautifulsoup4>=4.12
    echo lxml>=4.9
  )
  set "REQ=%REPO%\requirements.txt"
)

echo Using requirements: %REQ%
"%VENV%\Scripts\python.exe" -m pip install --upgrade pip || (echo [ERR] pip upgrade failed & pause & exit /b 1)
"%VENV%\Scripts\pip.exe" install -r "%REQ%" || (echo [ERR] pip install failed & pause & exit /b 1)

echo [OK] Virtual env ready at %VENV%
popd
exit /b 0
