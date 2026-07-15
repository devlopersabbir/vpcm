# VPSM Architecture (v0.0.1 Foundation)

Designed to endure a 10-20 year active development lifespan.

## Component Layout

```text
               +-----------------------+
               |        CLI/API        |
               +-----------+-----------+
                           |
                           v
               +-----------+-----------+
               |    Service Boundary   |
               +-----------+-----------+
                           |
                           v
               +-----------+-----------+
               | Repository Interface  |
               +-----------+-----------+
                           |
                           v
               +-----------+-----------+
               |   Database (SQLite)   |
               +-----------------------+
```

## Key Architectural Decoupling Rules

1. **Dependency Inversion**: Domain packages define interfaces for service and repository methods.
2. **Event Bus Decoupling**: Inter-module transactions are prohibited. Modules communicate state changes via the Event Bus.
3. **Multi-executable Entrypoints**:
   - `vpsm`: CLI.
   - `vpsmd`: Daemon (scheduler, async queues).
   - `vpsm-api`: HTTP REST API.
