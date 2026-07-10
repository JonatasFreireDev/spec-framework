---
name: new-framework-skill
description: Scaffold a new Codex skill in framework/skills/ following the repository's skill contract. Use when adding a specialist or orchestrator skill to the framework product, renaming an existing one, or bringing a skill's structure up to the canonical format.
---

# New Framework Skill Skill

## Purpose

Skills in `framework/skills/` are the framework's product: they ship to every adopter via `installFrameworkAssets` (into `.spec-framework/skills/` and `.codex/skills/`). A new skill that deviates from the canonical structure, hardcodes lab paths, or goes unregistered degrades every adopter repository at once. This skill scaffolds new framework skills correctly.

## Before creating

Answer these first; if any answer is unclear, stop and ask:

- **Specialist or orchestrator?** Specialists own exactly one canonical artifact's content. Orchestrators own flow, gates, sequencing, and handoff across artifacts — they never author artifact content themselves.
- **Which artifact does it own?** A specialist without an owned artifact, or an artifact already owned by another skill, means the skill should not exist — extend the existing owner instead.
- **Does the owned artifact have a template?** If the artifact is new, a matching template in `framework/template/` must be created in the same change.
- **Does this change the method?** Adding or reshaping a skill usually alters the framework contract — record an FDR (use the `fdr` skill) unless the change is purely editorial.

## Canonical SKILL.md structure

Create `framework/skills/<skill-name>/SKILL.md` (kebab-case folder, one skill per folder) with this exact section order, matching the existing skills:

```markdown
---
name: <skill-name>
description: "<Name> Skill. Use when Codex needs to <responsibility, lowercase continuation> in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# <Name> Skill

## Layer
<Discovery | Definition | Design | Planning | Execution | Validation | Governance>

## Responsibility
<One or two sentences. State what the skill owns and what it explicitly does not do.>

## Operating modes
- create: produce the first version of the artifact when this skill is generative.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact, finding, or decision in plain language.

## Inputs
<Semicolon-separated list of upstream artifacts and context.>

## Outputs
<Semicolon-separated list of artifacts or verdicts this skill produces.>

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.

## Workflow
<Numbered steps. Reference gates, statuses, and FDR routing rules where they apply.>

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] <Skill-specific checks.>
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: <skill or orchestrator that receives the output>.

Pass forward approved artifacts, findings, open questions, decisions, dependencies, risks, and required follow-up work.
```

## Path contract (non-negotiable)

Skills run in two environments with different layouts. Never hardcode either layout in the SKILL.md body:

- Write "the framework root's `FRAMEWORK.md`", "templates in framework/template/" (resolved via the path table in `framework/skills/README.md`) — not `.spec-framework/FRAMEWORK.md` or a lab-absolute path.
- Write "the active product root's `knowledge/conventions/gates.md`" — not `examples/events/...` or `product/...`.

## Registration checklist

After writing the SKILL.md:

1. Add the skill to the specialist or orchestrator list in `framework/skills/README.md`.
2. Check `FRAMEWORK.md` for skill rosters, flow diagrams, or step tables that must mention it, and `AGENTS.md` / `framework/AGENTS.framework.md` for flow references.
3. If it owns a new artifact: create the template in `framework/template/`, add the artifact to the canonical flow documentation, and update the readiness/validator rules if the artifact is gate-relevant.
4. Update the `Handoff` sections of neighboring skills whose flow now includes the new skill.
5. Record the FDR if the method changed.

No `package.json` or `framework-assets.mjs` changes are needed — `framework/skills/` is already copied wholesale.

## Verification

Run the `verify` skill. The path-contract test in `framework/tests/run-tests.mjs` will fail the build if any SKILL.md line mentions `FRAMEWORK.md` without "framework root", or mentions `knowledge/conventions/`, `knowledge/decisions/`, `.product/`, or `audits/security/` without "active product root". The description convention and README registration are not machine-checked — confirm them by inspection.
