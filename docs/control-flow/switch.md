# Switch Statement

The `switch` statement provides pattern matching and multi-way branching in Polyloft.

## Basic Syntax

```pf
switch value:
    case pattern1:
        // code
    case pattern2:
        // code
    default:
        // code
end
```

### Type Pattern Syntax
```pf
switch value:
    case (varName: TypeName):
        // varName is bound to value with type TypeName
    default:
        // code
end
```

## Examples

### Simple Switch
```pf
let day = 3

switch day:
    case 1:
        println("Monday")
    case 2:
        println("Tuesday")
    case 3:
        println("Wednesday")
    case 4:
        println("Thursday")
    case 5:
        println("Friday")
    default:
        println("Weekend")
end
```

### Switch with Strings
```pf
let command = "start"

switch command:
    case "start":
        println("Starting...")
        startProcess()
    case "stop":
        println("Stopping...")
        stopProcess()
    case "restart":
        println("Restarting...")
        stopProcess()
        startProcess()
    default:
        println("Unknown command")
end
```

### Switch with Return Values
```pf
def getGrade(score):
    switch true:
        case score >= 90:
            return "A"
        case score >= 80:
            return "B"
        case score >= 70:
            return "C"
        case score >= 60:
            return "D"
        default:
            return "F"
    end
end

println(getGrade(85))  // "B"
```

### Switch with Enums
```pf
enum Status
    PENDING
    ACTIVE
    COMPLETED
    CANCELLED
end

let status = Status.ACTIVE

switch status:
    case Status.PENDING:
        println("Waiting to start")
    case Status.ACTIVE:
        println("In progress")
    case Status.COMPLETED:
        println("Finished")
    case Status.CANCELLED:
        println("Cancelled")
end
```

### Multiple Cases (Fall-through Pattern)
```pf
let char = "a"

switch char:
    case "a":
    case "e":
    case "i":
    case "o":
    case "u":
        println("Vowel")
    default:
        println("Consonant")
end
```

### Switch with Expressions
```pf
let x = 5
let y = 10

switch x + y:
    case 10:
        println("Sum is 10")
    case 15:
        println("Sum is 15")
    case 20:
        println("Sum is 20")
    default:
        println("Other sum")
end
```

## Pattern Matching

### Type Matching with Variable Binding
```pf
def processValue(value):
    switch value:
        case (x: Int):
            println("Integer: #{x}")
            return x * 2
        case (s: String):
            println("String: #{s}")
            return s.toUpperCase()
        case (arr: Array):
            println("Array with #{arr.length()} items")
            return arr.length()
        case (m: Map):
            println("Map with #{m.size()} entries")
            return m.size()
        default:
            println("Unknown type")
            return nil
    end
end

processValue(42)        // Integer: 42, returns 84
processValue("hello")   // String: hello, returns "HELLO"
processValue([1, 2, 3]) // Array with 3 items, returns 3
```

### Type Matching Without Variable
```pf
def getTypeCategory(value):
    let typeName = Sys.type(value)
    switch typeName:
        case "Integer":
            return "number"
        case "Float":
            return "number"
        case "String":
            return "text"
        case "Array":
            return "collection"
        case "Map":
            return "collection"
        default:
            return "other"
    end
end
```

### Type Matching with Sys.type()
```pf
def processValue(value):
    let typeName = Sys.type(value)
    switch typeName:
        case "Integer":
            println("Processing integer: #{value}")
        case "String":
            println("Processing string: #{value}")
        case "Array":
            println("Processing array with #{value.length()} items")
        case "Map":
            println("Processing map with #{value.size()} entries")
        default:
            println("Unknown type")
    end
end

processValue(42)
processValue("hello")
processValue([1, 2, 3])
```

### Range-like Patterns
```pf
def categorizeAge(age):
    switch true:
        case age < 13:
            return "Child"
        case age < 20:
            return "Teenager"
        case age < 65:
            return "Adult"
        default:
            return "Senior"
    end
end

println(categorizeAge(15))  // "Teenager"
```

## Advanced Examples

### HTTP Status Codes
```pf
def handleResponse(statusCode):
    switch statusCode:
        case 200:
            println("OK")
        case 201:
            println("Created")
        case 204:
            println("No Content")
        case 400:
            println("Bad Request")
        case 401:
            println("Unauthorized")
        case 404:
            println("Not Found")
        case 500:
            println("Internal Server Error")
        default:
            println("Unknown status: #{statusCode}")
    end
end
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
    var state: State
    
    def init():
        this.state = State.IDLE
    end
    
    def handle(event):
        switch this.state:
            case State.IDLE:
                if event == "start":
                    this.state = State.RUNNING
                    println("Started")
                end
            case State.RUNNING:
                if event == "pause":
                    this.state = State.PAUSED
                    println("Paused")
                else:
                    if event == "stop":
                        this.state = State.STOPPED
                        println("Stopped")
                    end
                end
            case State.PAUSED:
                if event == "resume":
                    this.state = State.RUNNING
                    println("Resumed")
                else:
                    if event == "stop":
                        this.state = State.STOPPED
                        println("Stopped")
                    end
                end
            case State.STOPPED:
                println("Machine stopped")
        end
    end
end

let machine = StateMachine()
machine.handle("start")
machine.handle("pause")
machine.handle("resume")
```

### Polymorphic Processing
```pf
class Shape:
    def area():
        return 0
    end
end

class Circle < Shape:
    let radius: Float
    def init(r: Float):
        this.radius = r
    end
    def area():
        return Math.PI * this.radius * this.radius
    end
end

class Rectangle < Shape:
    let width: Float
    let height: Float
    def init(w: Float, h: Float):
        this.width = w
        this.height = h
    end
    def area():
        return this.width * this.height
    end
end

def calculateArea(shape):
    switch shape:
        case (c: Circle):
            println("Circle with radius #{c.radius}")
            return c.area()
        case (r: Rectangle):
            println("Rectangle #{r.width}x#{r.height}")
            return r.area()
        case (s: Shape):
            println("Generic shape")
            return s.area()
        default:
            return 0
    end
end

let circle = Circle(5.0)
let rect = Rectangle(4.0, 6.0)

println(calculateArea(circle))  // Circle with radius 5.0, returns ~78.54
println(calculateArea(rect))    // Rectangle 4.0x6.0, returns 24.0
```

### JSON Value Processing
```pf
def processJsonValue(value):
    switch value:
        case (n: Int):
            return {type: "integer", value: n}
        case (f: Float):
            return {type: "float", value: f}
        case (s: String):
            return {type: "string", value: s, length: s.length()}
        case (b: Bool):
            return {type: "boolean", value: b}
        case (arr: Array):
            return {type: "array", length: arr.length()}
        case (obj: Map):
            return {type: "object", keys: obj.keys()}
        default:
            return {type: "unknown"}
    end
end

println(processJsonValue(42))
println(processJsonValue("hello"))
println(processJsonValue([1, 2, 3]))
println(processJsonValue({a: 1, b: 2}))
```
```pf
def parseCommand(input):
    let parts = input.split(" ")
    let command = parts[0]
    
    switch command:
        case "add":
            if parts.length() >= 3:
                let a = parts[1]
                let b = parts[2]
                println("Result: #{a + b}")
            else:
                println("Usage: add <a> <b>")
            end
        case "multiply":
            if parts.length() >= 3:
                let a = parts[1]
                let b = parts[2]
                println("Result: #{a * b}")
            else:
                println("Usage: multiply <a> <b>")
            end
        case "help":
            println("Available commands: add, multiply, help, exit")
        case "exit":
            println("Goodbye!")
            return false
        default:
            println("Unknown command: #{command}")
    end
    
    return true
end
```

### File Extension Handler
```pf
def getFileType(filename):
    let parts = filename.split(".")
    if parts.length() < 2:
        return "unknown"
    end
    
    let ext = parts[parts.length() - 1].toLowerCase()
    
    switch ext:
        case "txt":
        case "md":
        case "log":
            return "text"
        case "jpg":
        case "jpeg":
        case "png":
        case "gif":
            return "image"
        case "mp3":
        case "wav":
        case "flac":
            return "audio"
        case "mp4":
        case "avi":
        case "mkv":
            return "video"
        case "pf":
            return "polyloft"
        default:
            return "unknown"
    end
end

println(getFileType("script.pf"))      // "polyloft"
println(getFileType("photo.jpg"))      // "image"
println(getFileType("document.txt"))   // "text"
```

## Best Practices

### ✅ DO - Use switch for multiple discrete values
```pf
switch status:
    case "pending":
        handlePending()
    case "active":
        handleActive()
    case "completed":
        handleCompleted()
end
```

### ✅ DO - Always include default case
```pf
switch value:
    case 1:
        doOne()
    case 2:
        doTwo()
    default:
        handleUnexpected()
end
```

### ✅ DO - Use enums with switch
```pf
enum Priority
    LOW
    MEDIUM
    HIGH
end

switch task.priority:
    case Priority.LOW:
        scheduleForLater()
    case Priority.MEDIUM:
        scheduleNormal()
    case Priority.HIGH:
        scheduleUrgent()
end
```

### ❌ DON'T - Use switch for simple if-else
```pf
// Bad: Too simple for switch
switch x == 5:
    case true:
        doSomething()
    case false:
        doOther()
end

// Good: Use if-else
if x == 5:
    doSomething()
else:
    doOther()
end
```

### ❌ DON'T - Nest switches deeply
```pf
// Bad: Too nested
switch type:
    case "A":
        switch subType:
            case "1":
                // ...
        end
end

// Better: Extract to functions
def handleTypeA(subType):
    switch subType:
        case "1":
            // ...
    end
end

switch type:
    case "A":
        handleTypeA(subType)
end
```

## Switch vs If-Else

### Use Switch When:
- Multiple discrete values to check
- Pattern matching needed
- Code clarity improved

### Use If-Else When:
- Simple binary conditions
- Complex boolean expressions
- Range checks

```pf
// Better with switch
switch day:
    case "Monday":
    case "Tuesday":
    case "Wednesday":
        println("Weekday")
end

// Better with if-else
if score >= 90 and attendance > 0.8:
    println("Excellent")
end
```

## See Also

- [Conditionals](conditionals.md)
- [Enums](../language/enums.md)
- [Pattern Matching](../advanced/patterns.md)
