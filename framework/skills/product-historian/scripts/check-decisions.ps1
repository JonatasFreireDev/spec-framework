[CmdletBinding()]
param([string]$ProductRoot = "product")
$ErrorActionPreference = "Stop"
& spec-framework decisions check --product-root $ProductRoot --json
exit $LASTEXITCODE