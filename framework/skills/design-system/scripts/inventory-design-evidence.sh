#!/usr/bin/env bash
set -euo pipefail
repository_root="."; product_root="product"; validate=false
while [[ $# -gt 0 ]]; do case "$1" in --repository-root) repository_root="$2"; shift 2;; --product-root) product_root="$2"; shift 2;; --validate) validate=true; shift;; *) echo "Usage: $0 [--repository-root <path>] [--product-root <path>] [--validate]" >&2; exit 2;; esac; done
repo="$(cd "$repository_root" && pwd)"; manifest="$repo/$product_root/.product/framework.json"
[[ -f "$manifest" ]] || { echo "Spec Framework manifest not found: $manifest" >&2; exit 1; }
echo "# Design evidence inventory"
roots="$(awk -F '"' '/^[[:space:]]*"path"[[:space:]]*:/ { path=$4 } /^[[:space:]]*"role"[[:space:]]*:/ { if (path != "") { print path "\t" $4; path="" } }' "$manifest")"
while IFS=$'\t' read -r path role; do
  [[ -n "$path" ]] || continue; root="$repo/$path"; [[ -d "$root" ]] || continue
  while IFS= read -r file; do echo "evidence: ${file#"$repo/"}"; done < <(find "$root" -type f \( -name package.json -o -name tailwind.config.js -o -name tailwind.config.ts -o -name '*.css' -o -name '*.scss' -o -name '*.sass' -o -name '*.svg' -o -name '*.fig' \) -print)
done <<< "$roots"
if $validate; then
  spec-framework design-system validate --product-root "$repo/$product_root"
fi
