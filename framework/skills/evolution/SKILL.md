---
name: evolution
description: "Evolution Skill. Use when an agent needs to Turn observations, metrics, feedback, and audit findings into explicit improvement candidates for approval in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Evolution Skill

## Layer
Audit

## Responsibility
Turn observations, metrics, feedback, and audit findings into explicit improvement candidates for approval.

## Operating modes
- create: produce the first version of the artifact when this skill is generative.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact, finding, or decision in plain language.

## Inputs
Feedback; metrics; QA findings; audit reports; roadmap; product principles; imported demand mappings; existing Feature and Use Case contexts.

## Outputs
Evolution proposals; demand classification; context relation updates; opportunity notes; expected impact; required approvals.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- The canonical template belongs to the skill that generates the artifact; read that skill's `assets/` directory.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before proposing or materially revising improvement candidates.

## Workflow
1. Read the relevant context and identify artifact status.
2. Read the complete parent chain and sibling Features, Use Cases, Specifications, decisions, Engineering System, and Design System before proposing a destination.
3. Classify the demand as an extension, new Use Case, new Feature, new Goal, new Domain, or non-delivery item. Keep the classification as a proposal until the human resolves ambiguity.
4. Compare the demand against the framework, template, approved decisions, and existing code evidence.
5. Separate verified facts from assumptions and recommendations; record `extends`, `reuses`, `depends_on`, `impacts`, and `supersedes` relations only with evidence.
6. Report gaps, conflicts, dependencies, and risks with file-level evidence when possible.
7. Ask for approval before changing canonical product artifacts.
8. Update the owning context.md, decision indexes, and import traceability only when the change is approved.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Distinguishes blockers from suggestions.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: documentation-writer.

Pass forward approved artifacts, findings, open questions, decisions, dependencies, risks, and required follow-up work.
