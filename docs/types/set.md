# Set Type

Sets in Polyloft are collections of unique values with fast membership testing.

## Creating Sets

### Set Constructor
```pf
let set = Set()
set.add(1)
set.add(2)
set.add(1)  // Duplicate ignored
println(set.size())  // 2
```

### Set with Initial Values
```pf
let set = Set(1, 2, 3, 4, 5)
println(set.size())  // 5
```

### Set from Array
```pf
let arr = [1, 2, 2, 3, 3, 3]
let set = Set(arr)
println(set.size())  // 3 (duplicates removed)
```

### Typed Set
```pf
let set = Set<String>("apple", "banana", "cherry")
```

## Methods

### `add(element)`
Adds an element to the set.

**Parameters:**
- `element` (Any): Element to add

**Returns:** void

```pf
let set = Set()
set.add("apple")
set.add("banana")
set.add("apple")  // Ignored (duplicate)
println(set.size())  // 2
```

### `remove(element)` / `delete(element)`
Removes an element from the set.

**Parameters:**
- `element` (Any): Element to remove

**Returns:** void

```pf
let set = Set(1, 2, 3, 4, 5)
set.remove(3)
println(set.has(3))  // false
```

### `has(element)` / `contains(element)`
Checks if element exists in set.

**Parameters:**
- `element` (Any): Element to check

**Returns:** Bool

```pf
let set = Set("a", "b", "c")
println(set.has("b"))  // true
println(set.has("d"))  // false
```

### `size()` / `length()`
Returns number of elements.

**Returns:** Int

```pf
let set = Set(1, 2, 3, 4, 5)
println(set.size())  // 5
```

### `isEmpty()`
Checks if set is empty.

**Returns:** Bool

```pf
println(Set().isEmpty())        // true
println(Set(1).isEmpty())       // false
```

### `clear()`
Removes all elements.

**Returns:** void

```pf
let set = Set(1, 2, 3)
set.clear()
println(set.isEmpty())  // true
```

### `toArray()`
Converts set to array.

**Returns:** Array

```pf
let set = Set("x", "y", "z")
let arr = set.toArray()
println(arr)  // ["x", "y", "z"]
```

## Set Operations

### Union
```pf
def union(set1, set2):
    let result = Set()
    for item in set1.toArray():
        result.add(item)
    end
    for item in set2.toArray():
        result.add(item)
    end
    return result
end

let a = Set(1, 2, 3)
let b = Set(3, 4, 5)
let u = union(a, b)
println(u.toArray())  // [1, 2, 3, 4, 5]
```

### Intersection
```pf
def intersection(set1, set2):
    let result = Set()
    for item in set1.toArray():
        if set2.has(item):
            result.add(item)
        end
    end
    return result
end

let a = Set(1, 2, 3, 4)
let b = Set(3, 4, 5, 6)
let i = intersection(a, b)
println(i.toArray())  // [3, 4]
```

### Difference
```pf
def difference(set1, set2):
    let result = Set()
    for item in set1.toArray():
        if not set2.has(item):
            result.add(item)
        end
    end
    return result
end

let a = Set(1, 2, 3, 4)
let b = Set(3, 4, 5)
let d = difference(a, b)
println(d.toArray())  // [1, 2]
```

### Subset Check
```pf
def isSubset(set1, set2):
    for item in set1.toArray():
        if not set2.has(item):
            return false
        end
    end
    return true
end

let a = Set(1, 2)
let b = Set(1, 2, 3, 4)
println(isSubset(a, b))  // true
```

## Examples

### Remove Duplicates
```pf
let numbers = [1, 2, 2, 3, 3, 3, 4, 4, 4, 4]
let unique = Set(numbers).toArray()
println(unique)  // [1, 2, 3, 4]
```

### Unique Words
```pf
def getUniqueWords(text):
    let words = text.toLowerCase().split(" ")
    let uniqueSet = Set()
    
    for word in words:
        uniqueSet.add(word)
    end
    
    return uniqueSet
end

let text = "the quick brown fox jumps over the lazy dog"
let unique = getUniqueWords(text)
println("Unique words: #{unique.size()}")
```

### Membership Testing
```pf
let validUsers = Set("alice", "bob", "charlie")

def checkAccess(username):
    if validUsers.has(username):
        println("Access granted")
        return true
    else:
        println("Access denied")
        return false
    end
end

checkAccess("alice")    // Access granted
checkAccess("mallory")  // Access denied
```

### Tag System
```pf
class Article:
    let title: String
    let tags: Set
    
    def init(title: String):
        this.title = title
        this.tags = Set()
    end
    
    def addTag(tag: String):
        this.tags.add(tag)
    end
    
    def hasTag(tag: String): Bool
        return this.tags.has(tag)
    end
    
    def getTags(): Array
        return this.tags.toArray()
    end
end

let article = Article("Polyloft Tutorial")
article.addTag("programming")
article.addTag("tutorial")
article.addTag("polyloft")

println(article.hasTag("tutorial"))  // true
println(article.getTags())
```

### Tracking Visited Items
```pf
let visited = Set()

def processItems(items):
    for item in items:
        if not visited.has(item):
            println("Processing: #{item}")
            processItem(item)
            visited.add(item)
        else:
            println("Skipping (already processed): #{item}")
        end
    end
end

processItems([1, 2, 3])
processItems([2, 3, 4])  // 2 and 3 skipped
```

## Best Practices

### ✅ DO - Use sets for uniqueness
```pf
let uniqueIds = Set()
for item in items:
    uniqueIds.add(item.id)
end
```

### ✅ DO - Use sets for fast lookup
```pf
let allowedValues = Set(1, 2, 3, 4, 5)
if allowedValues.has(value):
    process(value)
end
```

### ✅ DO - Remove duplicates with sets
```pf
let unique = Set(arrayWithDuplicates).toArray()
```

### ❌ DON'T - Use sets when order matters
```pf
// Bad: Sets don't preserve order
let set = Set(3, 1, 2)
// Order not guaranteed

// Good: Use array
let arr = [3, 1, 2]
```

### ❌ DON'T - Try to access by index
```pf
// Bad: Sets don't support indexing
// let first = set[0]  // Error

// Good: Convert to array first
let arr = set.toArray()
let first = arr[0]
```

## See Also

- [Array Type](array.md)
- [Map Type](map.md)
- [Loops](../control-flow/loops.md)
