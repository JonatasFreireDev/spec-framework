---
name: verify
description: Run the full mechanical gate suite for this repository after any change to scripts/, framework/ (skills, templates, validators, tools, tests), starter/, examples/events/, or package.json. Use before committing, when asked to verify a change, or when a validator/test failure needs to be diagnosed.
---

# Verify Skill

## Purpose

Run every mechanical gate this repository defines and report a clear pass/fail verdict. No change to framework assets, scripts, the starter, or the worked example should be considered done until this suite is green.

## Gates

Run from the repository root, in this order:

```bash
npm run check     # node --check on every CLI script, validator, and test runner
npm test          # framework test suite (framework/tests/run-tests.mjs)
npm run validate  # framework-validator against the worked example (--product-root examples/events --framework-root .)
```

When the change touches packaging surfaces (`package.json` `files`/`bin`, `scripts/framework-assets.mjs`, anything under `framework/` that ships to adopters), also run:

```bash
npm run pack:dry  # confirm the npm tarball still includes every shipped asset
```

## Interpreting results

- `npm test` runs the suite in `framework/tests/run-tests.mjs` and prints `ok - <test name>` lines plus a final `N/N tests passed`. Any count below the total is a failure even if the process exits quietly.
- `npm run validate` validates `examples/events` as the active product root. Failures usually mean the worked example drifted from a changed template, validator rule, or skill contract — fix the drift, do not weaken the validator to make it pass.
- The test suite includes packaging tests (asset inclusion, CLI install from a consumer project). If those fail after adding a new asset directory, check both the `files` list in `package.json` and the copy list in `installFrameworkAssets` in `scripts/framework-assets.mjs`.

## Release-level verification

For changes to `scripts/init-product.mjs`, `scripts/upgrade-product.mjs`, or `scripts/framework-assets.mjs`, additionally exercise the end-to-end consumer path in a temporary directory outside the repository:

```bash
npm pack
mkdir <tmp>/consumer && cd <tmp>/consumer
npm install <repo>/spec-framework-<version>.tgz --no-save
npx spec-framework init --target ../my-product
cd ../my-product && npx spec-framework validate
```

A freshly initialized product is expected to have validator findings for unfilled foundation content; the gate here is that `init` completes, `.spec-framework/` is fully populated, and `validate` runs rather than crashing.

## Reporting

Report each gate with its command, pass/fail status, and the failing output verbatim when something breaks. Never summarize a failure as "some tests failed" — name the failing test or validator rule.
