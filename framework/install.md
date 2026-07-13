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
irm https://raw.githubusercontent.com/JonatasFreireDev/spec-framework/master/scripts/init.ps1 | iex
```

Linux and macOS users can download and inspect `scripts/init.sh`, then execute it locally. Piping a remote shell script is not required.

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

Use `--starting-point existing-documents` with `--source-dir` or `--sources` to bootstrap from existing product material. This creates an analysis-only import run; it never treats source prose as approved product truth.

Draft materialization is a separate explicit command after the import mappings have been reviewed: `spec-framework import materialize --run IMPORT-001 --approved-by "Product Owner" --yes`.

The generated `product/BOOTSTRAP.md` pins the active import run id. `spec-framework work` remains blocked until that latest run records explicit materialization approval and at least one materialized draft path. This gate does not approve the resulting product artifacts.

Running `spec-framework init` interactively starts the Bubble Tea question
pipeline. The CLI shows an immutable installation plan and writes only after
confirmation.

Initialization adds only `product/` to the target repository. Framework assets are materialized under the operating system's user cache, and selected harnesses receive one user-scoped `spec-framework` dispatcher. The dispatcher activates exclusively from a valid `product/.product/framework.json`.
