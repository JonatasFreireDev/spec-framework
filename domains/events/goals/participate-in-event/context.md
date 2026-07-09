# Context: Participate In Event

```yaml
id: GOAL-001
type: goal
name: Participate in event
status: draft
owner_skill: user-goal
last_updated: 2026-07-09
delivery:
  level: L1
  priority: P0
  depends_on:
    - DOMAIN-001
  rationale: Participating in an event is the user goal that the QR check-in walking skeleton proves end to end.
```

## Purpose

Represents the user's goal of successfully attending an event they intend to join.

## Parent Artifacts

- DOMAIN-001 - ../../context.md - parent domain

## Child Artifacts

- FT-001 - features/qr-code-check-in/context.md - feature

## Dependencies

- DOMAIN-users - authenticated attendee identity - blocking: yes

## Related Artifacts

- FRAMEWORK.md - framework source of truth

## Canonical Documents

- Primary: goal.md
- Specification: N/A
- Implementation plan: N/A
- Execution graph: N/A
- Tasks: N/A

## Decisions

- None yet.

## Assumptions

- Attendees may need proof of admission at the event location.

## Open Questions

- Should check-in be organizer-driven, attendee-driven, or support both modes?

## Handoff

Next recommended skill: feature
Required reading before next step:
- goal.md
- features/qr-code-check-in/context.md
