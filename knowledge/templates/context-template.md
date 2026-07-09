# Context: [artifact name]

## 🧭 Snapshot

```yaml
id: [DOMAIN-001 | GOAL-001 | FT-001 | UC-001 | SPEC-001 | TK-001]
type: [domain | goal | feature | use-case | specification | implementation-plan | execution-graph | task]
name: [human readable name]
status: [draft | proposed | approved | in_progress | implemented | validated | released | deprecated | superseded]
owner_skill: [skill name]
last_updated: [YYYY-MM-DD]
delivery:
  level: [L0 | L1 | L2 | L3 | L4 | L5 | N/A]
  priority: [P0 | P1 | P2 | P3 | N/A]
  depends_on:
    - [artifact id/path]
  rationale: [why this artifact belongs here]
```

## 📌 Purpose

[One paragraph explaining why this artifact exists and what decision or work it enables.]

## 🗺️ Artifact Map

```mermaid
flowchart TD
  P["Parent artifact"] --> C["This context"]
  C --> CH["Child artifact"]
  C --> D["Decision"]
  C --> N["Next skill"]

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#fef3c7,stroke:#d97706,color:#78350f,stroke-width:3px;
  classDef pending fill:#f8fafc,stroke:#94a3b8,color:#334155;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class P done;
  class C current;
  class CH,D,N pending;
```

## 🔗 Relationships

| Type | Artifact | Path | Relationship |
| --- | --- | --- | --- |
| Parent | `[id]` | `[path]` | `[relationship]` |
| Child | `[id]` | `[path]` | `[relationship]` |
| Related | `[id]` | `[path]` | `[relationship]` |

## 🚧 Dependencies

| Dependency | Why Needed | Blocking | Status |
| --- | --- | --- | --- |
| `[id/path]` | `[reason]` | `[yes/no]` | `[open/ready/blocked]` |

## 📂 Canonical Documents

| Document | Path |
| --- | --- |
| Primary | [`[primary document]`]([path]) |
| Specification | [`[specification]`]([path or N/A]) |
| Design | [`[design]`]([path or N/A]) |
| Implementation plan | [`[implementation plan]`]([path or N/A]) |
| Execution graph | [`[execution graph]`]([path or N/A]) |
| Tasks | [`[tasks]`]([path or N/A]) |

## 🔐 Decisions

| Decision | Summary | Status |
| --- | --- | --- |
| `[DEC-XXX]` | `[summary]` | `[status]` |

## ⚠️ Assumptions And Open Questions

| Type | Item | Owner | Blocks |
| --- | --- | --- | --- |
| Assumption | `[assumption]` | `[role]` | `[artifact/status]` |
| Question | `[question]` | `[role]` | `[artifact/status]` |

## 🏁 Handoff

| Field | Value |
| --- | --- |
| Next recommended skill | `[skill]` |
| Required reading | [`[artifact]`]([path]) |
| Stop condition | `[approval gate/blocker]` |
