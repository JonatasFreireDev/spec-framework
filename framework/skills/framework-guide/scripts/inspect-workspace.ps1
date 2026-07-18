[CmdletBinding()]
param(
  [string]$ProductRoot = "product",
  [string]$Work = ""
)

$ErrorActionPreference = "Stop"
$args = @("--product-root", $ProductRoot)
if ($Work) { $args += @("--work", $Work) }
& spec-framework guide @args
if ($LASTEXITCODE) { exit $LASTEXITCODE }
& spec-framework status --graph @args
if ($LASTEXITCODE) { exit $LASTEXITCODE }
& spec-framework dashboard --json @args
exit $LASTEXITCODE
