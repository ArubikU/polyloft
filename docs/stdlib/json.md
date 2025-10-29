# JSON Support

Polyloft provides JSON serialization and deserialization through Map and Array methods.

## Map Serialization

### `map.serialize()`
Converts a Map to a JSON string.

**Returns:** String

```pf
let data = {
    name: "Alice",
    age: 25,
    active: true
}

let json = data.serialize()
println(json)  // {"name":"Alice","age":25,"active":true}
```

### `Map.deserialize(jsonString)`
Parses a JSON string into a Map.

**Parameters:**
- `jsonString` (String): JSON data to parse

**Returns:** Map

```pf
let jsonStr = '{"name":"Bob","age":30}'
let data = Map.deserialize(jsonStr)

println(data.name)  // "Bob"
println(data.age)   // 30
```

## Examples

### Serializing Nested Data
```pf
let user = {
    id: 123,
    profile: {
        name: "Alice",
        email: "alice@example.com"
    },
    tags: ["admin", "user"]
}

let json = user.serialize()
println(json)
```

### Parsing JSON Response
```pf
def fetchUser(userId):
    let response = Http.get("https://api.example.com/users/#{userId}")
    let data = Map.deserialize(response.body)
    return data
end

let user = fetchUser(123)
println("User: #{user.name}")
```

### Working with JSON Arrays
```pf
let users = [
    {name: "Alice", age: 25},
    {name: "Bob", age: 30}
]

// Serialize each user
for user in users:
    let json = user.serialize()
    println(json)
end
```

### JSON Configuration Files
```pf
def loadConfig(filename):
    let content = IO.readFile(filename)
    return Map.deserialize(content)
end

def saveConfig(filename, config):
    let json = config.serialize()
    IO.writeFile(filename, json)
end

let config = {
    host: "localhost",
    port: 8080,
    debug: true
}

saveConfig("config.json", config)
let loaded = loadConfig("config.json")
println("Port: #{loaded.port}")
```

### API Request with JSON
```pf
let requestData = {
    username: "alice",
    password: "secret123"
}

let response = Http.post(
    "https://api.example.com/login",
    requestData.serialize()
)

let result = Map.deserialize(response.body)
if result.success:
    println("Login successful!")
end
```

### Data Transformation
```pf
def transformUser(user):
    return {
        id: user.id,
        fullName: user.firstName + " " + user.lastName,
        email: user.email
    }
end

let rawUser = Map.deserialize('{"id":1,"firstName":"Alice","lastName":"Smith","email":"alice@example.com"}')
let transformed = transformUser(rawUser)
println(transformed.serialize())
```

## Best Practices

### ✅ DO - Validate JSON before parsing
```pf
def safeParseJSON(jsonStr):
    try:
        return Map.deserialize(jsonStr)
    catch e:
        println("Invalid JSON: #{e}")
        return {}
    end
end
```

### ✅ DO - Handle nested structures
```pf
let data = {
    user: {
        name: "Alice",
        settings: {
            theme: "dark"
        }
    }
}

let json = data.serialize()
```

### ❌ DON'T - Assume structure exists
```pf
// Bad: May error if key doesn't exist
let name = data.user.profile.name

// Good: Check first
if data.has("user") and data.user.has("profile"):
    let name = data.user.profile.name
end
```

## Helper Function

If you want a global `JSON.stringify` style function, you can create one:

```pf
class JSON:
    static def stringify(data):
        if Sys.type(data) == "Map":
            return data.serialize()
        else:
            // For other types, convert to map first
            let map = {value: data}
            return map.serialize()
        end
    end
    
    static def parse(jsonString):
        return Map.deserialize(jsonString)
    end
end

// Usage
let data = {name: "Alice", age: 25}
let json = JSON.stringify(data)
println(json)

let parsed = JSON.parse(json)
println(parsed.name)
```

## See Also

- [Map Type](../types/map.md)
- [Http Module](http.md)
- [IO Module](io.md)
