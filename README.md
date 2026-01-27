# MremoteGO

A Go implementation of mRemoteNG with git-compatible configuration files.

**Available in two flavors:**
- **CLI** (`mremotego`) - Command-line interface for terminal lovers
- **GUI** (`mremotego-gui`) - Graphical interface similar to mRemoteNG

## Features

- 🔐 Store connection information for multiple protocols (RDP, SSH, VNC, HTTP/HTTPS, Telnet)
- 📁 Organize connections in folders/groups
- 🔄 Git-friendly YAML configuration format (easy to diff and merge)
- 🖥️ Cross-platform (Windows, Linux, macOS)
- 💻 CLI and GUI interfaces
- 🔒 1Password integration for secure password storage
- 🔑 RDP auto-login using Windows Credential Manager
- 🎨 Custom application icon
- 🚀 Fast and lightweight
- 📂 Recent file tracking

## Screenshots

### GUI Version
```
┌─────────────────────────────────────────────────────────┐
│ File  Connection  Help                                   │
│ [+] [📁] [▶] [✏️] [🗑️] [🔄]                              │
├───────────────┬─────────────────────────────────────────┤
│ 📁 Production │ Connection Details                      │
│  🔐 Web1      │ 🔐 Web1                                 │
│  🔐 DB1       │ Protocol: ssh                           │
│ 📁 Development│ Host: web1.prod.com                     │
│  🔐 DevServer │ Port: 22                                │
└───────────────┴─────────────────────────────────────────┘
```

See [GUI-README.md](GUI-README.md) for GUI documentation.

## Installation

### Quick Start (GUI)

```bash
# Build GUI version (Windows with no console window)
go build -ldflags "-H windowsgui" -o MremoteGO.exe ./cmd/mremotego-gui

# Run
.\MremoteGO.exe
```

The GUI will automatically:
- Create a default config at `%APPDATA%\mremotego\config.yaml`
- Remember your last opened file
- Support drag-and-drop connection organization
- Hide console windows for background processes

### Quick Start (CLI)

```bash
# Build CLI version
go build -o mremotego.exe cmd/mremotego/main.go

# Initialize configuration
.\mremotego.exe init

# List connections
.\mremotego.exe list
```

## Usage

### Initialize configuration

```bash
mremotego init
```

### List all connections

```bash
mremotego list
```

### Add a new connection

```bash
# Add an SSH connection
mremotego add --name "Production Server" --protocol ssh --host 192.168.1.100 --port 22 --username admin

# Add an RDP connection
mremotego add --name "Windows Server" --protocol rdp --host server.example.com --port 3389 --username Administrator

# Add to a folder
mremotego add --name "Dev DB" --protocol ssh --host db.dev.local --folder "Development/Databases"
```

### Connect to a host

```bash
mremotego connect "Production Server"
```

### Edit a connection

```bash
mremotego edit "Production Server" --host 192.168.1.101 --port 2222
```

### Delete a connection

```bash
mremotego delete "Production Server"
```

### Export connections

```bash
mremotego export --output connections-backup.yaml
```

## Configuration

The configuration file is stored at `~/.config/mremotego/config.yaml` (Linux/Mac) or `%APPDATA%\mremotego\config.yaml` (Windows).

### Example Configuration

```yaml
version: "1.0"
connections:
  - name: "Production Servers"
    type: folder
    children:
      - name: "Web Server 1"
        type: connection
        protocol: ssh
        host: web1.prod.example.com
        port: 22
        username: deploy
        password: op://Private/Web Server 1/password  # 1Password reference
        description: "Primary web server"
        
      - name: "Database Server"
        type: connection
        protocol: rdp
        host: db.prod.example.com
        port: 3389
        username: Administrator
        password: op://Private/DB Server/password  # Secure password storage
        
  - name: "Development"
    type: folder
    children:
      - name: "Dev SSH"
        type: connection
        protocol: ssh
        host: dev.example.com
        port: 22
        username: developer
        password: plaintext_password_here  # Or plain text (not recommended)
```

### 1Password Integration

Store passwords securely in 1Password instead of config files:

```yaml
password: op://Private/Server Name/password
```

See [1PASSWORD-CLI-SETUP.md](1PASSWORD-CLI-SETUP.md) for setup instructions.

## Supported Protocols

- **SSH**: Secure Shell connections (uses PuTTY on Windows, native ssh on Mac/Linux)
- **RDP**: Remote Desktop Protocol (launches mstsc on Windows, xfreerdp on Linux)
- **VNC**: Virtual Network Computing (launches vncviewer)
- **HTTP/HTTPS**: Web interfaces (opens in default browser)
- **Telnet**: Legacy telnet connections

### Special Features

- **RDP Auto-Login**: Passwords stored in Windows Credential Manager for seamless login
- **1Password Integration**: Store passwords securely using `op://vault/item/field` references
- **PuTTY on Windows**: SSH connections use PuTTY with password auto-fill support

## Git-Friendly Format

Unlike mRemoteNG's XML format, MremoteGO uses YAML which provides:

- ✅ Clear diffs in version control
- ✅ Easy merge conflict resolution
- ✅ Human-readable format
- ✅ Comments support
- ✅ Better organization

## Development

### Project Structure

```
mremotego/
├── cmd/
│   ├── mremotego/          # CLI application
│   └── mremotego-gui/      # GUI application
├── internal/
│   ├── config/             # Configuration management
│   ├── gui/                # GUI components (Fyne)
│   ├── launcher/           # Protocol launchers
│   └── secrets/            # 1Password integration
├── pkg/
│   └── models/             # Data models
├── tools/
│   └── generate_icon.go    # Icon generator
└── go.mod
```

### Build

```bash
# GUI with hidden console window (Windows)
go build -ldflags "-H windowsgui" -o MremoteGO.exe ./cmd/mremotego-gui

# CLI
go build -o mremotego.exe ./cmd/mremotego
```

### Test

```bash
go test ./...
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
