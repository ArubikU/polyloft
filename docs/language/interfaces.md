# Interfaces

Interfaces in Polyloft define contracts that classes must implement. They specify methods that implementing classes must provide, enabling polymorphism and design by contract.

## Syntax

```pf
interface InterfaceName:
    methodSignature1()
    methodSignature2() -> ReturnType
    methodSignature3(param: Type) -> ReturnType
end
```

## Basic Interface Definition

```pf
interface Drawable:
    draw() -> void
    clear() -> void
end
```

## Implementing Interfaces

Classes implement interfaces using the `implements` keyword:

```pf
class Canvas implements Drawable:
    Canvas():
        // Constructor
    end
    
    def draw():
        println("Drawing on canvas")
    end
    
    def clear():
        println("Clearing canvas")
    end
end

let canvas = Canvas()
canvas.draw()   // Drawing on canvas
canvas.clear()  // Clearing canvas
```

## Multiple Interfaces

A class can implement multiple interfaces:

```pf
interface Saveable:
    save() -> void
end

interface Loadable:
    load() -> void
end

class Document implements Saveable, Loadable:
    let content
    
    Document():
        this.content = ""
    end
    
    def save():
        IO.writeFile("document.txt", this.content)
        println("Document saved")
    end
    
    def load():
        this.content = IO.readFile("document.txt")
        println("Document loaded")
    end
end
```

## Interface with Method Signatures

Interfaces can specify parameter types and return types:

```pf
interface Calculator:
    add(a: Float, b: Float) -> Float
    subtract(a: Float, b: Float) -> Float
    multiply(a: Float, b: Float) -> Float
    divide(a: Float, b: Float) -> Float
end

class BasicCalculator implements Calculator:
    def add(a: Float, b: Float): Float
        return a + b
    end
    
    def subtract(a: Float, b: Float): Float
        return a - b
    end
    
    def multiply(a: Float, b: Float): Float
        return a * b
    end
    
    def divide(a: Float, b: Float): Float
        if b == 0:
            throw "Division by zero"
        end
        return a / b
    end
end

let calc = BasicCalculator()
println(calc.add(10.0, 5.0))      // 15.0
println(calc.multiply(4.0, 3.0))   // 12.0
```

## Inheritance with Interfaces

Classes can extend other classes and implement interfaces:

```pf
interface Logger:
    log(message: String) -> void
end

class BaseService:
    protected var name
    
    BaseService(serviceName):
        this.name = serviceName
    end
end

class UserService < BaseService implements Logger:
    UserService(name):
        super(name)
    end
    
    def log(message: String):
        println("[#{this.name}] #{message}")
    end
    
    def getUser(id):
        this.log("Fetching user #{id}")
        // Fetch user logic
    end
end

let service = UserService("UserService")
service.log("Service started")
service.getUser(123)
```

## Interface Inheritance

Interfaces can extend other interfaces:

```pf
interface Printable:
    print() -> void
end

interface Exportable < Printable:
    export(format: String) -> void
end

class Report implements Exportable:
    let data
    
    Report(reportData):
        this.data = reportData
    end
    
    def print():
        println("Printing report...")
        println(this.data)
    end
    
    def export(format: String):
        println("Exporting report as #{format}")
        // Export logic
    end
end

let report = Report("Sales Data")
report.print()
report.export("PDF")
```

## Polymorphism with Interfaces

Interfaces enable polymorphism:

```pf
interface Shape:
    area() -> Float
    perimeter() -> Float
end

class Rectangle implements Shape:
    let width: Float
    let height: Float
    
    Rectangle(w: Float, h: Float):
        this.width = w
        this.height = h
    end
    
    def area(): Float
        return this.width * this.height
    end
    
    def perimeter(): Float
        return 2 * (this.width + this.height)
    end
end

class Circle implements Shape:
    let radius: Float
    
    Circle(r: Float):
        this.radius = r
    end
    
    def area(): Float
        return 3.14159 * this.radius * this.radius
    end
    
    def perimeter(): Float
        return 2 * 3.14159 * this.radius
    end
end

def printShapeInfo(shape: Shape):
    println("Area: #{shape.area()}")
    println("Perimeter: #{shape.perimeter()}")
end

let rect = Rectangle(10.0, 5.0)
let circle = Circle(7.0)

printShapeInfo(rect)
printShapeInfo(circle)
```

## Strategy Pattern

Interfaces are perfect for the Strategy pattern:

```pf
interface PaymentStrategy:
    pay(amount: Float) -> void
end

class CreditCardPayment implements PaymentStrategy:
    let cardNumber
    
    CreditCardPayment(number):
        this.cardNumber = number
    end
    
    def pay(amount: Float):
        println("Paying $#{amount} with credit card #{this.cardNumber}")
    end
end

class PayPalPayment implements PaymentStrategy:
    let email
    
    PayPalPayment(userEmail):
        this.email = userEmail
    end
    
    def pay(amount: Float):
        println("Paying $#{amount} via PayPal (#{this.email})")
    end
end

class ShoppingCart:
    var items = []
    var paymentStrategy: PaymentStrategy
    
    def setPaymentStrategy(strategy: PaymentStrategy):
        this.paymentStrategy = strategy
    end
    
    def checkout():
        let total = this.calculateTotal()
        this.paymentStrategy.pay(total)
    end
    
    def calculateTotal(): Float
        return 99.99  // Simplified
    end
end

let cart = ShoppingCart()
cart.setPaymentStrategy(CreditCardPayment("1234-5678"))
cart.checkout()

cart.setPaymentStrategy(PayPalPayment("user@example.com"))
cart.checkout()
```

## Observer Pattern

```pf
interface Observer:
    update(data: any) -> void
end

interface Subject:
    attach(observer: Observer) -> void
    detach(observer: Observer) -> void
    notify() -> void
end

class NewsPublisher implements Subject:
    var observers = []
    var latestNews = ""
    
    def attach(observer: Observer):
        this.observers = this.observers.concat([observer])
    end
    
    def detach(observer: Observer):
        // Remove observer logic
    end
    
    def notify():
        for obs in this.observers:
            obs.update(this.latestNews)
        end
    end
    
    def publishNews(news):
        this.latestNews = news
        this.notify()
    end
end

class EmailSubscriber implements Observer:
    let email
    
    EmailSubscriber(userEmail):
        this.email = userEmail
    end
    
    def update(data: any):
        println("Email to #{this.email}: #{data}")
    end
end

class SMSSubscriber implements Observer:
    let phone
    
    SMSSubscriber(phoneNumber):
        this.phone = phoneNumber
    end
    
    def update(data: any):
        println("SMS to #{this.phone}: #{data}")
    end
end

let publisher = NewsPublisher()
publisher.attach(EmailSubscriber("user@example.com"))
publisher.attach(SMSSubscriber("555-1234"))
publisher.publishNews("Breaking news!")
```

## Generic Interfaces

Interfaces can use generic type parameters:

```pf
interface Container<T>:
    add(item: T) -> void
    get(index: Int) -> T
    size() -> Int
end

class ArrayList<T> implements Container<T>:
    var items = []
    
    def add(item: T):
        this.items = this.items.concat([item])
    end
    
    def get(index: Int): T
        return this.items[index]
    end
    
    def size(): Int
        return this.items.length()
    end
end

let stringList = ArrayList<String>()
stringList.add("Hello")
stringList.add("World")
println(stringList.get(0))  // Hello
```

## Best Practices

### ✅ DO - Keep interfaces focused
```pf
// Good: Single responsibility
interface Readable:
    read() -> String
end

interface Writable:
    write(data: String) -> void
end

// Class can implement both if needed
class File implements Readable, Writable:
    def read(): String
        return IO.readFile(this.path)
    end
    
    def write(data: String):
        IO.writeFile(this.path, data)
    end
end
```

### ✅ DO - Use interfaces for abstraction
```pf
// Good: Program to interface, not implementation
interface DataStore:
    save(key: String, value: any) -> void
    load(key: String) -> any
end

class MemoryStore implements DataStore:
    var data = {}
    
    def save(key: String, value: any):
        this.data.set(key, value)
    end
    
    def load(key: String): any
        return this.data.get(key)
    end
end

class FileStore implements DataStore:
    def save(key: String, value: any):
        IO.writeFile(key, JSON.stringify(value))
    end
    
    def load(key: String): any
        return JSON.parse(IO.readFile(key))
    end
end

// Code works with any DataStore implementation
def processData(store: DataStore):
    store.save("key", "value")
    return store.load("key")
end
```

### ✅ DO - Name interfaces descriptively
```pf
// Good: Clear, descriptive names
interface Serializable:
    serialize() -> String
end

interface Comparable:
    compareTo(other: any) -> Int
end

interface Closeable:
    close() -> void
end
```

### ❌ DON'T - Create fat interfaces
```pf
// Bad: Too many responsibilities
interface Everything:
    save()
    load()
    print()
    export()
    validate()
    transform()
    // ... too many methods
end

// Good: Split into focused interfaces
interface Persistable:
    save()
    load()
end

interface Printable:
    print()
end

interface Exportable:
    export()
end
```

### ❌ DON'T - Add implementation to interfaces
```pf
// Bad: Interfaces should only declare contracts
interface Bad:
    def doSomething():
        println("Implementation")  // Wrong!
    end
end

// Good: Just the signature
interface Good:
    doSomething() -> void
end
```

## See Also

- [Classes](../definitions/class.md)
- [Built-in Interfaces](builtin-interfaces.md)
- [Access Modifiers](access-modifiers.md)
- [Annotations](annotations.md)
- [Sealed Classes](sealed.md)
