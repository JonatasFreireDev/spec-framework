[CmdletBinding()]
param([string]$RepositoryRoot = ".", [string]$ProductRoot = "product")

$ErrorActionPreference = "Stop"
$repo = (Resolve-Path $RepositoryRoot).Path
$manifestPath = Join-Path $repo (Join-Path $ProductRoot ".product/framework.json")
$manifest = Get-Content -Raw $manifestPath | ConvertFrom-Json
Write-Output "# Technical landscape inventory"
foreach ($root in $manifest.code_roots) {
  $path = Join-Path $repo $root.path
  Write-Output "`n## code-root: $($root.path)"
  Write-Output "role: $($root.role)"
  if (-not (Test-Path -LiteralPath $path -PathType Container)) { Write-Output "coverage: missing"; continue }
  Write-Output "coverage: inspect"
  Get-ChildItem -LiteralPath $path -Force -File -Recurse -ErrorAction SilentlyContinue |
    Where-Object { $_.Name -match '^(go\.mod|package\.json|pyproject\.toml|Cargo\.toml|pom\.xml|build\.gradle|settings\.gradle|Dockerfile|docker-compose.*|.*\.ya?ml|.*\.proto|.*\.graphql|.*\.tf)$' -or $_.Extension -in @(".sln", ".csproj", ".fsproj") } |
    Sort-Object FullName |
    ForEach-Object { Write-Output ("signal: " + $_.FullName.Substring($repo.Length + 1).Replace("\", "/")) }
}
