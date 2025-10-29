# Import and Module System

Polyloft provides a module system for organizing and reusing code across different files and projects.

## Built-in Modules 

These modules are automatically available without imports:

- **Math** - Mathematical functions and constants
- **Sys** - System operations (time, sleep, random, etc.)
- **IO** - File and directory operations  
- **Net** - TCP networking
- **Http** - HTTP client and server
- **List** - Generic list collection
- **Map** - Generic map/dictionary collection
- **Set** - Generic set collection

**Usage:**
```polyloft
// No import needed - use directly
let result = Math.sqrt(16)
let time = Sys.time()
let content = IO.readFile("data.txt")
```

## File Organization

### Single File Programs

For simple scripts, all code can be in one file:

```polyloft
// hello.pf
def greet(name: String):
    println("Hello, " + name + "!")
end

greet("World")
```

Run with:
```bash
polyloft run hello.pf
```

### Multi-File Projects

For larger projects, organize code into multiple files:

```
myproject/
├── main.pf          # Entry point
├── lib/
│   ├── utils.pf     # Utility functions
│   └── models.pf    # Data models
└── tests/
    └── test.pf      # Tests
```

## Modules and Namespaces

### Defining a Module

Create a module by defining classes and functions in a file:

```polyloft
// lib/math_utils.pf
class Calculator:
    Calculator():
        // Constructor
    end
    
    def add(a: Int, b: Int) -> Int:
        return a + b
    end
    
    def multiply(a: Int, b: Int) -> Int:
        return a * b
    end
end

def factorial(n: Int) -> Int:
    if n <= 1:
        return 1
    end
    return n * factorial(n - 1)
end
```

### Using Code from Other Files

Currently, Polyloft evaluates files independently. To share code:

1. **Define in the same file** for now
2. **Use the CLI to combine files** during build
3. **Future:** `import` statement (coming soon)

**Workaround - Include Pattern:**
```bash
# Concatenate files for now
cat lib/utils.pf main.pf > combined.pf
polyloft run combined.pf
```

## Static Modules (Built-in)

Built-in modules use static classes pattern:

```polyloft
// Math is a static module
class Math:
    static let PI: Float = 3.14159
    
    static def sqrt(x: Float) -> Float:
        // Implementation
    end
end

// Usage - no instance needed
let value = Math.sqrt(25)
let circle = 2 * Math.PI * radius
```

## Creating Your Own Static Modules

Define a class with only static members:

```polyloft
class StringUtils:
    static def reverse(s: String) -> String:
        let result = ""
        for i in range(s.length() - 1, -1, -1):
            result = result + s[i]
        end
        return result
    end
    
    static def capitalize(s: String) -> String:
        if s.length() == 0:
            return s
        end
        return s[0].toUpper() + s.substring(1)
    end
end

// Usage
let reversed = StringUtils.reverse("hello")
let capped = StringUtils.capitalize("world")
```

## Package Structure

### Recommended Project Structure

```
myproject/
├── main.pf              # Entry point
├── src/
│   ├── models/          # Data models
│   │   ├── user.pf
│   │   └── post.pf
│   ├── services/        # Business logic
│   │   ├── auth.pf
│   │   └── database.pf
│   └── utils/           # Utilities
│       ├── validators.pf
│       └── formatters.pf
├── tests/               # Tests
│   ├── test_models.pf
│   └── test_services.pf
└── README.md            # Documentation
```

### Entry Point (main.pf)

```polyloft
// main.pf
def main():
    println("Application started")
    
    // Initialize components
    let app = Application()
    app.run()
end

class Application:
    Application():
        // Setup
    end
    
    def run():
        println("Running application...")
    end
end

// Start the application
main()
```

## Future Import Syntax (Planned)

The import system is being developed. Expected syntax:

```polyloft
// Import specific items
import { Calculator, factorial } from "lib/math_utils"

// Import all
import * from "lib/math_utils"

// Import with alias
import Calculator as Calc from "lib/math_utils"

// Import module as namespace
import "lib/math_utils" as MathUtils
```

## Working with External Libraries

### Installing Packages

Use the CLI to manage dependencies:

```bash
# Initialize project
polyloft init myproject

# Add dependency (future feature)
polyloft install github.com/user/library

# Update dependencies
polyloft update
```

### Package Manifest (polyloft.toml)

```toml
[package]
name = "myproject"
version = "1.0.0"
author = "Your Name"
description = "My Polyloft project"

[dependencies]
# Future: external dependencies
# somelib = "1.0.0"

[scripts]
test = "polyloft test"
build = "polyloft build"
```

## Best Practices

### 1. Organize by Feature

Group related code together:
```
features/
├── user/
│   ├── user_model.pf
│   ├── user_service.pf
│   └── user_controller.pf
└── auth/
    ├── auth_service.pf
    └── token.pf
```

### 2. Use Clear Naming

```polyloft
// Good
class UserRepository:
    def findById(id: Int) -> User:
        // ...
    end
end

// Bad  
class UR:
    def fbi(i: Int) -> U:
        // ...
    end
end
```

### 3. Separate Concerns

Keep different responsibilities in separate files:
- Models: Data structures
- Services: Business logic
- Utils: Helper functions
- Constants: Configuration

### 4. Document Public APIs

```polyloft
// UserService provides user management functionality
class UserService:
    // Creates a new user with the given name and email
    // Returns the created user with assigned ID
    def createUser(name: String, email: String) -> User:
        // Implementation
    end
end
```

### 5. Keep Files Focused

Each file should have a single responsibility:
- One main class per file
- Related helper functions with the class
- Maximum 300-500 lines per file

## Module Resolution (Current)

Currently, Polyloft loads files explicitly:

```bash
# Run a single file
polyloft run main.pf

# Build project (processes all .pf files in directory)
polyloft build

# REPL with context
polyloft repl --load utils.pf
```

## Examples

### Example 1: Math Utilities Module

```polyloft
// mathutils.pf
class MathUtils:
    static def gcd(a: Int, b: Int) -> Int:
        if b == 0:
            return a
        end
        return MathUtils.gcd(b, a % b)
    end
    
    static def lcm(a: Int, b: Int) -> Int:
        return (a * b) / MathUtils.gcd(a, b)
    end
    
    static def isPrime(n: Int) -> Bool:
        if n <= 1:
            return false
        end
        for i in range(2, Math.sqrt(n).toInt() + 1):
            if n % i == 0:
                return false
            end
        end
        return true
    end
end

// Usage
let result = MathUtils.gcd(48, 18)
println("GCD: " + result.toString())
```

### Example 2: Configuration Module

```polyloft
// config.pf
class Config:
    static let APP_NAME: String = "MyApp"
    static let VERSION: String = "1.0.0"
    static let DEBUG: Bool = true
    static let PORT: Int = 8080
    
    static def getDatabaseUrl() -> String:
        if Config.DEBUG:
            return "localhost:5432/dev_db"
        else:
            return "prod.server.com:5432/prod_db"
        end
    end
end

// Usage
println(Config.APP_NAME + " v" + Config.VERSION)
let dbUrl = Config.getDatabaseUrl()
```

### Example 3: Validator Module

```polyloft
// validators.pf
class Validators:
    static def isEmail(email: String) -> Bool:
        return email.contains("@") and email.contains(".")
    end
    
    static def isStrongPassword(password: String) -> Bool:
        if password.length() < 8:
            return false
        end
        // Add more checks
        return true
    end
    
    static def isValidUrl(url: String) -> Bool:
        return url.startsWith("http://") or url.startsWith("https://")
    end
end

// Usage
if Validators.isEmail("user@example.com"):
    println("Valid email")
end
```

## See Also

- [CLI Tools](../cli/overview.md) - Command-line interface
- [Project Structure](../cli/init.md) - Initialize projects
- [Publishing](../cli/publishing.md) - Share your modules
- [Building](../cli/build.md) - Compile your code
