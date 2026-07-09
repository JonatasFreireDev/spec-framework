# Context: [artifact name]

```yaml
id: [DOMAIN-001 | GOAL-001 | FT-001 | UC-001 | SPEC-001 | TK-001]
type: [domain | goal | feature | use-case | specification | implementation-plan | execution-graph | task]
name: [human readable name]
status: [draft | proposed | approved | in_progress | implemented | validated | released | deprecated | superseded]
owner_skill: [skill file or role]
last_updated: [YYYY-MM-DD]
delivery:
  level: [L0 | L1 | L2 | L3 | L4 | L5 | N/A]
  priority: [P0 | P1 | P2 | P3 | N/A]
  rationale: [why this artifact belongs here]
```

## Purpose

[One paragraph explaining why this artifact exists and what decision or work it enables.]

## Parent Artifacts

- [id] - [path] - [relationship]

## Child Artifacts

- [id] - [path] - [relationship]

## Dependencies

- [id/path] - [why it is needed] - [blocking? yes/no]

## Related Artifacts

- [id/path] - [relationship]

## Canonical Documents

- Primary: [path]
- Specification: [path or N/A]
- Design: [path or N/A]
- Implementation plan: [path or N/A]
- Execution graph: [path or N/A]
- Tasks: [path or N/A]

## Decisions

- [DEC-XXX] - [summary] - [status]

## Assumptions

- [Assumption that must be validated or carried forward.]

## Open Questions

- [Question] - owner: [role] - blocks: [artifact/status]

## Handoff

Next recommended skill: [skill]
Required reading before next step:
- [path]
