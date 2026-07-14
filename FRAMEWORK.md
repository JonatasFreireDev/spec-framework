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

Defines the product we want to build, for whom, why now, and the durable outcome it should create. `vision.md` owns that direction and its boundaries; it links to, but does not duplicate, the canonical principles in `principles.md` or the canonical outcome metric and guardrails in `north-star.md`.

### Product Principles

Define the durable decision rules, trade-offs, examples, and anti-principles that guide product choices. They are canonical in `foundation/vision/principles.md`.

### North Star

Defines the durable user-value outcome, candidate metric, measurement notes, and guardrails. It is canonical in `foundation/vision/north-star.md`.

### Feature Brief

Provides the proportional entry contract for the `existing-feature` starting point. `foundation/feature-brief.md` combines the evidenced problem, desired outcome, scope, non-goals, constraints, success signal, and delivery strategy for one bounded feature. It replaces full Foundation approval only for that starting point; broad or uncertain product direction still requires Problem, Vision, Product Principles, North Star, and Strategy.

### Implementation Assessment

Provides the evidence boundary for the `existing-implementation` starting point. `knowledge/assessments/implementation-assessment.md` records observed behavior, architecture, data, integrations, tests, operational constraints, risks, and candidate product claims. Its individual approval confirms the assessment of current evidence; it does not approve inferred product intent. The full Foundation path remains required before workspace creation.

### Product Baseline

Provides the current-state contract for the `existing-product` starting point when code and operational evidence are more complete than product documentation. `foundation/product-baseline.md` consolidates the active audience, evidenced needs, delivered value, capabilities, operating model, current signals, decision rules, constraints, risks, and unknowns. It replaces separate Problem, Vision, Product Principles, and North Star approvals only for this starting point. Strategy remains separate and owns future direction.

### Strategy

Defines positioning, segments, metrics, trade-offs, roadmap, and criteria to advance or pause.

### Domain

Groups a coherent area of the business or product, such as `users`, `groups`, `events`, `friendship`, or `payments`. A domain is not the product name, a UI navigation section, or a catch-all container. It declares both what it owns and what it does not own, including cross-domain dependencies. The first modeled domain must continue through at least one User Goal, Feature, and Use Case as a walking skeleton; `domain.md` alone is not a delivery slice.

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

Translates the Specification into a verifiable user experience: visual flow, navigation, wireframes, mockups, states, accessibility, and alignment with the design system. When the feature has no interface, the artifact must use structured status `not_applicable` and a non-placeholder rationale.

Design declares two independent dimensions. `origin_mode` is `generate` when the solution is created from the Specification, `evolve` when an existing interface is intentionally changed, and `adopt` when a versioned external source is authoritative. `maturity` is `contract`, `wireframe`, `mockup`, or `prototype`. Existing descriptive designs are compatible as `generate/contract`.

Approved decisions and the Specification remain authoritative for behavior, security, privacy, and business rules. A source marked `visual_canonical` is authoritative only for presentation and interaction details that do not conflict with those higher-precedence contracts. Design tools and services are adapters; they never replace `design.md` or approve it.

### Design System

The Design System is an optional shared product artifact for products with recurring interface foundations, tokens, components, and patterns. It lives under `design/system/`, declares `generate`, `evolve`, or `adopt`, and uses semantic versioning. `design-system.md` is the human contract and `tokens/tokens.json` is the mechanical token source. Use-case Designs pin the system id/version and record consumed tokens, components, patterns, and deviations. A Design System does not replace the Specification or approve use-case Design.

### Engineering System

The optional Engineering System versions stable architecture, module and data ownership, integrations, standards, quality attributes, operations, and evidence under `engineering/`. Its detailed versioning, maturity, hashing, migration, and approval contract lives in [`docs/engineering-systems.md`](docs/engineering-systems.md). Specification and approved `DEC-*` records remain authoritative when contracts conflict.

### Engineering Quality System

The Engineering Quality System is the shared contract under `engineering/quality/` for quality attributes, test levels, risk-based coverage, environments, data, fitness functions, evidence, flaky tests, exceptions, and maturity. Its detailed policy, evidence, migration, compatibility, and approval contract lives in [`docs/engineering-systems.md`](docs/engineering-systems.md). It defines policy and capability, not delivery approval.

### Engineering Proposal

Translates an approved Specification, Design, Technical Discovery, and resolved Architecture Gate into the intended technical solution for one delivery. It describes boundaries, data ownership, integrations, dependencies, operations, tests, rollout, and deviations from the pinned Engineering System without sequencing implementation tasks.

### Engineering Review

Independently and read-only evaluates an Engineering Proposal before Implementation Planning. It verifies architecture, decisions, ownership, dependencies, quality attributes, operations, and testability. It does not edit the proposal, approve product decisions, or review implementation code.

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
- `L`: critical or sensitive flow; also requires Engineering Proposal, Engineering Review, analytics, audit, QA evidence, and Security Review.
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

An adopter repository owns one framework root, `product/`. Its stable ownership boundaries are:

| Area | Ownership |
| --- | --- |
| `.product/` | Manifest, artifact registry, decisions index, roadmap state, workspaces, claims, derivations, and approval history. |
| `foundation/` | Starting-point contracts and product direction. |
| `knowledge/` | Product rules, conventions, decisions, imports, and durable evidence. |
| `domains/` | Domain → User Goal → Feature → Use Case hierarchy and delivery artifacts. |
| `design/` | Shared Design System and product-owned design sources. |
| `engineering/` | Versioned Engineering System, quality contracts, operations, and evidence. |
| `audits/` | Readiness, consistency, QA, security, and release findings. |
| `releases/` | Product release records when materialized. |

A use-case bundle owns `context.md`, `use-case.md`, Specification and modular contracts, applicable Design, Technical Discovery, Engineering Proposal and Review, Implementation Plan, Execution Graph, canonical task files plus generated index, tests, analytics, and validation evidence. Exact paths and optional directories are defined by the starter, templates, registries, and declarative initialization contracts; this method defines ownership rather than duplicating that inventory.

The framework runtime remains outside the adopter repository:

```text
repo/
  product/
    BOOTSTRAP.md
    .product/
      framework.json
    foundation/
    knowledge/
    domains/
    design/
    engineering/
    audits/
    releases/
```

This repository keeps three explicit sources:

- `starter/` contains canonical product assets selected by initialization contracts.
- `examples/events/` contains the worked product instance used as learning material and validation fixture.
- `framework/` contains the executable framework core: audits, decisions, skills, templates, validators, distributable tools, framework-only tests, and adoption guidance. The repository root retains only entry points, packaging metadata, scripts, examples, and starter infrastructure.

New products never copy the repository root. `product/.product/framework.json` pins the method version and is the exclusive activation marker; the CLI materializes embedded method assets in the operating system's versioned user cache.

The modular artifact composition and approval-adapter map is maintained in [`docs/artifact-registry-modules.md`](docs/artifact-registry-modules.md). Starting-point contracts select artifact modules; the approval engine remains generic and only invokes adapters for composite or side-effectful contracts.

## 5. Context.md

Every `context.md` must let an AI understand where it is, what it needs to read, and what the safe next step is.

Product READMEs are navigation aids, not workflow or state authorities. Keep them only where they clarify an area boundary or a copyable template entry point, and keep them concise. `BOOTSTRAP.md` owns starting-point sequencing, `context.md` owns local state and handoff, skills own operating workflow, and templates own artifact structure. Do not ship a README only to preserve an empty directory; declarative initialization or the owning skill creates that directory when a real artifact requires it. Upgrade preserves adopter-owned READMEs even when newer starters omit them.

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
Feature -> Use Cases -> Specification -> Design -> Technical Discovery -> Engineering Proposal -> Engineering Review -> Implementation Plan -> Execution Graph -> Tasks -> Implementation -> QA Evidence -> Security Review -> Review -> Audit -> Release
```

The closed delivery flow is:

```text
Domain -> Domain Evolution -> Feature Selection -> New Feature -> Use Cases
-> Specification Contracts -> Design -> Technical Discovery -> Architecture Gate
-> Engineering Proposal -> Engineering Review -> Implementation Plan -> Execution Graph -> Tasks -> Code Runner
-> Code Review -> QA -> Commit Crafter -> PR Finalizer
```

`specification.md` remains the root contract. Large concerns live under `contracts/` (`product`, `behavior`, `ux`, `api`, `data`, `security`, `quality`, `observability`, and `rollout`) and use stable `REQ-*` and `AC-*` identifiers. Tier S requires behavior and quality; Tier M adds product, UX, API, data, rollout, and Technical Discovery; Tier L adds security and observability. An inapplicable contract must use structured status `not_applicable` and a non-placeholder rationale.

Design is mandatory for any use case with an interface. For deliveries without UI, `design.md` must exist with structured status `not_applicable`, a non-placeholder `Rationale`, and impacts on accessibility, observability, or operations when relevant. Free-text mentions of Not applicable never satisfy the gate. The same structured status and rationale contract applies to inapplicable modular Specification contracts.

QA Evidence and Security Review are validation gates. QA Evidence proves that acceptance criteria, tasks, flows, edge cases, regression, accessibility, observability, and security controls were verified. Security Review evaluates authentication, authorization, privacy, abuse, sensitive data, tokens, logs, dependencies, rollout, rollback, and residual risk. Security Review must also read the product's security baseline in `knowledge/conventions/security-baseline.md` and active threats in `audits/security/threat-register.md`. An artifact must not reach `validated` or `released` when there is a QA or security blocker.

When an Engineering Quality System is configured, `tests.md` must pin the consumed Engineering System id/version, reference the canonical quality policy, and select only environments, test-data classes, and platforms configured by its mechanical catalog. Its `Acceptance Traceability` table must map every acceptance criterion to meaningful risk, validation method, test level, evidence, and owner values, and it must declare either `None` or open, unexpired, in-scope `QEX-*` deviations explicitly. Approved QA evidence must pin the same version and record `passed` for Quality System conformity, environment and test data, and flaky-test and exception checks; `N/A` does not satisfy these configured gates. Legacy products without the new contract remain compatible until they explicitly migrate their Engineering System.

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

Documentary rigor is proportional to the use case's tier. Tier S avoids heavy artifacts when design, analytics, or audit use structured `not_applicable`; Tier L requires Engineering Proposal, Engineering Review, QA Evidence, and Security Review by default. Tier S and M require Engineering Proposal and Engineering Review when `context.md` declares at least one structured `engineering_trigger`: `new_dependency`, `external_integration`, `data_ownership_change`, `migration`, `architecture_boundary_change`, `deployment_change`, `security_boundary_change`, or `operational_change`. Trigger automation reads this closed list and never infers applicability from prose.

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
- A versioned source manifest, screen inventory, requirement mapping, and fidelity evidence when Design adopts or evolves an external source.

Gates:

- Without an approved Specification, do not generate design.
- Without an approved Design, or one using structured `not_applicable` with rationale, do not continue to Technical Discovery.
- Blocking UX findings go back to Specification or Design before proceeding.
- Canonical visual sources must have an immutable version or content hash. Missing required states, unresolved Specification conflicts, or unreviewed strict-fidelity deviations block Design from advancing.
- When the product declares a Design System, proposed-or-later use-case Design must pin its approved id/version and may not introduce shared tokens, components, or patterns silently.

## 7. Implementation Plan

The Implementation Plan is created after the Specification and the Design, and before the tasks. It must not write code. It must define the build strategy.

Before planning, `technical-discovery.md` must map applicable requirements to the real codebase and stable knowledge in `engineering/`. Its Architecture Gate must reference an approved decision or state `Not required` with concrete rationale.

When Engineering Proposal applies, `engineering-proposal.md` must pin the Engineering System version or explicitly state that no shared system is configured, distinguish existing evidence from the intended change, and receive a non-blocking `engineering-review.md` verdict before the Implementation Plan advances. A passed review records the SHA-256 hash of the current proposal; any proposal change makes the review stale. Proposed-or-later consumption of a configured Engineering System requires its current approved version and hash-matching approval evidence. Tier L always requires both artifacts. Engineering Review never substitutes for approval of a required product decision.

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

The Execution Graph is a DAG. It defines dependencies between tasks and enables parallel execution by agents. Each node references the task's canonical file by `path`. Runtime mechanics for workspaces, leases, scheduling, command plans, and integration are defined in [`docs/execution-runtime.md`](docs/execution-runtime.md).

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
- The `writeScope` check rolls out in two compatible phases: Phase A reports warnings; Phase B, after approved migration of existing graphs, can promote the same findings to errors.
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

Runtime v2 leases live in `.product/claims/<task-id>.json`; `.product/claims.json` remains a v1 compatibility index during migration. See [`docs/execution-runtime.md`](docs/execution-runtime.md) for the complete runtime contract.
## 9. Skills

Runtime v2 makes execution resumable and safely parallel. The complete workspace, lease, scheduler, command-plan, worktree, recovery, and integration contract is defined in [`docs/execution-runtime.md`](docs/execution-runtime.md). The runtime does not spawn agents or grant approval.

Runtime commands include `runtime`, `resume`, `handoff`, `checkpoint`, `lease`, `commands`, `schedule`, and `integrate`.

Skills are specialists. They can operate in modes such as `create`, `update`, `audit`, `evolve`, `explain`, `compare`, and `refactor`, but each must have a clear responsibility.

Definition and planning skills follow the shared Discovery and Challenge contract before substantive creation or material revision. They inspect repository and CLI evidence first, then use the harness-native structured question capability for human choices that cannot be discovered. Each round asks one to three focused questions; meaningful choices present concrete options, trade-offs, a recommendation, and a free-form path. Skills proactively warn about material scope, dependency, usability, security, operability, reversibility, approval, and delivery risks and propose safer alternatives. They must not finalize or hand off while a blocking question is unanswered, and conversational answers never grant formal approval. Harness adapters map the canonical `native_user_question` capability to their default question tool; a concise conversational question is the explicit fallback only when no structured tool is exposed.

| Layer | Specialist ownership |
| --- | --- |
| Foundation | Problem Discovery, Vision, Strategy, Domain Architect, and User Goal define evidenced product direction and boundaries. |
| Product Design | Journey, Feature, Use Case, UX/UI, UX Review, and Design System define value slices and verifiable experience. |
| Specification and planning | Specification, Engineering System, Technical Discovery, Engineering Proposal, Engineering Review, Implementation Planner, Execution Graph, and Task Generator turn approved intent into executable contracts. |
| Engineering and validation | Code Runner implements one approved task; QA, Code Review, and Security Review independently verify it; Threat Modeler maintains proactive security context; Commit Crafter and PR Finalizer prepare verified delivery without merging. |
| Audit and evolution | Gap, Conflict, Dependency, Impact, Evolution, Documentation, and Product Historian skills inspect health and route governed change. |

The versioned `framework/skills/` contracts are authoritative for inputs, procedure, outputs, limitations, and handoff. A specialist never assumes another skill's ownership or grants human approval.

## 10. Orchestrators

Orchestrators do not create primary artifacts. They control flow, order, gates, and handoffs.

| Orchestrator | Boundary |
| --- | --- |
| Framework Guide | Translates intent and current CLI state into the smallest safe route; it does not author artifacts or approve work. |
| Product | Coordinates the approved Foundation sequence and stops at each human gate. |
| Domain Evolution | Compares evidenced feature candidates and requires explicit human selection before handoff. |
| Existing Product Import | Moves sources through inventory, per-source traceability, conflicts, reviewed mappings, explicit materialization, and draft artifacts without treating sources as truth. |
| New Feature | Drives an approved feature through Use Cases, Specification, Design, engineering gates, Plan, Graph, and Tasks. |
| Audit and Evolution | Batch findings, compare improvement candidates, and route selected changes without silently expanding scope. |
| Documentation | Synchronizes contexts, indexes, templates, decisions, and derived documentation. |
| Delivery and Release | Coordinates implementation, independent review/QA/security, verified commits, PR preparation, integration, and release readiness without merging or bypassing gates. |

Framework Guide is the default entry route. Direct specialist routing requires current-session CLI evidence naming scope, gate, and owner, or an explicit human request naming specialist and scope. Persisted state must be revalidated; stale or ambiguous routes return to the Guide.

Failure routing is fixed. Defects with clear expected behavior route to Bug Fixer; missing or hollow coverage routes to QA or the test owner; incomplete or out-of-contract implementation routes to Code Runner; decision gaps route to Product Historian and a human. Code Review and QA remain independent and read-only, every code change re-enters QA, and three failed automated attempts require human escalation. Threat Modeler maintains proactive shared context; Security Review evaluates the concrete delivery. Commit Crafter does not push, PR Finalizer does not merge, and release readiness requires all applicable product, engineering, evidence, review, QA, and security gates.

## 11. Approval Gates

The detailed lifecycle, approval-record, staleness, authority, and failure-routing contract is maintained in [`docs/lifecycle-and-approvals.md`](docs/lifecycle-and-approvals.md). This section remains the canonical summary of the method-level gates and transitions.

Every step must end with a clear state:

```text
draft
proposed
materialized (Execution Graph only)
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

- Foundation ladder artifacts are registered from initialization: Problem, Vision, Product Principles, North Star, and Strategy. Their parent relationships are enforced by `spec-framework approve`, and approval updates the canonical artifact, applicable context status, registry status, and `.product/history/` evidence atomically.
- `validate --write-registry` preserves `parents`, `children`, `depends_on`, decisions, and delivery dependencies from structured companion `context.md` YAML. Starting-point-specific parents are additive and must not erase the existing product graph.
- `proposed`: does not require an approval record, but must not advance from an incomplete parent gate.
- `approved` and later states: require a corresponding approval record in `.product/history/`, with `artifact_id`, `path`, `content_hash`, `status_granted`, `approved_by`, `approved_at`, and `notes`. Validator and operational navigation both require that record to match the current artifact content; editing status prose alone never advances work.
- Human approval may be applied to one artifact with `approve` or to an explicit batch with `approve-batch`. Batch approval must preview the exact paths, IDs, hashes, ignored artifacts, blockers, and next gate first; it requires explicit scope, human identity, and `--yes`, and never includes stale or ineligible artifacts.
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

Framework or method changes are incorporated directly into this document, owning skill contracts, validators, and tests. Git history preserves their evolution. Skill contracts, gates, writeScope, QA policies, failure routing, commit policy, validators, and orchestration rules must not be recorded in `knowledge/decisions/`, because that folder is reserved for the adopter product.

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

## 15. How Agents Use The Framework

Agents in this repository maintain the framework; agents in adopter repositories operate only on the pinned instance and product-owned content under `product/`. Harness interaction details may differ, but method, gates, artifact contracts, and CLI semantics are agent-independent.

The runtime behavior shared by all agents is defined in the pinned `AGENTS.framework.md`. This section summarizes routing and authority only; specialized procedure belongs to the owning skill, and runtime mechanics belong to [`docs/execution-runtime.md`](docs/execution-runtime.md) and [`docs/lifecycle-and-approvals.md`](docs/lifecycle-and-approvals.md).

### Activation And Routing

- Activate only from a valid `product/.product/framework.json` with `framework: spec-framework`, a concrete version, and `activation.mode: manifest-only`.
- Resolve versioned skills through the user-scoped dispatcher. Do not copy specialist trees into adopter repositories.
- Route through Framework Guide unless current CLI evidence or an explicit human request names both specialist and concrete scope. Revalidate persisted handoffs with `guide`, `dashboard`, `status`, or `next`.
- Read `product/BOOTSTRAP.md`, relevant `context.md` files, applicable decisions, the owning skill, and its template before mutation.

### Initialization And Preservation

- `init` resolves one strict declarative starting-point contract, validates its complete materialization plan, stages it, and atomically publishes `product/`. Data contracts cannot execute arbitrary commands or escape the product root.
- A starting point changes the initial evidence, registry, bootstrap, and first gate; it never removes later rigor or approval requirements. Existing code and documents remain evidence, not approved truth. `audit-only` remains read-only until an explicit supported transition.
- `existing-documents` creates an analysis-only import run with `traceability.json` as the dedicated per-source ledger. The Artifact Importer agent reads each source, records evidence, extracted claims, candidate destinations, and unmapped gaps there, then proposes mappings. Human review of inventory, traceability, conflicts, and selected mappings is required before explicit draft materialization; imported artifacts retain their normal owners, parents, and individual approval gates.
- `init` never overwrites an existing `product/`. `upgrade` refreshes only the pinned runtime, manifest, and selected dispatchers; it never replays initialization over adopter-owned content.
- Starting-point details belong to `docs/starting-points.md`, `framework/init/`, and generated `BOOTSTRAP.md` rather than this operational summary.

### Operation And Authority

- Use concurrent `WORK-NNN` workspaces rather than a global active feature. Read CLI help and current mechanical output for command syntax and next actions.
- Definition and planning skills inspect evidence, use the harness-native question capability, compare options, recommend a path, warn about material risks, and stop on blocking questions. Execution and review skills follow approved contracts and route product ambiguity back to the owning definition skill.
- Human identity and explicit confirmation are required for approvals and consequential mutations. Conversation never creates approval evidence.
- External adapters are optional and supervised. Their output and availability never grant approval; provider failure does not roll back an initialized product.

### Lifecycle Boundary

Installation, CLI update, product upgrade, and removal are separate. Install scripts verify and install the CLI without running `init`; `update` changes the CLI binary; `upgrade` changes the pinned product runtime; `uninstall` removes only managed CLI paths and optional namespaced caches/dispatchers. Product repositories are never searched for or removed by CLI uninstall.

The CLI's generated help is authoritative for command syntax. Git history records how framework boundaries evolved; skills define specialist procedure; templates define artifact structure; `BOOTSTRAP.md` defines the initial route; `context.md` defines local state and handoff.

## 16. Final Rule

The framework must help agents think before they build.

If an AI cannot explain which problem, domain, goal, feature, use case, and Specification a task was born from, the task is not yet ready for implementation.
