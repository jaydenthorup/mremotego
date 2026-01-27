# MremoteGO - Project Summary

## What is MremoteGO?

MremoteGO is a **complete reimplementation of mRemoteNG in Go** with major improvements: **git-compatible configuration files** and **1Password integration**. Unlike mRemoteNG's XML format which is difficult to diff and merge, MremoteGO uses clean, human-readable YAML files that work perfectly with version control systems.

## Key Highlights

### 🎯 Two Interfaces, One Config

- **GUI Application** (`MremoteGO.exe`) - Full graphical interface with tree view, similar to mRemoteNG
  - Custom application icon
  - No console windows
  - Recent file tracking
  - File → Open Config menu
- **CLI Application** (`mremotego.exe`) - Command-line tool for automation and terminal workflows
- Both share the same YAML configuration file seamlessly!

### � Secure Password Management

**1Password Integration:**
```yaml
connections:
  - name: "Production Server"
    password: op://Private/Production Server/password  # Secure reference
```

- Store passwords in 1Password instead of config files
- Perfect for team sharing - commit configs to git, keep passwords secure
- Biometric unlock via 1Password desktop app
- Automatic password resolution at connect time

**RDP Auto-Login:**
- Windows Credential Manager integration
- No password prompts - automatic login
- Credentials cached per machine
- Works with both plain text and 1Password references

### �📝 Git-Friendly Configuration

**mRemoteNG (XML):**
```xml
<Node Name="Server1" Type="Connection" Descr="" Icon="mRemoteNG" Panel="General" Id="500e7d58-662a-44b4-aff0-3a4f547a3fee" Username="admin" Domain="" Password="encrypted" Protocol="SSH2" ...>
```

**MremoteGO (YAML):**
```yaml
connections:
  - name: "Server1"
    type: connection
    protocol: ssh
    host: server1.com
    port: 22
    username: admin
    password: op://Private/Server1/password  # 1Password reference
    description: "Production server"
```

### 🖥️ Cross-Platform

- ✅ Windows (PuTTY for SSH, mstsc for RDP)
- ✅ Linux (native ssh, xfreerdp for RDP)
- ✅ macOS (native ssh, Microsoft Remote Desktop)

### 🔌 Protocol Support

- **SSH** - Uses PuTTY on Windows with password auto-fill, native ssh elsewhere
- **RDP** - Windows Credential Manager for auto-login
- **VNC** - Launches vncviewer
- **HTTP/HTTPS** - Opens in default browser
- **Telnet** - Legacy terminal connections

## Project Structure

```
mremotego/
├── cmd/
│   ├── mremotego/          # CLI application
│   │   └── main.go
│   └── mremotego-gui/      # GUI application
│       └── main.go
├── internal/
│   ├── config/             # Configuration management
│   │   └── manager.go      # Load/save YAML config, recent file tracking
│   ├── gui/                # GUI components (Fyne)
│   │   ├── mainwindow.go   # Main window with tree view & menus
│   │   ├── dialogs.go      # Add/edit dialogs with 1Password support
│   │   └── icon.go         # Embedded application icon
│   ├── launcher/           # Protocol launchers
│   │   └── launcher.go     # Launch SSH, RDP, VNC with password handling
│   └── secrets/            # 1Password integration
│       └── onepassword.go  # 1Password CLI wrapper
├── pkg/
│   └── models/             # Data models
│       └── connection.go   # Connection & folder structs
├── tools/
│   └── generate_icon.go    # Icon generator script
├── icon.svg                # Application icon source
├── MremoteGO.exe           # GUI executable (built)
├── mremotego.exe           # CLI executable (built)
├── config.example.yaml     # Example configuration
├── README.md               # Main documentation
├── GUI-README.md           # GUI-specific docs
├── QUICKSTART.md           # Quick start guide
├── 1PASSWORD-CLI-SETUP.md  # 1Password setup guide
├── RDP-PASSWORD-ENCRYPTION.md  # RDP auto-login details
└── go.mod                  # Go dependencies
```

## Built With

- **Language**: Go 1.22+
- **GUI Framework**: Fyne v2.7.2 (cross-platform GUI toolkit)
- **CLI Framework**: Cobra v1.8.0 (CLI commands)
- **Config Format**: YAML v3 (human-readable, git-friendly)
- **Security**: 1Password CLI for password management
- **Windows Integration**: Credential Manager (cmdkey) for RDP auto-login

## Use Cases

### For System Administrators

```bash
# Add servers via CLI
mremotego add --name "web-prod-01" --protocol ssh --host 10.0.1.50 --username admin
mremotego add --name "web-prod-02" --protocol ssh --host 10.0.1.51 --username admin
mremotego add --name "db-prod" --protocol ssh --host 10.0.1.100 --folder "Production/Database"

# Quick connect
mremotego connect "web-prod-01"

# Or use GUI for visual management
mremotego-gui
```

### For Teams

```bash
# Store connections in git repository
cd %APPDATA%\mremotego
git init
git add config.yaml
git commit -m "Initial team connections"
git remote add origin git@github.com:company/connections.git
git push

# Team members clone the shared config
git clone git@github.com:company/connections.git ~/.config/mremotego
```

### For DevOps

```yaml
# Easy to template and generate programmatically
connections:
  - name: "{{ env }}-web-{{ index }}"
    protocol: ssh
    host: "{{ ip }}"
    username: deploy
    tags: [{{ env }}, web]
```

## Advantages Over mRemoteNG

| Feature | mRemoteNG | MremoteGO |
|---------|-----------|-----------|
| **Config Format** | XML (hard to diff) | YAML (clean diffs) |
| **Git Compatible** | ❌ | ✅ |
| **Cross-Platform** | Windows only | Windows, Linux, macOS |
| **CLI Interface** | ❌ | ✅ |
| **GUI Interface** | ✅ | ✅ |
| **Automation** | Limited | Full CLI support |
| **Modern Language** | C# (.NET) | Go (single binary) |

## Future Roadmap

### Near Term
- [ ] Drag & drop in GUI tree view
- [ ] Connection search/filter
- [ ] Import from mRemoteNG XML
- [ ] Password encryption

### Medium Term
- [ ] Embedded terminal tabs in GUI
- [ ] SSH key management
- [ ] Connection history
- [ ] Quick connect dialog

### Long Term
- [ ] Plugin system
- [ ] Shared team spaces
- [ ] Cloud sync integration
- [ ] Mobile companion app

## Building

### Build Both Versions

```bash
# CLI version
go build -o bin/mremotego.exe cmd/mremotego/main.go

# GUI version  
go build -o bin/mremotego-gui.exe cmd/mremotego-gui/main.go cmd/mremotego-gui/theme.go

# Or use build scripts
.\build-gui.ps1       # Windows
./build-gui.sh        # Linux/Mac
```

### Dependencies

CLI only requires:
- Go 1.21+
- Cobra & YAML libraries

GUI additionally requires:
- Fyne v2 (cross-platform GUI)
- Graphics libraries (automatically handled by Fyne)

## Real-World Example

```yaml
version: "1.0"
connections:
  - name: "Production"
    type: folder
    children:
      - name: "Load Balancer"
        type: connection
        protocol: ssh
        host: lb.prod.example.com
        port: 22
        username: admin
        tags: [production, critical, networking]
        
      - name: "Web Servers"
        type: folder
        children:
          - name: "web-01"
            type: connection
            protocol: ssh
            host: 10.0.1.10
            port: 22
            username: deploy
            
          - name: "web-02"
            type: connection
            protocol: ssh
            host: 10.0.1.11
            port: 22
            username: deploy
            
      - name: "Database"
        type: connection
        protocol: ssh
        host: db.prod.example.com
        port: 22
        username: dba
        description: "PostgreSQL primary"
        
  - name: "Windows Servers"
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

## Performance

- **Startup**: < 1 second
- **Config Load**: Nearly instant (YAML parsing)
- **Memory**: ~20-30 MB (CLI), ~40-60 MB (GUI)
- **Binary Size**: ~15 MB (CLI), ~25 MB (GUI with Fyne)

## Security Notes

⚠️ **Current Implementation**: Passwords are stored in plain text in the YAML file.

**Recommendations:**
- Use SSH keys instead of passwords where possible
- Set appropriate file permissions (0600)
- Consider password encryption (planned feature)
- Store sensitive configs in encrypted filesystems

## Contributing

We welcome contributions! Areas where help is needed:

1. **Password Encryption** - Implement secure password storage
2. **Import/Export** - Add mRemoteNG XML import
3. **Protocol Support** - Add more connection types
4. **Testing** - Platform-specific testing
5. **Documentation** - Improve guides and examples

## License

MIT License - See LICENSE file for details.

## Acknowledgments

- Inspired by mRemoteNG
- Built with Fyne GUI framework
- Uses Cobra for CLI commands
- YAML configuration via gopkg.in/yaml.v3

---

**Made with ❤️ and Go**
