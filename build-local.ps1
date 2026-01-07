# PowerShell version for Windows users without bash

Write-Host "======================================"
Write-Host "  envswitch - Local Build"
Write-Host "======================================"
Write-Host ""

# Detect Architecture
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

# Check for ARM
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
    $arch = "arm64"
}

Write-Host "Detected system:"
Write-Host "  OS:           windows"
Write-Host "  Architecture: $arch"
Write-Host ""

Write-Host "Building for windows/$arch..."
Write-Host ""

$env:GOOS = "windows"
$env:GOARCH = $arch

go build -o "envswitch.exe" .

if ($LASTEXITCODE -eq 0) {
    $size = (Get-Item "envswitch.exe").Length / 1MB
    Write-Host "✓ Build successful!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Binary created: envswitch.exe"
    Write-Host ("Size: {0:N2} MB" -f $size)
    Write-Host ""
    Write-Host "Usage:"
    Write-Host "  .\envswitch.exe --env test"
} else {
    Write-Host "✗ Build failed!" -ForegroundColor Red
    exit 1
}

