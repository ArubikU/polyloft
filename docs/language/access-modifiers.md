# Access Modifiers

Access modifiers in Polyloft control the visibility and accessibility of class members (fields and methods). They help enforce encapsulation and define the public interface of your classes.

## Available Modifiers

- `public` (or `pub`) - Accessible from anywhere
- `private` (or `priv`) - Accessible only within the class
- `protected` (or `prot`) - Accessible within the class and its subclasses
- `static` - Belongs to the class rather than instances

## Syntax

```pf
class ClassName:
    public let publicField
    private let privateField
    protected let protectedField
    static let classField
    
    public def publicMethod():
        // accessible from anywhere
    end
    
    private def privateMethod():
        // accessible only within this class
    end
    
    protected def protectedMethod():
        // accessible within this class and subclasses
    end
    
    static def staticMethod():
        // called on the class itself
    end
end
```

## Public Access

The `public` modifier (or its short form `pub`) makes members accessible from anywhere.

### Public Fields and Methods
```pf
class User:
    public let username
    public var status
    
    User(name):
        this.username = name
        this.status = "active"
    end
    
    public def getInfo() -> String
        return "#{this.username} (#{this.status})"
    end
end

let user = User("alice")
println(user.username)        // OK: public field
user.status = "inactive"      // OK: public field
println(user.getInfo())       // OK: public method
```

### Default Visibility

Fields and methods without an explicit modifier are **public by default**:

```pf
class Point:
    let x  // public by default
    let y  // public by default
    
    Point(px, py):
        this.x = px
        this.y = py
    end
    
    def distance() -> Float  // public by default
        return Math.sqrt(this.x * this.x + this.y * this.y)
    end
end
```

## Private Access

The `private` modifier (or `priv`) restricts access to within the class only.

### Private Fields
```pf
class BankAccount:
    private var balance
    
    BankAccount(initialBalance):
        this.balance = initialBalance
    end
    
    def deposit(amount):
        this.balance = this.balance + amount
    end
    
    def getBalance() -> Float
        return this.balance
    end
end

let account = BankAccount(1000)
account.deposit(500)
println(account.getBalance())  // OK: 1500
// println(account.balance)    // Error: balance is private
```

### Private Methods
```pf
class DataProcessor:
    private def validateData(data) -> Bool
        return data != nil and data.length() > 0
    end
    
    private def sanitizeData(data) -> String
        return data.trim()
    end
    
    public def processData(data) -> String
        if not this.validateData(data):
            throw "Invalid data"
        end
        return this.sanitizeData(data)
    end
end

let processor = DataProcessor()
processor.processData("  hello  ")  // OK: public method
// processor.validateData("test")   // Error: validateData is private
```

## Protected Access

The `protected` modifier (or `prot`) allows access within the class and its subclasses.

### Protected Fields and Methods
```pf
class Vehicle:
    protected var engineStatus
    
    Vehicle():
        this.engineStatus = "off"
    end
    
    protected def startEngine():
        this.engineStatus = "on"
    end
end

class Car < Vehicle:
    def drive():
        this.startEngine()  // OK: protected method accessible in subclass
        if this.engineStatus == "on":
            println("Driving...")
        end
    end
end

let car = Car()
car.drive()           // OK: public method
// car.startEngine()  // Error: startEngine is protected
```

### Template Method Pattern
```pf
class GameCharacter:
    protected var health
    
    GameCharacter():
        this.health = 100
    end
    
    // Public template method
    public def takeDamage(amount):
        this.health = this.health - amount
        this.onDamage(amount)  // Call protected hook
        if this.health <= 0:
            this.onDeath()
        end
    end
    
    // Protected hooks for subclasses
    protected def onDamage(amount):
        println("Took #{amount} damage")
    end
    
    protected def onDeath():
        println("Character died")
    end
end

class Player < GameCharacter:
    @Override
    protected def onDamage(amount):
        println("Player hurt! Health: #{this.health}")
    end
    
    @Override
    protected def onDeath():
        println("Game Over!")
    end
end

let player = Player()
player.takeDamage(30)  // OK: public method
// player.onDamage(10) // Error: onDamage is protected
```

## Static Members

The `static` modifier makes members belong to the class itself rather than instances.

### Static Fields and Methods
```pf
class MathHelper:
    static let PI = 3.14159
    
    static def square(x) -> Float
        return x * x
    end
    
    static def circleArea(radius) -> Float
        return MathHelper.PI * MathHelper.square(radius)
    end
end

println(MathHelper.square(5))      // 25
println(MathHelper.circleArea(10)) // 314.159
```

### Instance Counter Pattern
```pf
class Counter:
    static var totalInstances = 0
    var instanceId
    
    Counter():
        Counter.totalInstances = Counter.totalInstances + 1
        this.instanceId = Counter.totalInstances
    end
    
    static def getTotal() -> Int
        return Counter.totalInstances
    end
    
    def getInfo() -> String
        return "Instance #{this.instanceId} of #{Counter.totalInstances}"
    end
end

let c1 = Counter()
let c2 = Counter()
println(Counter.getTotal())  // 2
println(c1.getInfo())        // Instance 1 of 2
```

### Singleton Pattern
```pf
class Database:
    private static var instance = nil
    
    private Database():
        // Private constructor
    end
    
    static def getInstance() -> Database
        if Database.instance == nil:
            Database.instance = Database()
        end
        return Database.instance
    end
    
    def query(sql) -> String
        return "Results for: #{sql}"
    end
end

let db1 = Database.getInstance()
let db2 = Database.getInstance()
// db1 and db2 are the same instance
```

## Combining Modifiers

You can combine access modifiers with `static`:

```pf
class Config:
    private static var settings = {}
    
    private static def validate(key):
        return key != nil and key.length() > 0
    end
    
    public static def get(key) -> any
        if Config.validate(key):
            return Config.settings.get(key)
        end
        return nil
    end
    
    public static def set(key, value):
        if Config.validate(key):
            Config.settings.set(key, value)
        end
    end
end

Config.set("theme", "dark")
println(Config.get("theme"))  // dark
// Config.validate("key")     // Error: private static
```

## Access Levels Summary

| Modifier | Same Class | Subclass | Anywhere |
|----------|------------|----------|----------|
| `public` | ✅ | ✅ | ✅ |
| `protected` | ✅ | ✅ | ❌ |
| `private` | ✅ | ❌ | ❌ |
| (default) | ✅ | ✅ | ✅ |
| `static` | Class-level access, follows other modifiers |

## Short Forms

Polyloft supports shorter aliases:
- `pub` for `public`
- `priv` for `private`
- `prot` for `protected`

```pf
class Example:
    pub let publicField
    priv let privateField
    prot let protectedField
    
    pub def publicMethod():
        println("Public")
    end
end
```

## Best Practices

### ✅ DO - Use private by default
```pf
class UserService:
    private var cache = {}
    
    public def getUser(id):
        return this.fetchFromCache(id)
    end
    
    private def fetchFromCache(id):
        return this.cache.get(id)
    end
end
```

### ✅ DO - Use protected for extensibility
```pf
class HttpClient:
    protected def buildHeaders() -> Map
        return {"Content-Type": "application/json"}
    end
    
    public def post(url, data):
        let headers = this.buildHeaders()
        return Http.post(url, data, headers)
    end
end

class AuthHttpClient < HttpClient:
    private let token
    
    AuthHttpClient(authToken):
        this.token = authToken
    end
    
    @Override
    protected def buildHeaders() -> Map
        let headers = super.buildHeaders()
        headers.set("Authorization", "Bearer #{this.token}")
        return headers
    end
end
```

### ✅ DO - Use static for utilities
```pf
class StringUtils:
    static def capitalize(str) -> String
        if str.length() == 0:
            return str
        end
        return str.charAt(0).toUpperCase() + str.substring(1)
    end
end
```

### ❌ DON'T - Make everything public
```pf
// Bad: Exposes internals
class BadExample:
    public var internalState    // Should be private
    public var cacheData        // Should be private
    
    public def helperMethod():  // Should be private
        // ...
    end
end
```

## See Also

- [Classes](../definitions/class.md)
- [Annotations](annotations.md)
- [Interfaces](interfaces.md)
- [Sealed Classes](sealed.md)
