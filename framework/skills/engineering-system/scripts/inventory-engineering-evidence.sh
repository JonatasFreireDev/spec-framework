#!/usr/bin/env bash
set -euo pipefail
repository_root="."; product_root="product"; validate=false
while [[ $# -gt 0 ]]; do case "$1" in --repository-root) repository_root="$2"; shift 2;; --product-root) product_root="$2"; shift 2;; --validate) validate=true; shift;; *) echo "Usage: $0 [--repository-root <path>] [--product-root <path>] [--validate]" >&2; exit 2;; esac; done
repo="$(cd "$repository_root" && pwd)"; manifest="$repo/$product_root/.product/framework.json"
[[ -f "$manifest" ]] || { echo "Spec Framework manifest not found: $manifest" >&2; exit 1; }
echo "# Engineering evidence inventory"
roots="$(awk -F '"' '/^[[:space:]]*"path"[[:space:]]*:/ { path=$4 } /^[[:space:]]*"role"[[:space:]]*:/ { if (path != "") { print path "\t" $4; path="" } }' "$manifest")"
while IFS=$'\t' read -r path role; do
  [[ -n "$path" ]] || continue; root="$repo/$path"; echo; echo "## $path ($role)"
  [[ -d "$root" ]] || { echo "MISSING"; continue; }
  while IFS= read -r file; do echo "repository-contract: ${file#"$repo/"}"; done < <(find "$root" -maxdepth 1 -type f \( -name go.mod -o -name package.json -o -name pyproject.toml -o -name Cargo.toml -o -name pom.xml -o -name build.gradle -o -name settings.gradle -o -name '*.sln' -o -name '*.csproj' -o -name '*.fsproj' \) -print)
  while IFS= read -r file; do echo "evidence: ${file#"$repo/"}"; done < <(find "$root" -type f \( -name Dockerfile -o -name 'docker-compose*' -o -name Makefile -o -name README.md -o -name go.mod -o -name package.json -o -name pyproject.toml -o -name Cargo.toml -o -name pom.xml -o -name build.gradle -o -name settings.gradle -o -name '*.sln' -o -name '*.csproj' -o -name '*.fsproj' -o -name '*.yml' -o -name '*.yaml' -o -name '*.tf' -o -iname '*test*' -o -iname '*spec*' -o -ipath '*deploy*' -o -ipath '*infra*' -o -ipath '*runbook*' \) -print)
done <<< "$roots"
if $validate; then
  spec-framework engineering-system validate --product-root "$repo/$product_root"
fi
