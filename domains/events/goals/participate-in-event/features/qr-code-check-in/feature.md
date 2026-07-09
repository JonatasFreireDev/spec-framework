# Feature: QR Code Check-in

## Context

- ID: FT-001
- Status: draft
- Domain: DOMAIN-001
- User goal: GOAL-001
- Context file: context.md

## Summary

QR Code Check-in lets an attendee prove presence at an event and lets an organizer validate attendance quickly.

## Problem Fit

- User problem: attendees need a simple way to confirm they arrived.
- Business reason: actual attendance is a stronger value metric than RSVP.
- Evidence: assumed need for event operations; requires validation.

## Scope

### In Scope

- Generate a check-in QR code for an authenticated attendee.
- Validate the QR code for a specific event.
- Mark attendance as checked in.
- Prevent duplicate check-in side effects.

### Non-Goals

- Full offline mode.
- Payment or ticketing.
- Fraud scoring beyond basic expiration and single-use behavior.

## Use Cases

- UC-001 - Attendee checks in with QR code - draft

## UX Notes

- Entry points: event details and organizer check-in surface.
- Core states: QR available, QR expired, scan success, scan failed, already checked in.
- Empty/loading/error states: must be specified before implementation.
- Accessibility notes: QR status must have text equivalents.

## Data And Permissions

- Data touched: event attendance record, attendee id, event id, check-in timestamp.
- Permission model: attendee can generate own proof; organizer can validate for owned/managed event.
- Sensitive data or abuse risks: QR should not expose raw personal data.

## Analytics

- qr_check_in_generated - attendee generated QR.
- qr_check_in_validated - organizer validated QR.
- qr_check_in_failed - validation failed.

## Dependencies

- Authenticated user identity - blocking: yes
- Event attendance model - blocking: yes
- Organizer permission model - blocking: yes

## Acceptance Intent

The feature can move to use-case specification when QR ownership, expiration, validation authority, duplicate behavior, and analytics are defined.

## Open Questions

- Should QR payload be opaque token or signed payload?
- What exact roles can validate attendance?