# Install the Spec Framework CLI

Users install a precompiled `spec-framework` binary. Go and Node.js are not
runtime requirements after the Go cutover.

## Release archives

Download the archive for the operating system and architecture from the GitHub
release, verify it against `checksums.txt`, and place the executable on `PATH`.

| Operating system | Architectures | Executable |
| --- | --- | --- |
| Windows | amd64, arm64 | `spec-framework.exe` |
| Linux | amd64, arm64 | `spec-framework` |
| macOS | amd64, arm64 | `spec-framework` |

Verify the installation:

```bash
spec-framework version
spec-framework help
```

Alternatively, use the checksum-verifying interactive bootstrap:

```powershell
irm https://raw.githubusercontent.com/JonatasFreireDev/spec-framework/master/scripts/install.ps1 | iex
```

Linux and macOS users can download and inspect `scripts/install.sh`, then execute it locally. Piping a remote shell script is not required. Installation verifies the release checksum and prints the installed version; it does not run product initialization.

Manage the installed CLI independently from adopter products:

```bash
spec-framework update --check
spec-framework update --yes
spec-framework update --version 0.4.0 --yes
spec-framework uninstall
spec-framework uninstall --yes
spec-framework uninstall --purge --yes
```

The installer writes `install.json` beside the binary so lifecycle commands can identify installer-owned paths. `update` verifies the checksum and the candidate binary before changing the CLI. Set `GITHUB_TOKEN` when authenticated GitHub API access is required. `upgrade` changes only a product's pinned runtime, manifest, and dispatchers. `uninstall` never searches for or removes `product/`; `--purge` additionally removes framework caches and the namespaced Spec Framework dispatchers while preserving other agent skills.

Initialize without a TTY. The target may already contain product code; only an existing `product/` blocks initialization:

```bash
spec-framework init ../my-product --agents codex,cursor,claude --yes
```

| Starting point | Materialized entry contract |
| --- | --- |
| `new-product` | Full Foundation starter |
| `existing-product` | `foundation/product-baseline.md` plus Strategy |
| `existing-documents` | Immutable import run pinned in the product manifest |
| `existing-feature` | `foundation/feature-brief.md` with `Target Feature` |
| `existing-implementation` | `knowledge/assessments/implementation-assessment.md` plus full Foundation |
| `audit-only` | Read-only bootstrap; mutating CLI commands are refused |

Use `--starting-point existing-documents` with `--source-dir` or `--sources` to bootstrap from existing product material. This creates a scalable analysis-only import run with paged inventory and review chunks; use `--import-max-files`, `--import-max-total-bytes`, `--import-max-file-bytes`, and `--import-chunk-size` to set the explicit bootstrap budget. It never treats source prose as approved product truth. For a new demand, the Artifact Importer must compare the source with existing Features and Use Cases, propose a relation and destination, and preserve source traceability before materialization.

Before `init`, the agent inspects the complete repository and declares the semantic implementation map with `--code-roots <path:role,...>`, for example `--code-roots web:web,services/api:api,packages/core:library`. If the inspection confirms no implementation, use `--no-code-roots`. The CLI still detects common top-level markers when neither option is supplied, but records that result as `cli-fallback` / `needs-agent-review`; fallback candidates cannot unlock a Specification. Correct them later with `upgrade --code-roots ... --yes` or `upgrade --no-code-roots --yes`. The confirmed result is recorded in the product manifest and seeds `knowledge/assessments/product-landscape.md`. Product Landscape, Engineering System, and the shared Design System remain mandatory pre-Specification baselines; route Engineering System creation through `engineering-orchestrator` and its specialist owners. Legacy manifests remain compatible until they opt into this policy.

When an already policy-enabled manifest predates discovery provenance, `upgrade` marks it as `legacy-unclassified` / `needs-agent-review` without changing its roots or Product Landscape. The agent must reinspect and confirm the map with one of the same upgrade flags before a new Specification can advance.

For a large document set, create a bounded scalable run with `import create`, inspect it with `import status`, and resume one bounded chunk at a time. Draft materialization remains a separate explicit command after sources/chunks and mappings have been reviewed: `spec-framework import materialize --run IMPORT-001 --approved-by "Product Owner" --yes`.

The generated `product/BOOTSTRAP.md` is rendered from the starting-point map in `framework/init/bootstrap.json`; each step names the user goal, agent reading set, writable draft paths, prompt, gate, and next handoff. For imports it also pins the active run id. `spec-framework work` remains blocked until that latest run records explicit materialization approval and at least one materialized draft path. This gate does not approve the resulting product artifacts. After materialization, Evolution routes the demand to Feature, Use Case, Specification, or another owning skill based on the approved classification; it does not create a parallel hierarchy.

Running `spec-framework init` interactively starts the Bubble Tea question
pipeline. The CLI shows an immutable installation plan and writes only after
confirmation.

Initialization adds only `product/` to the target repository. The selected versioned JSON contract composes reusable starter asset sets, entry-specific files, registry transformations, and typed initialization actions; invalid sources, targets, patches, profiles, or actions fail before initialization completes. Framework assets are materialized under the operating system's user cache, including the initialization contracts and the `examples/events/` domain-modeling reference, and selected harnesses receive one user-scoped `spec-framework` dispatcher. The Codex dispatcher is written to `~/.agents/skills/spec-framework`; upgrade removes only the legacy namespaced dispatcher after installing the replacement. Every starting point that creates or revises domains must read that reference before its first domain change for business-area boundaries, explicit non-ownership, cross-domain dependencies, and a Domain -> User Goal -> Feature -> Use Case walking skeleton; `audit-only` uses it to assess existing boundaries without mutation. The dispatcher activates exclusively from a valid `product/.product/framework.json` and resolves Framework Guide first unless it has a verified specialist route with concrete scope. `upgrade --yes` refreshes the dispatcher for every agent selected in the manifest or `--agents`; it never replays initialization contracts over adopter-owned product files.

The CLI expands and validates the selected contract in memory, including explicit empty directories, writes the complete product to staging inside the target repository, and publishes `product/` atomically only after guides, manifest, runtime preparation, dispatchers, and starting-point actions succeed. File/directory collisions or unsafe directory paths fail during planning. Failed initialization removes staging and leaves no partial `product/`. An existing `product/` always blocks `init`; the compatibility `--force` flag never authorizes overwriting adopter-owned product content.

The initialized product includes [`tools/check-links.py`](../starter/product/tools/check-links.py). Run `python product/tools/check-links.py product` locally or in CI to verify relative Markdown links and section anchors. The script uses only the Python standard library and exits with status 1 when a link is broken.

Decision-specific CI diagnostics can run alongside it:

```bash
spec-framework decisions check --product-root product --strict
```

Use `--json` for CI annotations or `--fix-links --yes` only after reviewing the preview. The command never changes approval records, moves decision files, or rewrites ambiguous references.

The pinned runtime also materializes the shared contracts [`execution-runtime.md`](../docs/execution-runtime.md), [`engineering-systems.md`](../docs/engineering-systems.md), [`engineering-catalog-and-standards.md`](../docs/engineering-catalog-and-standards.md), and [`lifecycle-and-approvals.md`](../docs/lifecycle-and-approvals.md) alongside `FRAMEWORK.md` and `AGENTS.framework.md`. Skills and orchestrators use these documents for cross-cutting runtime, engineering-system, catalog, standards, and lifecycle rules.
