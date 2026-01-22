# RDP Password Encryption - Technical Details

## Current Status

**⚠️ KNOWN ISSUE**: RDP automatic password login is currently disabled due to compatibility issues with mstsc password format.

- RDP connections will **prompt for password** at connection time
- Password is securely stored in YAML using DPAPI encryption
- Technical issue: mstsc requires specific password encryption format that needs further investigation
- See issue tracker for updates on automatic password login feature

## How It Was Designed to Work

MremoteGO was designed to support **automatic RDP password login** using Windows Data Protection API (DPAPI), similar to mRemoteNG.

## Implementation Details

### 1. Password Storage (Windows DPAPI)

When you save an RDP connection with a password, MremoteGO:

1. **Encrypts the password** using Windows `CryptProtectData` API
   - Uses the current user's credentials as the encryption key
   - Password can only be decrypted by the same user on the same machine
   - No encryption keys stored in files

2. **Stores encrypted password** in the YAML config:
   ```yaml
   - name: "Windows Server"
     protocol: rdp
     host: server.example.com
     username: Administrator
     password: "base64encodedencryptedpassword"  # DPAPI encrypted
   ```

3. **Creates temporary .rdp file** when connecting (currently without password):
   ```ini
   full address:s:server.example.com:3389
   username:s:Administrator
   password 51:b:base64encryptedpassword
   ```

4. **Launches mstsc** with the .rdp file
   - Windows automatically decrypts and uses the password
   - No password prompt needed
   - Seamless login experience

### 2. Security

**Key Security Features:**

- ✅ **User-specific encryption** - Passwords encrypted per Windows user account
- ✅ **Machine-specific** - Encrypted passwords only work on the machine where they were created
- ✅ **No master password** - Uses Windows login credentials automatically
- ✅ **Same as mRemoteNG** - Industry-standard DPAPI encryption
- ✅ **Temporary files** - .rdp files created in `%TEMP%\mremotego\` and can be cleaned up

**Important Security Notes:**

⚠️ **Passwords are tied to your Windows user account**
- If you copy the config to another machine, passwords won't decrypt
- If you copy the config to another user account, passwords won't decrypt
- This is by design for security

⚠️ **Config file is still readable**
- While passwords are encrypted, the base64 string is visible in YAML
- Set appropriate file permissions: `icacls %APPDATA%\mremotego\config.yaml /inheritance:r /grant:r "%USERNAME%:F"`

⚠️ **Best Practices**
- Use SSH keys instead of passwords when possible
- Consider certificate-based RDP authentication for production
- Store sensitive configs in encrypted filesystems or use BitLocker
- Regular security audits of stored credentials

### 3. Cross-Platform Behavior

**Windows:**
- ✅ Full DPAPI encryption support
- ✅ Automatic password login
- ✅ Creates .rdp files with encrypted passwords

**Linux/macOS:**
- ⚠️ Passwords stored in plain YAML (xfreerdp limitation)
- ✅ Passed directly to xfreerdp command line
- 💡 Recommendation: Use SSH keys or certificate auth

### 4. Example Usage

**CLI:**
```bash
# Add RDP connection with password
mremotego add --name "WinServer" \
  --protocol rdp \
  --host 192.168.1.100 \
  --username Administrator \
  --password "MySecurePassword" \
  --domain CORP

# Connect (password automatically used)
mremotego connect "WinServer"
```

**GUI:**
1. Click **[+]** Add Connection
2. Fill in details including password
3. Click **Submit** - password is automatically encrypted
4. Click **Connect** - automatic login, no password prompt

### 5. Config File Example

**Before encryption (internal, not visible):**
```yaml
password: "MySecurePassword"
```

**After encryption (what's stored):**
```yaml
password: "AQAAANCMnd8BFdERjHoAwE/Cl+sBAAAA5xK7V9..."
```

The encrypted string is only decryptable by the same Windows user on the same machine.

### 6. Troubleshooting

**Password not working:**
- Config file was copied from another machine → Re-enter passwords
- Running as different Windows user → Re-enter passwords
- Windows profile was recreated → Re-enter passwords

**Password prompts still appear:**
- Windows may show security prompt first time
- Check Windows Credential Manager
- Verify .rdp file created in `%TEMP%\mremotego\`

**Security warnings:**
- Normal for untrusted certificates
- Can be disabled in connection settings (not recommended for production)

### 7. API Reference

**Encryption Functions (Windows only):**

```go
// Encrypt password using Windows DPAPI
encrypted, err := encryptPasswordWindows(password)

// Decrypt password using Windows DPAPI
decrypted, err := decryptPasswordWindows(encrypted)
```

**RDP File Creation:**
```go
// Creates temporary .rdp file with encrypted password
rdpFile, err := createRDPFile(connection, target)
```

## Comparison with mRemoteNG

| Feature | mRemoteNG | MremoteGO |
|---------|-----------|-----------|
| Encryption Method | Windows DPAPI | Windows DPAPI ✅ |
| Password Format | XML encrypted string | YAML base64 encrypted |
| Storage Location | confCons.xml | config.yaml |
| RDP File Creation | ✅ Temporary | ✅ Temporary |
| Auto-login | ✅ | ✅ |
| Cross-platform | ❌ Windows only | ✅ Win/Linux/Mac |

## Future Enhancements

Planned improvements:
- [ ] Option to use Windows Credential Manager
- [ ] Master password for additional encryption layer
- [ ] SSH key management integration
- [ ] Certificate-based authentication
- [ ] Cleanup old .rdp files automatically
- [ ] Audit log for password usage

---

**Security Recommendation:** For production environments, consider using certificate-based or Kerberos authentication instead of passwords where possible.
