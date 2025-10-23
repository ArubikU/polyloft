#!/bin/bash
# Polyloft Installation Script for Linux/macOS

set -e

echo "======================================"
echo "   Polyloft Installation Script"
echo "======================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

register_file_associations() {
    local app_dir="$HOME/.local/share/applications"
    local mime_dir="$HOME/.local/share/mime/packages"
    local source_desktop="$app_dir/polyloft-source.desktop"
    local binary_desktop="$app_dir/polyloft-binary.desktop"
    local mime_file="$mime_dir/polyloft.xml"

    mkdir -p "$app_dir" "$mime_dir"

    cat > "$source_desktop" <<EOF
[Desktop Entry]
Type=Application
Name=Polyloft Source Runner
Comment=Execute Polyloft source files with the Polyloft CLI
Exec="$GOPATH_BIN/polyloft" run %f
Terminal=true
Categories=Development;
MimeType=text/x-polyloft;
EOF

    cat > "$binary_desktop" <<'EOF'
[Desktop Entry]
Type=Application
Name=Polyloft Binary
Comment=Run Polyloft compiled binaries
Exec=%f
Terminal=true
Categories=Development;
MimeType=application/x-polyloft-binary;
EOF

    cat > "$mime_file" <<'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<mime-info xmlns="http://www.freedesktop.org/standards/shared-mime-info">
    <mime-type type="text/x-polyloft">
        <comment>Polyloft source file</comment>
        <glob pattern="*.hy"/>
    </mime-type>
    <mime-type type="application/x-polyloft-binary">
        <comment>Polyloft binary</comment>
        <glob pattern="*.hyx"/>
    </mime-type>
</mime-info>
EOF

        if command -v update-desktop-database >/dev/null 2>&1; then
                update-desktop-database "$app_dir" >/dev/null 2>&1 || true
        fi

        if command -v update-mime-database >/dev/null 2>&1; then
                update-mime-database "$(dirname "$mime_dir")" >/dev/null 2>&1 || true
        fi

        echo -e "${GREEN}[OK] Registered .pf and .pfx file associations${NC}"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go 1.22.0 or later from https://go.dev/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.22.0"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo -e "${RED}Error: Go version $GO_VERSION is too old${NC}"
    echo "Please upgrade to Go 1.22.0 or later from https://go.dev/dl/"
    exit 1
fi

echo -e "${GREEN}[OK] Go $GO_VERSION detected${NC}"
echo ""

# Install Polyloft
echo "Installing Polyloft..."
if go install github.com/ArubikU/polyloft/cmd/polyloft@latest; then
    echo -e "${GREEN}[OK] Polyloft installed successfully${NC}"
else
    echo -e "${RED}[ERROR] Installation failed${NC}"
    exit 1
fi
echo ""

# Get GOPATH/bin
GOPATH_BIN=$(go env GOPATH)/bin
POLYLOFT_PATH="$GOPATH_BIN/polyloft"

# Check if binary exists
if [ ! -f "$POLYLOFT_PATH" ]; then
    echo -e "${RED}Error: Polyloft binary not found at $POLYLOFT_PATH${NC}"
    exit 1
fi

# Check if GOPATH/bin is in PATH
if [[ ":$PATH:" != *":$GOPATH_BIN:"* ]]; then
    echo -e "${YELLOW}Warning: $GOPATH_BIN is not in your PATH${NC}"
    echo ""
    echo "Add this line to your shell configuration file:"
    echo ""
    
    # Detect shell
    if [ -n "$BASH_VERSION" ]; then
        SHELL_RC="$HOME/.bashrc"
    elif [ -n "$ZSH_VERSION" ]; then
        SHELL_RC="$HOME/.zshrc"
    else
        SHELL_RC="$HOME/.profile"
    fi
    
    echo -e "${GREEN}  export PATH=\"\$PATH:\$(go env GOPATH)/bin\"${NC}"
    echo ""
    echo "Then run: source $SHELL_RC"
    echo ""
    
    # Ask if user wants to add it automatically
    read -p "Would you like to add it automatically? (y/n): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "export PATH=\"\$PATH:\$(go env GOPATH)/bin\"" >> "$SHELL_RC"
        echo -e "${GREEN}[OK] Added to $SHELL_RC${NC}"
        echo "Please run: source $SHELL_RC"
        echo "Or restart your terminal"
    fi
else
    echo -e "${GREEN}[OK] GOPATH/bin is in your PATH${NC}"
fi

register_file_associations

echo ""
echo "======================================"
echo -e "${GREEN}Installation complete!${NC}"
echo "======================================"
echo ""
echo "Verify installation with:"
echo "  polyloft version"
echo ""
echo "Get started with:"
echo "  polyloft init          # Initialize new project"
echo "  polyloft run file.pf   # Run a Polyloft file"
echo "  polyloft repl          # Start interactive REPL"
echo ""
echo "For more information, visit:"
echo "  https://github.com/ArubikU/polyloft"
echo ""
