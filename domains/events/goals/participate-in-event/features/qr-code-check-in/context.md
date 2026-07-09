# Context: QR Code Check-in

```yaml
id: FT-001
type: feature
name: QR Code Check-in
status: draft
owner_skill: 07-feature.md
last_updated: 2026-07-09
```

## Purpose

Provides a concrete way for attendees and organizers to confirm event attendance.

## Parent Artifacts

- GOAL-001 - ../../context.md - parent user goal

## Child Artifacts

- UC-001 - use-cases/attendee-checks-in-with-qr-code/context.md - use case

## Dependencies

- DOMAIN-users - attendee identity - blocking: yes
- Event attendance data model - event participation state - blocking: yes

## Related Artifacts

- product/knowledge/templates/specification-template.md - expected specification shape

## Canonical Documents

- Primary: feature.md
- Specification: use-cases/attendee-checks-in-with-qr-code/specification.md
- Implementation plan: use-cases/attendee-checks-in-with-qr-code/implementation-plan.md
- Execution graph: use-cases/attendee-checks-in-with-qr-code/execution-graph.json
- Tasks: use-cases/attendee-checks-in-with-qr-code/tasks.md

## Decisions

- None yet.

## Assumptions

- Attendee and organizer both have authenticated sessions.
- QR codes should be time-bound to reduce abuse.

## Open Questions

- What is the expiration window for generated QR codes?
- Should organizers scan attendee QR codes, or should attendees scan event QR codes?

## Handoff

Next recommended skill: 08-use-case.md
Required reading before next step:
- feature.md
- use-cases/attendee-checks-in-with-qr-code/context.md