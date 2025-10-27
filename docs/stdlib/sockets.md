# Sockets

Las clases `Socket` y `ServerSocket` proporcionan funcionalidades para comunicación de red TCP, similar a Python.

## Socket - Cliente TCP

### Constructor

```polyloft
let socket = Socket()
```

### connect(host: String, port: Int, timeout: Int = 10) -> Bool

Conecta el socket a un host y puerto.

```polyloft
let socket = Socket()
if socket.connect("example.com", 80, 5) {
    Sys.println("Connected!")
} else {
    Sys.println("Connection failed")
}
```

### send(data: String) -> Int

Envía datos como string. Retorna el número de bytes enviados.

```polyloft
let socket = Socket()
socket.connect("localhost", 8080)
let sent = socket.send("GET / HTTP/1.1\r\n\r\n")
Sys.println("Sent " + sent + " bytes")
```

### sendBytes(data: Bytes) -> Int

Envía datos como bytes.

```polyloft
let socket = Socket()
socket.connect("localhost", 8080)
let data = Bytes.fromString("Hello")
socket.sendBytes(data)
```

### recv(size: Int = 1024, timeout: Int = 5) -> String

Recibe datos como string.

```polyloft
let socket = Socket()
socket.connect("localhost", 8080)
socket.send("HELLO\n")
let response = socket.recv(1024, 10)
Sys.println("Received: " + response)
```

### recvBytes(size: Int = 1024, timeout: Int = 5) -> Bytes

Recibe datos como bytes.

```polyloft
let socket = Socket()
socket.connect("localhost", 8080)
let data = socket.recvBytes(512, 5)
```

### close() -> Void

Cierra la conexión del socket.

```polyloft
socket.close()
```

### setReadTimeout(timeout: Int) -> Void

Establece el timeout de lectura en segundos.

```polyloft
socket.setReadTimeout(30)
```

### setWriteTimeout(timeout: Int) -> Void

Establece el timeout de escritura en segundos.

```polyloft
socket.setWriteTimeout(30)
```

### Campos Públicos

- `connected: Bool` - Indica si el socket está conectado
- `remoteAddr: String` - Dirección remota
- `localAddr: String` - Dirección local

## ServerSocket - Servidor TCP

### Constructor

```polyloft
let server = ServerSocket()
```

### bind(host: String, port: Int) -> Bool

Enlaza el servidor a un host y puerto para escuchar conexiones.

```polyloft
let server = ServerSocket()
if server.bind("0.0.0.0", 8080) {
    Sys.println("Server listening on port 8080")
} else {
    Sys.println("Failed to bind")
}
```

### accept() -> Socket

Acepta una conexión entrante y retorna un `Socket` conectado al cliente.

```polyloft
let server = ServerSocket()
server.bind("0.0.0.0", 8080)

Sys.println("Waiting for connection...")
let client = server.accept()
Sys.println("Client connected from: " + client.remoteAddr)
```

### close() -> Void

Cierra el servidor.

```polyloft
server.close()
```

### Campos Públicos

- `listening: Bool` - Indica si el servidor está escuchando
- `address: String` - Dirección del servidor

## Ejemplos

### Cliente Simple

```polyloft
class HttpClient {
    static func get(host: String, path: String): String {
        let socket = Socket()
        
        if !socket.connect(host, 80, 10) {
            return "Connection failed"
        }
        
        // Enviar petición HTTP
        let request = "GET " + path + " HTTP/1.1\r\n"
        request = request + "Host: " + host + "\r\n"
        request = request + "Connection: close\r\n\r\n"
        
        socket.send(request)
        
        // Recibir respuesta
        let response = ""
        let chunk = socket.recv(4096, 5)
        while chunk != "" {
            response = response + chunk
            chunk = socket.recv(4096, 1)
        }
        
        socket.close()
        return response
    }
}

// Usar el cliente
let response = HttpClient.get("example.com", "/")
Sys.println(response)
```

### Servidor Echo Simple

```polyloft
class EchoServer {
    static func start(port: Int): Void {
        let server = ServerSocket()
        
        if !server.bind("0.0.0.0", port) {
            Sys.println("Failed to start server")
            return
        }
        
        Sys.println("Echo server listening on port " + port)
        
        while true {
            // Aceptar cliente
            let client = server.accept()
            Sys.println("Client connected: " + client.remoteAddr)
            
            // Echo loop
            let data = client.recv(1024, 30)
            while data != "" {
                Sys.println("Received: " + data)
                client.send(data)  // Echo back
                data = client.recv(1024, 5)
            }
            
            Sys.println("Client disconnected")
            client.close()
        }
        
        server.close()
    }
}

// Iniciar servidor
EchoServer.start(9999)
```

### Cliente y Servidor con Bytes

```polyloft
// Servidor que recibe bytes
class ByteServer {
    static func start(port: Int): Void {
        let server = ServerSocket()
        server.bind("0.0.0.0", port)
        
        Sys.println("Byte server listening...")
        
        let client = server.accept()
        let data = client.recvBytes(100, 10)
        
        Sys.println("Received " + data.size() + " bytes")
        Sys.println("Hex: " + data.toHex())
        Sys.println("Text: " + data.toString())
        
        client.close()
        server.close()
    }
}

// Cliente que envía bytes
class ByteClient {
    static func send(host: String, port: Int): Void {
        let socket = Socket()
        socket.connect(host, port, 5)
        
        let data = Bytes.fromString("Hello from client!")
        socket.sendBytes(data)
        
        socket.close()
    }
}
```

### Chat Simple

```polyloft
class ChatServer {
    static func start(port: Int): Void {
        let server = ServerSocket()
        server.bind("0.0.0.0", port)
        
        Sys.println("Chat server started on port " + port)
        
        let client = server.accept()
        Sys.println("Client joined: " + client.remoteAddr)
        
        while true {
            let message = client.recv(1024, 60)
            if message == "" {
                break
            }
            
            Sys.println("Client: " + message)
            
            // Echo con prefijo
            client.send("Server received: " + message)
        }
        
        client.close()
        server.close()
    }
}
```

## Comparación con Python

La API es muy similar a la de Python:

```python
# Python
import socket

# Cliente
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect(("localhost", 8080))
s.send(b"Hello")
data = s.recv(1024)
s.close()

# Servidor
server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.bind(("0.0.0.0", 8080))
server.listen(1)
client, addr = server.accept()
data = client.recv(1024)
client.close()
server.close()
```

```polyloft
# Polyloft
// Cliente
let s = Socket()
s.connect("localhost", 8080)
s.send("Hello")
let data = s.recv(1024)
s.close()

// Servidor
let server = ServerSocket()
server.bind("0.0.0.0", 8080)
let client = server.accept()
let data = client.recv(1024)
client.close()
server.close()
```
