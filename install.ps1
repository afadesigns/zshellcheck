<#
.SYNOPSIS
    Install or uninstall ZShellCheck on Windows.

.DESCRIPTION
    Downloads the signed pre-built ZShellCheck binary from the latest
    GitHub release, drops it into %LOCALAPPDATA%\Programs\zshellcheck,
    and adds that directory to the user's PATH. No admin rights needed.

    The script verifies the cosign signature when cosign.exe is on PATH.
    When invoked with -Uninstall, it removes the binary, the directory,
    and the PATH entry it added — no leftover files.

.PARAMETER Version
    Tag to install (default: latest). Example: -Version v1.0.14

.PARAMETER Yes
    Skip the confirmation prompt before adding to PATH.

.PARAMETER Uninstall
    Remove every artifact this script wrote: the binary, the install
    directory, and the user-PATH entry.

.PARAMETER InstallDir
    Override the install directory. Default:
    %LOCALAPPDATA%\Programs\zshellcheck.

.EXAMPLE
    iwr -useb https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.ps1 | iex

.EXAMPLE
    .\install.ps1 -Version v1.0.14 -Yes

.EXAMPLE
    .\install.ps1 -Uninstall

.NOTES
    Compatible with Windows PowerShell 5.1 (built into Windows 10+) and
    PowerShell 7+. Tested on x64, ARM64, and x86. Mirrors install.sh.
#>
[CmdletBinding()]
param(
    [string]$Version = $env:ZSHELLCHECK_VERSION,
    [switch]$Yes,
    [switch]$Uninstall,
    [string]$InstallDir
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12 -bor [Net.SecurityProtocolType]::Tls13

# ----- Defaults & helpers -------------------------------------------------

$RepoOwner = 'afadesigns'
$RepoName  = 'zshellcheck'
$BinName   = 'zshellcheck.exe'

if (-not $Version) { $Version = 'latest' }
if (-not $InstallDir) {
    $InstallDir = Join-Path $env:LOCALAPPDATA "Programs\$RepoName"
}
$BinDir    = Join-Path $InstallDir 'bin'
$BinPath   = Join-Path $BinDir $BinName

function Write-Step($msg) { Write-Host "==> $msg" -ForegroundColor Cyan }
function Write-Ok($msg)   { Write-Host "    $msg" -ForegroundColor Green }
function Write-Warn($msg) { Write-Host "    $msg" -ForegroundColor Yellow }
function Write-Err($msg)  { Write-Host "    $msg" -ForegroundColor Red }

function Get-Arch {
    switch -Wildcard ($env:PROCESSOR_ARCHITECTURE) {
        'AMD64' { return 'x86_64' }
        'ARM64' { return 'arm64' }
        'X86'   { return 'i386' }
        default { throw "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE" }
    }
}

function Resolve-Version {
    param([string]$Tag)
    if ($Tag -ne 'latest') { return $Tag }
    $api = "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest"
    Write-Ok "Resolving latest tag from $api"
    $resp = Invoke-RestMethod -Uri $api -UseBasicParsing -Headers @{ 'User-Agent' = 'install.ps1' }
    return $resp.tag_name
}

function Get-UserPath {
    return [Environment]::GetEnvironmentVariable('Path', 'User')
}

function Set-UserPath {
    param([string]$NewPath)
    [Environment]::SetEnvironmentVariable('Path', $NewPath, 'User')
}

function Add-ToUserPath {
    param([string]$DirToAdd)
    $current = Get-UserPath
    if ([string]::IsNullOrWhiteSpace($current)) { $current = '' }
    $entries = $current -split ';' | Where-Object { $_ -ne '' }
    if ($entries -contains $DirToAdd) {
        Write-Ok "PATH already contains $DirToAdd"
        return
    }
    if (-not $Yes) {
        $reply = Read-Host "Add $DirToAdd to your user PATH? [Y/n]"
        if ($reply -match '^[nN]') { Write-Warn 'Skipped PATH update.'; return }
    }
    $newPath = ($entries + $DirToAdd) -join ';'
    Set-UserPath $newPath
    $env:Path = "$env:Path;$DirToAdd"
    Write-Ok "Added $DirToAdd to user PATH (new shells will pick this up)"
}

function Remove-FromUserPath {
    param([string]$DirToRemove)
    $current = Get-UserPath
    if ([string]::IsNullOrWhiteSpace($current)) { return }
    $entries = $current -split ';' | Where-Object { $_ -ne '' -and $_ -ne $DirToRemove }
    Set-UserPath ($entries -join ';')
    Write-Ok "Removed $DirToRemove from user PATH"
}

# ----- Uninstall path -----------------------------------------------------

if ($Uninstall) {
    Write-Step "Uninstalling $RepoName"
    if (Test-Path $InstallDir) {
        Remove-Item -Recurse -Force $InstallDir
        Write-Ok "Removed $InstallDir"
    } else {
        Write-Warn "Install directory not found: $InstallDir"
    }
    Remove-FromUserPath $BinDir
    Write-Step 'Uninstall complete.'
    exit 0
}

# ----- Install path -------------------------------------------------------

$arch = Get-Arch
$tag  = Resolve-Version $Version
$ver  = $tag.TrimStart('v')

$archive   = "${RepoName}_Windows_${arch}.zip"
$urlBase   = "https://github.com/$RepoOwner/$RepoName/releases/download/$tag"
$urlZip    = "$urlBase/$archive"
$urlSig    = "$urlZip.sig"
$urlPem    = "$urlZip.pem"
$urlSums   = "$urlBase/checksums.txt"

Write-Step "Installing $RepoName $tag ($arch) to $InstallDir"

$tmp = New-Item -ItemType Directory -Path (Join-Path $env:TEMP "zshellcheck-$([guid]::NewGuid().Guid)") -Force
try {
    $zipPath  = Join-Path $tmp.FullName $archive
    $sigPath  = "$zipPath.sig"
    $pemPath  = "$zipPath.pem"
    $sumsPath = Join-Path $tmp.FullName 'checksums.txt'

    Write-Ok "Downloading $archive"
    Invoke-WebRequest -Uri $urlZip  -OutFile $zipPath  -UseBasicParsing
    Write-Ok 'Downloading checksums.txt'
    Invoke-WebRequest -Uri $urlSums -OutFile $sumsPath -UseBasicParsing

    Write-Ok 'Verifying SHA-256 checksum'
    $expected = (Select-String -Path $sumsPath -Pattern ([regex]::Escape($archive)) | Select-Object -First 1).Line
    if (-not $expected) { throw "Archive $archive not listed in checksums.txt" }
    $expectedHash = ($expected -split '\s+')[0].ToLower()
    $actualHash   = (Get-FileHash $zipPath -Algorithm SHA256).Hash.ToLower()
    if ($expectedHash -ne $actualHash) {
        throw "Checksum mismatch for $archive (expected $expectedHash, got $actualHash)"
    }
    Write-Ok "SHA-256 OK: $actualHash"

    # Cosign signature verification (best-effort if cosign is available).
    $cosign = Get-Command cosign -ErrorAction SilentlyContinue
    if ($cosign) {
        Write-Ok "cosign found: $($cosign.Source) — verifying signature"
        Invoke-WebRequest -Uri $urlSig -OutFile $sigPath -UseBasicParsing
        Invoke-WebRequest -Uri $urlPem -OutFile $pemPath -UseBasicParsing
        & cosign verify-blob `
            --certificate $pemPath `
            --signature $sigPath `
            --certificate-identity-regexp "https://github.com/$RepoOwner/$RepoName/.*" `
            --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' `
            $zipPath
        if ($LASTEXITCODE -ne 0) { throw 'cosign verify-blob failed' }
        Write-Ok 'cosign signature OK'
    } else {
        Write-Warn 'cosign not on PATH — skipping signature verification (SHA-256 still validated)'
    }

    Write-Ok "Extracting to $InstallDir"
    if (-not (Test-Path $BinDir)) { New-Item -ItemType Directory -Path $BinDir -Force | Out-Null }
    $extract = Join-Path $tmp.FullName 'extract'
    Expand-Archive -Path $zipPath -DestinationPath $extract -Force
    $exe = Get-ChildItem -Path $extract -Recurse -Filter $BinName | Select-Object -First 1
    if (-not $exe) { throw "$BinName not found in $archive" }
    Copy-Item -Path $exe.FullName -Destination $BinPath -Force
    Write-Ok "Installed: $BinPath"

    # Confirm the binary actually runs.
    $version = & $BinPath -version 2>$null
    if ($LASTEXITCODE -ne 0) { throw "Installed binary does not run cleanly: $BinPath" }
    Write-Ok "Version reports: $version"

    Add-ToUserPath $BinDir

    Write-Step 'Install complete.'
    Write-Host "    Run:    zshellcheck path\to\script.zsh"
    Write-Host "    Remove: .\install.ps1 -Uninstall"
}
finally {
    Remove-Item -Recurse -Force $tmp.FullName -ErrorAction SilentlyContinue
}
