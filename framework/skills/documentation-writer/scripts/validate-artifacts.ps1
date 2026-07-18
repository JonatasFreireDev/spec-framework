[CmdletBinding()]
param([string]$ProductRoot = "product", [switch]$WriteRegistry)

$ErrorActionPreference = "Stop"
$args = @("validate", "--product-root", $ProductRoot)
if ($WriteRegistry) { $args += "--write-registry" }
& spec-framework @args
exit $LASTEXITCODE
