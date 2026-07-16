# Experimental ACP Runtime Adapter

The ACP runtime adapter is disabled by default. It dispatches only one task
whose readiness, lease, worktree, and write scope were already validated by the
framework runtime. It may produce a transcript or implementation evidence; it
cannot approve, create approval records, commit, push, merge, validate, or
release work.
