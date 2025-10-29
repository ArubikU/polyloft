# Async/Await and Promises

Polyloft provides Promise-based asynchronous programming using `async()` function and `.await()` method for non-blocking operations.

## Promise System

Polyloft uses JavaScript-style Promises for asynchronous operations.

### Creating Promises

### `async(function)`
Creates a Promise that executes the function asynchronously.

**Parameters:**
- `function`: Function to execute asynchronously

**Returns:** Promise

```pf
let promise = async(() => 42)
let result = promise.await()
println(result)  // 42
```

### Promise Constructor
```pf
let promise = Promise((resolve, reject) => do
    // Do async work
    if success:
        resolve(result)
    else:
        reject(error)
    end
end)
```

## Promise Methods

### `.await()`
Waits for the promise to resolve and returns the result.

**Returns:** Any (the resolved value)

```pf
let promise = async(() => do
    Sys.sleep(1000)
    return "Done!"
end)

let result = promise.await()
println(result)  // "Done!" (after 1 second)
```

### `.then(onFulfilled)`
Chains a callback to execute when the promise resolves.

**Parameters:**
- `onFulfilled` (Function): Callback function

**Returns:** Promise (for chaining)

```pf
let promise = async(() => 10)

promise.then((value) => do
    println("Got value: #{value}")
    return value * 2
end)

Sys.sleep(100)  // Wait for async operation
```

### `.catch(onRejected)`
Chains a callback to handle errors.

**Parameters:**
- `onRejected` (Function): Error handler

**Returns:** Promise

```pf
let promise = async(() => do
    throw RuntimeError("Something went wrong")
end)

promise.catch((err) => do
    println("Error: #{err}")
    return nil
end)

Sys.sleep(100)
```

### `.finally(handler)`
Executes code regardless of promise outcome.

**Parameters:**
- `handler` (Function): Cleanup function

**Returns:** Promise

```pf
let promise = async(() => fetchData())

promise
    .then((data) => processData(data))
    .catch((err) => println("Error: #{err}"))
    .finally(() => println("Cleanup done"))
```

## Examples

### Basic Async Operation
```pf
let promise = async(() => do
    Sys.sleep(1000)
    return 42
end)

let result = promise.await()
println("Result: #{result}")
```

### Promise Chaining
```pf
let result = nil

async(() => 5)
    .then((val) => val * 2)
    .then((val) => do
        result = val + 10
        return result
    end)

Sys.sleep(100)
println(result)  // 20 (5*2+10)
```

### Error Handling
```pf
let promise = async(() => do
    if someCondition:
        throw RuntimeError("error occurred")
    end
    return "success"
end)

let errorMsg = nil
promise.catch((err) => do
    errorMsg = err
    return nil
end)

Sys.sleep(100)
if errorMsg:
    println("Error: #{errorMsg}")
end
```

### Multiple Async Operations
```pf
// Start multiple operations
let p1 = async(() => fetchUser(1))
let p2 = async(() => fetchUser(2))
let p3 = async(() => fetchUser(3))

// Wait for all
let user1 = p1.await()
let user2 = p2.await()
let user3 = p3.await()

println("Loaded #{[user1, user2, user3].length()} users")
```

### Async with External Operations
```pf
let promise = async(() => do
    let response = Http.get("https://api.example.com/data")
    return response.body
end)

promise
    .then((data) => do
        println("Received: #{data}")
        return parseJSON(data)
    end)
    .then((parsed) => do
        println("Parsed: #{parsed}")
    end)
    .catch((err) => do
        println("Request failed: #{err}")
    end)

Sys.sleep(2000)  // Wait for completion
```

## CompletableFuture

Polyloft also provides Java-style CompletableFuture for more control.

### Creating CompletableFuture
```pf
let future = CompletableFuture()

// Complete it from another thread
thread spawn do
    Sys.sleep(100)
    future.complete(42)
end

let result = future.get()
println(result)  // 42
```

### CompletableFuture Methods

**`complete(value)`** - Complete the future with a value
```pf
let future = CompletableFuture()
future.complete(100)
```

**`get()`** - Wait for completion and get value
```pf
let value = future.get()
```

**`getTimeout(milliseconds)`** - Wait with timeout
```pf
try:
    let value = future.getTimeout(1000)
    println("Got: #{value}")
catch e:
    println("Timeout!")
end
```

**`isDone()`** - Check if completed
```pf
if future.isDone():
    let value = future.get()
end
```

**`cancel()`** - Cancel the future
```pf
future.cancel()
```

### CompletableFuture Example
```pf
let future = CompletableFuture()

thread spawn do
    Sys.sleep(500)
    let result = expensiveOperation()
    future.complete(result)
end

println("Waiting for result...")

try:
    let result = future.getTimeout(1000)
    println("Result: #{result}")
catch e:
    println("Operation timed out")
end
```

## Threading

Polyloft supports spawning threads for parallel execution.

### `thread spawn`
```pf
thread spawn do
    println("Running in background")
    Sys.sleep(1000)
    println("Background task done")
end

println("Main thread continues")
Sys.sleep(2000)
```

### With CompletableFuture
```pf
let future = CompletableFuture()

thread spawn do
    let result = heavyComputation()
    future.complete(result)
end

let result = future.get()
println("Computation result: #{result}")
```

## Common Patterns

### Parallel Data Fetching
```pf
def fetchAllUsers(ids):
    let promises = []
    
    for id in ids:
        let p = async(() => fetchUser(id))
        promises = promises.concat([p])
    end
    
    let users = []
    for p in promises:
        users = users.concat([p.await()])
    end
    
    return users
end

let userIds = [1, 2, 3, 4, 5]
let users = fetchAllUsers(userIds)
```

### Retry Pattern
```pf
def asyncWithRetry(operation, maxRetries):
    let attempts = 0
    
    loop attempts < maxRetries:
        let promise = async(operation)
        
        try:
            return promise.await()
        catch e:
            attempts = attempts + 1
            if attempts >= maxRetries:
                throw e
            end
            println("Retry #{attempts}/#{maxRetries}")
            Sys.sleep(1000 * attempts)
        end
    end
end
```

### Timeout with Promise
```pf
def withTimeout(promise, timeoutMs):
    let future = CompletableFuture()
    
    thread spawn do
        Sys.sleep(timeoutMs)
        if not future.isDone():
            future.complete(nil)
        end
    end
    
    thread spawn do
        let result = promise.await()
        if not future.isDone():
            future.complete(result)
        end
    end
    
    return future.get()
end
```

### Background Worker
```pf
let taskQueue = []
let running = true

def worker():
    loop running:
        if taskQueue.length() > 0:
            let task = taskQueue.shift()
            let promise = async(() => processTask(task))
            promise.then((result) => do
                println("Task completed: #{result}")
            end)
        else:
            Sys.sleep(100)
        end
    end
end

thread spawn do
    worker()
end

// Add tasks
taskQueue = taskQueue.concat(["task1", "task2", "task3"])
```

## Best Practices

### ✅ DO - Use async for I/O operations
```pf
let promise = async(() => do
    return IO.readFile("large.txt")
end)

// Continue with other work
doOtherWork()

// Get result when needed
let content = promise.await()
```

### ✅ DO - Chain promises for sequential operations
```pf
async(() => fetchUser(id))
    .then((user) => fetchProfile(user.id))
    .then((profile) => displayProfile(profile))
    .catch((err) => handleError(err))
```

### ✅ DO - Handle errors in promises
```pf
let promise = async(() => riskyOperation())
promise.catch((err) => do
    println("Error: #{err}")
    logError(err)
end)
```

### ❌ DON'T - Forget to wait for promises
```pf
// Bad: Promise never awaited
let promise = async(() => saveData())
// Missing: promise.await()

// Good: Wait for completion
let promise = async(() => saveData())
promise.await()
```

### ❌ DON'T - Block unnecessarily
```pf
// Bad: Blocking main thread
let result = expensiveOperation()

// Good: Use async
let promise = async(() => expensiveOperation())
// Do other work
let result = promise.await()
```

## See Also

- [Channels](channels.md)
- [Threading](threading.md)
- [Concurrency Patterns](../examples/concurrency.md)
- [Http Module](../stdlib/http.md)
