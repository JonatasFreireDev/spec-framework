# Context: QR Code Check-in

```yaml
id: FT-001
type: feature
name: QR Code Check-in
status: draft
owner_skill: feature
last_updated: 2026-07-09
```

## Purpose

Provides a concrete way for attendees and organizers to confirm event attendance.

## Parent Artifacts

- GOAL-001 - ../../context.md - parent user goal

## Child Artifacts

- UC-001 - use-cases/attendee-checks-in-with-qr-code/context.md - attendee-facing use case
- UC-002 - use-cases/organizer-validates-qr-code/context.md - organizer-facing use case

## Dependencies

- DOMAIN-users - attendee identity - blocking: yes
- Event attendance data model - event participation state - blocking: yes

## Related Artifacts

- knowledge/templates/specification-template.md - expected specification shape

## Canonical Documents

- Primary: feature.md
- Specifications:
  - use-cases/attendee-checks-in-with-qr-code/specification.md
  - use-cases/organizer-validates-qr-code/specification.md
- Implementation plans:
  - use-cases/attendee-checks-in-with-qr-code/implementation-plan.md
  - use-cases/organizer-validates-qr-code/implementation-plan.md
- Execution graphs:
  - use-cases/attendee-checks-in-with-qr-code/execution-graph.json
  - use-cases/organizer-validates-qr-code/execution-graph.json
- Tasks:
  - use-cases/attendee-checks-in-with-qr-code/tasks.md
  - use-cases/organizer-validates-qr-code/tasks.md

## Decisions

- None yet.

## Assumptions

- Attendee and organizer both have authenticated sessions.
- QR codes should be time-bound to reduce abuse.

## Open Questions

- What is the expiration window for generated QR codes?
- Should organizer validation remain online-only for L1?

## Handoff

Next recommended skill: use-case
Required reading before next step:
- feature.md
- use-cases/attendee-checks-in-with-qr-code/context.md
- use-cases/organizer-validates-qr-code/context.md
