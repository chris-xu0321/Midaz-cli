# install.ps1 — zero-dependency installer for the midaz CLI binary (Windows).
# Installs only the binary. Skills are installed separately via:
#   midaz skills install --yes
#
# Usage:
#   irm https://raw.githubusercontent.com/SparkssL/Midaz-cli/main/install.ps1 | iex
#   .\install.ps1 -Version 0.6.0 -InstallDir "$env:LOCALAPPDATA\Programs\Midaz\bin"
param(
    [string]$Version = "",
    [string]$InstallDir = ""
)

$ErrorActionPreference = "Stop"
$Repo = "SparkssL/Midaz-cli"
$Binary = "midaz"

if (-not $InstallDir) {
    $InstallDir = Join-Path $env:LOCALAPPDATA "Programs\Midaz\bin"
}

# Detect architecture
function Get-Arch {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default { throw "Unsupported architecture: $arch" }
    }
}

# Resolve latest version
function Resolve-Version {
    if ($Version) { return $Version }
    Write-Host "Fetching latest release..."
    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
    $tag = $release.tag_name
    if (-not $tag) { throw "Could not determine latest version" }
    return $tag -replace '^v', ''
}

# Download and extract
function Install-Binary {
    param([string]$Ver, [string]$Arch)

    $archive = "$Binary-$Ver-windows-$Arch.zip"
    $url = "https://github.com/$Repo/releases/download/v$Ver/$archive"
    $checksumsUrl = "https://github.com/$Repo/releases/download/v$Ver/checksums.txt"

    $tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) ([System.Guid]::NewGuid().ToString())
    New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

    try {
        $archivePath = Join-Path $tmpDir $archive
        Write-Host "Downloading $Binary v$Ver (windows/$Arch)..."
        Invoke-WebRequest -Uri $url -OutFile $archivePath -UseBasicParsing

        # Try to verify checksum
        try {
            $checksumsPath = Join-Path $tmpDir "checksums.txt"
            Invoke-WebRequest -Uri $checksumsUrl -OutFile $checksumsPath -UseBasicParsing
            $checksums = Get-Content $checksumsPath
            $expected = ($checksums | Where-Object { $_ -match $archive }) -replace '\s+.*', ''
            if ($expected) {
                $actual = (Get-FileHash -Path $archivePath -Algorithm SHA256).Hash.ToLower()
                if ($actual -ne $expected) {
                    throw "Checksum mismatch: expected $expected, got $actual"
                }
                Write-Host "Checksum verified."
            }
        } catch [System.Net.WebException] {
            Write-Host "Warning: could not download checksums, skipping verification"
        }

        # Extract
        Expand-Archive -Path $archivePath -DestinationPath $tmpDir -Force

        # Install
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        $dest = Join-Path $InstallDir "$Binary.exe"
        Copy-Item -Path (Join-Path $tmpDir "$Binary.exe") -Destination $dest -Force
        Write-Host "Installed $Binary to $dest"
    } finally {
        Remove-Item -Path $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}

# Add to user PATH if needed
function Ensure-Path {
    $userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($userPath -split ";" | Where-Object { $_ -eq $InstallDir }) {
        return
    }
    Write-Host ""
    Write-Host "$InstallDir is not on your PATH. Adding..."
    [Environment]::SetEnvironmentVariable("PATH", "$InstallDir;$userPath", "User")
    $env:PATH = "$InstallDir;$env:PATH"
    Write-Host "Added to user PATH. Restart your terminal for it to take effect."
}

# Main
$arch = Get-Arch
$ver = Resolve-Version
Install-Binary -Ver $ver -Arch $arch
Ensure-Path

Write-Host ""
Write-Host "Done! Run 'midaz version' to verify."
Write-Host ""
Write-Host "To install agent skills (Claude Code, Codex, etc.), run:"
Write-Host "  midaz skills install --yes"
