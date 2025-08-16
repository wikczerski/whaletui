@echo off
REM Build script for d5r (Windows Batch version)
REM This script helps with local testing of the build process on Windows

setlocal enabledelayedexpansion

REM Configuration
set BINARY_NAME=d5r
set BUILD_DIR=dist

REM Get version from git or use default
for /f "tokens=*" %%i in ('git describe --tags --always --dirty 2^>nul') do set VERSION=%%i
if "%VERSION%"=="" set VERSION=dev

echo Building %BINARY_NAME% version %VERSION%

REM Clean previous builds
echo Cleaning previous builds...
if exist "%BUILD_DIR%" rmdir /s /q "%BUILD_DIR%"
if exist "bin" rmdir /s /q "bin"

REM Create build directory
mkdir "%BUILD_DIR%"

REM Build variables
set LDFLAGS=-ldflags=-s -w -X main.Version=%VERSION%

echo Building for all platforms...

REM Linux builds
echo Building Linux binaries...
set GOOS=linux
set GOARCH=amd64
go build %LDFLAGS% -o "%BUILD_DIR%\%BINARY_NAME%-linux-amd64" main.go

set GOOS=linux
set GOARCH=arm64
go build %LDFLAGS% -o "%BUILD_DIR%\%BINARY_NAME%-linux-arm64" main.go

set GOOS=linux
set GOARCH=arm
set GOARM=7
go build %LDFLAGS% -o "%BUILD_DIR%\%BINARY_NAME%-linux-armv7" main.go

set GOOS=linux
set GOARCH=ppc64le
go build %LDFLAGS% -o "%BUILD_DIR%\%BINARY_NAME%-linux-ppc64le" main.go

set GOOS=linux
set GOARCH=s390x
go build %LDFLAGS% -o "%BUILD_DIR%\%BINARY_NAME%-linux-s390x" main.go

REM Darwin builds
echo Building Darwin binaries...
set GOOS=darwin
set GOARCH=amd64
go build %LDFLAGS% -o "%BUILD_DIR%\%BINARY_NAME%-darwin-amd64" main.go

set GOOS=darwin
set GOARCH=arm64
go build %LDFLAGS% -o "%BUILD_DIR%\%BINARY_NAME%-darwin-arm64" main.go

REM FreeBSD builds
echo Building FreeBSD binaries...
set GOOS=freebsd
set GOARCH=amd64
go build %LDFLAGS% -o "%BUILD_DIR%\%BINARY_NAME%-freebsd-amd64" main.go

set GOOS=freebsd
set GOARCH=arm64
go build %LDFLAGS% -o "%BUILD_DIR%\%BINARY_NAME%-freebsd-arm64" main.go

REM Windows builds
echo Building Windows binaries...
set GOOS=windows
set GOARCH=amd64
go build %LDFLAGS% -o "%BUILD_DIR%\%BINARY_NAME%-windows-amd64.exe" main.go

set GOOS=windows
set GOARCH=arm64
go build %LDFLAGS% -o "%BUILD_DIR%\%BINARY_NAME%-windows-arm64.exe" main.go

REM Clear environment variables
set GOOS=
set GOARCH=
set GOARM=

echo Build complete!
echo Binaries are in %BUILD_DIR%\

REM Create packages directory
mkdir "%BUILD_DIR%\packages"

echo Creating packages...

REM Linux packages
echo Creating Linux packages...
cd "%BUILD_DIR%"
tar -czf "packages\%BINARY_NAME%-linux-amd64.tar.gz" "%BINARY_NAME%-linux-amd64"
tar -czf "packages\%BINARY_NAME%-linux-arm64.tar.gz" "%BINARY_NAME%-linux-arm64"
tar -czf "packages\%BINARY_NAME%-linux-armv7.tar.gz" "%BINARY_NAME%-linux-armv7"
tar -czf "packages\%BINARY_NAME%-linux-ppc64le.tar.gz" "%BINARY_NAME%-linux-ppc64le"
tar -czf "packages\%BINARY_NAME%-linux-s390x.tar.gz" "%BINARY_NAME%-linux-s390x"

REM Darwin packages
echo Creating Darwin packages...
tar -czf "packages\%BINARY_NAME%-darwin-amd64.tar.gz" "%BINARY_NAME%-darwin-amd64"
tar -czf "packages\%BINARY_NAME%-darwin-arm64.tar.gz" "%BINARY_NAME%-darwin-arm64"

REM FreeBSD packages
echo Creating FreeBSD packages...
tar -czf "packages\%BINARY_NAME%-freebsd-amd64.tar.gz" "%BINARY_NAME%-freebsd-amd64"
tar -czf "packages\%BINARY_NAME%-freebsd-arm64.tar.gz" "%BINARY_NAME%-freebsd-arm64"

REM Windows packages
echo Creating Windows packages...
powershell -Command "Compress-Archive -Path '%BINARY_NAME%-windows-amd64.exe' -DestinationPath 'packages\%BINARY_NAME%-windows-amd64.zip' -Force"
powershell -Command "Compress-Archive -Path '%BINARY_NAME%-windows-arm64.exe' -DestinationPath 'packages\%BINARY_NAME%-windows-arm64.zip' -Force"

cd ..

echo Packages created!

REM Generate checksums
echo Generating checksums...
cd "%BUILD_DIR%\packages"

REM Use PowerShell to generate checksums since batch doesn't have built-in hash functions
powershell -Command "Get-ChildItem -File | ForEach-Object { $hash = Get-FileHash -Algorithm SHA256 -Path $_.FullName; Write-Output \"$($hash.Hash.ToLower())  $($_.Name)\" } | Out-File -FilePath \"../checksums.sha256\" -Encoding UTF8"

cd ..\..

echo Checksums generated: %BUILD_DIR%\checksums.sha256

REM Show results
echo Build Summary:
dir "%BUILD_DIR%" /b | find /c /v "" > temp_count.txt
set /p BINARY_COUNT=<temp_count.txt
set /a BINARY_COUNT-=2
echo Binaries: %BINARY_COUNT%
del temp_count.txt

dir "%BUILD_DIR%\packages" /b | find /c /v "" > temp_count.txt
set /p PACKAGE_COUNT=<temp_count.txt
echo Packages: %PACKAGE_COUNT%
del temp_count.txt

echo.
echo Files created:
dir "%BUILD_DIR%"

echo.
echo Packages created:
dir "%BUILD_DIR%\packages"

echo.
echo Checksums:
type "%BUILD_DIR%\checksums.sha256"

echo.
echo Build process completed successfully!
pause
