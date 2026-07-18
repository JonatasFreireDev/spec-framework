# Engineering Quality Model

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


| Attribute | Product expectation | Required evidence | Gate | Maturity |
| --- | --- | --- | --- | --- |
| Reliability | `[expectation]` | `[test/log/runbook]` | `[gate or Not configured]` | `baseline` |
| Security | `[expectation]` | `[test/review/decision]` | `[gate or Not configured]` | `baseline` |
| Observability | `[expectation]` | `[log/metric/trace]` | `[gate or Not configured]` | `baseline` |
| Performance | `[expectation]` | `[budget/benchmark]` | `[gate or Not configured]` | `baseline` |
| Evolvability | `[expectation]` | `[boundary/migration test]` | `[gate or Not configured]` | `baseline` |

Maturity records available evidence and does not approve risk or architecture.

## ✅ Agent Verification Checklist

- [ ] Each quality attribute states a product expectation, measurable evidence, gate, and maturity.
- [ ] Reliability, security, observability, performance, accessibility, and evolvability are considered.
- [ ] Missing capabilities and exceptions are explicit and owned.
- [ ] Maturity describes available evidence and does not approve architecture or residual risk.
