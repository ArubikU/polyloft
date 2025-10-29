# Classes

Classes in Polyloft provide object-oriented programming capabilities with fields, methods, constructors, and inheritance.

## Syntax

### Basic Class
```pf
class ClassName:
    // Fields
    let field1
    var field2
    
    // Constructor
    def init(param1, param2):
        this.field1 = param1
        this.field2 = param2
    end
    
    // Methods
    def methodName():
        // method body
    end
end
```

## Examples

### Simple Class
```pf
class Person:
    let name
    let age
    
    def init(name, age):
        this.name = name
        this.age = age
    end
    
    def greet():
        println("Hello, I'm #{this.name}")
    end
end

let person = Person("Alice", 25)
person.greet()  // Outputs: Hello, I'm Alice
```

### With Type Annotations
```pf
class Rectangle:
    let width: Float
    let height: Float
    
    def init(width: Float, height: Float):
        this.width = width
        this.height = height
    end
    
    def area(): Float
        return this.width * this.height
    end
    
    def perimeter(): Float
        return 2 * (this.width + this.height)
    end
end

let rect = Rectangle(10.0, 5.0)
println("Area: #{rect.area()}")
println("Perimeter: #{rect.perimeter()}")
```

### Private Fields
```pf
class BankAccount:
    private let balance
    
    def init(initialBalance):
        this.balance = initialBalance
    end
    
    def deposit(amount):
        this.balance = this.balance + amount
    end
    
    def getBalance():
        return this.balance
    end
end

let account = BankAccount(1000)
account.deposit(500)
println(account.getBalance())  // 1500
// println(account.balance)    // Error: balance is private
```

### Static Members
```pf
class MathHelper:
    static let PI = 3.14159
    
    static def square(x):
        return x * x
    end
    
    static def cube(x):
        return x * x * x
    end
end

println(MathHelper.PI)
println(MathHelper.square(5))  // 25
println(MathHelper.cube(3))    // 27
```

### Inheritance
```pf
class Animal:
    let name
    
    def init(name):
        this.name = name
    end
    
    def speak():
        println("Some sound")
    end
end

class Dog < Animal:
    def init(name):
        super.init(name)  // Call parent constructor
    end
    
    def speak():
        println("Woof! I'm #{this.name}")
    end
end

class Cat < Animal:
    def init(name):
        super.init(name)
    end
    
    def speak():
        println("Meow! I'm #{this.name}")
    end
end

let dog = Dog("Rex")
let cat = Cat("Whiskers")

dog.speak()  // Woof! I'm Rex
cat.speak()  // Meow! I'm Whiskers
```

### Method Overriding
```pf
class Shape:
    def area():
        return 0
    end
    
    def describe():
        println("This is a shape with area #{this.area()}")
    end
end

class Circle < Shape:
    let radius
    
    def init(radius):
        this.radius = radius
    end
    
    def area():
        return 3.14159 * this.radius * this.radius
    end
end

let circle = Circle(5)
circle.describe()  // This is a shape with area 78.53975
```

### Abstract Pattern
```pf
class Vehicle:
    let brand
    
    def init(brand):
        this.brand = brand
    end
    
    def startEngine():
        // To be overridden
        throw "startEngine must be implemented"
    end
end

class Car < Vehicle:
    def startEngine():
        println("#{this.brand} car engine started")
    end
end

let car = Car("Toyota")
car.startEngine()  // Toyota car engine started
```

### Getters and Setters
```pf
class Temperature:
    private var celsius
    
    def init(c):
        this.celsius = c
    end
    
    def getCelsius():
        return this.celsius
    end
    
    def setCelsius(c):
        this.celsius = c
    end
    
    def getFahrenheit():
        return this.celsius * 9 / 5 + 32
    end
    
    def setFahrenheit(f):
        this.celsius = (f - 32) * 5 / 9
    end
end

let temp = Temperature(25)
println("#{temp.getCelsius()}°C")
println("#{temp.getFahrenheit()}°F")
```

### Composition
```pf
class Engine:
    let horsepower
    
    def init(hp):
        this.horsepower = hp
    end
    
    def start():
        println("Engine with #{this.horsepower}hp started")
    end
end

class Car:
    let brand
    let engine
    
    def init(brand, hp):
        this.brand = brand
        this.engine = Engine(hp)
    end
    
    def start():
        println("Starting #{this.brand}...")
        this.engine.start()
    end
end

let car = Car("Ferrari", 500)
car.start()
```

### Class Variables vs Instance Variables
```pf
class Counter:
    static var totalInstances = 0
    var instanceId
    
    def init():
        Counter.totalInstances = Counter.totalInstances + 1
        this.instanceId = Counter.totalInstances
    end
    
    def getInfo():
        println("Instance #{this.instanceId} of #{Counter.totalInstances}")
    end
end

let c1 = Counter()
let c2 = Counter()
let c3 = Counter()

c1.getInfo()  // Instance 1 of 3
c2.getInfo()  // Instance 2 of 3
c3.getInfo()  // Instance 3 of 3
```

## Advanced Features

### Generic Classes
```pf
class Box<T>:
    let value: T
    
    def init(value: T):
        this.value = value
    end
    
    def getValue(): T
        return this.value
    end
end

let intBox = Box<Int>(42)
let strBox = Box<String>("Hello")

println(intBox.getValue())  // 42
println(strBox.getValue())  // Hello
```

### Interfaces
```pf
interface Drawable:
    def draw(): void
end

class Circle implements Drawable:
    let radius
    
    def init(r):
        this.radius = r
    end
    
    def draw():
        println("Drawing circle with radius #{this.radius}")
    end
end

let shape = Circle(10)
shape.draw()
```

### Sealed Classes
```pf
sealed class Result(Success, Failure)
end

class Success < Result:
    let value
    def init(v):
        this.value = v
    end
end

class Failure < Result:
    let error
    def init(e):
        this.error = e
    end
end

// Only Success and Failure can extend Result
```

## Best Practices

### ✅ DO - Use clear, descriptive class names
```pf
class UserAccount:     // Good
class ShoppingCart:    // Good
class PaymentProcessor: // Good
```

### ✅ DO - Encapsulate data
```pf
class Person:
    private let ssn
    let name
    
    def init(name, ssn):
        this.name = name
        this.ssn = ssn
    end
    
    def getRedactedSSN():
        return "***-**-" + this.ssn[-4:]
    end
end
```

### ✅ DO - Follow Single Responsibility Principle
```pf
class User:
    // Only user-related data and behavior
end

class UserRepository:
    // Only database operations
end

class UserValidator:
    // Only validation logic
end
```

### ❌ DON'T - Create god classes
```pf
// Bad: Too many responsibilities
class Application:
    def handleRequest()
    def validateData()
    def saveToDatabase()
    def sendEmail()
    def generateReport()
    // ... too much!
end
```

## Common Patterns

### Builder Pattern
```pf
class PersonBuilder:
    var name = ""
    var age = 0
    var email = ""
    
    def withName(n):
        this.name = n
        return this
    end
    
    def withAge(a):
        this.age = a
        return this
    end
    
    def withEmail(e):
        this.email = e
        return this
    end
    
    def build():
        return Person(this.name, this.age, this.email)
    end
end

let person = PersonBuilder()
    .withName("Alice")
    .withAge(25)
    .withEmail("alice@example.com")
    .build()
```

### Singleton Pattern
```pf
class Database:
    private static var instance = nil
    
    private def init():
        // Private constructor
    end
    
    static def getInstance():
        if Database.instance == nil:
            Database.instance = Database()
        end
        return Database.instance
    end
end

let db1 = Database.getInstance()
let db2 = Database.getInstance()
// db1 and db2 are the same instance
```

### Factory Pattern
```pf
class ShapeFactory:
    static def createShape(type, size):
        if type == "circle":
            return Circle(size)
        else:
            if type == "square":
                return Square(size)
            end
        end
        return nil
    end
end

let shape = ShapeFactory.createShape("circle", 10)
```

## See Also

- [Functions](function.md)
- [Interfaces](../language/interfaces.md)
- [Enums](../language/enums.md)
- [Records](../language/records.md)
- [Sealed Classes](../language/sealed.md)
- [Generics](../advanced/generics.md)
