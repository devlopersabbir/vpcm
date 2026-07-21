# Changelog

## [v0.1.12] - 2026-07-21

* feat: implement Docker containerization with compose support and automated CI/CD deployment pipeline (1774dcb)
* feat: add SSH-based inventory scanning and improve ServerDetail modal dismissibility (aa87c41)
* refactor: update ServerList UI with improved provider themes and expanded data models (0e2b3ad)
* feat: implement server favorite system with UI/CLI support and track recent connections (39a2ae9)
* Merge branch 'main' of github.com:devlopersabbir/vpcm into sabbir (44de21f)
* chore: bump version to v0.1.12 [skip ci] (8319d52)
* Merge branch 'main' of https://github.com/devlopersabbir/vpcm into sabbir (6ae9d6f)
* feat: add automated Windows and Unix installation/uninstallation scripts and update documentation (690a07d)
* feat: add session history component and update repository to support global connection logs (55a3bae)


## [v0.1.12] - 2026-07-20

* feat: add automated Windows and Unix installation/uninstallation scripts and update documentation (690a07d)


## [v0.1.12] - 2026-07-20

* fix: use GITHUB_WORKSPACE for absolute path resolution in macOS release build steps (73e3b80)
* refactor: overhaul frontend UI components with updated styling, border definitions, and linear gradients (9490e20)
* build: update Wails build command to include webkit2gtk_4_1 tags for Linux workflows (d3bfd95)
* build: upgrade libwebkit2gtk dependency from 4.0 to 4.1 in CI workflows (9ac9640)
* Merge branch 'dev' (2312ec8)
* feat: implement cross-platform Wails desktop application with build support and CI/CD integration (d296ae0)
* chore: remove unused package.json checksum file (c0d42cb)
* refactor: remove borders from glass UI components to adopt a cleaner, borderless design (5579bc9)
* feat: implement server management dashboard with Tailwind CSS and custom UI components (c165a8b)
* feat: initialize project structure with Wails desktop application boilerplate (f4e1e65)


## [v0.1.11] - 2026-07-20

* Merge branch 'dev' (f05ffa5)
* Merge branch 'dev' of github.com:devlopersabbir/vpcm into installer (6e4f011)
* fix: remove invalid local keyword in install.sh (300eb08)
* Merge branch 'dev' (b03bb7d)
* Merge branch 'sabbir' into dev (3b33a3c)
* Merge branch 'installer' into dev (d77adba)
* feat: improve installer with CLI arguments, robust error handling, shell wrapper support, and release workflow updates (4aa3de1)


## [v0.1.10] - 2026-07-19

* Merge branch 'dev' (f6b6c57)
* Merge branch 'sabbir' into dev (03d5100)
* feat: automate installation of vpsm, vpsmd, and vpsm-api binaries with non-interactive setup (5ef6f76)
* feat: add uninstall target to Makefile for removing binaries and shell wrappers (69937b3)


## [v0.1.9] - 2026-07-19

* Merge branch 'dev' (8d69f2c)
* Merge branch 'installer' into dev (5f1f8c1)
* feat: sanitize installation paths and force read input from tty in install and uninstall scripts (d2a417b)


## [v0.1.8] - 2026-07-19

* Merge branch 'dev' (3074f31)
* Merge branch 'sabbir' into dev (5bfc470)
* feat: add root endpoint providing API metadata, system uptime, and documentation links (ffe3e7d)
* feat: implement server metadata collection and update repository interfaces to support extended inventory views (42e8ce7)
* feat: add interactive installation enhancements and a new uninstaller script (53df2ee)


## [v0.1.7] - 2026-07-19

* Merge branch 'sabbir' (226b84b)
* docs: update project roadmap in TODO.md and document new API server CLI commands and response schemas (c9fa5d6)
* feat: expand server models with hardware metadata and implement connection logging functionality (27b0140)
* feat: implement local API server management commands and extend configuration for cloud-mode support (133a30f)


## [v0.1.6] - 2026-07-19

* Merge branch 'main' of github.com:devlopersabbir/vpcm (74adbff)
* ci: add github action workflow for cross-platform build validation on dev branch (50bf50d)


## [v0.1.5] - 2026-07-19

* Merge branch 'dev' (8ed8d9f)
* refactor: move version constant to cli.go and update release workflow path (28d7414)
* feat: implement CLI configuration management and automated installation script with GitHub Wiki sync automation (e8f28b5)
* feat: implement CLI command suite with SSH connection, data export, and project documentation (9b700bf)
* feat: add curl-based shell installer script and update README (a9d6df9)


## [v0.1.4] - 2026-07-19

* ci: fix undefined runTUI build error in release workflow (e09da60)
* Merge branch 'sabbir' (195e34d)
* bump: version to v0.1.3 (da9c9a8)
* docs: check off Interactive TUI Explorer in TODO.md (e975e8c)
* feat: add interactive Bubble Tea TUI explorer for server list and search (69356e4)


## [v0.1.3] - 2026-07-18

* Merge branch 'sabbir' (9e2f376)
* feat: add server export subcommand supporting ssh, json, csv, and yaml formats (f68a359)
* feat: add server rename subcommand and enforce connection details immutability (f9e4a49)


## [v0.1.2] - 2026-07-18

* Merge branch 'sabbir' (a43bc3b)
* docs: fully expand task items in TODO.md to match every roadmap detail (bca506b)
* docs: correct Focus Share weights in TODO.md status summary (1155362)
* docs: populate full 19-pillar platform vision in TODO.md (b2d9723)
* docs: append remote telemetry collection items to TODO.md (50cbf7e)
* docs: add TODO.md task board (72dcd27)
* docs: add MIT license (c5537eb)
* docs: add CONTRIBUTING.md guide (251e0dd)
* docs: improve README formatting and readability with additional spacing (d77a43f)
* docs: add premium architect-focused README.md (53c44ef)
* ci: compile multi-platform binaries and build contributors table in release notes (4106f2f)


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
