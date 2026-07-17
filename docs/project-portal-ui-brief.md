# Project Portal UI Brief

Create a clean, accessible, product-facing portal for reviewing a Spec Framework project. The interface is for people making product decisions; do not expose the worktree as a developer file browser by default.

## Primary journey

Start with pending decisions, navigate through the documentation hierarchy, open an artifact, read its content and context, then approve, reject with notes, or request a revision. Keep information progressive: summary first, details only when requested.

## Required screens

- Overview with decision-focused counts and priority blockers.
- Documentation explorer with hierarchy, search and filters.
- Artifact detail with tabs: Content, Details, Relations and History.
- Batch review with selected artifacts, blockers and confirmation.
- Worktree changes view for registered, changed, untracked, missing and stale files.

## Visual direction

Calm, professional, high-legibility interface inspired by Linear, GitHub and Notion. Use restrained, contextual color; every state must also have text and an icon. Prioritize keyboard navigation, visible focus, semantic landmarks and screen-reader announcements. The layout must adapt from desktop to mobile.

## Local API contract

- `GET /api/project-view`: complete view model for initial load.
- `GET /api/project-view/changes?since=<revision>&wait=<0-25>`: long-poll for worktree or Git changes.
- `POST /api/transition`: individual lifecycle decision.
- `POST /api/batch-approval-plan`: preflight for batch approval.
- `POST /api/batch-approve`: confirmed batch approval.
- `POST /api/batch-reject`: confirmed atomic batch request for changes.

`/api/project-view` returns artifacts, worktree files, Git state, metrics and the type configuration used to choose renderers and tabs. Do not invent responsible people, priorities, tags or due dates: those fields do not exist in this product.

## Artifact-specific rendering

- Markdown is the default content renderer.
- Execution graphs use a dependency graph renderer.
- Design systems expose overview, tokens, themes and components.
- Design artifacts expose maturity, sources, screens and mappings.
- Tasks expose dependencies, acceptance criteria, write scope and evidence.
- Decisions expose scope, affected artifacts and workflow effects.
- QA, security review and audit artifacts expose findings, evidence and verdict.

## State vocabulary

Lifecycle states: draft, proposed, approved, in progress, implemented, validated, released, rejected, deprecated and superseded. Derived states: current, stale, blocked and missing. Explain derived state and blockers in plain language.
