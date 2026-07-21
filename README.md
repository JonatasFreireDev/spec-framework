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

Install the checksum-verified CLI for your operating system. Installation does not create or modify a product repository.

**Windows (PowerShell)**

```powershell
irm https://raw.githubusercontent.com/JonatasFreireDev/spec-framework/master/scripts/install.ps1 | iex
```

**Linux / macOS**

```bash
curl -fsSL https://raw.githubusercontent.com/JonatasFreireDev/spec-framework/master/scripts/install.sh | sh
```

When you are ready to prepare a product, ask the agent to inspect the complete repository, identify every implementation root and semantic role, and initialize the product. The agent should execute an explicit command such as:

```bash
spec-framework init ../my-product --agents codex --code-roots web:web,services/api:api --yes
```

If inspection confirms that no implementation exists, the agent uses `--no-code-roots`. The interactive wizard remains available for manual use, but omitted root information invokes only a CLI fallback marked for agent review; it cannot unlock a Specification. The initialized repository receives **only `product/`** — the method (skills, templates, validators) lives in a versioned user cache, pinned by `product/.product/framework.json`. No `.spec-framework/`, generated skill trees, or CI files are added to the adopter repository.

For automation, the same choices are available through explicit flags. Run `spec-framework init --help` for the current options. A fallback root map can be corrected without overwriting product content through `spec-framework upgrade --code-roots path:role,... --yes`.

### Choose the right starting point

A starting point describes the evidence available today and the first contract that must be reviewed. It does not represent the application's technical entrypoint, approve existing material, or remove later gates.

Each choice resolves a versioned declarative contract that selects the required directories, files, registry entries, deterministic patches, and typed actions:

| Starting point | Use when | First path |
| --- | --- | --- |
| `new-product` | There is an idea or opportunity, but the product still needs definition | Problem → Vision → Principles/North Star → Strategy |
| `existing-product` | A real product is operating with users, releases, metrics, or support evidence | Product Baseline → Strategy |
| `existing-implementation` | Code exists, but its product intent or operating history is unclear | Implementation Assessment → full Foundation |
| `existing-documents` | PRDs, Jira, Confluence, wikis, or other documents are the main source | Inventory → mapping → conflict review → draft materialization |
| `existing-feature` | One small, well-bounded delivery is already understood | Feature Brief → shared baselines → target Feature |
| `audit-only` | The goal is to identify gaps without changing product state | Inspect → validate → report gaps |

Read the generated `product/BOOTSTRAP.md`; it names the current gate and next action. See the [starting-point guide](docs/starting-points.md) for examples, comparisons, and selection guidance.

### Manage the CLI

CLI lifecycle and product lifecycle are separate:

```bash
spec-framework update --check        # check for a CLI release
spec-framework update --yes          # update the CLI binary
spec-framework upgrade --yes         # update this product's pinned runtime
spec-framework uninstall             # preview local CLI removal
spec-framework uninstall --purge --yes
```

`uninstall` never searches for or removes product repositories. `--purge` additionally removes only the versioned runtime cache and namespaced Spec Framework dispatchers.

## How It Works

| Concept | In one line |
| --- | --- |
| **Specification** | A bounded root synthesis plus rigor-selected product, behavior, UX, API, data, security, quality, observability, and rollout contracts. New specs use semantic Contract v2 gates; existing specs migrate explicitly without `upgrade` rewriting adopter content. |
| **Design flow** | Specification → `design.md` (origin `generate`/`evolve`/`adopt`, versioned visual sources) → independent UX Review → human gate. |
| **Engineering flow** | Shared baseline: Engineering Orchestrator → Technical Landscape → Standards + Operations → Evidence → Engineering System approval. It runs sequentially by default or through bounded minimal-context native subagents. Delivery: Technical Discovery → Engineering Proposal → independent Engineering Review → Implementation Plan. |
| **Execution Graph** | Complete vertical task contracts as a DAG with explicit `writeScope`; parallel work never overlaps write paths. |
| **Approval gates** | `draft → proposed → approved → in_progress → implemented → validated → released`, each transition mechanically checked (content hashes, same-diff QA + Code Review, staleness detection). |
| **Shared systems** | Optional versioned Design System (`design/system/`) and scalable Engineering System (`engineering/`) whose graph, standards, operations, evidence, aggregate, and quality contracts have explicit specialist owners — pinned per delivery, never self-approving. |

## CLI At a Glance

The [Framework Guide skill](framework/skills/framework-guide/SKILL.md) is the default conversational front door when no verified specialist route exists — describe your goal and it recommends the smallest safe command. Current CLI guidance or an explicit human request naming both specialist and scope can route directly; persisted handoffs/checkpoints must first be revalidated with `dashboard`, `status`, `next`, or `guide`. The commands it routes to:

| Intent | Commands |
| --- | --- |
| Start or import a product | `init`, `import create/status/resume`, `import materialize` |
| Evolve an existing delivery | Read `context.md`, classify the demand, then route through `evolution` and the owning Feature/Use Case/Specification skill |
| Navigate and see state | `work`, `status`, `next`, `dashboard`, `guide` |
| Approve explicitly | `review`, `approve`, `approve-batch`, `approve-stage`, `gates` |
| Design workflows | `design init/import/register/map/inspect/audit`, `design-system init/inspect/validate/migrate` |
| Engineering workflows | `engineering-system inspect/validate/triggers/migrate` |
| Operate the graph | `graph ready/materialize/claim/complete`, `task readiness` |
| Governed execution | `resume`, `lease claim`, `schedule`, `commands plan` |
| Visual adapters (optional) | `adapters list/status/doctor/install` — version-pinned, preview first |
| Inspect decisions | `impact`, `decisions migrate` |
| Keep healthy | `validate`, `update`, `upgrade`, `migrate external-runtime`, `skill path <skill>` |
| Manage the local CLI | `update --check`, `update --yes`, `uninstall`, `uninstall --purge --yes` |

The CLI uses Cobra for its command tree and generated help. It deliberately does not load ambient user configuration: product manifests and explicit flags remain the source of truth.

All mutations preview before executing; approval commands require an explicit human identity and `--yes`. For example: `spec-framework approve-batch --foundation` previews the Foundation scope, and `spec-framework approve-batch --foundation --approved-by "Product Owner" --yes` applies it after human confirmation. Use `--all-eligible --until specification` to include the ordered Foundation, Domains, Feature, Use Case, and Specification stages. Direct terminal users may use `--interactive` for Bubble Tea confirmation; CI should use textual or `--json` output.

## Repository Map

| Area | Purpose |
| --- | --- |
| [FRAMEWORK.md](FRAMEWORK.md) | Canonical method and architecture. |
| [framework/](framework/) | Executable framework core: skills, templates, validators, and tools. |
| [framework/skills/](framework/skills/) | 49 skill folders, plus shared runtime contracts. |
| [starter/](starter/) | Canonical source assets selected by declarative initialization contracts. |
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

Method changes update `FRAMEWORK.md`, affected skills, templates, validators, documentation, and tests together. Git history is the maintenance record.

## License

[MIT](LICENSE)
