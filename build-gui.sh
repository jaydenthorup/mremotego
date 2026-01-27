#!/bin/bash
# Build script for MremoteGO GUI (Linux/Mac/WSL)

echo "Building MremoteGO GUI..."

# Build the GUI application
go build -o mremotego cmd/mremotego-gui/main.go cmd/mremotego-gui/theme.go

if [ $? -eq 0 ]; then
    echo "✓ Build successful: mremotego"
    
    # Make executable
    chmod +x mremotego
    
    # Show file size
    ls -lh mremotego | awk '{print "File size: " $5}'
else
    echo "✗ Build failed"
    exit 1
fi
