# Standard Library Overview

The Polyloft standard library provides a comprehensive set of modules for common programming tasks.

## Core Modules

### [Math Module](math.md)
Mathematical functions and constants for numerical computations.
- Trigonometric functions (sin, cos, tan)
- Exponential and logarithmic functions
- Constants (PI, E)
- Rounding and absolute values
- Min/max operations

### [Sys Module](sys.md)
System-level operations and utilities.
- Time and timing functions
- Random number generation
- User input
- Process control
- Sleep and delays

### [IO Module](io.md)
File and stream I/O operations.
- File reading and writing
- Directory operations
- Path manipulation
- File information and statistics

### [Net Module](net.md)
Network programming with TCP sockets.
- TCP server creation
- Client connections
- Data transmission
- Connection management

## Built-in Types

### [Collections](collections.md)
Generic collection types with rich methods.
- List<T> - Dynamic arrays
- Set<T> - Unique elements
- Map<K, V> - Key-value pairs
- Deque<T> - Double-ended queues

### [String Methods](string.md)
String manipulation and processing.
- Substring operations
- Search and replace
- Case conversion
- Splitting and joining

### [Number Methods](number.md)
Numeric type methods and conversions.
- Type conversions
- Formatting
- Parsing

### [Array Methods](array.md)
Array operations and transformations.
- Map, filter, reduce
- Sorting and searching
- Array manipulation

## Usage Examples

```polyloft
// Math operations
let area = Math.PI * radius * radius
let hypotenuse = Math.sqrt(a * a + b * b)

// System operations
let timestamp = Sys.time()
Sys.sleep(1000)
let random = Sys.random()

// File I/O
let content = IO.readFile("data.txt")
IO.writeFile("output.txt", result)

// Networking
let server = Net.listen(":8080")
let conn = server.accept()
conn.send("Hello!\n")

// Collections
let numbers = List<Int>(1, 2, 3, 4, 5)
let doubled = numbers.map((x) => x * 2)
let evens = numbers.filter((x) => x % 2 == 0)
```

## Module Organization

All standard library modules are:
- **Static classes**: Access via `Module.function()`
- **Type-safe**: Strong typing for all operations
- **Well-documented**: Examples and descriptions
- **Tested**: Comprehensive test coverage

## Importing

Standard library modules are automatically available:

```polyloft
// No import needed
println(Math.PI)
Sys.sleep(100)
IO.readFile("data.txt")
```

For third-party packages, use import statements:

```polyloft
import package from "package@author"
```

## See Also

- [Math Module](math.md)
- [Sys Module](sys.md)
- [IO Module](io.md)
- [Net Module](net.md)
- [Collections](collections.md)
