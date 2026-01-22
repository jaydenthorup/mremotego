# MremoteGO GUI Application

The GUI version of MremoteGO provides a graphical interface similar to mRemoteNG for managing your remote connections.

## Features

- **Tree View**: Organize connections in folders with an intuitive tree structure
- **Visual Icons**: Protocol-specific icons for easy identification (🔐 SSH, 🖥️ RDP, 📺 VNC, 🌐 HTTP/HTTPS)
- **Details Panel**: View connection details at a glance
- **Quick Connect**: Double-click or use the toolbar to launch connections
- **Drag & Drop**: *(Future feature)* Reorganize connections easily
- **Search**: *(Future feature)* Quickly find connections

## Installation

### Prerequisites

The GUI application requires:
- **Windows**: No additional requirements
- **Linux**: Install required libraries:
  ```bash
  sudo apt-get install libgl1-mesa-dev xorg-dev
  ```
- **macOS**: Xcode command line tools

### Building from Source

```bash
# Build the GUI application
go build -o mremotego-gui.exe cmd/mremotego-gui/main.go cmd/mremotego-gui/theme.go

# Or use the provided build script
# Windows
.\build-gui.ps1

# Linux/Mac
./build-gui.sh
```

## Usage

### Starting the Application

Simply run the executable:

```bash
# Windows
.\mremotego-gui.exe

# Linux/Mac
./mremotego-gui
```

On first run, if no configuration file exists, you'll be prompted to create one.

### Main Interface

```
┌─────────────────────────────────────────────────────────┐
│ File  Connection  Help                                   │
│ [+] [📁] [▶] [✏️] [🗑️] [🔄]                              │
├───────────────┬─────────────────────────────────────────┤
│ 📁 Production │ Connection Details                      │
│  🔐 Web1      │                                         │
│  🔐 Web2      │ 🔐 Web1                                 │
│  🔐 DB1       │ Protocol: ssh                           │
│               │ Host: web1.prod.com                     │
│ 📁 Dev        │ Port: 22                                │
│  🔐 DevServer │ Username: deploy                        │
│               │                                         │
│ 🖥️ Windows    │ Description: Production web server     │
│               │                                         │
│               │ [Connect]                               │
└───────────────┴─────────────────────────────────────────┘
```

### Toolbar Actions

- **[+]** Add New Connection - Opens dialog to create a new connection
- **[📁]** Add New Folder - Create a folder to organize connections
- **[▶]** Connect - Launch the selected connection
- **[✏️]** Edit - Modify connection details
- **[🗑️]** Delete - Remove the selected connection or folder
- **[🔄]** Refresh - Reload configuration from disk

### Managing Connections

#### Adding a Connection

1. Click the **[+]** button or **File → New Connection**
2. Fill in the connection details:
   - **Name**: Display name for the connection
   - **Protocol**: Choose from ssh, rdp, vnc, http, https, telnet
   - **Host**: Hostname or IP address
   - **Port**: Leave empty for default port
   - **Username**: Login username
   - **Password**: (Optional) Password stored in plain text
   - **Domain**: (For RDP) Windows domain
   - **Description**: Additional notes

3. Click **Submit**

#### Organizing with Folders

1. Click **[📁]** or **File → New Folder**
2. Enter folder name
3. *(Future)* Drag connections into folders

#### Connecting to a Host

1. Select a connection in the tree
2. Click **[▶]** button, or
3. Use **Connection → Connect** menu, or
4. Double-click the connection *(Future)*

The application will launch the appropriate client:
- **SSH**: Opens SSH client in terminal
- **RDP**: Launches mstsc (Windows) or xfreerdp (Linux/Mac)
- **VNC**: Opens VNC viewer
- **HTTP/HTTPS**: Opens in default browser

#### Editing Connections

1. Select a connection
2. Click **[✏️]** button or **Connection → Edit**
3. Modify details in the dialog
4. Click **Submit** to save changes

#### Deleting Items

1. Select a connection or folder
2. Click **[🗑️]** button or **Connection → Delete**
3. Confirm deletion

**Note**: Deleting a folder removes all connections within it.

## Configuration

The GUI application uses the same YAML configuration as the CLI version:

- **Windows**: `%APPDATA%\mremotego\config.yaml`
- **Linux/Mac**: `~/.config/mremotego/config.yaml`

You can edit the file manually or use the GUI exclusively - both work seamlessly together.

## Keyboard Shortcuts

*(Future feature)*

- `Ctrl+N` - New Connection
- `Ctrl+F` - New Folder
- `Enter` - Connect to selected
- `Ctrl+E` - Edit selected
- `Delete` - Delete selected
- `Ctrl+R` - Refresh
- `Ctrl+Q` - Quit

## Tips & Tricks

### Switching Between GUI and CLI

The same configuration file is used by both versions:

```bash
# Add connection via CLI
mremotego add --name "NewServer" --protocol ssh --host server.com

# Open GUI to see the new connection
mremotego-gui
```

### Git Integration

The configuration file is git-friendly:

```bash
cd %APPDATA%\mremotego  # Windows
cd ~/.config/mremotego   # Linux/Mac

git init
git add config.yaml
git commit -m "Initial connections"
```

Changes made in the GUI are immediately reflected in the YAML file with clean, readable diffs.

### Multi-User Setup

Share connections with your team:

1. Store the config in a shared git repository
2. Point the application to the shared config:
   ```bash
   # Set environment variable or use symbolic link
   mklink %APPDATA%\mremotego\config.yaml \\shared\team-connections\config.yaml
   ```

## Troubleshooting

### GUI Won't Start

**Windows**: Ensure you have graphics drivers installed

**Linux**: Install required dependencies:
```bash
sudo apt-get install libgl1-mesa-dev xorg-dev
```

**macOS**: Install Xcode command line tools:
```bash
xcode-select --install
```

### Connection Not Launching

Check that the required client is installed:
- **SSH**: `ssh` command available
- **RDP**: `mstsc` (Windows) or `xfreerdp` (Linux/Mac)
- **VNC**: VNC viewer installed
- **Telnet**: telnet client installed

### Configuration Not Saving

Check file permissions on the configuration directory:
```bash
# Windows
icacls %APPDATA%\mremotego

# Linux/Mac
ls -la ~/.config/mremotego
```

## Future Enhancements

Planned features for future releases:

- [ ] Drag & drop connection reorganization
- [ ] Search/filter connections
- [ ] Keyboard shortcuts
- [ ] Double-click to connect
- [ ] Connection grouping by tags
- [ ] Import from mRemoteNG XML files
- [ ] Encrypted password storage
- [ ] Session history
- [ ] Quick connect dialog
- [ ] Custom themes
- [ ] System tray integration
- [ ] Multi-tab connection launcher *(embedded terminals)*

## Comparison with mRemoteNG

| Feature | mRemoteNG | MremoteGO GUI |
|---------|-----------|---------------|
| Protocol Support | ✅ | ✅ |
| Connection Tree | ✅ | ✅ |
| Git-Friendly Config | ❌ (XML) | ✅ (YAML) |
| Cross-Platform | ❌ (Windows only) | ✅ |
| Embedded Terminals | ✅ | ⏳ (Planned) |
| Password Encryption | ✅ | ⏳ (Planned) |
| Import/Export | ✅ | ⏳ (Planned) |

## Contributing

Contributions are welcome! Areas where help is needed:

- Additional protocol support
- Embedded terminal tabs
- Password encryption
- UI/UX improvements
- Testing on different platforms

## License

MIT License - See LICENSE file for details
