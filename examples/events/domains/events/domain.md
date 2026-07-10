# Domain: Events

## Context

- ID: DOMAIN-001
- Status: draft
- Parent: product strategy
- Context file: context.md

## Definition

Events owns the knowledge and workflows for creating, discovering, joining, attending, and managing social events.

## User And Business Outcomes

- Users can find and attend relevant social experiences.
- Organizers can manage attendance with enough trust and operational control.
- The product can measure event participation as a core value signal.

## Boundaries

### Owns

- Event lifecycle and participation rules.
- Attendance state and check-in behavior.
- Event-facing analytics and operational states.

### Does Not Own

- Authentication and identity verification.
- Payment processing.
- Global notification delivery infrastructure.

## Core Concepts

- Event: a planned social gathering users may attend.
- Attendee: a user who participates in an event.
- Check-in: proof that an attendee arrived or was admitted.

## User Goals

- GOAL-001 - Participate in event - draft

## Cross-Domain Dependencies

- Users - identity, profile, and auth state.
- Notifications - reminders and event updates.

## Business Rules

- Attendees must be authenticated before check-in.
- Organizers need a way to validate attendance without exposing unnecessary user data.

## Metrics

- Event participation rate.
- Check-in completion rate.
- No-show rate.

## Risks And Open Questions

- Fraudulent check-ins may distort participation metrics.
- Offline or unstable network conditions may affect venue operations.

## Approval

- Approved by:
- Date:
- Notes: