#!/usr/bin/env sh
set -eu

printf '%s\n' "warning: scripts/init.sh is deprecated; use scripts/install.sh. Installation no longer runs product init automatically." >&2
if [ "${1:-}" ]; then
  printf '%s\n' "warning: the legacy target argument is ignored; run 'spec-framework init $1' explicitly after installation." >&2
fi
curl -fsSL https://raw.githubusercontent.com/JonatasFreireDev/spec-framework/master/scripts/install.sh | sh
