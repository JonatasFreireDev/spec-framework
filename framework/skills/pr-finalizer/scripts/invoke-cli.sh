#!/usr/bin/env bash
set -euo pipefail
# Read-only by default. Pass only the command flags required by the reviewed scope.
exec spec-framework review "$@"
