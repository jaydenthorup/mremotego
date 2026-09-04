# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Bitwarden integration**: connection passwords can be stored as `bw://<item-id>`
  references and are resolved through the Bitwarden CLI at connect time. Works
  with bitwarden.com, self-hosted Bitwarden and Vaultwarden.
  See [docs/BITWARDEN-SETUP.md](docs/BITWARDEN-SETUP.md).
- Optional field selector on references: `bw://<item-id>/username`, `/totp`,
  `/notes`; the password is the default.
- **Bitwarden item picker** in the add and edit connection dialogs, with search
  and a vault sync button, plus a "Store password in Bitwarden" option that
  creates a login item and replaces the password with its reference.
- Secret provider abstraction (`secrets.Provider`, `secrets.Registry`), so the
  configuration manager, launcher and GUI no longer depend on a single password
  manager.
- Unit tests for the secret providers, reference parsing and the encryption
  helper, plus a `go test ./...` step in CI.

### Changed
- The 1Password authentication warning at start-up is now a generic secret
  provider check, runs off the UI goroutine, and only asks about providers that
  the configuration actually references.

### Security
- The `bw serve` helper process is bound to loopback on a random port, started
  only when a Bitwarden reference is used, and terminated on exit. It is also
  placed in a Windows job object, and given a parent death signal on Linux, so
  it does not survive a crash.
- Passwords are no longer passed on a command line, where any local process
  could read them out of the process list:
  - SSH on Windows uses PuTTY `-pwfile` with a private temporary file that is
    deleted as soon as PuTTY has started, instead of `-pw`.
  - SSH on Linux and macOS passes the password to `sshpass` through the
    `SSHPASS` environment variable instead of `-p`; the temporary file is
    removed by the generated snippet before `ssh` starts.
  - RDP credentials are written to the Windows Credential Manager through the
    API instead of `cmdkey /pass:`.

### Known limitations
- On Linux, `xfreerdp` is still invoked with `/p:<password>`.

## [1.0.4] - 2026-01-28

### Fixed
- **Critical Fix**: 1Password references with special characters now work correctly
- Changed from `op read` to `op item get` for better special character handling
- Item names with parentheses, brackets, and other special chars now resolve properly
- Automatically decodes URL-encoded item names in references

### Technical Details
- `op://vault/(item)/field` references now work without encoding issues
- More robust error messages when 1Password retrieval fails

## [1.0.3] - 2026-01-28

### Fixed
- 1Password item names with special characters (parentheses, brackets, etc.) now work correctly
- Item names are now URL-encoded in references (e.g., `Server (Production)` → `Server%20%28Production%29`)
- GitHub Actions workflow now creates releases with proper binary assets when version tags are pushed

## [1.0.2] - 2026-01-28

### Added
- Recursive folder support in GUI dialogs - can now select nested folders at any depth
- Folder paths displayed with " / " separator (e.g., "Dev-Ops / Infrastructure / Builders")
- Add Folder dialog now allows creating subfolders within existing folders
- Helper functions: `collectAllFolders()`, `findConnectionParent()`, `findFolderByPath()`

### Changed
- Add Connection dialog now shows all nested folders in dropdown
- Edit Connection dialog correctly finds and displays deeply nested folder locations
- Can now move connections between folders at any nesting level

### Fixed
- GUI dialogs previously only supported root-level folders, now supports unlimited nesting depth

## [1.0.1] - 2026-01-28

### Changed
- Updated GUI with nested folder organization support
- Improved folder structure management

## [1.0.0] - 2026-01-27

### Added
- Initial release of MremoteGO
- Cross-platform connection manager (Windows, Linux, macOS)
- Support for SSH, RDP, VNC, HTTP, HTTPS, Telnet protocols
- 1Password integration for secure password management
- YAML-based connection storage
- GUI application using Fyne framework
- CLI commands: init, list, connect, add, edit, delete, export
- Folder organization for connections
- Connection search and filtering
- Password encryption support
- mRemoteNG XML import capability
