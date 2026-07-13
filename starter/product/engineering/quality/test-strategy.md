# Engineering Test Strategy

## Scope

Not configured. Define shared test levels and risk coverage from the adopter's real architecture, delivery risks, environments, platforms, and data constraints. Delivery-specific cases remain in each use case's `tests.md`.

## Test Levels

| Level | Purpose | Required when | Evidence |
| --- | --- | --- | --- |
| Unit/component | Isolated behavior and boundaries | Not configured | Not configured |
| Integration/contract | Data, API, dependency, and ownership boundaries | Not configured | Not configured |
| End-to-end | Critical user or operator flows | Not configured | Not configured |
| Manual/exploratory | Risks not economically automated | Not configured | Not configured |

## Flaky Tests And Exceptions

A flaky test is a quality finding, not a passing gate. Quarantine requires a tracked exception with owner, residual risk, mitigation, and expiry or review date.

## Delivery Application

Each `tests.md` pins the consumed Engineering System id/version, maps every `AC-*` to a validation method, identifies applicable risks, environments, test data, and platforms, and declares deviations or `None`.
