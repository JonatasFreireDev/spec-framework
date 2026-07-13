# Examples

## Purpose

Store examples that demonstrate how a product can use the framework.

Examples are not copied into new product repositories by default and are not product starter scope. `events/` is nevertheless the canonical reference for initial domain modeling: read it before creating the first domain to learn business-area boundaries, `Does Not Own` declarations, cross-domain dependencies, and a walking skeleton.

## Current Examples

| Example | Purpose |
| --- | --- |
| `events/` | QR-code check-in domain example extracted from the framework laboratory. |

## Validation

The `events/` example is a self-contained product root. Validate it with:

```bash
spec-framework validate --product-root examples/events --framework-root .
```

The Go validation workflow runs this on every pull request.
