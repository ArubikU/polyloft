# Polyloft Code Examples

This directory contains executable Polyloft code examples that demonstrate all language features.

## Running Examples

To run any example:

```bash
polyloft run <filename>.pf
```

## Available Examples

### Working Examples ✓

1. **01_variables.pf** - Variables (let, const, var) and type annotations
2. **02_classes.pf** - Class definition with correct constructor syntax
3. **03_modifiers.pf** - Access modifiers (private, public, protected)
4. **05_generics.pf** - Generic classes with type parameters
5. **06_math.pf** - Math module functions and constants
6. **07_sys.pf** - Sys module (time, sleep, random)
7. **08_collections.pf** - Array map, filter operations
8. **10_functions.pf** - Functions, lambdas, closures
9. **12_threads.pf** - Thread spawning and joining
10. **13_defer.pf** - Defer statement for cleanup
11. **15_string_interpolation.pf** - String operations
12. **19_list_generic.pf** - Generic List<T> operations
13. **21_loop_statement.pf** - Loop statement (infinite loop with break)
14. **22_http_client.pf** - HTTP Client module (documentation)
15. **23_http_server.pf** - HTTP Server module
16. **24_net_socket.pf** - Net Socket module (TCP client/server)

### Features Tested

- ✅ Variables: `let`, `const`, `var`
- ✅ Types: `Int`, `Float`, `String`, `Bool`
- ✅ Classes with constructor: `ClassName(params):`
- ✅ Access modifiers: `private`, `public`, `protected`
- ✅ Inheritance: `class Child < Parent`
- ✅ Generics: `class Box<T>`
- ✅ Functions: `def functionName(params) -> ReturnType`
- ✅ Lambdas: `(x: Int) -> Int => x * 2`
- ✅ Math module: `Math.PI`, `Math.sqrt()`, etc.
- ✅ Sys module: `Sys.time()`, `Sys.sleep()`, etc.
- ✅ HTTP module: `Http.get()`, `Http.createServer()`, etc.
- ✅ Net module: `Net.listen()`, `Net.connect()`, etc.
- ✅ Arrays: `[1, 2, 3]`
- ✅ Collection methods: `map()`, `filter()`, `reduce()`
- ✅ Control flow: `if/elif/else`, `for in range()`, `for in array`
- ✅ Loop statement: `loop` (infinite loop with `break`)
- ✅ Break and continue
- ✅ Threads: `thread spawn do ... end`, `thread join`
- ✅ Defer: `defer expression`
- ✅ Generic collections: `List<T>`, `Map<K,V>`

## Documentation References

Each example corresponds to documentation in `/docs/`:

- Variables → `/docs/language/variables-and-types.md`
- Classes → `/docs/language/classes-and-objects.md`
- Generics → `/docs/language/generics.md`
- Functions → `/docs/language/functions.md`
- Math → `/docs/stdlib/math.md`
- Sys → `/docs/stdlib/sys.md`
- HTTP → `/docs/stdlib/http.md`
- Net → `/docs/stdlib/net.md`
- Imports → `/docs/advanced/imports.md`
- Concurrency → `/docs/concurrency/`

## Testing All Examples

Run the test script:

```bash
cd /home/runner/work/polyloft/polyloft
for file in docs/examples/code/*.pf; do
    echo "Testing $(basename $file)..."
    ./bin/polyloft run "$file"
done
```

## Notes

- Constructor syntax is `ClassName(params):` not `def __init__(params):`
- All builtin modules (Math, Sys, IO, Net, Http) are automatically available
- Generic syntax uses `<T>` for type parameters
- **Use `loop` for infinite loops, not `while`**
- Loop statement requires explicit `break` to exit
- Channels use `channel[Type]()` syntax (advanced feature)
- Some advanced features may require specific syntax - refer to `/internal/e2e/` tests

## Contributing

To add new examples:

1. Create a `.pf` file with clear comments
2. Test it: `polyloft run yourfile.pf`
3. Update this README
4. Reference it in the relevant documentation

