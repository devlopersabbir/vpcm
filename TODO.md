# Project Task Board

This task board documents completed enhancements and future tasks for the VPS Connection Manager (VPSM) project.

## Initial Plan & Feature Board

### Dynamic Shell & Session Handling
- [x] **Terminal Resizing:** Implement `SIGWINCH` listener and propagate window size updates to remote terminals dynamically using `WindowChange` (on non-Windows hosts).
- [x] **Strict CLI Conventions:** Enforce standard command paradigms. Prevent unknown keywords from implicitly starting SSH connections; connections must explicitly start with `vpcm ssh`.
- [x] **Root-level `list` Command:** Support direct execution of `vpcm list` at root as a shortcut for `vpsm server list`.

### Database & Configuration Improvements
- [x] **Pluggable Database Engines:** Add support for pure-Go SQLite alongside MongoDB, complete with database setup auto-migrations.
- [x] **Custom Server Naming:** Prompt user for a custom name when dynamically connecting/saving a new server instead of forcing `user@ip`.
- [x] **Database Inventory Flush:** Implement `vpsm server flush` with a double-confirmation step to clear all server records.
- [x] **Interactive Config Setup:** Add `vpsm config init` to set up connections interactively and sync them to `~/.config/vpsm/config.yaml`.
- [x] **Interactive Config Edit:** Add `vpsm config edit` (registers conditionally if config exists) to edit settings.
- [x] **Styled Configuration Status:** Renders current configuration details using a formatted Lipgloss table.

### Cloud Integration & Automation
- [x] **Cloud Provider Detection:** Map IMDS endpoints, DMI identifiers, installed cloud agents, reverse DNS, and ASNs to automatically set the server's hosting provider field.
- [x] **Automated Release Workflows:** Add GitHub Actions pipeline to auto-bump patch versions, tag releases, and build release logs from git commits.
- [x] **Multi-Platform Assets:** Compile and bundle CLI binaries for Windows, Linux, and macOS platforms inside the release assets.
- [x] **Visual Contributors List:** Fetch contributor details and display avatar tables automatically inside the release notes.

### Project Setup
- [x] Create developer documentation ([README.md](file:///Users/sabbir/own/vpcm/README.md) and [CONTRIBUTING.md](file:///Users/sabbir/own/vpcm/CONTRIBUTING.md)).
- [x] Create project licensing ([LICENSE](file:///Users/sabbir/own/vpcm/LICENSE)).
