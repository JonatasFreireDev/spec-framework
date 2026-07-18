#!/usr/bin/env bash
set -euo pipefail
repository_root="."
product_root="product"
validate=""
while [[ $# -gt 0 ]]; do case "$1" in --repository-root) repository_root="$2"; shift 2;; --product-root) product_root="$2"; shift 2;; --validate) validate="--validate"; shift;; *) echo "Usage: $0 [--repository-root <path>] [--product-root <path>] [--validate]" >&2; exit 2;; esac; done
"$(dirname "$0")/../../framework-guide/scripts/inspect-code-roots.sh" --repository-root "$repository_root" --product-root "$product_root"
printf '\n## Required coverage\nMap every discovered module, user surface, data boundary, integration, business rule, test, configuration, design asset, and operational constraint into product/knowledge/assessments/product-landscape.md.\n'
[[ -z "$validate" ]] || spec-framework validate --product-root "$repository_root/$product_root"
