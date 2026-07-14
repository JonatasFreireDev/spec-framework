# Lifecycle And Approval Contract

This document contains the detailed lifecycle, approval-record, staleness, and authority rules. `FRAMEWORK.md` remains the method source of truth; this file groups the operational contract used by workflow and review skills.

## Artifact lifecycle

General artifact states are:

```text
draft -> proposed -> approved -> in_progress -> implemented -> validated -> released
```

`deprecated` and `superseded` are terminal product states. `materialized` is not a general artifact state; it is specific to the Execution Graph and means that canonical task files and its generated index exist.

Every transition must be valid for the artifact type and its starting-point contract. A downstream artifact cannot advance while a required parent gate, decision, evidence, or configured command is missing.

## Approval authority

- Agents may propose and produce scoped drafts or implementation evidence.
- Human identity and explicit confirmation are required for approvals and consequential mutations.
- Conversation never creates approval evidence.
- Agents must not create, edit, or repair approval records unless a human explicitly authorizes a named migration that includes approval-record generation.
- Review skills are read-only unless their contract explicitly grants a scoped correction authority; findings must be routed to the owning skill otherwise.

Approved and later artifacts require a matching record in `.product/history/` with `artifact_id`, `path`, `content_hash`, `status_granted`, `approved_by`, `approved_at`, and `notes`. The record must match the current normalized artifact content.

## Delivery evidence

An implemented executable task records immutable working-tree evidence: branch, base commit, changed paths, diff hash, tests, and applicable gate results. It does not require an early commit. Validation requires the applicable approved QA, Code Review, Security Review, and concrete evidence. Code Review and task QA must cover the same current diff hash.

## Staleness and derivation

Staleness is derived by the validator, not edited as a status. Downstream artifacts record source hashes in `.product/derivations.json`. When source content changes, the downstream artifact is stale and must be regenerated or re-approved before advancing through gates.

## Failure and stop routing

- Clear defect with expected behavior -> `bug-fixer`.
- Missing or hollow coverage -> `qa` or the test owner.
- Incomplete or out-of-contract implementation -> `code-runner`.
- Missing product or architecture decision -> `product-historian` plus a human.
- Security or privacy blocker -> `security-review` and the owning implementation or QA skill.
- Conflicting or stale state -> Framework Guide and the relevant owner.

Three failed automated attempts require human escalation. No route may bypass an approval, scope, evidence, or authority gate.

## Owning skills and orchestrators

- `framework-guide` and `delivery-orchestrator`: route lifecycle and authority decisions.
- `code-runner`, `code-review`, `qa`, and `security-review`: implement or verify delivery evidence.
- `commit-crafter`, `pr-finalizer`, `release-orchestrator`, and `integration-orchestrator`: advance only after their preceding gates are satisfied.
