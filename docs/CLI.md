# Polyloft CLI Reference

Complete reference for the Polyloft command-line interface.

## Installation

```bash
# Build from source
go build -o polyloft ./cmd/polyloft

# Add to PATH
export PATH=$PATH:/path/to/polyloft
```

## Commands

### `polyloft repl`

Start an interactive REPL (Read-Eval-Print Loop).

**Usage:**
```bash
polyloft repl [options]
```

**Options:**
- `--prompt <string>` - Custom REPL prompt (default: ">>> ")

**Examples:**
```bash
# Start REPL with default prompt
polyloft repl

# Start REPL with custom prompt
polyloft repl --prompt "pf> "
```

**REPL Session:**
```
>>> let x = 10
>>> let y = 20
>>> println(x + y)
30
>>> def greet(name)
...     return "Hello, #{name}"
... end
>>> greet("Alice")
"Hello, Alice"
>>> exit
```

### `polyloft run`

Execute a Polyloft source file.

**Usage:**
```bash
polyloft run [options] <file.pf>
```

**Options:**
- `--config <file>` - Configuration file (default: "polyloft.toml")

**Examples:**
```bash
# Run a script
polyloft run script.pf

# Run with custom config
polyloft run --config myconfig.toml app.pf
```

### `polyloft build`

Compile a Polyloft project to an executable or library.

**Usage:**
```bash
polyloft build [options]
```

**Options:**
- `-o <output>` - Output file name (defaults to project name)
- `--config <file>` - Configuration file (default: "polyloft.toml")

**Examples:**
```bash
# Build project
polyloft build

# Build with custom output name
polyloft build -o myapp

# Build with custom config
polyloft build --config build.toml -o release/app
```

### `polyloft init`

Initialize a new Polyloft project.

**Usage:**
```bash
polyloft init [project-name]
```

**Examples:**
```bash
# Initialize project in current directory
polyloft init

# Initialize named project
polyloft init myproject
```

**Generated Structure:**
```
myproject/
├── polyloft.toml
├── main.pf
├── lib/
└── README.md
```

### `polyloft install`

Install project dependencies.

**Usage:**
```bash
polyloft install [options]
```

**Options:**
- `--config <file>` - Configuration file (default: "polyloft.toml")
- `-g` - Install globally

**Examples:**
```bash
# Install dependencies from polyloft.toml
polyloft install

# Install package globally
polyloft install -g package-name
```

### `polyloft publish`

Publish a package to the Polyloft registry.

**Usage:**
```bash
polyloft publish [options]
```

**Options:**
- `--config <file>` - Configuration file (default: "polyloft.toml")

**Examples:**
```bash
# Publish package
polyloft publish

# Publish with custom config
polyloft publish --config release.toml
```

### `polyloft search`

Search for packages in the registry.

**Usage:**
```bash
polyloft search <query>
```

**Examples:**
```bash
# Search for packages
polyloft search http
polyloft search database
polyloft search "web framework"
```

### `polyloft update`

Update dependencies to latest versions.

**Usage:**
```bash
polyloft update
```

**Examples:**
```bash
# Update all dependencies
polyloft update
```

### `polyloft register`

Register a new account in the Polyloft registry.

**Usage:**
```bash
polyloft register
```

**Interactive Prompts:**
- Username
- Email
- Password

### `polyloft login`

Log in to the Polyloft registry.

**Usage:**
```bash
polyloft login
```

**Interactive Prompts:**
- Username or Email
- Password

### `polyloft logout`

Log out from the Polyloft registry.

**Usage:**
```bash
polyloft logout
```

### `polyloft generate-mappings`

Generate IDE mappings for better editor support.

**Usage:**
```bash
polyloft generate-mappings [options]
```

**Options:**
- `-o <file>` - Output file path (default: "mappings.json")
- `--root <dir>` - Root directory of the project (default: ".")

**Examples:**
```bash
# Generate mappings
polyloft generate-mappings

# Custom output location
polyloft generate-mappings -o .polyloft/mappings.json

# Specify project root
polyloft generate-mappings --root /path/to/project
```

### `polyloft version`

Display version information.

**Usage:**
```bash
polyloft version
```

**Example Output:**
```
Polyloft v1.0.0
Build: abc123def
Go: go1.21.0
```

### `polyloft help`

Display help information.

**Usage:**
```bash
polyloft help
polyloft -h
polyloft --help
```

## Configuration File

The `polyloft.toml` file configures your project:

```toml
[project]
name = "myapp"
version = "1.0.0"
author = "Your Name"
description = "My Polyloft application"

[dependencies]
http = "^1.0.0"
json = "^2.1.0"

[build]
entry = "main.pf"
output = "bin/myapp"
target = "native"

[dev]
port = 8080
reload = true
```

## Common Workflows

### Development

```bash
# Create new project
polyloft init myapp
cd myapp

# Start REPL for testing
polyloft repl

# Run during development
polyloft run main.pf

# Build for release
polyloft build -o release/myapp
```

### Package Management

```bash
# Install dependencies
polyloft install

# Add global tool
polyloft install -g linter

# Update packages
polyloft update

# Publish your package
polyloft login
polyloft publish
```

### Project Setup

```bash
# Initialize
polyloft init my-web-app

# Generate IDE support
polyloft generate-mappings

# Install deps
polyloft install

# Run
polyloft run main.pf
```

## Environment Variables

### `POLYLOFT_HOME`
Directory for global packages and configuration.

```bash
export POLYLOFT_HOME=~/.polyloft
```

### `POLYLOFT_PATH`
Additional paths for module resolution.

```bash
export POLYLOFT_PATH=/usr/local/lib/polyloft:/opt/polyloft/lib
```

## Exit Codes

- `0` - Success
- `1` - General error
- `2` - Parse error
- `3` - Runtime error
- `4` - Configuration error

## Examples

### Quick Script

```bash
# Create a simple script
echo 'println("Hello, Polyloft!")' > hello.pf

# Run it
polyloft run hello.pf
```

### Web Server

```bash
# Create server.pf
cat > server.pf << 'EOF'
let server = Http.createServer()

server.get("/", (req, res) => do
    res.send("Hello, World!")
end)

server.listen(8080)
println("Server running on http://localhost:8080")
EOF

# Run the server
polyloft run server.pf
```

### REPL Scripting

```bash
# Pass commands to REPL
echo 'let x = 10
println(x * x)
exit' | polyloft repl
```

## Debugging

### Verbose Output

```bash
# Enable verbose logging (if supported)
POLYLOFT_DEBUG=1 polyloft run app.pf
```

### Stack Traces

Errors include stack traces by default:

```
RuntimeError: Division by zero
  at calculateAverage (main.pf:15:20)
  at processData (main.pf:25:5)
  at main (main.pf:40:1)
```

## See Also

- [Quick Reference](QUICK_REFERENCE.md)
- [Project Structure](guides/project-structure.md)
- [Package Management](guides/packages.md)
