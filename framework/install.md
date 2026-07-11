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

Initialize without a TTY:

```bash
spec-framework init --target ../my-product --agents codex,cursor,claude --yes
```

Use `--starting-point existing-documents` with `--source-dir` or `--sources` to bootstrap from existing product material. This creates an analysis-only import run; it never treats source prose as approved product truth.

Draft materialization is a separate explicit command after the import mappings have been reviewed: `spec-framework import materialize --run IMPORT-001 --approved-by "Product Owner" --yes`.

Running `spec-framework init` interactively starts the Bubble Tea question
pipeline. The CLI shows an immutable installation plan and writes only after
confirmation.
