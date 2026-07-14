# Events Example

## Purpose

Demonstrate the framework with the domain `events`, goal `participate-in-event`, feature `qr-code-check-in`, and use cases `attendee-checks-in-with-qr-code` and `organizer-validates-qr-code`.

## Structure

This folder is a self-contained product root, mirroring the canonical structure from `FRAMEWORK.md` section 4:

| Artifact | Link |
| --- | --- |
| Events domain | [domains/events/context.md](domains/events/context.md) |
| QR code feature | [domains/events/goals/participate-in-event/features/qr-code-check-in/context.md](domains/events/goals/participate-in-event/features/qr-code-check-in/context.md) |
| Attendee use case | [domains/events/goals/participate-in-event/features/qr-code-check-in/use-cases/attendee-checks-in-with-qr-code/context.md](domains/events/goals/participate-in-event/features/qr-code-check-in/use-cases/attendee-checks-in-with-qr-code/context.md) |
| Organizer use case | [domains/events/goals/participate-in-event/features/qr-code-check-in/use-cases/organizer-validates-qr-code/context.md](domains/events/goals/participate-in-event/features/qr-code-check-in/use-cases/organizer-validates-qr-code/context.md) |
| Product decisions | [knowledge/decisions/](knowledge/decisions/) |
| Approval history | [.product/history/](.product/history/) |

Framework core assets (`FRAMEWORK.md`, `AGENTS.md`, `framework/skills/`, `framework/template/`, and `framework/validators/`) live outside this folder, primarily under `framework/`. Documents in this example that reference them use relative links that cross that boundary; the validator reports those as warnings because they point outside `--product-root`.

## Validation

```bash
spec-framework validate --product-root examples/events --framework-root .
```

The Go validation workflow runs this on every pull request.
