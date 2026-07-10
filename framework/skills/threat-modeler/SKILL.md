---
name: threat-modeler
description: "Threat Modeler Skill. Use when Codex needs to proactively model product-wide, domain-level, or cross-feature security threats; create or update security baselines; maintain a threat register; or analyze authorization, data integrity, privacy, abuse, replay, concurrency, idempotency, and stack-specific vulnerability classes before individual Security Review gates."
---

# Threat Modeler Skill

## Layer
Audit and Security Planning

## Responsibility
Model threats proactively across the product, a domain, a user goal, or a feature family before individual delivery artifacts reach Security Review.

Threat Modeler does not approve releases and does not replace Security Review. It creates reusable security context that Security Review must read.

## Operating Modes
- create: produce a first threat register or domain security baseline.
- update: revise threats, mitigations, owners, or status after scope, architecture, or stack changes.
- audit: inspect existing documentation for missing threat coverage.
- explain: summarize current threat posture, blockers, residual risks, and next owners.

## Inputs
the framework root's FRAMEWORK.md; relevant `context.md` files; domain, goal, feature, use case, specification, design, implementation plan, execution graph, tasks, QA evidence, Security Review, audit reports, approved decisions, stack conventions, and known incidents.

## Outputs
Security baseline updates; threat register entries; threat scenarios; impacted artifacts; mitigations; residual risks; required decisions; routing recommendations.

## Required Reading
- the framework root's `FRAMEWORK.md`.
- Relevant parent and local `context.md` files.
- the active product root's `knowledge/conventions/security-baseline.md`.
- `framework/template/security-baseline-template.md` when creating or normalizing a baseline.
- `framework/template/threat-register-template.md` when creating or normalizing a threat register.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.
- Framework security decisions in `framework/decisions/`, especially FDR-006 and FDR-009.
- Existing active product root `audits/security/threat-register.md` when present.
- Related Security Review, QA Evidence, and audit artifacts when present.

## Workflow
1. Define scope: product, domain, user goal, feature group, or high-risk use case.
2. Inventory assets, actors, trust boundaries, data classes, permissions, external systems, background jobs, logs, analytics, and operational controls.
3. Identify threat scenarios across authorization, authentication, privacy, data integrity, replay, abuse, concurrency, idempotency, secrets, dependencies, migrations, rollout, rollback, and observability.
4. Cross-check the relevant stack and product conventions for known vulnerability classes. Use web or official sources only when the user asks for current stack-specific research or the local docs are insufficient.
5. Map each threat to impacted artifacts and required mitigations. Do not invent product scope; mark missing behavior as an open question or decision candidate.
6. Update the security baseline when a rule should apply repeatedly across the product or domain.
7. Update the threat register when a concrete scenario needs tracking, mitigation, owner, evidence, or acceptance.
8. Route blockers through FDR-006. Vulnerability with expected behavior -> `bug-fixer`; missing security test -> `qa`; incomplete control -> `code-runner`; missing permission/privacy decision -> `product-historian` plus human approval.
9. Hand off to Security Review with baseline links, active threats, required evidence, and residual risks.

## Threat Lenses
- Authorization: role confusion, object ownership, cross-tenant access, admin bypass, direct object reference.
- Authentication and sessions: replay, token leakage, stale sessions, account takeover, weak recovery.
- Data integrity: tampering, race conditions, duplicate actions, partial writes, invalid state transitions.
- Privacy: PII exposure, over-collection, unsafe retention, analytics leakage, screenshot/log leakage.
- Abuse: enumeration, scraping, spam, fraud, harassment, content abuse, rate-limit bypass.
- Concurrency and idempotency: duplicate requests, retries, webhook replay, QR/token reuse, background job races.
- Secrets and dependencies: exposed credentials, insecure defaults, dependency CVEs, supply-chain risk.
- Operations: migrations, rollback, monitoring gaps, incident investigation gaps, unsafe feature flags.

## Output Rules
- Use concise tables with severity, likelihood, impact, affected artifacts, mitigation, owner, status, and evidence.
- Mark unverified assumptions explicitly.
- Never mark a threat as mitigated without evidence.
- Do not create or edit approval records.
- Do not implement application code.
- Do not downgrade a high or hard-to-reverse residual risk without human approval.

## Handoff
Next: Security Review, QA AI, Product Historian, Code Runner, Bug Fixer, or Audit Orchestrator depending on the finding route.

Pass forward the baseline path, threat register path, active blockers, required mitigations, required decisions, and residual risks.
