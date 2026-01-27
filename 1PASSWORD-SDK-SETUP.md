# 1Password CLI Setup

MremoteGO uses the **1Password CLI** to securely retrieve passwords from your 1Password vaults.

## Setup Requirements

### Install 1Password CLI

1. Install [1Password desktop app](https://1password.com/downloads) if you haven't already
2. Install [1Password CLI](https://developer.1password.com/docs/cli/get-started/)
   - Windows: Download from the link above or use `winget install 1Password.CLI`
3. Make sure you're signed in to the 1Password desktop app
4. Enable CLI integration in 1Password:
   - Open 1Password → Settings → Developer
   - Enable "Integrate with 1Password CLI"

### First Time Setup

The first time you use an `op://` reference in MremoteGO, 1Password will:
1. Prompt you to authorize the CLI integration
2. Request biometric authentication (Touch ID, Face ID, or Windows Hello)
3. Cache your session for future use

**That's it!** No tokens needed for personal use.

## Usage

Use `op://` reference format in your connection passwords:

```yaml
connections:
  - name: Production Server
    host: prod.example.com
    username: admin
    password: op://Private/Production Server/password
```

Reference format: `op://vault-name/item-name/field-name`

## Benefits

✅ **Secure** - Passwords never stored in config files  
✅ **Team sharing** - Share configs via git, passwords stay in 1Password  
✅ **Biometric unlock** - Uses your 1Password desktop app authentication  
✅ **No console windows** - CLI calls are hidden in the background  
✅ **Cross-platform** - Works on Windows, Mac, and Linux

## Troubleshooting

**"1Password CLI is not available"**
- Make sure `op` CLI is installed and in your PATH
- Run `op --version` in a terminal to verify installation

**"Failed to retrieve secret"**
- Make sure you're signed in to 1Password desktop app
- Verify "Integrate with 1Password CLI" is enabled in Settings → Developer
- Check that the vault name and item name in your reference are correct
- Ensure the item exists and you have access to the vault

**Still asking for authentication every time?**
- Make sure biometric unlock is enabled in 1Password settings
- The desktop app must be running when MremoteGO accesses secrets
