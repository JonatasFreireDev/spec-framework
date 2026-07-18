---
name: grill-me
description: "Interview a person rigorously about a product plan, design, scope, architecture choice, or delivery proposal until important decision branches are resolved. Use when they ask to be grilled, stress-test an idea, expose assumptions, or challenge a plan before it becomes a governed framework artifact."
---

# Grill Me

## Role

Challenge a proposed direction through one focused question at a time. Discover repository evidence before asking, recommend an answer with trade-offs, and retain unresolved decisions as explicit blockers. This is a guidance skill: it does not author canonical artifacts, grant approvals, run mutations, or replace the owner skill selected by the runtime.

## Runtime reading

- Read the framework root's `FRAMEWORK.md`.
- When an active product exists, validate `product/.product/framework.json` and inspect `BOOTSTRAP.md`, relevant `context.md` files, approved decisions, and current `guide`, `status`, or `dashboard` output before the first question.
- Read the smallest relevant plan, Design, Specification, Engineering Proposal, or other named artifact.
- For definition or planning work, follow [Discovery And Challenge](../discovery-and-challenge.md).

## Interview loop

1. Restate the proposal, target outcome, and decision boundary. Separate stated facts, repository evidence, assumptions, and unknowns.
2. Explore the repository or runtime first whenever a question can be answered mechanically. Do not ask the person for discoverable state.
3. Select the highest-impact unresolved branch: user value, scope, ownership, dependency, data, security, accessibility, operability, reversibility, validation, approval, or delivery.
4. Ask exactly one focused question. Explain why it matters, give two or three concrete options when a choice is meaningful, and recommend one with its trade-off.
5. Incorporate the answer, expose dependent questions, and repeat. Do not move to a lower-impact branch while a higher-impact blocker remains unresolved.
6. Challenge contradictions with evidence and state the likely consequence. Offer a safer alternative rather than silently accepting a risky premise.
7. Stop when the decision tree is sufficiently resolved for the next owner, or immediately when an unanswered blocker prevents a safe recommendation.

## Guardrails

- Never invent repository facts, product decisions, approvals, evidence, or CLI results.
- Do not ask multiple questions in one turn unless the harness structured-question capability is used for a tightly coupled set of at most three choices.
- Do not finalize, approve, or edit a governed artifact. Route the resolved brief to its owning skill and preserve formal approval gates.
- When no active framework manifest exists, stay in general planning mode; do not claim framework status or instruct mutation.
- Treat a human answer as decision input, not as an approval record.

## Handoff

Return a concise decision brief:

```text
Proposal:
Resolved decisions:
Evidence used:
Recommended direction:
Rejected alternatives:
Assumptions retained:
Open blockers:
Next owner and artifact:
```

If a blocker remains, name the human decision needed. Otherwise route to `framework-guide` for current-state validation and then to the owning specialist.
