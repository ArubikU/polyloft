# Documentation Quick Reference

Quick reference guide for Polyloft documentation.

## Essential Pages

### Getting Started
- [Installation](installation.md) - Install Polyloft
- [Getting Started](getting-started.md) - Your first program
- [Hello World Examples](examples/hello-world.md) - Various examples

### Language Basics
- [Syntax](language/syntax.md) - Language syntax rules
- [Variables and Types](language/variables-and-types.md) - Type system
- [Control Flow](language/control-flow.md) - if/for/while statements
- [Functions](language/functions.md) - Defining and using functions

### Object-Oriented Programming
- [Classes and Objects](language/classes-and-objects.md) - OOP basics
- [Interfaces](language/interfaces.md) - Interface definitions
- [Generics](language/generics.md) - Generic programming
- [Enums](language/enums.md) - Enumeration types

### Standard Library
- [Math](stdlib/math.md) - Mathematical functions (sqrt, pow, sin, cos, etc.)
- [Sys](stdlib/sys.md) - System operations (time, sleep, random, input)
- [IO](stdlib/io.md) - File operations (read, write, directories)
- [Net](stdlib/net.md) - Network programming (TCP servers/clients)

### Concurrency
- [Overview](concurrency/overview.md) - Concurrency models
- [Channels](concurrency/channels.md) - Message passing
- [Async/Await](concurrency/async-await.md) - Asynchronous operations
- [Threads](concurrency/threads.md) - Thread management
- [Defer](concurrency/defer.md) - Resource cleanup

### CLI Tools
- [Overview](cli/overview.md) - All CLI commands
- [run](cli/run.md) - Run programs
- [build](cli/build.md) - Build executables
- [repl](cli/repl.md) - Interactive shell
- [install](cli/dependencies.md) - Manage dependencies
- [publish](cli/publishing.md) - Publish packages

## Code Examples by Topic

### Math Operations
```polyloft
Math.sqrt(9)           // 3.0
Math.pow(2, 8)         // 256.0
Math.sin(Math.PI / 2)  // 1.0
Math.random()          // 0.0 to 1.0
```

### File I/O
```polyloft
IO.readFile("data.txt")
IO.writeFile("output.txt", content)
IO.listDir(".")
IO.exists("file.txt")
```

### Network
```polyloft
let server = Net.listen(":8080")
let conn = server.accept()
conn.send("Hello!\n")
let data = conn.recv()
conn.close()
```

### Concurrency
```polyloft
// Channels
let ch = Channel<Int>()
thread spawn do
    ch.send(42)
end
let value = ch.recv()

// Async/Await with Promises
let fetch = () => do
    return Promise((resolve, reject) => do
        let data = getData()
        resolve(data)
    end)
end

fetch().then((data) => do
    println(data)
end)
```

### Collections
```polyloft
let list = List<Int>(1, 2, 3)
list.map((x) => x * 2)
list.filter((x) => x > 1)
list.reduce(0, (a, b) => a + b)
```

## Common Tasks

### Read User Input
```polyloft
let name = Sys.input("Enter name: ")
println("Hello, #{name}!")
```

### Time Operations
```polyloft
let start = Sys.time()
Sys.sleep(1000)
let elapsed = Sys.time() - start
```

### Error Handling
```polyloft
try:
    riskyOperation()
catch error:
    println("Error: #{error}")
finally:
    cleanup()
end
```

### Class Definition
```polyloft
class Person:
    private let name: String
    private let age: Int
    
    Person(name: String, age: Int):
        this.name = name
        this.age = age
    end
    
    def greet():
        println("Hello, I'm #{this.name}")
    end
end
```

### Generic Class
```polyloft
class Box<T>:
    private let value: T
    
    Box(value: T):
        this.value = value
    end
    
    def get() -> T:
        return this.value
    end
end
```

## Keyboard Shortcuts (REPL)

- `Ctrl+C` - Exit REPL
- `Ctrl+D` - Exit REPL (Unix)
- `Ctrl+L` - Clear screen
- `↑`/`↓` - Command history

## File Extensions

- `.pf` - Polyloft source files
- `.pfx` - Polyloft executables (Windows)
- `polyloft.toml` - Project configuration

## Environment Variables

- `POLYLOFT_REGISTRY_URL` - Package registry URL

## Common Commands

```bash
# Run a file
polyloft run file.pf

# Start REPL
polyloft repl

# Build executable
polyloft build -o myapp

# Initialize project
polyloft init

# Install dependencies
polyloft install

# Publish package
polyloft publish

# Show version
polyloft version
```

## Troubleshooting

### Installation Issues
See [Installation Guide](installation.md#troubleshooting)

### Runtime Errors
- Check syntax in [Syntax Guide](language/syntax.md)
- Review error message carefully
- Test in REPL for quick debugging

### Documentation Issues
- Open issue on [GitHub](https://github.com/ArubikU/polyloft/issues)
- Check [Contributing Guide](contributing/development.md)

## Additional Resources

- [Full Documentation](SUMMARY.md) - Complete table of contents
- [GitHub Repository](https://github.com/ArubikU/polyloft) - Source code
- [Algorithm Examples](../algorithm_samples/) - Working code samples
- [VSCode Extension](vscode-extension.md) - IDE support

## Learning Path

1. **Beginner**: Start with [Getting Started](getting-started.md)
2. **Intermediate**: Learn [Classes](language/classes-and-objects.md) and [Generics](language/generics.md)
3. **Advanced**: Explore [Concurrency](concurrency/overview.md) and [Advanced Topics](advanced/type-system.md)

## Need Help?

- 📖 Read the documentation
- 💬 Ask in [Discussions](https://github.com/ArubikU/polyloft/discussions)
- 🐛 Report bugs in [Issues](https://github.com/ArubikU/polyloft/issues)
- 🤝 Contribute via [Pull Requests](contributing/development.md)
