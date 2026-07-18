[CmdletBinding()]
param(
  [string]$RepositoryRoot = ".",
  [string]$ProductRoot = "product",
  [switch]$Validate
)

$ErrorActionPreference = "Stop"
$repo = (Resolve-Path $RepositoryRoot).Path
$manifestPath = Join-Path $repo (Join-Path $ProductRoot ".product/framework.json")
if (-not (Test-Path $manifestPath)) { throw "Spec Framework manifest not found: $manifestPath" }
$manifest = Get-Content -Raw $manifestPath | ConvertFrom-Json

Write-Output "# Code-root inventory"
Write-Output "Repository: $repo"
if (-not $manifest.code_roots -or $manifest.code_roots.Count -eq 0) {
  Write-Output "No declared code roots. Confirm the intended stack and official scaffold command before creating implementation."
} else {
  foreach ($root in $manifest.code_roots) {
    $path = Join-Path $repo $root.path
    Write-Output "`n## $($root.path) ($($root.role))"
    if (-not (Test-Path $path)) { Write-Output "MISSING: declared root does not exist"; continue }
    Get-ChildItem -LiteralPath $path -Force -File -Recurse -ErrorAction SilentlyContinue |
      Where-Object { $_.Name -in @("package.json", "go.mod", "Cargo.toml", "pom.xml", "build.gradle", "build.gradle.kts", "Dockerfile") -or $_.Extension -in @(".csproj", ".sln") } |
      ForEach-Object { Write-Output ("evidence: " + $_.FullName.Substring($repo.Length + 1).Replace("\", "/")) }
  }
}
if ($Validate) { & spec-framework validate --product-root (Join-Path $repo $ProductRoot); exit $LASTEXITCODE }
