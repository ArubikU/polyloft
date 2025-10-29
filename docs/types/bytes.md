# Bytes Type

The Bytes type in Polyloft represents binary data and provides methods for byte manipulation and conversion.

## Creating Bytes

### Bytes Constructor
```pf
let bytes = Bytes()
```

### From String
```pf
let bytes = Bytes("hello")
println(bytes.size())  // 5
```

### From Array
```pf
let bytes = Bytes([72, 101, 108, 108, 111])  // "Hello" in ASCII
println(bytes.toString())  // "Hello"
```

## Methods

### `size()`
Returns the number of bytes.

**Returns:** Int

```pf
let bytes = Bytes("hello")
println(bytes.size())  // 5
```

### `get(index)`
Gets the byte value at index.

**Parameters:**
- `index` (Int): Position

**Returns:** Int (0-255)

```pf
let bytes = Bytes("ABC")
println(bytes.get(0))  // 65 (ASCII 'A')
println(bytes.get(1))  // 66 (ASCII 'B')
```

### `set(index, value)`
Sets the byte value at index.

**Parameters:**
- `index` (Int): Position
- `value` (Int): Byte value (0-255)

**Returns:** void

```pf
let bytes = Bytes("hello")
bytes.set(0, 72)  // 'H'
println(bytes.toString())  // "Hello"
```

### `toString()`
Converts bytes to string.

**Returns:** String

```pf
let bytes = Bytes([72, 101, 108, 108, 111])
println(bytes.toString())  // "Hello"
```

### `asString()`
Alias for toString().

**Returns:** String

```pf
let bytes = Bytes("data")
println(bytes.asString())  // "data"
```

### `asHex()`
Converts bytes to hexadecimal string.

**Returns:** String

```pf
let bytes = Bytes("hello")
println(bytes.asHex())  // "68656c6c6f"
```

### `asBinary()`
Converts bytes to binary string representation.

**Returns:** String

```pf
let bytes = Bytes([65])  // 'A'
println(bytes.asBinary())  // "01000001"
```

### `asArray()`
Converts bytes to array of integers.

**Returns:** Array

```pf
let bytes = Bytes("ABC")
let arr = bytes.asArray()
println(arr)  // [65, 66, 67]
```

### `asInt()`
Interprets bytes as integer.

**Returns:** Int

```pf
let bytes = Bytes([0, 0, 0, 42])
println(bytes.asInt())  // 42
```

### `asFloat()`
Interprets bytes as float.

**Returns:** Float

```pf
let bytes = Bytes([64, 73, 15, 219])
println(bytes.asFloat())  // Depends on encoding
```

### `asBool()`
Interprets bytes as boolean.

**Returns:** Bool

```pf
let bytes = Bytes([1])
println(bytes.asBool())  // true
```

### `slice(start, end?)`
Extracts a slice of bytes.

**Parameters:**
- `start` (Int): Start index
- `end` (Int, optional): End index

**Returns:** Bytes

```pf
let bytes = Bytes("hello world")
let slice = bytes.slice(0, 5)
println(slice.toString())  // "hello"
```

### `equals(other)`
Checks if two Bytes objects are equal.

**Parameters:**
- `other` (Bytes): Bytes to compare

**Returns:** Bool

```pf
let b1 = Bytes("hello")
let b2 = Bytes("hello")
let b3 = Bytes("world")

println(b1.equals(b2))  // true
println(b1.equals(b3))  // false
```

## Examples

### Binary Data Processing
```pf
def processBinaryData(data):
    let bytes = Bytes(data)
    println("Size: #{bytes.size()}")
    println("Hex: #{bytes.asHex()}")
    println("String: #{bytes.toString()}")
end

processBinaryData("Hello")
```

### Byte Manipulation
```pf
let bytes = Bytes([0, 0, 0, 0])
bytes.set(0, 255)
bytes.set(1, 128)
bytes.set(2, 64)
bytes.set(3, 32)

println("Hex: #{bytes.asHex()}")
println("Array: #{bytes.asArray()}")
```

### Encoding and Decoding
```pf
def encodeString(text):
    let bytes = Bytes(text)
    return bytes.asHex()
end

def decodeHex(hexStr):
    // In practice, you'd convert hex back to bytes
    return Crypto.hexDecode(hexStr)
end

let text = "Hello, World!"
let encoded = encodeString(text)
println("Encoded: #{encoded}")

let decoded = decodeHex(encoded)
println("Decoded: #{decoded}")
```

### Binary File Format
```pf
def createHeader(version, flags):
    let bytes = Bytes([0, 0, 0, 0])
    bytes.set(0, version)
    bytes.set(1, flags)
    return bytes
end

def readHeader(bytes):
    return {
        version: bytes.get(0),
        flags: bytes.get(1)
    }
end

let header = createHeader(1, 255)
println("Header hex: #{header.asHex()}")

let info = readHeader(header)
println("Version: #{info.version}")
println("Flags: #{info.flags}")
```

### Checksum Calculation
```pf
def calculateChecksum(bytes):
    let sum = 0
    for i in range(bytes.size()):
        sum = sum + bytes.get(i)
    end
    return sum % 256
end

let data = Bytes("Hello")
let checksum = calculateChecksum(data)
println("Checksum: #{checksum}")
```

### XOR Encryption
```pf
def xorBytes(bytes, key):
    let result = Bytes()
    for i in range(bytes.size()):
        let b = bytes.get(i)
        let k = key.get(i % key.size())
        result.set(i, b ^ k)  // XOR operation
    end
    return result
end

let plaintext = Bytes("secret")
let key = Bytes("key")

let encrypted = xorBytes(plaintext, key)
println("Encrypted: #{encrypted.asHex()}")

let decrypted = xorBytes(encrypted, key)
println("Decrypted: #{decrypted.toString()}")
```

### Binary Protocol
```pf
class BinaryMessage:
    let bytes: Bytes
    
    BinaryMessage(data: String):
        // Create message: [type][length][data]
        let dataBytes = Bytes(data)
        this.bytes = Bytes()
        this.bytes.set(0, 1)  // Type
        this.bytes.set(1, dataBytes.size())  // Length
        // Copy data bytes...
    end
    
    def getType(): Int
        return this.bytes.get(0)
    end
    
    def getLength(): Int
        return this.bytes.get(1)
    end
    
    def getData(): String
        return this.bytes.slice(2).toString()
    end
end

let msg = BinaryMessage("Hello")
println("Type: #{msg.getType()}")
println("Length: #{msg.getLength()}")
```

### Image Header Parsing
```pf
def parseImageHeader(bytes):
    // Example: simple image format
    return {
        signature: bytes.slice(0, 4).toString(),
        width: bytes.get(4) * 256 + bytes.get(5),
        height: bytes.get(6) * 256 + bytes.get(7),
        bpp: bytes.get(8)
    }
end

let header = Bytes([80, 78, 71, 0, 3, 232, 2, 88, 24])
let info = parseImageHeader(header)
println("Width: #{info.width}")
println("Height: #{info.height}")
println("BPP: #{info.bpp}")
```

### Network Packet
```pf
class NetworkPacket:
    let bytes: Bytes
    
    NetworkPacket():
        this.bytes = Bytes()
    end
    
    def setHeader(type, flags):
        this.bytes.set(0, type)
        this.bytes.set(1, flags)
    end
    
    def setPayload(data):
        let dataBytes = Bytes(data)
        let offset = 2
        for i in range(dataBytes.size()):
            this.bytes.set(offset + i, dataBytes.get(i))
        end
    end
    
    def toHex(): String
        return this.bytes.asHex()
    end
end

let packet = NetworkPacket()
packet.setHeader(1, 0)
packet.setPayload("DATA")
println("Packet: #{packet.toHex()}")
```

## Conversions

### String to Bytes
```pf
let text = "Hello"
let bytes = Bytes(text)
println(bytes.size())  // 5
```

### Bytes to String
```pf
let bytes = Bytes([72, 101, 108, 108, 111])
let text = bytes.toString()
println(text)  // "Hello"
```

### Bytes to Hex
```pf
let bytes = Bytes("ABC")
let hex = bytes.asHex()
println(hex)  // "414243"
```

### Array to Bytes
```pf
let arr = [65, 66, 67]
let bytes = Bytes(arr)
println(bytes.toString())  // "ABC"
```

## Best Practices

### ✅ DO - Use Bytes for binary data
```pf
let binaryData = Bytes(rawData)
let hex = binaryData.asHex()
```

### ✅ DO - Check bounds before access
```pf
if index >= 0 and index < bytes.size():
    let value = bytes.get(index)
end
```

### ✅ DO - Use proper conversions
```pf
let bytes = Bytes("data")
let hex = bytes.asHex()  // For hex
let str = bytes.toString()  // For text
```

### ❌ DON'T - Assume text encoding
```pf
// Bad: May not be valid UTF-8
let bytes = Bytes(binaryData)
let text = bytes.toString()

// Good: Check or use hex
let hex = bytes.asHex()
```

### ❌ DON'T - Modify bytes without bounds check
```pf
// Bad: May cause error
bytes.set(index, value)

// Good: Check first
if index < bytes.size():
    bytes.set(index, value)
end
```

## See Also

- [Crypto Module](../stdlib/crypto.md) - Encoding and hashing
- [String Type](string.md) - Text manipulation
- [Array Type](array.md) - Collections
- [IO Module](../stdlib/io.md) - File operations
