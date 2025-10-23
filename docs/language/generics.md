# Generics

Generics allow you to write flexible, reusable code that works with multiple types while maintaining type safety. Polyloft's generic system is inspired by Java's generics.

## Basic Generics

### Generic Classes

```polyloft
class Box<T>:
    private let value: T
    
    Box(value: T):
        this.value = value
    end
    
    def get() -> T:
        return this.value
    end
    
    def set(value: T):
        this.value = value
    end
end

let intBox = Box<Int>(42)
let strBox = Box<String>("Hello")

println(intBox.get())  // 42
println(strBox.get())  // "Hello"
```

### Generic Functions

```polyloft
def first<T>(list: List<T>) -> T:
    return list.get(0)
end

let numbers = List<Int>(1, 2, 3)
let firstNum = first(numbers)  // Type: Int

let words = List<String>("hello", "world")
let firstWord = first(words)  // Type: String
```

## Type Parameters

### Single Type Parameter

```polyloft
class Container<T>:
    private let items: List<T>
    
    Container():
        this.items = List<T>()
    end
    
    def add(item: T):
        this.items.add(item)
    end
end
```

### Multiple Type Parameters

```polyloft
class Pair<K, V>:
    private let key: K
    private let value: V
    
    Pair(key: K, value: V):
        this.key = key
        this.value = value
    end
    
    def getKey() -> K:
        return this.key
    end
    
    def getValue() -> V:
        return this.value
    end
end

let pair = Pair<String, Int>("age", 30)
```

## Bounded Type Parameters

### Upper Bounds (extends)

```polyloft
class NumberBox<T extends Number>:
    private let value: T
    
    NumberBox(value: T):
        this.value = value
    end
    
    def double() -> Number:
        return this.value * 2
    end
end

let intBox = NumberBox<Int>(42)       // OK
let floatBox = NumberBox<Float>(3.14) // OK
// let strBox = NumberBox<String>("hi") // Error: String doesn't extend Number
```

### Multiple Bounds

```polyloft
interface Comparable:
    def compareTo(other): Int
end

interface Serializable:
    def serialize(): String
end

class SortableContainer<T extends Comparable & Serializable>:
    def sort(items: List<T>):
        // Can call compareTo (from Comparable)
        // Can call serialize (from Serializable)
    end
end
```

## Wildcards

### Unbounded Wildcard (?)

```polyloft
def printList(list: List<?>):
    for item in list:
        println(item)
    end
end

let intList = List<Int>(1, 2, 3)
let strList = List<String>("a", "b", "c")

printList(intList)   // OK
printList(strList)   // OK
```

### Upper Bounded Wildcard (? extends T)

```polyloft
def sumNumbers(list: List<? extends Number>) -> Float:
    let sum = 0.0
    for num in list:
        sum = sum + num.toFloat()
    end
    return sum
end

let ints = List<Int>(1, 2, 3)
let floats = List<Float>(1.5, 2.5, 3.5)

println(sumNumbers(ints))    // 6.0
println(sumNumbers(floats))  // 7.5
```

### Lower Bounded Wildcard (? super T)

```polyloft
def addIntegers(list: List<? super Int>):
    list.add(1)
    list.add(2)
    list.add(3)
end

let intList = List<Int>()
let numList = List<Number>()
let objList = List<Object>()

addIntegers(intList)   // OK: List<Int>
addIntegers(numList)   // OK: List<Number> is super of Int
addIntegers(objList)   // OK: List<Object> is super of Int
```

## Variance

### Covariance (out)

Producer - can only produce values of type T:

```polyloft
class Producer<out T>:
    private let value: T
    
    Producer(value: T):
        this.value = value
    end
    
    def produce() -> T:
        return this.value
    end
    
    // Cannot consume: def consume(value: T) - Error!
end

let intProducer: Producer<Int> = Producer<Int>(42)
let numProducer: Producer<Number> = intProducer  // OK: covariant
```

### Contravariance (in)

Consumer - can only consume values of type T:

```polyloft
class Consumer<in T>:
    def consume(value: T):
        println("Consumed: #{value}")
    end
    
    // Cannot produce: def produce() -> T - Error!
end

let numConsumer: Consumer<Number> = Consumer<Number>()
let intConsumer: Consumer<Int> = numConsumer  // OK: contravariant
```

### Invariance (default)

Neither covariant nor contravariant:

```polyloft
class Container<T>:
    private let value: T
    
    Container(value: T):
        this.value = value
    end
    
    def set(value: T):
        this.value = value
    end
    
    def get() -> T:
        return this.value
    end
end

let intContainer: Container<Int> = Container<Int>(0)
// let numContainer: Container<Number> = intContainer  // Error: invariant
```

## Built-in Generic Collections

### List<T>

```polyloft
let numbers = List<Int>(1, 2, 3)
numbers.add(4)
numbers.add(5)

let first = numbers.get(0)  // Type: Int

// Generic methods
let doubled = numbers.map((x: Int) -> Int => x * 2)
let evens = numbers.filter((x: Int) -> Boolean => x % 2 == 0)
```

### Set<T>

```polyloft
let unique = Set<String>("apple", "banana", "apple")
println(unique.size())  // 2

unique.add("cherry")
println(unique.contains("banana"))  // true
```

### Map<K, V>

```polyloft
let scores = Map<String, Int>()
scores.set("Alice", 95)
scores.set("Bob", 87)

let aliceScore = scores.get("Alice")  // Type: Int
```

### Deque<T>

```polyloft
let deque = Deque<Int>()
deque.addFront(1)
deque.addBack(2)

let front = deque.removeFront()  // Type: Int
```

## Generic Constraints in Practice

### Example: Comparable Elements

```polyloft
def findMax<T extends Comparable>(list: List<T>) -> T:
    if list.isEmpty():
        throw "Empty list"
    end
    
    let max = list.get(0)
    for item in list:
        if item.compareTo(max) > 0:
            max = item
        end
    end
    return max
end
```

### Example: Numeric Operations

```polyloft
def average<T extends Number>(list: List<T>) -> Float:
    let sum = 0.0
    for num in list:
        sum = sum + num.toFloat()
    end
    return sum / list.size()
end

let ints = List<Int>(1, 2, 3, 4, 5)
println(average(ints))  // 3.0
```

### Example: Generic Builder

```polyloft
class Builder<T>:
    private let items: List<T>
    
    Builder():
        this.items = List<T>()
    end
    
    def add(item: T) -> Builder<T>:
        this.items.add(item)
        return this  // Method chaining
    end
    
    def build() -> List<T>:
        return this.items
    end
end

let numbers = Builder<Int>()
    .add(1)
    .add(2)
    .add(3)
    .build()
```

## Type Erasure

Like Java, Polyloft uses type erasure:

```polyloft
// At runtime, type parameters are erased
let intList = List<Int>(1, 2, 3)
let strList = List<String>("a", "b", "c")

// Both are just "List" at runtime
// Type checking happens at compile time
```

## Generic Best Practices

### 1. Use Descriptive Type Parameters

```polyloft
// Good
class Map<Key, Value>
class Result<Success, Error>

// Avoid
class Map<A, B>
```

### 2. Prefer Bounded Wildcards for APIs

```polyloft
// Producer: use extends
def processList(list: List<? extends Item>)

// Consumer: use super
def addItems(list: List<? super Item>)
```

### 3. Use Invariant Types for Mutable Data

```polyloft
// Mutable container should be invariant
class MutableBox<T>:
    def set(value: T)
    def get() -> T
end
```

### 4. Keep Generic Methods Simple

```polyloft
// Good: simple and clear
def first<T>(list: List<T>) -> T

// Avoid: overly complex
def complex<T extends A & B, U extends C<T>, V super U>()
```

## Common Generic Patterns

### Option/Maybe Type

```polyloft
class Option<T>:
    private let value: T
    private let present: Boolean
    
    Option(value: T):
        this.value = value
        this.present = value != null
    end
    
    def isPresent() -> Boolean:
        return this.present
    end
    
    def get() -> T:
        if not this.present:
            throw "No value present"
        end
        return this.value
    end
    
    def orElse(default: T) -> T:
        return this.present ? this.value : default
    end
end
```

### Result Type

```polyloft
class Result<T, E>:
    private let value: T
    private let error: E
    private let success: Boolean
    
    static def ok(value: T) -> Result<T, E>:
        let result = Result<T, E>()
        result.value = value
        result.success = true
        return result
    end
    
    static def error(error: E) -> Result<T, E>:
        let result = Result<T, E>()
        result.error = error
        result.success = false
        return result
    end
end
```

### Generic Factory

```polyloft
interface Factory<T>:
    def create() -> T
end

class IntFactory implements Factory<Int>:
    def create() -> Int:
        return 0
    end
end
```

## Limitations

1. Cannot create generic arrays: `new T[10]` - not allowed
2. Cannot use primitives as type arguments in some contexts
3. Cannot query runtime type of type parameters
4. Type erasure limits reflection capabilities

## See Also

- [Classes and Objects](classes-and-objects.md) - Class basics
- [Interfaces](interfaces.md) - Interface definitions
- [Variance](../advanced/variance.md) - Advanced variance concepts
- [Collections](../stdlib/collections.md) - Generic collections
