# Engineering Test Strategy

## 🧾 Generation And Agent Self-Check

> Complete this section when materializing the artifact. Keep unresolved items explicit in the relevant scope, findings, risks, or handoff section.

| Field | Value |
| --- | --- |
| Generated on | `YYYY-MM-DD` |
| Purpose | `[decision, evidence, contract, or handoff this artifact supports]` |
| Use when | `[workflow stage, trigger, or condition]` |
| Prepared by | `[owning skill, role, or accountable person]` |
| Scope covered | `[artifact, product area, use case, or review boundary]` |
| Required inputs and evidence | `[links to approved parents, documents, code, decisions, or observations]` |
| Ready when | `[artifact-specific completion, evidence, and gate conditions]` |
| Current status | `[status allowed by this artifact's owning workflow]` |


## Scope

[Define the shared validation strategy. Delivery-specific cases remain in each use case's `tests.md`.]

## Test Levels

| Level | Purpose | Required when | Evidence |
| --- | --- | --- | --- |
| Unit/component | Verify isolated behavior and boundaries | `[policy]` | `[runner output]` |
| Integration/contract | Verify data, API, dependency, and ownership boundaries | `[policy]` | `[runner output]` |
| End-to-end | Verify critical user or operator flows | `[policy]` | `[runner output/artifact]` |
| Manual/exploratory | Investigate risks not economically automated | `[policy]` | `[session notes/screenshots]` |

## Risk-Based Coverage

| Trigger | Required considerations |
| --- | --- |
| Permissions, security, privacy | Negative cases, least privilege, safe failure, Security Review |
| Data mutation or migration | Constraints, idempotency, rollback, partial failure |
| Visual surface | States, accessibility, responsive/platform coverage, visual evidence |
| External dependency | Contract failure, timeout, retry, observability |
| Performance-sensitive flow | Budget, representative load, failure threshold |

## Environments And Data

| Environment | Allowed data | Isolation/reset | Owner | Limitations |
| --- | --- | --- | --- | --- |
| `[environment]` | `[synthetic/anonymized/etc.]` | `[method]` | `[owner]` | `[limitations]` |

## Flaky Test Policy

A flaky test is a quality finding, not a passing gate. Quarantine requires a tracked exception with owner, residual risk, mitigation, and expiry or review date. Required gates remain blocking unless an authorized product decision explicitly accepts the risk.

## Delivery Application

Each `tests.md` pins the consumed Engineering System id/version, maps every `AC-*` to at least one validation method, records applicable risks and negative cases, identifies environment and data needs, and lists deviations or `None`.

## ✅ Agent Verification Checklist

- [ ] Test levels, risk triggers, environments, data, tools, and evidence expectations are explicit.
- [ ] Security, accessibility, performance, reliability, and operational validation are addressed.
- [ ] Flaky-test handling, exceptions, ownership, expiry, and escalation are defined.
- [ ] Delivery teams can apply the strategy without replacing use-case-specific test design.
