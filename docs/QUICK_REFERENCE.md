# Polyloft Quick Reference

Complete cheat sheet for Polyloft language features.

## Variables

```pf
var x = 10           // Mutable type
let y = 20           // Immutable type (recommended)
const PI = 3.14159   // Compile-time constant
final MAX = 100      // Runtime constant
```

## Data Types

```pf
// Primitives
let num = 42              // Int
let flt = 3.14            // Float
let str = "hello"         // String
let bool = true           // Bool

// Collections
let arr = [1, 2, 3]       // Array
let map = {a: 1, b: 2}    // Map
let set = Set(1, 2, 3)    // Set

// Special
let nothing = nil         // Nil
```

## String Interpolation

```pf
let name = "Alice"
let age = 25
println("My name is #{name} and I'm #{age}")
```

## Functions

```pf
// Basic function
def greet(name):
    return "Hello, #{name}"
end

// With types
def add(a: Int, b: Int): Int
    return a + b
end

// Lambda
let double = (x) => x * 2
```

## Classes

```pf
class Person:
    let name: String
    let age: Int
    
    Person(name: String, age: Int):
        this.name = name
        this.age = age
    end
    
    def greet():
        println("Hi, I'm #{this.name}")
    end
end

let p = Person("Alice", 25)
```

### Access Modifiers
```pf
class BankAccount:
    public let accountId        // Accessible everywhere
    private var balance         // Only in this class
    protected var lastAccess    // In class and subclasses
    
    BankAccount(id, initialBalance):
        this.accountId = id
        this.balance = initialBalance
        this.lastAccess = Sys.time()
    end
    
    public def getBalance() -> Float
        return this.balance
    end
    
    private def updateAccess():
        this.lastAccess = Sys.time()
    end
end

// Static members
class Config:
    static var appName = "MyApp"
    
    static def getSetting(key) -> any
        return settings.get(key)
    end
end
```

### Annotations
```pf
class Animal:
    def speak():
        println("Some sound")
    end
end

class Dog < Animal:
    @Override
    def speak():
        println("Woof!")
    end
end
```

### Interfaces
```pf
interface Drawable:
    draw() -> void
end

class Circle implements Drawable:
    let radius
    
    Circle(r):
        this.radius = r
    end
    
    @Override
    def draw():
        println("Drawing circle")
    end
end
```

## Control Flow

### If-Else
```pf
if condition:
    // code
else:
    if other:
        // code
    else:
        // code
    end
end
```

### For Loop
```pf
for i in range(10):
    println(i)
end

for item in array:
    println(item)
end
```

### Loop (while-like)
```pf
loop count < 10:
    println(count)
    count = count + 1
end

// Infinite loop
loop:
    if shouldStop:
        break
    end
end
```

### Do-Loop
```pf
do:
    println("At least once")
loop condition
```

### Switch
```pf
switch value:
    case 1:
        doOne()
    case 2:
        doTwo()
    default:
        doOther()
end
```

## Exception Handling

```pf
// Basic try-catch
try:
    riskyOperation()
catch e:
    println("Error: #{e}")
end

// Typed catch
try:
    processFile("data.txt")
catch e: FileNotFoundException:
    println("File not found")
catch e: IOException:
    println("IO error: #{e}")
end

// Throw exceptions
throw "Error message"

// Custom error types
class ValidationError < Exception:
    let field
    let message
    
    ValidationError(fieldName, msg):
        this.field = fieldName
        this.message = msg
    end
end

throw ValidationError("email", "Invalid format")
```

## Enums

```pf
enum Color
    RED
    GREEN
    BLUE
end

let c = Color.RED
println(c.name)      // "RED"
println(c.ordinal)   // 0
```

### Enum with Values
```pf
enum Planet
    EARTH(9.8)
    MARS(3.7)
    
    var gravity: Float
    
    Planet(g: Float):
        this.gravity = g
    end
end
```

## Records

```pf
record Point(x: Int, y: Int)
    def distance():
        return Math.sqrt(this.x * this.x + this.y * this.y)
    end
end

let p = Point(3, 4)
println(p.distance())  // 5.0
```

## Async/Await

```pf
// Create promise
let promise = async(() => fetchData())

// Wait for result
let result = promise.await()

// Chaining
async(() => 5)
    .then((val) => val * 2)
    .then((val) => println(val))

// CompletableFuture
let future = CompletableFuture()
thread spawn do
    future.complete(42)
end
let result = future.get()
```

## Standard Library

### Sys Module
```pf
Sys.time()              // Current timestamp (ms)
Sys.sleep(1000)         // Sleep 1 second
Sys.random()            // Random 0.0-1.0
Sys.type(value)         // Get type name
println("Hello")        // Print line
Sys.input("Prompt: ")   // Get user input
```

### Math Module
```pf
Math.PI                 // 3.14159...
Math.E                  // 2.71828...
Math.abs(-5)            // 5
Math.floor(3.7)         // 3
Math.ceil(3.2)          // 4
Math.round(3.5)         // 4
Math.sqrt(16)           // 4
Math.pow(2, 3)          // 8
Math.sin(0)             // 0
Math.cos(0)             // 1
Math.min(5, 3)          // 3
Math.max(5, 3)          // 5
```

### IO Module
```pf
IO.readFile("file.txt")
IO.writeFile("file.txt", content)
IO.appendFile("log.txt", line)
IO.exists("file.txt")
IO.delete("file.txt")
IO.readDir(".")
IO.mkdir("newdir")
```

### Http Module
```pf
// Client
let res = Http.get("https://api.example.com")
Http.post(url, data)
Http.put(url, data)
Http.delete(url)

// Server
let server = Http.createServer()
server.get("/users", (req, res) => do
    res.json({users: []})
end)
server.listen(8080)
```

## Array Methods

```pf
arr.length()            // Size
arr.push(item)          // Add to end
arr.pop()               // Remove from end
arr.shift()             // Remove from start
arr.unshift(item)       // Add to start
arr.indexOf(item)       // Find index
arr.contains(item)      // Check existence
arr.slice(start, end)   // Extract slice
arr.reverse()           // Reverse in place
arr.sort()              // Sort in place
arr.join(", ")          // Join to string

// Functional
arr.map((x) => x * 2)
arr.filter((x) => x > 5)
arr.reduce((acc, x) => acc + x, 0)
arr.find((x) => x > 5)
arr.every((x) => x > 0)
arr.some((x) => x > 10)
```

## String Methods

```pf
str.length()            // Length
str.charAt(0)           // Character at index
str.indexOf("sub")      // Find substring
str.substring(0, 5)     // Extract substring
str.toUpperCase()       // UPPERCASE
str.toLowerCase()       // lowercase
str.trim()              // Remove whitespace
str.startsWith("pre")   // Check prefix
str.endsWith("suf")     // Check suffix
str.contains("mid")     // Check contains
str.replace("a", "b")   // Replace all
str.split(",")          // Split to array
```

## Map Methods

```pf
map.get("key")          // Get value
map.set("key", value)   // Set value
map.has("key")          // Check key exists
map.remove("key")       // Remove key
map.keys()              // Array of keys
map.values()            // Array of values
map.entries()           // Array of [k,v] pairs
map.size()              // Number of entries
map.isEmpty()           // Check if empty
map.clear()             // Remove all
```

## Set Methods

```pf
set.add(item)           // Add item
set.remove(item)        // Remove item
set.has(item)           // Check existence
set.size()              // Number of items
set.isEmpty()           // Check if empty
set.toArray()           // Convert to array
set.clear()             // Remove all
```

## Operators

### Arithmetic
```pf
+  -  *  /  %           // Basic math
**                      // Exponentiation
```

### Comparison
```pf
==  !=                  // Equality
<  >  <=  >=           // Relational
```

### Logical
```pf
and  or  not           // Boolean logic
```

### Other
```pf
in                      // Membership
instanceof              // Type check
```

## Range

```pf
range(10)               // 0 to 9
range(1, 11)            // 1 to 10
range(0, 10, 2)         // 0, 2, 4, 6, 8
1...10                  // 1 to 10 (inclusive)
```

## Type Checking

```pf
Sys.type(value)         // Get type name
value instanceof Type   // Check instance
```

## Common Patterns

### Array Processing
```pf
// Sum
let sum = numbers.reduce((acc, n) => acc + n, 0)

// Filter and map
let result = numbers
    .filter((x) => x > 0)
    .map((x) => x * 2)

// Find
let item = items.find((i) => i.id == targetId)
```

### Error Handling
```pf
try:
    let result = operation()
    return result
catch e:
    println("Error: #{e}")
    return defaultValue
end
```

### Looping
```pf
// With index
for i in range(arr.length()):
    println("#{i}: #{arr[i]}")
end

// Break/continue
for item in items:
    if skip(item):
        continue
    end
    if done(item):
        break
    end
    process(item)
end
```

### Class Pattern
```pf
class Counter:
    var count: Int
    
    Counter():
        this.count = 0
    end
    
    def increment():
        this.count = this.count + 1
    end
    
    def getValue() -> Int
        return this.count
    end
end
```

### Async Pattern
```pf
let promises = []
for item in items:
    promises = promises.concat([async(() => process(item))])
end

for p in promises:
    let result = p.await()
    println(result)
end
```

## Comments

```pf
// Single line comment

/*
Multi-line
comment
*/
```

## Keywords

```
var let const final
def class record enum interface
if else switch case default
for loop do end break continue
try catch throw
async await
true false nil
this super
return implements
public private protected static
sealed abstract
and or not in instanceof
```

## Annotations

```
@Override    // Mark method overrides
// More annotations can be added in the future
```

## Best Practices

1. **Use `let` by default** - More type-safe than `var`
2. **Handle errors** - Always catch exceptions for I/O operations
3. **Validate input** - Check user input and external data
4. **Use descriptive names** - Make code self-documenting
5. **Keep functions small** - Single responsibility principle
6. **Prefer immutability** - Use records for data objects
7. **Document public APIs** - Add comments for complex logic
8. **Test edge cases** - Handle nil, empty, and boundary values

## See Also

- [Full Documentation](README.md)
- [Variables Guide](variables/README.md)
- [Control Flow](control-flow/README.md)
- [Standard Library](stdlib/README.md)
