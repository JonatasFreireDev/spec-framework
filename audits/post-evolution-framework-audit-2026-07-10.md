# Post-Evolution Framework Audit - 2026-07-10

## Executive snapshot

| Item | Status | Evidence |
| --- | --- | --- |
| Framework validator | ✅ ready | `node engineering\validators\framework-validator.mjs --write-registry --write-report` returned `Verdict: ✅ ready (0 errors, 0 warnings, 0 notes)`. |
| Syntax checks | ✅ passed | `node --check` passed for [framework-validator.mjs](../engineering/validators/framework-validator.mjs) and [move-artifact.mjs](../engineering/move-artifact.mjs). |
| Decisions indexed | ✅ complete | [.product/decisions.json](../.product/decisions.json) lists `DEC-001` through `DEC-009` as approved. |
| Artifact registry | ✅ generated | [.product/artifacts.json](../.product/artifacts.json) contains 60 artifacts after registry refresh. |
| Framework scope | 🟡 ready with notes | Framework mechanics are green; initial engineering tool tests now cover validator approval/staleness gates and move-tool rewriting. Product examples still include draft/blocked QA and security evidence by design. |

**Verdict:** 🟡 `ready_with_notes`

The framework evolution chain is coherent and mechanically enforced. Initial fixture tests now cover the validator and move tool; the remaining work is broader test coverage and keeping product examples from being mistaken for implementation-ready work.

## Evolution coverage

```mermaid
flowchart LR
  A["EV-001<br/>Approval records"] --> B["EV-004<br/>Derived staleness"]
  B --> C["EV-002<br/>Task records"]
  C --> D["EV-005<br/>Code evidence"]
  D --> E["EV-003<br/>Rigor tiers"]
  E --> F["EV-006<br/>Scoped identity + move tool"]
  F --> G["Audit<br/>current"]

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#dbeafe,stroke:#2563eb,color:#1e3a8a;
  classDef pending fill:#f3f4f6,stroke:#9ca3af,color:#374151;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class A,B,C,D,E,F done;
  class G current;
```

| Evolution | Audit result | Evidence |
| --- | --- | --- |
| EV-001 approval records | ✅ enforced | [AGENTS.md](../AGENTS.md) requires agents to stop on missing approval records; validator has approval record checks. |
| EV-004 staleness | ✅ enforced | [FRAMEWORK.md](../FRAMEWORK.md) defines stale as derived from [.product/derivations.json](../.product/derivations.json). |
| EV-002 task records | ✅ enforced | Task files are canonical and taskset indexes are generated/validated. |
| EV-005 code evidence | ✅ enforced | Implemented/validated task gates require branch, commits, PR, test evidence, and security evidence. |
| EV-003 rigor tiers | ✅ enforced | Validator checks `rigor_tier` and Tier L trigger rules. |
| EV-006 identity and move tooling | 🟡 usable with manual review | [.product/ids.json](../.product/ids.json) declares slug-scoped identity; move dry-run reports rewritten files and free-text mentions requiring review. |

## Findings

| Severity | Finding | Evidence | Recommendation |
| --- | --- | --- | --- |
| 🟢 Info | The framework validates cleanly after all approved evolutions. | Validator output: `0 errors, 0 warnings, 0 notes`. | Keep validator as required CI/local gate before commits. |
| 🟢 Resolved | Engineering scripts now have initial fixture tests. | [engineering/tests/run-tests.mjs](../engineering/tests/run-tests.mjs) covers approval-record blocking, derived staleness blocking, Phase A writeScope warnings, concrete QA evidence enforcement, Markdown link rewrite, JSON path rewrite, and free-text mention reporting. | Expand coverage for task-file validation, code-evidence gates, rigor-tier gates, Mermaid semantic bindings, and future Phase B writeScope errors. |
| 🟡 Medium | Move tooling intentionally reports free-text mentions instead of rewriting them, which creates a manual review step after moves. | `move-artifact --dry-run` reported 5 rewritten files and 90 free-text mentions requiring review for the organizer use case path. | Keep this behavior, but require the move report to be attached to the audit or PR when a move is executed. |
| 🟡 Medium | Product example artifacts are not implementation-ready, especially QA and security evidence. | Event use case QA/security files contain `blocked`, `not run`, and pending role/test evidence. | Treat them as framework examples until a product owner approves role, rollout, QA, and security evidence. |
| 🟢 Info | No mojibake pattern was found in Markdown files during this audit. | `rg "ð|âœ|âž|ï¸|Ÿ" --glob "*.md"` returned no matches. | No encoding cleanup is needed right now. |

## Dependency and gate view

```mermaid
flowchart TD
  A["Approved decisions<br/>DEC-001..DEC-009"] --> B["Framework rules<br/>FRAMEWORK.md + AGENTS.md"]
  B --> C["Validator<br/>approval, stale, tiers, tasks, code evidence"]
  C --> D["Generated registries<br/>artifacts + readiness reports"]
  C --> E["Product examples<br/>events QR check-in"]
  E --> F["QA evidence"]
  E --> G["Security review"]
  F --> H["Validation / release"]
  G --> H

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#dbeafe,stroke:#2563eb,color:#1e3a8a;
  classDef pending fill:#f3f4f6,stroke:#9ca3af,color:#374151;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class A,B,C,D done;
  class E current;
  class F,G,H blocked;
```

## Files changed by this audit

| File | Reason |
| --- | --- |
| [.product/artifacts.json](../.product/artifacts.json) | Regenerated by validator with current registry data. |
| [audits/framework-validation-report.md](framework-validation-report.md) | Regenerated validator report. |
| [audits/readiness/framework-readiness.md](readiness/framework-readiness.md) | Regenerated readiness report. |
| [audits/post-evolution-framework-audit-2026-07-10.md](post-evolution-framework-audit-2026-07-10.md) | New consolidated audit report. |

## Incomplete or intentionally blocked

| Area | Status | Notes |
| --- | --- | --- |
| Framework mechanics | ✅ complete for approved EVs | No blocking validator issues. |
| Script test harness | 🟡 partial | Initial fixture tests exist; broader gate coverage is still recommended. |
| Product QA evidence | 🔴 blocked | Event examples contain planned but not executed QA evidence. |
| Product security evidence | 🔴 blocked | Security reviews remain blocked until implementation evidence and role decisions exist. |
| Approval records | ✅ complete | Current approved+ artifacts have records; agents must not repair them without explicit migration approval. |

## Human approval questions

| Question | Recommendation |
| --- | --- |
| Should script fixture tests become mandatory before the next framework evolution? | ✅ Accept. Initial tests exist; require them before changing validator or move-tool behavior. |
| Should every executed move produce a saved audit artifact? | ✅ Accept with adjustment: require it for non-dry-run moves that touch more than one artifact subtree. |
| Should product examples advance toward implementation readiness now? | 🟡 Only if this repo is acting as the product repo. As a framework lab, keep examples demonstrative. |

## Recommended next steps

| Priority | Next step | Owner skill |
| --- | --- | --- |
| P1 | Expand fixture-based tests for task-file validation, code evidence, rigor tiers, and Mermaid semantic bindings. | `audit-orchestrator` + engineering |
| P2 | Add a generated move report option to `move-artifact.mjs` for executed moves. | `impact-analyzer` + engineering |
| P2 | Add CI guidance showing the exact validator and script checks expected before merge. | `documentation-orchestrator` |
| P3 | Decide whether the events examples should stay illustrative or become a real product slice. | `product-orchestrator` |

## Validation performed

| Command | Result |
| --- | --- |
| `node engineering\validators\framework-validator.mjs --write-registry --write-report` | ✅ `Verdict: ✅ ready (0 errors, 0 warnings, 0 notes)` |
| `node --check engineering\validators\framework-validator.mjs` | ✅ passed |
| `node --check engineering\move-artifact.mjs` | ✅ passed |
| `node --check engineering\tests\run-tests.mjs` | ✅ passed |
| `node engineering\tests\run-tests.mjs` | ✅ `11/11 tests passed` |
| `node engineering\move-artifact.mjs --from ... --to ... --dry-run` | ✅ passed; reported rewrites and free-text review items |
| `rg "ð\|âœ\|âž\|ï¸\|Ÿ" --glob "*.md"` | ✅ no matches |

## Final result

| Verdict | Blockers | Next owner |
| --- | --- | --- |
| 🟡 `ready_with_notes` | None for framework mechanics; product examples remain blocked for real implementation validation. | Audit Orchestrator, then engineering test hardening. |
