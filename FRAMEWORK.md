# Product Engineering Framework

Este documento consolida a arquitetura do framework de Product Engineering orientado por IA. Ele substitui o historico longo de conversa por uma fonte de verdade navegavel e executavel por agentes como Codex.

O objetivo nao e apenas organizar documentos. O objetivo e criar uma esteira em que produto, especificacao, planejamento, execucao e auditoria formam um sistema unico.

## 1. Tese

Este framework trata documentacao como infraestrutura de engenharia.

Em vez de pedir para uma IA "criar arquivos" ou "implementar uma feature" a partir de contexto solto, o produto passa por uma cadeia explicita:

```text
Problem -> Vision -> Strategy -> Domain -> User Goal -> Feature -> Use Case -> Specification -> Design -> Implementation Plan -> Execution Graph -> Tasks -> Code -> Validation -> Audit
```

A Specification e o contrato central. Ela une produto, UX, regras, arquitetura, dados, analytics, seguranca, testes e criterios de aceite. Tasks nao sao checklists soltos: sao unidades executaveis derivadas de uma Specification, ordenadas por um Execution Graph.

Design e um artefato de planejamento quando a entrega tiver interface. Ele traduz a Specification em fluxo visual, estados, wireframes ou mockups revisaveis antes do Implementation Plan, para que a engenharia nao descubra a experiencia somente durante a implementacao.

## 2. Principios

### Product first

Toda decisao tecnica deve manter rastreabilidade com uma necessidade de produto. A arvore sempre nasce do problema, passa pela visao e chega ao codigo por meio de dominios, objetivos, features e casos de uso.

### Domain driven

Dominios sao o centro da documentacao funcional. Nao existe uma pasta global de features como fonte principal. Uma feature pertence a um User Goal, e um User Goal pertence a um Domain.

### Specification driven

Historias podem existir para backlog ou comunicacao, mas nao sao o contrato principal para IA. O contrato principal e a Specification, porque ela descreve o comportamento completo e reduz ambiguidades antes da implementacao.

### Context driven

Todo nivel importante deve ter um `context.md`. Esse arquivo resume o objeto, seus pais, filhos, dependencias, decisoes, riscos e proximos documentos relevantes.

### Knowledge graph, not just folders

A arvore de pastas e apenas a interface humana. O modelo real deve funcionar como grafo: artefatos possuem ids, pais, filhos, dependencias, relacoes e consumidores.

### Approval gates

Agentes podem propor, mas mudancas relevantes devem passar por aprovacao explicita quando alteram escopo, arquitetura, regras de negocio, riscos, roadmap ou compromissos de entrega.

## 3. Modelo Conceitual

### Problem

Define a dor, a oportunidade e o contexto de mercado. E a raiz de justificativa do produto.

### Vision

Define o produto que queremos construir, para quem, por que agora e quais principios guiam as decisoes.

### Strategy

Define posicionamento, segmentos, metricas, trade-offs, roadmap e criterios para avancar ou pausar.

### Domain

Agrupa uma area coerente do negocio ou produto, como `users`, `groups`, `events`, `friendship` ou `payments`.

### User Goal

Substitui a nocao generica de "capability". Representa o objetivo estavel do usuario dentro de um dominio, por exemplo: "participar de um evento", "encontrar pessoas" ou "gerenciar perfil".

### Feature

E uma solucao concreta que ajuda um User Goal. Features podem entrar, sair, evoluir, ser fatiadas ou substituidas.

### Use Case

E uma interacao verificavel da feature. O Use Case e o ponto onde produto e engenharia se unem para gerar uma Specification implementavel.

### Delivery Level

Define o nivel de entrega em que um artefato deve entrar. O nivel responde "quando isso precisa existir no produto", sem substituir a prioridade dentro do nivel.

Niveis canonicos:

- `L0 Foundation`: base sem a qual o produto ou a esteira nao sustentam entregas seguras.
- `L1 Walking Skeleton`: menor fluxo ponta a ponta que prova o valor central.
- `L2 Core Loop`: ciclo principal que gera valor recorrente ao usuario.
- `L3 Trust, Safety and Quality`: confianca, seguranca, privacidade, moderacao, acessibilidade e qualidade de experiencia.
- `L4 Operations and Scale`: operacao, suporte, observabilidade, admin e escala.
- `L5 Growth and Optimization`: crescimento, experimentos, personalizacao e otimizacoes.

### Priority

Define a urgencia relativa dentro de um Delivery Level:

- `P0`: bloqueia o nivel atual.
- `P1`: necessario para considerar o nivel pronto.
- `P2`: importante, mas nao bloqueia a entrega do nivel.
- `P3`: melhoria, polish ou otimizacao.

### Specification

E a fonte de verdade para implementacao. Deve cobrir produto, fluxo, UI, APIs, dados, permissoes, analytics, seguranca, performance, acessibilidade, erros, edge cases, observabilidade e aceite.

### Design

Traduz a Specification em experiencia de usuario verificavel: fluxo visual, navegacao, wireframes, mockups, estados, acessibilidade e alinhamento com o design system. Quando a feature nao tiver interface, o artefato deve registrar explicitamente `Not applicable` e explicar por que.

### Implementation Plan

Traduz a Specification em estrategia tecnica. Pensa como um Tech Lead: sequencia, fases, dependencias, riscos, fatias, migracoes, backend, frontend, testes e rollout.

### Execution Graph

Representa as tarefas como um DAG. Cada no e uma unidade executavel com dependencias explicitas. Isso permite paralelismo seguro entre agentes e deixa claro quando algo esta bloqueado.

### Task

Unidade executavel derivada da Specification e do Execution Graph. Uma task deve ser pequena o suficiente para implementacao, teste, review e rollback.
## 4. Estrutura De Pastas

Estrutura canonica:

```text
product/
  .product/
    state.json
    decisions.json
    roadmap.json
    ids.json
    history/

  FRAMEWORK.md

  foundation/
    problem/
      problem.md
      opportunities.md
      researches/
      interviews/
      context.md
    vision/
      vision.md
      principles.md
      north-star.md
      context.md
    strategy/
      strategy.md
      personas.md
      competitors.md
      metrics.md
      roadmap.md
      context.md

  knowledge/
    glossary/
    business-rules/
    conventions/
    decisions/
    patterns/
    prompts/
    templates/
    examples/

  domains/
    <domain>/
      context.md
      domain.md
      decisions.md
      goals/
        <goal>/
          context.md
          goal.md
          journeys.md
          features/
            <feature>/
              context.md
              feature.md
              use-cases/
                <use-case>/
                  context.md
                  use-case.md
                  specification.md
                  implementation-plan.md
                  execution-graph.json
                  tasks.md
                  tests.md
                  analytics.md
                  design.md
                  audit.md

  design/
  engineering/
  audits/
  releases/
  skills/
```

## 5. Context.md

Todo `context.md` deve permitir que uma IA entenda onde esta, o que precisa ler e qual e o proximo passo seguro.

Template minimo:

```yaml
id: FT-023
type: feature
name: QR Code Check-in
status: draft
owner_skill: feature-ai

parents:
  - GOAL-003

children:
  - UC-001
  - UC-002

depends_on:
  - FT-008
  - DOMAIN-users

used_by:
  - RELEASE-001

related:
  - FT-055

documents:
  canonical: feature.md
  specification: use-cases/qr-code-check-in/specification.md
  design: use-cases/qr-code-check-in/design.md
  implementation_plan: use-cases/qr-code-check-in/implementation-plan.md
  execution_graph: use-cases/qr-code-check-in/execution-graph.json

delivery:
  level: L1
  priority: P0
  rationale: Sem check-in, o walking skeleton de evento nao fecha ponta a ponta.

open_questions:
  - Como expirar QR codes sem prejudicar usuarios offline?

decisions:
  - DEC-014
```

## 6. Specification Driven Development

O fluxo de uma nova feature deve ser:

```text
Feature -> Use Cases -> Specification -> Design -> Implementation Plan -> Execution Graph -> Tasks -> Implementation -> QA Evidence -> Security Review -> Review -> Audit -> Release
```

O Design e obrigatorio para qualquer use case com interface. Para entregas sem UI, `design.md` deve existir como artefato curto com `Not applicable`, justificativa e impactos para acessibilidade, observabilidade ou operacao quando houver.

QA Evidence e Security Review sao gates de validacao. QA Evidence comprova que os criterios de aceite, tasks, fluxos, bordas, regressao, acessibilidade, observabilidade e controles de seguranca foram verificados. Security Review avalia autenticacao, autorizacao, privacidade, abuso, dados sensiveis, tokens, logs, dependencias, rollout, rollback e risco residual. Um artefato nao deve chegar a `validated` ou `released` quando houver blocker de QA ou seguranca.

A Specification deve responder:

- O que exatamente deve acontecer?
- Qual problema do usuario isso resolve?
- O que esta dentro e fora de escopo?
- Quais fluxos, estados e erros existem?
- Quais regras de negocio se aplicam?
- Quais APIs, dados e permissoes sao necessarios?
- Quais eventos de analytics e logs devem existir?
- Quais riscos de seguranca, privacidade e abuso existem?
- Como a entrega sera testada e aceita?

Secoes obrigatorias:

```text
Product context
User goal
Delivery level
Priority
Feature scope
Non-goals
Use cases
Business rules
UX flow
UI states
API contracts
Data model
Events
Analytics
Permissions
Security
Performance
Accessibility
Error states
Edge cases
Observability
Rollout strategy
Feature flags
Acceptance criteria
Open questions
```

Security Review nao e uma promessa absoluta de ausencia de risco. O papel do gate e garantir que todos os controles definidos foram verificados com evidencia, que blockers estao resolvidos, e que riscos residuais estao documentados e aprovados por humanos quando forem relevantes.

## 6.1. Priorizacao De Entrega

Todo artefato executavel deve declarar `Delivery Level` e `Priority`. O nivel organiza o roadmap por maturidade do produto; a prioridade ordena o trabalho dentro do nivel.

Campos obrigatorios em Domain, User Goal, Feature, Use Case, Specification, Implementation Plan, Execution Graph e Task:

```yaml
delivery:
  level: L1
  priority: P0
  depends_on:
    - FT-008
  rationale: Explica por que esta entrega pertence a este nivel e por que esta prioridade foi atribuida.
```

Regras:

- `Delivery Level` nao e uma promessa de data.
- `Priority` so deve ser comparada dentro do mesmo nivel.
- Uma entrega `L3/P0` pode ser critica para confianca, mas ainda nao furar uma entrega `L1` que fecha o walking skeleton.
- Dependencias podem puxar uma entrega tecnica para um nivel anterior, desde que o `rationale` explique o motivo.
- Mudanca de `level` ou `priority` altera compromisso de entrega e deve passar por approval gate.

## 6.2. Design Driven Handoff

O Design nasce depois da Specification aprovada e antes do Implementation Plan.

Entradas:

- Specification aprovada.
- Design system, padroes de UX e telas vizinhas.
- Delivery Level e Priority da entrega.

Saidas:

- `design.md` no use case, com fluxo visual, estados, acessibilidade, componentes, dados exibidos e links para mockups quando existirem.
- Mockups ou wireframes em `product/design/` ou no diretorio de design canonico do produto, quando a entrega precisar de referencia visual.
- Revisao UX registrada antes do Implementation Plan quando a UI for relevante para o aceite.

Gates:

- Sem Specification aprovada, nao gere design.
- Sem Design aprovado ou marcado como `Not applicable`, nao gere Implementation Plan.
- Achado bloqueante de UX volta para Specification ou Design antes de seguir.

## 7. Implementation Plan

O Implementation Plan e criado depois da Specification e do Design e antes das tasks. Ele nao deve escrever codigo. Ele deve definir a estrategia de construcao.

Secoes recomendadas:

- Objetivo tecnico
- Escopo tecnico
- Delivery Level e Priority herdados ou ajustados com justificativa
- Dependencias
- Fases
- Sequencia de entrega
- Riscos
- Plano de testes
- Plano de rollback
- Arquivos ou modulos provaveis
- Decisoes que precisam de ADR
- Tasks candidatas

Exemplo de fases:

```text
1. Data model and migration
2. Server-side rules and permissions
3. Backend services and API
4. Frontend states and forms
5. Analytics and observability
6. Tests and fixtures
7. QA, review and release
```

## 8. Execution Graph

O Execution Graph e um DAG. Ele define dependencia entre tasks e permite execucao paralela por agentes.

Exemplo:

```json
{
  "id": "GRAPH-001",
  "sourceSpecification": "SPEC-001",
  "nodes": [
    {
      "id": "TK-001",
      "title": "Create event tables and policies",
      "type": "database",
      "dependsOn": []
    },
    {
      "id": "TK-002",
      "title": "Create event service",
      "type": "backend",
      "dependsOn": ["TK-001"]
    },
    {
      "id": "TK-003",
      "title": "Create event form UI",
      "type": "frontend",
      "dependsOn": ["TK-002"]
    },
    {
      "id": "TK-004",
      "title": "Instrument analytics",
      "type": "analytics",
      "dependsOn": ["TK-002", "TK-003"]
    }
  ]
}
```

Regras:

- Uma task so pode iniciar quando suas dependencias estao aprovadas.
- Tasks paralelas devem ter escopo de escrita separado.
- Toda task aponta para a Specification de origem.
- Toda mudanca de dependencia atualiza o grafo.
- Falhas de QA podem criar novos nos no grafo.
## 9. Skills

As skills sao especialistas. Elas podem operar em modos como `create`, `update`, `audit`, `evolve`, `explain`, `compare` e `refactor`, mas cada uma deve ter responsabilidade clara.

### Foundation

- Problem Discovery AI: descobre dores, oportunidades e evidencias.
- Vision AI: cria ou revisa visao, principios e north star.
- Strategy AI: define estrategia, segmentos, metricas e roadmap.
- Domain Architect AI: modela dominios e fronteiras.
- User Goal AI: modela objetivos do usuario dentro de dominios.

### Product Design

- Journey AI: mapeia jornadas.
- Feature AI: cria e evolui features.
- Use Case AI: detalha interacoes verificaveis.
- UX/UI AI: define fluxos, estados, wireframes, mockups, design system e acessibilidade.
- UX Review AI: revisa o design contra design system, principios de UX, acessibilidade e cobertura dos estados.

### Specification And Planning

- Specification AI: cria o contrato central de implementacao.
- Implementation Planner AI: transforma Specification em plano tecnico.
- Execution Graph AI: transforma plano em DAG de execucao.
- Task AI: gera tasks executaveis pequenas, testaveis e rastreaveis.

### Engineering And Validation

- Code Runner AI: implementa tasks.
- QA AI: valida comportamento, testes, bordas, performance e matriz de evidencias.
- Code Review AI: revisa qualidade, manutencao e regressao.
- Security Review AI: avalia autenticacao, autorizacao, privacidade, abuso, exposicao de dados, tokens, logs, dependencias, rollout e risco residual.

### Audit

- Gap Finder AI: procura lacunas.
- Conflict AI: procura contradicoes.
- Dependency AI: encontra dependencias implicitas.
- Impact Analysis AI: mede efeito de mudancas.
- Evolution AI: sugere melhorias.
- Documentation AI: atualiza docs.
- Product Historian AI: registra decisoes.

## 10. Orquestradores

Orquestradores nao criam artefatos primarios. Eles controlam fluxo, ordem, gates e handoffs.

### Product Orchestrator

Cria um produto do zero:

```text
Problem -> Vision -> Strategy -> Domains -> User Goals -> Roadmap
```

### New Feature Orchestrator

Recebe uma feature candidata e conduz:

```text
Impact -> Feature -> Use Cases -> Specification -> Design -> Plan -> Graph -> Tasks
```

### Audit Orchestrator

Executa auditorias em lote:

```text
Gap -> Conflict -> Dependency -> Impact -> Consistency
```

### Evolution Orchestrator

Agrupa melhorias candidatas, pergunta quais serao aprovadas e cria plano de evolucao.

### Documentation Orchestrator

Mantem `context.md`, indices, templates, decisoes e artefatos derivados sincronizados.

### Release Orchestrator

Antes de release, verifica:

- lacunas
- conflitos
- docs
- specs
- design
- tasks
- testes
- QA
- QA evidence
- review
- security review para entregas executaveis, com profundidade proporcional ao risco

## 11. Gates De Aprovacao

Cada etapa deve terminar com um estado claro:

```text
draft
proposed
approved
in_progress
implemented
validated
released
deprecated
superseded
```

Regras:

- `draft`: artefato criado, ainda incompleto.
- `proposed`: pronto para revisao humana ou de auditoria.
- `approved`: pode alimentar a proxima etapa.
- `in_progress`: esta sendo implementado.
- `implemented`: codigo ou artefato foi produzido.
- `validated`: passou por QA, review, Security Review quando aplicavel, e possui evidencias suficientes.
- `released`: chegou ao usuario ou ambiente alvo.
- `deprecated`: nao deve orientar novas implementacoes.
- `superseded`: substituido por outro artefato.

Transicoes obrigatorias:

- `proposed`: nao exige approval record, mas nao deve avancar a partir de um parent gate incompleto.
- `approved` e estados posteriores: exigem approval record correspondente em `.product/history/`, com `artifact_id`, `path`, `content_hash`, `status_granted`, `approved_by`, `approved_at` e `notes`.
- `implemented -> validated`: exige QA Evidence aprovada e sem blockers; exige Security Review aprovada quando houver codigo, dados, permissoes, tokens, API, pagamentos, uploads, mensagens, busca, admin, analytics sensivel ou qualquer risco de privacidade/abuso.
- `validated -> released`: exige Release Orchestrator, auditoria sem blockers, Security Review sem blockers, riscos residuais aceitos e rollback/monitoramento definidos.
- QA pode bloquear validacao quando qualquer criterio de aceite, task, controle de seguranca, regressao critica ou evidencia obrigatoria estiver ausente.
- Security Review pode bloquear validacao e release quando houver falha de autorizacao, vazamento de dados, decisao de permissao sem aprovacao, segredo exposto, abuso nao mitigado, logging inseguro ou risco residual alto sem decisao humana.

Approval records usam hash SHA-256 do arquivo inteiro com conteudo normalizado para LF e sem trailing whitespace por linha. Eles fornecem auditabilidade e gate mecanico, nao prova criptografica de aprovacao humana.

## 12. Decisoes

Decisoes relevantes devem ser registradas em `product/knowledge/decisions/` e indexadas em `.product/decisions.json`.

Uma decisao deve ser criada quando:

- muda arquitetura estrutural;
- muda regra de negocio importante;
- altera seguranca, privacidade, pagamento ou permissao;
- cria dependencia externa relevante;
- escolhe uma estrategia dificil de reverter;
- substitui uma decisao anterior.

## 13. Auditoria

Auditorias nao devem criar produto novo como comportamento padrao. Elas analisam e reportam.

Tipos de auditoria:

- Gap: o que falta?
- Conflict: o que contradiz outro artefato?
- Dependency: o que depende de que?
- Impact: o que muda se isso mudar?
- Consistency: os nomes, estados, ids e links batem?
- Security: ha risco de acesso indevido, abuso ou vazamento?
- UX: a experiencia fecha para a persona?

QA e Security Review devem produzir ou referenciar evidencias. Auditorias podem verificar a coerencia dessas evidencias, mas nao devem declarar uma entrega como segura quando os gates especializados estao ausentes ou bloqueados.

Saida esperada:

```text
Verdict: approved | approved_with_notes | blocked
Findings
Evidence
Required fixes
Suggested improvements
Residual risk
```

## 14. Evolution Engine

O framework deve permitir evolucao continua. Melhorias nao entram direto no produto; elas viram candidatas.

Fluxo:

```text
Observation -> Opportunity -> Proposal -> Impact Analysis -> Approval -> Updated Specification -> Updated Plan -> Updated Graph -> Tasks
```

Isso evita que sugestoes de IA virem escopo silencioso.

## 15. Como Usar Com Codex

Prompt recomendado para fase de arquitetura:

```text
Voce e um Software Architect colaborando no Product Engineering Framework.
Leia product/FRAMEWORK.md e os context.md relevantes.

Nesta etapa, nao crie arquivos e nao implemente.
Sua missao e criticar a arquitetura, encontrar ambiguidades, propor alternativas,
comparar trade-offs e perguntar o que precisa ser aprovado.

So implemente quando eu disser: CONGELAR ARQUITETURA.
```

Prompt recomendado para fase de geracao:

```text
Leia product/FRAMEWORK.md.
Use somente as decisoes aprovadas.
Nao invente novas camadas, nomes ou fluxos.
Converta a arquitetura aprovada em arquivos, templates e skills.
Preserve rastreabilidade entre Problem, Vision, Strategy, Domain, Goal, Feature,
Use Case, Specification, Implementation Plan, Execution Graph e Tasks.
```

Prompt recomendado para nova feature:

```text
Leia product/FRAMEWORK.md e o context.md do dominio.
Conduza a feature pelo fluxo:
Feature -> Use Cases -> Specification -> Design -> Implementation Plan -> Execution Graph -> Tasks.

Antes de persistir cada etapa, liste as decisoes, lacunas, conflitos e perguntas de aprovacao.
```

## 16. Roadmap Do Framework

### v0

- Estrutura de pastas.
- Templates basicos.
- `FRAMEWORK.md`.
- Lista de skills e orquestradores.

### v1

- Contextos canonicos em todos os niveis.
- IDs consistentes.
- Decision log.
- Templates completos.
- Auditorias basicas.

### v2

- Skills operacionais.
- Orquestradores com handoff.
- Execution Graph real.
- Geracao de tasks por DAG.
- Gates de aprovacao.

### v3

- Knowledge graph consultavel.
- Analise automatica de impacto.
- Execucao paralela de tasks por agentes.
- Replanejamento automatico apos falhas.

## 17. Regra Final

O framework deve ajudar agentes a pensar antes de construir.

Se uma IA nao consegue explicar de qual problema, dominio, objetivo, feature, caso de uso e Specification uma task nasceu, a task ainda nao esta pronta para implementacao.
