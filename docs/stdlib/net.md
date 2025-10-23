# Net Module

The Net module provides networking capabilities for building TCP servers and clients.

## Server Operations

### Net.listen(addr)
Creates a TCP server listening on the specified address.

**Parameters**:
- `addr: String` - Address to listen on (e.g., ":8080", "localhost:3000")

**Returns**: `Map` - Server object with methods

**Server Methods**:
- `addr: String` - The actual address the server is listening on
- `accept() -> Map` - Accepts a new connection (blocks until one arrives)
- `close() -> Void` - Closes the server

**Example**:
```polyloft
let server = Net.listen(":8080")
println("Server listening on " + server.addr)

loop
    let conn = server.accept()
    println("New connection from " + conn.remote)
    
    // Handle connection
    conn.send("Welcome!\n")
    conn.close()
    break  // Remove this in production
end
```

## Client Operations

### Net.connect(addr, timeout?)
Connects to a TCP server at the specified address.

**Parameters**:
- `addr: String` - Address to connect to (e.g., "localhost:8080")
- `timeout: Int` (optional) - Connection timeout in seconds (default: 10)

**Returns**: `Map` - Connection object with methods

**Connection Methods**:
- `remote: String` - Remote address
- `local: String` - Local address
- `send(data: String) -> Float` - Sends data, returns bytes sent
- `recv(size?: Int) -> String` - Receives data (default 1024 bytes)
- `close() -> Void` - Closes the connection

**Example**:
```polyloft
let conn = Net.connect("localhost:8080", 5)
println("Connected to #{conn.remote}")

conn.send("Hello, server!\n")
let response = conn.recv()
println("Server said: #{response}")

conn.close()
```

## Practical Examples

### Simple Echo Server
```polyloft
def runEchoServer(port: Int):
    let server = Net.listen(":#{port}")
    println("Echo server listening on #{server.addr}")
    
    loop
        let conn = server.accept()
        println("Client connected: #{conn.remote}")
        
        thread spawn do
            try:
                loop
                    let data = conn.recv()
                    if data == "":
                        break
                    end
                    println("Received: #{data}")
                    conn.send(data)
                end
            catch error:
                println("Error: #{error}")
            finally:
                conn.close()
                println("Client disconnected")
            end
        end
    end
end

runEchoServer(8080)
```

### HTTP-like Server
```polyloft
class SimpleHTTPServer:
    private let port: Int
    private let routes: Map
    
    SimpleHTTPServer(port: Int):
        this.port = port
        this.routes = Map()
    end
    
    def route(path: String, handler):
        this.routes.set(path, handler)
    end
    
    def start():
        let server = Net.listen(":#{this.port}")
        println("HTTP server on #{server.addr}")
        
        loop
            let conn = server.accept()
            thread spawn do
                this.handleRequest(conn)
            end
        end
    end
    
    def handleRequest(conn):
        defer conn.close()
        
        let request = conn.recv(4096)
        let lines = request.split("\n")
        let firstLine = lines[0]
        let parts = firstLine.split(" ")
        
        if parts.length() < 2:
            return
        end
        
        let method = parts[0]
        let path = parts[1]
        
        let handler = this.routes.get(path)
        if handler:
            let response = handler(method, path)
            conn.send(response)
        else:
            conn.send("HTTP/1.1 404 Not Found\r\n\r\nNot Found")
        end
    end
end

let server = SimpleHTTPServer(8080)

server.route("/", (method, path) => {
    return "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, World!"
})

server.route("/api/time", (method, path) => {
    let time = Sys.time()
    return "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n#{time}"
})

server.start()
```

### Chat Server
```polyloft
class ChatServer:
    private let clients: List
    
    ChatServer():
        this.clients = List()
    end
    
    def start(port: Int):
        let server = Net.listen(":#{port}")
        println("Chat server on #{server.addr}")
        
        loop
            let conn = server.accept()
            this.clients.add(conn)
            println("Client joined: #{conn.remote}")
            
            thread spawn do
                this.handleClient(conn)
            end
        end
    end
    
    def handleClient(conn):
        try:
            conn.send("Welcome to chat!\n")
            
            loop
                let msg = conn.recv()
                if msg == "":
                    break
                end
                
                this.broadcast("#{conn.remote}: #{msg}", conn)
            end
        catch error:
            println("Error: #{error}")
        finally:
            this.clients.remove(conn)
            conn.close()
            println("Client left: #{conn.remote}")
        end
    end
    
    def broadcast(message: String, sender):
        for client in this.clients:
            if client != sender:
                try:
                    client.send(message)
                catch error:
                    // Client disconnected
                end
            end
        end
    end
end

let chat = ChatServer()
chat.start(9000)
```

### TCP Client
```polyloft
class TCPClient:
    private let host: String
    private let port: Int
    private let conn
    
    TCPClient(host: String, port: Int):
        this.host = host
        this.port = port
        this.conn = null
    end
    
    def connect():
        this.conn = Net.connect("#{this.host}:#{this.port}")
        println("Connected to #{this.conn.remote}")
    end
    
    def disconnect():
        if this.conn:
            this.conn.close()
            this.conn = null
        end
    end
    
    def send(message: String):
        if not this.conn:
            throw "Not connected"
        end
        return this.conn.send(message)
    end
    
    def receive(size: Int) -> String:
        if not this.conn:
            throw "Not connected"
        end
        return this.conn.recv(size)
    end
end

let client = TCPClient("localhost", 8080)
client.connect()
client.send("Hello, server!\n")
let response = client.receive(1024)
println(response)
client.disconnect()
```

### Port Scanner
```polyloft
def scanPort(host: String, port: Int) -> Bool:
    try:
        let conn = Net.connect("#{host}:#{port}", 1)
        conn.close()
        return true
    catch error:
        return false
    end
end

def scanPorts(host: String, startPort: Int, endPort: Int):
    println("Scanning #{host}...")
    let openPorts = List()
    
    for port in range(startPort, endPort + 1):
        if scanPort(host, port):
            println("Port #{port}: OPEN")
            openPorts.add(port)
        end
    end
    
    println("\nFound #{openPorts.size()} open ports")
    return openPorts
end

scanPorts("localhost", 8000, 9000)
```

### Connection Pool
```polyloft
class ConnectionPool:
    private let host: String
    private let port: Int
    private let maxSize: Int
    private let available: List
    
    ConnectionPool(host: String, port: Int, maxSize: Int):
        this.host = host
        this.port = port
        this.maxSize = maxSize
        this.available = List()
    end
    
    def acquire():
        if this.available.size() > 0:
            return this.available.pop()
        end
        
        if this.available.size() < this.maxSize:
            return Net.connect("#{this.host}:#{this.port}")
        end
        
        throw "Connection pool exhausted"
    end
    
    def release(conn):
        if this.available.size() < this.maxSize:
            this.available.add(conn)
        else:
            conn.close()
        end
    end
    
    def closeAll():
        for conn in this.available:
            conn.close()
        end
        this.available.clear()
    end
end

let pool = ConnectionPool("localhost", 8080, 10)

def useConnection():
    let conn = pool.acquire()
    try:
        conn.send("Request\n")
        let response = conn.recv()
        return response
    finally:
        pool.release(conn)
    end
end
```

### Protocol Handler
```polyloft
class ProtocolHandler:
    def handleConnection(conn):
        defer conn.close()
        
        // Send greeting
        conn.send("HELLO\n")
        
        loop
            let command = conn.recv().trim()
            
            if command == "QUIT":
                conn.send("BYE\n")
                break
            elif command.startsWith("ECHO "):
                let msg = command.substring(5)
                conn.send("#{msg}\n")
            elif command == "TIME":
                let time = Sys.time()
                conn.send("#{time}\n")
            elif command == "PING":
                conn.send("PONG\n")
            else:
                conn.send("ERROR Unknown command\n")
            end
        end
    end
end

let server = Net.listen(":9000")
let handler = ProtocolHandler()

loop
    let conn = server.accept()
    thread spawn do
        handler.handleConnection(conn)
    end
end
```

## Network Best Practices

### Error Handling
Always wrap network operations in try-catch:
```polyloft
try:
    let conn = Net.connect("server.com:8080")
    defer conn.close()
    
    conn.send("data")
    let response = conn.recv()
catch error:
    println("Network error: #{error}")
end
```

### Timeouts
Use timeouts to prevent hanging:
```polyloft
// Connection timeout
let conn = Net.connect("slow-server.com:80", 5)

// Operation timeout using channels
let ch = Channel<String>()
thread spawn do
    let data = conn.recv()
    ch.send(data)
end

select:
    case ch.recv() as data:
        println("Received: #{data}")
    case Sys.sleep(5000):
        println("Timeout after 5 seconds")
end
```

### Resource Cleanup
Always close connections:
```polyloft
def processConnection(addr: String):
    let conn = Net.connect(addr)
    defer conn.close()  // Ensures cleanup
    
    // Work with connection
    conn.send("data")
end
```

### Concurrent Connections
Handle multiple clients with threads:
```polyloft
let server = Net.listen(":8080")

loop
    let conn = server.accept()
    
    thread spawn do
        handleClient(conn)
    end
end
```

## Notes

- TCP only (UDP not yet supported)
- Connections are blocking by default
- Use threads for concurrent connections
- Always close connections to free resources
- Network errors should be caught and handled
- Timeout is only for initial connection, not for I/O operations
- Use defer for guaranteed cleanup

## See Also

- [Channels](../concurrency/channels.md) - For async network patterns
- [Threads](../concurrency/threads.md) - Concurrent connection handling
- [IO Module](io.md) - File I/O operations
- [Exception Handling](../language/exceptions.md) - Error handling
