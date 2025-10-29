# Crypto Module

The `Crypto` module provides cryptographic hashing and encoding functions.

## Hash Functions

### `Crypto.md5(data)`
Computes MD5 hash of data.

**Parameters:**
- `data` (String): Data to hash

**Returns:** String (hex-encoded hash)

```pf
let hash = Crypto.md5("hello world")
println(hash)  // 5eb63bbbe01eeed093cb22bb8f5acdc3
```

### `Crypto.sha1(data)`
Computes SHA-1 hash of data.

**Parameters:**
- `data` (String): Data to hash

**Returns:** String (hex-encoded hash)

```pf
let hash = Crypto.sha1("hello world")
println(hash)  // 2aae6c35c94fcfb415dbe95f408b9ce91ee846ed
```

### `Crypto.sha256(data)`
Computes SHA-256 hash of data.

**Parameters:**
- `data` (String): Data to hash

**Returns:** String (hex-encoded hash)

```pf
let hash = Crypto.sha256("hello world")
println(hash)
// b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9
```

### `Crypto.sha512(data)`
Computes SHA-512 hash of data.

**Parameters:**
- `data` (String): Data to hash

**Returns:** String (hex-encoded hash)

```pf
let hash = Crypto.sha512("hello world")
println(hash)
// 309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f...
```

## Encoding Functions

### `Crypto.base64Encode(data)`
Encodes data to Base64.

**Parameters:**
- `data` (String): Data to encode

**Returns:** String (Base64-encoded)

```pf
let encoded = Crypto.base64Encode("hello world")
println(encoded)  // aGVsbG8gd29ybGQ=
```

### `Crypto.base64Decode(data)`
Decodes Base64 data.

**Parameters:**
- `data` (String): Base64-encoded data

**Returns:** String (decoded data)

```pf
let decoded = Crypto.base64Decode("aGVsbG8gd29ybGQ=")
println(decoded)  // hello world
```

### `Crypto.hexEncode(data)`
Encodes data to hexadecimal.

**Parameters:**
- `data` (String): Data to encode

**Returns:** String (hex-encoded)

```pf
let encoded = Crypto.hexEncode("hello")
println(encoded)  // 68656c6c6f
```

### `Crypto.hexDecode(data)`
Decodes hexadecimal data.

**Parameters:**
- `data` (String): Hex-encoded data

**Returns:** String (decoded data)

```pf
let decoded = Crypto.hexDecode("68656c6c6f")
println(decoded)  // hello
```

## Examples

### Password Hashing
```pf
def hashPassword(password, salt):
    let combined = password + salt
    return Crypto.sha256(combined)
end

let password = "mySecretPass"
let salt = "randomSalt123"
let hashed = hashPassword(password, salt)
println("Hashed password: #{hashed}")
```

### File Integrity Check
```pf
def getFileHash(filename):
    let content = IO.readFile(filename)
    return Crypto.sha256(content)
end

let hash1 = getFileHash("document.txt")
// ... file might be modified ...
let hash2 = getFileHash("document.txt")

if hash1 == hash2:
    println("File unchanged")
else:
    println("File was modified!")
end
```

### API Token Generation
```pf
def generateToken(userId, timestamp):
    let data = "#{userId}:#{timestamp}"
    let hash = Crypto.sha256(data)
    return Crypto.base64Encode(hash)
end

let token = generateToken(123, Sys.time())
println("API Token: #{token}")
```

### Data Encoding for URL
```pf
def encodeForUrl(data):
    return Crypto.base64Encode(data)
end

def decodeFromUrl(encoded):
    return Crypto.base64Decode(encoded)
end

let original = "user@example.com"
let encoded = encodeForUrl(original)
println("Encoded: #{encoded}")

let decoded = decodeFromUrl(encoded)
println("Decoded: #{decoded}")
```

### Checksum Verification
```pf
def verifyChecksum(data, expectedHash):
    let actualHash = Crypto.md5(data)
    return actualHash == expectedHash
end

let data = "important data"
let checksum = Crypto.md5(data)

// Later verification
if verifyChecksum(data, checksum):
    println("Checksum valid - data is intact")
else:
    println("Checksum invalid - data corrupted!")
end
```

### Session Token
```pf
def createSession(username):
    let timestamp = Sys.time()
    let sessionData = "#{username}:#{timestamp}"
    let hash = Crypto.sha256(sessionData)
    return {
        token: hash,
        username: username,
        created: timestamp
    }
end

let session = createSession("alice")
println("Session token: #{session.token}")
```

### Message Signing
```pf
def signMessage(message, secret):
    let data = message + secret
    return Crypto.sha256(data)
end

def verifySignature(message, signature, secret):
    let expectedSignature = signMessage(message, secret)
    return signature == expectedSignature
end

let secret = "mySecretKey"
let message = "Hello, World!"
let signature = signMessage(message, secret)

// Verify
if verifySignature(message, signature, secret):
    println("Signature valid")
else:
    println("Signature invalid")
end
```

### Binary Data Encoding
```pf
def encodeData(data):
    // Encode to hex for storage/transmission
    return Crypto.hexEncode(data)
end

def decodeData(encoded):
    return Crypto.hexDecode(encoded)
end

let data = "Binary data here"
let encoded = encodeData(data)
println("Encoded: #{encoded}")

let decoded = decodeData(encoded)
println("Decoded: #{decoded}")
```

## Hash Comparison

### Hash Strength
```pf
let data = "test data"

println("MD5:    #{Crypto.md5(data)}")
println("SHA1:   #{Crypto.sha1(data)}")
println("SHA256: #{Crypto.sha256(data)}")
println("SHA512: #{Crypto.sha512(data)}")
```

### Choosing Hash Function
- **MD5**: Fast but cryptographically broken. Use only for checksums, not security.
- **SHA1**: Deprecated for security. Use SHA-256 or higher.
- **SHA256**: Good balance of speed and security. Recommended for most use cases.
- **SHA512**: Maximum security but slower. Use for high-security applications.

## Best Practices

### ✅ DO - Use SHA-256 or SHA-512 for security
```pf
// Good: Strong hash
let hash = Crypto.sha256(password)
```

### ✅ DO - Add salt to password hashes
```pf
def hashPassword(password):
    let salt = generateSalt()
    let combined = password + salt
    return {
        hash: Crypto.sha256(combined),
        salt: salt
    }
end
```

### ✅ DO - Use Base64 for binary data in text
```pf
let binaryData = getBinaryData()
let encoded = Crypto.base64Encode(binaryData)
sendViaJson({data: encoded})
```

### ❌ DON'T - Use MD5 for passwords
```pf
// Bad: MD5 is not secure for passwords
let hash = Crypto.md5(password)

// Good: Use SHA-256 or better
let hash = Crypto.sha256(password)
```

### ❌ DON'T - Store passwords without hashing
```pf
// Bad: Plain text password
let user = {password: "secret123"}

// Good: Hashed password
let user = {passwordHash: Crypto.sha256("secret123")}
```

## Common Use Cases

### API Authentication
```pf
def createApiKey(userId, secret):
    let data = "#{userId}:#{Sys.time()}"
    let signature = Crypto.sha256(data + secret)
    return Crypto.base64Encode(signature)
end
```

### Data Integrity
```pf
def saveWithChecksum(filename, content):
    let checksum = Crypto.sha256(content)
    IO.writeFile(filename, content)
    IO.writeFile(filename + ".sha256", checksum)
end

def verifyIntegrity(filename):
    let content = IO.readFile(filename)
    let storedHash = IO.readFile(filename + ".sha256")
    let actualHash = Crypto.sha256(content)
    return storedHash == actualHash
end
```

### Cache Keys
```pf
def getCacheKey(params):
    let paramsStr = JSON.stringify(params)
    return Crypto.md5(paramsStr)
end

let cacheKey = getCacheKey({user: 123, page: 5})
```

## See Also

- [Bytes Type](../types/bytes.md)
- [String Type](../types/string.md)
- [IO Module](io.md)
