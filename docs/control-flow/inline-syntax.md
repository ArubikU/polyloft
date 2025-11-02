# Inline Control Flow Syntax

Polyloft supports inline (single-line) syntax for control flow statements, allowing you to write more concise code when the body is a single statement.

## Inline If Statement

```pf
// Inline if
if condition: statement

// Example
let x = 5
if x > 0: println("Positive")

// Inline if-elif-else chain
if x > 10: println("Greater than 10")
elif x > 5: println("Greater than 5")
else: println("5 or less")
```

### Mixed Inline and Multi-line

You can mix inline and multi-line syntax in the same if-elif-else chain:

```pf
let score = 85

if score >= 90:
    println("Grade: A")
    println("Excellent work!")
end
elif score >= 80: println("Grade: B")  // Inline elif
elif score >= 70: println("Grade: C")  // Inline elif
else:
    println("Grade: D or F")
    println("Please study more")
end
```

## Inline For Loop

```pf
// Inline for-in loop
for item in collection: statement

// Example
for i in [1, 2, 3]: println(i)

// With destructuring
for key, value in map: println("#{key}: #{value}")

// With where clause
for x in numbers where x > 0: println(x)
```

### Practical Examples

```pf
// Process each item
for file in files: processFile(file)

// Print array elements
for item in items: println("- #{item}")

// Iterate with index (using range)
for i in range(10): println("Number #{i}")
```

## Inline Loop Statement

```pf
// Inline loop with condition
loop condition: statement

// Example
let count = 0
loop count < 5: count = count + 1

// Inline infinite loop (use with break)
loop: if shouldStop(): break
```

### Loop Examples

```pf
// Simple counter
let i = 0
loop i < 10: i = i + 1

// With condition check
let running = true
loop running: running = checkStatus()

// Processing until done
loop hasMoreData(): processNext Data()
```

## Inline Switch Cases

```pf
// Inline case statements
switch value:
    case 1: println("One")
    case 2: println("Two")  
    case 3: println("Three")
    default: println("Other")
end
```

### Switch with Actions

```pf
switch command:
    case "start": startServer()
    case "stop": stopServer()
    case "restart": restartServer()
    case "status": showStatus()
    default: println("Unknown command")
end
```

## Multi-line vs Inline

Both syntaxes are supported and can be mixed as needed:

```pf
// Multi-line (requires 'end' for if/for/loop)
if condition:
    statement1
    statement2
    statement3
end

// Inline (no 'end' needed)
if condition: statement

// For loops - multi-line
for item in collection:
    process(item)
    validate(item)
end

// For loops - inline
for item in collection: process(item)
```

## When to Use Each Style

### Use Inline Syntax When:

- The body is a single, simple statement
- The line remains readable and not too long
- Quick actions or checks are needed
- Writing concise list processing code

```pf
// Good inline uses
if error: return nil
for x in list: println(x)
loop count < max: count = count + 1
```

### Use Multi-line Syntax When:

- The body has multiple statements
- Complex logic is involved
- Better readability is needed
- Debugging would be easier with separate lines

```pf
// Good multi-line uses
if isValid:
    prepareData()
    processData()
    saveResults()
    logSuccess()
end

for user in users:
    validateUser(user)
    updateDatabase(user)
    sendNotification(user)
end
```

## Chaining and Nesting

Inline syntax can be nested but should be used judiciously:

```pf
// Simple nesting is OK
for row in grid: for col in row: println(col)

// Complex nesting - prefer multi-line
for user in users:
    if user.isActive:
        for order in user.orders:
            processOrder(order)
        end
    end
end
```

## Benefits

- **Conciseness**: Single-statement bodies don't need multiple lines
- **Readability**: Short conditionals are easier to scan
- **Flexibility**: Mix inline and multi-line as needed
- **Familiarity**: Similar to Python's and other languages' inline syntax
- **Less Boilerplate**: Reduces visual clutter for simple cases

## Best Practices

### ✅ DO - Use inline for simple, single statements

```pf
if error: return nil
for item in list: process(item)
loop count < 10: count = count + 1
```

### ✅ DO - Keep inline statements readable

```pf
// Good: Clear and concise
if x > 0: println("Positive")

// Good: Still readable
for user in users where user.isActive: notifyUser(user)
```

### ✅ DO - Use multi-line for complex logic

```pf
if complexCondition:
    validateInput()
    processData()
    saveResults()
    logActivity()
end
```

### ❌ DON'T - Make lines too long

```pf
// Bad: Too long, hard to read
if user.isAuthenticated and user.hasPermission("admin") and user.isActive: performAdminOperation(user, data, options)

// Good: Use multi-line
if user.isAuthenticated and user.hasPermission("admin") and user.isActive:
    performAdminOperation(user, data, options)
end
```

### ❌ DON'T - Use inline for multiple statements

```pf
// Not supported: Can't do this
if condition: statement1; statement2

// Good: Use multi-line
if condition:
    statement1
    statement2
end
```

### ❌ DON'T - Overuse nesting with inline syntax

```pf
// Bad: Hard to read
for x in list: for y in x: if y > 0: println(y)

// Good: Use multi-line with proper indentation
for x in list:
    for y in x:
        if y > 0:
            println(y)
        end
    end
end
```

## Examples in Context

### Configuration Processing

```pf
// Process configuration with inline syntax
for key, value in config:
    if key.startsWith("app."): appConfig.set(key, value)
    elif key.startsWith("db."): dbConfig.set(key, value)
    else: defaultConfig.set(key, value)
end
```

### Data Validation

```pf
// Validate data with inline checks
for field in requiredFields:
    if not data.has(field): errors.add("Missing field: #{field}")
end

if errors.isEmpty(): return Success(data)
else: return Failure(errors)
```

### Simple Filtering

```pf
// Filter and process
let results = []
for item in items:
    if item.isValid: results.push(item)
end
```

### Quick Logging

```pf
// Conditional logging
if debugMode: println("[DEBUG] Processing item #{currentItem}")
if errorCount > 0: println("[WARN] #{errorCount} errors occurred")
```

## Compatibility Note

- Inline syntax does not require `end` keyword
- Multi-line syntax always requires `end` (except for inline cases within switch)
- Both styles can be mixed in the same program
- The parser automatically detects which style based on line breaks

## See Also

- [Conditionals](conditionals.md) - Full if/elif/else documentation
- [Loops](loops.md) - For and loop statement details
- [Switch Statement](switch.md) - Switch/case documentation
- [Quick Reference](../QUICK_REFERENCE.md) - Language syntax overview
