[CmdletBinding()]
param([string]$RepositoryRoot = ".", [string]$ProductRoot = "product", [switch]$Validate)

$ErrorActionPreference = "Stop"
$repo = (Resolve-Path $RepositoryRoot).Path
$manifest = Get-Content -Raw (Join-Path $repo (Join-Path $ProductRoot ".product/framework.json")) | ConvertFrom-Json
Write-Output "# Engineering evidence inventory"
foreach ($root in $manifest.code_roots) {
  $path = Join-Path $repo $root.path
  Write-Output "`n## $($root.path) ($($root.role))"
  if (-not (Test-Path $path)) { Write-Output "MISSING"; continue }
  Get-ChildItem -LiteralPath $path -Force -File -Recurse -ErrorAction SilentlyContinue |
    Where-Object { $_.Name -match '^(Dockerfile|docker-compose.*|Makefile|README\.md|.*\.ya?ml)$' -or $_.FullName -match '(test|spec|\.github|\.gitlab)' } |
    ForEach-Object { Write-Output ("evidence: " + $_.FullName.Substring($repo.Length + 1).Replace("\", "/")) }
}
if ($Validate) { & spec-framework engineering-system validate --product-root (Join-Path $repo $ProductRoot); exit $LASTEXITCODE }
