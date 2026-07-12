# Product Engineering Framework v2

This repository is the framework laboratory and reusable base for Specification Driven Development.

## Start Here

| Area | Purpose |
| --- | --- |
| [FRAMEWORK.md](FRAMEWORK.md) | Canonical method and architecture. |
| [framework/](framework/) | Framework core boundary and adoption guide. |
| [starter/](starter/) | Clean product skeleton for new repositories. |
| [examples/](examples/) | Worked examples and learning material. |
| [framework/skills/](framework/skills/) | Operational Codex skills. |
| [framework/template/](framework/template/) | Reusable artifact templates. |
| [framework/validators/](./framework/validators) | Mechanical validation gates. |

## Adoption

For a new product repository, start from [starter/](starter/) rather than copying this repository root.

The starter creates two explicit roots:

```text
.spec-framework/  # how the framework works
product/          # the product being built
```

Current recommended flow:

```text
spec-framework init --target ../my-product --agents codex --yes -> fill product/foundation -> create product/domains -> spec-framework validate
```

Local CLI form:

```bash
go run ./cmd/spec-framework init --target ../my-product --agents codex,cursor,claude --yes

# Bootstrap from existing epics, PRDs, or document directories
go run ./cmd/spec-framework init --target ../my-product --agents codex --starting-point existing-documents --source-dir ../product-docs --yes

# After Artifact Importer proposes mappings and a human approves them
spec-framework import materialize --run IMPORT-001 --approved-by "Product Owner" --yes

# Select and navigate one feature without creating a global active feature
spec-framework work --feature FT-001 --domain events --goal manage-event --created-by "Product Owner"
spec-framework work --feature FT-001 --use-case send-invitation --created-by "Product Owner"
spec-framework status --work WORK-001
spec-framework next --work WORK-001

# Record an explicit approval and inspect implementation readiness
spec-framework approve --artifact domains/events/context.md --grant approved --approved-by "Product Owner" --yes
spec-framework gates

# Operate an approved execution graph without automatically running agents
spec-framework graph ready --graph domains/events/goals/manage/features/invites/use-cases/send/execution-graph.json
spec-framework graph materialize --graph domains/events/goals/manage/features/invites/use-cases/send/execution-graph.json --yes
spec-framework task readiness --graph domains/events/goals/manage/features/invites/use-cases/send/execution-graph.json --task TK-001
spec-framework guide --work WORK-001
spec-framework review --work WORK-001 --stage tasks
spec-framework impact --decision DEC-021

# Resume, lease, schedule, and execute governed local commands
spec-framework resume --work WORK-001
spec-framework lease claim --work WORK-001 --graph domains/events/goals/manage/features/invites/use-cases/send/execution-graph.json --task TK-001 --agent codex --isolate
spec-framework schedule --work WORK-001 --graph domains/events/goals/manage/features/invites/use-cases/send/execution-graph.json --max-parallel 4
spec-framework commands plan --work WORK-001 --task TK-001 --risk R0 -- go test ./...
spec-framework validate
spec-framework upgrade --target ../my-product --agents codex,cursor,claude --yes
```

Release CLI installation:

```bash
Download the archive for Windows, Linux, or macOS from GitHub Releases, verify `checksums.txt`, and place `spec-framework` on `PATH`.
```

See [framework/adoption.md](framework/adoption.md).

## Ladder

```text
Problem -> Vision -> Strategy -> Domain -> User Goal -> Feature -> Use Case -> Specification -> Design -> Implementation Plan -> Execution Graph -> Tasks -> Code -> Validation -> Audit
```

## Repository Boundary

| Concern | Lives In |
| --- | --- |
| Framework method, validators, skills, templates, and FDRs | This repository core |
| New product state and product scope | `product/` in a repo created from `starter/` |
| Installed framework method assets | `.spec-framework/` in a repo created from `starter/` |
| Example domains and use cases | `examples/` |

## Quality Gates

Run:

```bash
go test ./...
go vet ./...
go test -race ./...
go run ./cmd/spec-framework validate --product-root examples/events --framework-root . --write-registry --write-report
```
