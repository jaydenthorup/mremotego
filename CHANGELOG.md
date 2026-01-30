# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-01-30

### 🎉 MAJOR: Native 1Password SDK Integration

**Breaking Change**: MremoteGO now uses the 1Password SDK instead of CLI for password management.

### Added
- **Native 1Password desktop app integration** - No CLI commands needed!
- **Biometric authentication support** - Unlock with Touch ID/Windows Hello
- **Seamless UX** - Just unlock 1Password app and MremoteGO automatically authenticates
- **Automatic re-authentication** - Handles locked 1Password app gracefully
- **Better error messages** - Clear feedback when authentication is needed

### Changed
- **Breaking**: Requires 1Password desktop app (BETA version for SDK support)
- **Breaking**: CGO now required for builds (for SDK integration)
- Switched from `op` CLI to 1Password SDK v0.4.0-beta.2
- Improved authentication flow with desktop app integration

### Migration Guide

#### For Users:
1. Install 1Password desktop app (BETA version)
2. Go to Settings → Developer in 1Password
3. Enable "Integrate with the 1Password SDKs"
4. Enable "Integrate with other apps"
5. Restart MremoteGO - it will prompt for biometric auth when needed!

#### For Developers:
- CGO must be enabled: `$env:CGO_ENABLED="1"` (Windows) or `export CGO_ENABLED=1` (Unix)
- C compiler required (gcc/clang on Linux/Mac, MinGW on Windows)
- SDK provides better API than CLI (structured data, proper error types)

### Technical Details
- Uses 1Password SDK v0.4.0-beta.2 with Windows desktop app support
- Session lasts 10 minutes, auto-expires for security
- Automatic re-authentication when 1Password app is locked
- All `op://vault/item/field` references continue to work
- SDK handles special characters better than CLI ever could

### Why This Change?

**Old (CLI-based)**:
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
