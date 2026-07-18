# Changelog

All notable changes to this project will be documented in this file.

## [0.1.0] - 2026-07-18

### Added
- **SQLite Database Support**: Support pure-Go CGO-free SQLite database driver alongside MongoDB, with auto-migrations on startup.
- **Interactive Configuration Subcommands**: Added interactive `vpsm config init` (to set up a new config file), `vpsm config edit` (to edit settings showing current values as defaults), and `vpsm config show` (renders active configuration in a beautiful table).
- **Auto Release Pipeline**: Integrated a GitHub Actions workflow to automate patch version bumping, CHANGELOG prepending, Git tagging, and GitHub Release generation on merge to main.
- **Automatic Cloud Provider Detection**: Integrated a multi-source detection engine that identifies whether a host is running on AWS, GCP, Azure, DigitalOcean, Hetzner, Linode, Vultr, or OVH. It uses remote IMDS requests, DMI/hardware vendor profiles, installed agent checks, reverse DNS, and local ASN lookups.
- **Database Inventory Flush Command**: Added `vpsm server flush` (`flash` alias) to securely clear all server records and reset auto-increment counters with a double-confirmation mechanism.
- **Root-level `list` Command**: Added a direct root-level command `vpcm list` as a fast shortcut for `vpcm server list`.
- **Dynamic Terminal Resizing**: Implemented SIGWINCH signals propagation using `WindowChange` to dynamically resize remote terminals upon host terminal resize.
- **Custom Server Naming**: Added interactive user prompt when registering a new server dynamically via SSH, preventing default non-descriptive names like `user@ip`.

### Changed
- **Strict Subcommand Conventions**: Removed the implicit automatic SSH-rewrite routing. Unknown commands will now correctly report a CLI syntax error rather than automatically attempting to SSH into the command argument name. You must explicitly start connections using `vpcm ssh`.
