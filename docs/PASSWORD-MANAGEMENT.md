# Password Management

MremoteGO provides flexible and secure password management with multiple options.

## Overview

| Method | Security | Team Sharing | Auto-Login | Best For |
|--------|----------|--------------|------------|----------|
| 1Password | ✅ High | ✅ Yes | ✅ Yes | **Recommended** - Teams |
| Bitwarden | ✅ High | ✅ Yes | ✅ Yes | **Recommended** - Teams, self-hosting |
| Plain Text | ⚠️ Low | ❌ No | ✅ Yes | Personal/testing only |
| No Password | ✅ Manual | N/A | ❌ No | SSH keys, certificates |

## 1Password Integration (Recommended)

Store passwords securely in 1Password and reference them in config files.

### Setup

See [1PASSWORD-SETUP.md](1PASSWORD-SETUP.md) for complete setup instructions.

### Usage

```yaml
connections:
  - name: "My Server"
    password: op://Private/My Server/password  # Secure reference
```

### Benefits
- ✅ Passwords never stored in config files
- ✅ Safe to commit configs to git
- ✅ Team password sharing
- ✅ Biometric unlock
- ✅ Automatic password rotation support
- ✅ Audit logs

## Bitwarden Integration (Recommended)

Store passwords in Bitwarden, bitwarden.com or self-hosted, and reference them
in config files.

### Setup

See [BITWARDEN-SETUP.md](BITWARDEN-SETUP.md) for complete setup instructions.

### Usage

```yaml
connections:
  - name: "My Server"
    password: bw://8f3c1d9a-4e2b-4c77-9f10-1a2b3c4d5e6f  # Secure reference
```

References use the item id, so renaming an item does not break the config. A
field can be selected explicitly with `bw://<item-id>/username`, `/totp` or
`/notes`.

### Benefits
- ✅ Passwords never stored in config files
- ✅ Safe to commit configs to git
- ✅ Team password sharing
- ✅ Works with self-hosted Bitwarden and Vaultwarden
- ✅ Item picker built into the connection dialog
- ✅ Master password never handled by MremoteGO

## Plain Text Passwords

For personal use or testing environments:

```yaml
connections:
  - name: "Dev Server"
    password: mypassword123  # Plain text (not recommended for production)
```

### Security Considerations
- ⚠️ Readable in config file
- ⚠️ Should not be committed to git
- ⚠️ No audit trail
- ✅ Still works with auto-login features

### Best Practices
- Use `.gitignore` to exclude config files
- Set file permissions: `icacls %APPDATA%\mremotego\config.yaml /inheritance:r /grant:r "%USERNAME%:F"`
- Consider encryption at rest (BitLocker, FileVault)

## RDP Auto-Login

MremoteGO uses **Windows Credential Manager** for seamless RDP connections.

### How It Works

1. **First Connection**: Password (from a password manager or plain text) is
   stored in Windows Credential Manager as a generic `TERMSRV/hostname`
   credential, written through the Credential Manager API rather than the
   `cmdkey` command, so the password never appears on a command line.

2. **Subsequent Connections**: Windows automatically retrieves credentials
   ```
   mstsc /v:hostname
   ```

3. **Auto-Login**: No password prompt needed

### Benefits
- ✅ Native Windows integration
- ✅ Persistent across sessions
- ✅ User-specific security
- ✅ Works with password manager references
- ✅ No passwords in temporary files

### Managing Stored Credentials

**View all RDP credentials:**
```powershell
cmdkey /list | Select-String "TERMSRV"
```

**Delete specific credential:**
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

## SSH Password Handling

### Windows (PuTTY)
MremoteGO uses PuTTY with password auto-fill. The password is written to a
private temporary file and passed with `-pwfile`, which PuTTY reads at start-up:
```
putty.exe -ssh -P 22 -l username -pwfile <temp file> hostname
```
The file is deleted as soon as PuTTY is running. PuTTY's `-pw` option is not
used because it puts the password on the command line, where any local process
can read it.

### Linux/Mac (Native SSH)
Uses native ssh client:
```
ssh username@hostname -p 22
```
When a password is set and `sshpass` is installed, the password is read from a
private temporary file into the `SSHPASS` environment variable and the file is
removed before `ssh` starts, so the password appears neither in the process
list nor on disk for longer than necessary.

### SSH Key Authentication (Recommended)
For better security, use SSH keys instead of passwords:

1. Generate key pair: `ssh-keygen -t ed25519`
2. Copy to server: `ssh-copy-id user@host`
3. Leave password field empty in MremoteGO
4. SSH will use key authentication automatically

## VNC Connections

VNC passwords are passed to vncviewer:
```
vncviewer hostname:port -password password
```

## Security Comparison

### 1Password References
```yaml
password: op://Private/Server/password
```
- ✅ Config file safe to commit to git
- ✅ Passwords centrally managed in 1Password
- ✅ Team sharing with access control
- ✅ Biometric unlock
- ✅ Audit logs
- ✅ Automatic rotation support

### Plain Text Passwords
```yaml
password: mypassword123
```
- ⚠️ Visible in config file
- ⚠️ Should not be in version control
- ⚠️ No audit trail
- ⚠️ Manual password rotation
- ✅ Simple for personal use
- ✅ Works offline

### Windows Credential Manager (RDP)
- ✅ OS-level secure storage
- ✅ User and machine specific
- ✅ Protected by Windows login
- ✅ Integrated with DPAPI
- ⚠️ Local to machine (not synced)

## Best Practices

### For Teams
1. ✅ Use 1Password or Bitwarden for all passwords
2. ✅ Store configs in git with `op://` or `bw://` references
3. ✅ Use shared vaults for team credentials
4. ✅ Enable biometric unlock
5. ✅ Regular access audits

### For Personal Use
1. ✅ Use 1Password or Bitwarden if you have one
2. ⚠️ Plain text is acceptable for local dev
3. ✅ Use SSH keys where possible
4. ✅ Keep config file permissions restricted
5. ✅ Don't commit passwords to public repos

### For Production
1. ✅ Use 1Password, Bitwarden or another enterprise password manager
2. ✅ Certificate-based authentication where possible
3. ✅ SSH keys instead of passwords
4. ✅ Regular credential rotation
5. ✅ Audit all password access
6. ✅ Multi-factor authentication

## Troubleshooting

### RDP asks for password despite stored credentials
- Check Windows Credential Manager: `cmdkey /list`
- Delete and reconnect to refresh: `cmdkey /delete:TERMSRV/hostname`
- Verify username format (use `DOMAIN\username` if needed)

### 1Password reference not working
- Verify 1Password desktop app is running and unlocked
- Check CLI integration is enabled: Settings → Developer
- Test reference: `op read "op://Private/Server/password"`
- Verify vault name and item name (case-sensitive)

### SSH password not auto-filling
- Windows: Verify PuTTY is installed and in PATH
- Check password field is not empty
- For 1Password: Verify reference is correct
- Consider using SSH keys instead

### Password visible in process list
- MremoteGO no longer passes passwords on a command line
- SSH uses PuTTY `-pwfile` on Windows and `SSHPASS` on Linux/Mac
- RDP uses the Windows Credential Manager API
- Known limitation: on Linux, `xfreerdp` is still invoked with `/p:`
