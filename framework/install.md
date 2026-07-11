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

Running `spec-framework init` interactively starts the Bubble Tea question
pipeline. The CLI shows an immutable installation plan and writes only after
confirmation.
