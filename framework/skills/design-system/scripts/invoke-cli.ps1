[CmdletBinding()]
param([Parameter(ValueFromRemainingArguments = $true)] [string[]]$Arguments)
$ErrorActionPreference = "Stop"
# Read-only by default. This wrapper never adds --yes or an approver identity.
& spec-framework design-system @Arguments
exit $LASTEXITCODE
