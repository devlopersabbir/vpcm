# PowerShell Installer Script for VPSM & VPCM on Windows
param (
    [string]$InstallDir = "$env:LocalAppData\Programs\vpsm",
    [switch]$NonInteractive = $false,
    [switch]$Yes = $false,
    [switch]$Help = $false
)

if ($Help) {
    Write-Host "Usage: install.ps1 [options]"
    Write-Host "Options:"
    Write-Host "  -InstallDir <path>   Installation directory (default: %LOCALAPPDATA%\Programs\vpsm)"
    Write-Host "  -NonInteractive, -Yes Non-interactive mode (accept defaults)"
    Write-Host "  -Help                Show this help message"
    exit 0
}

$ErrorActionPreference = 'Stop'
$Repo = "devlopersabbir/vpcm"

# Print Banner
Write-Host ""
Write-Host "  __   _____  ____  __  __ " -ForegroundColor Cyan
Write-Host "  \ \ / /  _ \/ ___||  \/  |" -ForegroundColor Cyan
Write-Host "   \ V /| |_) \___ \| |\/| |" -ForegroundColor Cyan
Write-Host "    \ / |  __/ ___) | |  | |" -ForegroundColor Cyan
Write-Host "     V  |_|   |____/|_|  |_|" -ForegroundColor Cyan
Write-Host ""
Write-Host "Welcome to the interactive VPSM & VPCM installer for Windows!" -ForegroundColor White
Write-Host "This script will retrieve, unpack, and configure the latest release.`n"

# Helper functions
function Write-Info($msg) { Write-Host "[info] $msg" -ForegroundColor Blue }
function Write-Success($msg) { Write-Host "[✓] $msg" -ForegroundColor Green }
function Write-Warn($msg) { Write-Host "[!] $msg" -ForegroundColor Yellow }
function Write-ErrorMsg($msg) { Write-Host "[error] $msg" -ForegroundColor Red }

# Detect Architecture
$Arch = "amd64"
$ProcArch = $env:PROCESSOR_ARCHITECTURE
if ($env:PROCESSOR_ARCHITEW6432) { $ProcArch = $env:PROCESSOR_ARCHITEW6432 }

switch -regex ($ProcArch) {
    "ARM64" { $Arch = "arm64" }
    "x86|i386" { $Arch = "386" }
    "AMD64|x86_64" { $Arch = "amd64" }
    default { $Arch = "amd64" }
}

# Resolve Latest Release Tag
Write-Info "Checking latest release of VPSM..."
$LatestRelease = $null

try {
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    $req = [System.Net.WebRequest]::Create("https://github.com/$Repo/releases/latest")
    $req.AllowAutoRedirect = $false
    $res = $req.GetResponse()
    $loc = $res.GetResponseHeader("Location")
    $res.Close()
    if ($loc -and $loc.Contains("/tag/")) {
        $LatestRelease = $loc.Substring($loc.LastIndexOf('/') + 1)
    }
} catch {
    # Fallback to GitHub API
}

if (-not $LatestRelease) {
    try {
        $apiRes = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing
        if ($apiRes -and $apiRes.tag_name) {
            $LatestRelease = $apiRes.tag_name
        }
    } catch {
        # Silent fallback failure
    }
}

if (-not $LatestRelease) {
    Write-ErrorMsg "Could not retrieve latest release version. Please check your network connection."
    exit 1
}

Write-Success "Latest release found: $LatestRelease"

# Interactive Path Selection
if (-not $NonInteractive -and -not $Yes -and [Environment]::UserInteractive) {
    Write-Host "`nStep 1: Choose Installation Directory" -ForegroundColor White
    $inputPath = Read-Host "Install path [default: $InstallDir]"
    if ($inputPath -and $inputPath.Trim() -ne "") {
        $InstallDir = $inputPath.Trim()
    }
} else {
    Write-Info "Non-interactive installation: using directory $InstallDir"
}

# Resolve Environment Variables in InstallDir if present (e.g. %USERPROFILE%)
$InstallDir = [System.Environment]::ExpandEnvironmentVariables($InstallDir)

$Filename = "vpsm-windows-$Arch.zip"
$DownloadUrl = "https://github.com/$Repo/releases/download/$LatestRelease/$Filename"

# Create Temporary Working Directory
$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("vpsm-install-" + [Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $TempDir -Force | Out-Null

try {
    $ZipPath = Join-Path $TempDir $Filename
    Write-Info "Downloading $Filename from GitHub..."
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $ZipPath -UseBasicParsing

    Write-Info "Extracting binaries..."
    Expand-Archive -Path $ZipPath -DestinationPath $TempDir -Force

    if (-not (Test-Path $InstallDir)) {
        Write-Info "Creating installation directory $InstallDir..."
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    $Binaries = @("vpsm.exe", "vpsmd.exe", "vpsm-api.exe")
    foreach ($bin in $Binaries) {
        $src = Join-Path $TempDir $bin
        if (-not (Test-Path $src)) {
            $found = Get-ChildItem -Path $TempDir -Recurse -Filter $bin -File | Select-Object -First 1
            if ($found) { $src = $found.FullName }
        }

        if ($src -and (Test-Path $src)) {
            $dest = Join-Path $InstallDir $bin
            Copy-Item -Path $src -Destination $dest -Force
            Write-Success "Installed $bin to $InstallDir"
        } else {
            Write-ErrorMsg "Failed to locate binary $bin in release archive."
            exit 1
        }
    }

    # Copy vpsm.exe to vpcm.exe
    $vpsmPath = Join-Path $InstallDir "vpsm.exe"
    $vpcmPath = Join-Path $InstallDir "vpcm.exe"
    Copy-Item -Path $vpsmPath -Destination $vpcmPath -Force
    Write-Success "Linked vpcm.exe to vpsm.exe"

    # Add to User PATH if missing
    $UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if (-not $UserPath) { $UserPath = "" }

    $CleanInstallDir = $InstallDir.TrimEnd('\', '/')
    $PathList = $UserPath -split ';' | ForEach-Object { $_.TrimEnd('\', '/') }

    if ($PathList -notcontains $CleanInstallDir) {
        $NewUserPath = if ([string]::IsNullOrWhiteSpace($UserPath)) { $InstallDir } else { "$UserPath;$InstallDir" }
        [Environment]::SetEnvironmentVariable("PATH", $NewUserPath, "User")
        $env:PATH = "$env:PATH;$InstallDir"
        Write-Success "Added $InstallDir to User PATH environment variable."
    }

    Write-Host "`n✨ VPSM (vpsm/vpcm) has been successfully installed to $InstallDir!" -ForegroundColor Green
    Write-Host "Run 'vpsm version' or 'vpcm version' in a new terminal window to verify your installation.`n" -ForegroundColor Cyan
} finally {
    if (Test-Path $TempDir) {
        Remove-Item -Path $TempDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}
