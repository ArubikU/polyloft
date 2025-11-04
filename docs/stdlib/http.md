# Http Module

The `Http` module provides comprehensive HTTP client and server functionality for making requests and building web services, including support for WebSockets, async/await, advanced routing, and middleware.

## HTTP Client

### Basic Methods

#### `Http.get(url, timeout?)`
Makes a GET request.

**Parameters:**
- `url` (String): URL to request
- `timeout` (Int, optional): Timeout in seconds

**Returns:** Map with response data including `status`, `ok`, `body`, `headers`

```pf
let response = Http.get("https://api.example.com/users", 5)
println(response["body"])
println(response["status"])
println(response["ok"])  // true for 2xx status codes
```

#### `Http.post(url, data, timeout?)`
Makes a POST request.

**Parameters:**
- `url` (String): URL to request
- `data` (Any): Request body (automatically serialized to JSON)
- `timeout` (Int, optional): Timeout in seconds

**Returns:** Map with response data

```pf
let data = {name: "Alice", age: 25}
let response = Http.post("https://api.example.com/users", data, 5)
println(response["status"])
```

#### `Http.put(url, data, timeout?)`
Makes a PUT request.

**Parameters:**
- `url` (String): URL to request
- `data` (Any): Request body
- `timeout` (Int, optional): Timeout in seconds

**Returns:** Map with response data

```pf
let updates = {age: 26}
let response = Http.put("https://api.example.com/users/1", updates, 5)
```

#### `Http.delete(url, timeout?)`
Makes a DELETE request.

**Parameters:**
- `url` (String): URL to request
- `timeout` (Int, optional): Timeout in seconds

**Returns:** Map with response data

```pf
let response = Http.delete("https://api.example.com/users/1", 5)
println("Deleted: #{response["status"]}")
```

### Simplified Request API

#### `Http.request(method, url, data?, timeout?, headers?)`
Makes a custom HTTP request with flexible parameters.

**Overloads:**
- `Http.request(method, url, data)` - Simplest form
- `Http.request(method, url, data, timeout)` - With timeout
- `Http.request(method, url, data, timeout, headers)` - With timeout and headers

**Parameters:**
- `method` (String): HTTP method (GET, POST, PUT, DELETE, PATCH, etc.)
- `url` (String): URL to request
- `data` (Any, optional): Request body
- `timeout` (Int, optional): Timeout in seconds
- `headers` (Map, optional): Custom headers

**Returns:** Map with response data

```pf
// Simple usage
let response = Http.request("PATCH", "https://api.example.com/users/1", {status: "active"})

// With timeout
let response = Http.request("POST", "https://api.example.com/data", body, 10)

// With timeout and headers
let headers = {"Authorization": "Bearer token123"}
let response = Http.request("POST", "https://api.example.com/data", body, 10, headers)
```

### Async HTTP Methods

All HTTP methods have async versions that return Promises, compatible with `await()` and `then()`.

#### `Http.getAsync(url, timeout?)`
Async GET request returning a Promise.

```pf
let promise = Http.getAsync("https://api.example.com/data")
let response = promise.await()
println(response["body"])
```

#### `Http.postAsync(url, data, timeout?)`
Async POST request returning a Promise.

```pf
let promise = Http.postAsync("https://api.example.com/users", {name: "Bob"})
promise.then((response) => do
    println("Created: #{response["status"]}")
end)
```

#### `Http.putAsync(url, data, timeout?)`
Async PUT request returning a Promise.

#### `Http.deleteAsync(url, timeout?)`
Async DELETE request returning a Promise.

#### `Http.requestAsync(method, url, data)`
Async custom request returning a Promise.

```pf
let promise = Http.requestAsync("PATCH", "https://api.example.com/users/1", {active: true})
let response = promise.await()
```

### Response Format

All HTTP client methods return a Map with the following structure:

```pf
{
    status: 200,              // HTTP status code
    ok: true,                 // true for 2xx status codes
    statusText: "OK",         // Status text
    body: {...},              // Response body (auto-parsed from JSON to Map)
    headers: {...}            // Response headers
}
```

## HTTP Server

### `Http.createServer(debug?)`
Creates an HTTP server instance.

**Parameters:**
- `debug` (Bool, optional): Enable debug mode

**Returns:** HttpServer

```pf
let server = Http.createServer()
```

### Server Configuration

#### `server.config(options)`
Configure server-wide settings.

**Parameters:**
- `options` (Map): Configuration options

**Available Options:**
- `cors` (Bool): Enable CORS
- `timeout` (Int): Request timeout in milliseconds
- `jsonLimit` (String): Max JSON body size (e.g., "5MB")
- `defaultHeaders` (Map): Default headers for all responses

```pf
server.config({
    cors: true,
    timeout: 5000,
    jsonLimit: "5MB",
    defaultHeaders: {"X-Powered-By": "Polyloft"}
})
```

### Logging System

#### `server.log(message, level?)`
Log messages with configurable levels.

**Parameters:**
- `message` (String): Message to log
- `level` (String, optional): Log level (debug, info, warn, error). Default: "info"

```pf
server.log("Server started", "info")
server.log("Processing request", "debug")
server.log("Warning: high load", "warn")
server.log("Database error", "error")
```

### Middleware System

#### Global Middleware

Apply middleware to all routes using `server.use()`.

```pf
// Logger middleware
server.use((req, res, next) => do
    server.log("#{req.method} #{req.path}")
    next()
end)
```

#### Route-Specific Middleware

Apply middleware to specific routes by passing an array of middleware functions.

```pf
def requireAuth(req, res, next):
    if req.headers["Authorization"]:
        next()
    else:
        res.status(401).json({error: "Unauthorized"})
    end
end

def logAccess(req, res, next):
    server.log("Access: #{req.path}", "info")
    next()
end

// Single middleware
server.get("/admin", [requireAuth], (req, res) => do
    res.ok({message: "Admin area"})
end)

// Multiple middlewares
server.get("/protected", [requireAuth, logAccess], (req, res) => do
    res.ok({data: "sensitive"})
end)
```

**Middleware Signature:**
Middlewares must have exactly 3 parameters: `(req, res, next)`

### Error Handling

#### `server.onError(handler)`
Register a global error handler.

**Parameters:**
- `handler` (Function): Error handler function `(err, req, res) => {...}`

```pf
server.onError((err, req, res) => do
    server.log("Error: #{err}", "error")
    res.error(500, "Internal Server Error")
end)
```

### Advanced Routing

#### Dynamic Parameters

Capture values from URL paths using `:paramName` syntax.

```pf
server.get("/users/:id", (req, res) => do
    let userId = req.params["id"]
    res.ok({userId: userId})
end)

// Multiple parameters
server.get("/posts/:postId/comments/:commentId", (req, res) => do
    res.ok({
        postId: req.params["postId"],
        commentId: req.params["commentId"]
    })
end)
```

#### Parameter Validation with Regex

Add regex patterns to validate parameters: `:paramName(regex)`

```pf
// Only numeric IDs
server.get("/products/:id([0-9]+)", (req, res) => do
    res.ok({productId: req.params["id"]})
end)

// Alphanumeric usernames
server.get("/usernames/:username([a-zA-Z0-9_]+)", (req, res) => do
    res.ok({username: req.params["username"]})
end)

// File extensions
server.get("/files/:filename([a-zA-Z0-9]+\\.txt)", (req, res) => do
    res.ok({filename: req.params["filename"]})
end)
```

**Validation Behavior:**
- Valid paths match the pattern and execute the handler
- Invalid paths return 404 Not Found

#### Wildcard Routes

Capture remaining path segments using `*paramName` syntax.

```pf
server.get("/static/*filepath", (req, res) => do
    let filepath = req.params["filepath"]
    res.ok({filepath: filepath})
end)

// Example: /static/css/style.css → filepath = "css/style.css"
```

#### Mixed Routes

Combine static and dynamic segments in any order.

```pf
server.get("/api/v1/users/:id/profile", (req, res) => do
    res.ok({userId: req.params["id"]})
end)
```

#### Route Priority

1. **Static routes** - Exact match (O(1) lookup)
2. **Dynamic routes** - Pattern matching (in registration order)
3. **404 Not Found** - No match

### Route Methods

All route methods support both simple and middleware variants:

**`get(path, handler)`** - Register GET route
**`get(path, middlewares, handler)`** - Register GET route with middlewares

```pf
server.get("/users", (req, res) => do
    res.json({users: ["Alice", "Bob"]})
end)

// With middlewares
server.get("/users", [auth, log], (req, res) => do
    res.json({users: ["Alice", "Bob"]})
end)
```

**`post(path, handler)`** / **`post(path, middlewares, handler)`** - Register POST route

**`put(path, handler)`** / **`put(path, middlewares, handler)`** - Register PUT route

**`delete(path, handler)`** / **`delete(path, middlewares, handler)`** - Register DELETE route

### WebSocket Support

#### `server.ws(path, handler)`
Register a WebSocket endpoint for real-time bidirectional communication.

**Parameters:**
- `path` (String): WebSocket route path
- `handler` (Function): WebSocket handler `(socket) => {...}`

```pf
server.ws("/chat", (socket) => do
    // Handle incoming messages
    socket.on("message", (msg) => do
        println("Received: #{msg}")
        socket.send("Echo: #{msg}")
    end)
    
    // Handle connection close
    socket.on("close", () => do
        println("Connection closed")
    end)
    
    // Send welcome message
    socket.send("Welcome to the chat!")
end)
```

#### WebSocket Object Methods

**`socket.send(message)`** - Send text message to client
```pf
socket.send("Hello from server!")
```

**`socket.broadcast(message)`** - Broadcast message to all clients
```pf
socket.broadcast("User joined the chat")
```

**`socket.on(event, handler)`** - Register event handler
```pf
socket.on("message", (msg) => do
    println("Got: #{msg}")
end)
```

**`socket.close()`** - Close the connection
```pf
socket.close()
```

#### WebSocket Events

- `"message"` - Text message received from client
- `"binary"` - Binary data received from client
- `"close"` - Connection closed

#### Multiple WebSocket Endpoints

```pf
// Chat endpoint
server.ws("/chat", chatHandler)

// Notifications endpoint
server.ws("/notifications", notificationHandler)

// Live data endpoint
server.ws("/live", liveHandler)
```

#### WebSocket Client (JavaScript)

```javascript
// Connect to WebSocket
const ws = new WebSocket("ws://localhost:8080/chat");

// Handle connection open
ws.onopen = () => {
    ws.send("Hello from client!");
};

// Handle incoming messages
ws.onmessage = (event) => {
    console.log("Received:", event.data);
};

// Handle connection close
ws.onclose = () => {
    console.log("Disconnected");
};
```

**`listen(port)`** - Start server
```pf
let info = server.listen(8080)
println("Server: #{info.message}")
println("Address: #{info.address}")
```

### Request Object

Properties available in request handlers:

- `req.method` - HTTP method (String)
- `req.url` - Full request URL (String)
- `req.path` - URL path (String)
- `req.query` - Query parameters (Map)
- `req.params` - Route parameters (Map) - captured from dynamic routes
- `req.body` - Request body (Any) - auto-parsed from JSON
- `req.headers` - Request headers (Map)

```pf
server.get("/users/:id", (req, res) => do
    let userId = req.params["id"]       // Route parameter
    let format = req.query["format"]    // Query parameter ?format=json
    let authHeader = req.headers["Authorization"]  // Header
    
    res.json({id: userId, format: format})
end)
```

### Response Object

#### Status and Headers

**`status(code)`** - Set status code (chainable)
```pf
res.status(404).json({error: "Not found"})
```

**`header(name, value)`** - Set header (chainable)
```pf
res.header("Content-Type", "application/json")
res.header("X-Custom", "value")
```

#### Response Methods

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

#### Response Shortcuts

**`res.ok(data)`** - Send 200 OK with JSON data
```pf
res.ok({message: "Success", items: [1, 2, 3]})
```

**`res.created(data)`** - Send 201 Created with JSON data
```pf
res.created({id: 1, name: "New Item"})
```

**`res.noContent()`** - Send 204 No Content
```pf
res.noContent()
```

**`res.notFound(message?)`** - Send 404 Not Found
```pf
res.notFound("Resource not found")
res.notFound()  // Default message
```

**`res.error(code, message)`** - Send custom error response
```pf
res.error(500, "Internal Server Error")
res.error(503, "Service Unavailable")
```

#### Template Rendering

**`res.render(template, data)`** - Render HTML template with data

```pf
server.get("/profile", (req, res) => do
    res.render("profile.html", {
        name: "Alice",
        age: 25,
        role: "Developer"
    })
end)
```

Template variables are replaced using `{{variableName}}` syntax.

## Examples

### Modern REST API with Middleware and Advanced Routing

```pf
let server = Http.createServer()

// Configure server
server.config({
    cors: true,
    timeout: 5000,
    defaultHeaders: {"X-API-Version": "1.0"}
})

// Global middleware
server.use((req, res, next) => do
    server.log("#{req.method} #{req.path}", "info")
    next()
end)

// Auth middleware
def requireAuth(req, res, next):
    let token = req.headers["Authorization"]
    if token:
        next()
    else:
        res.status(401).json({error: "Unauthorized"})
    end
end

let users = [
    {id: 1, name: "Alice", role: "admin"},
    {id: 2, name: "Bob", role: "user"}
]

// Public routes
server.get("/", (req, res) => do
    res.ok({message: "API v1.0", endpoints: ["/users", "/admin"]})
end)

// Get all users
server.get("/users", (req, res) => do
    res.ok(users)
end)

// Get user by ID with validation (numeric only)
server.get("/users/:id([0-9]+)", (req, res) => do
    let id = req.params["id"]
    let user = users.find((u) => u.id == id)
    
    if user:
        res.ok(user)
    else:
        res.notFound("User not found")
    end
end)

// Protected admin route
server.get("/admin", [requireAuth], (req, res) => do
    res.ok({message: "Admin area", users: users})
end)

// Create user
server.post("/users", (req, res) => do
    let newUser = {
        id: users.length() + 1,
        name: req.body["name"],
        role: req.body["role"]
    }
    users = users.concat([newUser])
    res.created(newUser)
end)

// Update user
server.put("/users/:id([0-9]+)", (req, res) => do
    let id = req.params["id"]
    let index = users.findIndex((u) => u.id == id)
    
    if index >= 0:
        users[index]["name"] = req.body["name"]
        res.ok(users[index])
    else:
        res.notFound()
    end
end)

// Delete user
server.delete("/users/:id([0-9]+)", (req, res) => do
    let id = req.params["id"]
    users = users.filter((u) => u.id != id)
    res.noContent()
end)

// Static file serving with wildcard
server.get("/static/*filepath", (req, res) => do
    let filepath = req.params["filepath"]
    res.ok({serving: filepath})
end)

// Error handler
server.onError((err, req, res) => do
    server.log("Error: #{err}", "error")
    res.error(500, "Internal Server Error")
end)

server.listen(8080)
println("API server running on http://localhost:8080")
```

### WebSocket Chat Application

```pf
let server = Http.createServer()

// HTTP endpoint for status
server.get("/", (req, res) => do
    res.html("<h1>Chat Server</h1><p>Connect via WebSocket at ws://localhost:8080/chat</p>")
end)

// WebSocket chat endpoint
server.ws("/chat", (socket) => do
    server.log("New chat connection", "info")
    
    // Send welcome message
    socket.send("Welcome to the chat! Say hello!")
    
    // Handle incoming messages
    socket.on("message", (msg) => do
        server.log("Message: #{msg}", "debug")
        
        // Echo back to sender
        socket.send("You said: #{msg}")
        
        // Broadcast to all clients
        socket.broadcast("Someone said: #{msg}")
    end)
    
    // Handle connection close
    socket.on("close", () => do
        server.log("Connection closed", "info")
    end)
end)

// WebSocket notifications endpoint
server.ws("/notifications", (socket) => do
    socket.send("Connected to notifications")
    
    socket.on("message", (msg) => do
        // Broadcast notification to all
        socket.broadcast("Notification: #{msg}")
    end)
end)

server.listen(8080)
println("Chat server running")
println("HTTP: http://localhost:8080")
println("WS:   ws://localhost:8080/chat")
```

### Async HTTP Client

```pf
// Using async/await
def fetchUserData(userId):
    let promise = Http.getAsync("https://api.example.com/users/#{userId}")
    let response = promise.await()
    
    if response["ok"]:
        return response["body"]
    else:
        println("Error: #{response["status"]}")
        return nil
    end
end

let user = fetchUserData(1)
if user:
    println("User: #{user}")
end

// Using promises with then()
Http.getAsync("https://api.example.com/data")
    .then((response) => do
        println("Status: #{response["status"]}")
        println("Data: #{response["body"]}")
    end)
    .catch((err) => do
        println("Error: #{err}")
    end)

// Parallel async requests
let usersPromise = Http.getAsync("https://api.example.com/users")
let postsPromise = Http.getAsync("https://api.example.com/posts")

let users = usersPromise.await()
let posts = postsPromise.await()

println("Got #{users["body"].length()} users")
println("Got #{posts["body"].length()} posts")
```

### Template Rendering

```pf
let server = Http.createServer()

server.get("/profile/:username", (req, res) => do
    let username = req.params["username"]
    
    res.render("profile.html", {
        username: username,
        role: "Developer",
        joinDate: "2024",
        skills: ["Polyloft", "JavaScript", "Go"]
    })
end)

server.get("/dashboard", (req, res) => do
    res.render("dashboard.html", {
        title: "Dashboard",
        stats: {
            users: 150,
            posts: 320,
            comments: 1250
        }
    })
end)

server.listen(8080)
```

### Wildcard Routes for File Serving

```pf
let server = Http.createServer()

// Serve static files
server.get("/assets/*filepath", (req, res) => do
    let filepath = req.params["filepath"]
    
    // In a real app, you would read the file
    res.ok({
        serving: filepath,
        type: "static asset"
    })
end)

// Examples:
// /assets/css/style.css → filepath = "css/style.css"
// /assets/js/app.js → filepath = "js/app.js"
// /assets/images/logo.png → filepath = "images/logo.png"

server.listen(8080)
```

### Multiple Parameter Routes

```pf
let server = Http.createServer()

// Blog posts with comments
server.get("/posts/:postId/comments/:commentId", (req, res) => do
    let postId = req.params["postId"]
    let commentId = req.params["commentId"]
    
    res.ok({
        post: postId,
        comment: commentId,
        content: "Comment content here"
    })
end)

// User posts with filters
server.get("/users/:userId/posts/:postId", (req, res) => do
    res.ok({
        userId: req.params["userId"],
        postId: req.params["postId"],
        format: req.query["format"]
    })
end)

server.listen(8080)
```

### Simplified HTTP Requests

```pf
// Simple PATCH request
let response = Http.request("PATCH", "https://api.example.com/users/1", {
    status: "active",
    lastLogin: Sys.time()
})

// POST with timeout
let response = Http.request("POST", "https://api.example.com/data", body, 10)

// Custom request with headers
let headers = {
    "Authorization": "Bearer token123",
    "Content-Type": "application/json"
}
let response = Http.request("PUT", "https://api.example.com/users/1", data, 5, headers)

// Response handling
if response["ok"]:
    println("Success: #{response["body"]}")
else:
    println("Failed: #{response["status"]}")
end
```

## Best Practices

### ✅ DO - Use async methods for non-blocking operations
```pf
// Parallel requests
let promise1 = Http.getAsync("https://api.example.com/users")
let promise2 = Http.getAsync("https://api.example.com/posts")

let users = promise1.await()
let posts = promise2.await()
```

### ✅ DO - Use response shortcuts for cleaner code
```pf
server.get("/users/:id", (req, res) => do
    let user = findUser(req.params["id"])
    
    if user:
        res.ok(user)  // Clean and semantic
    else:
        res.notFound("User not found")
    end
end)
```

### ✅ DO - Validate parameters with regex
```pf
// Only allow numeric IDs
server.get("/products/:id([0-9]+)", handler)

// Only allow valid usernames
server.get("/users/:username([a-zA-Z0-9_]{3,20})", handler)
```

### ✅ DO - Use middleware for cross-cutting concerns
```pf
// Authentication
def requireAuth(req, res, next):
    if isAuthenticated(req):
        next()
    else:
        res.status(401).json({error: "Unauthorized"})
    end
end

// Apply to protected routes
server.get("/admin", [requireAuth], handler)
```

### ✅ DO - Configure server-wide settings
```pf
server.config({
    cors: true,
    timeout: 5000,
    defaultHeaders: {"X-API-Version": "2.0"}
})
```

### ✅ DO - Use logging for debugging
```pf
server.log("Starting server", "info")
server.log("Debug info: #{data}", "debug")
server.log("Warning: high load", "warn")
server.log("Error occurred: #{err}", "error")
```

### ✅ DO - Handle errors with global error handler
```pf
server.onError((err, req, res) => do
    server.log("Error: #{err}", "error")
    res.error(500, "Internal Server Error")
end)
```

### ✅ DO - Use WebSockets for real-time features
```pf
// Real-time chat
server.ws("/chat", (socket) => do
    socket.on("message", (msg) => do
        socket.broadcast(msg)
    end)
end)
```

### ✅ DO - Handle errors properly
```pf
try:
    let response = Http.get(url, 5)
    processResponse(response)
catch e:
    println("Request failed: #{e}")
    useDefaultData()
end
```

### ✅ DO - Validate input
```pf
server.post("/users", (req, res) => do
    if not req.body["name"]:
        res.status(400).json({error: "Name required"})
        return
    end
    
    // Process valid input
end)
```

### ✅ DO - Set appropriate status codes
```pf
res.ok(data)              // 200 OK
res.created(created)      // 201 Created
res.noContent()           // 204 No Content
res.status(400).json(error)  // 400 Bad Request
res.notFound("Not found") // 404 Not Found
res.error(500, "Error")   // 500 Internal Error
```

### ✅ DO - Use wildcards for flexible routing
```pf
// Serve all files under /static
server.get("/static/*filepath", (req, res) => do
    serveFile(req.params["filepath"])
end)
```

### ❌ DON'T - Expose sensitive data
```pf
// Bad: Exposing passwords
res.json({user: user, password: user.password})

// Good: Filter sensitive data
res.json({user: {id: user.id, name: user.name}})
```

### ❌ DON'T - Forget to call next() in middleware
```pf
// Bad: Breaks the middleware chain
def myMiddleware(req, res, next):
    doSomething()
    // Missing next() - handler never runs!
end

// Good: Always call next()
def myMiddleware(req, res, next):
    doSomething()
    next()
end
```

### ❌ DON'T - Block with synchronous operations in WebSocket handlers
```pf
// Bad: Blocking operation
server.ws("/chat", (socket) => do
    socket.on("message", (msg) => do
        expensiveComputation()  // Blocks other messages
        socket.send(result)
    end)
end)

// Good: Use async operations
server.ws("/chat", (socket) => do
    socket.on("message", (msg) => do
        let promise = computeAsync()
        let result = promise.await()
        socket.send(result)
    end)
end)
```

## Performance Tips

### Use Route Priority Wisely
Static routes are faster than dynamic routes. Place specific routes before generic ones:

```pf
// Good order
server.get("/users/me", handlerMe)        // Static - checked first
server.get("/users/:id", handlerId)       // Dynamic - checked second

// Bad order
server.get("/users/:id", handlerId)       // Would match /users/me
server.get("/users/me", handlerMe)        // Never reached
```

### Cache Expensive Operations
```pf
let cache = {}

server.get("/data/:id", (req, res) => do
    let id = req.params["id"]
    
    if cache[id]:
        res.ok(cache[id])
    else:
        let data = fetchExpensiveData(id)
        cache[id] = data
        res.ok(data)
    end
end)
```

### Use Async Methods for I/O
```pf
// Better: Non-blocking
let promise = Http.getAsync(url)
let response = promise.await()

// Worse: Blocking
let response = Http.get(url)
```

## Security Recommendations

### Always Validate User Input
```pf
server.post("/users", (req, res) => do
    if not isValidEmail(req.body["email"]):
        res.status(400).json({error: "Invalid email"})
        return
    end
    
    // Process validated input
end)
```

### Use HTTPS in Production
```pf
// Configure HTTPS settings
server.config({
    ssl: true,
    cert: "/path/to/cert.pem",
    key: "/path/to/key.pem"
})
```

### Implement Rate Limiting
```pf
let requestCounts = {}

def rateLimiter(req, res, next):
    let ip = req.headers["X-Forwarded-For"]
    
    if not requestCounts[ip]:
        requestCounts[ip] = 0
    end
    
    requestCounts[ip] = requestCounts[ip] + 1
    
    if requestCounts[ip] > 100:
        res.status(429).json({error: "Too many requests"})
    else:
        next()
    end
end

server.use(rateLimiter)
```

### Sanitize Error Messages
```pf
server.onError((err, req, res) => do
    // Don't expose internal error details
    server.log("Internal error: #{err}", "error")
    res.error(500, "Something went wrong")  // Generic message
end)
```

## See Also

- [Async/Await](../advanced/async-await.md) - Working with Promises and async operations
- [Map Type](../types/map.md) - Understanding Maps for request/response data
- [Array Type](../types/array.md) - Working with arrays in HTTP responses
- [String Type](../types/string.md) - String interpolation and manipulation
- [IO Module](io.md) - File I/O operations

## Feature Summary

### HTTP Client
- ✅ Synchronous methods: `get`, `post`, `put`, `delete`, `request`
- ✅ Async methods: `getAsync`, `postAsync`, `putAsync`, `deleteAsync`, `requestAsync`
- ✅ Simplified request API with flexible parameters
- ✅ Auto-parse JSON responses to Polyloft Maps
- ✅ Standardized response format with `ok` field
- ✅ Timeout support

### HTTP Server
- ✅ Route methods: `get`, `post`, `put`, `delete`
- ✅ Advanced routing: dynamic params, regex validation, wildcards
- ✅ Middleware: global and route-specific
- ✅ Response shortcuts: `ok`, `created`, `noContent`, `notFound`, `error`
- ✅ Error handling: global error handler
- ✅ Configuration: server-wide settings
- ✅ Logging: configurable log levels
- ✅ Template rendering: HTML templates with variable substitution
- ✅ Request timeouts: configurable timeout limits

### WebSocket Support
- ✅ Real-time bidirectional communication
- ✅ Event-driven message handling
- ✅ Multiple WebSocket endpoints
- ✅ Broadcast support
- ✅ Seamless HTTP + WebSocket integration

## Quick Reference

### Client Methods
```pf
Http.get(url, timeout?)
Http.post(url, data, timeout?)
Http.put(url, data, timeout?)
Http.delete(url, timeout?)
Http.request(method, url, data?, timeout?, headers?)

Http.getAsync(url, timeout?)
Http.postAsync(url, data, timeout?)
Http.putAsync(url, data, timeout?)
Http.deleteAsync(url, timeout?)
Http.requestAsync(method, url, data)
```

### Server Methods
```pf
let server = Http.createServer()
server.config(options)
server.use(middleware)
server.onError(handler)
server.log(message, level?)

server.get(path, handler)
server.get(path, middlewares, handler)
server.post(path, handler)
server.post(path, middlewares, handler)
server.put(path, handler)
server.put(path, middlewares, handler)
server.delete(path, handler)
server.delete(path, middlewares, handler)

server.ws(path, handler)
server.listen(port)
```

### Response Methods
```pf
res.status(code)
res.header(name, value)
res.json(data)
res.send(text)
res.html(content)
res.render(template, data)

res.ok(data)
res.created(data)
res.noContent()
res.notFound(message?)
res.error(code, message)
```

### WebSocket Methods
```pf
socket.send(message)
socket.broadcast(message)
socket.on(event, handler)
socket.close()
```

### Route Patterns
```pf
"/users"                    // Static
"/users/:id"                // Dynamic parameter
"/users/:id([0-9]+)"        // With regex validation
"/posts/:pid/comments/:cid" // Multiple parameters
"/static/*filepath"         // Wildcard
"/api/v1/users/:id/profile" // Mixed static/dynamic
```
