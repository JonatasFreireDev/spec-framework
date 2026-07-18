---
name: dependency-analyzer
description: "Dependency Analyzer Skill. Use when an agent needs to Expose explicit and implicit dependencies across domains, features, use cases, specs, tasks, releases, and decisions in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Dependency Analyzer Skill

## Layer
Audit

## Responsibility
Expose explicit and implicit dependencies across domains, features, use cases, specs, tasks, releases, and decisions.

## Operating modes
- create: produce the first version of the artifact when this skill is generative.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact, finding, or decision in plain language.

## Inputs
Artifact subtree; context files; execution graphs; implementation plans; release notes.

## Outputs
Dependency report; blocked items; dependency graph notes; parallelization opportunities.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- The canonical template belongs to the skill that generates the artifact; read that skill's `assets/` directory.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Workflow
1. Read the relevant context and identify artifact status.
2. Compare the artifact against the framework, template, and approved decisions.
3. Separate verified facts from assumptions and recommendations.
4. Report gaps, conflicts, dependencies, and risks with file-level evidence when possible.
5. Ask for approval before changing canonical product artifacts.
6. Update context.md or decision indexes only when the change is approved.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Distinguishes blockers from suggestions.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: impact-analyzer.

Pass forward approved artifacts, findings, open questions, decisions, dependencies, risks, and required follow-up work.
