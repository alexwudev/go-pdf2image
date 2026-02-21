@echo off
cd /d "%~dp0\.."
echo === Building frontend ===
cd frontend
call npm install
call npm run build
cd ..

echo === Embedding icon resource ===
where go-winres >nul 2>nul
if %errorlevel% equ 0 (
    go-winres make --in platform\windows\winres.json --product-version 1.0.0.0 --file-version 1.0.0.0
) else (
    echo go-winres not found, skipping icon embed
)

echo === Building pdf2image.exe ===
go build -tags desktop,production -ldflags "-s -w -H windowsgui" -o pdf2image.exe .
if %errorlevel% equ 0 (
    del /q rsrc_windows_*.syso >nul 2>nul
    echo Build OK! -^> pdf2image.exe
) else (
    echo Build FAILED!
)
pause
