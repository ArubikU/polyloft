# Map Type

Maps in Polyloft are key-value data structures (dictionaries/hash tables).

## Creating Maps

### Map Literals
```pf
let empty = {}
let person = {name: "Alice", age: 25, city: "NYC"}
let mixed = {count: 42, active: true, items: [1, 2, 3]}
```

### Map Constructor
```pf
let map = Map()
map.put("key1", "value1")
map.put("key2", "value2")
```

### With Type Parameters
```pf
let map = Map<String, Int>()
map.put("one", 1)
map.put("two", 2)
```

## Access Methods

### `get(key, default?)`
Gets value for key, returns default if not found.

**Parameters:**
- `key` (Any): Key to lookup
- `default` (Any, optional): Default value

**Returns:** Any

```pf
let person = {name: "Alice", age: 25}
println(person.get("name"))           // "Alice"
println(person.get("email", "N/A"))   // "N/A"
```

### Bracket Access
```pf
let map = {x: 10, y: 20}
println(map["x"])  // 10
println(map.x)     // 10 (dot notation)
```

### `has(key)` / `hasKey(key)`
Checks if key exists.

**Parameters:**
- `key` (Any): Key to check

**Returns:** Bool

```pf
let map = {a: 1, b: 2}
println(map.has("a"))   // true
println(map.has("c"))   // false
```

## Modification Methods

### `set(key, value)` / `put(key, value)`
Sets a key-value pair.

**Parameters:**
- `key` (Any): Key
- `value` (Any): Value

**Returns:** void

```pf
let map = {}
map.set("name", "Alice")
map.put("age", 25)
println(map)  // {name: "Alice", age: 25}
```

### Bracket Assignment
```pf
let map = {}
map["x"] = 10
map.y = 20  // Dot notation
println(map)  // {x: 10, y: 20}
```

### `remove(key)` / `delete(key)`
Removes a key-value pair.

**Parameters:**
- `key` (Any): Key to remove

**Returns:** void

```pf
let map = {a: 1, b: 2, c: 3}
map.remove("b")
println(map)  // {a: 1, c: 3}
```

### `clear()`
Removes all entries.

**Returns:** void

```pf
let map = {a: 1, b: 2}
map.clear()
println(map.isEmpty())  // true
```

## Query Methods

### `size()` / `length()`
Returns number of entries.

**Returns:** Int

```pf
let map = {a: 1, b: 2, c: 3}
println(map.size())    // 3
println(map.length())  // 3
```

### `isEmpty()`
Checks if map is empty.

**Returns:** Bool

```pf
println({}.isEmpty())        // true
println({a: 1}.isEmpty())    // false
```

### `keys()`
Returns array of all keys.

**Returns:** Array

```pf
let map = {name: "Alice", age: 25, city: "NYC"}
let keys = map.keys()
println(keys)  // ["name", "age", "city"]
```

### `values()`
Returns array of all values.

**Returns:** Array

```pf
let map = {a: 1, b: 2, c: 3}
let vals = map.values()
println(vals)  // [1, 2, 3]
```

### `entries()`
Returns array of [key, value] pairs.

**Returns:** Array

```pf
let map = {x: 10, y: 20}
let entries = map.entries()
// [["x", 10], ["y", 20]]

for entry in entries:
    println("#{entry[0]}: #{entry[1]}")
end
```

## Iteration

### For-In Loop
```pf
let person = {name: "Alice", age: 25, city: "NYC"}

for key in person:
    println("#{key}: #{person[key]}")
end
// Outputs:
// name: Alice
// age: 25
// city: NYC
```

### Iterating Entries
```pf
let map = {a: 1, b: 2, c: 3}

for entry in map.entries():
    let key = entry[0]
    let value = entry[1]
    println("#{key} = #{value}")
end
```

## Examples

### Building Maps Dynamically
```pf
let config = {}
config.set("host", "localhost")
config.set("port", 8080)
config.set("debug", true)

println(config)
// {host: "localhost", port: 8080, debug: true}
```

### Merging Maps
```pf
def mergeMaps(map1, map2):
    let result = {}
    
    for key in map1:
        result[key] = map1[key]
    end
    
    for key in map2:
        result[key] = map2[key]
    end
    
    return result
end

let defaults = {timeout: 1000, retries: 3}
let custom = {timeout: 5000}
let config = mergeMaps(defaults, custom)
println(config)  // {timeout: 5000, retries: 3}
```

### Counting Occurrences
```pf
def countChars(text):
    let counts = {}
    
    for char in text:
        if counts.has(char):
            counts[char] = counts[char] + 1
        else:
            counts[char] = 1
        end
    end
    
    return counts
end

let text = "hello"
let counts = countChars(text)
println(counts)  // {h: 1, e: 1, l: 2, o: 1}
```

### Grouping Data
```pf
def groupBy(items, key):
    let groups = {}
    
    for item in items:
        let groupKey = item[key]
        
        if not groups.has(groupKey):
            groups[groupKey] = []
        end
        
        groups[groupKey] = groups[groupKey].concat([item])
    end
    
    return groups
end

let users = [
    {name: "Alice", role: "admin"},
    {name: "Bob", role: "user"},
    {name: "Charlie", role: "admin"}
]

let byRole = groupBy(users, "role")
println(byRole.admin.length())  // 2
println(byRole.user.length())   // 1
```

### Caching Pattern
```pf
let cache = {}

def fetchWithCache(id):
    if cache.has(id):
        println("Cache hit for #{id}")
        return cache[id]
    end
    
    println("Cache miss for #{id}, fetching...")
    let data = expensiveFetch(id)
    cache[id] = data
    return data
end

let data1 = fetchWithCache(1)  // Cache miss
let data2 = fetchWithCache(1)  // Cache hit
```

### Configuration Management
```pf
let config = {
    database: {
        host: "localhost",
        port: 5432,
        name: "mydb"
    },
    server: {
        port: 8080,
        host: "0.0.0.0"
    },
    features: {
        auth: true,
        logging: true
    }
}

// Access nested values
println(config.database.host)  // "localhost"
println(config.server.port)    // 8080
println(config.features.auth)  // true
```

### Transforming Maps
```pf
def mapValues(map, fn):
    let result = {}
    
    for key in map:
        result[key] = fn(map[key])
    end
    
    return result
end

let numbers = {a: 1, b: 2, c: 3}
let doubled = mapValues(numbers, (x) => x * 2)
println(doubled)  // {a: 2, b: 4, c: 6}
```

### Filtering Maps
```pf
def filterMap(map, predicate):
    let result = {}
    
    for key in map:
        if predicate(key, map[key]):
            result[key] = map[key]
        end
    end
    
    return result
end

let data = {a: 5, b: 12, c: 8, d: 3}
let filtered = filterMap(data, (k, v) => v > 5)
println(filtered)  // {b: 12, c: 8}
```

### Default Values Pattern
```pf
def getWithDefault(map, key, defaultValue):
    if map.has(key):
        return map[key]
    end
    return defaultValue
end

let settings = {theme: "dark"}
let theme = getWithDefault(settings, "theme", "light")
let language = getWithDefault(settings, "language", "en")
```

## Nested Maps

### Deep Access
```pf
let data = {
    user: {
        profile: {
            name: "Alice",
            age: 25
        },
        settings: {
            theme: "dark"
        }
    }
}

println(data.user.profile.name)     // "Alice"
println(data.user.settings.theme)   // "dark"
```

### Safe Navigation
```pf
def safeGet(map, path):
    let parts = path.split(".")
    let current = map
    
    for part in parts:
        if current.has(part):
            current = current[part]
        else:
            return nil
        end
    end
    
    return current
end

let data = {user: {name: "Alice"}}
println(safeGet(data, "user.name"))    // "Alice"
println(safeGet(data, "user.email"))   // nil
```

## Best Practices

### ✅ DO - Use consistent key types
```pf
// Good: All string keys
let map = {"name": "Alice", "age": 25}
```

### ✅ DO - Check for key existence
```pf
if map.has("key"):
    let value = map["key"]
    processValue(value)
end
```

### ✅ DO - Use descriptive keys
```pf
let user = {
    firstName: "Alice",
    lastName: "Smith",
    emailAddress: "alice@example.com"
}
```

### ❌ DON'T - Mix access patterns inconsistently
```pf
// Confusing
map["key1"] = value1
map.key2 = value2

// Better: Pick one style
map.key1 = value1
map.key2 = value2
```

### ❌ DON'T - Forget to handle missing keys
```pf
// Bad: May error if key doesn't exist
let value = map["missingKey"]

// Good: Check first or use default
let value = map.get("missingKey", defaultValue)
```

## See Also

- [Array Type](array.md)
- [String Type](string.md)
- [Set Type](set.md)
- [Loops](../control-flow/loops.md)
