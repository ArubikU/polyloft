# Functions

Functions in Polyloft are first-class citizens that can be defined, passed around, and invoked.

## Syntax

### Basic Function
```pf
def functionName(param1, param2):
    // function body
    return result
end
```

### With Type Annotations
```pf
def functionName(param1: Type1, param2: Type2): ReturnType
    return result
end
```

## Examples

### Simple Function
```pf
def greet(name):
    return "Hello, " + name
end

println(greet("Alice"))  // Outputs: Hello, Alice
```

### With Type Annotations
```pf
def add(a: Int, b: Int): Int
    return a + b
end

let sum = add(5, 3)
println(sum)  // Outputs: 8
```

### Multiple Parameters
```pf
def calculateArea(width: Float, height: Float): Float
    return width * height
end

let area = calculateArea(10.5, 20.3)
println("Area: #{area}")
```

### No Return Value (void)
```pf
def printMessage(msg: String): void
    println("Message: #{msg}")
end

printMessage("Hello World")
```

### Default Parameters
```pf
def greet(name, greeting="Hello"):
    return greeting + ", " + name
end

println(greet("Alice"))           // Hello, Alice
println(greet("Bob", "Hi"))       // Hi, Bob
```

### Variadic Parameters
```pf
def sum(numbers...):
    let total = 0
    for num in numbers:
        total = total + num
    end
    return total
end

println(sum(1, 2, 3))        // 6
println(sum(1, 2, 3, 4, 5))  // 15
```

### Nested Functions
```pf
def outer(x):
    def inner(y):
        return x + y
    end
    return inner(10)
end

println(outer(5))  // 15
```

### Functions as Values
```pf
def double(x):
    return x * 2
end

let f = double
println(f(5))  // 10

def apply(fn, value):
    return fn(value)
end

println(apply(double, 7))  // 14
```

### Anonymous Functions (Lambda)
```pf
let multiply = (x, y) => x * y
println(multiply(4, 5))  // 20

let numbers = [1, 2, 3, 4, 5]
let doubled = numbers.map((x) => x * 2)
println(doubled)  // [2, 4, 6, 8, 10]
```

### Recursive Functions
```pf
def factorial(n: Int): Int
    if n <= 1:
        return 1
    end
    return n * factorial(n - 1)
end

println(factorial(5))  // 120
```

```pf
def fibonacci(n: Int): Int
    if n <= 1:
        return n
    end
    return fibonacci(n - 1) + fibonacci(n - 2)
end

println(fibonacci(7))  // 13
```

### Higher-Order Functions
```pf
def makeMultiplier(factor):
    return (x) => x * factor
end

let double = makeMultiplier(2)
let triple = makeMultiplier(3)

println(double(5))  // 10
println(triple(5))  // 15
```

## Function Types

### Pure Functions
Functions with no side effects:
```pf
def add(a, b):
    return a + b  // Only returns a value
end
```

### Impure Functions
Functions with side effects:
```pf
def logAndAdd(a, b):
    println("Adding #{a} + #{b}")  // Side effect: I/O
    return a + b
end
```

## Best Practices

### ✅ DO - Use descriptive names
```pf
def calculateTotalPrice(items, taxRate):
    // ...
end
```

### ✅ DO - Keep functions small and focused
```pf
def validateEmail(email):
    return email.contains("@") and email.contains(".")
end

def sendEmail(to, subject, body):
    if not validateEmail(to):
        return false
    end
    // Send email logic
    return true
end
```

### ✅ DO - Add type annotations for clarity
```pf
def processUser(id: String, age: Int): Bool
    // ...
end
```

### ❌ DON'T - Create overly complex functions
```pf
// Bad: Too many responsibilities
def doEverything(data):
    let result = validateData(data)
    let processed = processData(result)
    let saved = saveToDatabase(processed)
    let email = sendNotification(saved)
    return email
end

// Better: Split into smaller functions
def processAndSaveData(data):
    let validated = validateData(data)
    let processed = processData(validated)
    return saveToDatabase(processed)
end
```

## Common Patterns

### Guard Clauses
```pf
def divide(a, b):
    if b == 0:
        return nil  // Early return
    end
    return a / b
end
```

### Factory Functions
```pf
def createPoint(x, y):
    return {x: x, y: y}
end

let p1 = createPoint(10, 20)
```

### Callback Pattern
```pf
def fetchData(url, onSuccess, onError):
    try:
        let data = makeRequest(url)
        onSuccess(data)
    catch e:
        onError(e)
    end
end

fetchData(
    "http://api.example.com/data",
    (data) => println("Success: #{data}"),
    (err) => println("Error: #{err}")
)
```

### Memoization
```pf
let cache = {}

def fibonacci(n):
    if n in cache:
        return cache[n]
    end
    
    if n <= 1:
        return n
    end
    
    let result = fibonacci(n - 1) + fibonacci(n - 2)
    cache[n] = result
    return result
end
```

## See Also

- [Classes](class.md)
- [Lambdas & Closures](../advanced/closures.md)
- [Generics](../advanced/generics.md)
