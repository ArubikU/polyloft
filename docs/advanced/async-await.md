# Async/Await

Polyloft provides async/await for asynchronous programming, allowing non-blocking operations.

## Syntax

### Async Function
```pf
async def functionName(params):
    // Asynchronous code
    return result
end
```

### Await Expression
```pf
let result = await asyncFunction()
```

## Basic Examples

### Simple Async Function
```pf
async def fetchData(url):
    let response = await Http.get(url)
    return response.body
end

let data = await fetchData("https://api.example.com/data")
println(data)
```

### Multiple Async Operations
```pf
async def getUserData(userId):
    let user = await fetchUser(userId)
    let posts = await fetchPosts(userId)
    let comments = await fetchComments(userId)
    
    return {
        user: user,
        posts: posts,
        comments: comments
    }
end
```

### Error Handling in Async
```pf
async def safeAsyncOperation():
    try:
        let result = await riskyAsyncCall()
        return result
    catch e:
        println("Async error: #{e}")
        return nil
    end
end
```

## Parallel Execution

### Launching Multiple Async Tasks
```pf
async def fetchMultiple():
    // Start all requests in parallel
    let task1 = async: fetchData("url1")
    let task2 = async: fetchData("url2")
    let task3 = async: fetchData("url3")
    
    // Wait for all to complete
    let data1 = await task1
    let data2 = await task2
    let data3 = await task3
    
    return [data1, data2, data3]
end
```

### Promise.all Pattern
```pf
async def fetchAll(urls):
    let tasks = []
    for url in urls:
        tasks = tasks.concat([async: Http.get(url)])
    end
    
    let results = []
    for task in tasks:
        results = results.concat([await task])
    end
    
    return results
end

let urls = ["url1", "url2", "url3"]
let responses = await fetchAll(urls)
```

## Async with Channels

```pf
async def worker(id, ch):
    loop:
        let task = await ch.receive()
        if task == nil:
            break
        end
        
        println("Worker #{id} processing #{task}")
        Sys.sleep(1000)
        
        await ch.send("Result from worker #{id}")
    end
end

async def main():
    let ch = Channel()
    
    // Start workers
    async: worker(1, ch)
    async: worker(2, ch)
    
    // Send tasks
    await ch.send("Task 1")
    await ch.send("Task 2")
    
    // Get results
    let result1 = await ch.receive()
    let result2 = await ch.receive()
    
    println(result1)
    println(result2)
end
```

## Examples

### Sequential vs Parallel
```pf
// Sequential (slow)
async def sequential():
    let r1 = await fetchData("url1")  // Wait 1s
    let r2 = await fetchData("url2")  // Wait 1s
    let r3 = await fetchData("url3")  // Wait 1s
    return [r1, r2, r3]  // Total: 3s
end

// Parallel (fast)
async def parallel():
    let t1 = async: fetchData("url1")
    let t2 = async: fetchData("url2")
    let t3 = async: fetchData("url3")
    
    return [
        await t1,  // All run concurrently
        await t2,  // Total: ~1s
        await t3
    ]
end
```

### Async File Operations
```pf
async def processFile(filename):
    let content = await IO.readFileAsync(filename)
    let processed = await processData(content)
    await IO.writeFileAsync("output.txt", processed)
    return "Done"
end

let result = await processFile("input.txt")
println(result)
```

### Async HTTP Server
```pf
async def handleRequest(request, response):
    let data = await fetchFromDatabase(request.params.id)
    response.json(data)
end

let server = Http.createServer()
server.get("/api/users/:id", handleRequest)
await server.listen(8080)
```

### Timeout Pattern
```pf
async def withTimeout(task, timeoutMs):
    let timer = async:
        Sys.sleep(timeoutMs)
        throw "Timeout"
    end
    
    // Race between task and timer
    try:
        return await task
    catch e:
        if e == "Timeout":
            println("Operation timed out")
            return nil
        end
        throw e
    end
end

let result = await withTimeout(
    async: slowOperation(),
    5000
)
```

### Retry with Backoff
```pf
async def retryWithBackoff(operation, maxRetries):
    let delay = 1000
    
    for attempt in range(maxRetries):
        try:
            return await operation()
        catch e:
            if attempt < maxRetries - 1:
                println("Attempt #{attempt + 1} failed, retrying...")
                await Sys.sleepAsync(delay)
                delay = delay * 2  // Exponential backoff
            else:
                throw e
            end
        end
    end
end
```

### Data Pipeline
```pf
async def pipeline(data):
    let step1 = await transform1(data)
    let step2 = await transform2(step1)
    let step3 = await transform3(step2)
    return step3
end

async def processBatch(items):
    let results = []
    for item in items:
        let result = await pipeline(item)
        results = results.concat([result])
    end
    return results
end
```

### Concurrent Workers
```pf
async def workerPool(tasks, numWorkers):
    let results = []
    let taskQueue = tasks
    let activeWorkers = 0
    
    async def processTask(task):
        let result = await executeTask(task)
        results = results.concat([result])
    end
    
    loop taskQueue.length() > 0 or activeWorkers > 0:
        if taskQueue.length() > 0 and activeWorkers < numWorkers:
            let task = taskQueue.shift()
            activeWorkers = activeWorkers + 1
            
            async:
                await processTask(task)
                activeWorkers = activeWorkers - 1
            end
        else:
            await Sys.sleepAsync(100)
        end
    end
    
    return results
end
```

### Event Loop Pattern
```pf
async def eventLoop():
    let running = true
    
    loop running:
        let event = await getNextEvent()
        
        switch event.type:
            case "click":
                await handleClick(event)
            case "keypress":
                await handleKeypress(event)
            case "quit":
                running = false
        end
    end
end
```

## Best Practices

### ✅ DO - Use await for async operations
```pf
async def fetchUserData():
    let user = await fetchUser()
    let profile = await fetchProfile(user.id)
    return {user: user, profile: profile}
end
```

### ✅ DO - Parallelize independent operations
```pf
// Good: Parallel execution
let task1 = async: operation1()
let task2 = async: operation2()
let r1 = await task1
let r2 = await task2
```

### ✅ DO - Handle async errors properly
```pf
async def safeOperation():
    try:
        return await riskyAsyncCall()
    catch e:
        logError(e)
        return defaultValue
    end
end
```

### ❌ DON'T - Await in loops unnecessarily
```pf
// Bad: Sequential (slow)
for item in items:
    await processItem(item)
end

// Good: Parallel (fast)
let tasks = items.map((item) => async: processItem(item))
for task in tasks:
    await task
end
```

### ❌ DON'T - Forget to await
```pf
// Bad: Returns promise, not value
async def getData():
    return Http.get(url)  // Missing await!
end

// Good: Returns actual data
async def getData():
    return await Http.get(url)
end
```

## See Also

- [Channels](channels.md)
- [Concurrency Patterns](../examples/concurrency.md)
- [Http Module](../stdlib/http.md)
