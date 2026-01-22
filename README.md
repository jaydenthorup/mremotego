# MremoteGO

A Go implementation of mRemoteNG with git-compatible configuration files.

**Available in two flavors:**
- **CLI** (`mremotego`) - Command-line interface for terminal lovers
- **GUI** (`mremotego-gui`) - Graphical interface similar to mRemoteNG

## Features

- 🔐 Store connection information for multiple protocols (RDP, SSH, VNC, HTTP/HTTPS)
- 📁 Organize connections in folders/groups
- 🔄 Git-friendly YAML configuration format (easy to diff and merge)
- 🖥️ Cross-platform (Windows, Linux, macOS)
- 💻 CLI and GUI interfaces
- 🔒 Support for credential inheritance
- 🚀 Fast and lightweight

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
# Build GUI version
go build -o mremotego-gui.exe cmd/mremotego-gui/main.go cmd/mremotego-gui/theme.go

# Run
.\mremotego-gui.exe
```

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
        description: "Primary web server"
        
      - name: "Database Server"
        type: connection
        protocol: ssh
        host: db.prod.example.com
        port: 22
        username: dbadmin
        
  - name: "Windows Machines"
    type: folder
    children:
      - name: "File Server"
        type: connection
        protocol: rdp
        host: files.example.com
        port: 3389
        username: Administrator
        domain: CORP
```

## Supported Protocols

- **SSH**: Secure Shell connections
- **RDP**: Remote Desktop Protocol (launches mstsc/xfreerdp)
- **VNC**: Virtual Network Computing
- **HTTP/HTTPS**: Web interfaces
- **Telnet**: Legacy telnet connections

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
│   └── mremotego/      # Main CLI application
├── internal/
│   ├── config/         # Configuration management
│   ├── connection/     # Connection handlers
│   └── launcher/       # Protocol launchers
├── pkg/
│   └── models/         # Data models
└── go.mod
```

### Build

```bash
go build -o mremotego cmd/mremotego/main.go
```

### Test

```bash
go test ./...
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
