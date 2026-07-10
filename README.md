# Product Engineering Framework v2

This repository is the framework laboratory and reusable base for Specification Driven Development.

## Start Here

| Area | Purpose |
| --- | --- |
| [FRAMEWORK.md](FRAMEWORK.md) | Canonical method and architecture. |
| [framework/](framework/) | Framework core boundary and adoption guide. |
| [starter/](starter/) | Clean product skeleton for new repositories. |
| [examples/](examples/) | Worked examples and learning material. |
| [.codex/skills/](.codex/skills/) | Operational Codex skills. |
| [knowledge/templates/](knowledge/templates/) | Reusable artifact templates. |
| [engineering/validators/](engineering/validators/) | Mechanical validation gates. |

## Adoption

For a new product repository, start from [starter/](starter/) rather than copying this repository root.

The starter creates two explicit roots:

```text
.spec-framework/  # how the framework works
product/          # the product being built
```

Current recommended flow:

```text
node scripts/init-product.mjs --target ../my-product -> fill product/foundation -> create product/domains -> node .spec-framework/tools/validate-product.mjs
```

Local CLI form:

```bash
node scripts/spec-framework.mjs init --target ../my-product
node scripts/spec-framework.mjs validate
node scripts/spec-framework.mjs upgrade --target ../my-product
```

See [framework/adoption.md](framework/adoption.md).

## Ladder

```text
Problem -> Vision -> Strategy -> Domain -> User Goal -> Feature -> Use Case -> Specification -> Design -> Implementation Plan -> Execution Graph -> Tasks -> Code -> Validation -> Audit
```

## Repository Boundary

| Concern | Lives In |
| --- | --- |
| Framework method, validators, skills, templates, and FDRs | This repository core |
| New product state and product scope | `product/` in a repo created from `starter/` |
| Installed framework method assets | `.spec-framework/` in a repo created from `starter/` |
| Example domains and use cases | `examples/` |

## Quality Gates

Run:

```bash
node engineering/validators/framework-validator.mjs --write-registry --write-report
node engineering/tests/run-tests.mjs
```
