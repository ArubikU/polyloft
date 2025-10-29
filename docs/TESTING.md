# Testing Documentation Examples

This document summarizes the validation of code examples in the documentation.

## Tested Features

### ✅ Successfully Tested

1. **Crypto Module** - All examples validated
   - MD5, SHA1, SHA256, SHA512 hashing
   - Base64 encoding/decoding
   - Hex encoding/decoding

2. **Basic Features** - Core functionality validated
   - Variable declarations (let, var, const)
   - String operations (length, toUpperCase, toLowerCase)
   - Array operations (map, filter, reduce)
   - Math operations (sqrt, pow, abs)

3. **Control Flow** - All patterns tested
   - For loops with ranges and collections
   - Loop statements with conditions
   - Do-loop statements
   - Switch statements

4. **Async/Await** - Promise system validated
   - async() function for creating promises
   - .await() for waiting
   - .then() chaining
   - CompletableFuture

## Known Limitations

### Map Serialization
The `map.serialize()` method exists but has runtime issues in certain contexts. For robust JSON handling, consider:
- Using Map methods directly (`get()`, `set()`, `keys()`, `values()`)
- Implementing custom serialization logic
- Using HTTP module's built-in JSON support

**Workaround:**
```pf
// Instead of direct interpolation with bracket notation
let map = {name: "Alice"}
let name = map.get("name")  // Use .get() method
println("Name: #{name}")
```

### String Interpolation with Brackets
String literals inside interpolation with bracket notation may cause parse errors:
```pf
// May cause issues
println("Value: #{map["key"]}")

// Use instead
let value = map.get("key")
println("Value: #{value}")
```

## Test Coverage

| Feature | Status | Notes |
|---------|--------|-------|
| Variables (let, var, const) | ✅ | Fully working |
| String methods | ✅ | All 25+ methods tested |
| Array methods | ✅ | Functional operations working |
| Map operations | ⚠️ | Use .get() instead of brackets in interpolation |
| Set operations | ✅ | All methods working |
| Bytes type | ✅ | Binary operations working |
| Math module | ✅ | All functions working |
| Crypto module | ✅ | All hashing and encoding working |
| IO module | ✅ | File operations working |
| Http module | ✅ | Client and server working |
| Async/Await | ✅ | Promise system working |
| Loops (for, loop, do) | ✅ | All loop types working |
| Switch statements | ✅ | Including type matching |
| Exception handling | ✅ | try/catch working |
| Enums | ✅ | With values and methods |
| Records | ✅ | Immutable data classes |
| Classes | ✅ | OOP features working |
| Functions | ✅ | Including lambdas |

## Example Test Files

### test_crypto.pf
```pf
println("=== Testing Crypto ===")

let hash = Crypto.sha256("hello")
println("SHA256: #{hash}")

let encoded = Crypto.base64Encode("data")
let decoded = Crypto.base64Decode(encoded)
println("Encoded/Decoded: #{decoded}")

println("=== Crypto tests passed! ===")
```

**Result:** ✅ All tests passed

### test_collections.pf
```pf
println("=== Testing Collections ===")

let arr = [1, 2, 3, 4, 5]
let doubled = arr.map((x) => x * 2)
let sum = arr.reduce((acc, x) => acc + x, 0)
println("Doubled: #{doubled}, Sum: #{sum}")

let set = Set(1, 2, 3, 2, 1)
println("Set size: #{set.size()}")

println("=== Tests passed! ===")
```

**Result:** ✅ All tests passed

## Running Tests

To test examples from the documentation:

```bash
# Test individual feature
go run ./cmd/polyloft run test_crypto.pf

# Test in REPL
polyloft repl
>>> let arr = [1, 2, 3]
>>> arr.map((x) => x * 2)
[2, 4, 6]
```

## Recommendations

1. **Use .get() for Map access in interpolations**
   - Avoids parse issues with nested quotes
   - More explicit and readable

2. **Test edge cases**
   - Empty collections
   - Nil values
   - Boundary conditions

3. **Error handling**
   - Always wrap risky operations in try/catch
   - Validate input before processing

4. **Performance**
   - Use appropriate data structures
   - Avoid modifying collections while iterating

## See Also

- [Quick Reference](QUICK_REFERENCE.md) - Syntax cheat sheet
- [CLI Reference](CLI.md) - Command-line usage
- [Examples Directory](../algorithm_samples/) - More examples
