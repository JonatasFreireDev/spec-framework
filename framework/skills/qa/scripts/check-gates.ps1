[CmdletBinding()]
param([string]$ProductRoot = "product")
$ErrorActionPreference = "Stop"
& spec-framework validate --product-root $ProductRoot
if ($LASTEXITCODE) { exit $LASTEXITCODE }
& spec-framework gates --product-root $ProductRoot
exit $LASTEXITCODE