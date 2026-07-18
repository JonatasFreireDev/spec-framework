[CmdletBinding()]
param([Parameter(ValueFromRemainingArguments = $true)] [string[]]$Arguments)
$ErrorActionPreference = "Stop"
& spec-framework schedule --dry-run @Arguments
exit $LASTEXITCODE