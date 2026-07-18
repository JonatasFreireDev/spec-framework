[CmdletBinding()]
param([string]$RepositoryRoot = ".", [string]$ProductRoot = "product", [switch]$Validate)

$ErrorActionPreference = "Stop"
$repo = (Resolve-Path $RepositoryRoot).Path
$manifest = Get-Content -Raw (Join-Path $repo (Join-Path $ProductRoot ".product/framework.json")) | ConvertFrom-Json
Write-Output "# Design evidence inventory"
foreach ($root in $manifest.code_roots) {
  $path = Join-Path $repo $root.path
  if (-not (Test-Path $path)) { continue }
  Get-ChildItem -LiteralPath $path -Force -File -Recurse -ErrorAction SilentlyContinue |
    Where-Object { $_.Name -in @("package.json", "tailwind.config.js", "tailwind.config.ts") -or $_.Extension -in @(".css", ".scss", ".sass", ".svg", ".fig") } |
    ForEach-Object { Write-Output ("evidence: " + $_.FullName.Substring($repo.Length + 1).Replace("\", "/")) }
}
if ($Validate) { & spec-framework design-system validate --product-root (Join-Path $repo $ProductRoot); exit $LASTEXITCODE }
