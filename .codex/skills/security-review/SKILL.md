---
name: security-review
description: "Security Review Skill. Use when Codex needs to evaluate authentication, authorization, privacy, abuse, secrets, data exposure, logging, dependency, rollout, and residual-risk controls before an executable artifact can be validated or released in the Spec Framework workflow."
---

# Security Review Skill

## Layer
Validation

## Responsibility
Evaluate whether an executable artifact can move forward without unacceptable security, privacy, permission, abuse, or operational-risk gaps.

Security Review does not replace QA. QA verifies the complete delivery contract and evidence matrix. Security Review owns security findings, required mitigations, and residual-risk classification.

## Operating Modes
- create: produce a first `security-review.md` from approved specification, design, implementation plan, execution graph, tasks, tests, code evidence, and known risks.
- update: revise a security review after fixes or scope changes.
- audit: inspect an artifact bundle for security gaps without changing product scope.
- explain: summarize security posture, blockers, residual risks, and required approvals.

## Inputs
Specification; design; implementation plan; execution graph; tasks; tests; QA evidence; implementation notes; dependency reports; audit logs; approved decisions.

## Outputs
Security verdict; threat model summary; control checklist; blocking findings; residual risks; required fixes; approval or release-blocking recommendation.

## Required Reading
- the framework root's `FRAMEWORK.md`.
- Relevant parent and local `context.md` files.
- `knowledge/templates/security-review-template.md`.
- the active product root's `knowledge/conventions/security-baseline.md` and any linked domain baseline.
- Existing active product root `audits/security/threat-register.md` entries that affect the artifact.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.
- Related `tests.md`, `qa-evidence.md`, and `audit.md` when present.
- `engineering/decisions/FDR-006-failure-routing-and-regression.md`.

## Workflow
1. Read the local context and identify artifact status, delivery level, priority, and release intent.
2. Read the product security baseline and active threat register entries that affect the artifact.
3. Confirm the Specification contains permissions, data classification, privacy, abuse, observability, error handling, rollout, and acceptance criteria.
4. Confirm the Design avoids unsafe data exposure in UI states, errors, empty states, permission prompts, and accessibility flows.
5. Confirm the Implementation Plan covers server-authoritative checks, secrets, dependency risk, migrations, rollback, observability, and security tests.
6. Confirm the Execution Graph and Tasks include explicit security work when the flow touches data, permissions, tokens, payments, uploads, messaging, search, public endpoints, or admin operations.
7. Review QA evidence and verify that security controls have evidence, not only intention.
8. Classify findings as blocker, required fix, note, or accepted residual risk.
9. Route blockers using FDR-006: security bug with clear expected behavior -> `bug-fixer`; missing security test -> `qa`; incomplete implementation -> `code-runner`; missing permission/privacy decision -> `product-historian` plus human approval.
10. Do not mark the artifact secure when a blocker remains. Request a decision for any accepted high or hard-to-reverse residual risk.

## Review Checklist
- [ ] Authentication and authorization are server-authoritative.
- [ ] Least privilege is explicit for all actors and roles.
- [ ] Sensitive data and PII are minimized, protected, and retained only as needed.
- [ ] Inputs are validated and unsafe outputs are escaped or avoided.
- [ ] Abuse cases, replay, enumeration, rate limits, and idempotency are addressed where relevant.
- [ ] Secrets, tokens, and credentials are not exposed in UI, logs, analytics, code, or documentation examples.
- [ ] Logs, analytics, and audit trails support investigation without leaking sensitive data.
- [ ] Dependencies, migrations, rollout, rollback, and monitoring have security-aware handling.
- [ ] Security tests or manual security evidence exist for every security acceptance criterion.
- [ ] Residual risks are documented with owner, severity, mitigation, and approval status.
- [ ] Blocking findings include route and owner.
- [ ] Security defects that escaped require permanent regression coverage before closure.

## Verdict Rules
- `passed`: no blocking findings; required controls have evidence; residual risks are low or explicitly approved.
- `passed_with_notes`: no blocker, but non-blocking fixes or monitoring actions remain.
- `blocked`: any high-risk gap, missing evidence for a required control, unapproved permission/privacy decision, or release-impacting unknown remains.

## Handoff
Next: bug-fixer, code-runner, QA AI, Product Historian, Audit Orchestrator, or Release Orchestrator depending on FDR-006 routing.

Pass forward the verdict, evidence links, blockers, residual risks, required decisions, and whether release is blocked.
