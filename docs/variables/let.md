# Variable Declaration: `let`

## Overview

`let` declares a **type-safe mutable variable**. The value can change, but the type is fixed after initialization.

## Syntax

```pf
let identifier = value            // Type inferred
let identifier: Type = value      // Type explicit
```

## Key Features

- ✅ Value can be reassigned
- ❌ Type cannot change (type-safe)
- ✅ Type inferred from first assignment
- ✅ **Recommended for most cases**

## Examples

### Basic Usage
```pf
let x = 10          // x is Int
x = 20              // OK - still Int
// x = "Hello"      // Error: cannot assign String to Int
```

### Type Inference
```pf
let name = "Alice"          // Inferred as String
let age = 25                // Inferred as Int
let height = 5.9            // Inferred as Float
let active = true           // Inferred as Bool
let items = [1, 2, 3]       // Inferred as Array
let data = {key: "value"}   // Inferred as Map
```

### Explicit Type Annotation
```pf
let count: Int = 0
let message: String = "Hello"
let temperature: Float = 98.6
```

### Reassignment Examples
```pf
let counter = 0
counter = counter + 1      // OK
counter = counter * 2      // OK
println(counter)           // Outputs: 2

let greeting = "Hello"
greeting = "Hi"            // OK
greeting = greeting + "!"  // OK
// greeting = 123          // Error: type mismatch
```

### Working with Collections
```pf
let numbers = [1, 2, 3]
numbers = [4, 5, 6]           // OK - same type
numbers = numbers.concat([7]) // OK - returns Array
// numbers = "not array"      // Error: type mismatch

let person = {name: "Alice", age: 25}
person = {name: "Bob", age: 30}       // OK
// person = [1, 2, 3]                 // Error: type mismatch
```

### Loop Variables
```pf
let sum = 0
for i in range(10):
    sum = sum + i
end
println(sum)  // Outputs: 45

let result = ""
for item in ["a", "b", "c"]:
    result = result + item
end
println(result)  // Outputs: abc
```

## Best Practices

### ✅ DO - Use `let` for most variables
```pf
let count = 0
let name = "Alice"
let items = []

for i in range(10):
    count = count + i
    items = items.concat([i])
end
```

### ✅ DO - Let the compiler infer types
```pf
let username = "alice"           // Good: clear type
let score = 100                  // Good: clear type
let factors = [1, 2, 3, 5]       // Good: clear type
```

## Common Patterns

### Accumulator Pattern
```pf
let total = 0
for num in numbers:
    total = total + num
end
```

### Builder Pattern
```pf
let query = "SELECT * FROM users"
if hasFilter:
    query = query + " WHERE active = true"
end
if hasLimit:
    query = query + " LIMIT 10"
end
```

## See Also

- [`var` - Dynamic typing](var.md)
- [`const` - Compile-time constant](const.md)
- [`final` - Runtime constant](final.md)
