# Built-in Interfaces

Polyloft provides several built-in interfaces that you can implement in your custom classes to enable specific behaviors and integrate seamlessly with language features like loops, indexing, and destructuring.

## Overview

Built-in interfaces allow your custom objects to:
- Be iterable in `for` loops
- Support destructuring assignments
- Be sliceable with bracket notation
- Work with index operators
- Integrate with collection utilities

## Iterable Interface

The `Iterable` interface allows objects to be used in `for...in` loops.

### Interface Definition
```pf
interface Iterable<T>:
    __length() -> Int
    __get(index: Int) -> T
end
```

### Methods

#### `__length() -> Int`
Returns the number of elements in the collection.

#### `__get(index: Int) -> T`
Returns the element at the specified index (0-based).

### Example Implementation
```pf
class CustomList implements Iterable<any>:
    var items = []
    
    def __length():
        return this.items.length()
    end
    
    def __get(index):
        return this.items[index]
    end
    
    def add(item):
        this.items = this.items.concat([item])
    end
end

// Usage
let list = CustomList()
list.add("apple")
list.add("banana")
list.add("cherry")

for item in list:
    println(item)
end
// Outputs: apple, banana, cherry
```

### Use Cases
- Custom collection classes
- Data structures (trees, graphs, etc.)
- Wrapper classes around external data
- Lazy evaluation sequences

## Unstructured Interface

The `Unstructured` interface enables objects to be destructured into multiple variables.

### Interface Definition
```pf
interface Unstructured:
    __pieces() -> Int
    __get_piece(index: Int) -> any
end
```

### Methods

#### `__pieces() -> Int`
Returns the number of pieces (components) the object can be destructured into.

#### `__get_piece(index: Int) -> any`
Returns the piece at the specified index.

### Example Implementation
```pf
class Point implements Unstructured:
    var x
    var y
    
    Point(px, py):
        this.x = px
        this.y = py
    end
    
    def __pieces():
        return 2
    end
    
    def __get_piece(index):
        if index == 0:
            return this.x
        else:
            if index == 1:
                return this.y
            else:
                throw "Index out of bounds"
            end
        end
    end
    
    // Alias for compatibility: some contexts use __getPiece (camelCase)
    def __getPiece(index):
        return this.__get_piece(index)
    end
end

// Usage with destructuring
let point = Point(10, 20)
let px, py = point
println("x: #{px}, y: #{py}")  // x: 10, y: 20

// Usage in loops with destructuring
let points = [Point(1, 2), Point(3, 4), Point(5, 6)]
for x, y in points:
    println("(#{x}, #{y})")
end
```

### Built-in Example: Pair Class

The built-in `Pair<K,V>` class implements `Unstructured`:

```pf
let pair = Pair("key", "value")
let k, v = pair
println("Key: #{k}, Value: #{v}")

// Works with map iteration
let map = {a: 1, b: 2, c: 3}
for key, value in map:
    // Each iteration destructures a Pair<String, Int>
    println("#{key} = #{value}")
end
```

### Use Cases
- Coordinate classes (Point, Vector)
- Complex number types
- Key-value pairs
- Tuples and multi-value returns
- RGB/RGBA color types

## Sliceable Interface

The `Sliceable` interface allows objects to support slice operations.

### Interface Definition
```pf
interface Sliceable<T>:
    __slice(start: Int, end: Int) -> T
end
```

### Methods

#### `__slice(start: Int, end: Int) -> T`
Returns a slice of the collection from `start` (inclusive) to `end` (exclusive).

### Example Implementation
```pf
class CustomArray implements Sliceable<CustomArray>:
    var data
    
    CustomArray(initialData):
        this.data = initialData
    end
    
    def __slice(start, end):
        let result = CustomArray([])
        for i in range(start, end):
            if i < this.data.length():
                result.data = result.data.concat([this.data[i]])
            end
        end
        return result
    end
    
    def toString():
        return this.data.toString()
    end
end

// Usage
let arr = CustomArray([1, 2, 3, 4, 5])
let sliced = arr.__slice(1, 4)
println(sliced)  // [2, 3, 4]
```

### Use Cases
- Custom array implementations
- String-like objects
- Buffer and byte array classes
- Sequence types

## Indexable Interface

The `Indexable` interface provides bracket operator support for getting, setting, and checking containment.

### Interface Definition
```pf
interface Indexable<K, V>:
    __get(key: K) -> V
    __set(key: K, value: V) -> Void
    __contains(key: K) -> Bool
end
```

### Methods

#### `__get(key: K) -> V`
Returns the value associated with the given key.

#### `__set(key: K, value: V) -> Void`
Sets the value for the given key.

#### `__contains(key: K) -> Bool`
Returns `true` if the key exists in the collection.

### Example Implementation
```pf
class Dictionary implements Indexable<String, any>:
    var entries = {}
    
    def __get(key):
        if this.__contains(key):
            return this.entries[key]
        end
        return nil
    end
    
    def __set(key, value):
        this.entries[key] = value
    end
    
    def __contains(key):
        // Check if key exists in entries
        for k in this.entries:
            if k == key:
                return true
            end
        end
        return false
    end
end

// Usage
let dict = Dictionary()
dict.__set("name", "Alice")
dict.__set("age", 30)

println(dict.__get("name"))        // Alice
println(dict.__contains("age"))    // true
println(dict.__contains("city"))   // false
```

### Use Cases
- Dictionary/map implementations
- Sparse arrays
- Cache implementations
- Configuration objects
- Database row wrappers

## Collection Interface

The `Collection` interface provides a standard API for collection types.

### Interface Definition
```pf
interface Collection<T>:
    size() -> Int
    isEmpty() -> Bool
    add(element: T) -> Void
    remove(element: T) -> Bool
    contains(element: T) -> Bool
    clear() -> Void
    asArray() -> Array<T>
end
```

### Methods

#### `size() -> Int`
Returns the number of elements in the collection.

#### `isEmpty() -> Bool`
Returns `true` if the collection is empty.

#### `add(element: T) -> Void`
Adds an element to the collection.

#### `remove(element: T) -> Bool`
Removes an element from the collection. Returns `true` if the element was found and removed.

#### `contains(element: T) -> Bool`
Returns `true` if the collection contains the specified element.

#### `clear() -> Void`
Removes all elements from the collection.

#### `asArray() -> Array<T>`
Returns an array containing all elements in the collection.

### Example Implementation
```pf
class SimpleSet implements Collection<any>:
    var items = []
    
    def size():
        return this.items.length()
    end
    
    def isEmpty():
        return this.size() == 0
    end
    
    def add(element):
        if !this.contains(element):
            this.items = this.items.concat([element])
        end
    end
    
    def remove(element):
        let newItems = []
        let found = false
        // Iterate through items, skipping the first occurrence of element
        for item in this.items:
            if !found && item == element:
                // Skip first occurrence
                found = true
            else:
                // Keep all other items
                newItems = newItems.concat([item])
            end
        end
        this.items = newItems
        return found
    end
    
    def contains(element):
        for item in this.items:
            if item == element:
                return true
            end
        end
        return false
    end
    
    def clear():
        this.items = []
    end
    
    def asArray():
        return this.items
    end
end

// Usage
let set = SimpleSet()
set.add(1)
set.add(2)
set.add(2)  // Won't add duplicate

println(set.size())           // 2
println(set.contains(1))      // true
println(set.asArray())        // [1, 2]
```

### Use Cases
- Custom collection implementations
- Set implementations
- Queue and stack classes
- Priority queues
- Sorted collections

## Combining Multiple Interfaces

Objects can implement multiple interfaces to provide richer functionality:

```pf
class FlexibleList implements Iterable<any>, Collection<any>, Sliceable<FlexibleList>:
    var data = []
    
    // Iterable methods
    def __length():
        return this.data.length()
    end
    
    def __get(index):
        return this.data[index]
    end
    
    // Collection methods
    def size():
        return this.data.length()
    end
    
    def isEmpty():
        return this.size() == 0
    end
    
    def add(element):
        this.data = this.data.concat([element])
    end
    
    def remove(element):
        let newData = []
        let found = false
        // Iterate through data, skipping the first occurrence of element
        for item in this.data:
            if !found && item == element:
                // Skip first occurrence
                found = true
            else:
                // Keep all other items
                newData = newData.concat([item])
            end
        end
        this.data = newData
        return found
    end
    
    def contains(element):
        for item in this.data:
            if item == element:
                return true
            end
        end
        return false
    end
    
    def clear():
        this.data = []
    end
    
    def asArray():
        return this.data
    end
    
    // Sliceable method
    def __slice(start, end):
        let result = FlexibleList()
        for i in range(start, end):
            if i < this.data.length():
                result.add(this.data[i])
            end
        end
        return result
    end
end

// Usage
let list = FlexibleList()
list.add(10)
list.add(20)
list.add(30)
list.add(40)

// Use as Iterable
for item in list:
    println(item)
end

// Use as Collection
println("Size: #{list.size()}")
println("Contains 20: #{list.contains(20)}")

// Use as Sliceable
let sliced = list.__slice(1, 3)
for item in sliced:
    println(item)  // 20, 30
end
```

## Best Practices

### ✅ DO - Implement complete interfaces
```pf
// Good: All required methods implemented
class MyList implements Iterable<any>:
    def __length():
        return this.items.length()
    end
    
    def __get(index):
        return this.items[index]
    end
end
```

### ✅ DO - Validate indices and bounds
```pf
// Good: Proper validation
func __get_piece(index):
    if index < 0 or index >= this.__pieces():
        throw "Index out of bounds"
    end
    // ... return piece
end
```

### ✅ DO - Return consistent types
```pf
// Good: Consistent return types
func __get(index):
    if index < 0 or index >= this.__length():
        return nil  // Consistent error handling
    end
    return this.items[index]
end
```

### ✅ DO - Document interface implementations
```pf
// Good: Clear documentation
// MyCustomList implements Iterable to enable for-in loops
class MyCustomList implements Iterable<String>:
    // Implementation...
end
```

### ❌ DON'T - Implement partial interfaces
```pf
// Bad: Missing required methods
class BadList implements Iterable<any>:
    def __length():
        return this.items.length()
    end
    // Missing __get method!
end
```

### ❌ DON'T - Return incorrect types
```pf
// Bad: __pieces should return Int
func __pieces():
    return "two"  // Wrong! Should return 2
end
```

### ❌ DON'T - Modify during iteration
```pf
// Bad: Modifying collection during iteration can cause issues
for item in myCollection:
    myCollection.remove(item)  // Dangerous!
end
```

## Common Patterns

### Read-only Wrapper
```pf
class ReadOnlyList implements Iterable<any>:
    final var source
    
    ReadOnlyList(data):
        this.source = data
    end
    
    def __length():
        return this.source.length()
    end
    
    def __get(index):
        return this.source[index]
    end
end
```

### Lazy Evaluation
```pf
class LazyRange implements Iterable<Int>:
    final var start
    final var end
    
    LazyRange(rangeStart, rangeEnd):
        this.start = rangeStart
        this.end = rangeEnd
    end
    
    def __length():
        return this.end - this.start
    end
    
    def __get(index):
        return this.start + index
    end
end
```

### Filtered View
```pf
class FilteredList implements Iterable<any>:
    final var source
    final var predicate
    var cached
    
    FilteredList(src, pred):
        this.source = src
        this.predicate = pred
        this.cached = nil
    end
    
    def buildCache():
        if this.cached == nil:
            this.cached = []
            for item in this.source:
                if this.predicate(item):
                    this.cached = this.cached.concat([item])
                end
            end
        end
    end
    
    def __length():
        this.buildCache()
        return this.cached.length()
    end
    
    def __get(index):
        this.buildCache()
        return this.cached[index]
    end
end
```

## See Also

- [Classes](../definitions/class.md)
- [Loops](../control-flow/loops.md)
- [Arrays](../types/array.md)
- [Maps](../types/map.md)
- [Records](records.md)
