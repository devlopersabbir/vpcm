# VPSM Wiki & Operations Guide

Welcome to the official operations guide for **VPSM (VPS Manager)**, a developer-first platform built to securely manage, document, and connect to remote servers.

---

## 💾 1. Installation

### One-liner Installer (Recommended)

This script automatically detects your Operating System (Linux/macOS) and CPU architecture, downloads the latest compiled binary release from GitHub, installs the executables, and configures the `vpcm` shell alias/link:

```bash
curl -fsSL https://raw.githubusercontent.com/devlopersabbir/vpcm/main/scripts/install.sh | bash
```

### Build from Source

Ensure you have Go 1.25+ installed:

```bash
git clone https://github.com/devlopersabbir/vpcm.git
cd vpcm
make install
```

---

## 🛠️ 2. Configuration & Initialization

VPSM can use either a local-first **SQLite** database (completely CGO-free) or a remote **MongoDB** cluster.

### Setup/Initialize Configuration

To configure your database driver, paths, and API settings:

```bash
vpsm config init
```

This will walk you through selecting the driver, database path, local REST API hosting options, and output a config file to `~/.config/vpsm/config.yaml`.

### Edit Existing Settings

If you already have a configuration file, run:

```bash
vpsm config edit
```

This command pre-fills prompts with your current active settings.

### View Active Configuration

To view active keys in a styled, developer-friendly table:

```bash
vpsm config show
```

---

## 🖥️ 3. Managing the Server Inventory

Once your database is configured, manage your server profiles using these commands:

### Add a Server

Add a host named `production-db` pointing to IP `10.0.0.5`:

```bash
vpcm server add production-db 10.0.0.5
```

### List Registered Servers

Display a formatted table of all servers:

```bash
vpcm list
```

### Rename a Server

Rename a profile without changing the underlying host credentials or access keys:

```bash
vpcm server rename production-db prod-db-east
```

### Remove a Server

Remove a profile from the active database by its Name or database ID:

```bash
vpcm server remove prod-db-east
```

### Flush all Records

Wipe your database completely. This action requires double confirmation (first prompting `[y/N]` and then requiring you to type `FLUSH`):

```bash
vpcm server flush
```

---

## 🔍 4. Interactive TUI Explorer & Live Search

If you have a large list of servers, run the interactive **Bubble Tea Explorer**:

```bash
vpcm list -i
```

- **Live Filtering:** Start typing to instantly search by name, username, IP, or cloud provider.
- **PTY Navigation:** Use Up/Down arrow keys (or Ctrl+N / Ctrl+P) to select a host.
- **Instant Connection:** Press **Enter** on any host to exit the list explorer and immediately launch a full interactive SSH terminal session.

---

## 📤 5. Exporting Server Data

Export your server inventory list in several formats to stdout or file backups:

```bash
# Export in JSON format
vpsm server export -f json

# Export in YAML format
vpsm server export -f yaml

# Export in CSV sheet format
vpsm server export -f csv

# Export in SSH config format (suitable for ~/.ssh/config)
vpsm server export -f ssh

# Save directly to a local file
vpsm server export -f csv -o ~/backups/servers.csv
```

---

## 🔌 6. Connecting via SSH

Establish interactive PTY terminal connections:

```bash
vpcm ssh root@10.0.0.5
```

- **Custom Name Prompts:** If the host is not registered, VPSM will prompt you to enter a custom name upon successful login.
- **Provider Detection:** The SSH connection agent automatically analyzes cloud metadata services (IMDS), DMI hardware vectors, DNS, and ASNs to determine the cloud provider (AWS, GCP, etc.) and save it to the host profile.
- **Terminal Resize Listener:** Dynamically monitors window changes (`SIGWINCH`) on your machine and propagates them to the remote server to prevent broken terminal layouts.

---

## 🖥️ 7. Desktop Application Dashboard

In addition to the CLI, VPSM includes a beautiful, cross-platform Wails desktop application. Refer to the [Desktop Application Documentation](file:///Users/sabbir/own/vpcm/wiki/Desktop-App.md) for full setup guides, system stat panels, and details.
