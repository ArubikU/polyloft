# Async/Await with Promises

Polyloft provides JavaScript-style Promises for asynchronous programming, allowing you to write non-blocking code with clean chaining syntax.

## Promise Basics

A `Promise` represents an asynchronous operation that will eventually complete or fail.

### Creating a Promise

```polyloft
let promise = Promise((resolve, reject) => do
    // Perform async operation
    if success:
        resolve(result)
    else:
        reject("Error message")
    end
end)
```

### Using Promises

```polyloft
promise.then((value) => do
    println("Success: #{value}")
    return processedValue
end).catch((error) => do
    println("Error: #{error}")
end).finally(() => do
    println("Cleanup")
end)
```

## Promise Methods

### then(onFulfilled)

Called when the promise is fulfilled successfully:

```polyloft
let promise = Promise((resolve, reject) => do
    resolve(42)
end)

promise.then((value) => do
    println("Got: #{value}")  // Got: 42
end)
```

### catch(onRejected)

Called when the promise is rejected:

```polyloft
let promise = Promise((resolve, reject) => do
    reject("Something went wrong")
end)

promise.catch((error) => do
    println("Error: #{error}")  // Error: Something went wrong
end)
```

### finally(onFinally)

Always called, regardless of success or failure:

```polyloft
let promise = Promise((resolve, reject) => do
    resolve("done")
end)

promise.finally(() => do
    println("Cleanup code here")
end)
```

## Promise Chaining

Promises can be chained to perform sequential operations:

```polyloft
Promise((resolve, reject) => do
    resolve(1)
end).then((value) => do
    println("Step 1: #{value}")  // Step 1: 1
    return value + 1
end).then((value) => do
    println("Step 2: #{value}")  // Step 2: 2
    return value + 1
end).then((value) => do
    println("Step 3: #{value}")  // Step 3: 3
end)
```

## Practical Examples

### HTTP Request Simulation

```polyloft
def fetchData(url: String):
    return Promise((resolve, reject) => do
        // Simulate HTTP request
        Sys.sleep(1000)
        if url.startsWith("https://"):
            resolve("Data from #{url}")
        else:
            reject("Invalid URL")
        end
    end)
end

fetchData("https://api.example.com").then((data) => do
    println(data)
end).catch((error) => do
    println("Failed: #{error}")
end)
```

### Multiple Async Operations

```polyloft
def loadUser(id: Int):
    return Promise((resolve, reject) => do
        Sys.sleep(500)
        resolve(Map("id" -> id, "name" -> "User #{id}"))
    end)
end

def loadPosts(userId: Int):
    return Promise((resolve, reject) => do
        Sys.sleep(500)
        resolve(List("Post 1", "Post 2", "Post 3"))
    end)
end

// Chain operations
loadUser(1).then((user) => do
    println("Loaded user: #{user.get('name')}")
    return loadPosts(user.get("id"))
end).then((posts) => do
    println("Loaded posts: #{posts}")
end)
```

### Error Handling

```polyloft
def riskyOperation():
    return Promise((resolve, reject) => do
        let random = Math.random()
        if random > 0.5:
            resolve("Success!")
        else:
            reject("Failed with random: #{random}")
        end
    end)
end

riskyOperation().then((result) => do
    println("Result: #{result}")
end).catch((error) => do
    println("Error: #{error}")
end).finally(() => do
    println("Operation completed")
end)
```

## CompletableFuture

Polyloft also provides Java-style `CompletableFuture` for more control:

```polyloft
let future = CompletableFuture()

// Complete it later
thread spawn do
    Sys.sleep(1000)
    future.complete(42)
end

// Wait for result
let result = future.get()
println("Result: #{result}")  // Result: 42
```

### CompletableFuture Methods

- `complete(value)` - Complete with a value
- `completeExceptionally(error)` - Complete with an error
- `get()` - Wait for completion and get value
- `getTimeout(ms)` - Get value with timeout
- `isDone()` - Check if completed
- `isCancelled()` - Check if cancelled
- `cancel()` - Cancel the future

### Example with Timeout

```polyloft
let future = CompletableFuture()

thread spawn do
    Sys.sleep(5000)
    future.complete("Too slow!")
end

try:
    let result = future.getTimeout(2000)  // 2 second timeout
    println(result)
catch error:
    println("Timeout!")  // This will be printed
end
```

## Best Practices

### 1. Always Handle Errors

```polyloft
promise.then((value) => do
    // Handle success
end).catch((error) => do
    // Always handle errors
    println("Error: #{error}")
end)
```

### 2. Use Finally for Cleanup

```polyloft
let resource = acquireResource()

performAsync(resource).finally(() => do
    resource.release()
end)
```

### 3. Chain Operations Logically

```polyloft
fetchUser()
    .then(validateUser)
    .then(loadData)
    .then(processData)
    .catch(handleError)
```

### 4. Return Promises for Chaining

```polyloft
def step1():
    return Promise((resolve, reject) => do
        resolve(1)
    end)
end

def step2(value: Int):
    return Promise((resolve, reject) => do
        resolve(value + 1)
    end)
end

step1().then(step2).then((result) => do
    println(result)  // 2
end)
```

## Common Patterns

### Parallel Execution

```polyloft
// Start multiple operations in parallel
let p1 = fetchData("url1")
let p2 = fetchData("url2")
let p3 = fetchData("url3")

// Wait for all (manual implementation)
let results = List()
p1.then((data) => results.add(data))
p2.then((data) => results.add(data))
p3.then((data) => results.add(data))
```

### Retry Logic

```polyloft
def retryOperation(maxRetries: Int):
    def attempt(retriesLeft: Int):
        return Promise((resolve, reject) => do
            try:
                let result = riskyCall()
                resolve(result)
            catch error:
                if retriesLeft > 0:
                    println("Retrying... #{retriesLeft} attempts left")
                    attempt(retriesLeft - 1).then(resolve).catch(reject)
                else:
                    reject(error)
                end
            end
        end)
    end
    return attempt(maxRetries)
end

retryOperation(3).then((result) => do
    println("Success: #{result}")
end).catch((error) => do
    println("All retries failed: #{error}")
end)
```

## See Also

- [Concurrency Overview](overview.md) - All concurrency models
- [Channels](channels.md) - Message passing concurrency
- [Threads](threads.md) - Thread-based concurrency
