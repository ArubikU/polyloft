# Variable Declaration: `var`

## Overview

`var` declares a **mutable variable** with a **mutable type**. This means both the value AND the type can change during execution.

## Syntax

```pf
var identifier = value
var identifier: Type = value
```

## Key Features

- ✅ Value can be reassigned
- ✅ Type can change
- ❌ No type safety enforcement
- ⚠️ Use with caution - can lead to runtime errors

## Examples

### Basic Usage
```pf
var x = 10          // x is Int
x = 20              // Still Int - OK
x = "Hello"         // Now String - Also OK!
println(x)          // Outputs: Hello
```

### Type Changes
```pf
var data = 42
println(Sys.type(data))  // Outputs: Int

data = "text"
println(Sys.type(data))  // Outputs: String

data = [1, 2, 3]
println(Sys.type(data))  // Outputs: Array
```

### With Explicit Type Annotation
```pf
// When you specify a type, it acts like 'let'
var count: Int = 0
count = 10          // OK
// count = "text"   // Error: type mismatch
```

### Common Use Cases
```pf
// 1. Dynamic data handling
var result = processData()
if result == nil:
    result = "default"
end

// 2. Flexible configuration
var config = loadConfig()
if config == nil:
    config = getDefaultConfig()
end

// 3. State machines
var state = "initial"
loop state != "done":
    if state == "initial":
        state = "processing"
    else:
        state = "done"
    end
end
```

## Comparison with Other Variable Types

| Feature | `var` | `let` | `const` | `final` |
|---------|-------|-------|---------|---------|
| Value mutable | ✅ Yes | ✅ Yes | ❌ No | ❌ No |
| Type mutable | ✅ Yes | ❌ No | ❌ No | ❌ No |
| Compile-time | ❌ No | ❌ No | ✅ Yes | ❌ No |
| Type safety | ❌ No | ✅ Yes | ✅ Yes | ✅ Yes |

## Best Practices

### ✅ DO
```pf
// Use var for truly dynamic scenarios
var response = makeRequest()
if response.isError():
    response = handleError()
end
```

### ❌ DON'T
```pf
// Don't use var when let is sufficient
var count = 0              // Bad: Should use 'let'
for i in range(10):
    count = count + i
end

// Better:
let count = 0
for i in range(10):
    count = count + i
end
```

## Performance Considerations

Using `var` can have slight performance overhead because:
- Type checking happens at runtime
- The compiler cannot optimize as aggressively
- Memory layout may change

**Recommendation**: Prefer `let` for better performance and type safety. Only use `var` when you truly need dynamic typing.

## See Also

- [`let` - Type-safe mutable variable](let.md)
- [`const` - Compile-time constant](const.md)
- [`final` - Runtime constant](final.md)
- [Type System](../advanced/types.md)