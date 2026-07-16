# Despacho supervisionado de subagents

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
