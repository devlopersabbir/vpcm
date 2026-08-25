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

### `vpsm server export`

Writes the whole inventory to stdout, or to a file with `--out` / `-o`.

```bash
vpsm server export --format <ssh|json|csv|yaml> [--out <file>]
```

_Example:_

```bash
vpsm server export -f json -o ~/backups/servers.json
```

### `vpsm server import`

Reads servers back in from a file produced by `vpsm server export`, or from an SSH config. The input path is taken from the argument or `--in` / `-i`, and defaults to stdin.

```bash
vpsm server import [file] [--format <ssh|json|csv|yaml|auto>] [--in <file>] [--on-conflict <skip|overwrite|rename|fail>] [--dry-run]
```

The format is detected from the file extension and contents unless `--format` says otherwise. Existing servers are matched by UUID and then by name, and `--on-conflict` decides whether those matches are skipped (default), updated, imported under a suffixed name, or treated as a fatal error.

_Example:_

```bash
vpsm server import ~/backups/servers.json --on-conflict overwrite
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
