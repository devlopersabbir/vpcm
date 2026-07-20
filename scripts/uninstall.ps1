# PowerShell Uninstaller Script for VPSM & VPCM on Windows
param (
    [string]$InstallDir = "$env:LocalAppData\Programs\vpsm",
    [switch]$Yes = $false,
    [switch]$NonInteractive = $false,
    [switch]$RemoveConfig = $false,
    [switch]$Help = $false
)

if ($Help) {
    Write-Host "Usage: uninstall.ps1 [options]"
    Write-Host "Options:"
    Write-Host "  -InstallDir <path>   Installation directory (default: %LOCALAPPDATA%\Programs\vpsm)"
    Write-Host "  -Yes, -NonInteractive Skip confirmation prompts"
    Write-Host "  -RemoveConfig        Remove configuration and data directory"
    Write-Host "  -Help                Show this help message"
    exit 0
}

$ErrorActionPreference = 'Continue'
$InstallDir = [System.Environment]::ExpandEnvironmentVariables($InstallDir)

# Helper functions
function Write-Info($msg) { Write-Host "[info] $msg" -ForegroundColor Blue }
function Write-Success($msg) { Write-Host "[✓] $msg" -ForegroundColor Green }
function Write-Warn($msg) { Write-Host "[!] $msg" -ForegroundColor Yellow }
function Write-ErrorMsg($msg) { Write-Host "[error] $msg" -ForegroundColor Red }

Write-Host ""
Write-Host "VPSM / VPCM Windows Uninstaller" -ForegroundColor Red
Write-Host "This script will remove VPSM binaries and PATH configurations.`n"

# Confirmation
if (-not $Yes -and -not $NonInteractive -and [Environment]::UserInteractive) {
    $confirm = Read-Host "Are you sure you want to uninstall VPSM from $InstallDir? [y/N]"
    if ($confirm -notmatch '^[Yy]$') {
        Write-Warn "Uninstall cancelled."
        exit 0
    }
}

# Remove Binaries
$Binaries = @("vpsm.exe", "vpcm.exe", "vpsmd.exe", "vpsm-api.exe")
foreach ($bin in $Binaries) {
    $binPath = Join-Path $InstallDir $bin
    if (Test-Path $binPath) {
        Remove-Item -Path $binPath -Force
        Write-Success "Removed $binPath"
    }
}

# Remove Directory if empty
if (Test-Path $InstallDir) {
    $items = Get-ChildItem -Path $InstallDir
    if (-not $items -or $items.Count -eq 0) {
        Remove-Item -Path $InstallDir -Recurse -Force
        Write-Success "Removed installation directory $InstallDir"
    } else {
        Write-Warn "Directory $InstallDir still contains other files; leaving directory."
    }
}

# Remove from User PATH
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath) {
    $CleanInstallDir = $InstallDir.TrimEnd('\', '/')
    $PathParts = $UserPath -split ';' | Where-Object { $_.TrimEnd('\', '/') -ne $CleanInstallDir }
    $NewUserPath = $PathParts -join ';'
    if ($NewUserPath -ne $UserPath) {
        [Environment]::SetEnvironmentVariable("PATH", $NewUserPath, "User")
        Write-Success "Removed $InstallDir from User PATH environment variable."
    }
}

# Config removal option
$ConfigDir = Join-Path $env:USERPROFILE ".config\vpsm"
if (Test-Path $ConfigDir) {
    if (-not $RemoveConfig -and -not $Yes -and -not $NonInteractive -and [Environment]::UserInteractive) {
        $rmConf = Read-Host "Remove config and data directory ($ConfigDir)? [y/N]"
        if ($rmConf -match '^[Yy]$') {
            $RemoveConfig = $true
        }
    }

    if ($RemoveConfig) {
        Remove-Item -Path $ConfigDir -Recurse -Force
        Write-Success "Removed config directory $ConfigDir"
    } else {
        Write-Info "Keeping config directory $ConfigDir"
    }
}

Write-Host "`n✨ VPSM has been successfully uninstalled from Windows.`n" -ForegroundColor Green
