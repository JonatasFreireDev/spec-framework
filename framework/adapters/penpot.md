# Penpot Design Source Adapter

Penpot is an optional `adopt` adapter. Register an immutable file version and selected object IDs; credentials remain outside the product repository.

```bash
spec-framework design register --type penpot --use-case <path> --source <penpot-url> --version <version> --nodes <object-ids> --authority visual-canonical
```

When live access is unavailable, export frames and use `design import --type penpot-export`. A URL without a version cannot be canonical.
