# Language Overview

Polyloft is a modern, hybrid programming language that combines object-oriented programming with scripting flexibility. This overview introduces the core concepts and philosophy of the language.

## Design Philosophy

Polyloft is built on these principles:

1. **Readability First**: Code should be clear and self-documenting
2. **Type Safety with Flexibility**: Strong typing when needed, inference when helpful
3. **Modern Features**: Contemporary programming patterns and constructs
4. **Educational Value**: Perfect for learning computer science concepts
5. **Practical Application**: Suitable for real-world projects

## Language Paradigms

Polyloft supports multiple programming paradigms:

### Object-Oriented Programming

```polyloft
class Animal:
    Animal(name: String):
        this.name = name
    end
    
    def speak():
        println("#{this.name} makes a sound")
    end
end

class Dog extends Animal:
    Dog(name: String):
        super(name)
    end
    
    def speak():
        println("#{this.name} barks")
    end
end
```

### Functional Programming

```polyloft
// Higher-order functions
let numbers = [1, 2, 3, 4, 5]
let squared = numbers.map((x: Int) -> Int => x * x)
let sum = squared.reduce(0, (a: Int, b: Int) -> Int => a + b)

// Function composition
def compose(f, g):
    return (x) => f(g(x))
end
```

### Procedural Programming

```polyloft
def calculateFactorial(n: Int) -> Int:
    let result = 1
    for i in range(1, n + 1):
        result = result * i
    end
    return result
end
```

## Core Features

### 1. Type System

Polyloft has a strong, gradual type system:

- **Type Inference**: Automatically deduces types
- **Explicit Typing**: Optional type annotations for clarity
- **Generic Types**: Parameterized types with bounds
- **Variance**: Covariance and contravariance support

```polyloft
// Type inference
let x = 42              // Int
let y = 3.14           // Float

// Explicit typing
let name: String = "Alice"
let age: Int = 30

// Generics
let list: List<Int> = List(1, 2, 3)
let map: Map<String, Int> = Map()
```

### 2. Collections

Rich collection types with functional methods:

```polyloft
// List
let numbers = [1, 2, 3, 4, 5]
numbers.map((x) => x * 2)
numbers.filter((x) => x % 2 == 0)
numbers.reduce(0, (a, b) => a + b)

// Map
let scores = Map<String, Int>()
scores.set("Alice", 95)
scores.get("Alice")

// Set
let unique = Set<String>("apple", "banana", "apple")
```

### 3. Functions

Multiple ways to define functions:

```polyloft
// Named function
def add(a: Int, b: Int) -> Int:
    return a + b
end

// Lambda expression
let multiply = (x: Int, y: Int) -> Int => x * y

// Closure
def makeCounter():
    let count = 0
    return () => {
        count = count + 1
        return count
    }
end
```

### 4. Classes and Interfaces

Full object-oriented support:

```polyloft
interface Drawable:
    def draw(): Void
end

class Circle implements Drawable:
    private let radius: Float
    
    Circle(radius: Float):
        this.radius = radius
    end
    
    def draw():
        println("Drawing circle with radius #{this.radius}")
    end
    
    def area() -> Float:
        return Math.PI * this.radius * this.radius
    end
end
```

### 5. Generics

Type-safe parameterized code:

```polyloft
class Box<T>:
    private let value: T
    
    Box(value: T):
        this.value = value
    end
    
    def get() -> T:
        return this.value
    end
    
    def set(value: T):
        this.value = value
    end
end

let intBox = Box<Int>(42)
let strBox = Box<String>("Hello")
```

### 6. Concurrency

Modern concurrency primitives:

```polyloft
// Async/Await with Promises
let promise = Promise((resolve, reject) => do
    let response = Http.get(url)
    resolve(response.body)
end)

promise.then((body) => do
    println("Response: #{body}")
end).catch((error) => do
    println("Error: #{error}")
end)

// Channels
let ch = Channel<Int>()
thread spawn do
    ch.send(42)
end
let value = ch.recv()

// Threads
let t = thread spawn do
    return performComputation()
end
let result = thread join t
```

### 7. Error Handling

Structured exception handling:

```polyloft
try:
    let file = IO.readFile("data.txt")
    let parsed = JSON.parse(file)
    processData(parsed)
catch error:
    println("Error: #{error}")
finally:
    cleanup()
end

// Custom exceptions
class ValidationError extends Exception:
    ValidationError(message: String):
        super(message)
    end
end

throw ValidationError("Invalid input")
```

### 8. Pattern Matching

Select statements for complex control flow:

```polyloft
select:
    case ch1.recv() as value:
        println("Received from ch1: #{value}")
    case ch2.send(42):
        println("Sent to ch2")
    default:
        println("No channel ready")
end
```

## Language Characteristics

### Static and Dynamic Features

Polyloft balances static and dynamic typing:

- **Compile-time checking** for type safety
- **Runtime reflection** for flexibility
- **Type inference** reduces verbosity
- **Optional typing** for quick prototyping

### Memory Management

- **Automatic garbage collection** (Go's GC)
- **No manual memory management** required
- **Efficient** for most use cases

### Performance

- **Compiled to native code** via Go backend
- **Fast startup time** for scripts
- **Optimized standard library**
- **Concurrent execution** with goroutines

## Syntax Highlights

### Indentation-Aware

Polyloft uses `end` keywords instead of braces:

```polyloft
if condition:
    doSomething()
    doSomethingElse()
end
```

### String Interpolation

Embedded expressions in strings:

```polyloft
let name = "Alice"
let age = 30
println("#{name} is #{age} years old")
```

### Range Expressions

Convenient iteration:

```polyloft
for i in range(10):      // 0 to 9
    println(i)
end

for i in range(5, 10):   // 5 to 9
    println(i)
end
```

### Defer Statement

Cleanup code execution:

```polyloft
def processFile(filename: String):
    let file = IO.openFile(filename)
    defer file.close()
    
    // Work with file
    // file.close() will be called when function returns
end
```

## Standard Library

Comprehensive built-in modules:

- **Math**: Mathematical functions and constants
- **Sys**: System operations and utilities
- **IO**: File and stream I/O
- **Net**: Network programming
- **String**: String manipulation
- **Collections**: List, Set, Map, Deque with rich methods

## Development Workflow

```bash
# Write code
echo 'println("Hello, World!")' > hello.pf

# Run immediately
polyloft run hello.pf

# Build executable
polyloft build -o hello

# Run executable
./hello
```

## Package System

```toml
[project]
name = "my-project"
version = "1.0.0"
entry_point = "src/main.pf"

[[dependencies.pf]]
name = "utils@author"
version = "1.0.0"
```

## Next Steps

Explore specific language features:

- [Basic Syntax](syntax.md) - Language syntax rules
- [Variables and Types](variables-and-types.md) - Type system details
- [Functions](functions.md) - Function definitions and usage
- [Classes and Objects](classes-and-objects.md) - OOP in Polyloft
- [Generics](generics.md) - Generic programming

Ready to learn more? Continue to [Basic Syntax](syntax.md)!
