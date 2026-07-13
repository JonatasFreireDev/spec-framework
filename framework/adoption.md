# Adoption Guide

## Goal

Create a new Specification Driven Development product repository that adds only `product/` and keeps the versioned framework runtime outside the repository.

## Recommended Path

Automated bootstrap:

```bash
spec-framework init ../my-product
```

CLI-style bootstrap from the framework repository:

```bash
spec-framework init ../my-product --agents codex,cursor,claude --yes
```

The wizard also asks for the starting point. All skills remain installed; the choice changes the generated bootstrap. For existing documents, use `--starting-point existing-documents` with `--source-dir` or `--sources`. The command inventories sources under `product/knowledge/imports/` but does not create Domains, User Goals, or Features without explicit approval.

After the Artifact Importer fills `mapping.json`, review the selected mappings and materialize them explicitly:

```bash
spec-framework import materialize --run IMPORT-001 --approved-by "Product Owner" --yes
```

For `existing-documents`, the latest run pinned in `product/.product/framework.json` must be materially complete before `spec-framework work` can create a workspace. Materialization authorizes selected draft writes only; review and approve each resulting product artifact through its normal owner and parent gates.

The command rejects missing evidence, paths outside `product/`, duplicate targets, non-draft content, and existing destination files.

Use `spec-framework work --feature <path-or-id>` to create an independent workspace, then `status` and `next` to see blockers and the canonical next skill. Use `approve` for human-reviewed status grants, `gates` before Code Runner, and `graph ready/claim/release/complete` to coordinate task ownership.

See [delivery-closure.md](delivery-closure.md) for the complete operational flow and command examples.

Install a versioned release binary as described in [install.md](install.md). Go and Node.js are not runtime requirements for adopters.

Initialization writes `product/.product/framework.json`, materializes the pinned embedded assets in the user cache, and installs one namespaced dispatcher for each selected agent in the user's harness directory. It does not create `.spec-framework/`, local agent trees, root guides, or CI workflows.

Activation is manifest-only. Mentions of Spec Framework do not activate the dispatcher when `product/.product/framework.json` is absent or invalid.

Manual development bootstrap:

1. Create or open the product repository.
2. From the framework source repository, run `go run ./cmd/spec-framework init <target>`.
3. Replace `product/` starter placeholders with product-specific content.
4. Run validation from the repository root.

```bash
spec-framework validate
```

Direct validator form when debugging:

```bash
spec-framework validate --product-root product --framework-root <framework-source-root> --write-registry --write-report
```

Upgrade an initialized product from the framework repository:

```bash
spec-framework upgrade --target ../my-product --agents codex --yes
```

## What Belongs To The Product

| Product-Owned Area | Purpose |
| --- | --- |
| `product/.product/` | Product state, registry, derivations, approval records, and adopted framework metadata. |
| `product/foundation/` | Problem, vision, and strategy for the product. |
| `product/domains/` | Product domains, goals, features, use cases, specifications, and tasks. |
| `product/knowledge/decisions/` | Product decisions only. |
| `product/knowledge/business-rules/` | Product business rules. |
| `product/audits/` | Product audits, readiness reports, QA evidence references, and threat register. |
| `product/releases/` | Product release notes and release readiness. |
| `product/design/` | Product design artifacts and mockups. |

## What Belongs To The Framework

| Framework-Owned Area | Purpose |
| --- | --- |
| Versioned user cache | Method, decisions, skills, templates, and validation assets resolved from the product manifest. |
| User-scoped `spec-framework` dispatcher | Manifest-gated agent integration that resolves specialized skills. |
| Installed `spec-framework` binary | Validation, bootstrap, upgrade, and migration tooling. |

## Non-Goals For Starter Repositories

- Do not copy framework FDRs into `product/knowledge/decisions/`.
- Do not inherit example domains as real product scope.
- Do not inherit retroactive approval records from the framework lab.
- Do not edit cached framework internals to encode product scope.

## Upgrade Direction

Stable commands:

```bash
spec-framework init ../my-product
spec-framework validate
spec-framework upgrade --target ../my-product
spec-framework dashboard --work WORK-001
spec-framework engineering-system inspect
spec-framework engineering-system validate
spec-framework engineering-system triggers
spec-framework engineering-system migrate --dry-run
spec-framework decisions migrate
```

Use `decisions migrate` as a preview first. Existing repositories should use `--interactive` to review ambiguous inferred types and scopes before applying the metadata migration.

Engineering System catalogs created before schema versioning must run `engineering-system migrate --dry-run` and then the same command without `--dry-run`. The migration only adds `schema_version: 1`, preserves product-owned fields, and never creates approval records. Approved systems must be re-approved by a human after any migrated content change.

Adoption is backed by the validator, package smoke tests, and the external-runtime / `product/` boundary.
