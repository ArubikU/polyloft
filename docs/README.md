# Polyloft Language Documentation

Welcome to the complete Polyloft language reference! This documentation covers all language features, standard library modules, and built-in types.

## Table of Contents

### 🚀 Getting Started
- [Quick Start Guide](quickstart.md)
- [Installation](../README.md#installation)
- [Your First Program](tutorial.md)

### 📖 Language Fundamentals

#### Variables & Constants
- [Variables: `var`, `let`, `const`, `final`](variables/README.md)
  - [`var` - Mutable type variable](variables/var.md)
  - [`let` - Type-inferred variable](variables/let.md)
  - [`const` - Compile-time constant](variables/const.md)
  - [`final` - Runtime constant](variables/final.md)

#### Functions & Classes
- [Functions](definitions/function.md)
- [Classes](definitions/class.md)
  - [Access Modifiers - public/private/protected/static](language/access-modifiers.md)
  - [Annotations - @Override and more](language/annotations.md)
- [Interfaces](language/interfaces.md)
- [Sealed Classes](language/sealed.md)
- [Enums](language/enums.md)
- [Records](language/records.md)
- [Import - Module system](language/import.md)

#### Control Flow
- [Conditionals: `if/else`](control-flow/conditionals.md)
- [Loops: `for`, `loop`, `do-loop`](control-flow/loops.md)
- [Switch Statement: Pattern matching](control-flow/switch.md)
- [Exception Handling: `try/catch`](control-flow/exceptions.md)

#### Advanced Features
- [Async/Await](advanced/async-await.md)
- [Promises & Futures](advanced/async-await.md#completablefuture)

### 📚 Standard Library

#### Core Modules
- [**Sys** - System utilities](stdlib/sys.md)
- [**Math** - Mathematical functions](stdlib/math.md)
- [**IO** - Input/Output operations](stdlib/io.md)
- [**Http** - HTTP client & server](stdlib/http.md)
- [**Crypto** - Cryptographic functions & encoding](stdlib/crypto.md)
- [**JSON** - JSON serialization & parsing](stdlib/json.md)

#### Built-in Types
- [**String** - Text manipulation](types/string.md)
- [**Array** - Dynamic arrays](types/array.md)
- [**Map** - Key-value mappings](types/map.md)
- [**Set** - Unique values](types/set.md)
- [**Bytes** - Binary data](types/bytes.md)

### 📝 Examples & Tutorials
- [Algorithm Examples](../algorithm_samples/README.md)
- [Common Patterns](examples/patterns.md)
- [Best Practices](examples/best-practices.md)

### 📚 Reference
- [**Quick Reference** - Complete cheat sheet](QUICK_REFERENCE.md)
- [**CLI Reference** - Command-line interface](CLI.md)
- [VSCode Extension](vscode-extension.md)
- [Build System](build-system.md)

## Quick Reference

### Hello World
```pf
println("Hello, World!")
```

### Variables
```pf
let name = "Alice"      // Type inferred as String
var age = 25            // Mutable, type can change
const PI = 3.14159      // Compile-time constant
final MAX_SIZE = 100    // Runtime constant
```

### Functions
```pf
def greet(name: String): String
    return "Hello, " + name
end

let result = greet("Alice")
println(result)  // Outputs: Hello, Alice
```

### Classes
```pf
class Person:
    let name: String
    let age: Int
    
    def init(name: String, age: Int):
        this.name = name
        this.age = age
    end
    
    def greet():
        println("Hello, I'm #{this.name}")
    end
end

let person = Person("Alice", 25)
person.greet()
```

### Control Flow
```pf
// If-else
if age >= 18:
    println("Adult")
else:
    println("Minor")
end

// For loop
for i in range(10):
    println(i)
end

// Loop statement (conditional loop)
let count = 0
loop count < 5:
    println(count)
    count = count + 1
end
```

### Error Handling
```pf
try:
    let result = riskyOperation()
    println("Success: #{result}")
catch e:
    println("Error: #{e}")
end
```

## Contributing

Found an issue or want to improve the documentation? See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

## License

Polyloft is licensed under the MIT License. See [LICENSE](../LICENSE) for details.
