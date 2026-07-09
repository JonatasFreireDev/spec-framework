# Context: Events

```yaml
id: DOMAIN-001
type: domain
name: Events
status: draft
owner_skill: 04-domain-architect.md
last_updated: 2026-07-09
```

## Purpose

Owns the product knowledge related to organizing and participating in events.

## Parent Artifacts

- strategy - product/foundation/strategy/README.md - product direction placeholder

## Child Artifacts

- GOAL-001 - goals/participate-in-event/context.md - user goal

## Dependencies

- DOMAIN-users - user identity and profile data - blocking: yes
- DOMAIN-notifications - event reminders and updates - blocking: no

## Related Artifacts

- product/FRAMEWORK.md - framework source of truth

## Canonical Documents

- Primary: domain.md
- Specification: N/A
- Implementation plan: N/A
- Execution graph: N/A
- Tasks: N/A

## Decisions

- None yet.

## Assumptions

- Events require authenticated users.
- Participation has permission and abuse-prevention implications.

## Open Questions

- Should venues and hosts be separate domains or concepts inside Events?

## Handoff

Next recommended skill: 05-user-goal.md
Required reading before next step:
- domain.md
- goals/participate-in-event/context.md