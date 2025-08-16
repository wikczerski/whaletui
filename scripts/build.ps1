# Build script for d5r (PowerShell version)
# This script helps with local testing of the build process on Windows

param(
    [string]$Version = ""
)

# Configuration
$BinaryName = "d5r"
$BuildDir = "dist"

if ([string]::IsNullOrEmpty($Version)) {
    try {
        $Version = git describe --tags --always --dirty 2>$null
        if ([string]::IsNullOrEmpty($Version)) {
            $Version = "dev"
        }
    }
    catch {
        $Version = "dev"
    }
}

Write-Host "Building $BinaryName version $Version" -ForegroundColor Green

# Clean previous builds
Write-Host "Cleaning previous builds..." -ForegroundColor Yellow
if (Test-Path $BuildDir) { Remove-Item -Recurse -Force $BuildDir }
if (Test-Path "bin") { Remove-Item -Recurse -Force "bin" }

# Create build directory
New-Item -ItemType Directory -Path $BuildDir -Force | Out-Null

# Build variables
$LDFLAGS = "-ldflags=`"-s -w -X main.Version=$Version`""

Write-Host "Building for all platforms..." -ForegroundColor Yellow

# Linux builds
Write-Host "Building Linux binaries..." -ForegroundColor Cyan
$env:GOOS = "linux"; $env:GOARCH = "amd64"; go build $LDFLAGS -o "$BuildDir/$BinaryName-linux-amd64" main.go
$env:GOOS = "linux"; $env:GOARCH = "arm64"; go build $LDFLAGS -o "$BuildDir/$BinaryName-linux-arm64" main.go
$env:GOOS = "linux"; $env:GOARCH = "arm"; $env:GOARM = "7"; go build $LDFLAGS -o "$BuildDir/$BinaryName-linux-armv7" main.go
$env:GOOS = "linux"; $env:GOARCH = "ppc64le"; go build $LDFLAGS -o "$BuildDir/$BinaryName-linux-ppc64le" main.go
$env:GOOS = "linux"; $env:GOARCH = "s390x"; go build $LDFLAGS -o "$BuildDir/$BinaryName-linux-s390x" main.go

# Darwin builds
Write-Host "Building Darwin binaries..." -ForegroundColor Cyan
$env:GOOS = "darwin"; $env:GOARCH = "amd64"; go build $LDFLAGS -o "$BuildDir/$BinaryName-darwin-amd64" main.go
$env:GOOS = "darwin"; $env:GOARCH = "arm64"; go build $LDFLAGS -o "$BuildDir/$BinaryName-darwin-arm64" main.go

# FreeBSD builds
Write-Host "Building FreeBSD binaries..." -ForegroundColor Cyan
$env:GOOS = "freebsd"; $env:GOARCH = "amd64"; go build $LDFLAGS -o "$BuildDir/$BinaryName-freebsd-amd64" main.go
$env:GOOS = "freebsd"; $env:GOARCH = "arm64"; go build $LDFLAGS -o "$BuildDir/$BinaryName-freebsd-arm64" main.go

# Windows builds
Write-Host "Building Windows binaries..." -ForegroundColor Cyan
$env:GOOS = "windows"; $env:GOARCH = "amd64"; go build $LDFLAGS -o "$BuildDir/$BinaryName-windows-amd64.exe" main.go
$env:GOOS = "windows"; $env:GOARCH = "arm64"; go build $LDFLAGS -o "$BuildDir/$BinaryName-windows-arm64.exe" main.go

# Clear environment variables
Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
Remove-Item Env:GOARM -ErrorAction SilentlyContinue

Write-Host "Build complete!" -ForegroundColor Green
Write-Host "Binaries are in $BuildDir/"

# Create packages directory
New-Item -ItemType Directory -Path "$BuildDir/packages" -Force | Out-Null

Write-Host "Creating packages..." -ForegroundColor Yellow

# Linux packages
Write-Host "Creating Linux packages..." -ForegroundColor Cyan
Set-Location $BuildDir
tar -czf "packages/$BinaryName-linux-amd64.tar.gz" "$BinaryName-linux-amd64"
tar -czf "packages/$BinaryName-linux-arm64.tar.gz" "$BinaryName-linux-arm64"
tar -czf "packages/$BinaryName-linux-armv7.tar.gz" "$BinaryName-linux-armv7"
tar -czf "packages/$BinaryName-linux-ppc64le.tar.gz" "$BinaryName-linux-ppc64le"
tar -czf "packages/$BinaryName-linux-s390x.tar.gz" "$BinaryName-linux-s390x"

# Darwin packages
Write-Host "Creating Darwin packages..." -ForegroundColor Cyan
tar -czf "packages/$BinaryName-darwin-amd64.tar.gz" "$BinaryName-darwin-amd64"
tar -czf "packages/$BinaryName-darwin-arm64.tar.gz" "$BinaryName-darwin-arm64"

# FreeBSD packages
Write-Host "Creating FreeBSD packages..." -ForegroundColor Cyan
tar -czf "packages/$BinaryName-freebsd-amd64.tar.gz" "$BinaryName-freebsd-amd64"
tar -czf "packages/$BinaryName-freebsd-arm64.tar.gz" "$BinaryName-freebsd-arm64"

# Windows packages
Write-Host "Creating Windows packages..." -ForegroundColor Cyan
Compress-Archive -Path "$BinaryName-windows-amd64.exe" -DestinationPath "packages/$BinaryName-windows-amd64.zip" -Force
Compress-Archive -Path "$BinaryName-windows-arm64.exe" -DestinationPath "packages/$BinaryName-windows-arm64.zip" -Force

Set-Location ..

Write-Host "Packages created!" -ForegroundColor Green

# Generate checksums
Write-Host "Generating checksums..." -ForegroundColor Yellow
Set-Location "$BuildDir/packages"

# PowerShell doesn't have a built-in sha256sum, so we'll use .NET
$checksums = @()
Get-ChildItem -File | ForEach-Object {
    $hash = Get-FileHash -Algorithm SHA256 -Path $_.FullName
    $checksums += "$($hash.Hash.ToLower())  $($_.Name)"
}

$checksums | Out-File -FilePath "../checksums.sha256" -Encoding UTF8

Set-Location ../..

Write-Host "Checksums generated: $BuildDir/checksums.sha256" -ForegroundColor Green

# Show results
Write-Host "Build Summary:" -ForegroundColor Green
$binaryCount = (Get-ChildItem "$BuildDir" -File | Where-Object { $_.Name -notlike "*packages*" }).Count
$packageCount = (Get-ChildItem "$BuildDir/packages" -File).Count
Write-Host "Binaries: $binaryCount"
Write-Host "Packages: $packageCount"
Write-Host ""

Write-Host "Files created:" -ForegroundColor Cyan
Get-ChildItem $BuildDir | Format-Table Name, Length, LastWriteTime

Write-Host "Packages created:" -ForegroundColor Cyan
Get-ChildItem "$BuildDir/packages" | Format-Table Name, Length, LastWriteTime

Write-Host "Checksums:" -ForegroundColor Cyan
Get-Content "$BuildDir/checksums.sha256"

Write-Host "Build process completed successfully!" -ForegroundColor Green
