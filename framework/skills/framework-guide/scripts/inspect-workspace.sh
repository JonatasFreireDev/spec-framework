#!/usr/bin/env bash
set -euo pipefail

product_root="product"
work=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --product-root) product_root="$2"; shift 2 ;;
    --work) work="$2"; shift 2 ;;
    *) echo "Usage: $0 [--product-root <path>] [--work <id>]" >&2; exit 2 ;;
  esac
done

args=(--product-root "$product_root")
[[ -n "$work" ]] && args+=(--work "$work")
spec-framework guide "${args[@]}"
spec-framework status --graph "${args[@]}"
spec-framework dashboard --json "${args[@]}"
