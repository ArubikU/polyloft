#!/bin/bash
# Polyloft Update Script for Linux/macOS

set -e

echo "======================================"
echo "    Polyloft Update Script"
echo "======================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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
Exec="$(go env GOPATH)/bin/polyloft" run %f
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
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go 1.22.0 or later from https://go.dev/dl/"
    exit 1
fi

# Check if Polyloft is currently installed
GOPATH_BIN=$(go env GOPATH)/bin
POLYLOFT_PATH="$GOPATH_BIN/polyloft"

if [ ! -f "$POLYLOFT_PATH" ]; then
    echo -e "${YELLOW}Polyloft is not currently installed${NC}"
    echo "Would you like to install it now? (y/n): "
    read -r -n 1 REPLY
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Redirecting to installation script..."
        curl -fsSL https://raw.githubusercontent.com/ArubikU/polyloft/main/scripts/install.sh | bash
        exit 0
    else
        echo "Update cancelled"
        exit 1
    fi
fi

# Get current version
if command -v polyloft &> /dev/null; then
    CURRENT_VERSION=$(polyloft version 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' || echo "unknown")
    echo -e "Current version: ${BLUE}$CURRENT_VERSION${NC}"
else
    echo -e "${YELLOW}Warning: polyloft command not in PATH${NC}"
    CURRENT_VERSION="unknown"
fi
echo ""

# Update Polyloft
echo "Updating Polyloft to latest version..."
if go install github.com/ArubikU/polyloft/cmd/polyloft@latest; then
    echo -e "${GREEN}[OK] Polyloft updated successfully${NC}"
else
    echo -e "${RED}[ERROR] Update failed${NC}"
    exit 1
fi
echo ""

# Get new version
if command -v polyloft &> /dev/null; then
    NEW_VERSION=$(polyloft version 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' || echo "unknown")
    echo -e "New version: ${GREEN}$NEW_VERSION${NC}"
    
    # Check if version changed
    if [ "$CURRENT_VERSION" = "$NEW_VERSION" ] && [ "$CURRENT_VERSION" != "unknown" ]; then
        echo -e "${YELLOW}You already had the latest version${NC}"
    fi
else
    echo -e "${YELLOW}Warning: Run 'source ~/.bashrc' or restart terminal if polyloft command not found${NC}"
fi

register_file_associations

echo ""
echo "======================================"
echo -e "${GREEN}Update complete!${NC}"
echo "======================================"
echo ""
echo "Verify with:"
echo "  polyloft version"
echo ""
