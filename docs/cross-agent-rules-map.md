# Cross-Agent Rules Map

## Objetivo

Identificar regras que devem ser aplicadas por todos os agentes, independentemente da skill especializada em execução. Essas regras devem ficar em um contrato comum do runtime e ser referenciadas pelas skills, evitando duplicação e divergência.

## Regra de arquitetura

As regras comuns orientam comportamento transversal. As skills continuam responsáveis pelo contrato do artefato, pelo escopo de escrita, pelos gates e pelas decisões específicas do domínio.

Fluxo recomendado:

```text
AGENTS.md local
    -> contrato comum do runtime
        -> Framework Guide / dispatcher
            -> skill especializada
```

## Candidatos a regras comuns

| Prioridade | Regra | Evidência atual | Destino sugerido | Observação |
| --- | --- | --- | --- | --- |
| P0 | Revisar a implementação em busca de blockers e gaps, corrigir os que estiverem no escopo e repetir a revisão | `AGENTS.md`, `framework/AGENTS.framework.md`, `framework/skills/code-review`, `framework/skills/gap-finder`, `framework/skills/qa` | `framework/AGENTS.framework.md` | Deve distinguir correção autorizada de aprovação formal. |
| P0 | Perguntar quando a dúvida puder alterar escopo, arquitetura, compatibilidade ou migração | `framework/skills/discovery-and-challenge.md`, `framework/AGENTS.framework.md` | Contrato comum do runtime | Usar a capacidade nativa de perguntas do harness. |
| P0 | Não inventar requisitos, decisões, dados, evidências ou aprovações | `FRAMEWORK.md`, `framework/AGENTS.framework.md`, várias skills | Contrato comum do runtime | Aprovação continua sendo humana e explícita. |
| P0 | Manter rastreabilidade por links para artifacts, documentos, seções, código e evidências | `FRAMEWORK.md`, templates, `starter/product/README.md`, `tools/check-links.py` | Contrato comum + validator/link checker | Referências identificáveis devem ser navegáveis ou reportadas como gap. |
| P0 | Produzir relatório com mudanças, decisões, validações, blockers, gaps, correções e riscos residuais | `AGENTS.md`, `framework/audits/README.md`, skills de auditoria/review | Contrato comum do runtime | Deve declarar explicitamente quando não houver gaps conhecidos. |
| P0 | Preservar escopo, histórico, aprovações e conteúdo do produto | `FRAMEWORK.md`, `framework/AGENTS.framework.md`, `starter/AGENTS.md` | Contrato comum do runtime | Não reparar approval records sem migração autorizada. |
| P1 | Avaliar impactos em testes, CI, instalação, upgrade, compatibilidade e distribuição | `AGENTS.md`, `FRAMEWORK.md`, regras de manutenção | Contrato comum para implementações | O nível de profundidade deve ser proporcional ao risco. |
| P1 | Planejar módulos plugáveis/desplugáveis com parâmetros explícitos | `AGENTS.md`, `FRAMEWORK.md` sobre módulos e adapters | Contrato comum de planejamento | A skill deve detalhar os módulos; a regra comum exige que a avaliação seja feita. |
| P1 | Definir comportamento padrão, dependências e testes das combinações modulares | `docs/artifact-registry-modules.md`, contratos de init, skills de engenharia | Contrato comum de planejamento | Evita módulos opcionais acoplados silenciosamente ao núcleo. |
| P1 | Ler a fonte de verdade, o contexto mais próximo e o contrato da skill antes de escrever | `framework/AGENTS.framework.md`, `starter/AGENTS.md`, `FRAMEWORK.md` | Contrato comum de execução | O dispatcher continua decidindo qual skill é dona da mudança. |
| P1 | Revalidar o estado depois de comandos e não assumir que a operação funcionou | `framework/skills/framework-guide`, `AGENTS.md` | Contrato comum de execução | Inclui status, diff, testes e saída do CLI. |
| P2 | Usar evidência imutável de diff/hash para implementação e revisão | `FRAMEWORK.md`, `starter/AGENTS.md`, `code-review`, `qa` | Contrato comum de entrega | Aplicável quando a mudança é executável; não deve bloquear documentação simples. |

## Regras que devem permanecer específicas

Estas regras não devem ser centralizadas integralmente porque dependem do artefato ou do domínio:

- contrato de campos de Domain, Goal, Feature, Use Case e Specification;
- gates de Design, Technical Discovery, Engineering Proposal e Engineering Review;
- análise de segurança, privacidade, abuso, permissões e ameaças;
- contratos de Design System, tokens, componentes e acessibilidade;
- geração de Tasks, DAG, leases, write scopes e evidência de código;
- adapters de aprovação e composição de artifacts por starting point;
- roteamento de falhas para Bug Fixer, QA, Code Runner, Product Historian ou Security Review;
- regras de status, aprovação e transição específicas de cada artefato.

## Forma de implementação sugerida

Usar `framework/AGENTS.framework.md` como o contrato comum do runtime, contendo somente as regras transversais. O dispatcher deve garantir que esse arquivo seja lido antes da skill especializada.

As skills não precisam repetir o texto completo. Cada uma deve apenas declarar:

1. quais regras comuns são relevantes;
2. quais regras específicas adiciona;
3. quais limites de escrita, aprovação e roteamento possui.

## Gaps identificados

1. Há regras comuns duplicadas em múltiplas skills, especialmente rastreabilidade, gaps, blockers, relatórios e perguntas.
2. O runtime possui `AGENTS.framework.md`, mas ainda não possui um contrato único e explícito para todas as regras transversais.
3. A exigência de links navegáveis existe em templates e no link checker, mas ainda não está declarada como regra operacional comum para todos os agentes.
4. A revisão pós-implementação e a avaliação de modularidade foram adicionadas ao `AGENTS.md` de manutenção, mas ainda não foram propagadas ao contrato do runtime dos repositórios adotantes.
5. O contrato comum deve ser lido pelo dispatcher como pré-requisito obrigatório, e também pelo Framework Guide antes do roteamento.

## Próxima decisão

Recomendação adotada: usar `AGENTS.framework.md` no runtime como contrato comum, torná-lo leitura obrigatória pelo Framework Guide/dispatcher e manter `AGENTS.md` como resumo local. Depois, reduzir duplicações nas skills sem remover regras específicas.

## Status da implementação

O contrato comum foi incorporado ao `framework/AGENTS.framework.md`. A ativação do runtime, os stop conditions, a autoridade de escrita, a política de links, a modularidade condicional e os formatos de relatório agora estão cobertos. Permanece apenas a redução futura de duplicidade nas skills.
