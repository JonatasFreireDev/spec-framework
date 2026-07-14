# FDR-045: Canonical Method Compression

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-045` |
| Status | `proposed` |
| Origin EV | `Agent context efficiency review` |
| Date | `2026-07-13` |
| Owner | `Framework Maintainer` |

## Context

`FRAMEWORK.md` is mandatory reading for framework skills and had grown beyond one thousand lines. It mixed method invariants with an exhaustive physical tree, command catalog, starting-point procedures, recommended prompts, and a roadmap that described already delivered capabilities. These details duplicated the CLI help, declarative init contracts, skills, templates, guides, and FDRs. The distributed agent instructions also concentrated several operational gates in one dense paragraph without stating source precedence.

## Decision

| Boundary | Contract |
| --- | --- |
| Canonical method | Keep `FRAMEWORK.md` as one mandatory canonical document; do not split the method into independently optional fragments. |
| Compression | Preserve concepts, ownership, gates, traceability, approval, safety, and preservation invariants while removing duplicated inventories and procedures. |
| Physical structure | Document stable ownership areas and use-case bundle responsibility. Exact paths belong to starter assets, registries, templates, and initialization contracts. |
| Commands | CLI generated help owns current command syntax; the method documents behavioral boundaries only. |
| Starting points | Detailed selection and procedures belong to `docs/starting-points.md`, `framework/init/`, generated `BOOTSTRAP.md`, and their FDRs. |
| Skills | `FRAMEWORK.md` defines layers and ownership boundaries; versioned `SKILL.md` contracts own specialist procedure and handoff. |
| Prompts | Remove recommended prompts when executable skill contracts express the stronger behavior. |
| Roadmap | Remove the stale version roadmap; planned method changes require an FDR or explicit evolution artifact. |
| Agent bootstrap | `AGENTS.framework.md` remains a concise operational bootstrap and declares authority order from method through current CLI evidence. |
| Semantics | This is a non-semantic compression. Existing gates and adopter preservation rules remain in force. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Mandatory method reading consumes less agent context and contains fewer competing instructions. |
| Positive | Volatile command and structure details remain with mechanically maintained sources. |
| Positive | Agent instructions expose precedence and operational gates more clearly. |
| Negative | Readers needing exact paths or command flags must follow the named authoritative source. |
| Negative | Future changes must resist copying detailed procedures back into the canonical method. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `4. Folder Structure` | Replace exhaustive inventory with ownership areas and source boundaries. |
| `9. Skills` and `10. Orchestrators` | Retain ownership and safety boundaries while delegating procedure to versioned skill contracts. |
| `15. How Agents Use The Framework` | Retain activation, routing, initialization, preservation, authority, and lifecycle invariants without duplicating guides or CLI help. |
| `16. Framework Roadmap` | Remove stale roadmap and renumber Final Rule. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |
| Distributed agent instructions | [../AGENTS.framework.md](../AGENTS.framework.md) |
| Skill catalog | [../skills/README.md](../skills/README.md) |
| Starting-point guide | [../../docs/starting-points.md](../../docs/starting-points.md) |
| Initialization contracts | [../init/contracts/new-product.json](../init/contracts/new-product.json) |

## Supersedes

- N/A; this refactors presentation without superseding prior method decisions.
