# build-and-install.ps1
# Pulls the latest code from GitHub, builds the Tauri app, and installs it on Windows.
# Run from the project root: .\build-and-install.ps1

param(
    [string]$Branch = "main"
)

$ErrorActionPreference = "Stop"

Write-Host "=== Timer Build & Install ===" -ForegroundColor Cyan

# Ensure we're in the repo root
if (-not (Test-Path "src-tauri/tauri.conf.json")) {
    Write-Host "Error: Run this script from the Timer project root." -ForegroundColor Red
    exit 1
}

# Pull latest
Write-Host "`n[1/4] Pulling latest from origin/$Branch..." -ForegroundColor Yellow
git checkout $Branch
git pull origin $Branch

# Install frontend dependencies
Write-Host "`n[2/4] Installing npm dependencies..." -ForegroundColor Yellow
npm install

# Build the Tauri app (produces MSI and NSIS installers)
Write-Host "`n[3/4] Building Tauri app..." -ForegroundColor Yellow
npx tauri build

# Find and run the installer
Write-Host "`n[4/4] Locating installer..." -ForegroundColor Yellow

$nsis = Get-ChildItem -Path "src-tauri/target/release/bundle/nsis" -Filter "*.exe" -ErrorAction SilentlyContinue | Select-Object -First 1
$msi  = Get-ChildItem -Path "src-tauri/target/release/bundle/msi"  -Filter "*.msi" -ErrorAction SilentlyContinue | Select-Object -First 1

if ($nsis) {
    Write-Host "Found NSIS installer: $($nsis.FullName)" -ForegroundColor Green
    Write-Host "Launching installer..." -ForegroundColor Yellow
    Start-Process -FilePath $nsis.FullName -Wait
    Write-Host "Done!" -ForegroundColor Green
} elseif ($msi) {
    Write-Host "Found MSI installer: $($msi.FullName)" -ForegroundColor Green
    Write-Host "Launching installer..." -ForegroundColor Yellow
    Start-Process -FilePath "msiexec.exe" -ArgumentList "/i", "`"$($msi.FullName)`"" -Wait
    Write-Host "Done!" -ForegroundColor Green
} else {
    Write-Host "No installer found. Check the build output above for errors." -ForegroundColor Red
    exit 1
}
