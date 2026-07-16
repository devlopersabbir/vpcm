# VPSM CLI Usage Guide

The `vpsm` CLI tool is the main interactive command-line interface for managing virtual private servers (VPS), checking configurations, diagnosing issues, and reviewing system integrity.

## Global Flags & Environment Overrides

The CLI respects standard configurations with the following priority order:

1. **CLI Flags** (e.g. `--config`)
2. **Environment Variables** (prefixed with `VPSM_`, e.g. `VPSM_LOGGING_LEVEL=debug`)
3. **Configuration File** (`~/.config/vpsm/config.yaml`)
4. **Hardcoded Defaults**

---

## Command Reference

### `vpsm doctor`

Performs a diagnostic check of the configuration files, database integrity, GORM migrations, and logger initialization. Use this command to verify if VPSM is correctly installed and ready.

```bash
vpsm doctor
```

### `vpsm version`

Prints the current version of the VPSM software.

```bash
vpsm version
```

### `vpsm config`

Prints the active configuration properties showing loaded database paths, API states, and logging priorities.

```bash
vpsm config
```

### `vpsm server add`

Registers a new host server into the database.

```bash
vpsm server add <name> <host>
```

_Example:_

```bash
vpsm server add web-prod 192.168.1.10
```

### `vpsm server list`

Lists all servers registered in the local inventory.

```bash
vpsm server list
```

---

## Configuration Reference (`~/.config/vpsm/config.yaml`)

```yaml
database:
  driver: sqlite
  path: ~/.local/share/vpsm/vpsm.db

logging:
  level: info
  format: pretty # Options: pretty, json

api:
  enabled: false
  host: 127.0.0.1
  port: 8080

ssh:
  timeout: 10s

collector:
  workers: 5

plugins:
  enabled: true
```
