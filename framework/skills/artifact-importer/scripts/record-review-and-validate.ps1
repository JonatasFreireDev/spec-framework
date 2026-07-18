[CmdletBinding()]
param(
  [Parameter(Mandatory)] [string]$Run,
  [Parameter(Mandatory)] [string]$Chunk,
  [Parameter(Mandatory)] [string]$Input,
  [Parameter(Mandatory)] [string]$Agent,
  [string]$ProductRoot = "product"
)

$ErrorActionPreference = "Stop"
& spec-framework import record-review --run $Run --chunk $Chunk --input $Input --agent $Agent --product-root $ProductRoot --yes
if ($LASTEXITCODE) { exit $LASTEXITCODE }
& spec-framework validate --product-root $ProductRoot --write-registry
exit $LASTEXITCODE
