#!/usr/bin/env bash
set -euo pipefail

product_root="product"; write_registry=false
while [[ $# -gt 0 ]]; do
  case "$1" in
    --product-root) product_root="$2"; shift 2 ;;
    --write-registry) write_registry=true; shift ;;
    *) echo "Usage: $0 [--product-root <path>] [--write-registry]" >&2; exit 2 ;;
  esac
done
args=(validate --product-root "$product_root")
$write_registry && args+=(--write-registry)
spec-framework "${args[@]}"
