# Product Engineering Framework

This document consolidates the architecture of the AI-driven Product Engineering framework. It replaces the long conversational history with a navigable source of truth that agents such as Codex can execute against.

The goal is not just to organize documents. The goal is to create a pipeline where product, specification, planning, execution, and audit form a single system.

## 1. Thesis

This framework treats documentation as engineering infrastructure.

Instead of asking an AI to "create files" or "implement a feature" from loose context, the product flows through an explicit chain:

```text
Problem -> Vision -> Strategy -> Domain -> User Goal -> Feature -> Use Case -> Specification -> Design -> Implementation Plan -> Execution Graph -> Tasks -> Code -> Validation -> Audit
```

The Specification is the central contract. It unites product, UX, rules, architecture, data, analytics, security, tests, and acceptance criteria. Tasks are not loose checklists: they are executable units derived from a Specification, ordered by an Execution Graph.

Design is a planning artifact when the delivery has an interface. It translates the Specification into visual flow, states, wireframes, or reviewable mockups before the Implementation Plan, so that engineering does not discover the experience only during implementation.

## 2. Principles

### Product first

Every technical decision must maintain traceability to a product need. The tree always originates from the problem, passes through the vision, and reaches the code through domains, goals, features, and use cases.

### Domain driven

Domains are the center of functional documentation. There is no global features folder as the primary source. A feature belongs to a User Goal, and a User Goal belongs to a Domain.

### Specification driven

Stories may exist for backlog or communication purposes, but they are not the primary contract for AI. The primary contract is the Specification, because it describes the complete behavior and reduces ambiguity before implementation.

### Context driven

Every important level must have a `context.md`. This file summarizes the object, its parents, children, dependencies, decisions, risks, and next relevant documents.

### Knowledge graph, not just folders

The folder tree is only the human-facing interface. The real model must function as a graph: artifacts have ids, parents, children, dependencies, relations, and consumers.

### Approval gates

Agents can propose, but relevant changes must go through explicit approval when they alter scope, architecture, business rules, risks, roadmap, or delivery commitments.

## 3. Conceptual Model

### Problem

Defines the pain, the opportunity, and the market context. It is the justification root of the product.

### Vision

Defines the product we want to build, for whom, why now, and which principles guide decisions.

### Strategy

Defines positioning, segments, metrics, trade-offs, roadmap, and criteria to advance or pause.

### Domain

Groups a coherent area of the business or product, such as `users`, `groups`, `events`, `friendship`, or `payments`.

### User Goal

Replaces the generic notion of "capability". Represents the user's stable objective within a domain, for example: "join an event", "find people", or "manage profile".

### Feature

Is a concrete solution that helps a User Goal. Features can enter, exit, evolve, be sliced, or be replaced.

### Use Case

Is a verifiable interaction of the feature. The Use Case is the point where product and engineering meet to generate an implementable Specification.

### Delivery Level

Defines the delivery level at which an artifact must enter. The level answers "when does this need to exist in the product", without replacing priority within the level.

Canonical levels:

- `L0 Foundation`: baseline without which the product or the pipeline cannot sustain safe deliveries.
- `L1 Walking Skeleton`: smallest end-to-end flow that proves the core value.
- `L2 Core Loop`: main cycle that generates recurring value for the user.
- `L3 Trust, Safety and Quality`: trust, security, privacy, moderation, accessibility, and experience quality.
- `L4 Operations and Scale`: operations, support, observability, admin, and scale.
- `L5 Growth and Optimization`: growth, experiments, personalization, and optimizations.

### Priority

Defines relative urgency within a Delivery Level:

- `P0`: blocks the current level.
- `P1`: required to consider the level ready.
- `P2`: important, but does not block the level's delivery.
- `P3`: improvement, polish, or optimization.

### Specification

Is the source of truth for implementation. It must cover product, flow, UI, APIs, data, permissions, analytics, security, performance, accessibility, errors, edge cases, observability, and acceptance.

### Design

Translates the Specification into a verifiable user experience: visual flow, navigation, wireframes, mockups, states, accessibility, and alignment with the design system. When the feature has no interface, the artifact must explicitly record `Not applicable` and explain why.

### Implementation Plan

Translates the Specification into technical strategy. Thinks like a Tech Lead: sequencing, phases, dependencies, risks, slices, migrations, backend, frontend, tests, and rollout.

### Execution Graph

Represents the tasks as a DAG. Each node is an executable unit with explicit dependencies and a `path` to its canonical file in `tasks/<task-id>.md`. This enables safe parallelism between agents and makes it clear when something is blocked.

### Task

Executable unit derived from the Specification and the Execution Graph. A task must be small enough for implementation, testing, review, and rollback. Each task lives in its own file at `tasks/<task-id>.md`; this file is the canonical source for status, contract, Delivery Level/Priority, links to code, and evidence. `tasks.md` is only a generated index for human navigation.

### Code Link

The framework uses a monorepo model for product delivery: documentation, code, and evidence live in the same product repository, and this `spec-framework` repository remains as a template/laboratory. Links between task and code must use paths relative to the repository when pointing to internal files; PRs can use a URL or an external identifier.

A task in `implemented` must record immutable working-tree evidence: `Branch`, `Base commit`, `Changed paths`, `Diff hash`, tests, and applicable gate results. Commits are intentionally deferred until independent Code Review and QA approve the same diff hash. A task in `validated` or `released` must record `Commits`, `Code paths`, `PR` when policy requires it, approved `Test status`, and concrete evidence such as gate logs, CI URL, screenshots, or QA evidence.

### Rigor Tier

Use cases declare `rigor_tier` in `context.md` to adjust documentary rigor to risk:

- `S`: small and low risk; requires specification, tasks, and tests.
- `M`: normal product flow; also requires design, implementation plan, and execution graph.
- `L`: critical or sensitive flow; also requires analytics, audit, QA evidence, and Security Review.
- `N/A`: structural example or placeholder without product scope.

Automatic Tier L triggers: auth, permissions, roles, payment, PII, upload, UGC, public surface, or migration that touches RLS/policies. A tier change requires an approval record for the use case, but not a new DEC when the policy remains the same.

### Identity And Moves

Every product object with its own folder must declare `slug` in `context.md`. The slug is born at artifact creation, matches the folder name, and remains immutable even if the human title changes.

IDs are unique within the parent's scope. When there is risk of ambiguity, references must combine ID and path. `.product/ids.json` records the identity policy; it must not be used as a global counter to allocate new artifacts.

Moving an artifact requires tooling:

```bash
spec-framework move --from <old-path> --to <new-path>
```

The move script rewrites Markdown links and paths in JSON that are mechanically resolvable. Mentions in free text are reported for human review, not rewritten automatically.
## 4. Folder Structure

Canonical structure:

```text
product/
  .product/
    state.json
    decisions.json
    roadmap.json
    ids.json
    history/

  foundation/
    problem/
      problem.md
      opportunities.md
      researches/
      interviews/
      context.md
    vision/
      vision.md
      principles.md
      north-star.md
      context.md
    strategy/
      strategy.md
      personas.md
      competitors.md
      metrics.md
      roadmap.md
      context.md

  knowledge/
    imports/
      sources/
      runs/
        IMPORT-NNN/
          inventory.json
          import-plan.json
          mapping.json
          conflicts.md
          import-report.md
    glossary/
    business-rules/
    conventions/
    decisions/
    patterns/
    prompts/
    templates/
    examples/

  domains/
    <domain>/
      context.md
      domain.md
      decisions.md
      goals/
        <goal>/
          context.md
          goal.md
          journeys.md
          features/
            <feature>/
              context.md
              feature.md
              use-cases/
                <use-case>/
                  context.md
                  use-case.md
                  specification.md
                  implementation-plan.md
                  execution-graph.json
                  tasks/
                    <task-id>.md
                  tasks.md
                  tests.md
                  analytics.md
                  design.md
                  audit.md

  design/
  engineering/
  audits/
  releases/
  skills/
```

In an adopter product repository, the structure above must live under `product/`. The assets that teach how to run the framework must live outside it, in `.spec-framework/`:

```text
repo/
  README.md
  BOOTSTRAP.md
  .spec-framework/
    FRAMEWORK.md
    AGENTS.framework.md
    decisions/
    skills/
    templates/
    manifest.json

  product/
    .product/
    foundation/
    knowledge/
    domains/
    design/
    engineering/
    audits/
    releases/
```

In this `spec-framework` repository, the structure exists in three explicit areas:

- `starter/` represents the clean skeleton that must be copied into new product repositories, already separated between `.spec-framework/` and `product/`.
- `examples/events/` contains the worked product instance used as learning material and validation fixture.
- `framework/` contains the executable framework core: audits, decisions, skills, templates, validators, distributable tools, framework-only tests, and adoption guidance. The repository root retains only entry points, packaging metadata, scripts, examples, and starter infrastructure.

New products must not copy the entire `spec-framework` root; they must start from `starter/` and install the framework assets into `.spec-framework/` per `framework/adoption.md`.

## 5. Context.md

Every `context.md` must let an AI understand where it is, what it needs to read, and what the safe next step is.

Minimal template:

```yaml
id: FT-023
type: feature
name: QR Code Check-in
status: draft
owner_skill: feature-ai
slug: qr-code-check-in
rigor_tier: L

parents:
  - GOAL-003

children:
  - UC-001
  - UC-002

depends_on:
  - FT-008
  - DOMAIN-users

used_by:
  - RELEASE-001

related:
  - FT-055

documents:
  canonical: feature.md
  specification: use-cases/qr-code-check-in/specification.md
  design: use-cases/qr-code-check-in/design.md
  implementation_plan: use-cases/qr-code-check-in/implementation-plan.md
  execution_graph: use-cases/qr-code-check-in/execution-graph.json

delivery:
  level: L1
  priority: P0
  rationale: Without check-in, the L1 event walking skeleton does not close end-to-end.

open_questions:
  - How to expire QR codes without hurting offline users?

decisions:
  - DEC-014
```

## 6. Specification Driven Development

The flow for a new feature must be:

```text
Feature -> Use Cases -> Specification -> Design -> Implementation Plan -> Execution Graph -> Tasks -> Implementation -> QA Evidence -> Security Review -> Review -> Audit -> Release
```

The closed delivery flow is:

```text
Domain -> Domain Evolution -> Feature Selection -> New Feature -> Use Cases
-> Specification Contracts -> Design -> Technical Discovery -> Architecture Gate
-> Implementation Plan -> Execution Graph -> Tasks -> Code Runner
-> Code Review -> QA -> Commit Crafter -> PR Finalizer
```

`specification.md` remains the root contract. Large concerns live under `contracts/` (`product`, `behavior`, `ux`, `api`, `data`, `security`, `quality`, `observability`, and `rollout`) and use stable `REQ-*` and `AC-*` identifiers. Tier S requires behavior and quality; Tier M adds product, UX, API, data, rollout, and Technical Discovery; Tier L adds security and observability. An inapplicable contract must say `Not applicable` with rationale.

Design is mandatory for any use case with an interface. For deliveries without UI, `design.md` must exist as a short artifact with `Not applicable`, justification, and impacts on accessibility, observability, or operations when relevant.

QA Evidence and Security Review are validation gates. QA Evidence proves that acceptance criteria, tasks, flows, edge cases, regression, accessibility, observability, and security controls were verified. Security Review evaluates authentication, authorization, privacy, abuse, sensitive data, tokens, logs, dependencies, rollout, rollback, and residual risk. Security Review must also read the product's security baseline in `knowledge/conventions/security-baseline.md` and active threats in `audits/security/threat-register.md`. An artifact must not reach `validated` or `released` when there is a QA or security blocker.

QA Evidence must bring the real evidence back to the use case: branch, commits, PR, code paths, commands or test methods, gate logs, CI URL when available, and screenshots when the delivery has a visual surface.

The Specification must answer:

- What exactly must happen?
- What user problem does this solve?
- What is in and out of scope?
- What flows, states, and errors exist?
- What business rules apply?
- What APIs, data, and permissions are needed?
- What analytics events and logs must exist?
- What security, privacy, and abuse risks exist?
- How will the delivery be tested and accepted?

Mandatory sections:

```text
Product context
User goal
Delivery level
Priority
Feature scope
Non-goals
Use cases
Business rules
UX flow
UI states
API contracts
Data model
Events
Analytics
Permissions
Security
Performance
Accessibility
Error states
Edge cases
Observability
Rollout strategy
Feature flags
Acceptance criteria
Open questions
```

Security Review is not an absolute promise of risk absence. The gate's role is to ensure that all defined controls were verified with evidence, that blockers are resolved, and that residual risks are documented and approved by humans when relevant.

Documentary rigor is proportional to the use case's tier. Tier S avoids heavy artifacts when design, analytics, or audit are `Not applicable`; Tier L requires QA Evidence and Security Review by default.

## 6.1. Delivery Prioritization

Every executable artifact must declare `Delivery Level` and `Priority`. The level organizes the roadmap by product maturity; priority orders work within the level.

Mandatory fields in Domain, User Goal, Feature, Use Case, Specification, Implementation Plan, Execution Graph, and Task:

```yaml
delivery:
  level: L1
  priority: P0
  depends_on:
    - FT-008
  rationale: Explains why this delivery belongs to this level and why this priority was assigned.
```

Rules:

- `Delivery Level` is not a date promise.
- `Priority` should only be compared within the same level.
- An `L3/P0` delivery can be critical for trust, but it still must not jump ahead of an `L1` delivery that closes the walking skeleton.
- Dependencies can pull a technical delivery to an earlier level, as long as the `rationale` explains why.
- Changing `level` or `priority` alters delivery commitment and must go through an approval gate.

## 6.2. Design Driven Handoff

Design is born after the Specification is approved and before the Implementation Plan.

Inputs:

- Approved Specification.
- Design system, UX patterns, and neighboring screens.
- Delivery Level and Priority of the delivery.

Outputs:

- `design.md` in the use case, with visual flow, states, accessibility, components, displayed data, and links to mockups when they exist.
- Mockups or wireframes in `product/design/` or the product's canonical design directory, when the delivery needs visual reference.
- UX review recorded before the Implementation Plan when UI is relevant to acceptance.

Gates:

- Without an approved Specification, do not generate design.
- Without an approved Design, or one marked `Not applicable`, do not generate the Implementation Plan.
- Blocking UX findings go back to Specification or Design before proceeding.

## 7. Implementation Plan

The Implementation Plan is created after the Specification and the Design, and before the tasks. It must not write code. It must define the build strategy.

Before planning, `technical-discovery.md` must map applicable requirements to the real codebase and stable knowledge in `engineering/`. Its Architecture Gate must reference an approved decision or state `Not required` with concrete rationale.

Recommended sections:

- Technical objective
- Technical scope
- Delivery Level and Priority, inherited or adjusted with justification
- Dependencies
- Phases
- Delivery sequence
- Risks
- Test plan
- Rollback plan
- Probable files or modules
- Decisions that need an ADR
- Candidate tasks

Example phases:

```text
1. Data model and migration
2. Server-side rules and permissions
3. Backend services and API
4. Frontend states and forms
5. Analytics and observability
6. Tests and fixtures
7. QA, review and release
```

## 8. Execution Graph

The Execution Graph is a DAG. It defines dependencies between tasks and enables parallel execution by agents. Each node references the task's canonical file by `path`.

Graph lifecycle is `draft → proposed → materialized → approved`. A draft or proposed graph validates its DAG and task contracts before task files exist. Confirmed materialization creates the missing canonical task files and generated index atomically. From `materialized` onward every node path must exist; approval happens only after Graph + Tasks validation. This avoids making task existence a precondition for approving the plan that defines those tasks.

Example:

```json
{
  "id": "GRAPH-001",
  "sourceSpecification": "SPEC-001",
  "nodes": [
    {
      "id": "TK-001",
      "path": "tasks/TK-001.md",
      "title": "Create event tables and policies",
      "type": "database",
      "dependsOn": [],
      "writeScope": ["supabase/migrations"],
      "sharedResources": ["local database schema"]
    },
    {
      "id": "TK-002",
      "path": "tasks/TK-002.md",
      "title": "Create event service",
      "type": "backend",
      "dependsOn": ["TK-001"],
      "writeScope": ["src/server/events"],
      "sharedResources": []
    },
    {
      "id": "TK-003",
      "path": "tasks/TK-003.md",
      "title": "Create event form UI",
      "type": "frontend",
      "dependsOn": ["TK-002"],
      "writeScope": ["src/app/events"],
      "sharedResources": []
    },
    {
      "id": "TK-004",
      "path": "tasks/TK-004.md",
      "title": "Instrument analytics",
      "type": "analytics",
      "dependsOn": ["TK-002", "TK-003"],
      "writeScope": ["src/analytics/events"],
      "sharedResources": ["analytics event catalog"]
    }
  ]
}
```

Rules:

- A task can only start when its dependencies are approved.
- Parallel tasks must have separate write scopes.
- Every node must declare `writeScope`: paths or modules the task may create or change.
- `sharedResources` declares generated or shared resources, such as indexes, locales, local schema, local database, generated contracts, or catalogs. Two parallel nodes that contend for the same resource must become a dependency or be merged.
- Two nodes are parallel when there is no dependency path between them in the DAG. Parallel nodes must not have overlapping `writeScope`.
- Path overlap is prefix-based: `src/` covers `src/foo.ts`.
- The rollout of the `writeScope` check follows FDR-003: Phase A reports warnings; Phase B, after approved migration of existing graphs, can promote the same findings to errors.
- Every task points to its source Specification.
- Every node must point to the canonical file at `tasks/<task-id>.md`.
- Snapshots in the graph, such as `title` and `type`, are allowed only when they match the referenced task file.
- Every dependency change updates the graph.
- QA failures can create new nodes in the graph.

The CLI can operate, but not automatically execute, the graph:

```text
spec-framework graph ready --graph <path>
spec-framework graph claim --graph <path> --task TK-001 --agent codex
spec-framework graph release --task TK-001 --agent codex
spec-framework graph complete --graph <path> --task TK-001 --agent codex
```

Runtime v2 leases live in `.product/claims/<task-id>.json`; `.product/claims.json` remains a v1 compatibility index during migration. Leases provide expiring operational ownership with heartbeat and recovery; they do not grant artifact approval or permission to exceed `writeScope`.
## 9. Skills

Runtime v2 makes execution resumable and safely parallel. Each `WORK-NNN` is a directory containing identity, state, structured handoffs, checkpoints, command plans, and evidence. Task ownership is a renewable lease with heartbeat and expiry; isolated tasks use one Git worktree under `.worktrees/WORK-NNN/TK-NNN`. The scheduler computes deterministic conflict-free waves from DAG dependencies, `writeScope`, and `sharedResources` but does not spawn agents. Command plans store argv rather than shell strings and initially permit only R0 read-only and R1 local-temporary operations. Validated task commits are integrated locally in DAG order; conflicts stop for human resolution and Integrated QA is mandatory.

Runtime commands include `runtime`, `resume`, `handoff`, `checkpoint`, `lease`, `commands`, `schedule`, and `integrate`.

Skills are specialists. They can operate in modes such as `create`, `update`, `audit`, `evolve`, `explain`, `compare`, and `refactor`, but each must have a clear responsibility.

### Foundation

- Problem Discovery AI: discovers pains, opportunities, and evidence.
- Vision AI: creates or revises vision, principles, and north star.
- Strategy AI: defines strategy, segments, metrics, and roadmap.
- Domain Architect AI: models domains and boundaries.
- User Goal AI: models user goals within domains.

### Product Design

- Journey AI: maps journeys.
- Feature AI: creates and evolves features.
- Use Case AI: details verifiable interactions.
- UX/UI AI: defines flows, states, wireframes, mockups, design system, and accessibility.
- UX Review AI: reviews design against the design system, UX principles, accessibility, and state coverage.

### Specification And Planning

- Specification AI: creates the central implementation contract.
- Implementation Planner AI: transforms Specification into a technical plan.
- Execution Graph AI: transforms the plan into an execution DAG.
- Task AI: generates small, testable, and traceable executable tasks.

### Engineering And Validation

- Code Runner AI: implements exactly one approved task per invocation, in TDD, respecting `writeScope`, reading gates from `knowledge/conventions/gates.md`, stopping when green, and without committing.
- QA AI: validates behavior, tests, edge cases, performance, and evidence matrix.
- Code Review AI: reviews implementation read-only through the lenses of completeness, adherence, and quality; findings use severity and follow FDR-006 routing.
- Security Review AI: evaluates authentication, authorization, privacy, abuse, data exposure, tokens, logs, dependencies, rollout, and residual risk.
- Threat Modeler AI: proactively models threats at the product, domain, or feature family level; maintains the security baseline and threat register to feed Security Review.
- Commit Crafter AI: turns verified changes into atomic local commits by concern, following `knowledge/conventions/commits.md`, without push.
- PR Finalizer AI: prepares or opens a PR with evidence and required links, following `knowledge/conventions/pull-requests.md`, without merge.

### Audit

- Gap Finder AI: looks for gaps.
- Conflict AI: looks for contradictions.
- Dependency AI: finds implicit dependencies.
- Impact Analysis AI: measures the effect of changes.
- Evolution AI: suggests improvements.
- Documentation AI: updates docs.
- Product Historian AI: records decisions.

## 10. Orchestrators

Orchestrators do not create primary artifacts. They control flow, order, gates, and handoffs.

### Domain Evolution Orchestrator

Coordinates approved goals, journeys, opportunity gaps, candidate features, delivery slices, dependency/impact analysis, and explicit human feature selection. It hands the selected feature to New Feature Orchestrator.

### Existing Product Import Orchestrator

Coordinates existing epics, PRDs, and other source documents through inventory, classification, reconciliation, approval, and draft materialization. Sources remain evidence, conflicts are never resolved silently, and canonical artifacts are not created before explicit human approval.

```text
Sources -> Inventory -> Candidates -> Conflicts -> Mapping -> Approval -> Draft artifacts
```

### Product Orchestrator

Creates a product from scratch:

```text
Problem -> Vision -> Strategy -> Domains -> User Goals -> Roadmap
```

### New Feature Orchestrator

Receives a candidate feature and drives:

```text
Impact -> Feature -> Use Cases -> Specification -> Design -> Plan -> Graph -> Tasks
```

### Audit Orchestrator

Runs batch audits:

```text
Gap -> Conflict -> Dependency -> Impact -> Consistency
```

### Evolution Orchestrator

Groups candidate improvements, asks which will be approved, and creates an evolution plan.

### Documentation Orchestrator

Keeps `context.md`, indexes, templates, decisions, and derived artifacts synchronized.

### QA Orchestration

QA is an independent, read-only verifier. QA re-executes the gates declared in `knowledge/conventions/gates.md` whenever possible, records real output or an explicit limitation, and routes blockers to the appropriate path instead of fixing code. Specialized routes such as bug-fixer and code-runner will be formalized by future evolutions.

### Implementation Orchestration

Code Runner receives one approved task at a time. If the task requires a change outside `writeScope`, an unapproved decision, a missing gate command, or a Specification gap, implementation stops and reports the blocker. After any code change, the flow returns to independent QA.

### Failure Routing

Failures follow FDR-006. A defect, regression, vulnerability with a clearly expected behavior, or production error goes to Bug Fixer. A missing test, hollow test, or missing negative/permission coverage goes back to QA or the test owner. Incomplete implementation or out-of-contract work goes back to Code Runner. A decision gap or ambiguous rule goes to Product Historian and a human. After any code change, the flow re-enters QA; a red gate cannot be skipped. Every gate or finding has a cap of three automated attempts before escalating to a human.

### Review Orchestration

Code Review is a read-only gate before validation and release of executable work. The review evaluates completeness against Specification/tasks, adherence to approved contracts, and implementation quality. `blocker` or `required_fix` findings need a route and owner via FDR-006; Code Review does not fix code or approve QA.

### Threat Modeling

Threat Modeler acts before and alongside Security Review. It models threats at the product, domain, goal, or feature family level, records recurring rules in the security baseline, and maintains `audits/security/threat-register.md` with scenarios, mitigations, owners, evidence, and residual risks. It does not validate release; Security Review consumes that context to assess a specific delivery.

### Delivery Orchestration

Commit Crafter creates atomic local commits only when explicitly invoked, separating concerns and following `knowledge/conventions/commits.md`; it does not push. PR Finalizer verifies gates, QA Evidence, Code Review, and Security Review when applicable, prepares or opens a PR with evidence links following `knowledge/conventions/pull-requests.md`, records the PR back onto the tasks when appropriate, and does not merge.

### Release Orchestrator

Before release, checks:

- gaps
- conflicts
- docs
- specs
- design
- tasks
- tests
- QA
- QA evidence
- review
- security review for executable deliveries, with depth proportional to risk

## 11. Approval Gates

Every step must end with a clear state:

```text
draft
proposed
approved
in_progress
implemented
validated
released
deprecated
superseded
```

Rules:

- `draft`: artifact created, still incomplete.
- `proposed`: ready for human or audit review.
- `approved`: can feed the next stage.
- `in_progress`: being implemented.
- `implemented`: code or artifact was produced.
- `validated`: passed QA, review, Security Review when applicable, and has sufficient evidence.
- `released`: reached the user or target environment.
- `deprecated`: must not guide new implementations.
- `superseded`: replaced by another artifact.

Mandatory transitions:

- `proposed`: does not require an approval record, but must not advance from an incomplete parent gate.
- `approved` and later states: require a corresponding approval record in `.product/history/`, with `artifact_id`, `path`, `content_hash`, `status_granted`, `approved_by`, `approved_at`, and `notes`.
- `approved -> in_progress`: requires an approved task or an explicit prototype/draft exception.
- `in_progress -> implemented`: requires structured working-tree evidence in the task file: branch, base commit, changed paths, diff hash, tests, and gate results. It does not require a commit.
- Code Runner can produce code and technical evidence, but does not commit, push, merge, create approval records, or approve QA.
- Commit Crafter can create local commits when explicitly invoked, but does not push, merge, or create approval records.
- Bug Fixer reproduces the defect with a failing test before fixing, fixes the root cause with a minimal change, leaves a permanent regression test, and returns to QA.
- `implemented -> validated`: requires approved QA Evidence with no blockers; requires approved Security Review when there is code, data, permissions, tokens, API, payments, uploads, messaging, search, admin, sensitive analytics, or any privacy/abuse risk.
- `implemented -> validated`: requires approved Code Review with no `blocker` or `required_fix` findings for executable deliveries.
- `implemented -> validated`: for an individual task, requires Code Review and QA approval over the same current diff hash, followed by Commit Crafter evidence, code paths, approved test status, and concrete evidence such as gate logs, CI URL, screenshots, or QA evidence. PR is required when repository policy declares it mandatory.
- PR Finalizer can prepare or open a PR when the hard gates are green or when the user explicitly requests a draft/prototype; it does not merge.
- Technical gates for the product live in `knowledge/conventions/gates.md`. Skills that execute or verify code must read that file and record the real output of applicable gates.
- `validated` and later states require non-placeholder QA Evidence for applicable gates. When a gate cannot be executed due to environment unavailability, QA must explicitly record the limitation instead of forging evidence.
- Deliveries with a visual surface require proportional visual evidence, such as a local screenshot or CI artifact, plus basic accessibility verification: role/label, focus, touch target, and contrast.
- `validated -> released`: requires the Release Orchestrator, an audit with no blockers, Security Review with no blockers, accepted residual risks, and defined rollback/monitoring.
- QA can block validation when any acceptance criterion, task, security control, critical regression, or mandatory evidence is missing.
- Security Review can block validation and release when there is an authorization failure, data leak, unapproved permission decision, exposed secret, unmitigated abuse, insecure logging, or high residual risk without a human decision.
- Blocking findings in QA Evidence, Security Review, or audit need a route and owner before approved artifacts can guide validation or release.

Approval records use a SHA-256 hash of the whole file with content normalized to LF and no trailing whitespace per line. They provide auditability and a mechanical gate, not cryptographic proof of human approval.

Staleness is a condition derived by the validator, not an editable status. Downstream artifacts record the hashes of source artifacts in `.product/derivations.json`; if the source content changes, the downstream becomes stale and must not advance through gates until it is regenerated or re-approved.

## 12. Decisions

Relevant product decisions must be recorded in `knowledge/decisions/` and indexed in `.product/decisions.json`.

`DEC-*` is the single product decision identity; an ADR is a DEC whose type is `architecture`. Decisions declare `type`, artifact/path `scope`, and optional structured `workflowEffects` (`requiredTaskTypes`, `requiredGates`, `requiredEvidence`, `requiredWriteScopes`, and `sharedResources`). A DEC can unblock an Architecture Gate only when it exists, is indexed, is approved, has a current hash-matching approval record, and applies to the affected scope. Decision prose is never executable. Effects constrain Graph, Tasks, Task Readiness, and configured gates; they never silently generate work.

Framework or method decisions live in `framework/decisions/FDR-*` or as explicit amendments to this document. Skill contracts, gates, writeScope, QA policies, failure routing, commit policy, validators, and orchestration rules must not be recorded in `knowledge/decisions/`, because that folder is reserved for the adopter product.

A decision must be created when it:

- changes structural architecture;
- changes an important business rule;
- alters security, privacy, payment, or permissions;
- creates a relevant external dependency;
- chooses a hard-to-reverse strategy;
- replaces a previous decision.

## 13. Audit

Audits must not create new product content as default behavior. They analyze and report.

Types of audit:

- Gap: what is missing?
- Conflict: what contradicts another artifact?
- Dependency: what depends on what?
- Impact: what changes if this changes?
- Consistency: do names, states, ids, and Markdown links match and resolve to existing files?
- Security: is there risk of improper access, abuse, or leakage?
- UX: does the experience close for the persona?

QA and Security Review must produce or reference evidence. Threat Modeler can maintain a living threat register for risks that span multiple deliveries. Audits can verify the coherence of that evidence, but must not declare a delivery safe when the specialized gates are missing or blocked.

Expected output:

```text
Verdict: approved | approved_with_notes | blocked
Findings
Evidence
Required fixes
Suggested improvements
Residual risk
```

## 14. Evolution Engine

The framework must allow continuous evolution. Improvements do not go directly into the product; they become candidates.

Flow:

```text
Observation -> Opportunity -> Proposal -> Impact Analysis -> Approval -> Updated Specification -> Updated Plan -> Updated Graph -> Tasks
```

This prevents AI suggestions from turning into silent scope.

## 15. How To Use With Codex

In this repository, Codex is working on the framework itself. In product repositories, Codex must operate on an instance created from `starter/`, reading the method in `.spec-framework/` and writing product artifacts in `product/`.

In a newly initialized product repository, read `BOOTSTRAP.md` first. It explains the ordered foundation gates and distinguishes a structurally valid starter from a product that is ready for implementation.

During `init`, choose the repository's starting point. This choice customizes `BOOTSTRAP.md`; it does not remove skills, orchestrators, artifacts, rigor requirements, or approval gates. When starting from existing documents, the CLI creates a source inventory and an analysis-only import run under `product/knowledge/imports/`. Review and explicitly approve mappings before materializing draft product artifacts.

Operational navigation uses concurrent workspaces rather than a global active feature:

```text
spec-framework work --feature <path-or-id> --created-by <human>
spec-framework status --work WORK-001
spec-framework next --work WORK-001
spec-framework approve --artifact <path> --grant approved --approved-by <human> --yes
spec-framework gates
spec-framework guide --work WORK-001
spec-framework graph materialize --graph <execution-graph.json> --yes
spec-framework task readiness --graph <execution-graph.json> --task TK-001 [--json]
spec-framework review --work WORK-001 --stage <stage>
spec-framework approve-stage --work WORK-001 --stage <stage> --approved-by <human> --yes
spec-framework impact --decision DEC-021 [--json]
spec-framework dashboard --work WORK-001 [--json]
spec-framework status --work WORK-001 --graph [--json]
spec-framework decisions migrate [--dry-run | --interactive | --yes]
```

Recommended prompt for the architecture phase:

```text
You are a Software Architect collaborating on the Product Engineering Framework.
Read .spec-framework/FRAMEWORK.md and the relevant product/**/context.md files.

At this stage, do not create files and do not implement.
Your mission is to critique the architecture, find ambiguities, propose alternatives,
compare trade-offs, and ask what needs to be approved.

Only implement when I say: FREEZE ARCHITECTURE.
```

Recommended prompt for the generation phase:

```text
Read .spec-framework/FRAMEWORK.md.
Use only approved decisions.
Do not invent new layers, names, or flows.
Convert the approved architecture into files, templates, and skills.
Preserve traceability between Problem, Vision, Strategy, Domain, Goal, Feature,
Use Case, Specification, Implementation Plan, Execution Graph, and Tasks.
```

Recommended prompt for a new feature:

```text
Read .spec-framework/FRAMEWORK.md and the domain's context.md.
Drive the feature through the flow:
Feature -> Use Cases -> Specification -> Design -> Implementation Plan -> Execution Graph -> Tasks.

Before persisting each stage, list the decisions, gaps, conflicts, and approval questions.
```

## 16. Framework Roadmap

### v0

- Folder structure.
- Basic templates.
- `FRAMEWORK.md`.
- List of skills and orchestrators.

### v1

- Canonical contexts at every level.
- Consistent IDs.
- Decision log.
- Complete templates.
- Basic audits.

### v2

- Operational skills.
- Orchestrators with handoff.
- Real Execution Graph.
- Task generation by DAG.
- Approval gates.

### v3

- Queryable knowledge graph.
- Automatic impact analysis.
- Supervised automatic agent spawning with provider adapters, `max_parallel`, leases, isolated worktrees, heartbeat, cancellation, and checkpoint recovery.
- Task-level review/QA followed by governed integration and Integrated QA for automatically spawned work.
- Automatic replanning after failures.

## 17. Final Rule

The framework must help agents think before they build.

If an AI cannot explain which problem, domain, goal, feature, use case, and Specification a task was born from, the task is not yet ready for implementation.
