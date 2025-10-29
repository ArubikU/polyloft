# IO Module

The `IO` module provides file and directory operations.

## File Operations

### `IO.readFile(path, encoding?)`
Reads entire file content.

**Parameters:**
- `path` (String): File path
- `encoding` (String, optional): Encoding (default: "utf-8")

**Returns:** String

```pf
let content = IO.readFile("data.txt")
println(content)

// With encoding
let content = IO.readFile("data.txt", "utf-8")
```

### `IO.writeFile(path, content, encoding?)`
Writes content to file.

**Parameters:**
- `path` (String): File path
- `content` (String): Content to write
- `encoding` (String, optional): Encoding (default: "utf-8")

**Returns:** void

```pf
IO.writeFile("output.txt", "Hello, World!")

let data = "Line 1\nLine 2\nLine 3"
IO.writeFile("data.txt", data)
```

### `IO.appendFile(path, content)`
Appends content to file.

**Parameters:**
- `path` (String): File path
- `content` (String): Content to append

**Returns:** void

```pf
IO.appendFile("log.txt", "New log entry\n")
```

### `IO.exists(path)`
Checks if file or directory exists.

**Parameters:**
- `path` (String): Path to check

**Returns:** Bool

```pf
if IO.exists("config.json"):
    let config = IO.readFile("config.json")
else:
    println("Config file not found")
end
```

### `IO.delete(path)`
Deletes a file.

**Parameters:**
- `path` (String): File path

**Returns:** void

```pf
if IO.exists("temp.txt"):
    IO.delete("temp.txt")
    println("File deleted")
end
```

## Directory Operations

### `IO.mkdir(path)`
Creates a directory.

**Parameters:**
- `path` (String): Directory path

**Returns:** void

```pf
IO.mkdir("logs")
IO.mkdir("data/cache")  // Creates parent directories
```

### `IO.readDir(path)`
Lists directory contents.

**Parameters:**
- `path` (String): Directory path

**Returns:** Array of Strings

```pf
let files = IO.readDir(".")
for file in files:
    println(file)
end
```

### `IO.isDir(path)`
Checks if path is a directory.

**Parameters:**
- `path` (String): Path to check

**Returns:** Bool

```pf
if IO.isDir("data"):
    let files = IO.readDir("data")
end
```

### `IO.isFile(path)`
Checks if path is a file.

**Parameters:**
- `path` (String): Path to check

**Returns:** Bool

```pf
if IO.isFile("data.txt"):
    let content = IO.readFile("data.txt")
end
```

## File Information

### `IO.getFileSize(path)`
Gets file size in bytes.

**Parameters:**
- `path` (String): File path

**Returns:** Int

```pf
let size = IO.getFileSize("data.txt")
println("File size: #{size} bytes")
```

### `IO.getFileInfo(path)`
Gets detailed file information.

**Parameters:**
- `path` (String): File path

**Returns:** Map

```pf
let info = IO.getFileInfo("data.txt")
println("Size: #{info.size}")
println("Modified: #{info.modTime}")
println("IsDir: #{info.isDir}")
```

## Classes

### `File`
Represents an open file handle.

#### Opening Files
```pf
let file = IO.openFile("data.txt")
let content = file.read()
file.close()
```

#### Methods

**`read(size?)`** - Read from file
```pf
let file = IO.openFile("data.txt")
let content = file.read()  // Read all
let chunk = file.read(100)  // Read 100 bytes
file.close()
```

**`write(content)`** - Write to file
```pf
let file = IO.openFile("output.txt", "w")
file.write("Hello, World!")
file.close()
```

**`readLine()`** - Read one line
```pf
let file = IO.openFile("data.txt")
loop:
    let line = file.readLine()
    if line == nil:
        break
    end
    println(line)
end
file.close()
```

**`close()`** - Close file
```pf
let file = IO.openFile("data.txt")
// Use file
file.close()  // Always close!
```

### `Buffer`
In-memory buffer for reading/writing.

```pf
// Create empty buffer
let buf = Buffer()
buf.write("Hello")
buf.write(" World")
let content = buf.toString()
println(content)  // "Hello World"

// Create buffer with data
let buf = Buffer("Initial data")
buf.write(" more data")
```

#### Methods

**`write(data)`** - Write to buffer
```pf
let buf = Buffer()
buf.write("Line 1\n")
buf.write("Line 2\n")
```

**`read(size)`** - Read from buffer
```pf
let data = buf.read(10)  // Read 10 bytes
```

**`toString()`** - Get buffer content
```pf
let content = buf.toString()
```

**`clear()`** - Clear buffer
```pf
buf.clear()
```

## Examples

### Reading Configuration
```pf
def loadConfig(filename):
    try:
        if IO.exists(filename):
            let content = IO.readFile(filename)
            return parseJSON(content)
        else:
            return getDefaultConfig()
        end
    catch e:
        println("Error loading config: #{e}")
        return getDefaultConfig()
    end
end

let config = loadConfig("config.json")
```

### Writing Log File
```pf
def log(message):
    let timestamp = Sys.time()
    let entry = "[#{timestamp}] #{message}\n"
    IO.appendFile("app.log", entry)
end

log("Application started")
log("Processing data")
log("Application stopped")
```

### Processing Large Files
```pf
def processLargeFile(filename):
    let file = IO.openFile(filename)
    let count = 0
    
    loop:
        let line = file.readLine()
        if line == nil:
            break
        end
        
        processLine(line)
        count = count + 1
    end
    
    file.close()
    return count
end
```

### Directory Listing
```pf
def listAllFiles(dir):
    let files = []
    let entries = IO.readDir(dir)
    
    for entry in entries:
        let path = dir + "/" + entry
        if IO.isDir(path):
            let subFiles = listAllFiles(path)
            files = files.concat(subFiles)
        else:
            files = files.concat([path])
        end
    end
    
    return files
end

let allFiles = listAllFiles(".")
for file in allFiles:
    println(file)
end
```

### File Backup
```pf
def backupFile(filename):
    if not IO.exists(filename):
        println("File not found")
        return false
    end
    
    let content = IO.readFile(filename)
    let backupName = filename + ".backup"
    IO.writeFile(backupName, content)
    println("Backup created: #{backupName}")
    return true
end
```

### CSV Reader
```pf
def readCSV(filename):
    let file = IO.openFile(filename)
    let rows = []
    
    loop:
        let line = file.readLine()
        if line == nil:
            break
        end
        
        let fields = line.trim().split(",")
        rows = rows.concat([fields])
    end
    
    file.close()
    return rows
end

let data = readCSV("data.csv")
for row in data:
    println(row)
end
```

### File Search
```pf
def findFiles(dir, extension):
    let results = []
    let entries = IO.readDir(dir)
    
    for entry in entries:
        let path = dir + "/" + entry
        if IO.isFile(path) and entry.endsWith(extension):
            results = results.concat([path])
        end
    end
    
    return results
end

let textFiles = findFiles(".", ".txt")
println("Found #{textFiles.length()} text files")
```

## Best Practices

### ✅ DO - Always close files
```pf
let file = IO.openFile("data.txt")
try:
    let content = file.read()
    processContent(content)
catch e:
    println("Error: #{e}")
end
file.close()  // Always close
```

### ✅ DO - Check file existence
```pf
if IO.exists("config.json"):
    let config = IO.readFile("config.json")
else:
    createDefaultConfig()
end
```

### ✅ DO - Handle errors
```pf
try:
    let content = IO.readFile("data.txt")
catch e:
    println("Could not read file: #{e}")
    content = ""
end
```

### ❌ DON'T - Read huge files into memory
```pf
// Bad: May run out of memory
let huge = IO.readFile("10GB.txt")

// Good: Process line by line
let file = IO.openFile("10GB.txt")
loop:
    let line = file.readLine()
    if line == nil:
        break
    end
    processLine(line)
end
file.close()
```

## See Also

- [Sys Module](sys.md)
- [String Type](../types/string.md)
- [Exception Handling](../control-flow/exceptions.md)
