# Importação escalável de documentos

Use este fluxo para importar PRDs, wikis, Jira, Confluence, pesquisas e acervo
legado sem transformar a fonte em verdade aprovada do produto.

## Como funciona

1. O CLI encontra os arquivos permitidos e valida limites antes de copiar.
2. Cada fonte recebe hash, tamanho e identificador; o inventário é paginado.
3. As fontes são separadas em lotes (`CHUNK-NNNN`) que podem ser retomados.
4. Um revisor reclama um lote, registra evidências por fonte e suas lacunas.
5. Somente depois da revisão de todos os lotes uma pessoa pode materializar
   os mapeamentos selecionados como rascunhos.

Importar nunca aprova escopo, decisão, arquitetura ou artefato. Rascunhos
continuam sujeitos aos donos e gates normais do framework.

## Casos de uso

| Caso | Como usar |
| --- | --- |
| Produto novo com muitos documentos | Inicialize com `existing-documents` e defina orçamento. |
| Importação interrompida | Consulte o status e retome o próximo lote. |
| Arquivo grande/binário | Limite tamanho; inventarie ou rejeite binários. |
| Novas fontes após a inicialização | Crie um novo import run, separado por assunto. |

## Inicializar a partir de documentos

```powershell
spec-framework init --target ./meu-produto --agents codex `
  --starting-point existing-documents --sources ../documentos `
  --import-max-files 1000 --import-max-total-bytes 500MB `
  --import-max-file-bytes 20MB --import-chunk-size 25 `
  --import-binary-policy inventory_only --yes
```

Os padrões são 500 arquivos, 200 MB totais, 10 MB por arquivo e 25 fontes por
lote. `.git` e `node_modules` são excluídos por padrão. Se o orçamento for
excedido, o run falha antes de publicar uma cópia parcial.

## Criar um run adicional

```powershell
spec-framework import create --product-root ./meu-produto/product `
  --include "**/*.md" --exclude "**/archive/**" `
  --max-files 500 --max-total-bytes 200MB --chunk-size 25 ../novas-fontes
```

Prefira vários runs coerentes, por exemplo produto/regras, pesquisas e
arquitetura, a um único acervo difícil de revisar.

## Revisar e retomar

```powershell
spec-framework import status --product-root ./meu-produto/product --run IMPORT-001
spec-framework import resume --product-root ./meu-produto/product `
  --run IMPORT-001 --agent artifact-importer
```

`resume` retorna um lote, como `CHUNK-0003`. A lease impede revisão concorrente
do mesmo lote e expira para permitir retomada segura.

Depois de analisar o lote, forneça ao menos uma evidência para cada fonte:

```json
{
  "source_evidence": {
    "SRC-000051": [{ "locator": "Seção 2.3", "claim": "Clientes exportam relatórios." }]
  },
  "gaps": { "SRC-000051": ["Formato de exportação não definido."] }
}
```

```powershell
spec-framework import record-review --product-root ./meu-produto/product `
  --run IMPORT-001 --chunk CHUNK-0003 --agent artifact-importer `
  --input ./review-chunk-0003.json --yes
```

Isso registra evidência e lacunas; não cria aprovações nem altera artefatos do
produto.

## Materializar rascunhos

Após revisar conflitos, lacunas e `mapping.json`, uma pessoa autoriza:

```powershell
spec-framework import materialize --product-root ./meu-produto/product `
  --run IMPORT-001 --approved-by "Product Owner" --yes
```

O comando bloqueia se existir lote pendente, em revisão ou bloqueado. Também
recusa destino duplicado, caminho fora do produto e sobrescrita de conteúdo.

## Onde ficam os dados

```text
product/knowledge/imports/
├── sources/IMPORT-001/                 # cópias preservadas
└── runs/IMPORT-001/
    ├── import-config.json               # filtros e limites
    ├── inventory/index.json             # índice compacto
    ├── inventory/pages/PAGE-0001.jsonl  # páginas do inventário
    ├── chunks/CHUNK-0001.json           # status e lease
    ├── traceability/CHUNK-0001.json     # evidências por fonte
    └── mapping.json                     # rascunhos propostos
```

Para dados sensíveis, decida previamente a retenção e quem pode ler as cópias
preservadas. Não apague import runs como limpeza: eles são evidência auditável.
