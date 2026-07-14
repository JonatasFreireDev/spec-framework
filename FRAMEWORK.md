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

The Engineering System is an optional shared product artifact for stable architecture, module and data ownership, integrations, standards, quality attributes, fitness functions, operations, and evidence. It lives under `engineering/`, declares `generate`, `evolve`, or `adopt`, and uses semantic versioning. `engineering-system.md` is the human contract and `engineering-system.yaml` is the mechanical catalog. Maturity is declared by area as `baseline`, `mapped`, `governed`, `verified`, or `operated`; it describes available evidence and never grants approval.

Approved decisions and the Specification remain authoritative. The Engineering System records and links established technical constraints but does not replace `DEC-*` records or their approval evidence. Its approval is composite: context, human contract, mechanical catalog, architecture, standards, quality, runbooks, and evidence are hashed in deterministic path order, so any shared engineering contract change requires human re-approval.

### Engineering Quality System

The Engineering Quality System is the shared, versioned quality contract within the Engineering System. It lives under `engineering/quality/` and defines product-wide quality attributes, test levels, risk-based coverage, environments, test data, fitness functions, evidence policy, flaky-test handling, exceptions, and maturity. Its canonical human contract is `quality-system.md`; `quality-system.yaml` is the mechanical catalog; `quality-model.md`, `test-strategy.md`, and `fitness-functions.yaml` are supporting contracts.

The Engineering Quality System defines policy and available capability; it does not grant delivery approval. A use-case `tests.md` applies the shared policy and pins the consumed Engineering System id/version. Delivery tasks implement tests. QA independently verifies the applicable policy, acceptance criteria, gates, and real evidence in `qa-evidence.md`. Security Review remains a separate specialized gate.

Gate commands remain canonical in the active product root's `knowledge/conventions/gates.md`. Quality maturity records available evidence and never waives a gate, changes an acceptance criterion, or approves residual risk. Human and mechanical capability maturity must agree; maturity above `baseline` requires a safe, resolvable path, URL, gate, CI, or command evidence reference. Exceptions require scope, owner, rationale, residual risk, mitigation, a valid expiry or review date, re-entry gate, and status. Only open, unexpired exceptions scoped to `product` or the consuming `domains/...` path may authorize deviations. Engineering System approval synchronizes status atomically across its context, human and mechanical system contracts, and both Quality System contracts before hashing the complete `engineering/` tree.

Legacy Engineering Systems whose quality area points to `quality/quality-model.md` remain valid until explicit migration. `spec-framework engineering-system migrate --dry-run` previews the additive Quality System migration; the applied command creates missing contracts before atomically replacing the catalog pointer, rolls back newly created files on failure, preserves adopter-owned files, and never creates approval evidence. Any previously approved system requires human re-approval because its composite hash changes.

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
                  technical-discovery.md
                  engineering-proposal.md
                  engineering-review.md
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
    system/
      context.md
      design-system.md
      foundations/
      tokens/
        tokens.json
        themes.json
      components/
      patterns/
      sources/
      evidence/
  engineering/
    context.md
    engineering-system.md
    engineering-system.yaml
    architecture/
    standards/
    quality/
      quality-system.md
      quality-system.yaml
      quality-model.md
      test-strategy.md
      fitness-functions.yaml
    runbooks/
    evidence/
  audits/
  releases/
  skills/
```

In an adopter product repository, the structure above must live under `product/`. The framework runtime and method assets are resolved outside the repository from the versioned user cache:

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

In this `spec-framework` repository, the structure exists in three explicit areas:

- `starter/` represents the clean `product/` skeleton copied into new product repositories.
- `examples/events/` contains the worked product instance used as learning material and validation fixture.
- `framework/` contains the executable framework core: audits, decisions, skills, templates, validators, distributable tools, framework-only tests, and adoption guidance. The repository root retains only entry points, packaging metadata, scripts, examples, and starter infrastructure.

New products must not copy the entire `spec-framework` root. They start from `starter/product/`; `product/.product/framework.json` pins the adopted version and is the exclusive activation marker. The CLI materializes embedded method assets in the operating system's versioned user cache.

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

Definition and planning skills follow the shared Discovery and Challenge contract before substantive creation or material revision. They inspect repository and CLI evidence first, then use the harness-native structured question capability for human choices that cannot be discovered. Each round asks one to three focused questions; meaningful choices present concrete options, trade-offs, a recommendation, and a free-form path. Skills proactively warn about material scope, dependency, usability, security, operability, reversibility, approval, and delivery risks and propose safer alternatives. They must not finalize or hand off while a blocking question is unanswered, and conversational answers never grant formal approval. Harness adapters map the canonical `native_user_question` capability to their default question tool; a concise conversational question is the explicit fallback only when no structured tool is exposed.

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
- Design System AI: creates, adopts, evolves, versions, and audits shared foundations, tokens, components, patterns, and sources before they are consumed by use-case Design.

### Specification And Planning

- Specification AI: creates the central implementation contract.
- Engineering System AI: creates, adopts, evolves, versions, and audits shared architecture, standards, quality attributes, fitness functions, operations, and evidence.
- Engineering Proposal AI: describes one delivery's intended technical change after discovery and before independent review.
- Engineering Review AI: independently reviews a delivery's Engineering Proposal before implementation planning without editing it or approving decisions.
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

### Framework Guide

Framework Guide is the conversational entry point to the CLI. It translates a person's goal into current mechanical state, the smallest safe command, and the correct specialist or approval handoff. It does not author canonical artifacts, approve work, or replace Command Planner and Command Executor.

The installed dispatcher routes framework-governed product operations through Framework Guide by default. It may route directly to a specialist only when current-session `guide`, `dashboard`, `status`, or `next` output names the workspace, concrete feature or use-case scope, current gate, and owner skill, or the human explicitly names both the specialist and the concrete scope. A persisted handoff or checkpoint identifies where to resume but must be revalidated by one of those read-only CLI commands before direct routing. A skill name or keyword without scope is not a verified route. A stale, ambiguous, or conflicting route returns to Framework Guide. This is an agent-routing rule, not a new product approval gate, and direct diagnostic CLI commands remain available.

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
Impact -> Feature -> Use Cases -> Specification -> Design -> Technical Discovery -> Architecture Gate -> Engineering Proposal -> Engineering Review -> Plan -> Graph -> Tasks
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

QA is an independent, read-only verifier. QA re-executes the gates declared in `knowledge/conventions/gates.md` whenever possible, records real output or an explicit limitation, and routes blockers to the appropriate Code Runner, Bug Fixer, test owner, Product Historian, or human path instead of fixing code.

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

- Foundation ladder artifacts are registered from initialization: Problem, Vision, Product Principles, North Star, and Strategy. Their parent relationships are enforced by `spec-framework approve`, and approval updates the canonical artifact, applicable context status, registry status, and `.product/history/` evidence atomically.
- `validate --write-registry` preserves `parents`, `children`, `depends_on`, decisions, and delivery dependencies from structured companion `context.md` YAML. Starting-point-specific parents are additive and must not erase the existing product graph.
- `proposed`: does not require an approval record, but must not advance from an incomplete parent gate.
- `approved` and later states: require a corresponding approval record in `.product/history/`, with `artifact_id`, `path`, `content_hash`, `status_granted`, `approved_by`, `approved_at`, and `notes`. Validator and operational navigation both require that record to match the current artifact content; editing status prose alone never advances work.
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

## 15. How Agents Use The Framework

In this repository, agents maintain the framework itself. In adopter repositories, supported agent harnesses operate on an instance composed from `starter/product/`, resolve the pinned method through the CLI runtime, and write framework documentation artifacts in `product/`. Harness-specific installation or interaction details may differ, but the manifest, method, gates, artifact contracts, and CLI behavior are agent-independent.

Spec Framework activates only when the current repository contains a valid `product/.product/framework.json` with `framework: spec-framework`, a concrete version, and `activation.mode: manifest-only`. A user mention, keyword, prompt, or similarly named file does not activate it. A single user-scoped, namespaced dispatcher resolves skills from the versioned cache; specialized skill trees are not copied into adopter repositories. `init` and `upgrade` install or refresh that dispatcher for every selected agent. After activation, the dispatcher resolves `framework-guide` first unless it has a verified direct route from current CLI guidance or an explicit human request that names both specialist and concrete scope; persisted runtime state must first be revalidated through the CLI.

In a newly initialized product repository, read `product/BOOTSTRAP.md` first. It explains the ordered foundation gates and distinguishes a structurally valid starter from a product that is ready for implementation.

The CLI command tree is implemented with Cobra. It owns command discovery, help, subcommand routing, and compatible exit handling; product manifests and explicit CLI flags remain the authoritative configuration contracts. Viper is not introduced unless the CLI gains a documented need for user-level configuration precedence across defaults, environment, and config files.

During `init`, choose the repository's starting point. This choice customizes `BOOTSTRAP.md` and the active artifact registry; it does not remove skills, orchestrators, rigor requirements, or approval gates. The `existing-feature` starting point replaces the full product Foundation package with one individually approved `foundation/feature-brief.md` for the bounded delivery. If product direction or scope is broad or uncertain, escalate to the full Problem, Vision, Product Principles, North Star, and Strategy path. Every starting point that creates or revises domains must read the pinned runtime's `examples/events/` reference before the first domain change and model a business-area boundary, explicit non-ownership, cross-domain dependencies, and one Domain -> User Goal -> Feature -> Use Case walking skeleton. `audit-only` uses the same reference to assess existing domain boundaries without changing product artifacts. When starting from existing documents, the CLI creates a source inventory and an analysis-only import run under `product/knowledge/imports/`. Review and explicitly approve mappings before materializing draft product artifacts.

Initialization is driven by versioned declarative contracts under `framework/init/contracts/`. Each starting point selects named asset sets, explicit directories including intentionally empty architecture, entry-specific template files, deterministic text patches, initial registry transformations, its bootstrap profile, and a closed list of CLI-owned actions. The CLI strictly parses the complete contract, expands every embedded source into an in-memory plan, rejects unknown fields, unsafe targets, file/directory collisions, ambiguous patches, and invalid artifact relationships, then materializes the verified plan in staging and atomically publishes `product/`. Contracts cannot execute arbitrary commands or escape `product/`. Entry-specific actions that require runtime state, such as creating an immutable import run, remain typed CLI behavior selected by contract name and complete in staging before publication. The contracts preserve the current starting-point method: reference skeletons required for explicit escalation remain present but inactive where the applicable registry contract excludes them. `init` never overwrites an existing `product/`, including with `--force`; `upgrade` refreshes the pinned runtime and manifest without replaying initialization contracts over adopter-owned content.

CLI lifecycle and adopter lifecycle are separate. The checksum-verifying `install` scripts install the released binary, record installer ownership in `install.json`, and do not run product initialization. `spec-framework update [--check | --version <version>]` checks or replaces that binary; replacement requires `--yes`, verifies the official release checksum and candidate version, stages beside the executable, and preserves or recovers the current binary on failure. `spec-framework upgrade` continues to update only an adopter's pinned runtime, manifest, and dispatchers. `spec-framework uninstall` previews ownership and removes the local binary, managed manifest, and installer-managed PATH entry; `--purge` additionally removes versioned runtime caches and only namespaced Spec Framework dispatchers. Neither uninstall mode searches for or removes `product/` repositories.

The `existing-implementation` starting point first materializes and individually approves `knowledge/assessments/implementation-assessment.md`. The assessment becomes a parent of Problem, after which Problem, Vision, Product Principles, North Star, and Strategy require their normal individual approvals before workspace creation. Implementation evidence may support Foundation drafts, but it never grants product approval by itself.

The `existing-product` starting point may treat code, tests, releases, telemetry, and operational history as primary evidence for one individually approved `foundation/product-baseline.md`. It then requires a separate individually approved Strategy whose parent is that baseline. Unknown audience, unclear delivered value, or material repositioning escalates to the full Foundation path.

For `existing-documents`, the latest import run declared in the canonical manifest is the entry gate. Review its inventory, conflicts, and selected mappings, then run `spec-framework import materialize` with explicit human identity. Materialization authorizes only the selected draft writes; each resulting product artifact retains its normal owner, parent, validation, and individual approval gates. Workspace creation remains blocked until that latest run is materially complete and approved for materialization.

When the canonical manifest declares `audit-only`, the CLI permits read-only validation, inspection, status, readiness, review, impact, and navigation, but blocks writes to product artifacts, registry, reports, approvals, workspaces, imports, migrations, design state, and delivery runtime state. `init` may establish the audit manifest and `upgrade` may maintain its pinned framework runtime. Turning findings into product work requires an explicit supported starting-point transition; agents must not bypass the guard by editing statuses or metadata manually.

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
spec-framework adapters list
spec-framework adapters status impeccable
spec-framework adapters doctor impeccable [--check-latest]
spec-framework adapters install impeccable --version <provider-cli-version> [--yes]
```

External adapters are optional. Read-only discovery and diagnosis may run without confirmation. Install and update must show the exact provider command, require an explicit version and `--yes`, execute with direct argv from the repository root, and never fabricate readiness after a provider failure. Removal is unsupported until the provider documents a reversible contract.

Interactive `init` may offer an explicit optional-adapter install choice. Headless installation requires adapter-specific selection, an explicit provider version, and `--yes`. Framework initialization completes before the external provider runs; provider failure is reported as partial success and never rolls back or deletes the initialized product.

Recommended prompt for the architecture phase:

```text
You are a Software Architect collaborating on the Product Engineering Framework.
Resolve the pinned framework through spec-framework skill path and read the framework root's FRAMEWORK.md plus the relevant product/**/context.md files.

At this stage, do not create files and do not implement.
Your mission is to critique the architecture, find ambiguities, propose alternatives,
compare trade-offs, and ask what needs to be approved.

Only implement when I say: FREEZE ARCHITECTURE.
```

Recommended prompt for the generation phase:

```text
Resolve and read the pinned framework root's FRAMEWORK.md.
Use only approved decisions.
Do not invent new layers, names, or flows.
Convert the approved architecture into files, templates, and skills.
Preserve traceability between Problem, Vision, Strategy, Domain, Goal, Feature,
Use Case, Specification, Implementation Plan, Execution Graph, and Tasks.
```

Recommended prompt for a new feature:

```text
Resolve and read the pinned framework root's FRAMEWORK.md and the domain's context.md.
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
