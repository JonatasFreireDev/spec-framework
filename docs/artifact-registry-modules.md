# Modular Artifact Registry

This document is the canonical maintenance map for composing the initial artifact registry by starting point. It separates installed asset sets from artifacts that receive identity, status, parent relationships, and approval evidence.

## Starting points and composition

| Starting point | Asset sets | Artifacts propres | Exclusões |
| --- | --- | --- | --- |
| `new-product` | Foundation completa, Knowledge, Domains, Design, Engineering, Governance | Problem, Vision, Principles, North Star, Strategy | Nenhuma |
| `existing-feature` | Core, Knowledge, Domains, Governance | Feature Brief; Features recebem `FEATURE-BRIEF-TBD` como parent | Problem, Vision, Principles, North Star, Strategy |
| `existing-product` | Core, Strategy, Knowledge, Domains, Design, Engineering, Governance | Product Baseline; Strategy recebe Baseline como parent | Problem, Vision, Principles, North Star |
| `existing-implementation` | Core, Foundation completa, Knowledge, Domains, Design, Engineering, Governance | Implementation Assessment; Problem recebe Assessment como parent | Nenhuma |
| `existing-documents` | Core, Knowledge, Governance | Registry mínimo e import run | Foundation e delivery artifacts |
| `audit-only` | Core, Security Baseline, Governance | Nenhum novo artifact de produto | Foundation e delivery artifacts |

## Artifact modules

| Module | Artifacts principais | Pode ser incluído quando |
| --- | --- | --- |
| `product-core` | Manifest, contexto raiz e ferramentas comuns | Sempre |
| `foundation-core` | Problem, Vision, Principles, North Star, Strategy | O produto ainda precisa definir direção |
| `foundation-feature` | Feature Brief | O ponto de entrada é uma feature existente |
| `foundation-baseline` | Product Baseline | O produto existente tem código e operação, mas pouca documentação |
| `foundation-assessment` | Implementation Assessment | É necessário mapear uma implementação antes da Foundation |
| `foundation-strategy` | Strategy | O starting point preserva apenas estratégia |
| `delivery-domain` | Domain, Goal, Feature, Use Case | O produto está pronto para modelar entrega |
| `delivery-design` | Design System, Design, componentes, padrões | O produto possui interface ou sistema visual |
| `delivery-engineering` | Engineering System, Quality System, Engineering Proposal, Review | O produto precisa de contratos técnicos compartilhados |
| `delivery-validation` | Tests, QA Evidence, Security Review, Audit | A entrega avançou para validação |
| `delivery-execution` | Execution Graph, Task Set, Task | A especificação está pronta para implementação |
| `governance` | Decisions, approvals, releases, audits | Sempre que houver governança aplicável |
| `knowledge` | Regras, convenções, decisões e imports | O produto precisa de conhecimento operacional |
| `knowledge-security-baseline` | Security Baseline | O produto está em modo de auditoria ou possui superfície sensível |

## Approval adapters

The approval engine remains generic. Adapters are only for side effects or composite contracts:

| Adapter | Responsibility |
| --- | --- |
| `generic` | Update document, registry, and approval record |
| `foundation-context` | Synchronize the canonical document status with its `context.md` |
| `engineering-system` | Synchronize context, YAML, Quality System, and composite hash |
| `design-system` | Validate tokens/catalog when applicable |
| `feature-brief` | Validate `targetFeature` and the selected feature |

Adding an artifact type must not require a new approval command. Add the type to the selected module, declare its registry entry and adapter, and add specialized validation only when its contract needs it.
