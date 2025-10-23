# Channels

Channels provide a way for concurrent code to communicate by sending and receiving values. They are inspired by Go's channels.

## Creating Channels

### Unbuffered Channel
```polyloft
let ch = Channel<Int>()
```

An unbuffered channel blocks on send until a receiver is ready, and blocks on receive until a sender is ready.

### Buffered Channel
```polyloft
// Not yet supported - use unbuffered channels
```

## Channel Operations

### Sending Values

```polyloft
let ch = Channel<String>()

thread spawn do
    ch.send("Hello")
    ch.send("World")
end
```

**Behavior**:
- Blocks until receiver is ready (unbuffered)
- Panics if channel is closed

### Receiving Values

```polyloft
let ch = Channel<Int>()

thread spawn do
    ch.send(42)
end

let value = ch.recv()
println(value)  // 42
```

**Behavior**:
- Blocks until sender is ready
- Returns `null` if channel is closed and empty

### Closing Channels

```polyloft
let ch = Channel<Int>()

thread spawn do
    ch.send(1)
    ch.send(2)
    ch.send(3)
    ch.close()  // Signal no more values
end

loop
    let value = ch.recv()
    if value == null:
        break  // Channel closed
    end
    println(value)
end
```

**Notes**:
- Sending on a closed channel panics
- Receiving from a closed channel returns `null`
- Closing a closed channel panics

## Select Statement

The select statement lets you wait on multiple channel operations.

### Basic Select
```polyloft
let ch1 = Channel<Int>()
let ch2 = Channel<String>()

select:
    case ch1.recv() as value:
        println("Got int: #{value}")
    case ch2.recv() as msg:
        println("Got string: #{msg}")
end
```

### Select with Send
```polyloft
select:
    case ch1.recv() as value:
        println("Received: #{value}")
    case ch2.send(42):
        println("Sent 42")
end
```

### Select with Default
```polyloft
select:
    case ch.recv() as value:
        println("Received: #{value}")
    default:
        println("No value available")
end
```

### Timeout with Select
```polyloft
let ch = Channel<Int>()

select:
    case ch.recv() as value:
        println("Got: #{value}")
    case Sys.sleep(5000):
        println("Timeout after 5 seconds")
end
```

## Patterns

### Producer-Consumer

```polyloft
class Producer:
    private let output: Channel
    
    Producer(output: Channel):
        this.output = output
    end
    
    def run():
        for i in range(10):
            this.output.send(i)
            Sys.sleep(100)
        end
        this.output.close()
    end
end

class Consumer:
    private let input: Channel
    
    Consumer(input: Channel):
        this.input = input
    end
    
    def run():
        loop
            let value = this.input.recv()
            if value == null:
                break
            end
            println("Consumed: #{value}")
        end
    end
end

let ch = Channel<Int>()
let producer = Producer(ch)
let consumer = Consumer(ch)

thread spawn do
    producer.run()
end

consumer.run()
```

### Pipeline

```polyloft
// Stage 1: Generate numbers
def generate(out: Channel):
    for i in range(1, 11):
        out.send(i)
    end
    out.close()
end

// Stage 2: Square numbers
def square(input: Channel, output: Channel):
    loop
        let value = input.recv()
        if value == null:
            break
        end
        output.send(value * value)
    end
    output.close()
end

// Stage 3: Print results
def print(input: Channel):
    loop
        let value = input.recv()
        if value == null:
            break
        end
        println(value)
    end
end

// Connect pipeline
let ch1 = Channel<Int>()
let ch2 = Channel<Int>()

thread spawn do
    generate(ch1)
end

thread spawn do
    square(ch1, ch2)
end

print(ch2)
```

### Fan-Out (Multiple Workers)

```polyloft
def worker(id: Int, jobs: Channel, results: Channel):
    loop
        let job = jobs.recv()
        if job == null:
            break
        end
        
        // Simulate work
        Sys.sleep(1000)
        let result = job * 2
        
        results.send(result)
    end
end

let jobs = Channel<Int>()
let results = Channel<Int>()

// Start workers
for i in range(3):
    thread spawn do
        worker(i, jobs, results)
    end
end

// Send jobs
thread spawn do
    for j in range(5):
        jobs.send(j)
    end
    jobs.close()
end

// Collect results
for i in range(5):
    let result = results.recv()
    println("Result: #{result}")
end
```

### Quit Channel

```polyloft
def worker(jobs: Channel, quit: Channel):
    loop
        select:
            case jobs.recv() as job:
                println("Processing: #{job}")
                processJob(job)
            case quit.recv():
                println("Worker stopping")
                return
        end
    end
end

let jobs = Channel<Int>()
let quit = Channel<Bool>()

thread spawn do
    worker(jobs, quit)
end

// Send work
for i in range(5):
    jobs.send(i)
end

// Stop worker
quit.send(true)
```

### Request-Response

```polyloft
class Request:
    let id: Int
    let data: String
    let response: Channel
    
    Request(id: Int, data: String):
        this.id = id
        this.data = data
        this.response = Channel<String>()
    end
end

def handler(requests: Channel):
    loop
        let req = requests.recv()
        if req == null:
            break
        end
        
        // Process request
        let result = "Processed: #{req.data}"
        
        // Send response
        req.response.send(result)
    end
end

let requests = Channel<Request>()

thread spawn do
    handler(requests)
end

// Make request
let req = Request(1, "test data")
requests.send(req)

// Wait for response
let response = req.response.recv()
println(response)
```

### Semaphore with Channels

```polyloft
class Semaphore:
    private let permits: Channel
    
    Semaphore(n: Int):
        this.permits = Channel<Bool>()
        for i in range(n):
            this.permits.send(true)
        end
    end
    
    def acquire():
        this.permits.recv()
    end
    
    def release():
        this.permits.send(true)
    end
end

let sem = Semaphore(3)

for i in range(10):
    thread spawn do
        sem.acquire()
        println("Worker #{i} running")
        Sys.sleep(1000)
        println("Worker #{i} done")
        sem.release()
    end
end

Sys.sleep(15000)
```

### Rate Limiter

```polyloft
class RateLimiter:
    private let tickets: Channel
    
    RateLimiter(rate: Int, interval: Int):
        this.tickets = Channel<Bool>()
        
        thread spawn do
            loop
                Sys.sleep(interval)
                
                // Try to add ticket (non-blocking)
                select:
                    case this.tickets.send(true):
                        // Ticket sent
                    default:
                        // Channel full, skip
                end
            end
        end
    end
    
    def wait():
        this.tickets.recv()
    end
end

let limiter = RateLimiter(5, 1000)  // 5 per second

for i in range(20):
    limiter.wait()
    println("Request #{i}")
end
```

## Best Practices

### 1. Close Channels from Sender
```polyloft
// Good
thread spawn do
    for i in range(10):
        ch.send(i)
    end
    ch.close()  // Sender closes
end
```

### 2. Check for Closed Channels
```polyloft
loop
    let value = ch.recv()
    if value == null:
        break  // Channel closed
    end
    process(value)
end
```

### 3. Use Select for Timeouts
```polyloft
select:
    case ch.recv() as value:
        handleValue(value)
    case Sys.sleep(5000):
        handleTimeout()
end
```

### 4. Avoid Channel Leaks
```polyloft
// Ensure all goroutines can exit
let done = Channel<Bool>()

thread spawn do
    select:
        case ch.recv() as value:
            process(value)
        case done.recv():
            return  // Exit goroutine
    end
end

// Later: signal done
done.send(true)
```

## Common Errors

### Deadlock
```polyloft
// Bad: Will deadlock
let ch = Channel<Int>()
ch.send(42)  // Blocks forever, no receiver

// Good: Receive in another goroutine
thread spawn do
    let value = ch.recv()
    println(value)
end
ch.send(42)
```

### Send on Closed Channel
```polyloft
let ch = Channel<Int>()
ch.close()
ch.send(42)  // Panic!
```

### Double Close
```polyloft
let ch = Channel<Int>()
ch.close()
ch.close()  // Panic!
```

## Performance Considerations

1. **Channel operations have overhead** - Don't use for very high-frequency operations
2. **Unbuffered channels require synchronization** - Can be slower than buffered
3. **Select statements have overhead** - Use simple receive when possible
4. **Channel GC** - Ensure channels can be garbage collected

## Type Safety

Channels are generic and type-safe:

```polyloft
let intCh = Channel<Int>()
let strCh = Channel<String>()

intCh.send(42)       // OK
// intCh.send("hi")  // Type error

let n: Int = intCh.recv()
```

## See Also

- [Concurrency Overview](overview.md) - All concurrency features
- [Threads](threads.md) - Thread management
- [Select Statement](overview.md#select-statement) - Multiplexing channels
- [Defer](defer.md) - Resource cleanup
