# Enums

Enums in Polyloft provide a way to define a type with a fixed set of named constants.

## Basic Enum

### Simple Enum
```pf
enum Color
    RED
    GREEN
    BLUE
end

let color = Color.RED
println(color.name)      // "RED"
println(color.ordinal)   // 0
```

### Using Enums
```pf
enum Status
    PENDING
    ACTIVE
    INACTIVE
    DELETED
end

let status = Status.ACTIVE
println(status.name)      // "ACTIVE"
println(status.ordinal)   // 1
```

## Enum with Values

### Enum with Constructor
```pf
enum Planet
    MERCURY(3.7)
    VENUS(8.87)
    EARTH(9.8)
    MARS(3.71)
    
    var gravity: Float
    
    Planet(g: Float):
        this.gravity = g
    end
    
    def weight(mass: Float): Float
        return mass * this.gravity
    end
end

let earth = Planet.EARTH
println(earth.gravity)         // 9.8
println(earth.weight(75.0))    // 735.0

let mars = Planet.MARS
println(mars.gravity)          // 3.71
println(mars.weight(75.0))     // 278.25
```

### Enum with Multiple Fields
```pf
enum HttpStatus
    OK(200, "OK")
    NOT_FOUND(404, "Not Found")
    SERVER_ERROR(500, "Internal Server Error")
    
    var code: Int
    var message: String
    
    HttpStatus(c: Int, msg: String):
        this.code = c
        this.message = msg
    end
    
    def isError(): Bool
        return this.code >= 400
    end
end

let status = HttpStatus.NOT_FOUND
println(status.code)           // 404
println(status.message)        // "Not Found"
println(status.isError())      // true
```

## Enum Methods

### Built-in Properties

Every enum constant has:
- `name` - The constant's name as a String
- `ordinal` - The constant's position (starting from 0)

```pf
enum Day
    MONDAY
    TUESDAY
    WEDNESDAY
    THURSDAY
    FRIDAY
    SATURDAY
    SUNDAY
end

let day = Day.WEDNESDAY
println(day.name)      // "WEDNESDAY"
println(day.ordinal)   // 2
```

### Static Methods

### `valueOf(name)`
Returns the enum constant with the specified name.

```pf
enum Color
    RED
    GREEN
    BLUE
end

let color = Color.valueOf("GREEN")
println(color.name)  // "GREEN"
```

### `values()`
Returns an array of all enum constants.

```pf
enum Color
    RED
    GREEN
    BLUE
end

let colors = Color.values()
for color in colors:
    println(color.name)
end
// Outputs: RED, GREEN, BLUE
```

### `names()`
Returns an array of all enum constant names.

```pf
enum Color
    RED
    GREEN
    BLUE
end

let names = Color.names()
println(names)  // ["RED", "GREEN", "BLUE"]
```

### `size()`
Returns the number of enum constants.

```pf
enum Color
    RED
    GREEN
    BLUE
end

println(Color.size())  // 3
```

## Examples

### Switch with Enum
```pf
enum TrafficLight
    RED
    YELLOW
    GREEN
end

def getAction(light: TrafficLight): String
    switch light:
        case TrafficLight.RED:
            return "Stop"
        case TrafficLight.YELLOW:
            return "Slow down"
        case TrafficLight.GREEN:
            return "Go"
    end
end

let light = TrafficLight.RED
println(getAction(light))  // "Stop"
```

### Enum in Classes
```pf
enum Priority
    LOW
    MEDIUM
    HIGH
    CRITICAL
end

class Task:
    let title: String
    let priority: Priority
    
    Task(title: String, priority: Priority):
        this.title = title
        this.priority = priority
    end
    
    def isUrgent(): Bool
        return this.priority == Priority.CRITICAL or 
               this.priority == Priority.HIGH
    end
end

let task = Task("Fix bug", Priority.CRITICAL)
println(task.isUrgent())  // true
```

### State Machine
```pf
enum State
    IDLE
    RUNNING
    PAUSED
    STOPPED
end

class StateMachine:
    var currentState: State
    
    StateMachine():
        this.currentState = State.IDLE
    end
    
    def start():
        if this.currentState == State.IDLE:
            this.currentState = State.RUNNING
            println("Started")
        end
    end
    
    def pause():
        if this.currentState == State.RUNNING:
            this.currentState = State.PAUSED
            println("Paused")
        end
    end
    
    def resume():
        if this.currentState == State.PAUSED:
            this.currentState = State.RUNNING
            println("Resumed")
        end
    end
    
    def stop():
        this.currentState = State.STOPPED
        println("Stopped")
    end
end

let machine = StateMachine()
machine.start()
machine.pause()
machine.resume()
machine.stop()
```

### Iterating Over Enums
```pf
enum Season
    SPRING
    SUMMER
    FALL
    WINTER
end

for season in Season.values():
    println("Season #{season.ordinal + 1}: #{season.name}")
end
```

### Enum Comparison
```pf
enum Level
    BEGINNER
    INTERMEDIATE
    ADVANCED
    EXPERT
end

let userLevel = Level.INTERMEDIATE
let requiredLevel = Level.ADVANCED

if userLevel.ordinal < requiredLevel.ordinal:
    println("Level too low")
else:
    println("Access granted")
end
```

### Enum with Methods
```pf
enum Operation
    ADD
    SUBTRACT
    MULTIPLY
    DIVIDE
    
    def apply(a: Float, b: Float): Float
        switch this:
            case Operation.ADD:
                return a + b
            case Operation.SUBTRACT:
                return a - b
            case Operation.MULTIPLY:
                return a * b
            case Operation.DIVIDE:
                if b != 0:
                    return a / b
                end
                return 0.0
        end
    end
end

let op = Operation.MULTIPLY
let result = op.apply(5.0, 3.0)
println(result)  // 15.0
```

### Configuration with Enums
```pf
enum Environment
    DEVELOPMENT("http://localhost:3000", true)
    STAGING("https://staging.example.com", true)
    PRODUCTION("https://example.com", false)
    
    var apiUrl: String
    var debugMode: Bool
    
    Environment(url: String, debug: Bool):
        this.apiUrl = url
        this.debugMode = debug
    end
    
    def log(message: String):
        if this.debugMode:
            println("[#{this.name}] #{message}")
        end
    end
end

let env = Environment.DEVELOPMENT
println(env.apiUrl)      // "http://localhost:3000"
env.log("Starting app")  // "[DEVELOPMENT] Starting app"
```

## Sealed Enums

Sealed enums restrict which classes can extend them.

```pf
sealed enum Result
    SUCCESS
    FAILURE
end

// Only SUCCESS and FAILURE can be used
```

## Best Practices

### ✅ DO - Use enums for fixed sets of constants
```pf
enum Direction
    NORTH
    SOUTH
    EAST
    WEST
end
```

### ✅ DO - Add behavior to enums
```pf
enum FileType
    TEXT(".txt")
    IMAGE(".jpg")
    VIDEO(".mp4")
    
    var extension: String
    
    FileType(ext: String):
        this.extension = ext
    end
    
    def isMedia(): Bool
        return this == FileType.IMAGE or this == FileType.VIDEO
    end
end
```

### ✅ DO - Use enums for type safety
```pf
def processOrder(status: OrderStatus):
    // Compiler ensures only valid OrderStatus values
end
```

### ❌ DON'T - Use strings instead of enums
```pf
// Bad: Easy to make typos
let status = "actve"  // Typo!

// Good: Type-safe
let status = Status.ACTIVE
```

## See Also

- [Records](records.md)
- [Sealed Classes](sealed.md)
- [Switch Statement](../control-flow/switch.md)
- [Classes](../definitions/class.md)
