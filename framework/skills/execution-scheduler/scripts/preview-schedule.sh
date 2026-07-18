#!/usr/bin/env bash
set -euo pipefail
exec spec-framework schedule --dry-run "$@"