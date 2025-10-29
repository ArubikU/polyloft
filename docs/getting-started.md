# Getting Started with Polyloft

This guide will help you write your first Polyloft program and understand the basics of the language.

## Your First Program

Create a file called `hello.pf`:

```polyloft
println("Hello, Polyloft!")
```

Run it with:

```bash
polyloft run hello.pf
```

You should see:
```
Hello, Polyloft!
```

## Basic Concepts

### 1. Variables

Polyloft has three types of variable declarations:

```polyloft
// Mutable variable with type inference
let x = 42

// Mutable variable with explicit type
let y: Int = 100

// Constant (immutable)
const PI = 3.14159

// Variable without initialization (requires type)
var name: String
name = "Alice"
```

### 2. Basic Types

```polyloft
// Numbers
let integer = 42           // Int
let float = 3.14          // Float
let double = 2.71828      // Double

// Strings
let text = "Hello"
let interpolated = "Value: #{integer}"

// Booleans
let isTrue = true
let isFalse = false

// Arrays
let numbers = [1, 2, 3, 4, 5]
let mixed = ["text", 42, true]
```

### 3. Functions

```polyloft
// Simple function
def greet(name: String):
    println("Hello, #{name}!")
end

greet("World")

// Function with return type
def add(a: Int, b: Int) -> Int:
    return a + b
end

let sum = add(5, 3)  // 8

// Lambda expressions
let multiply = (x: Int, y: Int) -> Int => x * y
println(multiply(4, 5))  // 20
```

### 4. Control Flow

```polyloft
// If-else
let age = 18
if age >= 18:
    println("Adult")
elif age >= 13:
    println("Teenager")
else:
    println("Child")
end

// For loops
for i in range(5):
    println(i)
end

// While loops
let count = 0
while count < 5:
    println(count)
    count = count + 1
end

// For-in loops
let fruits = ["apple", "banana", "orange"]
for fruit in fruits:
    println(fruit)
end
```

### 5. Classes

```polyloft
class Person:
    Person(name: String, age: Int):
        this.name = name
        this.age = age
    end
    
    def introduce():
        println("I'm #{this.name}, #{this.age} years old")
    end
    
    def birthday():
        this.age = this.age + 1
    end
end

let person = Person("Alice", 25)
person.introduce()  // I'm Alice, 25 years old
person.birthday()
person.introduce()  // I'm Alice, 26 years old
```

### 6. Lists and Collections

```polyloft
// Lists
let numbers = [1, 2, 3, 4, 5]

// Map
let doubled = numbers.map((x: Int) -> Int => x * 2)
println(doubled)  // [2, 4, 6, 8, 10]

// Filter
let evens = numbers.filter((x: Int) -> Boolean => x % 2 == 0)
println(evens)  // [2, 4]

// Reduce
let sum = numbers.reduce(0, (acc: Int, x: Int) -> Int => acc + x)
println(sum)  // 15
```

## Next Steps

Now that you understand the basics, explore:

- [Language Overview](language/overview.md) - Complete language reference
- [Standard Library](stdlib/overview.md) - Built-in modules and functions
- [CLI Tools](cli/overview.md) - Development tools
- [Examples](examples/hello-world.md) - More code examples

## Interactive Learning

Use the REPL to experiment:

```bash
polyloft repl
```

Try these commands:

```polyloft
>>> let x = 10
>>> let y = 20
>>> x + y
30
>>> def square(n: Int) -> Int: return n * n end
>>> square(7)
49
>>> [1, 2, 3].map((x: Int) -> Int => x * x)
[1, 4, 9]
```

## Common Patterns

### String Interpolation

```polyloft
let name = "Alice"
let age = 30
println("Name: #{name}, Age: #{age}")
```

### Error Handling

```polyloft
try:
    let result = riskyOperation()
    println("Success: #{result}")
catch error:
    println("Error: #{error}")
finally:
    println("Cleanup")
end
```

### Working with Files

```polyloft
// Read file
let content = IO.readFile("data.txt")
println(content)

// Write file
IO.writeFile("output.txt", "Hello, World!")
```

### Math Operations

```polyloft
let x = Math.sqrt(16)      // 4.0
let y = Math.pow(2, 8)     // 256.0
let z = Math.sin(Math.PI)  // 0.0 (approximately)
```

## Project Structure

For larger projects, create a `polyloft.toml` file:

```toml
[project]
name = "my-app"
version = "1.0.0"
entry_point = "src/main.pf"

[[dependencies.pf]]
name = "some-package@author"
version = "1.0.0"
```

Then use:

```bash
polyloft init        # Initialize project
polyloft install     # Install dependencies
polyloft run         # Run project
polyloft build       # Build executable
```

## Tips for Beginners

1. **Use the REPL**: Great for testing small code snippets
2. **Check types**: Use explicit types when unsure
3. **Read error messages**: They're usually helpful
4. **Start simple**: Begin with basic scripts before complex programs
5. **Use println**: Debug by printing values
6. **Explore examples**: Check `algorithm_samples/` directory

## Getting Help

- **Documentation**: Read this guide thoroughly
- **Examples**: Study code in `algorithm_samples/`
- **REPL**: Experiment interactively
- **Issues**: Report problems on GitHub
- **Community**: Ask questions in discussions

Ready to dive deeper? Continue to the [Language Overview](language/overview.md)!
