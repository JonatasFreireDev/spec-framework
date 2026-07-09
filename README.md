# Product Engineering Framework v2

This workspace turns product ideas into structured product engineering artifacts through Specification Driven Development.

Start here:

- FRAMEWORK.md is the canonical architecture for this framework.
- knowledge/templates/ contains the artifact templates.
- skills/v2/ contains the proposed specialist skills and orchestrators.

## Ladder

Problem -> Vision -> Strategy -> Domain -> User Goal -> Feature -> Use Case -> Specification -> Implementation Plan -> Execution Graph -> Tasks -> Code -> Validation -> Audit

## Source of truth

- Domains contain user goals.
- User goals contain features.
- Features contain use cases.
- Use cases contain specification, design, analytics, tests, audit, implementation plan, execution graph, and tasks.
- The Specification is the source of truth for downstream implementation artifacts.
- Tasks are generated from the Specification through the Implementation Plan and Execution Graph.
- Decisions are recorded in knowledge/decisions/ and indexed in .product/decisions.json.

## Quality gates

- audits/readiness/ defines the readiness gate for moving artifacts toward executable tasks.
- audits/readiness/UC-001-readiness.md applies the gate to the QR Code Check-in example.