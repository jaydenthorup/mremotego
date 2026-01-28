# Build script for MremoteGO (GUI + CLI combined)

Write-Host "Building MremoteGO..." -ForegroundColor Green

# Create bin directory if it doesn't exist
if (!(Test-Path -Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

# Build the application
# Note: Not using -H windowsgui so CLI mode works
# This means a console window will briefly appear when launching the GUI on Windows
if ($IsWindows -or $env:OS -match "Windows") {
    go build -o mremotego.exe cmd/mremotego-gui/main.go cmd/mremotego-gui/theme.go
    $outputFile = "mremotego.exe"
} else {
    go build -o mremotego cmd/mremotego-gui/main.go cmd/mremotego-gui/theme.go
    $outputFile = "mremotego"
}

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Build successful: $outputFile" -ForegroundColor Green
    Write-Host "  Run GUI: ./$outputFile" -ForegroundColor Cyan
    Write-Host "  Run CLI: ./$outputFile --help" -ForegroundColor Cyan
    
    # Show file size
    $fileSize = (Get-Item $outputFile).Length / 1MB
    Write-Host "File size: $([math]::Round($fileSize, 2)) MB" -ForegroundColor Gray
} else {
    Write-Host "✗ Build failed" -ForegroundColor Red
    exit 1
}
