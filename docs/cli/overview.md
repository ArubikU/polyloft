# CLI Overview

The Polyloft CLI provides tools for developing, building, and managing Polyloft projects.

## Commands

### polyloft version
Display the current version of Polyloft.

```bash
polyloft version
```

Output:
```
v0.2.6b built 2025-10-18
```

### polyloft run
Run a Polyloft file or project.

```bash
# Run a single file
polyloft run hello.pf

# Run current project (uses polyloft.toml)
polyloft run

# Run with arguments
polyloft run script.pf arg1 arg2
```

See [Running Programs](run.md) for details.

### polyloft build
Build an executable from your Polyloft project.

```bash
# Build with default name
polyloft build

# Build with custom name
polyloft build -o myapp

# Build for different platform
polyloft build -o myapp.exe
```

See [Building Executables](build.md) for details.

### polyloft init
Initialize a new Polyloft project.

```bash
# Create new project in current directory
polyloft init

# Create new project in specific directory
polyloft init my-project
```

Creates a `polyloft.toml` configuration file.

See [Project Initialization](init.md) for details.

### polyloft install
Install dependencies specified in `polyloft.toml`.

```bash
polyloft install
```

Downloads and installs both Polyloft packages and Go modules.

See [Dependency Management](dependencies.md) for details.

### polyloft repl
Start an interactive REPL (Read-Eval-Print Loop).

```bash
polyloft repl
```

See [REPL](repl.md) for details.

### polyloft search
Search for packages in the Polyloft registry.

```bash
polyloft search vector
polyloft search math
```

See [Package Publishing](publishing.md) for details.

### polyloft register
Register a new account on the package registry.

```bash
polyloft register
```

Prompts for username, email, and password.

### polyloft login
Authenticate with the package registry.

```bash
polyloft login
```

Prompts for username and password. Credentials are stored in `~/.polyloft/credentials.json`.

### polyloft logout
Clear authentication credentials.

```bash
polyloft logout
```

### polyloft publish
Publish your package to the registry.

```bash
polyloft publish
```

Requires authentication and a valid `polyloft.toml` configuration.

See [Package Publishing](publishing.md) for details.

### polyloft generate-mappings
Generate IDE mappings for enhanced IntelliSense.

```bash
polyloft generate-mappings
```

Creates `.polyloft/mappings.json` for use by the VSCode extension.

## Project Structure

### polyloft.toml

Configuration file for Polyloft projects:

```toml
[project]
name = "my-project"
version = "1.0.0"
entry_point = "src/main.pf"
author = "Your Name"
description = "Project description"

[[dependencies.pf]]
name = "package@author"
version = "1.0.0"

[[dependencies.go]]
name = "github.com/some/module"
version = "v1.0.0"
```

### Directory Structure

Recommended project layout:

```
my-project/
├── polyloft.toml          # Project configuration
├── src/
│   ├── main.pf            # Entry point
│   ├── lib/
│   │   └── utils.pf       # Library code
│   └── test/
│       └── test_main.pf   # Tests
├── .polyloft/             # IDE support files
│   └── mappings.json      # Generated mappings
└── README.md              # Documentation
```

## Global Options

### Environment Variables

- `POLYLOFT_REGISTRY_URL` - Package registry URL (default: https://registry.polyloft.org)

### Configuration Directory

Polyloft stores configuration in `~/.polyloft/`:
- `credentials.json` - Registry authentication
- `cache/` - Downloaded packages

## Common Workflows

### Creating a New Project

```bash
# Initialize project
polyloft init my-project
cd my-project

# Create source file
cat > src/main.pf << 'EOF'
println("Hello, Polyloft!")
EOF

# Run project
polyloft run

# Build executable
polyloft build -o myapp
./myapp
```

### Adding Dependencies

```bash
# Search for packages
polyloft search vector

# Edit polyloft.toml to add dependency
# Then install
polyloft install

# Use in your code
# import vectores from "vectores@arubiku"
```

### Publishing a Package

```bash
# Register account (once)
polyloft register

# Login
polyloft login

# Prepare package
# - Ensure polyloft.toml is complete
# - Test your code

# Publish
polyloft publish
```

### Development Cycle

```bash
# Edit code
vim src/main.pf

# Run immediately
polyloft run

# Or use REPL for testing
polyloft repl

# Build when ready
polyloft build
```

## Tips

1. **Use REPL for Experimentation**: Test code snippets quickly
2. **Build Incrementally**: Run frequently during development
3. **Version Your Packages**: Use semantic versioning
4. **Document Your Code**: Include README and comments
5. **Test Before Publishing**: Ensure code works on clean install

## Exit Codes

- `0` - Success
- `1` - Error (compilation, runtime, etc.)
- `2` - Invalid arguments
- Other codes indicate specific errors

## See Also

- [Project Initialization](init.md) - Creating projects
- [Running Programs](run.md) - Executing code
- [Building Executables](build.md) - Compilation
- [REPL](repl.md) - Interactive shell
- [Dependency Management](dependencies.md) - Managing packages
- [Package Publishing](publishing.md) - Sharing code
