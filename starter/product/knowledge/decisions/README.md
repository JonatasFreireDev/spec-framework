# Product Decisions

## Purpose

Record product or cross-cutting decisions that change product scope, business rules, data, permissions, privacy, payment, delivery commitments, or hard-to-reverse strategy. Design and engineering decisions have dedicated roots under `product/design/decisions/` and `product/engineering/decisions/`.

## Expected File

Use the product decision template after framework assets are installed:

```text
framework/template/decision-template.md
```

## Index

Approved decisions should be indexed in `.product/decisions.json`. Set `domain` to `product` or `cross-cutting` and keep `path` relative to `product/`.
