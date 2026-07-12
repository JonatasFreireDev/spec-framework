param(
    [string]$Version = $env:SPEC_FRAMEWORK_VERSION,
    [string]$Target = ""
)

$ErrorActionPreference = "Stop"
if (-not $Version) {
    $release = Invoke-RestMethod "https://api.github.com/repos/JonatasFreireDev/spec-framework/releases/latest"
    $Version = $release.tag_name.TrimStart("v")
}

$arch = if ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture -eq "Arm64") { "arm64" } else { "amd64" }
$base = "https://github.com/JonatasFreireDev/spec-framework/releases/download/v$Version"
$archive = "spec-framework_${Version}_windows_${arch}.zip"
$temp = Join-Path ([System.IO.Path]::GetTempPath()) "spec-framework-$Version-$arch"
$install = Join-Path $env:LOCALAPPDATA "spec-framework\bin"
New-Item -ItemType Directory -Force $temp, $install | Out-Null
Invoke-WebRequest "$base/$archive" -OutFile (Join-Path $temp $archive)
Invoke-WebRequest "$base/checksums.txt" -OutFile (Join-Path $temp "checksums.txt")

$expected = (Get-Content (Join-Path $temp "checksums.txt") | Where-Object { $_ -match [regex]::Escape($archive) }).Split()[0]
$actual = (Get-FileHash (Join-Path $temp $archive) -Algorithm SHA256).Hash.ToLowerInvariant()
if (-not $expected -or $actual -ne $expected.ToLowerInvariant()) { throw "Checksum verification failed for $archive" }

Expand-Archive (Join-Path $temp $archive) -DestinationPath $temp -Force
Copy-Item (Join-Path $temp "spec-framework.exe") (Join-Path $install "spec-framework.exe") -Force
$env:PATH = "$install;$env:PATH"
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (($userPath -split ";") -notcontains $install) {
    [Environment]::SetEnvironmentVariable("Path", (($userPath.TrimEnd(";"), $install) -join ";"), "User")
}

$arguments = @("init")
if ($Target) { $arguments += $Target }
& (Join-Path $install "spec-framework.exe") @arguments
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Spec Framework installed at $install and added to your user PATH."
