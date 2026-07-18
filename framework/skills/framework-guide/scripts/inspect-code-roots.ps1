[CmdletBinding()]
param(
  [string]$RepositoryRoot = ".",
  [string]$ProductRoot = "product",
  [switch]$Validate
)

$ErrorActionPreference = "Stop"
$repo = (Resolve-Path $RepositoryRoot).Path
$manifestPath = Join-Path $repo (Join-Path $ProductRoot ".product/framework.json")
$productPrefix = (Join-Path $repo $ProductRoot).TrimEnd("\", "/") + [System.IO.Path]::DirectorySeparatorChar
$markers = @("package.json", "go.mod", "Cargo.toml", "pom.xml", "build.gradle", "build.gradle.kts", "Dockerfile")

Write-Output "# Code-root inventory"
Write-Output "Repository: $repo"
if (Test-Path $manifestPath) {
  $manifest = Get-Content -Raw $manifestPath | ConvertFrom-Json
  Write-Output "Mode: post-init declared-root evidence"
  if ($manifest.code_root_discovery) {
    Write-Output "Discovery: $($manifest.code_root_discovery.mode) / $($manifest.code_root_discovery.status)"
  }
  if (-not $manifest.code_roots -or $manifest.code_roots.Count -eq 0) {
    Write-Output "No declared code roots. Confirm the intended stack and official scaffold command before creating implementation."
  } else {
    foreach ($root in $manifest.code_roots) {
      $path = Join-Path $repo $root.path
      Write-Output "`n## $($root.path) ($($root.role))"
      if (-not (Test-Path $path)) { Write-Output "MISSING: declared root does not exist"; continue }
      Get-ChildItem -LiteralPath $path -Force -File -Recurse -ErrorAction SilentlyContinue |
        Where-Object { $_.Name -in $markers -or $_.Extension -in @(".csproj", ".sln") } |
        ForEach-Object { Write-Output ("evidence: " + $_.FullName.Substring($repo.Length + 1).Replace("\", "/")) }
    }
  }
} else {
  Write-Output "Mode: pre-init candidate discovery"
  Write-Output "Candidates are evidence only. The agent must inspect boundaries and assign semantic roles before init."
  $candidates = Get-ChildItem -LiteralPath $repo -Force -File -Recurse -ErrorAction SilentlyContinue |
    Where-Object {
      ($_.Name -in $markers -or $_.Extension -in @(".csproj", ".sln")) -and
      -not $_.FullName.StartsWith($productPrefix, [System.StringComparison]::OrdinalIgnoreCase) -and
      $_.FullName -notmatch '[\\/](\.git|node_modules|vendor)[\\/]'
    } |
    ForEach-Object {
      $root = $_.Directory.FullName.Substring($repo.Length).TrimStart("\", "/").Replace("\", "/")
      [PSCustomObject]@{ Root = $(if ($root) { $root } else { "." }); Marker = $_.Name }
    } |
    Sort-Object Root, Marker -Unique
  if (-not $candidates) {
    Write-Output "No marker candidates found. The agent must still inspect repository contents before using --no-code-roots."
  } else {
    $candidates | ForEach-Object { Write-Output "candidate: $($_.Root) (marker: $($_.Marker))" }
  }
  Write-Output "Next: classify complete roots as web, api, worker, mobile, infrastructure, library, or another explicit role; then pass --code-roots path:role,... ."
}
if ($Validate) {
  if (-not (Test-Path $manifestPath)) { throw "--Validate requires an initialized product manifest: $manifestPath" }
  & spec-framework validate --product-root (Join-Path $repo $ProductRoot)
  exit $LASTEXITCODE
}
