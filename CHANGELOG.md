# Changelog

## [v0.1.1] - 2026-07-18

* docs: update CHANGELOG.md with SQLite and config features for v0.1.0 (49d0171)
* feat: conditionally register config init or edit subcommand based on file presence (4c28499)
* feat: add SQLite support and interactive configuration subcommands (a9d134a)
* ci: include git commit logs in CHANGELOG and GitHub Release description (8d0cb8f)
* ci: add auto release workflow on push to main (0d7f078)
* release: v0.1.0 (29cfd70)
* feat: add list command to root and remove implicit ssh command injection (589dc02)
* feat: implement multi-source cloud provider detection and integrate into server registration flow (cf18341)
* feat: add list command as alias for server list in CLI (785a091)
* feat: add server flush command with double confirmation to clear database inventory (01abb81)
* feat: add interactive prompt for custom server naming and document features in new roadmap file (c05a953)
* feat: add cross-platform SIGWINCH handling to propagate terminal resize events to remote sessions (457cacc)
* refactor: remove server management endpoints and add official API documentation (482860e)
* refactor: migrate database storage from GORM/SQL to MongoDB across repositories and services (9abd5fa)
* feat: add ssh command with credential storage and interactive shell support, including server management and CLI enhancements. (4d8273e)
* feat: initialize project structure with core services, database, and modular CLI/API architecture (3a7c2cf)


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
