#!/usr/bin/env bash
set -euo pipefail
repository_root="."; product_root="product"
while [[ $# -gt 0 ]]; do case "$1" in --repository-root) repository_root="$2"; shift 2;; --product-root) product_root="$2"; shift 2;; *) echo "Usage: $0 [--repository-root <path>] [--product-root <path>]" >&2; exit 2;; esac; done
repo="$(cd "$repository_root" && pwd)"; manifest="$repo/$product_root/.product/framework.json"
[[ -f "$manifest" ]] || { echo "Spec Framework manifest not found: $manifest" >&2; exit 1; }
echo "# Technical landscape inventory"
roots="$(awk -F '"' '/^[[:space:]]*"path"[[:space:]]*:/ { path=$4 } /^[[:space:]]*"role"[[:space:]]*:/ { if (path != "") { print path "\t" $4; path="" } }' "$manifest")"
while IFS=$'\t' read -r path role; do
  [[ -n "$path" ]] || continue; root="$repo/$path"; echo; echo "## code-root: $path"; echo "role: $role"
  [[ -d "$root" ]] || { echo "coverage: missing"; continue; }
  echo "coverage: inspect"
  find "$root" -type f \( -name go.mod -o -name package.json -o -name pyproject.toml -o -name Cargo.toml -o -name pom.xml -o -name build.gradle -o -name settings.gradle -o -name '*.sln' -o -name '*.csproj' -o -name '*.fsproj' -o -name Dockerfile -o -name 'docker-compose*' -o -name '*.yml' -o -name '*.yaml' -o -name '*.proto' -o -name '*.graphql' -o -name '*.tf' \) -print | sort | while IFS= read -r file; do echo "signal: ${file#"$repo/"}"; done
done <<< "$roots"
