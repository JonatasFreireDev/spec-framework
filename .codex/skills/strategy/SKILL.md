---
name: strategy
description: "Strategy Skill. Use when Codex needs to Translate vision into positioning, segments, metrics, roadmap logic, phase gates, Delivery Levels, and prioritization rules in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Strategy Skill

## Layer
Foundation

## Responsibility
Translate vision into positioning, segments, metrics, roadmap logic, phase gates, Delivery Levels, and prioritization rules.

## Operating modes
- create: produce the first version of the artifact.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact and why it exists.

## Inputs
Approved vision; personas or segment hypotheses; competitors; constraints; metric ideas.

## Outputs
strategy.md; personas.md; competitors.md; metrics.md; roadmap.md with Delivery Levels and Priorities; context.md updates.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in knowledge/templates/.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.

## Workflow
1. Read the parent context and confirm the artifact status.
2. Identify missing information, assumptions, conflicts, and dependencies.
3. Propose the artifact or revision using the matching template.
4. Record decision candidates for high-impact or hard-to-reverse choices.
5. Ask for approval before moving the artifact to the next ladder step.
6. Update context.md with new links, dependencies, questions, and status changes.

## Quality checklist
- [ ] Preserves traceability to the parent artifact.
- [ ] Uses the correct template and naming conventions.
- [ ] States scope, non-goals, assumptions, and open questions.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Defines how roadmap phases map to Delivery Levels and how P0/P1/P2/P3 priorities are assigned inside each level.
- [ ] Records meaningful decisions or decision candidates.
- [ ] Leaves a clear handoff for the next skill.

## Handoff
Next: 04-domain-architect.md

Pass forward approved artifacts, open questions, decisions, dependencies, risks, and any remaining audit findings.
