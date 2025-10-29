# IO Module

The IO module provides file and stream I/O operations for reading, writing, and manipulating files.

## File Reading

### IO.readFile(path)
Reads the entire contents of a file as a string.

**Parameters**:
- `path: String` - Path to the file

**Returns**: `String` - File contents

**Example**:
```polyloft
let content = IO.readFile("data.txt")
println(content)

// With error handling
try:
    let data = IO.readFile("config.json")
    processData(data)
catch error:
    println("Failed to read file: #{error}")
end
```

### IO.readLines(path)
Reads a file and returns an array of lines.

**Parameters**:
- `path: String` - Path to the file

**Returns**: `Array<String>` - Array of lines

**Example**:
```polyloft
let lines = IO.readLines("log.txt")
for line in lines:
    println(line)
end

// Filter non-empty lines
let nonEmpty = lines.filter((l: String) -> Bool => l.length() > 0)
```

## File Writing

### IO.writeFile(path, content, mode?)
Writes content to a file, creating it if it doesn't exist.

**Parameters**:
- `path: String` - Path to the file
- `content: String` - Content to write
- `mode: Int` (optional) - File permissions (Unix), default 0644

**Returns**: `Bool` - true on success

**Example**:
```polyloft
IO.writeFile("output.txt", "Hello, World!")

// With custom permissions
IO.writeFile("script.sh", "#!/bin/bash\necho 'Hi'", 0755)

// Write JSON
let data = {"name": "Alice", "age": 30}
IO.writeFile("data.json", JSON.stringify(data))
```

### IO.appendFile(path, content)
Appends content to a file, creating it if it doesn't exist.

**Parameters**:
- `path: String` - Path to the file
- `content: String` - Content to append

**Returns**: `Bool` - true on success

**Example**:
```polyloft
IO.appendFile("log.txt", "Log entry\n")

// Logging function
def log(message: String):
    let timestamp = Sys.time()
    IO.appendFile("app.log", "[#{timestamp}] #{message}\n")
end

log("Application started")
```

## File Operations

### IO.copyFile(source, destination)
Copies a file from source to destination.

**Parameters**:
- `source: String` - Source file path
- `destination: String` - Destination file path

**Returns**: `Float` - Number of bytes copied

**Example**:
```polyloft
let bytes = IO.copyFile("original.txt", "backup.txt")
println("Copied #{bytes} bytes")
```

### IO.moveFile(source, destination)
Moves/renames a file from source to destination.

**Parameters**:
- `source: String` - Source file path
- `destination: String` - Destination file path

**Returns**: `Bool` - true on success

**Example**:
```polyloft
IO.moveFile("temp.txt", "final.txt")

// Rename file
IO.moveFile("old_name.txt", "new_name.txt")
```

### IO.deleteFile(path)
Deletes a file.

**Parameters**:
- `path: String` - Path to the file

**Returns**: `Bool` - true on success

**Example**:
```polyloft
if IO.exists("temp.txt"):
    IO.deleteFile("temp.txt")
    println("Temporary file removed")
end
```

### IO.exists(path)
Checks if a file or directory exists.

**Parameters**:
- `path: String` - Path to check

**Returns**: `Bool` - true if exists

**Example**:
```polyloft
if IO.exists("config.json"):
    let config = IO.readFile("config.json")
else:
    println("Config file not found")
end
```

### IO.isFile(path)
Checks if path is a file.

**Parameters**:
- `path: String` - Path to check

**Returns**: `Bool` - true if it's a file

**Example**:
```polyloft
if IO.isFile("data.txt"):
    processFile("data.txt")
end
```

### IO.isDir(path)
Checks if path is a directory.

**Parameters**:
- `path: String` - Path to check

**Returns**: `Bool` - true if it's a directory

**Example**:
```polyloft
if IO.isDir("uploads"):
    let files = IO.listDir("uploads")
end
```

## Directory Operations

### IO.listDir(path)
Lists files and directories in a directory.

**Parameters**:
- `path: String` - Directory path

**Returns**: `Array<String>` - Array of file/directory names

**Example**:
```polyloft
let files = IO.listDir(".")
for file in files:
    println(file)
end

// Filter only .pf files
let pfFiles = files.filter((f: String) -> Bool => f.endsWith(".pf"))
```

### IO.mkdir(path)
Creates a directory.

**Parameters**:
- `path: String` - Directory path

**Returns**: `Bool` - true on success

**Example**:
```polyloft
IO.mkdir("output")
IO.mkdir("output/temp")
```

### IO.mkdirAll(path)
Creates a directory and all parent directories.

**Parameters**:
- `path: String` - Directory path

**Returns**: `Bool` - true on success

**Example**:
```polyloft
IO.mkdirAll("output/temp/nested/deep")
// Creates all intermediate directories
```

### IO.removeDir(path)
Removes an empty directory.

**Parameters**:
- `path: String` - Directory path

**Returns**: `Bool` - true on success

**Example**:
```polyloft
IO.removeDir("temp")
```

### IO.removeDirAll(path)
Recursively removes a directory and all its contents.

**Parameters**:
- `path: String` - Directory path

**Returns**: `Bool` - true on success

**Example**:
```polyloft
IO.removeDirAll("build")  // Removes directory and all files
```

## Path Operations

### IO.join(parts...)
Joins path segments into a single path.

**Parameters**:
- `parts: String...` - Path segments to join

**Returns**: `String` - Joined path

**Example**:
```polyloft
let path = IO.join("home", "user", "documents", "file.txt")
// Unix: "home/user/documents/file.txt"
// Windows: "home\user\documents\file.txt"
```

### IO.basename(path)
Returns the last element of path.

**Parameters**:
- `path: String` - File path

**Returns**: `String` - Base name

**Example**:
```polyloft
IO.basename("/home/user/file.txt")  // "file.txt"
IO.basename("C:\\Users\\user\\doc.pdf")  // "doc.pdf"
```

### IO.dirname(path)
Returns the directory part of path.

**Parameters**:
- `path: String` - File path

**Returns**: `String` - Directory name

**Example**:
```polyloft
IO.dirname("/home/user/file.txt")  // "/home/user"
IO.dirname("C:\\Users\\user\\doc.pdf")  // "C:\\Users\\user"
```

### IO.ext(path)
Returns the file extension.

**Parameters**:
- `path: String` - File path

**Returns**: `String` - File extension (including dot)

**Example**:
```polyloft
IO.ext("file.txt")      // ".txt"
IO.ext("archive.tar.gz") // ".gz"
IO.ext("README")        // ""
```

## File Information

### IO.stat(path)
Returns file information.

**Parameters**:
- `path: String` - File path

**Returns**: `Map` - File statistics

**Example**:
```polyloft
let info = IO.stat("data.txt")
println("Size: #{info.size} bytes")
println("Modified: #{info.modTime}")
println("Is directory: #{info.isDir}")
```

### IO.size(path)
Returns file size in bytes.

**Parameters**:
- `path: String` - File path

**Returns**: `Int` - File size

**Example**:
```polyloft
let bytes = IO.size("video.mp4")
let mb = bytes / (1024 * 1024)
println("Video is #{mb} MB")
```

## Practical Examples

### File Backup
```polyloft
def backupFile(filename: String):
    if IO.exists(filename):
        let timestamp = Sys.time()
        let backup = "#{filename}.#{timestamp}.bak"
        IO.copyFile(filename, backup)
        println("Backed up to #{backup}")
    end
end

backupFile("important.dat")
```

### Log Rotation
```polyloft
def rotateLog(logFile: String, maxSize: Int):
    if IO.exists(logFile) and IO.size(logFile) > maxSize:
        let i = 1
        while IO.exists("#{logFile}.#{i}"):
            i = i + 1
        end
        IO.moveFile(logFile, "#{logFile}.#{i}")
    end
end

rotateLog("app.log", 1024 * 1024)  // Rotate at 1MB
```

### Directory Iterator
```polyloft
def walkDir(dir: String, callback):
    if not IO.isDir(dir):
        return
    end
    
    let files = IO.listDir(dir)
    for file in files:
        let path = IO.join(dir, file)
        if IO.isFile(path):
            callback(path)
        elif IO.isDir(path):
            walkDir(path, callback)
        end
    end
end

// Find all .pf files
walkDir("src", (path: String) => {
    if path.endsWith(".pf"):
        println(path)
    end
})
```

### Configuration Manager
```polyloft
class Config:
    private let path: String
    private let data: Map
    
    Config(path: String):
        this.path = path
        this.load()
    end
    
    def load():
        if IO.exists(this.path):
            let content = IO.readFile(this.path)
            this.data = JSON.parse(content)
        else:
            this.data = Map()
        end
    end
    
    def save():
        let content = JSON.stringify(this.data)
        IO.writeFile(this.path, content)
    end
    
    def get(key: String):
        return this.data.get(key)
    end
    
    def set(key: String, value):
        this.data.set(key, value)
        this.save()
    end
end

let config = Config("settings.json")
config.set("theme", "dark")
```

## Notes

- All file paths can be absolute or relative
- Directory separators are platform-specific (/ on Unix, \ on Windows)
- Use `IO.join()` for cross-platform path handling
- File permissions (mode) only apply on Unix-like systems
- Always handle errors when performing I/O operations
- Large files may cause memory issues with `readFile()`; consider streaming

## See Also

- [Sys Module](sys.md) - System operations
- [Net Module](net.md) - Network I/O
- [Exception Handling](../language/exceptions.md) - Error handling
