# Build script for MremoteGO GUI (Windows)

Write-Host "Building MremoteGO GUI..." -ForegroundColor Green

# Create bin directory if it doesn't exist
if (!(Test-Path -Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

# Build the GUI application
# -ldflags "-H windowsgui" prevents the console window from appearing on Windows
if ($IsWindows -or $env:OS -match "Windows") {
    go build -ldflags "-H windowsgui" -o mremotego.exe cmd/mremotego-gui/main.go cmd/mremotego-gui/theme.go
    $outputFile = "mremotego.exe"
} else {
    go build -o mremotego cmd/mremotego-gui/main.go cmd/mremotego-gui/theme.go
    $outputFile = "mremotego"
}

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Build successful: $outputFile" -ForegroundColor Green
    
    # Show file size
    $fileSize = (Get-Item $outputFile).Length / 1MB
    Write-Host "File size: $([math]::Round($fileSize, 2)) MB" -ForegroundColor Cyan
} else {
    Write-Host "✗ Build failed" -ForegroundColor Red
    exit 1
}
