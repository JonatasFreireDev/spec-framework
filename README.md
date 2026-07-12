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

One-command interactive bootstrap on Windows:

```powershell
irm https://raw.githubusercontent.com/JonatasFreireDev/spec-framework/master/scripts/init.ps1 | iex
```

The script resolves a concrete release, verifies its checksum, installs the CLI for the current user, and opens the `init` wizard. The initialized repository receives only `product/`.

The starter adds one explicit root to the adopter repository:

```text
product/          # the product being built
```

Current recommended flow:

During initialization, Impeccable remains optional. Interactive `init` offers an install/skip choice. Headless initialization requires an explicit provider version:

```bash
spec-framework init ../my-product --agents codex --install-impeccable --impeccable-version latest --yes
```

The product is initialized first. If the optional provider installer fails, the CLI reports partial success and exits non-zero without deleting the initialized product.

`latest` is resolved through npm to a concrete semantic version before preview or execution. Use an exact value such as `2.3.2` when a fully reproducible build is required; the provider is never executed as an unresolved `@latest` package.

```text
spec-framework init ../my-product --agents codex --yes -> fill product/foundation -> create product/domains -> spec-framework validate
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
spec-framework dashboard --work WORK-001
spec-framework status --work WORK-001 --graph --json
spec-framework decisions migrate --product-root product
spec-framework decisions migrate --product-root product --interactive

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

## Visual Design Workflows

Design can be generated, evolved from an existing interface, or adopted from versioned Figma, Penpot, image, or other visual sources.

```bash
spec-framework design init --product-root product --use-case <path> --mode generate
spec-framework design import --product-root product --use-case <path> --type images --source <directory> --authority visual-canonical
spec-framework design register --product-root product --use-case <path> --type figma --source <url> --version <version> --nodes <ids> --authority visual-canonical
spec-framework design map --product-root product --use-case <path> --mappings mapping.json
spec-framework design inspect --product-root product --use-case <path>
spec-framework design audit --product-root product --use-case <path> --write-report
```

Impeccable, Figma, and Penpot are optional adapters. Imported or generated assets never approve Design; independent UX Review and the existing human approval gate still apply.

Inspect or install optional adapters explicitly:

```bash
spec-framework adapters list
spec-framework adapters status impeccable
spec-framework adapters doctor impeccable --check-latest
spec-framework adapters install impeccable --version 2.3.2
spec-framework adapters install impeccable --version 2.3.2 --yes
```

The first install command is a preview. `--yes` executes the official version-pinned provider command. Adapter removal is intentionally unavailable until the provider offers a documented reversible command.

For convenience, `--version latest` resolves and displays the current concrete version before execution.

Initialize and inspect the shared Design System when the product has recurring foundations, tokens, components, or patterns:

```bash
spec-framework design-system init --product-root product --mode generate
spec-framework design-system inspect --product-root product
spec-framework design-system validate --product-root product
spec-framework design-system migrate --product-root product --dry-run
```

Inspect the shared Engineering System and list the structured triggers that make Engineering Proposal and Review applicable to Tier S/M:

```bash
spec-framework engineering-system inspect --product-root product
spec-framework engineering-system validate --product-root product
spec-framework engineering-system triggers
spec-framework engineering-system migrate --product-root product --dry-run
```
```

## Repository Boundary

| Concern | Lives In |
| --- | --- |
| Framework method, validators, skills, templates, and FDRs | This repository core |
| New product state and product scope | `product/` in a repo created from `starter/` |
| Installed framework method assets | Versioned user cache resolved by `product/.product/framework.json` |
| Example domains and use cases | `examples/` |

## Quality Gates

Run:

```bash
go test ./...
go vet ./...
go test -race ./...
go run ./cmd/spec-framework validate --product-root examples/events --framework-root . --write-registry --write-report
```
