# Despacho supervisionado de subagents

## Engineering baseline specialists

Engineering Orchestrator can persist a baseline handoff in either `sequential`
or `delegated` mode. Sequential is the compatible default. In delegated mode,
the parent agent uses the harness-native subagent capability; the CLI only
validates and stores the envelope.

```powershell
spec-framework dispatch assign --product-root product --work WORK-001 --task .product/workspaces/WORK-001/engineering-handoff.json --role technical-landscape --agent landscape-1 --yes
spec-framework dispatch return --product-root product --work WORK-001 --id DISPATCH-... --agent landscape-1 --summary "technical graph mapped" --evidence "catalog and topology" --output-hashes "engineering/catalog/catalog.yaml=<sha256>" --blockers "" --decision-candidates "" --yes
```

The canonical phase order is Technical Landscape; Standards and Operations in
a bounded parallel phase; Evidence; then Engineering System aggregation. Later
assignments provide returned dispatch IDs through `--depends-on`. Every return
must identify product-relative outputs and their SHA-256 hashes. The runtime
rejects missing dependencies, stale or unreturned prerequisites, concurrency
above `max_parallel`, outputs outside the declared scope, and hash mismatches.
Blockers and decision candidates are persisted as structured envelope fields;
a blocked dependency cannot unlock a later specialist phase.

`dispatch` coordena subagents por envelopes persistidos. Ele não aprova,
commita, envia, faz merge ou publica.

```powershell
spec-framework dispatch plan --product-root product --graph <graph>
spec-framework dispatch assign --product-root product --work WORK-001 --graph <graph> --task TK-001 --agent runner --yes
spec-framework dispatch return --product-root product --work WORK-001 --id DISPATCH-... --agent runner --summary "feito" --diff-hash <hash> --evidence "test log" --yes
spec-framework dispatch assign --product-root product --work WORK-001 --role qa --parent DISPATCH-... --agent qa-1 --yes
spec-framework dispatch reconcile --product-root product --work WORK-001
```

QA, Code Review e Security Review recebem envelopes read-only presos ao mesmo
`diff_hash` retornado pelo Code Runner. A execução local é experimental e exige
`dispatch run --enable --yes`; ondas usam envelopes já atribuídos e
`dispatch wave --enable --yes`.

Antes de executar, habilite a capability no produto com `dispatch configure` e
liste os harnesses permitidos. `dispatch wave` recebe uma `--wave` persistida
do scheduler, não IDs arbitrários. Transcripts locais seguem a retenção definida
em `--transcript-retention`.

Também são suportados `--role artifact-importer --run IMPORT-NNN --chunk
CHUNK-NNNN`, `--role threat-modeler --task <boundary-path>` e `--role
technical-discovery --task <question-path>`. Esses três papéis são read-only:
retornam evidência, opções, ameaças ou lacunas, sem materializar drafts, criar
decisões, aceitar risco residual ou editar Engineering Proposals.

Para concluir um chunk despachado, use `dispatch return --review-input
<review.json> --yes`; o JSON é validado e registrado como evidência por fonte.
O executor também confere os caminhos alterados pelo Git contra o `writeScope`
da task e falha quando detectar escrita fora do contrato.
