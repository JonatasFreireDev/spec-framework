# Plano de feature — Importação escalável e retomável

## Objetivo

Permitir que `existing-documents` trate acervos grandes com previsibilidade de
volume, lotes retomáveis e rastreabilidade por arquivo e por trecho, sem mudar
a tese do framework: a fonte importada é evidência, rascunhos continuam sendo
do produto, e nenhuma aprovação ou materialização ocorre automaticamente.

## Estado atual e lacuna

Hoje `sourceimport.CreateRun` expande toda a árvore, copia cada arquivo,
calcula hash e grava um inventário completo. A cópia e o hash usam streaming,
mas a lista de arquivos e os ledgers são montados integralmente na memória.
Não há filtros, orçamento, exclusões padrão, lote, checkpoint ou retomada.
O Artefact Importer deve revisar cada fonte, o que torna um arquivo muito longo
ou milhares de documentos operacionalmente difícil.

## Princípios e limites invioláveis

- Nunca alterar, aprovar ou sobrescrever conteúdo adotante durante inventário,
  análise, retomada ou compactação.
- Tratar fontes como evidência; somente mapeamentos humanos selecionados podem
  gerar artefatos `draft` com `provenance.import_run`.
- Preservar hash, caminho original, tamanho, formato e evidência de cada fonte.
- Não usar configuração ambiente para semântica; o plano do lote é persistido
  dentro do próprio import run e flags explícitas apenas o criam.
- Não enviar conteúdo para serviços externos. Extratores e adaptadores devem
  ser declarados, opcionais e sem autoridade de aprovação.

## Proposta de arquitetura

### Módulo 1 — Plano de ingestão

Novo `import create` recebe fontes e cria `import-config.json` imutável após o
início:

```json
{
  "schema_version": 2,
  "include": ["docs/**/*.md", "research/**/*.pdf"],
  "exclude": ["**/node_modules/**", "**/.git/**", "**/*.zip"],
  "max_files": 500,
  "max_total_bytes": 209715200,
  "max_file_bytes": 10485760,
  "chunk_size": 25,
  "binary_policy": "inventory_only"
}
```

Padrão seguro: excluir `.git`, dependências e diretórios de build conhecidos;
inventariar binários sem tentar extrair conteúdo; falhar antes de copiar quando
um orçamento é excedido. `--allow-large-file` e `--allow-total-overflow` apenas
criam uma exceção explicitamente registrada, nunca silenciosa.

### Módulo 2 — Inventário paginado

Substituir o array único em memória por arquivos ordenados e paginados em
`runs/IMPORT-NNN/inventory/pages/PAGE-0001.jsonl`. Um índice compacto mantém
contagem, bytes, hash do plano e páginas concluídas. O inventário original v1
continua legível; novos runs usam v2. Cada linha contém `source_id`, caminho
original, cópia preservada, hash, tamanho, formato, classificação e motivo de
exclusão quando aplicável.

### Módulo 3 — Lotes de análise

`chunks/CHUNK-0001.json` referencia no máximo `chunk_size` source IDs, com
status `queued | analysing | reviewed | blocked | excluded`. A análise não
altera o inventário; produz `traceability/chunks/CHUNK-0001.json` e atualiza um
índice de cobertura. Para documentos longos, um `segment` contém um localizador
estável (página, capítulo, intervalo de linhas ou timestamp) e o hash do trecho
extraído. O texto completo não é duplicado no ledger.

### Módulo 4 — Retomada e consistência

`import status --run` é somente leitura. `import resume --run` escolhe apenas
o próximo lote `queued` ou um lote abandonado com lease expirada. Cada lote
recebe lease, heartbeat e checkpoint operacional equivalentes aos do runtime,
mas não recebe autorização de escrita em artefatos canônicos. Uma mudança de
hash da fonte marca seus lotes e mapeamentos como `stale`; não apaga análise
anterior.

### Módulo 5 — Materialização por seleção

`import materialize` aceita somente mappings cujas fontes estejam revisadas ou
explicitamente `not_applicable` com justificativa. A operação é atômica por
seleção: valida todos os destinos antes de escrever, não sobrescreve existentes
e atualiza traceabilidade somente após todos os rascunhos serem publicados.
O formato atual de `mapping.json` continua aceito; v2 adiciona `source_ids` e
`source_segments` em paralelo a `source_documents` durante a migração.

## CLI proposta

```text
spec-framework import create --source <dir|file>... --include <glob>...
  --exclude <glob>... --max-files 500 --max-total-bytes 200MB
  --max-file-bytes 10MB --chunk-size 25 --binary-policy inventory_only

spec-framework import status --run IMPORT-001 [--json]
spec-framework import resume --run IMPORT-001 --agent <id> [--chunk CHUNK-0003]
spec-framework import retry --run IMPORT-001 --chunk CHUNK-0003 --yes
spec-framework import materialize --run IMPORT-001 --approved-by <human> --yes
```

`resume` não executa um modelo nem materializa rascunhos: prepara e reclama o
lote para o Artefact Importer. A análise continua sendo responsabilidade do
skill e do humano que revisa os resultados.

## Dados, ativação e remoção

| Item | Padrão | Ativação | Remoção segura |
| --- | --- | --- | --- |
| Inventário v2 | Novo run | `import create` | manter run como evidência; não apagar em upgrade |
| Filtros/orçamentos | Conservadores | flags gravadas no config | criar novo run, não reescrever config |
| Lotes/leases | Desligados para v1 | run v2 | expirar/reconciliar, nunca apagar análise |
| Extrator de PDF/DOCX | Ausente | extensão declarada | voltar a `inventory_only` |
| Materialização v2 | Igual ao atual | aprovação + `--yes` | não reversível automaticamente; drafts seguem fluxo normal |

## Compatibilidade e upgrade

- Nenhum `IMPORT-NNN` existente é reescrito. O leitor detecta v1 e apresenta
  status legado; uma migração explícita cria um índice v2 derivado, preservando
  JSONs originais e hashes.
- Produtos novos usam v2 quando a feature estiver estável. Produtos existentes
  só adotam v2 em novos runs ou via `import migrate --run ... --yes`.
- `upgrade` nunca modifica fontes copiadas, traceabilidade, mapeamentos,
  aprovações, rascunhos ou checkpoints do adotante.

## Entregas incrementais

| Fase | Entrega | Esforço | Critério de aceite |
| --- | --- | --- | --- |
| 0 | DEC de política de limites e binários | P | decisão aprovada e contratos atualizados |
| 1 | Filtros, orçamentos e inventário v2 paginado | M | falha previsível antes de copiar além do orçamento |
| 2 | Chunks, status e retomada sem análise automática | M | interromper e retomar não duplica nem perde fontes |
| 3 | Segmentos e traceabilidade por trecho | M | conclusão aponta localizador estável da fonte |
| 4 | Mapeamento/materialização v2 atômicos | M | preserva regra de não sobrescrever e aprovação explícita |
| 5 | Extratores opcionais e observabilidade | G | extensão desabilitada não altera o fluxo padrão |

## Testes de combinação obrigatórios

- 10 mil arquivos pequenos: memória limitada, páginas determinísticas e ordem
  estável.
- Arquivo acima do limite; total acima do limite; exceção humana registrada.
- Glob include/exclude sobre caminhos Windows e Unix; `.git`, symlinks, binários
  e permissões negadas.
- Interrupção durante cópia, hash, criação de lote e atualização de checkpoint.
- Retomada com hash de fonte alterado, lease expirada e lote bloqueado.
- Dois agentes tentando o mesmo lote; lotes distintos em paralelo.
- Mistura de run v1 e v2; upgrade e migração sem reescrita de conteúdo adotante.
- Materialização com uma fonte não revisada, destino duplicado e falha no meio:
  nenhum artefato parcial pode permanecer.
- Confirmação de que nenhum caminho cria aprovação, commit, push, merge ou
  resolve review remoto.

## Superfícies afetadas

`internal/sourceimport`, `internal/cli`, `internal/validator`,
`framework/skills/artifact-importer`, templates de importação, starter,
`FRAMEWORK.md`, `docs/starting-points.md`, testes de CLI e de upgrade. Avaliar
uma extensão de extratores somente depois da Fase 4; ela não entra no núcleo.

## Decisões humanas pendentes

1. Quais limites padrão são aceitáveis (arquivos, bytes totais e arquivo único)?
2. A cópia integral para `knowledge/imports/sources/` continua obrigatória para
   todos os formatos, ou referências externas com hash são permitidas?
3. Quais formatos terão extrator oficial inicial (Markdown apenas, PDF, DOCX)?
4. Um lote pode ser marcado `not_applicable` por um humano sem leitura integral?
5. Qual política de retenção se aplica às fontes que contenham dados sensíveis?

## Não recomendações

Não fazer uma única chamada de agente para “entender toda a pasta”, não gravar
resumos sem fonte/trecho, não paralelizar dois agentes sobre o mesmo arquivo,
não transformar importação em aprovação e não introduzir banco de dados como
fonte canônica. JSON paginado e Markdown continuam auditáveis e preserváveis.
