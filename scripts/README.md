# Polyloft Installation Scripts

This directory contains installation and update scripts for Polyloft.

## Scripts

### Installation Scripts

- **`install.sh`** - Installation script for Linux/macOS (Bash) with desktop/MIME registration for `.pf` and `.pfx`
- **`install.ps1`** - Installation script for Windows (PowerShell) that also registers `.pf`/`.pfx` associations and appends `.PFX` to `PATHEXT`

### Update Scripts

- **`update.sh`** - Update script for Linux/macOS (Bash)
- **`update.ps1`** - Update script for Windows (PowerShell) that re-applies associations and verifies `.PFX` in `PATHEXT`

## Usage

### Linux/macOS

**Install:**
```bash
curl -fsSL https://raw.githubusercontent.com/ArubikU/polyloft/main/scripts/install.sh | bash
```

**Update:**
```bash
curl -fsSL https://raw.githubusercontent.com/ArubikU/polyloft/main/scripts/update.sh | bash
```

### Windows

**Install:**
```powershell
irm https://raw.githubusercontent.com/ArubikU/polyloft/main/scripts/install.ps1 | iex
```

**Update:**
```powershell
irm https://raw.githubusercontent.com/ArubikU/polyloft/main/scripts/update.ps1 | iex
```

## What the Scripts Do

### Installation Scripts (`install.sh` / `install.ps1`)

1. Check if Go is installed (requires Go 1.22.0+)
2. Install Polyloft via `go install github.com/ArubikU/polyloft/cmd/polyloft@latest`
3. Verify the binary was installed correctly
4. Check if `GOPATH/bin` is in PATH
5. Offer to add to PATH automatically (if not already in PATH)
6. Register `.pf` and `.pfx` file associations (Windows via registry, Linux/macOS via desktop/mime entries)
7. Ensure `.PFX` is part of `PATHEXT` on Windows so binaries execute like `.exe`
8. Display success message and usage instructions

### Update Scripts (`update.sh` / `update.ps1`)

1. Check if Go is installed
2. Check if Polyloft is currently installed
3. If not installed, offer to run installation script
4. Display current version
5. Update to latest version via `go install`
6. Display new version
7. Re-apply `.pf`/`.pfx` associations (Windows) and verify `.PFX` in `PATHEXT`
8. Confirm update was successful

## Manual Installation

If you prefer not to use the scripts, you can install manually:

```bash
# Install Polyloft
go install github.com/ArubikU/polyloft/cmd/polyloft@latest

# Add to PATH (Linux/macOS)
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc

# Add to PATH (Windows PowerShell)
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$env:USERPROFILE\go\bin", [EnvironmentVariableTarget]::User)
```

## Troubleshooting

See the [main README](../README.md#troubleshooting) for installation troubleshooting.

## For Maintainers

These scripts are referenced in:
- `README.md` - User installation instructions
- `PUBLISHING.md` - Publishing and release process

When updating scripts, ensure URLs in README.md point to the correct branch (typically `main`).
