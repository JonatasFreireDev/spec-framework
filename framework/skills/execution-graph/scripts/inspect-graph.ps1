[CmdletBinding()]
param([Parameter(Mandatory)] [string]$Graph, [string]$ProductRoot = "product")
$ErrorActionPreference = "Stop"
& spec-framework graph status --graph $Graph --product-root $ProductRoot
exit $LASTEXITCODE