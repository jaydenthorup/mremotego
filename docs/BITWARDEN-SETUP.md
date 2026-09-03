# Bitwarden Setup Guide

MremoteGO can read connection passwords from Bitwarden, so that:

- Passwords never appear in your config files
- Configs stay safe to commit and share via git
- Credentials can be rotated in one place for the whole team

This works with bitwarden.com, self-hosted Bitwarden and Vaultwarden.

## How it works

Bitwarden has no library interface for other applications, so MremoteGO uses
the official Bitwarden CLI. On first use it starts `bw serve` as a hidden child
process bound to `127.0.0.1` on a random free port, reads what it needs over
that local API, and stops the process when MremoteGO exits.

MremoteGO never asks for, sees or stores your master password. The child
process inherits the `BW_SESSION` variable, so you unlock the vault once in
your terminal and every lookup after that is authorised by that session.

## Quick Setup

### 1. Install the Bitwarden CLI

```powershell
winget install Bitwarden.CLI
```

```bash
# macOS
brew install bitwarden-cli

# Linux
npm install -g @bitwarden/cli
```

Verify with `bw --version`.

### 2. Point the CLI at your server

Only needed for self-hosted Bitwarden or Vaultwarden:

```powershell
bw config server https://vault.example.com
```

### 3. Log in and unlock

```powershell
bw login
$env:BW_SESSION = bw unlock --raw
bw status          # should report "unlocked"
```

```bash
bw login
export BW_SESSION="$(bw unlock --raw)"
bw status
```

### 4. Start MremoteGO from that same terminal

```powershell
.\mremotego.exe
```

The session key is inherited by MremoteGO and by the `bw serve` process it
starts. Launching MremoteGO from a desktop shortcut instead will leave the
vault locked, and it will show you these instructions.

## Using Bitwarden References

### In your config

Instead of a plain text password, store a reference to the vault item:

```yaml
connections:
  - name: "Production Server"
    protocol: ssh
    host: prod.example.com
    username: admin
    password: bw://8f3c1d9a-4e2b-4c77-9f10-1a2b3c4d5e6f
```

Reference format:

| Reference | Resolves to |
|-----------|-------------|
| `bw://<item-id>` | the item's password |
| `bw://<item-id>/password` | the item's password |
| `bw://<item-id>/username` | the item's username |
| `bw://<item-id>/totp` | the current TOTP code |
| `bw://<item-id>/notes` | the item's notes |

References use the item id rather than its name, so renaming an item in
Bitwarden does not break your config.

### Finding an item id

**In the GUI:** click **Bitwarden...** next to the password field when adding
or editing a connection. Search your vault, pick an item, and the reference is
filled in for you along with the username.

**On the command line:**

```powershell
bw list items --search "production" | ConvertFrom-Json | Select-Object id, name
```

### Creating items from MremoteGO

1. Type the password into the password field as usual
2. Tick **Store password in Bitwarden**
3. Save

MremoteGO creates a login item named after the connection, with the username
and a URI of `<protocol>://<host>`, then replaces the password in your config
with the new `bw://` reference.

## Security notes

- **`bw serve` has no authentication of its own.** Anything running as your
  user on your machine could talk to it while it is up. MremoteGO limits the
  exposure by binding it to loopback on a random port, starting it only when a
  Bitwarden reference is actually used, and terminating it on exit. On Windows
  the process is placed in a job object and on Linux it gets a parent death
  signal, so it also dies if MremoteGO crashes.
- **`BW_SESSION` unlocks your whole vault.** Treat it like a password: do not
  put it in a script that others can read, and close the terminal when done.
- **References are not encrypted at rest**, on purpose. They contain no secret
  material, and leaving them readable is what keeps configs diffable in git.
  This matches how `op://` references are handled.

## Troubleshooting

**"Bitwarden CLI (bw) is not installed or not in PATH"**
Install the CLI and make sure `bw --version` works in the same terminal you
start MremoteGO from.

**"bitwarden vault is locked"**
The session key is missing or expired. Run `bw unlock --raw` again, set
`BW_SESSION`, and restart MremoteGO from that terminal.

**"not logged in to bitwarden"**
Run `bw login`. For a self-hosted server, run `bw config server <url>` first.

**An item you just created elsewhere is not in the picker**
The CLI serves items from a local cache. Click **Sync vault** in the picker, or
run `bw sync`.

**RDP connects but prompts for the password**
RDP is deliberately allowed to continue when a reference cannot be resolved, so
you still get a login prompt instead of an error. Check the vault state as
above.

## See also

- [Password Management](PASSWORD-MANAGEMENT.md) - all password options
- [1Password Setup](1PASSWORD-SETUP.md) - the other supported password manager
- [Bitwarden CLI documentation](https://bitwarden.com/help/cli/)
