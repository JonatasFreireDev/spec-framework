# Examples

## Purpose

Store examples that demonstrate how a product can use the framework.

Examples are learning material. They are not copied into new product repositories by default and are not product starter scope.

## Current Examples

| Example | Purpose |
| --- | --- |
| `events/` | QR-code check-in domain example extracted from the framework laboratory. |

## Validation

The `events/` example is a self-contained product root. Validate it with:

```bash
node engineering/validators/framework-validator.mjs --product-root examples/events --framework-root .
```

`npm run validate` runs this by default.
