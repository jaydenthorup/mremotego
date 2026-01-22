# Build script for MremoteGO GUI (Windows)

Write-Host "Building MremoteGO GUI..." -ForegroundColor Green

# Create bin directory if it doesn't exist
if (!(Test-Path -Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

# Build the GUI application
go build -o bin/mremotego-gui.exe cmd/mremotego-gui/main.go cmd/mremotego-gui/theme.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Build successful: bin/mremotego-gui.exe" -ForegroundColor Green
    
    # Show file size
    $fileSize = (Get-Item "bin/mremotego-gui.exe").Length / 1MB
    Write-Host "File size: $([math]::Round($fileSize, 2)) MB" -ForegroundColor Cyan
} else {
    Write-Host "✗ Build failed" -ForegroundColor Red
    exit 1
}
