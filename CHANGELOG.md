# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
