# Records

Records in Polyloft provide a concise way to create immutable data classes with automatic constructor and methods.

## Basic Record

### Simple Record
```pf
record Point(x: Int, y: Int)
end

let p = Point(10, 20)
println("x: #{p.x}, y: #{p.y}")
```

### Record with Methods
```pf
record Point(x: Int, y: Int)
    def distance():
        return Math.sqrt(this.x * this.x + this.y * this.y)
    end
    
    def add(other: Point): Point
        return Point(this.x + other.x, this.y + other.y)
    end
end

let p1 = Point(3, 4)
let p2 = Point(1, 2)

println(p1.distance())     // 5.0
let p3 = p1.add(p2)
println("Sum: #{p3.x}, #{p3.y}")  // Sum: 4, 6
```

## Record Features

### Automatic Constructor
Records automatically generate a constructor with all fields as parameters.

```pf
record Person(name: String, age: Int, city: String)
end

let person = Person("Alice", 25, "NYC")
println(person.name)  // "Alice"
```

### Immutable Fields
Record fields are immutable by default.

```pf
record User(id: Int, username: String)
end

let user = User(1, "alice")
// user.id = 2  // Error: Cannot modify record field
```

### Type Checking
Records enforce types at construction.

```pf
record Product(name: String, price: Float, quantity: Int)
end

let product = Product("Widget", 19.99, 10)
println(Sys.type(product))  // "Product"
println(product instanceof Product)  // true
```

## Advanced Records

### Records with Computed Properties
```pf
record Rectangle(width: Float, height: Float)
    def area(): Float
        return this.width * this.height
    end
    
    def perimeter(): Float
        return 2 * (this.width + this.height)
    end
    
    def isSquare(): Bool
        return this.width == this.height
    end
end

let rect = Rectangle(10.0, 5.0)
println("Area: #{rect.area()}")           // 25.0
println("Perimeter: #{rect.perimeter()}") // 30.0
println("Is square: #{rect.isSquare()}")  // false
```

### Nested Records
```pf
record Address(street: String, city: String, zip: String)
end

record Person(name: String, age: Int, address: Address)
end

let addr = Address("123 Main St", "NYC", "10001")
let person = Person("Alice", 25, addr)

println(person.address.city)  // "NYC"
```

### Records with Validation
```pf
record Email(address: String)
    def init(address: String):
        if not address.contains("@"):
            throw ValidationError("Invalid email")
        end
        this.address = address
    end
    
    def domain(): String
        return this.address.split("@")[1]
    end
end

try:
    let email = Email("alice@example.com")
    println(email.domain())  // "example.com"
catch e:
    println("Error: #{e}")
end
```

## Common Patterns

### Data Transfer Objects (DTOs)
```pf
record UserDTO(
    id: Int,
    username: String,
    email: String,
    createdAt: Int
)
end

def getUserDTO(user):
    return UserDTO(
        user.id,
        user.username,
        user.email,
        Sys.time()
    )
end
```

### Value Objects
```pf
record Money(amount: Float, currency: String)
    def add(other: Money): Money
        if this.currency != other.currency:
            throw "Currency mismatch"
        end
        return Money(this.amount + other.amount, this.currency)
    end
    
    def multiply(factor: Float): Money
        return Money(this.amount * factor, this.currency)
    end
    
    def toString(): String
        return "#{this.amount} #{this.currency}"
    end
end

let price = Money(19.99, "USD")
let total = price.multiply(3)
println(total.toString())  // "59.97 USD"
```

### Coordinates and Geometry
```pf
record Vector2D(x: Float, y: Float)
    def magnitude(): Float
        return Math.sqrt(this.x * this.x + this.y * this.y)
    end
    
    def normalize(): Vector2D
        let mag = this.magnitude()
        if mag == 0:
            return Vector2D(0, 0)
        end
        return Vector2D(this.x / mag, this.y / mag)
    end
    
    def dot(other: Vector2D): Float
        return this.x * other.x + this.y * other.y
    end
end

let v1 = Vector2D(3, 4)
let v2 = Vector2D(1, 0)

println("Magnitude: #{v1.magnitude()}")  // 5.0
println("Dot product: #{v1.dot(v2)}")    // 3.0
```

### Configuration Objects
```pf
record DatabaseConfig(
    host: String,
    port: Int,
    database: String,
    username: String
)
    def connectionString(): String
        return "#{this.host}:#{this.port}/#{this.database}"
    end
end

let config = DatabaseConfig(
    "localhost",
    5432,
    "myapp",
    "admin"
)

println(config.connectionString())
```

### Comparison and Equality
```pf
record Range(start: Int, end: Int)
    def contains(value: Int): Bool
        return value >= this.start and value <= this.end
    end
    
    def overlaps(other: Range): Bool
        return this.start <= other.end and other.start <= this.end
    end
    
    def size(): Int
        return this.end - this.start + 1
    end
end

let r1 = Range(1, 10)
let r2 = Range(5, 15)

println(r1.contains(5))      // true
println(r1.overlaps(r2))     // true
println(r1.size())           // 10
```

## Records vs Classes

### When to Use Records
```pf
// ✅ Good: Simple data container
record Point(x: Int, y: Int)
end

// ✅ Good: Immutable value object
record Color(r: Int, g: Int, b: Int)
end
```

### When to Use Classes
```pf
// ❌ Bad: Mutable state needed
// Use class instead
class Counter:
    var count: Int
    
    def init():
        this.count = 0
    end
    
    def increment():
        this.count = this.count + 1
    end
end
```

## Best Practices

### ✅ DO - Use records for immutable data
```pf
record User(id: Int, name: String, email: String)
end
```

### ✅ DO - Add methods for computed properties
```pf
record Circle(radius: Float)
    def area(): Float
        return Math.PI * this.radius * this.radius
    end
end
```

### ✅ DO - Use descriptive field names
```pf
record Order(
    orderId: Int,
    customerId: Int,
    orderDate: Int,
    totalAmount: Float
)
end
```

### ❌ DON'T - Try to modify record fields
```pf
record Point(x: Int, y: Int)
end

let p = Point(1, 2)
// p.x = 5  // Error: Records are immutable
```

### ❌ DON'T - Use records for complex mutable state
```pf
// Bad: Complex state changes
// Use a class instead
record GameState(
    score: Int,
    level: Int,
    lives: Int,
    enemies: Array
)
end
```

## Comparison Table

| Feature | Record | Class |
|---------|--------|-------|
| Constructor | Automatic | Manual |
| Immutability | Immutable | Mutable |
| Fields | Public by default | Can be private |
| Inheritance | No | Yes |
| Use case | Data containers | Complex objects |

## See Also

- [Classes](../definitions/class.md)
- [Enums](enums.md)
- [Type System](../advanced/types.md)
