#!/bin/bash

# 02 - Classes
cat > 02_classes.pf << 'EOF'
// Test: Class constructors and methods
println("=== Classes Test ===")

class Person:
    private let name: String
    private let age: Int
    
    Person(name: String, age: Int):
        this.name = name
        this.age = age
    end
    
    def greet():
        println("Hello, I'm #{this.name}, #{this.age} years old")
    end
    
    def birthday():
        this.age = this.age + 1
    end
    
    def getName() -> String:
        return this.name
    end
end

let person = Person("Alice", 25)
person.greet()
person.birthday()
person.greet()
println("✓ Classes test passed")
EOF

# 03 - Modifiers
cat > 03_modifiers.pf << 'EOF'
// Test: Access modifiers
println("=== Access Modifiers Test ===")

class BankAccount:
    private let balance: Float
    public let accountNumber: String
    
    BankAccount(accountNumber: String, initialBalance: Float):
        this.accountNumber = accountNumber
        this.balance = initialBalance
    end
    
    public def getBalance() -> Float:
        return this.balance
    end
    
    public def deposit(amount: Float):
        this.balance = this.balance + amount
    end
end

let account = BankAccount("ACC123", 1000.0)
println("Account: " + account.accountNumber)
println("Balance: " + account.getBalance().toString())
account.deposit(500.0)
println("After deposit: " + account.getBalance().toString())
println("✓ Access modifiers test passed")
EOF

# 04 - Interfaces
cat > 04_interfaces.pf << 'EOF'
// Test: Interfaces
println("=== Interfaces Test ===")

interface Drawable:
    def draw(): Void
end

class Circle implements Drawable:
    private let radius: Float
    
    Circle(radius: Float):
        this.radius = radius
    end
    
    def draw():
        println("Drawing circle with radius #{this.radius}")
    end
end

let circle = Circle(5.0)
circle.draw()
println("✓ Interfaces test passed")
EOF

# 05 - Generics
cat > 05_generics.pf << 'EOF'
// Test: Generic classes
println("=== Generics Test ===")

class Box<T>:
    private let value: T
    
    Box(value: T):
        this.value = value
    end
    
    def get() -> T:
        return this.value
    end
end

let intBox = Box<Int>(42)
println("Int box: " + intBox.get().toString())

let strBox = Box<String>("Hello")
println("String box: " + strBox.get())
println("✓ Generics test passed")
EOF

# 06 - Math
cat > 06_math.pf << 'EOF'
// Test: Math module
println("=== Math Module Test ===")

println("Math.PI = " + Math.PI.toString())
println("Math.E = " + Math.E.toString())
println("Math.sqrt(9) = " + Math.sqrt(9).toString())
println("Math.pow(2, 3) = " + Math.pow(2, 3).toString())
println("Math.abs(-5.5) = " + Math.abs(-5.5).toString())
println("Math.floor(3.7) = " + Math.floor(3.7).toString())
println("Math.ceil(3.2) = " + Math.ceil(3.2).toString())
println("Math.min(5, 3) = " + Math.min(5, 3).toString())
println("Math.max(5, 3) = " + Math.max(5, 3).toString())
println("✓ Math module test passed")
EOF

# 07 - Sys
cat > 07_sys.pf << 'EOF'
// Test: Sys module
println("=== Sys Module Test ===")

let start = Sys.time()
println("Current time: " + start)

println("Sleeping 100ms...")
Sys.sleep(100)

let elapsed = Sys.time() - start
println("Elapsed: " + elapsed + "ms")

let random1 = Sys.random()
println("Random: " + random1.toString())
println("✓ Sys module test passed")
EOF

# 08 - Collections
cat > 08_collections.pf << 'EOF'
// Test: Arrays and collections
println("=== Collections Test ===")

let numbers = [1, 2, 3, 4, 5]
println("Original: " + numbers.toString())

let doubled = numbers.map((x: Int) -> Int => x * 2)
println("Doubled: " + doubled.toString())

let evens = numbers.filter((x: Int) -> Boolean => x % 2 == 0)
println("Evens: " + evens.toString())

println("✓ Collections test passed")
EOF

# 09 - Control Flow
cat > 09_control_flow.pf << 'EOF'
// Test: Control flow
println("=== Control Flow Test ===")

// if-elif-else
let age = 18
if age >= 18:
    println("Adult")
end

// for loop
for i in range(5):
    print(i.toString() + " ")
end
println("")

// while loop
let count = 0
while count < 3:
    println("Count: " + count.toString())
    count = count + 1
end

println("✓ Control flow test passed")
EOF

# 10 - Functions
cat > 10_functions.pf << 'EOF'
// Test: Functions
println("=== Functions Test ===")

def greet(name: String):
    println("Hello, #{name}!")
end

greet("World")

def add(a: Int, b: Int) -> Int:
    return a + b
end

let sum = add(5, 3)
println("5 + 3 = " + sum.toString())

let multiply = (x: Int, y: Int) -> Int => x * y
println("4 * 5 = " + multiply(4, 5).toString())
println("✓ Functions test passed")
EOF

echo "All test files created"
