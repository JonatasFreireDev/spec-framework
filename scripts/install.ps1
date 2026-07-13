param(
    [string]$Version = $env:SPEC_FRAMEWORK_VERSION
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

try {
    New-Item -ItemType Directory -Force $temp, $install | Out-Null
    Invoke-WebRequest "$base/$archive" -OutFile (Join-Path $temp $archive)
    Invoke-WebRequest "$base/checksums.txt" -OutFile (Join-Path $temp "checksums.txt")

    $line = Get-Content (Join-Path $temp "checksums.txt") | Where-Object { $_ -match [regex]::Escape($archive) } | Select-Object -First 1
    if (-not $line) { throw "Checksum for $archive was not published" }
    $expected = $line.Split()[0]
    $actual = (Get-FileHash (Join-Path $temp $archive) -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($actual -ne $expected.ToLowerInvariant()) { throw "Checksum verification failed for $archive" }

    Expand-Archive (Join-Path $temp $archive) -DestinationPath $temp -Force
    Copy-Item (Join-Path $temp "spec-framework.exe") (Join-Path $install "spec-framework.exe") -Force
    $env:PATH = "$install;$env:PATH"
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
	if (-not $userPath) { $userPath = "" }
    if (($userPath -split ";") -notcontains $install) {
		$newPath = (($userPath.TrimEnd(";"), $install) | Where-Object { $_ }) -join ";"
		[Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    }

    & (Join-Path $install "spec-framework.exe") version
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
	@{
		schema_version = 1
		managed_by = "spec-framework-installer"
		version = $Version
		executable = (Join-Path $install "spec-framework.exe")
		path_entry = $install
		installed_at = [DateTime]::UtcNow.ToString("o")
	} | ConvertTo-Json | Set-Content -LiteralPath (Join-Path $install "install.json") -Encoding UTF8
    Write-Host "Spec Framework installed at $install. Run 'spec-framework init' explicitly when you want to initialize a product."
}
finally {
    Remove-Item -LiteralPath $temp -Recurse -Force -ErrorAction SilentlyContinue
}
