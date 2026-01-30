# Quick Start Guide

Get started with MremoteGO in 5 minutes!

## Installation

### Windows

1. **Download or Build**
   ```powershell
   # Clone repository
   git clone https://github.com/jaydenthorup/mremotego
   cd mremotego
   
   # Build GUI (no console window)
   go build -ldflags "-H windowsgui" -o MremoteGO.exe ./cmd/mremotego-gui
   ```

2. **Run**
   ```powershell
   .\MremoteGO.exe
   ```

The app automatically creates a config file at `%APPDATA%\mremotego\config.yaml`

## First Connection

### 1. Add a Connection

Click **[+] Add** button and fill in:

```
Name:       My Server
Protocol:   SSH  (or RDP, VNC, HTTP, Telnet)
Host:       192.168.1.100
Port:       22
Username:   admin
Password:   your_password
```

Click **Submit**.

### 2. Connect

1. Select the connection in the tree
2. Click **[▶] Connect**
3. Auto-login happens automatically!

## Using 1Password (Recommended)

### Option 1: Desktop App Integration (Best Experience)

1. Install [1Password BETA desktop app](https://releases.1password.com/)
2. Enable SDK: Settings → Developer → Enable both SDK options
3. Configure in your `config.yaml`:
   ```yaml
   settings:
     onePasswordAccount: "your-account-name"
     vaultNameMappings:
       Personal: "vault-uuid-here"
       Work: "vault-uuid-here"
   ```
4. Get biometric authentication automatically!

### Option 2: CLI Fallback

1. Install [1Password CLI](https://developer.1password.com/docs/cli/get-started/)
2. Sign in: `op signin`
3. Launch MremoteGO from same terminal - automatic fallback!

### Add Connection with 1Password

When adding a connection:

**With vault mappings:**
```
password: op://Personal/My Server/password
```

**Or check "Store in 1Password"** to create the item automatically.

See [1PASSWORD-SETUP.md](1PASSWORD-SETUP.md) for complete details.

## Common Tasks

### Organize Connections

- **Create Folder**: Right-click tree → Add Folder
- **Drag & Drop**: Organize connections by dragging
- **Edit**: Double-click or right-click → Edit

### Open Different Config Files

- **File → Open Config...** to load another YAML file
- Perfect for work/personal separation
- Last opened file is remembered

### Share Configs with Team

1. Use 1Password for passwords
2. Config file looks like:
   ```yaml
   connections:
     - name: "Production"
       password: op://Shared/Production/password
   ```
3. Commit to git - safe to share!
4. Team members use their own 1Password for access

## Keyboard Shortcuts

| Action | Shortcut |
|--------|----------|
| Add Connection | Ctrl+N |
| Connect | Enter |
| Delete | Delete |
| Refresh | F5 |

## Protocols

### SSH
- **Windows**: Uses PuTTY (install separately)
- **Mac/Linux**: Native ssh
- **Auto-login**: Password passed with `-pw` flag

### RDP
- **Windows**: mstsc with Credential Manager
- **Linux**: xfreerdp
- **Auto-login**: ✅ Seamless

### VNC
- Uses vncviewer
- Install TigerVNC or similar

### HTTP/HTTPS
- Opens in default browser
- Good for web management panels

### Telnet
- Legacy protocol support
- Uses native telnet client

## Tips

### SSH Keys (Better than Passwords)
1. Generate: `ssh-keygen -t ed25519`
2. Copy to server: `ssh-copy-id user@host`
3. Leave password field empty in MremoteGO

### Separate Environments
```
C:\work\
  ├── production-config.yaml
  └── staging-config.yaml

C:\personal\
  └── home-servers.yaml
```

Load with: File → Open Config...

### Backup Your Config
Config location: `%APPDATA%\mremotego\config.yaml`

Backup options:
- Copy to cloud storage
- Commit to private git repo (with 1Password references)
- Regular file system backups

## Troubleshooting

### "PuTTY not found"
Install PuTTY: `winget install PuTTY.PuTTY`

### "1Password CLI not available"
Install CLI: `winget install 1Password.CLI`

### RDP asks for password
First connection stores credentials. Should auto-login after that.

### Config file not found
Default: `%APPDATA%\mremotego\config.yaml`
Create manually or use File → Open Config to specify location.

## Next Steps

- Read [README.md](../README.md) for full documentation
- Set up [1Password integration](docs/1PASSWORD-SETUP.md)
- Learn about [password management](docs/PASSWORD-MANAGEMENT.md)
- Explore [GUI features](docs/GUI-GUIDE.md)

## Getting Help

- Check documentation in `docs/` folder
- View example config: `config.example.yaml`
- Report issues: GitHub issues page
