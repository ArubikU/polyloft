# Exception Handling

Polyloft provides `try-catch` blocks for error handling and exception management.

## Syntax

### Basic Try-Catch
```pf
try:
    // Code that might throw an exception
catch e:
    // Handle the exception
end
```

### Try-Catch with Type
```pf
try:
    // Code that might throw
catch e: ErrorType:
    // Handle specific error type
end
```

## Examples

### Basic Exception Handling
```pf
try:
    let result = riskyOperation()
    println("Success: #{result}")
catch e:
    println("Error occurred: #{e}")
end
```

### Catching Specific Errors
```pf
try:
    let file = IO.readFile("data.txt")
    println(file)
catch e: FileNotFoundException:
    println("File not found!")
catch e: IOException:
    println("IO error: #{e}")
catch e:
    println("Unknown error: #{e}")
end
```

### Division by Zero
```pf
def safeDivide(a, b):
    try:
        return a / b
    catch e:
        println("Error: Cannot divide by zero")
        return nil
    end
end

println(safeDivide(10, 2))   // 5
println(safeDivide(10, 0))   // Error: Cannot divide by zero, nil
```

### Array Index Error
```pf
let arr = [1, 2, 3]

try:
    println(arr[10])  // Out of bounds
catch e:
    println("Index out of bounds")
end
```

### File Operations
```pf
def readConfig(filename):
    try:
        let content = IO.readFile(filename)
        return parseConfig(content)
    catch e:
        println("Could not read config: #{e}")
        return getDefaultConfig()
    end
end
```

### Network Operations
```pf
def fetchData(url):
    try:
        let response = Http.get(url)
        return response.body
    catch e: NetworkError:
        println("Network error: #{e}")
        return nil
    catch e: TimeoutError:
        println("Request timed out")
        return nil
    catch e:
        println("Unexpected error: #{e}")
        return nil
    end
end
```

## Throwing Exceptions

### Using `throw`
```pf
def validateAge(age):
    if age < 0:
        throw "Age cannot be negative"
    end
    if age > 150:
        throw "Age is too high"
    end
    return age
end

try:
    validateAge(-5)
catch e:
    println("Validation error: #{e}")
end
```

### Throwing Custom Errors
```pf
class ValidationError < Exception:
    let message
    
    def init(msg):
        this.message = msg
    end
end

def validateEmail(email):
    if not email.contains("@"):
        throw ValidationError("Invalid email format")
    end
    return email
end

try:
    validateEmail("invalid")
catch e: ValidationError:
    println("Validation failed: #{e.message}")
end
```

## Nested Try-Catch

```pf
def processData(data):
    try:
        let parsed = parseData(data)
        
        try:
            let validated = validateData(parsed)
            return saveData(validated)
        catch e:
            println("Validation error: #{e}")
            return nil
        end
        
    catch e:
        println("Parse error: #{e}")
        return nil
    end
end
```

## Exception Propagation

```pf
def innerFunction():
    throw "Something went wrong"
end

def middleFunction():
    innerFunction()  // Exception propagates up
end

def outerFunction():
    try:
        middleFunction()
    catch e:
        println("Caught in outer: #{e}")
    end
end

outerFunction()  // Outputs: Caught in outer: Something went wrong
```

## Finally-like Pattern

While Polyloft doesn't have a `finally` keyword, you can ensure cleanup:

```pf
def processWithCleanup():
    let resource = acquireResource()
    let success = false
    
    try:
        processResource(resource)
        success = true
    catch e:
        println("Error: #{e}")
    end
    
    // Cleanup always happens
    releaseResource(resource)
    
    return success
end
```

## Common Patterns

### Retry Pattern
```pf
def fetchWithRetry(url, maxRetries):
    let attempts = 0
    
    loop attempts < maxRetries:
        try:
            return Http.get(url)
        catch e:
            attempts = attempts + 1
            if attempts < maxRetries:
                println("Retry #{attempts}/#{maxRetries}")
                Sys.sleep(1000)
            end
        end
    end
    
    throw "Failed after #{maxRetries} attempts"
end
```

### Graceful Degradation
```pf
def getDataWithFallback():
    try:
        return fetchFromPrimarySource()
    catch e:
        println("Primary source failed, trying backup...")
        try:
            return fetchFromBackupSource()
        catch e2:
            println("Backup failed, using cache...")
            return getCachedData()
        end
    end
end
```

### Error Recovery
```pf
def parseWithRecovery(data):
    try:
        return strictParse(data)
    catch e:
        println("Strict parsing failed, trying lenient...")
        try:
            return lenientParse(data)
        catch e2:
            println("All parsing failed, using defaults")
            return getDefaultData()
        end
    end
end
```

### Validation Chain
```pf
def validateUser(user):
    try:
        validateName(user.name)
        validateEmail(user.email)
        validateAge(user.age)
        return true
    catch e:
        println("Validation failed: #{e}")
        return false
    end
end
```

### Safe Resource Access
```pf
def safeRead(filename):
    try:
        let file = IO.openFile(filename)
        let content = file.read()
        file.close()
        return content
    catch e:
        println("Error reading file: #{e}")
        return nil
    end
end
```

### Transaction Pattern
```pf
def performTransaction(operations):
    let completed = []
    
    try:
        for op in operations:
            op.execute()
            completed = completed.concat([op])
        end
        return true
    catch e:
        println("Transaction failed: #{e}")
        // Rollback completed operations
        for op in completed:
            try:
                op.rollback()
            catch e2:
                println("Rollback failed: #{e2}")
            end
        end
        return false
    end
end
```

## Built-in Exceptions

Common exception types in Polyloft:
- `Exception` - Base exception class
- `RuntimeError` - General runtime errors
- `TypeError` - Type mismatch errors
- `ArityError` - Wrong number of arguments
- `IndexError` - Array/string index errors
- `KeyError` - Map key errors
- `ValueError` - Invalid value errors
- `FileNotFoundException` - File not found
- `IOException` - I/O errors
- `NetworkError` - Network errors
- `TimeoutError` - Timeout errors

## Best Practices

### ✅ DO - Be specific with error handling
```pf
try:
    processData()
catch e: ValidationError:
    handleValidationError(e)
catch e: NetworkError:
    handleNetworkError(e)
end
```

### ✅ DO - Log errors appropriately
```pf
try:
    criticalOperation()
catch e:
    println("[ERROR] #{Sys.time()}: #{e}")
    logToFile(e)
    notifyAdmin(e)
end
```

### ✅ DO - Clean up resources
```pf
let file = IO.openFile("data.txt")
try:
    processFile(file)
catch e:
    println("Error: #{e}")
end
file.close()  // Always close
```

### ❌ DON'T - Catch and ignore all errors
```pf
// Bad: Silent failure
try:
    importantOperation()
catch e:
    // Ignoring error - bad!
end

// Good: Handle or propagate
try:
    importantOperation()
catch e:
    println("Operation failed: #{e}")
    throw e  // Re-throw if can't handle
end
```

### ❌ DON'T - Use exceptions for control flow
```pf
// Bad: Using exceptions for logic
try:
    let value = array[index]
catch e:
    value = defaultValue
end

// Good: Check first
let value = if index < array.length():
    array[index]
else:
    defaultValue
end
```

## See Also

- [Classes](../definitions/class.md)
- [Error Types](../reference/errors.md)
- [Conditionals](conditionals.md)
