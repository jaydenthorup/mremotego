# 1Password CLI Integration

MremoteGO supports storing passwords in 1Password and referencing them in your config file. This is perfect for team environments where you want to share connection configurations via git without exposing passwords.

## Prerequisites

1. Install [1Password CLI](https://developer.1password.com/docs/cli/get-started/)
2. Sign in to 1Password CLI: `op signin`
3. **Enable biometric unlock** (recommended for persistent sessions)

## Setup for Persistent Sessions

To avoid having to authenticate repeatedly, configure 1Password CLI for persistent sessions:

### Windows (Recommended Setup)

1. **Sign in with biometric unlock enabled:**
   ```powershell
   op signin
   ```

2. **Turn on 1Password app integration:**
   - Open 1Password app
   - Go to Settings → Developer
   - Enable "Connect with 1Password CLI"
   - Enable "Integrate with 1Password app"

3. **Verify integration:**
   ```powershell
   op account list
   op vault list
   ```

With this setup, 1Password CLI will use your existing 1Password app session and biometric unlock. No need to enter passwords repeatedly!

### Session Duration

- **With app integration**: Session lasts as long as your 1Password app is unlocked
- **Without app integration**: Session expires after inactivity (default: 30 minutes)
- **Biometric unlock**: Touch ID/Windows Hello unlocks CLI automatically

## Usage

### Storing Passwords in 1Password

Instead of entering a password directly, use a 1Password reference in this format:

```
op://vault-name/item-name/field-name
```

**Examples:**
- `op://Private/MyServer/password` - Gets the password field from "MyServer" item in "Private" vault
- `op://Work/Production-DB/admin-password` - Gets custom "admin-password" field
- `op://Shared/RDP-Servers/hwb-10-password` - Team-shared password from "Shared" vault

### In the GUI

When adding or editing a connection:
1. In the **Password** field, enter a 1Password reference like: `op://Private/MyServer/password`
2. Click **Save**
3. The reference is stored as-is in the YAML config (not encrypted)
4. When connecting, the password is automatically retrieved from 1Password

### In the Config File

You can manually edit your `config.yaml` to use 1Password references:

```yaml
connections:
  - name: Production Server
    protocol: rdp
    host: prod.example.com
    username: administrator
    password: op://Work/ProdServer/password  # 1Password reference
    port: 3389
```

### In the CLI

```bash
# Add a connection with 1Password reference
mremotego add "My Server" rdp server.example.com --username admin --password "op://Private/MyServer/password"

# Connect (password retrieved automatically)
mremotego connect "My Server"
```

## How It Works

1. **When saving**: 1Password references (starting with `op://`) are stored as-is in the YAML file
2. **When loading**: If 1Password CLI is available, references are automatically resolved using `op read`
3. **Fallback**: If `op` CLI is not available or authentication fails, the connection will prompt for password

## Benefits

✅ **Team Sharing**: Share connection configs via git without exposing passwords
✅ **Centralized Secrets**: All passwords stored securely in 1Password
✅ **Audit Trail**: 1Password tracks who accessed which secrets
✅ **Cross-Platform**: Works on Windows, Mac, and Linux
✅ **Git-Friendly**: No merge conflicts from encrypted passwords

## Password Storage Priority

MremoteGO supports three password storage methods (in order of priority):

1. **1Password References** (`op://...`) - Resolved from 1Password CLI
2. **DPAPI Encrypted** (Windows) - Encrypted using Windows Data Protection API
3. **Plain Text** (Fallback) - Stored as-is (not recommended for production)

## Example Workflow

### For Team Leads (Setting Up)

1. Create 1Password vault for your team (e.g., "DevOps-Servers")
2. Add server credentials as items in 1Password
3. Edit `config.yaml` to use 1Password references:
   ```yaml
   connections:
     - name: Production DB
       protocol: rdp
       host: db.prod.example.com
       username: dbadmin
       password: op://DevOps-Servers/ProdDB/password
   ```
4. Commit `config.yaml` to git (safe to commit - no actual passwords!)
5. Share vault access with team members

### For Team Members (Using)

1. Install 1Password CLI: `winget install 1Password.CLI`
2. Sign in: `op signin`
3. Clone the repo and run mremotego
4. Passwords are automatically retrieved from your 1Password account
5. Connect to servers without ever seeing the actual passwords!

## Troubleshooting

### "1Password CLI is not available"
- Install 1Password CLI: https://developer.1password.com/docs/cli/get-started/
- Verify installation: `op --version`

### "Failed to retrieve secret from 1Password"
- **Check authentication**: `op account list`
- **Sign in if needed**: `op signin`
- **Check vault access**: `op vault list`
- **Verify reference format**: `op://vault/item/field`
- **Test manually**: `op read "op://vault/item/field"`

### Session Keeps Expiring

**Best solution - Enable app integration:**
1. Open 1Password app → Settings → Developer
2. Enable "Connect with 1Password CLI"
3. Enable "Integrate with 1Password app"
4. Restart your terminal
5. CLI will now use your app session (stays signed in)

**Alternative - Use environment variable:**
```powershell
# Sign in and export session token
$env:OP_SESSION_my = $(op signin --raw)

# Or sign in normally (will prompt for password)
op signin
```

**Verify persistent session:**
```powershell
# Should show your accounts without prompting
op account list

# Should work without password prompt
op read "op://Private/test/password"
```

### Password prompt appears even with 1Password reference
- Check that the reference is correctly formatted in config
- Ensure 1Password CLI is authenticated: `op account get`
- Check that the item exists: `op item get "item-name" --vault "vault-name"`
- Look for error messages in the console output

### Biometric Unlock Not Working
- Ensure 1Password app is installed and unlocked
- Check Settings → Developer → "Use biometric unlock for 1Password CLI"
- Restart terminal after enabling
- On Windows: Ensure Windows Hello is configured

## Security Notes

⚠️ **1Password CLI Session**: The `op` CLI uses session tokens that expire. You may need to re-authenticate periodically using `op signin`.

✅ **Biometric Unlock**: 1Password CLI supports biometric authentication on supported devices.

✅ **Secret References**: 1Password references in the YAML file are safe to commit to git - they're just pointers, not actual passwords.

## Mixing Storage Methods

You can mix storage methods in the same config file:

```yaml
connections:
  # 1Password reference (recommended for teams)
  - name: Production Server
    password: op://Work/ProdServer/password
  
  # DPAPI encrypted (user-specific)
  - name: Personal Server
    password: AQAAANCMnd8BFdERjHoAwE...  # Long encrypted string
  
  # Plain text (legacy, not recommended)
  - name: Test Server
    password: test123
```

## Additional Resources

- [1Password CLI Documentation](https://developer.1password.com/docs/cli/)
- [1Password Secret References](https://developer.1password.com/docs/cli/secret-references/)
- [Managing Team Vaults](https://support.1password.com/create-share-vaults/)
