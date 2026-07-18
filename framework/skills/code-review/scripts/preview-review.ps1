[CmdletBinding()]
param([Parameter(Mandatory)] [string]$Work, [Parameter(Mandatory)] [string]$Stage, [string]$ProductRoot = "product")
$ErrorActionPreference = "Stop"
& spec-framework review --work $Work --stage $Stage --product-root $ProductRoot
exit $LASTEXITCODE