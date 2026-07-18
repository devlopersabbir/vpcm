# VPSM Platform Roadmap & Task Board

Welcome to the **Developer Infrastructure Platform** roadmap. VPSM is evolving from a simple SSH manager into a unified operations ecosystem.

---

## 🚀 Pillar Status Summary

| Pillar | Status | Progress |
| :--- | :---: | :---: |
| **1. SSH Management** | In Progress | 40% |
| **2. Infrastructure Inventory** | In Progress | 25% |
| **3. Lightweight Monitoring** | Planned | 0% |
| **4. Documentation & Operational Notes** | In Progress | 30% |
| **5. DevOps Utilities** | Planned | 0% |
| **6. Security Audits** | Planned | 0% |
| **7. Automation & Scheduler** | Planned | 0% |
| **8. Event Engine** | Planned | 0% |
| **9. Cloud Provider Integrations** | In Progress | 15% |
| **10. Terraform Sync** | Planned | 0% |
| **11. Kubernetes Inventory** | Planned | 0% |
| **12. Notifications (Slack/Telegram)** | Planned | 0% |
| **13. REST API Interface** | In Progress | 40% |
| **14. Web Dashboard** | Planned | 0% |
| **15. Team & Org Management** | Planned | 0% |
| **16. Plugin Architecture** | Planned | 0% |
| **17. AI Assistant** | Future | 0% |
| **18. Infrastructure Graph** | Future | 0% |
| **19. Model Context Protocol (MCP) Server** | Future | 0% |

---

## Detailed Task Board

### 1. SSH Management
- [x] **Add Server:** Register host credentials via interactive shell or configuration commands.
- [x] **Remove Server:** Safely delete server metadata and keys.
- [x] **Password Authentication:** Store and authenticate sessions using host passwords.
- [x] **SSH Key Management:** Load private keys for authentication.
- [x] **Interactive Shell:** Connect to remote servers with interactive shell access.
- [x] **Dynamic Terminal Resizing:** Propagate host resize events (`SIGWINCH`) to remote servers.
- [x] **Test Connection:** Doctor command to test host connections and settings.
- [ ] **Configure management:**
  - [ ] Rename server
  - [ ] Clone server profile
  - [ ] Duplicate configuration
  - [ ] Import/Export SSH config (`~/.ssh/config` sync)
  - [ ] Favorite / Recent server lists & connection history
- [ ] **Proxy & Tunneling:**
  - [ ] Jump host / Bastion / ProxyJump support
  - [ ] SOCKS proxy & Dynamic port forwarding (Local/Remote forwarding)
- [ ] **SSH Utilities:**
  - [ ] Non-interactive command execution
  - [ ] Parallel execution / Broadcast commands
  - [ ] Session recording & replay
  - [ ] SSH key generator & rotation
  - [ ] Keep-alive & Auto-reconnect on drop

---

### 2. Infrastructure Inventory
- [x] **Host Metadata Collection:**
  - [x] Hostname & Public IP detection
  - [x] Cloud provider auto-resolution (AWS, GCP, Azure, DigitalOcean, Hetzner, Linode, Vultr, OVH)
- [ ] **OS Details:**
  - [x] Distribution & version details
  - [ ] Timezone, locale, init system, and active package managers
- [ ] **Hardware Specs Collection:**
  - [ ] CPU core count & architecture
  - [ ] RAM, Swap, Disk layouts, and GPU specifications
  - [ ] Live temperature & Network link speed
- [ ] **Installed Software Auto-Detection:**
  - [ ] Detect runtime engines (Docker, Podman, Kubernetes, Node, Go, Rust, Python)
  - [ ] Detect web/proxy services (Caddy, Nginx, Apache)
  - [ ] Detect active databases (PostgreSQL, MySQL, Redis, MongoDB)
  - [ ] Detect message brokers & search (RabbitMQ, Kafka, Elasticsearch)
- [ ] **Containers & Services:**
  - [ ] List running images, containers, networks, and Docker Compose projects
  - [ ] Audit Systemd services (enabled, disabled, failed)
- [ ] **Network Audit:**
  - [ ] Interfaces, Routing tables, DNS config, and open listening ports
- [ ] **Search & Taxonomy:**
  - [x] Tags support (e.g. `production`, `database`)
  - [ ] Express CLI inventory search (e.g. `vpsm search docker`)

---

### 3. Lightweight Monitoring (CGO-Free)
- [ ] **Live Host Metrics:** Gather CPU, RAM, Disk, Swap, Load, and Processes.
- [ ] **Metrics History:** Store hour/day/week/month metrics locally.
- [ ] **Trigger Alerts:** Notify on CPU > 90%, Disk > 80%, Service down, or SSH connection failure.
- [ ] **Service & Container Status:** Check health states for Systemd services and Docker containers.

---

### 4. Documentation & Operational Notes
- [x] **Core Notes System:** Store runbooks, owner details, and notes per server.
- [ ] **Rich Formatting:** Render Markdown content (Mermaid diagrams, images, and code blocks).
- [ ] **File Attachments:** Link diagrams, PDFs, and screenshots to servers.
- [ ] **Server Changelogs:** Record chronological changes (e.g. "Installed Docker").

---

### 5. DevOps Utilities
- [ ] **Container Control:** CLI wrappers for `logs`, `restart`, `prune`, and `compose up/down`.
- [ ] **Service Management:** CLI wrappers for systemctl (`status`, `restart`, `logs`).
- [ ] **Log Tailer:** Tail and grep `journalctl` and docker logs.
- [ ] **Upgrade Audits:** Show pending updates (`apt`, `dnf`, `apk`, `brew`).
- [ ] **Backups & Secret Environment:**
  - [ ] Schedule database/volume backups
  - [ ] Manage `.env` secrets and environment variables

---

### 6. Security Engine
- [ ] **SSH Audit:** Warn on password authentication, root login, or weak ciphers.
- [ ] **Firewall Controls:** Manage configurations for UFW, firewalld, or iptables.
- [ ] **Authorized Keys Inventory:** Audit active SSH public keys on remote servers.
- [ ] **Vulnerability Audits:** Compare installed packages against known CVE databases.

---

### 7. Automation & Workflows
- [ ] **Event Scheduler:** Automate daily scans, audits, and backups.
- [ ] **IF/THEN Rules Engine:** Trigger alerts (e.g., IF Disk > 90% THEN alert).
- [ ] **Webhooks:** Hook outgoing events to server actions.

---

### 8. Events Timeline
- [ ] **System Event Logging:** Log connections, disconnects, alerts, and deployments.
- [ ] **Server Timeline View:** Check server history in chronological order.

---

### 9. Cloud Provider Integrations
- [x] **Automatic Cloud Provider Detection:** (Completed via IMDS, reverse DNS, and ASN lookups).
- [ ] **Cloud Discovery:** Automatically import instances from AWS, GCP, Azure, DO, etc.
- [ ] **Metadata Sync:** Import instance types, regions, volumes, and elastic IPs.

---

### 10. Terraform Sync
- [ ] **State Import:** Parse Terraform state files and auto-import servers/outputs.

---

### 11. Kubernetes Inventory
- [ ] **Cluster Resources:** Monitor pods, nodes, deployments, and Helm configurations.

---

### 12. Notifications System
- [ ] **Integrations:** Support notifications to Slack, Telegram, Discord, Email, and custom Webhooks.

---

### 13. REST API Interface
- [x] **Server Management API:** REST endpoints to add, list, and delete servers.
- [ ] **Extensive REST Coverage:** Add endpoints for notes, monitoring, events, and script runners.

---

### 14. Web Dashboard
- [ ] **Management UI:** Modern dashboard for servers, logs, deployments, and monitoring metrics.

---

### 15. Team & Org Management
- [ ] **Access Control:** User roles, permissions, projects, organizations, and audit logs.

---

### 16. Plugin Architecture
- [ ] **Custom Addons:** Extensibility framework to load external plugins (Cloudflare, AWS, etc.).

---

### 17. AI Assistant (Future)
- [ ] **Reasoning Engine:** Answer operational questions (e.g., "Why is production slow?").

---

### 18. Infrastructure Graph (Future)
- [ ] **Dependency Graph:** Visual map of dependencies (Customer -> Server -> Containers -> Databases).

---

### 19. Model Context Protocol (MCP) Server (Future)
- [ ] **AI-Agent Connection:** Expose safe operational APIs to AI coding assistants and agents.
