# 1Password SDK Setup

MremoteGO now uses the **1Password SDK** for much faster secret resolution!

## Setup Requirements

### Option 1: Desktop App Integration (Easiest - Recommended for Personal Use)

**Just install 1Password desktop app and you're done!**

1. Install [1Password desktop app](https://1password.com/downloads) if you haven't already
2. Make sure you're signed in to 1Password
3. Launch MremoteGO - it will automatically connect to your 1Password app!
4. First time you access a secret, 1Password will prompt for biometric unlock

**That's it!** No tokens, no CLI needed. The SDK talks directly to your 1Password app.

### Option 2: Service Account Token (For Teams/Automation)

For automated environments or team deployments:

1. Create a service account in 1Password:
   - Go to your 1Password account settings
   - Navigate to "Integrations" → "Service Accounts"
   - Create a new service account with vault access

2. Set the environment variable:
   ```powershell
   # Windows PowerShell
   $env:OP_SERVICE_ACCOUNT_TOKEN = "ops_your_token_here"
   
   # Or set it permanently in System Environment Variables
   [System.Environment]::SetEnvironmentVariable('OP_SERVICE_ACCOUNT_TOKEN', 'ops_your_token_here', 'User')
   ```

3. Launch MremoteGO - it will use the service account!

## Benefits of Using the SDK

✅ **10x faster** - No process spawning, direct API calls  
✅ **Persistent connection** - Keeps session alive  
✅ **Desktop app integration** - Uses biometric unlock from 1Password app  
✅ **No console windows** - Pure API, no CLI subprocess  
✅ **Secure** - No tokens exposed, uses OS keychain

## Usage

Just use the same `op://` reference format:

```yaml
connections:
  - name: Production Server
    host: prod.example.com
    username: admin
    password: op://Private/Production Server/password  # Works with SDK!
```

The SDK is **100% compatible** with existing configs - no changes needed!

## Troubleshooting

**"1Password SDK is not available"**
- Make sure 1Password desktop app is installed and you're signed in
- Or set `OP_SERVICE_ACCOUNT_TOKEN` for service account mode

**Biometric prompts not showing:**
- Check that 1Password app is running
- Verify biometric unlock is enabled in 1Password settings

**Still slow?**
- First access after app restart may be slower (SDK initialization)
- Subsequent accesses are instant (cached session)
