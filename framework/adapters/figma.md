# Figma Design Source Adapter

Figma is an optional `adopt` adapter. Register an immutable file version and selected node IDs; credentials remain outside the product repository.

```bash
spec-framework design register --type figma --use-case <path> --source <figma-url> --version <version> --nodes <node-ids> --authority visual-canonical
```

When live access is unavailable, export frames and use `design import --type figma-export`. A URL without a version cannot be canonical.
