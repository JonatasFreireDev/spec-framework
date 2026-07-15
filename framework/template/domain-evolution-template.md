# Domain Evolution: [cycle name]

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


## Snapshot

| Field | Value |
| --- | --- |
| ID | `EVOLUTION-NNN` |
| Status | `draft` |
| Domain | `[DOMAIN-ID/path]` |
| Owner | `domain-evolution-orchestrator` |

## Opportunity map

| Goal | Journey gap | Evidence | Candidate outcome |
| --- | --- | --- | --- |
| `[GOAL-*]` | `[gap]` | `[source]` | `[outcome]` |

## Candidate features

| Candidate | User value | Slice | Dependencies | Risk | Delivery |
| --- | --- | --- | --- | --- | --- |
| `[FT-*]` | `[value]` | `[entry → observable end]` | `[paths]` | `[risk]` | `[L*/P*]` |

## Selection

| Field | Value |
| --- | --- |
| Selected feature | `[path or pending]` |
| Approved by | `[human or pending]` |
| Rationale | `[comparison]` |
| Deferred/rejected | `[candidates and reasons]` |

## Handoff

Next: `new-feature-orchestrator` after explicit selection approval.

## ✅ Agent Verification Checklist

- [ ] Opportunities trace to goals, journey gaps, evidence, and candidate outcomes.
- [ ] Candidates are compared by value, dependencies, risks, level, priority, and slice.
- [ ] Selection records explicit human direction without treating rejected candidates as approved scope.
- [ ] The handoff names the selected feature boundary, owner, inputs, and next gate.
