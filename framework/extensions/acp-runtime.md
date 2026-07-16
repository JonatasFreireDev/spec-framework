# Experimental ACP Runtime Adapter

The ACP runtime adapter is disabled by default. Each invocation requires
explicit enablement and per-run acknowledgement. It claims only one ready task
for its named agent, stores a local transcript hash, and releases the temporary
lease at the end. It cannot approve, create approval records, resolve review
threads, commit, push, merge, validate, or release work.
