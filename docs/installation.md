# Installation Guide

This guide covers all the ways to install Polyloft on your system.

## Prerequisites

- **Go 1.22.0 or later**: Required for building from source
- **Operating System**: Linux, macOS, or Windows

### Installing Go

If you don't have Go installed:

1. Visit [go.dev/dl](https://go.dev/dl/)
2. Download the installer for your operating system
3. Follow the installation instructions
4. Verify installation: `go version`

## Quick Install

### Linux/macOS

Use our installation script:

```bash
curl -fsSL https://raw.githubusercontent.com/ArubikU/polyloft/main/scripts/install.sh | bash
```

Or download and run manually:

```bash
wget https://raw.githubusercontent.com/ArubikU/polyloft/main/scripts/install.sh
chmod +x install.sh
./install.sh
```

### Windows

Use PowerShell:

```powershell
irm https://raw.githubusercontent.com/ArubikU/polyloft/main/scripts/install.ps1 | iex
```

Or download and run:

```powershell
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/ArubikU/polyloft/main/scripts/install.ps1" -OutFile "install.ps1"
.\install.ps1
```

## Manual Installation

### Using Go Install

The easiest method if you have Go installed:

```bash
go install github.com/ArubikU/polyloft/cmd/polyloft@latest
```

This installs the `polyloft` binary to `$GOPATH/bin` (usually `~/go/bin` on Unix or `%USERPROFILE%\go\bin` on Windows).

### Building from Source

Clone and build the repository:

```bash
# Clone the repository
git clone https://github.com/ArubikU/polyloft.git
cd polyloft

# Build the CLI
go build -o bin/polyloft ./cmd/polyloft

# Optional: Install globally
sudo cp bin/polyloft /usr/local/bin/  # Linux/macOS
# or copy to a directory in your PATH on Windows
```

## Configure PATH

After installation, add Polyloft to your PATH:

### Linux/macOS

Add to your shell configuration file (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Apply the changes:

```bash
source ~/.bashrc  # or ~/.zshrc
```

### Windows

Using PowerShell:

```powershell
[Environment]::SetEnvironmentVariable(
    "Path", 
    $env:Path + ";$env:USERPROFILE\go\bin", 
    [EnvironmentVariableTarget]::User
)
```

Then restart your terminal.

## Verify Installation

Check that Polyloft is installed correctly:

```bash
polyloft version
```

You should see output like:
```
v0.2.6b built 2025-10-18
```

Try running the REPL:

```bash
polyloft repl
```

You should see:
```
Polyloft REPL v0.2.6b
Type 'exit' to quit
>>>
```

## File Associations

The installation scripts automatically set up file associations:

### Linux/macOS

- `.pf` files open with `polyloft run`
- `.pfx` executables run directly
- Desktop and MIME entries are created

### Windows

- `.pf` files open with `polyloft run`
- `.pfx` files execute directly (added to `PATHEXT`)

## IDE Setup

### VSCode Extension

Install the Polyloft VSCode extension for full IDE support:

```bash
cd polyloft/vscode-extension
npm install
npm run compile
code --install-extension .
```

Generate IntelliSense mappings:

```bash
polyloft generate-mappings
```

See [VSCode Extension](vscode-extension.md) for details.

## Updating Polyloft

### Using Update Scripts

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/ArubikU/polyloft/main/scripts/update.sh | bash
```

**Windows:**
```powershell
irm https://raw.githubusercontent.com/ArubikU/polyloft/main/scripts/update.ps1 | iex
```

### Manual Update

```bash
go install github.com/ArubikU/polyloft/cmd/polyloft@latest
```

Or rebuild from source:

```bash
cd polyloft
git pull
go build -o bin/polyloft ./cmd/polyloft
```

## Troubleshooting

### Command Not Found

**Problem**: `polyloft: command not found` after installation

**Solution**:

1. Verify Go is installed:
   ```bash
   go version
   ```

2. Check if binary exists:
   ```bash
   # Linux/macOS
   ls $(go env GOPATH)/bin/polyloft
   
   # Windows
   Test-Path "$env:USERPROFILE\go\bin\polyloft.exe"
   ```

3. Ensure `GOPATH/bin` is in PATH (see Configure PATH above)

4. Restart your terminal

### Module Not Found

**Problem**: Installation fails with "module not found"

**Solution**:

1. Check internet connection
2. Clear Go module cache:
   ```bash
   go clean -modcache
   ```
3. Retry installation:
   ```bash
   go install github.com/ArubikU/polyloft/cmd/polyloft@latest
   ```

### Permission Denied

**Problem**: Permission denied when installing

**Solution**:

**Linux/macOS:**
```bash
sudo chown -R $USER $(go env GOPATH)
```

**Windows:** Run PowerShell as Administrator

### Build Errors

**Problem**: Build fails with compilation errors

**Solution**:

1. Verify Go version is 1.22.0 or later:
   ```bash
   go version
   ```

2. Update Go if needed

3. Clean and rebuild:
   ```bash
   go clean -cache
   go build -o bin/polyloft ./cmd/polyloft
   ```

### VSCode Extension Issues

**Problem**: Extension not loading or not working

**Solution**:

1. Verify Polyloft is installed and in PATH
2. Generate mappings:
   ```bash
   polyloft generate-mappings
   ```
3. Reload VSCode: `Ctrl+Shift+P` → "Developer: Reload Window"
4. Check extension logs: View → Output → Polyloft Language Server

## Platform-Specific Notes

### Linux

- Tested on Ubuntu 20.04+, Debian, Fedora, Arch
- Requires `curl` or `wget` for installation scripts
- May need `sudo` for system-wide installation

### macOS

- Tested on macOS 11 (Big Sur) and later
- Works on both Intel and Apple Silicon
- May trigger Gatekeeper; use `xattr -d com.apple.quarantine` if needed

### Windows

- Tested on Windows 10 and 11
- Requires PowerShell 5.0 or later
- May need to enable script execution:
  ```powershell
  Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
  ```

## Docker Installation

Run Polyloft in Docker:

```dockerfile
FROM golang:1.22-alpine

RUN go install github.com/ArubikU/polyloft/cmd/polyloft@latest

WORKDIR /app

CMD ["polyloft", "repl"]
```

Build and run:

```bash
docker build -t polyloft .
docker run -it polyloft
```

## Next Steps

Once installed:

1. [Getting Started](getting-started.md) - Write your first program
2. [CLI Overview](cli/overview.md) - Learn the command-line tools
3. [REPL](cli/repl.md) - Interactive experimentation

## Uninstalling

To remove Polyloft:

```bash
# Remove binary
rm $(go env GOPATH)/bin/polyloft

# Remove configuration (optional)
rm -rf ~/.polyloft
```

## Getting Help

If you encounter issues:

1. Check this troubleshooting section
2. Search [GitHub Issues](https://github.com/ArubikU/polyloft/issues)
3. Open a new issue with details about your system and the error
