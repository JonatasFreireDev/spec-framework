# Impeccable Design Adapter

Impeccable is an optional generation/evolution adapter. Install it separately with `npx impeccable install`. The Spec Framework owns gates, paths, manifests, mappings, review, and approvals.

Generate a non-executing adapter plan:

```bash
spec-framework design adapter --adapter impeccable --use-case <path> --maturity wireframe
```

The plan never executes slash commands automatically. Generated assets must stay under `product/design/use-cases/<slug>/`, remain non-production, and pass independent UX Review.
