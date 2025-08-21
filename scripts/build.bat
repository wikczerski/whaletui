@echo off
REM Build script for whaletui with automatic version injection

REM Get version information from git
for /f "tokens=*" %%i in ('git describe --tags --always --dirty 2^>nul') do set VERSION=%%i
if "%VERSION%"=="" set VERSION=dev

for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul') do set COMMIT_SHA=%%i
if "%COMMIT_SHA%"=="" set COMMIT_SHA=unknown

for /f "tokens=*" %%i in ('powershell -Command "Get-Date -Format 'yyyy-MM-dd_HH:mm:ss_UTC'"') do set BUILD_DATE=%%i

echo Building whaletui...
echo Version: %VERSION%
echo Commit: %COMMIT_SHA%
echo Build Date: %BUILD_DATE%

REM Build with version information
go build -ldflags "-X github.com/wikczerski/whaletui/cmd.Version=%VERSION% -X github.com/wikczerski/whaletui/cmd.CommitSHA=%COMMIT_SHA% -X github.com/wikczerski/whaletui/cmd.BuildDate=%BUILD_DATE%" -o whaletui.exe .

echo Build complete: whaletui.exe
