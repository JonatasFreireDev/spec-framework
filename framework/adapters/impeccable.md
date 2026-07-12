# Impeccable Design Adapter

Impeccable is an optional generation/evolution adapter. The Spec Framework owns gates, paths, manifests, mappings, review, and approvals.

Use supervised discovery and a version-pinned installation:

```bash
spec-framework adapters doctor impeccable --check-latest
spec-framework adapters install impeccable --version <cli-version>
spec-framework adapters install impeccable --version <cli-version> --yes
```

The first install command previews the exact official `npx impeccable@<version> skills install` invocation. The second explicitly authorizes it. The upstream installer may still ask harness-specific questions.

Generate a non-executing adapter plan:

```bash
spec-framework design adapter --adapter impeccable --use-case <path> --maturity wireframe
```

The plan never executes slash commands automatically. Generated assets must stay under `product/design/use-cases/<slug>/`, remain non-production, and pass independent UX Review.
