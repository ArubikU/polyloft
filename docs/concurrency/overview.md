# Concurrency Overview

Polyloft provides modern concurrency primitives for building concurrent and parallel applications. The concurrency model is inspired by Go, JavaScript, and Java.

## Concurrency Models

Polyloft supports three main concurrency models:

### 1. Async/Await (JavaScript-style Promises)
```polyloft
let fetchData = (url: String) => do
    return Promise((resolve, reject) => do
        let response = httpGet(url)
        if response != null:
            resolve(response.body)
        else:
            reject("Failed to fetch data")
        end
    end)
end

fetchData("https://api.example.com/data").then((data) => do
    println("Received: #{data}")
end).catch((error) => do
    println("Error: #{error}")
end)
```

### 2. Channels (Go-style)
```polyloft
let ch = Channel<Int>()

thread spawn do
    ch.send(42)
end

let value = ch.recv()
println(value)  // 42
```

### 3. Threads
```polyloft
let t = thread spawn do
    return expensiveComputation()
end

let result = thread join t
println(result)
```

## Key Features

### Thread Spawning
Create lightweight concurrent tasks:
```polyloft
thread spawn do
    println("Running in background")
end
```

### Channel Communication
Safe message passing between threads:
```polyloft
let ch = Channel<String>()
ch.send("message")
let msg = ch.recv()
```

### Select Statement
Multiplexing multiple channel operations:
```polyloft
select:
    case ch1.recv() as value:
        println("Got from ch1: #{value}")
    case ch2.send(42):
        println("Sent to ch2")
    default:
        println("No operation ready")
end
```

### Defer Statement
Guaranteed cleanup execution:
```polyloft
def processFile(path: String):
    let file = IO.openFile(path)
    defer file.close()
    
    // Work with file
    // file.close() guaranteed to run
end
```

### Promises
JavaScript-style async operations:
```polyloft
let promise = Promise((resolve, reject) => {
    if success:
        resolve(result)
    else:
        reject(error)
    end
})

promise.then((value) => {
    println("Success: #{value}")
}).catch((error) => {
    println("Error: #{error}")
})
```

## Concurrency Patterns

### Producer-Consumer
```polyloft
let ch = Channel<Int>()

// Producer
thread spawn do
    for i in range(10):
        ch.send(i)
    end
    ch.close()
end

// Consumer
thread spawn do
    loop
        let value = ch.recv()
        if value == null:
            break
        end
        println("Consumed: #{value}")
    end
end
```

### Worker Pool
```polyloft
def workerPool(tasks: List, workers: Int):
    let taskCh = Channel<Any>()
    let resultCh = Channel<Any>()
    
    // Start workers
    for i in range(workers):
        thread spawn do
            loop
                let task = taskCh.recv()
                if task == null:
                    break
                end
                let result = processTask(task)
                resultCh.send(result)
            end
        end
    end
    
    // Send tasks
    thread spawn do
        for task in tasks:
            taskCh.send(task)
        end
        taskCh.close()
    end
    
    // Collect results
    let results = List()
    for i in range(tasks.size()):
        results.add(resultCh.recv())
    end
    
    return results
end
```

### Pipeline
```polyloft
def pipeline(input: Channel, stages: List) -> Channel:
    let current = input
    
    for stage in stages:
        let next = Channel()
        let prev = current
        
        thread spawn do
            loop
                let value = prev.recv()
                if value == null:
                    break
                end
                let result = stage(value)
                next.send(result)
            end
            next.close()
        end
        
        current = next
    end
    
    return current
end
```

### Fan-out/Fan-in
```polyloft
// Fan-out: Distribute work
def fanOut(input: Channel, n: Int) -> List:
    let outputs = List()
    
    for i in range(n):
        let ch = Channel()
        outputs.add(ch)
        
        thread spawn do
            loop
                let value = input.recv()
                if value == null:
                    break
                end
                ch.send(value)
            end
            ch.close()
        end
    end
    
    return outputs
end

// Fan-in: Merge results
def fanIn(inputs: List) -> Channel:
    let output = Channel()
    
    for input in inputs:
        thread spawn do
            loop
                let value = input.recv()
                if value == null:
                    break
                end
                output.send(value)
            end
        end
    end
    
    return output
end
```

## Synchronization

### Mutex (via Channels)
```polyloft
class Mutex:
    private let ch: Channel
    
    Mutex():
        this.ch = Channel<Bool>()
        this.ch.send(true)  // Initialize unlocked
    end
    
    def lock():
        this.ch.recv()
    end
    
    def unlock():
        this.ch.send(true)
    end
end

let mutex = Mutex()
mutex.lock()
// Critical section
mutex.unlock()
```

### Semaphore
```polyloft
class Semaphore:
    private let ch: Channel
    
    Semaphore(n: Int):
        this.ch = Channel<Bool>()
        for i in range(n):
            this.ch.send(true)
        end
    end
    
    def acquire():
        this.ch.recv()
    end
    
    def release():
        this.ch.send(true)
    end
end

let sem = Semaphore(3)  // Max 3 concurrent operations
sem.acquire()
performOperation()
sem.release()
```

## Best Practices

### 1. Always Close Channels
```polyloft
let ch = Channel<Int>()

thread spawn do
    for i in range(10):
        ch.send(i)
    end
    ch.close()  // Signal completion
end
```

### 2. Use Defer for Cleanup
```polyloft
def withResource():
    let resource = acquireResource()
    defer resource.release()
    
    // Use resource
end
```

### 3. Handle Channel Closure
```polyloft
loop
    let value = ch.recv()
    if value == null:
        break  // Channel closed
    end
    process(value)
end
```

### 4. Timeout Operations
```polyloft
select:
    case ch.recv() as value:
        println(value)
    case Sys.sleep(5000):
        println("Timeout after 5 seconds")
end
```

### 5. Avoid Blocking Main Thread
```polyloft
// Don't do this:
let result = blockingOperation()

// Do this:
thread spawn do
    let result = blockingOperation()
    handleResult(result)
end
```

## Performance Tips

1. **Use buffered channels** for better throughput
2. **Limit goroutine count** with worker pools
3. **Batch operations** to reduce overhead
4. **Prefer channels** over shared memory
5. **Profile before optimizing**

## Common Pitfalls

### Deadlock
```polyloft
// Bad: Deadlock
let ch = Channel<Int>()
ch.send(42)  // Blocks forever (no receiver)

// Good: Send in separate thread
let ch = Channel<Int>()
thread spawn do
    ch.send(42)
end
let value = ch.recv()
```

### Goroutine Leaks
```polyloft
// Bad: Goroutine never exits
thread spawn do
    loop
        // No way to stop
    end
end

// Good: Use channel for cancellation
let done = Channel<Bool>()

thread spawn do
    loop
        select:
            case done.recv():
                return
            default:
                doWork()
        end
    end
end

// Later: stop the goroutine
done.send(true)
```

### Race Conditions
```polyloft
// Bad: Race condition
let counter = 0

for i in range(100):
    thread spawn do
        counter = counter + 1  // Race!
    end
end

// Good: Use channel
let ch = Channel<Int>()

thread spawn do
    let counter = 0
    loop
        select:
            case ch.recv():
                counter = counter + 1
            case done.recv():
                return counter
        end
    end
end
```

## See Also

- [Async/Await](async-await.md) - Promise-based concurrency
- [Channels](channels.md) - Message passing
- [Threads](threads.md) - Thread management
- [Defer Statement](defer.md) - Cleanup guarantees
