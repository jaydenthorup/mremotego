# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.1-alpha] - 2026-01-30

### Added
- Enhanced 1Password authentication instructions showing both SDK and CLI options
- Real-time status display for SDK and CLI availability
- Authentication state indicators in setup dialog
- VS Code workspace configuration for easy development

### Changed
- Improved authentication instructions popup with clearer options
- Better guidance for users choosing between SDK and CLI methods

## [2.0.0] - 2026-01-30

### 🎉 MAJOR: Native 1Password SDK Integration with CLI Fallback

**Breaking Change**: MremoteGO now uses the 1Password SDK for enhanced desktop app integration.

### Added
- **Native 1Password desktop app integration** - Biometric auth with SDK
- **Automatic CLI fallback** - Falls back to `op` CLI if SDK unavailable
- **Vault name mapping** - Use friendly names instead of UUIDs in config
- **Biometric authentication support** - Touch ID/Face ID/Windows Hello
- **Seamless UX** - Just unlock 1Password app for automatic authentication
- **Automatic re-authentication** - Handles locked 1Password app gracefully
- **Better error messages** - Clear feedback with setup instructions
- **VS Code integration** - Build tasks, debug configs, recommended extensions
- **Comprehensive test suite** - 35+ unit tests with good coverage
- **CLI smoke tests** - Automated testing for CLI commands

### Changed
- **Breaking**: CGO now required for builds (for SDK and Fyne GUI)
- Switched from CLI-only to SDK with CLI fallback
- Improved authentication flow with dual-mode support
- Settings now persisted correctly (DeepCopy bug fixed)
- Enhanced copyright information in About dialog

### Migration Guide

#### For Users:

**Option 1: Desktop App (Recommended)**
1. Install 1Password desktop app (BETA version)
2. Go to Settings → Developer in 1Password
3. Enable "Integrate with the 1Password SDKs"
4. Enable "Integrate with other apps"
5. Configure in config.yaml:
   ```yaml
   settings:
     onePasswordAccount: "your-account-name"
     vaultNameMappings:
       Personal: "uuid-here"
   ```
6. Restart MremoteGO - biometric auth will prompt when needed!

**Option 2: CLI Fallback (Automatic)**
1. Install 1Password CLI: `op`
2. Sign in: `op signin`
3. Launch MremoteGO from same terminal - works automatically!

#### For Developers:
- CGO must be enabled: `$env:CGO_ENABLED="1"` (Windows) or `export CGO_ENABLED=1` (Unix)
- C compiler required (gcc/clang on Linux/Mac, MinGW/TDM-GCC on Windows)
- VS Code configuration included for easy setup
- Run tests: `go test ./...`

### Technical Details
- Uses 1Password SDK v0.4.0-beta.2 with desktop app support
- Automatic fallback to CLI provider if SDK initialization fails
- Vault name mapping stored in settings.vaultNameMappings
- Session lasts 10 minutes with SDK, auto-expires for security
- All `op://vault/item/field` references continue to work
- SDK handles special characters and URL encoding automatically
- DeepCopy implemented for Settings struct to prevent data loss
- 35+ unit tests covering models, config, secrets, and CLI commands

### Why This Change?

**New (SDK with CLI fallback)**:
- Required `op signin` in terminal before launching GUI
- Session tokens only worked when launched from same terminal
- Complex setup with environment variables
- No biometric support

**New (SDK-based)**:
- Just unlock 1Password app
- Biometric authentication (Touch ID/Windows Hello/Face ID)
- Works from any launcher (desktop shortcut, Start menu, etc.)
- Seamless UX matching native 1Password experience

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
