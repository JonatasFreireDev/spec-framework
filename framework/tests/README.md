# Engineering Tests

## Purpose

Go package tests create temporary repositories, execute real framework components, and remove fixtures afterward.

## Current Coverage

| Tool | Coverage |
| --- | --- |
| `internal/validator` | Gates, deterministic diagnostics, reports, registry, and Node parity fixtures. |
| `internal/moveartifact` | Planning, rollback, link/JSON rewrites, and mention reports. |
| `internal/install` | Init, upgrade, embedded assets, manifests, and multi-agent skills. |
| `internal/cli` | Command dispatch and end-to-end CLI flow. |
| `internal/wizard` | Bubble Tea state transitions and confirmation. |

## Run

```bash
go test ./...
```

Run these tests before changing validator gates, identity policy, staleness behavior, approval-record behavior, or artifact movement behavior.

## Next Step

Add fixtures for task-file validation, rigor-tier gates, Mermaid semantic bindings, anchor validation, bootstrap upgrade behavior, and future Phase B writeScope errors as those areas evolve.
