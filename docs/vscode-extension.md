# Polyloft VSCode Extension

The Polyloft VSCode extension provides comprehensive language support for developing with the Polyloft programming language.

## Features

### 1. Syntax Highlighting

Full syntax highlighting support for `.pf` files including:
- Keywords (var, let, const, class, def, if, for, etc.)
- Control flow statements
- Operators and delimiters
- String literals with interpolation (`#{expression}`)
- Numeric literals (integers, floats, floats)
- Comments (line and block)
- Class and interface definitions
- Function declarations
- Type annotations

### 2. Intelligent Linting

Real-time error detection and validation:
- Syntax errors (unclosed strings, unmatched brackets)
- Missing `end` keywords for blocks
- Invalid function declarations
- Use of `this` outside class context
- Return statements outside functions
- Convention warnings (e.g., class names should start with uppercase)

Configuration:
```json
{
  "polyloft.linting.enabled": true,
  "polyloft.linting.onType": true
}
```

### 3. Auto-completion

Context-aware code completion for:

**Keywords and Language Constructs:**
- All Polyloft keywords (var, let, const, class, def, if, for, etc.)
- Type names (String, Int, Float, Float, Bool, Void, Array, Map)

**Built-in Functions:**
- Global functions: `println`, `print`
- System functions: `Sys.time()`, `Sys.random()`, `Sys.sleep()`
- Math functions: `Math.sqrt()`, `Math.pow()`, `Math.sin()`, `Math.cos()`, etc.
- Math constants: `Math.PI`, `Math.E`

**User-defined Symbols:**
- Class names
- Function names
- Variable names
- Method names (after dot operator)
- Imported symbols

**Smart Completion:**
- Member completion after dot operator
- Parameter snippets for functions
- Import statement completion

### 4. Go to Definition

Navigate to symbol definitions:
- Jump to class definitions
- Jump to function/method definitions
- Jump to variable declarations
- Cross-file navigation for imports
- Resolve standard library imports

Press `F12` or `Ctrl+Click` on a symbol to jump to its definition.

### 5. Hover Information

Rich hover tooltips showing:
- Function signatures with parameter types and return types
- Type information for variables
- Documentation for built-in functions
- Class inheritance information
- Constant values

Hover over any symbol to see its information.

### 6. Multi-file Support

The extension understands project structure:
- Parse and resolve `import` statements
- Cross-file symbol resolution
- Support for standard library paths (`libs/`)
- Import path resolution strategies

Example imports:
```polyloft
import math.vector { Vec2, Vec3 }
import utils { add, greet }
```

## Built-in Packages

The extension includes comprehensive definitions for Polyloft's standard library:

### Global Functions
- `println(...args)` - Print to stdout with newline
- `print(...args)` - Print to stdout without newline

### Sys Package
System utilities:
- `Sys.time(format?)` - Current time in milliseconds (or Float if format="float")
- `Sys.random()` - Random float between 0.0 and 1.0
- `Sys.sleep(ms)` - Sleep for specified milliseconds

### Math Package
Mathematical functions and constants:

**Constants:**
- `Math.PI` - 3.141592653589793
- `Math.E` - 2.718281828459045

**Functions:**
- `Math.sqrt(x)` - Square root
- `Math.pow(base, exp)` - Power
- `Math.abs(x)` - Absolute value
- `Math.sin(x)`, `Math.cos(x)`, `Math.tan(x)` - Trigonometric functions
- `Math.floor(x)`, `Math.ceil(x)`, `Math.round(x)` - Rounding
- `Math.min(a, b)`, `Math.max(a, b)` - Min/Max

### Standard Library

**math.vector:**
- `Vec2` - 2D vector class
- `Vec3` - 3D vector class
- `Vec4` - 4D vector class

**utils:**
- `add(a, b)` - Add two integers
- `greet(name)` - Print greeting message

## Mappings Generation

The extension can use `mappings.json` files to provide enhanced IntelliSense for user libraries.

### Generating Mappings

Use the Polyloft CLI to generate mappings for your project:

```bash
polyloft generate-mappings -o mappings.json
```

This scans all `.pf` files in the `libs/` directory and generates a JSON file with:
- All class definitions with methods and fields
- All function definitions with parameters and return types
- Import/export information
- Type information
- File locations

### Mappings File Structure

```json
{
  "version": "1.0.0",
  "packages": {
    "package.name": {
      "name": "package.name",
      "path": "relative/path",
      "version": "1.0.0",
      "symbols": [
        {
          "name": "ClassName",
          "type": "class",
          "methods": [...],
          "fields": [...],
          "file": "path/to/file.pf",
          "line": 1
        }
      ],
      "imports": ["dependency1", "dependency2"],
      "exports": ["ClassName", "functionName"]
    }
  }
}
```

### Using Mappings in Projects

1. Generate mappings for your project libraries
2. Place `mappings.json` in your project root or `libs/` directory
3. The extension automatically loads and uses the mappings
4. Get accurate completions for all your project's symbols

## Language Configuration

The extension provides smart editing features:

- **Auto-closing pairs:** `{}`, `[]`, `()`, `""`, `''`
- **Bracket matching:** Matching pairs highlighted
- **Comment toggling:** `Ctrl+/` for line comments
- **Auto-indentation:** Smart indentation for blocks
- **Folding:** Code folding based on blocks and regions

## Installation

See [INSTALL.md](../vscode-extension/INSTALL.md) for detailed installation instructions.

Quick install for development:
```bash
cd vscode-extension
npm install
npm run compile
code --install-extension .
```

## Configuration Options

All configuration options are available in VSCode settings:

```json
{
  "polyloft.linting.enabled": true,
  "polyloft.linting.onType": true,
  "polyloft.completion.enabled": true,
  "polyloft.trace.server": "off"
}
```

### Configuration Details

- `polyloft.linting.enabled` - Enable/disable linting (default: true)
- `polyloft.linting.onType` - Enable linting while typing (default: true)
- `polyloft.completion.enabled` - Enable/disable auto-completion (default: true)
- `polyloft.trace.server` - Trace communication for debugging (default: "off")

## Development

### Extension Architecture

```
vscode-extension/
├── src/
│   ├── extension.ts       # Main extension entry point
│   ├── linter.ts          # Syntax validation and error reporting
│   ├── completion.ts      # Auto-completion provider
│   ├── definition.ts      # Go to definition provider
│   └── hover.ts           # Hover information provider
├── syntaxes/
│   └── polyloft.tmLanguage.json  # TextMate grammar
├── language-configuration/
│   └── language-configuration.json  # Language config
├── builtin-packages.json  # Built-in library definitions
├── package.json           # Extension manifest
└── tsconfig.json          # TypeScript configuration
```

### Building the Extension

```bash
cd vscode-extension
npm install          # Install dependencies
npm run compile      # Compile TypeScript
npm run watch        # Watch mode for development
```

### Testing

1. Open the `vscode-extension` folder in VSCode
2. Press `F5` to launch Extension Development Host
3. Open a `.pf` file in the new window
4. Test all features (highlighting, completion, linting, etc.)

### Adding New Built-in Functions

Edit `builtin-packages.json` to add new built-in functions:

```json
{
  "packages": {
    "YourPackage": {
      "type": "builtin",
      "description": "Your package description",
      "functions": [
        {
          "name": "yourFunction",
          "returnType": "String",
          "parameters": [
            { "name": "param1", "type": "Int" }
          ],
          "description": "Function description"
        }
      ]
    }
  }
}
```

## Roadmap

Future enhancements planned:

- [ ] Language Server Protocol (LSP) implementation
- [ ] Semantic highlighting
- [ ] Code actions and quick fixes
- [ ] Refactoring support (rename, extract method, etc.)
- [ ] Debugging support
- [ ] Test runner integration
- [ ] Snippets library
- [ ] Formatting provider
- [ ] Code lens for tests and benchmarks
- [ ] Integration with package registry

## Contributing

Contributions are welcome! The extension is part of the main Polyloft repository.

1. Clone the repository
2. Make changes in `vscode-extension/`
3. Test thoroughly
4. Submit a pull request

## License

The extension is licensed under the same license as Polyloft.
