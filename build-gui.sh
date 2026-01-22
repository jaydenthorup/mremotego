#!/bin/bash
# Build script for MremoteGO GUI (Linux/Mac)

echo "Building MremoteGO GUI..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Build the GUI application
go build -o bin/mremotego-gui cmd/mremotego-gui/main.go cmd/mremotego-gui/theme.go

if [ $? -eq 0 ]; then
    echo "✓ Build successful: bin/mremotego-gui"
    
    # Make executable
    chmod +x bin/mremotego-gui
    
    # Show file size
    ls -lh bin/mremotego-gui | awk '{print "File size: " $5}'
else
    echo "✗ Build failed"
    exit 1
fi
