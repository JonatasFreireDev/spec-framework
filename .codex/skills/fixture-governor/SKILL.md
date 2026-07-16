---
name: fixture-governor
description: Govern changes across Spec Framework starter assets, the Events worked-product fixture, and embedded runtime assets without mixing their ownership boundaries.
---

# Fixture Governor

## Purpose

Protect the distinction between reusable starter content and the Events validation fixture.

## Required reading

- `FRAMEWORK.md`
- `AGENTS.md`
- affected `starter/`, `examples/events/`, `framework/`, and `assets.go` files

## Workflow

1. Classify every changed artifact as framework-owned, clean starter content, generated target output, or Events product-owned fixture content.
2. Reject copying Events scope, approvals, decisions, or narratives into reusable framework or starter assets.
3. When a starter contract changes, verify initialization, upgrade preservation, embedded assets, and the matching fixture assertions.
4. When Events changes, read its relevant context and preserve its product lifecycle, approvals, and evidence.
5. Check whether empty directories can be created declaratively rather than retained through placeholder documentation.

## Output

Return the ownership map, required synchronization, preservation evidence, validation fixture impact, and unresolved boundary risks.
