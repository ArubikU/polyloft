# HTTP Module

The HTTP module provides both client and server functionality for HTTP communication.

## HTTP Client

### Static Methods

#### Http.get()
Performs an HTTP GET request.

```polyloft
let response = Http.get("https://api.example.com/data", 5000)
println(response)
```

**Parameters:**
- `url: String` - The URL to request
- `timeout: Int` - Timeout in milliseconds

**Returns:** Response object with status, headers, and body

#### Http.post()
Performs an HTTP POST request.

```polyloft
let data = {"name": "Alice", "age": 30}
let response = Http.post("https://api.example.com/users", data, 5000)
```

**Parameters:**
- `url: String` - The URL to request
- `data: any` - The data to send (will be JSON encoded)
- `timeout: Int` - Timeout in milliseconds

#### Http.put()
Performs an HTTP PUT request.

```polyloft
let data = {"id": 1, "name": "Alice Updated"}
let response = Http.put("https://api.example.com/users/1", data, 5000)
```

**Parameters:**
- `url: String` - The URL to request
- `data: any` - The data to send
- `timeout: Int` - Timeout in milliseconds

#### Http.delete()
Performs an HTTP DELETE request.

```polyloft
let response = Http.delete("https://api.example.com/users/1", 5000)
```

**Parameters:**
- `url: String` - The URL to request
- `timeout: Int` - Timeout in milliseconds

#### Http.request()
Performs a custom HTTP request with full control.

```polyloft
let headers = Map<String, String>()
headers.set("Authorization", "Bearer token123")
headers.set("Content-Type", "application/json")

let response = Http.request(
    "PATCH",
    "https://api.example.com/users/1",
    {"status": "active"},
    5000,
    headers
)
```

**Parameters:**
- `method: String` - HTTP method (GET, POST, PUT, DELETE, PATCH, etc.)
- `url: String` - The URL to request
- `data: any` - The data to send
- `timeout: Int` - Timeout in milliseconds
- `headers: Map` - Request headers

## HTTP Server

### Creating a Server

```polyloft
let server = Http.createServer()
```

Or with debug mode:

```polyloft
let server = Http.createServer(true)
```

### Defining Routes

#### GET Route
```polyloft
server.get("/users", (req, res) => {
    res.status(200).json({"users": ["Alice", "Bob"]})
})
```

#### POST Route
```polyloft
server.post("/users", (req, res) => {
    let name = req.body.get("name")
    res.status(201).send("User created: " + name)
})
```

#### PUT Route
```polyloft
server.put("/users/:id", (req, res) => {
    let id = req.params.get("id")
    res.status(200).send("User " + id + " updated")
})
```

#### DELETE Route
```polyloft
server.delete("/users/:id", (req, res) => {
    let id = req.params.get("id")
    res.status(204).send("")
})
```

### Starting the Server

```polyloft
server.listen(":8080")
println("Server running on http://localhost:8080")
```

## HttpRequest Object

Represents an incoming HTTP request.

**Properties:**
- `method: String` - HTTP method (GET, POST, etc.)
- `path: String` - Request path
- `url: String` - Full URL
- `headers: Map` - Request headers
- `query: Map` - Query parameters
- `body: any` - Request body (parsed if JSON)

## HttpResponse Object

Used to send responses to clients.

### Methods

#### status()
Set the HTTP status code.

```polyloft
res.status(200)
res.status(404)
res.status(500)
```

Returns the response object for chaining.

#### header()
Set a response header.

```polyloft
res.header("Content-Type", "application/json")
res.header("X-Custom-Header", "value")
```

Returns the response object for chaining.

#### send()
Send a text response.

```polyloft
res.send("Hello, World!")
```

#### json()
Send a JSON response.

```polyloft
let data = {"message": "Success", "code": 200}
res.json(data)
```

#### html()
Send an HTML response.

```polyloft
res.html("<h1>Welcome</h1><p>Hello, World!</p>")
```

## Complete Example: REST API

```polyloft
// Create server
let server = Http.createServer()

// In-memory data store
let users = List<Map>(
    {"id": 1, "name": "Alice"},
    {"id": 2, "name": "Bob"}
)

// GET all users
server.get("/users", (req, res) => {
    res.status(200).json(users)
})

// GET single user
server.get("/users/:id", (req, res) => {
    let id = req.params.get("id").toInt()
    let user = users.find((u) => u.get("id") == id)
    
    if user != null:
        res.status(200).json(user)
    else:
        res.status(404).json({"error": "User not found"})
    end
})

// POST create user
server.post("/users", (req, res) => {
    let newUser = req.body
    newUser.set("id", users.size() + 1)
    users.add(newUser)
    res.status(201).json(newUser)
})

// PUT update user
server.put("/users/:id", (req, res) => {
    let id = req.params.get("id").toInt()
    let index = users.findIndex((u) => u.get("id") == id)
    
    if index >= 0:
        users.set(index, req.body)
        res.status(200).json(req.body)
    else:
        res.status(404).json({"error": "User not found"})
    end
})

// DELETE user
server.delete("/users/:id", (req, res) => {
    let id = req.params.get("id").toInt()
    let index = users.findIndex((u) => u.get("id") == id)
    
    if index >= 0:
        users.remove(index)
        res.status(204).send("")
    else:
        res.status(404).json({"error": "User not found"})
    end
})

// Start server
server.listen(":3000")
println("REST API server running on http://localhost:3000")
```

## Error Handling

```polyloft
try:
    let response = Http.get("https://api.example.com/data", 5000)
    println("Success: " + response.toString())
catch error:
    println("HTTP request failed: " + error)
end
```

## Best Practices

1. **Always set timeouts** - Prevent hanging requests
2. **Use appropriate status codes** - Follow HTTP standards
3. **Validate input data** - Check req.body before using
4. **Handle errors gracefully** - Use try-catch blocks
5. **Set proper content types** - Use header() to specify content type
6. **Close connections** - Clean up resources when done

## Common Status Codes

- `200 OK` - Successful GET, PUT, PATCH
- `201 Created` - Successful POST
- `204 No Content` - Successful DELETE
- `400 Bad Request` - Invalid input
- `404 Not Found` - Resource doesn't exist
- `500 Internal Server Error` - Server error

## See Also

- [Net Module](net.md) - TCP socket communication
- [Concurrency](../concurrency/overview.md) - Async operations
- [Error Handling](../language/exceptions.md) - Try-catch blocks
