# Node CLI characterization

This directory freezes the externally observable Node CLI behavior that the Go
implementation must preserve during the migration described by FDR-013.

`node-contract.json` is the frozen pre-cutover baseline. The capture utility was
removed with the Node implementation after the Go black-box suite reached full
parity.

The snapshot normalizes temporary paths and line endings. It covers CLI help,
initialization, validation, move dry-run semantics, and upgrade preservation.
The future Go parity suite should compare semantic fields rather than depending
on incidental formatting that is intentionally changed by the new TUI.
