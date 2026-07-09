# Skills v2

This folder contains the proposed Product Engineering Framework skills. Each skill owns one step of the ladder described in `product/FRAMEWORK.md`.

## Ladder

File prefixes are stable identifiers. The operational order below names the file when it differs from the numeric prefix.

1. Problem Discovery
2. Vision
3. Strategy
4. Domain Architect
5. User Goal
6. Journey
7. Feature
8. Use Case
9. Specification (`09-specification.md`)
10. UX/UI (`13-ux-ui.md`)
11. Implementation Planner (`10-implementation-planner.md`)
12. Execution Graph (`11-execution-graph.md`)
13. Task Generator (`12-task-generator.md`)
14. QA (`14-qa.md`)
15. Gap Finder
16. Conflict Finder
17. Dependency Analyzer
18. Impact Analyzer
19. Evolution
20. Documentation Writer
21. Product Historian

## Orchestrators

Orchestrators live in `orchestrators/`. They do not own canonical artifact content. They route work, enforce approval gates, and keep handoffs explicit.

## Rule

No downstream task should be generated from an unapproved Specification. If a skill cannot trace its output back to a parent artifact and context file, it should stop and report the missing link.

For use cases with UI, no Implementation Plan should be generated until `design.md` is approved. For use cases without UI, `design.md` must explicitly say `Not applicable` and explain why.

Every executable artifact must carry `Delivery Level` and `Priority` from roadmap/feature through use case, specification, design, implementation plan, execution graph, and tasks.

## Operational Skills

- .agents/skills/readiness-validator/ is the first operational Codex skill extracted from this framework. It runs the executable readiness gate and reports whether a use case can move forward.
