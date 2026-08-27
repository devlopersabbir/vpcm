# PowerShell Uninstaller Script for VPSM & VPCM on Windows
param (
    [string]$InstallDir = "$env:LocalAppData\Programs\vpsm",
    [switch]$Yes = $false,
    [switch]$NonInteractive = $false,
    [switch]$PurgeDb = $false,
    [switch]$Help = $false
)

if ($Help) {
    Write-Host "Usage: uninstall.ps1 [options]"
    Write-Host "Options:"
    Write-Host "  -InstallDir <path>   Primary installation directory (default: %LOCALAPPDATA%\Programs\vpsm)"
    Write-Host "  -Yes, -NonInteractive Skip confirmation prompts and perform complete uninstallation"
    Write-Host "  -PurgeDb             Purge the SQLite database as well (by default it is preserved)"
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
Write-Host "Complete VPSM / VPCM Windows Uninstaller" -ForegroundColor Red
Write-Host "This script cleanly removes all VPSM binaries, background daemons, configurations,`nand PATH entries while safely preserving your SQLite database.`n"

# Confirmation
if (-not $Yes -and -not $NonInteractive -and [Environment]::UserInteractive) {
    Write-Host "Notice: This will remove VPSM/VPCM binaries, running daemons, configs, and shortcuts." -ForegroundColor Yellow
    Write-Host "Note: Your SQLite database will NOT be removed unless -PurgeDb is specified.`n" -ForegroundColor Green
    $confirm = Read-Host "Are you sure you want to proceed with uninstallation? [y/N]"
    if ($confirm -notmatch '^[Yy]$') {
        Write-Warn "Uninstall cancelled."
        exit 0
    }
}

# 1. Stop background processes and daemons
Write-Info "Stopping running background processes and daemons..."
$ProcessesToStop = @("vpsm-api", "vpsmd", "vpsm-desktop", "vpsm", "vpcm")
foreach ($proc in $ProcessesToStop) {
    Stop-Process -Name $proc -Force -ErrorAction SilentlyContinue
}
Write-Success "Terminated active VPSM processes and background daemons."

# 2. Remove Binaries
Write-Info "Scanning and removing binaries from installation paths..."
$SearchDirs = @(
    $InstallDir,
    "$env:LocalAppData\Programs\vpsm",
    "$env:USERPROFILE\bin",
    "$env:USERPROFILE\go\bin",
    "$env:ProgramFiles\vpsm"
)

$Binaries = @("vpsm.exe", "vpcm.exe", "vpsmd.exe", "vpsm-api.exe", "vpsm-desktop.exe")

foreach ($dir in $SearchDirs) {
    if (-not (Test-Path $dir)) { continue }
    foreach ($bin in $Binaries) {
        $binPath = Join-Path $dir $bin
        if (Test-Path $binPath) {
            Remove-Item -Path $binPath -Force -ErrorAction SilentlyContinue
            Write-Success "Removed $binPath"
        }
    }
    # If folder is now empty and not a root user folder, remove it
    if ($dir -match "vpsm$" -and (Test-Path $dir)) {
        $items = Get-ChildItem -Path $dir -ErrorAction SilentlyContinue
        if (-not $items -or $items.Count -eq 0) {
            Remove-Item -Path $dir -Recurse -Force -ErrorAction SilentlyContinue
            Write-Success "Removed installation directory $dir"
        }
    }
}

# 3. Remove from User PATH
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath) {
    $CleanInstallDir = $InstallDir.TrimEnd('\', '/')
    $PathParts = $UserPath -split ';' | Where-Object { $_.TrimEnd('\', '/') -ne $CleanInstallDir -and $_ -notmatch '(?i)[\\/]Programs[\\/]vpsm' }
    $NewUserPath = $PathParts -join ';'
    if ($NewUserPath -ne $UserPath) {
        [Environment]::SetEnvironmentVariable("PATH", $NewUserPath, "User")
        Write-Success "Removed VPSM from User PATH environment variable."
    }
}

# 4. Clean PowerShell Profiles (completions, aliases)
$ProfilePaths = @(
    $PROFILE,
    $PROFILE.CurrentUserCurrentHost,
    $PROFILE.CurrentUserAllHosts,
    $PROFILE.AllUsersCurrentHost,
    $PROFILE.AllUsersAllHosts
) | Where-Object { $_ -and (Test-Path $_) } | Select-Object -Unique

foreach ($prof in $ProfilePaths) {
    $content = Get-Content -Path $prof -Raw -ErrorAction SilentlyContinue
    if ($content -and ($content -match 'vpsm' -or $content -match 'vpcm')) {
        $lines = Get-Content -Path $prof | Where-Object { $_ -notmatch '(?i)(vpsm|vpcm)' }
        Set-Content -Path $prof -Value $lines
        Write-Success "Cleaned profile: $prof"
    }
}

# 5. Remove Configuration & Log Directories
Write-Info "Removing configuration directories and logs..."
$ConfigDirs = @(
    (Join-Path $env:USERPROFILE ".config\vpsm"),
    (Join-Path $env:USERPROFILE ".config\vpcm"),
    (Join-Path $env:APPDATA "vpsm"),
    (Join-Path $env:LOCALAPPDATA "vpsm")
)

foreach ($cdir in $ConfigDirs) {
    if (Test-Path $cdir) {
        Remove-Item -Path $cdir -Recurse -Force -ErrorAction SilentlyContinue
        Write-Success "Removed config directory $cdir"
    }
}

# 6. Remove Desktop & Start Menu Shortcuts
$ShortcutPaths = @(
    (Join-Path $env:APPDATA "Microsoft\Windows\Start Menu\Programs\VPSM*.lnk"),
    (Join-Path $env:USERPROFILE "Desktop\VPSM*.lnk"),
    (Join-Path $env:APPDATA "Microsoft\Windows\Start Menu\Programs\VPCM*.lnk"),
    (Join-Path $env:USERPROFILE "Desktop\VPCM*.lnk")
)
foreach ($sc in $ShortcutPaths) {
    if (Test-Path $sc) {
        Remove-Item -Path $sc -Force -ErrorAction SilentlyContinue
        Write-Success "Removed shortcut: $sc"
    }
}

# 7. SQLite database preservation / handling
$DefaultDbDir = Join-Path $env:USERPROFILE ".local\share\vpsm"
$DefaultDbPath = Join-Path $DefaultDbDir "vpsm.db"

if ($PurgeDb) {
    if (Test-Path $DefaultDbDir) {
        Remove-Item -Path $DefaultDbDir -Recurse -Force -ErrorAction SilentlyContinue
        Write-Warn "Purged SQLite database directory at: $DefaultDbDir"
    }
} else {
    if (Test-Path $DefaultDbPath) {
        Write-Success "Preserved SQLite database at: $DefaultDbPath"
        Write-Info "  (Your server inventory and data are safe and intact)"
    } else {
        Write-Info "No local SQLite database found at $DefaultDbPath."
    }
}

Write-Host "`n✨ VPSM / VPCM has been completely uninstalled from everywhere on Windows.`n" -ForegroundColor Green
