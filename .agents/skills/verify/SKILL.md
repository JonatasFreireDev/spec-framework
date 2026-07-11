---
name: verify
description: Run the complete Go mechanical gate suite after changes to the CLI, framework assets, starter, or worked example.
---

# Verify Skill

## Purpose

Run every mechanical gate defined by the Go repository and report exact failures.

## Gates

Run from the repository root:

```bash
gofmt -w assets.go cmd internal
go test ./...
go vet ./...
go run ./cmd/spec-framework validate --product-root examples/events --framework-root .
git diff --check
```

For concurrency or release changes, also run `go test -race ./...` on Linux CI. For packaging changes, cross-build Windows, Linux, and macOS on amd64 and arm64, then run the release smoke skill.

## Reporting

Report each command, status, and exact failing package or validator finding. Never hide a skipped sandbox or external gate.
