# MremoteGO# MremoteGO



> A modern, cross-platform remote connection manager. Like mRemoteNG, but with git-friendly YAML configs and 1Password integration.A modern, cross-platform remote connection manager written in Go. Think mRemoteNG, but with git-friendly YAML configs and 1Password integration.



[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)



## Why MremoteGO?## Why MremoteGO?



**The Problem**: mRemoteNG uses XML configs that are painful to diff, merge, and share with teams. Passwords are awkwardly encrypted per-machine.**Problem**: mRemoteNG uses XML configs that are painful to diff, merge, and share with teams.



**The Solution**: MremoteGO uses clean YAML configs that work beautifully with git, plus 1Password integration for secure team password sharing.**Solution**: MremoteGO uses clean YAML configs that work beautifully with git, plus 1Password integration for secure password sharing.



## Features### Key Features



- 🎨 **Modern GUI** - Clean interface with connection tree, drag-and-drop organization- 🎨 **Modern GUI** - Clean interface with connection tree and quick actions

- 🔐 **1Password Integration** - `op://` references keep passwords secure- 🔐 **1Password Integration** - Store passwords securely, share configs safely

- 🔑 **RDP Auto-Login** - Windows Credential Manager for seamless connections- � **RDP Auto-Login** - Windows Credential Manager integration

- 📝 **Git-Friendly** - YAML configs are easy to diff, merge, and review- � **Git-Friendly** - YAML configs are easy to diff and merge

- 🖥️ **Cross-Platform** - Windows, Linux, macOS- 🖥️ **Cross-Platform** - Windows, Linux, macOS

- ⚡ **Fast & Clean** - No console popups, instant connections- ⚡ **Fast** - No console window popups, instant connections

- 🚀 **Protocols** - SSH (PuTTY), RDP, VNC, HTTP/HTTPS, Telnet- � **Organize** - Folders, drag-and-drop, search

- 📁 **Organized** - Folders, search, recent files- � **Multiple Protocols** - SSH, RDP, VNC, HTTP/HTTPS, Telnet



## Quick Start## Screenshots



### 1. Install & Build### Main Interface

```

```powershell┌─────────────────────────────────────────────────────────┐

git clone https://github.com/yourusername/mremotego│ File  Connection  Help                                   │

cd mremotego│ [+] [📁] [▶] [✏️] [🗑️] [🔄]                              │

go build -ldflags "-H windowsgui" -o MremoteGO.exe ./cmd/mremotego-gui├───────────────┬─────────────────────────────────────────┤

.\MremoteGO.exe│ 📁 Production │ Connection Details                      │

```│  🔐 Web1      │ 🔐 Web1                                 │

│  🔐 DB1       │ Protocol: ssh                           │

### 2. Add a Connection│ 📁 Development│ Host: web1.prod.com                     │

│  🔐 DevServer │ Port: 22                                │

Click **[+] Add** → Fill in details → **Submit** → **[▶] Connect**└───────────────┴─────────────────────────────────────────┘

```

That's it! Auto-login works automatically.

See [GUI-README.md](GUI-README.md) for GUI documentation.

### 3. Optional: Set up 1Password

## Installation

```powershell

# Install 1Password CLI### Quick Start (GUI)

winget install 1Password.CLI

```bash

# Enable integration in 1Password → Settings → Developer# Build GUI version (Windows with no console window)

```go build -ldflags "-H windowsgui" -o MremoteGO.exe ./cmd/mremotego-gui



Use passwords like: `op://Private/Server Name/password`# Run

.\MremoteGO.exe

**📖 Full Guide**: [docs/QUICKSTART.md](docs/QUICKSTART.md)```



## Configuration ExampleThe GUI will automatically:

- Create a default config at `%APPDATA%\mremotego\config.yaml`

### YAML (Git-Friendly)- Remember your last opened file

- Support drag-and-drop connection organization

```yaml- Hide console windows for background processes

version: "1.0"

connections:### Quick Start (CLI)

  - name: "Production"

    type: folder```bash

    children:# Build CLI version

      - name: "Web Server"go build -o mremotego.exe cmd/mremotego/main.go

        type: connection

        protocol: ssh# Initialize configuration

        host: web.prod.com.\mremotego.exe init

        username: admin

        password: op://Shared/Web Server/password# List connections

      .\mremotego.exe list

      - name: "Windows RDP"```

        type: connection

        protocol: rdp## Usage

        host: win.prod.com

        username: Administrator### Initialize configuration

        password: op://Private/Windows Server/password

``````bash

mremotego init

### Comparison with mRemoteNG```



| Feature | mRemoteNG | MremoteGO |### List all connections

|---------|-----------|-----------|

| Config Format | XML | ✅ YAML |```bash

| Git Diffs | ❌ Messy | ✅ Clean |mremotego list

| Password Storage | Per-machine DPAPI | ✅ 1Password |```

| Team Sharing | ❌ Difficult | ✅ Easy |

| Auto-Login | ✅ | ✅ |### Add a new connection

| Cross-Platform | ❌ Windows only | ✅ All platforms |

```bash

## Supported Protocols# Add an SSH connection

mremotego add --name "Production Server" --protocol ssh --host 192.168.1.100 --port 22 --username admin

| Protocol | Windows | Linux/Mac | Auto-Login |

|----------|---------|-----------|------------|# Add an RDP connection

| **SSH** | PuTTY `-pw` | Native ssh | ✅ Yes |mremotego add --name "Windows Server" --protocol rdp --host server.example.com --port 3389 --username Administrator

| **RDP** | mstsc + CredMan | xfreerdp | ✅ Yes |

| **VNC** | vncviewer | vncviewer | ✅ Yes |# Add to a folder

| **HTTP/HTTPS** | Default browser | Default browser | N/A |mremotego add --name "Dev DB" --protocol ssh --host db.dev.local --folder "Development/Databases"

| **Telnet** | Native telnet | Native telnet | ✅ Yes |```



## 1Password Integration### Connect to a host



### Why 1Password?```bash

mremotego connect "Production Server"

- ✅ Passwords stay secure (not in config files)```

- ✅ Safe to commit configs to git

- ✅ Team sharing with access control### Edit a connection

- ✅ Biometric unlock

- ✅ Audit logs```bash

- ✅ Auto-rotation supportmremotego edit "Production Server" --host 192.168.1.101 --port 2222

```

### Example

### Delete a connection

```yaml

# Config file (safe to commit to git)```bash

connections:mremotego delete "Production Server"

  - name: "Production DB"```

    password: op://DevOps/Production DB/password

```### Export connections



When you connect:```bash

1. MremoteGO calls `op read op://...`mremotego export --output connections-backup.yaml

2. 1Password authenticates with biometric unlock```

3. Password is retrieved securely

4. Connection launches with auto-login## Configuration



**📖 Setup Guide**: [docs/1PASSWORD-SETUP.md](docs/1PASSWORD-SETUP.md)The configuration file is stored at `~/.config/mremotego/config.yaml` (Linux/Mac) or `%APPDATA%\mremotego\config.yaml` (Windows).



## Documentation### Example Configuration



| Document | Description |```yaml

|----------|-------------|version: "1.0"

| [Quick Start](docs/QUICKSTART.md) | Get started in 5 minutes |connections:

| [GUI Guide](docs/GUI-GUIDE.md) | Complete GUI reference |  - name: "Production Servers"

| [1Password Setup](docs/1PASSWORD-SETUP.md) | Secure password management |    type: folder

| [Password Management](docs/PASSWORD-MANAGEMENT.md) | Security details |    children:

      - name: "Web Server 1"

## Project Structure        type: connection

        protocol: ssh

```        host: web1.prod.example.com

mremotego/        port: 22

├── cmd/        username: deploy

│   ├── mremotego/          # CLI application        password: op://Private/Web Server 1/password  # 1Password reference

│   └── mremotego-gui/      # GUI application        description: "Primary web server"

├── internal/        

│   ├── config/             # YAML config management      - name: "Database Server"

│   ├── gui/                # Fyne GUI components        type: connection

│   ├── launcher/           # Protocol launchers        protocol: rdp

│   └── secrets/            # 1Password integration        host: db.prod.example.com

├── pkg/        port: 3389

│   └── models/             # Data models        username: Administrator

├── docs/                   # Documentation        password: op://Private/DB Server/password  # Secure password storage

├── MremoteGO.exe           # Built GUI app        

└── config.example.yaml     # Example config  - name: "Development"

```    type: folder

    children:

## Building      - name: "Dev SSH"

        type: connection

### GUI (Recommended)        protocol: ssh

        host: dev.example.com

```powershell        port: 22

# Windows (no console window)        username: developer

go build -ldflags "-H windowsgui" -o MremoteGO.exe ./cmd/mremotego-gui        password: plaintext_password_here  # Or plain text (not recommended)

```

# Linux/Mac

go build -o mremotego-gui ./cmd/mremotego-gui### 1Password Integration

```

Store passwords securely in 1Password instead of config files:

### CLI

```yaml

```powershellpassword: op://Private/Server Name/password

go build -o mremotego.exe ./cmd/mremotego```

```

See [1PASSWORD-CLI-SETUP.md](1PASSWORD-CLI-SETUP.md) for setup instructions.

## Requirements

## Supported Protocols

- Go 1.22 or higher

- **Windows**: PuTTY (for SSH)- **SSH**: Secure Shell connections (uses PuTTY on Windows, native ssh on Mac/Linux)

- **Optional**: 1Password desktop app + CLI- **RDP**: Remote Desktop Protocol (launches mstsc on Windows, xfreerdp on Linux)

- **VNC**: Virtual Network Computing (launches vncviewer)

## Use Cases- **HTTP/HTTPS**: Web interfaces (opens in default browser)

- **Telnet**: Legacy telnet connections

### System Administrators

### Special Features

```yaml

# production-servers.yaml (committed to git)- **RDP Auto-Login**: Passwords stored in Windows Credential Manager for seamless login

connections:- **1Password Integration**: Store passwords securely using `op://vault/item/field` references

  - name: "Web Cluster"- **PuTTY on Windows**: SSH connections use PuTTY with password auto-fill support

    type: folder

    children:## Git-Friendly Format

      - name: "web-01"

        host: 10.0.1.50Unlike mRemoteNG's XML format, MremoteGO uses YAML which provides:

        password: op://DevOps/web-01/password

      - name: "web-02"- ✅ Clear diffs in version control

        host: 10.0.1.51- ✅ Easy merge conflict resolution

        password: op://DevOps/web-02/password- ✅ Human-readable format

```- ✅ Comments support

- ✅ Better organization

Team shares config via git, passwords stay in 1Password.

## Development

### DevOps Teams

### Project Structure

- Separate configs per environment (dev/staging/prod)

- Shared 1Password vaults for team credentials```

- Git-based config versioningmremotego/

- Audit trail via 1Password logs├── cmd/

│   ├── mremotego/          # CLI application

### Personal Use│   └── mremotego-gui/      # GUI application

├── internal/

- Organize home lab connections│   ├── config/             # Configuration management

- SSH keys for personal servers│   ├── gui/                # GUI components (Fyne)

- Optional: Plain text passwords (not recommended for teams)│   ├── launcher/           # Protocol launchers

│   └── secrets/            # 1Password integration

## Security├── pkg/

│   └── models/             # Data models

### What's Secure├── tools/

│   └── generate_icon.go    # Icon generator

✅ 1Password integration - Passwords never in config files  └── go.mod

✅ RDP Credential Manager - OS-level secure storage  ```

✅ Process hiding - No console windows exposing commands  

✅ Biometric unlock - Touch ID/Windows Hello via 1Password  ### Build



### What's Not```bash

# GUI with hidden console window (Windows)

⚠️ Plain text passwords - Visible in config file  go build -ldflags "-H windowsgui" -o MremoteGO.exe ./cmd/mremotego-gui

⚠️ Config file permissions - User responsibility  

⚠️ Git commits - Don't commit plain text passwords  # CLI

go build -o mremotego.exe ./cmd/mremotego

**Recommendation**: Always use 1Password for team environments.```



**📖 Details**: [docs/PASSWORD-MANAGEMENT.md](docs/PASSWORD-MANAGEMENT.md)### Test



## Contributing```bash

go test ./...

Contributions welcome! Please:```



1. Fork the repository## License

2. Create a feature branch

3. Make your changesMIT License

4. Submit a pull request

## Contributing

## License

Contributions are welcome! Please feel free to submit a Pull Request.

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [mRemoteNG](https://mremoteng.org/)
- Built with [Fyne](https://fyne.io/) GUI toolkit
- 1Password integration via [1Password CLI](https://developer.1password.com/docs/cli/)

## Support

- 📖 Documentation: `docs/` folder
- 🐛 Issues: GitHub Issues
- 💬 Discussions: GitHub Discussions

---

**Made with ❤️ in Go**
