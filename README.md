# VPSM (VPS Manager)

A zero-dependency, CGO-free, single-binary CLI and REST API engine built in Go to seamlessly manage, document, and connect to your remote virtual private servers.

Designed for operations engineers and developers who want a local-first connection inventory without heavy external dependencies.

---

## Why VPSM?

Most VPS managers are either bloated web interfaces or simple SSH alias files that don't scale. VPSM provides a structured database inventory (SQLite or MongoDB) combined with a high-performance SSH connector.

- **Pluggable Database Architecture:** Run local-first with a pure-Go **SQLite** engine, or scale to a remote **MongoDB** instance.
- **Cloud Provider Auto-Detection:** Automatically identifies whether a host is running on AWS, GCP, Azure, DigitalOcean, Hetzner, Linode, Vultr, or OVH using IMDS endpoints, DMI identifiers, agent detection, DNS patterns, and ASN records.
- **Dynamic Terminal Resizing:** Propagates local window changes (`SIGWINCH`) to remote sessions dynamically.
- **API-Ready:** Includes a built-in REST API wrapper (`vpsm-api`) to integrate your server inventory with external automation pipelines.
- **Strict SSH Conventions:** Enforces explicit connection paradigms (`vpsm ssh <host>`) to prevent shell typos from translating into accidental connection requests.

---

### 1. Installation

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/devlopersabbir/vpcm/main/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/devlopersabbir/vpcm/main/scripts/install.ps1 | iex
```

**Build From Source:**
Alternatively, compile from source:
```bash
make install
```

### 2. Configure Your Universe

Initialize your configuration interactively:

```bash
vpsm config init
```

This creates/updates `~/.config/vpsm/config.yaml`. To view your active settings:

```bash
vpsm config show
```

### 3. Basic Commands

```bash
# List all registered servers
vpcm list

# Connect to a server (automatically prompts for password/identity and saves credentials)
vpcm ssh root@192.168.1.100

# Add a server manually to inventory
vpcm server add prod-db-1 192.168.1.101

# Remove a server
vpcm server remove prod-db-1

# Flush/wipe the inventory (requires double-confirmation)
vpcm server flush
```

---

## Architectural Principles

VPSM is built around clean architecture rules for maximum lifespan:

1. **Dependency Inversion:** Service domains bind to interfaces. Swapping database repositories is simple and doesn't affect connection execution flows.
2. **Event Bus Decoupling:** Inter-module communication is asynchronous and handled through an event broker, eliminating tight package coupling.
3. **Local-First, Scale-Second:** SQLite satisfies the immediate local requirement, while the Mongo driver serves cloud storage needs.
