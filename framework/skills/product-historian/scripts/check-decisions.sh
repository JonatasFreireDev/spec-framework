#!/usr/bin/env bash
set -euo pipefail
product_root="${1:-product}"
exec spec-framework decisions check --product-root "$product_root" --json