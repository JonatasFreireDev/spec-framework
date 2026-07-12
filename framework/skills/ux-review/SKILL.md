---
name: ux-review
description: "Independent read-only UX Review. Use to verify a Design against its Specification, visual sources, design system, accessibility, responsive behavior, states, and fidelity before approval."
---

# UX Review Skill

## Layer
Product Design Validation

## Responsibility
Review Design independently from the agent or adapter that created, evolved, or imported it. Produce findings and evidence; never fix or approve the artifact.

## Inputs
Approved Specification; `design.md`; source and use-case manifests; screen inventory; mappings; visual assets; design system; approved decisions.

## Outputs
UX review verdict; coverage and fidelity findings; accessibility and responsive findings; blockers; required fixes; evidence references; handoff to UX/UI or human approval.

## Required reading
- `FRAMEWORK.md` and FDR-021.
- Relevant parent and local `context.md` files.
- Specification contracts and `design.md`.
- Every source marked `visual_canonical` and its manifest.

## Workflow
1. Verify source identity, authority, version, and availability.
2. Verify every applicable REQ/AC is mapped to screens and states.
3. Verify the pinned Design System version, tokens, components, patterns, and deviations when a system is declared.
4. Review hierarchy, clarity, navigation, content, loading, empty, success, error, disabled, and permission states.
5. Review keyboard access, labels/roles, contrast, touch targets, reduced motion, and responsive coverage.
6. For strict fidelity, compare targets against the canonical source and classify deviations.
7. Return `approved`, `approved_with_notes`, or `blocked`; do not edit Design or approval records.

## Blocking findings
- Required behavior or state is missing or conflicts with the Specification.
- A canonical visual source is unversioned or unavailable without a snapshot.
- Security, privacy, permission, or accessibility requirements are violated.
- Strict-fidelity deviations lack explicit review.

## Handoff
Blocked findings return to UX/UI. A clean review goes to the human Design approval gate, then Implementation Planner.
