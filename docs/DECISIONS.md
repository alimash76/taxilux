# Decisions

Every dependency below is justified in one line, per AGENTS.md rule 6.

| Dependency | Reason |
|---|---|
| Go stdlib `net/http` | No heavy framework; stack-mandated in AGENTS.md |
| `modernc.org/sqlite` | Pure-Go SQLite driver, no CGO, simple static image |
| Astro | Static output, simple NL/EN i18n, shared layout; stack-mandated |
| nginx | Serve frontend static output in its own container; stack-mandated |
| Caddy | Reverse proxy with automatic HTTPS; stack-mandated |
| Docker Compose | Orchestration of backend, frontend, proxy; stack-mandated |
| Mailpit (dev only) | Local SMTP catcher for email development; stack-mandated |
| tini | Init process / signal forwarding in backend container |

## Proposed (not yet decided)

- Two calm colour palettes for Phase 1 (see docs/BACKLOG.md)
- Three "Services" bullets for Phase 1 (see docs/BACKLOG.md)
