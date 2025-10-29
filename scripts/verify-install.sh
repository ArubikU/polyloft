#!/bin/bash
# Polyloft Installation Verification Script for Linux/macOS

set +e  # Continue on errors

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
GRAY='\033[0;37m'
NC='\033[0m' # No Color

echo -e "${CYAN}======================================${NC}"
echo -e "${CYAN}   Polyloft Installation Verification${NC}"
echo -e "${CYAN}======================================${NC}"
echo ""

all_ok=true

# 1. Check if polyloft binary exists
echo -e "${CYAN}[1] Checking Polyloft binary...${NC}"
GOPATH_BIN="$(go env GOPATH)/bin"
POLYLOFT_PATH="$GOPATH_BIN/polyloft"

if [ -f "$POLYLOFT_PATH" ]; then
    echo -e "    ${GREEN}[OK] Binary found at: $POLYLOFT_PATH${NC}"
    
    # Check if executable
    if [ -x "$POLYLOFT_PATH" ]; then
        echo -e "    ${GREEN}[OK] Binary is executable${NC}"
    else
        echo -e "    ${RED}[ERROR] Binary is NOT executable${NC}"
        all_ok=false
    fi
else
    echo -e "    ${RED}[ERROR] Binary NOT found at: $POLYLOFT_PATH${NC}"
    all_ok=false
fi
echo ""

# 2. Check if GOPATH/bin is in PATH
echo -e "${CYAN}[2] Checking if GOPATH/bin is in PATH...${NC}"
if [[ ":$PATH:" == *":$GOPATH_BIN:"* ]]; then
    echo -e "    ${GREEN}[OK] GOPATH/bin is in PATH${NC}"
else
    echo -e "    ${RED}[ERROR] GOPATH/bin is NOT in PATH${NC}"
    all_ok=false
fi
echo ""

# 3. Try to execute polyloft
echo -e "${CYAN}[3] Testing Polyloft command...${NC}"
if command -v polyloft &> /dev/null; then
    polyloft_output=$(polyloft version 2>&1)
    exit_code=$?
    
    if [ $exit_code -eq 0 ]; then
        echo -e "    ${GREEN}[OK] Command executed successfully${NC}"
        echo -e "    ${GRAY}Output: $polyloft_output${NC}"
    else
        echo -e "    ${YELLOW}[WARNING] Command executed with exit code: $exit_code${NC}"
        echo -e "    ${GRAY}Output: $polyloft_output${NC}"
    fi
else
    echo -e "    ${RED}[ERROR] 'polyloft' command not found${NC}"
    all_ok=false
fi
echo ""

# 4. Check shell configuration files
echo -e "${CYAN}[4] Checking shell configuration...${NC}"
SHELL_NAME=$(basename "$SHELL")
echo -e "    Current shell: ${GRAY}$SHELL_NAME${NC}"

config_files=()
case "$SHELL_NAME" in
    bash)
        config_files=("$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.profile")
        ;;
    zsh)
        config_files=("$HOME/.zshrc" "$HOME/.zprofile")
        ;;
    fish)
        config_files=("$HOME/.config/fish/config.fish")
        ;;
    *)
        echo -e "    ${YELLOW}[WARNING] Unknown shell: $SHELL_NAME${NC}"
        ;;
esac

found_in_config=false
for config_file in "${config_files[@]}"; do
    if [ -f "$config_file" ]; then
        if grep -q "GOPATH" "$config_file" || grep -q "go/bin" "$config_file"; then
            echo -e "    ${GREEN}[OK] Found Go/GOPATH configuration in: $config_file${NC}"
            found_in_config=true
            break
        fi
    fi
done

if [ "$found_in_config" = false ]; then
    echo -e "    ${YELLOW}[WARNING] No Go/GOPATH configuration found in shell config files${NC}"
    echo -e "    ${YELLOW}Consider adding to your shell config (~/.bashrc, ~/.zshrc, etc.):${NC}"
    echo -e "    ${GRAY}export PATH=\$PATH:\$(go env GOPATH)/bin${NC}"
fi
echo ""

# 5. Check Go installation
echo -e "${CYAN}[5] Checking Go installation...${NC}"
if command -v go &> /dev/null; then
    go_version=$(go version)
    echo -e "    ${GREEN}[OK] $go_version${NC}"
    echo -e "    ${GRAY}GOPATH: $(go env GOPATH)${NC}"
    echo -e "    ${GRAY}GOROOT: $(go env GOROOT)${NC}"
else
    echo -e "    ${RED}[ERROR] Go is not installed${NC}"
    all_ok=false
fi
echo ""

# Summary
echo -e "${CYAN}======================================${NC}"
if [ "$all_ok" = true ]; then
    echo -e "${GREEN}Verification Result: ALL CHECKS PASSED${NC}"
    echo ""
    echo -e "${GREEN}Polyloft is installed correctly!${NC}"
    echo ""
    echo -e "${YELLOW}If commands still don't work, try:${NC}"
    echo -e "  ${NC}1. Reload your shell configuration:${NC}"
    echo -e "     ${CYAN}source ~/.bashrc${NC}  (for bash)"
    echo -e "     ${CYAN}source ~/.zshrc${NC}   (for zsh)"
    echo -e "  ${NC}2. Or restart your terminal${NC}"
else
    echo -e "${RED}Verification Result: SOME ISSUES FOUND${NC}"
    echo ""
    echo -e "${YELLOW}Recommended actions:${NC}"
    echo -e "  ${NC}1. Re-run the installation script: ./scripts/install.sh${NC}"
    echo -e "  ${NC}2. Add GOPATH/bin to your PATH by adding this line to your shell config:${NC}"
    echo -e "     ${CYAN}export PATH=\$PATH:\$(go env GOPATH)/bin${NC}"
    echo -e "  ${NC}3. Reload your shell configuration or restart your terminal${NC}"
    echo -e "  ${NC}4. Run this verification script again${NC}"
fi
echo -e "${CYAN}======================================${NC}"
echo ""

# Additional diagnostic info
echo -e "${CYAN}Diagnostic Information:${NC}"
echo -e "  ${GRAY}GOPATH: $(go env GOPATH)${NC}"
echo -e "  ${GRAY}GOPATH/bin: $GOPATH_BIN${NC}"
echo -e "  ${GRAY}Polyloft path: $POLYLOFT_PATH${NC}"
echo -e "  ${GRAY}Current PATH: $PATH${NC}"
echo ""
