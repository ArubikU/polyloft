# Hello World Examples

This page contains various "Hello World" examples demonstrating different Polyloft features.

## Basic Hello World

The simplest Polyloft program:

```polyloft
println("Hello, World!")
```

Run it:
```bash
polyloft run hello.pf
```

## Hello World with Variables

Using variables and string interpolation:

```polyloft
let greeting = "Hello"
let name = "World"
println("#{greeting}, #{name}!")
```

## Hello World Function

Defining and calling a function:

```polyloft
def greet(name: String):
    println("Hello, #{name}!")
end

greet("World")
greet("Polyloft")
```

## Hello World Class

Object-oriented approach:

```polyloft
class Greeter:
    private let greeting: String
    
    Greeter(greeting: String):
        this.greeting = greeting
    end
    
    def greet(name: String):
        println("#{this.greeting}, #{name}!")
    end
end

let greeter = Greeter("Hello")
greeter.greet("World")
```

## Hello World with Input

Interactive greeting:

```polyloft
let name = Sys.input("Enter your name: ")
println("Hello, #{name}!")
```

## Hello World to File

Writing to a file:

```polyloft
let message = "Hello, World!"
IO.writeFile("greeting.txt", message)
println("Greeting written to file!")

// Read it back
let content = IO.readFile("greeting.txt")
println("File contains: #{content}")
```

## Hello World with Concurrency

Using channels:

```polyloft
let ch = Channel<String>()

// Producer
thread spawn do
    ch.send("Hello")
    ch.send("World")
    ch.close()
end

// Consumer
loop
    let word = ch.recv()
    if word == null:
        break
    end
    println(word)
end
```

## Hello World with Collections

Using functional programming:

```polyloft
let words = ["Hello", "Beautiful", "World"]

// Map to uppercase
let upper = words.map((w: String) -> String => w.toUpper())

// Join with spaces
let greeting = upper.join(" ")
println(greeting + "!")
```

## Hello World with Generics

Generic greeting function:

```polyloft
class Message<T>:
    private let content: T
    
    Message(content: T):
        this.content = content
    end
    
    def display():
        println(this.content)
    end
end

let strMsg = Message<String>("Hello, World!")
let intMsg = Message<Int>(42)

strMsg.display()  // Hello, World!
intMsg.display()  // 42
```

## Next Steps

Now that you've seen Hello World in many forms, explore:

- [Getting Started Guide](../getting-started.md) - Learn more basics
- [Language Syntax](../language/syntax.md) - Complete syntax reference
- [Standard Library](../stdlib/overview.md) - Explore built-in modules
- [Examples](algorithms.md) - More complex examples

Try running these examples yourself:

```bash
# Save any example to a file
echo 'println("Hello, World!")' > hello.pf

# Run it
polyloft run hello.pf

# Or use the REPL
polyloft repl
>>> println("Hello, World!")
```
