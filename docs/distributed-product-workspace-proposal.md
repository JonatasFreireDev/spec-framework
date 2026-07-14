# Distributed Product Workspace

## Status

Proposta em exploração. Este documento registra a conversa sobre centralizar um produto distribuído em múltiplos repositórios e não altera, por si só, o contrato canônico do framework.

## Contexto

Um produto pode ser implementado em vários repositórios independentes, por exemplo:

- backend;
- frontend;
- BFF;
- workers;
- infraestrutura e deployment;
- contratos ou schemas compartilhados.

Cada repositório pode ter linguagem, framework, arquitetura, convenções, pipeline e configuração próprios. Ainda assim, a intenção do produto, as especificações, as dependências entre componentes e o fluxo de entrega precisam ser coordenados em um único lugar.

O framework atual parte de um modelo monorepo: documentação, código e evidências vivem no mesmo repositório. A proposta aqui é evoluir para um modelo opcional de polyrepo sem perder rastreabilidade, aprovação, execução, QA e auditoria.

## Ideia central

Criar um **Product Control Plane**, também chamado de **Distributed Product Workspace**.

```text
Product Control Plane
  ├── intenção do produto
  ├── especificações e decisões
  ├── contratos entre componentes
  ├── grafo de dependências
  ├── change sets multi-repo
  ├── ambientes e releases
  └── evidências agregadas

Repositórios de implementação
  ├── código
  ├── arquitetura local
  ├── configuração local
  ├── testes locais
  └── evidências do componente
```

A separação de autoridade seria:

| Área | Fonte de verdade |
| --- | --- |
| Intenção, escopo e comportamento do produto | Product Control Plane |
| Contratos públicos entre componentes | Product Control Plane, com cópias ou referências versionadas nos consumidores |
| Arquitetura interna do componente | Repositório do componente |
| Configuração e comandos locais | Repositório do componente |
| Execução de testes | Cada repositório, agregada pelo workspace |
| Versão efetivamente publicada | Release manifest central |
| Catálogo e navegação | Projeção opcional em Backstage ou ferramenta equivalente |

## Estrutura possível

```text
product-control-plane/
  product/
    foundation/
    domains/
    design/
    engineering/
    contracts/
      api/
      events/
      data/
    components/
    releases/
    workspaces/
    repositories.yaml
    environments.yaml

backend/
  .product/component.yaml
  .github/workflows/product-change.yml
  src/

frontend/
  .product/component.yaml
  .github/workflows/product-change.yml
  src/

bff/
  .product/component.yaml
  .github/workflows/product-change.yml
  src/
```

## Manifesto central de repositórios

O Product Control Plane teria um catálogo declarativo dos componentes:

```yaml
product: events

repositories:
  backend:
    url: github.com/acme/events-backend
    component: events-api
    default_branch: main
    adapter: go-service
    local_contract: .product/component.yaml

  frontend:
    url: github.com/acme/events-frontend
    component: events-web
    default_branch: main
    adapter: react-web
    local_contract: .product/component.yaml

  bff:
    url: github.com/acme/events-bff
    component: events-bff
    default_branch: main
    adapter: node-bff
    local_contract: .product/component.yaml
```

O manifesto local declararia como o componente implementa sua responsabilidade, sem impor sua arquitetura ao restante do produto:

```yaml
component: events-api
product: events

architecture:
  style: modular-monolith
  language: go
  framework: chi

owned:
  domains:
    - events
    - check-in

implements:
  - UC-EVENT-CHECKIN

gates:
  test: go test ./...
  lint: golangci-lint run
  build: go build ./...
```

## Change Set multi-repo

A unidade multi-repo não deve ser apenas uma task local. A proposta é introduzir um `CHANGESET-NNN`, que representa uma mudança de produto e coordena tarefas em diferentes repositórios.

```yaml
id: CHANGESET-001
feature: FT-023
title: Check-in por QR Code
status: proposed

repositories:
  - repository: backend
    base_ref: a13f9e2
    tasks:
      - create-checkin-endpoint
      - persist-checkin

  - repository: bff
    base_ref: 7b92aa1
    tasks:
      - expose-checkin-to-web

  - repository: frontend
    base_ref: 91cc4de
    tasks:
      - add-checkin-screen

contracts:
  - contracts/api/check-in.yaml
  - contracts/events/check-in-confirmed.yaml

integration:
  environment: ephemeral
  compatibility: backward-compatible
```

O fluxo seria:

```text
Feature aprovada
  -> Change Set
  -> impacto nos componentes
  -> refs-base fixadas
  -> branches e PRs por repositório
  -> gates locais
  -> ambiente integrado
  -> QA e revisão agregados
  -> Release Manifest
```

## Contratos entre componentes

O elemento compartilhado não deve ser a implementação. Deve ser o contrato:

- OpenAPI para APIs HTTP;
- AsyncAPI ou schemas versionados para eventos;
- JSON Schema, Avro ou Protobuf para mensagens;
- contratos de autenticação e autorização;
- contratos de observabilidade;
- testes de compatibilidade.

O backend pode usar uma arquitetura hexagonal, o BFF uma arquitetura orientada a composição e o frontend uma organização própria. A integração acontece pelas fronteiras aprovadas.

## Adapters por tipo de repositório

O workflow central não deve conter condicionais para todas as tecnologias. Cada componente seleciona um adapter, como:

```text
adapters/
  go-service/
  node-bff/
  react-web/
  python-worker/
  terraform/
```

Um adapter define preparação, comandos, testes, evidências, criação de branch, abertura de PR e limites de escrita. O framework opera sobre capacidades, como `api`, `ui`, `database`, `event-producer` e `event-consumer`.

## Release Manifest

Como commits em vários repositórios não são uma operação atômica, a versão do produto precisa fixar os componentes publicados:

```yaml
release: 2026.07.1
product: events

components:
  backend:
    repository: acme/events-backend
    commit: a13f9e2
    image: ghcr.io/acme/events-backend:a13f9e2
  bff:
    repository: acme/events-bff
    commit: 7b92aa1
    image: ghcr.io/acme/events-bff:7b92aa1
  frontend:
    repository: acme/events-frontend
    commit: 91cc4de
    image: ghcr.io/acme/events-frontend:91cc4de

contracts:
  api: check-in-v2
  events: check-in-events-v1
```

O sistema deve assumir compatibilidade progressiva, migrations backward-compatible, rollout em etapas e rollback por versão. Não deve prometer transação atômica entre repositórios.

## Orquestração

Em GitHub, a implementação inicial poderia usar:

- reusable workflows para compartilhar a lógica comum;
- `repository_dispatch` para disparar a execução em um repositório específico;
- uma GitHub App com permissões mínimas;
- PRs vinculadas pelo identificador do Change Set;
- CODEOWNERS e regras locais para aprovação.

O workflow central coordena. O workflow local executa os comandos e gates do componente.

## Quanto seria reutilizado

A maior parte do núcleo conceitual pode ser preservada. O principal impacto está no modelo de localização do código e na execução distribuída.

| Superfície atual | Reutilização estimada | Observação |
| --- | ---: | --- |
| Foundation, domínios, objetivos, features e use cases | 90–100% | Continuam sendo produto e não código |
| Specification e contratos modulares | 80–95% | Precisam ganhar escopo de componente e contrato externo |
| Design e Design System | 90–100% | Permanecem centralizados |
| Engineering System | 60–80% | Deve ser dividido em global e local |
| Quality System | 60–80% | Política central, comandos e evidências locais |
| Decisions, approvals e derivations | 75–90% | Precisam carregar repositório e commit afetado |
| Workspaces | 50–70% | O workspace passa a conter refs de vários repositórios |
| Execution Graph | 50–70% | O DAG precisa suportar nós por repositório |
| Tasks | 60–80% | O contrato permanece, mas cada task ganha escopo de repo |
| Code links | 20–40% | É a área com maior mudança |
| Delivery e Release | 40–60% | Precisa coordenar PRs, commits e release manifests |
| Validator | 50–70% | Ganha validação remota/local e consistência cross-repo |
| CLI/runtime | 40–60% | O runtime precisa clonar, consultar ou acionar vários repos |
| Starter e templates | 50–75% | Novos manifests e workflows locais serão necessários |
| Tests | 50–75% | Testes unitários permanecem; surgem fixtures distribuídas |

Como ordem de grandeza, assumindo que a primeira entrega seja uma versão funcional para GitHub e que não haja ainda uma UI própria:

| Fase | Esforço provável |
| --- | ---: |
| Modelo de dados, manifests e contratos | 1–2 semanas |
| Multi-repo workspace e change set | 2–4 semanas |
| Adapters e execução local | 2–4 semanas |
| PRs, dispatch, permissões e evidências | 2–3 semanas |
| Release manifest e integração | 2–4 semanas |
| Validator, testes e migração de fixtures | 2–4 semanas |
| **MVP técnico total** | **9–17 semanas** |

Para uma solução mais completa, com UI, múltiplos provedores Git, execução resiliente, retry, secrets, ambientes efêmeros, auditoria detalhada e migração compatível, a estimativa razoável seria de **4–8 meses**.

Essas estimativas são de engenharia e dependem principalmente de quantos provedores, linguagens, adapters e modos de deploy serão suportados inicialmente.

## O que mudaria muito

As mudanças estruturais mais relevantes seriam:

1. substituir o pressuposto `um produto = um repositório`;
2. criar identidade global para componente e repositório;
3. separar Engineering System global de Engineering System local;
4. permitir code links externos e versionados;
5. criar Change Set multi-repo;
6. fazer o Execution Graph atravessar repositórios;
7. agregar evidências de vários pipelines;
8. criar Release Manifest;
9. tratar integração e compatibilidade como gates próprios;
10. preservar autonomia dos repositórios de implementação.

## O que não deveria mudar

Os seguintes princípios podem continuar intactos:

- produto antes da implementação;
- Specification como contrato principal;
- decisões e aprovações explícitas;
- derivação e detecção de staleness;
- separação entre proposta, implementação, QA e revisão;
- tarefas pequenas e rastreáveis;
- auditoria baseada em evidência;
- agentes sem autoridade implícita para aprovar ou fazer merge;
- arquitetura local pertencendo ao time e ao repositório responsável.

## Estratégia recomendada

Não começar tentando transformar todo o framework de uma vez. A evolução mais segura seria:

1. documentar o modelo como proposta, como neste arquivo;
2. criar uma prova de conceito fora do caminho principal do monorepo atual;
3. suportar apenas GitHub e dois tipos de adapter;
4. criar um Product Control Plane mínimo;
5. testar um único Change Set envolvendo backend, BFF e frontend;
6. gerar um Release Manifest reproduzível;
7. comparar a experiência com o fluxo monorepo atual;
8. só então alterar `FRAMEWORK.md`, starter, CLI, validator e skills canônicas.

## Questões ainda abertas

- O Product Control Plane será um repositório separado ou uma modalidade dentro do atual product root?
- O framework deve suportar GitHub apenas no primeiro ciclo?
- A execução remota será feita por GitHub Actions, por um runner próprio ou por ambos?
- O agente poderá alterar código diretamente ou apenas abrir PRs?
- Contratos serão mantidos exclusivamente no control plane ou também versionados em pacotes publicados?
- O ambiente integrado será obrigatório para todo Change Set ou apenas para mudanças com trigger de integração?
- Como serão tratados repositórios pertencentes a equipes ou organizações diferentes?
- Qual será o mecanismo de autenticação e delegação de permissões?

## Conclusão provisória

Isso não é uma pequena extensão do fluxo atual. É uma evolução arquitetural importante, mas não exige reescrever todo o framework.

O núcleo de produto, especificação, decisão, aprovação, planejamento e auditoria é altamente reutilizável. O trabalho novo se concentra em transformar o monorepo implícito em um modelo distribuído: identidade de componentes, manifests, Change Sets, adapters, execução remota, evidências agregadas, contratos de integração e Release Manifests.

A recomendação é tratar essa mudança como uma nova capacidade opcional, com compatibilidade inicial para o modelo monorepo. Assim, o framework poderia oferecer:

```text
delivery_mode: monorepo
delivery_mode: polyrepo
```

em vez de substituir abruptamente o comportamento atual.
