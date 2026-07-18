#!/usr/bin/env bash
set -euo pipefail
[[ $# -ge 1 ]] || { echo "Usage: $0 <graph-path> [product-root]" >&2; exit 2; }
exec spec-framework graph status --graph "$1" --product-root "${2:-product}"