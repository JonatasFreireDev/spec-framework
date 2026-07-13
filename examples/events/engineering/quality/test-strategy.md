# Events Engineering Test Strategy

## Scope

This fixture applies shared planning policy to QR check-in documentation. It does not claim executable test infrastructure.

## Test Levels

| Level | Purpose | Required when | Current evidence |
| --- | --- | --- | --- |
| Unit/component | Token rules and isolated state transitions | Implementation exists | Not available |
| Integration/contract | Authorization, persistence, idempotency, expiry, event boundaries | QR validation is implemented | Not available |
| End-to-end | Attendee and organizer critical flows | UI and runtime exist | Not available |
| Manual/exploratory | Visual states, venue flow, and basic accessibility | Reviewable UI exists | Not available |

## Risk-Based Coverage

| Trigger | Required coverage |
| --- | --- |
| Permission boundary | Unauthorized and wrong-event organizer attempts are denied and safely logged. |
| Token lifecycle | Invalid, expired, replayed, and duplicate tokens fail safely or remain idempotent. |
| Data mutation | A successful scan creates one attendance result and partial failures do not duplicate it. |
| Visual surface | Loading, success, duplicate, invalid, expired, permission-denied, and offline states plus basic accessibility. |
| Observability | Outcomes are distinguishable without exposing token or attendee PII. |

## Environments And Data

Synthetic fixtures are required for events, organizers, attendees, valid/invalid/expired tokens, and permission boundaries. No executable environment is configured, so planned QA records `not run` and blocks validation.

## Flaky Tests And Exceptions

No flaky-test quarantine or policy exception is accepted in the fixture.

## Delivery Application

Each QR check-in `tests.md` pins `ENGSYS-EVENTS-001 @ 0.1.0`, maps acceptance criteria to planned methods, and records missing runtime evidence explicitly.
