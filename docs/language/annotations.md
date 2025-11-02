# Annotations

Annotations in Polyloft provide metadata about methods and classes. They are denoted by the `@` symbol followed by the annotation name. Annotations document intent, enable compile-time checks, and can influence behavior.

## Syntax

```pf
class ClassName:
    @AnnotationName
    def methodName():
        // method body
    end
end
```

Multiple annotations can be applied:

```pf
class Example:
    @Override
    @Deprecated
    def someMethod():
        // method body
    end
end
```

## @Override Annotation

The `@Override` annotation marks methods that override parent class methods or implement interface methods.

### Basic Override
```pf
class Animal:
    def speak():
        println("Some sound")
    end
end

class Dog < Animal:
    @Override
    def speak():
        println("Woof!")
    end
end

let dog = Dog()
dog.speak()  // Woof!
```

### With Return Types
```pf
class Shape:
    def area() -> Float
        return 0.0
    end
end

class Rectangle < Shape:
    let width: Float
    let height: Float
    
    Rectangle(w: Float, h: Float):
        this.width = w
        this.height = h
    end
    
    @Override
    def area() -> Float
        return this.width * this.height
    end
end
```

### Interface Implementation
```pf
interface Drawable:
    draw() -> void
    clear() -> void
end

class Canvas implements Drawable:
    @Override
    def draw():
        println("Drawing on canvas")
    end
    
    @Override
    def clear():
        println("Clearing canvas")
    end
end
```

### Template Method Pattern
```pf
class DataProcessor:
    def process(data):
        this.validate(data)
        let result = this.transform(data)
        this.save(result)
        return result
    end
    
    protected def validate(data):
        // Default validation
    end
    
    protected def transform(data):
        return data
    end
    
    protected def save(data):
        // Default save
    end
end

class JsonProcessor < DataProcessor:
    @Override
    protected def validate(data):
        if not data.contains("{"):
            throw "Invalid JSON"
        end
    end
    
    @Override
    protected def transform(data):
        return JSON.parse(data)
    end
    
    @Override
    protected def save(data):
        IO.writeFile("output.json", JSON.stringify(data))
    end
end
```

## Benefits of @Override

1. **Documentation**: Makes override intent clear
2. **Maintainability**: Helps when refactoring parent classes
3. **Code Review**: Easier to understand inheritance hierarchies
4. **Safety**: Can enable compiler checks (future feature)

## When to Use @Override

### ✅ DO use @Override when:

**Overriding parent methods:**
```pf
class Vehicle:
    def startEngine():
        println("Engine starting")
    end
end

class ElectricVehicle < Vehicle:
    @Override
    def startEngine():
        println("Electric motor starting")
    end
end
```

**Implementing interface methods:**
```pf
interface Closeable:
    close() -> void
end

class FileHandle implements Closeable:
    @Override
    def close():
        println("Closing file")
    end
end
```

**Overriding protected methods:**
```pf
class Base:
    protected def initialize():
        println("Base init")
    end
end

class Extended < Base:
    @Override
    protected def initialize():
        super.initialize()
        println("Extended init")
    end
end
```

### ❌ DON'T use @Override when:

**Defining new methods:**
```pf
class Example:
    // Wrong: Not overriding anything
    // @Override
    def newMethod():
        println("This is a new method")
    end
end
```

**Private methods (can't be overridden):**
```pf
class Example:
    // Wrong: Private methods can't be overridden
    // @Override
    private def helper():
        println("Private helper")
    end
end
```

## Annotation Naming

Annotations follow these conventions:

1. **PascalCase**: Start with capital letter
2. **No spaces**: Single or compound words
3. **Case-insensitive**: Normalized internally

```pf
@Override  // Standard form
@override  // Also works (normalized to Override)
@OVERRIDE  // Also works (normalized to Override)
```

## Multiple Inheritance with @Override

When inheriting from multiple sources:

```pf
interface Saveable:
    save() -> void
end

interface Loadable:
    load() -> void
end

class Document:
    def close():
        println("Closing document")
    end
end

class TextDocument < Document implements Saveable, Loadable:
    @Override
    def save():
        println("Saving text")
    end
    
    @Override
    def load():
        println("Loading text")
    end
    
    @Override
    def close():
        this.save()
        super.close()
    end
end
```

## Combining with Access Modifiers

```pf
class Parent:
    protected def process():
        println("Parent processing")
    end
    
    public def execute():
        println("Parent executing")
    end
end

class Child < Parent:
    @Override
    protected def process():  // Keep same visibility
        super.process()
        println("Child processing")
    end
    
    @Override
    public def execute():  // Keep same visibility
        super.execute()
        println("Child executing")
    end
end
```

## Common Patterns

### Strategy Pattern
```pf
interface PaymentStrategy:
    pay(amount: Float) -> void
end

class CreditCardPayment implements PaymentStrategy:
    let cardNumber
    
    CreditCardPayment(number):
        this.cardNumber = number
    end
    
    @Override
    def pay(amount: Float):
        println("Charging $#{amount} to card")
    end
end

class PayPalPayment implements PaymentStrategy:
    let email
    
    PayPalPayment(userEmail):
        this.email = userEmail
    end
    
    @Override
    def pay(amount: Float):
        println("Charging $#{amount} via PayPal")
    end
end
```

### Observer Pattern
```pf
interface Observer:
    update(data: any) -> void
end

class EmailObserver implements Observer:
    let email
    
    EmailObserver(userEmail):
        this.email = userEmail
    end
    
    @Override
    def update(data: any):
        println("Email to #{this.email}: #{data}")
    end
end
```

### Template Method
```pf
class GameLevel:
    def play():
        this.initialize()
        this.runGameLoop()
        this.cleanup()
    end
    
    protected def initialize():
        println("Level starting")
    end
    
    protected def runGameLoop():
        println("Running game")
    end
    
    protected def cleanup():
        println("Level ending")
    end
end

class BossLevel < GameLevel:
    @Override
    protected def initialize():
        super.initialize()
        println("Boss music playing")
    end
    
    @Override
    protected def runGameLoop():
        println("Epic boss battle!")
    end
    
    @Override
    protected def cleanup():
        println("Victory!")
        super.cleanup()
    end
end
```

## Best Practices

### ✅ DO - Always use @Override for clarity
```pf
class Manager < Employee:
    @Override  // Clear intent
    def getInfo() -> String
        return "Manager: #{this.name}"
    end
end
```

### ✅ DO - Use with super calls
```pf
class Extended < Base:
    @Override
    def initialize():
        super.initialize()  // Call parent first
        println("Extended initialization")
    end
end
```

### ✅ DO - Document complex overrides
```pf
class CachedStore < DataStore:
    @Override
    def query(criteria):
        // Check cache first, then delegate
        let cached = this.checkCache(criteria)
        if cached != nil:
            return cached
        end
        let result = super.query(criteria)
        this.cacheResult(criteria, result)
        return result
    end
end
```

### ❌ DON'T - Forget @Override in deep hierarchies
```pf
class A:
    def method():
        println("A")
    end
end

class B < A:
    @Override
    def method():
        println("B")
    end
end

class C < B:
    @Override  // Good: Shows this overrides B.method()
    def method():
        println("C")
    end
end
```

## Future Annotations

The annotation system is extensible. Potential future annotations include:

```pf
// Deprecation warnings
@Deprecated
def oldMethod():
    println("This method will be removed")
end

// Performance hints
@Inline
def fastOperation():
    return this.value * 2
end

// Documentation metadata
@Description("Processes user input and returns validated result")
def processInput(input) -> String
    return input.trim()
end
```

## See Also

- [Access Modifiers](access-modifiers.md)
- [Classes](../definitions/class.md)
- [Interfaces](interfaces.md)
- [Method Overriding](../definitions/class.md#method-overriding)
