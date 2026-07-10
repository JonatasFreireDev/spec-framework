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

Local linked CLI:

```bash
npm link
spec-framework init --target ../my-product
spec-framework validate
spec-framework upgrade --target ../my-product
```

Packaged CLI smoke path:

```bash
npm pack
mkdir ../my-consumer
cd ../my-consumer
npm install ../spec-framework/spec-framework-0.1.0.tgz --no-save
npx spec-framework init --target ../my-product
cd ../my-product
npx spec-framework validate
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
npm run check
npm test
npm run pack:dry
npm run validate
node engineering/validators/framework-validator.mjs --product-root examples/events --framework-root . --write-registry --write-report
```
