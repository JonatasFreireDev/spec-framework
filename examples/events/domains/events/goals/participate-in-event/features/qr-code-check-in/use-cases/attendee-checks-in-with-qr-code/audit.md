# Audit: Attendee Checks In With QR Code Example

## Context

- Scope: UC-001 example subtree
- Auditor skill: audit-orchestrator.md
- Date: 2026-07-09
- Verdict: approved_with_notes

## Summary

The example demonstrates the intended framework flow from domain to executable tasks. It is not approved for implementation because two product/security decisions remain open.

## Findings

### Medium QR Expiration Is Not Approved

- Evidence: specification.md > Open Questions
- Impact: token generation cannot be implemented safely without an expiration policy.
- Required fix: create and approve a decision record for QR expiration.
- Owner: product/security

### Medium Token Strategy Is Not Approved

- Evidence: implementation-plan.md > Decisions Needed
- Impact: data model and backend implementation differ depending on stored token vs signed token.
- Required fix: create and approve architecture/security decision.
- Owner: architecture/security

## Gaps

- Real code paths are placeholders until the app structure is inspected.
- Scanner dependency is not selected.

## Conflicts

- None found inside the example subtree.

## Dependencies

- Users/auth domain.
- Organizer permission model.
- Event attendance data model.

## Residual Risk

- Offline venue conditions are explicitly out of scope for v1 and may affect real-world usability.

## Approval

- Approved by:
- Date:
- Notes: