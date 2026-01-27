# MremoteGO# MremoteGO# MremoteGO# MremoteGO



> A modern, cross-platform remote connection manager with git-friendly YAML configs and 1Password integration.



[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)> A modern, cross-platform remote connection manager with git-friendly YAML configs and 1Password integration.

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)



## Why MremoteGO?

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)> A modern, cross-platform remote connection manager. Like mRemoteNG, but with git-friendly YAML configs and 1Password integration.A modern, cross-platform remote connection manager written in Go. Think mRemoteNG, but with git-friendly YAML configs and 1Password integration.

**Problem**: mRemoteNG uses XML configs that are painful to diff, merge, and share with teams.

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Solution**: MremoteGO uses clean YAML configs that work beautifully with git, plus 1Password integration for secure password sharing.



## ✨ Features

## Why MremoteGO?

- 🎨 **Modern GUI** - Clean interface with connection tree, search, and quick actions

- 🔐 **Password Encryption** - AES-256-GCM encryption at rest with master password[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)

- 🔑 **1Password Integration** - Store passwords securely using `op://` references

- 📝 **Git-Friendly** - YAML configs are easy to diff, merge, and review**Problem**: mRemoteNG uses XML configs that are painful to diff, merge, and share with teams.

- 🖥️ **Cross-Platform** - Windows, Linux, macOS

- ⚡ **Fast** - No console window popups, instant connections[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

- 🚀 **Multiple Protocols** - SSH, RDP, VNC, HTTP/HTTPS, Telnet

- 📁 **Organized** - Folders, drag-and-drop, search filtering**Solution**: MremoteGO uses clean YAML configs that work beautifully with git, plus 1Password integration for secure password sharing.

- 🔒 **Auto-Login** - Windows Credential Manager for RDP, password support for SSH



## 🚀 Quick Start

## ✨ Features

### Installation

## Why MremoteGO?## Why MremoteGO?

```bash

# Clone the repository- 🎨 **Modern GUI** - Clean interface with connection tree, search, and quick actions

git clone https://github.com/jaydenthorup/mremotego.git

cd mremotego- 🔐 **Password Encryption** - AES-256-GCM encryption at rest with master password



# Build GUI (Windows - no console window)- 🔑 **1Password Integration** - Store passwords securely using `op://` references

.\build-gui.ps1

- 📝 **Git-Friendly** - YAML configs are easy to diff, merge, and review**The Problem**: mRemoteNG uses XML configs that are painful to diff, merge, and share with teams. Passwords are awkwardly encrypted per-machine.**Problem**: mRemoteNG uses XML configs that are painful to diff, merge, and share with teams.

# Build GUI (Linux/Mac)

./build-gui.sh- 🖥️ **Cross-Platform** - Windows, Linux, macOS



# Run- ⚡ **Fast** - No console window popups, instant connections

.\mremotego.exe

```- 🚀 **Multiple Protocols** - SSH, RDP, VNC, HTTP/HTTPS, Telnet



### First Use- 📁 **Organized** - Folders, drag-and-drop, search filtering**The Solution**: MremoteGO uses clean YAML configs that work beautifully with git, plus 1Password integration for secure team password sharing.**Solution**: MremoteGO uses clean YAML configs that work beautifully with git, plus 1Password integration for secure password sharing.



1. Launch MremoteGO- 🔒 **Auto-Login** - Windows Credential Manager for RDP, password support for SSH

2. Enter a master password (optional - for encryption)

3. Click **[+]** to add a connection

4. Fill in host details

5. Click **[▶]** to connect## 🚀 Quick Start



That's it! 🎉## Features### Key Features



## 📖 Documentation### Installation



- [Quick Start Guide](docs/QUICKSTART.md) - Get started in 5 minutes

- [GUI Guide](docs/GUI-GUIDE.md) - Complete GUI reference

- [Encryption Guide](docs/ENCRYPTION.md) - Password encryption details```bash

- [1Password Setup](docs/1PASSWORD-SETUP.md) - Secure password management

- [Password Management](docs/PASSWORD-MANAGEMENT.md) - Security best practices# Clone the repository- 🎨 **Modern GUI** - Clean interface with connection tree, drag-and-drop organization- 🎨 **Modern GUI** - Clean interface with connection tree and quick actions



## 🔐 Password Storage Optionsgit clone https://github.com/jaydenthorup/mremotego.git



MremoteGO supports three ways to store passwords:cd mremotego- 🔐 **1Password Integration** - `op://` references keep passwords secure- 🔐 **1Password Integration** - Store passwords securely, share configs safely



### 1. Encrypted (Recommended)

```yaml

password: enc:AaBbCcDd...  # AES-256-GCM encrypted# Build GUI (Windows - no console window)- 🔑 **RDP Auto-Login** - Windows Credential Manager for seamless connections- � **RDP Auto-Login** - Windows Credential Manager integration

```

- ✅ Secure at rest.\build-gui.ps1

- ✅ Master password required to decrypt

- ✅ Safe for personal use- 📝 **Git-Friendly** - YAML configs are easy to diff, merge, and review- � **Git-Friendly** - YAML configs are easy to diff and merge



### 2. 1Password Reference (Best for Teams)# Build GUI (Linux/Mac)

```yaml

password: op://DevOps/server01/password./build-gui.sh- 🖥️ **Cross-Platform** - Windows, Linux, macOS- 🖥️ **Cross-Platform** - Windows, Linux, macOS

```

- ✅ Passwords never in config files

- ✅ Team sharing with access control

- ✅ Biometric unlock# Run- ⚡ **Fast & Clean** - No console popups, instant connections- ⚡ **Fast** - No console window popups, instant connections

- ✅ Audit logs

.\mremotego.exe

### 3. Plain Text (Not Recommended)

```yaml```- 🚀 **Protocols** - SSH (PuTTY), RDP, VNC, HTTP/HTTPS, Telnet- � **Organize** - Folders, drag-and-drop, search

password: mypassword123

```

- ⚠️ Visible in config file

- ⚠️ Not safe to commit to git### First Use- 📁 **Organized** - Folders, search, recent files- � **Multiple Protocols** - SSH, RDP, VNC, HTTP/HTTPS, Telnet



## 📋 Configuration Example



```yaml1. Launch MremoteGO

version: "1.0"

connections:2. Enter a master password (optional - for encryption)

  # SSH with encrypted password

  - name: Production Web Server3. Click **[+]** to add a connection## Quick Start## Screenshots

    type: connection

    protocol: ssh4. Fill in host details

    host: web1.prod.com

    port: 225. Click **[▶]** to connect

    username: admin

    password: enc:base64encrypteddata...

    description: Primary web server

That's it! 🎉### 1. Install & Build### Main Interface

  # RDP with 1Password reference

  - name: Windows Server

    type: connection

    protocol: rdp## 📖 Documentation```

    host: win.prod.com

    port: 3389

    username: Administrator

    password: op://DevOps/windows-server/password- [Quick Start Guide](docs/QUICKSTART.md) - Get started in 5 minutes```powershell┌─────────────────────────────────────────────────────────┐

    domain: MYDOMAIN

    resolution: 1920x1080- [GUI Guide](docs/GUI-GUIDE.md) - Complete GUI reference



  # Organized in folders- [Encryption Guide](docs/ENCRYPTION.md) - Password encryption detailsgit clone https://github.com/yourusername/mremotego│ File  Connection  Help                                   │

  - name: Development

    type: folder- [1Password Setup](docs/1PASSWORD-SETUP.md) - Secure password management

    children:

      - name: Dev Database- [Password Management](docs/PASSWORD-MANAGEMENT.md) - Security best practicescd mremotego│ [+] [📁] [▶] [✏️] [🗑️] [🔄]                              │

        type: connection

        protocol: ssh

        host: dev-db.local

        port: 22## 🔐 Password Storage Optionsgo build -ldflags "-H windowsgui" -o MremoteGO.exe ./cmd/mremotego-gui├───────────────┬─────────────────────────────────────────┤

        username: dbadmin

        password: op://DevOps/dev-db/password

```

MremoteGO supports three ways to store passwords:.\MremoteGO.exe│ 📁 Production │ Connection Details                      │

See [connections.example.yaml](connections.example.yaml) for more examples.



## 🌐 Supported Protocols

### 1. Encrypted (Recommended)```│  🔐 Web1      │ 🔐 Web1                                 │

| Protocol | Windows | Linux/Mac | Auto-Login |

|----------|---------|-----------|------------|```yaml

| **SSH** | PuTTY | Terminal | ✅ Yes |

| **RDP** | mstsc | xfreerdp | ✅ Yes |password: enc:AaBbCcDd...  # AES-256-GCM encrypted│  🔐 DB1       │ Protocol: ssh                           │

| **VNC** | vncviewer | vncviewer | ✅ Yes |

| **HTTP/HTTPS** | Browser | Browser | N/A |```

| **Telnet** | telnet | telnet | ✅ Yes |

- ✅ Secure at rest### 2. Add a Connection│ 📁 Development│ Host: web1.prod.com                     │

### Platform-Specific Features

- ✅ Master password required to decrypt

**Windows**:

- RDP: Stores credentials in Windows Credential Manager- ✅ Safe for personal use│  🔐 DevServer │ Port: 22                                │

- SSH: Launches PuTTY with `-pw` flag for auto-login

- GUI: No console window popups



**Linux**:### 2. 1Password Reference (Best for Teams)Click **[+] Add** → Fill in details → **Submit** → **[▶] Connect**└───────────────┴─────────────────────────────────────────┘

- SSH: Launches in terminal (gnome-terminal, xterm, konsole, etc.)

- Password authentication via sshpass```yaml



**macOS**:password: op://DevOps/server01/password```

- SSH: Launches in Terminal.app

- RDP: Opens Microsoft Remote Desktop via `rdp://` URL```



## 📊 Comparison with mRemoteNG- ✅ Passwords never in config filesThat's it! Auto-login works automatically.



| Feature | mRemoteNG | MremoteGO |- ✅ Team sharing with access control

|---------|-----------|-----------|

| Config Format | XML | ✅ YAML |- ✅ Biometric unlockSee [GUI-README.md](GUI-README.md) for GUI documentation.

| Git Diffs | ❌ Messy | ✅ Clean |

| Password Encryption | Per-machine DPAPI | ✅ AES-256-GCM |- ✅ Audit logs

| 1Password Integration | ❌ No | ✅ Yes |

| Team Sharing | ❌ Difficult | ✅ Easy |### 3. Optional: Set up 1Password

| Auto-Login | ✅ Yes | ✅ Yes |

| Cross-Platform | ❌ Windows only | ✅ All platforms |### 3. Plain Text (Not Recommended)



## 🛠️ Building from Source```yaml## Installation



### Requirementspassword: mypassword123

- Go 1.24 or higher

- Git``````powershell



### Build Commands- ⚠️ Visible in config file



```bash- ⚠️ Not safe to commit to git# Install 1Password CLI### Quick Start (GUI)

# Windows GUI (no console)

.\build-gui.ps1



# Linux/Mac GUI## 📋 Configuration Examplewinget install 1Password.CLI

./build-gui.sh



# CLI version

go build -o mremotego ./cmd/mremotego```yaml```bash



# Encryption toolversion: "1.0"

go build -o encrypt-passwords ./cmd/encrypt-passwords

```connections:# Enable integration in 1Password → Settings → Developer# Build GUI version (Windows with no console window)



## 🔧 CLI Tool  # SSH with encrypted password



MremoteGO also includes a CLI for automation:  - name: Production Web Server```go build -ldflags "-H windowsgui" -o MremoteGO.exe ./cmd/mremotego-gui



```bash    type: connection

# Initialize config

mremotego init    protocol: ssh



# List connections    host: web1.prod.com

mremotego list

    port: 22Use passwords like: `op://Private/Server Name/password`# Run

# Add connection

mremotego add --name "Server" --protocol ssh --host 192.168.1.100    username: admin



# Connect    password: enc:base64encrypteddata....\MremoteGO.exe

mremotego connect "Server"

    description: Primary web server

# Export

mremotego export --output backup.yaml**📖 Full Guide**: [docs/QUICKSTART.md](docs/QUICKSTART.md)```

```

  # RDP with 1Password reference

## 🏗️ Project Structure

  - name: Windows Server

```

mremotego/    type: connection

├── cmd/

│   ├── mremotego/          # CLI application    protocol: rdp## Configuration ExampleThe GUI will automatically:

│   ├── mremotego-gui/      # GUI application

│   └── encrypt-passwords/  # Password encryption tool    host: win.prod.com

├── internal/

│   ├── config/             # Configuration management    port: 3389- Create a default config at `%APPDATA%\mremotego\config.yaml`

│   ├── crypto/             # Encryption (AES-256-GCM)

│   ├── gui/                # Fyne GUI components    username: Administrator

│   ├── launcher/           # Protocol launchers

│   └── secrets/            # 1Password integration    password: op://DevOps/windows-server/password### YAML (Git-Friendly)- Remember your last opened file

├── pkg/

│   └── models/             # Data models    domain: MYDOMAIN

├── docs/                   # Documentation

├── build-gui.ps1          # Windows build script    resolution: 1920x1080- Support drag-and-drop connection organization

└── build-gui.sh           # Linux/Mac build script

```



## 🤝 Contributing  # Organized in folders```yaml- Hide console windows for background processes



Contributions are welcome! Please:  - name: Development



1. Fork the repository    type: folderversion: "1.0"

2. Create a feature branch (`git checkout -b feature/amazing-feature`)

3. Commit your changes (`git commit -m 'Add amazing feature'`)    children:

4. Push to the branch (`git push origin feature/amazing-feature`)

5. Open a Pull Request      - name: Dev Databaseconnections:### Quick Start (CLI)



## 📄 License        type: connection



MIT License - see [LICENSE](LICENSE) file for details.        protocol: ssh  - name: "Production"



Copyright © 2026 [Jayden Thorup](mailto:jayden.thorup@jayfiles.com)        host: dev-db.local



## 🙏 Acknowledgments        port: 22    type: folder```bash



- Inspired by [mRemoteNG](https://mremoteng.org/)        username: dbadmin

- Built with [Fyne](https://fyne.io/) GUI toolkit

- 1Password integration via [1Password CLI](https://developer.1password.com/docs/cli/)        password: op://DevOps/dev-db/password    children:# Build CLI version

- Encryption using Go's crypto libraries

```

## 💬 Support

      - name: "Web Server"go build -o mremotego.exe cmd/mremotego/main.go

- 📖 Documentation: [docs/](docs/)

- 🐛 Issues: [GitHub Issues](https://github.com/jaydenthorup/mremotego/issues)See [connections.example.yaml](connections.example.yaml) for more examples.

- 💡 Feature Requests: [GitHub Discussions](https://github.com/jaydenthorup/mremotego/discussions)

        type: connection

---

## 🌐 Supported Protocols

**Made with ❤️ in Go**

        protocol: ssh# Initialize configuration

| Protocol | Windows | Linux/Mac | Auto-Login |

|----------|---------|-----------|------------|        host: web.prod.com.\mremotego.exe init

| **SSH** | PuTTY | Terminal | ✅ Yes |

| **RDP** | mstsc | xfreerdp | ✅ Yes |        username: admin

| **VNC** | vncviewer | vncviewer | ✅ Yes |

| **HTTP/HTTPS** | Browser | Browser | N/A |        password: op://Shared/Web Server/password# List connections

| **Telnet** | telnet | telnet | ✅ Yes |

      .\mremotego.exe list

### Platform-Specific Features

      - name: "Windows RDP"```

**Windows**:

- RDP: Stores credentials in Windows Credential Manager        type: connection

- SSH: Launches PuTTY with `-pw` flag for auto-login

- GUI: No console window popups        protocol: rdp## Usage



**Linux**:        host: win.prod.com

- SSH: Launches in terminal (gnome-terminal, xterm, konsole, etc.)

- Password authentication via sshpass        username: Administrator### Initialize configuration



**macOS**:        password: op://Private/Windows Server/password

- SSH: Launches in Terminal.app

- RDP: Opens Microsoft Remote Desktop via `rdp://` URL``````bash



## 📊 Comparison with mRemoteNGmremotego init



| Feature | mRemoteNG | MremoteGO |### Comparison with mRemoteNG```

|---------|-----------|-----------|

| Config Format | XML | ✅ YAML |

| Git Diffs | ❌ Messy | ✅ Clean |

| Password Encryption | Per-machine DPAPI | ✅ AES-256-GCM || Feature | mRemoteNG | MremoteGO |### List all connections

| 1Password Integration | ❌ No | ✅ Yes |

| Team Sharing | ❌ Difficult | ✅ Easy ||---------|-----------|-----------|

| Auto-Login | ✅ Yes | ✅ Yes |

| Cross-Platform | ❌ Windows only | ✅ All platforms || Config Format | XML | ✅ YAML |```bash



## 🛠️ Building from Source| Git Diffs | ❌ Messy | ✅ Clean |mremotego list



### Requirements| Password Storage | Per-machine DPAPI | ✅ 1Password |```

- Go 1.24 or higher

- Git| Team Sharing | ❌ Difficult | ✅ Easy |



### Build Commands| Auto-Login | ✅ | ✅ |### Add a new connection



```bash| Cross-Platform | ❌ Windows only | ✅ All platforms |

# Windows GUI (no console)

.\build-gui.ps1```bash



# Linux/Mac GUI## Supported Protocols# Add an SSH connection

./build-gui.sh

mremotego add --name "Production Server" --protocol ssh --host 192.168.1.100 --port 22 --username admin

# CLI version

go build -o mremotego ./cmd/mremotego| Protocol | Windows | Linux/Mac | Auto-Login |



# Encryption tool|----------|---------|-----------|------------|# Add an RDP connection

go build -o encrypt-passwords ./cmd/encrypt-passwords

```| **SSH** | PuTTY `-pw` | Native ssh | ✅ Yes |mremotego add --name "Windows Server" --protocol rdp --host server.example.com --port 3389 --username Administrator



## 🔧 CLI Tool| **RDP** | mstsc + CredMan | xfreerdp | ✅ Yes |



MremoteGO also includes a CLI for automation:| **VNC** | vncviewer | vncviewer | ✅ Yes |# Add to a folder



```bash| **HTTP/HTTPS** | Default browser | Default browser | N/A |mremotego add --name "Dev DB" --protocol ssh --host db.dev.local --folder "Development/Databases"

# Initialize config

mremotego init| **Telnet** | Native telnet | Native telnet | ✅ Yes |```



# List connections

mremotego list

## 1Password Integration### Connect to a host

# Add connection

mremotego add --name "Server" --protocol ssh --host 192.168.1.100



# Connect### Why 1Password?```bash

mremotego connect "Server"

mremotego connect "Production Server"

# Export

mremotego export --output backup.yaml- ✅ Passwords stay secure (not in config files)```

```

- ✅ Safe to commit configs to git

## 🏗️ Project Structure

- ✅ Team sharing with access control### Edit a connection

```

mremotego/- ✅ Biometric unlock

├── cmd/

│   ├── mremotego/          # CLI application- ✅ Audit logs```bash

│   ├── mremotego-gui/      # GUI application

│   └── encrypt-passwords/  # Password encryption tool- ✅ Auto-rotation supportmremotego edit "Production Server" --host 192.168.1.101 --port 2222

├── internal/

│   ├── config/             # Configuration management```

│   ├── crypto/             # Encryption (AES-256-GCM)

│   ├── gui/                # Fyne GUI components### Example

│   ├── launcher/           # Protocol launchers

│   └── secrets/            # 1Password integration### Delete a connection

├── pkg/

│   └── models/             # Data models```yaml

├── docs/                   # Documentation

├── build-gui.ps1          # Windows build script# Config file (safe to commit to git)```bash

└── build-gui.sh           # Linux/Mac build script

```connections:mremotego delete "Production Server"



## 🤝 Contributing  - name: "Production DB"```



Contributions are welcome! Please:    password: op://DevOps/Production DB/password



1. Fork the repository```### Export connections

2. Create a feature branch (`git checkout -b feature/amazing-feature`)

3. Commit your changes (`git commit -m 'Add amazing feature'`)

4. Push to the branch (`git push origin feature/amazing-feature`)

5. Open a Pull RequestWhen you connect:```bash



## 📄 License1. MremoteGO calls `op read op://...`mremotego export --output connections-backup.yaml



MIT License - see [LICENSE](LICENSE) file for details.2. 1Password authenticates with biometric unlock```



Copyright © 2026 [Jayden Thorup](mailto:jayden.thorup@jayfiles.com)3. Password is retrieved securely



## 🙏 Acknowledgments4. Connection launches with auto-login## Configuration



- Inspired by [mRemoteNG](https://mremoteng.org/)

- Built with [Fyne](https://fyne.io/) GUI toolkit

- 1Password integration via [1Password CLI](https://developer.1password.com/docs/cli/)**📖 Setup Guide**: [docs/1PASSWORD-SETUP.md](docs/1PASSWORD-SETUP.md)The configuration file is stored at `~/.config/mremotego/config.yaml` (Linux/Mac) or `%APPDATA%\mremotego\config.yaml` (Windows).

- Encryption using Go's crypto libraries



## 💬 Support

## Documentation### Example Configuration

- 📖 Documentation: [docs/](docs/)

- 🐛 Issues: [GitHub Issues](https://github.com/jaydenthorup/mremotego/issues)

- 💡 Feature Requests: [GitHub Discussions](https://github.com/jaydenthorup/mremotego/discussions)

| Document | Description |```yaml

---

|----------|-------------|version: "1.0"

**Made with ❤️ in Go**

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
