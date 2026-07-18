[CmdletBinding()]
param([string]$RepositoryRoot = ".", [string]$ProductRoot = "product", [switch]$Validate)

$ErrorActionPreference = "Stop"
$script = Join-Path $PSScriptRoot "../../framework-guide/scripts/inspect-code-roots.ps1"
& $script -RepositoryRoot $RepositoryRoot -ProductRoot $ProductRoot
Write-Output "`n## Required coverage"
Write-Output "Map every discovered module, user surface, data boundary, integration, business rule, test, configuration, design asset, and operational constraint into product/knowledge/assessments/product-landscape.md."
if ($Validate) { & spec-framework validate --product-root (Join-Path (Resolve-Path $RepositoryRoot) $ProductRoot); exit $LASTEXITCODE }
