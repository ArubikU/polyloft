# Constant Declaration: `const`

## Overview

`const` declares a **compile-time constant**. The value is evaluated at compile time and cannot be changed.

## Syntax

```pf
const IDENTIFIER = value
```

## Key Features

- ❌ Value cannot be reassigned
- ❌ Type cannot change
- ✅ Evaluated at compile time
- ✅ Must be a constant expression
- ✅ Convention: Use UPPER_CASE names

## Examples

### Basic Usage
```pf
const PI = 3.14159
const MAX_SIZE = 100
const APP_NAME = "Polyloft"

println(PI)  // Outputs: 3.14159
// PI = 3.14  // Error: cannot reassign constant
```

### Compile-Time Expressions
```pf
const SECONDS_PER_MINUTE = 60
const MINUTES_PER_HOUR = 60
const SECONDS_PER_HOUR = SECONDS_PER_MINUTE * MINUTES_PER_HOUR

println(SECONDS_PER_HOUR)  // Outputs: 3600
```

### Type Inference
```pf
const MAX_VALUE = 100        // Type: Int
const GREETING = "Hello"     // Type: String
const ENABLED = true         // Type: Bool
const RATE = 1.5             // Type: Float
```

### Common Use Cases
```pf
// Configuration
const DEFAULT_PORT = 8080
const DEFAULT_HOST = "localhost"
const TIMEOUT_MS = 5000

// Mathematical constants
const PI = 3.14159265359
const E = 2.71828182846
const GOLDEN_RATIO = 1.61803398875

// Application constants
const VERSION = "1.0.0"
const BUILD_DATE = "2024-01-01"
const AUTHOR = "Polyloft Team"
```

## What Can Be `const`

### ✅ Allowed (Compile-time values)
```pf
const NUM = 42
const TEXT = "hello"
const BOOL = true
const CALC = 10 * 5 + 2
```

### ❌ Not Allowed (Runtime values)
```pf
// const TIME = Sys.time()          // Error: runtime function
// const INPUT = Sys.input()        // Error: runtime function
// const RANDOM = Math.random()     // Error: runtime function
```

## Comparison with Other Variable Types

| Feature | `var` | `let` | `const` | `final` |
|---------|-------|-------|---------|---------|
| Value mutable | ✅ Yes | ✅ Yes | ❌ No | ❌ No |
| Type mutable | ✅ Yes | ❌ No | ❌ No | ❌ No |
| Compile-time | ❌ No | ❌ No | ✅ Yes | ❌ No |
| Type safety | ❌ No | ✅ Yes | ✅ Yes | ✅ Yes |
| Use case | Dynamic | Variables | Config | Computed |

## Best Practices

### ✅ DO - Use for true constants
```pf
const PI = 3.14159
const MAX_RETRIES = 3
const ERROR_MESSAGE = "Operation failed"
```

### ✅ DO - Use UPPER_CASE naming
```pf
const DEFAULT_TIMEOUT = 5000    // Good
const API_ENDPOINT = "/api/v1"  // Good
```

### ❌ DON'T - Use for runtime values
```pf
// Bad: This should be 'final'
// const CURRENT_TIME = Sys.time()

// Better:
final CURRENT_TIME = Sys.time()
```

## See Also

- [`let` - Mutable variable](let.md)
- [`final` - Runtime constant](final.md)
- [`var` - Dynamic variable](var.md)
