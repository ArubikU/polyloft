# Sealed Classes

Sealed classes in Polyloft restrict which classes can extend them. They define a closed hierarchy where all subclasses must be declared within the sealed class definition or in the same file.

## Syntax

```pf
sealed class ClassName(Subclass1, Subclass2, ...)
end
```

The sealed class declaration specifies all permitted subclasses in parentheses.

## Basic Sealed Class

```pf
sealed class Result(Success, Failure)
end

class Success < Result:
    let value
    
    Success(v):
        this.value = v
    end
end

class Failure < Result:
    let error
    
    Failure(e):
        this.error = e
    end
end

// Only Success and Failure can extend Result
// Any other class attempting to extend Result will be rejected
```

## Use Cases

### Representing State

Sealed classes are perfect for representing a fixed set of states:

```pf
sealed class ConnectionState(Connected, Disconnected, Connecting, Error)
end

class Connected < ConnectionState:
    let timestamp
    
    Connected():
        this.timestamp = Sys.time()
    end
end

class Disconnected < ConnectionState:
    Disconnected():
    end
end

class Connecting < ConnectionState:
    let attempts
    
    Connecting(attemptCount):
        this.attempts = attemptCount
    end
end

class Error < ConnectionState:
    let message
    
    Error(msg):
        this.message = msg
    end
end

def handleConnection(state: ConnectionState):
    if state instanceof Connected:
        println("Connected at #{state.timestamp}")
    else:
        if state instanceof Disconnected:
            println("Disconnected")
        else:
            if state instanceof Connecting:
                println("Connecting (attempt #{state.attempts})")
            else:
                if state instanceof Error:
                    println("Error: #{state.message}")
                end
            end
        end
    end
end
```

### Result Types

Sealed classes work well for result/option types:

```pf
sealed class Option(Some, None)
end

class Some < Option:
    let value
    
    Some(v):
        this.value = v
    end
    
    def get():
        return this.value
    end
    
    def getOrElse(defaultValue):
        return this.value
    end
end

class None < Option:
    None():
    end
    
    def get():
        throw "Cannot get value from None"
    end
    
    def getOrElse(defaultValue):
        return defaultValue
    end
end

def divide(a: Float, b: Float): Option
    if b == 0:
        return None()
    end
    return Some(a / b)
end

let result = divide(10.0, 2.0)
if result instanceof Some:
    println("Result: #{result.get()}")
else:
    println("No result")
end
```

### API Responses

```pf
sealed class ApiResponse(Success, Error, Loading)
end

class Success < ApiResponse:
    let data
    
    Success(responseData):
        this.data = responseData
    end
end

class Error < ApiResponse:
    let message
    let code
    
    Error(msg, errorCode):
        this.message = msg
        this.code = errorCode
    end
end

class Loading < ApiResponse:
    Loading():
    end
end

def handleApiResponse(response: ApiResponse):
    if response instanceof Success:
        println("Data: #{response.data}")
    else:
        if response instanceof Error:
            println("Error #{response.code}: #{response.message}")
        else:
            if response instanceof Loading:
                println("Loading...")
            end
        end
    end
end
```

### Algebraic Data Types

Sealed classes enable algebraic data type patterns:

```pf
sealed class Expression(Number, Add, Multiply, Variable)
end

class Number < Expression:
    let value: Float
    
    Number(v: Float):
        this.value = v
    end
end

class Add < Expression:
    let left: Expression
    let right: Expression
    
    Add(l: Expression, r: Expression):
        this.left = l
        this.right = r
    end
end

class Multiply < Expression:
    let left: Expression
    let right: Expression
    
    Multiply(l: Expression, r: Expression):
        this.left = l
        this.right = r
    end
end

class Variable < Expression:
    let name: String
    
    Variable(varName: String):
        this.name = varName
    end
end

def evaluate(expr: Expression, vars: Map): Float
    if expr instanceof Number:
        return expr.value
    else:
        if expr instanceof Add:
            return evaluate(expr.left, vars) + evaluate(expr.right, vars)
        else:
            if expr instanceof Multiply:
                return evaluate(expr.left, vars) * evaluate(expr.right, vars)
            else:
                if expr instanceof Variable:
                    return vars.get(expr.name)
                end
            end
        end
    end
    return 0.0
end

// Example: (5 + x) * 2
let expr = Multiply(
    Add(Number(5.0), Variable("x")),
    Number(2.0)
)

let vars = {"x": 3.0}
println(evaluate(expr, vars))  // (5 + 3) * 2 = 16
```

## Pattern Matching with Switch

Sealed classes work well with switch statements:

```pf
sealed class PaymentMethod(CreditCard, PayPal, BankTransfer)
end

class CreditCard < PaymentMethod:
    let cardNumber
    
    CreditCard(number):
        this.cardNumber = number
    end
end

class PayPal < PaymentMethod:
    let email
    
    PayPal(userEmail):
        this.email = userEmail
    end
end

class BankTransfer < PaymentMethod:
    let accountNumber
    
    BankTransfer(account):
        this.accountNumber = account
    end
end

def processPayment(method: PaymentMethod, amount: Float):
    switch method:
        case CreditCard:
            println("Charging $#{amount} to card #{method.cardNumber}")
        case PayPal:
            println("Charging $#{amount} via PayPal (#{method.email})")
        case BankTransfer:
            println("Transferring $#{amount} from account #{method.accountNumber}")
    end
end

let payment = CreditCard("1234-5678-9012-3456")
processPayment(payment, 99.99)
```

## Sealed Enums

Sealed classes can represent enhanced enums:

```pf
sealed class LogLevel(Debug, Info, Warning, Error, Fatal)
end

class Debug < LogLevel:
    Debug():
    end
    
    def getSeverity(): Int
        return 0
    end
end

class Info < LogLevel:
    Info():
    end
    
    def getSeverity(): Int
        return 1
    end
end

class Warning < LogLevel:
    Warning():
    end
    
    def getSeverity(): Int
        return 2
    end
end

class Error < LogLevel:
    let stackTrace
    
    Error(trace):
        this.stackTrace = trace
    end
    
    def getSeverity(): Int
        return 3
    end
end

class Fatal < LogLevel:
    let stackTrace
    let context
    
    Fatal(trace, ctx):
        this.stackTrace = trace
        this.context = ctx
    end
    
    def getSeverity(): Int
        return 4
    end
end

def log(level: LogLevel, message: String):
    let prefix = ""
    if level instanceof Debug:
        prefix = "[DEBUG]"
    else:
        if level instanceof Info:
            prefix = "[INFO]"
        else:
            if level instanceof Warning:
                prefix = "[WARN]"
            else:
                if level instanceof Error:
                    prefix = "[ERROR]"
                else:
                    if level instanceof Fatal:
                        prefix = "[FATAL]"
                    end
                end
            end
        end
    end
    println("#{prefix} #{message}")
end
```

## Benefits of Sealed Classes

### 1. Exhaustive Checking
The compiler knows all possible subclasses, enabling exhaustive checking in pattern matching.

### 2. Controlled Inheritance
Prevents external code from creating unauthorized subclasses.

### 3. Clear Intent
Documents that the hierarchy is complete and closed.

### 4. Better Tooling
IDEs can provide better autocomplete and navigation for sealed hierarchies.

## Best Practices

### ✅ DO - Use for closed hierarchies
```pf
// Good: Fixed set of related types
sealed class HttpMethod(GET, POST, PUT, DELETE, PATCH)
end

class GET < HttpMethod:
    GET():
    end
end

class POST < HttpMethod:
    let body
    
    POST(requestBody):
        this.body = requestBody
    end
end

// ... other methods
```

### ✅ DO - Combine with pattern matching
```pf
sealed class AsyncResult(Pending, Complete, Failed)
end

class Pending < AsyncResult:
    Pending():
    end
end

class Complete < AsyncResult:
    let data
    
    Complete(result):
        this.data = result
    end
end

class Failed < AsyncResult:
    let error
    
    Failed(err):
        this.error = err
    end
end

def handleAsync(result: AsyncResult):
    if result instanceof Pending:
        println("Still processing...")
    else:
        if result instanceof Complete:
            println("Success: #{result.data}")
        else:
            if result instanceof Failed:
                println("Failed: #{result.error}")
            end
        end
    end
end
```

### ✅ DO - Use for state machines
```pf
sealed class TrafficLight(Red, Yellow, Green)
end

class Red < TrafficLight:
    Red():
    end
    
    def next(): TrafficLight
        return Green()
    end
end

class Yellow < TrafficLight:
    Yellow():
    end
    
    def next(): TrafficLight
        return Red()
    end
end

class Green < TrafficLight:
    Green():
    end
    
    def next(): TrafficLight
        return Yellow()
    end
end
```

### ❌ DON'T - Use when inheritance should be open
```pf
// Bad: Animals should allow extension
// sealed class Animal(Dog, Cat)
// end

// Good: Open hierarchy
class Animal:
    def makeSound():
        println("Some sound")
    end
end

class Dog < Animal:
    def makeSound():
        println("Woof")
    end
end

class Cat < Animal:
    def makeSound():
        println("Meow")
    end
end

// Users can add new animals
class Bird < Animal:
    def makeSound():
        println("Chirp")
    end
end
```

### ❌ DON'T - Overuse sealed classes
```pf
// Bad: Unnecessary sealing
// sealed class Configuration(AppConfig)
// end

// Good: No need to seal if only one subclass makes sense
class Configuration:
    let settings
    
    Configuration():
        this.settings = {}
    end
end
```

## Sealed Classes vs Enums

**Use Enums when:**
- Simple constant values
- No associated data
- No behavior per variant

**Use Sealed Classes when:**
- Each variant has different data
- Each variant has different behavior
- Complex state representation

```pf
// Enum: Simple constants
enum Color
    RED
    GREEN
    BLUE
end

// Sealed: Complex variants
sealed class Shape(Circle, Rectangle, Triangle)
end

class Circle < Shape:
    let radius: Float
    
    Circle(r: Float):
        this.radius = r
    end
    
    def area(): Float
        return 3.14159 * this.radius * this.radius
    end
end

class Rectangle < Shape:
    let width: Float
    let height: Float
    
    Rectangle(w: Float, h: Float):
        this.width = w
        this.height = h
    end
    
    def area(): Float
        return this.width * this.height
    end
end
```

## See Also

- [Classes](../definitions/class.md)
- [Enums](enums.md)
- [Interfaces](interfaces.md)
- [Switch Statement](../control-flow/switch.md)
- [Type Checking](../reference/types.md)
