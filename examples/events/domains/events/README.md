# Events Domain

## Purpose

The Events domain owns product knowledge for organizing, joining, attending, validating, and measuring event participation.

## When To Use

Use this domain when work concerns event lifecycle, attendee participation, organizer validation, check-in rules, or event-facing analytics. Do not use it for global authentication, payments, or notification infrastructure except as dependencies.

## Expected Files

| Artifact | Link | Purpose |
| --- | --- | --- |
| Context | [context.md](context.md) | Current domain context and dependencies. |
| Domain | [domain.md](domain.md) | Events domain definition and boundaries. |
| Participate in event | [goals/participate-in-event](goals/participate-in-event/goal.md) | User goal for attending and completing event participation. |

## Responsible Skill

Primary owner: Domain Architect AI.

Supporting skills: User Goal AI, Feature AI, Use Case AI, Specification AI, UX/UI AI.

## Next Step

Use the `participate-in-event` goal to refine event attendance workflows and approve the QR code check-in use cases before implementation planning.
