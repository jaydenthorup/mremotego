# Quick Start Guide

This guide will help you get started with MremoteGO in just a few minutes.

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/mremotego
cd mremotego

# Build the application
go build -o mremotego cmd/mremotego/main.go

# (Optional) Move to a directory in your PATH
# Linux/Mac:
sudo mv mremotego /usr/local/bin/
# Windows: Move to a directory in your PATH or use full path
```

### Using Go Install

```bash
go install github.com/yourusername/mremotego@latest
```

## First Steps

### 1. Initialize Configuration

Create a new configuration file with example connections:

```bash
mremotego init
```

This creates a configuration file at:
- **Windows**: `%APPDATA%\mremotego\config.yaml`
- **Linux/Mac**: `~/.config/mremotego/config.yaml`

### 2. List Connections

View all your connections:

```bash
mremotego list
```

Output:
```
📁 Examples
  🔐 Example SSH (ssh://example.com:22)
     └─ Example SSH connection
  🖥️ Example RDP (rdp://windows-server.local:3389)
     └─ Example RDP connection
```

### 3. Add a New Connection

Add an SSH connection:

```bash
mremotego add --name "My Server" \
  --protocol ssh \
  --host 192.168.1.100 \
  --port 22 \
  --username admin \
  --description "Production web server"
```

Add an RDP connection:

```bash
mremotego add --name "Windows Dev" \
  --protocol rdp \
  --host dev.windows.local \
  --username Developer \
  --folder "Development"
```

### 4. Connect to a Host

Launch a connection:

```bash
mremotego connect "My Server"
```

This will:
- **SSH**: Launch your SSH client with the configured settings
- **RDP**: Launch mstsc (Windows) or xfreerdp (Linux/Mac)
- **HTTP/HTTPS**: Open in your default browser
- **VNC**: Launch your VNC viewer

### 5. Edit a Connection

Update connection details:

```bash
mremotego edit "My Server" --host 192.168.1.101 --port 2222
```

### 6. Delete a Connection

Remove a connection:

```bash
mremotego delete "My Server"
```

## Organizing Connections

### Using Folders

Organize connections in folders using the `--folder` flag:

```bash
mremotego add --name "Web1" --protocol ssh --host web1.prod.com --folder "Production/Web"
mremotego add --name "Web2" --protocol ssh --host web2.prod.com --folder "Production/Web"
mremotego add --name "DB1" --protocol ssh --host db1.prod.com --folder "Production/Database"
```

Your structure will look like:

```
📁 Production
  📁 Web
    🔐 Web1 (ssh://web1.prod.com:22)
    🔐 Web2 (ssh://web2.prod.com:22)
  📁 Database
    🔐 DB1 (ssh://db1.prod.com:22)
```

### Using Tags

Add tags for better organization:

```bash
mremotego add --name "Server1" \
  --protocol ssh \
  --host server1.com \
  --tags "production,web,critical"
```

## Advanced Usage

### Custom Config Location

Use a different config file:

```bash
mremotego --config /path/to/custom-config.yaml list
```

### Protocol-Specific Options

#### SSH with Custom Port
```bash
mremotego add --name "SSH Custom" --protocol ssh --host example.com --port 2222 --username deploy
```

#### RDP with Domain
```bash
mremotego add --name "Corp Server" --protocol rdp --host server.corp.local --username john --domain CORP
```

#### HTTP/HTTPS
```bash
mremotego add --name "Admin Panel" --protocol https --host admin.example.com --port 8443
```

### Export Configuration

Backup your configuration:

```bash
mremotego export --output backup-$(date +%Y%m%d).yaml
```

## Git Integration

### Initialize Git Repository

```bash
cd ~/.config/mremotego  # Linux/Mac
# or
cd %APPDATA%\mremotego  # Windows

git init
git add config.yaml
git commit -m "Initial commit"
```

### Track Changes

```bash
# After adding/editing connections
git diff config.yaml

# Commit changes
git add config.yaml
git commit -m "Added production servers"
```

### Share with Team

```bash
# Add remote repository
git remote add origin git@github.com:yourorg/connections.git
git push -u origin main

# Team members can clone
git clone git@github.com:yourorg/connections.git ~/.config/mremotego
```

## Example Workflows

### Quick SSH Connection

```bash
# Add
mremotego add --name prod-web --protocol ssh --host 10.0.1.50 --username deploy

# Connect
mremotego connect prod-web
```

### Managing Multiple Environments

```bash
# Development
mremotego add --name dev-api --protocol ssh --host dev.api.local --folder "Dev"
mremotego add --name dev-db --protocol ssh --host dev.db.local --folder "Dev"

# Production
mremotego add --name prod-api --protocol ssh --host prod.api.com --folder "Production"
mremotego add --name prod-db --protocol ssh --host prod.db.com --folder "Production"

# List organized by folders
mremotego list
```

### Team Onboarding

```bash
# Initialize with team config
mremotego --config team-connections.yaml init

# Edit with team connections
# Then commit to shared repository
git add team-connections.yaml
git commit -m "Team connection template"
git push
```

## Tips and Tricks

1. **Use Tab Completion**: Generate shell completion with `mremotego completion bash/zsh/fish/powershell`

2. **Alias Common Commands**: Add to your shell profile:
   ```bash
   alias mrc='mremotego connect'
   alias mrl='mremotego list'
   alias mra='mremotego add'
   ```

3. **Keep Sensitive Data Separate**: Don't commit passwords. Use SSH keys or prompt for passwords.

4. **Regular Backups**: Set up automated exports:
   ```bash
   # Add to crontab
   0 0 * * * mremotego export --output ~/backups/mremotego-$(date +\%Y\%m\%d).yaml
   ```

5. **Version Control Best Practices**:
   - Use meaningful commit messages
   - Review diffs before committing
   - Use branches for testing new configurations

## Troubleshooting

### Connection Not Found
```bash
# List all connections to check the exact name
mremotego list
```

### Config File Not Found
```bash
# Initialize a new config
mremotego init

# Or specify custom location
mremotego --config /path/to/config.yaml list
```

### Protocol Client Not Found

Make sure the required client is installed:
- **SSH**: OpenSSH client (usually pre-installed)
- **RDP**: mstsc (Windows) or xfreerdp (Linux/Mac)
- **VNC**: vncviewer or any VNC client
- **Telnet**: telnet client

## Next Steps

- Check out the [README.md](README.md) for complete documentation
- View example configurations in [config.example.yaml](config.example.yaml)
- Report issues on GitHub
- Contribute improvements via pull requests

## Support

For help and support:
- GitHub Issues: https://github.com/yourusername/mremotego/issues
- Documentation: https://github.com/yourusername/mremotego
