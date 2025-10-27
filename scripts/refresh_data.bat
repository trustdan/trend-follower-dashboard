@echo off
setlocal
pushd "%~dp0\.."
set REPO=%CD%
set VENV=%REPO%\venv
set PY="%VENV%\Scripts\python.exe"
set OUTDIR=%REPO%\data
if not exist "%OUTDIR%" mkdir "%OUTDIR%"

%PY% "%REPO%\python\finviz_scraper.py" ^
  --preset-file "%REPO%\config\presets.json" ^
  --out "%OUTDIR%\screened.csv" || exit /b 1

echo [OK] Data refreshed at %OUTDIR%\screened.csv
exit /b 0
