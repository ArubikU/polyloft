# Sys Module

The `Sys` module provides system-level utilities for time, I/O, type inspection, and program control.

## Functions

### `Sys.time(mode?)`
Returns current timestamp in milliseconds since Unix epoch.

**Parameters:**
- `mode` (optional): "float" to return Float, otherwise returns Int64

**Returns:** Int64 or Float

**Examples:**
```pf
let now = Sys.time()
println("Timestamp: #{now}")  // Outputs: 1234567890123

let nowFloat = Sys.time("float")
println("Float time: #{nowFloat}")  // Outputs: 1234567890123.0
```

### `Sys.sleep(milliseconds)`
Pauses execution for specified milliseconds.

**Parameters:**
- `milliseconds` (Int): Time to sleep

**Returns:** void

**Examples:**
```pf
println("Starting...")
Sys.sleep(1000)  // Sleep for 1 second
println("1 second later")

Sys.sleep(500)   // Sleep for 0.5 seconds
println("Done!")
```

### `Sys.random()`
Returns a random float between 0.0 and 1.0.

**Returns:** Float

**Examples:**
```pf
let rand = Sys.random()
println(rand)  // Outputs: 0.42387...

// Random integer between 0-99
let randInt = (Sys.random() * 100)
println(randInt)
```

### `Sys.seed(value)`
Sets the random number generator seed.

**Parameters:**
- `value` (Int): Seed value

**Returns:** void

**Examples:**
```pf
Sys.seed(42)
println(Sys.random())  // Always same with same seed
println(Sys.random())
```

### `Sys.type(value)`
Returns the type name of a value as a string.

**Parameters:**
- `value` (Any): Value to inspect

**Returns:** String

**Examples:**
```pf
println(Sys.type(42))           // Outputs: Int
println(Sys.type("hello"))      // Outputs: String
println(Sys.type([1, 2, 3]))    // Outputs: Array
println(Sys.type({a: 1}))       // Outputs: Map
println(Sys.type(true))         // Outputs: Bool
println(Sys.type(3.14))         // Outputs: Float
```

### `Sys.print(values...)`
Prints values without newline.

**Parameters:**
- `values...` (variadic): Values to print

**Returns:** void

**Examples:**
```pf
Sys.print("Hello")
Sys.print(" ")
Sys.print("World")  // Outputs: Hello World (no newline)
```

### `Sys.println(values...)`
Prints values with newline.

**Parameters:**
- `values...` (variadic): Values to print

**Returns:** void

**Examples:**
```pf
Sys.println("Hello", "World")  // Outputs: Hello World
Sys.println("Line 1")
Sys.println("Line 2")
```

**Note:** The global `println()` function is an alias for `Sys.println()`.

### `Sys.input(prompt?, default?, type?)`
Reads user input from console.

**Parameters:**
- `prompt` (optional String): Message to display
- `default` (optional): Default value if empty input
- `type` (optional String): Cast type ("int", "float", "bool")

**Returns:** String (or cast type)

**Examples:**
```pf
let name = Sys.input("Enter your name: ")
println("Hello, #{name}!")

let age = Sys.input("Enter age: ", 0, "int")
println("Age: #{age}")

let temp = Sys.input("Temperature: ", 0.0, "float")
println("Temp: #{temp}")
```

### `Sys.format(format, values...)`
Formats a string with placeholders.

**Parameters:**
- `format` (String): Format string with `{}` placeholders
- `values...` (variadic): Values to insert

**Returns:** String

**Examples:**
```pf
let msg = Sys.format("Hello, {}!", "Alice")
println(msg)  // Outputs: Hello, Alice!

let info = Sys.format("{} is {} years old", "Bob", 25)
println(info)  // Outputs: Bob is 25 years old
```

### `Sys.exit(args...)`
Terminates the program with error message.

**Parameters:**
- `args...` (variadic): Error message components

**Returns:** Never returns (throws exception)

**Examples:**
```pf
if criticalError:
    Sys.exit("Fatal error occurred!")
end

Sys.exit("Error code:", errorCode)
```

## Classes

### `Cronometer`
High-precision timer for measuring elapsed time.

#### Constructor
```pf
let timer = Cronometer()
```

#### Methods

**`start()`** - Start/restart the timer
```pf
timer.start()
```

**`stop()`** - Stop the timer
```pf
timer.stop()
```

**`elapsedMilliseconds()`** - Get elapsed time in milliseconds
```pf
let ms = timer.elapsedMilliseconds()
println("Elapsed: #{ms} ms")
```

**`elapsedFormatted()`** - Get formatted elapsed time (HH:MM:SS.mmm)
```pf
let formatted = timer.elapsedFormatted()
println("Time: #{formatted}")  // Outputs: 00:01:23.456
```

#### Complete Example
```pf
let timer = Cronometer()
timer.start()

// Do some work
for i in range(1000000):
    let x = i * i
end

timer.stop()
println("Operation took: #{timer.elapsedFormatted()}")
println("Milliseconds: #{timer.elapsedMilliseconds()}")
```

## Common Patterns

### Benchmarking
```pf
let start = Sys.time()
performOperation()
let end = Sys.time()
println("Took #{end - start} ms")
```

### Using Cronometer for Timing
```pf
let timer = Cronometer()
timer.start()

// Your code here
processData()

timer.stop()
println("Processing time: #{timer.elapsedFormatted()}")
```

### Type-Safe Input
```pf
let age = Sys.input("Enter age: ", 0, "int")
if age < 18:
    println("Minor")
else:
    println("Adult")
end
```

### Random Number Generation
```pf
// Random integer between min and max
def randInt(min, max):
    return min + (Sys.random() * (max - min + 1))
end

let dice = randInt(1, 6)
println("You rolled: #{dice}")
```

## See Also

- [Math Module](math.md) - Mathematical functions
- [IO Module](io.md) - File I/O operations
- [String Type](../types/string.md) - String manipulation
