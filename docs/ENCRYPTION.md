# Password Encryption

MremoteGO supports encrypting passwords at rest in the configuration file using AES-256-GCM encryption.

## How It Works

- **Master Password**: You set a master password when launching the GUI
- **Key Derivation**: PBKDF2 with 100,000 iterations derives an AES-256 key from your master password
- **Encrypted Format**: Passwords are stored as `enc:base64(salt+nonce+ciphertext)`
- **No Encryption**: Leave the master password blank to store passwords in plain text

## Security Features

- **AES-256-GCM**: Industry-standard authenticated encryption
- **Unique Salt**: Each password gets a random 16-byte salt
- **Random Nonce**: Each encryption uses a unique nonce
- **1Password Integration**: 1Password references (`op://...`) are NOT encrypted (no need)

## Using Encryption in the GUI

1. Launch MremoteGO
2. Enter a master password (or leave blank for no encryption)
3. Add/edit connections - passwords are automatically encrypted when saved
4. Next time you launch, enter the same master password to decrypt

## Command Line Tool: encrypt-passwords

Encrypt or decrypt an existing config file:

### Encrypt Passwords

```bash
# Windows
.\encrypt-passwords.exe -encrypt

# Linux/Mac
./encrypt-passwords -encrypt

# Specify config path
./encrypt-passwords -encrypt -config path/to/config.yaml
```

### Decrypt Passwords

```bash
# Windows
.\encrypt-passwords.exe -decrypt

# Linux/Mac
./encrypt-passwords -decrypt
```

**Warning**: Decrypting saves passwords in plain text!

## Password Storage Options

You have three options for storing passwords:

1. **Plain Text** (no master password)
   ```yaml
   password: mypassword123
   ```

2. **Encrypted** (with master password)
   ```yaml
   password: enc:base64encodedencrypteddata
   ```

3. **1Password Reference** (most secure)
   ```yaml
   password: op://DevOps/server01/password
   ```

## Best Practices

- **1Password**: Use for maximum security (passwords stored in 1Password vault)
- **Encrypted**: Good for standalone use without 1Password
- **Plain Text**: Only use for non-sensitive connections or development

## Technical Details

- **Algorithm**: AES-256-GCM
- **Key Derivation**: PBKDF2-HMAC-SHA256
- **Iterations**: 100,000
- **Salt Size**: 16 bytes (random per password)
- **Nonce Size**: 12 bytes (random per encryption)

## Migration

If you have an existing unencrypted config:

1. Option A: Use the GUI (next launch will prompt for password)
2. Option B: Use `encrypt-passwords -encrypt` tool

If you want to remove encryption:

1. Option A: Leave password blank on next launch (won't re-encrypt)
2. Option B: Use `encrypt-passwords -decrypt` tool

## Troubleshooting

**"Failed to load config (wrong password?)"**
- You entered the wrong master password
- Config was encrypted with a different password
- Config file is corrupted

**"Passwords not staying encrypted"**
- Make sure you're entering a master password when launching
- If you leave it blank, passwords save as plain text

**"Can I change my master password?"**
1. Decrypt with old password: `encrypt-passwords -decrypt`
2. Encrypt with new password: `encrypt-passwords -encrypt`

## Security Considerations

- Master password is **never stored** - you must remember it
- Lost master password = **lost access** to encrypted passwords
- Config file permissions are set to `0600` (owner read/write only)
- Consider using 1Password integration for better security
