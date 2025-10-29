# Polyloft Update Script for Windows PowerShell

$ErrorActionPreference = "Stop"

function Register-PolyloftFileAssociations {
    param(
        [Parameter(Mandatory=$true)][string]$PolyloftPath
    )

    try {
        $classesRoot = "HKCU:\Software\Classes"

        $hyExtKey = Join-Path $classesRoot ".pf"
        New-Item $hyExtKey -Force | Out-Null
        Set-ItemProperty -Path $hyExtKey -Name "(default)" -Value "Polyloft.Source"

        $hyClassKey = Join-Path $classesRoot "Polyloft.Source"
        New-Item $hyClassKey -Force | Out-Null
        Set-ItemProperty -Path $hyClassKey -Name "(default)" -Value "Polyloft Source File"
        New-Item (Join-Path $hyClassKey "DefaultIcon") -Force | Out-Null
        Set-ItemProperty -Path (Join-Path $hyClassKey "DefaultIcon") -Name "(default)" -Value "$PolyloftPath,0"
        $hyOpenCommand = Join-Path $hyClassKey "shell\open\command"
        New-Item $hyOpenCommand -Force | Out-Null
        Set-ItemProperty -Path $hyOpenCommand -Name "(default)" -Value "`"$PolyloftPath`" run `"%1`""

        $hyxExtKey = Join-Path $classesRoot ".pfx"
        New-Item $hyxExtKey -Force | Out-Null
        Set-ItemProperty -Path $hyxExtKey -Name "(default)" -Value "Polyloft.Binary"

        $hyxClassKey = Join-Path $classesRoot "Polyloft.Binary"
        New-Item $hyxClassKey -Force | Out-Null
        Set-ItemProperty -Path $hyxClassKey -Name "(default)" -Value "Polyloft Binary"
        New-Item (Join-Path $hyxClassKey "DefaultIcon") -Force | Out-Null
        Set-ItemProperty -Path (Join-Path $hyxClassKey "DefaultIcon") -Name "(default)" -Value "$PolyloftPath,0"
        $hyxOpenCommand = Join-Path $hyxClassKey "shell\open\command"
        New-Item $hyxOpenCommand -Force | Out-Null
        Set-ItemProperty -Path $hyxOpenCommand -Name "(default)" -Value "`"%1`" %*"

        Write-Host "[OK] Verified .pf and .pfx file associations" -ForegroundColor Green
    } catch {
        Write-Host "Warning: Could not register file associations automatically: $_" -ForegroundColor Yellow
    }
}

function Ensure-PolyloftPathext {
    $currentPathext = [Environment]::GetEnvironmentVariable("PATHEXT", "User")
    if ([string]::IsNullOrEmpty($currentPathext)) {
        $currentPathext = [Environment]::GetEnvironmentVariable("PATHEXT", "Machine")
    }

    if ($currentPathext -notmatch "(?i)(^|;)\.PFX(;|$)") {
        if ([string]::IsNullOrEmpty($currentPathext)) {
            $currentPathext = ".COM;.EXE;.BAT;.CMD"
        }

        if ($currentPathext[-1] -ne ';') {
            $currentPathext += ';'
        }
        $newPathext = $currentPathext + ".PFX"
        [Environment]::SetEnvironmentVariable("PATHEXT", $newPathext, [EnvironmentVariableTarget]::User)
        $env:PATHEXT = $newPathext
        Write-Host "[OK] Added .PFX to PATHEXT" -ForegroundColor Green
    } else {
        Write-Host "[OK] .PFX already present in PATHEXT" -ForegroundColor Green
    }
}

Write-Host "======================================" -ForegroundColor Cyan
Write-Host "    Polyloft Update Script" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

# Check if Go is installed
try {
    $null = go version 2>$null
    if ($LASTEXITCODE -ne 0) { throw }
} catch {
    Write-Host "Error: Go is not installed" -ForegroundColor Red
    Write-Host "Please install Go 1.22.0 or later from https://go.dev/dl/" -ForegroundColor Yellow
    exit 1
}

# Check if Polyloft is currently installed
$goPathBin = "$(go env GOPATH)\bin"
$polyloftPath = "$goPathBin\polyloft.exe"

if (-not (Test-Path $polyloftPath)) {
    Write-Host "Polyloft is not currently installed" -ForegroundColor Yellow
    $response = Read-Host "Would you like to install it now? (y/n)"
    if ($response -match '^[Yy]$') {
        Write-Host "Redirecting to installation script..." -ForegroundColor Cyan
        Invoke-Expression (Invoke-WebRequest -Uri "https://raw.githubusercontent.com/ArubikU/polyloft/main/scripts/install.ps1" -UseBasicParsing).Content
        exit 0
    } else {
        Write-Host "Update cancelled" -ForegroundColor Yellow
        exit 1
    }
}

# Get current version
try {
    $currentVersion = polyloft version 2>$null | Select-String -Pattern 'v\d+\.\d+\.\d+' | ForEach-Object { $_.Matches.Value }
    if (-not $currentVersion) { $currentVersion = "unknown" }
    Write-Host "Current version: $currentVersion" -ForegroundColor Blue
} catch {
    Write-Host "Warning: Could not determine current version" -ForegroundColor Yellow
    $currentVersion = "unknown"
}
Write-Host ""

# Update Polyloft
Write-Host "Updating Polyloft to latest version..." -ForegroundColor Cyan
try {
    go install github.com/ArubikU/polyloft/cmd/polyloft@latest
    if ($LASTEXITCODE -ne 0) { throw }
    Write-Host "[OK] Polyloft updated successfully" -ForegroundColor Green
} catch {
    Write-Host "[ERROR] Update failed" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Get new version
try {
    # Refresh PATH for current session
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","User") + ";" + [System.Environment]::GetEnvironmentVariable("Path","Machine")
    
    $newVersion = polyloft version 2>$null | Select-String -Pattern 'v\d+\.\d+\.\d+' | ForEach-Object { $_.Matches.Value }
    if (-not $newVersion) { $newVersion = "unknown" }
    Write-Host "New version: $newVersion" -ForegroundColor Green
    
    # Check if version changed
    if ($currentVersion -eq $newVersion -and $currentVersion -ne "unknown") {
        Write-Host "You already had the latest version" -ForegroundColor Yellow
    }
} catch {
    Write-Host "Warning: Could not verify new version. Try restarting your terminal." -ForegroundColor Yellow
}

Register-PolyloftFileAssociations -PolyloftPath $polyloftPath
Ensure-PolyloftPathext

Write-Host ""
Write-Host "======================================" -ForegroundColor Cyan
Write-Host "Update complete!" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Verify with:"
Write-Host "  polyloft version" -ForegroundColor Yellow
Write-Host ""
