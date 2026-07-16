# Plano de evolução — Subagents governados

## Objetivo

Evoluir dos contratos atuais de tasks, chunks, leases, grafo e evidências para
uma orquestração assistida. O objetivo é paralelismo controlado sem delegar
decisões de produto, aprovação humana ou autoridade de entrega ao runtime.

## Limites permanentes

- Grafo e chunks são as unidades canônicas; o orquestrador não inventa trabalho.
- Cada subagent recebe escopo, leitura obrigatória, lease e evidência esperada.
- QA, Code Review e Security Review continuam independentes e read-only.
- Nenhum subagent aprova, commita, envia, faz merge, publica ou resolve review.
- Artefatos do produto permanecem auditáveis e pertencem ao adotante.

## Fase 0 — Fundamentos

Definir um `DispatchEnvelope` com unidade, agente, hash de entrada, leituras,
write scope, recursos compartilhados, evidências, status e operações proibidas.
Estados: `candidate`, `assigned`, `running`, `returned`, `reviewed`, `closed`,
`blocked`, `expired` e `cancelled`.

Aceite: atribuição e retorno são reconstruíveis sem conversa anterior.

## Fase 1 — Planejamento read-only

Criar `dispatch plan --work WORK-NNN`. Mostra tasks/chunks elegíveis,
dependências, ondas, conflitos, risco e especialista recomendado. Não cria
agente, lease ou arquivo de produto.

## Fase 2 — Atribuição explícita

Criar `dispatch assign --unit <id> --agent <id> --yes`. Cria lease, envelope
imutável, handoff inicial e evento operacional. Falha para unidade não pronta,
dependência pendente, lease existente ou conflito de escopo.

Padrão desligado por configuração do produto. Remoção segura: impedir novas
atribuições e preservar envelopes/evidências existentes.

## Fase 3 — Retorno e observação

Criar `dispatch return`, `dispatch observe` e `dispatch reconcile`. O retorno
exige resumo, hashes, evidências, bloqueios e rota sugerida. Observação mostra
leases, trabalho ativo, worktrees e gates pendentes. Reconciliação só reporta
estado órfão; não repara nada.

## Fase 4 — Execução local supervisionada

Iniciar somente um harness explicitamente configurado para um envelope já
atribuído. Código exige worktree; importação exige chunk. O processo recebe
envelope e diretórios permitidos, não uma instrução livre. Logs e transcrições
são evidência operacional.

Guardrails: limite de concorrência, timeout, heartbeat, cancelamento explícito,
allowlist de harnesses e nenhuma credencial externa por padrão.

## Fase 5 — Ondas supervisionadas

Criar `dispatch wave --wave WAVE-003 --max-parallel 3 --yes`. Inicia apenas
unidades atribuídas, prontas e sem escopo/recurso compartilhado sobreposto.
Novas unidades param quando surge blocker; as já iniciadas não são apagadas.

Primeiro rollout: chunks de importação. Segundo rollout: tasks em worktrees.

## Fase 6 — Retorno para gates independentes

Retornos viram candidatos a QA, Code Review e Security Review, nunca aprovação.
O orquestrador verifica hashes/evidências, mas não altera estados canônicos.

## Fase 7 — Recomendações adaptativas

Depois de uso comprovado, recomendar tamanho de lote, capacidade, gargalos e
prioridade por dependência. Recomendações exigem confirmação e não mudam escopo
ou critérios de aprovação.

## Contrato mínimo

```json
{
  "dispatch_id": "DISPATCH-001",
  "unit": { "kind": "task", "id": "TK-014" },
  "agent": "code-runner-1",
  "input_hash": "sha256:...",
  "required_reading": ["specification.md", "tasks/TK-014.md"],
  "write_scope": ["src/billing"],
  "expected_evidence": ["test log", "diff hash"],
  "status": "assigned",
  "forbidden": ["approval", "commit", "push", "merge", "release"]
}
```

## Testes obrigatórios

- Duas atribuições concorrentes para a mesma unidade.
- Dependência pendente, conflito de escopo e recurso compartilhado.
- Processo perdido, lease expirada, retorno duplicado e envelope órfão.
- Hash incompatível, evidência ausente e escrita fora do escopo.
- Wave parcialmente iniciada, cancelamento e retomada.
- Tentativas de aprovação, commit, push, merge, release e review externo.
- Upgrade e remoção sem apagar histórico operacional.

## Métricas

- Tempo entre unidade pronta e atribuída.
- Taxa de lease expirada, conflito e retorno bloqueado.
- Cobertura de evidência por retorno.
- Tempo até as revisões independentes.

## Decisões humanas necessárias

1. Harnesses locais suportados inicialmente.
2. Limite de concorrência por máquina e workspace.
3. Tipos de task permitidos após o piloto de importação.
4. Retenção e acesso a transcrições.
5. Política de cancelamento e recuperação de leases expiradas.

## Não recomendações

Não criar agentes ilimitados, não automatizar decisões de produto, não permitir
memória compartilhada sem proveniência e não automatizar merge/aprovação. Cada
fase só avança após demonstrar auditabilidade e critérios de aceite.
