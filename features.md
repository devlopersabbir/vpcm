# VPSM Features & Roadmap

This file documents all currently implemented features and future plans for the VPS Connection Manager (VPSM) project, organized date-wise using a checklist format.

## Features & Plan

### 2026-07-18

- [x] **Dynamic Terminal Resizing**: Support interactive shell window resizing via SIGWINCH (unix platforms) and update SSH sessions dynamically via `WindowChange`.
- [x] **New Server Custom Naming**: Prompt user for a human-readable name when dynamically connecting and registering a new server in the database (instead of default `user@ip`).
- [x] **Database Flush/Flash Command**: Add `vpsm server flush` (`flash` alias) to clear all servers from the database with a double confirmation step.
- [x] **Automatic Cloud Provider Detection**: Implement multi-step provider detection (IMDS, DMI, cloud agents, reverse DNS, and ASN lookup) to automatically set the `Provider` field in the database.

### Foundation Features (v0.0.1)

- [x] **Configuration Validation**: Verify setup configurations using the `doctor` command.
- [x] **Config Inspection**: Show loaded database driver, path, API, and logger status using the `config` command.
- [x] **Server Inventory**: Add, list, and remove servers registered in the database inventory (`vpsm server`).
- [x] **SSH Connection Session**: Seamlessly connect to remote servers (`vpsm ssh` / `vpsm <host>`) using key-based or password-based authentication.
- [x] **Auto-saving Credentials**: Automatically save credentials to the database on successful connection.
