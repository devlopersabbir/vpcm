# Windows Code Signing & Authenticode Helper Script
# Signs Windows .exe binaries using signtool.exe or certutil so Windows Defender / SmartScreen trusts the binary.

param (
    [string]$FilePath = "app\vpsm-desktop\build\bin\vpsm-desktop.exe",
    [string]$CertPath = $env:WIN_CODE_SIGN_CERT,
    [string]$CertPassword = $env:WIN_CODE_SIGN_PASS
)

if (-not (Test-Path $FilePath)) {
    Write-Host "[!] Executable not found at '$FilePath'. Please build first." -ForegroundColor Yellow
    exit 1
}

Write-Host "[info] Checking Windows Code Signing Certificate..." -ForegroundColor Blue

if ($CertPath -and (Test-Path $CertPath)) {
    Write-Host "[info] Signing $FilePath with certificate: $CertPath" -ForegroundColor Blue
    signtool sign /f $CertPath /p $CertPassword /tr http://timestamp.digicert.com /td sha256 /fd sha256 $FilePath
    Write-Host "[✓] Successfully signed $FilePath with Authenticode signature!" -ForegroundColor Green
} else {
    Write-Host "[!] WIN_CODE_SIGN_CERT environment variable not set." -ForegroundColor Yellow
    Write-Host "[info] To establish full Windows SmartScreen trust, sign your release using a PFX code signing certificate:" -ForegroundColor Cyan
    Write-Host "       signtool sign /f mycert.pfx /p secret /tr http://timestamp.digicert.com /td sha256 /fd sha256 $FilePath" -ForegroundColor Gray
}
