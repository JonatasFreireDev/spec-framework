param(
    [string]$Version = $env:SPEC_FRAMEWORK_VERSION,
    [string]$Target = ""
)

$ErrorActionPreference = "Stop"
Write-Warning "scripts/init.ps1 is deprecated; use scripts/install.ps1. Installation no longer runs product init automatically."
if ($Target) { Write-Warning "The legacy Target argument is ignored. Run 'spec-framework init $Target' explicitly after installation." }
if ($Version) { $env:SPEC_FRAMEWORK_VERSION = $Version }
$script = Invoke-RestMethod "https://raw.githubusercontent.com/JonatasFreireDev/spec-framework/master/scripts/install.ps1"
& ([scriptblock]::Create($script))
