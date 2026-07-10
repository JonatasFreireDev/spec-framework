# Framework Validators

## Purpose

Store installed framework validators.

Validators should run against `product/` by default.

Expected future command:

```bash
node .spec-framework/validators/framework-validator.mjs --product-root product
```

Use `--framework-root .spec-framework` when running from the repository root and the validator cannot infer the framework asset location.
