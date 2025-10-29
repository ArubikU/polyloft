# Polyloft Installation Script for Windows PowerShell

$ErrorActionPreference = "Stop"

function Register-PolyloftFileAssociations {
    param(
        [Parameter(Mandatory=$true)][string]$PolyloftPath
    )

    try {
        $classesRoot = "HKCU:\Software\Classes"

        # .pf source files
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

        # .pfx binaries
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

        Write-Host "[OK] Registered .pf and .pfx file associations" -ForegroundColor Green
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
Write-Host "   Polyloft Installation Script" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

# Check if Go is installed
try {
    $goVersion = go version 2>$null
    if ($LASTEXITCODE -ne 0) { throw }
} catch {
    Write-Host "Error: Go is not installed" -ForegroundColor Red
    Write-Host "Please install Go 1.22.0 or later from https://go.dev/dl/" -ForegroundColor Yellow
    exit 1
}

# Extract and check Go version
$versionMatch = $goVersion -match "go(\d+\.\d+\.\d+)"
if ($versionMatch) {
    $currentVersion = [version]$matches[1]
    $requiredVersion = [version]"1.22.0"
    
    if ($currentVersion -lt $requiredVersion) {
        Write-Host "Error: Go version $currentVersion is too old" -ForegroundColor Red
        Write-Host "Please upgrade to Go 1.22.0 or later from https://go.dev/dl/" -ForegroundColor Yellow
        exit 1
    }
}

Write-Host "[OK] Go $($matches[1]) detected" -ForegroundColor Green
Write-Host ""

# Install Polyloft
Write-Host "Installing Polyloft..." -ForegroundColor Cyan
try {
    go install github.com/ArubikU/polyloft/cmd/polyloft@latest
    if ($LASTEXITCODE -ne 0) { throw }
    Write-Host "[OK] Polyloft installed successfully" -ForegroundColor Green
} catch {
    Write-Host "[ERROR] Installation failed" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Get GOPATH/bin
$goPathBin = "$(go env GOPATH)\bin"
$polyloftPath = "$goPathBin\polyloft.exe"

# Check if binary exists
if (-not (Test-Path $polyloftPath)) {
    Write-Host "Error: Polyloft binary not found at $polyloftPath" -ForegroundColor Red
    exit 1
}

# Check if GOPATH/bin is in PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$goPathBin*") {
    Write-Host "Warning: $goPathBin is not in your PATH" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Adding to PATH..." -ForegroundColor Cyan
    
    try {
        # Add to user PATH
        $newPath = $currentPath + ";$goPathBin"
        [Environment]::SetEnvironmentVariable("Path", $newPath, [EnvironmentVariableTarget]::User)
        
        # Update current session
        $env:Path = [System.Environment]::GetEnvironmentVariable("Path","User") + ";" + [System.Environment]::GetEnvironmentVariable("Path","Machine")
        
        Write-Host "[OK] Added $goPathBin to PATH" -ForegroundColor Green
        Write-Host ""
        Write-Host "Note: You may need to restart your terminal for changes to take effect" -ForegroundColor Yellow
    } catch {
        Write-Host "[ERROR] Failed to add to PATH automatically" -ForegroundColor Red
        Write-Host ""
        Write-Host "Please add manually by running:" -ForegroundColor Yellow
        Write-Host '  [Environment]::SetEnvironmentVariable("Path", $env:Path + ";' + $goPathBin + '", [EnvironmentVariableTarget]::User)' -ForegroundColor Green
    }
} else {
    Write-Host "[OK] GOPATH\bin is in your PATH" -ForegroundColor Green
}

Register-PolyloftFileAssociations -PolyloftPath $polyloftPath
Ensure-PolyloftPathext

Write-Host ""
Write-Host "======================================" -ForegroundColor Cyan
Write-Host "Installation complete!" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Verify installation with:"
Write-Host "  polyloft version" -ForegroundColor Yellow
Write-Host ""
Write-Host "Get started with:"
Write-Host "  polyloft init          # Initialize new project" -ForegroundColor Yellow
Write-Host "  polyloft run file.pf   # Run a Polyloft file" -ForegroundColor Yellow
Write-Host "  polyloft repl          # Start interactive REPL" -ForegroundColor Yellow
Write-Host ""
Write-Host "For more information, visit:"
Write-Host "  https://github.com/ArubikU/polyloft" -ForegroundColor Cyan
Write-Host ""
