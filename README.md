# Spec Framework

**Specification Driven Development for AI agents** — a method and a CLI that turn product, specification, planning, execution, and audit into one traceable pipeline, instead of asking an AI to "implement a feature" from loose context.

📖 **[Read the visual guide](https://jonatasfreiredev.github.io/spec-framework/)** · [Canonical method (FRAMEWORK.md)](FRAMEWORK.md) · [Adoption guide](framework/adoption.md)

```text
Problem → Vision → Strategy → Domain → User Goal → Feature → Use Case
   → Specification → Design → Engineering → Plan → Execution Graph
   → Tasks → Code → Validation → Audit
```

Every artifact carries an id, parents, dependencies, and an approval state. Agents propose; humans approve. If an AI cannot explain which problem a task was born from, the task is not ready for implementation.

## Quick Start

One command installs the CLI (checksum-verified) and opens the interactive `init` wizard:

**Windows (PowerShell)**

```powershell
irm https://raw.githubusercontent.com/JonatasFreireDev/spec-framework/master/scripts/init.ps1 | iex
```

**Linux / macOS**

```bash
curl -fsSL https://raw.githubusercontent.com/JonatasFreireDev/spec-framework/master/scripts/init.sh | sh
```

The initialized repository receives **only `product/`** — the method (skills, templates, validators) lives in a versioned user cache, pinned by `product/.product/framework.json`. No `.spec-framework/`, no generated skill trees, no CI files polluting your repo.

Already have the CLI? Then:

```bash
spec-framework init ../my-product --agents codex --yes   # new product
spec-framework init ../my-product --agents codex \
  --starting-point existing-documents --source-dir ../product-docs --yes   # import PRDs/epics
```

The starting point selects the active entry contract and gate:

| Starting point | Entry contract | Gate before `work` |
| --- | --- | --- |
| `new-product` | Problem → Vision → Principles/North Star → Strategy | Full Foundation approved |
| `existing-product` | Product Baseline → Strategy | Both individually approved |
| `existing-documents` | Latest import run | Selected mappings materially complete; outputs remain draft |
| `existing-feature` | Feature Brief with one `Target Feature` | Current approval matching the selected feature |
| `existing-implementation` | Implementation Assessment → full Foundation | Assessment and Foundation approved |
| `audit-only` | Read-only inspection | Mutations and workspace creation blocked |

Read the generated `product/BOOTSTRAP.md`; it names the current gate and exact next command. Every starting point that creates or revises domains uses the installed `examples/events/` reference before its first domain change to model a business area with explicit ownership, non-ownership, dependencies, and one Domain -> User Goal -> Feature -> Use Case walking skeleton; `audit-only` uses it for assessment only. The complete ladder above is the default `new-product` path, while proportional starting points rejoin it at their governed handoff.

## How It Works

| Concept | In one line |
| --- | --- |
| **Specification** | The central contract: flow, UI, APIs, data, permissions, analytics, security, and acceptance — written before any code. |
| **Design flow** | Specification → `design.md` (origin `generate`/`evolve`/`adopt`, versioned visual sources) → independent UX Review → human gate. |
| **Engineering flow** | Technical Discovery → Engineering Proposal → independent Engineering Review → Implementation Plan. |
| **Execution Graph** | Tasks as a DAG with explicit `writeScope`; parallel work never overlaps write paths. |
| **Approval gates** | `draft → proposed → approved → in_progress → implemented → validated → released`, each transition mechanically checked (content hashes, same-diff QA + Code Review, staleness detection). |
| **Shared systems** | Optional versioned Design System (`design/system/`) and Engineering System (`engineering/`) — pinned per delivery, never self-approving. |

## CLI At a Glance

The [Framework Guide skill](framework/skills/framework-guide/SKILL.md) is the default conversational front door when no verified specialist route exists — describe your goal and it recommends the smallest safe command. Current CLI guidance or an explicit human request naming both specialist and scope can route directly; persisted handoffs/checkpoints must first be revalidated with `dashboard`, `status`, `next`, or `guide`. The commands it routes to:

| Intent | Commands |
| --- | --- |
| Start or import a product | `init`, `import materialize` |
| Navigate and see state | `work`, `status`, `next`, `dashboard`, `guide` |
| Approve explicitly | `review`, `approve`, `approve-stage`, `gates` |
| Design workflows | `design init/import/register/map/inspect/audit`, `design-system init/inspect/validate/migrate` |
| Engineering workflows | `engineering-system inspect/validate/triggers/migrate` |
| Operate the graph | `graph ready/materialize/claim/complete`, `task readiness` |
| Governed execution | `resume`, `lease claim`, `schedule`, `commands plan` |
| Visual adapters (optional) | `adapters list/status/doctor/install` — version-pinned, preview first |
| Inspect decisions | `impact`, `decisions migrate` |
| Keep healthy | `validate`, `upgrade`, `migrate external-runtime`, `skill path <skill>` |

All mutations preview before executing; approval commands require an explicit human identity and `--yes`.

## Repository Map

| Area | Purpose |
| --- | --- |
| [FRAMEWORK.md](FRAMEWORK.md) | Canonical method and architecture. |
| [framework/](framework/) | Executable framework core: skills, templates, validators, FDRs. |
| [framework/skills/](framework/skills/) | 50 specialist and orchestrator skill contracts. |
| [starter/](starter/) | Clean `product/` skeleton copied into new repositories. |
| [examples/events/](examples/events/) | Worked product instance used as learning material and validation fixture. |
| [docs/](docs/) | Published visual guide (GitHub Pages). |
| [cmd/](cmd/) · [internal/](internal/) | The Go CLI. |

**Boundary rule:** framework method assets live in this repository and ship inside the CLI binary; an adopter repository owns only `product/`. See [framework/adoption.md](framework/adoption.md).

## Development

```bash
go test ./...
go vet ./...
go test -race ./...
go run ./cmd/spec-framework validate --product-root examples/events --framework-root . --write-registry --write-report
```

Method changes are recorded as [Framework Decision Records](framework/decisions/) (`FDR-*`).

## License

[MIT](LICENSE)
