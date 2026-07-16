# Compozy-Inspired Improvement Backlog

## Purpose and status

This is a framework-maintenance backlog derived from a read-only comparison with [compozy/compozy](https://github.com/compozy/compozy). It records candidates; it does not approve architecture, alter Framework contracts, or authorize implementation.

Every topic preserves these constraints:

- Product-owned artifacts remain the source of truth.
- Human approval and current hash-matching approval records remain mandatory.
- Code Review, QA, and Security Review remain independent and read-only.
- An operational runtime may produce evidence but cannot validate, release, commit, push, merge, or create approval records.
- Upgrade must preserve adopter-owned product content.

## Working agreement

Before implementation, create a bounded proposal naming affected contracts, CLI surface, storage, migration, default and removal behavior, plus combination tests. Obtain an approved framework decision first whenever the topic changes architecture, an external dependency, security/privacy, or a hard-to-reverse policy.

| ID | Topic | Disposition | Priority | Status |
| --- | --- | --- | --- | --- |
| CI-01 | Operational event ledger and replay | Adopt | P1 | Candidate |
| CI-02 | Local status/watch surface | Adopt | P1 | Candidate |
| CI-03 | Review-finding schema validation | Adopt | P1 | Candidate |
| CI-04 | Governed workflow memory | Adapt | P1 | Candidate |
| CI-05 | Normalized external-review ingestion | Adapt | P1 | Candidate |
| CI-06 | Explicit wave worktrees and cleanup | Adapt | P2 | Candidate |
| CI-07 | ACP runtime adapter | Adapt | P2 | Candidate |
| CI-08 | Operational-state reconciler | Investigate | P2 | Candidate |
| CI-09 | Extension manifest and capabilities | Investigate | P3 | Candidate |
| CI-10 | Ambient global configuration | Reject | — | Closed |
| CI-11 | Automatic commits and review remediation | Reject | — | Closed |
| CI-12 | Task-frontmatter-only graph | Reject | — | Closed |

## CI-01 — Operational event ledger and replay

**Problem.** Workspaces have state, checkpoints, handoffs, leases, and command evidence, but no uniform timeline to reconstruct a failed run.

**Proposed outcome.** Add an optional append-only ledger under the workspace runtime area, such as .product/workspaces/WORK-NNN/events. It records operational facts: lease claimed, command started/finished, checkpoint written, integration planned, and task blocked. It is not a product artifact and cannot change lifecycle state.

**Modules and controls.** Event envelope, append-only writer, reader/replay API, retention policy, redaction policy, and CLI formatter. Start opt-in or only with a new runtime version; stopping or deleting the ledger must not invalidate product artifacts or approval records.

**Risks.** Command-output leakage, unbounded disk growth, and treating the ledger as truth.

**Done when.** Ordering is deterministic; appends survive restart; redaction is tested; replay reports corrupt gaps; an event cannot advance implemented, validated, or released.

## CI-02 — Local status and watch surface

**Problem.** Navigation exposes product and workspace state but lacks an operator view of a live or interrupted execution.

**Proposed outcome.** Add read-only runtime status and runtime watch views over the ledger, checkpoints, leases, and command evidence. Watch should begin as local polling; a daemon is not required.

**Boundary.** These views may show blockers and suggested owners, but never execute commands, recover leases, or approve artifacts.

**Dependencies.** CI-01 is preferred, although an initial snapshot can read existing workspace files.

**Done when.** Idle, active, expired-lease, failed-command, and stale-input fixtures render correctly; read-only mode makes no writes; JSON output is stable.

## CI-03 — Schema validation for imported review findings

**Problem.** Code Review and QA templates exist, but externally sourced feedback has no provider-neutral, mechanically validated format.

**Proposed outcome.** Define a review-finding schema and validator with stable ID, source, external reference, observed commit/diff hash when available, severity, description, status, affected scope, evidence, and suggested owner. A finding is evidence and routing input, never approval.

**Scope.** Template, validator rule, JSON output, and docs only. Do not add a provider client in this topic.

**Done when.** Missing provenance and invalid severity fail; manual findings are accepted; diff mismatch is reported; a finding cannot make a task validated.

## CI-04 — Governed workflow memory

**Problem.** Handoffs and checkpoints preserve snapshots, but related task runs can lose concise operational learning.

**Proposed outcome.** Add optional two-tier memory: task-local notes at memory/tasks/TK-NNN.md and shared cross-task context at memory/shared.md. It may contain links to approved decisions, verified discoveries, active risks, and handoff facts.

**Guardrails.** It must not contain credentials, invented requirements, approval claims, or replacements for Specifications, DEC records, QA Evidence, or Code Review. Promoted shared entries need source links; compaction preserves active risks and references; canonical repository state wins on conflict.

**Activation and removal.** Disabled by default and enabled per workspace through explicit runtime configuration. Disabling stops writes; removal changes no canonical product state.

**Done when.** Links/provenance validate; compaction preserves active risks; contradictions are reported; memory never writes approval history.

## CI-05 — Normalized ingestion of external reviews

**Problem.** GitHub, CodeRabbit, CI, and AI feedback need an auditable entry path without provider authority over delivery state.

**Proposed outcome.** Add an optional reviews import adapter family. It reads a named source, normalizes data into CI-03 findings, and proposes routing: behavior defect to Bug Fixer/Code Runner; coverage to QA; security to Security Review; ambiguous requirement/decision to Product Historian and a human.

**Non-goals.** No automatic remote thread resolution, code edits, commits, pushes, merges, approvals, or validation. Provider failure cannot roll back product state.

**Dependencies.** CI-03 and explicit decisions on credentials, network access, and data handling.

**Done when.** Every imported finding remains traceable to its source; malformed payload fails closed; manual and provider findings share one schema; a fix requires new QA and Code Review over the current diff hash.

## CI-06 — Explicit worktree execution by scheduler wave

**Problem.** The runtime can create task worktrees and compute safe waves but does not yet join them in an explicit operator-selected execution mode.

**Proposed outcome.** Add wave isolation: allocate worktrees only after a per-run choice and only where no dependency, writeScope, or sharedResources conflict exists. Plan integration in DAG order and require integrated QA afterward.

**Activation and removal.** Default off. Require an explicit isolate mode plus a current approved graph. Cleanup must be bounded and preserve failed worktrees according to retention policy.

**Risks.** Merge conflicts, disk growth, mistaken parallelism, and implicit auto-execution.

**Done when.** Overlapping paths never run concurrently; no worktree appears without activation; conflicts stop integration; cleanup is covered for success, failure, cancellation, and restart.

## CI-07 — Governed ACP runtime adapter

**Problem.** The runtime deliberately does not spawn agents. ACP could offer portable execution across runtimes, but crossing that boundary is architectural.

**Proposed outcome.** An experimental agent-runtime module dispatches exactly one ready, leased task to a selected ACP-compatible runtime, using isolation when necessary. It captures transcript and implementation evidence for normal QA and Code Review.

**Hard limits.** It cannot approve, create approval records, alter scope, bypass writeScope, commit, push, merge, mark work validated/released, or run remote/destructive actions without the existing human gate.

**Activation.** Disabled by default; explicit product/runtime configuration and per-run acknowledgement. Runtime/model selection must use manifest or flags, never ambient user configuration.

**Dependencies.** CI-01/CI-02 are strongly recommended; CI-06 supports parallel execution. An approved framework decision on ACP support and distribution is mandatory first.

**Done when.** It refuses unready, stale, or unleased tasks; records transcript and diff evidence; rejects scope escape; interruption resumes or yields a clear blocker; post-run state is at most implemented.

## CI-08 — Operational-state reconciler

**Problem.** State is distributed among task files, graph, leases, checkpoints, command evidence, worktrees, and future event records; drift needs diagnosis without automatic repair.

**Proposed outcome.** A read-only runtime reconcile command reports expired leases, orphaned worktrees, stale plans, missing evidence, diff-hash mismatch, and task/graph conflicts.

**Boundary.** It cannot repair approval records, change task status, delete worktrees, or invent evidence. Every corrective follow-up stays explicit and previewable.

**Done when.** Fixtures cover each drift type; default execution changes no files; JSON output is deterministic; every finding names an owner.

## CI-09 — Extension manifest and capability model

**Problem.** Visual adapters exist, but future execution and review adapters need a uniform declaration of version, scope, dependencies, and permissions.

**Proposed outcome.** Investigate a versioned extension manifest for read-only observers, controlled importers, and runtime adapters. Capabilities must map to actual boundaries such as artifacts.read, reviews.import, and runtime.dispatch; declaring a capability is not a sandbox.

**Required decisions.** Trust model, signing/verification, installation source, network policy, subprocess isolation, version pinning, compatibility, and revocation. Do not build an SDK before these decisions.

**Done when for a prototype.** Disabled third-party extensions do nothing; invalid capability fails; adapters cannot escape write scope; upgrade preserves adopter content when an extension is absent or incompatible.

## CI-10 — Do not add ambient global configuration

**Decision.** Rejected. The framework deliberately uses product manifests and explicit flags rather than inherited user configuration.

**Reason.** Global precedence reduces reproducibility and weakens the relationship between a pinned runtime and observed behavior.

**Allowed alternative.** An operator-local cache may hold non-semantic data such as CLI paths or event retention, but cannot alter gates, task behavior, agent selection, approvals, or external access.

## CI-11 — Do not auto-commit or auto-remediate review feedback

**Decision.** Rejected. Automation may propose a fix or produce implementation evidence, but cannot commit, push, merge, resolve a remote thread, or make a task validated.

**Reason.** The framework requires independent QA and Code Review over the same current diff hash, followed by explicit Commit Crafter and PR Finalizer stages.

**Allowed alternative.** CI-05 may route an imported finding into an approved task, which then follows the normal delivery lifecycle.

## CI-12 — Do not collapse the Execution Graph into task frontmatter

**Decision.** Rejected. Task-local frontmatter can hold metadata, but the canonical Execution Graph owns dependencies, writeScope, shared resources, source sections, requirements, acceptance checks, delivery level, and priority.

**Reason.** A frontmatter-only graph cannot preserve the framework's traceability and conflict-aware scheduling contract.

**Allowed alternative.** Improve task-file validation or derive a read-only task index, while retaining execution-graph.json as canonical.

## Suggested sequencing

1. CI-03 and CI-08: bounded validation work without new execution authority.
2. CI-01 and CI-02: observability foundation.
3. CI-04 and CI-05: controlled operational context and review ingestion.
4. CI-06: explicit scheduler-to-worktree integration.
5. Decide CI-07 and CI-09 only after operational boundaries and extension trust are approved.
