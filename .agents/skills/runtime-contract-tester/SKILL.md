---
name: runtime-contract-tester
description: Test CLI and runtime contract changes in the Spec Framework, including optional capabilities, dispatch, imports, persistence, concurrency, and failure behavior.
---

# Runtime Contract Tester

## Purpose

Turn runtime changes into explicit, compatible contracts. Do not declare an unavailable external or sandboxed command as passing.

## Required reading

- `FRAMEWORK.md`
- changed `cmd/` and `internal/` packages
- related CLI help, JSON contracts, and tests

## Workflow

1. List changed commands, flags, persisted files, defaults, and external process boundaries.
2. Verify new optional capabilities are disabled by default and can be enabled, disabled, and removed safely.
3. Test success, invalid input, interrupted execution, resume/idempotency, and backward-compatible absence of configuration.
4. For concurrent work, test lease ownership, dependency ordering, write scopes, and deterministic reconciliation.
5. Combine changed features with import, dispatch, review, approval, and upgrade boundaries when applicable.

## Output

Return the contract matrix, commands run, compatible defaults, combination-test evidence, failures, and residual runtime risks.
