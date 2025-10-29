# Sys Module

The Sys module provides system-level operations, including timing, random number generation, and process control.

## Time Functions

### Sys.time(mode?)
Returns the current time in milliseconds since Unix epoch.

**Parameters**:
- `mode: String` (optional) - "float" to return Float, otherwise returns Int

**Returns**: `Int` or `Float` - Milliseconds since epoch

**Example**:
```polyloft
let timestamp = Sys.time()          // Int64: 1698765432100
let floatTime = Sys.time("float")   // Float: 1698765432100.0

// Measure execution time
let start = Sys.time()
someOperation()
let elapsed = Sys.time() - start
println("Took #{elapsed}ms")
```

### Sys.sleep(milliseconds)
Pauses execution for the specified duration.

**Parameters**:
- `milliseconds: Int` - Duration to sleep in milliseconds

**Returns**: `Void`

**Example**:
```polyloft
println("Starting...")
Sys.sleep(1000)  // Sleep for 1 second
println("Done!")

// Animated countdown
for i in range(5, 0, -1):
    println(i)
    Sys.sleep(1000)
end
println("Blast off!")
```

## Random Number Generation

### Sys.random()
Returns a random floating-point number between 0.0 (inclusive) and 1.0 (exclusive).

**Returns**: `Float` - Random number

**Example**:
```polyloft
let r = Sys.random()  // e.g., 0.7234512

// Random integer in range
def randomInt(min: Int, max: Int) -> Int:
    return Math.floor(Sys.random() * (max - min + 1)) + min
end

let dice = randomInt(1, 6)
```

### Sys.seed(value)
Seeds the random number generator with the specified value.

**Parameters**:
- `value: Int` - Seed value

**Returns**: `Void`

**Example**:
```polyloft
// Reproducible random numbers
Sys.seed(12345)
let r1 = Sys.random()  // Always same for seed 12345
let r2 = Sys.random()

// Reset with same seed
Sys.seed(12345)
let r3 = Sys.random()  // Same as r1
```

## User Input

### Sys.input(prompt?, defaultValue?, castType?)
Reads user input from standard input.

**Parameters**:
- `prompt: String` (optional) - Prompt message to display
- `defaultValue: Any` (optional) - Default value if input is empty
- `castType: String` (optional) - Type to cast input to ("int", "float", "bool")

**Returns**: `String` or cast type - User input

**Example**:
```polyloft
// Simple input
let name = Sys.input("Enter your name: ")
println("Hello, #{name}!")

// With default value
let age = Sys.input("Enter age: ", 18, "int")

// Boolean input
let confirm = Sys.input("Continue? (yes/no): ")
if confirm == "yes":
    proceed()
end
```

## Process Control

### Sys.exit(message?)
Terminates the program immediately.

**Parameters**:
- `message: String` (optional) - Exit message

**Returns**: Never returns (exits process)

**Example**:
```polyloft
if errorCondition:
    Sys.exit("Critical error occurred")
end

// Exit with status
Sys.exit()
```

## Practical Examples

### Execution Timer
```polyloft
class Timer:
    private let startTime: Int
    
    Timer():
        this.startTime = Sys.time()
    end
    
    def elapsed() -> Int:
        return Sys.time() - this.startTime
    end
    
    def reset():
        this.startTime = Sys.time()
    end
end

let timer = Timer()
performTask()
println("Task took #{timer.elapsed()}ms")
```

### Rate Limiting
```polyloft
class RateLimiter:
    private let interval: Int
    private let lastCall: Int
    
    RateLimiter(interval: Int):
        this.interval = interval
        this.lastCall = 0
    end
    
    def canCall() -> Bool:
        let now = Sys.time()
        if now - this.lastCall >= this.interval:
            this.lastCall = now
            return true
        end
        return false
    end
end

let limiter = RateLimiter(1000)  // Max once per second
if limiter.canCall():
    performAction()
end
```

### Random Data Generation
```polyloft
// Random string generator
def randomString(length: Int) -> String:
    let chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
    let result = ""
    
    for i in range(length):
        let index = Math.floor(Sys.random() * chars.length())
        result = result + chars[index]
    end
    
    return result
end

let token = randomString(16)
```

### Retry Logic
```polyloft
def retryWithBackoff(operation, maxAttempts: Int):
    let attempts = 0
    
    while attempts < maxAttempts:
        try:
            return operation()
        catch error:
            attempts = attempts + 1
            if attempts >= maxAttempts:
                throw error
            end
            
            // Exponential backoff
            let delay = Math.pow(2, attempts) * 1000
            println("Retry #{attempts} after #{delay}ms...")
            Sys.sleep(delay)
        end
    end
end
```

### Timeout Wrapper
```polyloft
def withTimeout(operation, timeoutMs: Int):
    let done = Channel<Bool>()
    let result = Channel<Any>()
    
    thread spawn do
        try:
            let r = operation()
            result.send(r)
        catch e:
            result.send(e)
        end
    end
    
    select:
        case result.recv() as value:
            return value
        case Sys.sleep(timeoutMs):
            throw "Operation timed out after #{timeoutMs}ms"
    end
end
```

### Progress Indicator
```polyloft
def showProgress(total: Int):
    for i in range(total + 1):
        let percent = (i * 100) / total
        let bar = "=" * (percent / 2)
        print("\r[#{bar}] #{percent}%")
        Sys.sleep(100)
    end
    println("")
end

showProgress(50)
```

### Benchmark Function
```polyloft
def benchmark(name: String, operation, iterations: Int):
    println("Benchmarking: #{name}")
    
    let total = 0
    for i in range(iterations):
        let start = Sys.time()
        operation()
        total = total + (Sys.time() - start)
    end
    
    let avg = total / iterations
    println("Average: #{avg}ms over #{iterations} iterations")
    return avg
end

benchmark("Sort", () => sortArray(data), 100)
```

### Countdown Timer
```polyloft
def countdown(seconds: Int):
    for i in range(seconds, 0, -1):
        println("#{i}...")
        Sys.sleep(1000)
    end
    println("Time's up!")
end

countdown(10)
```

## Notes

- `Sys.time()` returns milliseconds since Unix epoch (January 1, 1970 UTC)
- Random numbers are seeded automatically but can be controlled with `Sys.seed()`
- `Sys.sleep()` suspends the current goroutine, not the entire program
- `Sys.exit()` terminates immediately without running deferred functions
- For file I/O, see the [IO Module](io.md)
- For network operations, see the [Net Module](net.md)

## See Also

- [Math Module](math.md) - Mathematical functions including Math.random()
- [IO Module](io.md) - File and stream operations
- [Concurrency](../concurrency/overview.md) - Async operations and timing
