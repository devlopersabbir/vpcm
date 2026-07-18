# VPSM Platform Roadmap & Task Board

Welcome to the **Developer Infrastructure Platform** roadmap. This document serves as the master task list mapping every pillar, subsystem, and feature requirement of the VPSM ecosystem.

---

## 🚀 Platform Focus & Architecture Summary

| Pillar | Platform Share (Weight) | Focus Area | Status |
| :--- | :---: | :--- | :---: |
| **1. SSH Management** | 20% | Core secure connectivity & terminal shell | In Progress |
| **2. Infrastructure Inventory** | 30% | Deep host metadata, hardware, and runtime scans | In Progress |
| **3. Lightweight Monitoring** | 20% | Agentless metrics, historical tracking, and health checks | Planned |
| **4. Documentation** | 15% | Operational knowledge base, runbooks, and changelogs | In Progress |
| **5. DevOps Utilities** | 15% | Application control, Docker, logs, backups, and deploys | Planned |
| **6. Security Engine** | Add-on | SSH audit, firewalls, user inventory, CVE checks | Planned |
| **7. Automation Workflows** | Add-on | Scheduler, IF/THEN rules, webhooks, and hooks | Planned |
| **8. Event Engine** | Core Utility | Central timeline event logging & event timeline | Planned |
| **9. Cloud Integrations** | Core Utility | Multi-cloud discovery & automatic metadata sync | In Progress |
| **10. Terraform Integration** | Integration | Infrastructure sync & state/outputs parsing | Planned |
| **11. Kubernetes Support** | Integration | Cluster pods, services, nodes, and Helm inventory | Planned |
| **12. Notifications Engine** | Integration | Telegram, Slack, Discord, SMS, Webhooks | Planned |
| **13. REST API Interface** | Interface | Complete programmatic control over platform features | In Progress |
| **14. Web Dashboard** | Interface | Visual operations board & management console | Planned |
| **15. Team Management** | Enterprise | Role-based access control, organizations, audit logs | Planned |
| **16. Plugin Framework** | Extensibility | Open addon engine (AWS, Docker, Cloudflare plugins) | Planned |
| **17. AI Assistant** | Future | Infrastructure reasoning, troubleshooting, comparisons | Planned |
| **18. Infrastructure Graph** | Future | Visual dependency modeling (Customer -> DB) | Planned |
| **19. MCP Server** | Future | Model Context Protocol for safe AI agent operations | Planned |

---

## Detailed Task Board

### 1. SSH Management (20% Focus Share)

#### Connection Management
- [x] **Add Server:** Interactively register remote hosts into the database inventory.
- [x] **Remove Server:** Safely delete server credentials and profiles.
- [x] **Password Authentication:** Store and connect to hosts using passwords.
- [x] **SSH Key Management:** Load and authenticate using private keys.
- [x] **Auto Detect New Server:** Provider and capabilities detection on connection.
- [x] **Auto Save Connection:** Automatically record successfully established connections to database.
- [x] **Rename Server:** Rename existing server profile.
- [ ] **Clone Server Profile:** Clone configurations to spawn new profiles.
- [ ] **Duplicate Configuration:** Quick configuration cloning.
- [ ] **Import SSH Config:** Parse and import `~/.ssh/config` files.
- [ ] **Export SSH Config:** Export current database records to `~/.ssh/config`.
- [ ] **Auto Sync `~/.ssh/config`:** Automatically synchronize changes bidirectionally.
- [ ] **Favorite Servers:** Mark and group frequently accessed servers.
- [ ] **Recent Servers:** List recently connected hosts.
- [ ] **Connection History:** Track connection timestamp logs.
- [ ] **Multiple SSH Identities:** Support multi-key chains per host.
- [ ] **Proxy & Tunneling:**
  - [ ] Jump Host / Bastion support
  - [ ] ProxyJump integration
  - [ ] SOCKS Proxy tunneling
  - [ ] Dynamic Port Forwarding
  - [ ] Local Port Forwarding
  - [ ] Remote Port Forwarding

#### SSH Terminal
- [x] **Interactive Shell:** Full interactive PTY terminal access.
- [x] **Dynamic Terminal Resizing:** Propagate window size changes (`SIGWINCH`) to remote terminal.
- [x] **Verify Fingerprints:** Check and save host fingerprints.
- [ ] **Non-interactive Execution:** Run commands remotely and return standard output.
- [ ] **Multiple Terminals:** Side-by-side terminal session layouts.
- [ ] **Broadcast Commands:** Execute single command across multiple terminals simultaneously.
- [ ] **Parallel Execution:** Execute commands concurrently on multiple hosts in the background.
- [ ] **Session Recording:** Capture terminal input/output logs.
- [ ] **Session Replay:** Play back recorded terminal sessions.
- [ ] **Clipboard Support:** Host-to-guest copy/paste bridging.
- [ ] **Terminal Themes:** Custom styles, colors, and fonts.

#### File Transfer
- [ ] **SCP Support:** Standard secure copy file transfer.
- [ ] **SFTP Engine:** Full SFTP client.
- [ ] **Upload:** Copy files from local host to guest.
- [ ] **Download:** Retrieve files from guest to local host.
- [ ] **Directory Sync:** Bidirectional directory synchronization.
- [ ] **Resume Transfer:** Automatically pick up aborted file transfers.
- [ ] **Compression:** Compress files during transfers.

#### SSH Utilities
- [x] **Test Connection:** Doctor check to verify connectivity parameters.
- [ ] **Generate SSH Key:** Create new keys on-the-fly.
- [ ] **Rotate Keys:** Automatically replace old authorized keys with new keys.
- [ ] **Benchmark Latency:** Measure host ping and network response times.
- [ ] **Keep Alive:** Send keep-alive packets to prevent SSH idle drops.
- [ ] **Auto Reconnect:** Reconnect instantly when connection drops.

---

### 2. Infrastructure Inventory (30% Focus Share)

#### Server Information
- [x] **Hostname:** Retrieve the actual hostname.
- [x] **Public IP:** Fetch the public IPv4 address.
- [x] **Provider:** Auto-resolve VPS hosting provider (AWS, GCP, etc.).
- [ ] **Private IP:** Scan for internal network IPs.
- [ ] **MAC Address:** Read MAC addresses.
- [ ] **Region:** Identify hosting region.
- [ ] **Availability Zone:** Find specific provider datacenter zone.
- [ ] **Virtualization:** Identify hypervisor type (KVM, Xen, Hyper-V).
- [ ] **Instance Type:** Match instance size names (e.g. `t3.medium`).
- [ ] **Serial Number:** Read bios/system serial numbers.
- [ ] **BIOS:** Identify BIOS version.
- [ ] **Uptime:** Retrieve host system uptime.

#### Hardware specs
- [ ] **CPU:** Read CPU model, cores, and threads.
- [ ] **RAM:** Scan total and available memory.
- [ ] **Disk:** Audit storage partitions, types, and free space.
- [ ] **Swap:** Check Swap partition sizes.
- [ ] **GPU:** Check for graphic accelerator units.
- [ ] **Architecture:** Identify CPU architecture (`amd64`, `arm64`, etc.).
- [ ] **Temperature:** Read sensor temperatures.
- [ ] **Network Speed:** Measure link speed capacity.

#### Operating System
- [x] **Distribution:** Detect OS name (Ubuntu, Debian, CentOS, Alpine).
- [x] **Version:** Read OS release version.
- [ ] **Kernel:** Read running kernel version.
- [ ] **Timezone:** Retrieve system timezone setting.
- [ ] **Locale:** Identify language and locale parameters.
- [ ] **Init System:** Identify init manager (Systemd, OpenRC, SysVInit).
- [ ] **Package Manager:** Determine active package managers (`apt`, `dnf`, `apk`).

#### Installed Software Detection
- [ ] **Detect Key Platforms:**
  - [ ] Runtimes: Node, Python, Go, Java, Rust, PHP
  - [ ] Containers: Docker, Podman, Kubernetes, k3s
  - [ ] Web Servers: Caddy, Nginx, Apache
  - [ ] Databases: PostgreSQL, MySQL, Redis, MongoDB
  - [ ] Brokers: RabbitMQ, Kafka, Elasticsearch
  - [ ] CI/CD: Jenkins, GitLab Runner

#### Containers Inventory
- [ ] **Images:** List cached images.
- [ ] **Containers:** Audit active/stopped containers.
- [ ] **Networks:** Audit custom networks.
- [ ] **Volumes:** Monitor mapped storage volumes.
- [ ] **Compose Projects:** Group running containers by Docker Compose definitions.

#### Services Inventory
- [ ] **Systemd Audit:** List services and filter by:
  - [ ] Enabled
  - [ ] Disabled
  - [ ] Failed

#### Network Configuration
- [ ] **Interfaces:** List network adapters.
- [ ] **Routing:** Audit IP routing tables.
- [ ] **DNS:** Read `/etc/resolv.conf` settings.
- [ ] **Firewall:** Identify active configurations (UFW, iptables).
- [ ] **Open Ports:** List active TCP/UDP ports.
- [ ] **Listening Services:** Match open ports to active processes.

#### Search & Classification
- [x] **Tags System:** Add labels to group server types.
- [ ] **Search Command:** Query inventory profiles (e.g. `vpsm search docker`).

---

### 3. Lightweight Monitoring (20% Focus Share)

- [ ] **Live Metrics:** Real-time stream of CPU, RAM, Disk, Swap, Load, Temperature, Network, and Processes.
- [ ] **History Logs:** Store hourly, daily, weekly, and monthly metric snapshots.
- [ ] **Alerts:** Trigger alerts on:
  - [ ] CPU > 90%
  - [ ] Disk > 80%
  - [ ] RAM > 90%
  - [ ] Service Down
  - [ ] SSH Connection Failed
- [ ] **Health Status Classification:** Mark servers as `Online`, `Offline`, `Warning`, or `Critical`.
- [ ] **Service Monitoring:** Continually inspect running Docker containers, databases, systemd services, and web servers.
- [ ] **Terminal Dashboard:** Print real-time dashboard inside terminal (CPU, RAM, Disk, Top Processes, Containers, Services).

---

### 4. Documentation (15% Focus Share)

- [x] **Server Notes:** Add markdown notes to servers (purpose, owner, credentials).
- [ ] **Markdown Renderer:** Support styled terminal/web rendering with images, Mermaid diagrams, code blocks, and links.
- [ ] **Attachments:** Link architectural diagrams, PDFs, and screenshots.
- [ ] **Runbooks:** Create operational runbooks (e.g. "Restart API", "Recover Database", "Deploy Backend").
- [ ] **Architecture maps:** Store network topology and server dependency maps.
- [ ] **Server Changelogs:** Save server modification logs (e.g. "Updated PostgreSQL", "Migrated Database").

---

### 5. DevOps Utilities (15% Focus Share)

- [ ] **Docker integration:** Wrappers for `ps`, `logs`, `exec`, `restart`, `pull`, `prune`, and `compose up/down`.
- [ ] **Service Controls:** CLI interfaces for `systemctl` restart, status, logs, enable, and disable.
- [ ] **Cron Manager:** Monitor and schedule cronjobs and systemd timers.
- [ ] **Log Engine:** Quick views for `journalctl`, `docker logs`, `tail`, and `grep`.
- [ ] **System Updates:** List pending OS packages to update.
- [ ] **Backups:** Trigger and schedule backups for databases, files, and Docker volumes.
- [ ] **Restore Engine:** Recover servers from backup files.
- [ ] **Deployments:** Run Git deployment routines (`git pull`, `docker compose`, `systemctl`).
- [ ] **Script Store:** Save and run reusable operational scripts (`cleanup.sh`, `backup.sh`).
- [ ] **Secrets & Environments:** Manage host-specific `.env` secrets and environment variables.

---

### 6. Security Engine

- [ ] **SSH Audit:** Report on weak ciphers, active password auth, and root login states.
- [ ] **Firewall Controls:** Manage firewalls (UFW, Firewalld, iptables, nftables).
- [ ] **Fail2Ban Audits:** Monitor service status, bans, and attempts.
- [ ] **User Audit:** Monitor active users, sudo privileges, and groups.
- [ ] **Keys Audit:** Scan authorized SSH keys.
- [ ] **Vulnerability (CVE) Check:** Scan installed package versions against known vulnerabilities.
- [ ] **CIS Compliance:** Basic CIS benchmark checks.

---

### 7. Automation

- [ ] **Automation Scheduler:** Schedule tasks (daily scans, weekly audits, monthly backups).
- [ ] **Rules Workflows:** Build actions using triggers (e.g. IF Disk > 90% THEN send notification).
- [ ] **Hooks Engine:** Trigger workflows on server lifecycle changes (Server Added, Deleted, Deployment Finished).
- [ ] **Outgoing Webhooks:** Push system updates to target URLs.

---

### 8. Events Timeline

- [ ] **Event Timeline Logging:** Chronologically record server events (Server Added, Connected, Alert, Audit).
- [ ] **Timeline View:** Audit timeline history from the CLI or dashboard.

---

### 9. Cloud Integrations

- [x] **Automatic Cloud Provider Detection:** (Completed via IMDS, reverse DNS, and ASN lookups).
- [ ] **Instance Discovery:** Connect and discover VPS listings from AWS, Azure, GCP, DO, Hetzner, Linode, Vultr, and OVH.
- [ ] **Metadata Synchronization:** Sync instance types, regions, volumes, and elastic IPs.

---

### 10. Terraform Integration

- [ ] **Terraform Parser:** Read state files and auto-import server configs and outputs.

---

### 11. Kubernetes Support

- [ ] **Cluster Inventory:** Scan for active pods, nodes, deployments, services, ingress configurations, and Helm charts.

---

### 12. Notifications Engine

- [ ] **Integrations:** Push alerts to Telegram, Slack, Discord, Email, Webhooks, and SMS.

---

### 13. REST API Interface

- [x] **Server Management API:** Endpoints to add, list, and delete servers.
- [ ] **API Surface Expansion:** Programmatic access for monitoring metrics, notes, runbooks, and event timelines.

---

### 14. Web Dashboard

- [ ] **Visual Management UI:** Design interactive panels for dashboard stats, server status, containers, logs, notes, and settings.

---

### 15. Team Management

- [ ] **Access Controls:** User authentication, organizations, projects, roles, and administrative audit logs.

---

### 16. Plugin Framework

- [ ] **Extensibility Engine:** Define and load external plugins for third-party cloud tools and platforms.

---

### 17. AI Assistant (Future)

- [ ] **Reasoning Engine:** Leverage LLMs to troubleshoot slow servers, parse logs, recommend fixes, and compare settings.

---

### 18. Infrastructure Graph (Future)

- [ ] **Dependency Maps:** Draw maps linking Customers, Projects, Environments, Servers, Containers, and Databases.

---

### 19. Model Context Protocol (MCP) Server (Future)

- [ ] **AI-Agent Connection:** Expose standard context endpoints for AI systems to query inventory, inspect metrics, and trigger operations.
