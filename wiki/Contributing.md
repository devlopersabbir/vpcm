# Contributing to VPSM

Thank you for your interest in contributing to VPSM! We welcome contributions from developers of all skill levels. To maintain code quality and long-term project viability, we follow clean-code principles and a structured development workflow.

---

## Technical Stack & Prerequisites

- **Go Version:** Go 1.25.x or higher.
- **Database Backend:** Local SQLite or a running MongoDB instance.
- **Tools:** `make`, `git`, and standard terminal utilities (`zip`, `tar` for releases).

---

## Project Structure & Architecture

VPSM follows **Clean Architecture** patterns to decouple business logic from framework dependencies:

- **Service Interfaces (`internal/inventory/interfaces.go`):** Define what operations are possible. All execution code must bind to interfaces rather than concrete structs.
- **Pluggable Repositories (`internal/inventory/`):** Both `mongoServerRepository` and `sqliteServerRepository` satisfy the same database interface. If you add database features, make sure they are implemented in both SQLite and MongoDB backends.
- **Event Bus Decoupling (`internal/events/`):** Do not call modules directly from other modules. Use `events.Publish()` to communicate changes asynchronously (e.g., triggering background inventory checks).

---

## Local Development Workflow

### 1. Build and Run

We use a simple `Makefile` to automate builds and configuration setup.

```bash
# Compile and build binaries locally
make build

# Run formatting check
make lint

# Run the test suite
make test

# Install vpsm/vpcm binaries to /usr/local/bin
make install
```

### 2. Testing Configuration

To test changes against config files, you can initialize a local test workspace and run:

```bash
./bin/vpsm config init
```

---

## Coding Standards & Guidelines

- **Zero-CGO Policy:** To ensure seamless cross-compilation across platforms (macOS, Windows, Linux), we only accept pure-Go dependencies. (e.g., we use `modernc.org/sqlite` instead of `go-sqlite3`).
- **Semantic Commits:** We enforce structured commit messages to auto-generate changelogs:
  - `feat:` for new capabilities.
  - `fix:` for bug fixes.
  - `docs:` for documentation updates.
  - `chore:` for general cleanup/formatting/dependencies.
  - `ci:` for CI/CD pipeline modifications.

---

## Submitting Pull Requests

1. **Fork the Repo:** Create a branch for your feature or bug fix (`git checkout -b feature/cool-new-thing`).
2. **Implement & Test:** Ensure all local tests pass and standard Go formatting is applied (`make lint` and `make test`).
3. **Changelog:** If adding a user-facing feature, update `CHANGELOG.md` under the unreleased changes block.
4. **Push & PR:** Push your branch and open a Pull Request against the `main` branch.
