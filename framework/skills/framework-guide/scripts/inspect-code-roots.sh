#!/usr/bin/env bash
set -euo pipefail

repository_root="."
product_root="product"
validate=false
while [[ $# -gt 0 ]]; do
  case "$1" in
    --repository-root) repository_root="$2"; shift 2 ;;
    --product-root) product_root="$2"; shift 2 ;;
    --validate) validate=true; shift ;;
    *) echo "Usage: $0 [--repository-root <path>] [--product-root <path>] [--validate]" >&2; exit 2 ;;
  esac
done
repo="$(cd "$repository_root" && pwd)"
manifest="$repo/$product_root/.product/framework.json"
[[ -f "$manifest" ]] || { echo "Spec Framework manifest not found: $manifest" >&2; exit 1; }

echo "# Code-root inventory"
echo "Repository: $repo"
roots="$(awk -F '"' '/^[[:space:]]*"path"[[:space:]]*:/ { path=$4 } /^[[:space:]]*"role"[[:space:]]*:/ { if (path != "") { print path "\t" $4; path="" } }' "$manifest")"
if [[ -z "$roots" ]]; then
  echo "No declared code roots. Confirm the intended stack and official scaffold command before creating implementation."
else
  while IFS=$'\t' read -r path role; do
    root="$repo/$path"
    echo
    echo "## $path ($role)"
    [[ -d "$root" ]] || { echo "MISSING: declared root does not exist"; continue; }
    while IFS= read -r file; do echo "evidence: ${file#"$repo/"}"; done < <(find "$root" -type f \( -name package.json -o -name go.mod -o -name Cargo.toml -o -name pom.xml -o -name build.gradle -o -name build.gradle.kts -o -name Dockerfile -o -name '*.csproj' -o -name '*.sln' \) -print)
  done <<< "$roots"
fi
if $validate; then spec-framework validate --product-root "$repo/$product_root"; fi
