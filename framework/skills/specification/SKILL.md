---
name: specification
description: "Specification Skill. Use when an agent needs to Create the implementation contract that unifies product, UX, rules, data, APIs, analytics, security, and acceptance criteria in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Specification Skill

## Layer
Specification

## Responsibility
Create the implementation contract that unifies product, UX, rules, data, APIs, analytics, security, and acceptance criteria.

## Operating modes
- create: produce the first version of the artifact.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact and why it exists.

## Inputs
Approved use case; local and parent context.md files; sibling Specifications; business rules; design notes; Engineering System; Design System; decisions; approved demand classification.

## Outputs
`specification.md` root contract; applicable `contracts/*.md`; stable `REQ-*` and `AC-*` coverage; Delivery Level/Priority rationale; unresolved questions; decision candidates; context.md updates.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- This skill owns its generation resources: `assets/specification-template.md`, the concern-specific templates under `assets/*-contract-template.md`, and the machine-readable module registry at `references/contracts.yaml`. `assets/specification-contract-template.md` is retained only as the legacy shared base.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before substantive creation or material revision.

## Workflow
1. Read the local Use Case context, the complete parent chain, sibling Use Cases and Specifications, and confirm the artifact status and `specification_contract_version`.
2. Load reusable approved decisions, Engineering System contracts, Design System references, and code evidence named by the context. Separate observations, approved facts, inferences, and unresolved choices.
3. Confirm that the work covers one bounded interaction and whether it extends an existing contract; do not duplicate an existing requirement.
4. Read `references/contracts.yaml`, select modules from rigor and applicability, and record a concrete rationale for every inapplicable module.
5. Execute the selected modules in bounded passes: evidence and scope; product and behavior; UX, API, and data; security, quality, and observability; rollout and reversibility.
6. For each pass, use its concern-specific template, resolve discoverable facts, ask only focused blocking questions, and stop that pass when a product or architecture decision remains unresolved.
7. Give every testable requirement and acceptance criterion a stable ID; link sources, dependencies, risks, and verification methods; reject duplicate or orphan IDs.
8. Run an adversarial audit across contracts for contradictions, missing states, unsafe assumptions, uncovered failure modes, and inconsistent terminology. Route correctable gaps back to the owning pass.
9. Keep `specification.md` as the canonical index and cross-contract synthesis. Do not duplicate the detailed modular contracts in the root document.
10. Record decision candidates for high-impact or hard-to-reverse choices and keep the Specification `draft` while a blocking question remains.
11. Ask for approval before moving the artifact to the next ladder step.
12. Update context.md with the contract version, relations, reuse references, impacts, links, dependencies, questions, audit result, and status changes.

## Execution modes

- `sequential` is the default and works in every harness.
- `delegated` may assign disjoint contract passes to native subagents when the harness supports them. Each assignment receives only the approved evidence and declared module scope, writes only its own contract, and returns unresolved questions and requirement mappings. The parent Specification skill performs the final adversarial audit and remains accountable for the bundle.
- Falling back from delegated to sequential never weakens applicability, traceability, review, or approval gates.

## Quality checklist
- [ ] Preserves traceability to the parent artifact.
- [ ] Uses the correct template and naming conventions.
- [ ] States scope, non-goals, assumptions, and open questions.
- [ ] Includes Delivery Level and Priority and explains any change from the source feature/use case.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records meaningful decisions or decision candidates.
- [ ] Uses the v2 module registry and concern-specific templates when `specification_contract_version: 2`.
- [ ] Separates evidence, inference, assumptions, and decisions.
- [ ] Completes an adversarial cross-contract audit with no known material gap before `proposed`.
- [ ] Leaves a clear handoff for the next skill.

## Handoff
Next: ux-ui.

Pass forward approved specification, Delivery Level, Priority, open questions, decisions, dependencies, risks, and any remaining audit findings.
