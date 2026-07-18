#!/usr/bin/env bash
set -euo pipefail
product_root="${1:-product}"
spec-framework validate --product-root "$product_root" --strict
exec spec-framework gates --product-root "$product_root"