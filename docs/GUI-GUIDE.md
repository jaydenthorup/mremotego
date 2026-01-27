# GUI Guide

Complete guide to the MremoteGO graphical interface.

## Overview

MremoteGO provides a clean, modern GUI similar to mRemoteNG but with better usability.

### Main Window

```
┌─────────────────────────────────────────────────────────┐
│ File  Connection  Help                                   │
├───────────────────────────────────────────────────────────┤
│ [+] [📁] [▶] [✏️] [🗑️] [🔄]                              │
├───────────────┬─────────────────────────────────────────┤
│ 📁 Production │ Connection Details                      │
│  🔐 Web1      │ 🔐 Web1                                 │
│  🔐 DB1       │ Protocol: SSH                           │
│  🔐 RDP1      │ Host: web1.prod.com                     │
│ 📁 Development│ Port: 22                                │
│  🔐 DevBox    │ Username: admin                         │
│  🔐 TestDB    │ Password: op://Private/Web1/password    │
└───────────────┴─────────────────────────────────────────┘
```

## Toolbar Buttons

| Button | Action | Keyboard |
|--------|--------|----------|
| **[+]** Add | Add new connection | Ctrl+N |
| **[📁]** Folder | Create new folder | Ctrl+Shift+N |
| **[▶]** Connect | Launch connection | Enter |
| **[✏️]** Edit | Edit connection | F2 |
| **[🗑️]** Delete | Delete connection/folder | Delete |
| **[🔄]** Refresh | Reload config | F5 |

## Menu Bar

### File Menu

**Open Config...** (Ctrl+O)
- Load a different configuration file
- Supports `.yaml` and `.yml` files
- Last opened file is remembered

**Save** (Ctrl+S)
- Saves current configuration
- Auto-saves on connection changes

**Exit** (Alt+F4)
- Close application
- Config automatically saved

### Connection Menu

**Add Connection** (Ctrl+N)
- Opens add connection dialog

**Edit Connection** (F2)
- Edit selected connection

**Delete** (Delete)
- Removes connection or folder
- Prompts for confirmation

**Connect** (Enter)
- Launches selected connection

### Help Menu

**About**
- Version information
- Build details

## Connection Tree

### Organization

**Folders**
- Right-click tree → **Add Folder**
- Nest folders for organization
- Drag connections to folders

**Drag & Drop**
- Reorganize connections
- Move between folders
- Visual feedback during drag

**Icons**
- 📁 Folder
- 🔐 Connection
- Protocol-specific icons (future)

### Context Menu

Right-click on connection:
- **Connect** - Launch connection
- **Edit** - Edit details
- **Duplicate** - Copy connection
- **Delete** - Remove connection

Right-click on folder:
- **Add Connection** - Add to folder
- **Add Folder** - Create subfolder
- **Rename** - Change folder name
- **Delete** - Remove folder

## Add/Edit Connection Dialog

### Basic Tab

**Required Fields:**
- **Name**: Display name in tree
- **Protocol**: SSH, RDP, VNC, HTTP, HTTPS, Telnet
- **Host**: Hostname or IP address

**Optional Fields:**
- **Port**: Default per protocol if empty
- **Username**: Login username
- **Password**: Plain text or `op://` reference
- **Description**: Notes about connection

### 1Password Integration

**Password Field Options:**

1. **Plain Text**
   ```
   mypassword123
   ```

2. **1Password Reference**
   ```
   op://Private/Server Name/password
   ```

3. **Store in 1Password** (Checkbox)
   - Creates new 1Password item automatically
   - Select vault from dropdown
   - Item named after connection name

### Advanced Tab (Future)

- Extra arguments per protocol
- Custom color coding
- Connection-specific settings
- Terminal preferences

## Details Panel

Shows information about selected connection:

```
Connection: Web Server 1
Protocol:   SSH
Host:       web1.prod.com
Port:       22
Username:   admin
Password:   ••••••••••••••••
```

Click **Connect** button to launch.

## Keyboard Shortcuts

### Global

| Shortcut | Action |
|----------|--------|
| Ctrl+N | Add Connection |
| Ctrl+Shift+N | Add Folder |
| Ctrl+O | Open Config |
| Ctrl+S | Save Config |
| F5 | Refresh Tree |
| Ctrl+Q | Quit |

### Tree Navigation

| Shortcut | Action |
|----------|--------|
| ↑↓ | Navigate connections |
| Enter | Connect to selected |
| F2 | Edit selected |
| Delete | Delete selected |
| Ctrl+C | Copy connection |
| Ctrl+V | Paste connection |

### Connection

| Shortcut | Action |
|----------|--------|
| Enter | Connect |
| Ctrl+Enter | Connect in new session |
| Esc | Cancel dialog |

## Configuration File

### Location

**Windows**: `%APPDATA%\mremotego\config.yaml`

**Linux/Mac**: `~/.config/mremotego/config.yaml`

### Recent Files

Last opened file tracked at: `%APPDATA%\mremotego\recent.txt`

GUI automatically opens last file on startup.

### Format

```yaml
version: "1.0"
connections:
  - name: "Production"
    type: folder
    children:
      - name: "Web Server"
        type: connection
        protocol: ssh
        host: web.prod.com
        port: 22
        username: admin
        password: op://Private/Web/password
```

## Tips & Tricks

### Quick Connect

Double-click a connection to launch immediately.

### Multi-Config Workflow

```
Work:     C:\work\production.yaml
Personal: C:\users\me\Documents\home-servers.yaml
Testing:  C:\dev\test-connections.yaml
```

Use **File → Open Config** to switch between them.

### Search (Future Feature)

Press **Ctrl+F** to search connections by name, host, or description.

### Color Coding (Future Feature)

Assign colors to connections:
- 🔴 Production (be careful!)
- 🟡 Staging
- 🟢 Development

### Templates (Future Feature)

Save connection templates for common setups:
- Standard SSH Server
- Windows RDP Server
- VNC Desktop

## Troubleshooting

### Tree not refreshing after config change
Press **F5** to manually refresh.

### Can't drag and drop
Make sure you're not in edit mode. Save or cancel any open dialogs.

### Wrong connection selected
Click on connection name in tree to select it properly.

### Password field shows reference instead of dots
This is intentional for `op://` references - so you can see what's configured.

### Application icon not showing
Icon is embedded in the executable. May not show in some environments.

## Customization

### Window Size

Window size and position are remembered between sessions.

### Theme (Future)

- Light/Dark mode toggle
- Custom color schemes
- Accessibility options

## Performance

### Large Configuration Files

MremoteGO handles thousands of connections efficiently:
- Lazy loading of tree nodes
- Efficient YAML parsing
- Fast search indexing

### Startup Time

First launch may be slower due to:
- Config file loading
- 1Password CLI availability check
- Window restoration

Subsequent launches are faster.

## See Also

- [Quick Start](QUICKSTART.md) - Getting started guide
- [1Password Setup](1PASSWORD-SETUP.md) - Password management
- [Password Management](PASSWORD-MANAGEMENT.md) - Security details
- [README](../README.md) - Full documentation
