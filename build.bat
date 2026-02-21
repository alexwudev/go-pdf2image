@echo off
cd /d "%~dp0"
echo === Building frontend ===
cd frontend
call npm install
call npm run build
cd ..
echo === Building pdf2image.exe ===
go build -tags desktop,production -ldflags "-H windowsgui" -o pdf2image.exe .
if %errorlevel% equ 0 (
    echo Build OK!
) else (
    echo Build FAILED!
)
pause
