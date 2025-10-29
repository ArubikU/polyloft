# String Type

Strings in Polyloft are immutable sequences of characters with extensive manipulation methods.

## Creating Strings

### String Literals
```pf
let str1 = "Hello, World!"
let str2 = 'Single quotes also work'
let empty = ""
```

### String Interpolation
```pf
let name = "Alice"
let age = 25
let message = "My name is #{name} and I'm #{age} years old"
println(message)  // My name is Alice and I'm 25 years old
```

### Multiline Strings
```pf
let multiline = "Line 1
Line 2
Line 3"
```

## Methods

### `length()`
Returns the number of characters in the string.

**Returns:** Int

```pf
let text = "Hello"
println(text.length())  // 5
```

### `isEmpty()`
Checks if the string is empty.

**Returns:** Bool

```pf
println("".isEmpty())        // true
println("hello".isEmpty())   // false
```

### `charAt(index)`
Returns the character at the specified index.

**Parameters:**
- `index` (Int): Position of character

**Returns:** String

```pf
let text = "Hello"
println(text.charAt(0))  // "H"
println(text.charAt(4))  // "o"
```

### `indexOf(substring)`
Returns the index of the first occurrence of substring, or -1 if not found.

**Parameters:**
- `substring` (String): String to search for

**Returns:** Int

```pf
let text = "Hello, World!"
println(text.indexOf("World"))  // 7
println(text.indexOf("xyz"))    // -1
```

### `substring(start, end?)`
Extracts a substring from start to end (exclusive).

**Parameters:**
- `start` (Int): Starting index
- `end` (Int, optional): Ending index (exclusive)

**Returns:** String

```pf
let text = "Hello, World!"
println(text.substring(0, 5))   // "Hello"
println(text.substring(7))      // "World!"
```

### `toUpperCase()`
Converts string to uppercase.

**Returns:** String

```pf
println("hello".toUpperCase())  // "HELLO"
```

### `toLowerCase()`
Converts string to lowercase.

**Returns:** String

```pf
println("HELLO".toLowerCase())  // "hello"
```

### `trim()`
Removes leading and trailing whitespace.

**Returns:** String

```pf
println("  hello  ".trim())  // "hello"
```

### `startsWith(prefix)`
Checks if string starts with prefix.

**Parameters:**
- `prefix` (String): Prefix to check

**Returns:** Bool

```pf
let text = "Hello, World!"
println(text.startsWith("Hello"))  // true
println(text.startsWith("World"))  // false
```

### `endsWith(suffix)`
Checks if string ends with suffix.

**Parameters:**
- `suffix` (String): Suffix to check

**Returns:** Bool

```pf
let text = "Hello, World!"
println(text.endsWith("World!"))  // true
println(text.endsWith("Hello"))   // false
```

### `contains(substring)`
Checks if string contains substring.

**Parameters:**
- `substring` (String): String to search for

**Returns:** Bool

```pf
let text = "Hello, World!"
println(text.contains("World"))  // true
println(text.contains("xyz"))    // false
```

### `replace(old, new)`
Replaces all occurrences of old with new.

**Parameters:**
- `old` (String): String to replace
- `new` (String): Replacement string

**Returns:** String

```pf
let text = "Hello, World!"
println(text.replace("World", "Polyloft"))  // "Hello, Polyloft!"
```

### `split(delimiter)`
Splits string into array using delimiter.

**Parameters:**
- `delimiter` (String): Split delimiter

**Returns:** Array

```pf
let text = "apple,banana,cherry"
let fruits = text.split(",")
println(fruits)  // ["apple", "banana", "cherry"]
```

### `repeat(count)`
Repeats the string count times.

**Parameters:**
- `count` (Int): Number of repetitions

**Returns:** String

```pf
println("Ha".repeat(3))  // "HaHaHa"
```

### `padStart(length, padString)`
Pads string to length at the start.

**Parameters:**
- `length` (Int): Target length
- `padString` (String): Padding string

**Returns:** String

```pf
println("5".padStart(3, "0"))  // "005"
```

### `padEnd(length, padString)`
Pads string to length at the end.

**Parameters:**
- `length` (Int): Target length
- `padString` (String): Padding string

**Returns:** String

```pf
println("5".padEnd(3, "0"))  // "500"
```

## Examples

### String Concatenation
```pf
let first = "Hello"
let second = "World"
let combined = first + " " + second
println(combined)  // "Hello World"
```

### String Comparison
```pf
let a = "apple"
let b = "banana"

println(a == b)  // false
println(a < b)   // true (alphabetical)
```

### Building Strings
```pf
let parts = ["Hello", "beautiful", "world"]
let sentence = ""

for part in parts:
    sentence = sentence + part + " "
end

println(sentence.trim())  // "Hello beautiful world"
```

### Parsing
```pf
let csv = "Alice,25,NYC"
let fields = csv.split(",")
let name = fields[0]
let age = fields[1]
let city = fields[2]

println("Name: #{name}")
println("Age: #{age}")
println("City: #{city}")
```

### Validation
```pf
def isValidEmail(email):
    return email.contains("@") and 
           email.contains(".") and
           not email.isEmpty()
end

println(isValidEmail("user@example.com"))  // true
println(isValidEmail("invalid"))           // false
```

### Text Processing
```pf
def titleCase(text):
    let words = text.toLowerCase().split(" ")
    let result = []
    
    for word in words:
        if word.length() > 0:
            let first = word.charAt(0).toUpperCase()
            let rest = word.substring(1)
            result = result.concat([first + rest])
        end
    end
    
    return result.join(" ")
end

println(titleCase("hello world"))  // "Hello World"
```

### Counting Characters
```pf
def countChar(text, target):
    let count = 0
    for char in text:
        if char == target:
            count = count + 1
        end
    end
    return count
end

println(countChar("hello", "l"))  // 2
```

## String Operators

### Indexing
```pf
let text = "Hello"
println(text[0])  // "H"
println(text[4])  // "o"
```

### Slicing
```pf
let text = "Hello, World!"
println(text[0:5])   // "Hello"
println(text[7:12])  // "World"
```

### Iteration
```pf
let text = "Hello"
for char in text:
    println(char)
end
// Outputs: H, e, l, l, o
```

## Best Practices

### ✅ DO - Use string interpolation
```pf
let name = "Alice"
let age = 25
let message = "Hello, I'm #{name} and I'm #{age}"  // Good
```

### ✅ DO - Check before operations
```pf
if not text.isEmpty():
    let first = text.charAt(0)
end
```

### ❌ DON'T - Concatenate in loops (inefficient)
```pf
// Bad: Slow for large loops
let result = ""
for i in range(1000):
    result = result + "x"
end

// Better: Use array and join
let parts = []
for i in range(1000):
    parts = parts.concat(["x"])
end
let result = parts.join("")
```

## See Also

- [Array Type](array.md)
- [String Interpolation](../language/interpolation.md)
