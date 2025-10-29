# Polyloft Installation Verification Script for Windows PowerShell

$ErrorActionPreference = "Continue"

Write-Host "======================================" -ForegroundColor Cyan
Write-Host "   Polyloft Installation Verification" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

$allOk = $true

# 1. Check if polyloft.exe exists
Write-Host "[1] Checking Polyloft binary..." -ForegroundColor Cyan
$goPathBin = "$(go env GOPATH)\bin"
$polyloftPath = "$goPathBin\polyloft.exe"

if (Test-Path $polyloftPath) {
    Write-Host "    [OK] Binary found at: $polyloftPath" -ForegroundColor Green
} else {
    Write-Host "    [ERROR] Binary NOT found at: $polyloftPath" -ForegroundColor Red
    $allOk = $false
}
Write-Host ""

# 2. Check if GOPATH\bin is in PATH
Write-Host "[2] Checking if GOPATH\bin is in PATH..." -ForegroundColor Cyan
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$machinePath = [Environment]::GetEnvironmentVariable("Path", "Machine")
$currentSessionPath = $env:Path

Write-Host "    User PATH contains GOPATH\bin: " -NoNewline
if ($userPath -like "*$goPathBin*") {
    Write-Host "[OK]" -ForegroundColor Green
} else {
    Write-Host "[MISSING]" -ForegroundColor Red
    $allOk = $false
}

Write-Host "    Current session PATH contains GOPATH\bin: " -NoNewline
if ($currentSessionPath -like "*$goPathBin*") {
    Write-Host "[OK]" -ForegroundColor Green
} else {
    Write-Host "[MISSING]" -ForegroundColor Yellow
    Write-Host "    Note: You may need to restart your terminal" -ForegroundColor Yellow
}
Write-Host ""

# 3. Try to execute polyloft
Write-Host "[3] Testing Polyloft command..." -ForegroundColor Cyan
try {
    $polyloftOutput = & polyloft version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "    [OK] Command executed successfully" -ForegroundColor Green
        Write-Host "    Output: $polyloftOutput" -ForegroundColor Gray
    } else {
        Write-Host "    [WARNING] Command executed with exit code: $LASTEXITCODE" -ForegroundColor Yellow
        Write-Host "    Output: $polyloftOutput" -ForegroundColor Gray
    }
} catch {
    Write-Host "    [ERROR] Command failed: $_" -ForegroundColor Red
    $allOk = $false
}
Write-Host ""

# 4. Check file associations
Write-Host "[4] Checking file associations..." -ForegroundColor Cyan
$classesRoot = "HKCU:\Software\Classes"

# Check .pf
$pfExtKey = Join-Path $classesRoot ".pf"
if (Test-Path $pfExtKey) {
    $pfValue = (Get-ItemProperty -Path $pfExtKey).'(default)'
    Write-Host "    [OK] .pf association: $pfValue" -ForegroundColor Green
} else {
    Write-Host "    [ERROR] .pf association NOT found" -ForegroundColor Red
    $allOk = $false
}

# Check .pfx
$pfxExtKey = Join-Path $classesRoot ".pfx"
if (Test-Path $pfxExtKey) {
    $pfxValue = (Get-ItemProperty -Path $pfxExtKey).'(default)'
    Write-Host "    [OK] .pfx association: $pfxValue" -ForegroundColor Green
} else {
    Write-Host "    [ERROR] .pfx association NOT found" -ForegroundColor Red
    $allOk = $false
}
Write-Host ""

# 5. Check PATHEXT
Write-Host "[5] Checking PATHEXT..." -ForegroundColor Cyan
$userPathext = [Environment]::GetEnvironmentVariable("PATHEXT", "User")
$machinePathext = [Environment]::GetEnvironmentVariable("PATHEXT", "Machine")
$currentPathext = $env:PATHEXT

Write-Host "    User PATHEXT contains .PFX: " -NoNewline
if ($userPathext -match "(?i)(^|;)\.PFX(;|$)") {
    Write-Host "[OK]" -ForegroundColor Green
} else {
    Write-Host "[MISSING]" -ForegroundColor Red
    $allOk = $false
}

Write-Host "    Current session PATHEXT contains .PFX: " -NoNewline
if ($currentPathext -match "(?i)(^|;)\.PFX(;|$)") {
    Write-Host "[OK]" -ForegroundColor Green
} else {
    Write-Host "[MISSING]" -ForegroundColor Yellow
    Write-Host "    Note: You may need to restart your terminal" -ForegroundColor Yellow
}
Write-Host ""

# Summary
Write-Host "======================================" -ForegroundColor Cyan
if ($allOk) {
    Write-Host "Verification Result: ALL CHECKS PASSED" -ForegroundColor Green
    Write-Host ""
    Write-Host "Polyloft is installed correctly!" -ForegroundColor Green
    Write-Host ""
    Write-Host "If commands still don't work, try:" -ForegroundColor Yellow
    Write-Host "  1. Close and reopen your terminal" -ForegroundColor White
    Write-Host "  2. Or run: " -NoNewline -ForegroundColor White
    Write-Host "`$env:Path = [System.Environment]::GetEnvironmentVariable('Path','User') + ';' + [System.Environment]::GetEnvironmentVariable('Path','Machine')" -ForegroundColor Cyan
} else {
    Write-Host "Verification Result: SOME ISSUES FOUND" -ForegroundColor Red
    Write-Host ""
    Write-Host "Recommended actions:" -ForegroundColor Yellow
    Write-Host "  1. Re-run the installation script: .\scripts\install.ps1" -ForegroundColor White
    Write-Host "  2. Restart your terminal" -ForegroundColor White
    Write-Host "  3. Run this verification script again" -ForegroundColor White
}
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

# Additional diagnostic info
Write-Host "Diagnostic Information:" -ForegroundColor Cyan
Write-Host "  GOPATH: $(go env GOPATH)" -ForegroundColor Gray
Write-Host "  GOPATH\bin: $goPathBin" -ForegroundColor Gray
Write-Host "  Polyloft path: $polyloftPath" -ForegroundColor Gray
Write-Host ""
