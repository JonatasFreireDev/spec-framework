# Agent Instructions

## Purpose

This repository is the laboratory, implementation, validation fixture, and distribution source for the Spec Framework. Agents working here maintain the framework rather than adopting it as an ordinary product.

Always read `FRAMEWORK.md` before changing framework behavior.

## Repository Boundaries

| Area | Responsibility |
| --- | --- |
| `FRAMEWORK.md` | Canonical framework method and architecture. |
| `framework/` | Framework-owned skills, templates, validators, adapters, audits, engineering guidance, and FDRs. |
| `starter/` | Clean skeleton copied into new adopter repositories. |
| `examples/events/` | Worked product fixture used for learning and validation; never use it as the starter. |
| `cmd/`, `internal/`, `assets.go` | Go CLI, embedded assets, installer, runtime, workflow, and validator implementation. |
| `.agents/`, `.claude/`, `.codex/` | Repository-maintenance skills and agent integrations. |

Do not encode Events product scope into reusable framework assets. Do not treat `examples/events/` as the clean source for new products.

## Sources Of Truth

- Method and gates: `FRAMEWORK.md`.
- Framework decisions: `framework/decisions/FDR-*`.
- Canonical shipped skills: `framework/skills/`.
- Canonical templates: `framework/template/`.
- Adopter skeleton: `starter/`.
- CLI behavior: `cmd/spec-framework/` and `internal/`.
- Worked product state: `examples/events/` and its `.product/` metadata.

Generated agent trees are derived copies. Do not edit them as the canonical source of framework skills.

## Maintenance Skills

Use repository-local maintenance skills when applicable:

- `fdr`: framework-method decisions.
- `new-framework-skill`: new or normalized shipped skills.
- `sync-framework-assets`: synchronization across embedded assets and agent targets.
- `verify`: complete mechanical gate suite.
- `release-publisher` and `release-smoke`: approved releases.

Framework specialist and orchestrator skills live under `framework/skills/`. When modifying their contracts, follow the repository skill scaffolding, registration, handoff, and validation rules.

## Change Rules

1. Record changes to method, skill contracts, validator behavior, gates, or delivery workflow as an FDR unless the change is purely editorial synchronization.
2. Keep `FRAMEWORK.md`, FDRs, skills, templates, validators, starter assets, examples, installer behavior, and tests synchronized for the affected surface.
3. Preserve adopter-owned product content during `upgrade`; never solve an upgrade by overwriting product scope or approval history.
4. Keep FDRs in `framework/decisions/` and product decisions in the active product root's `knowledge/decisions/`.
5. Do not create, edit, or repair product approval records unless the human explicitly authorizes a migration that names approval-record generation.
6. Use `spec-framework move` for artifact moves governed by the framework and review reported free-text mentions.
7. Do not implement application code as part of framework documentation bootstrap, planning, readiness, or maintenance work.
8. Preserve user changes in a dirty worktree and avoid unrelated rewrites.

## Product Fixture Work

When modifying product artifacts under `examples/events/`, treat it like an adopter product:

- read relevant parent and local `context.md` files;
- use the owning skill and matching template;
- respect parent approvals, rigor tiers, Delivery Level, Priority, decisions, derivations, and staleness;
- do not advance statuses or fabricate evidence;
- validate with `--product-root examples/events --framework-root .`.

The detailed product delivery contract is canonical in `FRAMEWORK.md`; do not duplicate it here.

## Distribution And Synchronization

The released CLI embeds `FRAMEWORK.md`, `starter/`, framework agent instructions, decisions, skills, and templates. When changing shipped assets:

1. verify the embedded asset boundary in `assets.go`;
2. verify `init` output for applicable Codex, Cursor, and Claude targets;
3. verify `upgrade` preserves adopter-owned files;
4. update installation/adoption documentation and release smoke coverage when needed;
5. keep target-specific extensions out of unsupported agent trees.

## Verification

Run from the repository root:

```bash
gofmt -w assets.go cmd internal
go test ./...
go vet ./...
go run ./cmd/spec-framework validate --product-root examples/events --framework-root .
git diff --check
```

For release or packaging changes, also run the release smoke workflow and required cross-builds.

## Reporting

Report changed framework surfaces, starter/example synchronization, decisions, validation, compatibility, migration, release impact, and the recommended next owner or command.

Use concise status tables and Mermaid only when they materially improve understanding. Use valid UTF-8 status icons: ✅, 🟡, 🔴, and ➖.
