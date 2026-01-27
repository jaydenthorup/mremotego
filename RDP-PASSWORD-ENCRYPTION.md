# RDP Password Management - Technical Details

## Current Implementation

MremoteGO uses **Windows Credential Manager** (cmdkey) for RDP password management, providing seamless auto-login without storing passwords in the config file.

## How It Works

### 1. Password Storage (Windows Credential Manager)

When you connect to an RDP server with a password, MremoteGO:

1. **Stores credentials** in Windows Credential Manager using `cmdkey`:
   ```
   cmdkey /generic:TERMSRV/hostname /user:username /pass:password
   ```

2. **Launches mstsc** which automatically retrieves credentials:
   ```
   mstsc /v:hostname:port
   ```

3. **Auto-login happens** - Windows uses stored credentials automatically
   - No password prompt needed
   - Seamless login experience
   - Credentials persist until explicitly deleted

### 2. 1Password Integration

For secure team-shareable passwords, use 1Password references:

```yaml
- name: "Windows Server"
  protocol: rdp
  host: server.example.com
  username: Administrator
  password: op://Private/Windows Server/password  # 1Password reference
```

When connecting:
1. MremoteGO calls `op read op://vault/item/field` via CLI
2. 1Password authenticates using biometric unlock
3. Password is retrieved and passed to Windows Credential Manager
4. mstsc connects automatically

See [1PASSWORD-CLI-SETUP.md](1PASSWORD-CLI-SETUP.md) for setup.

### 3. Security

**Key Security Features:**

- ✅ **Windows Credential Manager** - OS-level secure storage
- ✅ **User-specific** - Credentials tied to Windows user account
- ✅ **Machine-specific** - Only accessible on the machine where stored
- ✅ **1Password integration** - Team-shareable passwords via op:// references
- ✅ **No plaintext storage** - Passwords never stored in config files
- ✅ **Biometric unlock** - 1Password uses Touch ID/Windows Hello

**Important Security Notes:**

⚠️ **Credentials are tied to your Windows user account**
- Stored in Windows Credential Manager, not config files
- If you copy the config to another machine, you can still connect (using 1Password)
- Each machine stores its own local credentials cache

⚠️ **Config file security**
- Plain text passwords: Not recommended (will be stored in Credential Manager on connect)
- 1Password references: Safe to commit to git (actual passwords stay in 1Password)

⚠️ **Best Practices**
- Use 1Password references (`op://`) for team-shared passwords
- Use SSH keys instead of passwords when possible
- Consider certificate-based RDP authentication for production
- Regular security audits of stored credentials
- Clear credential cache: `cmdkey /delete:TERMSRV/hostname`

### 4. Cross-Platform Behavior

**Windows:**
- ✅ Windows Credential Manager for RDP auto-login
- ✅ PuTTY for SSH with `-pw` flag
- ✅ 1Password CLI integration

**Linux:**
- ✅ xfreerdp with `/u:` and `/p:` flags
- ✅ Native ssh client
- ✅ 1Password CLI integration

**macOS:**
- ✅ Microsoft Remote Desktop or xfreerdp
- ✅ Native ssh client
- ✅ 1Password CLI integration

### 5. Example Usage

**GUI:**
1. Add new connection
2. Enter password or 1Password reference: `op://Private/Server/password`
3. Optional: Check "Store in 1Password" to save new passwords
4. Connect - auto-login happens automatically
**CLI:**
```bash
# Add RDP connection with 1Password reference
mremotego add --name "WinServer" \
  --protocol rdp \
  --host 192.168.1.100 \
  --username Administrator \
  --password "op://Private/WinServer/password"

# Connect (password automatically retrieved from 1Password)
mremotego connect "WinServer"
```

### 6. Credential Management Commands

**View stored credentials:**
```powershell
cmdkey /list | Select-String "TERMSRV"
```

**Manually delete credentials:**
```powershell
cmdkey /delete:TERMSRV/hostname
```

**Delete all RDP credentials:**
```powershell
cmdkey /list | Select-String "TERMSRV" | ForEach-Object { 
  $target = ($_ -split " ")[1]
  cmdkey /delete:$target
}
```

### 7. Troubleshooting

**Password not working:**
- Check 1Password CLI is signed in: `op whoami`
- Verify the op:// reference is correct
- Ensure 1Password desktop app is running
- Check "Integrate with 1Password CLI" is enabled in 1Password settings

**RDP still prompts for password:**
- Credentials may not be stored in Credential Manager yet
- Connect once to store credentials
- Check Windows Credential Manager for TERMSRV/hostname entry

**1Password authentication fails:**
- Sign in to 1Password desktop app
- Enable biometric unlock in 1Password settings
- Verify vault access permissions

## Comparison with mRemoteNG

| Feature | mRemoteNG | MremoteGO |
|---------|-----------|-----------|
| Encryption Method | XML with master password | Windows Credential Manager |
| Password Format | XML encrypted string | 1Password references or plain text |
| Storage Location | confCons.xml | Windows Credential Manager |
| RDP File Creation | ❌ Uses mstsc directly | ✅ Temporary .rdp files |
| Auto-login | ✅ | ✅ |
| Team Sharing | ❌ Difficult | ✅ 1Password integration |
| Cross-platform | ❌ Windows only | ✅ Win/Linux/Mac |

## Security Architecture

```
Config File (git-shareable)
├── Connections with op:// references
└── No plaintext passwords

Connection Flow:
1. User clicks "Connect"
2. MremoteGO detects op:// reference
3. Calls: op read op://vault/item/field
4. 1Password authenticates (biometric)
5. Password retrieved securely
6. Stored in Windows Credential Manager
7. mstsc launches with auto-login
```

---

**Security Recommendation:** Use 1Password references (`op://`) for all passwords to enable secure team sharing and keep passwords out of config files.
