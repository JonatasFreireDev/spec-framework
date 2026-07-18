#!/usr/bin/env bash
set -euo pipefail
work=""; stage=""; product_root="product"
while [[ $# -gt 0 ]]; do case "$1" in --work) work="$2"; shift 2;; --stage) stage="$2"; shift 2;; --product-root) product_root="$2"; shift 2;; *) exit 2;; esac; done
[[ -n "$work" && -n "$stage" ]] || { echo "Usage: $0 --work <id> --stage <stage> [--product-root <path>]" >&2; exit 2; }
exec spec-framework review --work "$work" --stage "$stage" --product-root "$product_root"