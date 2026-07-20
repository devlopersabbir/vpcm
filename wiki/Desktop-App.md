# VPSM Desktop Application

Ultimate VPSM Solutions includes a custom Wails desktop application for visualizing server specs, installed packages, and session logs inside a premium, borderless React dashboard.

## Features

- **Borderless Glassmorphism UI**: Beautiful, dark-themed dashboard powered by Tailwind CSS v4 and Lucide-React.
- **Server Inventory Stats**: Instant overview of your total node counts, active auth types, and cloud providers.
- **Auditing & Live Specs**: Deep-dive specifications organized into tabs:
  - **Overview**: UUID, creation dates, last seen timestamps, and provider tags.
  - **Hardware Specs**: CPU models, core count, RAM, swap details, storage utilization, virtualization technology, and system uptime.
  - **Operating System**: Linux distribution info, kernel version, package managers, and timezone.
  - **Network Config**: System hostnames, public and private IPs, MAC addresses, and availability zones.
  - **Installed Apps**: Lists all detected software packages on the target host.
  - **Sessions/Connection Logs**: Audit history for remote SSH logins.

## Installation & Running

### Build from Source
Ensure you have Wails CLI, Go 1.25+, and Node.js installed:

1. Clone the repository and navigate to the project directory:
   ```bash
   git clone https://github.com/devlopersabbir/vpcm.git
   cd vpcm
   ```
2. Build the application using the Makefile:
   ```bash
   make build-desktop
   ```
   This compiles the frontend assets and packages the app binary under `app/vpsm-desktop/build/bin/`.

### Development Mode
To run the desktop application locally with hot-reloading:
```bash
make dev-desktop
```
