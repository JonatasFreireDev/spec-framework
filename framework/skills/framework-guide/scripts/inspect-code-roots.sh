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

echo "# Code-root inventory"
echo "Repository: $repo"
if [[ -f "$manifest" ]]; then
  echo "Mode: post-init declared-root evidence"
  discovery="$(awk -F '"' '/^[[:space:]]*"mode"[[:space:]]*:/ { mode=$4 } /^[[:space:]]*"status"[[:space:]]*:/ { if (mode != "") { print mode " / " $4; exit } }' "$manifest")"
  [[ -z "$discovery" ]] || echo "Discovery: $discovery"
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
else
  echo "Mode: pre-init candidate discovery"
  echo "Candidates are evidence only. The agent must inspect boundaries and assign semantic roles before init."
  candidates="$(find "$repo" \( -path "$repo/.git" -o -path '*/node_modules' -o -path '*/vendor' -o -path "$repo/$product_root" \) -prune -o -type f \( -name package.json -o -name go.mod -o -name Cargo.toml -o -name pom.xml -o -name build.gradle -o -name build.gradle.kts -o -name Dockerfile -o -name '*.csproj' -o -name '*.sln' \) -print | while IFS= read -r file; do root="$(dirname "${file#"$repo/"}")"; [[ "$root" == "." ]] || root="${root#./}"; printf '%s\t%s\n' "$root" "$(basename "$file")"; done | sort -u)"
  if [[ -z "$candidates" ]]; then
    echo "No marker candidates found. The agent must still inspect repository contents before using --no-code-roots."
  else
    while IFS=$'\t' read -r root marker; do echo "candidate: $root (marker: $marker)"; done <<< "$candidates"
  fi
  echo "Next: classify complete roots as web, api, worker, mobile, infrastructure, library, or another explicit role; then pass --code-roots path:role,... ."
fi
if $validate; then
  [[ -f "$manifest" ]] || { echo "--validate requires an initialized product manifest: $manifest" >&2; exit 1; }
  spec-framework validate --product-root "$repo/$product_root"
fi
