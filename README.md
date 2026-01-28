# MremoteGO# MremoteGO# MremoteGO# MremoteGO# MremoteGO# MremoteGO



> A modern, cross-platform remote connection manager with git-friendly YAML configs and 1Password integration.



[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)> A modern, cross-platform remote connection manager with git-friendly YAML configs and 1Password integration.

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)](https://github.com/jaydenthorup/mremotego)



## Why MremoteGO?[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)> A modern, cross-platform remote connection manager with git-friendly YAML configs and 1Password integration.



**The Problem**: mRemoteNG uses XML configs that are painful to diff, merge, and share with teams. Passwords are awkwardly encrypted per-machine.[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)



**The Solution**: MremoteGO uses clean YAML configs that work beautifully with git, plus optional 1Password integration for secure team password sharing.



## ✨ Features## Why MremoteGO?



- 🎨 **Modern GUI** - Clean interface with connection tree, search, and quick actions[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)> A modern, cross-platform remote connection manager with git-friendly YAML configs and 1Password integration.

- 🔐 **Password Encryption** - AES-256-GCM encryption at rest with master password

- 🔑 **1Password Integration** - Store passwords securely using `op://` references**The Problem**: mRemoteNG uses XML configs that are painful to diff, merge, and share with teams. Passwords are awkwardly encrypted per-machine.

- 📝 **Git-Friendly** - YAML configs are easy to diff, merge, and review

- 🖥️ **Cross-Platform** - Windows, Linux, macOS (AMD64 & ARM64)[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

- ⚡ **Fast** - Native GUI with instant connections

- 🚀 **Multiple Protocols** - SSH, RDP, VNC, HTTP/HTTPS, Telnet**The Solution**: MremoteGO uses clean YAML configs that work beautifully with git, plus 1Password integration for secure team password sharing.

- 📁 **Organized** - Folders, drag-and-drop, search filtering

- 🔒 **Auto-Login** - Credential manager integration for seamless connections



## 🚀 Quick Start## ✨ Features



### Download## Why MremoteGO?



Download the latest release for your platform from the [Releases](https://github.com/jaydenthorup/mremotego/releases) page.- 🎨 **Modern GUI** - Clean interface with connection tree, search, and quick actions



### Build from Source- 🔐 **Password Encryption** - AES-256-GCM encryption at rest with master password[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)> A modern, cross-platform remote connection manager. Like mRemoteNG, but with git-friendly YAML configs and 1Password integration.A modern, cross-platform remote connection manager written in Go. Think mRemoteNG, but with git-friendly YAML configs and 1Password integration.



```bash- 🔑 **1Password Integration** - Store passwords securely using `op://` references

# Clone the repository

git clone https://github.com/jaydenthorup/mremotego.git- 📝 **Git-Friendly** - YAML configs are easy to diff, merge, and review**Problem**: mRemoteNG uses XML configs that are painful to diff, merge, and share with teams.

cd mremotego

- 🖥️ **Cross-Platform** - Windows, Linux, macOS

# Build (all platforms)

go build -o mremotego ./cmd/mremotego-gui- ⚡ **Fast** - No console window popups, instant connections[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)



# Or use platform-specific build scripts- 🚀 **Multiple Protocols** - SSH, RDP, VNC, HTTP/HTTPS, Telnet

# Windows: .\build-gui.ps1

# Linux/Mac: ./build-gui.sh- 📁 **Organized** - Folders, drag-and-drop, search filtering**Solution**: MremoteGO uses clean YAML configs that work beautifully with git, plus 1Password integration for secure password sharing.

```

- 🔒 **Auto-Login** - Windows Credential Manager for RDP, password support for SSH

### First Run



1. Launch `mremotego` (or `mremotego.exe` on Windows)

2. Optionally set a master password for encryption## 🚀 Quick Start

3. Create your first connection or import from mRemoteNG

## ✨ Features

## 📖 Usage

### Installation

### Creating Connections

## Why MremoteGO?

**GUI Method:**

1. Click **[+]** or press `Ctrl+N````bash

2. Fill in connection details (name, protocol, host, credentials)

3. Optionally push password to 1Password# Clone the repository- 🎨 **Modern GUI** - Clean interface with connection tree, search, and quick actions

4. Click **Save**

git clone https://github.com/jaydenthorup/mremotego.git

**YAML Method:**

```yamlcd mremotego- 🔐 **Password Encryption** - AES-256-GCM encryption at rest with master password[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)

connections:

  - name: Production Server

    type: connection

    protocol: ssh# Build GUI (Windows - no console window)- 🔑 **1Password Integration** - Store passwords securely using `op://` references

    host: prod.example.com

    port: 22.\build-gui.ps1

    username: admin

    password: op://DevOps/prod-server/password  # 1Password reference- 📝 **Git-Friendly** - YAML configs are easy to diff, merge, and review**Problem**: mRemoteNG uses XML configs that are painful to diff, merge, and share with teams.

    

  - name: Dev Folder# Build GUI (Linux/Mac)

    type: folder

    children:./build-gui.sh- 🖥️ **Cross-Platform** - Windows, Linux, macOS

      - name: Dev Server

        type: connection

        protocol: rdp

        host: dev.example.com# Run- ⚡ **Fast** - No console window popups, instant connections[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

        username: developer

        password: encrypted:AES256:...  # Encrypted with master password.\mremotego.exe

```

```- 🚀 **Multiple Protocols** - SSH, RDP, VNC, HTTP/HTTPS, Telnet

### Connecting



- **Double-click** a connection in the tree

- **Right-click** → **Connect**### First Use- 📁 **Organized** - Folders, drag-and-drop, search filtering**Solution**: MremoteGO uses clean YAML configs that work beautifully with git, plus 1Password integration for secure password sharing.

- Select and press **Enter**



### Searching

1. Launch MremoteGO- 🔒 **Auto-Login** - Windows Credential Manager for RDP, password support for SSH

- Use the search box at the top

- Filter by connection name, host, or protocol2. Enter a master password (optional - for encryption)

- Results update in real-time

3. Click **[+]** to add a connection

### Command Line

4. Fill in host details

MremoteGO also supports CLI mode:

5. Click **[▶]** to connect## 🚀 Quick Start

```bash

# List all connections

mremotego list

That's it! 🎉## ✨ Features

# Connect to a specific host

mremotego connect "Production Server"



# Export connections## 📖 Documentation### Installation

mremotego export --output connections-backup.yaml



# Encrypt passwords in config

mremotego-encrypt- [Quick Start Guide](docs/QUICKSTART.md) - Get started in 5 minutes## Why MremoteGO?## Why MremoteGO?

```

- [GUI Guide](docs/GUI-GUIDE.md) - Complete GUI reference

## 🔐 Security

- [Encryption Guide](docs/ENCRYPTION.md) - Password encryption details```bash

### Password Encryption

- [1Password Setup](docs/1PASSWORD-SETUP.md) - Secure password management

MremoteGO supports two password storage methods:

- [Password Management](docs/PASSWORD-MANAGEMENT.md) - Security best practices# Clone the repository- 🎨 **Modern GUI** - Clean interface with connection tree, search, and quick actions

1. **Encrypted** (recommended for local use):

   - AES-256-GCM encryption with PBKDF2 key derivation

   - Master password required on startup

   - Passwords stored as `password: encrypted:AES256:...`## 🔐 Password Storage Optionsgit clone https://github.com/jaydenthorup/mremotego.git



2. **1Password** (recommended for teams):

   - Uses 1Password CLI for password resolution

   - Store passwords as `password: op://Vault/Item/field`MremoteGO supports three ways to store passwords:cd mremotego- 🔐 **Password Encryption** - AES-256-GCM encryption at rest with master password

   - Requires 1Password CLI installed and authenticated



### SSH Host Keys

### 1. Encrypted (Recommended)

- **New hosts**: Automatically accepted on first connection (`StrictHostKeyChecking=accept-new`)

- **Changed hosts**: Shows warning with clear instructions to verify and fix

- **Known hosts**: Standard SSH verification

```yaml# Build GUI (Windows - no console window)- 🔑 **1Password Integration** - Store passwords securely using `op://` references

## 🔧 Configuration

password: enc:AaBbCcDd...  # AES-256-GCM encrypted

### Config File Location

```.\build-gui.ps1

- **Windows**: `%APPDATA%\mremotego\config.yaml` or `.\config.yaml`

- **Linux/Mac**: `~/.config/mremotego/config.yaml` or `./config.yaml`



### Example Config- ✅ Secure at rest- 📝 **Git-Friendly** - YAML configs are easy to diff, merge, and review**The Problem**: mRemoteNG uses XML configs that are painful to diff, merge, and share with teams. Passwords are awkwardly encrypted per-machine.**Problem**: mRemoteNG uses XML configs that are painful to diff, merge, and share with teams.



```yaml- ✅ Master password required to decrypt

connections:

  - name: Infrastructure- ✅ Safe for personal use# Build GUI (Linux/Mac)

    type: folder

    children:

      - name: web-01

        type: connection### 2. 1Password Reference (Best for Teams)./build-gui.sh- 🖥️ **Cross-Platform** - Windows, Linux, macOS

        protocol: ssh

        host: web-01.example.com

        port: 22

        username: admin```yaml

        password: op://DevOps/web-01/password

        password: op://DevOps/server01/password

      - name: db-01

        type: connection```# Run- ⚡ **Fast** - No console window popups, instant connections

        protocol: ssh

        host: db-01.example.com

        port: 22

        username: postgres- ✅ Passwords never in config files.\mremotego.exe

        password: encrypted:AES256:base64encodeddata

```- ✅ Team sharing with access control



### Importing from mRemoteNG- ✅ Biometric unlock```- 🚀 **Multiple Protocols** - SSH, RDP, VNC, HTTP/HTTPS, Telnet



Use the PowerShell conversion script:- ✅ Audit logs



```powershell

.\convert-mremoteng-to-yaml.ps1 -InputXml "confCons.xml" -OutputYaml "connections.yaml"

```### 3. Plain Text (Not Recommended)



## 📚 Documentation### First Use- 📁 **Organized** - Folders, drag-and-drop, search filtering**The Solution**: MremoteGO uses clean YAML configs that work beautifully with git, plus 1Password integration for secure team password sharing.**Solution**: MremoteGO uses clean YAML configs that work beautifully with git, plus 1Password integration for secure password sharing.



- [Quick Start Guide](docs/QUICKSTART.md)```yaml

- [Password Management](docs/PASSWORD-MANAGEMENT.md)

- [1Password Setup](docs/1PASSWORD-SETUP.md)password: mypassword123

- [Encryption Guide](docs/ENCRYPTION.md)

- [GUI Guide](docs/GUI-GUIDE.md)```



## 🛠️ Development1. Launch MremoteGO- 🔒 **Auto-Login** - Windows Credential Manager for RDP, password support for SSH



### Prerequisites- ⚠️ Visible in config file



- Go 1.23 or later- ⚠️ Not safe to commit to git2. Enter a master password (optional - for encryption)

- Platform-specific requirements:

  - **Windows**: Go toolchain (no additional dependencies)

  - **Linux**: `libgl1-mesa-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libxxf86vm-dev`

  - **macOS**: Xcode command line tools## 📋 Configuration Example3. Click **[+]** to add a connection



### Building



```bash```yaml4. Fill in host details

# Development build

go build ./cmd/mremotego-guiversion: "1.0"



# Release build (Windows)connections:5. Click **[▶]** to connect## 🚀 Quick Start

go build -ldflags "-s -w" -o mremotego.exe ./cmd/mremotego-gui

  # SSH with encrypted password

# Cross-platform builds handled by GitHub Actions

```  - name: Production Web Server



### Project Structure    type: connection



```    protocol: sshThat's it! 🎉## Features### Key Features

mremotego/

├── cmd/    host: web1.prod.com

│   ├── mremotego-gui/    # GUI+CLI entry point

│   ├── mremotego/        # CLI commands    port: 22

│   └── encrypt-passwords/ # Password encryption tool

├── internal/    username: admin

│   ├── config/           # Config management

│   ├── crypto/           # Encryption/decryption    password: enc:base64encrypteddata...## 📖 Documentation### Installation

│   ├── gui/              # Fyne GUI components

│   ├── launcher/         # Protocol launchers    description: Primary web server

│   └── secrets/          # 1Password integration

├── pkg/

│   └── models/           # Data models

└── docs/                 # Documentation  # RDP with 1Password reference

```

  - name: Windows Server- [Quick Start Guide](docs/QUICKSTART.md) - Get started in 5 minutes

## 🤝 Contributing

    type: connection

Contributions are welcome! Please:

    protocol: rdp- [GUI Guide](docs/GUI-GUIDE.md) - Complete GUI reference

1. Fork the repository

2. Create a feature branch (`git checkout -b feature/amazing-feature`)    host: win.prod.com

3. Commit your changes (`git commit -m 'Add amazing feature'`)

4. Push to the branch (`git push origin feature/amazing-feature`)    port: 3389- [Encryption Guide](docs/ENCRYPTION.md) - Password encryption details```bash

5. Open a Pull Request

    username: Administrator

## 📝 License

    password: op://DevOps/windows-server/password- [1Password Setup](docs/1PASSWORD-SETUP.md) - Secure password management

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

    domain: MYDOMAIN

## 🙏 Acknowledgments

    resolution: 1920x1080- [Password Management](docs/PASSWORD-MANAGEMENT.md) - Security best practices# Clone the repository- 🎨 **Modern GUI** - Clean interface with connection tree, drag-and-drop organization- 🎨 **Modern GUI** - Clean interface with connection tree and quick actions

- Inspired by [mRemoteNG](https://mremoteng.org/)

- Built with [Fyne](https://fyne.io/) GUI framework

- 1Password integration via [1Password CLI](https://developer.1password.com/docs/cli/)

  # Organized in folders

## 📬 Support

  - name: Development

- 🐛 [Report Issues](https://github.com/jaydenthorup/mremotego/issues)

- 💬 [Discussions](https://github.com/jaydenthorup/mremotego/discussions)    type: folder## 🔐 Password Storage Optionsgit clone https://github.com/jaydenthorup/mremotego.git

- 📧 [Contact](mailto:jayden.thorup@jayfiles.com)

    children:

---

      - name: Dev Database

**Note**: This project is not affiliated with or endorsed by mRemoteNG.

        type: connection

        protocol: sshMremoteGO supports three ways to store passwords:cd mremotego- 🔐 **1Password Integration** - `op://` references keep passwords secure- 🔐 **1Password Integration** - Store passwords securely, share configs safely

        host: dev-db.local

        port: 22

        username: dbadmin

        password: op://DevOps/dev-db/password### 1. Encrypted (Recommended)

```

```yaml

See [connections.example.yaml](connections.example.yaml) for more examples.

password: enc:AaBbCcDd...  # AES-256-GCM encrypted# Build GUI (Windows - no console window)- 🔑 **RDP Auto-Login** - Windows Credential Manager for seamless connections- � **RDP Auto-Login** - Windows Credential Manager integration

## 🌐 Supported Protocols

```

| Protocol | Windows | Linux/Mac | Auto-Login |

|----------|---------|-----------|------------|- ✅ Secure at rest.\build-gui.ps1

| **SSH** | PuTTY | Terminal | ✅ Yes |

| **RDP** | mstsc | xfreerdp | ✅ Yes |- ✅ Master password required to decrypt

| **VNC** | vncviewer | vncviewer | ✅ Yes |

| **HTTP/HTTPS** | Browser | Browser | N/A |- ✅ Safe for personal use- 📝 **Git-Friendly** - YAML configs are easy to diff, merge, and review- � **Git-Friendly** - YAML configs are easy to diff and merge

| **Telnet** | telnet | telnet | ✅ Yes |



### Platform-Specific Features

### 2. 1Password Reference (Best for Teams)# Build GUI (Linux/Mac)

**Windows**:

- RDP: Stores credentials in Windows Credential Manager```yaml

- SSH: Launches PuTTY with `-pw` flag for auto-login

- GUI: No console window popupspassword: op://DevOps/server01/password./build-gui.sh- 🖥️ **Cross-Platform** - Windows, Linux, macOS- 🖥️ **Cross-Platform** - Windows, Linux, macOS



**Linux**:```

- SSH: Launches in terminal (gnome-terminal, xterm, konsole, etc.)

- Password authentication via sshpass- ✅ Passwords never in config files



**macOS**:- ✅ Team sharing with access control

- SSH: Launches in Terminal.app

- RDP: Opens Microsoft Remote Desktop via `rdp://` URL- ✅ Biometric unlock# Run- ⚡ **Fast & Clean** - No console popups, instant connections- ⚡ **Fast** - No console window popups, instant connections



## 📊 Comparison with mRemoteNG- ✅ Audit logs



| Feature | mRemoteNG | MremoteGO |.\mremotego.exe

|---------|-----------|-----------|

| Config Format | XML | ✅ YAML |### 3. Plain Text (Not Recommended)

| Git Diffs | ❌ Messy | ✅ Clean |

| Password Encryption | Per-machine DPAPI | ✅ AES-256-GCM |```yaml```- 🚀 **Protocols** - SSH (PuTTY), RDP, VNC, HTTP/HTTPS, Telnet- � **Organize** - Folders, drag-and-drop, search

| 1Password Integration | ❌ No | ✅ Yes |

| Team Sharing | ❌ Difficult | ✅ Easy |password: mypassword123

| Auto-Login | ✅ Yes | ✅ Yes |

| Cross-Platform | ❌ Windows only | ✅ All platforms |```



## 🛠️ Building from Source- ⚠️ Visible in config file



### Requirements- ⚠️ Not safe to commit to git### First Use- 📁 **Organized** - Folders, search, recent files- � **Multiple Protocols** - SSH, RDP, VNC, HTTP/HTTPS, Telnet



- Go 1.24 or higher

- Git

## 📋 Configuration Example

### Build Commands



```bash

# Windows GUI (no console)```yaml1. Launch MremoteGO

.\build-gui.ps1

version: "1.0"

# Linux/Mac GUI

./build-gui.shconnections:2. Enter a master password (optional - for encryption)



# CLI version  # SSH with encrypted password

go build -o mremotego ./cmd/mremotego

  - name: Production Web Server3. Click **[+]** to add a connection## Quick Start## Screenshots

# Encryption tool

go build -o encrypt-passwords ./cmd/encrypt-passwords    type: connection

```

    protocol: ssh4. Fill in host details

## 🔧 CLI Tool

    host: web1.prod.com

MremoteGO also includes a CLI for automation:

    port: 225. Click **[▶]** to connect

```bash

# Initialize config    username: admin

mremotego init

    password: enc:base64encrypteddata...

# List connections

mremotego list    description: Primary web server



# Add connectionThat's it! 🎉### 1. Install & Build### Main Interface

mremotego add --name "Server" --protocol ssh --host 192.168.1.100

  # RDP with 1Password reference

# Connect

mremotego connect "Server"  - name: Windows Server



# Export    type: connection

mremotego export --output backup.yaml

```    protocol: rdp## 📖 Documentation```



## 🏗️ Project Structure    host: win.prod.com



```    port: 3389

mremotego/

├── cmd/    username: Administrator

│   ├── mremotego/          # CLI application

│   ├── mremotego-gui/      # GUI application    password: op://DevOps/windows-server/password- [Quick Start Guide](docs/QUICKSTART.md) - Get started in 5 minutes```powershell┌─────────────────────────────────────────────────────────┐

│   └── encrypt-passwords/  # Password encryption tool

├── internal/    domain: MYDOMAIN

│   ├── config/             # Configuration management

│   ├── crypto/             # Encryption (AES-256-GCM)    resolution: 1920x1080- [GUI Guide](docs/GUI-GUIDE.md) - Complete GUI reference

│   ├── gui/                # Fyne GUI components

│   ├── launcher/           # Protocol launchers

│   └── secrets/            # 1Password integration

├── pkg/  # Organized in folders- [Encryption Guide](docs/ENCRYPTION.md) - Password encryption detailsgit clone https://github.com/yourusername/mremotego│ File  Connection  Help                                   │

│   └── models/             # Data models

├── docs/                   # Documentation  - name: Development

├── build-gui.ps1          # Windows build script

└── build-gui.sh           # Linux/Mac build script    type: folder- [1Password Setup](docs/1PASSWORD-SETUP.md) - Secure password management

```

    children:

## 🤝 Contributing

      - name: Dev Database- [Password Management](docs/PASSWORD-MANAGEMENT.md) - Security best practicescd mremotego│ [+] [📁] [▶] [✏️] [🗑️] [🔄]                              │

Contributions are welcome! Please:

        type: connection

1. Fork the repository

2. Create a feature branch (`git checkout -b feature/amazing-feature`)        protocol: ssh

3. Commit your changes (`git commit -m 'Add amazing feature'`)

4. Push to the branch (`git push origin feature/amazing-feature`)        host: dev-db.local

5. Open a Pull Request

        port: 22## 🔐 Password Storage Optionsgo build -ldflags "-H windowsgui" -o MremoteGO.exe ./cmd/mremotego-gui├───────────────┬─────────────────────────────────────────┤

## 📄 License

        username: dbadmin

MIT License - see [LICENSE](LICENSE) file for details.

        password: op://DevOps/dev-db/password

Copyright © 2026 [Jayden Thorup](mailto:jayden.thorup@jayfiles.com)

```

## 🙏 Acknowledgments

MremoteGO supports three ways to store passwords:.\MremoteGO.exe│ 📁 Production │ Connection Details                      │

- Inspired by [mRemoteNG](https://mremoteng.org/)

- Built with [Fyne](https://fyne.io/) GUI toolkitSee [connections.example.yaml](connections.example.yaml) for more examples.

- 1Password integration via [1Password CLI](https://developer.1password.com/docs/cli/)

- Encryption using Go's crypto libraries



## 💬 Support## 🌐 Supported Protocols



- 📖 Documentation: [docs/](docs/)### 1. Encrypted (Recommended)```│  🔐 Web1      │ 🔐 Web1                                 │

- 🐛 Issues: [GitHub Issues](https://github.com/jaydenthorup/mremotego/issues)

- 💡 Feature Requests: [GitHub Discussions](https://github.com/jaydenthorup/mremotego/discussions)| Protocol | Windows | Linux/Mac | Auto-Login |



---|----------|---------|-----------|------------|```yaml



**Made with ❤️ in Go**| **SSH** | PuTTY | Terminal | ✅ Yes |


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
