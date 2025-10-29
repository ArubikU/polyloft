# Http Module

The `Http` module provides HTTP client and server functionality for making requests and building web services.

## HTTP Client

### `Http.get(url, options?)`
Makes a GET request.

**Parameters:**
- `url` (String): URL to request
- `options` (Map, optional): Request options

**Returns:** Map with response data

```pf
let response = Http.get("https://api.example.com/users")
println(response.body)
println(response.status)
```

### `Http.post(url, body, options?)`
Makes a POST request.

**Parameters:**
- `url` (String): URL to request
- `body` (Any): Request body
- `options` (Map, optional): Request options

**Returns:** Map with response data

```pf
let data = {name: "Alice", age: 25}
let response = Http.post("https://api.example.com/users", data)
println(response.status)
```

### `Http.put(url, body, options?)`
Makes a PUT request.

**Parameters:**
- `url` (String): URL to request
- `body` (Any): Request body
- `options` (Map, optional): Request options

**Returns:** Map with response data

```pf
let updates = {age: 26}
let response = Http.put("https://api.example.com/users/1", updates)
```

### `Http.delete(url, options?)`
Makes a DELETE request.

**Parameters:**
- `url` (String): URL to request
- `options` (Map, optional): Request options

**Returns:** Map with response data

```pf
let response = Http.delete("https://api.example.com/users/1")
println("Deleted: #{response.status}")
```

### `Http.request(method, url, options)`
Makes a custom HTTP request.

**Parameters:**
- `method` (String): HTTP method
- `url` (String): URL to request
- `options` (Map): Request options

**Returns:** Map with response data

```pf
let response = Http.request("PATCH", "https://api.example.com/users/1", {
    body: {status: "active"},
    headers: {"Authorization": "Bearer token123"}
})
```

## HTTP Server

### `Http.createServer()`
Creates an HTTP server instance.

**Returns:** HttpServer

```pf
let server = Http.createServer()
```

### HttpServer Methods

**`get(path, handler)`** - Register GET route
```pf
server.get("/users", (req, res) => do
    res.json({users: ["Alice", "Bob"]})
end)
```

**`post(path, handler)`** - Register POST route
```pf
server.post("/users", (req, res) => do
    let user = req.body
    // Save user
    res.status(201).json({id: 1, name: user.name})
end)
```

**`put(path, handler)`** - Register PUT route
```pf
server.put("/users/:id", (req, res) => do
    let id = req.params.id
    let updates = req.body
    // Update user
    res.json({id: id, updated: true})
end)
```

**`delete(path, handler)`** - Register DELETE route
```pf
server.delete("/users/:id", (req, res) => do
    let id = req.params.id
    // Delete user
    res.status(204).send()
end)
```

**`listen(port)`** - Start server
```pf
let info = server.listen(8080)
println("Server running on port #{info.port}")
```

### Request Object

Properties available in request handlers:

- `req.method` - HTTP method
- `req.url` - Request URL
- `req.path` - URL path
- `req.query` - Query parameters (Map)
- `req.params` - Route parameters (Map)
- `req.body` - Request body
- `req.headers` - Request headers (Map)

```pf
server.get("/users/:id", (req, res) => do
    let userId = req.params.id
    let format = req.query.format
    let authHeader = req.headers.Authorization
    
    // Process request
    res.json({id: userId, format: format})
end)
```

### Response Object

**`status(code)`** - Set status code
```pf
res.status(404).json({error: "Not found"})
```

**`header(name, value)`** - Set header
```pf
res.header("Content-Type", "application/json")
res.header("X-Custom", "value")
```

**`json(data)`** - Send JSON response
```pf
res.json({message: "Success", data: items})
```

**`send(content)`** - Send text response
```pf
res.send("Hello, World!")
```

**`html(content)`** - Send HTML response
```pf
res.html("<h1>Welcome</h1><p>Hello from Polyloft!</p>")
```

## Examples

### Simple API Server
```pf
let server = Http.createServer()

let users = [
    {id: 1, name: "Alice"},
    {id: 2, name: "Bob"}
]

// Get all users
server.get("/users", (req, res) => do
    res.json(users)
end)

// Get single user
server.get("/users/:id", (req, res) => do
    let id = req.params.id
    let user = users.find((u) => u.id == id)
    
    if user:
        res.json(user)
    else:
        res.status(404).json({error: "User not found"})
    end
end)

// Create user
server.post("/users", (req, res) => do
    let newUser = {
        id: users.length() + 1,
        name: req.body.name
    }
    users = users.concat([newUser])
    res.status(201).json(newUser)
end)

server.listen(8080)
println("Server running on http://localhost:8080")
```

### HTTP Client Example
```pf
def fetchUserData(userId):
    try:
        let response = Http.get("https://api.example.com/users/#{userId}")
        
        if response.status == 200:
            return response.body
        else:
            println("Error: #{response.status}")
            return nil
        end
    catch e:
        println("Request failed: #{e}")
        return nil
    end
end

let user = fetchUserData(1)
if user:
    println("User: #{user.name}")
end
```

### REST API with CRUD
```pf
let server = Http.createServer()
let data = []

// Create
server.post("/items", (req, res) => do
    let item = {
        id: data.length() + 1,
        name: req.body.name,
        value: req.body.value
    }
    data = data.concat([item])
    res.status(201).json(item)
end)

// Read all
server.get("/items", (req, res) => do
    res.json(data)
end)

// Read one
server.get("/items/:id", (req, res) => do
    let id = req.params.id
    let item = data.find((i) => i.id == id)
    
    if item:
        res.json(item)
    else:
        res.status(404).json({error: "Not found"})
    end
end)

// Update
server.put("/items/:id", (req, res) => do
    let id = req.params.id
    let index = data.findIndex((i) => i.id == id)
    
    if index >= 0:
        data[index].name = req.body.name
        data[index].value = req.body.value
        res.json(data[index])
    else:
        res.status(404).json({error: "Not found"})
    end
end)

// Delete
server.delete("/items/:id", (req, res) => do
    let id = req.params.id
    let index = data.findIndex((i) => i.id == id)
    
    if index >= 0:
        data = data.filter((i) => i.id != id)
        res.status(204).send()
    else:
        res.status(404).json({error: "Not found"})
    end
end)

server.listen(3000)
```

### Middleware Pattern
```pf
let server = Http.createServer()

// Logger middleware
def logRequest(req, res, next):
    println("[#{Sys.time()}] #{req.method} #{req.path}")
    next()
end

// Auth middleware
def requireAuth(req, res, next):
    let token = req.headers.Authorization
    
    if token and validateToken(token):
        next()
    else:
        res.status(401).json({error: "Unauthorized"})
    end
end

// Apply middleware
server.use(logRequest)

// Protected routes
server.get("/admin", requireAuth, (req, res) => do
    res.json({message: "Admin area"})
end)

server.listen(8080)
```

### Async HTTP Requests
```pf
let promise = async(() => do
    return Http.get("https://api.example.com/data")
end)

promise.then((response) => do
    println("Status: #{response.status}")
    println("Data: #{response.body}")
end).catch((err) => do
    println("Error: #{err}")
end)
```

## Best Practices

### ✅ DO - Handle errors properly
```pf
try:
    let response = Http.get(url)
    processResponse(response)
catch e:
    println("Request failed: #{e}")
    useDefaultData()
end
```

### ✅ DO - Validate input
```pf
server.post("/users", (req, res) => do
    if not req.body.name:
        res.status(400).json({error: "Name required"})
        return
    end
    
    // Process valid input
end)
```

### ✅ DO - Set appropriate status codes
```pf
res.status(200).json(data)      // OK
res.status(201).json(created)   // Created
res.status(204).send()          // No Content
res.status(400).json(error)     // Bad Request
res.status(404).json(error)     // Not Found
res.status(500).json(error)     // Internal Error
```

### ❌ DON'T - Expose sensitive data
```pf
// Bad: Exposing passwords
res.json({user: user, password: user.password})

// Good: Filter sensitive data
res.json({user: {id: user.id, name: user.name}})
```

## See Also

- [Async/Await](../advanced/async-await.md)
- [Map Type](../types/map.md)
- [IO Module](io.md)
