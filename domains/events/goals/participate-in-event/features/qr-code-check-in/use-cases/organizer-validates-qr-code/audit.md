# Audit: Organizer Validates QR Code

## Context

- ID: AUD-002
- Status: draft
- Scope: UC-002, SPEC-002, DES-002, PLAN-002, GRAPH-002, TASKSET-002, TEST-002, ANA-002

## Verdict

Approved with notes.

## Findings

### F-001 Online Versus Offline Validation Is Not Approved

- Severity: medium
- Evidence: `specification.md` lists offline validation as non-goal and open question.
- Required fix: confirm online-only validation for L1 or approve an offline validation decision.

### F-002 Organizer Permission Roles Need Human Approval

- Severity: high
- Evidence: `context.md`, `specification.md`, and `implementation-plan.md` all list organizer roles as open.
- Required fix: approve which roles can validate check-in.

### F-003 Manual Fallback Is Out Of Scope

- Severity: medium
- Evidence: `design.md` and `implementation-plan.md` leave camera failure fallback as an open question.
- Required fix: decide whether L1 can ship without manual fallback.

## Evidence

- Use case contains main, alternate, error, and edge flows.
- Specification includes required framework sections.
- Design covers scanner states and accessibility.
- Plan avoids application code and defines sequencing.
- Execution graph is a DAG.
- Tasks point to specification-derived work.
- Tests cover behavior, permissions, data, UX, analytics, and accessibility.

## Required Fixes Before Approval

- Approve organizer roles.
- Approve online-only or offline validation scope.
- Approve whether manual fallback is required for camera failure.

## Suggested Improvements

- Add a future UX review after mockups exist.
- Add a release readiness report before shipping.
- Link this use case to a concrete release candidate when roadmap is approved.

## Residual Risk

Venue operations may suffer if network connectivity is poor and online-only validation remains the L1 choice.

## Next Recommended Skill

Impact Analysis AI for the open decisions, then Product Historian AI to record approved decisions.
