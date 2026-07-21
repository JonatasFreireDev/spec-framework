# Specification Contract v2 Migration

## Boundary

Migration is an explicit adopter-owned content change. `init` uses v2 for newly
generated Specifications, but `upgrade` never adds the activation field, rewrites
an existing Specification, changes product status, or creates approval evidence.

## Procedure

1. Validate the current product and retain the legacy result as migration evidence.
2. Select one bounded use case. Read its complete parent chain, decisions,
   derivations, current Specification, contracts, tests, and downstream consumers.
3. If the Specification is approved, obtain explicit human direction to revise
   it. Do not edit or repair its approval record. Material revision makes existing
   downstream derivations stale and the revised Specification returns to `draft`.
4. Add `specification_contract_version: 2` to the use-case `context.md`.
5. Convert `specification.md` to the v2 index and synthesis template. Preserve
   approved intent and links; do not summarize away requirements.
6. Convert each rigor-applicable module with the concern-specific template from
   `contracts.yaml`. Record evidence-backed `not_applicable` only where allowed.
7. Resolve deterministic validator findings, then run the Specification skill's
   adversarial audit. Keep unresolved product or architecture decisions blocking.
8. Run normal and strict product validation. Request new human approval only
   after the v2 bundle is proposed and all applicable downstream artifacts have
   been assessed for staleness.

## Compatibility And Rollback

- Removing the activation field returns validation to the legacy contract; it
  does not restore content or approval state.
- Use version control to restore adopter-owned files when abandoning a migration.
- Never delete or synthesize `.product/history/` records during rollback.
- Migrate use cases independently so a broad product can use legacy and v2
  contracts concurrently.
