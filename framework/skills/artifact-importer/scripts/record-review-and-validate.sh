#!/usr/bin/env bash
set -euo pipefail

run=""; chunk=""; input=""; agent=""; product_root="product"
while [[ $# -gt 0 ]]; do
  case "$1" in
    --run) run="$2"; shift 2 ;;
    --chunk) chunk="$2"; shift 2 ;;
    --input) input="$2"; shift 2 ;;
    --agent) agent="$2"; shift 2 ;;
    --product-root) product_root="$2"; shift 2 ;;
    *) echo "Usage: $0 --run <id> --chunk <id> --input <json> --agent <id> [--product-root <path>]" >&2; exit 2 ;;
  esac
done
[[ -n "$run" && -n "$chunk" && -n "$input" && -n "$agent" ]] || { echo "run, chunk, input, and agent are required" >&2; exit 2; }

spec-framework import record-review --run "$run" --chunk "$chunk" --input "$input" --agent "$agent" --product-root "$product_root" --yes
spec-framework validate --product-root "$product_root" --write-registry
